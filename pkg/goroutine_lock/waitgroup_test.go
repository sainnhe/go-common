package goroutinelock // nolint:testpackage

import (
	"testing"
	"time"
)

// TestWaitGroup_AddAndDone tests the basic functionality of the Add and Done methods.
func TestWaitGroup_AddAndDone(t *testing.T) {
	t.Parallel()

	wg := &waitGroupImpl{}

	// Test initial got.
	if got := wg.GetCount(); got != 0 {
		t.Errorf("Want initial count to be 0, got %d", got)
	}

	// Increase count.
	wg.Add(1)
	if got := wg.GetCount(); got != 1 {
		t.Errorf("Want count after Add(1) to be 1, got %d", got)
	}

	// Decrease count.
	wg.Done()
	if got := wg.GetCount(); got != 0 {
		t.Errorf("Want count after Done() to be 0, got %d", got)
	}
}

// TestWaitGroup_Concurrency tests the WaitGroup in a concurrent environment.
func TestWaitGroup_Concurrency(t *testing.T) {
	t.Parallel()

	wg := &waitGroupImpl{}
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

	// Verify that the got has returned to zero.
	if got := wg.GetCount(); got != 0 {
		t.Errorf("Want count to be 0 after all goroutines done, got %d", got)
	}
}

// TestWaitGroup_StartShutdown tests the StartShutdown method.
func TestWaitGroup_StartShutdown(t *testing.T) {
	t.Parallel()

	wg := &waitGroupImpl{}

	// Start shutdown.
	wg.StartShutdown()

	// Attempt to add count after shutdown.
	wg.Add(1)
	if got := wg.GetCount(); got != 1 {
		t.Errorf("Want count to be 1 after Add(1) post-shutdown, got %d", got)
	}

	// Should receive warning logs. We can't capture log content here, but we can ensure the code path is executed.

	// Complete work.
	wg.Done()
	if got := wg.GetCount(); got != 0 {
		t.Errorf("Want count to be 0 after Done(), got %d", got)
	}
}

// TestWaitGroup_AddAfterShutdown tests adding new counts after shutdown has started.
func TestWaitGroup_AddAfterShutdown(t *testing.T) {
	t.Parallel()

	wg := &waitGroupImpl{}

	// Start shutdown.
	wg.StartShutdown()

	// Add count after shutdown.
	wg.Add(1)
	if got := wg.GetCount(); got != 1 {
		t.Errorf("Want count to be 1 after Add(1) post-shutdown, got %d", got)
	}

	// Complete work.
	wg.Done()
	if got := wg.GetCount(); got != 0 {
		t.Errorf("Want count to be 0 after Done(), got %d", got)
	}
}

// TestWaitGroup_GetCount tests the accuracy of the GetCount method.
func TestWaitGroup_GetCount(t *testing.T) {
	t.Parallel()

	wg := &waitGroupImpl{}

	wg.Add(5)
	if got := wg.GetCount(); got != 5 {
		t.Errorf("Want count to be 5, got %d", got)
	}

	wg.Done()
	if got := wg.GetCount(); got != 4 {
		t.Errorf("Want count to be 4 after Done(), got %d", got)
	}
}
