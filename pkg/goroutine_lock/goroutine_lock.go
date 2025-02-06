// Package goroutinelock implements goroutine lock.
package goroutinelock

import (
	"context"
	"log/slog"
	"time"

	"github.com/teamsorghum/go-common/pkg/log"
)

// wg is the waitgroup for implementing goroutine locks.
var wg = &waitGroupImpl{}

// Lock locks goroutine to ensure that the task won't be interrupted.
func Lock() {
	wg.Add(1)
}

// Unlock releases a goroutine lock.
// NOTE: This function must be used via defer to avoid panic in the middle and causing the lock to not be released.
func Unlock() {
	wg.Done()
}

// Wait waits for all goroutine locks to be released, or the timeout period has been exceeded.
func Wait(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	wgDone := make(chan struct{})
	go func() {
		wg.StartShutdown()
		if count := wg.GetCount(); count > 0 {
			slog.Info("Waiting for goroutine locks to be released...", "remain", count)
			wg.Wait()
		}
		close(wgDone)
	}()
	select {
	case <-ctx.Done():
		if log.DefaultLogger != nil {
			log.DefaultLogger.Warn("The timeout period has been exceeded, however there are still some goroutine"+
				"locks that have not been released.", "remain", wg.GetCount())
		}
	case <-wgDone:
	}
}
