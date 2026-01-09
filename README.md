# TTS Bulk Generator

[![CI](https://github.com/aashish-joshi/tts-bulk/workflows/CI/badge.svg)](https://github.com/aashish-joshi/tts-bulk/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/aashish-joshi/tts-bulk)](https://goreportcard.com/report/github.com/aashish-joshi/tts-bulk)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A modern, high-performance bulk text-to-speech (TTS) generator written in Go. Convert large batches of text to audio files efficiently using various TTS providers.

## Features

- 🚀 **High Performance**: Concurrent processing with configurable parallelism
- 🔌 **Modular Architecture**: Easy to add new TTS providers
- 📊 **Progress Tracking**: Real-time status updates and detailed logging
- 🎯 **Multiple Formats**: Support for MP3 and WAV audio formats
- 🛡️ **Robust Error Handling**: Comprehensive error handling and reporting
- 📝 **Structured Logging**: Clear, informative logs for debugging and monitoring
- ✅ **Well Tested**: Comprehensive test suite with high coverage
- 🏗️ **Modern Go Layout**: Follows Go project layout best practices

## Supported Providers

- **Deepgram**: High-quality neural TTS with multiple models and voices

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/aashish-joshi/tts-bulk/releases).

### Build from Source

Requires Go 1.22 or later.

```bash
git clone https://github.com/aashish-joshi/tts-bulk.git
cd tts-bulk
make build
```

Or using Go directly:

```bash
go install github.com/aashish-joshi/tts-bulk/cmd/tts-bulk@latest
```

## Quick Start

1. **Get an API key** from [Deepgram](https://www.deepgram.com/) and set it as an environment variable:
   ```bash
   export DEEPGRAM_API_KEY=your_api_key_here
   ```

2. **Create a CSV file** with your scripts (see `sample-scripts.csv` for an example):
   ```csv
   "hello-world","Hello World, how are you doing?"
   "greeting","Welcome to TTS Bulk Generator!"
   ```

3. **Run the tool**:
   ```bash
   ./tts-bulk
   ```

Audio files will be generated in the `audio/` directory by default.

## Usage

### Basic Usage

```bash
# Generate MP3 files from scripts.csv
./tts-bulk

# Generate WAV files
./tts-bulk -format=wav

# Use a specific model
./tts-bulk -model=aura-asteria-en

# Custom output directory and CSV file
./tts-bulk -csv=my-scripts.csv -output=my-audio
```

### Command-Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-csv` | Path to CSV file containing scripts | `scripts.csv` |
| `-output` | Output directory for audio files | `audio` |
| `-format` | Audio format (mp3, wav) | `mp3` |
| `-model` | TTS model name | `aura-asteria-en` |
| `-verbose` | Enable verbose logging | `false` |
| `-version` | Show version information | - |

### CSV File Format

The CSV file should have two columns:
1. **Label**: Used as the filename for the generated audio
2. **Script**: The text to convert to speech

Example:
```csv
"intro","Welcome to our application"
"tutorial-step1","First, click on the settings button"
"tutorial-step2","Then, adjust your preferences"
```

### Available Models

For Deepgram, see the [available models documentation](https://developers.deepgram.com/docs/tts-models). Some popular options:

- `aura-asteria-en` (default)
- `aura-luna-en`
- `aura-stella-en`
- `aura-athena-en`
- `aura-hera-en`

## Architecture

### Project Structure

```
tts-bulk/
├── cmd/tts-bulk/          # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── logger/            # Structured logging
│   ├── processor/         # Batch processing logic
│   └── provider/          # TTS provider implementations
│       └── deepgram/      # Deepgram provider
├── pkg/types/             # Public types and interfaces
└── Makefile               # Build automation
```

### Key Components

- **Config**: Manages application configuration from environment and flags
- **Logger**: Provides structured logging with configurable levels
- **Processor**: Handles batch processing with concurrency control
- **Provider**: Abstraction layer for TTS services (currently Deepgram)

### Design Principles

- **Separation of Concerns**: Each package has a single, well-defined responsibility
- **Dependency Injection**: Components receive dependencies through constructors
- **Interface-based Design**: Provider abstraction allows easy addition of new TTS services
- **Error Propagation**: Errors are wrapped with context for better debugging
- **Testability**: Mock-friendly design with comprehensive tests

## Adding a New Provider

The modular architecture makes it easy to add support for new TTS providers:

1. Create a new package under `internal/provider/`
2. Implement the `types.Provider` interface:
   ```go
   type Provider interface {
       GenerateAudio(ctx context.Context, req TTSRequest, outputPath string) error
       Name() string
       Close() error
   }
   ```
3. Add configuration support
4. Update the factory pattern
5. Add tests

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed instructions.

## Development

### Prerequisites

- Go 1.22 or later
- Make (optional, for using Makefile commands)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linters (requires golangci-lint)
make lint

# Run all checks
make check
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

## Performance

- **Concurrent Processing**: Processes multiple TTS requests in parallel (default: 3 concurrent requests)
- **Rate Limiting**: Built-in rate limiting to respect API limits
- **Efficient I/O**: Streaming writes for large audio files
- **Resource Management**: Proper cleanup and resource pooling

## Troubleshooting

### API Key Issues

If you see `DEEPGRAM_API_KEY environment variable is not set`:
```bash
export DEEPGRAM_API_KEY=your_key_here
```

### File Not Found

Ensure your CSV file exists:
```bash
ls -la scripts.csv
```

### Permission Errors

Make sure the output directory is writable:
```bash
chmod 755 audio/
```

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Changelog

### Version 2.0.0 (2024)

- Complete rewrite with modern Go project structure
- Modular architecture with provider abstraction
- Improved error handling and logging
- Comprehensive test suite
- CI/CD pipeline with GitHub Actions
- Better concurrency control and performance
- Enhanced documentation

### Version 1.0.0 (Initial)

- Basic TTS generation with Deepgram
- Simple CSV processing
- MP3 and WAV format support

## Acknowledgments

- [Deepgram](https://www.deepgram.com/) for providing the TTS API
- The Go community for excellent tools and libraries

## Support

- 📫 [Report Issues](https://github.com/aashish-joshi/tts-bulk/issues)
- 💬 [Discussions](https://github.com/aashish-joshi/tts-bulk/discussions)
- 📖 [Documentation](https://github.com/aashish-joshi/tts-bulk/wiki)

---

Made with ❤️ by the TTS Bulk community
