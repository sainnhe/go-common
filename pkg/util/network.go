package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

// DownloadFile downloads a file from the given url, with supports for progress bar, context canceling, redirecting.
func DownloadFile(ctx context.Context, url, output string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close() // nolint:errcheck

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("response: %+v", rsp)
	}

	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0600) // nolint:gosec,mnd
	if err != nil {
		return err
	}
	defer f.Close() // nolint:errcheck

	bar := progressbar.DefaultBytes(
		rsp.ContentLength,
		output,
	)
	_, err = io.Copy(io.MultiWriter(f, bar), rsp.Body)
	return err
}
