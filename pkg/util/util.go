// Package util implements some utility functions.
package util

import (
	"runtime/debug"

	"github.com/teamsorghum/go-common/pkg/log"
)

// Recover allow the program to recover from panic and print logs. It should be used via defer.
func Recover() {
	if err := recover(); err != nil && log.DefaultLogger != nil {
		log.DefaultLogger.Errorf("Recovered from panic. Error = %+v\n%s", err, debug.Stack())
	}
}
