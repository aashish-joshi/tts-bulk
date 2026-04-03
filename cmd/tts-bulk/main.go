// Package main is the entry point for the TTS Bulk application.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aashish-joshi/tts-bulk/internal/config"
	"github.com/aashish-joshi/tts-bulk/internal/logger"
	"github.com/aashish-joshi/tts-bulk/internal/processor"
	"github.com/aashish-joshi/tts-bulk/internal/provider"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

const (
	appName    = "tts-bulk"
	appVersion = "2.0.0"
)

func main() {
	model := flag.String("model", "aura-asteria-en", "TTS model name (default: aura-asteria-en)")
	format := flag.String("format", "mp3", "Audio format: mp3 or wav (default: mp3)")
	output := flag.String("output", "audio", "Output directory for audio files (default: audio)")
	csvPath := flag.String("csv", "scripts.csv", "Path to CSV file (default: scripts.csv)")
	providerName := flag.String("provider", "", "TTS provider: deepgram or local (default: deepgram, or TTS_PROVIDER env var)")
	providerURL := flag.String("provider-url", "http://localhost:8000", "Base URL for local/OpenAI-compatible provider (default: http://localhost:8000)")
	voice := flag.String("voice", "", "Voice for TTS, used by local/OpenAI-compatible provider (e.g. alloy)")
	rateLimitMs := flag.Int("rate-limit-ms", -1, "Milliseconds to wait between requests (-1 = provider default)")
	version := flag.Bool("version", false, "Show version information")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		return
	}

	log := logger.Default()
	if *verbose {
		log.SetLevel(logger.LevelDebug)
		log.Debug("Verbose logging enabled")
	}

	cfg, err := config.Load(*csvPath, *output, *format, *model, *providerName, *providerURL, *voice, *rateLimitMs)
	if err != nil {
		log.Error("Configuration error: %v", err)
		os.Exit(1)
	}

	os.Exit(run(log, cfg))
}

func run(log *logger.Logger, cfg *config.Config) int {
	log.Info("Starting %s v%s", appName, appVersion)
	log.Debug("Provider: %s, Model: %s, Format: %s", cfg.ProviderType, cfg.Model, cfg.AudioFormat)
	log.Debug("CSV: %s, Output: %s", cfg.CSVPath, cfg.OutputDir)

	p, err := provider.New(cfg)
	if err != nil {
		log.Error("Failed to create provider: %v", err)
		return 1
	}
	defer p.Close()

	log.Info("Provider initialized: %s", p.Name())

	proc, err := processor.New(processor.Config{
		Provider:      p,
		OutputDir:     cfg.OutputDir,
		AudioFormat:   cfg.AudioFormat,
		MaxConcurrent: cfg.MaxConcurrent,
	})
	if err != nil {
		log.Error("Failed to create processor: %v", err)
		return 1
	}
	defer proc.Close()

	log.Info("Processor initialized with max concurrency: %d", cfg.MaxConcurrent)

	ctx := context.Background()
	results, err := proc.ProcessCSV(ctx, cfg.CSVPath)
	if err != nil {
		log.Error("Processing failed: %v", err)
		return 1
	}

	return reportResults(log, results)
}

func reportResults(log *logger.Logger, results []types.TTSResult) int {
	successCount, failureCount := 0, 0
	for _, result := range results {
		if result.Error != nil {
			failureCount++
			log.Error("Failed: %s - %v", result.Label, result.Error)
		} else {
			successCount++
		}
	}

	log.Info("Processing complete: %d successful, %d failed", successCount, failureCount)

	if failureCount > 0 {
		return 1
	}
	return 0
}
