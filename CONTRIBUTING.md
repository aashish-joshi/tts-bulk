# Contributing to TTS Bulk

Thank you for your interest in contributing to TTS Bulk! This document provides guidelines and instructions for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- A text editor or IDE
- (Optional) golangci-lint for linting

### Setting Up the Development Environment

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR-USERNAME/tts-bulk.git
   cd tts-bulk
   ```

2. Install dependencies:
   ```bash
   make install
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

## Project Structure

```
tts-bulk/
├── cmd/
│   └── tts-bulk/          # Main application entry point
│       └── main.go
├── internal/              # Private application code
│   ├── config/            # Configuration management
│   ├── logger/            # Structured logging
│   ├── processor/         # Batch processing logic
│   └── provider/          # TTS provider implementations
│       └── deepgram/      # Deepgram provider
├── pkg/
│   └── types/             # Public types and interfaces
├── .github/
│   └── workflows/         # CI/CD pipelines
├── Makefile               # Build automation
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

## Development Workflow

### Making Changes

1. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding standards below.

3. Add tests for your changes:
   ```bash
   # Run tests to ensure they pass
   make test
   ```

4. Format and lint your code:
   ```bash
   make fmt
   make vet
   make lint  # Requires golangci-lint
   ```

5. Commit your changes with a clear commit message:
   ```bash
   git commit -m "Add feature: description of your changes"
   ```

6. Push to your fork and submit a pull request.

## Coding Standards

### Go Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code (run `make fmt`)
- Use meaningful variable and function names
- Add comments for exported functions, types, and packages
- Keep functions small and focused on a single responsibility

### Code Organization

- **cmd/**: Contains main applications
- **internal/**: Contains private application code that shouldn't be imported by other projects
- **pkg/**: Contains public libraries that can be imported by other projects
- Place tests in `*_test.go` files alongside the code they test

### Testing

- Write unit tests for all new functionality
- Aim for high test coverage (>80%)
- Use table-driven tests where appropriate
- Mock external dependencies in tests
- Run tests before submitting a PR: `make test`

### Error Handling

- Always check and handle errors
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Use custom error types for domain-specific errors when appropriate

### Logging

- Use the `internal/logger` package for logging
- Use appropriate log levels:
  - `Debug`: Detailed information for debugging
  - `Info`: General informational messages
  - `Warn`: Warning messages for potentially problematic situations
  - `Error`: Error messages for failures

## Adding a New TTS Provider

To add a new TTS provider, follow the comprehensive guide in [docs/ADDING_PROVIDERS.md](docs/ADDING_PROVIDERS.md).

**Quick steps:**

1. Create a new package under `internal/provider/yourprovider`
2. Implement the `types.Provider` interface
3. Add tests for your provider implementation
4. Update `cmd/tts-bulk/main.go` to support your provider
5. Update documentation

**Detailed instructions with examples:** See [docs/ADDING_PROVIDERS.md](docs/ADDING_PROVIDERS.md)

## Using as a Go Module

For information on using TTS Bulk as a library in your Go applications, see [docs/USAGE.md](docs/USAGE.md).

## Pull Request Process

1. Ensure your code passes all tests and lints cleanly
2. Update documentation if you're changing functionality
3. Add a clear description of your changes in the PR
4. Link any related issues in the PR description
5. Wait for review from maintainers
6. Address any feedback and push updates to your branch

## Reporting Issues

When reporting issues, please include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Your environment (OS, Go version, etc.)
- Any relevant logs or error messages

## Questions?

If you have questions about contributing, feel free to:

- Open an issue with the "question" label
- Start a discussion in the Discussions tab

## License

By contributing to TTS Bulk, you agree that your contributions will be licensed under the Apache License 2.0.

Thank you for contributing! 🎉
