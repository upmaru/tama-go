package module_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/perception/module"
)

// createMockServer creates a test HTTP server with the given handler.
func createMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// ValidateInputResponse validates that actual input matches expected input.
func ValidateInputResponse(t *testing.T, actual, expected module.Input) {
	if actual.ID != expected.ID {
		t.Errorf("Expected input ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.Type != expected.Type {
		t.Errorf("Expected input type %s, got %s", expected.Type, actual.Type)
	}

	if actual.ThoughtModuleID != expected.ThoughtModuleID {
		t.Errorf("Expected input thought_module_id %s, got %s", expected.ThoughtModuleID, actual.ThoughtModuleID)
	}

	if actual.ThoughtID != expected.ThoughtID {
		t.Errorf("Expected input thought_id %s, got %s", expected.ThoughtID, actual.ThoughtID)
	}

	if actual.ClassCorpusID != expected.ClassCorpusID {
		t.Errorf("Expected input class_corpus_id %s, got %s", expected.ClassCorpusID, actual.ClassCorpusID)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected input provision_state %s, got %s", expected.ProvisionState, actual.ProvisionState)
	}
}

func TestModuleGetInput(t *testing.T) {
	expectedInput := module.Input{
		ID:              "input-123",
		Type:            "text",
		ThoughtModuleID: "module-123",
		ThoughtID:       "thought-123",
		ClassCorpusID:   "corpus-123",
		ProvisionState:  "active",
	}

	expectedResponse := module.InputResponse{
		Data: expectedInput,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/module/inputs/input-123" {
			t.Errorf("Expected path /provision/perception/module/inputs/input-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())
	input, err := service.GetInput("input-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)
}

func TestModuleGetInputError(t *testing.T) {
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())
	input, err := service.GetInput("invalid-id")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if input != nil {
		t.Error("Expected nil input on error")
	}
}

func TestModuleCreateInput(t *testing.T) {
	expectedInput := module.Input{
		ID:              "input-123",
		Type:            "text",
		ThoughtModuleID: "module-123",
		ThoughtID:       "thought-123",
		ClassCorpusID:   "corpus-123",
		ProvisionState:  "active",
	}

	expectedResponse := module.InputResponse{
		Data: expectedInput,
	}

	var receivedRequest *module.CreateInputRequest

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/module/inputs" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123/module/inputs, got %s", r.URL.Path)
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
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())

	req := module.CreateInputRequest{
		Input: module.CreateInputData{
			Type:          "text",
			ClassCorpusID: "corpus-123",
		},
	}

	input, err := service.CreateInput("thought-123", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)

	// Validate request payload
	if receivedRequest.Input.Type != req.Input.Type {
		t.Errorf("Expected request type %s, got %s", req.Input.Type, receivedRequest.Input.Type)
	}

	if receivedRequest.Input.ClassCorpusID != req.Input.ClassCorpusID {
		t.Errorf("Expected request class_corpus_id %s, got %s",
			req.Input.ClassCorpusID, receivedRequest.Input.ClassCorpusID)
	}
}

func TestModuleCreateInputValidation(t *testing.T) {
	config := tama.Config{
		BaseURL:        "http://example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())

	req := module.UpdateInputRequest{
		Input: module.UpdateInputData{
			Type:          "text",
			ClassCorpusID: "corpus-456",
		},
	}

	_, err = service.UpdateInput("", req)
	if err == nil || err.Error() != "input ID is required" {
		t.Errorf("Expected 'input ID is required' error, got %v", err)
	}

	// Test empty type
	createReq := module.CreateInputRequest{
		Input: module.CreateInputData{
			Type:          "",
			ClassCorpusID: "corpus-123",
		},
	}

	_, err = service.CreateInput("thought-123", createReq)
	if err == nil || err.Error() != "input type is required" {
		t.Errorf("Expected 'input type is required' error, got %v", err)
	}

	// Test empty class corpus ID
	createReq2 := module.CreateInputRequest{
		Input: module.CreateInputData{
			Type:          "text",
			ClassCorpusID: "",
		},
	}

	_, err = service.CreateInput("thought-123", createReq2)
	if err == nil || err.Error() != "class corpus ID is required" {
		t.Errorf("Expected 'class corpus ID is required' error, got %v", err)
	}
}

func TestModuleUpdateInput(t *testing.T) {
	expectedInput := module.Input{
		ID:              "input-123",
		Type:            "image",
		ThoughtModuleID: "module-123",
		ThoughtID:       "thought-456",
		ClassCorpusID:   "corpus-456",
		ProvisionState:  "active",
	}

	expectedResponse := module.InputResponse{
		Data: expectedInput,
	}

	var receivedRequest *module.UpdateInputRequest

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/module/inputs/input-123" {
			t.Errorf("Expected path /provision/perception/module/inputs/input-123, got %s", r.URL.Path)
		}

		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())

	req := module.UpdateInputRequest{
		Input: module.UpdateInputData{
			Type:          "image",
			ClassCorpusID: "corpus-456",
		},
	}

	input, err := service.UpdateInput("input-123", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)

	// Validate request payload
	if receivedRequest.Input.Type != req.Input.Type {
		t.Errorf("Expected request type %s, got %s", req.Input.Type, receivedRequest.Input.Type)
	}

	if receivedRequest.Input.ClassCorpusID != req.Input.ClassCorpusID {
		t.Errorf("Expected request class_corpus_id %s, got %s",
			req.Input.ClassCorpusID, receivedRequest.Input.ClassCorpusID)
	}
}

func TestModuleReplaceInput(t *testing.T) {
	expectedInput := module.Input{
		ID:              "input-123",
		Type:            "audio",
		ThoughtModuleID: "module-123",
		ThoughtID:       "thought-789",
		ClassCorpusID:   "corpus-789",
		ProvisionState:  "active",
	}

	expectedResponse := module.InputResponse{
		Data: expectedInput,
	}

	var receivedRequest *module.UpdateInputRequest

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/module/inputs/input-123" {
			t.Errorf("Expected path /provision/perception/module/inputs/input-123, got %s", r.URL.Path)
		}

		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())

	req := module.UpdateInputRequest{
		Input: module.UpdateInputData{
			Type:          "audio",
			ClassCorpusID: "corpus-789",
		},
	}

	input, err := service.ReplaceInput("input-123", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)

	// Validate request payload
	if receivedRequest.Input.Type != req.Input.Type {
		t.Errorf("Expected request type %s, got %s", req.Input.Type, receivedRequest.Input.Type)
	}

	if receivedRequest.Input.ClassCorpusID != req.Input.ClassCorpusID {
		t.Errorf("Expected request class_corpus_id %s, got %s",
			req.Input.ClassCorpusID, receivedRequest.Input.ClassCorpusID)
	}
}

func TestModuleDeleteInput(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/module/inputs/input-123" {
			t.Errorf("Expected path /provision/perception/module/inputs/input-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())
	err = service.DeleteInput("input-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestModuleGetInputEmptyID(t *testing.T) {
	config := tama.Config{
		BaseURL:        "http://example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())
	_, err = service.GetInput("")

	if err == nil || err.Error() != "input ID is required" {
		t.Errorf("Expected 'input ID is required' error, got %v", err)
	}
}

func TestModuleUpdateInputEmptyID(t *testing.T) {
	config := tama.Config{
		BaseURL:        "http://example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())

	req := module.UpdateInputRequest{
		Input: module.UpdateInputData{
			Type: "text",
		},
	}

	_, err = service.UpdateInput("", req)

	if err == nil || err.Error() != "input ID is required" {
		t.Errorf("Expected 'input ID is required' error, got %v", err)
	}
}

func TestModuleReplaceInputEmptyID(t *testing.T) {
	config := tama.Config{
		BaseURL:        "http://example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())

	req := module.UpdateInputRequest{
		Input: module.UpdateInputData{
			Type: "text",
		},
	}

	_, err = service.ReplaceInput("", req)

	if err == nil || err.Error() != "input ID is required" {
		t.Errorf("Expected 'input ID is required' error, got %v", err)
	}
}

func TestModuleDeleteInputEmptyID(t *testing.T) {
	config := tama.Config{
		BaseURL:        "http://example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())
	err = service.DeleteInput("")

	if err == nil || err.Error() != "input ID is required" {
		t.Errorf("Expected 'input ID is required' error, got %v", err)
	}
}

func TestModuleCreateInputWithFieldErrors(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := map[string]interface{}{
			"errors": map[string][]string{
				"type":            {"is required"},
				"class_corpus_id": {"is invalid"},
			},
		}
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	service := module.NewService(client.GetHTTPClient())

	req := module.CreateInputRequest{
		Input: module.CreateInputData{
			Type:          "text",
			ClassCorpusID: "invalid-id",
		},
	}

	input, err := service.CreateInput("thought-123", req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if input != nil {
		t.Error("Expected nil input on error")
	}

	// Check that error message contains expected validation errors
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "type is required") {
		t.Errorf("Expected error to contain 'type is required', got: %s", errorMsg)
	}

	if !strings.Contains(errorMsg, "class_corpus_id is invalid") {
		t.Errorf("Expected error to contain 'class_corpus_id is invalid', got: %s", errorMsg)
	}
}
