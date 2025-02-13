package glock

import (
	"sync"
	"sync/atomic"

	"github.com/teamsorghum/go-common/pkg/log"
)

// waitGroupImpl defines a wait group with counter and status.
type waitGroupImpl struct {
	sync.WaitGroup
	count           int64
	shutdownStarted uint32
}

// Add adds delta and bumps counter.
func (w *waitGroupImpl) Add(delta int) {
	atomic.AddInt64(&w.count, int64(delta))
	w.WaitGroup.Add(delta)
	if atomic.LoadUint32(&w.shutdownStarted) > 0 {
		log.GetDefault().Warn("Shutdown started but a new goroutine lock is added")
	}
}

// Done decrease wait group counter by 1.
func (w *waitGroupImpl) Done() {
	atomic.AddInt64(&w.count, -1)
	w.WaitGroup.Done()
	if atomic.LoadUint32(&w.shutdownStarted) > 0 {
		count := w.GetCount()
		if count > 0 {
			log.GetDefault().Info("Waiting for goroutine locks to be released...", "remain", count)
		} else {
			log.GetDefault().Info("All goroutine locks have been released")
		}
	}
}

// GetCount gets wait group counter.
func (w *waitGroupImpl) GetCount() int {
	return int(atomic.LoadInt64(&w.count))
}

// StartShutdown sets the shutdown status to true.
func (w *waitGroupImpl) StartShutdown() {
	atomic.AddUint32(&w.shutdownStarted, 1)
}
