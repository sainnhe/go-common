// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

// nolint:goconst
package limiter_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/rueidis"
	"github.com/teamsorghum/go-common/pkg/limiter"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/valkey-io/valkey-go"
)

// This example demonstrates how to perform rate limit using Redis.
// It assumes you have a working Redis server listened on localhost:6379 with empty username and password.
func Example_rateLimitRedis() {
	logger := log.Global()

	// Initialize a rueidis client.
	rueidisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a new rate limit service.
	s, err := limiter.NewRedisRateLimitService(&limiter.RateLimitConfig{
		Enable:   true,            // Enable rate limit.
		Prefix:   "redis_example", // Prefix for keys used to describe current business and avoid conflicts.
		Limit:    1,               // Limit of requests in a given time window.
		WindowMs: 1000,            // Time window for measurement in milliseconds.
	}, rueidisClient, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Let's launch 3 goroutines:
	//
	// 1. One for allowing 1 request immediately, which increases the counter by 1.
	// 2. One for allowing 3 requests after sleeping for 200 milliseconds, which increases the counter by 3.
	// 3. One for checking if a request is allowed under the limit without incrementing the counter, after sleeping for
	//    500 milliseconds.
	//
	// Since we limit 1 request in 1 second, there should be 3 failed requests and the check should also fail.

	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	identifier := "test"
	failedCount := int32(0)
	checkSuccess := true
	ctx := context.Background()
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		result, err := s.Allow(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			atomic.AddInt32(&failedCount, 1)
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Duration(200) * time.Millisecond)
		result, err := s.AllowN(ctx, identifier, 3)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			atomic.AddInt32(&failedCount, 3)
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		result, err := s.Check(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		mu.Lock()
		checkSuccess = result.Allowed
		mu.Unlock()
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("failedCount = %d, checkSuccess = %t\n", failedCount, checkSuccess)

	// Output: failedCount = 3, checkSuccess = false
}

// This example demonstrates how to perform rate limit using Valkey.
// It assumes you have a working Valkey server listened on localhost:6379 with empty username and password.
func Example_rateLimitValkey() {
	logger := log.Global()

	// Initialize a valkey client.
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a new rate limit service.
	s, err := limiter.NewValkeyRateLimitService(&limiter.RateLimitConfig{
		Enable:   true,             // Enable rate limit.
		Prefix:   "valkey_example", // Prefix for keys used to describe current business and avoid conflicts.
		Limit:    1,                // Limit of requests in a given time window.
		WindowMs: 1000,             // Time window for measurement in milliseconds.
	}, valkeyClient, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Let's launch 3 goroutines:
	//
	// 1. One for allowing 1 request immediately, which increases the counter by 1.
	// 2. One for allowing 3 requests after sleeping for 200 milliseconds, which increases the counter by 3.
	// 3. One for checking if a request is allowed under the limit without incrementing the counter, after sleeping for
	//    500 milliseconds.
	//
	// Since we limit 1 request in 1 second, there should be 3 failed requests and the check should also fail.

	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	identifier := "test"
	failedCount := int32(0)
	checkSuccess := true
	ctx := context.Background()
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		result, err := s.Allow(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			atomic.AddInt32(&failedCount, 1)
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Duration(200) * time.Millisecond)
		result, err := s.AllowN(ctx, identifier, 3)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			atomic.AddInt32(&failedCount, 3)
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Duration(500) * time.Millisecond)
		result, err := s.Check(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		mu.Lock()
		checkSuccess = result.Allowed
		mu.Unlock()
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("failedCount = %d, checkSuccess = %t\n", failedCount, checkSuccess)

	// Output: failedCount = 3, checkSuccess = false
}

// This example demonstrates how to perform peak shaving using Redis.
// It assumes you have a working Redis server listened on localhost:6379 with empty username and password.
func Example_peakShavingRedis() {
	logger := log.Global()

	// Initialize a rueidis client.
	rueidisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a new peak shaving service.
	s, err := limiter.NewRedisPeakShavingService(&limiter.PeakShavingConfig{
		Enable:            true,            // Enable peak shaving.
		Prefix:            "redis_example", // Prefix for keys used to describe current business and avoid conflicts.
		Limit:             3,               // Limit of requests in a given time window.
		WindowMs:          500,             // Time window for measurement in milliseconds.
		MaxAttempts:       3,               // max number of attempts
		AttemptIntervalMs: 500,             // interval between each attempt in milliseconds
	}, rueidisClient, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Let's launch 3 goroutines:
	//
	// 1. One for allowing 1 request immediately, which increases the counter by 1.
	// 2. One for allowing 3 requests after sleeping for 250 milliseconds, which increases the counter by 3.
	// 3. One for checking immediately if a request is allowed under the limit without incrementing the counter.
	//
	// According to our config, the first request should success, and the following 3 requests should fail at the first
	// time, but retry successfully after waiting 500ms. The check should success, and the total time should be greater
	// than 500 milliseconds.

	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	identifier := "test"
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(3)
	startTime := time.Now()

	go func() {
		result, err := s.Allow(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			logger.Error(fmt.Sprintf("Expect allowed, got %+v", result))
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Duration(250) * time.Millisecond)
		result, err := s.AllowN(ctx, identifier, 3)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			logger.Error(fmt.Sprintf("Expect allowed, got %+v", result))
		}
		wg.Done()
	}()

	go func() {
		result, err := s.Check(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			logger.Error(fmt.Sprintf("Expect allowed, got %+v", result))
		}
		wg.Done()
	}()

	wg.Wait()

	if time.Since(startTime) > time.Duration(500)*time.Millisecond {
		fmt.Println("Total time > 500 milliseconds")
	} else {
		fmt.Println("Total time <= 500 milliseconds")
	}

	// Output: Total time > 500 milliseconds
}

// This example demonstrates how to perform peak shaving using Valkey.
// It assumes you have a working Valkey server listened on localhost:6379 with empty username and password.
func Example_peakShavingValkey() {
	logger := log.Global()

	// Initialize a valkey client.
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a new peak shaving service.
	s, err := limiter.NewValkeyPeakShavingService(&limiter.PeakShavingConfig{
		Enable:            true,             // Enable peak shaving.
		Prefix:            "valkey_example", // Prefix for keys used to describe current business and avoid conflicts.
		Limit:             3,                // Limit of requests in a given time window.
		WindowMs:          500,              // Time window for measurement in milliseconds.
		MaxAttempts:       3,                // max number of attempts
		AttemptIntervalMs: 500,              // interval between each attempt in milliseconds
	}, valkeyClient, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Let's launch 3 goroutines:
	//
	// 1. One for allowing 1 request immediately, which increases the counter by 1.
	// 2. One for allowing 3 requests after sleeping for 250 milliseconds, which increases the counter by 3.
	// 3. One for checking immediately if a request is allowed under the limit without incrementing the counter.
	//
	// According to our config, the first request should success, and the following 3 requests should fail at the first
	// time, but retry successfully after waiting 500ms. The check should success, and the total time should be greater
	// than 500 milliseconds.

	// The identifier is used to group traffics. Requests with the same identifier share the same counter.
	identifier := "test"
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(3)
	startTime := time.Now()

	go func() {
		result, err := s.Allow(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			logger.Error(fmt.Sprintf("Expect allowed, got %+v", result))
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Duration(250) * time.Millisecond)
		result, err := s.AllowN(ctx, identifier, 3)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			logger.Error(fmt.Sprintf("Expect allowed, got %+v", result))
		}
		wg.Done()
	}()

	go func() {
		result, err := s.Check(ctx, identifier)
		if err != nil {
			logger.Error(err.Error())
		}
		if !result.Allowed {
			logger.Error(fmt.Sprintf("Expect allowed, got %+v", result))
		}
		wg.Done()
	}()

	wg.Wait()

	if time.Since(startTime) > time.Duration(500)*time.Millisecond {
		fmt.Println("Total time > 500 milliseconds")
	} else {
		fmt.Println("Total time <= 500 milliseconds")
	}

	// Output: Total time > 500 milliseconds
}
