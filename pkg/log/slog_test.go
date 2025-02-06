package log_test

import (
	"context"
	"testing"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/log"
	ctxutil "github.com/teamsorghum/go-common/pkg/util/ctx"
)

func TestSlog_NewSlog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cfg       *log.Config
		expectErr bool
	}{
		{
			"InvalidFilePath",
			&log.Config{
				"slog",
				"debug",
				"/invalid/file/path",
			},
			true,
		},
		{
			"InvalidLogLevel",
			&log.Config{
				"slog",
				"none",
				"/tmp/test.log",
			},
			true,
		},
		{
			"Success",
			&log.Config{
				"slog",
				"debug",
				"/tmp/test.log",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, cleanup, err := log.NewSlog(tt.cfg)
			if err == nil {
				defer cleanup()
			}
			if tt.expectErr != (err != nil) {
				t.Fatalf("Expect error: %+v, get %+v", tt.expectErr, err)
			}
		})
	}
}

func TestSlog_Logger(t *testing.T) {
	t.Parallel()

	cfg := &log.Config{
		"slog",
		"debug",
		"/tmp/test.log",
	}
	logger, cleanup, err := log.NewSlog(cfg)
	if err != nil {
		t.Errorf("Slog init error: %+v", err)
	}
	defer cleanup()

	tests := []struct {
		name        string
		level       string
		msg         string
		attrs       []any
		args        []any
		withAttrs   []any
		withContext context.Context
	}{
		{
			"TestWithAttrs",
			constant.LogLevelDebug,
			"Test WithAttrs",
			[]any{"key2", "value2"},
			[]any{},
			[]any{"key1", "value1"},
			nil,
		},
		{
			"TestWithContext",
			constant.LogLevelDebug,
			"Test WithContext",
			[]any{"key2", "value2"},
			[]any{},
			[]any{"key1", "value1"},
			ctxutil.PutContextFields(context.Background(), map[any]any{"key": "value"}),
		},
		{
			"TestDebug",
			constant.LogLevelDebug,
			"Test debug %s",
			[]any{"key", "value"},
			[]any{"foo"},
			nil,
			nil,
		},
		{
			"TestInfo",
			constant.LogLevelInfo,
			"Test info %s",
			[]any{"key", "value"},
			[]any{"foo"},
			nil,
			nil,
		},
		{
			"TestWarn",
			constant.LogLevelWarn,
			"Test warn %s",
			[]any{"key", "value"},
			[]any{"foo"},
			nil,
			nil,
		},
		{
			"TestError",
			constant.LogLevelError,
			"Test error %s",
			[]any{"key", "value"},
			[]any{"foo"},
			nil,
			nil,
		},
	}

	for _, tt := range tests { // nolint:paralleltest
		t.Run(tt.name, func(_ *testing.T) {
			l := logger
			if tt.withAttrs != nil {
				l = l.WithAttrs(tt.withAttrs...)
			}
			if tt.withContext != nil {
				l = l.WithContext(tt.withContext)
			}

			switch tt.level {
			case constant.LogLevelDebug:
				l.Debug(tt.msg, tt.attrs...)
				l.Debugf(tt.msg, tt.args...)
			case constant.LogLevelInfo:
				l.Info(tt.msg, tt.attrs...)
				l.Infof(tt.msg, tt.args...)
			case constant.LogLevelWarn:
				l.Warn(tt.msg, tt.attrs...)
				l.Warnf(tt.msg, tt.args...)
			case constant.LogLevelError:
				l.Error(tt.msg, tt.attrs...)
				l.Errorf(tt.msg, tt.args...)
			}
		})
	}
}
