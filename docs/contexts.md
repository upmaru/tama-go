# Contexts Service

Access via `client.Contexts.*`

## Table of Contents

- [Input Operations](#input-operations)

## Input Operations

### GetInput(id string) (*Input, error)

Retrieves a specific input by ID.

**Endpoint:** `GET /provision/contexts/inputs/:id`

**Parameters:**
- `id` (string): Input ID (required)

**Returns:**
- `*Input`: Input object with ID, Type, ThoughtContextID, ClassCorpusID, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
input, err := client.Contexts.GetInput("input-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Input: %s (%s)", input.Type, input.ProvisionState)
```

### CreateInput(thoughtContextID string, req CreateInputRequest) (*Input, error)

Creates a new contexts input within a thought context.

**Endpoint:** `POST /provision/contexts/:thought_context_id/inputs`

**Parameters:**
- `thoughtContextID` (string): Thought context ID (required)
- `req` (CreateInputRequest): Input creation request (required)

```go
type CreateInputRequest struct {
    Input CreateInputData `json:"input"`
}

type CreateInputData struct {
    Type          string `json:"type"`
    ClassCorpusID string `json:"class_corpus_id"`
}
```

**Returns:**
- `*Input`: Created input object
- `error`: Error if request fails

**Example:**
```go
input, err := client.Contexts.CreateInput("thought-context-123", contexts.CreateInputRequest{
    Input: contexts.CreateInputData{
        Type:          "text",
        ClassCorpusID: "corpus-456",
    },
})
```

### UpdateInput(id string, req UpdateInputRequest) (*Input, error)

Updates an existing input using PATCH.

**Endpoint:** `PATCH /provision/contexts/inputs/:id`

```go
type UpdateInputRequest struct {
    Input UpdateInputData `json:"input"`
}

type UpdateInputData struct {
    Type          string `json:"type,omitempty"`
    ClassCorpusID string `json:"class_corpus_id,omitempty"`
}
```

### ReplaceInput(id string, req UpdateInputRequest) (*Input, error)

Replaces an existing input using PUT.

**Endpoint:** `PUT /provision/contexts/inputs/:id`

### DeleteInput(id string) error

Deletes an input by ID.

**Endpoint:** `DELETE /provision/contexts/inputs/:id`
