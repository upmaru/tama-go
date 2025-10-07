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

func TestToolsGetOutput(t *testing.T) {
	expected := tools.Output{
		ID:             "output-123",
		ThoughtToolID:  "tool-456",
		ClassCorpusID:  "corpus-789",
		ProvisionState: "active",
	}

	expectedResp := tools.OutputResponse{Data: expected}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/outputs/output-123" {
			t.Errorf("Expected path /provision/tools/outputs/output-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	out, err := client.Tools.GetOutput("output-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOutputResponse(t, *out, expected)
}

func TestToolsGetOutputError(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := tools.Error{StatusCode: 404, Errors: map[string][]string{"output": {"not found"}}}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	_, err = client.Tools.GetOutput("missing")
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

func TestToolsCreateOutput(t *testing.T) {
	request := tools.CreateOutputRequest{Output: tools.OutputRequestData{ClassCorpusID: "corpus-789"}}
	expected := tools.Output{
		ID:             "output-456",
		ThoughtToolID:  "tool-123",
		ClassCorpusID:  "corpus-789",
		ProvisionState: "pending",
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/tool-123/outputs" {
			t.Errorf("Expected path /provision/tools/tool-123/outputs, got %s", r.URL.Path)
		}
		var received tools.CreateOutputRequest
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if received.Output.ClassCorpusID != request.Output.ClassCorpusID {
			t.Errorf("Expected class_corpus_id %s, got %s", request.Output.ClassCorpusID, received.Output.ClassCorpusID)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tools.OutputResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	out, err := client.Tools.CreateOutput("tool-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOutputResponse(t, *out, expected)
}

func TestToolsCreateOutputValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	// Missing thought tool id
	_, err = client.Tools.CreateOutput(
		"",
		tools.CreateOutputRequest{Output: tools.OutputRequestData{ClassCorpusID: "c1"}},
	)
	if err == nil {
		t.Error("Expected validation error for empty thought tool ID")
	}
	// Missing class corpus id
	_, err = client.Tools.CreateOutput(
		"tool-1",
		tools.CreateOutputRequest{Output: tools.OutputRequestData{ClassCorpusID: ""}},
	)
	if err == nil {
		t.Error("Expected validation error for empty class corpus ID")
	}
}

func TestToolsUpdateOutput(t *testing.T) {
	request := tools.UpdateOutputRequest{Output: tools.UpdateOutputData{ClassCorpusID: "corpus-updated"}}
	expected := tools.Output{
		ID:             "output-123",
		ThoughtToolID:  "tool-456",
		ClassCorpusID:  "corpus-updated",
		ProvisionState: "active",
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/outputs/output-123" {
			t.Errorf("Expected path /provision/tools/outputs/output-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.OutputResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	out, err := client.Tools.UpdateOutput("output-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOutputResponse(t, *out, expected)
}

func TestToolsReplaceOutput(t *testing.T) {
	request := tools.UpdateOutputRequest{Output: tools.UpdateOutputData{ClassCorpusID: "corpus-replaced"}}
	expected := tools.Output{
		ID:             "output-123",
		ThoughtToolID:  "tool-456",
		ClassCorpusID:  "corpus-replaced",
		ProvisionState: "active",
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/outputs/output-123" {
			t.Errorf("Expected path /provision/tools/outputs/output-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools.OutputResponse{Data: expected})
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	out, err := client.Tools.ReplaceOutput("output-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ValidateOutputResponse(t, *out, expected)
}

func TestToolsDeleteOutput(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/tools/outputs/output-123" {
			t.Errorf("Expected path /provision/tools/outputs/output-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client, err := tama.NewClient(tama.Config{
		BaseURL:        server.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	err = client.Tools.DeleteOutput("output-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
