package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// NewLight initializes a light logger.
func NewLight() Logger {
	slogLogger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelDebug,
		TimeFormat: time.StampMilli,
		NoColor:    false,
	}))
	return &slogImpl{
		slogLogger,
		slogLogger,
		nil,
		[]any{},
	}
}
