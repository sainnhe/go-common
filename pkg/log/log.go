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

type loggerTypeT int

const (
	loggerTypeLight = loggerTypeT(0)
	loggerTypeLocal = loggerTypeT(1)
	loggerTypeOTel  = loggerTypeT(2)
)

var gCfg *Config
var gLogLevel slog.Level
var gLoggerType loggerTypeT
var gLogger *slog.Logger
var gWriter io.Writer
var mu sync.Mutex
var defaultCfg = &Config{
	"light",
	"debug",
	LocalConfig{},
}

func handleSetGlobalConfig(cfg *Config) (cleanup func(), err error) {
	// Init a non-nil cleanup function to avoid panic on calling it.
	cleanup = func() {}

	// Check if cfg is nil.
	if cfg == nil {
		err = constant.ErrNilDeps
		return
	}

	// Check the log level.
	var logLevel slog.Level
	switch cfg.Level {
	case "debug", "":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		err = errors.New("invalid log level")
		return
	}

	// Check the logger type and set global logger
	var loggerType loggerTypeT
	switch cfg.Type {
	case "light", "":
		loggerType = loggerTypeLight
	case "local":
		loggerType = loggerTypeLocal
		cleanup = initMultiWriter(&cfg.Local)
	case "otel":
		loggerType = loggerTypeOTel
	default:
		err = errors.New("invalid logger type")
		return
	}

	// Set global variables.
	gLogLevel = logLevel
	gLoggerType = loggerType
	gCfg = cfg
	gLogger = handleNewLogger("global")

	return
}

func initMultiWriter(cfg *LocalConfig) (cleanup func()) {
	consoleWriter := os.Stderr
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
	}
	gWriter = io.MultiWriter(consoleWriter, fileWriter)
	return func() {
		if err := errors.Join(consoleWriter.Close(), fileWriter.Close()); err != nil {
			GetGlobalLogger().Error("Close logger writer failed.", constant.LogAttrError, err)
		}
		// syscall.Sync() returns an error on macOS but doesn't return anything on Linux, so let's disable errcheck here
		syscall.Sync() // nolint:errcheck
	}
}

func handleNewLogger(pkgName string) *slog.Logger {
	switch gLoggerType {
	case loggerTypeLocal:
		return slog.New(tint.NewHandler(gWriter, &tint.Options{
			AddSource:  true,
			Level:      gLogLevel,
			TimeFormat: time.StampMilli,
			NoColor:    false,
		})).With(constant.LogAttrPackage, pkgName)
	case loggerTypeOTel:
		return otelslog.NewLogger(pkgName, otelslog.WithSource(true))
	default:
		return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      gLogLevel,
			TimeFormat: time.StampMilli,
			NoColor:    false,
		})).With(constant.LogAttrPackage, pkgName)
	}
}

// SetGlobalConfig sets a global config that will be used every time a new logger is initialized, and returns a cleanup
// hook that cleans resources used by loggers.
//
// Note that calling this function will also sets a global logger based on the given config.
func SetGlobalConfig(cfg *Config) (cleanup func(), err error) {
	mu.Lock()
	defer mu.Unlock()

	return handleSetGlobalConfig(cfg)
}

// GetGlobalLogger returns the global logger.
// If the global logger is not set, initialize the global logger based on a default config and return it.
func GetGlobalLogger() *slog.Logger {
	mu.Lock()
	defer mu.Unlock()

	if gCfg == nil || gLogger == nil {
		_, _ = handleSetGlobalConfig(defaultCfg)
	}

	return gLogger
}

// NewLogger initializes a new logger with the given package name.
func NewLogger(pkgName string) *slog.Logger {
	mu.Lock()
	if gCfg == nil {
		_, _ = handleSetGlobalConfig(defaultCfg)
	}
	mu.Unlock()

	return handleNewLogger(pkgName)
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
