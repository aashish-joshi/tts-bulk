package config

import (
	"os"
	"testing"

	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		csvPath   string
		outputDir string
		format    string
		model     string
		apiKey    string
		wantErr   bool
	}{
		{
			name:      "valid configuration",
			csvPath:   "test.csv",
			outputDir: "output",
			format:    "mp3",
			model:     "test-model",
			apiKey:    "test-key",
			wantErr:   false,
		},
		{
			name:      "valid wav format",
			csvPath:   "test.csv",
			outputDir: "output",
			format:    "wav",
			model:     "test-model",
			apiKey:    "test-key",
			wantErr:   false,
		},
		{
			name:      "invalid format",
			csvPath:   "test.csv",
			outputDir: "output",
			format:    "ogg",
			model:     "test-model",
			apiKey:    "test-key",
			wantErr:   true,
		},
		{
			name:      "missing API key",
			csvPath:   "test.csv",
			outputDir: "output",
			format:    "mp3",
			model:     "test-model",
			apiKey:    "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set API key environment variable
			if tt.apiKey != "" {
				os.Setenv("DEEPGRAM_API_KEY", tt.apiKey)
			} else {
				os.Unsetenv("DEEPGRAM_API_KEY")
			}

			cfg, err := Load(tt.csvPath, tt.outputDir, tt.format, tt.model)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error: %v", err)
				return
			}

			if cfg.CSVPath != tt.csvPath {
				t.Errorf("Load() CSVPath = %v, want %v", cfg.CSVPath, tt.csvPath)
			}

			if cfg.Model != tt.model {
				t.Errorf("Load() Model = %v, want %v", cfg.Model, tt.model)
			}
		})
	}
}

func TestGetDeepgramConfig(t *testing.T) {
	tests := []struct {
		name          string
		audioFormat   types.AudioFormat
		wantContainer string
		wantEncoding  string
	}{
		{
			name:          "mp3 format",
			audioFormat:   types.FormatMP3,
			wantContainer: "",
			wantEncoding:  "mp3",
		},
		{
			name:          "wav format",
			audioFormat:   types.FormatWAV,
			wantContainer: "wav",
			wantEncoding:  "linear16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				AudioFormat: tt.audioFormat,
			}

			container, encoding := cfg.GetDeepgramConfig()

			if container != tt.wantContainer {
				t.Errorf("GetDeepgramConfig() container = %v, want %v", container, tt.wantContainer)
			}

			if encoding != tt.wantEncoding {
				t.Errorf("GetDeepgramConfig() encoding = %v, want %v", encoding, tt.wantEncoding)
			}
		})
	}
}
