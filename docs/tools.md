# Tools Service

The Tools service provisions inputs, initializers, outputs, output options, and
trusted thought-tool modifiers. Its methods are exposed directly through
`Client.Tools`:

```go
client := tama.NewClient(config)

input, err := client.Tools.GetInput("input-123")
modifier, err := client.Tools.GetModifier("modifier-123")
```

Methods do not take a `context.Context`; request timeouts are configured on the
root `tama.Client`.

## Supported operations

| Resource | Get | Create parent | Update | Replace | Delete |
| --- | --- | --- | --- | --- | --- |
| Input | `GetInput` | `CreateInput(thoughtToolID, ...)` | `UpdateInput` | `ReplaceInput` | `DeleteInput` |
| Initializer | `GetInitializer` | `CreateInitializer(thoughtToolID, ...)` | `UpdateInitializer` | `ReplaceInitializer` | `DeleteInitializer` |
| Output | `GetOutput` | `CreateOutput(thoughtToolID, ...)` | `UpdateOutput` | `ReplaceOutput` | `DeleteOutput` |
| Option | `GetOption` | `CreateOption(outputID, ...)` | `UpdateOption` | `ReplaceOption` | `DeleteOption` |
| Modifier | `GetModifier` | `CreateModifier(thoughtToolID, ...)` | `UpdateModifier` | `ReplaceModifier` | `DeleteModifier` |

The Tama provision API does not currently expose Tools list endpoints. Phoenix
routes both `PATCH` and `PUT` to each resource's update action. Consequently,
the `Replace*` methods use the same request fields as `Update*`; callers should
not assume that PUT can change server-owned or immutable fields.

## Inputs

An input associates a thought tool with a class corpus and input type.

```go
created, err := client.Tools.CreateInput("tool-123", tools.CreateInputRequest{
    Input: tools.InputRequestData{
        Type:          "text",
        ClassCorpusID: "corpus-789",
    },
})

updated, err := client.Tools.UpdateInput(created.ID, tools.UpdateInputRequest{
    Input: tools.UpdateInputData{Type: "json"},
})
```

## Initializers

An initializer supplies a reference, optional index, and arbitrary parameters
used to initialize a thought tool.

```go
index := 0
created, err := client.Tools.CreateInitializer("tool-123", tools.CreateInitializerRequest{
    Initializer: tools.InitializerRequestData{
        Reference:  "tama/example/initializer",
        Index:      &index,
        Parameters: map[string]any{"temperature": 0.2},
    },
})
```

## Outputs and options

An output associates a thought tool with a class corpus. An option associates
that output with a Motor action modifier.

```go
output, err := client.Tools.CreateOutput("tool-123", tools.CreateOutputRequest{
    Output: tools.OutputRequestData{ClassCorpusID: "corpus-789"},
})

option, err := client.Tools.CreateOption(output.ID, tools.CreateOptionRequest{
    Option: tools.OptionRequestData{ActionModifierID: "action-modifier-123"},
})
```

`tools.Option` and `tools.Modifier` are different resources. An option refers to
an action-level `motor.Modifier`; a trusted tool modifier populates request
arguments from runtime metadata.

## Trusted tool modifiers

A trusted tool modifier reserves a request field for a value obtained from
authoritative runtime metadata instead of model-generated arguments.

### Types and constants

```go
type ModifierSource struct {
    Type string `json:"type"`
    Path string `json:"path"`
}

type Modifier struct {
    ID              string
    Index           int
    Target          string
    OnMissingParent string
    OnMissingSource string
    Source          ModifierSource
    ThoughtToolID   string
    ProvisionState  string
}
```

Use the exported wire-value constants:

```go
tools.ModifierMissingPolicyError
tools.ModifierMissingPolicySkip
tools.ModifierSourceTypeMetadata
tools.ModifierSourcePathActorIdentifier
tools.ModifierSourcePathOriginEntityIdentifier
tools.ModifierSourcePathCurrentTimestamp
```

### Create or reactivate

```go
modifier, err := client.Tools.CreateModifier("tool-123", tools.CreateModifierRequest{
    Modifier: tools.ModifierRequestData{
        Index:           0,
        Target:          "/body/search/scope/user_id",
        OnMissingParent: tools.ModifierMissingPolicySkip,
        OnMissingSource: tools.ModifierMissingPolicyError,
        Source: tools.ModifierSource{
            Type: tools.ModifierSourceTypeMetadata,
            Path: tools.ModifierSourcePathActorIdentifier,
        },
    },
})
```

Create sends `POST /provision/tools/:thought_tool_id/modifiers`. Tama activates a
new modifier and returns `201`. If the exact configuration already exists in an
inactive state, Tama reactivates it and returns the same ID.

All create fields are required. The client rejects empty structural fields and
negative indexes before sending a request. Tama remains authoritative for
allowed values, JSON Pointer syntax, callable-schema compatibility, and active
modifier conflicts.

### Get

```go
modifier, err := client.Tools.GetModifier("modifier-123")
```

Get sends `GET /provision/tools/modifiers/:id`. Only active modifiers are
visible; invalid, unknown, and inactive IDs return `404`.

### Update

```go
source := tools.ModifierSource{
    Type: tools.ModifierSourceTypeMetadata,
    Path: tools.ModifierSourcePathOriginEntityIdentifier,
}

modifier, err := client.Tools.UpdateModifier("modifier-123", tools.UpdateModifierRequest{
    Modifier: tools.UpdateModifierData{
        OnMissingParent: tools.ModifierMissingPolicyError,
        Source:          &source,
    },
})
```

Update sends `PATCH /provision/tools/modifiers/:id` and supports partial
updates. `Source` is a pointer so it can be omitted; when supplied, both its
`Type` and `Path` are required. `Index` is not part of the update type because
Tama makes it immutable.

### Replace

```go
modifier, err := client.Tools.ReplaceModifier("modifier-123", tools.UpdateModifierRequest{
    Modifier: tools.UpdateModifierData{
        Target:          "/body/search/scope/account_id",
        OnMissingParent: tools.ModifierMissingPolicyError,
        OnMissingSource: tools.ModifierMissingPolicyError,
        Source:          &source,
    },
})
```

Replace sends `PUT /provision/tools/modifiers/:id`. Tama routes PUT to the same
update action as PATCH, so the mutable request type is shared and `index`
remains unavailable.

### Delete

```go
err := client.Tools.DeleteModifier("modifier-123")
```

Delete sends `DELETE /provision/tools/modifiers/:id`. Tama deactivates the
resource and returns it with `provision_state` set to `inactive`; the client
discards that response body.

## Errors

Parsed API failures are returned as `*tools.Error`. This permits callers to
distinguish refresh-time not-found responses and inspect flattened validation
fields:

```go
var apiErr *tools.Error
if errors.As(err, &apiErr) {
    if apiErr.StatusCode == http.StatusNotFound {
        // The active resource no longer exists.
    }

    sourcePathMessages := apiErr.Errors["source.path"]
}
```

Local structural validation errors are ordinary Go errors. Server responses,
including `404`, `401`, and `422`, retain their HTTP status through
`tools.Error.StatusCode`.

## Modifier authorization

A credential needs both narrow capabilities for a complete lifecycle:

```text
provision.tools.modifier.read
provision.tools.modifier.manage
```

`read` authorizes get. `manage` authorizes create, update/replace, and delete.
Either `provision.tools.all` or `provision.all` grants all modifier operations.
