package limiter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/rueidis"
	"github.com/sainnhe/go-common/pkg/limiter"
)

func TestLimiter_nilDependency(t *testing.T) {
	t.Parallel()

	s, err := limiter.NewService(nil, nil)
	if s != nil || err == nil {
		t.Fatalf("Got service = %+v, err = %+v", s, err)
	}
}

func TestLimiter_disable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	identifier := "test_disable"

	rueidisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := limiter.NewService(
		&limiter.Config{Enable: false, EnableLog: true}, rueidisClient)
	if s == nil || err != nil {
		t.Fatalf("Got service = %+v, err = %+v", s, err)
	}

	result1, err1 := s.Allow(ctx, identifier)
	result2, err2 := s.AllowN(ctx, identifier, 3)
	result3, err3 := s.Check(ctx, identifier)

	if !result1.Allowed || !result2.Allowed || !result3.Allowed {
		t.Fatalf("Expect all allowed, got result1 = %+v, result2 = %+v, result3 = %+v", result1, result2, result3)
	}
	if errors.Join(err1, err2, err3) != nil {
		t.Fatalf("Expect no error, got err1 = %+v, err2 = %+v, err3 = %+v", err1, err2, err3)
	}
}

func TestLimiter_peakShavingFailed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	identifier := "test_failed"

	rueidisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := limiter.NewService(
		&limiter.Config{
			Enable:            true,
			Prefix:            "*",
			Limit:             2,
			WindowMs:          500,
			MaxAttempts:       2,
			AttemptIntervalMs: 500,
			EnableLog:         true,
		}, rueidisClient)
	if s == nil || err != nil {
		t.Fatalf("Got service = %+v, err = %+v", s, err)
	}

	result, err := s.AllowN(ctx, identifier, 3)

	if result.Allowed || err != nil {
		t.Fatalf("Expect not allowed and nil error, got result = %+v, err = %+v", result, err)
	}
}
