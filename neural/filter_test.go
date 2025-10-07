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

func TestNeuralGetFilter(t *testing.T) {
	expected := neural.Filter{
		ID:             "filter-123",
		ListenerID:     "listener-123",
		ChainID:        "chain-456",
		ProvisionState: "active",
	}
	expectedResp := neural.FilterResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/filters/filter-123" {
			t.Errorf("Expected path /provision/neural/filters/filter-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	filter, err := client.Neural.GetFilter("filter-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if filter.ID != expected.ID {
		t.Errorf("Expected filter ID %s, got %s", expected.ID, filter.ID)
	}
	if filter.ListenerID != expected.ListenerID {
		t.Errorf("Expected listener ID %s, got %s", expected.ListenerID, filter.ListenerID)
	}
	if filter.ChainID != expected.ChainID {
		t.Errorf("Expected chain ID %s, got %s", expected.ChainID, filter.ChainID)
	}
	if filter.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expected.ProvisionState, filter.ProvisionState)
	}
}

func TestNeuralGetFilterError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"filter": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Neural.GetFilter("missing")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	var nerr *neural.Error
	if errors.As(err, &nerr) {
		if nerr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", nerr.StatusCode)
		}
		if nerr.Errors == nil || len(nerr.Errors["filter"]) == 0 || nerr.Errors["filter"][0] != "not found" {
			t.Errorf("Expected 'filter not found', got %v", nerr.Errors)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralCreateFilter(t *testing.T) {
	expected := neural.Filter{
		ID:             "filter-789",
		ListenerID:     "listener-123",
		ChainID:        "chain-456",
		ProvisionState: "active",
	}
	expectedResp := neural.FilterResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/listeners/listener-123/filters" {
			t.Errorf("Expected path /provision/neural/listeners/listener-123/filters, got %s", r.URL.Path)
		}
		var req neural.CreateFilterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Filter.ChainID != "chain-456" {
			t.Errorf("Expected chain_id 'chain-456', got %s", req.Filter.ChainID)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	createReq := neural.CreateFilterRequest{Filter: neural.FilterRequestData{ChainID: "chain-456"}}
	filter, err := client.Neural.CreateFilter("listener-123", createReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if filter.ID != expected.ID {
		t.Errorf("Expected filter ID %s, got %s", expected.ID, filter.ID)
	}
	if filter.ChainID != expected.ChainID {
		t.Errorf("Expected chain ID %s, got %s", expected.ChainID, filter.ChainID)
	}
}

func TestNeuralCreateFilterValidation(t *testing.T) {
client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Empty listener ID
_, err = client.Neural.CreateFilter(
		"",
		neural.CreateFilterRequest{
			Filter: neural.FilterRequestData{ChainID: "chain-1"},
		},
	)
	if err == nil {
		t.Error("Expected validation error for empty listener ID")
	}
	// Empty chain ID
	_, err = client.Neural.CreateFilter("listener-1", neural.CreateFilterRequest{Filter: neural.FilterRequestData{}})
	if err == nil {
		t.Error("Expected validation error for empty chain ID")
	}
}

func TestNeuralUpdateFilter(t *testing.T) {
	expected := neural.Filter{
		ID:             "filter-123",
		ListenerID:     "listener-123",
		ChainID:        "chain-999",
		ProvisionState: "active",
	}
	expectedResp := neural.FilterResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/filters/filter-123" {
			t.Errorf("Expected path /provision/neural/filters/filter-123, got %s", r.URL.Path)
		}
		var req neural.UpdateFilterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Filter.ChainID != "chain-999" {
			t.Errorf("Expected chain_id 'chain-999', got %s", req.Filter.ChainID)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	updateReq := neural.UpdateFilterRequest{Filter: neural.UpdateFilterData{ChainID: "chain-999"}}
	filter, err := client.Neural.UpdateFilter("filter-123", updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if filter.ChainID != expected.ChainID {
		t.Errorf("Expected chain ID %s, got %s", expected.ChainID, filter.ChainID)
	}
}

func TestNeuralUpdateFilterValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Empty filter ID
	_, err = client.Neural.UpdateFilter("", neural.UpdateFilterRequest{Filter: neural.UpdateFilterData{ChainID: "x"}})
	if err == nil {
		t.Error("Expected validation error for empty filter ID")
	}
	// Empty chain ID
	_, err = client.Neural.UpdateFilter("filter-1", neural.UpdateFilterRequest{Filter: neural.UpdateFilterData{}})
	if err == nil {
		t.Error("Expected validation error for empty chain ID")
	}
}

func TestNeuralReplaceFilter(t *testing.T) {
	expected := neural.Filter{
		ID:             "filter-123",
		ListenerID:     "listener-123",
		ChainID:        "chain-abc",
		ProvisionState: "active",
	}
	expectedResp := neural.FilterResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/filters/filter-123" {
			t.Errorf("Expected path /provision/neural/filters/filter-123, got %s", r.URL.Path)
		}
		var req neural.UpdateFilterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Filter.ChainID != "chain-abc" {
			t.Errorf("Expected chain_id 'chain-abc', got %s", req.Filter.ChainID)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	replaceReq := neural.UpdateFilterRequest{Filter: neural.UpdateFilterData{ChainID: "chain-abc"}}
	filter, err := client.Neural.ReplaceFilter("filter-123", replaceReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if filter.ChainID != expected.ChainID {
		t.Errorf("Expected chain ID %s, got %s", expected.ChainID, filter.ChainID)
	}
}

func TestNeuralDeleteFilter(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/neural/filters/filter-123" {
			t.Errorf("Expected path /provision/neural/filters/filter-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Neural.DeleteFilter("filter-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNeuralDeleteFilterValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Neural.DeleteFilter("")
	if err == nil {
		t.Error("Expected validation error for empty filter ID")
	}
}
