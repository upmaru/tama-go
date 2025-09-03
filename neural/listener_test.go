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

func TestNeuralGetListener(t *testing.T) {
	expected := neural.Listener{
		ID:             "listener-123",
		SpaceID:        "space-123",
		Endpoint:       "https://example.com/hook",
		ProvisionState: "active",
	}

	expectedResp := neural.ListenerResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/listeners/listener-123" {
			t.Errorf("Expected path /provision/neural/listeners/listener-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	config := tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second}
	client := tama.NewClient(config)

	listener, err := client.Neural.GetListener("listener-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if listener.ID != expected.ID {
		t.Errorf("Expected listener ID %s, got %s", expected.ID, listener.ID)
	}
	if listener.SpaceID != expected.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expected.SpaceID, listener.SpaceID)
	}
	if listener.Endpoint != expected.Endpoint {
		t.Errorf("Expected endpoint %s, got %s", expected.Endpoint, listener.Endpoint)
	}
	if listener.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expected.ProvisionState, listener.ProvisionState)
	}
}


func TestNeuralGetListenerError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"listener": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second}
	client := tama.NewClient(config)

	_, err := client.Neural.GetListener("missing")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	var nerr *neural.Error
	if errors.As(err, &nerr) {
		if nerr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", nerr.StatusCode)
		}
		if nerr.Errors == nil || len(nerr.Errors["listener"]) == 0 || nerr.Errors["listener"][0] != "not found" {
			t.Errorf("Expected 'listener not found', got %v", nerr.Errors)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralCreateListener(t *testing.T) {
	expected := neural.Listener{
		ID:             "listener-789",
		SpaceID:        "space-123",
		Endpoint:       "https://example.com/hook",
		ProvisionState: "active",
	}
	expectedResp := neural.ListenerResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/spaces/space-123/listeners" {
			t.Errorf("Expected path /provision/neural/spaces/space-123/listeners, got %s", r.URL.Path)
		}
		var req neural.CreateListenerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Listener.Endpoint != "https://example.com/hook" {
			t.Errorf("Expected endpoint 'https://example.com/hook', got %s", req.Listener.Endpoint)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})

	createReq := neural.CreateListenerRequest{Listener: neural.ListenerRequestData{Endpoint: "https://example.com/hook"}}
	listener, err := client.Neural.CreateListener("space-123", createReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if listener.ID != expected.ID {
		t.Errorf("Expected listener ID %s, got %s", expected.ID, listener.ID)
	}
	if listener.Endpoint != expected.Endpoint {
		t.Errorf("Expected endpoint %s, got %s", expected.Endpoint, listener.Endpoint)
	}
}

func TestNeuralCreateListenerValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", APIKey: "test-key", Timeout: 10 * time.Second})

	// Empty space ID
	_, err := client.Neural.CreateListener("", neural.CreateListenerRequest{Listener: neural.ListenerRequestData{Endpoint: "https://example.com/hook"}})
	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Empty endpoint
	_, err = client.Neural.CreateListener("space-123", neural.CreateListenerRequest{Listener: neural.ListenerRequestData{}})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestNeuralUpdateListener(t *testing.T) {
	expected := neural.Listener{
		ID:             "listener-123",
		SpaceID:        "space-123",
		Endpoint:       "https://example.com/new",
		ProvisionState: "active",
	}
	expectedResp := neural.ListenerResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/listeners/listener-123" {
			t.Errorf("Expected path /provision/neural/listeners/listener-123, got %s", r.URL.Path)
		}
		var req neural.UpdateListenerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Listener.Endpoint != "https://example.com/new" {
			t.Errorf("Expected new endpoint, got %s", req.Listener.Endpoint)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})

	updateReq := neural.UpdateListenerRequest{Listener: neural.UpdateListenerData{Endpoint: "https://example.com/new"}}
	listener, err := client.Neural.UpdateListener("listener-123", updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if listener.Endpoint != expected.Endpoint {
		t.Errorf("Expected endpoint %s, got %s", expected.Endpoint, listener.Endpoint)
	}
}

func TestNeuralUpdateListenerValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", APIKey: "test-key", Timeout: 10 * time.Second})

	// Empty listener ID
	_, err := client.Neural.UpdateListener("", neural.UpdateListenerRequest{Listener: neural.UpdateListenerData{Endpoint: "x"}})
	if err == nil {
		t.Error("Expected validation error for empty listener ID")
	}
	// Empty endpoint
	_, err = client.Neural.UpdateListener("listener-123", neural.UpdateListenerRequest{Listener: neural.UpdateListenerData{}})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestNeuralReplaceListener(t *testing.T) {
	expected := neural.Listener{
		ID:             "listener-123",
		SpaceID:        "space-123",
		Endpoint:       "https://example.com/replace",
		ProvisionState: "active",
	}
	expectedResp := neural.ListenerResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/listeners/listener-123" {
			t.Errorf("Expected path /provision/neural/listeners/listener-123, got %s", r.URL.Path)
		}
		var req neural.UpdateListenerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Listener.Endpoint != "https://example.com/replace" {
			t.Errorf("Expected replace endpoint, got %s", req.Listener.Endpoint)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})

	replaceReq := neural.UpdateListenerRequest{Listener: neural.UpdateListenerData{Endpoint: "https://example.com/replace"}}
	listener, err := client.Neural.ReplaceListener("listener-123", replaceReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if listener.Endpoint != expected.Endpoint {
		t.Errorf("Expected endpoint %s, got %s", expected.Endpoint, listener.Endpoint)
	}
}

func TestNeuralReplaceListenerValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", APIKey: "test-key", Timeout: 10 * time.Second})

	// Empty listener ID
	_, err := client.Neural.ReplaceListener("", neural.UpdateListenerRequest{Listener: neural.UpdateListenerData{Endpoint: "x"}})
	if err == nil {
		t.Error("Expected validation error for empty listener ID")
	}
	// Empty endpoint
	_, err = client.Neural.ReplaceListener("listener-123", neural.UpdateListenerRequest{Listener: neural.UpdateListenerData{}})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestNeuralDeleteListener(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/listeners/listener-123" {
			t.Errorf("Expected path /provision/neural/listeners/listener-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})

	if err := client.Neural.DeleteListener("listener-123"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNeuralDeleteListenerValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", APIKey: "test-key", Timeout: 10 * time.Second})
	if err := client.Neural.DeleteListener(""); err == nil {
		t.Error("Expected validation error for empty listener ID")
	}
}

