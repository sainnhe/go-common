package glock_test

import (
	"testing"
	"time"

	"github.com/teamsorghum/go-common/pkg/glock"
)

func TestGlock(t *testing.T) {
	t.Parallel()

	sleepTime := time.Duration(500) * time.Millisecond
	startTime := time.Now()
	glock.Lock()
	go func() {
		time.Sleep(sleepTime)
		glock.Unlock()
	}()
	glock.Wait()
	if time.Since(startTime) < sleepTime {
		t.Fatalf("Expect duration < %+v", sleepTime)
	}
}
