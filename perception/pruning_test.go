package perception_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/perception"
)

func TestPerceptionGetPruning(t *testing.T) {
	expectedPruning := perception.Pruning{
		ID:                    "pruning-123",
		ThoughtID:             "thought-456",
		PreviousVersionsCount: 2,
	}

	expectedResponse := perception.PruningResponse{
		Data: expectedPruning,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/prunings/pruning-123" {
			t.Errorf("Expected path /provision/perception/prunings/pruning-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	pruning, err := client.Perception.GetPruning("pruning-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidatePruningResponse(t, *pruning, expectedPruning)
}

func TestPerceptionCreatePruning(t *testing.T) {
	request := perception.CreatePruningRequest{
		Pruning: perception.CreatePruningData{
			PreviousVersionsCount: 3,
		},
	}

	expectedPruning := perception.Pruning{
		ID:                    "pruning-456",
		ThoughtID:             "thought-123",
		PreviousVersionsCount: 3,
	}

	expectedResponse := perception.PruningResponse{
		Data: expectedPruning,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/prunings" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/prunings, got %s",
				r.URL.Path,
			)
		}

		var received perception.CreatePruningRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if received.Pruning.PreviousVersionsCount != request.Pruning.PreviousVersionsCount {
			t.Errorf(
				"Expected previous_versions_count %d, got %d",
				request.Pruning.PreviousVersionsCount,
				received.Pruning.PreviousVersionsCount,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	pruning, err := client.Perception.CreatePruning("thought-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidatePruningResponse(t, *pruning, expectedPruning)
}

func TestPerceptionUpdatePruning(t *testing.T) {
	count := 5
	request := perception.UpdatePruningRequest{
		Pruning: perception.UpdatePruningData{
			PreviousVersionsCount: &count,
		},
	}

	expectedPruning := perception.Pruning{
		ID:                    "pruning-789",
		ThoughtID:             "thought-456",
		PreviousVersionsCount: 5,
	}

	expectedResponse := perception.PruningResponse{
		Data: expectedPruning,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/prunings/pruning-789" {
			t.Errorf("Expected path /provision/perception/prunings/pruning-789, got %s", r.URL.Path)
		}

		var received perception.UpdatePruningRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if received.Pruning.PreviousVersionsCount == nil {
			t.Fatal("Expected previous_versions_count to be present in request")
		}

		if *received.Pruning.PreviousVersionsCount != *request.Pruning.PreviousVersionsCount {
			t.Errorf(
				"Expected previous_versions_count %d, got %d",
				*request.Pruning.PreviousVersionsCount,
				*received.Pruning.PreviousVersionsCount,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	pruning, err := client.Perception.UpdatePruning("pruning-789", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidatePruningResponse(t, *pruning, expectedPruning)
}

func TestPerceptionReplacePruning(t *testing.T) {
	count := 7
	request := perception.UpdatePruningRequest{
		Pruning: perception.UpdatePruningData{
			PreviousVersionsCount: &count,
		},
	}

	expectedPruning := perception.Pruning{
		ID:                    "pruning-999",
		ThoughtID:             "thought-456",
		PreviousVersionsCount: 7,
	}

	expectedResponse := perception.PruningResponse{
		Data: expectedPruning,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/prunings/pruning-999" {
			t.Errorf("Expected path /provision/perception/prunings/pruning-999, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	pruning, err := client.Perception.ReplacePruning("pruning-999", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidatePruningResponse(t, *pruning, expectedPruning)
}

func TestPerceptionDeletePruning(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/prunings/pruning-456" {
			t.Errorf("Expected path /provision/perception/prunings/pruning-456, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Perception.DeletePruning("pruning-456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionCreatePruningValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.CreatePruning("", perception.CreatePruningRequest{
		Pruning: perception.CreatePruningData{PreviousVersionsCount: 1},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought ID")
	}
}

func TestPerceptionGetPruningEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.GetPruning("")
	if err == nil {
		t.Error("Expected validation error for empty pruning ID in GetPruning")
	}
}

func TestPerceptionUpdatePruningEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	count := 1
	_, err = client.Perception.UpdatePruning("", perception.UpdatePruningRequest{
		Pruning: perception.UpdatePruningData{PreviousVersionsCount: &count},
	})
	if err == nil {
		t.Error("Expected validation error for empty pruning ID in UpdatePruning")
	}
}

func TestPerceptionReplacePruningEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	count := 1
	_, err = client.Perception.ReplacePruning("", perception.UpdatePruningRequest{
		Pruning: perception.UpdatePruningData{PreviousVersionsCount: &count},
	})
	if err == nil {
		t.Error("Expected validation error for empty pruning ID in ReplacePruning")
	}
}

func TestPerceptionDeletePruningEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Perception.DeletePruning("")
	if err == nil {
		t.Error("Expected validation error for empty pruning ID in DeletePruning")
	}
}
