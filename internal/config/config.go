// Package config handles application configuration.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

// Config holds the application configuration.
type Config struct {
	// Provider settings
	ProviderType string
	APIKey       string
	ProviderURL  string
	Voice        string
	Model        string
	RateLimitMs  int

	// Input/Output settings
	CSVPath     string
	OutputDir   string
	AudioFormat types.AudioFormat

	// Processing settings
	MaxConcurrent int
	RetryAttempts int
}

// Load loads configuration from environment variables and command-line flags.
func Load(csvPath, outputDir, format, model, provider, providerURL, voice string, rateLimitMs int) (*Config, error) {
	// Provider selection: param → TTS_PROVIDER env var → default "deepgram"
	if provider == "" {
		provider = os.Getenv("TTS_PROVIDER")
	}
	if provider == "" {
		provider = "deepgram"
	}

	cfg := &Config{
		ProviderType:  provider,
		ProviderURL:   providerURL,
		Voice:         voice,
		CSVPath:       csvPath,
		OutputDir:     outputDir,
		Model:         model,
		RateLimitMs:   rateLimitMs,
		MaxConcurrent: 3,
		RetryAttempts: 2,
	}

	// Validate and set audio format
	formatLower := strings.ToLower(format)
	switch formatLower {
	case "mp3":
		cfg.AudioFormat = types.FormatMP3
	case "wav":
		cfg.AudioFormat = types.FormatWAV
	default:
		return nil, fmt.Errorf("unsupported audio format: %s (supported: mp3, wav)", format)
	}

	// Get API key from environment based on provider
	switch strings.ToLower(provider) {
	case "deepgram":
		cfg.APIKey = os.Getenv("DEEPGRAM_API_KEY")
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("DEEPGRAM_API_KEY environment variable is not set")
		}
	default:
		// For local/openai-compatible providers the API key is optional
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
	}

	// Validate required fields
	if cfg.CSVPath == "" {
		return nil, fmt.Errorf("CSV path is required")
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = "audio"
	}

	if cfg.Model == "" {
		cfg.Model = "aura-asteria-en"
	}

	return cfg, nil
}
