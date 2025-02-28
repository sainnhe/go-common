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
}

// Add increments the counter by delta and adds to the internal WaitGroup.
// If Logger is set, logger warns when adding positive delta after [WaitGroup.Wait] has been called.
// Refer to [sync.WaitGroup.Add] for more information.
func (w *WaitGroup) Add(delta int) {
	w.wg.Add(delta)
	count := atomic.AddInt64(&w.count, int64(delta))

	if w.shouldLog() {
		if delta > 0 {
			if len(w.Name) > 0 {
				w.Logger.Warn("[WaitGroup] Added tasks after Wait call.",
					logAttrName, w.Name, logAttrDelta, delta, logAttrCount, count)
			} else {
				w.Logger.Warn("[WaitGroup] Added tasks after Wait call.", logAttrDelta, delta, logAttrCount, count)
			}
		} else if delta < 0 {
			w.logCompletion(count)
		}
	}
}

// Done decrements the counter by 1 and signals completion to the internal WaitGroup.
// If Logger is set, logger tracks counter decrements after [WaitGroup.Wait] has been called.
// Refer to [sync.WaitGroup.Done] for more information.
func (w *WaitGroup) Done() {
	w.wg.Done()
	count := atomic.AddInt64(&w.count, -1)

	if w.shouldLog() {
		w.logCompletion(count)
	}
}

// Wait blocks until the counter reaches zero. Enables logging of subsequent operations.
// Refer to [sync.WaitGroup.Wait] for more information.
func (w *WaitGroup) Wait() {
	atomic.StoreUint32(&w.waitStarted, 1)
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

// logCompletion handles common completion logging logic
func (w *WaitGroup) logCompletion(count int64) {
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
