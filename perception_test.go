package tama_test

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

func TestPerceptionGetChain(t *testing.T) {
	expectedChain := perception.Chain{
		ID:             "chain-123",
		SpaceID:        "space-123",
		Name:           "test-chain",
		Slug:           "test-chain-slug",
		ProvisionState: "active",
	}

	expectedResponse := perception.ChainResponse{
		Data: expectedChain,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123" {
			t.Errorf("Expected path /provision/perception/chains/chain-123, got %s", r.URL.Path)
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
	chain, err := client.Perception.GetChain("chain-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if chain.ID != expectedChain.ID {
		t.Errorf("Expected chain ID %s, got %s", expectedChain.ID, chain.ID)
	}

	if chain.Name != expectedChain.Name {
		t.Errorf("Expected chain name %s, got %s", expectedChain.Name, chain.Name)
	}

	if chain.SpaceID != expectedChain.SpaceID {
		t.Errorf("Expected chain space_id %s, got %s", expectedChain.SpaceID, chain.SpaceID)
	}

	if chain.Slug != expectedChain.Slug {
		t.Errorf("Expected chain slug %s, got %s", expectedChain.Slug, chain.Slug)
	}

	if chain.ProvisionState != expectedChain.ProvisionState {
		t.Errorf("Expected chain provision_state %s, got %s", expectedChain.ProvisionState, chain.ProvisionState)
	}
}

func TestPerceptionGetChainError(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := perception.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"chain": {"not found"},
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
	_, err := client.Perception.GetChain("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var perceptionErr *perception.Error
	if errors.As(err, &perceptionErr) {
		if perceptionErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", perceptionErr.StatusCode)
		}
		if perceptionErr.Errors == nil || len(perceptionErr.Errors["chain"]) == 0 ||
			perceptionErr.Errors["chain"][0] != "not found" {
			t.Errorf("Expected error 'chain not found', got %v", perceptionErr.Errors)
		}
	} else {
		t.Errorf("Expected perception.Error, got %T", err)
	}
}

func TestPerceptionCreateChain(t *testing.T) {
	request := perception.CreateChainRequest{
		Chain: perception.ChainRequestData{
			Name: "new-chain",
		},
	}

	expectedChain := perception.Chain{
		ID:             "chain-456",
		SpaceID:        "space-123",
		Name:           "new-chain",
		Slug:           "new-chain-slug",
		ProvisionState: "pending",
	}

	expectedResponse := perception.ChainResponse{
		Data: expectedChain,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/spaces/space-123/chains" {
			t.Errorf("Expected path /provision/perception/spaces/space-123/chains, got %s", r.URL.Path)
		}

		var receivedRequest perception.CreateChainRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Chain.Name != request.Chain.Name {
			t.Errorf("Expected chain name %s, got %s", request.Chain.Name, receivedRequest.Chain.Name)
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
	chain, err := client.Perception.CreateChain("space-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if chain.ID != expectedChain.ID {
		t.Errorf("Expected chain ID %s, got %s", expectedChain.ID, chain.ID)
	}

	if chain.Name != expectedChain.Name {
		t.Errorf("Expected chain name %s, got %s", expectedChain.Name, chain.Name)
	}

	if chain.SpaceID != expectedChain.SpaceID {
		t.Errorf("Expected chain space_id %s, got %s", expectedChain.SpaceID, chain.SpaceID)
	}
}

func TestPerceptionCreateChainValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty space ID
	_, err := client.Perception.CreateChain("", perception.CreateChainRequest{
		Chain: perception.ChainRequestData{Name: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty chain name
	_, err = client.Perception.CreateChain("space-123", perception.CreateChainRequest{
		Chain: perception.ChainRequestData{Name: ""},
	})
	if err == nil {
		t.Error("Expected validation error for empty chain name")
	}
}

func TestPerceptionCreateChainNameValidationDelegated(t *testing.T) {
	// Test that API response with name validation errors is handled correctly
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResp := perception.Error{
			StatusCode: 422,
			Errors: map[string][]string{
				"name": {"is required", "must be at least 3 characters"},
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
	_, err := client.Perception.CreateChain("space-123", perception.CreateChainRequest{
		Chain: perception.ChainRequestData{Name: "ab"},
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var perceptionErr *perception.Error
	if errors.As(err, &perceptionErr) {
		if perceptionErr.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code 422, got %d", perceptionErr.StatusCode)
		}
		if len(perceptionErr.Errors["name"]) != 2 {
			t.Errorf("Expected 2 name errors, got %d", len(perceptionErr.Errors["name"]))
		}
	} else {
		t.Errorf("Expected perception.Error, got %T", err)
	}
}

func TestPerceptionUpdateChain(t *testing.T) {
	request := perception.UpdateChainRequest{
		Chain: perception.UpdateChainData{
			Name: "updated-chain",
		},
	}

	expectedChain := perception.Chain{
		ID:             "chain-123",
		SpaceID:        "space-123",
		Name:           "updated-chain",
		Slug:           "updated-chain-slug",
		ProvisionState: "active",
	}

	expectedResponse := perception.ChainResponse{
		Data: expectedChain,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123" {
			t.Errorf("Expected path /provision/perception/chains/chain-123, got %s", r.URL.Path)
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
	chain, err := client.Perception.UpdateChain("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if chain.Name != expectedChain.Name {
		t.Errorf("Expected chain name %s, got %s", expectedChain.Name, chain.Name)
	}
}

func TestPerceptionReplaceChain(t *testing.T) {
	request := perception.UpdateChainRequest{
		Chain: perception.UpdateChainData{
			Name: "replaced-chain",
		},
	}

	expectedChain := perception.Chain{
		ID:             "chain-123",
		SpaceID:        "space-123",
		Name:           "replaced-chain",
		Slug:           "replaced-chain-slug",
		ProvisionState: "active",
	}

	expectedResponse := perception.ChainResponse{
		Data: expectedChain,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123" {
			t.Errorf("Expected path /provision/perception/chains/chain-123, got %s", r.URL.Path)
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
	chain, err := client.Perception.ReplaceChain("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if chain.Name != expectedChain.Name {
		t.Errorf("Expected chain name %s, got %s", expectedChain.Name, chain.Name)
	}
}

func TestPerceptionDeleteChain(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123" {
			t.Errorf("Expected path /provision/perception/chains/chain-123, got %s", r.URL.Path)
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
	err := client.Perception.DeleteChain("chain-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetChainEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.GetChain("")
	if err == nil {
		t.Error("Expected validation error for empty chain ID in GetChain")
	}
}

func TestPerceptionUpdateChainEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.UpdateChain("", perception.UpdateChainRequest{
		Chain: perception.UpdateChainData{Name: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty chain ID in UpdateChain")
	}
}

func TestPerceptionReplaceChainEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.ReplaceChain("", perception.UpdateChainRequest{
		Chain: perception.UpdateChainData{Name: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty chain ID in ReplaceChain")
	}
}

func TestPerceptionDeleteChainEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Perception.DeleteChain("")
	if err == nil {
		t.Error("Expected validation error for empty chain ID in DeleteChain")
	}
}

func TestPerceptionCreateChainWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/spaces/space-123/chains" {
			t.Errorf("Expected path /provision/perception/spaces/space-123/chains, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		errorResponse := map[string]any{
			"errors": map[string][]string{
				"name": {"can't be blank", "is too short (minimum is 3 characters)"},
				"base": {"space is not active"},
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
	request := perception.CreateChainRequest{
		Chain: perception.ChainRequestData{
			Name: "ab", // Short name to trigger server-side validation
		},
	}

	_, err := client.Perception.CreateChain("space-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "name can't be blank") {
		t.Errorf("Expected error to contain 'name can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "name is too short") {
		t.Errorf("Expected error to contain 'name is too short', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "base space is not active") {
		t.Errorf("Expected error to contain 'base space is not active', got %s", errorMsg)
	}
}

func TestPerceptionGetThought(t *testing.T) {
	expectedThought := perception.Thought{
		ID:            "thought-123",
		ChainID:       "chain-123",
		OutputClassID: "class-123",
		Module: perception.Module{
			ID:        "module-123",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.7,
				"max_tokens":  100,
			},
		},
		ProvisionState: "active",
		Relation:       "description",
		Index:          1,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	thought, err := client.Perception.GetThought("thought-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionGetThoughtError(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, _ *http.Request) {
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	_, err := client.Perception.GetThought("nonexistent")

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
			Module: perception.Module{
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
		Module: perception.Module{
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf("Expected path /provision/perception/chains/chain-123/thoughts, got %s", r.URL.Path)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf("Expected thought relation %s, got %s", request.Thought.Relation, receivedRequest.Thought.Relation)
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionCreateThoughtWithOutputClassID(t *testing.T) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation:      "description",
			OutputClassID: "class-456",
			Module: perception.Module{
				Reference: "tama/agentic/generate",
				Parameters: map[string]any{
					"temperature": 0.8,
					"max_tokens":  150,
				},
			},
		},
	}

	expectedThought := perception.Thought{
		ID:            "thought-789",
		ChainID:       "chain-123",
		OutputClassID: "class-456",
		Module: perception.Module{
			ID:        "module-789",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.8,
				"max_tokens":  150,
			},
		},
		ProvisionState: "pending",
		Relation:       "description",
		Index:          3,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf("Expected path /provision/perception/chains/chain-123/thoughts, got %s", r.URL.Path)
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionCreateThoughtWithIndex(t *testing.T) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Index:    5,
			Module: perception.Module{
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
		Module: perception.Module{
			ID:        "module-456",
			Reference: "tama/agentic/generate",
			Parameters: map[string]any{
				"temperature": 0.8,
				"max_tokens":  150,
			},
		},
		ProvisionState: "pending",
		Relation:       "description",
		Index:          5,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf("Expected path /provision/perception/chains/chain-123/thoughts, got %s", r.URL.Path)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Validate that the index is properly included in the request body
		if receivedRequest.Thought.Index != request.Thought.Index {
			t.Errorf("Expected thought index %d, got %d", request.Thought.Index, receivedRequest.Thought.Index)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf("Expected thought relation %s, got %s", request.Thought.Relation, receivedRequest.Thought.Relation)
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionCreateThoughtValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty chain ID
	_, err := client.Perception.CreateThought("", perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Module: perception.Module{
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
			Module: perception.Module{
				Reference:  "tama/agentic/generate",
				Parameters: map[string]any{},
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
			Module: perception.Module{
				Reference:  "",
				Parameters: map[string]any{},
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty module reference")
	}
}

func TestPerceptionUpdateThought(t *testing.T) {
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation:      "updated-description",
			OutputClassID: "class-789",
			Module: perception.Module{
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
		Module: perception.Module{
			ID:        "module-123",
			Reference: "tama/agentic/analyze",
			Parameters: map[string]any{
				"depth": 3,
			},
		},
		ProvisionState: "active",
		Relation:       "updated-description",
		Index:          1,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	thought, err := client.Perception.UpdateThought("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thought.Relation != expectedThought.Relation {
		t.Errorf("Expected thought relation %s, got %s", expectedThought.Relation, thought.Relation)
	}

	if thought.Module.Reference != expectedThought.Module.Reference {
		t.Errorf("Expected module reference %s, got %s", expectedThought.Module.Reference, thought.Module.Reference)
	}
}

func TestPerceptionUpdateThoughtWithIndex(t *testing.T) {
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation:      "updated-description",
			OutputClassID: "class-789",
			Index:         3,
			Module: perception.Module{
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
		Module: perception.Module{
			ID:        "module-123",
			Reference: "tama/agentic/analyze",
			Parameters: map[string]any{
				"depth": 3,
			},
		},
		ProvisionState: "active",
		Relation:       "updated-description",
		Index:          3,
	}

	expectedResponse := perception.ThoughtResponse{
		Data: expectedThought,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		if receivedRequest.Thought.Index != request.Thought.Index {
			t.Errorf("Expected thought index %d, got %d", request.Thought.Index, receivedRequest.Thought.Index)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf("Expected thought relation %s, got %s", request.Thought.Relation, receivedRequest.Thought.Relation)
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	thought, err := client.Perception.UpdateThought("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thought.Relation != expectedThought.Relation {
		t.Errorf("Expected thought relation %s, got %s", expectedThought.Relation, thought.Relation)
	}

	if thought.Module.Reference != expectedThought.Module.Reference {
		t.Errorf("Expected module reference %s, got %s", expectedThought.Module.Reference, thought.Module.Reference)
	}

	if thought.Index != expectedThought.Index {
		t.Errorf("Expected thought index %d, got %d", expectedThought.Index, thought.Index)
	}
}

func TestPerceptionCreateThoughtWithZeroIndex(t *testing.T) {
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description",
			Index:    0, // Explicitly set zero index
			Module: perception.Module{
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
		Module: perception.Module{
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf("Expected path /provision/perception/chains/chain-123/thoughts, got %s", r.URL.Path)
		}

		var receivedRequest perception.CreateThoughtRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Validate that zero index is properly included in the request body
		if receivedRequest.Thought.Index != 0 {
			t.Errorf("Expected thought index 0, got %d", receivedRequest.Thought.Index)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf("Expected thought relation %s, got %s", request.Thought.Relation, receivedRequest.Thought.Relation)
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
	thought, err := client.Perception.CreateThought("chain-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thought.Index != 0 {
		t.Errorf("Expected thought index 0, got %d", thought.Index)
	}

	validateThoughtResponse(t, *thought, expectedThought)
}

func TestPerceptionDeleteThought(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	err := client.Perception.DeleteThought("thought-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetThoughtEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.GetThought("")
	if err == nil {
		t.Error("Expected validation error for empty thought ID in GetThought")
	}
}

func TestPerceptionUpdateThoughtEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.UpdateThought("", perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{Relation: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought ID in UpdateThought")
	}
}

func TestPerceptionDeleteThoughtEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Perception.DeleteThought("")
	if err == nil {
		t.Error("Expected validation error for empty thought ID in DeleteThought")
	}
}

func TestPerceptionUpdateThoughtWithZeroIndex(t *testing.T) {
	request := perception.UpdateThoughtRequest{
		Thought: perception.UpdateThoughtData{
			Relation:      "updated-description",
			OutputClassID: "class-789",
			Index:         0, // Explicitly set zero index
			Module: perception.Module{
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
		Module: perception.Module{
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		if receivedRequest.Thought.Index != 0 {
			t.Errorf("Expected thought index 0, got %d", receivedRequest.Thought.Index)
		}

		if receivedRequest.Thought.Relation != request.Thought.Relation {
			t.Errorf("Expected thought relation %s, got %s", request.Thought.Relation, receivedRequest.Thought.Relation)
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
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf("Expected path /provision/perception/chains/chain-123/thoughts, got %s", r.URL.Path)
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "invalid-relation", // Valid relation to bypass client validation
			Module: perception.Module{
				Reference:  "invalid/reference", // Valid reference to bypass client validation
				Parameters: map[string]any{},
			},
		},
	}

	_, err := client.Perception.CreateThought("chain-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "relation can't be blank") {
		t.Errorf("Expected error to contain 'relation can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "relation is not included in the list") {
		t.Errorf("Expected error to contain 'relation is not included in the list', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "module reference is invalid") {
		t.Errorf("Expected error to contain 'module reference is invalid', got %s", errorMsg)
	}
}

func validateThoughtResponse(t *testing.T, actual, expected perception.Thought) {
	if actual.ID != expected.ID {
		t.Errorf("Expected thought ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ChainID != expected.ChainID {
		t.Errorf("Expected thought chain_id %s, got %s", expected.ChainID, actual.ChainID)
	}

	if actual.OutputClassID != expected.OutputClassID {
		t.Errorf("Expected thought output_class_id %s, got %s", expected.OutputClassID, actual.OutputClassID)
	}

	if actual.Module.ID != expected.Module.ID {
		t.Errorf("Expected module ID %s, got %s", expected.Module.ID, actual.Module.ID)
	}

	if actual.Module.Reference != expected.Module.Reference {
		t.Errorf("Expected module reference %s, got %s", expected.Module.Reference, actual.Module.Reference)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected thought provision_state %s, got %s", expected.ProvisionState, actual.ProvisionState)
	}

	if actual.Relation != expected.Relation {
		t.Errorf("Expected thought relation %s, got %s", expected.Relation, actual.Relation)
	}

	if actual.Index != expected.Index {
		t.Errorf("Expected thought index %d, got %d", expected.Index, actual.Index)
	}
}

func TestPerceptionNestedErrorParsing(t *testing.T) {
	// Test API response with nested validation errors (e.g., module.reference)
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/chains/chain-123/thoughts" {
			t.Errorf("Expected path /provision/perception/chains/chain-123/thoughts, got %s", r.URL.Path)
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
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	request := perception.CreateThoughtRequest{
		Thought: perception.ThoughtRequestData{
			Relation: "description", // Valid to bypass client validation
			Module: perception.Module{
				Reference:  "tama/some/invalid", // Valid format to bypass client validation
				Parameters: map[string]any{},
			},
		},
	}

	_, err := client.Perception.CreateThought("chain-123", request)

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
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, perceptionErr.StatusCode)
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
				expectedField, getKeys(perceptionErr.Errors))
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

// Helper function to get keys from a map for debugging.
func getKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Path tests.
func TestPerceptionGetPath(t *testing.T) {
	expectedPath := perception.Path{
		ID:             "path-123",
		ThoughtID:      "thought-123",
		ProvisionState: "active",
		TargetClassID:  "class-123",
		Parameters: map[string]any{
			"threshold": 0.8,
			"enabled":   true,
		},
	}

	expectedResponse := perception.PathResponse{
		Data: expectedPath,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/paths/path-123" {
			t.Errorf("Expected path /provision/perception/paths/path-123, got %s", r.URL.Path)
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
	path, err := client.Perception.GetPath("path-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validatePathResponse(t, *path, expectedPath)
}

func TestPerceptionCreatePath(t *testing.T) {
	request := perception.CreatePathRequest{
		Path: perception.PathRequestData{
			TargetClassID: "class-456",
			Parameters: map[string]any{
				"threshold": 0.9,
				"enabled":   false,
			},
		},
	}

	expectedPath := perception.Path{
		ID:             "path-456",
		ThoughtID:      "thought-123",
		ProvisionState: "pending",
		TargetClassID:  "class-456",
		Parameters: map[string]any{
			"threshold": 0.9,
			"enabled":   false,
		},
	}

	expectedResponse := perception.PathResponse{
		Data: expectedPath,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/paths" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123/paths, got %s", r.URL.Path)
		}

		var receivedRequest perception.CreatePathRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Path.TargetClassID != request.Path.TargetClassID {
			t.Errorf("Expected target_class_id %s, got %s",
				request.Path.TargetClassID, receivedRequest.Path.TargetClassID)
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
	path, err := client.Perception.CreatePath("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validatePathResponse(t, *path, expectedPath)
}

func TestPerceptionCreatePathValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty thought ID
	_, err := client.Perception.CreatePath("", perception.CreatePathRequest{
		Path: perception.PathRequestData{TargetClassID: "class-123"},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought ID")
	}

	// Test empty target class ID
	_, err = client.Perception.CreatePath("thought-123", perception.CreatePathRequest{
		Path: perception.PathRequestData{TargetClassID: ""},
	})
	if err == nil {
		t.Error("Expected validation error for empty target class ID")
	}
}

func TestPerceptionUpdatePath(t *testing.T) {
	request := perception.UpdatePathRequest{
		Path: perception.UpdatePathData{
			TargetClassID: "class-789",
			Parameters: map[string]any{
				"threshold": 0.7,
				"enabled":   true,
			},
		},
	}

	expectedPath := perception.Path{
		ID:             "path-123",
		ThoughtID:      "thought-123",
		ProvisionState: "active",
		TargetClassID:  "class-789",
		Parameters: map[string]any{
			"threshold": 0.7,
			"enabled":   true,
		},
	}

	expectedResponse := perception.PathResponse{
		Data: expectedPath,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/paths/path-123" {
			t.Errorf("Expected path /provision/perception/paths/path-123, got %s", r.URL.Path)
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
	path, err := client.Perception.UpdatePath("path-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if path.TargetClassID != expectedPath.TargetClassID {
		t.Errorf("Expected target_class_id %s, got %s", expectedPath.TargetClassID, path.TargetClassID)
	}
}

func TestPerceptionReplacePath(t *testing.T) {
	request := perception.UpdatePathRequest{
		Path: perception.UpdatePathData{
			TargetClassID: "class-replaced",
			Parameters: map[string]any{
				"threshold": 0.5,
			},
		},
	}

	expectedPath := perception.Path{
		ID:             "path-123",
		ThoughtID:      "thought-123",
		ProvisionState: "active",
		TargetClassID:  "class-replaced",
		Parameters: map[string]any{
			"threshold": 0.5,
		},
	}

	expectedResponse := perception.PathResponse{
		Data: expectedPath,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/paths/path-123" {
			t.Errorf("Expected path /provision/perception/paths/path-123, got %s", r.URL.Path)
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
	path, err := client.Perception.ReplacePath("path-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if path.TargetClassID != expectedPath.TargetClassID {
		t.Errorf("Expected target_class_id %s, got %s", expectedPath.TargetClassID, path.TargetClassID)
	}
}

func TestPerceptionDeletePath(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/paths/path-123" {
			t.Errorf("Expected path /provision/perception/paths/path-123, got %s", r.URL.Path)
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
	err := client.Perception.DeletePath("path-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetPathEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.GetPath("")
	if err == nil {
		t.Error("Expected validation error for empty path ID in GetPath")
	}
}

func TestPerceptionUpdatePathEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.UpdatePath("", perception.UpdatePathRequest{
		Path: perception.UpdatePathData{TargetClassID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty path ID in UpdatePath")
	}
}

func TestPerceptionReplacePathEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.ReplacePath("", perception.UpdatePathRequest{
		Path: perception.UpdatePathData{TargetClassID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty path ID in ReplacePath")
	}
}

func TestPerceptionDeletePathEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Perception.DeletePath("")
	if err == nil {
		t.Error("Expected validation error for empty path ID in DeletePath")
	}
}

func validatePathResponse(t *testing.T, actual, expected perception.Path) {
	if actual.ID != expected.ID {
		t.Errorf("Expected path ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ThoughtID != expected.ThoughtID {
		t.Errorf("Expected path thought_id %s, got %s", expected.ThoughtID, actual.ThoughtID)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected path provision_state %s, got %s", expected.ProvisionState, actual.ProvisionState)
	}

	if actual.TargetClassID != expected.TargetClassID {
		t.Errorf("Expected path target_class_id %s, got %s", expected.TargetClassID, actual.TargetClassID)
	}
}

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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

	validateContextResponse(t, *context, expectedContext)
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/contexts" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123/contexts, got %s", r.URL.Path)
		}

		var receivedRequest perception.CreateContextRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Context.PromptID != request.Context.PromptID {
			t.Errorf("Expected prompt_id %s, got %s", request.Context.PromptID, receivedRequest.Context.PromptID)
		}

		if receivedRequest.Context.Layer != request.Context.Layer {
			t.Errorf("Expected layer %d, got %d", request.Context.Layer, receivedRequest.Context.Layer)
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

	validateContextResponse(t, *context, expectedContext)
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

func validateContextResponse(t *testing.T, actual, expected perception.Context) {
	if actual.ID != expected.ID {
		t.Errorf("Expected context ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ThoughtID != expected.ThoughtID {
		t.Errorf("Expected context thought_id %s, got %s", expected.ThoughtID, actual.ThoughtID)
	}

	if actual.PromptID != expected.PromptID {
		t.Errorf("Expected context prompt_id %s, got %s", expected.PromptID, actual.PromptID)
	}

	if actual.Layer != expected.Layer {
		t.Errorf("Expected context layer %d, got %d", expected.Layer, actual.Layer)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected context provision_state %s, got %s", expected.ProvisionState, actual.ProvisionState)
	}
}

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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
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

	if int(processor.Configuration["batch_size"].(float64)) != expectedProcessor.Configuration["batch_size"].(int) {
		t.Errorf(
			"Expected batch_size %v, got %v",
			expectedProcessor.Configuration["batch_size"],
			processor.Configuration["batch_size"],
		)
	}
}

func TestPerceptionGetProcessorError(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, _ *http.Request) {
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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

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
	_, err = client.Perception.CreateProcessor("thought-123", "completion", perception.CreateProcessorRequest{
		Processor: perception.ProcessorRequestData{
			Configuration: map[string]any{"test": "value"},
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty model ID")
	}

	// Test valid model ID
	_, err = client.Perception.CreateProcessor("thought-123", "completion", perception.CreateProcessorRequest{
		Processor: perception.ProcessorRequestData{
			ModelID:       "valid-model-123",
			Configuration: map[string]any{"test": "value"},
		},
	})
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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

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

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

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
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	err := client.Perception.DeleteProcessor("thought-123", "completion")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetProcessorEmptyThoughtID(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	_, err := client.Perception.GetProcessor("", "completion")

	if err == nil || !strings.Contains(err.Error(), "thought ID is required") {
		t.Errorf("Expected 'thought ID is required' error, got %v", err)
	}
}

func TestPerceptionGetProcessorEmptyType(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	_, err := client.Perception.GetProcessor("thought-123", "")

	if err == nil || !strings.Contains(err.Error(), "processor type is required") {
		t.Errorf("Expected 'processor type is required' error, got %v", err)
	}
}

func TestPerceptionUpdateProcessorEmptyThoughtID(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

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
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	err := client.Perception.DeleteProcessor("", "completion")

	if err == nil || !strings.Contains(err.Error(), "thought ID is required") {
		t.Errorf("Expected 'thought ID is required' error, got %v", err)
	}
}
