package models

import (
	"os"
	"path/filepath"
)

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "TheDictationApp", "models")
}

func Path(name string) string { return filepath.Join(Dir(), name+".bin") }

func Missing(name string) bool {
	_, err := os.Stat(Path(name))
	return err != nil
}
