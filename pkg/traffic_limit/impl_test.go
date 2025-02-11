// Package trafficlimit_test is the test package for package trafficlimit.
package trafficlimit_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	loadconfig "github.com/teamsorghum/go-common/pkg/load_config"
	"github.com/teamsorghum/go-common/pkg/log"
	trafficlimit "github.com/teamsorghum/go-common/pkg/traffic_limit"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeyotel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/mock/gomock"
)

func TestTrafficLimitService(t *testing.T) {
	t.Parallel()
	metricExporter, _ := stdoutmetric.New()
	otel.SetMeterProvider(metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter, metric.WithInterval(time.Duration(3)*time.Second))),
	))

	ctx := context.Background()

	rateLimitCfg, _ := loadconfig.Load[trafficlimit.RateLimitConfig](nil, loadconfig.TypeNil)
	rateLimitCfg.Prefix = "test"
	peakShavingCfg, _ := loadconfig.Load[trafficlimit.PeakShavingConfig](nil, loadconfig.TypeNil)
	peakShavingCfg.MaxAttempts = 2
	peakShavingCfg.AttemptIntervalMs = 1000
	peakShavingCfg.Prefix = "test"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.NewMockLogger(ctrl)
	logger.EXPECT().WithAttrs(gomock.Any()).Return(logger).AnyTimes()
	logger.EXPECT().WithContext(gomock.Any()).Return(logger).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	valkeyClient, err := valkeyotel.NewClient(valkey.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		t.Fatalf("Init valkey client failed: %+v", err)
	}

	rateLimitProxy, _ := trafficlimit.NewRateLimitImpl(rateLimitCfg, valkeyClient, log.GetDefault())
	peakShavingProxy, _ := trafficlimit.NewPeakShavingImpl(peakShavingCfg, valkeyClient, log.GetDefault())

	t.Run("Rate limit", func(t *testing.T) {
		t.Parallel()

		wg := &sync.WaitGroup{}
		failedCount := int32(0)
		sleepUntilNextSecond()
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result, err := rateLimitProxy.Allow(ctx, "test")
				if err != nil {
					t.Errorf("Rate limit error: %+v", err)
				}
				if !result.Allowed {
					atomic.AddInt32(&failedCount, 1)
				}
			}()
		}
		wg.Wait()
		// Two concurrent requests should have one fail.
		if failedCount != 1 {
			t.Errorf("Failed count = %d, expect to be 1.", failedCount)
		}
	})

	t.Run("Peak shaving", func(t *testing.T) {
		t.Parallel()

		wg := &sync.WaitGroup{}
		failedCount := int32(0)
		sleepUntilNextSecond()
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result, err := peakShavingProxy.Allow(ctx, "test")
				if err != nil {
					t.Errorf("Peak shaving error: %+v", err)
				}
				if !result.Allowed {
					atomic.AddInt32(&failedCount, 1)
				}
			}()
		}
		wg.Wait()
		// One of the three requests should succeed, one should wait for 1s and then succeed, and one should fail after
		// reaching the maximum number of retries.
		if failedCount != 1 {
			t.Errorf("Failed count = %d, expect to be 1.", failedCount)
		}
	})

}

func sleepUntilNextSecond() {
	now := time.Now()
	sleepDuration := time.Second - time.Duration(now.Nanosecond())*time.Nanosecond
	time.Sleep(sleepDuration)
}
