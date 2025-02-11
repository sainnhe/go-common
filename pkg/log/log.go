// Package log implements a common logging interface.
package log

import (
	"errors"
	"sync"
)

var defaultLogger Logger
var mu sync.Mutex

// ProvideLogger initializes a logger based on Config and provides logger dependency. It'll also set this logger as the
// default logger.
func ProvideLogger(cfg *Config) (logger Logger, cleanup func(), err error) {
	if cfg == nil {
		return nil, func() {}, errors.New("Config is nil")
	}
	switch cfg.Type {
	case "light":
		logger, cleanup, err = NewLight(), nil, nil
	case "slog":
		logger, cleanup, err = NewSlog(cfg)
	case "loki":
		logger, cleanup, err = NewLoki(cfg)
	default:
		err = errors.New("invalid logger type")
	}
	SetDefault(logger)
	return
}

// GetDefault gets the default logger.
func GetDefault() Logger {
	mu.Lock()
	defer mu.Unlock()
	// Return default logger if it exists.
	if defaultLogger != nil {
		return defaultLogger
	}
	// Otherwise set default logger and return it.
	defaultLogger = NewLight()
	return defaultLogger
}

// SetDefault sets the default logger.
func SetDefault(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	if logger == nil {
		defaultLogger = NewLight()
	} else {
		defaultLogger = logger
	}
}
