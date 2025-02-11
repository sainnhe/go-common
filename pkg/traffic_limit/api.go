//go:generate mockgen -typed -write_package_comment=false -source=api.go -destination=api_mock.go -package trafficlimit

// Package trafficlimit implements traffic limitation, including rate limit and peak shaving.
// Rate limit performs rate limitation. If the current traffic exceeds the specified limitation, return failure.
// Peak shaving shaves the peak traffic. If the current traffic exceeds the specified limitation, sleep for a while and
// retry for N times. If the traffic still exceeds the specified limitation, return failure.
package trafficlimit

import (
	"context"

	"github.com/valkey-io/valkey-go/valkeylimiter"
)

// Proxy is the traffic limitation proxy.
type Proxy interface {
	// Check checks if a request is allowed under the limit without incrementing the count.
	Check(ctx context.Context, identifier string,
		options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error)
	// Allow allows a single request, incrementing the counter if allowed.
	Allow(ctx context.Context, identifier string,
		options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error)
	// AllowN allows n requests, incrementing the counter accordingly if allowed.
	AllowN(ctx context.Context, identifier string, n int64,
		options ...valkeylimiter.RateLimitOption) (valkeylimiter.Result, error)
}
