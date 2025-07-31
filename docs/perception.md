## Perception Service

Access via `client.Perception.*`

### Chain Operations

#### GetChain(id string) (*Chain, error)

Retrieves a specific chain by ID.

**Endpoint:** `GET /provision/perception/chains/:id`

**Parameters:**
- `id` (string): Chain ID (required)

**Returns:**
- `*Chain`: Chain object with ID, SpaceID, Name, Slug, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
chain, err := client.Perception.GetChain("chain-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Chain: %s (%s)", chain.Name, chain.ProvisionState)
```

#### CreateChain(spaceID string, req CreateChainRequest) (*Chain, error)

Creates a new chain within a space.

**Endpoint:** `POST /provision/perception/spaces/:space_id/chains`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `req` (CreateChainRequest): Chain creation request (required)

```go
type CreateChainRequest struct {
    Chain ChainRequestData `json:"chain"`
}

type ChainRequestData struct {
    Name string `json:"name"`
}
```

**Returns:**
- `*Chain`: Created chain object
- `error`: Error if request fails

**Example:**
```go
chain, err := client.Perception.CreateChain("space-123", perception.CreateChainRequest{
    Chain: perception.ChainRequestData{
        Name: "Processing Chain",
    },
})
```

#### UpdateChain(id string, req UpdateChainRequest) (*Chain, error)

Updates an existing chain using PATCH.

**Endpoint:** `PATCH /provision/perception/chains/:id`

```go
type UpdateChainRequest struct {
    Chain UpdateChainData `json:"chain"`
}

type UpdateChainData struct {
    Name string `json:"name,omitempty"`
}
```

#### ReplaceChain(id string, req UpdateChainRequest) (*Chain, error)

Replaces an existing chain using PUT.

**Endpoint:** `PUT /provision/perception/chains/:id`

#### DeleteChain(id string) error

Deletes a chain by ID.

**Endpoint:** `DELETE /provision/perception/chains/:id`

### Thought Operations

#### GetThought(id string) (*Thought, error)

Retrieves a specific thought by ID.

**Endpoint:** `GET /provision/perception/thoughts/:id`

**Parameters:**
- `id` (string): Thought ID (required)

**Returns:**
- `*Thought`: Thought object with ID, ChainID, Module, Relation, and Index
- `error`: Error if request fails

**Example:**
```go
thought, err := client.Perception.GetThought("thought-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Thought: %s (%s)", thought.Relation, thought.ProvisionState)
```

#### CreateThought(chainID string, req CreateThoughtRequest) (*Thought, error)

Creates a new thought within a chain.

**Endpoint:** `POST /provision/perception/chains/:chain_id/thoughts`

**Parameters:**
- `chainID` (string): Chain ID (required)
- `req` (CreateThoughtRequest): Thought creation request (required)

```go
type CreateThoughtRequest struct {
    Thought ThoughtRequestData `json:"thought"`
}

type ThoughtRequestData struct {
    Relation      string `json:"relation"`
    OutputClassID string `json:"output_class_id,omitempty"`
    Module        Module `json:"module"`
}

type Module struct {
    Reference  string         `json:"reference"`
    Parameters map[string]any `json:"parameters"`
}
```

**Returns:**
- `*Thought`: Created thought object
- `error`: Error if request fails

**Example:**
```go
thought, err := client.Perception.CreateThought("chain-123", perception.CreateThoughtRequest{
    Thought: perception.ThoughtRequestData{
        Relation:      "description",
        OutputClassID: "class-123",
        Module: perception.Module{
            Reference: "tama/agentic/generate",
            Parameters: map[string]any{
                "temperature": 0.7,
                "max_tokens":  150,
            },
        },
    },
})
```

#### UpdateThought(id string, req UpdateThoughtRequest) (*Thought, error)

Updates an existing thought using PATCH.

**Endpoint:** `PATCH /provision/perception/thoughts/:id`

```go
type UpdateThoughtRequest struct {
    Thought UpdateThoughtData `json:"thought"`
}

type UpdateThoughtData struct {
    Relation      string `json:"relation,omitempty"`
    OutputClassID string `json:"output_class_id,omitempty"`
    Module        Module `json:"module,omitempty"`
}
```

#### DeleteThought(id string) error

Deletes a thought by ID.

**Endpoint:** `DELETE /provision/perception/thoughts/:id`

### Path Operations

#### GetPath(id string) (*Path, error)

Retrieves a specific path by ID.

**Endpoint:** `GET /provision/perception/paths/:id`

**Parameters:**
- `id` (string): Path ID (required)

**Returns:**
- `*Path`: Path object with ID, ThoughtID, TargetClassID, Parameters, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
path, err := client.Perception.GetPath("path-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Path: %s (%s)", path.TargetClassID, path.ProvisionState)
```

#### CreatePath(thoughtID string, req CreatePathRequest) (*Path, error)

Creates a new path within a thought.

**Endpoint:** `POST /provision/perception/thoughts/:thought_id/paths`

**Parameters:**
- `thoughtID` (string): Thought ID (required)
- `req` (CreatePathRequest): Path creation request (required)

```go
type CreatePathRequest struct {
    Path PathRequestData `json:"path"`
}

type PathRequestData struct {
    TargetClassID string         `json:"target_class_id"`
    Parameters    map[string]any `json:"parameters,omitempty"`
}
```

**Returns:**
- `*Path`: Created path object
- `error`: Error if request fails

**Example:**
```go
path, err := client.Perception.CreatePath("thought-123", perception.CreatePathRequest{
    Path: perception.PathRequestData{
        TargetClassID: "class-456",
        Parameters: map[string]any{
            "threshold":    0.8,
            "max_results":  10,
            "output_format": "json",
        },
    },
})
```

#### UpdatePath(id string, req UpdatePathRequest) (*Path, error)

Updates an existing path using PATCH.

**Endpoint:** `PATCH /provision/perception/paths/:id`

```go
type UpdatePathRequest struct {
    Path UpdatePathData `json:"path"`
}

type UpdatePathData struct {
    TargetClassID string         `json:"target_class_id,omitempty"`
    Parameters    map[string]any `json:"parameters,omitempty"`
}
```

#### ReplacePath(id string, req UpdatePathRequest) (*Path, error)

Replaces an existing path using PUT.

**Endpoint:** `PUT /provision/perception/paths/:id`

#### DeletePath(id string) error

Deletes a path by ID.

**Endpoint:** `DELETE /provision/perception/paths/:id`

### Context Operations

#### GetContext(id string) (*Context, error)

Retrieves a specific context by ID.

**Endpoint:** `GET /provision/perception/contexts/:id`

**Parameters:**
- `id` (string): Context ID (required)

**Returns:**
- `*Context`: Context object with ID, ThoughtID, PromptID, Layer, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
context, err := client.Perception.GetContext("context-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Context: %s layer %d (%s)", context.PromptID, context.Layer, context.ProvisionState)
```

#### CreateContext(thoughtID string, req CreateContextRequest) (*Context, error)

Creates a new context within a thought.

**Endpoint:** `POST /provision/perception/thoughts/:thought_id/contexts`

**Parameters:**
- `thoughtID` (string): Thought ID (required)
- `req` (CreateContextRequest): Context creation request (required)

```go
type CreateContextRequest struct {
    Context ContextRequestData `json:"context"`
}

type ContextRequestData struct {
    PromptID string `json:"prompt_id"`
    Layer    int    `json:"layer"`
}
```

**Returns:**
- `*Context`: Created context object
- `error`: Error if request fails

**Example:**
```go
context, err := client.Perception.CreateContext("thought-123", perception.CreateContextRequest{
    Context: perception.ContextRequestData{
        PromptID: "prompt-456",
        Layer:    2,
    },
})
```

#### UpdateContext(id string, req UpdateContextRequest) (*Context, error)

Updates an existing context using PATCH.

**Endpoint:** `PATCH /provision/perception/contexts/:id`

```go
type UpdateContextRequest struct {
    Context UpdateContextData `json:"context"`
}

type UpdateContextData struct {
    PromptID string `json:"prompt_id,omitempty"`
    Layer    int    `json:"layer,omitempty"`
}
```

#### ReplaceContext(id string, req UpdateContextRequest) (*Context, error)

Replaces an existing context using PUT.

**Endpoint:** `PUT /provision/perception/contexts/:id`

#### DeleteContext(id string) error

Deletes a context by ID.

**Endpoint:** `DELETE /provision/perception/contexts/:id`
