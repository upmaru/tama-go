# Tools Service

The **Tools** service provides two main categories of operations:
- **Input** – Manage tool inputs (type, class corpus ID, CRUD).
- **Initializer** – Manage tool initializers (reference, index, parameters, CRUD).

These are exposed via the `Client.Tools` field and can be accessed as follows:

```go
client := tama.NewClient(config)
tools := client.Tools

// Input operations
input, err := tools.Input.Create(ctx, newInput)

// Initializer operations
initializer, err := tools.Initializer.Create(ctx, newInitializer)
```

## Table of Contents

- [Overview](#overview)
- [Input Operations](#input-operations)
  - [Create](#create-input)
  - [Get](#get-input)
  - [List](#list-inputs)
  - [Update](#update-input)
  - [Delete](#delete-input)
- [Initializer Operations](#initializer-operations)
  - [Create](#create-initializer)
  - [Get](#get-initializer)
  - [List](#list-initializers)
  - [Update](#update-initializer)
  - [Delete](#delete-initializer)

## Overview

The Tools service allows clients to store and retrieve metadata about tools that can be used by other services such as **Perception**.  
- An **Input** represents a set of parameters that a tool expects.
- An **Initializer** contains the logic (reference, index, parameters) required to instantiate or invoke a tool.

Both resources support full CRUD operations with pagination for listing endpoints.

## Input Operations

### Create
```go
func (s *InputService) Create(ctx context.Context, input *tools.Input) (*tools.Input, error)
```
Creates a new tool input. Required fields: `Type`, `ClassCorpusID`.

### Get
```go
func (s *InputService) Get(ctx context.Context, id string) (*tools.Input, error)
```
Retrieves an existing input by its unique identifier.

### List
```go
func (s *InputService) List(ctx context.Context, opts *tools.ListOptions) ([]*tools.Input, error)
```
Returns a paginated list of inputs. Supports filtering by `Type` or `ClassCorpusID`.

### Update
```go
func (s *InputService) Update(ctx context.Context, id string, input *tools.InputUpdate) (*tools.Input, error)
```
Updates mutable fields of an existing input.

### Delete
```go
func (s *InputService) Delete(ctx context.Context, id string) error
```
Deletes the specified input.

## Initializer Operations

### Create
```go
func (s *InitializerService) Create(ctx context.Context, init *tools.Initializer) (*tools.Initializer, error)
```
Creates a new tool initializer. Required fields: `Reference`, `Index`, `Parameters`.

### Get
```go
func (s *InitializerService) Get(ctx context.Context, id string) (*tools.Initializer, error)
```
Retrieves an existing initializer by its unique identifier.

### List
```go
func (s *InitializerService) List(ctx context.Context, opts *tools.ListOptions) ([]*tools.Initializer, error)
```
Returns a paginated list of initializers. Supports filtering by `Reference`.

### Update
```go
func (s *InitializerService) Update(ctx context.Context, id string, init *tools.InitializerUpdate) (*tools.Initializer, error)
```
Updates mutable fields of an existing initializer.

### Delete
```go
func (s *InitializerService) Delete(ctx context.Context, id string) error
```
Deletes the specified initializer.

## Usage Example

```go
import (
    "context"
    "github.com/upmaru/tama-go"
)

func main() {
    cfg := tama.Config{BaseURL: "https://api.tama.io", APIKey: "YOUR_KEY"}
    client := tama.NewClient(cfg)
    ctx := context.Background()

    // Create an input
    in, _ := client.Tools.Input.Create(ctx, &tama.tools.Input{
        Type:          "text",
        ClassCorpusID: "corpus-123",
    })

    // Create an initializer referencing the input
    init, _ := client.Tools.Initializer.Create(ctx, &tama.tools.Initializer{
        Reference:  in.ID,
        Index:      0,
        Parameters: []string{"param1", "param2"},
    })
}
```

> **Tip:**  
> Use the listing endpoints with pagination options (`Limit`, `Offset`) to efficiently navigate large sets of inputs or initializers.

--- 

For more detailed information on each service, refer to the generated SDK documentation in the `tools` package source.