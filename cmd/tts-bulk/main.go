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
	"github.com/aashish-joshi/tts-bulk/internal/provider/deepgram"
)

const (
	appName    = "tts-bulk"
	appVersion = "2.0.0"
)

func main() {
	// Define command-line flags
	model := flag.String("model", "aura-asteria-en", "TTS model name (default: aura-asteria-en)")
	format := flag.String("format", "mp3", "Audio format: mp3 or wav (default: mp3)")
	output := flag.String("output", "audio", "Output directory for audio files (default: audio)")
	csvPath := flag.String("csv", "scripts.csv", "Path to CSV file (default: scripts.csv)")
	version := flag.Bool("version", false, "Show version information")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	// Show version if requested
	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		return
	}

	// Set up logger
	log := logger.Default()
	if *verbose {
		log.SetLevel(logger.LevelDebug)
		log.Debug("Verbose logging enabled")
	}

	log.Info("Starting %s v%s", appName, appVersion)

	// Load configuration
	cfg, err := config.Load(*csvPath, *output, *format, *model)
	if err != nil {
		log.Error("Configuration error: %v", err)
		os.Exit(1)
	}

	log.Info("Configuration loaded successfully")
	log.Debug("Provider: %s, Model: %s, Format: %s", cfg.ProviderType, cfg.Model, cfg.AudioFormat)
	log.Debug("CSV: %s, Output: %s", cfg.CSVPath, cfg.OutputDir)

	// Create provider
	container, encoding := cfg.GetDeepgramConfig()
	dgProvider, err := deepgram.New(deepgram.Config{
		APIKey:    cfg.APIKey,
		Model:     cfg.Model,
		Container: container,
		Encoding:  encoding,
	})
	if err != nil {
		log.Error("Failed to create provider: %v", err)
		os.Exit(1)
	}
	defer dgProvider.Close()

	log.Info("Provider initialized: %s", dgProvider.Name())

	// Create processor
	proc, err := processor.New(processor.Config{
		Provider:      dgProvider,
		OutputDir:     cfg.OutputDir,
		AudioFormat:   cfg.AudioFormat,
		MaxConcurrent: cfg.MaxConcurrent,
	})
	if err != nil {
		log.Error("Failed to create processor: %v", err)
		os.Exit(1)
	}
	defer proc.Close()

	log.Info("Processor initialized with max concurrency: %d", cfg.MaxConcurrent)

	// Process CSV
	ctx := context.Background()
	results, err := proc.ProcessCSV(ctx, cfg.CSVPath)
	if err != nil {
		log.Error("Processing failed: %v", err)
		os.Exit(1)
	}

	// Report results
	successCount := 0
	failureCount := 0
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
		os.Exit(1)
	}
}
