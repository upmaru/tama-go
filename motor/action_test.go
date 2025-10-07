package motor_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestMotorGetActionByPathAndMethod(t *testing.T) {
	actionPath := "/3/movie/{movie_id}"
	actionMethod := "GET"
	encodedPath := base64.URLEncoding.EncodeToString([]byte(actionPath))

	expectedAction := motor.Action{
		ID:              "action-456",
		Identifier:      "movie-details",
		Path:            actionPath,
		Method:          actionMethod,
		SpecificationID: "spec-789",
	}

	expectedResponse := motor.ActionResponse{
		Data: expectedAction,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedURLPath := "/provision/motor/specifications/spec-789/actions/" + encodedPath
		if r.URL.Path != expectedURLPath {
			t.Errorf("Expected path %s, got %s", expectedURLPath, r.URL.Path)
		}

		// Check for method query parameter (should be lowercase)
		methodParam := r.URL.Query().Get("method")
		expectedMethodParam := "get" // Method should be lowercased
		if methodParam != expectedMethodParam {
			t.Errorf("Expected method query parameter %s, got %s", expectedMethodParam, methodParam)
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
	action, err := client.Motor.GetActionByPathAndMethod("spec-789", actionPath, actionMethod)

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

func TestMotorGetActionByPathAndMethodError(t *testing.T) {
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
	_, err = client.Motor.GetActionByPathAndMethod("spec-789", "/nonexistent/path", "GET")

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

func TestMotorGetActionByPathAndMethodValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty specification ID
	_, err = client.Motor.GetActionByPathAndMethod("", "/some/path", "GET")
	if err == nil {
		t.Error("Expected validation error for empty specification ID")
	}
	expectedMsg := "specification ID, path, and method are required"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test empty path
	_, err = client.Motor.GetActionByPathAndMethod("spec-789", "", "GET")
	if err == nil {
		t.Error("Expected validation error for empty path")
	}
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test empty method
	_, err = client.Motor.GetActionByPathAndMethod("spec-789", "/some/path", "")
	if err == nil {
		t.Error("Expected validation error for empty method")
	}
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test all empty
	_, err = client.Motor.GetActionByPathAndMethod("", "", "")
	if err == nil {
		t.Error("Expected validation error for all empty parameters")
	}
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestMotorGetActionByPathAndMethodValidationNoHTTPCall(t *testing.T) {
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
	_, err = client.Motor.GetActionByPathAndMethod("", "/some/path", "GET")
	if err == nil {
		t.Error("Expected validation error for empty specification ID")
	}
	if serverCalled {
		t.Error("HTTP call was made despite validation error")
	}

	// Reset and test empty path - should not make HTTP call
	serverCalled = false
	_, err = client.Motor.GetActionByPathAndMethod("spec-789", "", "GET")
	if err == nil {
		t.Error("Expected validation error for empty path")
	}
	if serverCalled {
		t.Error("HTTP call was made despite validation error")
	}

	// Reset and test empty method - should not make HTTP call
	serverCalled = false
	_, err = client.Motor.GetActionByPathAndMethod("spec-789", "/some/path", "")
	if err == nil {
		t.Error("Expected validation error for empty method")
	}
	if serverCalled {
		t.Error("HTTP call was made despite validation error")
	}
}

func TestMotorGetActionByPathAndMethodEncoding(t *testing.T) {
	// Test that the path is correctly encoded and method is lowercased
	testCases := []struct {
		name         string
		path         string
		method       string
		expectedPath string
	}{
		{
			name:         "simple path with GET",
			path:         "/api/v1/test",
			method:       "GET",
			expectedPath: base64.URLEncoding.EncodeToString([]byte("/api/v1/test")),
		},
		{
			name:         "path with parameters and POST",
			path:         "/3/movie/{movie_id}",
			method:       "POST",
			expectedPath: base64.URLEncoding.EncodeToString([]byte("/3/movie/{movie_id}")),
		},
		{
			name:         "path with special characters and PUT",
			path:         "/api/users/{id}/posts?filter=active",
			method:       "PUT",
			expectedPath: base64.URLEncoding.EncodeToString([]byte("/api/users/{id}/posts?filter=active")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedAction := motor.Action{
				ID:              "test-action",
				Identifier:      "test-identifier",
				Path:            tc.path,
				Method:          tc.method,
				SpecificationID: "test-spec",
			}

			expectedResponse := motor.ActionResponse{
				Data: expectedAction,
			}

			server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
				expectedURLPath := "/provision/motor/specifications/test-spec/actions/" + tc.expectedPath
				if r.URL.Path != expectedURLPath {
					t.Errorf("Expected path %s, got %s", expectedURLPath, r.URL.Path)
				}

				// Check for method query parameter (should be lowercase)
				methodParam := r.URL.Query().Get("method")
				expectedMethodParam := strings.ToLower(tc.method)
				if methodParam != expectedMethodParam {
					t.Errorf("Expected method query parameter %s, got %s", expectedMethodParam, methodParam)
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
			action, err := client.Motor.GetActionByPathAndMethod("test-spec", tc.path, tc.method)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if action.Path != tc.path {
				t.Errorf("Expected action path %s, got %s", tc.path, action.Path)
			}

			if action.Method != tc.method {
				t.Errorf("Expected action method %s, got %s", tc.method, action.Method)
			}
		})
	}
}
