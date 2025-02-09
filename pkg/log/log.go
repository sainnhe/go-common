// Package log implements a common logging interface.
package log

import "errors"

var defaultLogger Logger

// ProvideLogger initializes a logger based on Config and provides logger dependency. It'll also set this logger as the
// default logger.
func ProvideLogger(cfg *Config) (logger Logger, cleanup func(), err error) {
	if cfg == nil {
		return nil, func() {}, errors.New("Config is nil")
	}
	switch cfg.Type {
	case "light":
		logger, cleanup, err = NewLight(), func() {}, nil
	case "slog":
		logger, cleanup, err = NewSlog(cfg)
	default:
		err = errors.New("invalid logger type")
	}
	SetDefault(logger)
	return
}

// GetDefault gets the default logger.
func GetDefault() Logger {
	// Return default logger if it exists.
	if defaultLogger != nil {
		return defaultLogger
	}
	// Otherwise set default logger and return it.
	SetDefault(nil)
	return GetDefault()
}

// SetDefault sets the default logger.
func SetDefault(logger Logger) {
	if logger == nil {
		defaultLogger = NewLight()
	} else {
		defaultLogger = logger
	}
}
