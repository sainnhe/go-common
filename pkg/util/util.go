// Package util implements some utility functions.
package util

import (
	"fmt"
	"runtime/debug"

	"github.com/teamsorghum/go-common/pkg/log"
)

// Recover allow the program to recover from panic and print logs using [log.GetDefault]. It should be used via defer,
// otherwise panics can't be captured.
func Recover() {
	if err := recover(); err != nil {
		log.GetDefault().Errorf("Recovered from panic: %+v\n%s", err, debug.Stack())
	}
}

// ToStr converts a variable to a string. It's basically a convenient wrapper of [fmt.Sprintf] that uses %+v as
// placeholder.
func ToStr(v any) string {
	return fmt.Sprintf("%+v", v)
}
