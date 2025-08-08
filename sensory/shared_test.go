package sensory_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/upmaru/tama-go/sensory"
)

// CreateMockServer creates a test HTTP server with the given handler.
func CreateMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// validateModelRequest validates model request data.
func validateModelRequest(t *testing.T, req sensory.CreateModelRequest, expectedIdentifier, expectedPath string) {
	if req.Model.Identifier != expectedIdentifier {
		t.Errorf("Expected request identifier '%s', got %s", expectedIdentifier, req.Model.Identifier)
	}
	if req.Model.Path != expectedPath {
		t.Errorf("Expected request path '%s', got %s", expectedPath, req.Model.Path)
	}
}

// validateModelResponse validates model response data.
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

// validateModelParameters validates model parameters.
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

// validateComplexParameters validates complex model parameters.
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

// validateArrayParameter validates array parameters.
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

// validateObjectParameter validates object parameters.
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

// createTestModelWithParameters creates a test model with complex parameters.
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

// createTestModelRequest creates a test model request with complex parameters.
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

// validateParametersRequest validates complex parameters request.
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

// validateRequestBasicParams validates basic parameters in request.
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

// validateRequestArrayParam validates array parameters in request.
func validateRequestArrayParam(t *testing.T, params map[string]any) {
	stop, ok := params["stop"].([]any)
	if !ok {
		t.Errorf("Expected stop to be an array, got %T", params["stop"])
	} else if len(stop) != 2 || stop[0] != "\\n" || stop[1] != "###" {
		t.Errorf("Expected stop array ['\\n', '###'], got %v", stop)
	}
}

// validateRequestObjectParam validates object parameters in request.
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
