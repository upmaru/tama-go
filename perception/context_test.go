package perception_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/perception"
)

// Context tests.
func TestPerceptionGetContext(t *testing.T) {
	expectedContext := perception.Context{
		ID:             "context-123",
		ThoughtID:      "thought-123",
		PromptID:       "prompt-123",
		Layer:          2,
		ProvisionState: "active",
	}

	expectedResponse := perception.ContextResponse{
		Data: expectedContext,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/contexts/context-123" {
			t.Errorf("Expected path /provision/perception/contexts/context-123, got %s", r.URL.Path)
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
	context, err := client.Perception.GetContext("context-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateContextResponse(t, *context, expectedContext)
}

func TestPerceptionCreateContext(t *testing.T) {
	request := perception.CreateContextRequest{
		Context: perception.ContextRequestData{
			PromptID: "prompt-456",
			Layer:    3,
		},
	}

	expectedContext := perception.Context{
		ID:             "context-456",
		ThoughtID:      "thought-123",
		PromptID:       "prompt-456",
		Layer:          3,
		ProvisionState: "pending",
	}

	expectedResponse := perception.ContextResponse{
		Data: expectedContext,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/contexts" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/contexts, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateContextRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Context.PromptID != request.Context.PromptID {
			t.Errorf(
				"Expected prompt_id %s, got %s",
				request.Context.PromptID,
				receivedRequest.Context.PromptID,
			)
		}

		if receivedRequest.Context.Layer != request.Context.Layer {
			t.Errorf(
				"Expected layer %d, got %d",
				request.Context.Layer,
				receivedRequest.Context.Layer,
			)
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
	context, err := client.Perception.CreateContext("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateContextResponse(t, *context, expectedContext)
}

func TestPerceptionCreateContextValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty thought ID
	_, err := client.Perception.CreateContext("", perception.CreateContextRequest{
		Context: perception.ContextRequestData{PromptID: "prompt-123", Layer: 1},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought ID")
	}

	// Test empty prompt ID
	_, err = client.Perception.CreateContext("thought-123", perception.CreateContextRequest{
		Context: perception.ContextRequestData{PromptID: "", Layer: 1},
	})
	if err == nil {
		t.Error("Expected validation error for empty prompt ID")
	}
}

func TestPerceptionUpdateContext(t *testing.T) {
	request := perception.UpdateContextRequest{
		Context: perception.UpdateContextData{
			PromptID: "prompt-789",
			Layer:    5,
		},
	}

	expectedContext := perception.Context{
		ID:             "context-123",
		ThoughtID:      "thought-123",
		PromptID:       "prompt-789",
		Layer:          5,
		ProvisionState: "active",
	}

	expectedResponse := perception.ContextResponse{
		Data: expectedContext,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/contexts/context-123" {
			t.Errorf("Expected path /provision/perception/contexts/context-123, got %s", r.URL.Path)
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
	context, err := client.Perception.UpdateContext("context-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if context.PromptID != expectedContext.PromptID {
		t.Errorf("Expected prompt_id %s, got %s", expectedContext.PromptID, context.PromptID)
	}

	if context.Layer != expectedContext.Layer {
		t.Errorf("Expected layer %d, got %d", expectedContext.Layer, context.Layer)
	}
}

func TestPerceptionReplaceContext(t *testing.T) {
	request := perception.UpdateContextRequest{
		Context: perception.UpdateContextData{
			PromptID: "prompt-replaced",
			Layer:    1,
		},
	}

	expectedContext := perception.Context{
		ID:             "context-123",
		ThoughtID:      "thought-123",
		PromptID:       "prompt-replaced",
		Layer:          1,
		ProvisionState: "active",
	}

	expectedResponse := perception.ContextResponse{
		Data: expectedContext,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/contexts/context-123" {
			t.Errorf("Expected path /provision/perception/contexts/context-123, got %s", r.URL.Path)
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
	context, err := client.Perception.ReplaceContext("context-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if context.PromptID != expectedContext.PromptID {
		t.Errorf("Expected prompt_id %s, got %s", expectedContext.PromptID, context.PromptID)
	}

	if context.Layer != expectedContext.Layer {
		t.Errorf("Expected layer %d, got %d", expectedContext.Layer, context.Layer)
	}
}

func TestPerceptionDeleteContext(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/contexts/context-123" {
			t.Errorf("Expected path /provision/perception/contexts/context-123, got %s", r.URL.Path)
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
	err := client.Perception.DeleteContext("context-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetContextEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.GetContext("")
	if err == nil {
		t.Error("Expected validation error for empty context ID in GetContext")
	}
}

func TestPerceptionUpdateContextEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.UpdateContext("", perception.UpdateContextRequest{
		Context: perception.UpdateContextData{PromptID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty context ID in UpdateContext")
	}
}

func TestPerceptionReplaceContextEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.ReplaceContext("", perception.UpdateContextRequest{
		Context: perception.UpdateContextData{PromptID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty context ID in ReplaceContext")
	}
}

func TestPerceptionDeleteContextEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Perception.DeleteContext("")
	if err == nil {
		t.Error("Expected validation error for empty context ID in DeleteContext")
	}
}
