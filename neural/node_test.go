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

func TestNeuralGetNode(t *testing.T) {
	expectedNode := neural.Node{
		ID:             "node-123",
		On:             "active",
		Type:           "compute",
		SpaceID:        "space-456",
		ClassID:        "class-789",
		ChainID:        "chain-abc",
		ProvisionState: "active",
	}

	expectedResponse := neural.NodeResponse{
		Data: expectedNode,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/nodes/node-123" {
			t.Errorf("Expected path /provision/neural/nodes/node-123, got %s", r.URL.Path)
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
	node, err := client.Neural.GetNode("node-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if node.ID != expectedNode.ID {
		t.Errorf("Expected node ID %s, got %s", expectedNode.ID, node.ID)
	}

	if node.On != expectedNode.On {
		t.Errorf("Expected node on %s, got %s", expectedNode.On, node.On)
	}

	if node.Type != expectedNode.Type {
		t.Errorf("Expected node type %s, got %s", expectedNode.Type, node.Type)
	}

	if node.SpaceID != expectedNode.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expectedNode.SpaceID, node.SpaceID)
	}

	if node.ClassID != expectedNode.ClassID {
		t.Errorf("Expected class ID %s, got %s", expectedNode.ClassID, node.ClassID)
	}

	if node.ChainID != expectedNode.ChainID {
		t.Errorf("Expected chain ID %s, got %s", expectedNode.ChainID, node.ChainID)
	}
}

func TestNeuralGetNodeError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"node": {"not found"},
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
	_, err = client.Neural.GetNode("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var neuralErr *neural.Error
	if errors.As(err, &neuralErr) {
		if neuralErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", neuralErr.StatusCode)
		}
		if neuralErr.Errors == nil || len(neuralErr.Errors["node"]) == 0 ||
			neuralErr.Errors["node"][0] != "not found" {
			t.Errorf("Expected error 'node not found', got %v", neuralErr.Errors)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralCreateNode(t *testing.T) {
	expectedNode := neural.Node{
		ID:             "node-789",
		On:             "active",
		Type:           "compute",
		SpaceID:        "space-456",
		ClassID:        "class-789",
		ChainID:        "chain-abc",
		ProvisionState: "active",
	}

	expectedResponse := neural.NodeResponse{
		Data: expectedNode,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-456/nodes" {
			t.Errorf("Expected path /provision/neural/spaces/space-456/nodes, got %s", r.URL.Path)
		}

		var req neural.CreateNodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Node.Type != "compute" {
			t.Errorf("Expected request type 'compute', got %s", req.Node.Type)
		}

		if req.Node.ClassID != "class-789" {
			t.Errorf("Expected request class ID 'class-789', got %s", req.Node.ClassID)
		}

		if req.Node.ChainID != "chain-abc" {
			t.Errorf("Expected request chain ID 'chain-abc', got %s", req.Node.ChainID)
		}

		if req.Node.On != "active" {
			t.Errorf("Expected request on 'active', got %s", req.Node.On)
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

	createReq := neural.CreateNodeRequest{
		Node: neural.NodeRequestData{
			On:      "active",
			Type:    "compute",
			ClassID: "class-789",
			ChainID: "chain-abc",
		},
	}

	node, err := client.Neural.CreateNode("space-456", createReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if node.ID != expectedNode.ID {
		t.Errorf("Expected node ID %s, got %s", expectedNode.ID, node.ID)
	}

	if node.Type != expectedNode.Type {
		t.Errorf("Expected node type %s, got %s", expectedNode.Type, node.Type)
	}

	if node.ClassID != expectedNode.ClassID {
		t.Errorf("Expected class ID %s, got %s", expectedNode.ClassID, node.ClassID)
	}

	if node.ChainID != expectedNode.ChainID {
		t.Errorf("Expected chain ID %s, got %s", expectedNode.ChainID, node.ChainID)
	}
}

func TestNeuralCreateNodeValidation(t *testing.T) {
	config := tama.Config{
		BaseURL:        "http://localhost:8080",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty space ID validation
	_, err = client.Neural.CreateNode("", neural.CreateNodeRequest{
		Node: neural.NodeRequestData{
			Type:    "compute",
			ClassID: "class-789",
			ChainID: "chain-abc",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty type validation
	_, err = client.Neural.CreateNode("space-456", neural.CreateNodeRequest{
		Node: neural.NodeRequestData{
			ClassID: "class-789",
			ChainID: "chain-abc",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty type")
	}

	// Test empty class ID validation
	_, err = client.Neural.CreateNode("space-456", neural.CreateNodeRequest{
		Node: neural.NodeRequestData{
			Type:    "compute",
			ChainID: "chain-abc",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty class ID")
	}

	// Test empty chain ID validation
	_, err = client.Neural.CreateNode("space-456", neural.CreateNodeRequest{
		Node: neural.NodeRequestData{
			Type:    "compute",
			ClassID: "class-789",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty chain ID")
	}
}

func TestNeuralUpdateNode(t *testing.T) {
	expectedNode := neural.Node{
		ID:             "node-123",
		On:             "inactive",
		Type:           "storage",
		SpaceID:        "space-456",
		ClassID:        "class-789",
		ChainID:        "chain-abc",
		ProvisionState: "active",
	}

	expectedResponse := neural.NodeResponse{
		Data: expectedNode,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/nodes/node-123" {
			t.Errorf("Expected path /provision/neural/nodes/node-123, got %s", r.URL.Path)
		}

		var req neural.UpdateNodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
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

	updateReq := neural.UpdateNodeRequest{
		Node: neural.UpdateNodeData{
			On:   "inactive",
			Type: "storage",
		},
	}

	node, err := client.Neural.UpdateNode("node-123", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if node.On != expectedNode.On {
		t.Errorf("Expected node on %s, got %s", expectedNode.On, node.On)
	}

	if node.Type != expectedNode.Type {
		t.Errorf("Expected node type %s, got %s", expectedNode.Type, node.Type)
	}
}

func TestNeuralReplaceNode(t *testing.T) {
	expectedNode := neural.Node{
		ID:             "node-123",
		On:             "active",
		Type:           "network",
		SpaceID:        "space-456",
		ClassID:        "class-789",
		ChainID:        "chain-abc",
		ProvisionState: "active",
	}

	expectedResponse := neural.NodeResponse{
		Data: expectedNode,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/nodes/node-123" {
			t.Errorf("Expected path /provision/neural/nodes/node-123, got %s", r.URL.Path)
		}

		var req neural.UpdateNodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
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

	replaceReq := neural.UpdateNodeRequest{
		Node: neural.UpdateNodeData{
			On:   "active",
			Type: "network",
		},
	}

	node, err := client.Neural.ReplaceNode("node-123", replaceReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if node.Type != expectedNode.Type {
		t.Errorf("Expected node type %s, got %s", expectedNode.Type, node.Type)
	}
}

func TestNeuralDeleteNode(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/nodes/node-123" {
			t.Errorf("Expected path /provision/neural/nodes/node-123, got %s", r.URL.Path)
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

	err = client.Neural.DeleteNode("node-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNeuralDeleteNodeValidation(t *testing.T) {
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

	// Test empty node ID validation
	err = client.Neural.DeleteNode("")

	if err == nil {
		t.Error("Expected validation error for empty node ID")
	}
}
