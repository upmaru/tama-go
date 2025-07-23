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
    "github.com/upmaru/tama-go/neural"
    "github.com/upmaru/tama-go/perception"
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
    space, err := client.Neural.CreateSpace(neural.CreateSpaceRequest{
        Space: neural.SpaceRequestData{
            Name: "My Neural Space",
            Type: "root",
        },
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created space: ID=%s, Name=%s, Type=%s, State=%s\n", 
        space.ID, space.Name, space.Type, space.CurrentState)
    
    // Create a source in the space
    source, err := client.Sensory.CreateSource(space.ID, sensory.CreateSourceRequest{
        Source: sensory.SourceRequestData{
            Name:     "AI Model Source",
            Type:     "model",
            Endpoint: "https://api.example.com/v1",
            Credential: sensory.SourceCredential{
                APIKey: "source-api-key",
            },
        },
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created source: ID=%s, Name=%s, Endpoint=%s, SpaceID=%s, State=%s\n", 
        source.ID, source.Name, source.Endpoint, source.SpaceID, source.CurrentState)
    
    // Create a limit for the source
    limit, err := client.Sensory.CreateLimit(source.ID, sensory.CreateLimitRequest{
        Limit: sensory.LimitRequestData{
            ScaleUnit:  "minutes",
            ScaleCount: 1,
            Count:      100,
        },
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created limit: ID=%s, SourceID=%s, Count=%d, State=%s\n", 
        limit.ID, limit.SourceID, limit.Count, limit.CurrentState)
    
    // Create a perception chain
    chain, err := client.Perception.CreateChain(space.ID, perception.CreateChainRequest{
        Chain: perception.ChainRequestData{
            Name: "AI Processing Chain",
        },
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created chain: ID=%s, Name=%s, SpaceID=%s, State=%s\n", 
        chain.ID, chain.Name, chain.SpaceID, chain.CurrentState)
    
    // Create a thought in the chain
    thought, err := client.Perception.CreateThought(chain.ID, perception.CreateThoughtRequest{
        Thought: perception.ThoughtRequestData{
            Relation: "description",
            Module: perception.Module{
                Reference: "tama/agentic/generate",
                Parameters: map[string]any{
                    "temperature": 0.7,
                    "max_tokens":  150,
                    "model":       "gpt-4",
                },
            },
        },
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created thought: ID=%s, ChainID=%s, Relation=%s, State=%s\n", 
        thought.ID, thought.ChainID, thought.Relation, thought.CurrentState)
}
```

## Project Structure

The client library is organized into the following packages:

### Main Package
- `client.go` - Main client configuration and initialization
- `neural.go` - Neural service wrapper that uses the neural package
- `sensory.go` - Sensory service wrapper that uses the sensory package
- `perception.go` - Perception service wrapper that uses the perception package
- `types.go` - Shared types and documentation

### Neural Package (`neural/`)
- `service.go` - Service definition and neural-related types
- `space.go` - Space operations (GET, POST, PATCH, PUT, DELETE)
- `processor.go` - Processor operations (GET, POST, PATCH, PUT, DELETE)

### Sensory Package (`sensory/`)
- `service.go` - Service definition and sensory-related types
- `source.go` - Source operations (GET, POST, PATCH, PUT, DELETE)
- `model.go` - Model operations (GET, POST, PATCH, PUT, DELETE)
- `limit.go` - Limit operations (GET, POST, PATCH, PUT, DELETE)

### Perception Package (`perception/`)
- `service.go` - Service definition and perception-related types
- `chain.go` - Chain operations (GET, POST, PATCH, PUT, DELETE)
- `thought.go` - Thought operations (GET, POST, PATCH, DELETE)

### Examples
- `example/` - Working examples demonstrating all features

This modular structure separates concerns into different packages, making the codebase easier to navigate, maintain, and extend. Each service package encapsulates its related functionality with its own types and operations.

## API Coverage

The client provides comprehensive coverage of the Tama API endpoints, organized by resource type:

### Neural Resources (`/provision/neural`)

#### Spaces
- `GET /provision/neural/spaces/:id` - Get space by ID
- `POST /provision/neural/spaces` - Create new space
- `PATCH /provision/neural/spaces/:id` - Update space
- `PUT /provision/neural/spaces/:id` - Replace space
- `DELETE /provision/neural/spaces/:id` - Delete space

#### Processors
- `GET /provision/neural/spaces/:space_id/models/:model_id/processor` - Get processor
- `POST /provision/neural/spaces/:space_id/models/:model_id/processor` - Create processor
- `PATCH /provision/neural/spaces/:space_id/models/:model_id/processor` - Update processor
- `PUT /provision/neural/spaces/:space_id/models/:model_id/processor` - Replace processor
- `DELETE /provision/neural/spaces/:space_id/models/:model_id/processor` - Delete processor

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

Note: Limits are associated with sources via the `source_id` field and track resource usage counts with current state.

### Perception Resources (`/provision/perception`)

#### Chains
- `GET /provision/perception/chains/:id` - Get chain by ID
- `POST /provision/perception/spaces/:space_id/chains` - Create chain in space
- `PATCH /provision/perception/chains/:id` - Update chain
- `PUT /provision/perception/chains/:id` - Replace chain
- `DELETE /provision/perception/chains/:id` - Delete chain

#### Thoughts
- `GET /provision/perception/thoughts/:id` - Get thought by ID
- `POST /provision/perception/chains/:chain_id/thoughts` - Create thought in chain
- `PATCH /provision/perception/thoughts/:id` - Update thought
- `DELETE /provision/perception/thoughts/:id` - Delete thought

Note: Thoughts are associated with chains and contain module configurations for AI processing operations.

## Usage Examples

### Neural Service - Spaces

```go
import "github.com/upmaru/tama-go/neural"

// Create a space
space, err := client.Neural.CreateSpace(neural.CreateSpaceRequest{
    Space: neural.SpaceRequestData{
        Name: "Production Space",
        Type: "root",
    },
})
// space will have ID, Name, Slug, Type, and CurrentState populated

// Get a space
space, err := client.Neural.GetSpace("space-123")

// Update a space (partial update)
space, err := client.Neural.UpdateSpace("space-123", neural.UpdateSpaceRequest{
    Space: neural.UpdateSpaceData{
        Name: "Updated Production Space",
        Type: "component",
    },
})
// CurrentState cannot be updated via API - it's managed server-side

// Replace a space (full replacement)
space, err := client.Neural.ReplaceSpace("space-123", neural.UpdateSpaceRequest{
    Space: neural.UpdateSpaceData{
        Name: "New Production Space",
        Type: "root",
    },
})

// Delete a space
err := client.Neural.DeleteSpace("space-123")
```

### Neural Service - Processors

```go
import "github.com/upmaru/tama-go/neural"

// Create a processor
processor, err := client.Neural.CreateProcessor("space-123", "model-123", neural.CreateProcessorRequest{
    Processor: neural.ProcessorRequestData{
        Type: "completion",
        Configuration: map[string]any{
            "temperature":  0.8,
            "tool_choice": "required",
            "role_mappings": []map[string]any{
                {"from": "user", "to": "human"},
                {"from": "assistant", "to": "ai"},
            },
        },
    },
})
// processor will have ID, SpaceID, ModelID, Type, Configuration, and CurrentState populated

// Get a processor
processor, err := client.Neural.GetProcessor("space-123", "model-123")

// Update a processor (partial update)
processor, err := client.Neural.UpdateProcessor("space-123", "model-123", neural.UpdateProcessorRequest{
    Processor: neural.UpdateProcessorData{
        Type: "embedding",
        Configuration: map[string]any{
            "max_tokens": 512,
            "templates": []map[string]any{
                {"type": "query", "content": "Query: {text}"},
            },
        },
    },
})
// CurrentState cannot be updated via API - it's managed server-side

// Replace a processor (full replacement)
processor, err := client.Neural.ReplaceProcessor("space-123", "model-123", neural.UpdateProcessorRequest{
    Processor: neural.UpdateProcessorData{
        Type: "reranking",
        Configuration: map[string]any{
            "top_n": 3,
        },
    },
})

// Delete a processor
err := client.Neural.DeleteProcessor("space-123", "model-123")
```

#### Processor Types and Configuration

Processors support three types: `"completion"`, `"embedding"`, and `"reranking"`. Each type has its own configuration schema:

##### Completion Type

For text completion and chat completion tasks:

```go
processor, err := client.Neural.CreateProcessor("space-123", "model-123", neural.CreateProcessorRequest{
    Processor: neural.ProcessorRequestData{
        Type: "completion",
        Configuration: map[string]any{
            "temperature":  0.8,  // decimal, default: 0.8
            "tool_choice": "required", // enum: "required", "auto", "any", default: "required"
            "role_mappings": []map[string]any{
                {
                    "from": "user",
                    "to":   "human",
                },
                {
                    "from": "assistant", 
                    "to":   "ai",
                },
            },
        },
    },
})
```

##### Embedding Type

For text embedding and vector generation:

```go
processor, err := client.Neural.CreateProcessor("space-123", "model-123", neural.CreateProcessorRequest{
    Processor: neural.ProcessorRequestData{
        Type: "embedding",
        Configuration: map[string]any{
            "max_tokens": 512, // integer, default: 512
            "templates": []map[string]any{
                {
                    "type":    "query",
                    "content": "Query: {text}",
                },
                {
                    "type":    "document",
                    "content": "Document: {text}",
                },
            },
        },
    },
})
```

##### Reranking Type

For document reranking and relevance scoring:

```go
processor, err := client.Neural.CreateProcessor("space-123", "model-123", neural.CreateProcessorRequest{
    Processor: neural.ProcessorRequestData{
        Type: "reranking",
        Configuration: map[string]any{
            "top_n": 3, // integer, default: 3
        },
    },
})
```

##### Configuration Field Details

**Completion Configuration:**
- `temperature` (decimal): Controls randomness in generation, default: 0.8
- `tool_choice` (string): Tool selection strategy - "required", "auto", or "any", default: "required"
- `role_mappings` (array): Maps input roles to model-specific roles
  - `from` (string): Input role name
  - `to` (string): Model role name

**Embedding Configuration:**
- `max_tokens` (integer): Maximum tokens to process, default: 512
- `templates` (array): Text templates for different embedding types
  - `type` (string): Template type - "query" or "document"
  - `content` (string): Template string with {text} placeholder

**Reranking Configuration:**
- `top_n` (integer): Number of top results to return, default: 3

### Sensory Service - Sources

```go
import "github.com/upmaru/tama-go/sensory"

// Create a source in a space
source, err := client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
    Source: sensory.SourceRequestData{
        Name: "Mistral Source",
        Type: "model",
        Endpoint: "https://api.mistral.ai/v1",
        Credential: sensory.SourceCredential{
            APIKey: "your-api-key",
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
            APIKey: "your-updated-api-key",
        },
    },
})

// Delete a source
err := client.Sensory.DeleteSource("source-123")
```

### Sensory Service - Models

```go
import "github.com/upmaru/tama-go/sensory"

// Create a model for a source
model, err := client.Sensory.CreateModel("source-123", sensory.CreateModelRequest{
    Model: sensory.ModelRequestData{
        Identifier: "mistral-small-latest",
        Path:       "/chat/completions",
        Parameters: map[string]any{
            "reasoning_effort": "low",
            "temperature":      1.0,
            "max_tokens":       2000,
            "stream":           true,
            "stop":             []string{"\n", "###"},
            "config": map[string]any{
                "timeout":      30,
                "enable_cache": true,
            },
        },
    },
})

// Get a model
model, err := client.Sensory.GetModel("model-123")

// Update a model
model, err := client.Sensory.UpdateModel("model-123", sensory.UpdateModelRequest{
    Model: sensory.UpdateModelData{
        Identifier: "mistral-large-latest",
        Path:       "/chat/completions",
        Parameters: map[string]any{
            "temperature": 0.8,
            "max_tokens":  1500,
        },
    },
})

// Delete a model
err := client.Sensory.DeleteModel("model-123")
```

#### Model Parameters

The `Parameters` field in models accepts any valid JSON values, allowing flexible configuration:

```go
Parameters: map[string]any{
    // String values
    "reasoning_effort": "low",
    
    // Numeric values
    "temperature":       0.8,
    "max_tokens":        1500,
    "frequency_penalty": 0.1,
    
    // Boolean values
    "stream": true,
    
    // Array values
    "stop": []string{"\n", "###", "END"},
    
    // Object values
    "config": map[string]any{
        "timeout":      30,
        "enable_cache": true,
        "retries":      3,
    },
}
```

### Sensory Service - Limits

```go
import "github.com/upmaru/tama-go/sensory"

// Create a limit for a source
limit, err := client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
    Limit: sensory.LimitRequestData{
        ScaleUnit:  "seconds",
        ScaleCount: 1,
        Count:      32,
    },
})

// Get a limit
limit, err := client.Sensory.GetLimit("limit-123")

// Update a limit
limit, err := client.Sensory.UpdateLimit("limit-123", sensory.UpdateLimitRequest{
    Limit: sensory.UpdateLimitData{
        ScaleUnit:    "minutes",
        ScaleCount:   5,
        Count:        100,
        CurrentState: "active",
    },
})

// Delete a limit
err := client.Sensory.DeleteLimit("limit-123")
```

### Perception Service - Chains

```go
// Create a chain
chain, err := client.Perception.CreateChain("space-123", perception.CreateChainRequest{
    Chain: perception.ChainRequestData{
        Name: "Processing Chain",
    },
})

// Get a chain
chain, err := client.Perception.GetChain("chain-123")

// Update a chain
chain, err := client.Perception.UpdateChain("chain-123", perception.UpdateChainRequest{
    Chain: perception.UpdateChainData{
        Name: "Updated Processing Chain",
    },
})

// Delete a chain
err := client.Perception.DeleteChain("chain-123")
```

### Perception Service - Thoughts

```go
// Create a thought
thought, err := client.Perception.CreateThought("chain-123", perception.CreateThoughtRequest{
    Thought: perception.ThoughtRequestData{
        Relation: "description",
        Module: perception.Module{
            Reference: "tama/agentic/generate",
            Parameters: map[string]any{
                "temperature": 0.7,
                "max_tokens":  150,
                "model":       "gpt-4",
            },
        },
    },
})

// Get a thought
thought, err := client.Perception.GetThought("thought-123")

// Update a thought
thought, err := client.Perception.UpdateThought("thought-123", perception.UpdateThoughtRequest{
    Thought: perception.UpdateThoughtData{
        Relation: "analysis",
        Module: perception.Module{
            Reference: "tama/agentic/analyze",
            Parameters: map[string]any{
                "depth":       3,
                "focus_areas": []string{"sentiment", "intent", "entities"},
            },
        },
    },
})

// Delete a thought
err := client.Perception.DeleteThought("thought-123")
```

#### Thought Module Configuration

Thoughts contain module configurations that define AI processing operations:

##### Generate Module
```go
Module: perception.Module{
    Reference: "tama/agentic/generate",
    Parameters: map[string]any{
        "temperature": 0.7,
        "max_tokens":  150,
        "model":       "gpt-4",
        "prompt":      "Generate a summary of the input text",
    },
}
```

##### Analyze Module
```go
Module: perception.Module{
    Reference: "tama/agentic/analyze",
    Parameters: map[string]any{
        "depth":       3,
        "focus_areas": []string{"sentiment", "intent", "entities"},
        "output_format": "structured",
    },
}
```

##### Preprocess Module
```go
Module: perception.Module{
    Reference: "tama/agentic/preprocess",
    Parameters: map[string]any{
        "clean_text":  true,
        "normalize":   true,
        "remove_html": true,
    },
}
```

##### Validate Module
```go
Module: perception.Module{
    Reference: "tama/agentic/validate",
    Parameters: map[string]any{
        "strict_mode":    false,
        "schema_version": "v2",
        "required_fields": []string{"input", "output"},
    },
}
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

The client provides structured error handling with service-specific error types:

### Neural Service Errors

```go
import "github.com/upmaru/tama-go/neural"

space, err := client.Neural.GetSpace("invalid-id")
if err != nil {
    if apiErr, ok := err.(*neural.Error); ok {
        fmt.Printf("Neural API Error %d\n", apiErr.StatusCode)
    } else {
        fmt.Printf("Client Error: %v\n", err)
    }
}
```

### Sensory Service Errors

```go
import "github.com/upmaru/tama-go/sensory"

source, err := client.Sensory.GetSource("invalid-id")
if err != nil {
    if apiErr, ok := err.(*sensory.Error); ok {
        fmt.Printf("Sensory API Error %d\n", apiErr.StatusCode)
    } else {
        fmt.Printf("Client Error: %v\n", err)
    }
}
```

### Perception Service Errors

```go
import "github.com/upmaru/tama-go/perception"

chain, err := client.Perception.GetChain("invalid-id")
if err != nil {
    if apiErr, ok := err.(*perception.Error); ok {
        fmt.Printf("Perception API Error %d\n", apiErr.StatusCode)
    } else {
        fmt.Printf("Client Error: %v\n", err)
    }
}
```

## Data Types

### Neural Package Types

- **neural.Space**: Neural space resource with configuration, type, and current state
- **neural.Processor**: Neural processor resource with type-specific configuration
- **neural.CreateSpaceRequest**: For creating new spaces
- **neural.UpdateSpaceRequest**: For updating existing spaces
- **neural.CreateProcessorRequest**: For creating new processors
- **neural.UpdateProcessorRequest**: For updating existing processors
- **neural.SpaceRequestData**: Space data in create requests
- **neural.UpdateSpaceData**: Space data in update requests
- **neural.ProcessorRequestData**: Processor data in create requests
- **neural.UpdateProcessorData**: Processor data in update requests
- **neural.SpaceResponse**: API response wrapper for space operations
- **neural.ProcessorResponse**: API response wrapper for processor operations
- **neural.Error**: Neural service specific error type

### Sensory Package Types

- **sensory.Source**: Sensory data source with type and connection details
- **sensory.Model**: Machine learning model with identifier, path, and configurable parameters
- **sensory.Limit**: Resource limits with counts, scale units, current state, and source association
- **sensory.CreateSourceRequest**: For creating new sources
- **sensory.UpdateSourceRequest**: For updating existing sources
- **sensory.CreateModelRequest**: For creating new models
- **sensory.UpdateModelRequest**: For updating existing models
- **sensory.CreateLimitRequest**: For creating new limits
- **sensory.UpdateLimitRequest**: For updating existing limits
- **sensory.Error**: Sensory service specific error type

### Perception Package Types

- **perception.Chain**: Perception chain resource with name, slug, and current state
- **perception.Thought**: Thought resource with module configuration, relation, and index
- **perception.Module**: Module configuration with reference and parameters
- **perception.CreateChainRequest**: For creating new chains
- **perception.UpdateChainRequest**: For updating existing chains
- **perception.CreateThoughtRequest**: For creating new thoughts
- **perception.UpdateThoughtRequest**: For updating existing thoughts
- **perception.ChainRequestData**: Chain data in create requests
- **perception.UpdateChainData**: Chain data in update requests
- **perception.ThoughtRequestData**: Thought data in create requests
- **perception.UpdateThoughtData**: Thought data in update requests
- **perception.ChainResponse**: API response wrapper for chain operations
- **perception.ThoughtResponse**: API response wrapper for thought operations
- **perception.Error**: Perception service specific error type

## Examples

See the [example/main.go](example/main.go) file for a complete working example demonstrating all client features.

## Requirements

- Go 1.23 or later
- Active Tama API credentials

## Dependencies

- [go-resty/resty](https://github.com/go-resty/resty) - HTTP client library

## Testing

Run the test suite:

```bash
go test -v
```

Run integration tests (requires API credentials):

```bash
export TAMA_BASE_URL="https://api.tama.io"
export TAMA_API_KEY="your-api-key"
go test -tags=integration -v
```

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and questions:
- Create an issue on GitHub
- Check the API documentation
- Review the examples in this repository