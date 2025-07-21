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
	server := createMockServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"space": {"not found"},
			},
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
		if neuralErr.Errors == nil || len(neuralErr.Errors["space"]) == 0 ||
			neuralErr.Errors["space"][0] != "not found" {
			t.Errorf("Expected error 'space not found', got %v", neuralErr.Errors)
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

func TestNeuralGetProcessor(t *testing.T) {
	expectedProcessor := neural.Processor{
		ID:           "processor-123",
		SpaceID:      "space-123",
		ModelID:      "model-123",
		Type:         "completion",
		CurrentState: "active",
		Configuration: map[string]any{
			"batch_size":    32,
			"learning_rate": 0.001,
			"epochs":        100,
		},
	}

	expectedResponse := neural.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123/types/completion/processor" {
			t.Errorf("Expected path /provision/neural/spaces/space-123/types/completion/processor, got %s", r.URL.Path)
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
	processor, err := client.Neural.GetProcessor("space-123", "completion")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if processor.ID != expectedProcessor.ID {
		t.Errorf("Expected processor ID %s, got %s", expectedProcessor.ID, processor.ID)
	}

	if processor.Type != expectedProcessor.Type {
		t.Errorf("Expected processor type %s, got %s", expectedProcessor.Type, processor.Type)
	}

	if int(processor.Configuration["batch_size"].(float64)) != expectedProcessor.Configuration["batch_size"].(int) {
		t.Errorf(
			"Expected batch_size %v, got %v",
			expectedProcessor.Configuration["batch_size"],
			processor.Configuration["batch_size"],
		)
	}
}

func TestNeuralGetProcessorError(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"processor": {"not found"},
			},
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
	_, err := client.Neural.GetProcessor("space-123", "nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var neuralErr *neural.Error
	if errors.As(err, &neuralErr) {
		if neuralErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", neuralErr.StatusCode)
		}
		if neuralErr.Errors == nil || len(neuralErr.Errors["processor"]) == 0 ||
			neuralErr.Errors["processor"][0] != "not found" {
			t.Errorf("Expected error 'processor not found', got %v", neuralErr.Errors)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralCreateProcessor(t *testing.T) {
	expectedProcessor := neural.Processor{
		ID:           "processor-789",
		SpaceID:      "space-123",
		ModelID:      "model-123",
		Type:         "completion",
		CurrentState: "active",
		Configuration: map[string]any{
			"batch_size":    64,
			"learning_rate": 0.01,
		},
	}

	expectedResponse := neural.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123/types/completion/processor" {
			t.Errorf("Expected path /provision/neural/spaces/space-123/types/completion/processor, got %s", r.URL.Path)
		}

		var req neural.CreateProcessorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Processor.Type != "completion" {
			t.Errorf("Expected request type 'completion', got %s", req.Processor.Type)
		}

		if req.Processor.Configuration["temperature"] != 0.8 {
			t.Errorf("Expected temperature 0.8, got %v", req.Processor.Configuration["temperature"])
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

	createReq := neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			Type: "completion",
			Configuration: map[string]any{
				"temperature": 0.8,
				"tool_choice": "required",
			},
		},
	}

	processor, err := client.Neural.CreateProcessor("space-123", "completion", createReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if processor.ID != expectedProcessor.ID {
		t.Errorf("Expected processor ID %s, got %s", expectedProcessor.ID, processor.ID)
	}

	if processor.Type != expectedProcessor.Type {
		t.Errorf("Expected processor type %s, got %s", expectedProcessor.Type, processor.Type)
	}
}

func TestNeuralCreateProcessorValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	// Test empty space ID validation
	_, err := client.Neural.CreateProcessor("", "completion", neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			Type: "completion",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty processor type validation
	_, err = client.Neural.CreateProcessor("space-123", "", neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			Type: "completion",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty processor type")
	}

	// Test empty type validation
	_, err = client.Neural.CreateProcessor("space-123", "completion", neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			Configuration: map[string]any{"test": "value"},
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty type")
	}

	// Test invalid type validation
	_, err = client.Neural.CreateProcessor("space-123", "invalid-type", neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			Type:          "invalid-type",
			Configuration: map[string]any{"test": "value"},
		},
	})

	if err == nil {
		t.Error("Expected validation error for invalid type")
	}

	// Test valid types
	validTypes := []string{"completion", "embedding", "reranking"}
	for _, validType := range validTypes {
		_, err = client.Neural.CreateProcessor("space-123", validType, neural.CreateProcessorRequest{
			Processor: neural.ProcessorRequestData{
				Type:          validType,
				Configuration: map[string]any{"test": "value"},
			},
		})
		// We expect a network error since we're not mocking this, but not a validation error
		if err != nil && err.Error() == "processor type must be 'completion', 'embedding', or 'reranking'" {
			t.Errorf("Valid type %s should not cause validation error", validType)
		}
	}
}

func TestNeuralUpdateProcessor(t *testing.T) {
	expectedProcessor := neural.Processor{
		ID:           "processor-123",
		SpaceID:      "space-123",
		ModelID:      "model-123",
		Type:         "embedding",
		CurrentState: "active",
		Configuration: map[string]any{
			"max_tokens": 512,
		},
	}

	expectedResponse := neural.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123/types/embedding/processor" {
			t.Errorf("Expected path /provision/neural/spaces/space-123/types/embedding/processor, got %s", r.URL.Path)
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

	updateReq := neural.UpdateProcessorRequest{
		Processor: neural.UpdateProcessorData{
			Type: "embedding",
			Configuration: map[string]any{
				"max_tokens": 512,
			},
		},
	}

	processor, err := client.Neural.UpdateProcessor("space-123", "embedding", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if processor.Type != expectedProcessor.Type {
		t.Errorf("Expected processor type %s, got %s", expectedProcessor.Type, processor.Type)
	}
}

func TestNeuralReplaceProcessor(t *testing.T) {
	expectedProcessor := neural.Processor{
		ID:           "processor-123",
		SpaceID:      "space-123",
		ModelID:      "model-123",
		Type:         "reranking",
		CurrentState: "active",
		Configuration: map[string]any{
			"top_n": 3,
		},
	}

	expectedResponse := neural.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123/types/reranking/processor" {
			t.Errorf("Expected path /provision/neural/spaces/space-123/types/reranking/processor, got %s", r.URL.Path)
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

	replaceReq := neural.UpdateProcessorRequest{
		Processor: neural.UpdateProcessorData{
			Type: "reranking",
			Configuration: map[string]any{
				"top_n": 3,
			},
		},
	}

	processor, err := client.Neural.ReplaceProcessor("space-123", "reranking", replaceReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if processor.Type != expectedProcessor.Type {
		t.Errorf("Expected processor type %s, got %s", expectedProcessor.Type, processor.Type)
	}
}

func TestNeuralDeleteProcessor(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123/types/completion/processor" {
			t.Errorf("Expected path /provision/neural/spaces/space-123/types/completion/processor, got %s", r.URL.Path)
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

	err := client.Neural.DeleteProcessor("space-123", "completion")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
