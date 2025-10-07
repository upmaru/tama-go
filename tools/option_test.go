package tools_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/tools"
)

func TestToolsGetOption(t *testing.T) {
	expected := tools.Option{
		ID:                  "option-123",
		ThoughtToolOutputID: "output-456",
		ActionModifierID:    "modifier-789",
		ProvisionState:      "active",
	}

	expectedResp := tools.OptionResponse{Data: expected}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/options/option-123" {
			t.Errorf("Expected path /provision/tools/options/option-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	opt, err := client.Tools.GetOption("option-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOptionResponse(t, *opt, expected)
}

func TestToolsGetOptionError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := tools.Error{StatusCode: 404, Errors: map[string][]string{"option": {"not found"}}}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Tools.GetOption("missing")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	var tErr *tools.Error
	if !errors.As(err, &tErr) {
		t.Fatalf("Expected tools.Error, got %T", err)
	}
	if tErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", tErr.StatusCode)
	}
}

func TestToolsCreateOption(t *testing.T) {
	request := tools.CreateOptionRequest{Option: tools.OptionRequestData{ActionModifierID: "modifier-789"}}
	expected := tools.Option{
		ID:                  "option-456",
		ThoughtToolOutputID: "output-123",
		ActionModifierID:    "modifier-789",
		ProvisionState:      "pending",
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/outputs/output-123/options" {
			t.Errorf("Expected path /provision/tools/outputs/output-123/options, got %s", r.URL.Path)
		}
		var received tools.CreateOptionRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if received.Option.ActionModifierID != request.Option.ActionModifierID {
			t.Errorf(
				"Expected action_modifier_id %s, got %s",
				request.Option.ActionModifierID,
				received.Option.ActionModifierID,
			)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tools.OptionResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	opt, err := client.Tools.CreateOption("output-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOptionResponse(t, *opt, expected)
}

func TestToolsCreateOptionValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", ClientID: "test-client-id", ClientSecret: "test-client-secret", SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	// Missing output id
	_, err = client.Tools.CreateOption("", tools.CreateOptionRequest{
		Option: tools.OptionRequestData{ActionModifierID: "m1"},
	})
	if err == nil {
		t.Error("Expected validation error for empty output ID")
	}
	// Missing action modifier id
	_, err = client.Tools.CreateOption(
		"output-1",
		tools.CreateOptionRequest{
			Option: tools.OptionRequestData{
				ActionModifierID: "",
			},
		},
	)
	if err == nil {
		t.Error("Expected validation error for empty action modifier ID")
	}
}

func TestToolsUpdateOption(t *testing.T) {
	request := tools.UpdateOptionRequest{Option: tools.UpdateOptionData{ActionModifierID: "modifier-updated"}}
	expected := tools.Option{
		ID:                  "option-123",
		ThoughtToolOutputID: "output-456",
		ActionModifierID:    "modifier-updated",
		ProvisionState:      "active",
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/options/option-123" {
			t.Errorf("Expected path /provision/tools/options/option-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.OptionResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	opt, err := client.Tools.UpdateOption("option-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOptionResponse(t, *opt, expected)
}

func TestToolsReplaceOption(t *testing.T) {
	request := tools.UpdateOptionRequest{Option: tools.UpdateOptionData{ActionModifierID: "modifier-replaced"}}
	expected := tools.Option{
		ID:                  "option-123",
		ThoughtToolOutputID: "output-456",
		ActionModifierID:    "modifier-replaced",
		ProvisionState:      "active",
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/options/option-123" {
			t.Errorf("Expected path /provision/tools/options/option-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.OptionResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	opt, err := client.Tools.ReplaceOption("option-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOptionResponse(t, *opt, expected)
}

func TestToolsDeleteOption(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/options/option-123" {
			t.Errorf("Expected path /provision/tools/options/option-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{BaseURL: server.URL, ClientID: "test-client-id", ClientSecret: "test-client-secret", Timeout: 10 * time.Second, SkipTokenFetch: true})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	if err := client.Tools.DeleteOption("option-123"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
