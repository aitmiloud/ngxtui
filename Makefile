.PHONY: build run clean install test lint fmt help

# Binary name
BINARY_NAME=ngxtui
BINARY_PATH=./bin/$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Main package path
MAIN_PATH=./cmd/ngxtui

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

# Run the application (requires sudo)
run: build
	@echo "Running $(BINARY_NAME)..."
	@sudo $(BINARY_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies installed"

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Lint the code
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2" && exit 1)
	golangci-lint run ./...

# Format the code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Format complete"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "Build complete for all platforms"

# Development mode with auto-reload (requires entr)
dev:
	@echo "Starting development mode..."
	@which entr > /dev/null || (echo "entr not installed. Install with: apt-get install entr (Linux) or brew install entr (macOS)" && exit 1)
	@find . -name "*.go" | entr -r make run

# Show help
help:
	@echo "NgxTUI Makefile Commands:"
	@echo ""
	@echo "  make build         - Build the application"
	@echo "  make run           - Build and run the application (requires sudo)"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make deps          - Install dependencies"
	@echo "  make install       - Install binary to GOPATH/bin"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"
	@echo "  make build-all     - Build for multiple platforms"
	@echo "  make dev           - Run in development mode with auto-reload"
	@echo "  make help          - Show this help message"
	@echo ""

# Default target
.DEFAULT_GOAL := help
