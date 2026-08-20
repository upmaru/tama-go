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

func TestPerceptionGetThought(t *testing.T) {
	expectedThought := perception.Thought{
		ID:            "thought-123",
		ChainID:       "chain-123",
		OutputClassID: "class-123",
		Module: &perception.Module{
			ID:        "module-123",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.7,
				"max_tokens":  100,
			},
		},
		Faculty: &perception.Faculty{
			QueueID:  "queue-abc",
			Priority: 2,
		},
		ProvisionState: "active",
		Relation:       "description",
		Index:          1,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123, got %s", r.URL.Path)
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
	thought, err := client.Perception.GetThought("thought-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionGetThoughtError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := perception.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"thought": {"not found"},
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
	_, err = client.Perception.GetThought("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var perceptionErr *perception.Error
	if errors.As(err, &perceptionErr) {
		if perceptionErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", perceptionErr.StatusCode)
		}
		if perceptionErr.Errors == nil || len(perceptionErr.Errors["thought"]) == 0 ||
			perceptionErr.Errors["thought"][0] != "not found" {
			t.Errorf("Expected error 'thought not found', got %v", perceptionErr.Errors)
		}
	} else {
		t.Errorf("Expected perception.Error, got %T", err)
	}
}

func TestPerceptionCreateThought(t *testing.T) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
				Parameters: map[string]any{
					"temperature": 0.8,
					"max_tokens":  150,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-456",
		ChainID:       "chain-123",
		OutputClassID: "class-123",
		Module: &perception.Module{
			ID:        "module-456",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.8,
				"max_tokens":  150,
			},
		},
		ProvisionState: "pending",
		Relation:       "description",
		Index:          2,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf(
				"Expected thought relation %s, got %s",
				request.Thought.Relation,
				receivedRequest.Thought.Relation,
			)
		}

		if receivedRequest.Thought.Module.Reference != request.Thought.Module.Reference {
			t.Errorf(
				"Expected module reference %s, got %s",
				request.Thought.Module.Reference,
				receivedRequest.Thought.Module.Reference,
			)
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
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionCreateThoughtWithOutputClassID(t *testing.T) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation:      "description",
			OutputClassID: "class-456",
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
				Parameters: map[string]any{
					"temperature": 0.9,
					"max_tokens":  200,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-789",
		ChainID:       "chain-123",
		OutputClassID: "class-456",
		Module: &perception.Module{
			ID:        "module-789",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.9,
				"max_tokens":  200,
			},
		},
		ProvisionState: "pending",
		Relation:       "description",
		Index:          3,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Thought.OutputClassID != request.Thought.OutputClassID {
			t.Errorf(
				"Expected output_class_id %s, got %s",
				request.Thought.OutputClassID,
				receivedRequest.Thought.OutputClassID,
			)
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
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionCreateThoughtWithIndex(t *testing.T) {
	indexValue := 5
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Index:    &indexValue,
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
				Parameters: map[string]any{
					"temperature": 0.6,
					"max_tokens":  120,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-101",
		ChainID:       "chain-123",
		OutputClassID: "class-123",
		Module: &perception.Module{
			ID:        "module-101",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.6,
				"max_tokens":  120,
			},
		},
		ProvisionState: "pending",
		Relation:       "description",
		Index:          5,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Validate that the index is properly included in the request body
		if receivedRequest.Thought.Index == nil || *receivedRequest.Thought.Index != *request.Thought.Index {
			var receivedIndex int
			if receivedRequest.Thought.Index != nil {
				receivedIndex = *receivedRequest.Thought.Index
			}
			t.Errorf(
				"Expected thought index %d, got %d",
				*request.Thought.Index,
				receivedIndex,
			)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf(
				"Expected thought relation %s, got %s",
				request.Thought.Relation,
				receivedRequest.Thought.Relation,
			)
		}

		if receivedRequest.Thought.Module.Reference != request.Thought.Module.Reference {
			t.Errorf(
				"Expected module reference %s, got %s",
				request.Thought.Module.Reference,
				receivedRequest.Thought.Module.Reference,
			)
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
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionCreateThoughtValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty chain ID
	_, err = client.Perception.CreateThought("", perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Module: &perception.Module{
				Reference:  "tama/agentic/generate",
				Parameters: map[string]any{},
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty chain ID")
	}

	// Test empty relation
	_, err = client.Perception.CreateThought("chain-123", perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "",
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty relation")
	}

	// Test empty module reference
	_, err = client.Perception.CreateThought("chain-123", perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Module: &perception.Module{
				Reference: "",
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty module reference")
	}

	// Test faculty with empty queue ID
	_, err = client.Perception.CreateThought("chain-123", perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
			},
			Faculty: &perception.Faculty{
				QueueID:  "",
				Priority: 1,
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty faculty queue ID")
	}
}

func TestPerceptionUpdateThought(t *testing.T) {
	indexValue := 3
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation: "updated-description",
			Index:    &indexValue,
			Module: &perception.Module{
				Reference: "tama/agentic/analyze",
				Parameters: map[string]any{
					"depth": 3,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-111",
		ChainID:       "chain-123",
		OutputClassID: "class-789",
		Module: &perception.Module{
			ID:        "module-111",
			Reference: "tama/agentic/analyze",
			Parameters: map[string]any{
				"depth": 3,
				"scope": "global",
			},
		},
		ProvisionState: "active",
		Relation:       "analysis",
		Index:          7,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123, got %s", r.URL.Path)
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
	thought, err := client.Perception.UpdateThought("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thought.Relation != expectedThought.Relation {
		t.Errorf("Expected thought relation %s, got %s", expectedThought.Relation, thought.Relation)
	}

	if thought.Module.Reference != expectedThought.Module.Reference {
		t.Errorf(
			"Expected module reference %s, got %s",
			expectedThought.Module.Reference,
			thought.Module.Reference,
		)
	}
}

//nolint:gocognit
func TestPerceptionUpdateThoughtWithIndex(t *testing.T) {
	indexValue := 3
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation:      "updated-description",
			OutputClassID: "class-789",
			Index:         &indexValue,
			Module: &perception.Module{
				Reference: "tama/agentic/analyze",
				Parameters: map[string]any{
					"depth": 3,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-222",
		ChainID:       "chain-123",
		OutputClassID: "class-456",
		Module: &perception.Module{
			ID:        "module-222",
			Reference: "tama/agentic/transform",
			Parameters: map[string]any{
				"format": "json",
				"schema": "v2",
			},
		},
		ProvisionState: "active",
		Relation:       "transformation",
		Index:          10,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123, got %s", r.URL.Path)
		}

		var receivedRequest perception.UpdateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Validate that the index is properly included in the request body
		if receivedRequest.Thought.Index == nil || *receivedRequest.Thought.Index != *request.Thought.Index {
			var receivedIndex int
			if receivedRequest.Thought.Index != nil {
				receivedIndex = *receivedRequest.Thought.Index
			}
			t.Errorf(
				"Expected thought index %d, got %d",
				*request.Thought.Index,
				receivedIndex,
			)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf(
				"Expected thought relation %s, got %s",
				request.Thought.Relation,
				receivedRequest.Thought.Relation,
			)
		}

		if receivedRequest.Thought.Module.Reference != request.Thought.Module.Reference {
			t.Errorf(
				"Expected module reference %s, got %s",
				request.Thought.Module.Reference,
				receivedRequest.Thought.Module.Reference,
			)
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
	thought, err := client.Perception.UpdateThought("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thought.Relation != expectedThought.Relation {
		t.Errorf("Expected thought relation %s, got %s", expectedThought.Relation, thought.Relation)
	}

	if thought.Module.Reference != expectedThought.Module.Reference {
		t.Errorf(
			"Expected module reference %s, got %s",
			expectedThought.Module.Reference,
			thought.Module.Reference,
		)
	}

	if thought.Index != expectedThought.Index {
		t.Errorf("Expected thought index %d, got %d", expectedThought.Index, thought.Index)
	}
}

func TestPerceptionCreateThoughtWithZeroIndex(t *testing.T) {
	zeroIndex := 0
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Index:    &zeroIndex, // Explicitly set zero index
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
				Parameters: map[string]any{
					"temperature": 0.8,
					"max_tokens":  150,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-456",
		ChainID:       "chain-123",
		OutputClassID: "class-123",
		Module: &perception.Module{
			ID:        "module-456",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.8,
				"max_tokens":  150,
			},
		},
		ProvisionState: "pending",
		Relation:       "description",
		Index:          0,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	requestCount := 0
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		requestCount++

		// Validate index based on request count
		ValidateRequestIndex(t, &receivedRequest, requestCount)

		if receivedRequest.Thought.Relation != "description" {
			t.Errorf(
				"Expected thought relation description, got %s",
				receivedRequest.Thought.Relation,
			)
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
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thought.Index != 0 {
		t.Errorf("Expected thought index 0, got %d", thought.Index)
	}

	// Test with a nil index
	requestNilIndex := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Index:    nil, // Do not set the index
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
				Parameters: map[string]any{
					"temperature": 0.6,
					"max_tokens":  100,
				},
			},
		},
	}

	_, err = client.Perception.CreateThought("chain-123", requestNilIndex)
	if err != nil {
		t.Fatalf("Expected no error for nil index, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionOptionalIndexBehaviorSimple(t *testing.T) {
	testIndexBehavior(t, "DelegatedThoughtWithoutIndex", nil, "thought-nil-index", "target-thought-123", 5)

	zeroIndex := 0
	testIndexBehavior(t, "DelegatedThoughtWithZeroIndex", &zeroIndex, "thought-zero-index", "target-thought-456", 0)

	positiveIndex := 3
	testIndexBehavior(t, "DelegatedThoughtWithPositiveIndex", &positiveIndex,
		"thought-positive-index", "target-thought-789", 3)
}

func TestPerceptionOptionalIndexBehavior(t *testing.T) {
	testIndexBehavior(t, "DelegatedThoughtWithoutIndex", nil, "thought-nil-index", "target-thought-123", 5)

	zeroIndex := 0
	testIndexBehavior(t, "DelegatedThoughtWithZeroIndex", &zeroIndex, "thought-zero-index", "target-thought-456", 0)

	positiveIndex := 3
	testIndexBehavior(t, "DelegatedThoughtWithPositiveIndex", &positiveIndex,
		"thought-positive-index", "target-thought-789", 3)
}

func TestPerceptionDeleteThought(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123, got %s", r.URL.Path)
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
	err = client.Perception.DeleteThought("thought-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetThoughtEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.GetThought("")
	if err == nil {
		t.Error("Expected validation error for empty thought ID in GetThought")
	}
}

func TestPerceptionUpdateThoughtEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.UpdateThought("", perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{Relation: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought ID in UpdateThought")
	}
}

func TestPerceptionDeleteThoughtEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Perception.DeleteThought("")
	if err == nil {
		t.Error("Expected validation error for empty thought ID in DeleteThought")
	}
}

func TestPerceptionUpdateThoughtWithZeroIndex(t *testing.T) {
	zeroIndex := 0
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation: "updated-description",
			Index:    &zeroIndex, // Explicitly set zero index
			Module: &perception.Module{
				Reference: "tama/agentic/analyze",
				Parameters: map[string]any{
					"depth": 3,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-123",
		ChainID:       "chain-123",
		OutputClassID: "class-789",
		Module: &perception.Module{
			ID:        "module-123",
			Reference: "tama/agentic/analyze",
			Parameters: map[string]any{
				"depth": 3,
			},
		},
		ProvisionState: "active",
		Relation:       "updated-description",
		Index:          0,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123, got %s", r.URL.Path)
		}

		var receivedRequest perception.UpdateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Validate that zero index is properly included in the request body
		if receivedRequest.Thought.Index == nil || *receivedRequest.Thought.Index != 0 {
			var receivedIndex int
			if receivedRequest.Thought.Index != nil {
				receivedIndex = *receivedRequest.Thought.Index
			}
			t.Errorf("Expected thought index 0, got %d", receivedIndex)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf(
				"Expected thought relation %s, got %s",
				request.Thought.Relation,
				receivedRequest.Thought.Relation,
			)
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
	thought, err := client.Perception.UpdateThought("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thought.Index != 0 {
		t.Errorf("Expected thought index 0, got %d", thought.Index)
	}

	if thought.Relation != expectedThought.Relation {
		t.Errorf("Expected thought relation %s, got %s", expectedThought.Relation, thought.Relation)
	}
}

func TestPerceptionCreateThoughtWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		errorResponse := map[string]any{
			"errors": map[string][]string{
				"relation": {"can't be blank", "is not included in the list"},
				"module":   {"reference is invalid"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
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
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "invalid-relation", // Valid relation to bypass client validation
			Module: &perception.Module{
				Reference:  "invalid/reference", // Valid reference to bypass client validation
				Parameters: map[string]any{},
			},
		},
	}

	_, err = client.Perception.CreateThought("chain-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "relation can't be blank") {
		t.Errorf("Expected error to contain 'relation can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "relation is not included in the list") {
		t.Errorf(
			"Expected error to contain 'relation is not included in the list', got %s",
			errorMsg,
		)
	}
	if !strings.Contains(errorMsg, "module reference is invalid") {
		t.Errorf("Expected error to contain 'module reference is invalid', got %s", errorMsg)
	}
}

func TestPerceptionCreateThoughtWithDelegation(t *testing.T) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "delegation",
			Delegation: &perception.Delegation{
				TargetThoughtID: "target-thought-123",
			},
		},
	}

	expectedThought := perception.Thought{
		ID:      "thought-delegation-456",
		ChainID: "chain-123",
		Delegation: &perception.Delegation{
			TargetThoughtID: "target-thought-123",
		},
		ProvisionState: "pending",
		Relation:       "delegation",
		Index:          1,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf(
				"Expected thought relation %s, got %s",
				request.Thought.Relation,
				receivedRequest.Thought.Relation,
			)
		}

		if receivedRequest.Thought.Delegation == nil {
			t.Error("Expected delegation to be present")
		} else if receivedRequest.Thought.Delegation.TargetThoughtID != request.Thought.Delegation.TargetThoughtID {
			t.Errorf(
				"Expected delegation target_thought_id %s, got %s",
				request.Thought.Delegation.TargetThoughtID,
				receivedRequest.Thought.Delegation.TargetThoughtID,
			)
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
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionCreateThoughtWithFaculty(t *testing.T) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "faculty",
			Module: &perception.Module{
				Reference: "tama/agentic/generate",
			},
			Faculty: &perception.Faculty{
				QueueID:  "queue-123",
				Priority: 0,
			},
		},
	}

	expectedThought := perception.Thought{
		ID:      "thought-faculty-456",
		ChainID: "chain-123",
		Module: &perception.Module{
			ID:        "module-456",
			Reference: "tama/agentic/generate",
		},
		Faculty: &perception.Faculty{
			QueueID:  "queue-123",
			Priority: 0,
		},
		ProvisionState: "pending",
		Relation:       "faculty",
		Index:          2,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Thought.Faculty == nil {
			t.Fatal("Expected faculty to be present in the request")
		}

		if receivedRequest.Thought.Faculty.QueueID != request.Thought.Faculty.QueueID {
			t.Errorf(
				"Expected faculty queue_id %s, got %s",
				request.Thought.Faculty.QueueID,
				receivedRequest.Thought.Faculty.QueueID,
			)
		}

		if receivedRequest.Thought.Faculty.Priority != request.Thought.Faculty.Priority {
			t.Errorf(
				"Expected faculty priority %d, got %d",
				request.Thought.Faculty.Priority,
				receivedRequest.Thought.Faculty.Priority,
			)
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
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionUpdateThoughtWithDelegation(t *testing.T) {
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation: "updated-delegation",
			Delegation: &perception.Delegation{
				TargetThoughtID: "new-target-thought-789",
			},
		},
	}

	expectedThought := perception.Thought{
		ID:      "thought-123",
		ChainID: "chain-123",
		Delegation: &perception.Delegation{
			TargetThoughtID: "new-target-thought-789",
		},
		ProvisionState: "active",
		Relation:       "updated-delegation",
		Index:          1,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123, got %s", r.URL.Path)
		}

		var receivedRequest perception.UpdateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf(
				"Expected thought relation %s, got %s",
				request.Thought.Relation,
				receivedRequest.Thought.Relation,
			)
		}

		if receivedRequest.Thought.Delegation == nil {
			t.Error("Expected delegation to be present")
		} else if receivedRequest.Thought.Delegation.TargetThoughtID != request.Thought.Delegation.TargetThoughtID {
			t.Errorf(
				"Expected delegation target_thought_id %s, got %s",
				request.Thought.Delegation.TargetThoughtID,
				receivedRequest.Thought.Delegation.TargetThoughtID,
			)
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
	thought, err := client.Perception.UpdateThought("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionUpdateThoughtWithFaculty(t *testing.T) {
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation: "updated-faculty",
			Faculty: &perception.Faculty{
				QueueID:  "queue-789",
				Priority: 3,
			},
		},
	}

	expectedThought := perception.Thought{
		ID:      "thought-123",
		ChainID: "chain-123",
		Faculty: &perception.Faculty{
			QueueID:  "queue-789",
			Priority: 3,
		},
		ProvisionState: "active",
		Relation:       "updated-faculty",
		Index:          1,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123, got %s", r.URL.Path)
		}

		var receivedRequest perception.UpdateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Thought.Faculty == nil {
			t.Fatal("Expected faculty to be present in the request")
		}

		if receivedRequest.Thought.Faculty.QueueID != request.Thought.Faculty.QueueID {
			t.Errorf(
				"Expected faculty queue_id %s, got %s",
				request.Thought.Faculty.QueueID,
				receivedRequest.Thought.Faculty.QueueID,
			)
		}

		if receivedRequest.Thought.Faculty.Priority != request.Thought.Faculty.Priority {
			t.Errorf(
				"Expected faculty priority %d, got %d",
				request.Thought.Faculty.Priority,
				receivedRequest.Thought.Faculty.Priority,
			)
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
	thought, err := client.Perception.UpdateThought("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateThoughtResponse(t, *thought, expectedThought)
}

//nolint:gocognit
func TestPerceptionNestedErrorParsing(t *testing.T) {
	// Test API response with nested validation errors (e.g., module.reference)
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf(
				"Expected path /provision/perception/chains/chain-123/thoughts, got %s",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		// Nested error structure like the real API returns
		errorResponse := map[string]any{
			"errors": map[string]any{
				"module": map[string]any{
					"reference": []string{"is invalid"},
					"parameters": map[string]any{
						"temperature": []string{"must be between 0 and 1"},
						"config": map[string]any{
							"database": map[string]any{
								"connection": []string{"cannot be established"},
							},
						},
					},
				},
				"relation": []string{"can't be blank"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
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
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description", // Valid to bypass client validation
			Module: &perception.Module{
				Reference:  "tama/some/invalid", // Valid format to bypass client validation
				Parameters: map[string]any{},
			},
		},
	}

	_, err = client.Perception.CreateThought("chain-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that it's our custom perception error type
	var perceptionErr *perception.Error
	if !errors.As(err, &perceptionErr) {
		t.Fatalf("Expected *perception.Error, got %T", err)
	}

	// Verify status code
	if perceptionErr.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf(
			"Expected status code %d, got %d",
			http.StatusUnprocessableEntity,
			perceptionErr.StatusCode,
		)
	}

	// Check that nested fields are flattened with dot notation
	expectedFields := map[string][]string{
		"module.reference":                             {"is invalid"},
		"module.parameters.temperature":                {"must be between 0 and 1"},
		"module.parameters.config.database.connection": {"cannot be established"},
		"relation": {"can't be blank"},
	}

	for expectedField, expectedMessages := range expectedFields {
		actualMessages, exists := perceptionErr.Errors[expectedField]
		if !exists {
			t.Errorf("Expected field '%s' not found in error. Available fields: %v",
				expectedField, GetKeys(perceptionErr.Errors))
			continue
		}

		if len(actualMessages) != len(expectedMessages) {
			t.Errorf("Field '%s': expected %d messages, got %d",
				expectedField, len(expectedMessages), len(actualMessages))
			continue
		}

		for i, expectedMsg := range expectedMessages {
			if actualMessages[i] != expectedMsg {
				t.Errorf("Field '%s' message %d: expected '%s', got '%s'",
					expectedField, i, expectedMsg, actualMessages[i])
			}
		}
	}

	// Check the complete error message
	errorMsg := perceptionErr.Error()
	expectedSubstrings := []string{
		"API error 422:",
		"module.reference is invalid",
		"module.parameters.temperature must be between 0 and 1",
		"module.parameters.config.database.connection cannot be established",
		"relation can't be blank",
	}

	for _, substring := range expectedSubstrings {
		if !strings.Contains(errorMsg, substring) {
			t.Errorf("Expected error message to contain '%s', got: %s", substring, errorMsg)
		}
	}
}

// Helper functions for thought tests

func createIndexTestData(
	index *int, thoughtID, targetID string, expectedIndex int,
) (perception.CreateThoughtRequest, perception.ThoughtResponse) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "delegation",
			Index:    index,
			Delegation: &perception.Delegation{
				TargetThoughtID: targetID,
			},
		},
	}

	expectedThought := perception.Thought{
		ID:      thoughtID,
		ChainID: "chain-123",
		Delegation: &perception.Delegation{
			TargetThoughtID: targetID,
		},
		ProvisionState: "pending",
		Relation:       "delegation",
		Index:          expectedIndex,
	}

	return request, perception.ThoughtResponse{Data: expectedThought}
}

func validateIndexRequest(t *testing.T, r *http.Request, expectedIndex *int, expectedTargetID string) {
	if r.Method != http.MethodPost {
		t.Errorf("Expected POST request, got %s", r.Method)
	}

	var receivedRequest perception.CreateThoughtRequest
	if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
		t.Fatalf("Failed to decode request body: %v", err)
	}

	if expectedIndex == nil {
		if receivedRequest.Thought.Index != nil {
			t.Errorf("Expected thought index to be nil (omitted), got %d", *receivedRequest.Thought.Index)
		}
	} else {
		if receivedRequest.Thought.Index == nil {
			t.Error("Expected thought index to be present (not nil), but got nil")
		} else if *receivedRequest.Thought.Index != *expectedIndex {
			t.Errorf("Expected thought index %d, got %d", *expectedIndex, *receivedRequest.Thought.Index)
		}
	}

	if receivedRequest.Thought.Delegation == nil {
		t.Error("Expected delegation to be present")
	} else if receivedRequest.Thought.Delegation.TargetThoughtID != expectedTargetID {
		t.Errorf(
			"Expected target thought ID %s, got %s",
			expectedTargetID,
			receivedRequest.Thought.Delegation.TargetThoughtID,
		)
	}
}

func testIndexBehavior(t *testing.T, testName string, index *int, thoughtID, targetID string, expectedIndex int) {
	t.Run(testName, func(t *testing.T) {
		request, expectedResponse := createIndexTestData(index, thoughtID, targetID, expectedIndex)

		server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
			validateIndexRequest(t, r, index, targetID)
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
		thought, err := client.Perception.CreateThought("chain-123", request)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if thought.Index != expectedIndex {
			t.Errorf("Expected index %d, got %d", expectedIndex, thought.Index)
		}
	})
}
