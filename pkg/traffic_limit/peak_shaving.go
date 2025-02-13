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

type peakShavingImpl struct {
	limiter valkeylimiter.RateLimiterClient
	l       log.Logger
	cfg     *PeakShavingConfig
	prefix  string
	meter   metric.Meter
}

// NewPeakShavingImpl initializes a new peak shaving proxy.
func NewPeakShavingImpl(
	cfg *PeakShavingConfig, valkeyClient valkey.Client, logger log.Logger) (proxy Proxy, err error) {
	if cfg == nil || logger == nil || valkeyClient == nil {
		err = fmt.Errorf("nil dependency: cfg = %+v, valkeyClient = %+v, logger = %+v", cfg, valkeyClient, logger)
		return
	}
	limiter, err := valkeylimiter.NewRateLimiter(valkeylimiter.RateLimiterOption{
		ClientBuilder: func(_ valkey.ClientOption) (valkey.Client, error) { return valkeyClient, nil },
		KeyPrefix:     "peak_" + cfg.Prefix,
		Limit:         cfg.Limit,
		Window:        time.Duration(cfg.WindowMs) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	return &peakShavingImpl{
		limiter,
		logger.WithAttrs(constant.LogAttrAPI, "peak_shaving"),
		cfg,
		"*",
		otel.Meter("github.com/teamsorghum/go-common/pkg/traffic_limit"),
	}, nil
}

func (p *peakShavingImpl) Check(ctx context.Context, identifier string,
	options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error) {
	l := p.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "Check", "identifier", identifier)
	if !p.cfg.Enable {
		l.Debug("Peak shaving disabled. Skipping...")
		return valkeylimiter.Result{Allowed: true}, nil
	}
	return p.limiter.Check(ctx, identifier, options...)
}

func (p *peakShavingImpl) Allow(ctx context.Context, identifier string,
	options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error) {
	l := p.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "Allow", "identifier", identifier)
	return p.allowN(ctx, identifier, 1, l, options...)
}

func (p *peakShavingImpl) AllowN(ctx context.Context, identifier string, n int64,
	options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error) {
	l := p.l.WithContext(ctx).WithAttrs(constant.LogAttrMethod, "AllowN", "identifier", identifier)
	return p.allowN(ctx, identifier, n, l, options...)
}

func (p *peakShavingImpl) allowN(ctx context.Context, identifier string, n int64, logger log.Logger,
	options ...valkeylimiter.RateLimitOption) (result valkeylimiter.Result, err error) {
	if !p.cfg.Enable {
		logger.Debug("Peak shaving disabled. Skipping...")
		return valkeylimiter.Result{Allowed: true}, nil
	}
	mc, err := p.meter.Int64Counter("trafficlimit.peakshaving.failed",
		metric.WithDescription("Peak shaving failed"), metric.WithUnit("1"))
	if err != nil {
		logger.Error("Get meter failed.", constant.LogAttrError, err.Error())
	}
	for i := 0; i < p.cfg.MaxAttempts; i++ {
		tmpLogger := logger.WithAttrs("attempt", i+1)
		result, err = p.limiter.AllowN(ctx, identifier, n, options...)
		if err != nil {
			mc.Add(ctx, 1)
			tmpLogger.Error("Peak shaving error.", constant.LogAttrError, err.Error())
			return
		}
		if result.Allowed {
			tmpLogger.Debug("Peak shaving allowed.")
			return
		}
		tmpLogger.Warn("Reaches peak shaving limit. Sleep and retry.")
		time.Sleep(time.Duration(p.cfg.AttemptIntervalMs) * time.Millisecond)
	}
	mc.Add(ctx, 1)
	logger.Error("Peak shaving hits max attempts.")
	return
}
