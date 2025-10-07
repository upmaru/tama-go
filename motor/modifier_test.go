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

// createModifierMockServer creates a test HTTP server with the given handler for modifier tests.
func createModifierMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestMotorGetModifier(t *testing.T) {
	expected := motor.Modifier{
		ID:             "modifier-123",
		Name:           "sanitize",
		ActionID:       "action-456",
		Schema:         map[string]any{"rule": "trim"},
		ProvisionState: "active",
	}

	expectedResp := motor.ModifierResponse{Data: expected}

	server := createModifierMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/motor/modifiers/modifier-123" {
			t.Errorf("Expected path /provision/motor/modifiers/modifier-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	mod, err := client.Motor.GetModifier("modifier-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if mod.ID != expected.ID {
		t.Errorf("Expected modifier ID %s, got %s", expected.ID, mod.ID)
	}
	if mod.Name != expected.Name {
		t.Errorf("Expected modifier name %s, got %s", expected.Name, mod.Name)
	}
	if mod.ActionID != expected.ActionID {
		t.Errorf("Expected action_id %s, got %s", expected.ActionID, mod.ActionID)
	}
}

func TestMotorGetModifierError(t *testing.T) {
	server := createModifierMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := motor.Error{
			StatusCode: 404,
			Errors:     map[string][]string{"modifier": {"not found"}},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Motor.GetModifier("missing")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	var mErr *motor.Error
	if !errors.As(err, &mErr) {
		t.Fatalf("Expected motor.Error, got %T", err)
	}
	if mErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", mErr.StatusCode)
	}
}

func TestMotorCreateModifier(t *testing.T) {
	actionID := "action-456"
	expected := motor.Modifier{
		ID:             "modifier-789",
		Name:           "normalize",
		ActionID:       actionID,
		Schema:         map[string]any{"rule": "lowercase"},
		ProvisionState: "active",
	}
	server := createModifierMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/motor/actions/"+actionID+"/modifiers" {
			t.Errorf("Expected path /provision/motor/actions/%s/modifiers, got %s", actionID, r.URL.Path)
		}
		var req motor.CreateModifierRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if req.Modifier.Name == "" {
			t.Error("Expected name in request")
		}
		if req.Modifier.Schema == nil {
			t.Error("Expected schema in request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(motor.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	created, err := client.Motor.CreateModifier(
		actionID, motor.CreateModifierRequest{Modifier: motor.ModifierRequestData{
			Name:   "normalize",
			Schema: map[string]any{"rule": "lowercase"},
		}})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if created.ID != expected.ID {
		t.Errorf("Expected modifier ID %s, got %s", expected.ID, created.ID)
	}
}

func TestMotorCreateModifierValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", APIKey: "test-key"})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Motor.CreateModifier("",
		motor.CreateModifierRequest{Modifier: motor.ModifierRequestData{Name: "x", Schema: map[string]any{"a": 1}}})
	if err == nil {
		t.Error("Expected validation error for empty action ID")
	}
	_, err = client.Motor.CreateModifier("action-1",
		motor.CreateModifierRequest{Modifier: motor.ModifierRequestData{Schema: map[string]any{"a": 1}}})
	if err == nil {
		t.Error("Expected validation error for empty name")
	}
	_, err = client.Motor.CreateModifier("action-1",
		motor.CreateModifierRequest{Modifier: motor.ModifierRequestData{Name: "x"}})
	if err == nil {
		t.Error("Expected validation error for nil schema")
	}
}

func TestMotorUpdateModifier(t *testing.T) {
	modifierID := "modifier-123"
	expected := motor.Modifier{
		ID:             modifierID,
		Name:           "updated",
		ActionID:       "action-456",
		Schema:         map[string]any{"rule": "strip"},
		ProvisionState: "active",
	}
	server := createModifierMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/motor/modifiers/"+modifierID {
			t.Errorf("Expected path /provision/motor/modifiers/%s, got %s", modifierID, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(motor.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	updated, err := client.Motor.UpdateModifier(
		modifierID,
		motor.UpdateModifierRequest{
			Modifier: motor.UpdateModifierData{Name: "updated", Schema: map[string]any{"rule": "strip"}},
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if updated.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, updated.Name)
	}
}

func TestMotorReplaceModifier(t *testing.T) {
	modifierID := "modifier-123"
	expected := motor.Modifier{
		ID:             modifierID,
		Name:           "replaced",
		ActionID:       "action-456",
		Schema:         map[string]any{"rule": "noop"},
		ProvisionState: "active",
	}
	server := createModifierMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/motor/modifiers/"+modifierID {
			t.Errorf("Expected path /provision/motor/modifiers/%s, got %s", modifierID, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(motor.ModifierResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	replaced, err := client.Motor.ReplaceModifier(
		modifierID,
		motor.UpdateModifierRequest{
			Modifier: motor.UpdateModifierData{Name: "replaced", Schema: map[string]any{"rule": "noop"}},
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if replaced.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, replaced.Name)
	}
}

func TestMotorDeleteModifier(t *testing.T) {
	modifierID := "modifier-123"
	server := createModifierMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/motor/modifiers/"+modifierID {
			t.Errorf("Expected path /provision/motor/modifiers/%s, got %s", modifierID, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	if err := client.Motor.DeleteModifier(modifierID); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
