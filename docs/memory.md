## Memory Service

Access via `client.Memory.*`

### Prompt Operations

#### GetPrompt(id string) (*Prompt, error)

Retrieves a specific prompt by ID.

**Parameters:**
- `id` (string): Prompt ID (required)

**Returns:**
- `*Prompt`: The prompt resource
- `error`: Error if request fails

**Example:**
```go
prompt, err := client.Memory.GetPrompt("prompt-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Prompt: %s - %s\n", prompt.Name, prompt.Role)
```

#### CreatePrompt(spaceID string, req CreatePromptRequest) (*Prompt, error)

Creates a new prompt in a specific space.

**Parameters:**
- `spaceID` (string): Space ID where the prompt will be created (required)
- `req` (CreatePromptRequest): Prompt creation request (required)

**Returns:**
- `*Prompt`: The created prompt resource
- `error`: Error if request fails

**Example:**
```go
req := memory.CreatePromptRequest{
    Prompt: memory.PromptRequestData{
        Name:    "Assistant Prompt",
        Content: "You are a helpful assistant.",
        Role:    "system",
    },
}

prompt, err := client.Memory.CreatePrompt("space-123", req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created prompt: %s\n", prompt.ID)
```

#### UpdatePrompt(id string, req UpdatePromptRequest) (*Prompt, error)

Updates an existing prompt using PATCH (partial update).

**Parameters:**
- `id` (string): Prompt ID (required)
- `req` (UpdatePromptRequest): Prompt update request (required)

**Returns:**
- `*Prompt`: The updated prompt resource
- `error`: Error if request fails

#### ReplacePrompt(id string, req UpdatePromptRequest) (*Prompt, error)

Replaces an existing prompt using PUT (full replacement).

**Parameters:**
- `id` (string): Prompt ID (required)
- `req` (UpdatePromptRequest): Prompt replacement request (required)

**Returns:**
- `*Prompt`: The replaced prompt resource
- `error`: Error if request fails

#### DeletePrompt(id string) error

Deletes a prompt by ID.

**Parameters:**
- `id` (string): Prompt ID (required)

**Returns:**
- `error`: Error if request fails
