package dlock_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/redis/rueidis"
	"github.com/sainnhe/go-common/pkg/dlock"
)

func TestDlock_nilDeps(t *testing.T) {
	t.Parallel()

	s, e := dlock.NewService(nil, nil)
	if s != nil || e == nil {
		t.Fatal("Expect s == nil and e != nil")
	}
}

func TestDlock(t *testing.T) {
	t.Parallel()

	// Init rueidis client
	rc, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Config
	cfg := &dlock.Config{
		Prefix:       "test_dlock",
		ExpireMs:     2000,
		RetryAfterMs: 30,
	}

	// Init service
	locker, err := dlock.NewService(cfg, rc)
	if err != nil {
		t.Fatal(err)
	}

	// Init wait group
	wg := &sync.WaitGroup{}
	key1 := "foo" // will expire
	key2 := "bar" // will be released

	// Results
	errs := []error{}

	// Try to acquire key1 and key2 immediately, should succeed
	wg.Add(1)
	go func() {
		defer wg.Done()

		// key1
		success, e := locker.TryAcquire(context.Background(), key1)
		if e != nil {
			errs = append(errs, fmt.Errorf("[1] Expect nil error, got %w", e))
		}
		if !success {
			errs = append(errs, fmt.Errorf("[1] Expect success = true, got false"))
		}

		// key2
		success, e = locker.TryAcquire(context.Background(), key2)
		if e != nil {
			errs = append(errs, fmt.Errorf("[1] Expect nil error, got %w", e))
		}
		if !success {
			errs = append(errs, fmt.Errorf("[1] Expect success = true, got false"))
		}
	}()

	// Acquire key1 and key2 after 500 ms, should succeed
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(500) * time.Millisecond)

		// key1
		e := locker.Acquire(context.Background(), key1)
		if e != nil {
			errs = append(errs, fmt.Errorf("[2] lock foo failed, err = %w", e))
		}

		// key2
		e = locker.Acquire(context.Background(), key2)
		if e != nil {
			errs = append(errs, fmt.Errorf("[2] lock bar failed, err = %w", e))
		}
	}()

	// Try to acquire key1 after 1000ms, should fail because it has not been expired
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(1000) * time.Millisecond)
		success, e := locker.TryAcquire(context.Background(), key1)
		if e != nil {
			errs = append(errs, fmt.Errorf("[3] Expect nil error, got %w", e))
		}
		if success {
			errs = append(errs, fmt.Errorf("[3] Expect success = false, got true"))
		}
	}()

	// Acquire key1 after 1000ms with a 500 ms context, should fail
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(1000) * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(500)*time.Millisecond)
		defer cancel()
		e := locker.Acquire(ctx, key1)
		if !errors.Is(e, context.DeadlineExceeded) {
			errs = append(errs, fmt.Errorf("[4] Expect DeadlineExceeded, got %+v", e))
		}
	}()

	// Try to acquire key1 after 3000ms, should succeed because key has expired.
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(3000) * time.Millisecond)
		success, e := locker.TryAcquire(context.Background(), key1)
		if e != nil {
			errs = append(errs, fmt.Errorf("[5] Expect nil error, got %w", e))
		}
		if !success {
			errs = append(errs, fmt.Errorf("[5] Expect success = true, got false"))
		}
	}()

	// Release after 3000ms, should fail because key has expired.
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(3000) * time.Millisecond)
		e := locker.Release(context.Background(), key1)
		if !errors.Is(e, dlock.ErrKeyNotExists) {
			errs = append(errs, fmt.Errorf("[6] Expect ErrKeyNotExists, got %+v", e))
		}
	}()

	// After 1500ms, release key2, should succeed
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(1500) * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200)*time.Millisecond)
		defer cancel()

		// Release
		e := locker.Release(ctx, key2)
		if e != nil {
			errs = append(errs, fmt.Errorf("[7] Expect nil error, got %w", e))
			return
		}
	}()

	// Check results
	wg.Wait()
	if len(errs) != 0 {
		t.Fatalf("Errors: %+v", errs)
	}
}
