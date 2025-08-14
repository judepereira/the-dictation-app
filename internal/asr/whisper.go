package asr

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"dictation/internal/config"
	"github.com/mutablelogic/go-whisper/sys/whisper"
)

func Run(ctx context.Context, cfg *config.Config, audioIn <-chan []float32, textOut chan<- string) {
	defer close(textOut)

	ctxParams := whisper.DefaultContextParams()
	model := whisper.Whisper_init_from_file_with_params(cfg.Model, ctxParams)
	if model == nil {
		log.Printf("ASR: failed to load model: %s", cfg.Model)
		return
	}
	defer whisper.Whisper_free(model)

	full := whisper.DefaultFullParams(whisper.SAMPLING_GREEDY)
	full.SetTranslate(false)
	if cfg.Language != "" && cfg.Language != "auto" {
		full.SetLanguage(cfg.Language)
	}

	// Keep a rolling window of audio to provide context for streaming updates.
	const (
		sampleRate = 16000           // expected by Whisper
		windowSec  = 60              // rolling window duration
		maxSamples = sampleRate * 30 // 30 seconds
	)
	_ = windowSec // retained for clarity if you adjust the window elsewhere

	var (
		mu      sync.Mutex
		samples []float32
		prev    string
		wg      sync.WaitGroup
	)

	// Periodically run decoding on the rolling window and emit only the suffix.
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				buf := make([]float32, len(samples))
				copy(buf, samples)
				mu.Unlock()

				if len(buf) == 0 {
					continue
				}

				// Run full decode using the current buffer and state.
				if err := whisper.Whisper_full(model, full, buf); err != nil {
					log.Printf("ASR: decode error: %v", err)
					continue
				}

				var b strings.Builder
				n := model.NumSegments()
				for i := 0; i < n; i++ {
					text := strings.TrimSpace(model.Segment(i).Text)
					if text[:1] == "[" {
						samples = samples[0:]
						continue
					}
					b.WriteString(text + " ")
				}
				curr := strings.TrimSpace(b.String()) + " "

				if prev != curr {
					select {
					case textOut <- curr:
						prev = curr
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case chunk, ok := <-audioIn:
			if !ok {
				wg.Wait()
				return
			}
			if len(chunk) == 0 {
				continue
			}

			mu.Lock()
			if len(samples)+len(chunk) > maxSamples {
				drop := len(samples) + len(chunk) - maxSamples
				if drop > len(samples) {
					drop = len(samples)
				}
				samples = append(samples[drop:], chunk...)
			} else {
				samples = append(samples, chunk...)
			}
			mu.Unlock()
		}
	}
}
