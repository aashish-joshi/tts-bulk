package provider

import (
	"fmt"
	"strings"

	"github.com/aashish-joshi/tts-bulk/internal/config"
	"github.com/aashish-joshi/tts-bulk/internal/provider/deepgram"
	"github.com/aashish-joshi/tts-bulk/internal/provider/openaispeech"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

// New creates and returns the TTS provider specified in cfg.
func New(cfg *config.Config) (types.Provider, error) {
	switch strings.ToLower(cfg.ProviderType) {
	case string(ProviderDeepgram):
		return deepgram.New(deepgram.Config{
			APIKey:      cfg.APIKey,
			Model:       cfg.Model,
			Format:      string(cfg.AudioFormat),
			RateLimitMs: cfg.RateLimitMs,
		})
	case string(ProviderLocal), string(ProviderOpenAI):
		return openaispeech.New(openaispeech.Config{
			BaseURL:     cfg.ProviderURL,
			APIKey:      cfg.APIKey,
			Model:       cfg.Model,
			Voice:       cfg.Voice,
			Format:      string(cfg.AudioFormat),
			RateLimitMs: cfg.RateLimitMs,
		})
	default:
		return nil, fmt.Errorf("unknown provider %q (supported: deepgram, local, openai)", cfg.ProviderType)
	}
}
