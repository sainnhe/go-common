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
		err = constant.ErrNilDep
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
		logger = NewOTel(cfg.OTel)
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

// NewOTel initializes a new OTel based logger.
func NewOTel(cfg OTelConfig) *slog.Logger {
	return otelslog.NewLogger(cfg.Name,
		otelslog.WithVersion(cfg.Version),
		otelslog.WithSource(true))
}
