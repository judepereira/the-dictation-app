// internal/config/config.go
package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Model          string `json:"model"`
	Language       string `json:"language"`
	AutoType       bool   `json:"auto_type"`
	HotkeyWindowMs int    `json:"hotkey_window_ms"`
	VADEnabled     bool   `json:"vad_enabled"`
	Version        string `json:"version"`
}

func path() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "DictationApp", "config.json")
}

func defaultConfig() *Config {
	return &Config{Model: "/Users/jude/developer/whisper.cpp/models/ggml-base.en.bin", Language: "en", AutoType: true, HotkeyWindowMs: 300, VADEnabled: true, Version: "0.1.0"}
}

func Load() (*Config, error) {
	p := path()
	b, err := os.ReadFile(p)
	log.Printf("Config file: %s", p)
	if err != nil {
		cfg := defaultConfig()
		_ = Save(cfg)
		return cfg, nil
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		cfg = *defaultConfig()
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	p := path()
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(p, b, 0o644)
}
