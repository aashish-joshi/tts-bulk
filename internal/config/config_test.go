package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		csvPath     string
		outputDir   string
		format      string
		model       string
		provider    string
		providerURL string
		voice       string
		rateLimitMs int
		apiKey      string
		envProvider string
		wantErr     bool
	}{
		{
			name:        "valid deepgram configuration",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "mp3",
			model:       "test-model",
			provider:    "deepgram",
			rateLimitMs: -1,
			apiKey:      "test-key",
			wantErr:     false,
		},
		{
			name:        "valid wav format",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "wav",
			model:       "test-model",
			provider:    "deepgram",
			rateLimitMs: -1,
			apiKey:      "test-key",
			wantErr:     false,
		},
		{
			name:        "invalid format",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "ogg",
			model:       "test-model",
			provider:    "deepgram",
			rateLimitMs: -1,
			apiKey:      "test-key",
			wantErr:     true,
		},
		{
			name:        "missing API key for deepgram",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "mp3",
			model:       "test-model",
			provider:    "deepgram",
			rateLimitMs: -1,
			apiKey:      "",
			wantErr:     true,
		},
		{
			name:        "local provider no API key required",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "mp3",
			model:       "tts-1",
			provider:    "local",
			providerURL: "http://localhost:8000",
			voice:       "alloy",
			rateLimitMs: 0,
			apiKey:      "",
			wantErr:     false,
		},
		{
			name:        "openai alias no API key required",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "mp3",
			model:       "tts-1",
			provider:    "openai",
			providerURL: "http://localhost:8000",
			voice:       "alloy",
			rateLimitMs: 0,
			apiKey:      "",
			wantErr:     false,
		},
		{
			name:        "provider from TTS_PROVIDER env var",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "mp3",
			model:       "tts-1",
			provider:    "",
			providerURL: "http://localhost:8000",
			voice:       "alloy",
			rateLimitMs: 0,
			envProvider: "local",
			apiKey:      "",
			wantErr:     false,
		},
		{
			name:        "default provider is deepgram",
			csvPath:     "test.csv",
			outputDir:   "output",
			format:      "mp3",
			model:       "test-model",
			provider:    "",
			rateLimitMs: -1,
			apiKey:      "test-key",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars before each test
			os.Unsetenv("DEEPGRAM_API_KEY")
			os.Unsetenv("TTS_PROVIDER")
			os.Unsetenv("OPENAI_API_KEY")

			if tt.apiKey != "" {
				os.Setenv("DEEPGRAM_API_KEY", tt.apiKey)
			}
			if tt.envProvider != "" {
				os.Setenv("TTS_PROVIDER", tt.envProvider)
			}

			cfg, err := Load(tt.csvPath, tt.outputDir, tt.format, tt.model, tt.provider, tt.providerURL, tt.voice, tt.rateLimitMs)

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

			if cfg.ProviderURL != tt.providerURL {
				t.Errorf("Load() ProviderURL = %v, want %v", cfg.ProviderURL, tt.providerURL)
			}

			if cfg.Voice != tt.voice {
				t.Errorf("Load() Voice = %v, want %v", cfg.Voice, tt.voice)
			}

			if cfg.RateLimitMs != tt.rateLimitMs {
				t.Errorf("Load() RateLimitMs = %v, want %v", cfg.RateLimitMs, tt.rateLimitMs)
			}
		})
	}
}
