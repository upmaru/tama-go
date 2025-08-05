package perception_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/perception"
)

func TestPerceptionToolEndpoints(t *testing.T) {
	// This is a basic placeholder test to ensure the file compiles
	// Real implementation would require mocking HTTP requests

	// Check that Tool struct has the right fields
	tool := perception.Tool{
		ID:             "test-id",
		ThoughtID:      "test-thought-id",
		ActionID:       "test-action-id",
		ProvisionState: "test-state",
	}

	if tool.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", tool.ID)
	}

	if tool.ThoughtID != "test-thought-id" {
		t.Errorf("Expected ThoughtID 'test-thought-id', got '%s'", tool.ThoughtID)
	}

	if tool.ActionID != "test-action-id" {
		t.Errorf("Expected ActionID 'test-action-id', got '%s'", tool.ActionID)
	}

	if tool.ProvisionState != "test-state" {
		t.Errorf("Expected ProvisionState 'test-state', got '%s'", tool.ProvisionState)
	}
}

func TestPerceptionCreateTool(t *testing.T) {
	request := perception.CreateToolRequest{
		Tool: perception.CreateToolData{
			ActionID: "action-123",
		},
	}

	expectedTool := perception.Tool{
		ID:             "tool-456",
		ThoughtID:      "thought-123",
		ActionID:       "action-123",
		ProvisionState: "pending",
	}

	expectedResponse := perception.ToolResponse{
		Data: expectedTool,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/thoughts/thought-123/tools" {
			t.Errorf("Expected path /provision/perception/thoughts/thought-123/tools, got %s", r.URL.Path)
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
	tool, err := client.Perception.CreateTool("thought-123", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tool.ID != expectedTool.ID {
		t.Errorf("Expected tool ID %s, got %s", expectedTool.ID, tool.ID)
	}

	if tool.ThoughtID != expectedTool.ThoughtID {
		t.Errorf("Expected tool ThoughtID %s, got %s", expectedTool.ThoughtID, tool.ThoughtID)
	}

	if tool.ActionID != expectedTool.ActionID {
		t.Errorf("Expected tool ActionID %s, got %s", expectedTool.ActionID, tool.ActionID)
	}

	if tool.ProvisionState != expectedTool.ProvisionState {
		t.Errorf(
			"Expected tool ProvisionState %s, got %s",
			expectedTool.ProvisionState,
			tool.ProvisionState,
		)
	}
}

func TestPerceptionGetTool(t *testing.T) {
	expectedTool := perception.Tool{
		ID:             "tool-123",
		ThoughtID:      "thought-456",
		ActionID:       "action-789",
		ProvisionState: "active",
	}

	expectedResponse := perception.ToolResponse{
		Data: expectedTool,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/tools/tool-123" {
			t.Errorf("Expected path /provision/perception/tools/tool-123, got %s", r.URL.Path)
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
	tool, err := client.Perception.GetTool("tool-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tool.ID != expectedTool.ID {
		t.Errorf("Expected tool ID %s, got %s", expectedTool.ID, tool.ID)
	}
}

func TestPerceptionUpdateTool(t *testing.T) {
	request := perception.UpdateToolRequest{
		Tool: perception.UpdateToolData{
			ActionID: "updated-action",
		},
	}

	expectedTool := perception.Tool{
		ID:             "tool-789",
		ThoughtID:      "thought-456",
		ActionID:       "updated-action",
		ProvisionState: "updated",
	}

	expectedResponse := perception.ToolResponse{
		Data: expectedTool,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/tools/tool-789" {
			t.Errorf("Expected path /provision/perception/tools/tool-789, got %s", r.URL.Path)
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
	tool, err := client.Perception.UpdateTool("tool-789", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tool.ActionID != expectedTool.ActionID {
		t.Errorf("Expected tool ActionID %s, got %s", expectedTool.ActionID, tool.ActionID)
	}

	if tool.ProvisionState != expectedTool.ProvisionState {
		t.Errorf(
			"Expected tool ProvisionState %s, got %s",
			expectedTool.ProvisionState,
			tool.ProvisionState,
		)
	}
}

func TestPerceptionDeleteTool(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/tools/tool-456" {
			t.Errorf("Expected path /provision/perception/tools/tool-456, got %s", r.URL.Path)
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
	err := client.Perception.DeleteTool("tool-456")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPerceptionReplaceTool(t *testing.T) {
	request := perception.UpdateToolRequest{
		Tool: perception.UpdateToolData{
			ActionID: "replaced-action",
		},
	}

	expectedTool := perception.Tool{
		ID:             "tool-999",
		ThoughtID:      "thought-456",
		ActionID:       "replaced-action",
		ProvisionState: "replaced",
	}

	expectedResponse := perception.ToolResponse{
		Data: expectedTool,
	}

	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/perception/tools/tool-999" {
			t.Errorf("Expected path /provision/perception/tools/tool-999, got %s", r.URL.Path)
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
	tool, err := client.Perception.ReplaceTool("tool-999", request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tool.ActionID != expectedTool.ActionID {
		t.Errorf("Expected tool ActionID %s, got %s", expectedTool.ActionID, tool.ActionID)
	}

	if tool.ProvisionState != expectedTool.ProvisionState {
		t.Errorf(
			"Expected tool ProvisionState %s, got %s",
			expectedTool.ProvisionState,
			tool.ProvisionState,
		)
	}
}

func TestPerceptionCreateToolValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	// Test empty thought ID
	_, err := client.Perception.CreateTool("", perception.CreateToolRequest{
		Tool: perception.CreateToolData{
			ActionID: "action-123",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty thought ID")
	}

	// Test empty action ID
	_, err = client.Perception.CreateTool("thought-123", perception.CreateToolRequest{
		Tool: perception.CreateToolData{
			ActionID: "",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty action ID")
	}
}

func TestPerceptionGetToolEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.GetTool("")
	if err == nil {
		t.Error("Expected validation error for empty tool ID in GetTool")
	}
}

func TestPerceptionUpdateToolEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.UpdateTool("", perception.UpdateToolRequest{
		Tool: perception.UpdateToolData{ActionID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty tool ID in UpdateTool")
	}
}

func TestPerceptionReplaceToolEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Perception.ReplaceTool("", perception.UpdateToolRequest{
		Tool: perception.UpdateToolData{ActionID: "test"},
	})
	if err == nil {
		t.Error("Expected validation error for empty tool ID in ReplaceTool")
	}
}

func TestPerceptionDeleteToolEmptyID(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	err := client.Perception.DeleteTool("")
	if err == nil {
		t.Error("Expected validation error for empty tool ID in DeleteTool")
	}
}

// Helper to test URL construction (basic validation).
func TestToolURLConstruction(t *testing.T) {
	testCases := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "valid tool ID",
			id:       "tool-123",
			expected: "/provision/perception/tools/tool-123",
		},
		{
			name:     "empty tool ID",
			id:       "",
			expected: "/provision/perception/tools/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(subT *testing.T) {
			// Basic validation that we can construct these URLs
			url := "/provision/perception/tools/" + tc.id
			if len(url) == 0 {
				subT.Error("URL should not be empty")
			}
			// Verify the URL contains expected components
			if tc.id != "" && !strings.Contains(url, tc.id) {
				subT.Errorf("Expected URL to contain ID %s, got %s", tc.id, url)
			}
		})
	}
}
