//go:generate mockgen -typed -write_package_comment=false -source=limiter.go -destination=limiter_mock.go -package limiter

// Package limiter implements traffic limiter, including rate limit and peak shaving.
//
// Rate limit performs rate limitation. If the current traffic exceeds the specified limitation, return failure.
//
// Peak shaving shaves the peak traffic. If the current traffic exceeds the specified limitation, sleep for a while and
// retry for N times. If the traffic still exceeds the specified limitation, return failure.
//
// This implementation supports both [Redis] and [Valkey]. You can use [rueidis.Client] or [valkey.Client] to initialize
// a new limiter [Proxy].
//
// Metrics will be collected using the global [metric.MeterProvider].
// For example, you can use the following code to register a global MeterProvider that outputs metrics to stdout every 3
// seconds:
//
//	metricExporter, err := stdoutmetric.New()
//	if err != nil {
//	    logger.Fatal(err.Error())
//	}
//	otel.SetMeterProvider(metric.NewMeterProvider(
//	    metric.WithReader(
//	        metric.NewPeriodicReader(
//	            metricExporter, metric.WithInterval(time.Duration(3)*time.Second))),
//	))
//
// There are 2 metric counters:
//
//   - "limiter.ratelimit.failed": Indicates failed rate limit in Allow and AllowN.
//   - "limiter.peakshaving.failed": Indicates failed peak shaving in Allow and AllowN.
//
// [Redis]: https://redis.io/
// [Valkey]: https://valkey.io/
package limiter

import "context"

// Result is the result of a limiter operation.
type Result struct {
	// Allowed indicates whether the request is allowed.
	Allowed bool

	// Remaining is the number of remaining requests in the current window.
	Remaining int64

	// ResetAtMs is the Unix timestamp in milliseconds at which the rate limit will reset.
	ResetAtMs int64
}

// Proxy defines the interface for a limiter.
type Proxy interface {
	// Check checks if a request is allowed under the limit without incrementing the counter.
	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	Check(ctx context.Context, identifier string) (*Result, error)

	// Allow allows a single request, incrementing the counter if allowed.
	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	Allow(ctx context.Context, identifier string) (*Result, error)

	// AllowN allows n requests, incrementing the counter accordingly if allowed.
	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	AllowN(ctx context.Context, identifier string, n int64) (*Result, error)
}
