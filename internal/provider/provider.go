// Package provider defines the TTS provider interface and factory.
package provider

import (
	"fmt"

	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

// ProviderType represents the type of TTS provider.
type ProviderType string

const (
	// ProviderDeepgram represents the Deepgram TTS provider.
	ProviderDeepgram ProviderType = "deepgram"
)

// Config holds configuration for creating a provider.
type Config struct {
	Type   ProviderType
	APIKey string
	Model  string
}

// Factory creates TTS providers based on configuration.
type Factory interface {
	Create(config Config) (types.Provider, error)
}

type factory struct{}

// NewFactory creates a new provider factory.
func NewFactory() Factory {
	return &factory{}
}

// Create creates a new TTS provider based on the configuration.
func (f *factory) Create(config Config) (types.Provider, error) {
	switch config.Type {
	case ProviderDeepgram:
		// Import is handled in the specific provider package
		// to avoid circular dependencies
		return nil, fmt.Errorf("provider creation should be done through specific provider package")
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}
}
