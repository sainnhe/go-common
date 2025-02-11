package trafficlimit

import (
	"context"
	"fmt"
	"time"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeylimiter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type rateLimitImpl struct {
	limiter valkeylimiter.RateLimiterClient
	l       log.Logger
	cfg     *RateLimitConfig
	prefix  string
	meter   metric.Meter
}

// NewRateLimitImpl initializes a new rate limit proxy.
func NewRateLimitImpl(
	cfg *RateLimitConfig, valkeyClient valkey.Client, logger log.Logger) (proxy Proxy, cleanup func(), err error) {
	if cfg == nil || logger == nil || valkeyClient == nil {
		err = fmt.Errorf("nil dependency: cfg = %+v, valkeyClient == %+v, logger = %+v", cfg, valkeyClient, logger)
		return
	}
	limiter, err := valkeylimiter.NewRateLimiter(valkeylimiter.RateLimiterOption{
		ClientBuilder: func(_ valkey.ClientOption) (valkey.Client, error) { return valkeyClient, nil },
		KeyPrefix:     "rate_" + cfg.Prefix,
		Limit:         cfg.Limit,
		Window:        time.Duration(cfg.WindowMs) * time.Millisecond,
	})
	if err != nil {
		return nil, func() {}, err
	}
	return &rateLimitImpl{
		limiter,
		logger.WithAttrs(constant.LogAttrAPI, "rate_limit"),
		cfg,
		"*",
		otel.Meter("github.com/teamsorghum/go-common/pkg/traffic_limit"),
	}, func() {}, nil
}

func (r *rateLimitImpl) Check(ctx context.Context, identifier string,
	options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error) {
	l := r.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "Check", "identifier", identifier)
	if !r.cfg.Enable {
		l.Debug("Rate limit disabled. Skipping...")
		return valkeylimiter.Result{Allowed: true}, nil
	}
	return r.limiter.Check(ctx, identifier, options...)
}

func (r *rateLimitImpl) Allow(ctx context.Context, identifier string,
	options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error) {
	l := r.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "Allow", "identifier", identifier)
	return r.allowN(ctx, identifier, 1, l, options...)
}

func (r *rateLimitImpl) AllowN(ctx context.Context, identifier string, n int64,
	options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error) {
	l := r.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "AllowN", "identifier", identifier)
	return r.allowN(ctx, identifier, n, l, options...)
}

func (r *rateLimitImpl) allowN(ctx context.Context, identifier string, n int64, logger log.Logger,
	options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error) {
	if !r.cfg.Enable {
		logger.Debug("Rate limit disabled. Skipping...")
		return valkeylimiter.Result{Allowed: true}, nil
	}
	mc, err := r.meter.Int64Counter("trafficlimit.ratelimit.failed",
		metric.WithDescription("Rate limit failed"), metric.WithUnit("1"))
	if err != nil {
		logger.Error("Get meter failed.", constant.LogAttrError, err.Error())
	}
	result, err := r.limiter.AllowN(ctx, identifier, n, options...)
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
