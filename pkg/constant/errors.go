package constant

import "errors"

const (
	// ErrCodeUnknown indicates an unknown error.
	ErrCodeUnknown = 1
)

// ErrNilDeps indicates that there exists nil dependencies.
var ErrNilDeps = errors.New("nil dependencies")
