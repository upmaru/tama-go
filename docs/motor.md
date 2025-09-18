# Motor Service

Access via `client.Motor.*`

## Table of Contents

- [Action Operations](#action-operations)
- [Action Structure](#action-structure)
- [Modifier Operations](#modifier-operations)
- [Modifier Structure](#modifier-structure)

## Action Operations

### GetAction(specID, id string) (*Action, error)

Retrieves a specific motor action by specification ID and action ID.

**Endpoint:** `GET /provision/motor/specifications/:spec_id/actions/:id`

**Parameters:**
- `specID` (string): Specification ID (required)
- `id` (string): Action ID (required)

**Returns:**
- `*Action`: Action object with ID, Identifier, Path, Method, and SpecificationID
- `error`: Error if request fails

**Example:**
```go
action, err := client.Motor.GetAction("spec-123", "action-456")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Action: %+v\n", action)
```

### GetActionByPathAndMethod(specID, path, method string) (*Action, error)

Retrieves a specific motor action by specification ID, path, and method.

**Endpoint:** `GET /provision/motor/specifications/:spec_id/actions/:encoded_path`

**Parameters:**
- `specID` (string): Specification ID (required)
- `path` (string): Action path to search for (required)
- `method` (string): HTTP method to search for (required)

**Returns:**
- `*Action`: Action object with ID, Identifier, Path, Method, and SpecificationID
- `error`: Error if request fails

**Notes:**
- The `path` argument is automatically **URL‑safe Base64 encoded** by the client before it is appended to the request URL.
- The `method` string is converted to **lowercase** and sent as the `method` query parameter.
- These transformations are handled internally, so callers provide raw values (e.g., `/api/endpoint`, `"POST"`).

**Example:**
```go
action, err := client.Motor.GetActionByPathAndMethod("spec-123", "/api/endpoint", "POST")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Action: %+v\n", action)
```

### Execute() error

Executes the motor action.

**Parameters:**
- None

**Returns:**
- `error`: Error if execution fails

*Note:* The current implementation of `Execute` is a placeholder and simply returns `nil`. It does not perform any network call or side effect.

**Example:**
```go
action := &tama.Action{
    ID: "action-456",
    Path: "/api/endpoint",
    Method: "POST",
}
err := action.Execute()
if err != nil {
    log.Fatal(err)
}
```

## Action Structure

### Action

The `Action` struct represents a motor action with the following fields:

```go
type Action struct {
    ID              string `json:"id"`
    Identifier      string `json:"identifier"`
    Path            string `json:"path"`
    Method          string `json:"method"`
    SpecificationID string `json:"specification_id"`
}
```

**Fields:**
- `ID` (string): Unique identifier for the action
- `Identifier` (string): Human-readable identifier for the action
- `Path` (string): API endpoint path to execute
- `Method` (string): HTTP method to use for execution (GET, POST, PUT, DELETE, etc.)
- `SpecificationID` (string): ID of the specification this action belongs to

### ActionResponse

The `ActionResponse` struct wraps the Action data in API responses:

```go
type ActionResponse struct {
    Data Action `json:"data"`
}

## JSON Payloads

The API responses are wrapped in an `ActionResponse` object.
Example response:

```go
{
  "data": {
    "id": "action-456",
    "identifier": "CreateUser",
    "path": "/api/endpoint",
    "method": "POST",
    "specification_id": "spec-123"
  }
}
```

The request payloads for motor actions are not required for the current client methods, as all data is retrieved from the server.

**Fields:**
- `Data` (Action): The action data returned by the API

## Modifier Operations

### GetModifier(id string) (*Modifier, error)

Retrieves a specific motor modifier by ID.

**Endpoint:** `GET /provision/motor/modifiers/:id`

**Parameters:**
- `id` (string): Modifier ID (required)

**Returns:**
- `*Modifier`: Modifier object with ID, Name, ActionID, Schema, and ProvisionState
- `error`: Error if request fails

**Example:**
```go
modifier, err := client.Motor.GetModifier("modifier-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Modifier: %+v\n", modifier)
```

### CreateModifier(actionID string, req CreateModifierRequest) (*Modifier, error)

Creates a new modifier under a specific action.

**Endpoint:** `POST /provision/motor/actions/:action_id/modifiers`

**Parameters:**
- `actionID` (string): Parent action ID (required)
- `req` (CreateModifierRequest): Modifier creation request (required)
  - `Modifier` (ModifierRequestData): Modifier data (required)
    - `Name` (string): Modifier name (required)
    - `Schema` (map[string]any): Arbitrary configuration (required)

**Returns:**
- `*Modifier`: Created modifier object
- `error`: Error if request fails

**Request Structure:**
```go
type CreateModifierRequest struct {
    Modifier ModifierRequestData `json:"modifier"`
}

type ModifierRequestData struct {
    Name   string         `json:"name"`
    Schema map[string]any `json:"schema"`
}
```

**Example:**
```go
req := motor.CreateModifierRequest{
    Modifier: motor.ModifierRequestData{
        Name:   "sanitize",
        Schema: map[string]any{"rule": "trim"},
    },
}
created, err := client.Motor.CreateModifier("action-456", req)
```

### UpdateModifier(id string, req UpdateModifierRequest) (*Modifier, error)

Updates an existing modifier using PATCH (partial update).

**Endpoint:** `PATCH /provision/motor/modifiers/:id`

**Parameters:**
- `id` (string): Modifier ID (required)
- `req` (UpdateModifierRequest): Modifier update request (required)
  - `Modifier` (UpdateModifierData): Updatable fields
    - `Name` (string, optional)
    - `Schema` (map[string]any, optional)

**Returns:**
- `*Modifier`: Updated modifier object
- `error`: Error if request fails

**Request Structure:**
```go
type UpdateModifierRequest struct {
    Modifier UpdateModifierData `json:"modifier"`
}

type UpdateModifierData struct {
    Name   string         `json:"name,omitempty"`
    Schema map[string]any `json:"schema,omitempty"`
}
```

**Example:**
```go
update := motor.UpdateModifierRequest{
    Modifier: motor.UpdateModifierData{
        Name:   "normalize",
        Schema: map[string]any{"rule": "lowercase"},
    },
}
updated, err := client.Motor.UpdateModifier("modifier-123", update)
```

### ReplaceModifier(id string, req UpdateModifierRequest) (*Modifier, error)

Replaces an existing modifier using PUT (full replacement semantics on the server side).

**Endpoint:** `PUT /provision/motor/modifiers/:id`

**Parameters:**
- `id` (string): Modifier ID (required)
- `req` (UpdateModifierRequest): Replacement request (required)

**Returns:**
- `*Modifier`: Replaced modifier object
- `error`: Error if request fails

**Example:**
```go
replace := motor.UpdateModifierRequest{
    Modifier: motor.UpdateModifierData{
        Name:   "sanitize",
        Schema: map[string]any{"rule": "strip"},
    },
}
replaced, err := client.Motor.ReplaceModifier("modifier-123", replace)
```

### DeleteModifier(id string) error

Deletes a modifier by ID.

**Endpoint:** `DELETE /provision/motor/modifiers/:id`

**Parameters:**
- `id` (string): Modifier ID (required)

**Returns:**
- `error`: Error if request fails

**Example:**
```go
if err := client.Motor.DeleteModifier("modifier-123"); err != nil {
    log.Fatal(err)
}
```

## Modifier Structure

### Modifier

The `Modifier` struct represents a motor modifier with the following fields:

```go
type Modifier struct {
    ID             string         `json:"id,omitempty"`
    Name           string         `json:"name"`
    ActionID       string         `json:"action_id,omitempty"`
    Schema         map[string]any `json:"schema"`
    ProvisionState string         `json:"provision_state,omitempty"`
}
```

### ModifierResponse

Wraps a Modifier in API responses:

```go
type ModifierResponse struct {
    Data Modifier `json:"data"`
}
```
