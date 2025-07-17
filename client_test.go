package tama_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
)

// createMockServer creates a test HTTP server with the given handler.
func createMockServer(_ *testing.T, handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestNewClient(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	// Client is created successfully - internal fields are not accessible from external package

	if client.Neural == nil {
		t.Error("Expected Neural service to be initialized")
	}

	if client.Sensory == nil {
		t.Error("Expected Sensory service to be initialized")
	}
}

func TestNewClientDefaultTimeout(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		// No timeout specified
	}

	client := tama.NewClient(config)

	// We can't directly test the timeout, but we can ensure the client was created
	if client == nil {
		t.Fatal("Expected client to be created with default timeout")
	}
}

func TestSetAPIKey(_ *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "original-key",
	})

	newAPIKey := "new-api-key"
	client.SetAPIKey(newAPIKey)

	// API key is set successfully - internal field is not accessible from external package
}

func TestSetDebug(_ *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	client.SetDebug(true)
	// Note: In a real implementation, you might want to verify that debug mode is actually set
	// This would depend on how debug mode is implemented in the resty client
}

func TestErrorStruct(t *testing.T) {
	// Test neural Error with only status code
	neuralErr := &neural.Error{
		StatusCode: 404,
	}

	expectedErrorMsg := "API error 404"
	if neuralErr.Error() != expectedErrorMsg {
		t.Errorf("Expected error message %s, got %s", expectedErrorMsg, neuralErr.Error())
	}

	// Test error without status code
	neuralErrNoStatus := &neural.Error{}

	expectedErrorMsgNoStatus := "API error"
	if neuralErrNoStatus.Error() != expectedErrorMsgNoStatus {
		t.Errorf("Expected error message %s, got %s", expectedErrorMsgNoStatus, neuralErrNoStatus.Error())
	}

	// Test field-specific errors
	fieldErr := &neural.Error{
		StatusCode: 422,
		Errors: map[string][]string{
			"source_id": {"has already been taken"},
			"name":      {"is required", "must be at least 3 characters"},
		},
	}

	errorMsg := fieldErr.Error()
	// Check that all field errors are included
	if !contains(errorMsg, "source_id has already been taken") {
		t.Errorf("Expected error message to contain 'source_id has already been taken', got %s", errorMsg)
	}
	if !contains(errorMsg, "name is required") {
		t.Errorf("Expected error message to contain 'name is required', got %s", errorMsg)
	}
	if !contains(errorMsg, "name must be at least 3 characters") {
		t.Errorf("Expected error message to contain 'name must be at least 3 characters', got %s", errorMsg)
	}
	if !contains(errorMsg, "API error 422:") {
		t.Errorf("Expected error message to contain status code, got %s", errorMsg)
	}

	// Test field-specific errors without status code
	fieldErrNoStatus := &neural.Error{
		Errors: map[string][]string{
			"email": {"is invalid"},
		},
	}

	errorMsgNoStatus := fieldErrNoStatus.Error()
	expectedNoStatus := "API error: email is invalid"
	if errorMsgNoStatus != expectedNoStatus {
		t.Errorf("Expected error message %s, got %s", expectedNoStatus, errorMsgNoStatus)
	}
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestEmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
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
}
