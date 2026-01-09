# Examples

This directory contains examples demonstrating various use cases of the TTS Bulk Generator.

## Basic Example

Generate audio files from a simple CSV:

```bash
# Set your API key
export DEEPGRAM_API_KEY=your_key_here

# Run with default settings
./tts-bulk -csv=examples/basic-example.csv
```

## Custom Configuration Example

Use custom model and output settings:

```bash
# Use a different voice model and WAV format
./tts-bulk \
  -csv=examples/custom-config.csv \
  -model=aura-luna-en \
  -format=wav \
  -output=examples/output
```

## Verbose Logging Example

See detailed processing information:

```bash
# Enable verbose logging to see debug information
./tts-bulk -csv=examples/basic-example.csv -verbose
```

## Programmatic Usage

### Using the Library in Your Go Code

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aashish-joshi/tts-bulk/internal/config"
    "github.com/aashish-joshi/tts-bulk/internal/processor"
    "github.com/aashish-joshi/tts-bulk/internal/provider/deepgram"
)

func main() {
    // Create Deepgram provider
    provider, err := deepgram.New(deepgram.Config{
        APIKey:    "your_api_key",
        Model:     "aura-asteria-en",
        Container: "",
        Encoding:  "mp3",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer provider.Close()

    // Create processor
    proc, err := processor.New(processor.Config{
        Provider:      provider,
        OutputDir:     "audio",
        AudioFormat:   "mp3",
        MaxConcurrent: 3,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer proc.Close()

    // Process CSV
    ctx := context.Background()
    results, err := proc.ProcessCSV(ctx, "scripts.csv")
    if err != nil {
        log.Fatal(err)
    }

    // Check results
    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("Failed: %s - %v\n", result.Label, result.Error)
        } else {
            fmt.Printf("Success: %s -> %s\n", result.Label, result.FilePath)
        }
    }
}
```

### Custom Provider Implementation

```go
package main

import (
    "context"
    "fmt"

    "github.com/aashish-joshi/tts-bulk/pkg/types"
)

// CustomProvider implements the types.Provider interface
type CustomProvider struct {
    apiKey string
}

func (p *CustomProvider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
    // Implement your TTS logic here
    fmt.Printf("Generating audio for: %s\n", req.Label)
    return nil
}

func (p *CustomProvider) Name() string {
    return "custom-provider"
}

func (p *CustomProvider) Close() error {
    return nil
}
```

## Error Handling Example

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"

    "github.com/aashish-joshi/tts-bulk/internal/processor"
    "github.com/aashish-joshi/tts-bulk/internal/provider/deepgram"
)

func main() {
    provider, err := deepgram.New(deepgram.Config{
        APIKey:    "your_api_key",
        Model:     "aura-asteria-en",
        Container: "",
        Encoding:  "mp3",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer provider.Close()

    proc, err := processor.New(processor.Config{
        Provider:      provider,
        OutputDir:     "audio",
        AudioFormat:   "mp3",
        MaxConcurrent: 3,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer proc.Close()

    ctx := context.Background()
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

    // Count and report results
    successCount := 0
    failureCount := 0
    
    for _, result := range results {
        if result.Error != nil {
            failureCount++
            log.Printf("Failed to process '%s': %v", result.Label, result.Error)
        } else {
            successCount++
        }
    }

    fmt.Printf("Complete: %d successful, %d failed\n", successCount, failureCount)
}
```

## Files in This Directory

- `basic-example.csv`: Simple CSV with a few entries
- `custom-config.csv`: Example with varied content types
- `README.md`: This file
