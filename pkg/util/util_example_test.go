package util_test

import (
	"fmt"

	"github.com/sainnhe/go-common/pkg/util"
)

func ExampleRecover() {
	// Let's run a function that will panic and see if it'll be recovered.
	func() {
		// Use defer to capture panics.
		defer util.Recover()

		// Manually trigger a panic.
		panic("test panic")
	}()

	// If recovery from panic is successful, this message will be printed.
	fmt.Println("Recovered.")

	// Output: Recovered.
}
