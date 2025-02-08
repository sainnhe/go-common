//go:generate mockgen -typed -write_package_comment=false -source=api.go -destination=api_mock.go -package trafficlimit

// Package trafficlimit implements traffic limitation.
package trafficlimit

import (
	"context"
)

// Proxy is the traffic limitation proxy.
type Proxy interface {
	// RateLimit performs rate limitation. If the current traffic exceeds the specified QPS, return failure.
	RateLimit(ctx context.Context) error
	// PeakShaving shaves the peak traffic. If the current traffic exceeds the specified QPS, sleep for a while and
	// retry for N times. If the traffic still exceeds the specified QPS, return failure.
	PeakShaving(ctx context.Context) error
	// SetPrefix sets prefix for keys used in cache, which can be used to describe current business and avoid conflicts.
	SetPrefix(prefix string)
}
