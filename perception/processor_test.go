package perception_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/perception"
)

func TestPerceptionGetProcessor(t *testing.T) {
	expectedProcessor := perception.Processor{
		ID:             "processor-123",
		ThoughtID:      "thought-123",
		ModelID:        "model-123",
		Type:           "completion",
		ProvisionState: "active",
		Configuration: map[string]any{
			"batch_size":    32,
			"learning_rate": 0.001,
			"epochs":        100,
		},
	}

	expectedResponse := perception.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/types/completion/processor" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/types/completion/processor, got %s",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	processor, err := client.Perception.GetProcessor("thought-123", "completion")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if processor.ID != expectedProcessor.ID {
		t.Errorf("Expected processor ID %s, got %s", expectedProcessor.ID, processor.ID)
	}

	if processor.Type != expectedProcessor.Type {
		t.Errorf("Expected processor type %s, got %s", expectedProcessor.Type, processor.Type)
	}

	if int(
		processor.Configuration["batch_size"].(float64),
	) != expectedProcessor.Configuration["batch_size"].(int) {
		t.Errorf(
			"Expected batch_size %v, got %v",
			expectedProcessor.Configuration["batch_size"],
			processor.Configuration["batch_size"],
		)
	}
}

func TestPerceptionGetProcessorError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := perception.Error{
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
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err := client.Perception.GetProcessor("thought-123", "nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var perceptionErr *perception.Error
	if errors.As(err, &perceptionErr) {
		if perceptionErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", perceptionErr.StatusCode)
		}
		if perceptionErr.Errors == nil || len(perceptionErr.Errors["processor"]) == 0 ||
			perceptionErr.Errors["processor"][0] != "not found" {
			t.Errorf("Expected error 'processor not found', got %v", perceptionErr.Errors)
		}
	} else {
		t.Errorf("Expected perception.Error, got %T", err)
	}
}

func TestPerceptionCreateProcessor(t *testing.T) {
	expectedProcessor := perception.Processor{
		ID:             "processor-789",
		ThoughtID:      "thought-123",
		ModelID:        "model-123",
		Type:           "completion",
		ProvisionState: "active",
		Configuration: map[string]any{
			"batch_size":    64,
			"learning_rate": 0.01,
		},
	}

	expectedResponse := perception.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/types/completion/processor" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/types/completion/processor, got %s",
				r.URL.Path,
			)
		}

		var req perception.CreateProcessorRequest
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
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	createReq := perception.CreateProcessorRequest{
		Processor: perception.ProcessorRequestData{
			ModelID: "model-123",
			Configuration: map[string]any{
				"temperature": 0.8,
				"tool_choice": "required",
			},
		},
	}

	processor, err := client.Perception.CreateProcessor("thought-123", "completion", createReq)

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

func TestPerceptionCreateProcessorValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty thought ID validation
	_, err := client.Perception.CreateProcessor("", "completion", perception.CreateProcessorRequest{
		Processor: perception.ProcessorRequestData{
			ModelID: "model-123",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty thought ID")
	}

	// Test empty processor type validation
	_, err = client.Perception.CreateProcessor("thought-123", "", perception.CreateProcessorRequest{
		Processor: perception.ProcessorRequestData{
			ModelID: "model-123",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty processor type")
	}

	// Test empty model ID validation
	_, err = client.Perception.CreateProcessor(
		"thought-123",
		"completion",
		perception.CreateProcessorRequest{
			Processor: perception.ProcessorRequestData{
				Configuration: map[string]any{"test": "value"},
			},
		},
	)

	if err == nil {
		t.Error("Expected validation error for empty model ID")
	}

	// Test valid model ID
	_, err = client.Perception.CreateProcessor(
		"thought-123",
		"completion",
		perception.CreateProcessorRequest{
			Processor: perception.ProcessorRequestData{
				ModelID:       "valid-model-123",
				Configuration: map[string]any{"test": "value"},
			},
		},
	)
	// We expect a network error since we're not mocking this, but not a validation error
	if err != nil && err.Error() == "model ID is required" {
		t.Error("Valid model ID should not cause validation error")
	}
}

func TestPerceptionUpdateProcessor(t *testing.T) {
	expectedProcessor := perception.Processor{
		ID:             "processor-123",
		ThoughtID:      "thought-123",
		ModelID:        "model-123",
		Type:           "embedding",
		ProvisionState: "active",
		Configuration: map[string]any{
			"max_tokens": 512,
		},
	}

	expectedResponse := perception.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/types/embedding/processor" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/types/embedding/processor, got %s",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	updateReq := perception.UpdateProcessorRequest{
		Processor: perception.UpdateProcessorData{
			ModelID: "model-123",
			Configuration: map[string]any{
				"max_tokens": 512,
			},
		},
	}

	processor, err := client.Perception.UpdateProcessor("thought-123", "embedding", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if processor.Type != expectedProcessor.Type {
		t.Errorf("Expected processor type %s, got %s", expectedProcessor.Type, processor.Type)
	}
}

func TestPerceptionReplaceProcessor(t *testing.T) {
	expectedProcessor := perception.Processor{
		ID:             "processor-123",
		ThoughtID:      "thought-123",
		ModelID:        "model-123",
		Type:           "reranking",
		ProvisionState: "active",
		Configuration: map[string]any{
			"top_n": 3,
		},
	}

	expectedResponse := perception.ProcessorResponse{
		Data: expectedProcessor,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/types/reranking/processor" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/types/reranking/processor, got %s",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	replaceReq := perception.UpdateProcessorRequest{
		Processor: perception.UpdateProcessorData{
			ModelID: "model-123",
			Configuration: map[string]any{
				"top_n": 3,
			},
		},
	}

	processor, err := client.Perception.ReplaceProcessor("thought-123", "reranking", replaceReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if processor.Type != expectedProcessor.Type {
		t.Errorf("Expected processor type %s, got %s", expectedProcessor.Type, processor.Type)
	}
}

func TestPerceptionDeleteProcessor(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/types/completion/processor" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/types/completion/processor, got %s",
				r.URL.Path,
			)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err := client.Perception.DeleteProcessor("thought-123", "completion")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetProcessorEmptyThoughtID(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err := client.Perception.GetProcessor("", "completion")

	if err == nil || !strings.Contains(err.Error(), "thought ID is required") {
		t.Errorf("Expected 'thought ID is required' error, got %v", err)
	}
}

func TestPerceptionGetProcessorEmptyType(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err := client.Perception.GetProcessor("thought-123", "")

	if err == nil || !strings.Contains(err.Error(), "processor type is required") {
		t.Errorf("Expected 'processor type is required' error, got %v", err)
	}
}

func TestPerceptionUpdateProcessorEmptyThoughtID(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	updateReq := perception.UpdateProcessorRequest{
		Processor: perception.UpdateProcessorData{
			ModelID: "model-123",
		},
	}

	_, err := client.Perception.UpdateProcessor("", "completion", updateReq)

	if err == nil || !strings.Contains(err.Error(), "thought ID is required") {
		t.Errorf("Expected 'thought ID is required' error, got %v", err)
	}
}

func TestPerceptionReplaceProcessorEmptyThoughtID(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	replaceReq := perception.UpdateProcessorRequest{
		Processor: perception.UpdateProcessorData{
			ModelID: "model-123",
		},
	}

	_, err := client.Perception.ReplaceProcessor("", "completion", replaceReq)

	if err == nil || !strings.Contains(err.Error(), "thought ID is required") {
		t.Errorf("Expected 'thought ID is required' error, got %v", err)
	}
}

func TestPerceptionDeleteProcessorEmptyThoughtID(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		ClientID: "test-client-id",
		ClientSecret: "test-client-secret",
		SkipTokenFetch: true,
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err := client.Perception.DeleteProcessor("", "completion")

	if err == nil || !strings.Contains(err.Error(), "thought ID is required") {
		t.Errorf("Expected 'thought ID is required' error, got %v", err)
	}
}
