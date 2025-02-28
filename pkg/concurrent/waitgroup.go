package concurrent

import (
	"sync"
	"sync/atomic"

	"github.com/teamsorghum/go-common/pkg/log"
)

const (
	logAttrName  = "name"
	logAttrDelta = "delta"
	logAttrCount = "count"
)

// WaitGroup wraps [sync.WaitGroup] with counter tracking and logging capabilities.
// When Logger is set to a non-nil value, it enables logging of state transitions.
type WaitGroup struct {
	wg          sync.WaitGroup
	count       int64      // Current counter value (atomic)
	waitStarted uint32     // Atomic flag indicating if Wait was called
	Name        string     // Optional identifier that will be used in logging
	Logger      log.Logger // Optional logger instance
	mu          sync.Mutex
}

// Add increments the counter by delta and adds to the internal WaitGroup.
// This method has no effect if [WaitGroup.Wait] has been called.
// Refer to [sync.WaitGroup.Add] for more information.
func (w *WaitGroup) Add(delta int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.WaitStarted() {
		if w.Logger != nil {
			w.Logger.Warn("Attempting to add count after Wait.")
		}
		return
	}

	w.wg.Add(delta)
	atomic.AddInt64(&w.count, int64(delta))
}

// Done decrements the counter by 1 and signals completion to the internal WaitGroup.
// If Logger is set, logger tracks counter decrements after [WaitGroup.Wait] has been called.
// Refer to [sync.WaitGroup.Done] for more information.
func (w *WaitGroup) Done() {
	w.wg.Done()
	count := atomic.AddInt64(&w.count, -1)

	if w.shouldLog() {
		if count > 0 {
			if len(w.Name) > 0 {
				w.Logger.Info("[WaitGroup] Counter updated.", logAttrName, w.Name, logAttrCount, count)
			} else {
				w.Logger.Info("[WaitGroup] Counter updated.", logAttrCount, count)
			}
		} else {
			if len(w.Name) > 0 {
				w.Logger.Info("[WaitGroup] Completed.", logAttrName, w.Name)
			} else {
				w.Logger.Info("[WaitGroup] Completed.")
			}
		}
	}
}

// Wait blocks until the counter reaches zero. Enables logging of subsequent operations.
// Refer to [sync.WaitGroup.Wait] for more information.
func (w *WaitGroup) Wait() {
	w.mu.Lock()
	atomic.StoreUint32(&w.waitStarted, 1)
	w.mu.Unlock()
	w.wg.Wait()
}

// GetCount returns the current counter value atomically.
func (w *WaitGroup) GetCount() int64 {
	return atomic.LoadInt64(&w.count)
}

// WaitStarted checks if [WaitGroup.Wait] has been called.
func (w *WaitGroup) WaitStarted() bool {
	return atomic.LoadUint32(&w.waitStarted) > 0
}

// shouldLog determines if logging is enabled based on Logger presence
func (w *WaitGroup) shouldLog() bool {
	return w.Logger != nil && w.WaitStarted()
}
