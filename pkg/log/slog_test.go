package log_test

import (
	"context"
	"testing"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/log"
	ctxutil "github.com/teamsorghum/go-common/pkg/util/ctx"
)

const logPath = "/tmp/teamsorghum-go-common-test/log/1.log"

func TestSlog_NewSlog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *log.Config
		wantErr bool
	}{
		{
			"Invalid log level",
			&log.Config{
				"slog",
				"none",
				&log.File{
					logPath,
					1,
					1,
				},
			},
			true,
		},
		{
			"Success",
			&log.Config{
				"slog",
				"debug",
				&log.File{
					logPath,
					1,
					1,
				},
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
			if tt.wantErr != (err != nil) {
				t.Fatalf("Want error: %+v, got %+v", tt.wantErr, err)
			}
		})
	}
}

func TestSlog_Logger(t *testing.T) {
	t.Parallel()

	cfg := &log.Config{
		"slog",
		"debug",
		&log.File{
			logPath,
			1,
			1,
		},
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
			"With attrs",
			constant.LogLevelDebug,
			"Test WithAttrs",
			[]any{"key2", "value2"},
			[]any{},
			[]any{"key1", "value1"},
			nil,
		},
		{
			"With context",
			constant.LogLevelDebug,
			"Test WithContext",
			[]any{"key2", "value2"},
			[]any{},
			[]any{"key1", "value1"},
			ctxutil.PutFields(context.Background(), map[any]any{"key": "value"}),
		},
		{
			"Debug",
			constant.LogLevelDebug,
			"Test debug %s",
			[]any{"key", "value"},
			[]any{"foo"},
			nil,
			nil,
		},
		{
			"Info",
			constant.LogLevelInfo,
			"Test info %s",
			[]any{"key", "value"},
			[]any{"foo"},
			nil,
			nil,
		},
		{
			"Warn",
			constant.LogLevelWarn,
			"Test warn %s",
			[]any{"key", "value"},
			[]any{"foo"},
			nil,
			nil,
		},
		{
			"Error",
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
