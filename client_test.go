package tama_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
)

// createMockServer creates a test HTTP server with the given handler.
func createMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
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

func TestSetAPIKey(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "original-key",
	})

	newAPIKey := "new-api-key"
	client.SetAPIKey(newAPIKey)

	// API key is set successfully - internal field is not accessible from external package
}

func TestSetDebug(t *testing.T) {
	client := tama.NewClient(tama.Config{
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
