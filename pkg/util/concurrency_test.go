package util_test

import (
	"testing"

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
