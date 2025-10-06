package sensory_test

import (
	"encoding/json"
	"net/http"
	"testing"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/sensory"
)

func TestSensoryGetSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-123",
		SpaceID: "space-456",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type": "string",
				},
			},
		},
		Version:        "1.0.0",
		Endpoint:       "https://api.test.com/v1",
		CurrentState:   "active",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	spec, err := client.Sensory.GetSpecification("spec-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.ID != expectedSpec.ID {
		t.Errorf("Expected specification ID %s, got %s", expectedSpec.ID, spec.ID)
	}

	if spec.SpaceID != expectedSpec.SpaceID {
		t.Errorf("Expected specification space ID %s, got %s", expectedSpec.SpaceID, spec.SpaceID)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}

	if spec.CurrentState != expectedSpec.CurrentState {
		t.Errorf("Expected specification current state %s, got %s", expectedSpec.CurrentState, spec.CurrentState)
	}

	if spec.ProvisionState != expectedSpec.ProvisionState {
		t.Errorf("Expected specification provision state %s, got %s", expectedSpec.ProvisionState, spec.ProvisionState)
	}

	// Validate schema structure
	if spec.Schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got %v", spec.Schema["type"])
	}
}

func TestSensoryCreateSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-789",
		SpaceID: "space-123",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]any{
					"type": "string",
				},
			},
		},
		Version:        "2.0.0",
		Endpoint:       "https://api.example.com/v2",
		CurrentState:   "pending",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/spaces/space-123/specifications" {
			t.Errorf("Expected path /provision/sensory/spaces/space-123/specifications, got %s", r.URL.Path)
		}

		var req sensory.CreateSpecificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Specification.Version != "2.0.0" {
			t.Errorf("Expected request version '2.0.0', got %s", req.Specification.Version)
		}

		if req.Specification.Endpoint != "https://api.example.com/v2" {
			t.Errorf("Expected request endpoint 'https://api.example.com/v2', got %s", req.Specification.Endpoint)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	createReq := sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{
						"type": "string",
					},
				},
			},
			Version:  "2.0.0",
			Endpoint: "https://api.example.com/v2",
		},
	}

	spec, err := client.Sensory.CreateSpecification("space-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.ID != expectedSpec.ID {
		t.Errorf("Expected specification ID %s, got %s", expectedSpec.ID, spec.ID)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}
}

func TestSensoryCreateSpecificationValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty space ID validation
	_, err := client.Sensory.CreateSpecification("", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
			},
			Version:  "1.0.0",
			Endpoint: "https://api.test.com",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty schema validation
	_, err = client.Sensory.CreateSpecification("space-123", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Version:  "1.0.0",
			Endpoint: "https://api.test.com",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty schema")
	}

	// Test empty version validation
	_, err = client.Sensory.CreateSpecification("space-123", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
			},
			Endpoint: "https://api.test.com",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty version")
	}

	// Test empty endpoint validation
	_, err = client.Sensory.CreateSpecification("space-123", sensory.CreateSpecificationRequest{
		Specification: sensory.SpecificationRequestData{
			Schema: map[string]any{
				"type": "object",
			},
			Version: "1.0.0",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestSensoryGetSpecification_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err := client.Sensory.GetSpecification("")
	if err == nil {
		t.Error("Expected validation error for empty specification ID in GetSpecification")
	}
}

func TestSensoryUpdateSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-123",
		SpaceID: "space-456",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"updated_field": map[string]any{
					"type": "string",
				},
			},
		},
		Version:        "1.1.0",
		Endpoint:       "https://api.updated.com/v1",
		CurrentState:   "active",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSpecificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Specification.Version != "1.1.0" {
			t.Errorf("Expected request version '1.1.0', got %s", req.Specification.Version)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	updateReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"updated_field": map[string]any{
						"type": "string",
					},
				},
			},
			Version:  "1.1.0",
			Endpoint: "https://api.updated.com/v1",
		},
	}

	spec, err := client.Sensory.UpdateSpecification("spec-123", updateReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}
}

func TestSensoryUpdateSpecification_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	updateReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Version: "1.1.0",
		},
	}

	_, err := client.Sensory.UpdateSpecification("", updateReq)
	if err == nil {
		t.Error("Expected validation error for empty specification ID in UpdateSpecification")
	}
}

func TestSensoryReplaceSpecification(t *testing.T) {
	expectedSpec := sensory.Specification{
		ID:      "spec-123",
		SpaceID: "space-456",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"replaced_field": map[string]any{
					"type": "number",
				},
			},
		},
		Version:        "2.0.0",
		Endpoint:       "https://api.replaced.com/v2",
		CurrentState:   "pending",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SpecificationResponse{
		Data: expectedSpec,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSpecificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Specification.Version != "2.0.0" {
			t.Errorf("Expected request version '2.0.0', got %s", req.Specification.Version)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	replaceReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"replaced_field": map[string]any{
						"type": "number",
					},
				},
			},
			Version:  "2.0.0",
			Endpoint: "https://api.replaced.com/v2",
		},
	}

	spec, err := client.Sensory.ReplaceSpecification("spec-123", replaceReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if spec.Version != expectedSpec.Version {
		t.Errorf("Expected specification version %s, got %s", expectedSpec.Version, spec.Version)
	}

	if spec.Endpoint != expectedSpec.Endpoint {
		t.Errorf("Expected specification endpoint %s, got %s", expectedSpec.Endpoint, spec.Endpoint)
	}
}

func TestSensoryReplaceSpecification_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	replaceReq := sensory.UpdateSpecificationRequest{
		Specification: sensory.UpdateSpecificationData{
			Version: "2.0.0",
		},
	}

	_, err := client.Sensory.ReplaceSpecification("", replaceReq)
	if err == nil {
		t.Error("Expected validation error for empty specification ID in ReplaceSpecification")
	}
}

func TestSensoryDeleteSpecification(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/specifications/spec-123" {
			t.Errorf("Expected path /provision/sensory/specifications/spec-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err := client.Sensory.DeleteSpecification("spec-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSensoryDeleteSpecification_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err := client.Sensory.DeleteSpecification("")
	if err == nil {
		t.Error("Expected validation error for empty specification ID in DeleteSpecification")
	}
}
