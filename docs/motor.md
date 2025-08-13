# Motor Service

Access via `client.Motor.*`

## Table of Contents

- [Action Operations](#action-operations)
- [Action Structure](#action-structure)

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
