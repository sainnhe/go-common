package constant

import "errors"

// ErrNilDep indicates that the given dependency is nil.
var ErrNilDep = errors.New("nil dependency")
