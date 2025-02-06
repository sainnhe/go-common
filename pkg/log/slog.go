package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	ctxutil "github.com/teamsorghum/go-common/pkg/util/ctx"
)

type slogImpl struct {
	logger         *slog.Logger
	originalLogger *slog.Logger
	ctx            context.Context
	attrs          []any
}

// NewSlog initializes a slog based logger.
func NewSlog(cfg *Config) (logger Logger, cleanup func(), err error) {
	var fileMode fs.FileMode = 0666
	logFile, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fileMode)
	if err != nil {
		return
	}
	cleanup = func() {
		if logFile != nil {
			logFile.Close()
		}
	}
	consoleWriter := os.Stderr
	fileWriter := io.Writer(logFile)
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
		slogLogger,
		nil,
		[]any{},
	}
	logger.Info("Slog initialized.", "config", cfg)
	return
}

func buildAttrs(ctx context.Context, attrs []any) []any {
	ctxAttrs := []any{}
	fields := ctxutil.GetContextFields(ctx)
	for k, v := range fields {
		ctxAttrs = append(ctxAttrs, fmt.Sprintf("ctx_%+v", k), v)
	}
	return append(ctxAttrs, attrs...)
}

func (s *slogImpl) WithAttrs(attrs ...any) Logger {
	if len(attrs) == 0 {
		return s
	}
	newAttrs := []any{}
	copy(newAttrs, s.attrs)
	newAttrs = append(newAttrs, attrs...)
	return &slogImpl{
		s.originalLogger.With(buildAttrs(s.ctx, newAttrs)...),
		s.originalLogger,
		s.ctx,
		newAttrs,
	}
}

func (s *slogImpl) WithContext(ctx context.Context) Logger {
	return &slogImpl{
		s.originalLogger.With(buildAttrs(ctx, s.attrs)...),
		s.originalLogger,
		ctx,
		s.attrs,
	}
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
