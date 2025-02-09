//go:generate mockgen -typed -write_package_comment=false -source=api.go -destination=api_mock.go -package log

package log

import "context"

// Logger defines a common logging interface.
type Logger interface {
	// WithAttrs returns a Logger that includes the given attributes in each output operation.
	WithAttrs(attrs ...any) Logger
	// WithContext returns a Logger that includes the given context in each output operation.
	WithContext(ctx context.Context) Logger
	// Debug outputs a debug level message with additional attributes.
	Debug(msg string, attrs ...any)
	// Debugf outputs a debug level message of a formatted string.
	Debugf(msg string, args ...any)
	// Info outputs a info level message with additional attributes.
	Info(msg string, attrs ...any)
	// Infof outputs a info level message of a formatted string.
	Infof(msg string, args ...any)
	// Warn outputs a warn level message with additional attributes.
	Warn(msg string, attrs ...any)
	// Warnf outputs a warn level message of a formatted string.
	Warnf(msg string, args ...any)
	// Error outputs a error level message with additional attributes.
	Error(msg string, attrs ...any)
	// Errorf outputs a error level message of a formatted string.
	Errorf(msg string, args ...any)
	// Fatal outputs a fatal level message with additional attributes.
	Fatal(msg string, attrs ...any)
	// Fatalf outputs a fatal level message of a formatted string.
	Fatalf(msg string, args ...any)
}
