package main

import (
	"fmt"
	"log"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
	"github.com/upmaru/tama-go/sensory"
)

func main() {
	// Initialize the client
	config := tama.Config{
		BaseURL: "https://api.tama.io", // Replace with your actual API base URL
		APIKey:  "your-api-key",        // Replace with your actual API key
		Timeout: 30 * time.Second,
	}

	client := tama.NewClient(config)

	// Enable debug mode to see HTTP requests/responses (optional)
	client.SetDebug(true)

	// Example: Neural Space operations
	fmt.Println("=== Neural Space Operations ===")

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
		fmt.Printf("Created space: %+v\n", space)
	}

	// Get a space by ID (replace with actual ID)
	spaceID := "space-123"
	space, err = client.Neural.GetSpace(spaceID)
	if err != nil {
		log.Printf("Error getting space: %v", err)
	} else {
		fmt.Printf("Retrieved space: %+v\n", space)
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
		fmt.Printf("Updated space: %+v\n", space)
	}

	// Example: Sensory Source operations
	fmt.Println("\n=== Sensory Source Operations ===")

	// Create a new source
	newSource := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "My Data Source",
			Type:     "model",
			Endpoint: "https://api.mistral.ai/v1",
			Credential: sensory.SourceCredential{
				ApiKey: "your-api-key-here",
			},
		},
	}

	source, err := client.Sensory.CreateSource(spaceID, newSource)
	if err != nil {
		log.Printf("Error creating source: %v", err)
	} else {
		fmt.Printf("Created source: %+v\n", source)
	}

	// Get a source by ID (replace with actual ID)
	sourceID := "source-123"
	source, err = client.Sensory.GetSource(sourceID)
	if err != nil {
		log.Printf("Error getting source: %v", err)
	} else {
		fmt.Printf("Retrieved source: %+v\n", source)
	}

	// Update a source
	updateSource := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name:     "Updated Data Source",
			Type:     "model",
			Endpoint: "https://api.openai.com/v1",
			Credential: &sensory.SourceCredential{
				ApiKey: "your-updated-api-key",
			},
		},
	}

	source, err = client.Sensory.UpdateSource(sourceID, updateSource)
	if err != nil {
		log.Printf("Error updating source: %v", err)
	} else {
		fmt.Printf("Updated source: %+v\n", source)
	}

	// Example: Sensory Model operations
	fmt.Println("\n=== Sensory Model Operations ===")

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
		fmt.Printf("Created model: %+v\n", model)
	}

	// Get a model by ID (replace with actual ID)
	modelID := "model-123"
	model, err = client.Sensory.GetModel(modelID)
	if err != nil {
		log.Printf("Error getting model: %v", err)
	} else {
		fmt.Printf("Retrieved model: %+v\n", model)
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
		fmt.Printf("Updated model: %+v\n", model)
	}

	// Example: Sensory Limit operations
	fmt.Println("\n=== Sensory Limit Operations ===")

	// Create a new limit
	newLimit := sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			ScaleUnit:  "seconds",
			ScaleCount: 1,
			Count:      32,
		},
	}

	limit, err := client.Sensory.CreateLimit(sourceID, newLimit)
	if err != nil {
		log.Printf("Error creating limit: %v", err)
	} else {
		fmt.Printf("Created limit: %+v\n", limit)
	}

	// Get a limit by ID (replace with actual ID)
	limitID := "limit-123"
	limit, err = client.Sensory.GetLimit(limitID)
	if err != nil {
		log.Printf("Error getting limit: %v", err)
	} else {
		fmt.Printf("Retrieved limit: %+v\n", limit)
	}

	// Update a limit
	updateLimit := sensory.UpdateLimitRequest{
		Limit: sensory.UpdateLimitData{
			ScaleUnit:    "minutes",
			ScaleCount:   5,
			Count:        100,
			CurrentState: "active",
		},
	}

	limit, err = client.Sensory.UpdateLimit(limitID, updateLimit)
	if err != nil {
		log.Printf("Error updating limit: %v", err)
	} else {
		fmt.Printf("Updated limit: %+v\n", limit)
	}

	// Example: Delete operations
	fmt.Println("\n=== Delete Operations ===")

	// Delete resources (uncomment to test)
	/*
		err = client.Sensory.DeleteLimit(limitID)
		if err != nil {
			log.Printf("Error deleting limit: %v", err)
		} else {
			fmt.Println("Deleted limit successfully")
		}

		err = client.Sensory.DeleteModel(modelID)
		if err != nil {
			log.Printf("Error deleting model: %v", err)
		} else {
			fmt.Println("Deleted model successfully")
		}

		err = client.Sensory.DeleteSource(sourceID)
		if err != nil {
			log.Printf("Error deleting source: %v", err)
		} else {
			fmt.Println("Deleted source successfully")
		}

		err = client.Neural.DeleteSpace(spaceID)
		if err != nil {
			log.Printf("Error deleting space: %v", err)
		} else {
			fmt.Println("Deleted space successfully")
		}
	*/

	fmt.Println("\nExample completed!")
}
