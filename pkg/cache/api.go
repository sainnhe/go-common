//go:generate mockgen -typed -write_package_comment=false -source=api.go -destination=api_mock.go -package cache

// Package cache implements valkey cache.
package cache

import (
	"context"

	"github.com/valkey-io/valkey-go"
)

// Proxy is the valkey proxy.
type Proxy interface {
	// Set sets a key value pair.
	Set(ctx context.Context, key, value string) error
	// Setex sets a key value pair and the expiration time.
	Setex(ctx context.Context, key, value string, seconds int64) error
	// Get gets a result of a key.
	Get(ctx context.Context, key string) valkey.ValkeyResult
	// Delete deletes a key.
	Delete(ctx context.Context, key string) error
	// Expire sets the expiration time of a key.
	Expire(ctx context.Context, key string, seconds int64) error
	// Incr increases the value of a key by 1 and returns the value after increasing.
	Incr(ctx context.Context, key string) (int64, error)
	// IncrBy increases the value of a key by N and returns the value after increasing.
	IncrBy(ctx context.Context, key string, increment int64) (int64, error)
}
