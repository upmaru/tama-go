package sensory_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/sensory"
)

func TestSensoryGetSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-123",
		Name:           "Test Source",
		Endpoint:       "https://api.test.com/v1",
		SpaceID:        "space-456",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123" {
			t.Errorf("Expected path /provision/sensory/sources/source-123, got %s", r.URL.Path)
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

	source, err := client.Sensory.GetSource("source-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryGetSourceBySpecificationAndSlug(t *testing.T) {
	specID := "spec-123"
	slug := "slugged-source"

	expectedSource := sensory.Source{
		ID:             "source-xyz",
		Name:           "Slugged Source",
		Slug:           slug,
		Endpoint:       "https://api.slugged.com/v1",
		SpaceID:        "space-999",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SourceResponse{Data: expectedSource}

	expectedPath := "/provision/sensory/specifications/" + specID + "/sources/" + slug

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	source, err := client.Sensory.GetSourceBySpecificationAndSlug(specID, slug)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}
	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}
	if source.Slug != expectedSource.Slug {
		t.Errorf("Expected source slug %s, got %s", expectedSource.Slug, source.Slug)
	}
	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}
	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}
}

func TestSensoryCreateSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-789",
		Name:           "New Source",
		Endpoint:       "https://api.mistral.ai/v1",
		SpaceID:        "space-123",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/spaces/space-123/sources" {
			t.Errorf("Expected path /provision/sensory/spaces/space-123/sources, got %s", r.URL.Path)
		}

		var req sensory.CreateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Source.Name != "New Source" {
			t.Errorf("Expected request name 'New Source', got %s", req.Source.Name)
		}

		if req.Source.Type != "model" {
			t.Errorf("Expected request type 'model', got %s", req.Source.Type)
		}

		if req.Source.Endpoint != "https://api.mistral.ai/v1" {
			t.Errorf("Expected request endpoint 'https://api.mistral.ai/v1', got %s", req.Source.Endpoint)
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

	createReq := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "New Source",
			Type:     "model",
			Endpoint: "https://api.mistral.ai/v1",
			Credential: sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	}

	source, err := client.Sensory.CreateSource("space-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryCreateSourceValidation(t *testing.T) {
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
	_, err = client.Sensory.CreateSource("", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "Test",
			Type:     "model",
			Endpoint: "https://api.test.com",
			Credential: sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty space ID")
	}

	// Test empty name validation
	_, err = client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Type:     "model",
			Endpoint: "https://api.test.com",
			Credential: sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	// Test empty type validation
	_, err = client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "Test",
			Endpoint: "https://api.test.com",
			Credential: sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty type")
	}

	// Test empty endpoint validation
	_, err = client.Sensory.CreateSource("space-123", sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name: "Test",
			Type: "model",
			Credential: sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty endpoint")
	}
}

func TestSensoryGetSource_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

_, err = client.Sensory.GetSource("")
	if err == nil {
		t.Error("Expected validation error for empty source ID in GetSource")
	}
}

func TestSensoryUpdateSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-123",
		Name:           "Updated Source",
		Endpoint:       "https://api.updated.com/v1",
		SpaceID:        "space-456",
		ProvisionState: "active",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123" {
			t.Errorf("Expected path /provision/sensory/sources/source-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Source.Name != "Updated Source" {
			t.Errorf("Expected request name 'Updated Source', got %s", req.Source.Name)
		}

		if req.Source.Endpoint != "https://api.updated.com/v1" {
			t.Errorf("Expected request endpoint 'https://api.updated.com/v1', got %s", req.Source.Endpoint)
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

	updateReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name:     "Updated Source",
			Endpoint: "https://api.updated.com/v1",
			Credential: &sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	}

	source, err := client.Sensory.UpdateSource("source-123", updateReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryUpdateSource_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	updateReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name: "Updated Source",
		},
	}

	_, err = client.Sensory.UpdateSource("", updateReq)
	if err == nil {
		t.Error("Expected validation error for empty source ID in UpdateSource")
	}
}

func TestSensoryReplaceSource(t *testing.T) {
	expectedSource := sensory.Source{
		ID:             "source-123",
		Name:           "Replaced Source",
		Endpoint:       "https://api.replaced.com/v1",
		SpaceID:        "space-456",
		ProvisionState: "pending",
	}

	expectedResponse := sensory.SourceResponse{
		Data: expectedSource,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123" {
			t.Errorf("Expected path /provision/sensory/sources/source-123, got %s", r.URL.Path)
		}

		var req sensory.UpdateSourceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Source.Name != "Replaced Source" {
			t.Errorf("Expected request name 'Replaced Source', got %s", req.Source.Name)
		}

		if req.Source.Endpoint != "https://api.replaced.com/v1" {
			t.Errorf("Expected request endpoint 'https://api.replaced.com/v1', got %s", req.Source.Endpoint)
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

	replaceReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name:     "Replaced Source",
			Type:     "model",
			Endpoint: "https://api.replaced.com/v1",
			Credential: &sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	}

	source, err := client.Sensory.ReplaceSource("source-123", replaceReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if source.ID != expectedSource.ID {
		t.Errorf("Expected source ID %s, got %s", expectedSource.ID, source.ID)
	}

	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}

	if source.Endpoint != expectedSource.Endpoint {
		t.Errorf("Expected source endpoint %s, got %s", expectedSource.Endpoint, source.Endpoint)
	}

	if source.ProvisionState != expectedSource.ProvisionState {
		t.Errorf("Expected source provision state %s, got %s", expectedSource.ProvisionState, source.ProvisionState)
	}

	if source.SpaceID != expectedSource.SpaceID {
		t.Errorf("Expected source space ID %s, got %s", expectedSource.SpaceID, source.SpaceID)
	}
}

func TestSensoryReplaceSource_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	replaceReq := sensory.UpdateSourceRequest{
		Source: sensory.UpdateSourceData{
			Name: "Replaced Source",
		},
	}

	_, err = client.Sensory.ReplaceSource("", replaceReq)
	if err == nil {
		t.Error("Expected validation error for empty source ID in ReplaceSource")
	}
}

func TestSensoryCreateSourceWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/spaces/space-123/sources" {
			t.Errorf("Expected path /provision/sensory/spaces/space-123/sources, got %s", r.URL.Path)
		}

		// Return field validation errors
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := map[string]any{
			"errors": map[string][]string{
				"name":     {"is required"},
				"endpoint": {"is invalid URL", "must use HTTPS"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
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

	createReq := sensory.CreateSourceRequest{
		Source: sensory.SourceRequestData{
			Name:     "test-source", // Valid name to bypass client validation
			Type:     "ollama",
			Endpoint: "https://valid-endpoint.com", // Valid endpoint to bypass client validation
			Credential: sensory.SourceCredential{
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				SkipTokenFetch: true,
			},
		},
	}

_, err = client.Sensory.CreateSource("space-123", createReq)
	if err == nil {
		t.Fatal("Expected error for invalid source data")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "name is required") {
		t.Errorf("Expected error to contain 'name is required', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "endpoint is invalid URL") {
		t.Errorf("Expected error to contain 'endpoint is invalid URL', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "endpoint must use HTTPS") {
		t.Errorf("Expected error to contain 'endpoint must use HTTPS', got %s", errorMsg)
	}
}
