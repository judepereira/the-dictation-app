// internal/setup/checks.go
package setup

import (
	"errors"
	"os"
)

func HasEnoughDisk(bytesNeeded int64, path string) error {
	// TODO: implement statfs; return nil for now.
	if bytesNeeded <= 0 {
		return errors.New("bytesNeeded must be positive")
	}
	_, err := os.Stat(path)
	return err
}

func ModelRequired(missing bool) bool { return missing }
