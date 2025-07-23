package main

import (
	"errors"
	"log"
	"strings"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/memory"
	"github.com/upmaru/tama-go/neural"
	"github.com/upmaru/tama-go/perception"
	"github.com/upmaru/tama-go/sensory"
)

const (
	exampleSpaceID    = "space-123"
	exampleSourceID   = "source-123"
	exampleModelID    = "model-123"
	exampleLimitID    = "limit-123"
	examplePromptID   = "prompt-123"
	exampleChainID    = "chain-123"
	exampleThoughtID  = "thought-123"
	defaultTimeout    = 30
	defaultLimitCount = 32
	scaleCountValue   = 5
	limitCountValue   = 100
	// AI model parameters.
	defaultTemperature = 0.7
	defaultMaxTokens   = 150
	defaultDepth       = 3
)

func main() {
	client := initializeClient()

	// Run examples in separate functions to reduce complexity
	// Run examples
	runNeuralSpaceOperations(client)
	runNeuralClassOperations(client)
	runMemoryPromptOperations(client)
	runSensorySourceOperations(client)
	runSensoryModelOperations(client)
	runSensoryLimitOperations(client)
	runPerceptionChainOperations(client)
	runPerceptionThoughtOperations(client)

	// Demonstrate enhanced error handling
	demonstrateErrorHandling(client)
	runDeleteOperations(client)

	log.Printf("Example completed!")
}

// initializeClient creates and configures the Tama client.
func initializeClient() *tama.Client {
	config := tama.Config{
		BaseURL: "http://localhost:4000", // Local development server
		APIKey:  "your-api-key",          // Replace with your actual API key
		Timeout: defaultTimeout * time.Second,
	}

	client := tama.NewClient(config)
	client.SetDebug(true) // Enable debug mode to see HTTP requests/responses (optional)
	return client
}

// runNeuralSpaceOperations demonstrates neural space operations.
func runNeuralSpaceOperations(client *tama.Client) {
	log.Printf("=== Neural Space Operations ===")

	// Create a new space
	newSpace := neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "My Neural Space",
			Type: "root",
		},
	}

	space, err := client.Neural.CreateSpace(newSpace)
	if err != nil {
		log.Printf("Error creating space: %v", err)
	} else {
		log.Printf("Created space: ID=%s, Name=%s, Type=%s, State=%s",
			space.ID, space.Name, space.Type, space.CurrentState)
	}

	// Get a space by ID (replace with actual ID)
	spaceID := exampleSpaceID
	space, err = client.Neural.GetSpace(spaceID)
	if err != nil {
		log.Printf("Error getting space: %v", err)
	} else {
		log.Printf("Retrieved space: ID=%s, Name=%s, Type=%s, State=%s",
			space.ID, space.Name, space.Type, space.CurrentState)
	}

	// Update a space
	updateSpace := neural.UpdateSpaceRequest{
		Space: neural.UpdateSpaceData{
			Name: "Updated Neural Space",
			Type: "component",
		},
	}

	space, err = client.Neural.UpdateSpace(spaceID, updateSpace)
	if err != nil {
		log.Printf("Error updating space: %v", err)
	} else {
		log.Printf("Updated space: %+v", space)
	}
}

// runNeuralClassOperations demonstrates neural class operations.
func runNeuralClassOperations(client *tama.Client) {
	log.Printf("=== Neural Class Operations ===")

	spaceID := exampleSpaceID
	classID := "class-123"

	// Create, Get, Update, and Replace operations
	demoCreateClass(client, spaceID)
	demoGetClass(client, classID)
	demoUpdateClass(client, classID)
	demoReplaceClass(client, classID)
}

// demoCreateClass demonstrates creating a class with a real-world schema.
func demoCreateClass(client *tama.Client, spaceID string) {
	newClass := neural.CreateClassRequest{
		Class: neural.ClassRequestData{
			Schema: createActionCallSchema(),
		},
	}

	class, err := client.Neural.CreateClass(spaceID, newClass)
	if err != nil {
		log.Printf("Error creating class: %v", err)
	} else {
		log.Printf("Created class: ID=%s, SpaceID=%s, Name=%s, State=%s",
			class.ID, class.SpaceID, class.Name, class.CurrentState)
		log.Printf("Description: %s", class.Description)
		if title, ok := class.Schema["title"].(string); ok {
			log.Printf("Schema title: %s", title)
		}
	}
}

// demoGetClass demonstrates retrieving a class by ID.
func demoGetClass(client *tama.Client, classID string) {
	class, err := client.Neural.GetClass(classID)
	if err != nil {
		log.Printf("Error getting class: %v", err)
	} else {
		log.Printf("Retrieved class: ID=%s, SpaceID=%s, Name=%s, State=%s",
			class.ID, class.SpaceID, class.Name, class.CurrentState)
		log.Printf("Description: %s", class.Description)
	}
}

// demoUpdateClass demonstrates updating a class with an enhanced schema.
func demoUpdateClass(client *tama.Client, classID string) {
	updateClass := neural.UpdateClassRequest{
		Class: neural.UpdateClassData{
			Schema: createEnhancedActionCallSchema(),
		},
	}

	class, err := client.Neural.UpdateClass(classID, updateClass)
	if err != nil {
		log.Printf("Error updating class: %v", err)
	} else {
		log.Printf("Updated class: ID=%s, Name=%s", class.ID, class.Name)
		if desc, ok := class.Schema["description"].(string); ok {
			log.Printf("Updated schema description: %s", desc)
		}
	}
}

// demoReplaceClass demonstrates replacing a class with a completely new schema.
func demoReplaceClass(client *tama.Client, classID string) {
	replaceClass := neural.UpdateClassRequest{
		Class: neural.UpdateClassData{
			Schema: createSimpleMessageSchema(),
		},
	}

	class, err := client.Neural.ReplaceClass(classID, replaceClass)
	if err != nil {
		log.Printf("Error replacing class: %v", err)
	} else {
		log.Printf("Replaced class: ID=%s, Name=%s", class.ID, class.Name)
		if title, ok := class.Schema["title"].(string); ok {
			log.Printf("New schema title: %s", title)
		}
	}
}

// createActionCallSchema creates the action-call schema.
func createActionCallSchema() map[string]any {
	return map[string]any{
		"title":       "action-call",
		"description": "An action call is a request to execute an action.",
		"type":        "object",
		"properties": map[string]any{
			"code": map[string]any{
				"description": "The status of the action call",
				"type":        "integer",
			},
			"tool_id": map[string]any{
				"description": "The ID of the tool to execute",
				"type":        "string",
			},
			"parameters": map[string]any{
				"description": "The parameters to pass to the action",
				"type":        "object",
			},
			"content_type": map[string]any{
				"description": "The content type of the response",
				"type":        "string",
			},
			"content": map[string]any{
				"description": "The response from the action",
				"type":        "object",
			},
		},
		"required": []any{"tool_id", "parameters", "code", "content_type", "content"},
	}
}

// createEnhancedActionCallSchema creates an enhanced action-call schema with timestamp.
func createEnhancedActionCallSchema() map[string]any {
	schema := createActionCallSchema()
	schema["description"] = "An updated action call schema with additional fields."

	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		return schema
	}
	properties["timestamp"] = map[string]any{
		"description": "When the action was called",
		"type":        "string",
		"format":      "date-time",
	}

	schema["required"] = []any{"tool_id", "parameters", "code", "content_type", "content", "timestamp"}
	return schema
}

// createSimpleMessageSchema creates a simple message schema.
func createSimpleMessageSchema() map[string]any {
	return map[string]any{
		"title":       "simple-message",
		"description": "A simple message schema",
		"type":        "object",
		"properties": map[string]any{
			"text": map[string]any{
				"description": "The message text",
				"type":        "string",
			},
			"sender": map[string]any{
				"description": "Who sent the message",
				"type":        "string",
			},
		},
		"required": []any{"text", "sender"},
	}
}

// runMemoryPromptOperations demonstrates memory prompt operations.
func runMemoryPromptOperations(client *tama.Client) {
	log.Printf("=== Memory Prompt Operations ===")

	spaceID := exampleSpaceID

	// Create a new prompt
	newPrompt := memory.CreatePromptRequest{
		Prompt: memory.PromptRequestData{
			Name:    "Assistant Helper",
			Content: "You are a helpful assistant that provides clear and concise answers.",
			Role:    "system",
		},
	}

	prompt, err := client.Memory.CreatePrompt(spaceID, newPrompt)
	if err != nil {
		log.Printf("Error creating prompt: %v", err)
	} else {
		log.Printf("Created prompt: ID=%s, Name=%s, Role=%s, SpaceID=%s, State=%s",
			prompt.ID, prompt.Name, prompt.Role, prompt.SpaceID, prompt.CurrentState)
		log.Printf("Content: %s", prompt.Content)
	}

	// Get a prompt by ID (replace with actual ID)
	promptID := examplePromptID
	prompt, err = client.Memory.GetPrompt(promptID)
	if err != nil {
		log.Printf("Error getting prompt: %v", err)
	} else {
		log.Printf("Retrieved prompt: ID=%s, Name=%s, Role=%s, SpaceID=%s, State=%s",
			prompt.ID, prompt.Name, prompt.Role, prompt.SpaceID, prompt.CurrentState)
		log.Printf("Slug: %s, Content: %s", prompt.Slug, prompt.Content)
	}

	// Update a prompt
	updatePrompt := memory.UpdatePromptRequest{
		Prompt: memory.UpdatePromptData{
			Name:    "Updated Assistant Helper",
			Content: "You are an updated helpful assistant that provides detailed explanations.",
		},
	}

	prompt, err = client.Memory.UpdatePrompt(promptID, updatePrompt)
	if err != nil {
		log.Printf("Error updating prompt: %v", err)
	} else {
		log.Printf("Updated prompt: ID=%s, Name=%s, Content=%s",
			prompt.ID, prompt.Name, prompt.Content)
	}

	// Replace a prompt (full replacement)
	replacePrompt := memory.UpdatePromptRequest{
		Prompt: memory.UpdatePromptData{
			Name:    "Completely New Assistant",
			Content: "You are a completely new assistant with different capabilities.",
			Role:    "assistant",
		},
	}

	prompt, err = client.Memory.ReplacePrompt(promptID, replacePrompt)
	if err != nil {
		log.Printf("Error replacing prompt: %v", err)
	} else {
		log.Printf("Replaced prompt: ID=%s, Name=%s, Role=%s, Content=%s",
			prompt.ID, prompt.Name, prompt.Role, prompt.Content)
	}
}

// runSensorySourceOperations demonstrates sensory source operations.
func runSensorySourceOperations(client *tama.Client) {
	log.Printf("=== Sensory Source Operations ===")

	spaceID := exampleSpaceID

	// Create a new source
	newSource := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "My Data Source",
			Type:     "model",
			Endpoint: "https://api.mistral.ai/v1",
			Credential: sensory.SourceCredential{
				APIKey: "your-api-key-here",
			},
		},
	}

	source, err := client.Sensory.CreateSource(spaceID, newSource)
	if err != nil {
		log.Printf("Error creating source: %v", err)
	} else {
		log.Printf("Created source: ID=%s, Name=%s, Endpoint=%s, SpaceID=%s, State=%s",
			source.ID, source.Name, source.Endpoint, source.SpaceID, source.CurrentState)
	}

	// Get a source by ID (replace with actual ID)
	sourceID := exampleSourceID
	source, err = client.Sensory.GetSource(sourceID)
	if err != nil {
		log.Printf("Error getting source: %v", err)
	} else {
		log.Printf("Retrieved source: ID=%s, Name=%s, Endpoint=%s, SpaceID=%s, State=%s",
			source.ID, source.Name, source.Endpoint, source.SpaceID, source.CurrentState)
	}

	// Update a source
	updateSource := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name:     "Updated Data Source",
			Type:     "model",
			Endpoint: "https://api.openai.com/v1",
			Credential: &sensory.SourceCredential{
				APIKey: "your-updated-api-key",
			},
		},
	}

	source, err = client.Sensory.UpdateSource(sourceID, updateSource)
	if err != nil {
		log.Printf("Error updating source: %v", err)
	} else {
		log.Printf("Updated source: ID=%s, Name=%s, Endpoint=%s, SpaceID=%s, State=%s",
			source.ID, source.Name, source.Endpoint, source.SpaceID, source.CurrentState)
	}
}

// runSensoryModelOperations demonstrates sensory model operations.
func runSensoryModelOperations(client *tama.Client) {
	log.Printf("=== Sensory Model Operations ===")

	sourceID := exampleSourceID

	// Create a new model
	newModel := sensory.CreateModelRequest{
		Model: sensory.ModelRequestData{
			Identifier: "mistral-small-latest",
			Path:       "/chat/completions",
		},
	}

	model, err := client.Sensory.CreateModel(sourceID, newModel)
	if err != nil {
		log.Printf("Error creating model: %v", err)
	} else {
		log.Printf("Created model: %+v", model)
	}

	// Get a model by ID (replace with actual ID)
	modelID := exampleModelID
	model, err = client.Sensory.GetModel(modelID)
	if err != nil {
		log.Printf("Error getting model: %v", err)
	} else {
		log.Printf("Retrieved model: %+v", model)
	}

	// Update a model
	updateModel := sensory.UpdateModelRequest{
		Model: sensory.UpdateModelData{
			Identifier: "mistral-large-latest",
			Path:       "/chat/completions",
		},
	}

	model, err = client.Sensory.UpdateModel(modelID, updateModel)
	if err != nil {
		log.Printf("Error updating model: %v", err)
	} else {
		log.Printf("Updated model: %+v", model)
	}
}

// runSensoryLimitOperations demonstrates sensory limit operations.
func runSensoryLimitOperations(client *tama.Client) {
	log.Printf("=== Sensory Limit Operations ===")

	sourceID := exampleSourceID

	// Create a new limit
	newLimit := sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			ScaleUnit:  "seconds",
			ScaleCount: 1,
			Count:      defaultLimitCount,
		},
	}

	limit, err := client.Sensory.CreateLimit(sourceID, newLimit)
	if err != nil {
		log.Printf("Error creating limit: %v", err)
	} else {
		log.Printf("Created limit: %+v", limit)
	}

	// Get a limit by ID (replace with actual ID)
	limitID := exampleLimitID
	limit, err = client.Sensory.GetLimit(limitID)
	if err != nil {
		log.Printf("Error getting limit: %v", err)
	} else {
		log.Printf("Retrieved limit: %+v", limit)
	}

	// Update a limit
	updateLimit := sensory.UpdateLimitRequest{
		Limit: sensory.UpdateLimitData{
			ScaleUnit:    "minutes",
			ScaleCount:   scaleCountValue,
			Count:        limitCountValue,
			CurrentState: "active",
		},
	}

	limit, err = client.Sensory.UpdateLimit(limitID, updateLimit)
	if err != nil {
		log.Printf("Error updating limit: %v", err)
	} else {
		log.Printf("Updated limit: %+v", limit)
	}
}

// demonstrateErrorHandling shows examples of the enhanced error handling
// for both general API errors and field-specific validation errors.
func demonstrateErrorHandling(client *tama.Client) {
	log.Printf("=== Enhanced Error Handling Examples ===")

	// Example 1: Field validation errors
	log.Printf("--- Example 1: Field Validation Errors ---")
	invalidSource := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "",             // Invalid: empty name
			Type:     "invalid-type", // Invalid: unsupported type
			Endpoint: "not-a-url",    // Invalid: malformed URL
			Credential: sensory.SourceCredential{
				APIKey: "", // Invalid: empty API key
			},
		},
	}

	_, err := client.Sensory.CreateSource("invalid-space-id", invalidSource)
	if err != nil {
		handleEnhancedError("CreateSource", err)
	}

	// Example 2: General API error (resource not found)
	log.Printf("--- Example 2: General API Error ---")
	_, err = client.Sensory.GetSource("non-existent-source-id")
	if err != nil {
		handleEnhancedError("GetSource", err)
	}

	// Example 3: Neural service field validation
	log.Printf("--- Example 3: Neural Service Validation ---")
	invalidSpace := neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "",        // Invalid: empty name
			Type: "invalid", // Invalid: unsupported type
		},
	}

	_, err = client.Neural.CreateSpace(invalidSpace)
	if err != nil {
		handleEnhancedError("CreateSpace", err)
	}

	// Example 3a: Neural class validation
	log.Printf("--- Example 3a: Neural Class Validation ---")
	invalidClass := neural.CreateClassRequest{
		Class: neural.ClassRequestData{
			Schema: nil, // Invalid: missing schema
		},
	}

	_, err = client.Neural.CreateClass("invalid-space-id", invalidClass)
	if err != nil {
		handleEnhancedError("CreateClass", err)
	}

	// Example 4: Memory service field validation
	log.Printf("--- Example 4: Memory Service Validation ---")
	invalidPrompt := memory.CreatePromptRequest{
		Prompt: memory.PromptRequestData{
			Name:    "",        // Invalid: empty name
			Content: "Short",   // Invalid: too short
			Role:    "invalid", // Invalid: unsupported role
		},
	}

	_, err = client.Memory.CreatePrompt("invalid-space-id", invalidPrompt)
	if err != nil {
		handleEnhancedError("CreatePrompt", err)
	}
}

// handleEnhancedError demonstrates comprehensive error handling
// for the new API error structure with field-specific validation.
func handleEnhancedError(operation string, err error) {
	log.Printf("Error in %s operation:", operation)

	// Check if it's a sensory API error
	var sensoryErr *sensory.Error
	if errors.As(err, &sensoryErr) {
		handleAPIError(sensoryErr.StatusCode, sensoryErr.Errors)
		return
	}

	// Check if it's a neural API error
	var neuralErr *neural.Error
	if errors.As(err, &neuralErr) {
		handleAPIError(neuralErr.StatusCode, neuralErr.Errors)
		return
	}

	// Check if it's a memory API error
	var memoryErr *memory.Error
	if errors.As(err, &memoryErr) {
		handleAPIError(memoryErr.StatusCode, memoryErr.Errors)
		return
	}

	// Handle client/network errors
	log.Printf("  Client/Network Error: %v", err)
}

// handleAPIError processes API errors with field validation support.
func handleAPIError(statusCode int, errors map[string][]string) {
	if len(errors) > 0 {
		// Handle field validation errors
		log.Printf("  Field validation errors (Status: %d):", statusCode)
		for field, messages := range errors {
			log.Printf("    %s: %s", field, strings.Join(messages, ", "))
		}
	} else {
		// Handle general API errors
		log.Printf("  API Error %d", statusCode)
	}
}

// runDeleteOperations demonstrates delete operations.
func runDeleteOperations(_ *tama.Client) {
	log.Printf("=== Delete Operations ===")

	// Delete resources (uncomment to test)
	/*
		promptID := examplePromptID
		limitID := exampleLimitID
		modelID := exampleModelID
		sourceID := exampleSourceID
		spaceID := exampleSpaceID
		classID := "class-123"

		err := client.Memory.DeletePrompt(promptID)
		if err != nil {
			log.Printf("Error deleting prompt: %v", err)
		} else {
			log.Printf("Deleted prompt successfully")
		}

		err = client.Sensory.DeleteLimit(limitID)
		if err != nil {
			log.Printf("Error deleting limit: %v", err)
		} else {
			log.Printf("Deleted limit successfully")
		}

		err = client.Sensory.DeleteModel(modelID)
		if err != nil {
			log.Printf("Error deleting model: %v", err)
		} else {
			log.Printf("Deleted model successfully")
		}

		err = client.Sensory.DeleteSource(sourceID)
		if err != nil {
			log.Printf("Error deleting source: %v", err)
		} else {
			log.Printf("Deleted source successfully")
		}

		err = client.Neural.DeleteClass(classID)
		if err != nil {
			log.Printf("Error deleting class: %v", err)
		} else {
			log.Printf("Deleted class successfully")
		}

		err = client.Neural.DeleteSpace(spaceID)
		if err != nil {
			log.Printf("Error deleting space: %v", err)
		} else {
			log.Printf("Deleted space successfully")
		}
	*/
}

// runPerceptionChainOperations demonstrates perception chain operations.
func runPerceptionChainOperations(client *tama.Client) {
	log.Printf("--- Perception Chain Operations ---")

	// Create a chain
	createChain := perception.CreateChainRequest{
		Chain: perception.ChainRequestData{
			Name: "Example Processing Chain",
		},
	}

	chain, err := client.Perception.CreateChain(exampleSpaceID, createChain)
	if err != nil {
		log.Printf("Error creating chain: %v", err)
		return
	}

	log.Printf("Created chain: ID=%s, Name=%s, State=%s",
		chain.ID, chain.Name, chain.CurrentState)

	// Get the chain
	retrievedChain, err := client.Perception.GetChain(chain.ID)
	if err != nil {
		log.Printf("Error getting chain: %v", err)
	} else {
		log.Printf("Retrieved chain: ID=%s, SpaceID=%s, Name=%s, Slug=%s, State=%s",
			retrievedChain.ID, retrievedChain.SpaceID, retrievedChain.Name,
			retrievedChain.Slug, retrievedChain.CurrentState)
	}

	// Update the chain
	updateChain := perception.UpdateChainRequest{
		Chain: perception.UpdateChainData{
			Name: "Updated Processing Chain",
		},
	}

	updatedChain, err := client.Perception.UpdateChain(chain.ID, updateChain)
	if err != nil {
		log.Printf("Error updating chain: %v", err)
	} else {
		log.Printf("Updated chain: ID=%s, Name=%s", updatedChain.ID, updatedChain.Name)
	}

	log.Printf("Chain operations completed successfully")
}

// runPerceptionThoughtOperations demonstrates perception thought operations.
func runPerceptionThoughtOperations(client *tama.Client) {
	log.Printf("--- Perception Thought Operations ---")

	// Create a thought
	createThought := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Module: perception.Module{
				Reference: "tama/agentic/generate",
				Parameters: map[string]any{
					"temperature": defaultTemperature,
					"max_tokens":  defaultMaxTokens,
					"model":       "gpt-4",
				},
			},
		},
	}

	thought, err := client.Perception.CreateThought(exampleChainID, createThought)
	if err != nil {
		log.Printf("Error creating thought: %v", err)
		return
	}

	log.Printf("Created thought: ID=%s, ChainID=%s, Relation=%s, State=%s, Index=%d",
		thought.ID, thought.ChainID, thought.Relation, thought.CurrentState, thought.Index)
	log.Printf("Module: Reference=%s, ID=%s",
		thought.Module.Reference, thought.Module.ID)

	// Get the thought
	retrievedThought, err := client.Perception.GetThought(thought.ID)
	if err != nil {
		log.Printf("Error getting thought: %v", err)
	} else {
		log.Printf("Retrieved thought: ID=%s, ChainID=%s, OutputClassID=%s",
			retrievedThought.ID, retrievedThought.ChainID, retrievedThought.OutputClassID)
		log.Printf("Retrieved thought module: Reference=%s, Parameters=%v",
			retrievedThought.Module.Reference, retrievedThought.Module.Parameters)
	}

	// Update the thought
	updateThought := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation: "analysis",
			Module: perception.Module{
				Reference: "tama/agentic/analyze",
				Parameters: map[string]any{
					"depth":       defaultDepth,
					"focus_areas": []string{"sentiment", "intent", "entities"},
				},
			},
		},
	}

	updatedThought, err := client.Perception.UpdateThought(thought.ID, updateThought)
	if err != nil {
		log.Printf("Error updating thought: %v", err)
	} else {
		log.Printf("Updated thought: ID=%s, Relation=%s, ModuleRef=%s",
			updatedThought.ID, updatedThought.Relation, updatedThought.Module.Reference)
	}

	// Demonstrate creating multiple thoughts in a chain
	log.Printf("Creating additional thoughts for the chain...")

	thoughts := []perception.CreateThoughtRequest{
		{
			Thought: perception.ThoughtRequestData{
				Relation: "preprocessing",
				Module: perception.Module{
					Reference: "tama/agentic/preprocess",
					Parameters: map[string]any{
						"clean_text": true,
						"normalize":  true,
					},
				},
			},
		},
		{
			Thought: perception.ThoughtRequestData{
				Relation: "validation",
				Module: perception.Module{
					Reference: "tama/agentic/validate",
					Parameters: map[string]any{
						"strict_mode":    false,
						"schema_version": "v2",
					},
				},
			},
		},
	}

	for i, thoughtReq := range thoughts {
		newThought, createErr := client.Perception.CreateThought(exampleChainID, thoughtReq)
		if createErr != nil {
			log.Printf("Error creating thought %d: %v", i+1, createErr)
		} else {
			log.Printf("Created thought %d: ID=%s, Relation=%s, Index=%d",
				i+1, newThought.ID, newThought.Relation, newThought.Index)
		}
	}

	log.Printf("Thought operations completed successfully")
}
