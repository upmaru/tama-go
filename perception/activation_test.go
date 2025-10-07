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

func TestPerceptionGetActivation(t *testing.T) {
	expectedActivation := perception.Activation{
		ID:             "activation-123",
		ThoughtPathID:  "path-123",
		ChainID:        "chain-123",
		ProvisionState: "active",
	}

	expectedResponse := perception.ActivationResponse{
		Data: expectedActivation,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/activations/activation-123" {
			t.Errorf("Expected path /provision/perception/activations/activation-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	activation, err := client.Perception.GetActivation("activation-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateActivationResponse(t, *activation, expectedActivation)
}

func TestPerceptionGetActivationError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := perception.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"activation": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Perception.GetActivation("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var perceptionErr *perception.Error
	if errors.As(err, &perceptionErr) {
		if perceptionErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", perceptionErr.StatusCode)
		}
		if perceptionErr.Errors == nil || len(perceptionErr.Errors["activation"]) == 0 ||
			perceptionErr.Errors["activation"][0] != "not found" {
			t.Errorf("Expected error 'activation not found', got %v", perceptionErr.Errors)
		}
	} else {
		t.Errorf("Expected perception.Error, got %T", err)
	}
}

func TestPerceptionCreateActivation(t *testing.T) {
	request := perception.CreateActivationRequest{
		Activation: perception.ActivationRequestData{
			ChainID: "chain-123",
		},
	}

	expectedActivation := perception.Activation{
		ID:             "activation-456",
		ThoughtPathID:  "path-123",
		ChainID:        "chain-123",
		ProvisionState: "pending",
	}

	expectedResponse := perception.ActivationResponse{
		Data: expectedActivation,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/paths/path-123/activations" {
			t.Errorf("Expected path /provision/perception/paths/path-123/activations, got %s", r.URL.Path)
		}

		var receivedRequest perception.CreateActivationRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Activation.ChainID != request.Activation.ChainID {
			t.Errorf("Expected chain ID %s, got %s", request.Activation.ChainID, receivedRequest.Activation.ChainID)
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
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	activation, err := client.Perception.CreateActivation("path-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateActivationResponse(t, *activation, expectedActivation)
}

func TestPerceptionCreateActivationValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty path ID
	_, err = client.Perception.CreateActivation("", perception.CreateActivationRequest{
		Activation: perception.ActivationRequestData{ChainID: "chain-123"},
	})
	if err == nil {
		t.Error("Expected validation error for empty path ID")
	}

	// Test empty chain ID
	_, err = client.Perception.CreateActivation("path-123", perception.CreateActivationRequest{
		Activation: perception.ActivationRequestData{ChainID: ""},
	})
	if err == nil {
		t.Error("Expected validation error for empty chain ID")
	}
}

func TestPerceptionCreateActivationChainIDValidationDelegated(t *testing.T) {
	// Test that API response with chain_id validation errors is handled correctly
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResp := perception.Error{
			StatusCode: 422,
			Errors: map[string][]string{
				"chain_id": {"is required", "must be a valid chain"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Perception.CreateActivation("path-123", perception.CreateActivationRequest{
		Activation: perception.ActivationRequestData{ChainID: "invalid"},
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var perceptionErr *perception.Error
	if errors.As(err, &perceptionErr) {
		if perceptionErr.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code 422, got %d", perceptionErr.StatusCode)
		}
		if len(perceptionErr.Errors["chain_id"]) != 2 {
			t.Errorf("Expected 2 chain_id errors, got %d", len(perceptionErr.Errors["chain_id"]))
		}
	} else {
		t.Errorf("Expected perception.Error, got %T", err)
	}
}

func TestPerceptionUpdateActivation(t *testing.T) {
	request := perception.UpdateActivationRequest{
		Activation: perception.UpdateActivationData{
			ChainID: "updated-chain-123",
		},
	}

	expectedActivation := perception.Activation{
		ID:             "activation-123",
		ThoughtPathID:  "path-123",
		ChainID:        "updated-chain-123",
		ProvisionState: "active",
	}

	expectedResponse := perception.ActivationResponse{
		Data: expectedActivation,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/activations/activation-123" {
			t.Errorf("Expected path /provision/perception/activations/activation-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	activation, err := client.Perception.UpdateActivation("activation-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if activation.ChainID != expectedActivation.ChainID {
		t.Errorf("Expected activation chain_id %s, got %s", expectedActivation.ChainID, activation.ChainID)
	}
}

func TestPerceptionReplaceActivation(t *testing.T) {
	request := perception.UpdateActivationRequest{
		Activation: perception.UpdateActivationData{
			ChainID: "replaced-chain-123",
		},
	}

	expectedActivation := perception.Activation{
		ID:             "activation-123",
		ThoughtPathID:  "path-123",
		ChainID:        "replaced-chain-123",
		ProvisionState: "active",
	}

	expectedResponse := perception.ActivationResponse{
		Data: expectedActivation,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/activations/activation-123" {
			t.Errorf("Expected path /provision/perception/activations/activation-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	activation, err := client.Perception.ReplaceActivation("activation-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if activation.ChainID != expectedActivation.ChainID {
		t.Errorf("Expected activation chain_id %s, got %s", expectedActivation.ChainID, activation.ChainID)
	}
}

func TestPerceptionDeleteActivation(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/activations/activation-123" {
			t.Errorf("Expected path /provision/perception/activations/activation-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	err = client.Perception.DeleteActivation("activation-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetActivationEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.GetActivation("")
	if err == nil {
		t.Error("Expected validation error for empty activation ID in GetActivation")
	}
}

func TestPerceptionUpdateActivationEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.UpdateActivation("", perception.UpdateActivationRequest{
		Activation: perception.UpdateActivationData{ChainID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty activation ID in UpdateActivation")
	}
}

func TestPerceptionReplaceActivationEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.ReplaceActivation("", perception.UpdateActivationRequest{
		Activation: perception.UpdateActivationData{ChainID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty activation ID in ReplaceActivation")
	}
}

func TestPerceptionDeleteActivationEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Perception.DeleteActivation("")
	if err == nil {
		t.Error("Expected validation error for empty activation ID in DeleteActivation")
	}
}

func TestPerceptionCreateActivationWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/paths/path-123/activations" {
			t.Errorf("Expected path /provision/perception/paths/path-123/activations, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		errorResponse := map[string]any{
			"errors": map[string][]string{
				"chain_id": {"can't be blank", "must be a valid chain ID"},
				"base":     {"path is not active"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	request := perception.CreateActivationRequest{
		Activation: perception.ActivationRequestData{
			ChainID: "invalid", // Invalid chain ID to trigger server-side validation
		},
	}

	_, err = client.Perception.CreateActivation("path-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "chain_id can't be blank") {
		t.Errorf("Expected error to contain 'chain_id can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "chain_id must be a valid chain ID") {
		t.Errorf("Expected error to contain 'chain_id must be a valid chain ID', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "base path is not active") {
		t.Errorf("Expected error to contain 'base path is not active', got %s", errorMsg)
	}
}
