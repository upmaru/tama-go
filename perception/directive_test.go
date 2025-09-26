package perception_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/perception"
)

func TestPerceptionGetDirective(t *testing.T) {
	expected := perception.Directive{
		ID:              "directive-123",
		ThoughtPathID:   "path-123",
		PromptID:        "prompt-123",
		TargetThoughtID: "thought-999",
		ProvisionState:  "active",
	}

	expectedResp := perception.DirectiveResponse{Data: expected}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/perception/directives/directive-123" {
			t.Errorf("Expected path /provision/perception/directives/directive-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	directive, err := client.Perception.GetDirective("directive-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if directive.ID != expected.ID ||
		directive.ThoughtPathID != expected.ThoughtPathID ||
		directive.PromptID != expected.PromptID ||
		directive.TargetThoughtID != expected.TargetThoughtID ||
		directive.ProvisionState != expected.ProvisionState {
		t.Errorf("Directive mismatch: got %+v, expected %+v", directive, expected)
	}
}

func TestPerceptionGetDirectiveError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := perception.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"directive": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	_, err := client.Perception.GetDirective("nonexistent")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	var pErr *perception.Error
	if !errors.As(err, &pErr) {
		t.Fatalf("Expected perception.Error, got %T", err)
	}
	if pErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", pErr.StatusCode)
	}
}

func TestPerceptionCreateDirective(t *testing.T) {
	request := perception.CreateDirectiveRequest{
		Directive: perception.DirectiveRequestData{
			PromptID:        "prompt-123",
			TargetThoughtID: "thought-999",
		},
	}
	expected := perception.Directive{
		ID:              "directive-456",
		ThoughtPathID:   "path-123",
		PromptID:        "prompt-123",
		TargetThoughtID: "thought-999",
		ProvisionState:  "pending",
	}
	expectedResp := perception.DirectiveResponse{Data: expected}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/perception/paths/path-123/directives" {
			t.Errorf("Expected path /provision/perception/paths/path-123/directives, got %s", r.URL.Path)
		}
		var recv perception.CreateDirectiveRequest
		if err := json.NewDecoder(r.Body).Decode(&recv); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if recv.Directive.PromptID != request.Directive.PromptID {
			t.Errorf(
				"Expected prompt_id %s, got %s",
				request.Directive.PromptID,
				recv.Directive.PromptID,
			)
		}
		if recv.Directive.TargetThoughtID != request.Directive.TargetThoughtID {
			t.Errorf(
				"Expected target_thought_id %s, got %s",
				request.Directive.TargetThoughtID,
				recv.Directive.TargetThoughtID,
			)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	directive, err := client.Perception.CreateDirective("path-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if directive.PromptID != expected.PromptID {
		t.Errorf("Expected prompt_id %s, got %s", expected.PromptID, directive.PromptID)
	}
}

func TestPerceptionCreateDirectiveValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", APIKey: "test-key"})
	// Empty path ID
	_, err := client.Perception.CreateDirective(
		"",
		perception.CreateDirectiveRequest{Directive: perception.DirectiveRequestData{PromptID: "p"}},
	)
	if err == nil {
		t.Error("Expected validation error for empty path ID")
	}
	// Empty prompt ID
	_, err = client.Perception.CreateDirective(
		"path-123",
		perception.CreateDirectiveRequest{
			Directive: perception.DirectiveRequestData{PromptID: ""},
		},
	)
	if err == nil {
		t.Error("Expected validation error for empty prompt ID")
	}

	// Empty target_thought_id
	_, err = client.Perception.CreateDirective(
		"path-123",
		perception.CreateDirectiveRequest{
			Directive: perception.DirectiveRequestData{PromptID: "p", TargetThoughtID: ""},
		},
	)
	if err == nil {
		t.Error("Expected validation error for empty target thought ID")
	}
}

func TestPerceptionCreateDirectiveWithFieldErrors(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := map[string]any{
			"errors": map[string][]string{
				"prompt_id": {"can't be blank", "is invalid"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	_, err := client.Perception.CreateDirective(
		"path-123",
		perception.CreateDirectiveRequest{
			Directive: perception.DirectiveRequestData{PromptID: "invalid", TargetThoughtID: "thought-999"},
		},
	)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	errMsg := err.Error()
	if !strings.Contains(errMsg, "prompt_id can't be blank") || !strings.Contains(errMsg, "prompt_id is invalid") {
		t.Errorf("Expected prompt_id field errors in error message, got %s", errMsg)
	}
}

func TestPerceptionUpdateDirective(t *testing.T) {
	request := perception.UpdateDirectiveRequest{
		Directive: perception.UpdateDirectiveData{
			PromptID:        "prompt-updated",
			TargetThoughtID: "thought-999",
		},
	}
	expected := perception.Directive{
		ID:              "directive-123",
		ThoughtPathID:   "path-123",
		PromptID:        "prompt-updated",
		TargetThoughtID: "thought-999",
		ProvisionState:  "active",
	}
	expectedResp := perception.DirectiveResponse{Data: expected}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/perception/directives/directive-123" {
			t.Errorf("Expected path /provision/perception/directives/directive-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	directive, err := client.Perception.UpdateDirective("directive-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if directive.PromptID != expected.PromptID {
		t.Errorf("Expected prompt_id %s, got %s", expected.PromptID, directive.PromptID)
	}
	if directive.TargetThoughtID != expected.TargetThoughtID {
		t.Errorf("Expected target_thought_id %s, got %s", expected.TargetThoughtID, directive.TargetThoughtID)
	}
}

func TestPerceptionReplaceDirective(t *testing.T) {
	request := perception.UpdateDirectiveRequest{
		Directive: perception.UpdateDirectiveData{
			PromptID:        "prompt-replaced",
			TargetThoughtID: "thought-999",
		},
	}
	expected := perception.Directive{
		ID:              "directive-123",
		ThoughtPathID:   "path-123",
		PromptID:        "prompt-replaced",
		TargetThoughtID: "thought-999",
		ProvisionState:  "active",
	}
	expectedResp := perception.DirectiveResponse{Data: expected}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/perception/directives/directive-123" {
			t.Errorf("Expected path /provision/perception/directives/directive-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	directive, err := client.Perception.ReplaceDirective("directive-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if directive.PromptID != expected.PromptID {
		t.Errorf("Expected prompt_id %s, got %s", expected.PromptID, directive.PromptID)
	}
	if directive.TargetThoughtID != expected.TargetThoughtID {
		t.Errorf("Expected target_thought_id %s, got %s", expected.TargetThoughtID, directive.TargetThoughtID)
	}
}

func TestPerceptionDeleteDirective(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/perception/directives/directive-123" {
			t.Errorf("Expected path /provision/perception/directives/directive-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{BaseURL: server.URL, APIKey: "test-key", Timeout: 10 * time.Second})
	if err := client.Perception.DeleteDirective("directive-123"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionDirectiveEmptyIDValidations(t *testing.T) {
	client := tama.NewClient(tama.Config{BaseURL: "https://api.example.com", APIKey: "test-key"})
	if _, err := client.Perception.GetDirective(""); err == nil {
		t.Error("Expected validation error for empty directive ID in GetDirective")
	}
	if _, err := client.Perception.UpdateDirective("", perception.UpdateDirectiveRequest{Directive: perception.UpdateDirectiveData{PromptID: "p"}}); err == nil {
		t.Error("Expected validation error for empty directive ID in UpdateDirective")
	}
	if _, err := client.Perception.ReplaceDirective("", perception.UpdateDirectiveRequest{Directive: perception.UpdateDirectiveData{PromptID: "p"}}); err == nil {
		t.Error("Expected validation error for empty directive ID in ReplaceDirective")
	}
	if err := client.Perception.DeleteDirective(""); err == nil {
		t.Error("Expected validation error for empty directive ID in DeleteDirective")
	}
}
