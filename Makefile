.PHONY: build test clean install lint fmt vet run help

# Variables
BINARY_NAME=tts-bulk
MAIN_PATH=./cmd/tts-bulk
BUILD_DIR=builds
GO=go
GOFLAGS=-v

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

## build-all: Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/darwin-amd64/$(BINARY_NAME) $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/darwin-arm64/$(BINARY_NAME) $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "Build complete!"

## test: Run tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

## install: Install dependencies
install:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

## lint: Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## run: Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

## check: Run fmt, vet, and lint
check: fmt vet lint
	@echo "All checks passed!"
