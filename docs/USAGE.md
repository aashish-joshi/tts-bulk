# Using TTS Bulk as a Go Module

This guide explains how to use TTS Bulk as a library in your Go applications.

## Installation

```bash
go get github.com/aashish-joshi/tts-bulk
```

## Quick Start

### Basic Usage

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/aashish-joshi/tts-bulk/internal/processor"
	"github.com/aashish-joshi/tts-bulk/internal/provider/deepgram"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

func main() {
	// Create a Deepgram provider
	provider, err := deepgram.New(deepgram.Config{
		APIKey:    os.Getenv("DEEPGRAM_API_KEY"),
		Model:     "aura-asteria-en",
		Container: "",
		Encoding:  "mp3",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Close()

	// Create a processor
	proc, err := processor.New(processor.Config{
		Provider:      provider,
		OutputDir:     "audio",
		AudioFormat:   types.FormatMP3,
		MaxConcurrent: 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer proc.Close()

	// Process a CSV file
	ctx := context.Background()
	results, err := proc.ProcessCSV(ctx, "scripts.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Check results
	for _, result := range results {
		if result.Error != nil {
			log.Printf("Failed: %s - %v\n", result.Label, result.Error)
		} else {
			log.Printf("Success: %s -> %s\n", result.Label, result.FilePath)
		}
	}
}
```

## Component Overview

### 1. Providers (`pkg/types` and `internal/provider`)

Providers implement the TTS generation logic for different services.

#### Available Providers

**Deepgram** (`internal/provider/deepgram`)

```go
import "github.com/aashish-joshi/tts-bulk/internal/provider/deepgram"

provider, err := deepgram.New(deepgram.Config{
	APIKey:    "your-api-key",
	Model:     "aura-asteria-en",
	Container: "",        // "" for MP3, "wav" for WAV
	Encoding:  "mp3",     // "mp3" or "linear16"
})
```

#### Using the Provider Directly

```go
import (
	"context"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

// Generate a single audio file
req := types.TTSRequest{
	Text:   "Hello, world!",
	Label:  "greeting",
	Format: "mp3",
}

err := provider.GenerateAudio(context.Background(), req, "output/greeting.mp3")
```

### 2. Processor (`internal/processor`)

The processor handles batch operations and concurrency.

#### Creating a Processor

```go
import "github.com/aashish-joshi/tts-bulk/internal/processor"

proc, err := processor.New(processor.Config{
	Provider:      provider,        // types.Provider interface
	OutputDir:     "audio",         // Output directory
	AudioFormat:   types.FormatMP3, // Audio format
	MaxConcurrent: 3,               // Max concurrent requests
})
```

#### Processing CSV Files

```go
results, err := proc.ProcessCSV(ctx, "scripts.csv")
```

CSV format: `"label","text"`

```csv
"intro","Welcome to our application"
"step1","Click the settings button"
"step2","Adjust your preferences"
```

#### Processing Results

```go
for _, result := range results {
	if result.Error != nil {
		// Handle error
		fmt.Printf("Failed: %s - %v\n", result.Label, result.Error)
	} else {
		// Success
		fmt.Printf("Generated: %s\n", result.FilePath)
	}
}
```

### 3. Configuration (`internal/config`)

The config package provides configuration management for the CLI. When using as a library, you typically create configurations directly.

```go
import "github.com/aashish-joshi/tts-bulk/pkg/types"

// Audio formats
format := types.FormatMP3  // or types.FormatWAV
```

### 4. Logging (`internal/logger`)

Structured logging for your application.

```go
import "github.com/aashish-joshi/tts-bulk/internal/logger"

// Get default logger
log := logger.Default()

// Set log level
log.SetLevel(logger.LevelDebug)

// Log messages
log.Debug("Debug message")
log.Info("Info message")
log.Warn("Warning message")
log.Error("Error message")
```

## Advanced Usage

### Custom Provider Implementation

Implement the `types.Provider` interface:

```go
package main

import (
	"context"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

type MyProvider struct {
	apiKey string
}

func (p *MyProvider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
	// Your implementation
	return nil
}

func (p *MyProvider) Name() string {
	return "myprovider"
}

func (p *MyProvider) Close() error {
	return nil
}
```

### Error Handling

```go
import (
	"context"
	"errors"
	"fmt"
)

results, err := proc.ProcessCSV(ctx, "scripts.csv")

// Handle processing errors
if err != nil {
	if errors.Is(err, context.Canceled) {
		fmt.Println("Processing was canceled")
	} else {
		fmt.Printf("Processing failed: %v\n", err)
	}
	return
}

// Check individual results
successCount := 0
failureCount := 0

for _, result := range results {
	if result.Error != nil {
		failureCount++
		// Log or handle specific failures
	} else {
		successCount++
	}
}

fmt.Printf("Completed: %d success, %d failed\n", successCount, failureCount)
```

### Context Management

```go
import (
	"context"
	"time"
)

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

results, err := proc.ProcessCSV(ctx, "scripts.csv")

// With cancellation
ctx, cancel := context.WithCancel(context.Background())

// Cancel from another goroutine
go func() {
	// Wait for signal
	<-stopSignal
	cancel()
}()

results, err := proc.ProcessCSV(ctx, "scripts.csv")
```

### Custom Concurrency Control

```go
// Lower concurrency for rate-limited APIs
proc, err := processor.New(processor.Config{
	Provider:      provider,
	OutputDir:     "audio",
	AudioFormat:   types.FormatMP3,
	MaxConcurrent: 2,  // Only 2 concurrent requests
})

// Higher concurrency for faster processing
proc, err := processor.New(processor.Config{
	Provider:      provider,
	OutputDir:     "audio",
	AudioFormat:   types.FormatMP3,
	MaxConcurrent: 10, // 10 concurrent requests
})
```

## Complete Example Application

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aashish-joshi/tts-bulk/internal/logger"
	"github.com/aashish-joshi/tts-bulk/internal/processor"
	"github.com/aashish-joshi/tts-bulk/internal/provider/deepgram"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

func main() {
	// Parse command-line flags
	csvPath := flag.String("csv", "scripts.csv", "Path to CSV file")
	outputDir := flag.String("output", "audio", "Output directory")
	model := flag.String("model", "aura-asteria-en", "TTS model")
	verbose := flag.Bool("verbose", false, "Verbose logging")
	flag.Parse()

	// Setup logging
	appLogger := logger.Default()
	if *verbose {
		appLogger.SetLevel(logger.LevelDebug)
	}

	// Get API key
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPGRAM_API_KEY environment variable not set")
	}

	// Create provider
	provider, err := deepgram.New(deepgram.Config{
		APIKey:    apiKey,
		Model:     *model,
		Container: "",
		Encoding:  "mp3",
	})
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}
	defer provider.Close()

	appLogger.Info("Provider initialized: %s", provider.Name())

	// Create processor
	proc, err := processor.New(processor.Config{
		Provider:      provider,
		OutputDir:     *outputDir,
		AudioFormat:   types.FormatMP3,
		MaxConcurrent: 3,
	})
	if err != nil {
		log.Fatalf("Failed to create processor: %v", err)
	}
	defer proc.Close()

	// Process CSV
	ctx := context.Background()
	results, err := proc.ProcessCSV(ctx, *csvPath)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	// Report results
	successCount := 0
	failureCount := 0

	for _, result := range results {
		if result.Error != nil {
			failureCount++
			appLogger.Error("Failed: %s - %v", result.Label, result.Error)
		} else {
			successCount++
			appLogger.Debug("Generated: %s", result.FilePath)
		}
	}

	appLogger.Info("Complete: %d successful, %d failed", successCount, failureCount)

	if failureCount > 0 {
		os.Exit(1)
	}
}
```

## API Reference

### Types (`pkg/types`)

```go
// Provider interface
type Provider interface {
	GenerateAudio(ctx context.Context, req TTSRequest, outputPath string) error
	Name() string
	Close() error
}

// TTSRequest represents a TTS generation request
type TTSRequest struct {
	Text     string
	Label    string
	Model    string
	Format   string
	Encoding string
}

// TTSResult represents the result of TTS generation
type TTSResult struct {
	Label    string
	FilePath string
	Error    error
}

// AudioFormat constants
const (
	FormatMP3 AudioFormat = "mp3"
	FormatWAV AudioFormat = "wav"
)
```

### Processor (`internal/processor`)

```go
// Create new processor
func New(cfg Config) (*Processor, error)

// Process CSV file
func (p *Processor) ProcessCSV(ctx context.Context, csvPath string) ([]types.TTSResult, error)

// Clean up resources
func (p *Processor) Close() error
```

### Logger (`internal/logger`)

```go
// Create logger
func New(level Level, out io.Writer) *Logger

// Get default logger
func Default() *Logger

// Set log level
func (l *Logger) SetLevel(level Level)

// Log methods
func (l *Logger) Debug(format string, v ...interface{})
func (l *Logger) Info(format string, v ...interface{})
func (l *Logger) Warn(format string, v ...interface{})
func (l *Logger) Error(format string, v ...interface{})
```

## Testing

### Unit Testing with Mocks

```go
import (
	"context"
	"testing"
	"github.com/aashish-joshi/tts-bulk/pkg/types"
)

type MockProvider struct{}

func (m *MockProvider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
	// Mock implementation
	return nil
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) Close() error {
	return nil
}

func TestYourFunction(t *testing.T) {
	provider := &MockProvider{}
	// Use provider in your tests
}
```

## Best Practices

1. **Always close resources**:
   ```go
   defer provider.Close()
   defer proc.Close()
   ```

2. **Use contexts for cancellation**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
   defer cancel()
   ```

3. **Handle errors properly**:
   ```go
   if err != nil {
       return fmt.Errorf("failed to process: %w", err)
   }
   ```

4. **Use appropriate log levels**:
   - Debug: Detailed diagnostic information
   - Info: General informational messages
   - Warn: Warning conditions
   - Error: Error conditions

## Getting Help

- **Documentation**: See `docs/` directory
- **Examples**: Check `examples/` directory
- **Issues**: Open an issue on GitHub
- **Reference**: See `cmd/tts-bulk/main.go` for CLI implementation

## License

Apache License 2.0 - see LICENSE file for details.
