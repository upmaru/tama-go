# Tools Service

The **Tools** service provides three main categories of operations:
- **Input** – Manage tool inputs (type, class corpus ID, CRUD).
- **Initializer** – Manage tool initializers (reference, index, parameters, CRUD).
- **Output** – Manage tool outputs (class corpus ID, CRUD).

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
- [Output Operations](#output-operations)

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

## Output Operations

### GetOutput(id string) (*Output, error)

Retrieves a specific tools output by ID.

**Endpoint:** `GET /provision/tools/outputs/:id`

**Parameters:**
- `id` (string): Output ID (required)

**Returns:**
- `*Output`: Output object with ID, ThoughtToolID, ClassCorpusID, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
output, err := client.Tools.GetOutput("output-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Output: %+v\n", output)
```

### CreateOutput(thoughtToolID string, req CreateOutputRequest) (*Output, error)

Creates a new output under a specific thought tool.

**Endpoint:** `POST /provision/tools/:thought_tool_id/outputs`

**Parameters:**
- `thoughtToolID` (string): Parent thought tool ID (required)
- `req` (CreateOutputRequest): Output creation request (required)
  - `Output` (OutputRequestData): Output data (required)
    - `ClassCorpusID` (string): Class corpus ID (required)

**Returns:**
- `*Output`: Created output object
- `error`: Error if request fails

**Request Structure:**
```go
type CreateOutputRequest struct {
    Output OutputRequestData `json:"output"`
}

type OutputRequestData struct {
    ClassCorpusID string `json:"class_corpus_id"`
}
```

**Example:**
```go
req := tools.CreateOutputRequest{
    Output: tools.OutputRequestData{
        ClassCorpusID: "corpus-789",
    },
}
created, err := client.Tools.CreateOutput("tool-123", req)
```

### UpdateOutput(id string, req UpdateOutputRequest) (*Output, error)

Updates an existing output using PATCH (partial update).

**Endpoint:** `PATCH /provision/tools/outputs/:id`

**Parameters:**
- `id` (string): Output ID (required)
- `req` (UpdateOutputRequest): Output update request (required)
  - `Output` (UpdateOutputData): Updatable fields
    - `ClassCorpusID` (string, optional)

**Returns:**
- `*Output`: Updated output object
- `error`: Error if request fails

**Request Structure:**
```go
type UpdateOutputRequest struct {
    Output UpdateOutputData `json:"output"`
}

type UpdateOutputData struct {
    ClassCorpusID string `json:"class_corpus_id,omitempty"`
}
```

**Example:**
```go
update := tools.UpdateOutputRequest{
    Output: tools.UpdateOutputData{
        ClassCorpusID: "corpus-updated",
    },
}
updated, err := client.Tools.UpdateOutput("output-123", update)
```

### ReplaceOutput(id string, req UpdateOutputRequest) (*Output, error)

Replaces an existing output using PUT (full replacement semantics on the server side).

**Endpoint:** `PUT /provision/tools/outputs/:id`

**Parameters:**
- `id` (string): Output ID (required)
- `req` (UpdateOutputRequest): Replacement request (required)

**Returns:**
- `*Output`: Replaced output object
- `error`: Error if request fails

**Example:**
```go
replace := tools.UpdateOutputRequest{
    Output: tools.UpdateOutputData{
        ClassCorpusID: "corpus-replaced",
    },
}
replaced, err := client.Tools.ReplaceOutput("output-123", replace)
```

### DeleteOutput(id string) error

Deletes an output by ID.

**Endpoint:** `DELETE /provision/tools/outputs/:id`

**Parameters:**
- `id` (string): Output ID (required)

**Returns:**
- `error`: Error if request fails

**Example:**
```go
if err := client.Tools.DeleteOutput("output-123"); err != nil {
    log.Fatal(err)
}
```

## Output Structure

### Output

The `Output` struct represents a tools output with the following fields:

```go
type Output struct {
    ID             string `json:"id,omitempty"`
    ThoughtToolID  string `json:"thought_tool_id,omitempty"`
    ClassCorpusID  string `json:"class_corpus_id"`
    ProvisionState string `json:"provision_state"`
}
```

### OutputResponse

Wraps an Output in API responses:

```go
type OutputResponse struct {
    Data Output `json:"data"`
}
```
