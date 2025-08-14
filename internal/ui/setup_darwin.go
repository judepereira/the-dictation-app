//go:build darwin

package ui

/*
#cgo darwin CFLAGS: -x objective-c -fmodules -mmacosx-version-min=10.13
#cgo darwin LDFLAGS: -framework Cocoa
void ShowSetupWindow(void);
void SetupSetIndeterminate(_Bool indeterminate);
void SetupUpdateProgress(double percent);
void CloseSetupWindow(void);
*/
import "C"

import (
	"context"
	"dictation/internal/models"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func ShowSetupWindow() {
	C.ShowSetupWindow()
}

func SetSetupIndeterminate(ind bool) {
	C.SetupSetIndeterminate(C._Bool(ind))
}

func UpdateSetupProgress(percent float64) {
	C.SetupUpdateProgress(C.double(percent))
}

func CloseSetupWindow() {
	C.CloseSetupWindow()
}

func DownloadFileWithProgress(ctx context.Context, url, dstPath string) error {
	ShowSetupWindow()
	UpdateSetupProgress(0)

	log.Printf("Downloading %s to %s\n", url, dstPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	os.MkdirAll(models.Dir(), 0755)

	out, err := os.Create(dstPath)

	if err != nil {
		return fmt.Errorf("create %s: %w", dstPath, err)
	}
	defer func() {
		_ = out.Close()
	}()

	var total int64 = -1
	if resp.ContentLength > 0 {
		total = resp.ContentLength
		SetSetupIndeterminate(false)
	} else {
		SetSetupIndeterminate(true)
	}

	const bufSize = 256 * 1024
	buf := make([]byte, bufSize)
	var written int64
	lastUI := time.Now()

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}

		nr, rerr := resp.Body.Read(buf)
		if nr > 0 {
			nw, werr := out.Write(buf[:nr])
			if werr != nil {
				return fmt.Errorf("write: %w", werr)
			}
			if nw != nr {
				return errors.New("short write")
			}
			written += int64(nw)

			if total > 0 && time.Since(lastUI) >= 50*time.Millisecond {
				p := float64(written) / float64(total) * 100.0
				if p < 0 {
					p = 0
				} else if p > 100 {
					p = 100
				}
				UpdateSetupProgress(p)
				lastUI = time.Now()
			}
		}

		if rerr != nil {
			if rerr == io.EOF {
				break
			}
			return fmt.Errorf("read: %w", rerr)
		}
	}

	if total > 0 {
		UpdateSetupProgress(100)
	}
	time.AfterFunc(300*time.Millisecond, func() {
		CloseSetupWindow()
	})

	return nil
}
