# Perception Service

Access via `client.Perception.*`

## Table of Contents

- [Delegation Operations](#delegation-operations)
- [Processor Operations](#processor-operations)
- [Chain Operations](#chain-operations)
- [Thought Operations](#thought-operations)
- [Path Operations](#path-operations)
- [Context Operations](#context-operations)
- [Tool Operations](#tool-operations)
- [Initializer Operations](#initializer-operations)
- [Activation Operations](#activation-operations)

## Delegation Operations

### GetDelegation(thoughtID string) (*Delegation, error)

Retrieves a specific delegation by thought ID.

**Endpoint:** `GET /provision/perception/thoughts/:thought_id/delegation`

**Parameters:**
- `thoughtID` (string): Thought ID (required)

**Returns:**
- `*Delegation`: Delegation object with ID, ThoughtID, TargetThoughtID, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
delegation, err := client.Perception.GetDelegation("thought-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Delegation: %s -> %s (%s)", delegation.ID, delegation.TargetThoughtID, delegation.ProvisionState)
```

### CreateDelegation(thoughtID string, req CreateDelegationRequest) (*Delegation, error)

Creates a new delegation within a thought.

**Endpoint:** `POST /provision/perception/thoughts/:thought_id/delegation`

**Parameters:**
- `thoughtID` (string): Thought ID (required)
- `req` (CreateDelegationRequest): Delegation creation request (required)

```go
type CreateDelegationRequest struct {
    Delegation DelegationRequestData `json:"delegation"`
}

type DelegationRequestData struct {
    TargetThoughtID string `json:"target_thought_id"`
}
```

**Returns:**
- `*Delegation`: Created delegation object
- `error`: Error if request fails

**Example:**
```go
delegation, err := client.Perception.CreateDelegation("thought-123", perception.CreateDelegationRequest{
    Delegation: perception.DelegationRequestData{
        TargetThoughtID: "target-thought-123",
    },
})
```

### UpdateDelegation(thoughtID string, req UpdateDelegationRequest) (*Delegation, error)

Updates an existing delegation using PATCH.

**Endpoint:** `PATCH /provision/perception/thoughts/:thought_id/delegation`

```go
type UpdateDelegationRequest struct {
    Delegation UpdateDelegationData `json:"delegation"`
}

type UpdateDelegationData struct {
    TargetThoughtID string `json:"target_thought_id,omitempty"`
}
```

### ReplaceDelegation(thoughtID string, req UpdateDelegationRequest) (*Delegation, error)

Replaces an existing delegation using PUT.

**Endpoint:** `PUT /provision/perception/thoughts/:thought_id/delegation`

### DeleteDelegation(thoughtID string) error

Deletes a delegation by thought ID.

**Endpoint:** `DELETE /provision/perception/thoughts/:thought_id/delegation`

## Processor Operations

### GetProcessor(thoughtID, processorType string) (*Processor, error)

Retrieves a specific processor by thought ID and type.

**Endpoint:** `GET /provision/perception/thoughts/:thought_id/types/:type/processor`

**Parameters:**
- `thoughtID` (string): Thought ID (required)
- `processorType` (string): Processor type (required)

**Returns:**
- `*Processor`: Processor object with ID, ThoughtID, ModelID, Configuration, ProvisionState, and Type
- `error`: Error if request fails

**Example:**
```go
processor, err := client.Perception.GetProcessor("thought-123", "text-generation")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Processor: %s (%s)", processor.ModelID, processor.ProvisionState)
```

### CreateProcessor(thoughtID, processorType string, req CreateProcessorRequest) (*Processor, error)

Creates a new processor within a thought.

**Endpoint:** `POST /provision/perception/thoughts/:thought_id/types/:type/processor`

**Parameters:**
- `thoughtID` (string): Thought ID (required)
- `processorType` (string): Processor type (required)
- `req` (CreateProcessorRequest): Processor creation request (required)

```go
type CreateProcessorRequest struct {
    Processor ProcessorRequestData `json:"processor"`
}

type ProcessorRequestData struct {
    ModelID       string         `json:"model_id"`
    Configuration map[string]any `json:"configuration"`
}
```

**Returns:**
- `*Processor`: Created processor object
- `error`: Error if request fails

**Example:**
```go
processor, err := client.Perception.CreateProcessor("thought-123", "text-generation", perception.CreateProcessorRequest{
    Processor: perception.ProcessorRequestData{
        ModelID: "gpt-4",
        Configuration: map[string]any{
            "temperature":   0.7,
            "max_tokens":    150,
            "top_p":         0.9,
            "frequency_penalty": 0.0,
        },
    },
})
```

### UpdateProcessor(thoughtID, processorType string, req UpdateProcessorRequest) (*Processor, error)

Updates an existing processor using PATCH.

**Endpoint:** `PATCH /provision/perception/thoughts/:thought_id/types/:type/processor`

```go
type UpdateProcessorRequest struct {
    Processor UpdateProcessorData `json:"processor"`
}

type UpdateProcessorData struct {
    ModelID       string         `json:"model_id,omitempty"`
    Configuration map[string]any `json:"configuration,omitempty"`
}
```

### ReplaceProcessor(thoughtID, processorType string, req UpdateProcessorRequest) (*Processor, error)

Replaces an existing processor using PUT.

**Endpoint:** `PUT /provision/perception/thoughts/:thought_id/types/:type/processor`

### DeleteProcessor(thoughtID, processorType string) error

Deletes a processor by thought ID and type.

**Endpoint:** `DELETE /provision/perception/thoughts/:thought_id/types/:type/processor`

## Chain Operations

### GetChain(id string) (*Chain, error)

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

### CreateChain(spaceID string, req CreateChainRequest) (*Chain, error)

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

### UpdateChain(id string, req UpdateChainRequest) (*Chain, error)

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

### ReplaceChain(id string, req UpdateChainRequest) (*Chain, error)

Replaces an existing chain using PUT.

**Endpoint:** `PUT /provision/perception/chains/:id`

### DeleteChain(id string) error

Deletes a chain by ID.

**Endpoint:** `DELETE /provision/perpection/chains/:id`

## Thought Operations

### GetThought(id string) (*Thought, error)

Retrieves a specific thought by ID.

**Endpoint:** `GET /provision/perception/thoughts/:id`

**Parameters:**
- `id` (string): Thought ID (required)

**Returns:**
- `*Thought`: Thought object with ID, ChainID, Module, Delegation, Relation, and Index
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

### CreateThought(chainID string, req CreateThoughtRequest) (*Thought, error)

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
    Relation      string      `json:"relation"`
    OutputClassID string      `json:"output_class_id,omitempty"`
    Module        *Module     `json:"module,omitempty"`
    Delegation    *Delegation `json:"delegation,omitempty"`
}

type Module struct {
    Reference  string         `json:"reference"`
    Parameters map[string]any `json:"parameters"`
}

type Delegation struct {
    TargetThoughtID string `json:"target_thought_id"`
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
        Module: &perception.Module{
            Reference: "tama/agentic/generate",
            Parameters: map[string]any{
                "temperature": 0.7,
                "max_tokens":  150,
            },
        },
    },
})
```

### UpdateThought(id string, req UpdateThoughtRequest) (*Thought, error)

Updates an existing thought using PATCH.

**Endpoint:** `PATCH /provision/perception/thoughts/:id`

```go
type UpdateThoughtRequest struct {
    Thought UpdateThoughtData `json:"thought"`
}

type UpdateThoughtData struct {
    Relation      string      `json:"relation,omitempty"`
    OutputClassID string      `json:"output_class_id,omitempty"`
    Module        *Module     `json:"module,omitempty"`
    Delegation    *Delegation `json:"delegation,omitempty"`
}
```

### DeleteThought(id string) error

Deletes a thought by ID.

**Endpoint:** `DELETE /provision/perception/thoughts/:id`

## Path Operations

### GetPath(id string) (*Path, error)

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

### CreatePath(thoughtID string, req CreatePathRequest) (*Path, error)

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

### UpdatePath(id string, req UpdatePathRequest) (*Path, error)

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

### ReplacePath(id string, req UpdatePathRequest) (*Path, error)

Replaces an existing path using PUT.

**Endpoint:** `PUT /provision/perception/paths/:id`

### DeletePath(id string) error

Deletes a path by ID.

**Endpoint:** `DELETE /provision/perception/paths/:id`

## Context Operations

### GetContext(id string) (*Context, error)

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

### CreateContext(thoughtID string, req CreateContextRequest) (*Context, error)

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

### UpdateContext(id string, req UpdateContextRequest) (*Context, error)

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

### ReplaceContext(id string, req UpdateContextRequest) (*Context, error)

Replaces an existing context using PUT.

**Endpoint:** `PUT /provision/perception/contexts/:id`

### DeleteContext(id string) error

Deletes a context by ID.

**Endpoint:** `DELETE /provision/perception/contexts/:id`

## Tool Operations

### GetTool(id string) (*Tool, error)

Retrieves a specific tool by ID.

**Endpoint:** `GET /provision/perception/tools/:id`

**Parameters:**
- `id` (string): Tool ID (required)

**Returns:**
- `*Tool`: Tool object with ID, ThoughtID, ActionID, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
tool, err := client.Perception.GetTool("tool-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Tool: %s (%s)", tool.ActionID, tool.ProvisionState)
```

### CreateTool(thoughtID string, req CreateToolRequest) (*Tool, error)

Creates a new tool within a thought.

**Endpoint:** `POST /provision/perception/thoughts/:thought_id/tools`

**Parameters:**
- `thoughtID` (string): Thought ID (required)
- `req` (CreateToolRequest): Tool creation request (required)

```go
type CreateToolRequest struct {
    Tool CreateToolData `json:"tool"`
}

type CreateToolData struct {
    ActionID string `json:"action_id"`
}
```

**Returns:**
- `*Tool`: Created tool object
- `error`: Error if request fails

**Example:**
```go
tool, err := client.Perception.CreateTool("thought-123", perception.CreateToolRequest{
    Tool: perception.CreateToolData{
        ActionID: "action-456",
    },
})
```

### UpdateTool(id string, req UpdateToolRequest) (*Tool, error)

Updates an existing tool using PATCH.

**Endpoint:** `PATCH /provision/perception/tools/:id`

```go
type UpdateToolRequest struct {
    Tool UpdateToolData `json:"tool"`
}

type UpdateToolData struct {
    ActionID string `json:"action_id,omitempty"`
}
```

### ReplaceTool(id string, req UpdateToolRequest) (*Tool, error)

Replaces an existing tool using PUT.

**Endpoint:** `PUT /provision/perception/tools/:id`

### DeleteTool(id string) error

Deletes a tool by ID.

**Endpoint:** `DELETE /provision/perception/tools/:id`

## Initializer Operations

### GetInitializer(id string) (*Initializer, error)

Retrieves a specific initializer by ID.

**Endpoint:** `GET /provision/perception/initializers/:id`

**Parameters:**
- `id` (string): Initializer ID (required)

**Returns:**
- `*Initializer`: Initializer object with ID, Parameters, Index, ProvisionState, ThoughtID, ClassID, and Reference
- `error`: Error if request fails

**Example:**
```go
initializer, err := client.Perception.GetInitializer("initializer-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Initializer: %s (%s)", initializer.ClassID, initializer.ProvisionState)
```

### CreateInitializer(thoughtID string, req CreateInitializerRequest) (*Initializer, error)

Creates a new initializer within a thought.

**Endpoint:** `POST /provision/perception/thoughts/:thought_id/initializers`

**Parameters:**
- `thoughtID` (string): Thought ID (required)
- `req` (CreateInitializerRequest): Initializer creation request (required)

```go
type CreateInitializerRequest struct {
    Initializer InitializerRequestData `json:"initializer"`
}

type InitializerRequestData struct {
    Parameters map[string]any `json:"parameters,omitempty"`
    Index      *int           `json:"index,omitempty"`
    ClassID    string         `json:"class_id"`
    Reference  string         `json:"reference"`
}
```

**Returns:**
- `*Initializer`: Created initializer object
- `error`: Error if request fails

**Example:**
```go
initializer, err := client.Perception.CreateInitializer("thought-123", perception.CreateInitializerRequest{
    Initializer: perception.InitializerRequestData{
        ClassID:   "class-456",
        Reference: "reference-789",
        Parameters: map[string]any{
            "param1": "value1",
            "param2": 42,
        },
    },
})
```

### UpdateInitializer(id string, req UpdateInitializerRequest) (*Initializer, error)

Updates an existing initializer using PATCH.

**Endpoint:** `PATCH /provision/perception/initializers/:id`

```go
type UpdateInitializerRequest struct {
    Initializer UpdateInitializerData `json:"initializer"`
}

type UpdateInitializerData struct {
    Parameters map[string]any `json:"parameters,omitempty"`
    Index      *int           `json:"index,omitempty"`
    ClassID    string         `json:"class_id,omitempty"`
    Reference  string         `json:"reference,omitempty"`
}
```

### ReplaceInitializer(id string, req UpdateInitializerRequest) (*Initializer, error)

Replaces an existing initializer using PUT.

**Endpoint:** `PUT /provision/perception/initializers/:id`

### DeleteInitializer(id string) error

Deletes an initializer by ID.

**Endpoint:** `DELETE /provision/perception/initializers/:id`

## Activation Operations

### GetActivation(id string) (*Activation, error)

Retrieves a specific activation by ID.

**Endpoint:** `GET /provision/perception/activations/:id`

**Parameters:**
- `id` (string): Activation ID (required)

**Returns:**
- `*Activation`: Activation object with ID, ThoughtPathID, ChainID, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
activation, err := client.Perception.GetActivation("activation-123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Printf("Activation: %s (%s)", activation.ChainID, activation.ProvisionState)
```

### CreateActivation(pathID string, req CreateActivationRequest) (*Activation, error)

Creates a new activation within a path.

**Endpoint:** `POST /provision/perception/paths/:path_id/activations`

**Parameters:**
- `pathID` (string): Path ID (required)
- `req` (CreateActivationRequest): Activation creation request (required)

```go
type CreateActivationRequest struct {
    Activation ActivationRequestData `json:"activation"`
}

type ActivationRequestData struct {
    ChainID string `json:"chain_id"`
}
```

**Returns:**
- `*Activation`: Created activation object
- `error`: Error if request fails

**Example:**
```go
activation, err := client.Perception.CreateActivation("path-123", perception.CreateActivationRequest{
    Activation: perception.ActivationRequestData{
        ChainID: "chain-456",
    },
})
```

### UpdateActivation(id string, req UpdateActivationRequest) (*Activation, error)

Updates an existing activation using PATCH.

**Endpoint:** `PATCH /provision/perception/activations/:id`

```go
type UpdateActivationRequest struct {
    Activation UpdateActivationData `json:"activation"`
}

type UpdateActivationData struct {
    ChainID string `json:"chain_id,omitempty"`
}
```

### ReplaceActivation(id string, req UpdateActivationRequest) (*Activation, error)

Replaces an existing activation using PUT.

**Endpoint:** `PUT /provision/perception/activations/:id`

### DeleteActivation(id string) error

Deletes an activation by ID.

**Endpoint:** `DELETE /provision/perception/activations/:id`
