package util_test

import (
	"context"
	"testing"

	"github.com/teamsorghum/go-common/pkg/util"
)

func TestDownloadFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		url         string
		output      string
		expectError bool
	}{
		{
			"Success",
			"https://golang.org/robots.txt",
			"/tmp/teamsorghum-go-common-test/robots.txt",
			false,
		},
		{
			"Invalid Protocol",
			"unknown://golang.org/robots.txt",
			"/tmp/teamsorghum-go-common-test/robots.txt",
			true,
		},
		{
			"404",
			"https://golang.org/robots.txt.bak",
			"/tmp/teamsorghum-go-common-test/robots.txt",
			true,
		},
		{
			"Output permission denied",
			"https://golang.org/robots.txt",
			"/not_exist/robots.txt",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := util.DownloadFile(context.Background(), tt.url, tt.output)
			if (err == nil) == tt.expectError {
				t.Fatalf("Expect err = %+v, got %+v", tt.expectError, err)
			}
		})
	}
}
