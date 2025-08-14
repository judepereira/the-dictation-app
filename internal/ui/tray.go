package ui

import (
	"log"
	"os"
	"sync"

	"github.com/getlantern/systray"

	"dictation/internal/config"
)

var (
	onQuit   func()
	trayOnce sync.Once
	quitOnce sync.Once
)

func StartTray(cfg *config.Config) {
	trayOnce.Do(func() {
		systray.Run(func() {
			if iconBytes, err := os.ReadFile("assets/icon.png"); err == nil {
				systray.SetTemplateIcon(iconBytes, iconBytes)
			} else {
				log.Printf("Failed to load tray icon: %v", err)
			}

			systray.SetTitle("Dictation")
			systray.SetTooltip("Dictation is running")

			statusModel := systray.AddMenuItem("Model: "+cfg.Model, "")
			statusModel.Disable()
			statusLang := systray.AddMenuItem("Language: "+cfg.Language, "")
			statusLang.Disable()

			systray.AddSeparator()

			itemSetup := systray.AddMenuItem("Setup…", "Open first-time setup")
			itemAbout := systray.AddMenuItem("About", "About this app")

			systray.AddSeparator()

			itemQuit := systray.AddMenuItem("Quit", "Quit the application")

			go func() {
				for {
					select {
					case <-itemSetup.ClickedCh:
						ShowSetup()
					case <-itemAbout.ClickedCh:
						ShowAbout()
					case <-itemQuit.ClickedCh:
						if onQuit != nil {
							go onQuit()
						}
						systray.Quit()
						return
					}
				}
			}()

			log.Printf("Tray started. Lang=%s Model=%s\n", cfg.Language, cfg.Model)
		}, func() {
			log.Println("Tray stopped.")
		})
	})
}

func StopTray() { quitOnce.Do(func() { systray.Quit() }) }

func ShowAbout() { showAboutNative() }

func ShowSetup() { showSetupNative() }
