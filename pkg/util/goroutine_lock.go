package util

import (
	"context"
	"log/slog"
	"time"

	goroutinelock "github.com/sainnhe/go-common/internal/goroutine_lock"
)

// GoroutineLock locks goroutine to ensure that the task won't be interrupted.
func GoroutineLock() {
	goroutinelock.Wg.Add(1)
}

// GoroutineUnlock releases a goroutine lock.
// NOTE: This function must be used via defer to avoid panic in the middle and causing the lock to not be released.
func GoroutineUnlock() {
	goroutinelock.Wg.Done()
}

// GoroutineWait waits for all goroutine locks to be released, or the timeout period has been exceeded.
func GoroutineWait(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	wgDone := make(chan struct{})
	go func() {
		goroutinelock.Wg.StartShutdown()
		if count := goroutinelock.Wg.GetCount(); count > 0 {
			slog.Info("Waiting for goroutine locks to be released...", "remain", count)
			goroutinelock.Wg.Wait()
		}
		close(wgDone)
	}()
	select {
	case <-ctx.Done():
		slog.Warn("The timeout period has been exceeded, however there are still some goroutine locks that have"+
			" not been released.", "remain", goroutinelock.Wg.GetCount())
	case <-wgDone:
	}
}
