// cmd/app/main.go
package main

import "C"
import (
	"context"
	"dictation/internal/models"
	"log"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"dictation/internal/asr"
	"dictation/internal/audio"
	"dictation/internal/config"
	"dictation/internal/hotkey"
	"dictation/internal/keyout"
	"dictation/internal/ui"
)

func main() {
	cfg, _ := config.Load()

	if !keyout.EnsureAccessibility() {
		log.Println("Accessibility not granted.")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if models.Missing(cfg.Model) {
		err := InstallModel(ctx,
			"https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-"+cfg.Model+".bin",
			models.Path(cfg.Model))

		if err != nil {
			log.Printf("Failed to install model: %v", err)
			return
		}
	}

	toggleCh := make(chan bool, 1)

	audioCtx, audioCancel := context.WithCancel(ctx)
	var capturing atomic.Bool

	go hotkey.Run(toggleCh)
	go func() {
		for toggle := range toggleCh {
			if toggle && !capturing.Swap(true) {
				audioCtx, audioCancel = context.WithCancel(ctx)
				audioCh := make(chan []float32, 8)
				textCh := make(chan string, 8)
				keyout.Start()
				go audio.Capture(audioCtx, audioCh)
				go asr.Run(audioCtx, cfg, audioCh, textCh)
				go func() {
					for t := range textCh {
						keyout.TypeString(t)
					}
				}()
			} else if !toggle && capturing.Swap(false) {
				audioCancel()
				keyout.Stop()
			}
		}
	}()

	ui.StartTray(cfg)

	<-ctx.Done()
	ui.StopTray()
}

func InstallModel(ctx context.Context, url, dst string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	return ui.DownloadFileWithProgress(ctx, url, dst)
}
