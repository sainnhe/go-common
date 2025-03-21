// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

// Package log implements a common logger based on [slog].
package log

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
	"github.com/teamsorghum/go-common/pkg/constant"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/natefinch/lumberjack.v2"
)

// global is the global logger.
var global *slog.Logger
var mu sync.Mutex

// SetGlobal sets the global logger.
func SetGlobal(logger *slog.Logger) {
	mu.Lock()
	global = logger
	mu.Unlock()
}

// Global returns the global logger.
func Global() *slog.Logger {
	mu.Lock()
	defer mu.Unlock()
	if global == nil {
		global = NewLight(slog.LevelDebug)
	}
	return global
}

// NewLogger initializes a new [slog.Logger] based on the given [Config].
func NewLogger(cfg *Config) (logger *slog.Logger, cleanup func(), err error) {
	if cfg == nil {
		err = constant.ErrNilDeps
		return
	}

	// Let's set a default cleanup function to avoid nil pointer panic.
	cleanup = func() {}

	var level slog.Level
	switch cfg.Level {
	case "debug", "":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		err = errors.New("invalid log level")
		return
	}

	switch cfg.Type {
	case "light", "":
		logger = NewLight(level)
		return
	case "local":
		logger, cleanup = NewLocal(cfg.Local, level)
		return
	case "otel":
		logger = NewOTel(cfg.OTel.Name)
		return
	default:
		err = errors.New("invalid log")
		return
	}
}

// NewLight initializes a new light logger.
func NewLight(level slog.Level) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		Level:      level,
		TimeFormat: time.StampMilli,
		NoColor:    false,
	}))
}

// NewLocal initializes a new local logger.
func NewLocal(cfg LocalConfig, level slog.Level) (logger *slog.Logger, cleanup func()) {
	consoleWriter := os.Stderr
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
	}
	multiWriter := io.MultiWriter(consoleWriter, fileWriter)
	logger = slog.New(tint.NewHandler(multiWriter, &tint.Options{
		AddSource:  true,
		Level:      level,
		TimeFormat: time.StampMilli,
		NoColor:    false,
	}))
	cleanup = func() {
		if err := errors.Join(consoleWriter.Close(), fileWriter.Close()); err != nil {
			logger.Error("Close writer failed.", constant.LogAttrError, err.Error())
		}
		// syscall.Sync() returns an error on macOS but doesn't return anything on Linux, so let's disable errcheck here
		syscall.Sync() // nolint:errcheck
	}
	return
}

// NewOTel returns a Logger from the global LoggerProvider. The name must be the
// name of the library providing instrumentation. This name may be the same as
// the instrumented code only if that code provides built-in instrumentation.
// If the name is empty, then a implementation defined default name will be
// used instead.
func NewOTel(name string) *slog.Logger {
	if len(name) == 0 {
		name = "slog"
	}
	return otelslog.NewLogger(name, otelslog.WithSource(true))
}

// WithOTelAttrs returns a new logger with OpenTelemetry attributes.
func WithOTelAttrs(logger *slog.Logger, attrs ...attribute.KeyValue) *slog.Logger {
	if logger == nil {
		return nil
	}
	args := make([]any, 0, 2*len(attrs)) // nolint:mnd
	for _, attr := range attrs {
		args = append(args, string(attr.Key), attr.Value.AsString())
	}
	return logger.With(args...)
}
