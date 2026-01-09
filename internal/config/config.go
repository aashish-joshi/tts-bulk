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
	Model        string

	// Input/Output settings
	CSVPath     string
	OutputDir   string
	AudioFormat types.AudioFormat

	// Processing settings
	MaxConcurrent int
	RetryAttempts int
}

// Load loads configuration from environment variables and command-line flags.
func Load(csvPath, outputDir, format, model string) (*Config, error) {
	cfg := &Config{
		ProviderType:  "deepgram", // Default provider
		CSVPath:       csvPath,
		OutputDir:     outputDir,
		Model:         model,
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

	// Get API key from environment
	cfg.APIKey = os.Getenv("DEEPGRAM_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("DEEPGRAM_API_KEY environment variable is not set")
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

// GetDeepgramConfig returns Deepgram-specific configuration.
func (c *Config) GetDeepgramConfig() (container, encoding string) {
	if c.AudioFormat == types.FormatWAV {
		return "wav", "linear16"
	}
	return "", "mp3"
}
