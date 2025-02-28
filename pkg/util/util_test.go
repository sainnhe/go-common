package util_test

import (
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
