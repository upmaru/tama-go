# Neural Service

Access via `client.Neural.*`

## Table of Contents

- [Space](#space)
- [Processor](#processor)
- [Node](#node)
- [Bridge](#bridge)
- [Class](#class)
- [Corpus](#corpus)

## Space

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

## Processor Operations

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

## Node Operations

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

## Bridge Operations

#### GetBridge(id string) (*Bridge, error)

Retrieves a specific bridge by ID.

**Endpoint:** `GET /provision/neural/bridges/:id`

**Parameters:**
- `id` (string): Bridge ID (required)

**Returns:**
- `*Bridge`: Bridge object with ID, SpaceID, TargetSpaceID, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
bridge, err := client.Neural.GetBridge("bridge-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Bridge: %+v\n", bridge)
```

#### CreateBridge(spaceID string, req CreateBridgeRequest) (*Bridge, error)

Creates a new bridge within a space.

**Endpoint:** `POST /provision/neural/spaces/:space_id/bridges`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `req` (CreateBridgeRequest): Bridge creation request
  - `Bridge` (BridgeRequestData): Bridge data (required)
    - `TargetSpaceID` (string): Target space ID (required)

**Returns:**
- `*Bridge`: Created bridge object with ID, SpaceID, TargetSpaceID, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
req := tama.CreateBridgeRequest{
    Bridge: tama.BridgeRequestData{
        TargetSpaceID: "target-space-456",
    },
}
bridge, err := client.Neural.CreateBridge("space-123", req)
```

#### UpdateBridge(id string, req UpdateBridgeRequest) (*Bridge, error)

Updates an existing bridge using PATCH (partial update).

**Endpoint:** `PATCH /provision/neural/bridges/:id`

**Parameters:**
- `id` (string): Bridge ID (required)
- `req` (UpdateBridgeRequest): Bridge update request
  - `Bridge` (UpdateBridgeData): Bridge update data (required)
    - `TargetSpaceID` (string): New target space ID (optional)

**Returns:**
- `*Bridge`: Updated bridge object with all fields including server-managed ProvisionState
- `error`: Error if request fails

#### ReplaceBridge(id string, req UpdateBridgeRequest) (*Bridge, error)

Replaces an existing bridge using PUT (full replacement).

**Endpoint:** `PUT /provision/neural/bridges/:id`

**Parameters:**
- `id` (string): Bridge ID (required)
- `req` (UpdateBridgeRequest): Bridge replacement request
  - `Bridge` (UpdateBridgeData): Bridge replacement data (required)
    - `TargetSpaceID` (string): New target space ID (required)

**Returns:**
- `*Bridge`: Replaced bridge object with all fields including server-managed ProvisionState
- `error`: Error if request fails

#### DeleteBridge(id string) error

Deletes a bridge by ID.

**Endpoint:** `DELETE /provision/neural/bridges/:id`

**Parameters:**
- `id` (string): Bridge ID (required)

**Returns:**
- `error`: Error if request fails

## Class Operations

#### GetClass(id string) (*Class, error)

Retrieves a specific class by ID.

**Endpoint:** `GET /provision/neural/classes/:id`

**Parameters:**
- `id` (string): Class ID (required)

**Returns:**
- `*Class`: Class object with ID, SpaceID, ProvisionState, Schema, Name, and Description
- `error`: Error if request fails

**Example:**
```go
class, err := client.Neural.GetClass("class-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Class: %+v\n", class)
```

#### CreateClass(spaceID string, req CreateClassRequest) (*Class, error)

Creates a new class within a space.

**Endpoint:** `POST /provision/neural/spaces/:space_id/classes`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `req` (CreateClassRequest): Class creation request
  - `Class` (ClassRequestData): Class data (required)
    - `Schema` (map[string]any): Class schema (required)

**Returns:**
- `*Class`: Created class object with ID, SpaceID, ProvisionState, Schema, Name, and Description
- `error`: Error if request fails

**Example:**
```go
req := tama.CreateClassRequest{
    Class: tama.ClassRequestData{
        Schema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "name": map[string]any{
                    "type": "string",
                },
            },
        },
    },
}
class, err := client.Neural.CreateClass("space-123", req)
``

#### UpdateClass(id string, req UpdateClassRequest) (*Class, error)

Updates an existing class using PATCH (partial update).

**Endpoint:** `PATCH /provision/neural/classes/:id`

**Parameters:**
- `id` (string): Class ID (required)
- `req` (UpdateClassRequest): Class update request
  - `Class` (UpdateClassData): Class update data (required)
    - `Schema` (map[string]any): New class schema (optional)

**Returns:**
- `*Class`: Updated class object with all fields including server-managed ProvisionState
- `error`: Error if request fails

#### ReplaceClass(id string, req UpdateClassRequest) (*Class, error)

Replaces an existing class using PUT (full replacement).

**Endpoint:** `PUT /provision/neural/classes/:id`

**Parameters:**
- `id` (string): Class ID (required)
- `req` (UpdateClassRequest): Class replacement request
  - `Class` (UpdateClassData): Class replacement data (required)
    - `Schema` (map[string]any): New class schema (required)

**Returns:**
- `*Class`: Replaced class object with all fields including server-managed ProvisionState
- `error`: Error if request fails

#### DeleteClass(id string) error

Deletes a class by ID.

**Endpoint:** `DELETE /provision/neural/classes/:id`

**Parameters:**
- `id` (string): Class ID (required)

**Returns:**
- `error`: Error if request fails

## Corpus Operations

#### GetCorpus(id string) (*Corpus, error)

Retrieves a specific corpus by ID.

**Endpoint:** `GET /provision/neural/corpora/:id`

**Parameters:**
- `id` (string): Corpus ID (required)

**Returns:**
- `*Corpus`: Corpus object with ID, Main, Name, Slug, Template, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
corpus, err := client.Neural.GetCorpus("corpus-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Corpus: %+v\n", corpus)
``

#### CreateCorpus(classID string, req CreateCorpusRequest) (*Corpus, error)

Creates a new corpus within a class.

**Endpoint:** `POST /provision/neural/classes/:class_id/corpora`

**Parameters:**
- `classID` (string): Class ID (required)
- `req` (CreateCorpusRequest): Corpus creation request
  - `Corpus` (CorpusRequestData): Corpus data (required)
    - `Main` (bool): Whether this is the main corpus (required)
    - `Name` (string): Corpus name (required)
    - `Template` (string): Corpus template (required)

**Returns:**
- `*Corpus`: Created corpus object with ID, Main, Name, Slug, Template, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
req := tama.CreateCorpusRequest{
    Corpus: tama.CorpusRequestData{
        Main:     true,
        Name:     "My Corpus",
        Template: "my-template",
    },
}
corpus, err := client.Neural.CreateCorpus("class-123", req)
``

#### UpdateCorpus(id string, req UpdateCorpusRequest) (*Corpus, error)

Updates an existing corpus using PATCH (partial update).

**Endpoint:** `PATCH /provision/neural/corpora/:id`

**Parameters:**
- `id` (string): Corpus ID (required)
- `req` (UpdateCorpusRequest): Corpus update request
  - `Corpus` (UpdateCorpusData): Corpus update data (required)
    - `Main` (*bool): New main flag (optional)
    - `Name` (string): New corpus name (optional)
    - `Template` (string): New corpus template (optional)

**Returns:**
- `*Corpus`: Updated corpus object with all fields including server-managed ProvisionState
- `error`: Error if request fails

#### ReplaceCorpus(id string, req UpdateCorpusRequest) (*Corpus, error)

Replaces an existing corpus using PUT (full replacement).

**Endpoint:** `PUT /provision/neural/corpora/:id`

**Parameters:**
- `id` (string): Corpus ID (required)
- `req` (UpdateCorpusRequest): Corpus replacement request
  - `Corpus` (UpdateCorpusData): Corpus replacement data (required)
    - `Main` (*bool): New main flag (required)
    - `Name` (string): New corpus name (required)
    - `Template` (string): New corpus template (required)

**Returns:**
- `*Corpus`: Replaced corpus object with all fields including server-managed ProvisionState
- `error`: Error if request fails

#### DeleteCorpus(id string) error

Deletes a corpus by ID.

**Endpoint:** `DELETE /provision/neural/corpora/:id`

**Parameters:**
- `id` (string): Corpus ID (required)

**Returns:**
- `error`: Error if request fails
