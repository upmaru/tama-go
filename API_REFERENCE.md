# Tama Go Client API Reference

This document provides a comprehensive reference for all methods available in the Tama Go client library.

## Table of Contents

- [Client Configuration](#client-configuration)
- [Neural Service](#neural-service)
- [Sensory Service](#sensory-service)
- [Error Handling](#error-handling)
- [Data Types](#data-types)

## Client Configuration

### NewClient(config Config) *Client

Creates a new Tama API client with the provided configuration.

**Parameters:**
- `config` (Config): Configuration object containing:
  - `BaseURL` (string): The base URL of the Tama API (required)
  - `APIKey` (string): Your API authentication key (required)
  - `Timeout` (time.Duration): Request timeout (optional, default: 30s)

**Returns:**
- `*Client`: Configured client instance

**Example:**
```go
config := tama.Config{
    BaseURL: "https://api.tama.io",
    APIKey:  "your-api-key",
    Timeout: 30 * time.Second,
}
client := tama.NewClient(config)
```

### Client Methods

#### SetAPIKey(apiKey string)

Updates the API key for authentication.

**Parameters:**
- `apiKey` (string): New API key

#### SetDebug(debug bool)

Enables or disables debug mode for HTTP requests.

**Parameters:**
- `debug` (bool): Enable/disable debug mode

## Neural Service

Access via `client.Neural.*`

### GetSpace(id string) (*Space, error)

Retrieves a specific neural space by ID.

**Endpoint:** `GET /provision/neural/spaces/:id`

**Parameters:**
- `id` (string): Space ID (required)

**Returns:**
- `*Space`: Space object with ID, Name, Slug, Type, and CurrentState
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
    - `Type` (string): Space type - "root" or "component" (required)

**Returns:**
- `*Space`: Created space object with ID, Name, Slug, Type, and CurrentState
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
    - `Type` (string): New space type - "root" or "component" (optional)

**Returns:**
- `*Space`: Updated space object with all fields including server-managed CurrentState
- `error`: Error if request fails

### ReplaceSpace(id string, req UpdateSpaceRequest) (*Space, error)

Replaces an existing space using PUT (full replacement).

**Endpoint:** `PUT /provision/neural/spaces/:id`

**Parameters:**
- `id` (string): Space ID (required)
- `req` (UpdateSpaceRequest): Replacement request

**Returns:**
- `*Space`: Updated space object with all fields including server-managed CurrentState
- `error`: Error if request fails

### DeleteSpace(id string) error

Deletes a space by ID.

**Endpoint:** `DELETE /provision/neural/spaces/:id`

**Parameters:**
- `id` (string): Space ID (required)

**Returns:**
- `error`: Error if request fails

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
- `*Source`: Created source object with ID, Name, Endpoint, SpaceID, and server-managed CurrentState
- `error`: Error if request fails

**Note:** The `CurrentState` and `SpaceID` fields are managed server-side and cannot be set during creation.

#### UpdateSource(id string, req UpdateSourceRequest) (*Source, error)

Updates an existing source using PATCH.

**Endpoint:** `PATCH /provision/sensory/sources/:id`

**Parameters:**
- `id` (string): Source ID (required)
- `req` (UpdateSourceRequest): Update request

**Returns:**
- `*Source`: Updated source object with all fields including server-managed CurrentState and SpaceID
- `error`: Error if request fails

**Note:** The `CurrentState` and `SpaceID` fields cannot be updated via API calls - they are managed server-side.

#### ReplaceSource(id string, req UpdateSourceRequest) (*Source, error)

Replaces an existing source using PUT.

**Endpoint:** `PUT /provision/sensory/sources/:id`

**Parameters:**
- `id` (string): Source ID (required)
- `req` (UpdateSourceRequest): Replacement request

**Returns:**
- `*Source`: Updated source object with all fields including server-managed CurrentState and SpaceID
- `error`: Error if request fails

**Note:** The `CurrentState` and `SpaceID` fields cannot be updated via API calls - they are managed server-side.

#### DeleteSource(id string) error

Deletes a source by ID.

**Endpoint:** `DELETE /provision/sensory/sources/:id`

### Model Operations

#### GetModel(id string) (*Model, error)

Retrieves a specific model by ID.

**Endpoint:** `GET /provision/sensory/models/:id`

**Parameters:**
- `id` (string): Model ID (required)

**Returns:**
- `*Model`: Model object with ID, Identifier, Path, and server-managed CurrentState
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

**Returns:**
- `*Model`: Created model object with ID, Identifier, Path, and server-managed CurrentState
- `error`: Error if request fails

**Note:** The `CurrentState` field is managed server-side and cannot be set during creation.

**Example:**
```go
req := sensory.CreateModelRequest{
    Model: sensory.ModelRequestData{
        Identifier: "mistral-large-latest",
        Path:       "/chat/completions",
    },
}
model, err := client.Sensory.CreateModel("source-123", req)
```

#### UpdateModel(id string, req UpdateModelRequest) (*Model, error)

Updates an existing model using PATCH.

**Endpoint:** `PATCH /provision/sensory/models/:id`

**Parameters:**
- `id` (string): Model ID (required)
- `req` (UpdateModelRequest): Update request
  - `Model` (UpdateModelData): Model update data (required)
    - `Identifier` (string): New model identifier (optional)
    - `Path` (string): New model path (optional)

**Returns:**
- `*Model`: Updated model object with all fields including server-managed CurrentState
- `error`: Error if request fails

**Note:** The `CurrentState` field cannot be updated via API calls - it is managed server-side.

#### ReplaceModel(id string, req UpdateModelRequest) (*Model, error)

Replaces an existing model using PUT.

**Endpoint:** `PUT /provision/sensory/models/:id`

**Parameters:**
- `id` (string): Model ID (required)
- `req` (UpdateModelRequest): Replacement request
  - `Model` (UpdateModelData): Model update data (required)
    - `Identifier` (string): New model identifier (optional)
    - `Path` (string): New model path (optional)

**Returns:**
- `*Model`: Updated model object with all fields including server-managed CurrentState
- `error`: Error if request fails

**Note:** The `CurrentState` field cannot be updated via API calls - it is managed server-side.

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
    - `ScaleCount` (int): Scale count (required, must be > 0)
    - `Count` (int): Count value (required, must be > 0)

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

## Error Handling

### Error Type

The client returns structured errors of type `*Error`:

```go
type Error struct {
    StatusCode int    `json:"status_code"`
    Message    string `json:"message"`
    Details    string `json:"details,omitempty"`
}
```

### Error Handling Example

```go
space, err := client.Neural.GetSpace("invalid-id")
if err != nil {
    if apiErr, ok := err.(*tama.Error); ok {
        // API error
        fmt.Printf("API Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
        if apiErr.Details != "" {
            fmt.Printf("Details: %s\n", apiErr.Details)
        }
    } else {
        // Client/network error
        fmt.Printf("Client Error: %v\n", err)
    }
}
```

## Data Types

### Core Resources

#### Space

```go
type Space struct {
    ID           string `json:"id,omitempty"`
    Name         string `json:"name"`
    Slug         string `json:"slug,omitempty"`
    Type         string `json:"type"`
    CurrentState string `json:"current_state"`
}
```

#### Source

```go
type Source struct {
    ID           string `json:"id,omitempty"`
    Name         string `json:"name"`
    Endpoint     string `json:"endpoint"`
    SpaceID      string `json:"space_id"`
    CurrentState string `json:"current_state"`
}
```

#### Model

```go
type Model struct {
    ID           string `json:"id,omitempty"`
    Identifier   string `json:"identifier"`
    Path         string `json:"path"`
    CurrentState string `json:"current_state"`
}
```

#### Limit

```go
type Limit struct {
    ID           string `json:"id,omitempty"`
    SourceID     string `json:"source_id"`
    Count        int    `json:"count"`
    ScaleUnit    string `json:"scale_unit"`
    ScaleCount   int    `json:"scale_count"`
    CurrentState string `json:"current_state"`
}
```

### Request Types

#### CreateSpaceRequest

```go
type CreateSpaceRequest struct {
    Space SpaceRequest `json:"space"`
}
```

#### UpdateSpaceRequest

```go
type UpdateSpaceRequest struct {
    Space UpdateSpaceData `json:"space"`
}
```

#### CreateSourceRequest

```go
type CreateSourceRequest struct {
    Source SourceRequestData `json:"source"`
}
```

#### UpdateSourceRequest

```go
type UpdateSourceRequest struct {
    Source UpdateSourceData `json:"source"`
}
```

#### CreateModelRequest

```go
type CreateModelRequest struct {
    Model ModelRequestData `json:"model"`
}
```

#### UpdateModelRequest

```go
type UpdateModelRequest struct {
    Model UpdateModelData `json:"model"`
}
```

#### CreateLimitRequest

```go
type CreateLimitRequest struct {
    Limit LimitRequestData `json:"limit"`
}
```

#### UpdateLimitRequest

```go
type UpdateLimitRequest struct {
    Limit UpdateLimitData `json:"limit"`
}
```

## Configuration Types

#### Config

```go
type Config struct {
    BaseURL string
    APIKey  string
    Timeout time.Duration
}
```

#### Response

```go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *Error      `json:"error,omitempty"`
}
```

#### SpaceRequest

```go
type SpaceRequest struct {
    Name string `json:"name"`
    Type string `json:"type"` // "root" or "component"
}
```

#### UpdateSpaceData

```go
type UpdateSpaceData struct {
    Name string `json:"name,omitempty"`
    Type string `json:"type,omitempty"` // "root" or "component"
}
```

**Note:** The `CurrentState` field cannot be updated via API calls - it is managed server-side.

#### SpaceResponse

```go
type SpaceResponse struct {
    Data Space `json:"data"`
}
```

#### SourceRequestData

```go
type SourceRequestData struct {
    Name       string           `json:"name"`
    Type       string           `json:"type"`
    Endpoint   string           `json:"endpoint"`
    Credential SourceCredential `json:"credential"`
}
```

#### UpdateSourceData

```go
type UpdateSourceData struct {
    Name       string            `json:"name,omitempty"`
    Type       string            `json:"type,omitempty"`
    Endpoint   string            `json:"endpoint,omitempty"`
    Credential *SourceCredential `json:"credential,omitempty"`
}
```

**Note:** The `CurrentState` and `SpaceID` fields cannot be updated via API calls - they are managed server-side.

#### SourceResponse

```go
type SourceResponse struct {
    Data Source `json:"data"`
}
```

#### SourceCredential

```go
type SourceCredential struct {
    ApiKey string `json:"api_key"`
}
```

#### ModelRequestData

```go
type ModelRequestData struct {
    Identifier string `json:"identifier"`
    Path       string `json:"path"`
}
```

#### UpdateModelData

```go
type UpdateModelData struct {
    Identifier string `json:"identifier,omitempty"`
    Path       string `json:"path,omitempty"`
}
```

**Note:** The `CurrentState` field cannot be updated via API calls - it is managed server-side.

#### ModelResponse

```go
type ModelResponse struct {
    Data Model `json:"data"`
}
```

#### LimitRequestData

```go
type LimitRequestData struct {
    ScaleUnit  string `json:"scale_unit"`
    ScaleCount int    `json:"scale_count"`
    Count      int    `json:"count"`
}
```

#### UpdateLimitData

```go
type UpdateLimitData struct {
    ScaleUnit    string `json:"scale_unit,omitempty"`
    ScaleCount   int    `json:"scale_count,omitempty"`
    Count        int    `json:"count,omitempty"`
    CurrentState string `json:"current_state,omitempty"`
}
```

#### LimitResponse

```go
type LimitResponse struct {
    Data Limit `json:"data"`
}
```

#### SpaceRequest