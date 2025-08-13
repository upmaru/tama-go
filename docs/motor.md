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
- The path is URL-safe base64 encoded before being used in the endpoint
- The method is converted to lowercase before being passed as a query parameter

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
```

**Fields:**
- `Data` (Action): The action data returned by the API