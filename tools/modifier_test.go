package tools_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/tools"
)

func TestToolsGetModifier(t *testing.T) {
	expected := modifierFixture()
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/modifiers/modifier-123" {
			t.Errorf("Expected modifier path, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	modifier, err := client.Tools.GetModifier("modifier-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateModifierResponse(t, *modifier, expected)
}

func TestToolsCreateModifier(t *testing.T) {
	expected := modifierFixture()
	request := validCreateModifierRequest()

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/tool-456/modifiers" {
			t.Errorf("Expected thought-tool modifier path, got %s", r.URL.Path)
		}

		modifier := decodeModifierEnvelope(t, r)
		if len(modifier) != 5 {
			t.Errorf("Expected 5 modifier fields, got %d: %v", len(modifier), modifier)
		}
		index, present := modifier["index"]
		if !present {
			t.Fatal("Expected index to be present in create payload")
		}
		if index != float64(0) {
			t.Errorf("Expected index 0, got %v", index)
		}
		if modifier["target"] != request.Modifier.Target {
			t.Errorf("Expected target %s, got %v", request.Modifier.Target, modifier["target"])
		}
		validateModifierSourcePayload(t, modifier["source"], request.Modifier.Source)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tools.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	modifier, err := client.Tools.CreateModifier("tool-456", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateModifierResponse(t, *modifier, expected)
}

func TestToolsUpdateModifierOmitsUnsetFields(t *testing.T) {
	expected := modifierFixture()

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/modifiers/modifier-123" {
			t.Errorf("Expected modifier path, got %s", r.URL.Path)
		}

		modifier := decodeModifierEnvelope(t, r)
		if len(modifier) != 0 {
			t.Errorf("Expected unset fields to be omitted, got %v", modifier)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	modifier, err := client.Tools.UpdateModifier("modifier-123", tools.UpdateModifierRequest{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateModifierResponse(t, *modifier, expected)
}

func TestToolsUpdateModifierSendsCompleteSource(t *testing.T) {
	expected := modifierFixture()
	source := tools.ModifierSource{
		Type: tools.ModifierSourceTypeMetadata,
		Path: tools.ModifierSourcePathOriginEntityIdentifier,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		modifier := decodeModifierEnvelope(t, r)
		if len(modifier) != 1 {
			t.Errorf("Expected only source in update payload, got %v", modifier)
		}
		if _, present := modifier["index"]; present {
			t.Error("Did not expect immutable index in update payload")
		}
		validateModifierSourcePayload(t, modifier["source"], source)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	_, err := client.Tools.UpdateModifier("modifier-123", tools.UpdateModifierRequest{
		Modifier: tools.UpdateModifierData{Source: &source},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestToolsReplaceModifier(t *testing.T) {
	expected := modifierFixture()
	source := tools.ModifierSource{
		Type: tools.ModifierSourceTypeMetadata,
		Path: tools.ModifierSourcePathCurrentTimestamp,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/modifiers/modifier-123" {
			t.Errorf("Expected modifier path, got %s", r.URL.Path)
		}

		modifier := decodeModifierEnvelope(t, r)
		if _, present := modifier["index"]; present {
			t.Error("Did not expect immutable index in replace payload")
		}
		validateModifierSourcePayload(t, modifier["source"], source)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	modifier, err := client.Tools.ReplaceModifier("modifier-123", tools.UpdateModifierRequest{
		Modifier: tools.UpdateModifierData{
			OnMissingParent: tools.ModifierMissingPolicyError,
			Source:          &source,
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateModifierResponse(t, *modifier, expected)
}

func TestToolsDeleteModifierAcceptsInactiveResponse(t *testing.T) {
	inactive := modifierFixture()
	inactive.ProvisionState = "inactive"

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/modifiers/modifier-123" {
			t.Errorf("Expected modifier path, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.ModifierResponse{Data: inactive})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	if err := client.Tools.DeleteModifier("modifier-123"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestToolsModifierValidationDoesNotMakeHTTPRequest(t *testing.T) {
	var requestCount atomic.Int32
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	client := modifierClient(server.URL)
	validCreate := validCreateModifierRequest()
	validSource := validCreate.Modifier.Source

	tests := []struct {
		name        string
		expectedErr string
		call        func() error
	}{
		{
			name:        "get empty modifier ID",
			expectedErr: "modifier ID is required",
			call: func() error {
				_, err := client.Tools.GetModifier("")
				return err
			},
		},
		{
			name:        "create empty thought tool ID",
			expectedErr: "thought tool ID is required",
			call: func() error {
				_, err := client.Tools.CreateModifier("", validCreate)
				return err
			},
		},
		{
			name:        "create negative index",
			expectedErr: "modifier index must be non-negative",
			call: func() error {
				req := validCreate
				req.Modifier.Index = -1
				_, err := client.Tools.CreateModifier("tool-456", req)
				return err
			},
		},
		{
			name:        "create missing target",
			expectedErr: "modifier target is required",
			call: func() error {
				req := validCreate
				req.Modifier.Target = ""
				_, err := client.Tools.CreateModifier("tool-456", req)
				return err
			},
		},
		{
			name:        "create missing parent policy",
			expectedErr: "modifier on missing parent policy is required",
			call: func() error {
				req := validCreate
				req.Modifier.OnMissingParent = ""
				_, err := client.Tools.CreateModifier("tool-456", req)
				return err
			},
		},
		{
			name:        "create missing source policy",
			expectedErr: "modifier on missing source policy is required",
			call: func() error {
				req := validCreate
				req.Modifier.OnMissingSource = ""
				_, err := client.Tools.CreateModifier("tool-456", req)
				return err
			},
		},
		{
			name:        "create missing source type",
			expectedErr: "modifier source type is required",
			call: func() error {
				req := validCreate
				req.Modifier.Source.Type = ""
				_, err := client.Tools.CreateModifier("tool-456", req)
				return err
			},
		},
		{
			name:        "create missing source path",
			expectedErr: "modifier source path is required",
			call: func() error {
				req := validCreate
				req.Modifier.Source.Path = ""
				_, err := client.Tools.CreateModifier("tool-456", req)
				return err
			},
		},
		{
			name:        "update empty modifier ID",
			expectedErr: "modifier ID is required",
			call: func() error {
				_, err := client.Tools.UpdateModifier("", tools.UpdateModifierRequest{})
				return err
			},
		},
		{
			name:        "update missing source type",
			expectedErr: "modifier source type is required",
			call: func() error {
				source := validSource
				source.Type = ""
				_, err := client.Tools.UpdateModifier("modifier-123", tools.UpdateModifierRequest{
					Modifier: tools.UpdateModifierData{Source: &source},
				})
				return err
			},
		},
		{
			name:        "update missing source path",
			expectedErr: "modifier source path is required",
			call: func() error {
				source := validSource
				source.Path = ""
				_, err := client.Tools.UpdateModifier("modifier-123", tools.UpdateModifierRequest{
					Modifier: tools.UpdateModifierData{Source: &source},
				})
				return err
			},
		},
		{
			name:        "replace empty modifier ID",
			expectedErr: "modifier ID is required",
			call: func() error {
				_, err := client.Tools.ReplaceModifier("", tools.UpdateModifierRequest{})
				return err
			},
		},
		{
			name:        "replace missing source path",
			expectedErr: "modifier source path is required",
			call: func() error {
				source := validSource
				source.Path = ""
				_, err := client.Tools.ReplaceModifier("modifier-123", tools.UpdateModifierRequest{
					Modifier: tools.UpdateModifierData{Source: &source},
				})
				return err
			},
		},
		{
			name:        "delete empty modifier ID",
			expectedErr: "modifier ID is required",
			call: func() error {
				return client.Tools.DeleteModifier("")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.call()
			if err == nil || err.Error() != test.expectedErr {
				t.Errorf("Expected %q, got %v", test.expectedErr, err)
			}
		})
	}

	if count := requestCount.Load(); count != 0 {
		t.Errorf("Expected no HTTP requests, got %d", count)
	}
}

func TestToolsGetModifierNotFoundReturnsTypedError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string]string{"detail": "Not Found"},
		})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	_, err := client.Tools.GetModifier("modifier-123")
	assertModifierAPIError(t, err, http.StatusNotFound, "detail", "Not Found")
}

func TestToolsCreateModifierPreservesNestedValidationError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string]any{
				"source": map[string][]string{"path": {"can't be blank"}},
			},
		})
	})
	defer server.Close()

	client := modifierClient(server.URL)
	_, err := client.Tools.CreateModifier("tool-456", validCreateModifierRequest())
	assertModifierAPIError(t, err, http.StatusUnprocessableEntity, "source.path", "can't be blank")
}

func TestToolsModifierTransportErrorsIncludeOperation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	baseURL := server.URL
	server.Close()

	client := modifierClient(baseURL)
	source := validCreateModifierRequest().Modifier.Source
	tests := []struct {
		name      string
		operation string
		call      func() error
	}{
		{
			name:      "get",
			operation: "failed to get modifier",
			call: func() error {
				_, err := client.Tools.GetModifier("modifier-123")
				return err
			},
		},
		{
			name:      "create",
			operation: "failed to create modifier",
			call: func() error {
				_, err := client.Tools.CreateModifier("tool-456", validCreateModifierRequest())
				return err
			},
		},
		{
			name:      "update",
			operation: "failed to update modifier",
			call: func() error {
				_, err := client.Tools.UpdateModifier("modifier-123", tools.UpdateModifierRequest{
					Modifier: tools.UpdateModifierData{Source: &source},
				})
				return err
			},
		},
		{
			name:      "delete",
			operation: "failed to delete modifier",
			call: func() error {
				return client.Tools.DeleteModifier("modifier-123")
			},
		},
		{
			name:      "replace",
			operation: "failed to replace modifier",
			call: func() error {
				_, err := client.Tools.ReplaceModifier("modifier-123", tools.UpdateModifierRequest{
					Modifier: tools.UpdateModifierData{Source: &source},
				})
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.call()
			if err == nil || !strings.Contains(err.Error(), test.operation) {
				t.Errorf("Expected error containing %q, got %v", test.operation, err)
			}
		})
	}
}

func modifierClient(baseURL string) *tama.Client {
	return tama.NewClient(tama.Config{
		BaseURL: baseURL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	})
}

func modifierFixture() tools.Modifier {
	return tools.Modifier{
		ID:              "modifier-123",
		Index:           0,
		Target:          "/body/search/scope/user_id",
		OnMissingParent: tools.ModifierMissingPolicySkip,
		OnMissingSource: tools.ModifierMissingPolicyError,
		Source: tools.ModifierSource{
			Type: tools.ModifierSourceTypeMetadata,
			Path: tools.ModifierSourcePathActorIdentifier,
		},
		ThoughtToolID:  "tool-456",
		ProvisionState: "active",
	}
}

func validCreateModifierRequest() tools.CreateModifierRequest {
	modifier := modifierFixture()

	return tools.CreateModifierRequest{
		Modifier: tools.ModifierRequestData{
			Index:           modifier.Index,
			Target:          modifier.Target,
			OnMissingParent: modifier.OnMissingParent,
			OnMissingSource: modifier.OnMissingSource,
			Source:          modifier.Source,
		},
	}
}

func decodeModifierEnvelope(t *testing.T, r *http.Request) map[string]any {
	t.Helper()

	var envelope map[string]map[string]any
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode modifier request: %v", err)
	}

	modifier, present := envelope["modifier"]
	if !present {
		t.Fatal("Expected modifier envelope")
	}
	if len(envelope) != 1 {
		t.Errorf("Expected only modifier envelope, got %v", envelope)
	}

	return modifier
}

func validateModifierSourcePayload(t *testing.T, actual any, expected tools.ModifierSource) {
	t.Helper()

	source, ok := actual.(map[string]any)
	if !ok {
		t.Fatalf("Expected source object, got %T", actual)
	}
	if len(source) != 2 {
		t.Errorf("Expected complete source with 2 fields, got %v", source)
	}
	if source["type"] != expected.Type {
		t.Errorf("Expected source type %s, got %v", expected.Type, source["type"])
	}
	if source["path"] != expected.Path {
		t.Errorf("Expected source path %s, got %v", expected.Path, source["path"])
	}
}

func validateModifierResponse(t *testing.T, actual, expected tools.Modifier) {
	t.Helper()

	if actual != expected {
		t.Errorf("Expected modifier %+v, got %+v", expected, actual)
	}
}

func assertModifierAPIError(
	t *testing.T,
	err error,
	statusCode int,
	field string,
	message string,
) {
	t.Helper()

	if err == nil {
		t.Fatal("Expected API error, got nil")
	}

	var apiErr *tools.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected *tools.Error, got %T", err)
	}
	if apiErr.StatusCode != statusCode {
		t.Errorf("Expected status code %d, got %d", statusCode, apiErr.StatusCode)
	}
	if len(apiErr.Errors[field]) != 1 || apiErr.Errors[field][0] != message {
		t.Errorf("Expected %s error %q, got %v", field, message, apiErr.Errors)
	}
}
