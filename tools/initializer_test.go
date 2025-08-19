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

func TestToolsGetInitializer(t *testing.T) {
	expectedInitializer := tools.Initializer{
		ID:             "initializer-123",
		Reference:      "test-reference",
		Index:          0,
		Parameters:     map[string]any{"temperature": 0.7, "max_tokens": 100},
		ThoughtToolID:  "tool-456",
		ProvisionState: "active",
	}

	expectedResponse := tools.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/initializers/initializer-123" {
			t.Errorf("Expected path /provision/tools/initializers/initializer-123, got %s", r.URL.Path)
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

	client := tama.NewClient(config)
	initializer, err := client.Tools.GetInitializer("initializer-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestToolsGetInitializerError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := tools.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"initializer": {"not found"},
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

	client := tama.NewClient(config)
	_, err := client.Tools.GetInitializer("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var toolsErr *tools.Error
	if errors.As(err, &toolsErr) {
		if toolsErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", toolsErr.StatusCode)
		}
		if toolsErr.Errors == nil || len(toolsErr.Errors["initializer"]) == 0 ||
			toolsErr.Errors["initializer"][0] != "not found" {
			t.Errorf("Expected error 'initializer not found', got %v", toolsErr.Errors)
		}
	} else {
		t.Errorf("Expected tools.Error, got %T", err)
	}
}

func TestToolsCreateInitializer(t *testing.T) {
	indexValue := 0
	request := tools.CreateInitializerRequest{
		Initializer: tools.InitializerRequestData{
			Reference:  "test-reference",
			Index:      &indexValue,
			Parameters: map[string]any{"temperature": 0.7, "max_tokens": 100},
		},
	}

	expectedInitializer := tools.Initializer{
		ID:             "initializer-456",
		Reference:      "test-reference",
		Index:          0,
		Parameters:     map[string]any{"temperature": 0.7, "max_tokens": 100},
		ThoughtToolID:  "tool-123",
		ProvisionState: "pending",
	}

	expectedResponse := tools.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/tool-123/initializers" {
			t.Errorf("Expected path /provision/tools/tool-123/initializers, got %s", r.URL.Path)
		}

		var receivedRequest tools.CreateInitializerRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Initializer.Reference != request.Initializer.Reference {
			t.Errorf(
				"Expected initializer reference %s, got %s",
				request.Initializer.Reference,
				receivedRequest.Initializer.Reference,
			)
		}

		if (receivedRequest.Initializer.Index == nil && request.Initializer.Index != nil) ||
			(receivedRequest.Initializer.Index != nil && request.Initializer.Index == nil) ||
			(receivedRequest.Initializer.Index != nil && request.Initializer.Index != nil &&
				*receivedRequest.Initializer.Index != *request.Initializer.Index) {
			var receivedIndex, expectedIndex interface{}
			if receivedRequest.Initializer.Index != nil {
				receivedIndex = *receivedRequest.Initializer.Index
			}
			if request.Initializer.Index != nil {
				expectedIndex = *request.Initializer.Index
			}
			t.Errorf(
				"Expected initializer index %v, got %v",
				expectedIndex,
				receivedIndex,
			)
		}

		// Validate parameters
		if len(receivedRequest.Initializer.Parameters) != len(request.Initializer.Parameters) {
			t.Errorf(
				"Expected %d parameters, got %d",
				len(request.Initializer.Parameters),
				len(receivedRequest.Initializer.Parameters),
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	initializer, err := client.Tools.CreateInitializer("tool-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestToolsCreateInitializerValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	indexValue := 0
	// Test empty thought tool ID
	_, err := client.Tools.CreateInitializer("", tools.CreateInitializerRequest{
		Initializer: tools.InitializerRequestData{
			Reference: "test-ref",
			Index:     &indexValue,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought tool ID")
	}

	// Test empty initializer reference
	_, err = client.Tools.CreateInitializer("tool-123", tools.CreateInitializerRequest{
		Initializer: tools.InitializerRequestData{
			Reference: "",
			Index:     &indexValue,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty initializer reference")
	}
}

func TestToolsCreateInitializerWithNilParameters(t *testing.T) {
	indexValue := 0
	request := tools.CreateInitializerRequest{
		Initializer: tools.InitializerRequestData{
			Reference:  "test-reference",
			Index:      &indexValue,
			Parameters: nil, // Will be initialized to empty map
		},
	}

	expectedInitializer := tools.Initializer{
		ID:             "initializer-456",
		Reference:      "test-reference",
		Index:          0,
		Parameters:     map[string]any{},
		ThoughtToolID:  "tool-123",
		ProvisionState: "pending",
	}

	expectedResponse := tools.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		var receivedRequest tools.CreateInitializerRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Due to JSON omitempty tag, nil parameters map will be omitted from JSON
		// and when unmarshaled it will remain nil, which is expected behavior
		if len(receivedRequest.Initializer.Parameters) != 0 {
			t.Errorf("Expected empty parameters map, got %d items", len(receivedRequest.Initializer.Parameters))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	initializer, err := client.Tools.CreateInitializer("tool-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestToolsCreateInitializerFieldValidationDelegated(t *testing.T) {
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
				"reference": {"is required", "must be valid"},
				"index":     {"can't be blank"},
				"base":      {"thought tool is not active"},
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

	client := tama.NewClient(config)
	invalidIndex := -1 // Use an invalid int value instead of string
	_, err := client.Tools.CreateInitializer("tool-123", tools.CreateInitializerRequest{
		Initializer: tools.InitializerRequestData{
			Reference: "invalid",
			Index:     &invalidIndex,
		},
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var toolsErr *tools.Error
	if errors.As(err, &toolsErr) {
		if toolsErr.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code 422, got %d", toolsErr.StatusCode)
		}
		if len(toolsErr.Errors["reference"]) != 2 {
			t.Errorf("Expected 2 reference errors, got %d", len(toolsErr.Errors["reference"]))
		}
	} else {
		t.Errorf("Expected tools.Error, got %T", err)
	}
}

func TestToolsUpdateInitializer(t *testing.T) {
	indexValue := 1
	request := tools.UpdateInitializerRequest{
		Initializer: tools.UpdateInitializerData{
			Reference:  "updated-reference",
			Index:      &indexValue,
			Parameters: map[string]any{"temperature": 0.8},
		},
	}

	expectedInitializer := tools.Initializer{
		ID:             "initializer-123",
		Reference:      "updated-reference",
		Index:          1,
		Parameters:     map[string]any{"temperature": 0.8},
		ThoughtToolID:  "tool-456",
		ProvisionState: "active",
	}

	expectedResponse := tools.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/initializers/initializer-123" {
			t.Errorf("Expected path /provision/tools/initializers/initializer-123, got %s", r.URL.Path)
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

	client := tama.NewClient(config)
	initializer, err := client.Tools.UpdateInitializer("initializer-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestToolsReplaceInitializer(t *testing.T) {
	indexValue := 2
	request := tools.UpdateInitializerRequest{
		Initializer: tools.UpdateInitializerData{
			Reference:  "replaced-reference",
			Index:      &indexValue,
			Parameters: map[string]any{"max_tokens": 200},
		},
	}

	expectedInitializer := tools.Initializer{
		ID:             "initializer-123",
		Reference:      "replaced-reference",
		Index:          2,
		Parameters:     map[string]any{"max_tokens": 200},
		ThoughtToolID:  "tool-456",
		ProvisionState: "active",
	}

	expectedResponse := tools.InitializerResponse{
		Data: expectedInitializer,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/initializers/initializer-123" {
			t.Errorf("Expected path /provision/tools/initializers/initializer-123, got %s", r.URL.Path)
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

	client := tama.NewClient(config)
	initializer, err := client.Tools.ReplaceInitializer("initializer-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ValidateInitializerResponse(t, *initializer, expectedInitializer)
}

func TestToolsDeleteInitializer(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/initializers/initializer-123" {
			t.Errorf("Expected path /provision/tools/initializers/initializer-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	err := client.Tools.DeleteInitializer("initializer-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestToolsGetInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Tools.GetInitializer("")
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in GetInitializer")
	}
}

func TestToolsUpdateInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Tools.UpdateInitializer("", tools.UpdateInitializerRequest{
		Initializer: tools.UpdateInitializerData{Reference: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in UpdateInitializer")
	}
}

func TestToolsReplaceInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Tools.ReplaceInitializer("", tools.UpdateInitializerRequest{
		Initializer: tools.UpdateInitializerData{Reference: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in ReplaceInitializer")
	}
}

func TestToolsDeleteInitializerEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Tools.DeleteInitializer("")
	if err == nil {
		t.Error("Expected validation error for empty initializer ID in DeleteInitializer")
	}
}

func TestToolsCreateInitializerWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/tools/tool-123/initializers" {
			t.Errorf("Expected path /provision/tools/tool-123/initializers, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		errorResponse := map[string]any{
			"errors": map[string][]string{
				"reference":  {"can't be blank", "is not valid"},
				"index":      {"can't be blank", "must be numeric"},
				"parameters": {"invalid format"},
				"base":       {"thought tool is not active"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	invalidIndex := -1
	request := tools.CreateInitializerRequest{
		Initializer: tools.InitializerRequestData{
			Reference:  "invalid-reference",
			Index:      &invalidIndex,
			Parameters: map[string]any{"invalid": "value"},
		},
	}

	_, err := client.Tools.CreateInitializer("tool-123", request)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "reference can't be blank") {
		t.Errorf("Expected error to contain 'reference can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "reference is not valid") {
		t.Errorf("Expected error to contain 'reference is not valid', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "index can't be blank") {
		t.Errorf("Expected error to contain 'index can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "parameters invalid format") {
		t.Errorf("Expected error to contain 'parameters invalid format', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "base thought tool is not active") {
		t.Errorf("Expected error to contain 'base thought tool is not active', got %s", errorMsg)
	}
}
