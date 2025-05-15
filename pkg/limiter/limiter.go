//go:generate mockgen -write_package_comment=false -source=limiter.go -destination=limiter_mock.go -package limiter

/*
Package limiter implements a rate limiter with support for peak shaving.

In a traditional rate limiter, if the current request volume exceeds a specified threshold, it will return a failure.

In peak shaving however, if the threshold is exceeded, the limiter will sleep for a while and retry for N times.
If the threshold is still exceeded, a failure will be returned.

So, rate limit is a special variant of peak shaving when N is 0.
You can use this package as a generic rate limiter with additional support for peak shaving.
*/
package limiter

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislimiter"
	"github.com/sainnhe/go-common/pkg/constant"
	"github.com/sainnhe/go-common/pkg/log"
)

const pkgName = "github.com/sainnhe/go-common/pkg/limiter"

// Service is the limiter service.
type Service interface {
	// Check checks if a request is allowed under the limit without incrementing the counter.
	//
	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	Check(ctx context.Context, identifier string, options ...rueidislimiter.RateLimitOption) (
		rueidislimiter.Result, error)

	// Allow allows a single request, incrementing the counter if allowed, sleeping and retrying otherwise.
	//
	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	//
	// If the maximum number of attempts is reached, the result will be not allowed and the error will be nil.
	Allow(ctx context.Context, identifier string, options ...rueidislimiter.RateLimitOption) (
		rueidislimiter.Result, error)

	// AllowN allows n requests, incrementing the counter accordingly if allowed, sleeping and retrying otherwise.
	//
	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	//
	// If the maximum number of attempts is reached, the result will be not allowed and the error will be nil.
	AllowN(ctx context.Context, identifier string, n int64, options ...rueidislimiter.RateLimitOption) (
		rueidislimiter.Result, error)
}

type serviceImpl struct {
	rl  rueidislimiter.RateLimiterClient
	l   *slog.Logger
	cfg *Config
}

// NewService initializes a new limiter service.
func NewService(cfg *Config, rc rueidis.Client) (Service, error) {
	// Check arguments
	if cfg == nil || rc == nil {
		return nil, constant.ErrNilDeps
	}

	// Initialize rueidis limiter
	rl, _ := rueidislimiter.NewRateLimiter(rueidislimiter.RateLimiterOption{
		ClientBuilder: func(_ rueidis.ClientOption) (rueidis.Client, error) { return rc, nil },
		KeyPrefix:     "peak_" + cfg.Prefix,
		Limit:         cfg.Limit,
		Window:        time.Duration(cfg.WindowMs) * time.Millisecond,
	})

	// Initialize service
	return &serviceImpl{
		rl,
		log.NewLogger(pkgName),
		cfg,
	}, nil
}

func (s *serviceImpl) Check(ctx context.Context, identifier string, options ...rueidislimiter.RateLimitOption) (
	rueidislimiter.Result, error) {
	// Return if limiter is disabled
	if !s.cfg.Enable {
		if s.cfg.EnableLog {
			s.l.DebugContext(ctx, "Limiter disabled. Skipping...")
		}
		return rueidislimiter.Result{Allowed: true}, nil
	}
	return s.rl.Check(ctx, identifier, options...)
}

func (s *serviceImpl) Allow(ctx context.Context, identifier string, options ...rueidislimiter.RateLimitOption) (
	rueidislimiter.Result, error) {
	logger := s.l.With(constant.LogAttrMethod, "Allow", "identifier", identifier)
	return s.allowN(ctx, identifier, 1, logger, options...)
}

func (s *serviceImpl) AllowN(ctx context.Context, identifier string, n int64,
	options ...rueidislimiter.RateLimitOption) (rueidislimiter.Result, error) {
	logger := s.l.With(constant.LogAttrMethod, "AllowN", "identifier", identifier, "n", n)
	return s.allowN(ctx, identifier, n, logger, options...)
}

func (s *serviceImpl) allowN(ctx context.Context, identifier string, n int64, logger *slog.Logger,
	options ...rueidislimiter.RateLimitOption) (result rueidislimiter.Result, err error) {
	// Return if limiter is disabled
	if !s.cfg.Enable {
		if s.cfg.EnableLog {
			logger.DebugContext(ctx, "Limiter disabled. Skipping...")
		}
		return rueidislimiter.Result{Allowed: true}, nil
	}

	// If peak shaving is disabled
	if s.cfg.MaxAttempts == 0 {
		result, err = s.rl.AllowN(ctx, identifier, n, options...)
		if s.cfg.EnableLog {
			if err != nil {
				logger.ErrorContext(ctx, "Rate limit failed.", "result", result, constant.LogAttrError, err)
			} else {
				logger.DebugContext(ctx, "Rate limit allowed.", "result", result)
			}
		}
		return
	}

	// Attempt for N times
	for i := range s.cfg.MaxAttempts {
		result, err = s.rl.AllowN(ctx, identifier, n, options...)
		if err != nil {
			if s.cfg.EnableLog {
				logger.ErrorContext(ctx, "Peak shaving error.",
					constant.LogAttrAttempt, i+1,
					constant.LogAttrError, err)
			}
			return
		}
		if result.Allowed {
			if s.cfg.EnableLog {
				logger.DebugContext(ctx, "Peak shaving allowed.",
					constant.LogAttrAttempt, i+1,
					constant.LogAttrResult, result,
				)
			}
			return
		}
		if s.cfg.EnableLog {
			logger.WarnContext(ctx, "Reached peak shaving limit. Sleep and retry.",
				constant.LogAttrAttempt, i+1,
				constant.LogAttrResult, result,
			)
		}
		time.Sleep(time.Duration(s.cfg.AttemptIntervalMs) * time.Millisecond)
	}
	if s.cfg.EnableLog {
		logger.ErrorContext(ctx, "Peak shaving hits max attempts.", constant.LogAttrResult, result)
	}

	return
}
