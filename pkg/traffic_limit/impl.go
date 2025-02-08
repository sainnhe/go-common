package trafficlimit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/teamsorghum/go-common/pkg/cache"
	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/log"
)

type proxyImpl struct {
	c   cache.Proxy
	l   log.Logger
	cfg *Config
}

// NewProxyImpl initializes a new traffic limit proxy.
func NewProxyImpl(cfg *Config, logger log.Logger, cacheProxy cache.Proxy) (proxy Proxy, cleanup func(), err error) {
	if cfg == nil || cfg.RateLimit == nil || cfg.PeakShaving == nil || logger == nil || cacheProxy == nil {
		err = fmt.Errorf("nil dependency: cfg = %+v, logger = %+v, cacheProxy == %+v", cfg, logger, cacheProxy)
		return
	}
	return &proxyImpl{
		cacheProxy,
		logger.WithAttrs(constant.LogAttrAPI, "trafficlimit"),
		cfg,
	}, func() {}, nil
}

func (p *proxyImpl) RateLimit(ctx context.Context) error {
	key := p.getKey("rate")
	l := p.l.WithContext(ctx).WithAttrs(
		constant.LogAttrMethod, "RateLimit", "qps", p.cfg.RateLimit.QPS, "key", key)
	if !p.cfg.RateLimit.Enable {
		l.Debug("Rate limit disabled. Skipping...")
		return nil
	}
	reply, err := p.c.Incr(ctx, key)
	go func() {
		if err := p.c.Expire(ctx, key, constant.TrafficLimitExpirationTimeSec); err != nil {
			l.Error("Expire key failed.", constant.LogAttrError, err)
		}
	}()
	if err != nil {
		l.Error("Incr failed.", constant.LogAttrError, err)
		return err
	}
	if reply > p.cfg.RateLimit.QPS {
		l.Error("Exceeds rate limit.", "reply", reply)
		return errors.New("exceeds rate limit")
	}
	l.Debug("Pass rate limit.", "reply", reply)
	return nil
}

func (p *proxyImpl) PeakShaving(ctx context.Context) error {
	l := p.l.WithContext(ctx).WithAttrs(
		constant.LogAttrMethod, "PeakShaving",
		"qps", p.cfg.PeakShaving.QPS,
		"max_attempts", p.cfg.PeakShaving.MaxAttempts)
	if !p.cfg.PeakShaving.Enable {
		l.Debug("Peak shaving disabled. Skipping...")
		return nil
	}
	startTime := time.Now().UnixMilli()
	for i := 0; i < p.cfg.PeakShaving.MaxAttempts; i++ {
		key := p.getKey("peak")
		tmpLogger := l.WithAttrs("key", key, "attempt", i+1)
		reply, err := p.c.Incr(ctx, key)
		if err != nil {
			tmpLogger.Error("Incr failed.", constant.LogAttrError, err, "cost_time_ms", time.Now().UnixMilli()-startTime)
			return err
		}
		go func() {
			if err := p.c.Expire(ctx, key, constant.TrafficLimitExpirationTimeSec); err != nil {
				tmpLogger.Error("Expire key failed.", constant.LogAttrError, err)
			}
		}()
		if reply <= p.cfg.PeakShaving.QPS {
			tmpLogger.Debug("Peak shaving success.", "reply", reply, "cost_time_ms", time.Now().UnixMilli()-startTime)
			return nil
		}
		tmpLogger.Warn("Reach peak shaving limit, sleeping...", "reply", reply)
		time.Sleep(time.Duration(p.cfg.PeakShaving.AttemptIntervalMs) * time.Millisecond)
	}
	l.Error("Peak shaving hits max retry.", "cost_time_ms", time.Now().UnixMilli()-startTime)
	return errors.New("peak shaving hits max retry")
}

func (p *proxyImpl) getKey(operation string) string {
	return fmt.Sprintf("%s_%s_%d",
		p.cfg.Prefix,
		operation,
		time.Now().Unix())
}
