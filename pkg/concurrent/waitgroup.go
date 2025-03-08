// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package concurrent

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
)

const (
	logAttrName  = "name"
	logAttrDelta = "delta"
	logAttrCount = "count"
)

// WaitGroup is another implementation of [sync.WaitGroup] with support for counter tracking, status tracking, logging
// capabilities and adding after [WaitGroup.Wait] has been called.
// Since this implementation supports logging and uses [sync.Mutex] for synchronization, the performance is lower than
// [sync.WaitGroup]. Use it only in performance-insensitive scenarios.
//
// By default, logging is disabled. To enable logging, set Logger to a non-nil value.
//
// NOTE: A WaitGroup must not be copied after first use.
type WaitGroup struct {
	// Name is the optional identifier that will be used in logging.
	Name string

	// Logger is the optional logger instance. Set this field to a non-nil value will enable logging.
	Logger *slog.Logger

	count       int64
	waitStarted bool
	ch          chan struct{}
	mu          sync.Mutex
}

// Add adds delta, which may be negative, to the [WaitGroup] counter.
// If the counter becomes zero or negative, all goroutines blocked on [WaitGroup.Wait] are released.
//
// Unlike [sync.WaitGroup.Add], this implementation supports adding after [WaitGroup.Wait] has been called.
func (w *WaitGroup) Add(delta int) {
	// Update status
	count, waitStarted := w.updateStatus(delta, false)

	// Handle logging
	if waitStarted && w.Logger != nil {
		if delta > 0 {
			w.Logger.Warn(w.addLogPrefix("Attempt to add count after Wait()."),
				logAttrDelta, delta, logAttrCount, count)
		} else if delta < 0 {
			w.logCompletion(count)
		}
	}
}

// Done decrements the counter by 1.
func (w *WaitGroup) Done() {
	// Update status
	count, waitStarted := w.updateStatus(-1, false)

	// Handle logging
	if waitStarted && w.Logger != nil {
		w.logCompletion(count)
	}
}

// Wait blocks until the counter reaches zero or negative.
// The subsequent operations will be logged if Logger is a non-nil value.
func (w *WaitGroup) Wait() {
	// Update status
	count, _ := w.updateStatus(0, true)

	if w.Logger != nil {
		w.Logger.Info(w.addLogPrefix("Wait started."), logAttrCount, count)
	}

	// Since w.mu is not initialized, we need to initialize it here. We must add lock to avoid data race.
	if w.ch == nil {
		w.mu.Lock()
		w.ch = make(chan struct{})
		w.mu.Unlock()
	}

	// Blocks until the counter reaches zero or negative
	<-w.ch
}

// GetCount returns the current counter value, which may be negative.
func (w *WaitGroup) GetCount() int64 {
	return atomic.LoadInt64(&w.count)
}

// WaitStarted checks if [WaitGroup.Wait] has been called.
func (w *WaitGroup) WaitStarted() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.waitStarted
}

// logCompletion handles common completion logging logic.
func (w *WaitGroup) logCompletion(count int64) {
	if count > 0 {
		w.Logger.Info(w.addLogPrefix("Counter updated."), logAttrCount, count)
	} else {
		w.Logger.Info(w.addLogPrefix("Completed."))
	}
}

// updateStatus updates the status, including w.count, w.waitStarted and w.ch.
func (w *WaitGroup) updateStatus(delta int, startWait bool) (count int64, waitStarted bool) {
	// Add lock
	w.mu.Lock()
	defer w.mu.Unlock()

	// Update wait status
	if startWait {
		w.waitStarted = startWait
	}
	waitStarted = w.waitStarted

	// Update counter
	w.count += int64(delta)
	count = w.count

	// Update channel status
	if w.waitStarted && w.count <= 0 {
		close(w.ch)
	}

	return
}

// addLogPrefix adds prefix for log message based on w.Name.
func (w *WaitGroup) addLogPrefix(msg string) string {
	if len(w.Name) > 0 {
		return fmt.Sprintf("[WaitGroup-%s] %s", w.Name, msg)
	}
	return fmt.Sprintf("[WaitGroup] %s", msg)
}
