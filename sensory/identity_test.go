package sensory_test

import (
	"encoding/json"
	"net/http"
	"testing"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/sensory"
)

func TestSensoryGetIdentity(t *testing.T) {
	expectedIdentity := sensory.Identity{
		ID:              "identity-123",
		SpecificationID: "spec-456",
		ProvisionState:  "active",
		CurrentState:    "running",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/health",
			Method: "GET",
			Codes:  []int{200},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
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

	identity, err := client.Sensory.GetIdentity("identity-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.SpecificationID != expectedIdentity.SpecificationID {
		t.Errorf("Expected specification ID %s, got %s", expectedIdentity.SpecificationID, identity.SpecificationID)
	}

	if identity.ProvisionState != expectedIdentity.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expectedIdentity.ProvisionState, identity.ProvisionState)
	}

	if identity.CurrentState != expectedIdentity.CurrentState {
		t.Errorf("Expected current state %s, got %s", expectedIdentity.CurrentState, identity.CurrentState)
	}

	if identity.Identifier != expectedIdentity.Identifier {
		t.Errorf("Expected identifier %s, got %s", expectedIdentity.Identifier, identity.Identifier)
	}

	if identity.Validation.Path != expectedIdentity.Validation.Path {
		t.Errorf("Expected validation path %s, got %s", expectedIdentity.Validation.Path, identity.Validation.Path)
	}

	if identity.Validation.Method != expectedIdentity.Validation.Method {
		t.Errorf("Expected validation method %s, got %s",
			expectedIdentity.Validation.Method, identity.Validation.Method)
	}

	if len(identity.Validation.Codes) != len(expectedIdentity.Validation.Codes) {
		t.Errorf("Expected %d validation codes, got %d",
			len(expectedIdentity.Validation.Codes), len(identity.Validation.Codes))
	}

	if identity.Validation.Codes[0] != expectedIdentity.Validation.Codes[0] {
		t.Errorf("Expected validation code %d, got %d",
			expectedIdentity.Validation.Codes[0], identity.Validation.Codes[0])
	}
}

func TestSensoryCreateIdentity(t *testing.T) {
	requestData := sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-api-key",
			Validation: sensory.Validation{
				Path:   "/health",
				Method: "GET",
				Codes:  []int{200},
			},
		},
	}

	expectedIdentity := sensory.Identity{
		ID:              "identity-456",
		SpecificationID: "spec-123",
		ProvisionState:  "pending",
		CurrentState:    "initializing",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/health",
			Method: "GET",
			Codes:  []int{200},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		expectedPath := "/provision/sensory/specifications/spec-123/identifiers/test-identifier/identities"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var receivedRequest sensory.CreateIdentityRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Identity.APIKey != requestData.Identity.APIKey {
			t.Errorf("Expected API key %s, got %s", requestData.Identity.APIKey, receivedRequest.Identity.APIKey)
		}

		if receivedRequest.Identity.Validation.Path != requestData.Identity.Validation.Path {
			t.Errorf("Expected validation path %s, got %s",
				requestData.Identity.Validation.Path, receivedRequest.Identity.Validation.Path)
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

	identity, err := client.Sensory.CreateIdentity("spec-123", "test-identifier", requestData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.SpecificationID != expectedIdentity.SpecificationID {
		t.Errorf("Expected specification ID %s, got %s", expectedIdentity.SpecificationID, identity.SpecificationID)
	}

	if identity.ProvisionState != expectedIdentity.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expectedIdentity.ProvisionState, identity.ProvisionState)
	}
}

func TestSensoryCreateIdentityValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "http://localhost:8080",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Test empty specification ID
	_, err = client.Sensory.CreateIdentity("", "identifier", sensory.CreateIdentityRequest{})
	if err == nil || err.Error() != "specification ID is required" {
		t.Errorf("Expected 'specification ID is required' error, got: %v", err)
	}

	// Test empty identifier
	_, err = client.Sensory.CreateIdentity("spec-123", "", sensory.CreateIdentityRequest{})
	if err == nil || err.Error() != "identifier is required" {
		t.Errorf("Expected 'identifier is required' error, got: %v", err)
	}

	// Test missing authentication (no API key or client credentials)
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{},
	})
	if err == nil || err.Error() != "either API key or client credentials (client_id and client_secret) are required" {
		t.Errorf(
			"Expected 'either API key or client credentials (client_id and client_secret) are required' error, got: %v",
			err)
	}

	// Test providing both API key and client credentials (should fail)
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey:       "test-api-key",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	})
	if err == nil || err.Error() != "provide either API key or client credentials, not both" {
		t.Errorf("Expected 'provide either API key or client credentials, not both' error, got: %v", err)
	}

	// Test incomplete client credentials (only client_id)
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			ClientID: "test-client-id",
		},
	})
	if err == nil || err.Error() != "either API key or client credentials (client_id and client_secret) are required" {
		t.Errorf(
			"Expected 'either API key or client credentials (client_id and client_secret) are required' error, got: %v",
			err)
	}

	// Test incomplete client credentials (only client_secret)
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			ClientSecret: "test-client-secret",
		},
	})
	if err == nil || err.Error() != "either API key or client credentials (client_id and client_secret) are required" {
		t.Errorf(
			"Expected 'either API key or client credentials (client_id and client_secret) are required' error, got: %v",
			err)
	}

	// Test empty validation path
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-key",
			Validation: sensory.Validation{
				Method: "GET",
				Codes:  []int{200},
			},
		},
	})
	if err == nil || err.Error() != "validation path is required" {
		t.Errorf("Expected 'validation path is required' error, got: %v", err)
	}

	// Test empty validation method
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-key",
			Validation: sensory.Validation{
				Path:  "/health",
				Codes: []int{200},
			},
		},
	})
	if err == nil || err.Error() != "validation method is required" {
		t.Errorf("Expected 'validation method is required' error, got: %v", err)
	}

	// Test empty validation codes
	_, err = client.Sensory.CreateIdentity("spec-123", "identifier", sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			APIKey: "test-key",
			Validation: sensory.Validation{
				Path:   "/health",
				Method: "GET",
			},
		},
	})
	if err == nil || err.Error() != "validation codes are required" {
		t.Errorf("Expected 'validation codes are required' error, got: %v", err)
	}
}

func TestSensoryGetIdentity_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "http://localhost:8080",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Sensory.GetIdentity("")
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}

func TestSensoryUpdateIdentity(t *testing.T) {
	requestData := sensory.UpdateIdentityRequest{
		Identity: sensory.UpdateIdentityData{
			APIKey: "updated-api-key",
			Validation: &sensory.Validation{
				Path:   "/health/check",
				Method: "POST",
				Codes:  []int{200, 201},
			},
		},
	}

	expectedIdentity := sensory.Identity{
		ID:              "identity-123",
		SpecificationID: "spec-456",
		ProvisionState:  "active",
		CurrentState:    "running",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/health/check",
			Method: "POST",
			Codes:  []int{200, 201},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
		}

		var receivedRequest sensory.UpdateIdentityRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Identity.APIKey != requestData.Identity.APIKey {
			t.Errorf("Expected API key %s, got %s", requestData.Identity.APIKey, receivedRequest.Identity.APIKey)
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

	identity, err := client.Sensory.UpdateIdentity("identity-123", requestData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.Validation.Path != expectedIdentity.Validation.Path {
		t.Errorf("Expected validation path %s, got %s", expectedIdentity.Validation.Path, identity.Validation.Path)
	}

	if identity.Validation.Method != expectedIdentity.Validation.Method {
		t.Errorf("Expected validation method %s, got %s",
			expectedIdentity.Validation.Method, identity.Validation.Method)
	}
}

func TestSensoryUpdateIdentity_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "http://localhost:8080",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Sensory.UpdateIdentity("", sensory.UpdateIdentityRequest{})
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}

func TestSensoryReplaceIdentity(t *testing.T) {
	requestData := sensory.UpdateIdentityRequest{
		Identity: sensory.UpdateIdentityData{
			APIKey: "replaced-api-key",
			Validation: &sensory.Validation{
				Path:   "/status",
				Method: "GET",
				Codes:  []int{200, 204},
			},
		},
	}

	expectedIdentity := sensory.Identity{
		ID:              "identity-123",
		SpecificationID: "spec-456",
		ProvisionState:  "active",
		CurrentState:    "running",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/status",
			Method: "GET",
			Codes:  []int{200, 204},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
		}

		var receivedRequest sensory.UpdateIdentityRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Identity.APIKey != requestData.Identity.APIKey {
			t.Errorf("Expected API key %s, got %s", requestData.Identity.APIKey, receivedRequest.Identity.APIKey)
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

	identity, err := client.Sensory.ReplaceIdentity("identity-123", requestData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.Validation.Path != expectedIdentity.Validation.Path {
		t.Errorf("Expected validation path %s, got %s", expectedIdentity.Validation.Path, identity.Validation.Path)
	}
}

func TestSensoryReplaceIdentity_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "http://localhost:8080",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.Sensory.ReplaceIdentity("", sensory.UpdateIdentityRequest{})
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}

func TestSensoryDeleteIdentity(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/identities/identity-123" {
			t.Errorf("Expected path /provision/sensory/identities/identity-123, got %s", r.URL.Path)
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

	err = client.Sensory.DeleteIdentity("identity-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSensoryDeleteIdentity_EmptyIDValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "http://localhost:8080",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.Sensory.DeleteIdentity("")
	if err == nil || err.Error() != "identity ID is required" {
		t.Errorf("Expected 'identity ID is required' error, got: %v", err)
	}
}

func TestSensoryCreateIdentityWithClientCredentials(t *testing.T) {
	requestData := sensory.CreateIdentityRequest{
		Identity: sensory.IdentityRequestData{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Validation: sensory.Validation{
				Path:   "/health",
				Method: "GET",
				Codes:  []int{200},
			},
		},
	}

	expectedIdentity := sensory.Identity{
		ID:              "identity-789",
		SpecificationID: "spec-123",
		ProvisionState:  "pending",
		CurrentState:    "initializing",
		Identifier:      "test-identifier",
		Validation: sensory.Validation{
			Path:   "/health",
			Method: "GET",
			Codes:  []int{200},
		},
	}

	expectedResponse := sensory.IdentityResponse{
		Data: expectedIdentity,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		expectedPath := "/provision/sensory/specifications/spec-123/identifiers/test-identifier/identities"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var receivedRequest sensory.CreateIdentityRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Identity.ClientID != requestData.Identity.ClientID {
			t.Errorf("Expected client ID %s, got %s", requestData.Identity.ClientID, receivedRequest.Identity.ClientID)
		}

		if receivedRequest.Identity.ClientSecret != requestData.Identity.ClientSecret {
			t.Errorf(
				"Expected client secret %s, got %s",
				requestData.Identity.ClientSecret,
				receivedRequest.Identity.ClientSecret)
		}

		if receivedRequest.Identity.APIKey != "" {
			t.Errorf("Expected empty API key, got %s", receivedRequest.Identity.APIKey)
		}

		if receivedRequest.Identity.Validation.Path != requestData.Identity.Validation.Path {
			t.Errorf("Expected validation path %s, got %s",
				requestData.Identity.Validation.Path, receivedRequest.Identity.Validation.Path)
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

	identity, err := client.Sensory.CreateIdentity("spec-123", "test-identifier", requestData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if identity.ID != expectedIdentity.ID {
		t.Errorf("Expected identity ID %s, got %s", expectedIdentity.ID, identity.ID)
	}

	if identity.SpecificationID != expectedIdentity.SpecificationID {
		t.Errorf("Expected specification ID %s, got %s", expectedIdentity.SpecificationID, identity.SpecificationID)
	}

	if identity.ProvisionState != expectedIdentity.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expectedIdentity.ProvisionState, identity.ProvisionState)
	}
}
