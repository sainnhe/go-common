// Package util implements some utility functions.
package util

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/teamsorghum/go-common/pkg/log"
)

// Recover allow the program to recover from panic and print logs. It should be used via defer.
func Recover() {
	if err := recover(); err != nil && log.DefaultLogger != nil {
		log.DefaultLogger.Errorf("Recovered from panic. Error = %+v\n%s", err, debug.Stack())
	}
}

// ToStr converts a variable to a string.
func ToStr(v any) string {
	return fmt.Sprintf("%+v", v)
}

/*
MergeErrors merges multiple errors and returns a merged error.

The logic of merge is:
  - If all the incoming errors are nil, return nil.
  - If 1 or more of the incoming errors are not nil, then merge the error message and return.

Params:
  - errs ...error: Errors to be merged.

Returns:
  - error: Merged error.
*/
func MergeErrors(errs ...error) error {
	nonNilErrors := []error{}
	for _, err := range errs {
		if err != nil {
			nonNilErrors = append(nonNilErrors, err)
		}
	}
	switch len(nonNilErrors) {
	case 0:
		return nil
	default:
		msgs := make([]string, len(nonNilErrors))
		for i, err := range nonNilErrors {
			msgs[i] = err.Error()
		}
		return errors.New(strings.Join(msgs, "; "))
	}
}
