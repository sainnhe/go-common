package util_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/teamsorghum/go-common/pkg/util"
)

func TestConcurrentRun(t *testing.T) {
	t.Parallel()

	// Define Arg and Result types
	type Arg int
	type Result int

	// Test cases
	tests := []struct {
		name        string
		concurrency int32
		args        []Arg
		f           func(Arg) Result
		want        []Result
	}{
		{
			name:        "Simple addition",
			concurrency: 2,
			args:        []Arg{1, 2, 3, 4, 5},
			f: func(a Arg) Result {
				return Result(a + 10)
			},
			want: []Result{11, 12, 13, 14, 15},
		},
		{
			name:        "Square numbers",
			concurrency: 3,
			args:        []Arg{1, 2, 3, 4},
			f: func(a Arg) Result {
				return Result(a * a)
			},
			want: []Result{1, 4, 9, 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.ConcurrentRun(tt.concurrency, tt.args, tt.f)
			if len(got) != len(tt.want) {
				t.Fatalf("Want results of length %d, got %d", len(tt.want), len(got))
			}
			for i, res := range got {
				if res != tt.want[i] {
					t.Errorf("Want result[%d] = %v, got %v", i, tt.want[i], res)
				}
			}
		})
	}
}

func ExampleConcurrentRun() {
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
	results := util.ConcurrentRun(concurrency, args, f)
	fmt.Println(results)

	// Output: [3 7 11 15 19]
}
