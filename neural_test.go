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
		ID:             "space-123",
		Name:           "test-space",
		Slug:           "test-space-slug",
		Type:           "root",
		ProvisionState: "active",
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
		ID:             "space-789",
		Name:           "new-space",
		Slug:           "new-space-slug",
		Type:           "root",
		ProvisionState: "active",
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

func TestNeuralCreateSpaceTypeValidationDelegated(t *testing.T) {
	// This test verifies that type validation is now delegated to the API
	// Previously invalid types should now pass client validation
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	// Test that previously invalid type values now pass client validation
	// The API will handle the actual validation
	_, err := client.Neural.CreateSpace(neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "test-space",
			Type: "invalid-type", // This should not trigger client validation error
		},
	})

	// Should not get a client validation error for type
	// Any error should come from the actual API call (network/API error)
	if err != nil && err.Error() == "space type must be 'root' or 'component'" {
		t.Error("Client should not validate specific type values - validation should be delegated to API")
	}

	// Test another previously invalid type
	_, err = client.Neural.CreateSpace(neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "test-space",
			Type: "custom-type",
		},
	})

	// Should not get the old client validation error
	if err != nil && err.Error() == "space type must be 'root' or 'component'" {
		t.Error("Client should not validate specific type values - validation should be delegated to API")
	}
}

func TestNeuralUpdateSpace(t *testing.T) {
	expectedSpace := neural.Space{
		ID:             "space-123",
		Name:           "updated-space",
		Slug:           "updated-slug",
		Type:           "component",
		ProvisionState: "active",
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
		ID:             "processor-123",
		SpaceID:        "space-123",
		ModelID:        "model-123",
		Type:           "completion",
		ProvisionState: "active",
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
		ID:             "processor-789",
		SpaceID:        "space-123",
		ModelID:        "model-123",
		Type:           "completion",
		ProvisionState: "active",
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

		if req.Processor.ModelID != "model-123" {
			t.Errorf("Expected request model ID 'model-123', got %s", req.Processor.ModelID)
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
			ModelID: "model-123",
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
			ModelID: "model-123",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty processor type validation
	_, err = client.Neural.CreateProcessor("space-123", "", neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			ModelID: "model-123",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty processor type")
	}

	// Test empty model ID validation
	_, err = client.Neural.CreateProcessor("space-123", "completion", neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			Configuration: map[string]any{"test": "value"},
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty model ID")
	}

	// Test valid model ID
	_, err = client.Neural.CreateProcessor("space-123", "completion", neural.CreateProcessorRequest{
		Processor: neural.ProcessorRequestData{
			ModelID:       "valid-model-123",
			Configuration: map[string]any{"test": "value"},
		},
	})
	// We expect a network error since we're not mocking this, but not a validation error
	if err != nil && err.Error() == "model ID is required" {
		t.Error("Valid model ID should not cause validation error")
	}
}

func TestNeuralUpdateProcessor(t *testing.T) {
	expectedProcessor := neural.Processor{
		ID:             "processor-123",
		SpaceID:        "space-123",
		ModelID:        "model-123",
		Type:           "embedding",
		ProvisionState: "active",
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
			ModelID: "model-123",
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
		ID:             "processor-123",
		SpaceID:        "space-123",
		ModelID:        "model-123",
		Type:           "reranking",
		ProvisionState: "active",
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
			ModelID: "model-123",
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

func TestNeuralGetClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-123",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
		},
		Name:        "TestClass",
		Description: "Test class for schema validation",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
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

	class, err := client.Neural.GetClass("class-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.SpaceID != expectedClass.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expectedClass.SpaceID, class.SpaceID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralGetClassError(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string][]string{
				"class": {"not found"},
			},
		})
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	_, err := client.Neural.GetClass("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiError *neural.Error
	if !errors.As(err, &apiError) {
		t.Fatalf("Expected neural.Error, got %T", err)
	}

	if apiError.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", apiError.StatusCode)
	}
}

func TestNeuralCreateClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-789",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
		},
		Name:        "NewClass",
		Description: "A new class for testing",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-456/classes" {
			t.Errorf("Expected path /provision/neural/spaces/space-456/classes, got %s", r.URL.Path)
		}

		var req neural.CreateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Class.Schema == nil {
			t.Error("Expected schema in request, got nil")
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

	createReq := neural.CreateClassRequest{
		Class: neural.ClassRequestData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{"type": "string"},
				},
			},
		},
	}

	class, err := client.Neural.CreateClass("space-456", createReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.SpaceID != expectedClass.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expectedClass.SpaceID, class.SpaceID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralCreateClassValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "http://localhost:8080",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	// Test missing space ID
	createReq := neural.CreateClassRequest{
		Class: neural.ClassRequestData{
			Schema: map[string]any{"type": "object"},
		},
	}

	_, err := client.Neural.CreateClass("", createReq)
	if err == nil {
		t.Error("Expected error for missing space ID, got nil")
	}

	// Test missing schema
	createReqNoSchema := neural.CreateClassRequest{
		Class: neural.ClassRequestData{},
	}

	_, err = client.Neural.CreateClass("space-123", createReqNoSchema)
	if err == nil {
		t.Error("Expected error for missing schema, got nil")
	}
}

func TestNeuralUpdateClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-123",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
				"age":  map[string]any{"type": "number"},
			},
		},
		Name:        "TestClass",
		Description: "Updated test class",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
		}

		var req neural.UpdateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
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

	updateReq := neural.UpdateClassRequest{
		Class: neural.UpdateClassData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{"type": "string"},
					"age":  map[string]any{"type": "number"},
				},
			},
		},
	}

	class, err := client.Neural.UpdateClass("class-123", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralReplaceClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-123",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title": map[string]any{"type": "string"},
			},
		},
		Name:        "ReplacedClass",
		Description: "A completely replaced class",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
		}

		var req neural.UpdateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
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

	replaceReq := neural.UpdateClassRequest{
		Class: neural.UpdateClassData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{"type": "string"},
				},
			},
		},
	}

	class, err := client.Neural.ReplaceClass("class-123", replaceReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralDeleteClass(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
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

	err := client.Neural.DeleteClass("class-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNeuralCreateClassWithRealWorldSchema(t *testing.T) {
	expectedClass := createActionCallClass()
	expectedResponse := neural.ClassResponse{Data: expectedClass}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		validateCreateClassRequest(t, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := createTestClient(server.URL)
	createReq := neural.CreateClassRequest{
		Class: neural.ClassRequestData{Schema: createActionCallSchema()},
	}

	class, err := client.Neural.CreateClass("space-123", createReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateClassResponse(t, *class, expectedClass)
}

func createActionCallClass() neural.Class {
	return neural.Class{
		ID:             "class-action-call",
		SpaceID:        "space-123",
		ProvisionState: "active",
		Schema:         createActionCallSchema(),
		Name:           "ActionCall",
		Description:    "Schema for action call requests",
	}
}

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

func validateCreateClassRequest(t *testing.T, r *http.Request) {
	if r.Method != http.MethodPost {
		t.Errorf("Expected POST request, got %s", r.Method)
	}
	if r.URL.Path != "/provision/neural/spaces/space-123/classes" {
		t.Errorf("Expected path /provision/neural/spaces/space-123/classes, got %s", r.URL.Path)
	}

	var req neural.CreateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatalf("Failed to decode request body: %v", err)
	}

	validateSchemaStructure(t, req.Class.Schema)
}

func validateSchemaStructure(t *testing.T, schema map[string]any) {
	if schema == nil {
		t.Error("Expected schema in request, got nil")
		return
	}

	title, ok := schema["title"].(string)
	if !ok || title != "action-call" {
		t.Errorf("Expected schema title 'action-call', got %v", schema["title"])
	}

	expectedDesc := "An action call is a request to execute an action."
	desc, ok := schema["description"].(string)
	if !ok || desc != expectedDesc {
		t.Errorf("Expected specific description, got %v", schema["description"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Error("Expected properties to be a map")
		return
	}

	requiredProps := []string{"code", "tool_id", "parameters", "content_type", "content"}
	for _, prop := range requiredProps {
		if _, exists := props[prop]; !exists {
			t.Errorf("Expected property %s to exist in schema", prop)
		}
	}
}

func createTestClient(baseURL string) *tama.Client {
	config := tama.Config{
		BaseURL: baseURL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}
	return tama.NewClient(config)
}

func validateClassResponse(t *testing.T, actual, expected neural.Class) {
	if actual.ID != expected.ID {
		t.Errorf("Expected class ID %s, got %s", expected.ID, actual.ID)
	}
	if actual.SpaceID != expected.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expected.SpaceID, actual.SpaceID)
	}
	if actual.Name != expected.Name {
		t.Errorf("Expected class name %s, got %s", expected.Name, actual.Name)
	}

	title, ok := actual.Schema["title"].(string)
	if !ok || title != "action-call" {
		t.Errorf("Expected schema title 'action-call', got %v", actual.Schema["title"])
	}

	expectedDesc := "An action call is a request to execute an action."
	desc, ok := actual.Schema["description"].(string)
	if !ok || desc != expectedDesc {
		t.Errorf("Expected specific description, got %v", actual.Schema["description"])
	}
}
