package limiter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislimiter"
	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeylimiter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type peakShavingImpl struct {
	rl     rueidislimiter.RateLimiterClient
	vl     valkeylimiter.RateLimiterClient
	rp     sync.Pool
	vp     sync.Pool
	l      *slog.Logger
	cfg    *PeakShavingConfig
	prefix string
	meter  metric.Meter
}

// NewRedisPeakShavingProxy initializes a new peak shaving proxy using Redis.
func NewRedisPeakShavingProxy(
	cfg *PeakShavingConfig, rueidisClient rueidis.Client, logger *slog.Logger) (proxy Proxy, err error) {
	// Check arguments
	if cfg == nil || logger == nil || rueidisClient == nil {
		err = constant.ErrNilDep
		return
	}

	// Initialize rueidis limiter
	limiter, err := rueidislimiter.NewRateLimiter(rueidislimiter.RateLimiterOption{
		ClientBuilder: func(_ rueidis.ClientOption) (rueidis.Client, error) { return rueidisClient, nil },
		KeyPrefix:     "peak_" + cfg.Prefix,
		Limit:         cfg.Limit,
		Window:        time.Duration(cfg.WindowMs) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	// Initialize peakShavingImpl
	return &peakShavingImpl{
		limiter,
		nil,
		sync.Pool{
			New: func() any {
				return new(rueidislimiter.Result)
			},
		},
		sync.Pool{},
		logger.With(constant.LogAttrAPI, "peak_shaving"),
		cfg,
		"*",
		otel.Meter("github.com/teamsorghum/go-common/pkg/limiter"),
	}, nil
}

// NewValkeyPeakShavingProxy initializes a new peak shaving proxy using Valkey.
func NewValkeyPeakShavingProxy(
	cfg *PeakShavingConfig, valkeyClient valkey.Client, logger *slog.Logger) (proxy Proxy, err error) {
	// Check arguments
	if cfg == nil || logger == nil || valkeyClient == nil {
		err = constant.ErrNilDep
		return
	}

	// Initialize valkey limiter
	limiter, err := valkeylimiter.NewRateLimiter(valkeylimiter.RateLimiterOption{
		ClientBuilder: func(_ valkey.ClientOption) (valkey.Client, error) { return valkeyClient, nil },
		KeyPrefix:     "peak_" + cfg.Prefix,
		Limit:         cfg.Limit,
		Window:        time.Duration(cfg.WindowMs) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	// Initialize peakShavingImpl
	return &peakShavingImpl{
		nil,
		limiter,
		sync.Pool{},
		sync.Pool{
			New: func() any {
				return new(valkeylimiter.Result)
			},
		},
		logger.With(constant.LogAttrAPI, "peak_shaving"),
		cfg,
		"*",
		otel.Meter("github.com/teamsorghum/go-common/pkg/limiter"),
	}, nil
}

func (p *peakShavingImpl) Check(ctx context.Context, identifier string) (result *Result, err error) {
	l := p.l.With(constant.LogAttrMethod, "Check", "identifier", identifier)

	// Return Allowed if peak shaving is disabled.
	if !p.cfg.Enable {
		l.DebugContext(ctx, "Peak shaving disabled. Skipping...")
		return &Result{Allowed: true}, nil
	}

	// If Redis is used
	if p.rl != nil {
		rr := p.rp.Get().(*rueidislimiter.Result)
		defer p.rp.Put(rr)
		*rr, err = p.rl.Check(ctx, identifier)
		return convertResult(rr, nil), err
	}

	// If Valkey is used
	vr := p.vp.Get().(*valkeylimiter.Result)
	defer p.vp.Put(vr)
	*vr, err = p.vl.Check(ctx, identifier)
	return convertResult(nil, vr), err
}

func (p *peakShavingImpl) Allow(ctx context.Context, identifier string) (*Result, error) {
	l := p.l.With(constant.LogAttrMethod, "Allow", "identifier", identifier)
	return p.allowN(ctx, identifier, 1, l)
}

func (p *peakShavingImpl) AllowN(ctx context.Context, identifier string, n int64) (*Result, error) {
	l := p.l.With(constant.LogAttrMethod, "AllowN", "identifier", identifier)
	return p.allowN(ctx, identifier, n, l)
}

func (p *peakShavingImpl) allowN(
	ctx context.Context, identifier string, n int64, logger *slog.Logger) (result *Result, err error) {

	// Return Allowed if peak shaving is disabled.
	if !p.cfg.Enable {
		logger.DebugContext(ctx, "Peak shaving disabled. Skipping...")
		return &Result{Allowed: true}, nil
	}

	// Initialize metric counter
	mc, err := p.meter.Int64Counter("limiter.peakshaving.failed",
		metric.WithDescription("Peak shaving failed"), metric.WithUnit("1"))
	if err != nil {
		logger.ErrorContext(ctx, "Get meter failed.", constant.LogAttrError, err.Error())
	}

	// Try for MaxAttempts times
	for i := range p.cfg.MaxAttempts {
		tmpLogger := logger.With("attempt", i+1)

		// AllowN
		if p.rl != nil { // If Redis is used
			rr := p.rp.Get().(*rueidislimiter.Result)
			*rr, err = p.rl.AllowN(ctx, identifier, n)
			result = convertResult(rr, nil)
			p.rp.Put(rr)
		} else { // If Valkey is used
			vr := p.vp.Get().(*valkeylimiter.Result)
			*vr, err = p.vl.AllowN(ctx, identifier, n)
			result = convertResult(nil, vr)
			p.vp.Put(vr)
		}

		// Handle result
		if err != nil {
			mc.Add(ctx, 1)
			tmpLogger.ErrorContext(ctx, "Peak shaving error.", constant.LogAttrError, err.Error())
			return
		}
		if result.Allowed {
			tmpLogger.DebugContext(ctx, "Peak shaving allowed.")
			return
		}
		tmpLogger.WarnContext(ctx, "Reaches peak shaving limit. Sleep and retry.")
		time.Sleep(time.Duration(p.cfg.AttemptIntervalMs) * time.Millisecond)
	}
	mc.Add(ctx, 1)
	logger.ErrorContext(ctx, "Peak shaving hits max attempts.")
	return
}
