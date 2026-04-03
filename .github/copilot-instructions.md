# Copilot Instructions

## Project Overview

`tts-bulk` is a CLI tool that reads a 2-column CSV (label, text), converts each row to audio via the Deepgram TTS API, and writes the output files concurrently with configurable parallelism (default: 3) and a 1-second rate limit per request.

## Build, Test, and Lint

```bash
make build          # Build binary for current platform → ./tts-bulk
make test           # Run all tests with race detector + coverage
make check          # fmt + vet + lint
make lint           # golangci-lint run ./...
make clean          # Remove build artifacts
```

Run a single test:

```bash
go test -v ./internal/config/... -run TestLoad
go test -v -race ./...                         # All tests with race detector
```

Linter is `golangci-lint` with shadow checking, `goimports` (local-prefix: `github.com/aashish-joshi/tts-bulk`), `errcheck`, `revive`, and `misspell` (US locale).

## Architecture

```
cmd/tts-bulk/main.go          CLI entry point, flag parsing, top-level orchestration
internal/config/              Config loading from env vars + CLI flags, validation
internal/logger/              Structured singleton logger (4 levels, thread-safe)
internal/processor/           Batch CSV processing, semaphore-based concurrency
internal/provider/deepgram/   Deepgram TTS API implementation
pkg/types/                    Public Provider interface, TTSRequest/TTSResult, AudioFormat constants
```

### Data Flow

```
CLI flags + DEEPGRAM_API_KEY env var
→ config.Load()          validates inputs, sets defaults
→ deepgram.New()         initializes SDK client
→ processor.New()        creates output directory
→ processor.ProcessCSV() reads CSV, spawns goroutines (bounded by semaphore)
  → provider.GenerateAudio() calls Deepgram API, sleeps 1s
  → writes file to output dir
→ collect []TTSResult, log summary, exit 1 if any failures
```

### Provider Interface (`pkg/types/types.go`)

```go
type Provider interface {
    GenerateAudio(ctx context.Context, req TTSRequest, outputPath string) error
    Name() string
    Close() error
}
```

New providers implement this interface and are wired in `main.go`. See `docs/ADDING_PROVIDERS.md`.

## Key Conventions

### Error Handling

- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- Validate eagerly in `config.Load()` — all checks happen at startup.
- Use result collection (store `err` in `TTSResult.Error`) rather than panicking in goroutines.

### Logging

- Use package-level functions: `logger.Info("message: %s", val)`, `logger.Error("failed: %v", err)`.
- Levels: Debug (verbose flag), Info, Warn, Error — configured globally, not per call-site.
- Format: `[HH:MM:SS] [LEVEL] message`. Never log API keys.

### Concurrency

- Semaphore pattern via buffered channel: `sem := make(chan struct{}, maxConcurrent)`.
- Goroutines acquire with `sem <- struct{}{}` and release with `defer func() { <-sem }()`.
- `sync.WaitGroup` for synchronization; results collected into a shared slice (goroutines write to pre-indexed slots).

### Configuration

- Single source of truth: `config.Load()` — config is immutable after creation.
- Defaults: provider=deepgram, format=mp3, output=audio, model=aura-asteria-en, maxConcurrent=3.
- `config.Config.GetDeepgramConfig()` maps user-facing format string to Deepgram container/encoding.

### Tests

- Table-driven tests with `[]struct{ name, input, want, wantErr }`.
- Mock providers implement `types.Provider`; use `t.TempDir()` for file I/O.
- Tests live alongside source as `*_test.go` files.

### Naming

- Interfaces: short and minimal (`Provider`).
- Log level constants: `LevelDebug`, `LevelInfo`, etc.
- Audio format constants: `FormatMP3`, `FormatWAV`.
- Unexported struct fields for internal state (`.mu`, `.level`, `.logger`).

## Documentation (Markdown)

Codacy enforces markdownlint rules on all `.md` files. Follow these rules when writing or editing documentation:

- **No tabs** (MD010): Use spaces for indentation inside fenced code blocks, not tab characters.
- **Blank line before lists** (MD032): Always add a blank line before ordered (`1.`) and unordered (`-`, `*`) lists.

Example — correct:

```markdown
The following formats are supported:

- mp3
- wav
```

Example — incorrect (no blank line before list, tabs in code block):

```markdown
The following formats are supported:
- mp3
- wav
```
