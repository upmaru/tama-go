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

func TestPerceptionGetInitializer(t *testing.T) {
	expectedInitializer := perception.Initializer{
		ID: "initializer-123",
		Parameters: map[string]any{
			"max_tokens":  150,
			"temperature": 0.8,
		},
		Index:          intPtr(1),
		ProvisionState: "active",
		ThoughtID:      "thought-123",
		ClassID:        "class-123",
		Reference:      "tama/agentic/init",
	}

	expectedResponse := perception.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/initializers/initializer-123" {
			t.Errorf("Expected path /provision/perception/initializers/initializer-123, got %s", r.URL.Path)
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
	initializer, err := client.Perception.GetInitializer("initializer-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestPerceptionGetInitializerError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := perception.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"initializer": {"not found"},
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
	_, err := client.Perception.GetInitializer("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var perceptionErr *perception.Error
	if errors.As(err, &perceptionErr) {
		if perceptionErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", perceptionErr.StatusCode)
		}
		if perceptionErr.Errors == nil || len(perceptionErr.Errors["initializer"]) == 0 ||
			perceptionErr.Errors["initializer"][0] != "not found" {
			t.Errorf("Expected error 'initializer not found', got %v", perceptionErr.Errors)
		}
	} else {
		t.Errorf("Expected perception.Error, got %T", err)
	}
}

func TestPerceptionCreateInitializer(t *testing.T) {
	request := perception.CreateInitializerRequest{
		Initializer: perception.InitializerRequestData{
			Parameters: map[string]any{
				"temperature": 0.7,
				"max_tokens":  100,
			},
			ClassID:   "class-456",
			Reference: "tama/agentic/init",
		},
	}

	expectedInitializer := perception.Initializer{
		ID: "initializer-456",
		Parameters: map[string]any{
			"temperature": 0.7,
			"max_tokens":  100,
		},
		Index:          intPtr(2),
		ProvisionState: "pending",
		ThoughtID:      "thought-123",
		ClassID:        "class-456",
		Reference:      "tama/agentic/init",
	}

	expectedResponse := perception.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/initializers" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/initializers, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateInitializerRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Initializer.ClassID != request.Initializer.ClassID {
			t.Errorf(
				"Expected initializer class_id %s, got %s",
				request.Initializer.ClassID,
				receivedRequest.Initializer.ClassID,
			)
		}

		if receivedRequest.Initializer.Reference != request.Initializer.Reference {
			t.Errorf(
				"Expected initializer reference %s, got %s",
				request.Initializer.Reference,
				receivedRequest.Initializer.Reference,
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
	initializer, err := client.Perception.CreateInitializer("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestPerceptionCreateInitializerWithIndex(t *testing.T) {
	indexValue := 5
	request := perception.CreateInitializerRequest{
		Initializer: perception.InitializerRequestData{
			Parameters: map[string]any{
				"temperature": 0.6,
				"max_tokens":  120,
			},
			Index:     &indexValue,
			ClassID:   "class-789",
			Reference: "tama/agentic/init",
		},
	}

	expectedInitializer := perception.Initializer{
		ID: "initializer-789",
		Parameters: map[string]any{
			"temperature": 0.6,
			"max_tokens":  120,
		},
		Index:          intPtr(5),
		ProvisionState: "pending",
		ThoughtID:      "thought-123",
		ClassID:        "class-789",
		Reference:      "tama/agentic/init",
	}

	expectedResponse := perception.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/initializers" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/initializers, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateInitializerRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Validate that the index is properly included in the request body
		if receivedRequest.Initializer.Index == nil || *receivedRequest.Initializer.Index != *request.Initializer.Index {
			var receivedIndex int
			if receivedRequest.Initializer.Index != nil {
				receivedIndex = *receivedRequest.Initializer.Index
			}
			t.Errorf(
				"Expected initializer index %d, got %d",
				*request.Initializer.Index,
				receivedIndex,
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
	initializer, err := client.Perception.CreateInitializer("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestPerceptionCreateInitializerValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty thought ID
	_, err := client.Perception.CreateInitializer("", perception.CreateInitializerRequest{
		Initializer: perception.InitializerRequestData{
			ClassID:   "class-123",
			Reference: "tama/agentic/init",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought ID")
	}

	// Test empty class ID
	_, err = client.Perception.CreateInitializer("thought-123", perception.CreateInitializerRequest{
		Initializer: perception.InitializerRequestData{
			ClassID:   "",
			Reference: "tama/agentic/init",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty class ID")
	}

	// Test empty reference
	_, err = client.Perception.CreateInitializer("thought-123", perception.CreateInitializerRequest{
		Initializer: perception.InitializerRequestData{
			ClassID:   "class-123",
			Reference: "",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty reference")
	}
}

func TestPerceptionUpdateInitializer(t *testing.T) {
	indexValue := 3
	request := perception.UpdateInitializerRequest{
		Initializer: perception.UpdateInitializerData{
			Parameters: map[string]any{
				"depth": 3,
			},
			Index:     &indexValue,
			ClassID:   "class-updated",
			Reference: "tama/agentic/updated",
		},
	}

	expectedInitializer := perception.Initializer{
		ID: "initializer-111",
		Parameters: map[string]any{
			"depth": 3,
			"scope": "global",
		},
		Index:          intPtr(7),
		ProvisionState: "active",
		ThoughtID:      "thought-123",
		ClassID:        "class-updated",
		Reference:      "tama/agentic/updated",
	}

	expectedResponse := perception.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/initializers/initializer-123" {
			t.Errorf("Expected path /provision/perception/initializers/initializer-123, got %s", r.URL.Path)
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
	initializer, err := client.Perception.UpdateInitializer("initializer-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if initializer.ClassID != expectedInitializer.ClassID {
		t.Errorf("Expected initializer class_id %s, got %s", expectedInitializer.ClassID, initializer.ClassID)
	}

	if initializer.Reference != expectedInitializer.Reference {
		t.Errorf(
			"Expected initializer reference %s, got %s",
			expectedInitializer.Reference,
			initializer.Reference,
		)
	}
}

func TestPerceptionReplaceInitializer(t *testing.T) {
	request := perception.UpdateInitializerRequest{
		Initializer: perception.UpdateInitializerData{
			Parameters: map[string]any{
				"format": "json",
			},
			ClassID:   "class-replaced",
			Reference: "tama/agentic/replaced",
		},
	}

	expectedInitializer := perception.Initializer{
		ID: "initializer-222",
		Parameters: map[string]any{
			"format": "json",
		},
		Index:          intPtr(10),
		ProvisionState: "active",
		ThoughtID:      "thought-123",
		ClassID:        "class-replaced",
		Reference:      "tama/agentic/replaced",
	}

	expectedResponse := perception.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/initializers/initializer-123" {
			t.Errorf("Expected path /provision/perception/initializers/initializer-123, got %s", r.URL.Path)
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
	initializer, err := client.Perception.ReplaceInitializer("initializer-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestPerceptionDeleteInitializer(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/initializers/initializer-123" {
			t.Errorf("Expected path /provision/perception/initializers/initializer-123, got %s", r.URL.Path)
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
	err := client.Perception.DeleteInitializer("initializer-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.GetInitializer("")
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in GetInitializer")
	}
}

func TestPerceptionUpdateInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.UpdateInitializer("", perception.UpdateInitializerRequest{
		Initializer: perception.UpdateInitializerData{ClassID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in UpdateInitializer")
	}
}

func TestPerceptionReplaceInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.ReplaceInitializer("", perception.UpdateInitializerRequest{
		Initializer: perception.UpdateInitializerData{ClassID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in ReplaceInitializer")
	}
}

func TestPerceptionDeleteInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Perception.DeleteInitializer("")
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in DeleteInitializer")
	}
}

func TestPerceptionCreateInitializerWithZeroIndex(t *testing.T) {
	zeroIndex := 0
	request := perception.CreateInitializerRequest{
		Initializer: perception.InitializerRequestData{
			Parameters: map[string]any{
				"temperature": 0.8,
				"max_tokens":  150,
			},
			Index:     &zeroIndex, // Explicitly set zero index
			ClassID:   "class-456",
			Reference: "tama/agentic/init",
		},
	}

	expectedInitializer := perception.Initializer{
		ID: "initializer-456",
		Parameters: map[string]any{
			"temperature": 0.8,
			"max_tokens":  150,
		},
		Index:          intPtr(0),
		ProvisionState: "pending",
		ThoughtID:      "thought-123",
		ClassID:        "class-456",
		Reference:      "tama/agentic/init",
	}

	expectedResponse := perception.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/initializers" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/initializers, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateInitializerRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Validate that zero index is properly included in the request body
		if receivedRequest.Initializer.Index == nil || *receivedRequest.Initializer.Index != 0 {
			var receivedIndex int
			if receivedRequest.Initializer.Index != nil {
				receivedIndex = *receivedRequest.Initializer.Index
			}
			t.Errorf("Expected initializer index 0, got %d", receivedIndex)
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
	initializer, err := client.Perception.CreateInitializer("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if initializer.Index == nil || *initializer.Index != 0 {
		var actualIndex int
		if initializer.Index != nil {
			actualIndex = *initializer.Index
		}
		t.Errorf("Expected initializer index 0, got %d", actualIndex)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestPerceptionCreateInitializerWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/initializers" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/initializers, got %s",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		errorResponse := map[string]any{
			"errors": map[string][]string{
				"class_id":  {"can't be blank", "is not included in the list"},
				"reference": {"is invalid"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	request := perception.CreateInitializerRequest{
		Initializer: perception.InitializerRequestData{
			ClassID:   "invalid-class",   // Valid format to bypass client validation
			Reference: "invalid/reference", // Valid format to bypass client validation
		},
	}

	_, err := client.Perception.CreateInitializer("thought-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "class_id can't be blank") {
		t.Errorf("Expected error to contain 'class_id can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "class_id is not included in the list") {
		t.Errorf(
			"Expected error to contain 'class_id is not included in the list', got %s",
			errorMsg,
		)
	}
	if !strings.Contains(errorMsg, "reference is invalid") {
		t.Errorf("Expected error to contain 'reference is invalid', got %s", errorMsg)
	}
}

// Helper functions for initializer tests

// intPtr is a helper function to create a pointer to an int.
func intPtr(i int) *int {
	return &i
}
