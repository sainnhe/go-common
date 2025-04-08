// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

// Package limiter implements traffic limiter, including rate limit and peak shaving.
//
// Rate limit performs rate limitation. If the current traffic exceeds the specified limitation, return failure.
//
// Peak shaving shaves the peak traffic. If the current traffic exceeds the specified limitation, sleep for a while and
// retry for N times. If the traffic still exceeds the specified limitation, return failure.
//
// This implementation supports both [Redis] and [Valkey]. You can use [rueidis.Client] or [valkey.Client] to initialize
// a new [RateLimitService] or [PeakShavingService].
//
// Metrics will be collected using the global [metric.MeterProvider].
// There are 2 metric counters that will be collected:
//
//   - "limiter.ratelimit.failed": Indicates failed rate limit in Allow and AllowN.
//   - "limiter.peakshaving.failed": Indicates failed peak shaving in Allow and AllowN.
//
// [Redis]: https://redis.io/
// [Valkey]: https://valkey.io/
package limiter

// Result is the result of a limiter operation.
type Result struct {
	// Allowed indicates whether the request is allowed.
	Allowed bool

	// Remaining is the number of remaining requests in the current window.
	Remaining int64

	// ResetAtMs is the Unix timestamp in milliseconds at which the rate limit will reset.
	ResetAtMs int64
}

const pkgName = "github.com/teamsorghum/go-common/pkg/limiter"
