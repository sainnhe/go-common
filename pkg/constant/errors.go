// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package constant

import "errors"

// ErrNilDeps indicates that there exists nil dependencies.
var ErrNilDeps = errors.New("nil dependencies")
