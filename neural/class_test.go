package neural_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
)

func TestNeuralGetClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-123",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
		},
		Name:        "TestClass",
		Description: "Test class for schema validation",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
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

	class, err := client.Neural.GetClass("class-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.SpaceID != expectedClass.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expectedClass.SpaceID, class.SpaceID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralGetClassError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string][]string{
				"class": {"not found"},
			},
		})
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	_, err := client.Neural.GetClass("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiError *neural.Error
	if !errors.As(err, &apiError) {
		t.Fatalf("Expected neural.Error, got %T", err)
	}

	if apiError.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", apiError.StatusCode)
	}
}

func TestNeuralCreateClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-789",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
		},
		Name:        "NewClass",
		Description: "A new class for testing",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/spaces/space-456/classes" {
			t.Errorf("Expected path /provision/neural/spaces/space-456/classes, got %s", r.URL.Path)
		}

		var req neural.CreateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Class.Schema == nil {
			t.Error("Expected schema in request, got nil")
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

	createReq := neural.CreateClassRequest{
		Class: neural.ClassRequestData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{"type": "string"},
				},
			},
		},
	}

	class, err := client.Neural.CreateClass("space-456", createReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.SpaceID != expectedClass.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expectedClass.SpaceID, class.SpaceID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralCreateClassValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "http://localhost:8080",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	// Test missing space ID
	createReq := neural.CreateClassRequest{
		Class: neural.ClassRequestData{
			Schema: map[string]any{"type": "object"},
		},
	}

	_, err := client.Neural.CreateClass("", createReq)
	if err == nil {
		t.Error("Expected error for missing space ID, got nil")
	}

	// Test missing schema
	createReqNoSchema := neural.CreateClassRequest{
		Class: neural.ClassRequestData{},
	}

	_, err = client.Neural.CreateClass("space-123", createReqNoSchema)
	if err == nil {
		t.Error("Expected error for missing schema, got nil")
	}
}

func TestNeuralUpdateClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-123",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
				"age":  map[string]any{"type": "number"},
			},
		},
		Name:        "TestClass",
		Description: "Updated test class",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
		}

		var req neural.UpdateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
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

	updateReq := neural.UpdateClassRequest{
		Class: neural.UpdateClassData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{"type": "string"},
					"age":  map[string]any{"type": "number"},
				},
			},
		},
	}

	class, err := client.Neural.UpdateClass("class-123", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralReplaceClass(t *testing.T) {
	expectedClass := neural.Class{
		ID:             "class-123",
		SpaceID:        "space-456",
		ProvisionState: "active",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title": map[string]any{"type": "string"},
			},
		},
		Name:        "ReplacedClass",
		Description: "A completely replaced class",
	}

	expectedResponse := neural.ClassResponse{
		Data: expectedClass,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
		}

		var req neural.UpdateClassRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
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

	replaceReq := neural.UpdateClassRequest{
		Class: neural.UpdateClassData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{"type": "string"},
				},
			},
		},
	}

	class, err := client.Neural.ReplaceClass("class-123", replaceReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class.ID != expectedClass.ID {
		t.Errorf("Expected class ID %s, got %s", expectedClass.ID, class.ID)
	}

	if class.Name != expectedClass.Name {
		t.Errorf("Expected class name %s, got %s", expectedClass.Name, class.Name)
	}
}

func TestNeuralDeleteClass(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123" {
			t.Errorf("Expected path /provision/neural/classes/class-123, got %s", r.URL.Path)
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

	err := client.Neural.DeleteClass("class-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNeuralCreateClassWithRealWorldSchema(t *testing.T) {
	expectedClass := createActionCallClass()
	expectedResponse := neural.ClassResponse{Data: expectedClass}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		validateCreateClassRequest(t, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := createTestClient(server.URL)
	createReq := neural.CreateClassRequest{
		Class: neural.ClassRequestData{Schema: createActionCallSchema()},
	}

	class, err := client.Neural.CreateClass("space-123", createReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateClassResponse(t, *class, expectedClass)
}

func createActionCallClass() neural.Class {
	return neural.Class{
		ID:             "class-action-call",
		SpaceID:        "space-123",
		ProvisionState: "active",
		Schema:         createActionCallSchema(),
		Name:           "ActionCall",
		Description:    "Schema for action call requests",
	}
}

func createActionCallSchema() map[string]any {
	return map[string]any{
		"title":       "action-call",
		"description": "An action call is a request to execute an action.",
		"type":        "object",
		"properties": map[string]any{
			"code": map[string]any{
				"description": "The status of the action call",
				"type":        "integer",
			},
			"tool_id": map[string]any{
				"description": "The ID of the tool to execute",
				"type":        "string",
			},
			"parameters": map[string]any{
				"description": "The parameters to pass to the action",
				"type":        "object",
			},
			"content_type": map[string]any{
				"description": "The content type of the response",
				"type":        "string",
			},
			"content": map[string]any{
				"description": "The response from the action",
				"type":        "object",
			},
		},
		"required": []any{"tool_id", "parameters", "code", "content_type", "content"},
	}
}

func validateCreateClassRequest(t *testing.T, r *http.Request) {
	if r.Method != http.MethodPost {
		t.Errorf("Expected POST request, got %s", r.Method)
	}
	if r.URL.Path != "/provision/neural/spaces/space-123/classes" {
		t.Errorf("Expected path /provision/neural/spaces/space-123/classes, got %s", r.URL.Path)
	}

	var req neural.CreateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatalf("Failed to decode request body: %v", err)
	}

	validateSchemaStructure(t, req.Class.Schema)
}

func validateSchemaStructure(t *testing.T, schema map[string]any) {
	if schema == nil {
		t.Error("Expected schema in request, got nil")
		return
	}

	title, ok := schema["title"].(string)
	if !ok || title != "action-call" {
		t.Errorf("Expected schema title 'action-call', got %v", schema["title"])
	}

	expectedDesc := "An action call is a request to execute an action."
	desc, ok := schema["description"].(string)
	if !ok || desc != expectedDesc {
		t.Errorf("Expected specific description, got %v", schema["description"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Error("Expected properties to be a map")
		return
	}

	requiredProps := []string{"code", "tool_id", "parameters", "content_type", "content"}
	for _, prop := range requiredProps {
		if _, exists := props[prop]; !exists {
			t.Errorf("Expected property %s to exist in schema", prop)
		}
	}
}

func validateClassResponse(t *testing.T, actual, expected neural.Class) {
	if actual.ID != expected.ID {
		t.Errorf("Expected class ID %s, got %s", expected.ID, actual.ID)
	}
	if actual.SpaceID != expected.SpaceID {
		t.Errorf("Expected space ID %s, got %s", expected.SpaceID, actual.SpaceID)
	}
	if actual.Name != expected.Name {
		t.Errorf("Expected class name %s, got %s", expected.Name, actual.Name)
	}

	title, ok := actual.Schema["title"].(string)
	if !ok || title != "action-call" {
		t.Errorf("Expected schema title 'action-call', got %v", actual.Schema["title"])
	}

	expectedDesc := "An action call is a request to execute an action."
	desc, ok := actual.Schema["description"].(string)
	if !ok || desc != expectedDesc {
		t.Errorf("Expected specific description, got %v", actual.Schema["description"])
	}
}
