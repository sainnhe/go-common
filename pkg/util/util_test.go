package util_test

import (
	"errors"
	"testing"

	"github.com/teamsorghum/go-common/pkg/util"
)

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

func TestToStr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "Int", input: 42, want: "42"},
		{name: "String", input: "hello", want: "hello"},
		{name: "Struct", input: struct{ Name string }{Name: "Alice"}, want: "{Name:Alice}"},
		{name: "Slice", input: []int{1, 2, 3}, want: "[1 2 3]"},
		{name: "Map", input: map[string]int{"a": 1, "b": 2}, want: "map[a:1 b:2]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.ToStr(tt.input)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
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
			name: "All nil errors",
			errs: []error{nil, nil, nil},
			want: nil,
		},
		{
			name: "No errors",
			errs: []error{},
			want: nil,
		},
		{
			name: "Single error",
			errs: []error{err1},
			want: err1,
		},
		{
			name: "Multiple errors",
			errs: []error{err1, err2, err3},
			want: errors.New("error one; error two; error three"),
		},
		{
			name: "Mixed nil and non-nil errors",
			errs: []error{err1, nil, err3},
			want: errors.New("error one; error three"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.MergeErrors(tt.errs...)
			if tt.want == nil && got != nil {
				t.Errorf("Want nil error, got %q", got)
			} else if tt.want != nil {
				if got == nil {
					t.Errorf("Want error %q, got nil", tt.want)
				} else if got.Error() != tt.want.Error() {
					t.Errorf("Want error %q, got %q", tt.want.Error(), got.Error())
				}
			}
		})
	}
}
