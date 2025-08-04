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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
		t.Errorf(
			"Expected chain provision_state %s, got %s",
			expectedChain.ProvisionState,
			chain.ProvisionState,
		)
	}
}

func TestPerceptionGetChainError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
