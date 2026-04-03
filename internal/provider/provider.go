// Package provider defines common types and constants for TTS providers.
// Each specific provider (e.g., Deepgram) should implement the types.Provider interface
// defined in pkg/types/types.go.
//
// To add a new provider:
// 1. Create a new package under internal/provider/ (e.g., internal/provider/awspolly)
// 2. Implement the types.Provider interface
// 3. Add a New() constructor function that returns the provider
// 4. Update main.go to support the new provider
//
// Example:
//
//	package awspolly
//
//	import (
//	    "context"
//	    "github.com/aashish-joshi/tts-bulk/pkg/types"
//	)
//
//	type Provider struct {
//	    client *polly.Client
//	}
//
//	func New(apiKey string) (*Provider, error) {
//	    // Initialize AWS Polly client
//	    return &Provider{...}, nil
//	}
//
//	func (p *Provider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
//	    // Implementation
//	}
//
//	func (p *Provider) Name() string { return "awspolly" }
//	func (p *Provider) Close() error { return nil }
package provider

// ProviderType represents the type of TTS provider.
type ProviderType string

const (
	// ProviderDeepgram represents the Deepgram TTS provider.
	ProviderDeepgram ProviderType = "deepgram"
	// ProviderLocal represents a locally hosted OpenAI-compatible TTS provider.
	ProviderLocal ProviderType = "local"
	// ProviderOpenAI is an alias for ProviderLocal (OpenAI-compatible API).
	ProviderOpenAI ProviderType = "openai"
)
