//go:build integration
// +build integration

package tama

import (
	"os"
	"testing"
	"time"

	"github.com/upmaru/tama-go/neural"
	"github.com/upmaru/tama-go/sensory"
)

// Integration tests require actual API credentials and endpoint
// Run with: go test -tags=integration -v

func TestIntegrationNeuralSpaceLifecycle(t *testing.T) {
	baseURL := os.Getenv("TAMA_BASE_URL")
	apiKey := os.Getenv("TAMA_API_KEY")

	if baseURL == "" || apiKey == "" {
		t.Skip("Skipping integration test: TAMA_BASE_URL and TAMA_API_KEY environment variables must be set")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second,
	})

	// Test space creation
	createReq := neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "Integration Test Space",
			Type: "root",
		},
	}

	space, err := client.Neural.CreateSpace(createReq)
	if err != nil {
		t.Fatalf("Failed to create space: %v", err)
	}

	if space.ID == "" {
		t.Fatal("Created space should have an ID")
	}

	t.Logf("Created space with ID: %s", space.ID)

	// Test space retrieval
	retrievedSpace, err := client.Neural.GetSpace(space.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve space: %v", err)
	}

	if retrievedSpace.ID != space.ID {
		t.Errorf("Retrieved space ID %s doesn't match created space ID %s", retrievedSpace.ID, space.ID)
	}

	// Test space update
	updateReq := neural.UpdateSpaceRequest{
		Space: neural.UpdateSpaceData{
			Name: "Updated Integration Test Space",
			Type: "component",
		},
	}

	updatedSpace, err := client.Neural.UpdateSpace(space.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update space: %v", err)
	}

	if updatedSpace.Name != updateReq.Space.Name {
		t.Errorf("Expected name '%s', got '%s'", updateReq.Space.Name, updatedSpace.Name)
	}

	// Test space deletion
	err = client.Neural.DeleteSpace(space.ID)
	if err != nil {
		t.Fatalf("Failed to delete space: %v", err)
	}

	// Verify space is deleted
	_, err = client.Neural.GetSpace(space.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted space")
	}

	t.Log("Space lifecycle test completed successfully")
}

func TestIntegrationSensorySourceLifecycle(t *testing.T) {
	baseURL := os.Getenv("TAMA_BASE_URL")
	apiKey := os.Getenv("TAMA_API_KEY")
	spaceID := os.Getenv("TAMA_TEST_SPACE_ID")

	if baseURL == "" || apiKey == "" || spaceID == "" {
		t.Skip("Skipping integration test: TAMA_BASE_URL, TAMA_API_KEY, and TAMA_TEST_SPACE_ID environment variables must be set")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second,
	})

	// Test source creation
	createReq := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "Integration Test Source",
			Type:     "model",
			Endpoint: "https://api.test.com/v1",
			Credential: sensory.SourceCredential{
				ApiKey: "test-key-12345",
			},
		},
	}

	source, err := client.Sensory.CreateSource(spaceID, createReq)
	if err != nil {
		t.Fatalf("Failed to create source: %v", err)
	}

	if source.ID == "" {
		t.Fatal("Created source should have an ID")
	}

	t.Logf("Created source with ID: %s", source.ID)

	// Test source retrieval
	retrievedSource, err := client.Sensory.GetSource(source.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve source: %v", err)
	}

	if retrievedSource.ID != source.ID {
		t.Errorf("Retrieved source ID %s doesn't match created source ID %s", retrievedSource.ID, source.ID)
	}

	// Test source update
	updateReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name:     "Updated Integration Test Source",
			Type:     "model",
			Endpoint: "https://api.updated-test.com/v1",
			Credential: &sensory.SourceCredential{
				ApiKey: "updated-test-key-67890",
			},
		},
	}

	updatedSource, err := client.Sensory.UpdateSource(source.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update source: %v", err)
	}

	if updatedSource.Name != updateReq.Source.Name {
		t.Errorf("Expected name '%s', got '%s'", updateReq.Source.Name, updatedSource.Name)
	}

	// Test source deletion
	err = client.Sensory.DeleteSource(source.ID)
	if err != nil {
		t.Fatalf("Failed to delete source: %v", err)
	}

	// Verify source is deleted
	_, err = client.Sensory.GetSource(source.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted source")
	}

	t.Log("Source lifecycle test completed successfully")
}

func TestIntegrationErrorHandling(t *testing.T) {
	baseURL := os.Getenv("TAMA_BASE_URL")
	apiKey := os.Getenv("TAMA_API_KEY")

	if baseURL == "" || apiKey == "" {
		t.Skip("Skipping integration test: TAMA_BASE_URL and TAMA_API_KEY environment variables must be set")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Timeout: 30 * time.Second,
	})

	// Test 404 error
	_, err := client.Neural.GetSpace("nonexistent-space-id")
	if err == nil {
		t.Error("Expected error for nonexistent space")
	}

	if apiErr, ok := err.(*Error); ok {
		if apiErr.StatusCode != 404 {
			t.Errorf("Expected 404 status code, got %d", apiErr.StatusCode)
		}
	} else {
		t.Errorf("Expected *Error type, got %T", err)
	}

	t.Log("Error handling test completed successfully")
}

func TestIntegrationAuthentication(t *testing.T) {
	baseURL := os.Getenv("TAMA_BASE_URL")

	if baseURL == "" {
		t.Skip("Skipping integration test: TAMA_BASE_URL environment variable must be set")
	}

	// Test with invalid API key
	client := NewClient(Config{
		BaseURL: baseURL,
		APIKey:  "invalid-api-key",
		Timeout: 30 * time.Second,
	})

	_, err := client.Neural.GetSpace("any-space-id")
	if err == nil {
		t.Error("Expected authentication error with invalid API key")
	}

	if apiErr, ok := err.(*Error); ok {
		if apiErr.StatusCode != 401 && apiErr.StatusCode != 403 {
			t.Errorf("Expected 401 or 403 status code for auth error, got %d", apiErr.StatusCode)
		}
	}

	t.Log("Authentication test completed successfully")
}
