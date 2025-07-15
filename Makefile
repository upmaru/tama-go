.PHONY: help test test-integration test-coverage build clean fmt lint mod-tidy example

# Default target
help:
	@echo "Available targets:"
	@echo "  test               Run unit tests"
	@echo "  test-integration   Run integration tests (requires env vars)"
	@echo "  test-coverage      Run tests with coverage report"
	@echo "  build              Build the example application"
	@echo "  clean              Clean build artifacts"
	@echo "  fmt                Format Go code"
	@echo "  lint               Run golangci-lint"
	@echo "  mod-tidy           Tidy Go modules"
	@echo "  example            Run the example application"

# Run unit tests
test:
	go test -v ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@echo "Make sure to set TAMA_BASE_URL, TAMA_API_KEY, and TAMA_TEST_SPACE_ID environment variables"
	go test -tags=integration -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the example application
build:
	go build -o bin/tama-example ./example

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format Go code
fmt:
	go fmt ./...

# Run golangci-lint (requires golangci-lint to be installed)
lint:
	golangci-lint run

# Tidy Go modules
mod-tidy:
	go mod tidy

# Run the example application
example: build
	./bin/tama-example

# Install development dependencies
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all checks (format, lint, test)
check: fmt lint test

# Release preparation
release-check: mod-tidy fmt lint test test-coverage
	@echo "Release checks completed successfully"
