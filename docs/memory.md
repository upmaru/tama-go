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

### Topic Operations

#### GetTopic(id string) (*Topic, error)

Retrieves a specific topic by ID.

**Endpoint:** `GET /provision/memory/topics/:id`

**Parameters:**
- `id` (string): Topic ID (required)

**Returns:**
- `*Topic`: The topic resource (ID, ListenerID, ClassID, ProvisionState)
- `error`: Error if request fails

**Example:**
```go
topic, err := client.Memory.GetTopic("topic-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Topic: %s -> class %s\n", topic.ID, topic.ClassID)
```

#### CreateTopic(listenerID string, req CreateTopicRequest) (*Topic, error)

Creates a new topic under a listener.

**Endpoint:** `POST /provision/memory/listeners/:listener_id/topics`

**Parameters:**
- `listenerID` (string): Listener ID (required)
- `req` (CreateTopicRequest): Topic creation request (required)
  - `Topic` (TopicRequestData): Topic data (required)
    - `ClassID` (string): Class ID for the topic (required)

**Returns:**
- `*Topic`: The created topic resource
- `error`: Error if request fails

**Example:**
```go
req := memory.CreateTopicRequest{
    Topic: memory.TopicRequestData{
        ClassID: "class-456",
    },
}

topic, err := client.Memory.CreateTopic("listener-123", req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created topic: %s\n", topic.ID)
```

#### UpdateTopic(id string, req UpdateTopicRequest) (*Topic, error)

Updates an existing topic using PATCH (partial update).

**Endpoint:** `PATCH /provision/memory/topics/:id`

**Parameters:**
- `id` (string): Topic ID (required)
- `req` (UpdateTopicRequest): Topic update request (required)
  - `Topic` (UpdateTopicData): Topic update data (required)
    - `ClassID` (string): New class ID (required)

**Returns:**
- `*Topic`: The updated topic resource
- `error`: Error if request fails

#### ReplaceTopic(id string, req UpdateTopicRequest) (*Topic, error)

Replaces an existing topic using PUT (full replacement).

**Endpoint:** `PUT /provision/memory/topics/:id`

**Parameters:**
- `id` (string): Topic ID (required)
- `req` (UpdateTopicRequest): Topic replacement request (required)
  - `Topic` (UpdateTopicData): Topic replacement data (required)
    - `ClassID` (string): Class ID (required)

**Returns:**
- `*Topic`: The replaced topic resource
- `error`: Error if request fails

#### DeleteTopic(id string) error

Deletes a topic by ID.

**Endpoint:** `DELETE /provision/memory/topics/:id`

**Parameters:**
- `id` (string): Topic ID (required)

**Returns:**
- `error`: Error if request fails
