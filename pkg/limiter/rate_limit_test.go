// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package limiter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/rueidis"
	"github.com/teamsorghum/go-common/pkg/limiter"
	"github.com/valkey-io/valkey-go"
)

func TestRateLimit_nilDependency(t *testing.T) {
	t.Parallel()

	// Redis
	s, err := limiter.NewRedisRateLimitService(nil, nil)
	if s != nil || err == nil {
		t.Fatalf("Got service = %+v, err = %+v", s, err)
	}

	// Valkey
	s, err = limiter.NewValkeyRateLimitService(nil, nil)
	if s != nil || err == nil {
		t.Fatalf("Got service = %+v, err = %+v", s, err)
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

	s, err := limiter.NewRedisRateLimitService(
		&limiter.RateLimitConfig{}, rueidisClient)
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

	// Valkey
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err = limiter.NewValkeyRateLimitService(
		&limiter.RateLimitConfig{}, valkeyClient)
	if s == nil || err != nil {
		t.Fatalf("Got service = %+v, err = %+v", s, err)
	}

	result1, err1 = s.Allow(ctx, identifier)
	result2, err2 = s.AllowN(ctx, identifier, 3)
	result3, err3 = s.Check(ctx, identifier)

	if !result1.Allowed || !result2.Allowed || !result3.Allowed {
		t.Fatalf("Expect all allowed, got result1 = %+v, result2 = %+v, result3 = %+v", result1, result2, result3)
	}
	if errors.Join(err1, err2, err3) != nil {
		t.Fatalf("Expect no error, got err1 = %+v, err2 = %+v, err3 = %+v", err1, err2, err3)
	}
}
