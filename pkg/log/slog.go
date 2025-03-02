package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
	"github.com/teamsorghum/go-common/pkg/constant"
	"gopkg.in/natefinch/lumberjack.v2"
)

type slogImpl struct {
	logger         *slog.Logger
	originalLogger *slog.Logger
	ctx            context.Context
	attrs          []any
}

// NewSlog initializes a slog based logger.
func NewSlog(cfg *Config) (logger Logger, cleanup func(), err error) {
	if cfg == nil {
		err = errors.New("nil dependency")
		return
	}
	consoleWriter := os.Stderr
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Slog.Path,
		MaxSize:    cfg.Slog.MaxSizeMB,
		MaxBackups: cfg.Slog.MaxBackups,
	}
	cleanup = func() {
		if err := fileWriter.Close(); err != nil {
			logger.Errorf("Close fileWriter failed: %+v", err)
		}
		syscall.Sync() // nolint:errcheck
	}
	multiWriter := io.MultiWriter(consoleWriter, fileWriter)
	var logLevel slog.Level
	switch cfg.Level {
	case "debug":
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
	slogLogger := slog.New(tint.NewHandler(multiWriter, &tint.Options{
		AddSource:  true,
		Level:      logLevel,
		TimeFormat: time.StampMilli,
		NoColor:    false,
	}))
	logger = &slogImpl{
		slogLogger,
		slogLogger,
		nil,
		[]any{},
	}
	return
}

// buildLogger builds a new logger from ctx and attrs
func (s *slogImpl) buildLogger() {
	// Build all attrs
	ctxFields := GetCtxFields(s.ctx)
	resultAttrs := make([]any, 0, 2*len(ctxFields)+len(s.attrs))
	for k, v := range ctxFields {
		resultAttrs = append(resultAttrs, fmt.Sprintf("ctx_%+v", k), v)
	}
	resultAttrs = append(resultAttrs, s.attrs...)
	// Update logger
	s.logger = s.originalLogger.With(resultAttrs...)
}

func (s *slogImpl) WithAttrs(attrs ...any) Logger {
	if len(attrs) == 0 {
		return s
	}
	// Create a new slice and append new attrs to the end of s.attrs
	newAttrs := make([]any, 0, len(attrs)+len(s.attrs))
	newAttrs = append(newAttrs, s.attrs...)
	newAttrs = append(newAttrs, attrs...)
	// Build a new logger
	newLogger := &slogImpl{
		s.logger,
		s.originalLogger,
		s.ctx,
		newAttrs,
	}
	newLogger.buildLogger()
	return newLogger
}

func (s *slogImpl) WithContext(ctx context.Context) Logger {
	if len(GetCtxFields(ctx)) == 0 {
		return s
	}
	newLogger := &slogImpl{
		s.logger,
		s.originalLogger,
		ctx,
		s.attrs,
	}
	newLogger.buildLogger()
	return newLogger
}

func (s *slogImpl) Debug(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Debug(msg, attrs...)
	} else {
		s.logger.DebugContext(s.ctx, msg, attrs...)
	}
}

func (s *slogImpl) Debugf(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Debug(fmt.Sprintf(msg, args...))
	} else {
		s.logger.DebugContext(s.ctx, fmt.Sprintf(msg, args...))
	}
}

func (s *slogImpl) Info(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Info(msg, attrs...)
	} else {
		s.logger.InfoContext(s.ctx, msg, attrs...)
	}
}

func (s *slogImpl) Infof(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Info(fmt.Sprintf(msg, args...))
	} else {
		s.logger.InfoContext(s.ctx, fmt.Sprintf(msg, args...))
	}
}

func (s *slogImpl) Warn(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Warn(msg, attrs...)
	} else {
		s.logger.WarnContext(s.ctx, msg, attrs...)
	}
}

func (s *slogImpl) Warnf(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Warn(fmt.Sprintf(msg, args...))
	} else {
		s.logger.WarnContext(s.ctx, fmt.Sprintf(msg, args...))
	}
}

func (s *slogImpl) Error(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Error(msg, attrs...)
	} else {
		s.logger.ErrorContext(s.ctx, msg, attrs...)
	}
}

func (s *slogImpl) Errorf(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Error(fmt.Sprintf(msg, args...))
	} else {
		s.logger.ErrorContext(s.ctx, fmt.Sprintf(msg, args...))
	}
}

func (s *slogImpl) Fatal(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Error(msg, attrs...)
	} else {
		s.logger.ErrorContext(s.ctx, msg, attrs...)
	}
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		s.Error("Kill process failed.", constant.LogAttrError, err.Error())
	}
}

func (s *slogImpl) Fatalf(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Error(fmt.Sprintf(msg, args...))
	} else {
		s.logger.ErrorContext(s.ctx, fmt.Sprintf(msg, args...))
	}
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		s.Error("Kill process failed.", constant.LogAttrError, err.Error())
	}
}
