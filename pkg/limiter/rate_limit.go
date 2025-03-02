package limiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislimiter"
	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeylimiter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type rateLimitImpl struct {
	rl     rueidislimiter.RateLimiterClient
	vl     valkeylimiter.RateLimiterClient
	rp     sync.Pool
	vp     sync.Pool
	l      log.Logger
	cfg    *RateLimitConfig
	prefix string
	meter  metric.Meter
}

// NewRedisRateLimitProxy initializes a new rate limit proxy using Redis.
func NewRedisRateLimitProxy(
	cfg *RateLimitConfig, rueidisClient rueidis.Client, logger log.Logger) (proxy Proxy, err error) {
	// Check arguments
	if cfg == nil || logger == nil || rueidisClient == nil {
		err = fmt.Errorf("nil dependency: cfg = %+v, rueidisClient = %+v, logger = %+v", cfg, rueidisClient, logger)
		return
	}

	// Initialize rueidis limiter
	limiter, err := rueidislimiter.NewRateLimiter(rueidislimiter.RateLimiterOption{
		ClientBuilder: func(_ rueidis.ClientOption) (rueidis.Client, error) { return rueidisClient, nil },
		KeyPrefix:     "rate_" + cfg.Prefix,
		Limit:         cfg.Limit,
		Window:        time.Duration(cfg.WindowMs) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	// Initialize rateLimitImpl
	return &rateLimitImpl{
		limiter,
		nil,
		sync.Pool{
			New: func() any {
				return new(rueidislimiter.Result)
			},
		},
		sync.Pool{},
		logger.WithAttrs(constant.LogAttrAPI, "rate_limit"),
		cfg,
		"*",
		otel.Meter("github.com/teamsorghum/go-common/pkg/limiter"),
	}, nil
}

// NewValkeyRateLimitProxy initializes a new rate limit proxy using Valkey.
func NewValkeyRateLimitProxy(
	cfg *RateLimitConfig, valkeyClient valkey.Client, logger log.Logger) (proxy Proxy, err error) {
	// Check arguments
	if cfg == nil || logger == nil || valkeyClient == nil {
		err = fmt.Errorf("nil dependency: cfg = %+v, valkeyClient = %+v, logger = %+v", cfg, valkeyClient, logger)
		return
	}

	// Initialize valkey limiter
	limiter, err := valkeylimiter.NewRateLimiter(valkeylimiter.RateLimiterOption{
		ClientBuilder: func(_ valkey.ClientOption) (valkey.Client, error) { return valkeyClient, nil },
		KeyPrefix:     "rate_" + cfg.Prefix,
		Limit:         cfg.Limit,
		Window:        time.Duration(cfg.WindowMs) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	// Initialize rateLimitImpl
	return &rateLimitImpl{
		nil,
		limiter,
		sync.Pool{},
		sync.Pool{
			New: func() any {
				return new(valkeylimiter.Result)
			},
		},
		logger.WithAttrs(constant.LogAttrAPI, "rate_limit"),
		cfg,
		"*",
		otel.Meter("github.com/teamsorghum/go-common/pkg/limiter"),
	}, nil
}

func (r *rateLimitImpl) Check(ctx context.Context, identifier string) (result *Result, err error) {
	l := r.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "Check", "identifier", identifier)

	// Return Allowed if rate limit is disabled.
	if !r.cfg.Enable {
		l.Debug("Rate limit disabled. Skipping...")
		return &Result{Allowed: true}, nil
	}

	// If Redis is used
	if r.rl != nil {
		rr := r.rp.Get().(*rueidislimiter.Result)
		defer r.rp.Put(rr)
		*rr, err = r.rl.Check(ctx, identifier)
		return convertResult(rr, nil), err
	}

	// If Valkey is used
	vr := r.vp.Get().(*valkeylimiter.Result)
	defer r.vp.Put(vr)
	*vr, err = r.vl.Check(ctx, identifier)
	return convertResult(nil, vr), err
}

func (r *rateLimitImpl) Allow(ctx context.Context, identifier string) (*Result, error) {
	l := r.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "Allow", "identifier", identifier)
	return r.allowN(ctx, identifier, 1, l)
}

func (r *rateLimitImpl) AllowN(ctx context.Context, identifier string, n int64) (*Result, error) {
	l := r.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "AllowN", "identifier", identifier)
	return r.allowN(ctx, identifier, n, l)
}

func (r *rateLimitImpl) allowN(
	ctx context.Context, identifier string, n int64, logger log.Logger) (result *Result, err error) {
	// Return Allowed if rate limit is disabled.
	if !r.cfg.Enable {
		logger.Debug("Rate limit disabled. Skipping...")
		return &Result{Allowed: true}, nil
	}

	// Initialize metric counter
	mc, err := r.meter.Int64Counter("limiter.ratelimit.failed",
		metric.WithDescription("Rate limit failed"), metric.WithUnit("1"))
	if err != nil {
		logger.Error("Get meter failed.", constant.LogAttrError, err.Error())
	}

	// AllowN
	if r.rl != nil { // If Redis is used
		rr := r.rp.Get().(*rueidislimiter.Result)
		*rr, err = r.rl.AllowN(ctx, identifier, n)
		result = convertResult(rr, nil)
		r.rp.Put(rr)
	} else { // If Valkey is used
		vr := r.vp.Get().(*valkeylimiter.Result)
		*vr, err = r.vl.AllowN(ctx, identifier, n)
		result = convertResult(nil, vr)
		r.vp.Put(vr)
	}

	// Handle result
	if err != nil {
		mc.Add(ctx, 1)
		logger.Error("Rate limit error.", constant.LogAttrError, err.Error())
		return result, err
	}
	if !result.Allowed {
		mc.Add(ctx, 1)
		logger.Error("Exceeds rate limit.")
	} else {
		logger.Debug("Rate limit allowed.")
	}
	return result, err
}

// convertResult converts rueidis or valkey result to Result.
func convertResult(rueidisResult *rueidislimiter.Result, valkeyResult *valkeylimiter.Result) *Result {
	if rueidisResult == nil && valkeyResult == nil {
		return nil
	}
	if rueidisResult != nil {
		return &Result{
			rueidisResult.Allowed,
			rueidisResult.Remaining,
			rueidisResult.ResetAtMs,
		}
	}
	return &Result{
		valkeyResult.Allowed,
		valkeyResult.Remaining,
		valkeyResult.ResetAtMs,
	}
}
