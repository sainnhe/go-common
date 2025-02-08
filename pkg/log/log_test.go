package log_test

import (
	"context"
	"testing"

	loadconfig "github.com/teamsorghum/go-common/pkg/load_config"
	"github.com/teamsorghum/go-common/pkg/log"
	ctxutil "github.com/teamsorghum/go-common/pkg/util/ctx"
)

const logPath = "/tmp/test/log"

func TestLog_NewLog(t *testing.T) {
	t.Parallel()

	defaultCfg, err := loadconfig.Load[log.Config]("", "")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		cfg     *log.Config
		wantErr bool
	}{
		{
			"Default config",
			defaultCfg,
			false,
		},
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
			"Slog",
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

			_, cleanup, err := log.ProvideLogger(tt.cfg)
			if err == nil {
				defer cleanup()
			}
			if tt.wantErr != (err != nil) {
				t.Fatalf("Want error: %+v, got %+v", tt.wantErr, err)
			}
		})
	}
}

// nolint:paralleltest
func TestLog_Logger(t *testing.T) {
	defaultCfg, err := loadconfig.Load[log.Config]("", "")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		cfg  *log.Config
	}{
		{
			"Default config",
			defaultCfg,
		},
		{
			"Slog",
			&log.Config{
				"slog",
				"debug",
				&log.File{
					logPath,
					1,
					1,
				},
			},
		},
	}

	output := func(logger log.Logger, msg string, attrs, args []any) {
		logger.Debug(msg, attrs...)
		logger.Debugf(msg, args...)
		logger.Info(msg, attrs...)
		logger.Infof(msg, args...)
		logger.Warn(msg, attrs...)
		logger.Warnf(msg, args...)
		logger.Error(msg, attrs...)
		logger.Errorf(msg, args...)
	}

	msg := "Test %s"
	attrs := []any{"attr1", "attr1", "attr2", "attr2"}
	args := []any{"arg"}

	for _, tt := range tests {
		logger, cleanup, err := log.ProvideLogger(tt.cfg)
		if err != nil {
			t.Fatal(err)
		}
		defer cleanup() // nolint:gocritic

		t.Run(tt.name+" default output", func(_ *testing.T) {
			output(logger, msg, attrs, args)
		})

		t.Run(tt.name+" with attrs", func(_ *testing.T) {
			output(
				logger.WithAttrs("k1", "v1").WithAttrs("k2", "v2"),
				msg, attrs, args)
		})

		t.Run(tt.name+" with context", func(_ *testing.T) {
			wrongCtx := ctxutil.PutFields(context.Background(), map[any]any{"wrong": "wrong"})
			ctx := ctxutil.PutFields(context.Background(), map[any]any{"k": "v"})
			output(
				logger.WithContext(wrongCtx).WithContext(ctx),
				msg, attrs, args)
		})

		t.Run(tt.name+" with attrs and context", func(_ *testing.T) {
			output(
				logger.WithAttrs("k1", "v1").WithAttrs("k2", "v2").WithContext(
					ctxutil.PutFields(context.Background(), map[any]any{"k": "v"})),
				msg, attrs, args)
		})
	}
}
