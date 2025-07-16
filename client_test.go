package tama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/upmaru/tama-go/neural"
	"github.com/upmaru/tama-go/sensory"
)

func TestNewClient(t *testing.T) {
	config := Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.baseURL != config.BaseURL {
		t.Errorf("Expected baseURL %s, got %s", config.BaseURL, client.baseURL)
	}

	if client.apiKey != config.APIKey {
		t.Errorf("Expected apiKey %s, got %s", config.APIKey, client.apiKey)
	}

	if client.Neural == nil {
		t.Error("Expected Neural service to be initialized")
	}

	if client.Sensory == nil {
		t.Error("Expected Sensory service to be initialized")
	}
}

func TestNewClientDefaultTimeout(t *testing.T) {
	config := Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		// No timeout specified
	}

	client := NewClient(config)
	// We can't directly test the timeout, but we can ensure the client was created
	if client == nil {
		t.Fatal("Expected client to be created with default timeout")
	}
}

func TestSetAPIKey(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "initial-key",
	})

	newAPIKey := "new-api-key"
	client.SetAPIKey(newAPIKey)

	if client.apiKey != newAPIKey {
		t.Errorf("Expected apiKey %s, got %s", newAPIKey, client.apiKey)
	}
}

func TestSetDebug(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// This should not panic
	client.SetDebug(true)
	client.SetDebug(false)
}

func TestErrorStruct(t *testing.T) {
	err := &Error{
		StatusCode: 404,
		Message:    "Not Found",
		Details:    "Resource does not exist",
	}

	expected := "API error 404: Not Found - Resource does not exist"
	if err.Error() != expected {
		t.Errorf("Expected error message %s, got %s", expected, err.Error())
	}

	// Test without details
	err.Details = ""
	expected = "API error 404: Not Found"
	if err.Error() != expected {
		t.Errorf("Expected error message %s, got %s", expected, err.Error())
	}
}

// Mock server helpers for testing HTTP interactions

func createMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestNeuralGetSpace(t *testing.T) {
	expectedSpace := neural.Space{
		ID:   "space-123",
		Name: "Test Space",
		Slug: "test-space",
	}

	expectedResponse := neural.SpaceResponse{
		Data: expectedSpace,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123" {
			t.Errorf("Expected path /provision/neural/spaces/space-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	space, err := client.Neural.GetSpace("space-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if space.ID != expectedSpace.ID {
		t.Errorf("Expected space ID %s, got %s", expectedSpace.ID, space.ID)
	}

	if space.Name != expectedSpace.Name {
		t.Errorf("Expected space name %s, got %s", expectedSpace.Name, space.Name)
	}
}

func TestNeuralGetSpaceError(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(neural.Error{
			StatusCode: 404,
			Message:    "Space not found",
		})
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	_, err := client.Neural.GetSpace("nonexistent")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if apiErr, ok := err.(*neural.Error); ok {
		if apiErr.StatusCode != 404 {
			t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
		}
		if apiErr.Message != "Space not found" {
			t.Errorf("Expected message 'Space not found', got %s", apiErr.Message)
		}
	} else {
		t.Errorf("Expected *neural.Error type, got %T", err)
	}
}



func TestNeuralCreateSpace(t *testing.T) {
	expectedSpace := neural.Space{
		ID:   "space-456",
		Name: "New Space",
		Slug: "new-space",
	}

	expectedResponse := neural.SpaceResponse{
		Data: expectedSpace,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces" {
			t.Errorf("Expected path /provision/neural/spaces, got %s", r.URL.Path)
		}

		var req neural.CreateSpaceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Space.Name != "New Space" {
			t.Errorf("Expected request name 'New Space', got %s", req.Space.Name)
		}

		if req.Space.Type != "root" {
			t.Errorf("Expected request type 'root', got %s", req.Space.Type)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: "New Space",
			Type: "root",
		},
	}

	space, err := client.Neural.CreateSpace(createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if space.ID != expectedSpace.ID {
		t.Errorf("Expected space ID %s, got %s", expectedSpace.ID, space.ID)
	}
}

func TestNeuralCreateSpaceValidation(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty name validation
	_, err := client.Neural.CreateSpace(neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{Type: "root"},
	})
	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	// Test empty type validation
	_, err = client.Neural.CreateSpace(neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{Name: "Test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty type")
	}

	// Test invalid type validation
	_, err = client.Neural.CreateSpace(neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{Name: "Test", Type: "invalid"},
	})
	if err == nil {
		t.Error("Expected validation error for invalid type")
	}
}

func TestNeuralUpdateSpace(t *testing.T) {
	expectedSpace := neural.Space{
		ID:   "space-123",
		Name: "Updated Space",
		Slug: "updated-space",
	}

	expectedResponse := neural.SpaceResponse{
		Data: expectedSpace,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	updateReq := neural.UpdateSpaceRequest{
		Space: neural.UpdateSpaceData{
			Name: "Updated Space",
			Type: "component",
		},
	}

	space, err := client.Neural.UpdateSpace("space-123", updateReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if space.Name != expectedSpace.Name {
		t.Errorf("Expected space name %s, got %s", expectedSpace.Name, space.Name)
	}
}

func TestNeuralDeleteSpace(t *testing.T) {
	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-123" {
			t.Errorf("Expected path /provision/neural/spaces/space-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	err := client.Neural.DeleteSpace("space-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSensoryGetSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:   "source-123",
		Name: "Test Source",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123" {
			t.Errorf("Expected path /provision/sensory/sources/source-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	source, err := client.Sensory.GetSource("source-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}
}

func TestSensoryCreateSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:   "source-789",
		Name: "New Source",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/spaces/space-123/sources" {
			t.Errorf("Expected path /provision/sensory/spaces/space-123/sources, got %s", r.URL.Path)
		}

		var req sensory.CreateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Source.Name != "New Source" {
			t.Errorf("Expected request name 'New Source', got %s", req.Source.Name)
		}

		if req.Source.Type != "model" {
			t.Errorf("Expected request type 'model', got %s", req.Source.Type)
		}

		if req.Source.Endpoint != "https://api.mistral.ai/v1" {
			t.Errorf("Expected request endpoint 'https://api.mistral.ai/v1', got %s", req.Source.Endpoint)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "New Source",
			Type:     "model",
			Endpoint: "https://api.mistral.ai/v1",
			Credential: sensory.SourceCredential{
				ApiKey: "test-key",
			},
		},
	}

	source, err := client.Sensory.CreateSource("space-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}
}

func TestSensoryCreateSourceValidation(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty space ID validation
	_, err := client.Sensory.CreateSource("", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "Test",
			Type:     "model",
			Endpoint: "https://api.test.com",
			Credential: sensory.SourceCredential{
				ApiKey: "test-key",
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty name validation
	_, err = client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Type:     "model",
			Endpoint: "https://api.test.com",
			Credential: sensory.SourceCredential{
				ApiKey: "test-key",
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	// Test empty type validation
	_, err = client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "Test",
			Endpoint: "https://api.test.com",
			Credential: sensory.SourceCredential{
				ApiKey: "test-key",
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty type")
	}

	// Test empty endpoint validation
	_, err = client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name: "Test",
			Type: "model",
			Credential: sensory.SourceCredential{
				ApiKey: "test-key",
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestEmptyIDValidation(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test Neural service validations
	_, err := client.Neural.GetSpace("")
	if err == nil {
		t.Error("Expected validation error for empty space ID in GetSpace")
	}

	_, err = client.Neural.UpdateSpace("", neural.UpdateSpaceRequest{})
	if err == nil {
		t.Error("Expected validation error for empty space ID in UpdateSpace")
	}

	err = client.Neural.DeleteSpace("")
	if err == nil {
		t.Error("Expected validation error for empty space ID in DeleteSpace")
	}

	// Test Sensory service validations
	_, err = client.Sensory.GetSource("")
	if err == nil {
		t.Error("Expected validation error for empty source ID in GetSource")
	}

	_, err = client.Sensory.GetModel("")
	if err == nil {
		t.Error("Expected validation error for empty model ID in GetModel")
	}
}

func TestSensoryGetModel(t *testing.T) {
	expectedModel := sensory.Model{
		ID:         "model-123",
		Identifier: "mistral-small-latest",
	}

	expectedResponse := sensory.ModelResponse{
		Data: expectedModel,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/models/model-123" {
			t.Errorf("Expected path /provision/sensory/models/model-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	model, err := client.Sensory.GetModel("model-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if model.ID != expectedModel.ID {
		t.Errorf("Expected model ID %s, got %s", expectedModel.ID, model.ID)
	}

	if model.Identifier != expectedModel.Identifier {
		t.Errorf("Expected model identifier %s, got %s", expectedModel.Identifier, model.Identifier)
	}
}

func TestSensoryCreateModel(t *testing.T) {
	expectedModel := sensory.Model{
		ID:         "model-789",
		Identifier: "mistral-large-latest",
	}

	expectedResponse := sensory.ModelResponse{
		Data: expectedModel,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123/models" {
			t.Errorf("Expected path /provision/sensory/sources/source-123/models, got %s", r.URL.Path)
		}

		var req sensory.CreateModelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Model.Identifier != "mistral-large-latest" {
			t.Errorf("Expected request identifier 'mistral-large-latest', got %s", req.Model.Identifier)
		}

		if req.Model.Path != "/chat/completions" {
			t.Errorf("Expected request path '/chat/completions', got %s", req.Model.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := sensory.CreateModelRequest{
		Model: sensory.ModelRequestData{
			Identifier: "mistral-large-latest",
			Path:       "/chat/completions",
		},
	}

	model, err := client.Sensory.CreateModel("source-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if model.ID != expectedModel.ID {
		t.Errorf("Expected model ID %s, got %s", expectedModel.ID, model.ID)
	}
}

func TestSensoryCreateModelValidation(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty source ID validation
	_, err := client.Sensory.CreateModel("", sensory.CreateModelRequest{
		Model: sensory.ModelRequestData{
			Identifier: "test-model",
			Path:       "/chat/completions",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty source ID")
	}

	// Test empty identifier validation
	_, err = client.Sensory.CreateModel("source-123", sensory.CreateModelRequest{
		Model: sensory.ModelRequestData{
			Path: "/chat/completions",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty identifier")
	}

	// Test empty path validation
	_, err = client.Sensory.CreateModel("source-123", sensory.CreateModelRequest{
		Model: sensory.ModelRequestData{
			Identifier: "test-model",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty path")
	}
}

func TestEmptyIDValidationLimits(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.GetLimit("")
	if err == nil {
		t.Error("Expected validation error for empty limit ID in GetLimit")
	}
}

func TestSensoryGetLimit(t *testing.T) {
	expectedLimit := sensory.Limit{
		ID:           "limit-123",
		SourceID:     "source-456",
		Count:        32,
		ScaleUnit:    "seconds",
		ScaleCount:   1,
		CurrentState: "active",
	}

	expectedResponse := sensory.LimitResponse{
		Data: expectedLimit,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/limits/limit-123" {
			t.Errorf("Expected path /provision/sensory/limits/limit-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	limit, err := client.Sensory.GetLimit("limit-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if limit.ID != expectedLimit.ID {
		t.Errorf("Expected limit ID %s, got %s", expectedLimit.ID, limit.ID)
	}

	if limit.Count != expectedLimit.Count {
		t.Errorf("Expected count %d, got %d", expectedLimit.Count, limit.Count)
	}

	if limit.ScaleUnit != expectedLimit.ScaleUnit {
		t.Errorf("Expected scale unit %s, got %s", expectedLimit.ScaleUnit, limit.ScaleUnit)
	}

	if limit.ScaleCount != expectedLimit.ScaleCount {
		t.Errorf("Expected scale count %d, got %d", expectedLimit.ScaleCount, limit.ScaleCount)
	}
}

func TestSensoryCreateLimit(t *testing.T) {
	expectedLimit := sensory.Limit{
		ID:           "limit-789",
		SourceID:     "source-123",
		Count:        64,
		ScaleUnit:    "minutes",
		ScaleCount:   5,
		CurrentState: "active",
	}

	expectedResponse := sensory.LimitResponse{
		Data: expectedLimit,
	}

	server := createMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123/limits" {
			t.Errorf("Expected path /provision/sensory/sources/source-123/limits, got %s", r.URL.Path)
		}

		var req sensory.CreateLimitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Limit.Count != 64 {
			t.Errorf("Expected count 64, got %d", req.Limit.Count)
		}

		if req.Limit.ScaleUnit != "minutes" {
			t.Errorf("Expected request scale unit 'minutes', got %s", req.Limit.ScaleUnit)
		}

		if req.Limit.ScaleCount != 5 {
			t.Errorf("Expected request scale count 5, got %d", req.Limit.ScaleCount)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      64,
			ScaleUnit:  "minutes",
			ScaleCount: 5,
		},
	}

	limit, err := client.Sensory.CreateLimit("source-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if limit.ID != expectedLimit.ID {
		t.Errorf("Expected limit ID %s, got %s", expectedLimit.ID, limit.ID)
	}
}

func TestSensoryCreateLimitValidation(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty source ID validation
	_, err := client.Sensory.CreateLimit("", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      32,
			ScaleUnit:  "seconds",
			ScaleCount: 1,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty source ID")
	}

	// Test empty scale_unit validation
	_, err = client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      32,
			ScaleCount: 1,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty scale_unit")
	}

	// Test invalid scale_count validation
	_, err = client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      32,
			ScaleUnit:  "seconds",
			ScaleCount: 0,
		},
	})
	if err == nil {
		t.Error("Expected validation error for zero scale_count")
	}

	// Test invalid limit value validation
	_, err = client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      0,
			ScaleUnit:  "seconds",
			ScaleCount: 1,
		},
	})
	if err == nil {
		t.Error("Expected validation error for zero limit value")
	}
}
