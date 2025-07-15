# Tama Go Client Library

A Go client library for the Tama API, providing easy access to Neural and Sensory provisioning endpoints.

## Installation

```bash
go get github.com/upmaru/tama-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    tama "github.com/upmaru/tama-go"
    "github.com/upmaru/tama-go/sensory"
)

func main() {
    // Initialize the client
    config := tama.Config{
        BaseURL: "https://api.tama.io",
        APIKey:  "your-api-key",
        Timeout: 30 * time.Second,
    }
    
    client := tama.NewClient(config)
    
    // Create a neural space
    space, err := client.Neural.CreateSpace(tama.CreateSpaceRequest{
        Space: tama.SpaceRequest{
            Name: "My Neural Space",
            Type: "root",
        },
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created space: %+v\n", space)
}
```

## API Coverage

The client provides comprehensive coverage of the Tama API endpoints, organized by resource type:

### Neural Resources (`/provision/neural`)

#### Spaces
- `GET /provision/neural/spaces/:id` - Get space by ID
- `POST /provision/neural/spaces` - Create new space
- `PATCH /provision/neural/spaces/:id` - Update space
- `PUT /provision/neural/spaces/:id` - Replace space
- `DELETE /provision/neural/spaces/:id` - Delete space

### Sensory Resources (`/provision/sensory`)

#### Sources
- `GET /provision/sensory/sources/:id` - Get source by ID
- `POST /provision/sensory/spaces/:space_id/sources` - Create source in space
- `PATCH /provision/sensory/sources/:id` - Update source
- `PUT /provision/sensory/sources/:id` - Replace source
- `DELETE /provision/sensory/sources/:id` - Delete source

#### Models
- `GET /provision/sensory/models/:id` - Get model by ID
- `POST /provision/sensory/sources/:source_id/models` - Create model for source
- `PATCH /provision/sensory/models/:id` - Update model
- `PUT /provision/sensory/models/:id` - Replace model
- `DELETE /provision/sensory/models/:id` - Delete model

#### Limits
- `GET /provision/sensory/limits/:id` - Get limit by ID
- `POST /provision/sensory/sources/:source_id/limits` - Create limit for source
- `PATCH /provision/sensory/limits/:id` - Update limit
- `PUT /provision/sensory/limits/:id` - Replace limit
- `DELETE /provision/sensory/limits/:id` - Delete limit

## Usage Examples

### Neural Service - Spaces

```go
// Create a space
space, err := client.Neural.CreateSpace(tama.CreateSpaceRequest{
    Space: tama.SpaceRequest{
        Name: "Production Space",
        Type: "root",
    },
})

// Get a space
space, err := client.Neural.GetSpace("space-123")

// Update a space (partial update)
space, err := client.Neural.UpdateSpace("space-123", tama.UpdateSpaceRequest{
    Space: tama.UpdateSpaceData{
        Name: "Updated Production Space",
        Type: "component",
    },
})

// Replace a space (full replacement)
space, err := client.Neural.ReplaceSpace("space-123", tama.UpdateSpaceRequest{
    Space: tama.UpdateSpaceData{
        Name: "New Production Space",
        Type: "root",
    },
})

// Delete a space
err := client.Neural.DeleteSpace("space-123")
```

### Sensory Service - Sources

```go
// Create a source in a space
source, err := client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
    Source: sensory.SourceRequestData{
        Name: "Mistral Source",
        Type: "model",
        Endpoint: "https://api.mistral.ai/v1",
        Credential: sensory.SourceCredential{
            ApiKey: "your-api-key",
        },
    },
})

// Get a source
source, err := client.Sensory.GetSource("source-123")

// Update a source
source, err := client.Sensory.UpdateSource("source-123", sensory.UpdateSourceRequest{
    Source: sensory.UpdateSourceData{
        Name: "Updated Mistral Source",
        Endpoint: "https://api.mistral.ai/v2",
        Credential: &sensory.SourceCredential{
            ApiKey: "your-updated-api-key",
        },
    },
})

// Delete a source
err := client.Sensory.DeleteSource("source-123")
```

### Sensory Service - Models

```go
// Create a model for a source
model, err := client.Sensory.CreateModel("source-123", sensory.CreateModelRequest{
    Model: sensory.ModelRequestData{
        Identifier: "mistral-small-latest",
        Path:       "/chat/completions",
    },
})

// Get a model
model, err := client.Sensory.GetModel("model-123")

// Update a model
model, err := client.Sensory.UpdateModel("model-123", sensory.UpdateModelRequest{
    Model: sensory.UpdateModelData{
        Identifier: "mistral-large-latest",
        Path:       "/chat/completions",
    },
})

// Delete a model
err := client.Sensory.DeleteModel("model-123")
```

### Sensory Service - Limits

```go
// Create a limit for a source
limit, err := client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
    Limit: sensory.LimitRequestData{
        ScaleUnit:  "seconds",
        ScaleCount: 1,
        Limit:      32,
    },
})

// Get a limit
limit, err := client.Sensory.GetLimit("limit-123")

// Update a limit
limit, err := client.Sensory.UpdateLimit("limit-123", sensory.UpdateLimitRequest{
    Limit: sensory.UpdateLimitData{
        ScaleUnit:  "minutes",
        ScaleCount: 5,
        Limit:      100,
    },
})

// Delete a limit
err := client.Sensory.DeleteLimit("limit-123")
```

## Configuration

### Client Configuration

```go
config := tama.Config{
    BaseURL: "https://api.tama.io",  // Required: API base URL
    APIKey:  "your-api-key",         // Required: API authentication key
    Timeout: 30 * time.Second,       // Optional: Request timeout (default: 30s)
}

client := tama.NewClient(config)
```

### Authentication

The client supports API key authentication. Set your API key in the config:

```go
client.SetAPIKey("your-new-api-key")
```

### Debug Mode

Enable debug mode to see HTTP request/response details:

```go
client.SetDebug(true)
```

## Error Handling

The client provides structured error handling:

```go
space, err := client.Neural.GetSpace("invalid-id")
if err != nil {
    if apiErr, ok := err.(*tama.Error); ok {
        fmt.Printf("API Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
        if apiErr.Details != "" {
            fmt.Printf("Details: %s\n", apiErr.Details)
        }
    } else {
        fmt.Printf("Client Error: %v\n", err)
    }
}
```

## Data Types

### Core Resources

- **Space**: Neural space resource with configuration
- **Source**: Sensory data source with type and connection details
- **Model**: Machine learning model with version and parameters
- **Limit**: Resource limits with values and units

### Request Types

- **CreateSpaceRequest**: For creating new spaces
- **UpdateSpaceRequest**: For updating existing spaces
- **CreateSourceRequest**: For creating new sources
- **UpdateSourceRequest**: For updating existing sources
- **CreateModelRequest**: For creating new models
- **UpdateModelRequest**: For updating existing models
- **CreateLimitRequest**: For creating new limits
- **UpdateLimitRequest**: For updating existing limits

## Examples

See the [example/main.go](example/main.go) file for a complete working example demonstrating all client features.

## Requirements

- Go 1.23 or later
- Active Tama API credentials

## Dependencies

- [go-resty/resty](https://github.com/go-resty/resty) - HTTP client library

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Project Structure

The client library is organized into the following files and packages:

### Main Package
- `client.go` - Main client configuration and initialization
- `neural.go` - Neural service for space operations
- `sensory.go` - SensoryService wrapper that uses the sensory package
- `types.go` - Core data structures (Space, Error, etc.)

### Sensory Package (`sensory/`)
- `service.go` - Service definition and all sensory-related types
- `source.go` - Source operations (GET, POST, PATCH, PUT, DELETE)
- `model.go` - Model operations (GET, POST, PATCH, PUT, DELETE)
- `limit.go` - Limit operations (GET, POST, PATCH, PUT, DELETE)

### Examples
- `example/` - Working examples demonstrating all features

This modular structure separates concerns into different packages, making the codebase easier to navigate, maintain, and extend. The sensory package encapsulates all sensory-related functionality with its own types and operations.

## Support

For issues and questions:
- Create an issue on GitHub
- Check the API documentation
- Review the examples in this repository