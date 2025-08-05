# Tama Go Client API Documentation

This directory contains the API documentation for the Tama Go client library, organized by service and function.

## Service Documentation

### Core Services
- [Client Configuration](client_configuration.md) - How to configure and initialize the Tama client
- [Neural Service](neural.md) - Space, Processor, Class, Corpus, and Bridge operations
- [Memory Service](memory.md) - Prompt operations and memory management
- [Sensory Service](sensory.md) - Source, Model, Limit, Specification, and Identity operations
- [Motor Service](motor.md) - Action operations
- [Perception Service](perception.md) - Chain, Thought, Path, and Context operations

## Quick Start

1. **Installation**: Add the Tama Go client to your project
2. **Configuration**: See [Client Configuration](client_configuration.md) for setup instructions
3. **Service Usage**: Choose the service you need from the list above
4. **Error Handling**: All operations return errors that should be properly handled

## Example Usage

```go
package main

import (
    "log"
    "time"
    "github.com/upmaru/tama-go"
)

func main() {
    // Configure the client
    config := tama.Config{
        BaseURL: "https://api.tama.io",
        APIKey:  "your-api-key",
        Timeout: 30 * time.Second,
    }
    client := tama.NewClient(config)

    // Use any service
    space, err := client.Neural.GetSpace("space-123")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Space: %+v", space)
}
```

## Services Overview

- **Neural Service**: Manages neural spaces, processors, classes, corpora, and bridges
- **Memory Service**: Handles prompt operations and memory management
- **Sensory Service**: Manages sources, models, limits, specifications, and identities
- **Motor Service**: Handles action operations
- **Perception Service**: Handles chains, thoughts, paths, and contexts

Each service provides comprehensive CRUD operations with proper error handling and validation.
