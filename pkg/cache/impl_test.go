package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/teamsorghum/go-common/pkg/cache"
	"github.com/teamsorghum/go-common/pkg/log"
	"github.com/teamsorghum/go-common/pkg/util"
	"go.uber.org/mock/gomock"
)

func TestProxy(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.NewMockLogger(ctrl)
	logger.EXPECT().WithAttrs(gomock.Any()).Return(logger).AnyTimes()
	logger.EXPECT().WithContext(gomock.Any()).Return(logger).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	// Read config in environment variables.
	cfg, err := util.LoadConfigFile[cache.Config]("", "")
	if err != nil {
		t.Fatalf("Load config error: %+v", err)
	}

	// Initialize Proxy
	proxy, cleanup, err := cache.NewProxyImpl(cfg, logger)
	if err != nil {
		t.Fatalf("Initialize Proxy failed: %v", err)
	}
	defer cleanup()

	t.Run("set and get", func(t *testing.T) { // nolint:paralleltest
		key := "test:set_and_get"
		value := "test_value"

		err := proxy.Set(ctx, key, value)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		result := proxy.Get(ctx, key)
		if result.Error() != nil {
			t.Errorf("Get failed: %v", result.Error())
		} else {
			val, err := result.AsBytes()
			if err != nil || string(val) != value {
				t.Errorf("Expect %s, get %s, err = %+v", value, string(val), err)
			}
		}

		// Cleanup
		err = proxy.Delete(ctx, key)
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}

		result = proxy.Get(ctx, key)
		if result.Error() == nil {
			t.Errorf("Getting a deleted key should return error.")
		}
	})

	t.Run("setex and expire", func(t *testing.T) { // nolint:paralleltest
		key := "test:setex_and_expire"
		value := "test_value"

		err := proxy.Setex(ctx, key, value, 1)
		if err != nil {
			t.Errorf("Setex failed: %v", err)
		}

		result := proxy.Get(ctx, key)
		if result.Error() != nil {
			t.Errorf("Get failed: %v", result.Error())
		} else {
			val, err := result.AsBytes()
			if err != nil || string(val) != value {
				t.Errorf("Expect %s, get %s, err = %+v", value, string(val), err)
			}
		}

		// Waiting for key to expire.
		time.Sleep(2 * time.Second)

		result = proxy.Get(ctx, key)
		if result.Error() == nil {
			t.Errorf("Getting a deleted key should return error.")
		}
	})

	t.Run("incr and incrby", func(t *testing.T) { // nolint:paralleltest
		key := "test:incr_and_incr_by"

		// Ensure key does not exist.
		_ = proxy.Delete(ctx, key)

		val, err := proxy.Incr(ctx, key)
		if err != nil {
			t.Errorf("Incr failed: %v", err)
		}
		if val != 1 {
			t.Errorf("Expect 1, get %d", val)
		}

		val, err = proxy.IncrBy(ctx, key, 5)
		if err != nil {
			t.Errorf("IncrBy failed: %v", err)
		}
		if val != 6 {
			t.Errorf("Expect 6, get %d", val)
		}

		// Cleanup
		err = proxy.Delete(ctx, key)
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})
}
