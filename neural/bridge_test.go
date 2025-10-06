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

func TestNeuralGetBridge(t *testing.T) {
	expectedBridge := neural.Bridge{
		ID:             "bridge-123",
		SpaceID:        "space-123",
		TargetSpaceID:  "space-456",
		ProvisionState: "active",
	}

	expectedResponse := neural.BridgeResponse{
		Data: expectedBridge,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/bridges/bridge-123" {
			t.Errorf("Expected path /provision/neural/bridges/bridge-123, got %s", r.URL.Path)
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
	bridge, err := client.Neural.GetBridge("bridge-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bridge.ID != expectedBridge.ID {
		t.Errorf("Expected bridge ID %s, got %s", expectedBridge.ID, bridge.ID)
	}

	if bridge.SpaceID != expectedBridge.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expectedBridge.SpaceID, bridge.SpaceID)
	}

	if bridge.TargetSpaceID != expectedBridge.TargetSpaceID {
		t.Errorf("Expected target space ID %s, got %s", expectedBridge.TargetSpaceID, bridge.TargetSpaceID)
	}

	if bridge.ProvisionState != expectedBridge.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expectedBridge.ProvisionState, bridge.ProvisionState)
	}
}

func TestNeuralGetBridgeError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"bridge": {"not found"},
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
	_, err = client.Neural.GetBridge("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var neuralErr *neural.Error
	if errors.As(err, &neuralErr) {
		if neuralErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", neuralErr.StatusCode)
		}
		if neuralErr.Errors == nil || len(neuralErr.Errors["bridge"]) == 0 ||
			neuralErr.Errors["bridge"][0] != "not found" {
			t.Errorf("Expected error 'bridge not found', got %v", neuralErr.Errors)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralCreateBridge(t *testing.T) {
	expectedBridge := neural.Bridge{
		ID:             "bridge-789",
		SpaceID:        "space-123",
		TargetSpaceID:  "space-456",
		ProvisionState: "active",
	}

	expectedResponse := neural.BridgeResponse{
		Data: expectedBridge,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123/bridges" {
			t.Errorf("Expected path /provision/neural/spaces/space-123/bridges, got %s", r.URL.Path)
		}

		var req neural.CreateBridgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Bridge.TargetSpaceID != "space-456" {
			t.Errorf("Expected request target space ID 'space-456', got %s", req.Bridge.TargetSpaceID)
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

	createReq := neural.CreateBridgeRequest{
		Bridge: neural.BridgeRequestData{
			TargetSpaceID: "space-456",
		},
	}

	bridge, err := client.Neural.CreateBridge("space-123", createReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bridge.ID != expectedBridge.ID {
		t.Errorf("Expected bridge ID %s, got %s", expectedBridge.ID, bridge.ID)
	}

	if bridge.SpaceID != expectedBridge.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expectedBridge.SpaceID, bridge.SpaceID)
	}

	if bridge.TargetSpaceID != expectedBridge.TargetSpaceID {
		t.Errorf("Expected target space ID %s, got %s", expectedBridge.TargetSpaceID, bridge.TargetSpaceID)
	}
}

func TestNeuralCreateBridgeValidation(t *testing.T) {
	config := tama.Config{
		BaseURL:        "http://example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty space ID validation
	_, err = client.Neural.CreateBridge("", neural.CreateBridgeRequest{
		Bridge: neural.BridgeRequestData{
			TargetSpaceID: "space-456",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty target space ID validation
	_, err = client.Neural.CreateBridge("space-123", neural.CreateBridgeRequest{
		Bridge: neural.BridgeRequestData{},
	})

	if err == nil {
		t.Error("Expected validation error for empty target space ID")
	}
}

func TestNeuralUpdateBridge(t *testing.T) {
	expectedBridge := neural.Bridge{
		ID:             "bridge-123",
		SpaceID:        "space-123",
		TargetSpaceID:  "space-789",
		ProvisionState: "active",
	}

	expectedResponse := neural.BridgeResponse{
		Data: expectedBridge,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/bridges/bridge-123" {
			t.Errorf("Expected path /provision/neural/bridges/bridge-123, got %s", r.URL.Path)
		}

		var req neural.UpdateBridgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Bridge.TargetSpaceID != "space-789" {
			t.Errorf("Expected request target space ID 'space-789', got %s", req.Bridge.TargetSpaceID)
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

	updateReq := neural.UpdateBridgeRequest{
		Bridge: neural.UpdateBridgeData{
			TargetSpaceID: "space-789",
		},
	}

	bridge, err := client.Neural.UpdateBridge("bridge-123", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bridge.TargetSpaceID != expectedBridge.TargetSpaceID {
		t.Errorf("Expected target space ID %s, got %s", expectedBridge.TargetSpaceID, bridge.TargetSpaceID)
	}
}

func TestNeuralUpdateBridgeValidation(t *testing.T) {
	config := tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty bridge ID validation
	_, err = client.Neural.UpdateBridge("", neural.UpdateBridgeRequest{
		Bridge: neural.UpdateBridgeData{
			TargetSpaceID: "space-456",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty bridge ID")
	}

	// Test empty target space ID validation
	_, err = client.Neural.UpdateBridge("bridge-123", neural.UpdateBridgeRequest{
		Bridge: neural.UpdateBridgeData{},
	})

	if err == nil {
		t.Error("Expected validation error for empty target space ID")
	}
}

func TestNeuralReplaceBridge(t *testing.T) {
	expectedBridge := neural.Bridge{
		ID:             "bridge-123",
		SpaceID:        "space-123",
		TargetSpaceID:  "space-999",
		ProvisionState: "active",
	}

	expectedResponse := neural.BridgeResponse{
		Data: expectedBridge,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/bridges/bridge-123" {
			t.Errorf("Expected path /provision/neural/bridges/bridge-123, got %s", r.URL.Path)
		}

		var req neural.UpdateBridgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Bridge.TargetSpaceID != "space-999" {
			t.Errorf("Expected request target space ID 'space-999', got %s", req.Bridge.TargetSpaceID)
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

	replaceReq := neural.UpdateBridgeRequest{
		Bridge: neural.UpdateBridgeData{
			TargetSpaceID: "space-999",
		},
	}

	bridge, err := client.Neural.ReplaceBridge("bridge-123", replaceReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bridge.TargetSpaceID != expectedBridge.TargetSpaceID {
		t.Errorf("Expected target space ID %s, got %s", expectedBridge.TargetSpaceID, bridge.TargetSpaceID)
	}
}

func TestNeuralReplaceBridgeValidation(t *testing.T) {
	config := tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty bridge ID validation
	_, err = client.Neural.ReplaceBridge("", neural.UpdateBridgeRequest{
		Bridge: neural.UpdateBridgeData{
			TargetSpaceID: "space-456",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty bridge ID")
	}

	// Test empty target space ID validation
	_, err = client.Neural.ReplaceBridge("bridge-123", neural.UpdateBridgeRequest{
		Bridge: neural.UpdateBridgeData{},
	})

	if err == nil {
		t.Error("Expected validation error for empty target space ID")
	}
}

func TestNeuralDeleteBridge(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/bridges/bridge-123" {
			t.Errorf("Expected path /provision/neural/bridges/bridge-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Neural.DeleteBridge("bridge-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNeuralDeleteBridgeValidation(t *testing.T) {
	config := tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty bridge ID validation
	err = client.Neural.DeleteBridge("")
	if err == nil {
		t.Error("Expected validation error for empty bridge ID")
	}
}
