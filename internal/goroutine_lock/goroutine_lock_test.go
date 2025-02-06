// Package goroutinelock_test is the test package for goroutinelock.
package goroutinelock_test

import (
	"testing"
	"time"

	goroutinelock "github.com/teamsorghum/go-common/internal/goroutine_lock"
)

// TestWaitGroup_AddAndDone tests the basic functionality of the Add and Done methods.
func TestWaitGroup_AddAndDone(t *testing.T) {
	t.Parallel()

	wg := &goroutinelock.WaitGroup{}

	// Test initial count.
	if count := wg.GetCount(); count != 0 {
		t.Errorf("Expected initial count to be 0, got %d", count)
	}

	// Increase count.
	wg.Add(1)
	if count := wg.GetCount(); count != 1 {
		t.Errorf("Expected count after Add(1) to be 1, got %d", count)
	}

	// Decrease count.
	wg.Done()
	if count := wg.GetCount(); count != 0 {
		t.Errorf("Expected count after Done() to be 0, got %d", count)
	}
}

// TestWaitGroup_Concurrency tests the WaitGroup in a concurrent environment.
func TestWaitGroup_Concurrency(t *testing.T) {
	t.Parallel()

	wg := &goroutinelock.WaitGroup{}
	goroutineCount := 100
	doneCh := make(chan struct{})

	// Launch multiple goroutines, each calling Add and Done.
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			// Simulate some work.
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()
	}

	// Wait for all goroutines to complete.
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	<-doneCh

	// Verify that the count has returned to zero.
	if count := wg.GetCount(); count != 0 {
		t.Errorf("Expected count to be 0 after all goroutines done, got %d", count)
	}
}

// TestWaitGroup_StartShutdown tests the StartShutdown method.
func TestWaitGroup_StartShutdown(t *testing.T) {
	t.Parallel()

	wg := &goroutinelock.WaitGroup{}

	// Start shutdown.
	wg.StartShutdown()

	// Attempt to add count after shutdown.
	wg.Add(1)
	if count := wg.GetCount(); count != 1 {
		t.Errorf("Expected count to be 1 after Add(1) post-shutdown, got %d", count)
	}

	// Should receive warning logs. We can't capture log content here, but we can ensure the code path is executed.

	// Complete work.
	wg.Done()
	if count := wg.GetCount(); count != 0 {
		t.Errorf("Expected count to be 0 after Done(), got %d", count)
	}
}

// TestWaitGroup_AddAfterShutdown tests adding new counts after shutdown has started.
func TestWaitGroup_AddAfterShutdown(t *testing.T) {
	t.Parallel()

	wg := &goroutinelock.WaitGroup{}

	// Start shutdown.
	wg.StartShutdown()

	// Add count after shutdown.
	wg.Add(1)
	if count := wg.GetCount(); count != 1 {
		t.Errorf("Expected count to be 1 after Add(1) post-shutdown, got %d", count)
	}

	// Complete work.
	wg.Done()
	if count := wg.GetCount(); count != 0 {
		t.Errorf("Expected count to be 0 after Done(), got %d", count)
	}
}

// TestWaitGroup_GetCount tests the accuracy of the GetCount method.
func TestWaitGroup_GetCount(t *testing.T) {
	t.Parallel()

	wg := &goroutinelock.WaitGroup{}

	wg.Add(5)
	if count := wg.GetCount(); count != 5 {
		t.Errorf("Expected count to be 5, got %d", count)
	}

	wg.Done()
	if count := wg.GetCount(); count != 4 {
		t.Errorf("Expected count to be 4 after Done(), got %d", count)
	}
}
