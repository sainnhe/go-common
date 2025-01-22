// Package goroutinelock implements goroutine locks.
package goroutinelock

import (
	"log/slog"
	"sync"
	"sync/atomic"
)

// Wg is used to implement goroutine locks.
var Wg = &WaitGroup{}

// WaitGroup defines a wait group with counter and status.
type WaitGroup struct {
	sync.WaitGroup
	count           int64
	shutdownStarted uint32
}

// Add adds delta and bumps counter.
func (wg *WaitGroup) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
	if atomic.LoadUint32(&wg.shutdownStarted) > 0 {
		slog.Warn("Shutdown started but a new goroutine lock is added")
	}
}

// Done decrease wait group counter by 1.
func (wg *WaitGroup) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
	if atomic.LoadUint32(&wg.shutdownStarted) > 0 {
		count := wg.GetCount()
		if count > 0 {
			slog.Info("Waiting for goroutine locks to be released...", "remain", count)
		} else {
			slog.Info("All goroutine locks have been released")
		}
	}
}

// GetCount gets wait group counter.
func (wg *WaitGroup) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}

// StartShutdown sets the shutdown status to true.
func (wg *WaitGroup) StartShutdown() {
	atomic.AddUint32(&wg.shutdownStarted, 1)
}
