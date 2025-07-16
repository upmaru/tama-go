package tama

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/upmaru/tama-go/neural"
)

// createMockServer creates a test HTTP server with the given handler
func createMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

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
		APIKey:  "original-key",
	})

	newAPIKey := "new-api-key"
	client.SetAPIKey(newAPIKey)

	if client.apiKey != newAPIKey {
		t.Errorf("Expected API key %s, got %s", newAPIKey, client.apiKey)
	}
}

func TestSetDebug(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	client.SetDebug(true)
	// Note: In a real implementation, you might want to verify that debug mode is actually set
	// This would depend on how debug mode is implemented in the resty client
}

func TestErrorStruct(t *testing.T) {
	// Test neural Error
	neuralErr := &neural.Error{
		StatusCode: 404,
		Message:    "Not found",
		Details:    "Resource does not exist",
	}

	expectedErrorMsg := "API error 404: Not found - Resource does not exist"
	if neuralErr.Error() != expectedErrorMsg {
		t.Errorf("Expected error message %s, got %s", expectedErrorMsg, neuralErr.Error())
	}

	// Test error without details
	neuralErrNoDetails := &neural.Error{
		StatusCode: 500,
		Message:    "Internal server error",
	}

	expectedErrorMsgNoDetails := "API error 500: Internal server error"
	if neuralErrNoDetails.Error() != expectedErrorMsgNoDetails {
		t.Errorf("Expected error message %s, got %s", expectedErrorMsgNoDetails, neuralErrNoDetails.Error())
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
}
