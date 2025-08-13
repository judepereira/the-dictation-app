// internal/setup/downloader.go
package setup

import (
	"context"
	"io"
	"net/http"
	"os"
)

type Progress struct {
	Bytes int64
	Total int64
}

func Download(ctx context.Context, url, dest string, total int64, prog chan<- Progress) error {
	tmp := dest + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer out.Close()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var written int64
	buf := make([]byte, 64*1024)
	for {
		n, rErr := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := out.Write(buf[:n]); wErr != nil {
				return wErr
			}
			written += int64(n)
			select {
			case prog <- Progress{Bytes: written, Total: total}:
			default:
			}
		}
		if rErr == io.EOF {
			break
		}
		if rErr != nil {
			return rErr
		}
	}
	out.Close()
	return os.Rename(tmp, dest)
}
