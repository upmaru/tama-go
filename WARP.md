# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

Tama Go Client Library is a comprehensive Go client for the Tama API, providing access to Neural, Sensory, Memory, Perception, Motor, and Contexts provisioning endpoints. The library is built with a modular service-oriented architecture.

## Common Development Commands

### Building and Testing
```bash
# Run unit tests
make test

# Run tests with coverage report
make test-coverage

# Run integration tests (requires environment variables)
export TAMA_BASE_URL="https://api.tama.io"
export TAMA_API_KEY="your-api-key"
export TAMA_TEST_SPACE_ID="space-123"
make test-integration

# Build the example application
make build

# Run the example application
make example
```

### Code Quality and Formatting
```bash
# Format code
make fmt

# Run linter (uses strict golangci-lint configuration)
make lint

# Tidy Go modules
make mod-tidy

# Run all quality checks (format, lint, test)
make check

# Run complete CI checks locally
make ci-check

# Prepare for release (all checks including security)
make release-check
```

### Development Tools
```bash
# Install development tools (golangci-lint, gosec, govulncheck)
make install-tools

# Security scanning
make security-scan

# Vulnerability checking
make vulnerability-check
```

### Running Specific Tests
```bash
# Run tests for a specific package
go test -v ./neural
go test -v ./sensory
go test -v ./perception

# Run a single test function
go test -v ./neural -run TestCreateSpace

# Run tests with race detection
go test -race ./...
```

## Architecture Overview

### Service-Oriented Design
The codebase follows a modular service pattern where each API domain is encapsulated in its own package:

- **Main Client** (`client.go`): Central client that orchestrates all services
- **Neural Service** (`neural/`): Manages spaces, processors, classes, corpora, bridges, and nodes
- **Sensory Service** (`sensory/`): Handles sources, models, limits, specifications, and identities
- **Memory Service** (`memory/`): Manages prompts and memory operations
- **Perception Service** (`perception/`): Handles chains, thoughts, paths, contexts, processors, activations, and initializers
- **Motor Service** (`motor/`): Manages actions
- **Contexts Service** (`contexts/`): Handles input operations

### Error Handling Architecture
Each service implements sophisticated error handling with nested error parsing:
- Service-specific error types (`neural.Error`, `sensory.Error`, etc.)
- Support for field validation errors with dot-notation keys (e.g., `module.config.database.connection`)
- Fallback error handling for different error response formats
- Unified error interface across all services

### HTTP Client Pattern
All services use a shared HTTP client pattern:
- Resty-based HTTP client with configurable timeouts
- Bearer token authentication
- Debug mode support
- Consistent request/response handling across services

### Testing Structure
- Unit tests for each service component (`*_test.go`)
- Shared test utilities (`shared_test.go`)
- Integration tests requiring live API credentials
- Mock-friendly architecture for testing

## Key Development Patterns

### Service Initialization
Services are initialized through the main client:
```go
client := tama.NewClient(config)
// Services are automatically available:
// client.Neural, client.Sensory, client.Memory, etc.
```

### Request/Response Pattern
All services follow consistent request/response patterns:
- Create operations: `Create{Resource}Request` with `{Resource}RequestData`
- Update operations: `Update{Resource}Request` with `Update{Resource}Data`
- Response wrappers: `{Resource}Response` containing the resource data
- Resources have consistent fields: `ID`, `ProvisionState`, service-specific data

### Module System (Perception Service)
The perception service has a nested module system:
- Thoughts contain modules with references and parameters
- Module types: `tama/agentic/generate`, `tama/agentic/analyze`, `tama/agentic/preprocess`, `tama/agentic/validate`
- Module inputs are managed through the `perception/module/` subpackage

### Configuration Flexibility
Resources support flexible configuration through `map[string]any` fields:
- Neural processors have type-specific configurations (completion, embedding, reranking)
- Sensory models accept arbitrary parameters for different AI providers
- Perception modules support various parameter sets
- All configurations are validated server-side

## Code Quality Standards

### Linting Configuration
The project uses a strict golangci-lint configuration (`.golangci.yml`):
- 50+ enabled linters covering security, performance, style, and correctness
- Custom rules for deprecated packages and import organization
- Local prefix configuration for import ordering: `github.com/upmaru/tama-go`
- Maximum line length: 120 characters
- Function complexity limits and cognitive complexity checks

### Go Version Requirements
- Minimum Go version: 1.23
- Uses modern Go features like range-over-int (Go 1.22+)
- Dependencies managed through Go modules

### Testing Requirements
- Unit tests for all public functions
- Integration tests for API operations (tagged with `integration`)
- Coverage reporting through `make test-coverage`
- Race detection enabled in CI

## Important Implementation Details

### Error Response Parsing
Services implement multi-format error parsing to handle various API error structures:
1. Nested errors with arbitrary depth (e.g., `module.reference.parameters.temperature`)
2. Array format errors (`map[string][]string`)
3. Single string format errors (converted to array format)
4. Fallback error handling for unknown formats

### Resource Relationships
Understanding the hierarchical relationships is crucial:
- **Spaces** contain classes, processors, bridges, chains, and sources
- **Classes** contain corpora and are referenced by thoughts
- **Sources** contain models and limits
- **Chains** contain thoughts, which contain paths and contexts
- **Thoughts** contain modules and can reference output classes

### Provision State Management
All resources have `ProvisionState` fields that track server-side status:
- States are read-only from the client perspective
- Common states: `pending`, `provisioning`, `active`, `failed`
- State transitions are managed server-side

### Credential Handling
Sensitive data (API keys, credentials) should be handled through:
- Environment variables for configuration
- Secure credential structures in request data
- Never log or expose credentials in debug output

## Working with Examples

The `example/main.go` file provides comprehensive usage examples for all services. Key patterns demonstrated:
- Client initialization and configuration
- CRUD operations for each resource type
- Error handling best practices
- Schema creation for complex data structures
- Integration between services (e.g., creating chains with thoughts that reference classes)

When adding new features, update the example file to demonstrate proper usage patterns.
