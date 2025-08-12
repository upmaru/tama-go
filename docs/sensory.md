## Sensory Service

Access via `client.Sensory.*`

### Source Operations

#### GetSource(id string) (*Source, error)

Retrieves a specific source by ID.

**Endpoint:** `GET /provision/sensory/sources/:id`

**Parameters:**
- `id` (string): Source ID (required)

**Returns:**
- `*Source`: Source object
- `error`: Error if request fails

#### GetSourceBySpecificationAndSlug(specificationID string, slug string) (*Source, error)

Retrieves a source by specification ID and source slug.

**Endpoint:** `GET /provision/sensory/specifications/:specification_id/sources/:id`

**Parameters:**
- `specificationID` (string): Specification ID (required)
- `slug` (string): Source slug (required)

**Returns:**
- `*Source`: Source object
- `error`: Error if request fails

#### CreateSource(spaceID string, req CreateSourceRequest) (*Source, error)

Creates a new source in a specific space.

**Endpoint:** `POST /provision/sensory/spaces/:space_id/sources`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `req` (CreateSourceRequest): Source creation request
  - `Source` (SourceRequestData): Source data (required)
    - `Name` (string): Source name (required)
    - `Type` (string): Source type (required)
    - `Endpoint` (string): Source endpoint URL (required)
    - `Credential` (SourceCredential): Source credentials (required)

**Returns:**
- `*Source`: Created source object with ID, Name, Endpoint, SpaceID, and server-managed ProvisionState
- `error`: Error if request fails

**Note:** The `ProvisionState` and `SpaceID` fields are managed server-side and cannot be set during creation.

#### UpdateSource(id string, req UpdateSourceRequest) (*Source, error)

Updates an existing source using PATCH.

**Endpoint:** `PATCH /provision/sensory/sources/:id`

**Parameters:**
- `id` (string): Source ID (required)
- `req` (UpdateSourceRequest): Update request

**Returns:**
- `*Source`: Updated source object with all fields including server-managed ProvisionState and SpaceID
- `error`: Error if request fails

**Note:** The `ProvisionState` and `SpaceID` fields cannot be updated via API calls - they are managed server-side.

#### ReplaceSource(id string, req UpdateSourceRequest) (*Source, error)

Replaces an existing source using PUT.

**Endpoint:** `PUT /provision/sensory/sources/:id`

**Parameters:**
- `id` (string): Source ID (required)
- `req` (UpdateSourceRequest): Replacement request

**Returns:**
- `*Source`: Updated source object with all fields including server-managed ProvisionState and SpaceID
- `error`: Error if request fails

**Note:** The `ProvisionState` and `SpaceID` fields cannot be updated via API calls - they are managed server-side.

#### DeleteSource(id string) error

Deletes a source by ID.

**Endpoint:** `DELETE /provision/sensory/sources/:id`

**Parameters:**
- `id` (string): Source ID (required)

### Model Operations

#### GetModel(id string) (*Model, error)

Retrieves a specific model by ID.

**Endpoint:** `GET /provision/sensory/models/:id`

**Parameters:**
- `id` (string): Model ID (required)

**Returns:**
- `*Model`: Model object with ID, Identifier, Path, Parameters, and server-managed ProvisionState
- `error`: Error if request fails

**Example:**
```go
model, err := client.Sensory.GetModel("model-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Model: %+v\n", model)
```

#### CreateModel(sourceID string, req CreateModelRequest) (*Model, error)

Creates a new model for a specific source.

**Endpoint:** `POST /provision/sensory/sources/:source_id/models`

**Parameters:**
- `sourceID` (string): Source ID (required)
- `req` (CreateModelRequest): Model creation request
  - `Model` (ModelRequestData): Model data (required)
    - `Identifier` (string): Model identifier (required)
    - `Path` (string): Model path (required)
    - `Parameters` (map[string]any): Model parameters (optional)

**Returns:**
- `*Model`: Created model object with ID, Identifier, Path, Parameters, and server-managed ProvisionState
- `error`: Error if request fails

**Note:** The `ProvisionState` field is managed server-side and cannot be set during creation.

**Example:**
```go
req := sensory.CreateModelRequest{
    Model: sensory.ModelRequestData{
        Identifier: "mistral-large-latest",
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
}
model, err := client.Sensory.CreateModel("source-123", req)
```

**Parameters Field:**
The `Parameters` field accepts any valid JSON values:
- Strings: `"reasoning_effort": "low"`
- Numbers: `"temperature": 1.0`, `"max_tokens": 2000`
- Booleans: `"stream": true`
- Arrays: `"stop": []string{"\n", "###"}`
- Objects: `"config": map[string]any{...}`

#### UpdateModel(id string, req UpdateModelRequest) (*Model, error)

Updates an existing model using PATCH.

**Endpoint:** `PATCH /provision/sensory/models/:id`

**Parameters:**
- `id` (string): Model ID (required)
- `req` (UpdateModelRequest): Update request
  - `Model` (UpdateModelData): Model update data (required)
    - `Identifier` (string): New model identifier (optional)
    - `Path` (string): New model path (optional)
    - `Parameters` (map[string]any): New model parameters (optional)

**Returns:**
- `*Model`: Updated model object with all fields including Parameters and server-managed ProvisionState
- `error`: Error if request fails

**Note:** The `ProvisionState` field cannot be updated via API calls - it is managed server-side.

#### ReplaceModel(id string, req UpdateModelRequest) (*Model, error)

Replaces an existing model using PUT.

**Endpoint:** `PUT /provision/sensory/models/:id`

**Parameters:**
- `id` (string): Model ID (required)
- `req` (UpdateModelRequest): Replacement request
  - `Model` (UpdateModelData): Model update data (required)
    - `Identifier` (string): New model identifier (optional)
    - `Path` (string): New model path (optional)
    - `Parameters` (map[string]any): New model parameters (optional)

**Returns:**
- `*Model`: Updated model object with all fields including Parameters and server-managed ProvisionState
- `error`: Error if request fails

**Note:** The `ProvisionState` field cannot be updated via API calls - it is managed server-side.

#### DeleteModel(id string) error

Deletes a model by ID.

**Endpoint:** `DELETE /provision/sensory/models/:id`

**Parameters:**
- `id` (string): Model ID (required)

**Returns:**
- `error`: Error if request fails

### Limit Operations

#### GetLimit(id string) (*Limit, error)

Retrieves a specific limit by ID.

**Endpoint:** `GET /provision/sensory/limits/:id`

#### CreateLimit(sourceID string, req CreateLimitRequest) (*Limit, error)

Creates a new limit for a specific source.

**Endpoint:** `POST /provision/sensory/sources/:source_id/limits`

**Parameters:**
- `sourceID` (string): Source ID (required)
- `req` (CreateLimitRequest): Limit creation request
  - `Limit` (LimitRequestData): Limit data (required)
    - `ScaleUnit` (string): Scale unit (required)
    - `ScaleCount` (int): Scale count (required, must be  0)
    - `Count` (int): Count value (required, must be  0)

**Note:** The created limit will automatically be associated with the specified source via its `source_id` field.

#### UpdateLimit(id string, req UpdateLimitRequest) (*Limit, error)

Updates an existing limit using PATCH.

**Endpoint:** `PATCH /provision/sensory/limits/:id`

#### ReplaceLimit(id string, req UpdateLimitRequest) (*Limit, error)

Replaces an existing limit using PUT.

**Endpoint:** `PUT /provision/sensory/limits/:id`

#### DeleteLimit(id string) error

Deletes a limit by ID.

**Endpoint:** `DELETE /provision/sensory/limits/:id`

### Specification Operations

#### GetSpecification(id string) (*Specification, error)

Retrieves a specific specification by ID.

**Endpoint:** `GET /provision/sensory/specifications/:id`

**Parameters:**
- `id` (string): Specification ID (required)

**Returns:**
- `*Specification`: Specification object with ID, SpaceID, Schema, Version, Endpoint, CurrentState, and server-managed ProvisionState
- `error`: Error if request fails

**Example:**
```go
spec, err := client.Sensory.GetSpecification("spec-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Specification: %+v\n", spec)
```

#### CreateSpecification(spaceID string, req CreateSpecificationRequest) (*Specification, error)

Creates a new specification in a specific space.

**Endpoint:** `POST /provision/sensory/spaces/:space_id/specifications`

**Parameters:**
- `spaceID` (string): Space ID (required)
- `req` (CreateSpecificationRequest): Specification creation request
  - `Specification` (SpecificationRequestData): Specification data (required)
    - `Schema` (map[string]any): JSON schema definition (required)
    - `Version` (string): Specification version (required)
    - `Endpoint` (string): Specification endpoint URL (required)

**Returns:**
- `*Specification`: Created specification object with ID, SpaceID, Schema, Version, Endpoint, and server-managed CurrentState and ProvisionState
- `error`: Error if request fails

**Note:** The `CurrentState` and `ProvisionState` fields are managed server-side and cannot be set during creation.

**Example:**
```go
createReq := sensory.CreateSpecificationRequest{
    Specification: sensory.SpecificationRequestData{
        Schema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "message": map[string]any{
                    "type": "string",
                },
            },
        },
        Version:  "1.0.0",
        Endpoint: "https://api.example.com/v1",
    },
}

spec, err := client.Sensory.CreateSpecification("space-123", createReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created specification: %+v\n", spec)
```

#### UpdateSpecification(id string, req UpdateSpecificationRequest) (*Specification, error)

Updates an existing specification using PATCH.

**Endpoint:** `PATCH /provision/sensory/specifications/:id`

**Parameters:**
- `id` (string): Specification ID (required)
- `req` (UpdateSpecificationRequest): Update request
  - `Specification` (UpdateSpecificationData): Specification update data
    - `Schema` (map[string]any): JSON schema definition (optional)
    - `Version` (string): Specification version (optional)
    - `Endpoint` (string): Specification endpoint URL (optional)

**Returns:**
- `*Specification`: Updated specification object with all fields including server-managed CurrentState and ProvisionState
- `error`: Error if request fails

**Note:** The `CurrentState` and `ProvisionState` fields are managed server-side and cannot be updated via API calls.

#### ReplaceSpecification(id string, req UpdateSpecificationRequest) (*Specification, error)

Replaces an existing specification using PUT.

**Endpoint:** `PUT /provision/sensory/specifications/:id`

**Parameters:**
- `id` (string): Specification ID (required)
- `req` (UpdateSpecificationRequest): Replacement request

**Returns:**
- `*Specification`: Updated specification object with all fields
- `error`: Error if request fails

#### DeleteSpecification(id string) error

Deletes a specification by ID.

**Endpoint:** `DELETE /provision/sensory/specifications/:id`

**Parameters:**
- `id` (string): Specification ID (required)

**Returns:**
- `error`: Error if request fails

### Identity Operations

#### GetIdentity(id string) (*Identity, error)

Retrieves a specific identity by ID.

**Endpoint:** `GET /provision/sensory/identities/:id`

**Parameters:**
- `id` (string): Identity ID (required)

**Returns:**
- `*Identity`: Identity data
- `error`: Error if request fails

**Response:**
```go
{
  "data": {
    "id": "identity-123",
    "specification_id": "spec-456",
    "provision_state": "active",
    "current_state": "running",
    "identifier": "test-identifier",
    "validation": {
      "path": "/health",
      "method": "GET",
      "codes": [200]
    }
  }
}
```

#### CreateIdentity(specificationID, identifier string, req CreateIdentityRequest) (*Identity, error)

Creates a new identity for a specific specification and identifier.

**Endpoint:** `POST /provision/sensory/specifications/:specification_id/identifiers/:identifier/identities`

**Parameters:**
- `specificationID` (string): Specification ID (required)
- `identifier` (string): Identifier (required)
- `req` (CreateIdentityRequest): Identity data (required)

**Request:**
```go
type CreateIdentityRequest struct {
  Identity IdentityRequestData `json:"identity"`
}

type IdentityRequestData struct {
  APIKey     string     `json:"api_key"`
  Validation Validation `json:"validation"`
}

type Validation struct {
  Path   string `json:"path"`
  Method string `json:"method"`
  Codes  []int  `json:"codes"`
}
```

**Returns:**
- `*Identity`: Created identity data
- `error`: Error if request fails

#### UpdateIdentity(id string, req UpdateIdentityRequest) (*Identity, error)

Updates an existing identity using PATCH.

**Endpoint:** `PATCH /provision/sensory/identities/:id`

**Parameters:**
- `id` (string): Identity ID (required)
- `req` (UpdateIdentityRequest): Update data

**Request:**
```go
type UpdateIdentityRequest struct {
  Identity UpdateIdentityData `json:"identity"`
}

type UpdateIdentityData struct {
  APIKey     string      `json:"api_key,omitempty"`
  Validation *Validation `json:"validation,omitempty"`
}
```

**Returns:**
- `*Identity`: Updated identity data
- `error`: Error if request fails

#### ReplaceIdentity(id string, req UpdateIdentityRequest) (*Identity, error)

Replaces an existing identity using PUT.

**Endpoint:** `PUT /provision/sensory/identities/:id`

**Parameters:**
- `id` (string): Identity ID (required)
- `req` (UpdateIdentityRequest): Replacement data

**Returns:**
- `*Identity`: Replaced identity data
- `error`: Error if request fails

#### DeleteIdentity(id string) error

Deletes an identity by ID.

**Endpoint:** `DELETE /provision/sensory/identities/:id`

**Parameters:**
- `id` (string): Identity ID (required)

**Returns:**
- `error`: Error if request fails
