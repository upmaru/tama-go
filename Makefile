.PHONY: help test test-integration test-coverage build clean fmt lint mod-tidy example security-scan vulnerability-check ci-check install-tools check release-check

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
	@echo "  security-scan      Run Gosec security scanner"
	@echo "  vulnerability-check Run govulncheck for vulnerabilities"
	@echo "  ci-check           Run all CI checks locally"
	@echo "  install-tools      Install development tools"
	@echo "  check              Run format, lint, and test"
	@echo "  release-check      Run all checks for release preparation"

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

# Run security scan with Gosec
security-scan:
	@echo "Running Gosec security scan..."
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; }
	gosec -no-fail -fmt sarif -out gosec-results.sarif ./...
	gosec ./...

# Run vulnerability check
vulnerability-check:
	@echo "Running govulncheck..."
	@command -v govulncheck >/dev/null 2>&1 || { echo "Installing govulncheck..."; go install golang.org/x/vuln/cmd/govulncheck@latest; }
	govulncheck ./...

# Install development dependencies
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "Development tools installed successfully"

# Run all checks (format, lint, test)
check: fmt lint test

# Run all CI checks locally
ci-check: mod-tidy fmt lint test security-scan vulnerability-check
	@echo "All CI checks completed successfully"

# Release preparation
release-check: mod-tidy fmt lint test test-coverage security-scan vulnerability-check
	@echo "Release checks completed successfully"
