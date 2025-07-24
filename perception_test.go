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
