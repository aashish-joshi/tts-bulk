// Package processor handles batch TTS processing.
package processor

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aashish-joshi/tts-bulk/internal/logger"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

// Processor handles batch TTS processing.
type Processor struct {
	provider      types.Provider
	outputDir     string
	audioFormat   types.AudioFormat
	maxConcurrent int
	logger        *logger.Logger
}

// Config holds processor configuration.
type Config struct {
	Provider      types.Provider
	OutputDir     string
	AudioFormat   types.AudioFormat
	MaxConcurrent int
}

// New creates a new batch processor.
func New(cfg Config) (*Processor, error) {
	if cfg.Provider == nil {
		return nil, fmt.Errorf("provider is required")
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = "audio"
	}

	if cfg.MaxConcurrent <= 0 {
		cfg.MaxConcurrent = 3
	}

	// Create output directory if it doesn't exist
	outputDir := strings.ToLower(cfg.OutputDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &Processor{
		provider:      cfg.Provider,
		outputDir:     outputDir,
		audioFormat:   cfg.AudioFormat,
		maxConcurrent: cfg.MaxConcurrent,
		logger:        logger.Default(),
	}, nil
}

// ProcessCSV processes a CSV file and generates TTS audio for each row.
func (p *Processor) ProcessCSV(ctx context.Context, csvPath string) ([]types.TTSResult, error) {
	// Read CSV file
	records, err := p.readCSV(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	p.logger.Info("Processing %d records from %s", len(records), csvPath)

	// Process records concurrently
	results := p.processRecords(ctx, records)

	// Count successes and failures
	successCount := 0
	failureCount := 0
	for _, result := range results {
		if result.Error == nil {
			successCount++
		} else {
			failureCount++
		}
	}

	p.logger.Info("Processing complete: %d successful, %d failed", successCount, failureCount)

	return results, nil
}

func (p *Processor) readCSV(csvPath string) ([][2]string, error) {
	// Check if file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", csvPath)
	}

	// Open CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read CSV
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Convert to records
	records := make([][2]string, 0, len(rows))
	for i, row := range rows {
		if len(row) < 2 {
			p.logger.Warn("Skipping invalid record at row %d: insufficient columns", i+1)
			continue
		}
		records = append(records, [2]string{row[0], row[1]})
	}

	return records, nil
}

func (p *Processor) processRecords(ctx context.Context, records [][2]string) []types.TTSResult {
	results := make([]types.TTSResult, len(records))
	var wg sync.WaitGroup

	// Semaphore for concurrency control
	semaphore := make(chan struct{}, p.maxConcurrent)

	for i, record := range records {
		wg.Add(1)

		go func(idx int, label, text string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Process the record
			result := p.processRecord(ctx, label, text)
			results[idx] = result
		}(i, record[0], record[1])
	}

	wg.Wait()
	return results
}

func (p *Processor) processRecord(ctx context.Context, label, text string) types.TTSResult {
	result := types.TTSResult{
		Label: label,
	}

	// Build output path
	outputPath := filepath.Join(p.outputDir, fmt.Sprintf("%s.%s", label, p.audioFormat))
	result.FilePath = outputPath

	// Create TTS request
	req := types.TTSRequest{
		Text:   text,
		Label:  label,
		Format: string(p.audioFormat),
	}

	// Generate audio
	err := p.provider.GenerateAudio(ctx, req, outputPath)
	if err != nil {
		result.Error = err
		p.logger.Error("Failed to generate TTS for '%s': %v", label, err)
	} else {
		p.logger.Info("Generated TTS for '%s' -> %s", label, outputPath)
	}

	return result
}

// Close cleans up processor resources.
func (p *Processor) Close() error {
	if p.provider != nil {
		return p.provider.Close()
	}
	return nil
}
