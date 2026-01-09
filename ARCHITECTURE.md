# Architecture Documentation

## Overview

TTS Bulk is designed with a modern, modular Go architecture that emphasizes:
- **Separation of concerns**: Each component has a single responsibility
- **Interface-based design**: Easy to extend and test
- **Dependency injection**: Components receive dependencies explicitly
- **Error handling**: Comprehensive error wrapping and context

## Project Structure

```
tts-bulk/
├── cmd/                        # Application entry points
│   └── tts-bulk/              # Main CLI application
│       └── main.go            # Application bootstrap
├── internal/                   # Private application code
│   ├── config/                # Configuration management
│   │   ├── config.go          # Config loading and validation
│   │   └── config_test.go     # Config tests
│   ├── logger/                # Structured logging
│   │   ├── logger.go          # Logger implementation
│   │   └── logger_test.go     # Logger tests
│   ├── processor/             # Batch processing
│   │   ├── processor.go       # Main processor logic
│   │   └── processor_test.go  # Processor tests
│   └── provider/              # TTS provider abstraction
│       ├── provider.go        # Provider interface
│       └── deepgram/          # Deepgram implementation
│           └── deepgram.go    # Deepgram provider
├── pkg/                       # Public library code
│   └── types/                 # Public types
│       └── types.go           # Common types and interfaces
├── examples/                  # Usage examples
│   ├── README.md              # Examples documentation
│   ├── basic-example.csv      # Simple example
│   └── custom-config.csv      # Advanced example
├── .github/                   # GitHub configuration
│   └── workflows/             # CI/CD workflows
│       └── ci.yml             # Main CI pipeline
├── Makefile                   # Build automation
├── .golangci.yml             # Linter configuration
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── README.md                  # Project documentation
├── CONTRIBUTING.md            # Contribution guidelines
└── ARCHITECTURE.md            # This file
```

## Core Components

### 1. Types Package (`pkg/types`)

**Purpose**: Defines public interfaces and types used throughout the application.

**Key Types**:
- `Provider`: Interface that all TTS providers must implement
- `TTSRequest`: Request structure for TTS generation
- `TTSResult`: Result structure containing output and errors
- `AudioFormat`: Enum for supported audio formats

**Design Principles**:
- Public API that can be imported by external packages
- No dependencies on internal packages
- Well-documented interfaces for extensibility

### 2. Configuration (`internal/config`)

**Purpose**: Centralized configuration management.

**Responsibilities**:
- Load configuration from environment variables
- Parse and validate command-line flags
- Provide typed configuration to other components
- Set sensible defaults

**Design Decisions**:
- Single source of truth for configuration
- Validation happens at load time
- Immutable after creation

### 3. Logger (`internal/logger`)

**Purpose**: Provides structured logging throughout the application.

**Features**:
- Multiple log levels (Debug, Info, Warn, Error)
- Thread-safe operations
- Configurable output destination
- Consistent format across application

**Design Decisions**:
- Singleton pattern with default instance
- Level-based filtering
- No external dependencies

### 4. Processor (`internal/processor`)

**Purpose**: Orchestrates batch TTS processing.

**Responsibilities**:
- Read and parse CSV files
- Manage concurrent processing
- Track progress and results
- Handle errors gracefully

**Key Features**:
- Configurable concurrency with semaphore pattern
- Goroutine-based parallel processing
- Comprehensive error collection
- Progress reporting

**Design Decisions**:
- Semaphore for concurrency control
- WaitGroup for synchronization
- Non-blocking error collection

### 5. Provider Abstraction (`internal/provider`)

**Purpose**: Abstract TTS provider implementations behind a common interface.

**Benefits**:
- Easy to add new providers
- Testable with mock providers
- Provider-agnostic business logic

**Implementation**: Deepgram Provider (`internal/provider/deepgram`)

**Features**:
- Implements the `Provider` interface
- Handles Deepgram-specific configuration
- Rate limiting to respect API limits
- Error wrapping with context

## Data Flow

```
1. User runs CLI with flags
   ↓
2. Configuration loaded and validated
   ↓
3. Provider initialized (e.g., Deepgram)
   ↓
4. Processor created with provider
   ↓
5. CSV file read and parsed
   ↓
6. Records processed concurrently
   ├→ Generate TTS for record 1
   ├→ Generate TTS for record 2
   └→ Generate TTS for record 3
   ↓
7. Results collected and reported
   ↓
8. Resources cleaned up
```

## Concurrency Model

### Processing Pipeline

The processor uses a semaphore-based concurrency model:

```go
// Semaphore limits concurrent goroutines
semaphore := make(chan struct{}, maxConcurrent)

for each record {
    wg.Add(1)
    go func(record) {
        // Acquire semaphore slot
        semaphore <- struct{}{}
        defer func() { <-semaphore }()
        
        // Process record
        generateTTS(record)
        
        wg.Done()
    }(record)
}

wg.Wait() // Wait for all goroutines
```

**Benefits**:
- Controlled parallelism
- Efficient resource utilization
- Protection against overwhelming the API
- Simple error handling

### Rate Limiting

Provider implementations include rate limiting:

```go
// In Deepgram provider
func (p *Provider) GenerateAudio(...) error {
    // Generate audio
    _, err := p.client.ToSave(...)
    
    // Rate limit: 1 second between requests
    time.Sleep(time.Second)
    
    return err
}
```

## Error Handling

### Error Wrapping

Errors are wrapped with context at each layer:

```go
// In processor
if err := provider.GenerateAudio(...); err != nil {
    return fmt.Errorf("failed to generate TTS: %w", err)
}

// In provider
if cfg.APIKey == "" {
    return nil, fmt.Errorf("API key is required")
}
```

### Error Collection

The processor collects errors without stopping:

```go
for _, result := range results {
    if result.Error != nil {
        // Log error but continue
        logger.Error("Failed: %s - %v", result.Label, result.Error)
    }
}
```

## Testing Strategy

### Unit Tests

Each component has comprehensive unit tests:

- **Config**: Test configuration loading and validation
- **Logger**: Test log levels and output
- **Processor**: Test CSV parsing, processing, error handling
- **Provider**: Mock provider for testing without API calls

### Test Structure

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

### Mocking

The `Provider` interface enables easy mocking:

```go
type MockProvider struct {
    failOnGenerate bool
}

func (m *MockProvider) GenerateAudio(...) error {
    if m.failOnGenerate {
        return errors.New("mock error")
    }
    return nil
}
```

## Extensibility

### Adding a New Provider

1. **Create provider package**:
   ```
   internal/provider/newprovider/
   └── newprovider.go
   ```

2. **Implement Provider interface**:
   ```go
   type Provider struct {
       // Provider-specific fields
   }
   
   func (p *Provider) GenerateAudio(...) error {
       // Implementation
   }
   
   func (p *Provider) Name() string {
       return "newprovider"
   }
   
   func (p *Provider) Close() error {
       // Cleanup
   }
   ```

3. **Add configuration support**:
   Update `internal/config/config.go` to support new provider

4. **Update factory pattern**:
   Add provider creation logic in main.go or factory

5. **Add tests**:
   Create `newprovider_test.go` with comprehensive tests

### Adding New Features

The modular design makes it easy to add features:

- **Progress tracking**: Add to processor without changing providers
- **Retry logic**: Implement in processor layer
- **Custom output formats**: Add to config and provider interface
- **Metrics**: Add new logger or metrics package

## Performance Considerations

### Optimization Strategies

1. **Concurrent Processing**:
   - Multiple goroutines process records in parallel
   - Configurable concurrency level

2. **Efficient I/O**:
   - Direct file writes from provider
   - No intermediate buffering

3. **Rate Limiting**:
   - Prevents API throttling
   - Maintains good API citizenship

4. **Resource Management**:
   - Proper cleanup with defer
   - Semaphore-based goroutine control

### Scalability

Current architecture supports:
- **Hundreds of concurrent requests**: Adjust `maxConcurrent`
- **Large CSV files**: Streaming CSV parsing
- **Long-running jobs**: Graceful shutdown support (future)

## Security

### Best Practices

1. **API Key Management**:
   - Never hardcode API keys
   - Use environment variables
   - No API keys in logs

2. **Input Validation**:
   - Validate all configuration inputs
   - Sanitize file paths
   - Check CSV format

3. **Error Messages**:
   - Don't leak sensitive information
   - Log safely without exposing credentials

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
jobs:
  - test: Run all unit tests with race detector
  - lint: Run golangci-lint for code quality
  - build: Build for multiple platforms (Linux, macOS, Windows)
```

### Quality Gates

- All tests must pass
- Linter must report no issues
- Build must succeed for all platforms

## Future Improvements

### Potential Enhancements

1. **Additional Providers**:
   - AWS Polly
   - Google Cloud TTS
   - Azure Cognitive Services

2. **Advanced Features**:
   - Retry logic with exponential backoff
   - Progress bar for CLI
   - Configuration files (YAML/JSON)
   - Dry-run mode
   - Webhooks for completion notifications

3. **Performance**:
   - Adaptive rate limiting
   - Connection pooling
   - Caching for repeated text

4. **Observability**:
   - Metrics export (Prometheus)
   - Distributed tracing
   - Health checks

## Conclusion

The architecture balances simplicity with extensibility. It follows Go best practices and makes it easy to:
- Add new TTS providers
- Test components in isolation
- Scale processing capacity
- Maintain and evolve the codebase

For questions or suggestions about the architecture, please open an issue or discussion on GitHub.
