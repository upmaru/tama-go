package neural_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
)

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

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
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
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Neural.GetProcessor("space-123", "nonexistent")

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

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

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
		BaseURL:        "http://localhost:8080",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty space ID validation
	_, err = client.Neural.CreateProcessor("", "completion", neural.CreateProcessorRequest{
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

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

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

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

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
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Neural.DeleteProcessor("space-123", "completion")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
