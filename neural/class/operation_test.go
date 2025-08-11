package class_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural/class"
)

// createMockServer creates a test HTTP server with the given handler.
func createMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// ValidateOperationResponse validates that actual operation matches expected operation.
func ValidateOperationResponse(t *testing.T, actual, expected class.Operation) {
	if actual.ID != expected.ID {
		t.Errorf("Expected operation ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.CurrentState != expected.CurrentState {
		t.Errorf("Expected operation current_state %s, got %s", expected.CurrentState, actual.CurrentState)
	}

	if actual.ClassID != expected.ClassID {
		t.Errorf("Expected operation class_id %s, got %s", expected.ClassID, actual.ClassID)
	}

	if len(actual.NodeIDs) != len(expected.NodeIDs) {
		t.Errorf("Expected %d node IDs, got %d", len(expected.NodeIDs), len(actual.NodeIDs))
		return
	}

	for i, nodeID := range expected.NodeIDs {
		if actual.NodeIDs[i] != nodeID {
			t.Errorf("Expected node ID %s at index %d, got %s", nodeID, i, actual.NodeIDs[i])
		}
	}
}

func TestClassGetOperation(t *testing.T) {
	expectedOperation := class.Operation{
		ID:           "operation-123",
		CurrentState: "completed",
		ClassID:      "class-123",
		NodeIDs:      []string{"node-1", "node-2"},
	}

	expectedResponse := class.OperationResponse{
		Data: expectedOperation,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123/operations/operation-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123/operations/operation-123, got %s", r.URL.Path)
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
	service := class.NewService(client.GetHTTPClient())
	operation, err := service.GetOperation("class-123", "operation-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateOperationResponse(t, *operation, expectedOperation)
}

func TestClassGetOperationError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{
			"errors": map[string][]string{
				"id": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	service := class.NewService(client.GetHTTPClient())
	operation, err := service.GetOperation("class-123", "invalid-id")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if operation != nil {
		t.Error("Expected nil operation on error")
	}
}

func TestClassCreateOperation(t *testing.T) {
	expectedOperation := class.Operation{
		ID:           "operation-123",
		CurrentState: "pending",
		ClassID:      "class-123",
		NodeIDs:      []string{},
	}

	expectedResponse := class.OperationResponse{
		Data: expectedOperation,
	}

	var receivedRequest *class.CreateOperationRequest

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123/operations" {
			t.Errorf("Expected path /provision/neural/classes/class-123/operations, got %s", r.URL.Path)
		}

		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
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
	service := class.NewService(client.GetHTTPClient())

	nodeType := "worker"
	req := class.CreateOperationRequest{
		Operation: class.CreateOperationData{
			ChainIDs: []string{"chain-1", "chain-2"},
			NodeType: &nodeType,
		},
	}

	operation, err := service.CreateOperation("class-123", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateOperationResponse(t, *operation, expectedOperation)

	// Validate request payload
	if len(receivedRequest.Operation.ChainIDs) != 2 {
		t.Errorf("Expected 2 chain IDs, got %d", len(receivedRequest.Operation.ChainIDs))
	}

	if receivedRequest.Operation.ChainIDs[0] != "chain-1" {
		t.Errorf("Expected first chain ID 'chain-1', got %s", receivedRequest.Operation.ChainIDs[0])
	}

	if receivedRequest.Operation.ChainIDs[1] != "chain-2" {
		t.Errorf("Expected second chain ID 'chain-2', got %s", receivedRequest.Operation.ChainIDs[1])
	}

	if receivedRequest.Operation.NodeType == nil || *receivedRequest.Operation.NodeType != "worker" {
		t.Errorf("Expected node type 'worker', got %v", receivedRequest.Operation.NodeType)
	}
}

func TestClassCreateOperationValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "http://example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	service := class.NewService(client.GetHTTPClient())

	// Test empty class ID
	req := class.CreateOperationRequest{
		Operation: class.CreateOperationData{
			ChainIDs: []string{"chain-1"},
		},
	}

	_, err := service.CreateOperation("", req)
	if err == nil || err.Error() != "class ID is required" {
		t.Errorf("Expected 'class ID is required' error, got %v", err)
	}

	// Test empty chain IDs
	req = class.CreateOperationRequest{
		Operation: class.CreateOperationData{
			ChainIDs: []string{},
		},
	}

	_, err = service.CreateOperation("class-123", req)
	if err == nil || err.Error() != "chain IDs are required" {
		t.Errorf("Expected 'chain IDs are required' error, got %v", err)
	}
}

func TestClassGetOperationEmptyClassID(t *testing.T) {
	config := tama.Config{
		BaseURL: "http://example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	service := class.NewService(client.GetHTTPClient())
	_, err := service.GetOperation("", "operation-123")

	if err == nil || err.Error() != "class ID is required" {
		t.Errorf("Expected 'class ID is required' error, got %v", err)
	}
}

func TestClassGetOperationEmptyOperationID(t *testing.T) {
	config := tama.Config{
		BaseURL: "http://example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	service := class.NewService(client.GetHTTPClient())
	_, err := service.GetOperation("class-123", "")

	if err == nil || err.Error() != "operation ID is required" {
		t.Errorf("Expected 'operation ID is required' error, got %v", err)
	}
}

func TestClassCreateOperationWithFieldErrors(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{
			"errors": map[string][]string{
				"chain_ids": {"is required"},
				"node_type": {"is invalid"},
			},
		}
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	service := class.NewService(client.GetHTTPClient())

	nodeType := "invalid-type"
	req := class.CreateOperationRequest{
		Operation: class.CreateOperationData{
			ChainIDs: []string{"chain-1"},
			NodeType: &nodeType,
		},
	}

	operation, err := service.CreateOperation("class-123", req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if operation != nil {
		t.Error("Expected nil operation on error")
	}

	// Check that error message contains expected validation errors
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "chain_ids is required") {
		t.Errorf("Expected error to contain 'chain_ids is required', got: %s", errorMsg)
	}

	if !strings.Contains(errorMsg, "node_type is invalid") {
		t.Errorf("Expected error to contain 'node_type is invalid', got: %s", errorMsg)
	}
}
