package main

import (
	"errors"
	"log"
	"strings"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/memory"
	"github.com/upmaru/tama-go/neural"
	"github.com/upmaru/tama-go/sensory"
)

const (
	exampleSpaceID    = "space-123"
	exampleSourceID   = "source-123"
	exampleModelID    = "model-123"
	exampleLimitID    = "limit-123"
	examplePromptID   = "prompt-123"
	defaultTimeout    = 30
	defaultLimitCount = 32
	scaleCountValue   = 5
	limitCountValue   = 100
)

func main() {
	client := initializeClient()

	// Run examples in separate functions to reduce complexity
	// Run examples
	runNeuralSpaceOperations(client)
	runMemoryPromptOperations(client)
	runSensorySourceOperations(client)
	runSensoryModelOperations(client)
	runSensoryLimitOperations(client)

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

		err = client.Neural.DeleteSpace(spaceID)
		if err != nil {
			log.Printf("Error deleting space: %v", err)
		} else {
			log.Printf("Deleted space successfully")
		}
	*/
}
