# Neural Service

Access via `client.Neural.*`

### GetSpace(id string) (*Space, error)

Retrieves a specific neural space by ID.

**Endpoint:** `GET /provision/neural/spaces/:id`

**Parameters:**
- `id` (string): Space ID (required)

**Returns:**
- `*Space`: Space object with ID, Name, Slug, Type, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
space, err := client.Neural.GetSpace("space-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Space: %+v\n", space)
```

### CreateSpace(req CreateSpaceRequest) (*Space, error)

Creates a new neural space.

**Endpoint:** `POST /provision/neural/spaces`

**Parameters:**
- `req` (CreateSpaceRequest): Space creation request
  - `Space` (SpaceRequest): Space data (required)
    - `Name` (string): Space name (required)
    - `Type` (string): Space type (required)

**Returns:**
- `*Space`: Created space object with ID, Name, Slug, Type, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
req := tama.CreateSpaceRequest{
    Space: tama.SpaceRequest{
        Name: "My Neural Space",
        Type: "root",
    },
}
space, err := client.Neural.CreateSpace(req)
```

### UpdateSpace(id string, req UpdateSpaceRequest) (*Space, error)

Updates an existing space using PATCH (partial update).

**Endpoint:** `PATCH /provision/neural/spaces/:id`

**Parameters:**
- `id` (string): Space ID (required)
- `req` (UpdateSpaceRequest): Update request
  - `Space` (UpdateSpaceData): Space update data (required)
    - `Name` (string): New space name (optional)
    - `Type` (string): New space type (optional)

**Returns:**
- `*Space`: Updated space object with all fields including server-managed ProvisionState
- `error`: Error if request fails

### ReplaceSpace(id string, req UpdateSpaceRequest) (*Space, error)

Replaces an existing space using PUT (full replacement).

**Endpoint:** `PUT /provision/neural/spaces/:id`

**Parameters:**
- `id` (string): Space ID (required)
- `req` (UpdateSpaceRequest): Replacement request

**Returns:**
- `*Space`: Updated space object with all fields including server-managed ProvisionState
- `error`: Error if request fails

### DeleteSpace(id string) error

Deletes a space by ID.

**Endpoint:** `DELETE /provision/neural/spaces/:id`

**Parameters:**
- `id` (string): Space ID (required)

**Returns:**
- `error`: Error if request fails

### Processor Operations

#### GetProcessor(spaceID, processorType string) (*Processor, error)

Retrieves a specific processor by space ID and processor type.

**Endpoint:** `GET /provision/neural/spaces/:space_id/types/:type/processor`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `processorType` (string): Processor type (required)

**Returns:**
- `*Processor`: Processor data
- `error`: Error if request fails

**Example:**
```go
processor, err := client.Neural.GetProcessor("space-123", "chat")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Processor: %+v\n", processor)
```

#### CreateProcessor(spaceID, processorType string, req CreateProcessorRequest) (*Processor, error)

Creates a new processor for a specific space and type.

**Endpoint:** `POST /provision/neural/spaces/:space_id/types/:type/processor`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `processorType` (string): Processor type (required)
- `req` (CreateProcessorRequest): Processor creation data (required)

**Request Structure:**
```go
type CreateProcessorRequest struct {
    Processor ProcessorRequestData `json:"processor"`
}

type ProcessorRequestData struct {
    ModelID       string         `json:"model_id"`      // Model identifier for the processor
    Configuration map[string]any `json:"configuration"` // Type-specific configuration
}
```

**Returns:**
- `*Processor`: Created processor object
- `error`: Error if request fails

**Example:**
```go
req := neural.CreateProcessorRequest{
    Processor: neural.ProcessorRequestData{
        ModelID: "model-456",
        Configuration: map[string]any{
            "temperature":  0.8,
            "tool_choice": "required",
            "role_mappings": []map[string]any{
                {"from": "user", "to": "human"},
                {"from": "assistant", "to": "ai"},
            },
        },
    },
}

processor, err := client.Neural.CreateProcessor("space-123", "chat", req)
```

#### UpdateProcessor(spaceID, processorType string, req UpdateProcessorRequest) (*Processor, error)

Updates an existing processor using PATCH (partial update).

**Endpoint:** `PATCH /provision/neural/spaces/:space_id/types/:type/processor`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `processorType` (string): Processor type (required)
- `req` (UpdateProcessorRequest): Processor update data (required)

**Request Structure:**
```go
type UpdateProcessorRequest struct {
    Processor UpdateProcessorData `json:"processor"`
}

type UpdateProcessorData struct {
    ModelID       string         `json:"model_id,omitempty"`
    Configuration map[string]any `json:"configuration,omitempty"`
}
```

**Returns:**
- `*Processor`: Updated processor
- `error`: Error if request fails

#### ReplaceProcessor(spaceID, processorType string, req UpdateProcessorRequest) (*Processor, error)

Replaces an existing processor using PUT (full replacement).

**Endpoint:** `PUT /provision/neural/spaces/:space_id/types/:type/processor`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `processorType` (string): Processor type (required)
- `req` (UpdateProcessorRequest): Processor replacement data (required)

**Returns:**
- `*Processor`: Replaced processor
- `error`: Error if request fails

#### DeleteProcessor(spaceID, processorType string) error

Deletes a processor by space ID and processor type.

**Endpoint:** `DELETE /provision/neural/spaces/:space_id/types/:type/processor`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `processorType` (string): Processor type (required)

**Returns:**
- `error`: Error if request fails

### Node Operations

#### GetNode(id string) (*Node, error)

Retrieves a specific neural node by ID.

**Endpoint:** `GET /provision/neural/nodes/:id`

**Parameters:**
- `id` (string): Node ID (required)

**Returns:**
- `*Node`: Node object with relevant node details
- `error`: Error if request fails

**Example:**
```go
node, err := client.Neural.GetNode("node-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Node: %+v\n", node)
```

#### CreateNode(req CreateNodeRequest) (*Node, error)

Creates a new neural node.

**Endpoint:** `POST /provision/neural/nodes`

**Parameters:**
- `req` (CreateNodeRequest): Node creation request

**Returns:**
- `*Node`: Created node object with all fields
- `error`: Error if request fails

#### UpdateNode(id string, req UpdateNodeRequest) (*Node, error)

Updates an existing node using PATCH (partial update).

**Endpoint:** `PATCH /provision/neural/nodes/:id`

**Parameters:**
- `id` (string): Node ID (required)
- `req` (UpdateNodeRequest): Update request

**Returns:**
- `*Node`: Updated node object with all fields
- `error`: Error if request fails

#### ReplaceNode(id string, req UpdateNodeRequest) (*Node, error)

Replaces an existing node using PUT (full replacement).

**Endpoint:** `PUT /provision/neural/nodes/:id`

**Parameters:**
- `id` (string): Node ID (required)
- `req` (UpdateNodeRequest): Replacement request

**Returns:**
- `*Node`: Updated node object with all fields
- `error`: Error if request fails

#### DeleteNode(id string) error

Deletes a node by ID.

**Endpoint:** `DELETE /provision/neural/nodes/:id`

**Parameters:**
- `id` (string): Node ID (required)

**Returns:**
- `error`: Error if request fails
