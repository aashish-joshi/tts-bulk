// Package types defines common types used across the application.
package types

import "context"

// TTSRequest represents a text-to-speech conversion request.
type TTSRequest struct {
	Text     string
	Label    string
	Model    string
	Format   string
	Encoding string
}

// TTSResult represents the result of a TTS conversion.
type TTSResult struct {
	Label    string
	FilePath string
	Error    error
}

// AudioFormat represents supported audio formats.
type AudioFormat string

const (
	// FormatMP3 represents MP3 audio format.
	FormatMP3 AudioFormat = "mp3"
	// FormatWAV represents WAV audio format.
	FormatWAV AudioFormat = "wav"
)

// Provider represents a TTS provider interface.
type Provider interface {
	// GenerateAudio generates audio from text and saves it to the specified path.
	GenerateAudio(ctx context.Context, req TTSRequest, outputPath string) error
	// Name returns the provider's name.
	Name() string
	// Close cleans up any resources used by the provider.
	Close() error
}
