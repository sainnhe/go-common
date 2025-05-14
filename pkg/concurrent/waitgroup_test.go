package concurrent_test

import (
	"testing"
	"time"

	"github.com/sainnhe/go-common/pkg/concurrent"
	"github.com/sainnhe/go-common/pkg/log"
)

func TestWaitGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		wg   *concurrent.WaitGroup
	}{
		{
			"No name, no logger",
			&concurrent.WaitGroup{},
		},
		{
			"No logger",
			&concurrent.WaitGroup{
				Name: "test",
			},
		},
		{
			"No name",
			&concurrent.WaitGroup{
				Logger: log.GetGlobalLogger(),
			},
		},
		{
			"Has both name and logger",
			&concurrent.WaitGroup{
				Name:   "test",
				Logger: log.GetGlobalLogger(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			go func() {
				time.Sleep(time.Duration(200) * time.Millisecond)
				tt.wg.Wait()
				log.GetGlobalLogger().Info("Wait completed.")
			}()

			tt.wg.Add(3)
			tt.wg.Done()

			if count := tt.wg.GetCount(); count != 2 {
				t.Fatalf("Expect count = 2, got %d", count)
			}
			if tt.wg.WaitStarted() {
				t.Fatal("Expect WaitStarted = false, got true")
			}

			time.Sleep(time.Duration(400) * time.Millisecond)

			if !tt.wg.WaitStarted() {
				t.Fatal("Expect WaitStarted = true, got false")
			}

			tt.wg.Add(1)
			tt.wg.Done()
			tt.wg.Add(-2)

			if count := tt.wg.GetCount(); count != 0 {
				t.Fatalf("Expect count = 0, got %d", count)
			}
		})
	}
}
