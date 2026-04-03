// Package deepgram implements the Deepgram TTS provider.
package deepgram

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aashish-joshi/tts-bulk/pkg/types"
	api "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
)

const defaultRateLimitMs = 1000

// Provider implements the TTS provider interface for Deepgram.
type Provider struct {
	client      *api.Client
	apiKey      string
	model       string
	options     *interfaces.SpeakOptions
	rateLimitMs int
}

// Config holds Deepgram-specific configuration.
type Config struct {
	APIKey      string
	Model       string
	Format      string
	RateLimitMs int
}

// formatToContainerEncoding maps an audio format string to the Deepgram
// container and encoding values required by the API.
func formatToContainerEncoding(format string) (container, encoding string) {
	if strings.ToLower(format) == "wav" {
		return "wav", "linear16"
	}
	return "", "mp3"
}

// New creates a new Deepgram TTS provider.
func New(cfg Config) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Set API key in environment for Deepgram SDK
	os.Setenv("DEEPGRAM_API_KEY", cfg.APIKey)

	// Initialize the Deepgram SDK
	client.Init(client.InitLib{
		LogLevel: client.LogLevelErrorOnly,
	})

	// Create REST client
	c := client.NewRESTWithDefaults()
	dgClient := api.New(c)

	container, encoding := formatToContainerEncoding(cfg.Format)

	options := &interfaces.SpeakOptions{
		Model:     strings.ToLower(cfg.Model),
		Container: container,
		Encoding:  encoding,
	}

	rateLimitMs := cfg.RateLimitMs
	if rateLimitMs < 0 {
		rateLimitMs = defaultRateLimitMs
	}

	return &Provider{
		client:      dgClient,
		apiKey:      cfg.APIKey,
		model:       cfg.Model,
		options:     options,
		rateLimitMs: rateLimitMs,
	}, nil
}

// GenerateAudio generates audio from text and saves it to the specified path.
func (p *Provider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
	if req.Text == "" {
		return fmt.Errorf("text cannot be empty")
	}

	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Generate TTS and save to file
	_, err := p.client.ToSave(ctx, outputPath, req.Text, p.options)
	if err != nil {
		return fmt.Errorf("failed to generate TTS: %w", err)
	}

	// Rate limiting
	if p.rateLimitMs > 0 {
		time.Sleep(time.Duration(p.rateLimitMs) * time.Millisecond)
	}

	return nil
}

// Name returns the provider's name.
func (p *Provider) Name() string {
	return "deepgram"
}

// Close cleans up any resources used by the provider.
func (p *Provider) Close() error {
	// Deepgram SDK doesn't require explicit cleanup
	return nil
}
