package sensory_test

import (
	"encoding/json"
	"net/http"
	"testing"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/sensory"
)

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

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

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

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

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
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty source ID validation
	_, err = client.Sensory.CreateModel("", sensory.CreateModelRequest{
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
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Sensory.GetModel("")
	if err == nil {
		t.Error("Expected validation error for empty model ID in GetModel")
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

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

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
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	updateReq := sensory.UpdateModelRequest{
		Model: sensory.UpdateModelData{
			Identifier: "updated-model",
		},
	}

	_, err = client.Sensory.UpdateModel("", updateReq)
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

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

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

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	createReq := createTestModelRequest()
	model, err := client.Sensory.CreateModel("source-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	validateModelResponse(t, *model, expectedModel)
	validateComplexParameters(t, model.Parameters, expectedModel.Parameters)
}

func TestSensoryReplaceModel_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	replaceReq := sensory.UpdateModelRequest{
		Model: sensory.UpdateModelData{
			Identifier: "replaced-model",
		},
	}

	_, err = client.Sensory.ReplaceModel("", replaceReq)
	if err == nil {
		t.Error("Expected validation error for empty model ID in ReplaceModel")
	}
}
