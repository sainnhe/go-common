// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package limiter

import (
	"context"
	"log/slog"
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

// RateLimitService is the rate limit service.
type RateLimitService struct {
	rl     rueidislimiter.RateLimiterClient
	vl     valkeylimiter.RateLimiterClient
	rp     sync.Pool
	vp     sync.Pool
	l      *slog.Logger
	cfg    *RateLimitConfig
	prefix string
	meter  metric.Meter
}

// NewRedisRateLimitService initializes a new rate limit service using Redis.
func NewRedisRateLimitService(cfg *RateLimitConfig, rueidisClient rueidis.Client) (*RateLimitService, error) {
	// Check arguments
	if cfg == nil || rueidisClient == nil {
		return nil, constant.ErrNilDeps
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

	// Initialize RateLimitService
	return &RateLimitService{
		limiter,
		nil,
		sync.Pool{
			New: func() any {
				return new(rueidislimiter.Result)
			},
		},
		sync.Pool{},
		log.NewLogger(pkgName),
		cfg,
		"*",
		otel.Meter(pkgName),
	}, nil
}

// NewValkeyRateLimitService initializes a new rate limit service using Valkey.
func NewValkeyRateLimitService(cfg *RateLimitConfig, valkeyClient valkey.Client) (*RateLimitService, error) {
	// Check arguments
	if cfg == nil || valkeyClient == nil {
		return nil, constant.ErrNilDeps
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

	// Initialize RateLimitService
	return &RateLimitService{
		nil,
		limiter,
		sync.Pool{},
		sync.Pool{
			New: func() any {
				return new(valkeylimiter.Result)
			},
		},
		log.NewLogger(pkgName),
		cfg,
		"*",
		otel.Meter(pkgName),
	}, nil
}

// Check checks if a request is allowed under the limit without incrementing the counter.
// The identifier is used to group traffics. Requests with the same identifier share the same counter.
func (r *RateLimitService) Check(ctx context.Context, identifier string) (result *Result, err error) {
	l := r.l.With(constant.LogAttrMethod, "Check", "identifier", identifier)

	// Return Allowed if rate limit is disabled.
	if !r.cfg.Enable {
		l.DebugContext(ctx, "Rate limit disabled. Skipping...")
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

// Allow allows a single request, incrementing the counter if allowed.
// The identifier is used to group traffics. Requests with the same identifier share the same counter.
func (r *RateLimitService) Allow(ctx context.Context, identifier string) (*Result, error) {
	l := r.l.With(constant.LogAttrMethod, "Allow", "identifier", identifier)
	return r.allowN(ctx, identifier, 1, l)
}

// AllowN allows n requests, incrementing the counter accordingly if allowed.
// The identifier is used to group traffics. Requests with the same identifier share the same counter.
func (r *RateLimitService) AllowN(ctx context.Context, identifier string, n int64) (*Result, error) {
	l := r.l.With(constant.LogAttrMethod, "AllowN", "identifier", identifier)
	return r.allowN(ctx, identifier, n, l)
}

// allowN is the actual underlying implementation of Allow and AllowN.
func (r *RateLimitService) allowN(
	ctx context.Context, identifier string, n int64, logger *slog.Logger) (result *Result, err error) {
	// Return Allowed if rate limit is disabled.
	if !r.cfg.Enable {
		logger.DebugContext(ctx, "Rate limit disabled. Skipping...")
		return &Result{Allowed: true}, nil
	}

	// Initialize metric counter
	mc, err := r.meter.Int64Counter("limiter.ratelimit.failed",
		metric.WithDescription("Rate limit failed"), metric.WithUnit("1"))
	if err != nil {
		logger.ErrorContext(ctx, "Get meter failed.", constant.LogAttrError, err.Error())
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
		logger.ErrorContext(ctx, "Rate limit error.", constant.LogAttrError, err.Error())
		return result, err
	}
	if !result.Allowed {
		mc.Add(ctx, 1)
		logger.ErrorContext(ctx, "Exceeds rate limit.")
	} else {
		logger.DebugContext(ctx, "Rate limit allowed.")
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
