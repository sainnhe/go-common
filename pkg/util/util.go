// Package util implements some utility functions.
package util

import (
	"fmt"
	"runtime/debug"

	"github.com/teamsorghum/go-common/pkg/log"
)

// Recover allow the program to recover from panic and print logs. It should be used via defer.
func Recover() {
	if err := recover(); err != nil {
		log.GetDefault().Errorf("Recovered from panic: %+v\n%s", err, debug.Stack())
	}
}

// ToStr converts a variable to a string.
func ToStr(v any) string {
	return fmt.Sprintf("%+v", v)
}
