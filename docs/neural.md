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

#### GetProcessor(spaceID, modelID string) (*Processor, error)

Retrieves a specific processor by space ID and model ID.

**Endpoint:** `GET /provision/neural/spaces/:space_id/models/:model_id/processor`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `modelID` (string): Model ID (required)

**Returns:**
- `*Processor`: Processor data
- `error`: Error if request fails

**Example:**
```go
processor, err := client.Neural.GetProcessor("space-123", "model-456")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Processor: %+v\n", processor)
```

#### CreateProcessor(spaceID, modelID string, req CreateProcessorRequest) (*Processor, error)

Creates a new processor for a specific space and model.

**Endpoint:** `POST /provision/neural/spaces/:space_id/models/:model_id/processor`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `modelID` (string): Model ID (required)
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
