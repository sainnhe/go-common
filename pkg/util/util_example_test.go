package util_test

import (
	"fmt"

	"github.com/teamsorghum/go-common/pkg/util"
)

func ExampleRecover() {
	// Let's run a function that will panic and see if it'll be recovered.
	func() {
		// Use defer to capture panics.
		defer util.Recover()

		// Manually trigger a panic.
		panic("test panic")
	}()

	// If recovery from panic is successful, this line will be printed.
	fmt.Println("Recovered.")

	// Output: Recovered.
}
