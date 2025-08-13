// internal/ui/tray.go
package ui

import (
	"log"

	"dictation/internal/config"
)

var onQuit func()

func StartTray(cfg *config.Config) {
	// TODO: integrate github.com/getlantern/systray.
	log.Printf("Tray started. Lang=%s Model=%s\n", cfg.Language, cfg.Model)
}

func StopTray() { log.Println("Tray stopped.") }

func ShowAbout() { showAboutNative() }

func ShowSetup() { showSetupNative() }
