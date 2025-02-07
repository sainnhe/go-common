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
	ctxutil "github.com/teamsorghum/go-common/pkg/util/ctx"
	"gopkg.in/natefinch/lumberjack.v2"
)

type slogImpl struct {
	logger *slog.Logger
	ctx    context.Context
	attrs  []any
}

// NewSlog initializes a slog based logger.
func NewSlog(cfg *Config) (logger Logger, cleanup func(), err error) {
	if cfg == nil {
		err = errors.New("nil dependency")
		return
	}
	consoleWriter := os.Stderr
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.File.Path,
		MaxSize:    cfg.File.MaxSizeMB,
		MaxBackups: cfg.File.MaxBackups,
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
		return nil, nil, errors.New("invalid log level")
	}
	slogLogger := slog.New(tint.NewHandler(multiWriter, &tint.Options{
		AddSource:  true,
		Level:      logLevel,
		TimeFormat: time.StampMilli,
		NoColor:    false,
	}))
	logger = &slogImpl{
		slogLogger,
		nil,
		[]any{},
	}
	cleanup = func() {
		if err := fileWriter.Close(); err != nil {
			logger.Errorf("Close fileWriter failed: %+v", err)
		}
		syscall.Sync() // nolint:errcheck
	}
	return
}

func (s *slogImpl) buildAttrs(attrs ...any) []any {
	fields := ctxutil.GetFields(s.ctx)
	resultAttrs := make([]any, 0, 2*len(fields)+len(s.attrs)+len(attrs))
	for k, v := range fields {
		resultAttrs = append(resultAttrs, fmt.Sprintf("ctx_%+v", k), v)
	}
	resultAttrs = append(resultAttrs, s.attrs...)
	return append(resultAttrs, attrs...)
}

func (s *slogImpl) WithAttrs(attrs ...any) Logger {
	if len(attrs) == 0 {
		return s
	}
	newAttrs := make([]any, 0, len(attrs)+len(s.attrs))
	copy(newAttrs, s.attrs)
	newAttrs = append(newAttrs, attrs...)
	return &slogImpl{
		s.logger,
		s.ctx,
		newAttrs,
	}
}

func (s *slogImpl) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return s
	}
	return &slogImpl{
		s.logger,
		ctx,
		s.attrs,
	}
}

func (s *slogImpl) Debug(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Debug(msg, s.buildAttrs(attrs...)...)
	} else {
		s.logger.DebugContext(s.ctx, msg, s.buildAttrs(attrs...)...)
	}
}

func (s *slogImpl) Debugf(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Debug(fmt.Sprintf(msg, args...), s.buildAttrs()...)
	} else {
		s.logger.DebugContext(s.ctx, fmt.Sprintf(msg, args...), s.buildAttrs()...)
	}
}

func (s *slogImpl) Info(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Info(msg, s.buildAttrs(attrs...)...)
	} else {
		s.logger.InfoContext(s.ctx, msg, s.buildAttrs(attrs...)...)
	}
}

func (s *slogImpl) Infof(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Info(fmt.Sprintf(msg, args...), s.buildAttrs()...)
	} else {
		s.logger.InfoContext(s.ctx, fmt.Sprintf(msg, args...), s.buildAttrs()...)
	}
}

func (s *slogImpl) Warn(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Warn(msg, s.buildAttrs(attrs...)...)
	} else {
		s.logger.WarnContext(s.ctx, msg, s.buildAttrs(attrs...)...)
	}
}

func (s *slogImpl) Warnf(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Warn(fmt.Sprintf(msg, args...), s.buildAttrs()...)
	} else {
		s.logger.WarnContext(s.ctx, fmt.Sprintf(msg, args...), s.buildAttrs()...)
	}
}

func (s *slogImpl) Error(msg string, attrs ...any) {
	if s.ctx == nil {
		s.logger.Error(msg, s.buildAttrs(attrs...)...)
	} else {
		s.logger.ErrorContext(s.ctx, msg, s.buildAttrs(attrs...)...)
	}
}

func (s *slogImpl) Errorf(msg string, args ...any) {
	if s.ctx == nil {
		s.logger.Error(fmt.Sprintf(msg, args...), s.buildAttrs()...)
	} else {
		s.logger.ErrorContext(s.ctx, fmt.Sprintf(msg, args...), s.buildAttrs()...)
	}
}
