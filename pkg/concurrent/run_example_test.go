// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package concurrent_test

import (
	"fmt"
	"time"

	"github.com/teamsorghum/go-common/pkg/concurrent"
)

func ExampleRun() {
	// Define Arg and Result struct.
	type Arg struct {
		x int
		y int
	}
	type Result int

	// Define concurrency.
	concurrency := int32(3)

	// Define arguments.
	args := []Arg{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10},
	}

	// Define task handler function.
	f := func(arg Arg) Result {
		time.Sleep(100 * time.Millisecond)
		return Result(arg.x + arg.y)
	}

	// Run tasks concurrently.
	results := concurrent.Run(concurrency, args, f)
	fmt.Println(results)

	// Output: [3 7 11 15 19]
}
