package limiter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/rueidis"
	"github.com/teamsorghum/go-common/pkg/limiter"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/valkey-io/valkey-go"
)

func TestRateLimit_nilDependency(t *testing.T) {
	t.Parallel()

	// Redis
	proxy, err := limiter.NewRedisRateLimitProxy(nil, nil, nil)
	if proxy != nil || err == nil {
		t.Fatalf("Got proxy = %+v, err = %+v", proxy, err)
	}

	// Valkey
	proxy, err = limiter.NewValkeyRateLimitProxy(nil, nil, nil)
	if proxy != nil || err == nil {
		t.Fatalf("Got proxy = %+v, err = %+v", proxy, err)
	}
}

func TestRateLimit_disable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	identifier := "test_disable"

	// Redis
	rueidisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	proxy, err := limiter.NewRedisRateLimitProxy(
		&limiter.RateLimitConfig{}, rueidisClient, log.Global())
	if proxy == nil || err != nil {
		t.Fatalf("Got proxy = %+v, err = %+v", proxy, err)
	}

	result1, err1 := proxy.Allow(ctx, identifier)
	result2, err2 := proxy.AllowN(ctx, identifier, 3)
	result3, err3 := proxy.Check(ctx, identifier)

	if !result1.Allowed || !result2.Allowed || !result3.Allowed {
		t.Fatalf("Expect all allowed, got result1 = %+v, result2 = %+v, result3 = %+v", result1, result2, result3)
	}
	if errors.Join(err1, err2, err3) != nil {
		t.Fatalf("Expect no error, got err1 = %+v, err2 = %+v, err3 = %+v", err1, err2, err3)
	}

	// Valkey
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	proxy, err = limiter.NewValkeyRateLimitProxy(
		&limiter.RateLimitConfig{}, valkeyClient, log.Global())
	if proxy == nil || err != nil {
		t.Fatalf("Got proxy = %+v, err = %+v", proxy, err)
	}

	result1, err1 = proxy.Allow(ctx, identifier)
	result2, err2 = proxy.AllowN(ctx, identifier, 3)
	result3, err3 = proxy.Check(ctx, identifier)

	if !result1.Allowed || !result2.Allowed || !result3.Allowed {
		t.Fatalf("Expect all allowed, got result1 = %+v, result2 = %+v, result3 = %+v", result1, result2, result3)
	}
	if errors.Join(err1, err2, err3) != nil {
		t.Fatalf("Expect no error, got err1 = %+v, err2 = %+v, err3 = %+v", err1, err2, err3)
	}
}
