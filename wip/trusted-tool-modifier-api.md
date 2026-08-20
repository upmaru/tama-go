# Trusted Tool Modifier API Client Support

## Status

The `tama-go` client described here is implemented in this checkout. The
Terraform resource remains a follow-up and the client has not been released.

The server-side API that introduces `Tama.Tools.Modifier` is merged into Tama's
`develop` branch. Confirm that it is deployed in the target environment before
using the client from `terraform-provider-tama`.

## Goal

Add typed support for the trusted thought-tool modifier provision API to the
existing `tools.Service`. The resulting `tama-go` release must expose the CRUD
operations and error information needed to implement a
`tama_thought_tool_modifier` resource in `terraform-provider-tama`.

The client is only a provisioning transport. Modifier execution, trusted
metadata resolution, JSON Schema compatibility, lifecycle transitions, and
conflict enforcement remain Tama server responsibilities.

## Motivation

Trusted tool modifiers let Tama populate selected action request arguments from
authoritative runtime metadata instead of model-generated text. The first
consumer needs to provision an actor identifier into a nested request field:

```text
/body/search/scope/user_id
```

Terraform needs a stable resource for this configuration. The provider already
uses `tama-go` for thought-tool inputs, outputs, initializers, and options, so
the modifier resource should use the same boundary rather than constructing
HTTP requests inside the provider.

## Authoritative Tama API contract

All endpoints use the authenticated `/provision` API pipeline.

| Operation | Method and path | Success | Notes |
| --- | --- | --- | --- |
| Create or reactivate | `POST /provision/tools/:thought_tool_id/modifiers` | `201` | Creates an inactive record and activates it. An exact inactive configuration is reactivated with the same ID. |
| Show | `GET /provision/tools/modifiers/:id` | `200` | Only active modifiers are visible. Invalid, unknown, and inactive IDs return `404`. |
| Update | `PATCH /provision/tools/modifiers/:id` | `200` | Can update active or inactive records. `index` is immutable. Active updates revalidate compatibility and conflicts. |
| Replace fields | `PUT /provision/tools/modifiers/:id` | `200` | Phoenix routes PUT to the same update action as PATCH. It accepts the same mutable fields and cannot replace `index`. |
| Deactivate | `DELETE /provision/tools/modifiers/:id` | `200` | Deactivates the active record and returns it with `provision_state = "inactive"`. A later show returns `404`. |

There is no list endpoint. Do not add `ListModifier` unless Tama adds that
route. The generated Phoenix resource routes do include PUT, matching adjacent
Tools resources and their `Replace*` client methods.

### Authorization

A Terraform credential that manages the complete resource lifecycle needs
read and manage access. Any one of the broad scopes below covers all actions:

```text
provision.all
provision.tools.all
```

The narrow scopes are:

```text
provision.tools.modifier.read
provision.tools.modifier.manage
```

`read` authorizes show. `manage` authorizes create, update, and delete; it does
not authorize show by itself. A narrowly scoped provider credential therefore
needs both capabilities.

### Request shape

Create wraps the complete configuration under `modifier`:

```json
{
  "modifier": {
    "index": 0,
    "target": "/body/search/scope/user_id",
    "on_missing_parent": "skip",
    "on_missing_source": "error",
    "source": {
      "type": "metadata",
      "path": "actor_identifier"
    }
  }
}
```

The create fields are all required. `provision_state`, `id`, and
`thought_tool_id` are server-owned and must not be sent in the body.

Update uses the same envelope with only mutable fields:

```json
{
  "modifier": {
    "target": "/body/search/scope/user_id",
    "on_missing_parent": "skip",
    "on_missing_source": "error",
    "source": {
      "type": "metadata",
      "path": "actor_identifier"
    }
  }
}
```

The client update type should omit `index`. Tama rejects changing it, and the
Terraform resource must model an index change as replacement.

### Response shape

Every successful endpoint returns the resource in a `data` envelope:

```json
{
  "data": {
    "id": "019c...",
    "index": 0,
    "target": "/body/search/scope/user_id",
    "on_missing_parent": "skip",
    "on_missing_source": "error",
    "source": {
      "type": "metadata",
      "path": "actor_identifier"
    },
    "thought_tool_id": "019c...",
    "provision_state": "active"
  }
}
```

The delete response has the same shape with `provision_state` set to
`inactive`. The client delete method discards this body, consistent with the
existing Tools CRUD methods; Terraform only needs to know whether deactivation
succeeded.

### Allowed values and server validation

Version 1 allows these values:

```text
on_missing_parent: error | skip
on_missing_source: error | skip
source.type: metadata
source.path: actor_identifier | origin_entity_identifier | current_timestamp
```

The server also enforces:

- `index` is non-negative, immutable after creation, and unique among active
  modifiers owned by one thought tool;
- `target` is an RFC 6901 JSON Pointer encoded in at most 512 bytes;
- the decoded target contains 2 through 16 non-empty tokens;
- the first token is `body`, `path`, or `query`;
- traversal is limited to ordinary object properties in version 1;
- array, `oneOf`, `anyOf`, `allOf`, and dynamic `additionalProperties`
  traversal is unsupported;
- the target must resolve in at least one complete callable schema owned by the
  thought tool; and
- active targets cannot be equal or have an ancestor/descendant relationship.

Do not duplicate callable-schema resolution or active conflict detection in
`tama-go`. Those checks depend on current server state and imported action
specifications. The client may reject structurally empty required fields and
negative indexes, but Tama remains authoritative and returns field errors for
invalid targets or configurations.

## Proposed `tama-go` public API

Add `tools/modifier.go`. The resource belongs in the existing `tools` package,
not `motor`; `motor.Modifier` is a different action-level schema modifier.

Use string fields, matching adjacent `tama-go` resources, and export constants
to prevent consumers from repeating wire literals:

```go
const (
    ModifierMissingPolicyError = "error"
    ModifierMissingPolicySkip  = "skip"

    ModifierSourceTypeMetadata = "metadata"

    ModifierSourcePathActorIdentifier        = "actor_identifier"
    ModifierSourcePathOriginEntityIdentifier = "origin_entity_identifier"
    ModifierSourcePathCurrentTimestamp       = "current_timestamp"
)

type ModifierSource struct {
    Type string `json:"type"`
    Path string `json:"path"`
}

type Modifier struct {
    ID              string         `json:"id,omitempty"`
    Index           int            `json:"index"`
    Target          string         `json:"target"`
    OnMissingParent string         `json:"on_missing_parent"`
    OnMissingSource string         `json:"on_missing_source"`
    Source          ModifierSource `json:"source"`
    ThoughtToolID   string         `json:"thought_tool_id,omitempty"`
    ProvisionState  string         `json:"provision_state"`
}

type ModifierResponse struct {
    Data Modifier `json:"data"`
}

type CreateModifierRequest struct {
    Modifier ModifierRequestData `json:"modifier"`
}

type ModifierRequestData struct {
    Index           int            `json:"index"`
    Target          string         `json:"target"`
    OnMissingParent string         `json:"on_missing_parent"`
    OnMissingSource string         `json:"on_missing_source"`
    Source          ModifierSource `json:"source"`
}

type UpdateModifierRequest struct {
    Modifier UpdateModifierData `json:"modifier"`
}

type UpdateModifierData struct {
    Target          string          `json:"target,omitempty"`
    OnMissingParent string          `json:"on_missing_parent,omitempty"`
    OnMissingSource string          `json:"on_missing_source,omitempty"`
    Source          *ModifierSource `json:"source,omitempty"`
}
```

`Source` is a pointer only on update so an omitted source is distinguishable
from a supplied source. If it is supplied, send both `type` and `path`.

Expose these methods on `*tools.Service`, which automatically makes them
available through `client.Tools`:

```go
func (s *Service) GetModifier(id string) (*Modifier, error)

func (s *Service) CreateModifier(
    thoughtToolID string,
    req CreateModifierRequest,
) (*Modifier, error)

func (s *Service) UpdateModifier(
    id string,
    req UpdateModifierRequest,
) (*Modifier, error)

func (s *Service) ReplaceModifier(
    id string,
    req UpdateModifierRequest,
) (*Modifier, error)

func (s *Service) DeleteModifier(id string) error
```

Method behavior should follow `tools/input.go`, `tools/output.go`, and
`tools/initializer.go`:

1. reject an empty resource ID before making an HTTP request;
2. reject an empty thought-tool ID on create;
3. validate the complete create payload's required fields and non-negative
   index;
4. preserve partial PATCH and PUT semantics, validating both nested source
   fields only when `source` is supplied;
5. use `SetBody`, `SetResult`, and the shared Resty client;
6. call `handleAPIError` for every non-success response;
7. wrap transport failures with operation-specific context; and
8. return the decoded `data` value for get, create, and update.

An update or replacement with no supplied fields is sent as
`{"modifier": {}}`, matching the existing Tools convention and Tama's no-op
update behavior. `ReplaceModifier` shares `UpdateModifierRequest` because PUT
routes to the same server action and still cannot change `index`.

Do not add a second service constructor or a new root client field. The current
`ToolsService` embeds `*tools.Service`, so methods added to that service are
already promoted to `client.Tools`.

### Error contract needed by Terraform

The existing `tools.Error` must remain the returned error for parsed API
failures. In particular, the provider needs to detect an inactive or externally
deleted modifier during refresh:

```go
var apiErr *tools.Error
if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
    // Remove the Terraform resource from state.
}
```

Validation failures should retain Tama's flattened field keys and messages,
for example `target`, `index`, `source.type`, or `source.path`. Do not replace a
server `422` with an untyped formatted string.

## `tama-go` test plan

Add `tools/modifier_test.go` using the existing `httptest` helpers. Cover:

- get uses `GET /provision/tools/modifiers/:id` and decodes every response
  field, including the nested source;
- create uses `POST /provision/tools/:thought_tool_id/modifiers`, preserves
  `index: 0`, sends the complete envelope, accepts `201`, and decodes the active
  resource; inspect the raw JSON object so omission cannot be mistaken for a
  decoded zero value;
- update uses `PATCH /provision/tools/modifiers/:id`, omits unset fields, and
  sends a complete nested source when present; inspect the raw JSON object to
  prove absent fields and the immutable index are not serialized;
- replace uses `PUT /provision/tools/modifiers/:id`, shares the mutable update
  request, and cannot serialize the immutable index;
- delete uses `DELETE /provision/tools/modifiers/:id` and treats the `200`
  inactive response as success;
- empty IDs, empty thought-tool IDs, negative indexes, missing targets,
  missing policies, and incomplete sources fail locally without an HTTP call;
- Tama's `{"errors":{"detail":"Not Found"}}` response becomes `*tools.Error`
  with `StatusCode == 404`;
- a real nested changeset shape such as
  `{"errors":{"source":{"path":["can't be blank"]}}}` preserves the
  `source.path` field and its messages; and
- transport errors retain get/create/update/replace/delete operation context.

Run:

```bash
go test ./tools
make check
go build ./...
go test -race ./...
go mod tidy
git diff --check
```

The final diff after `go mod tidy` must not contain unintended module changes.
A live integration test is optional and must use a disposable thought tool
whose action schema contains the chosen target; unit tests must not depend on
deployed API state.

Update the Tools service comments, `docs/tools.md`, the documentation index, and
the root README so the public client surface includes modifiers. Remove stale
examples that use nonexistent nested services, contexts, or list endpoints.
Document that Tools `Replace*` methods route PUT to the same update actions as
PATCH and therefore accept only their mutable request fields.

## Terraform provider handoff

After a released `tama-go` version contains the client methods, implement a
`tama_thought_tool_modifier` resource in `terraform-provider-tama`. Do not use
the existing `tama_motor_modifier` resource as the API model; it represents a
different resource.

The expected Terraform model is:

| Attribute | Mode | Lifecycle |
| --- | --- | --- |
| `id` | computed | Server-generated modifier ID. |
| `thought_tool_id` | required | `RequiresReplace`; ownership cannot move. |
| `index` | required integer | Non-negative and `RequiresReplace`; Tama makes it immutable. |
| `target` | required string | Mutable; server validates pointer syntax and callable compatibility. |
| `on_missing_parent` | required string | Mutable; validate `error` or `skip`. |
| `on_missing_source` | required string | Mutable; validate `error` or `skip`. |
| `source` | required single nested object | Mutable; contains required `type` and `path`. |
| `provision_state` | computed | Expected to be `active` while managed. |

The nested source validators should allow only `metadata` for `type` and the
three version 1 metadata paths. If the provider validates the target length, it
must count UTF-8 bytes rather than characters; otherwise leave the 512-byte
check to Tama to avoid a mismatched rule.

Provider lifecycle requirements:

- Create maps the Terraform plan to `tools.CreateModifierRequest` and accepts
  the ID returned by create or exact reactivation.
- Read maps all server fields back into state. A typed `404` must call
  `resp.State.RemoveResource(ctx)` instead of producing a permanent diagnostic.
- Update sends `target`, both missing-value policies, and the complete source.
  It must not send `index`.
- Delete calls `DeleteModifier`; Terraform then removes the resource from
  state. Deactivation, rather than physical deletion, is the intended server
  behavior.
- Import accepts the modifier ID and reads the active resource. Importing an
  inactive ID must fail as not found.

Register the resource in `tama/provider.go`, add resource and acceptance tests,
and generate documentation and examples following the adjacent thought-tool
resources. Acceptance coverage should include create, update, import, delete,
exact reactivation, a negative index validation, an invalid target API error,
and refresh after external deactivation.

During local provider development, a temporary Go module replacement may point
`github.com/upmaru/tama-go` at `../tama-go`. Remove that replacement before
release. The provider's pinned `tama-go` version must ultimately be advanced to
the released version containing this API.

## Sequencing

1. Merge and deploy the Tama modifier API where integration or acceptance tests
   will run.
2. Implement and validate `tools/modifier.go` and its unit tests in `tama-go`.
3. Update client documentation and examples.
4. Release `tama-go`; do not assume a version number before checking registry
   and repository state.
5. Update `terraform-provider-tama` to the released client version.
6. Implement, register, document, and acceptance-test
   `tama_thought_tool_modifier`.
7. Release the provider only after its generated documentation and full test
   gates pass.

## Non-goals

This client phase does not:

- implement modifier execution or trusted metadata resolution;
- inspect action schemas locally;
- add a list endpoint that Tama does not expose;
- implement the Terraform resource in `tama-go`;
- deploy the modifier to any Memovee graph; or
- remove prompt-based identifier instructions before the modifier is
  provisioned in the target environment.

## Completion criteria

The `tama-go` phase is complete only when:

- all request, response, resource, source, and constant types are public;
- get, create, update, replace, and delete call the exact Tama routes;
- zero is preserved as a valid required create index;
- update cannot accidentally request an index change;
- nested source values round-trip without a generic `map[string]any`;
- `404` and `422` responses remain typed `*tools.Error` values;
- modifier unit tests cover methods, payloads, validation, and API errors;
- the full `tama-go` quality gate passes; and
- public Tools documentation reflects the implemented surface.

Provider work may begin against a local client checkout after these contracts
pass, but a provider release must depend on a released `tama-go` version.
