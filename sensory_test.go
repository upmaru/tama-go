package tama_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/sensory"
)

func TestSensoryGetSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-123",
		Name:           "Test Source",
		Endpoint:       "https://api.test.com/v1",
		SpaceID:        "space-456",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123" {
			t.Errorf("Expected path /provision/sensory/sources/source-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
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

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryCreateSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-789",
		Name:           "New Source",
		Endpoint:       "https://api.mistral.ai/v1",
		SpaceID:        "space-123",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
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

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "New Source",
			Type:     "model",
			Endpoint: "https://api.mistral.ai/v1",
			Credential: sensory.SourceCredential{
				APIKey: "test-api-key",
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

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryCreateSourceValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
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
				APIKey: "test-api-key",
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
				APIKey: "test-key",
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
				APIKey: "test-key",
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
				APIKey: "test-key",
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestSensoryGetSource_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.GetSource("")
	if err == nil {
		t.Error("Expected validation error for empty source ID in GetSource")
	}
}

func TestSensoryGetModel(t *testing.T) {
	expectedModel := sensory.Model{
		ID:         "model-123",
		Identifier: "mistral-small-latest",
		Path:       "/chat/completions",
		Parameters: map[string]any{
			"temperature": 0.7,
			"max_tokens":  1000.0,
		},
		ProvisionState: "active",
	}

	expectedResponse := sensory.ModelResponse{
		Data: expectedModel,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/models/model-123" {
			t.Errorf("Expected path /provision/sensory/models/model-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
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

	if model.Path != expectedModel.Path {
		t.Errorf("Expected model path %s, got %s", expectedModel.Path, model.Path)
	}

	if model.ProvisionState != expectedModel.ProvisionState {
		t.Errorf("Expected model provision state %s, got %s", expectedModel.ProvisionState, model.ProvisionState)
	}

	if len(model.Parameters) != len(expectedModel.Parameters) {
		t.Errorf("Expected %d parameters, got %d", len(expectedModel.Parameters), len(model.Parameters))
	}

	for key, expectedValue := range expectedModel.Parameters {
		if actualValue, exists := model.Parameters[key]; !exists {
			t.Errorf("Expected parameter %s not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected parameter %s to be %v, got %v", key, expectedValue, actualValue)
		}
	}
}

func TestSensoryCreateModel(t *testing.T) {
	expectedModel := sensory.Model{
		ID:         "model-789",
		Identifier: "mistral-large-latest",
		Path:       "/chat/completions",
		Parameters: map[string]any{
			"reasoning_effort": "low",
			"temperature":      1.0,
		},
		ProvisionState: "active",
	}

	expectedResponse := sensory.ModelResponse{
		Data: expectedModel,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123/models" {
			t.Errorf("Expected path /provision/sensory/sources/source-123/models, got %s", r.URL.Path)
		}

		var req sensory.CreateModelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		validateModelRequest(t, req, "mistral-large-latest", "/chat/completions")

		expectedParams := map[string]any{
			"reasoning_effort": "low",
			"temperature":      1.0,
		}
		if len(req.Model.Parameters) != len(expectedParams) {
			t.Errorf("Expected %d parameters, got %d", len(expectedParams), len(req.Model.Parameters))
		}
		validateModelParameters(t, req.Model.Parameters, expectedParams)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := sensory.CreateModelRequest{
		Model: sensory.ModelRequestData{
			Identifier: "mistral-large-latest",
			Path:       "/chat/completions",
			Parameters: map[string]any{
				"reasoning_effort": "low",
				"temperature":      1.0,
			},
		},
	}

	model, err := client.Sensory.CreateModel("source-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validateModelResponse(t, *model, expectedModel)
	validateModelParameters(t, model.Parameters, expectedModel.Parameters)
}

func TestSensoryCreateModelValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
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

func TestSensoryGetModel_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.GetModel("")
	if err == nil {
		t.Error("Expected validation error for empty model ID in GetModel")
	}
}

func TestSensoryGetLimit(t *testing.T) {
	expectedLimit := sensory.Limit{
		ID:             "limit-123",
		SourceID:       "source-456",
		Count:          32,
		ScaleUnit:      "seconds",
		ScaleCount:     1,
		ProvisionState: "active",
	}

	expectedResponse := sensory.LimitResponse{
		Data: expectedLimit,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/limits/limit-123" {
			t.Errorf("Expected path /provision/sensory/limits/limit-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
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
		ID:             "limit-789",
		SourceID:       "source-123",
		Count:          64,
		ScaleUnit:      "minutes",
		ScaleCount:     5,
		ProvisionState: "active",
	}

	expectedResponse := sensory.LimitResponse{
		Data: expectedLimit,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
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

	client := tama.NewClient(tama.Config{
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
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
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

	// Test invalid count value validation
	_, err = client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      0,
			ScaleUnit:  "seconds",
			ScaleCount: 1,
		},
	})
	if err == nil {
		t.Error("Expected validation error for zero count value")
	}
}

func TestSensoryGetLimit_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.GetLimit("")
	if err == nil {
		t.Error("Expected validation error for empty limit ID in GetLimit")
	}
}

func TestSensoryUpdateSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-123",
		Name:           "Updated Source",
		Endpoint:       "https://api.updated.com/v1",
		SpaceID:        "space-456",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123" {
			t.Errorf("Expected path /provision/sensory/sources/source-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Source.Name != "Updated Source" {
			t.Errorf("Expected request name 'Updated Source', got %s", req.Source.Name)
		}

		if req.Source.Endpoint != "https://api.updated.com/v1" {
			t.Errorf("Expected request endpoint 'https://api.updated.com/v1', got %s", req.Source.Endpoint)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	updateReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name:     "Updated Source",
			Endpoint: "https://api.updated.com/v1",
			Credential: &sensory.SourceCredential{
				APIKey: "updated-api-key",
			},
		},
	}

	source, err := client.Sensory.UpdateSource("source-123", updateReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryUpdateSource_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	updateReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name: "Updated Source",
		},
	}

	_, err := client.Sensory.UpdateSource("", updateReq)
	if err == nil {
		t.Error("Expected validation error for empty source ID in UpdateSource")
	}
}

func TestSensoryReplaceSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-123",
		Name:           "Replaced Source",
		Endpoint:       "https://api.replaced.com/v1",
		SpaceID:        "space-456",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123" {
			t.Errorf("Expected path /provision/sensory/sources/source-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Source.Name != "Replaced Source" {
			t.Errorf("Expected request name 'Replaced Source', got %s", req.Source.Name)
		}

		if req.Source.Endpoint != "https://api.replaced.com/v1" {
			t.Errorf("Expected request endpoint 'https://api.replaced.com/v1', got %s", req.Source.Endpoint)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	replaceReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name:     "Replaced Source",
			Type:     "model",
			Endpoint: "https://api.replaced.com/v1",
			Credential: &sensory.SourceCredential{
				APIKey: "replaced-api-key",
			},
		},
	}

	source, err := client.Sensory.ReplaceSource("source-123", replaceReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryReplaceSource_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	replaceReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name: "Replaced Source",
		},
	}

	_, err := client.Sensory.ReplaceSource("", replaceReq)
	if err == nil {
		t.Error("Expected validation error for empty source ID in ReplaceSource")
	}
}

func TestSensoryUpdateModel(t *testing.T) {
	expectedModel := sensory.Model{
		ID:         "model-123",
		Identifier: "mistral-large-updated",
		Path:       "/v1/chat/completions",
		Parameters: map[string]any{
			"max_tokens": 2000.0,
			"top_p":      0.95,
		},
		ProvisionState: "active",
	}

	expectedResponse := sensory.ModelResponse{
		Data: expectedModel,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/models/model-123" {
			t.Errorf("Expected path /provision/sensory/models/model-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateModelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Model.Identifier != "mistral-large-updated" {
			t.Errorf("Expected request identifier 'mistral-large-updated', got %s", req.Model.Identifier)
		}

		if req.Model.Path != "/v1/chat/completions" {
			t.Errorf("Expected request path '/v1/chat/completions', got %s", req.Model.Path)
		}

		expectedParams := map[string]any{
			"max_tokens": 2000.0,
			"top_p":      0.95,
		}
		validateModelParameters(t, req.Model.Parameters, expectedParams)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	updateReq := sensory.UpdateModelRequest{
		Model: sensory.UpdateModelData{
			Identifier: "mistral-large-updated",
			Path:       "/v1/chat/completions",
			Parameters: map[string]any{
				"max_tokens": 2000.0,
				"top_p":      0.95,
			},
		},
	}

	model, err := client.Sensory.UpdateModel("model-123", updateReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validateModelResponse(t, *model, expectedModel)
	validateModelParameters(t, model.Parameters, expectedModel.Parameters)
}

func TestSensoryUpdateModel_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	updateReq := sensory.UpdateModelRequest{
		Model: sensory.UpdateModelData{
			Identifier: "updated-model",
		},
	}

	_, err := client.Sensory.UpdateModel("", updateReq)
	if err == nil {
		t.Error("Expected validation error for empty model ID in UpdateModel")
	}
}

func TestSensoryReplaceModel(t *testing.T) {
	expectedModel := sensory.Model{
		ID:         "model-123",
		Identifier: "mistral-large-replaced",
		Path:       "/v2/chat/completions",
		Parameters: map[string]any{
			"stream":      true,
			"temperature": 0.5,
		},
		ProvisionState: "active",
	}

	expectedResponse := sensory.ModelResponse{
		Data: expectedModel,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/models/model-123" {
			t.Errorf("Expected path /provision/sensory/models/model-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateModelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Model.Identifier != "mistral-large-replaced" {
			t.Errorf("Expected request identifier 'mistral-large-replaced', got %s", req.Model.Identifier)
		}

		if req.Model.Path != "/v2/chat/completions" {
			t.Errorf("Expected request path '/v2/chat/completions', got %s", req.Model.Path)
		}

		expectedParams := map[string]any{
			"stream":      true,
			"temperature": 0.5,
		}
		validateModelParameters(t, req.Model.Parameters, expectedParams)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	replaceReq := sensory.UpdateModelRequest{
		Model: sensory.UpdateModelData{
			Identifier: "mistral-large-replaced",
			Path:       "/v2/chat/completions",
			Parameters: map[string]any{
				"stream":      true,
				"temperature": 0.5,
			},
		},
	}

	model, err := client.Sensory.ReplaceModel("model-123", replaceReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validateModelResponse(t, *model, expectedModel)
	validateModelParameters(t, model.Parameters, expectedModel.Parameters)
}

func TestSensoryModelParameters(t *testing.T) {
	expectedModel := createTestModelWithParameters()
	expectedResponse := sensory.ModelResponse{Data: expectedModel}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		validateParametersRequest(t, r, w, expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := createTestModelRequest()
	model, err := client.Sensory.CreateModel("source-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validateModelResponse(t, *model, expectedModel)
	validateComplexParameters(t, model.Parameters, expectedModel.Parameters)
}

func TestSensoryReplaceModel_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	replaceReq := sensory.UpdateModelRequest{
		Model: sensory.UpdateModelData{
			Identifier: "replaced-model",
		},
	}

	_, err := client.Sensory.ReplaceModel("", replaceReq)
	if err == nil {
		t.Error("Expected validation error for empty model ID in ReplaceModel")
	}
}

func TestSensoryFieldSpecificErrors(t *testing.T) {
	// Test sensory field-specific errors
	fieldErr := &sensory.Error{
		StatusCode: 422,
		Errors: map[string][]string{
			"source_id": {"has already been taken"},
			"name":      {"is required", "must be at least 3 characters"},
		},
	}

	errorMsg := fieldErr.Error()
	// Check that all field errors are included
	if !strings.Contains(errorMsg, "source_id has already been taken") {
		t.Errorf("Expected error message to contain 'source_id has already been taken', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "name is required") {
		t.Errorf("Expected error message to contain 'name is required', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "name must be at least 3 characters") {
		t.Errorf("Expected error message to contain 'name must be at least 3 characters', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "API error 422:") {
		t.Errorf("Expected error message to contain status code, got %s", errorMsg)
	}

	// Test error with only status code
	statusOnlyErr := &sensory.Error{
		StatusCode: 404,
	}

	expectedStatusMsg := "API error 404"
	if statusOnlyErr.Error() != expectedStatusMsg {
		t.Errorf("Expected error message %s, got %s", expectedStatusMsg, statusOnlyErr.Error())
	}

	// Test field-specific errors without status code
	fieldErrNoStatus := &sensory.Error{
		Errors: map[string][]string{
			"endpoint": {"is invalid URL"},
		},
	}

	errorMsgNoStatus := fieldErrNoStatus.Error()
	expectedNoStatus := "API error: endpoint is invalid URL"
	if errorMsgNoStatus != expectedNoStatus {
		t.Errorf("Expected error message %s, got %s", expectedNoStatus, errorMsgNoStatus)
	}
}

// Helper functions to reduce cognitive complexity

func validateModelRequest(t *testing.T, req sensory.CreateModelRequest, expectedIdentifier, expectedPath string) {
	if req.Model.Identifier != expectedIdentifier {
		t.Errorf("Expected request identifier '%s', got %s", expectedIdentifier, req.Model.Identifier)
	}
	if req.Model.Path != expectedPath {
		t.Errorf("Expected request path '%s', got %s", expectedPath, req.Model.Path)
	}
}

func validateModelResponse(t *testing.T, actual, expected sensory.Model) {
	if actual.ID != expected.ID {
		t.Errorf("Expected model ID %s, got %s", expected.ID, actual.ID)
	}
	if actual.Identifier != expected.Identifier {
		t.Errorf("Expected model identifier %s, got %s", expected.Identifier, actual.Identifier)
	}
	if actual.Path != expected.Path {
		t.Errorf("Expected model path %s, got %s", expected.Path, actual.Path)
	}
	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected model provision state %s, got %s", expected.ProvisionState, actual.ProvisionState)
	}
}

func validateModelParameters(t *testing.T, actual map[string]any, expected map[string]any) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %d parameters, got %d", len(expected), len(actual))
	}
	for key, expectedValue := range expected {
		if actualValue, exists := actual[key]; !exists {
			t.Errorf("Expected parameter %s not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected parameter %s to be %v, got %v", key, expectedValue, actualValue)
		}
	}
}

func validateComplexParameters(t *testing.T, actual map[string]any, expected map[string]any) {
	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			t.Errorf("Expected parameter %s not found in response", key)
			continue
		}

		switch key {
		case "stop":
			validateArrayParameter(t, key, actualValue, expectedValue)
		case "config":
			validateObjectParameter(t, key, actualValue, expectedValue)
		default:
			if actualValue != expectedValue {
				t.Errorf("Expected parameter %s to be %v, got %v", key, expectedValue, actualValue)
			}
		}
	}
}

func validateArrayParameter(t *testing.T, key string, actual any, expected any) {
	expectedSlice := expected.([]string)
	actualSlice, ok := actual.([]any)
	if !ok {
		t.Errorf("Expected %s to be array, got %T", key, actual)
		return
	}
	if len(actualSlice) != len(expectedSlice) {
		t.Errorf("Expected %s array length %d, got %d", key, len(expectedSlice), len(actualSlice))
		return
	}
	for i, expectedItem := range expectedSlice {
		if actualSlice[i] != expectedItem {
			t.Errorf("Expected %s[%d] to be %v, got %v", key, i, expectedItem, actualSlice[i])
		}
	}
}

func validateObjectParameter(t *testing.T, key string, actual any, expected any) {
	expectedMap := expected.(map[string]any)
	actualMap, ok := actual.(map[string]any)
	if !ok {
		t.Errorf("Expected %s to be object, got %T", key, actual)
		return
	}
	for configKey, configExpected := range expectedMap {
		if actualMap[configKey] != configExpected {
			t.Errorf("Expected %s.%s to be %v, got %v", key, configKey, configExpected, actualMap[configKey])
		}
	}
}

func createTestModelWithParameters() sensory.Model {
	return sensory.Model{
		ID:         "model-params-123",
		Identifier: "test-model-with-params",
		Path:       "/test/completions",
		Parameters: map[string]any{
			"temperature":       0.8,
			"max_tokens":        1500.0,
			"top_p":             0.9,
			"frequency_penalty": 0.1,
			"presence_penalty":  0.2,
			"stream":            true,
			"stop":              []string{"\\n", "###"},
			"reasoning_effort":  "medium",
			"config": map[string]any{
				"enable_cache": true,
				"timeout":      30.0,
			},
		},
		ProvisionState: "active",
	}
}

func createTestModelRequest() sensory.CreateModelRequest {
	return sensory.CreateModelRequest{
		Model: sensory.ModelRequestData{
			Identifier: "test-model-with-params",
			Path:       "/test/completions",
			Parameters: map[string]any{
				"temperature":       0.8,
				"max_tokens":        1500.0,
				"top_p":             0.9,
				"frequency_penalty": 0.1,
				"presence_penalty":  0.2,
				"stream":            true,
				"stop":              []string{"\\n", "###"},
				"reasoning_effort":  "medium",
				"config": map[string]any{
					"enable_cache": true,
					"timeout":      30.0,
				},
			},
		},
	}
}

func validateParametersRequest(t *testing.T, r *http.Request, w http.ResponseWriter,
	expectedResponse sensory.ModelResponse) {
	if r.Method != http.MethodPost {
		t.Errorf("Expected POST request, got %s", r.Method)
	}

	var req sensory.CreateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatalf("Failed to decode request body: %v", err)
	}

	validateRequestBasicParams(t, req.Model.Parameters)
	validateRequestArrayParam(t, req.Model.Parameters)
	validateRequestObjectParam(t, req.Model.Parameters)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expectedResponse)
}

func validateRequestBasicParams(t *testing.T, params map[string]any) {
	basicParams := map[string]any{
		"temperature":      0.8,
		"max_tokens":       1500.0,
		"stream":           true,
		"reasoning_effort": "medium",
	}
	for key, expected := range basicParams {
		if params[key] != expected {
			t.Errorf("Expected %s %v, got %v", key, expected, params[key])
		}
	}
}

func validateRequestArrayParam(t *testing.T, params map[string]any) {
	stop, ok := params["stop"].([]any)
	if !ok {
		t.Errorf("Expected stop to be an array, got %T", params["stop"])
	} else if len(stop) != 2 || stop[0] != "\\n" || stop[1] != "###" {
		t.Errorf("Expected stop array ['\\n', '###'], got %v", stop)
	}
}

func validateRequestObjectParam(t *testing.T, params map[string]any) {
	config, ok := params["config"].(map[string]any)
	if !ok {
		t.Errorf("Expected config to be an object, got %T", params["config"])
		return
	}
	if config["enable_cache"] != true {
		t.Errorf("Expected config.enable_cache true, got %v", config["enable_cache"])
	}
	if config["timeout"] != 30.0 {
		t.Errorf("Expected config.timeout 30.0, got %v", config["timeout"])
	}
}

func TestSensoryCreateSourceWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/spaces/space-123/sources" {
			t.Errorf("Expected path /provision/sensory/spaces/space-123/sources, got %s", r.URL.Path)
		}

		// Return field validation errors
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := map[string]any{
			"errors": map[string][]string{
				"name":     {"is required"},
				"endpoint": {"is invalid URL", "must use HTTPS"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "test-source", // Valid name to bypass client validation
			Type:     "ollama",
			Endpoint: "https://valid-endpoint.com", // Valid endpoint to bypass client validation
			Credential: sensory.SourceCredential{
				APIKey: "test-key",
			},
		},
	}

	_, err := client.Sensory.CreateSource("space-123", createReq)
	if err == nil {
		t.Fatal("Expected error for invalid source data")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "name is required") {
		t.Errorf("Expected error to contain 'name is required', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "endpoint is invalid URL") {
		t.Errorf("Expected error to contain 'endpoint is invalid URL', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "endpoint must use HTTPS") {
		t.Errorf("Expected error to contain 'endpoint must use HTTPS', got %s", errorMsg)
	}
}

func TestSensoryGetSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-123",
		SpaceID: "space-456",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type": "string",
				},
			},
		},
		Version:        "1.0.0",
		Endpoint:       "https://api.test.com/v1",
		CurrentState:   "active",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	spec, err := client.Sensory.GetSpecification("spec-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.ID != expectedSpec.ID {
		t.Errorf("Expected specification ID %s, got %s", expectedSpec.ID, spec.ID)
	}

	if spec.SpaceID != expectedSpec.SpaceID {
		t.Errorf("Expected specification space ID %s, got %s", expectedSpec.SpaceID, spec.SpaceID)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}

	if spec.CurrentState != expectedSpec.CurrentState {
		t.Errorf("Expected specification current state %s, got %s", expectedSpec.CurrentState, spec.CurrentState)
	}

	if spec.ProvisionState != expectedSpec.ProvisionState {
		t.Errorf("Expected specification provision state %s, got %s", expectedSpec.ProvisionState, spec.ProvisionState)
	}

	// Validate schema structure
	if spec.Schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got %v", spec.Schema["type"])
	}
}

func TestSensoryCreateSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-789",
		SpaceID: "space-123",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]any{
					"type": "string",
				},
			},
		},
		Version:        "2.0.0",
		Endpoint:       "https://api.example.com/v2",
		CurrentState:   "pending",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/spaces/space-123/specifications" {
			t.Errorf("Expected path /provision/sensory/spaces/space-123/specifications, got %s", r.URL.Path)
		}

		var req sensory.CreateSpecificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Specification.Version != "2.0.0" {
			t.Errorf("Expected request version '2.0.0', got %s", req.Specification.Version)
		}

		if req.Specification.Endpoint != "https://api.example.com/v2" {
			t.Errorf("Expected request endpoint 'https://api.example.com/v2', got %s", req.Specification.Endpoint)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{
						"type": "string",
					},
				},
			},
			Version:  "2.0.0",
			Endpoint: "https://api.example.com/v2",
		},
	}

	spec, err := client.Sensory.CreateSpecification("space-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.ID != expectedSpec.ID {
		t.Errorf("Expected specification ID %s, got %s", expectedSpec.ID, spec.ID)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}
}

func TestSensoryCreateSpecificationValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty space ID validation
	_, err := client.Sensory.CreateSpecification("", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
			},
			Version:  "1.0.0",
			Endpoint: "https://api.test.com",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty schema validation
	_, err = client.Sensory.CreateSpecification("space-123", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Version:  "1.0.0",
			Endpoint: "https://api.test.com",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty schema")
	}

	// Test empty version validation
	_, err = client.Sensory.CreateSpecification("space-123", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
			},
			Endpoint: "https://api.test.com",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty version")
	}

	// Test empty endpoint validation
	_, err = client.Sensory.CreateSpecification("space-123", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
			},
			Version: "1.0.0",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestSensoryGetSpecification_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.GetSpecification("")
	if err == nil {
		t.Error("Expected validation error for empty specification ID in GetSpecification")
	}
}

func TestSensoryUpdateSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-123",
		SpaceID: "space-456",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"updated_field": map[string]any{
					"type": "string",
				},
			},
		},
		Version:        "1.1.0",
		Endpoint:       "https://api.updated.com/v1",
		CurrentState:   "active",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSpecificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Specification.Version != "1.1.0" {
			t.Errorf("Expected request version '1.1.0', got %s", req.Specification.Version)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	updateReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"updated_field": map[string]any{
						"type": "string",
					},
				},
			},
			Version:  "1.1.0",
			Endpoint: "https://api.updated.com/v1",
		},
	}

	spec, err := client.Sensory.UpdateSpecification("spec-123", updateReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}
}

func TestSensoryUpdateSpecification_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	updateReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Version: "1.1.0",
		},
	}

	_, err := client.Sensory.UpdateSpecification("", updateReq)
	if err == nil {
		t.Error("Expected validation error for empty specification ID in UpdateSpecification")
	}
}

func TestSensoryReplaceSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-123",
		SpaceID: "space-456",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"replaced_field": map[string]any{
					"type": "number",
				},
			},
		},
		Version:        "2.0.0",
		Endpoint:       "https://api.replaced.com/v2",
		CurrentState:   "pending",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSpecificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Specification.Version != "2.0.0" {
			t.Errorf("Expected request version '2.0.0', got %s", req.Specification.Version)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	replaceReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"replaced_field": map[string]any{
						"type": "number",
					},
				},
			},
			Version:  "2.0.0",
			Endpoint: "https://api.replaced.com/v2",
		},
	}

	spec, err := client.Sensory.ReplaceSpecification("spec-123", replaceReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}
}

func TestSensoryReplaceSpecification_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	replaceReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Version: "2.0.0",
		},
	}

	_, err := client.Sensory.ReplaceSpecification("", replaceReq)
	if err == nil {
		t.Error("Expected validation error for empty specification ID in ReplaceSpecification")
	}
}

func TestSensoryDeleteSpecification(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	err := client.Sensory.DeleteSpecification("spec-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSensoryDeleteSpecification_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Sensory.DeleteSpecification("")
	if err == nil {
		t.Error("Expected validation error for empty specification ID in DeleteSpecification")
	}
}

func TestSensoryGetIdentity(t *testing.T) {
	expectedIdentity := sensory.Identity{
		ID:              "identity-123",
		SpecificationID: "spec-456",
		ProvisionState:  "active",
		CurrentState:    "running",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/health",
			Method: "GET",
			Codes:  []int{200},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	identity, err := client.Sensory.GetIdentity("identity-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.SpecificationID != expectedIdentity.SpecificationID {
		t.Errorf("Expected specification ID %s, got %s", expectedIdentity.SpecificationID, identity.SpecificationID)
	}

	if identity.ProvisionState != expectedIdentity.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expectedIdentity.ProvisionState, identity.ProvisionState)
	}

	if identity.CurrentState != expectedIdentity.CurrentState {
		t.Errorf("Expected current state %s, got %s", expectedIdentity.CurrentState, identity.CurrentState)
	}

	if identity.Identifier != expectedIdentity.Identifier {
		t.Errorf("Expected identifier %s, got %s", expectedIdentity.Identifier, identity.Identifier)
	}

	if identity.Validation.Path != expectedIdentity.Validation.Path {
		t.Errorf("Expected validation path %s, got %s", expectedIdentity.Validation.Path, identity.Validation.Path)
	}

	if identity.Validation.Method != expectedIdentity.Validation.Method {
		t.Errorf("Expected validation method %s, got %s",
			expectedIdentity.Validation.Method, identity.Validation.Method)
	}

	if len(identity.Validation.Codes) != len(expectedIdentity.Validation.Codes) {
		t.Errorf("Expected %d validation codes, got %d",
			len(expectedIdentity.Validation.Codes), len(identity.Validation.Codes))
	}

	if identity.Validation.Codes[0] != expectedIdentity.Validation.Codes[0] {
		t.Errorf("Expected validation code %d, got %d",
			expectedIdentity.Validation.Codes[0], identity.Validation.Codes[0])
	}
}

func TestSensoryCreateIdentity(t *testing.T) {
	requestData := sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-api-key",
			Validation: sensory.Validation{
				Path:   "/health",
				Method: "GET",
				Codes:  []int{200},
			},
		},
	}

	expectedIdentity := sensory.Identity{
		ID:              "identity-456",
		SpecificationID: "spec-123",
		ProvisionState:  "pending",
		CurrentState:    "initializing",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/health",
			Method: "GET",
			Codes:  []int{200},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		expectedPath := "/provision/sensory/specifications/spec-123/identifiers/test-identifier/identities"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var receivedRequest sensory.CreateIdentityRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Identity.APIKey != requestData.Identity.APIKey {
			t.Errorf("Expected API key %s, got %s", requestData.Identity.APIKey, receivedRequest.Identity.APIKey)
		}

		if receivedRequest.Identity.Validation.Path != requestData.Identity.Validation.Path {
			t.Errorf("Expected validation path %s, got %s",
				requestData.Identity.Validation.Path, receivedRequest.Identity.Validation.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	identity, err := client.Sensory.CreateIdentity("spec-123", "test-identifier", requestData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.SpecificationID != expectedIdentity.SpecificationID {
		t.Errorf("Expected specification ID %s, got %s", expectedIdentity.SpecificationID, identity.SpecificationID)
	}

	if identity.ProvisionState != expectedIdentity.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expectedIdentity.ProvisionState, identity.ProvisionState)
	}
}

func TestSensoryCreateIdentityValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost:8080",
		APIKey:  "test-key",
	})

	// Test empty specification ID
	_, err := client.Sensory.CreateIdentity("", "identifier", sensory.CreateIdentityRequest{})
	if err == nil || err.Error() != "specification ID is required" {
		t.Errorf("Expected 'specification ID is required' error, got: %v", err)
	}

	// Test empty identifier
	_, err = client.Sensory.CreateIdentity("spec-123", "", sensory.CreateIdentityRequest{})
	if err == nil || err.Error() != "identifier is required" {
		t.Errorf("Expected 'identifier is required' error, got: %v", err)
	}

	// Test empty API key
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{},
	})
	if err == nil || err.Error() != "API key is required" {
		t.Errorf("Expected 'API key is required' error, got: %v", err)
	}

	// Test empty validation path
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-key",
			Validation: sensory.Validation{
				Method: "GET",
				Codes:  []int{200},
			},
		},
	})
	if err == nil || err.Error() != "validation path is required" {
		t.Errorf("Expected 'validation path is required' error, got: %v", err)
	}

	// Test empty validation method
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-key",
			Validation: sensory.Validation{
				Path:  "/health",
				Codes: []int{200},
			},
		},
	})
	if err == nil || err.Error() != "validation method is required" {
		t.Errorf("Expected 'validation method is required' error, got: %v", err)
	}

	// Test empty validation codes
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-key",
			Validation: sensory.Validation{
				Path:   "/health",
				Method: "GET",
			},
		},
	})
	if err == nil || err.Error() != "validation codes are required" {
		t.Errorf("Expected 'validation codes are required' error, got: %v", err)
	}
}

func TestSensoryGetIdentity_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost:8080",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.GetIdentity("")
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}

func TestSensoryUpdateIdentity(t *testing.T) {
	requestData := sensory.UpdateIdentityRequest{
		Identity: sensory.UpdateIdentityData{
			APIKey: "updated-api-key",
			Validation: &sensory.Validation{
				Path:   "/health/check",
				Method: "POST",
				Codes:  []int{200, 201},
			},
		},
	}

	expectedIdentity := sensory.Identity{
		ID:              "identity-123",
		SpecificationID: "spec-456",
		ProvisionState:  "active",
		CurrentState:    "running",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/health/check",
			Method: "POST",
			Codes:  []int{200, 201},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
		}

		var receivedRequest sensory.UpdateIdentityRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Identity.APIKey != requestData.Identity.APIKey {
			t.Errorf("Expected API key %s, got %s", requestData.Identity.APIKey, receivedRequest.Identity.APIKey)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	identity, err := client.Sensory.UpdateIdentity("identity-123", requestData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.Validation.Path != expectedIdentity.Validation.Path {
		t.Errorf("Expected validation path %s, got %s", expectedIdentity.Validation.Path, identity.Validation.Path)
	}

	if identity.Validation.Method != expectedIdentity.Validation.Method {
		t.Errorf("Expected validation method %s, got %s",
			expectedIdentity.Validation.Method, identity.Validation.Method)
	}
}

func TestSensoryUpdateIdentity_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost:8080",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.UpdateIdentity("", sensory.UpdateIdentityRequest{})
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}

func TestSensoryReplaceIdentity(t *testing.T) {
	requestData := sensory.UpdateIdentityRequest{
		Identity: sensory.UpdateIdentityData{
			APIKey: "replaced-api-key",
			Validation: &sensory.Validation{
				Path:   "/status",
				Method: "GET",
				Codes:  []int{200, 204},
			},
		},
	}

	expectedIdentity := sensory.Identity{
		ID:              "identity-123",
		SpecificationID: "spec-456",
		ProvisionState:  "active",
		CurrentState:    "running",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/status",
			Method: "GET",
			Codes:  []int{200, 204},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
		}

		var receivedRequest sensory.UpdateIdentityRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Identity.APIKey != requestData.Identity.APIKey {
			t.Errorf("Expected API key %s, got %s", requestData.Identity.APIKey, receivedRequest.Identity.APIKey)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	identity, err := client.Sensory.ReplaceIdentity("identity-123", requestData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.Validation.Path != expectedIdentity.Validation.Path {
		t.Errorf("Expected validation path %s, got %s", expectedIdentity.Validation.Path, identity.Validation.Path)
	}
}

func TestSensoryReplaceIdentity_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost:8080",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.ReplaceIdentity("", sensory.UpdateIdentityRequest{})
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}

func TestSensoryDeleteIdentity(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	err := client.Sensory.DeleteIdentity("identity-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSensoryDeleteIdentity_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost:8080",
		APIKey:  "test-key",
	})

	err := client.Sensory.DeleteIdentity("")
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}
