// Package trafficlimit_test is the test package for package trafficlimit.
package trafficlimit_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/teamsorghum/go-common/pkg/cache"
	loadconfig "github.com/teamsorghum/go-common/pkg/load_config"
	"github.com/teamsorghum/go-common/pkg/log"
	trafficlimit "github.com/teamsorghum/go-common/pkg/traffic_limit"
	"go.uber.org/mock/gomock"
)

func TestTrafficLimitService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cacheConfig, err := loadconfig.Load[cache.Config]("", "")
	if cacheConfig == nil || err != nil {
		t.Fatalf("Load config failed: config = %+v, err = %+v", cacheConfig, err)
	}
	trafficLimitConfig := &trafficlimit.Config{
		RateLimit: &trafficlimit.RateLimit{
			Enable: true,
			QPS:    1,
		},
		PeakShaving: &trafficlimit.PeakShaving{
			Enable:            true,
			QPS:               1,
			MaxAttempts:       2,
			AttemptIntervalMs: 1000,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.NewMockLogger(ctrl)
	logger.EXPECT().WithAttrs(gomock.Any()).Return(logger).AnyTimes()
	logger.EXPECT().WithContext(gomock.Any()).Return(logger).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	cacheProxy, cacheCleanup, err := cache.NewProxyImpl(cacheConfig, logger)
	if err != nil {
		t.Fatalf("Init cache proxy failed: %+v", err)
	}
	defer cacheCleanup()
	proxy, cleanup, err := trafficlimit.NewProxyImpl(trafficLimitConfig, logger, cacheProxy)
	if err != nil {
		t.Fatalf("Init traffic limit proxy failed: %+v", err)
	}
	defer cleanup()
	proxy.SetPrefix("test")

	t.Run("Rate limit", func(t *testing.T) { // nolint:paralleltest
		wg := &sync.WaitGroup{}
		errCount := int32(0)
		sleepUntilNextSecond()
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := proxy.RateLimit(ctx); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}()
		}
		wg.Wait()
		// Two concurrent requests should have one fail.
		if errCount != 1 {
			t.Errorf("Error count = %d, expect to be 1.", errCount)
		}
	})

	t.Run("Peak shaving", func(t *testing.T) { // nolint:paralleltest
		wg := &sync.WaitGroup{}
		errCount := int32(0)
		sleepUntilNextSecond()
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := proxy.PeakShaving(ctx); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}()
		}
		wg.Wait()
		// One of the three requests should succeed, one should wait for 1s and then succeed, and one should fail after
		// reaching the maximum number of retries.
		if errCount != 1 {
			t.Errorf("Error count = %d, expect to be 1.", errCount)
		}
	})

}

func sleepUntilNextSecond() {
	now := time.Now()
	sleepDuration := time.Second - time.Duration(now.Nanosecond())*time.Nanosecond
	time.Sleep(sleepDuration)
}
