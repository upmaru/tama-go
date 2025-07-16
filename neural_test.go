package tama_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
)

func TestNeuralGetSpace(t *testing.T) {
	expectedSpace := neural.Space{
		ID:           "space-123",
		Name:         "test-space",
		Slug:         "test-space-slug",
		Type:         "root",
		CurrentState: "active",
	}

	expectedResponse := neural.SpaceResponse{
		Data: expectedSpace,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123" {
			t.Errorf("Expected path /provision/neural/spaces/space-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	space, err := client.Neural.GetSpace("space-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if space.ID != expectedSpace.ID {
		t.Errorf("Expected space ID %s, got %s", expectedSpace.ID, space.ID)
	}

	if space.Name != expectedSpace.Name {
		t.Errorf("Expected space name %s, got %s", expectedSpace.Name, space.Name)
	}
}

func TestNeuralGetSpaceError(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Message:    "Space not found",
			Details:    "The requested space does not exist",
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	_, err := client.Neural.GetSpace("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var neuralErr *neural.Error
	if errors.As(err, &neuralErr) {
		if neuralErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", neuralErr.StatusCode)
		}
		if neuralErr.Message != "Space not found" {
			t.Errorf("Expected message 'Space not found', got %s", neuralErr.Message)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralCreateSpace(t *testing.T) {
	expectedSpace := neural.Space{
		ID:           "space-789",
		Name:         "new-space",
		Slug:         "new-space-slug",
		Type:         "root",
		CurrentState: "active",
	}

	expectedResponse := neural.SpaceResponse{
		Data: expectedSpace,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces" {
			t.Errorf("Expected path /provision/neural/spaces, got %s", r.URL.Path)
		}

		var req neural.CreateSpaceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Space.Name != "new-space" {
			t.Errorf("Expected request name 'new-space', got %s", req.Space.Name)
		}

		if req.Space.Type != "root" {
			t.Errorf("Expected request type 'root', got %s", req.Space.Type)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	createReq := neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "new-space",
			Type: "root",
		},
	}

	space, err := client.Neural.CreateSpace(createReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if space.ID != expectedSpace.ID {
		t.Errorf("Expected space ID %s, got %s", expectedSpace.ID, space.ID)
	}

	if space.Name != expectedSpace.Name {
		t.Errorf("Expected space name %s, got %s", expectedSpace.Name, space.Name)
	}
}

func TestNeuralCreateSpaceValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	// Test empty name validation
	_, err := client.Neural.CreateSpace(neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Type: "root",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	// Test empty type validation
	_, err = client.Neural.CreateSpace(neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "test-name",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty type")
	}
}

func TestNeuralUpdateSpace(t *testing.T) {
	expectedSpace := neural.Space{
		ID:           "space-123",
		Name:         "updated-space",
		Slug:         "updated-slug",
		Type:         "component",
		CurrentState: "active",
	}

	expectedResponse := neural.SpaceResponse{
		Data: expectedSpace,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123" {
			t.Errorf("Expected path /provision/neural/spaces/space-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	updateReq := neural.UpdateSpaceRequest{
		Space: neural.UpdateSpaceData{
			Name: "updated-space",
			Type: "component",
		},
	}

	space, err := client.Neural.UpdateSpace("space-123", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if space.Name != expectedSpace.Name {
		t.Errorf("Expected space name %s, got %s", expectedSpace.Name, space.Name)
	}
}

func TestNeuralDeleteSpace(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123" {
			t.Errorf("Expected path /provision/neural/spaces/space-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	err := client.Neural.DeleteSpace("space-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
