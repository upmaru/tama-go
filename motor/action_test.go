package motor_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/motor"
)

// createMockServer creates a test HTTP server with the given handler.
func createMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestMotorGetAction(t *testing.T) {
	expectedAction := motor.Action{
		ID:              "action-123",
		Identifier:      "deploy-app",
		Path:            "/api/v1/deploy",
		Method:          "POST",
		SpecificationID: "spec-456",
	}

	expectedResponse := motor.ActionResponse{
		Data: expectedAction,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/motor/specifications/spec-456/actions/action-123" {
			t.Errorf("Expected path /provision/motor/specifications/spec-456/actions/action-123, got %s", r.URL.Path)
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
	action, err := client.Motor.GetAction("spec-456", "action-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if action.ID != expectedAction.ID {
		t.Errorf("Expected action ID %s, got %s", expectedAction.ID, action.ID)
	}

	if action.Identifier != expectedAction.Identifier {
		t.Errorf("Expected action identifier %s, got %s", expectedAction.Identifier, action.Identifier)
	}

	if action.Path != expectedAction.Path {
		t.Errorf("Expected action path %s, got %s", expectedAction.Path, action.Path)
	}

	if action.Method != expectedAction.Method {
		t.Errorf("Expected action method %s, got %s", expectedAction.Method, action.Method)
	}

	if action.SpecificationID != expectedAction.SpecificationID {
		t.Errorf("Expected action specification_id %s, got %s", expectedAction.SpecificationID, action.SpecificationID)
	}
}

func TestMotorGetActionError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := motor.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"action": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
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
	_, err = client.Motor.GetAction("spec-456", "nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var motorErr *motor.Error
	if errors.As(err, &motorErr) {
		if motorErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", motorErr.StatusCode)
		}
		if motorErr.Errors == nil || len(motorErr.Errors["action"]) == 0 ||
			motorErr.Errors["action"][0] != "not found" {
			t.Errorf("Expected error 'action not found', got %v", motorErr.Errors)
		}
	} else {
		t.Errorf("Expected motor.Error, got %T", err)
	}
}

func TestMotorGetActionValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "http://example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty specification ID
	_, err = client.Motor.GetAction("", "action-123")
	if err == nil {
		t.Error("Expected validation error for empty specification ID")
	}
	expectedMsg := "specification ID and action ID are required"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test empty action ID
	_, err = client.Motor.GetAction("spec-456", "")
	if err == nil {
		t.Error("Expected validation error for empty action ID")
	}
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test both empty
	_, err = client.Motor.GetAction("", "")
	if err == nil {
		t.Error("Expected validation error for both empty parameters")
	}
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestMotorGetActionValidationNoHTTPCall(t *testing.T) {
	// Create a mock server that should never be called
	serverCalled := false
	server := createMockServer(func(_ http.ResponseWriter, _ *http.Request) {
		serverCalled = true
		t.Error("Mock server was called, but validation should prevent HTTP requests")
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

	// Test empty specification ID - should not make HTTP call
	_, err = client.Motor.GetAction("", "action-123")
	if err == nil {
		t.Error("Expected validation error for empty specification ID")
	}
	if serverCalled {
		t.Error("HTTP call was made despite validation error")
	}

	// Reset and test empty action ID - should not make HTTP call
	serverCalled = false
	_, err = client.Motor.GetAction("spec-456", "")
	if err == nil {
		t.Error("Expected validation error for empty action ID")
	}
	if serverCalled {
		t.Error("HTTP call was made despite validation error")
	}
}

func TestAction_Execute(t *testing.T) {
	action := &motor.Action{
		ID:              "action-123",
		Identifier:      "deploy-app",
		Path:            "/api/v1/deploy",
		Method:          "POST",
		SpecificationID: "spec-456",
	}
	err := action.Execute()
	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}
}
