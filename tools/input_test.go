package tools_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/tools"
)

func TestToolsGetInput(t *testing.T) {
	expectedInput := tools.Input{
		ID:             "input-123",
		Type:           "text",
		ThoughtToolID:  "tool-456",
		ClassCorpusID:  "corpus-789",
		ProvisionState: "active",
	}

	expectedResponse := tools.InputResponse{
		Data: expectedInput,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/inputs/input-123" {
			t.Errorf("Expected path /provision/tools/inputs/input-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	input, err := client.Tools.GetInput("input-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)
}

func TestToolsGetInputError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := tools.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"input": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Tools.GetInput("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var toolsErr *tools.Error
	if errors.As(err, &toolsErr) {
		if toolsErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", toolsErr.StatusCode)
		}
		if toolsErr.Errors == nil || len(toolsErr.Errors["input"]) == 0 ||
			toolsErr.Errors["input"][0] != "not found" {
			t.Errorf("Expected error 'input not found', got %v", toolsErr.Errors)
		}
	} else {
		t.Errorf("Expected tools.Error, got %T", err)
	}
}

func TestToolsCreateInput(t *testing.T) {
	request := tools.CreateInputRequest{
		Input: tools.InputRequestData{
			Type:          "text",
			ClassCorpusID: "corpus-789",
		},
	}

	expectedInput := tools.Input{
		ID:             "input-456",
		Type:           "text",
		ThoughtToolID:  "tool-123",
		ClassCorpusID:  "corpus-789",
		ProvisionState: "pending",
	}

	expectedResponse := tools.InputResponse{
		Data: expectedInput,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/tool-123/inputs" {
			t.Errorf("Expected path /provision/tools/tool-123/inputs, got %s", r.URL.Path)
		}

		var receivedRequest tools.CreateInputRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Input.Type != request.Input.Type {
			t.Errorf("Expected input type %s, got %s", request.Input.Type, receivedRequest.Input.Type)
		}

		if receivedRequest.Input.ClassCorpusID != request.Input.ClassCorpusID {
			t.Errorf(
				"Expected input class_corpus_id %s, got %s",
				request.Input.ClassCorpusID,
				receivedRequest.Input.ClassCorpusID,
			)
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
	input, err := client.Tools.CreateInput("tool-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)
}

func TestToolsCreateInputValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty thought tool ID
	_, err = client.Tools.CreateInput("", tools.CreateInputRequest{
		Input: tools.InputRequestData{Type: "text", ClassCorpusID: "corpus-123"},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought tool ID")
	}

	// Test empty input type
	_, err = client.Tools.CreateInput("tool-123", tools.CreateInputRequest{
		Input: tools.InputRequestData{Type: "", ClassCorpusID: "corpus-123"},
	})
	if err == nil {
		t.Error("Expected validation error for empty input type")
	}

	// Test empty class corpus ID
	_, err = client.Tools.CreateInput("tool-123", tools.CreateInputRequest{
		Input: tools.InputRequestData{Type: "text", ClassCorpusID: ""},
	})
	if err == nil {
		t.Error("Expected validation error for empty class corpus ID")
	}
}

func TestToolsCreateInputFieldValidationDelegated(t *testing.T) {
	// Test that API response with field validation errors is handled correctly
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResp := tools.Error{
			StatusCode: 422,
			Errors: map[string][]string{
				"type":            {"is required", "must be valid"},
				"class_corpus_id": {"can't be blank"},
				"base":            {"thought tool is not active"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Tools.CreateInput("tool-123", tools.CreateInputRequest{
		Input: tools.InputRequestData{Type: "invalid", ClassCorpusID: "invalid"},
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var toolsErr *tools.Error
	if errors.As(err, &toolsErr) {
		if toolsErr.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code 422, got %d", toolsErr.StatusCode)
		}
		if len(toolsErr.Errors["type"]) != 2 {
			t.Errorf("Expected 2 type errors, got %d", len(toolsErr.Errors["type"]))
		}
	} else {
		t.Errorf("Expected tools.Error, got %T", err)
	}
}

func TestToolsUpdateInput(t *testing.T) {
	request := tools.UpdateInputRequest{
		Input: tools.UpdateInputData{
			Type:          "updated-text",
			ClassCorpusID: "corpus-updated",
		},
	}

	expectedInput := tools.Input{
		ID:             "input-123",
		Type:           "updated-text",
		ThoughtToolID:  "tool-456",
		ClassCorpusID:  "corpus-updated",
		ProvisionState: "active",
	}

	expectedResponse := tools.InputResponse{
		Data: expectedInput,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/inputs/input-123" {
			t.Errorf("Expected path /provision/tools/inputs/input-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	input, err := client.Tools.UpdateInput("input-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)
}

func TestToolsReplaceInput(t *testing.T) {
	request := tools.UpdateInputRequest{
		Input: tools.UpdateInputData{
			Type:          "replaced-text",
			ClassCorpusID: "corpus-replaced",
		},
	}

	expectedInput := tools.Input{
		ID:             "input-123",
		Type:           "replaced-text",
		ThoughtToolID:  "tool-456",
		ClassCorpusID:  "corpus-replaced",
		ProvisionState: "active",
	}

	expectedResponse := tools.InputResponse{
		Data: expectedInput,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/inputs/input-123" {
			t.Errorf("Expected path /provision/tools/inputs/input-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	input, err := client.Tools.ReplaceInput("input-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInputResponse(t, *input, expectedInput)
}

func TestToolsDeleteInput(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/inputs/input-123" {
			t.Errorf("Expected path /provision/tools/inputs/input-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	err = client.Tools.DeleteInput("input-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestToolsGetInputEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Tools.GetInput("")
	if err == nil {
		t.Error("Expected validation error for empty input ID in GetInput")
	}
}

func TestToolsUpdateInputEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Tools.UpdateInput("", tools.UpdateInputRequest{
		Input: tools.UpdateInputData{Type: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty input ID in UpdateInput")
	}
}

func TestToolsReplaceInputEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Tools.ReplaceInput("", tools.UpdateInputRequest{
		Input: tools.UpdateInputData{Type: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty input ID in ReplaceInput")
	}
}

func TestToolsDeleteInputEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Tools.DeleteInput("")
	if err == nil {
		t.Error("Expected validation error for empty input ID in DeleteInput")
	}
}

func TestToolsCreateInputWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/tool-123/inputs" {
			t.Errorf("Expected path /provision/tools/tool-123/inputs, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		errorResponse := map[string]any{
			"errors": map[string][]string{
				"type":            {"can't be blank", "is not valid"},
				"class_corpus_id": {"can't be blank", "must exist"},
				"base":            {"thought tool is not active"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
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
	request := tools.CreateInputRequest{
		Input: tools.InputRequestData{
			Type:          "invalid-type",   // Invalid type to trigger server-side validation
			ClassCorpusID: "invalid-corpus", // Invalid corpus ID to trigger server-side validation
		},
	}

	_, err = client.Tools.CreateInput("tool-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "type can't be blank") {
		t.Errorf("Expected error to contain 'type can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "type is not valid") {
		t.Errorf("Expected error to contain 'type is not valid', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "class_corpus_id can't be blank") {
		t.Errorf("Expected error to contain 'class_corpus_id can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "base thought tool is not active") {
		t.Errorf("Expected error to contain 'base thought tool is not active', got %s", errorMsg)
	}
}
