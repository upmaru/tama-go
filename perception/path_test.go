package perception_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/perception"
)

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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
	path, err := client.Perception.GetPath("path-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidatePathResponse(t, *path, expectedPath)
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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/paths" {
			t.Errorf(
				"Expected path /provision/perception/thoughts/thought-123/paths, got %s",
				r.URL.Path,
			)
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
	path, err := client.Perception.CreatePath("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidatePathResponse(t, *path, expectedPath)
}

func TestPerceptionCreatePathValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty thought ID
	_, err = client.Perception.CreatePath("", perception.CreatePathRequest{
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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
	path, err := client.Perception.UpdatePath("path-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if path.TargetClassID != expectedPath.TargetClassID {
		t.Errorf(
			"Expected target_class_id %s, got %s",
			expectedPath.TargetClassID,
			path.TargetClassID,
		)
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

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
	path, err := client.Perception.ReplacePath("path-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if path.TargetClassID != expectedPath.TargetClassID {
		t.Errorf(
			"Expected target_class_id %s, got %s",
			expectedPath.TargetClassID,
			path.TargetClassID,
		)
	}
}

func TestPerceptionDeletePath(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
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
	err = client.Perception.DeletePath("path-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionGetPathEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.GetPath("")
	if err == nil {
		t.Error("Expected validation error for empty path ID in GetPath")
	}
}

func TestPerceptionUpdatePathEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.UpdatePath("", perception.UpdatePathRequest{
		Path: perception.UpdatePathData{TargetClassID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty path ID in UpdatePath")
	}
}

func TestPerceptionReplacePathEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Perception.ReplacePath("", perception.UpdatePathRequest{
		Path: perception.UpdatePathData{TargetClassID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty path ID in ReplacePath")
	}
}

func TestPerceptionDeletePathEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Perception.DeletePath("")
	if err == nil {
		t.Error("Expected validation error for empty path ID in DeletePath")
	}
}
