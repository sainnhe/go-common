package util_test

import (
	"errors"
	"testing"

	"github.com/teamsorghum/go-common/pkg/util"
)

func TestToStr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "int", input: 42, want: "42"},
		{name: "string", input: "hello", want: "hello"},
		{name: "struct", input: struct{ Name string }{Name: "Alice"}, want: "{Name:Alice}"},
		{name: "slice", input: []int{1, 2, 3}, want: "[1 2 3]"},
		{name: "map", input: map[string]int{"a": 1, "b": 2}, want: "map[a:1 b:2]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.ToStr(tt.input)
			if got != tt.want {
				t.Errorf("ToStr(%v) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMergeErrors(t *testing.T) {
	t.Parallel()

	err1 := errors.New("error one")
	err2 := errors.New("error two")
	err3 := errors.New("error three")

	tests := []struct {
		name string
		errs []error
		want error
	}{
		{
			name: "all nil errors",
			errs: []error{nil, nil, nil},
			want: nil,
		},
		{
			name: "no errors",
			errs: []error{},
			want: nil,
		},
		{
			name: "single error",
			errs: []error{err1},
			want: err1,
		},
		{
			name: "multiple errors",
			errs: []error{err1, err2, err3},
			want: errors.New("error one; error two; error three"),
		},
		{
			name: "mixed nil and non-nil errors",
			errs: []error{err1, nil, err3},
			want: errors.New("error one; error three"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.MergeErrors(tt.errs...)
			if tt.want == nil && got != nil {
				t.Errorf("%s: expected nil error, got %q", tt.name, got)
			} else if tt.want != nil {
				if got == nil {
					t.Errorf("%s: expected error %q, got nil", tt.name, tt.want)
				} else if got.Error() != tt.want.Error() {
					t.Errorf("%s: expected error %q, got %q", tt.name, tt.want.Error(), got.Error())
				}
			}
		})
	}
}

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
			name:        "simple addition",
			concurrency: 2,
			args:        []Arg{1, 2, 3, 4, 5},
			f: func(a Arg) Result {
				return Result(a + 10)
			},
			want: []Result{11, 12, 13, 14, 15},
		},
		{
			name:        "square numbers",
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

			results := util.ConcurrentRun(tt.concurrency, tt.args, tt.f)
			if len(results) != len(tt.want) {
				t.Fatalf("%s: Expected results of length %d, got %d", tt.name, len(tt.want), len(results))
			}
			for i, res := range results {
				if res != tt.want[i] {
					t.Errorf("%s: Expected result[%d] = %v, got %v", tt.name, i, tt.want[i], res)
				}
			}
		})
	}
}

func TestRecover(t *testing.T) {
	t.Parallel()

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Test panicked: %v", r)
			}
		}()

		func() {
			defer util.Recover()
			panic("test panic")
		}()
		// If the panic is recovered, execution continues here
	}()
	// Reaching this point means the panic was recovered
}
