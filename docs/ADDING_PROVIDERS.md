# Adding a New TTS Provider

This guide walks you through adding a new TTS provider to the TTS Bulk project.

## Overview

The TTS Bulk architecture uses an interface-based design that makes it easy to add new TTS providers. Each provider must implement the `types.Provider` interface defined in `pkg/types/types.go`.

## Prerequisites

- Understanding of Go interfaces
- API credentials for the TTS service you want to integrate
- Familiarity with the TTS service's SDK or API

## Step-by-Step Guide

### Step 1: Create Provider Package

Create a new directory under `internal/provider/` for your provider:

```bash
mkdir -p internal/provider/yourprovider
```

### Step 2: Implement the Provider Interface

Create `internal/provider/yourprovider/yourprovider.go`:

```go
package yourprovider

import (
 "context"
 "fmt"
 "os"
 
 "github.com/aashish-joshi/tts-bulk/pkg/types"
)

// Provider implements the TTS provider interface.
type Provider struct {
 client  *YourClient
 apiKey  string
 model   string
}

// Config holds provider-specific configuration.
type Config struct {
 APIKey string
 Model  string
}

// New creates a new provider instance.
func New(cfg Config) (*Provider, error) {
 if cfg.APIKey == "" {
  return nil, fmt.Errorf("API key is required")
 }
 
 // Initialize your client
 client := initializeClient(cfg.APIKey)
 
 return &Provider{
  client: client,
  apiKey: cfg.APIKey,
  model:  cfg.Model,
 }, nil
}

// GenerateAudio generates audio from text and saves it to outputPath.
func (p *Provider) GenerateAudio(ctx context.Context, req types.TTSRequest, outputPath string) error {
 // Validate input
 if req.Text == "" {
  return fmt.Errorf("text cannot be empty")
 }
 
 // Call your provider's API
 audioData, err := p.client.Synthesize(ctx, req.Text, p.model)
 if err != nil {
  return fmt.Errorf("failed to synthesize: %w", err)
 }
 
 // Save to file
 if err := os.WriteFile(outputPath, audioData, 0644); err != nil {
  return fmt.Errorf("failed to write file: %w", err)
 }
 
 return nil
}

// Name returns the provider's name.
func (p *Provider) Name() string {
 return "yourprovider"
}

// Close cleans up resources.
func (p *Provider) Close() error {
 return nil
}
```

**Reference**: See `internal/provider/deepgram/deepgram.go` for a complete implementation example.

### Step 3: Add Tests

Create `internal/provider/yourprovider/yourprovider_test.go`:

```go
package yourprovider

import (
 "context"
 "testing"
)

func TestNew(t *testing.T) {
 // Test valid configuration
 provider, err := New(Config{
  APIKey: "test-key",
  Model:  "test-model",
 })
 if err != nil {
  t.Fatalf("New() failed: %v", err)
 }
 if provider == nil {
  t.Fatal("New() returned nil provider")
 }
}
```

**Reference**: See `internal/processor/processor_test.go` for mock provider examples.

### Step 4: Register in the Provider Factory

Add your provider to the switch statement in `internal/provider/factory.go`:

```go
case "yourprovider":
    return yourprovider.New(yourprovider.Config{
        APIKey: cfg.APIKey,
        Model:  cfg.Model,
    })
```

Also add a constant in `internal/provider/provider.go`:

```go
// ProviderYours represents your TTS provider.
ProviderYours ProviderType = "yourprovider"
```

That is all that needs to change — the CLI flag (`-provider=yourprovider`) is already
wired through the factory, so no edits to `main.go` are required.

### Step 5: Update Documentation

1. **Add to README.md**:

   ```markdown
   ## Supported Providers
   - **Deepgram**: High-quality neural TTS
   - **YourProvider**: Description of your provider
   ```

2. **Update examples/** if needed

3. **Add environment variable to documentation**:

   ```markdown
   export YOUR_PROVIDER_API_KEY=your_key_here
   ```

### Step 6: Test Your Implementation

```bash
# Run tests
go test ./internal/provider/yourprovider/

# Test with the CLI
export YOUR_PROVIDER_API_KEY=your_key
./tts-bulk -provider=yourprovider -csv=sample-scripts.csv
```

## The Provider Interface

All providers must implement this interface from `pkg/types/types.go`:

```go
type Provider interface {
 // GenerateAudio generates audio from text and saves to outputPath
 GenerateAudio(ctx context.Context, req TTSRequest, outputPath string) error
 
 // Name returns the provider's name
 Name() string
 
 // Close cleans up resources
 Close() error
}
```

## Best Practices

### 1. Error Handling

Always wrap errors with context:

```go
return fmt.Errorf("failed to generate audio: %w", err)
```

### 2. Rate Limiting

Respect API rate limits:

```go
import "time"

func (p *Provider) GenerateAudio(...) error {
 // Generate audio
 err := p.callAPI()
 
 // Rate limit
 time.Sleep(time.Second)
 
 return err
}
```

### 3. Resource Cleanup

Always clean up in Close():

```go
func (p *Provider) Close() error {
 if p.client != nil {
  return p.client.Close()
 }
 return nil
}
```

### 4. Input Validation

Validate all inputs:

```go
if req.Text == "" {
 return fmt.Errorf("text cannot be empty")
}
if outputPath == "" {
 return fmt.Errorf("output path cannot be empty")
}
```

### 5. Logging

Use the internal logger:

```go
import "github.com/aashish-joshi/tts-bulk/internal/logger"

logger.Debug("Generating audio with model: %s", p.model)
logger.Info("Successfully generated audio for: %s", req.Label)
```

## Reference Implementation

The Deepgram provider (`internal/provider/deepgram/`) is a complete, production-ready reference implementation that demonstrates all best practices.

## Getting Help

- Check `internal/provider/deepgram/` for reference
- Review `pkg/types/types.go` for interface details
- See `internal/processor/processor_test.go` for mocking examples
- Open an issue on GitHub for questions

## Contributing

Once your provider is ready:

1. Ensure all tests pass
2. Update documentation
3. Add examples
4. Submit a pull request

Thank you for contributing! 🎉
