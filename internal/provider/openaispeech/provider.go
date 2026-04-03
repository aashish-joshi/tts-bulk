// Package openaispeech implements an OpenAI Speech API-compatible TTS provider.
// It works with any locally hosted server that exposes the OpenAI Speech API
// (e.g. Kokoro-FastAPI, OpenedAI Speech, AllTalk).
package openaispeech

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

const (
	defaultBaseURL     = "http://localhost:8000"
	defaultModel       = "tts-1"
	defaultVoice       = "alloy"
	defaultFormat      = "mp3"
	speechEndpoint     = "/v1/audio/speech"
)

// Config holds configuration for the OpenAI-compatible TTS provider.
type Config struct {
	// BaseURL is the base URL of the OpenAI-compatible TTS server.
	BaseURL string
	// APIKey is the optional API key. Not required for local servers.
	APIKey string
	// Model is the TTS model name (e.g. "tts-1").
	Model string
	// Voice is the voice to use (e.g. "alloy", "echo", "nova").
	Voice string
	// Format is the audio format: "mp3" or "wav".
	Format string
	// RateLimitMs is the delay in milliseconds between requests.
	// Use 0 for no rate limiting.
	RateLimitMs int
}

// Provider implements types.Provider for OpenAI-compatible TTS servers.
type Provider struct {
	cfg    Config
	client *http.Client
}

// speechRequest mirrors the OpenAI Speech API request body.
type speechRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	Voice          string `json:"voice"`
	ResponseFormat string `json:"response_format"`
}

// New creates a new OpenAI-compatible TTS provider.
func New(cfg Config) (*Provider, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	if cfg.Model == "" {
		cfg.Model = defaultModel
	}
	if cfg.Voice == "" {
		cfg.Voice = defaultVoice
	}
	if cfg.Format == "" {
		cfg.Format = defaultFormat
	}

	return &Provider{
		cfg:    cfg,
		client: &http.Client{},
	}, nil
}

// GenerateAudio calls the OpenAI-compatible speech endpoint and writes the
// audio response to outputPath.
func (p *Provider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
	if req.Text == "" {
		return fmt.Errorf("text cannot be empty")
	}
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	body, err := json.Marshal(speechRequest{
		Model:          p.cfg.Model,
		Input:          req.Text,
		Voice:          p.cfg.Voice,
		ResponseFormat: p.cfg.Format,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.BaseURL+speechEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if p.cfg.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("TTS server returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(errBody)))
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := os.WriteFile(outputPath, audioData, 0o600); err != nil {
		return fmt.Errorf("failed to write audio file: %w", err)
	}

	if p.cfg.RateLimitMs > 0 {
		time.Sleep(time.Duration(p.cfg.RateLimitMs) * time.Millisecond)
	}

	return nil
}

// Name returns the provider's name.
func (p *Provider) Name() string {
	return "openaispeech"
}

// Close is a no-op for this provider.
func (p *Provider) Close() error {
	return nil
}
