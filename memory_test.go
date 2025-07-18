package tama_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/memory"
)

func createMockServerForMemory(_ *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

// TestMemoryGetPrompt tests retrieving a prompt by ID.
func TestMemoryGetPrompt(t *testing.T) {
	expectedPrompt := memory.Prompt{
		ID:           "prompt-123",
		Name:         "Test Prompt",
		Slug:         "test-prompt",
		Content:      "You are a helpful assistant.",
		Role:         "system",
		SpaceID:      "space-456",
		CurrentState: "active",
	}

	expectedResponse := memory.PromptResponse{
		Data: expectedPrompt,
	}

	server := createMockServerForMemory(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/memory/prompts/prompt-123" {
			t.Errorf("Expected path /provision/memory/prompts/prompt-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	prompt, err := client.Memory.GetPrompt("prompt-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if prompt.ID != expectedPrompt.ID {
		t.Errorf("Expected prompt ID %s, got %s", expectedPrompt.ID, prompt.ID)
	}

	if prompt.Name != expectedPrompt.Name {
		t.Errorf("Expected prompt name %s, got %s", expectedPrompt.Name, prompt.Name)
	}

	if prompt.Content != expectedPrompt.Content {
		t.Errorf("Expected prompt content %s, got %s", expectedPrompt.Content, prompt.Content)
	}

	if prompt.Role != expectedPrompt.Role {
		t.Errorf("Expected prompt role %s, got %s", expectedPrompt.Role, prompt.Role)
	}

	if prompt.SpaceID != expectedPrompt.SpaceID {
		t.Errorf("Expected prompt space ID %s, got %s", expectedPrompt.SpaceID, prompt.SpaceID)
	}

	if prompt.CurrentState != expectedPrompt.CurrentState {
		t.Errorf("Expected prompt current state %s, got %s", expectedPrompt.CurrentState, prompt.CurrentState)
	}
}

// TestMemoryCreatePrompt tests creating a new prompt.
func TestMemoryCreatePrompt(t *testing.T) {
	expectedPrompt := memory.Prompt{
		ID:           "prompt-789",
		Name:         "New Prompt",
		Slug:         "new-prompt",
		Content:      "You are a coding assistant.",
		Role:         "assistant",
		SpaceID:      "space-123",
		CurrentState: "pending",
	}

	expectedResponse := memory.PromptResponse{
		Data: expectedPrompt,
	}

	server := createMockServerForMemory(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/memory/spaces/space-123/prompts" {
			t.Errorf("Expected path /provision/memory/spaces/space-123/prompts, got %s", r.URL.Path)
		}

		var req memory.CreatePromptRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Prompt.Name != "New Prompt" {
			t.Errorf("Expected request name 'New Prompt', got %s", req.Prompt.Name)
		}

		if req.Prompt.Content != "You are a coding assistant." {
			t.Errorf("Expected request content 'You are a coding assistant.', got %s", req.Prompt.Content)
		}

		if req.Prompt.Role != "assistant" {
			t.Errorf("Expected request role 'assistant', got %s", req.Prompt.Role)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := memory.CreatePromptRequest{
		Prompt: memory.PromptRequestData{
			Name:    "New Prompt",
			Content: "You are a coding assistant.",
			Role:    "assistant",
		},
	}

	prompt, err := client.Memory.CreatePrompt("space-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if prompt.ID != expectedPrompt.ID {
		t.Errorf("Expected prompt ID %s, got %s", expectedPrompt.ID, prompt.ID)
	}

	if prompt.Name != expectedPrompt.Name {
		t.Errorf("Expected prompt name %s, got %s", expectedPrompt.Name, prompt.Name)
	}

	if prompt.Content != expectedPrompt.Content {
		t.Errorf("Expected prompt content %s, got %s", expectedPrompt.Content, prompt.Content)
	}

	if prompt.Role != expectedPrompt.Role {
		t.Errorf("Expected prompt role %s, got %s", expectedPrompt.Role, prompt.Role)
	}

	if prompt.SpaceID != expectedPrompt.SpaceID {
		t.Errorf("Expected prompt space ID %s, got %s", expectedPrompt.SpaceID, prompt.SpaceID)
	}

	if prompt.CurrentState != expectedPrompt.CurrentState {
		t.Errorf("Expected prompt current state %s, got %s", expectedPrompt.CurrentState, prompt.CurrentState)
	}
}

// TestMemoryCreatePromptValidation tests validation for creating a prompt.
func TestMemoryCreatePromptValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	tests := []struct {
		name    string
		spaceID string
		req     memory.CreatePromptRequest
		wantErr string
	}{
		{
			name:    "empty space ID",
			spaceID: "",
			req: memory.CreatePromptRequest{
				Prompt: memory.PromptRequestData{
					Name:    "Test Prompt",
					Content: "Test content",
					Role:    "user",
				},
			},
			wantErr: "space ID is required",
		},
		{
			name:    "empty prompt name",
			spaceID: "space-123",
			req: memory.CreatePromptRequest{
				Prompt: memory.PromptRequestData{
					Name:    "",
					Content: "Test content",
					Role:    "user",
				},
			},
			wantErr: "prompt name is required",
		},
		{
			name:    "empty prompt content",
			spaceID: "space-123",
			req: memory.CreatePromptRequest{
				Prompt: memory.PromptRequestData{
					Name:    "Test Prompt",
					Content: "",
					Role:    "user",
				},
			},
			wantErr: "prompt content is required",
		},
		{
			name:    "empty prompt role",
			spaceID: "space-123",
			req: memory.CreatePromptRequest{
				Prompt: memory.PromptRequestData{
					Name:    "Test Prompt",
					Content: "Test content",
					Role:    "",
				},
			},
			wantErr: "prompt role is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Memory.CreatePrompt(tt.spaceID, tt.req)
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// TestMemoryGetPrompt_EmptyIDValidation tests validation for empty ID.
func TestMemoryGetPrompt_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	_, err := client.Memory.GetPrompt("")
	if err == nil {
		t.Errorf("Expected error for empty ID, got nil")
	}
	if err.Error() != "prompt ID is required" {
		t.Errorf("Expected 'prompt ID is required', got %q", err.Error())
	}
}

// TestMemoryUpdatePrompt tests updating a prompt.
func TestMemoryUpdatePrompt(t *testing.T) {
	expectedPrompt := memory.Prompt{
		ID:           "prompt-123",
		Name:         "Updated Prompt",
		Slug:         "updated-prompt",
		Content:      "You are an updated assistant.",
		Role:         "system",
		SpaceID:      "space-456",
		CurrentState: "active",
	}

	expectedResponse := memory.PromptResponse{
		Data: expectedPrompt,
	}

	server := createMockServerForMemory(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/memory/prompts/prompt-123" {
			t.Errorf("Expected path /provision/memory/prompts/prompt-123, got %s", r.URL.Path)
		}

		var req memory.UpdatePromptRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Prompt.Name != "Updated Prompt" {
			t.Errorf("Expected request name 'Updated Prompt', got %s", req.Prompt.Name)
		}

		if req.Prompt.Content != "You are an updated assistant." {
			t.Errorf("Expected request content 'You are an updated assistant.', got %s", req.Prompt.Content)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	updateReq := memory.UpdatePromptRequest{
		Prompt: memory.UpdatePromptData{
			Name:    "Updated Prompt",
			Content: "You are an updated assistant.",
		},
	}

	prompt, err := client.Memory.UpdatePrompt("prompt-123", updateReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if prompt.Name != expectedPrompt.Name {
		t.Errorf("Expected prompt name %s, got %s", expectedPrompt.Name, prompt.Name)
	}

	if prompt.Content != expectedPrompt.Content {
		t.Errorf("Expected prompt content %s, got %s", expectedPrompt.Content, prompt.Content)
	}
}

// TestMemoryUpdatePrompt_EmptyIDValidation tests validation for empty ID in update.
func TestMemoryUpdatePrompt_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	updateReq := memory.UpdatePromptRequest{
		Prompt: memory.UpdatePromptData{
			Name: "Updated Prompt",
		},
	}

	_, err := client.Memory.UpdatePrompt("", updateReq)
	if err == nil {
		t.Errorf("Expected error for empty ID, got nil")
	}
	if err.Error() != "prompt ID is required" {
		t.Errorf("Expected 'prompt ID is required', got %q", err.Error())
	}
}

// TestMemoryReplacePrompt tests replacing a prompt.
func TestMemoryReplacePrompt(t *testing.T) {
	expectedPrompt := memory.Prompt{
		ID:           "prompt-123",
		Name:         "Replaced Prompt",
		Slug:         "replaced-prompt",
		Content:      "You are a completely new assistant.",
		Role:         "assistant",
		SpaceID:      "space-456",
		CurrentState: "active",
	}

	expectedResponse := memory.PromptResponse{
		Data: expectedPrompt,
	}

	server := createMockServerForMemory(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/memory/prompts/prompt-123" {
			t.Errorf("Expected path /provision/memory/prompts/prompt-123, got %s", r.URL.Path)
		}

		var req memory.UpdatePromptRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Prompt.Name != "Replaced Prompt" {
			t.Errorf("Expected request name 'Replaced Prompt', got %s", req.Prompt.Name)
		}

		if req.Prompt.Content != "You are a completely new assistant." {
			t.Errorf("Expected request content 'You are a completely new assistant.', got %s", req.Prompt.Content)
		}

		if req.Prompt.Role != "assistant" {
			t.Errorf("Expected request role 'assistant', got %s", req.Prompt.Role)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	replaceReq := memory.UpdatePromptRequest{
		Prompt: memory.UpdatePromptData{
			Name:    "Replaced Prompt",
			Content: "You are a completely new assistant.",
			Role:    "assistant",
		},
	}

	prompt, err := client.Memory.ReplacePrompt("prompt-123", replaceReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if prompt.Name != expectedPrompt.Name {
		t.Errorf("Expected prompt name %s, got %s", expectedPrompt.Name, prompt.Name)
	}

	if prompt.Content != expectedPrompt.Content {
		t.Errorf("Expected prompt content %s, got %s", expectedPrompt.Content, prompt.Content)
	}

	if prompt.Role != expectedPrompt.Role {
		t.Errorf("Expected prompt role %s, got %s", expectedPrompt.Role, prompt.Role)
	}
}

// TestMemoryReplacePrompt_EmptyIDValidation tests validation for empty ID in replace.
func TestMemoryReplacePrompt_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	replaceReq := memory.UpdatePromptRequest{
		Prompt: memory.UpdatePromptData{
			Name: "Replaced Prompt",
		},
	}

	_, err := client.Memory.ReplacePrompt("", replaceReq)
	if err == nil {
		t.Errorf("Expected error for empty ID, got nil")
	}
	if err.Error() != "prompt ID is required" {
		t.Errorf("Expected 'prompt ID is required', got %q", err.Error())
	}
}

// TestMemoryDeletePrompt tests deleting a prompt.
func TestMemoryDeletePrompt(t *testing.T) {
	server := createMockServerForMemory(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/memory/prompts/prompt-123" {
			t.Errorf("Expected path /provision/memory/prompts/prompt-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	err := client.Memory.DeletePrompt("prompt-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// TestMemoryDeletePrompt_EmptyIDValidation tests validation for empty ID in delete.
func TestMemoryDeletePrompt_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "http://localhost",
		APIKey:  "test-key",
	})

	err := client.Memory.DeletePrompt("")
	if err == nil {
		t.Errorf("Expected error for empty ID, got nil")
	}
	if err.Error() != "prompt ID is required" {
		t.Errorf("Expected 'prompt ID is required', got %q", err.Error())
	}
}

// TestMemoryFieldSpecificErrors tests field-specific error handling.
func TestMemoryFieldSpecificErrors(t *testing.T) {
	// Test memory field-specific errors
	fieldErr := &memory.Error{
		StatusCode: 422,
		Errors: map[string][]string{
			"name":    {"can't be blank"},
			"content": {"is too short (minimum is 10 characters)"},
			"role":    {"must be one of: system, user, assistant"},
		},
	}

	errorMsg := fieldErr.Error()
	// Check that all field errors are included
	if !strings.Contains(errorMsg, "name can't be blank") {
		t.Errorf("Expected error message to contain 'name can't be blank', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "content is too short (minimum is 10 characters)") {
		t.Errorf("Expected error message to contain 'content is too short (minimum is 10 characters)', got %s",
			errorMsg)
	}
	if !strings.Contains(errorMsg, "role must be one of: system, user, assistant") {
		t.Errorf("Expected error message to contain 'role must be one of: system, user, assistant', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "API error 422:") {
		t.Errorf("Expected error message to contain status code, got %s", errorMsg)
	}

	// Test error with only status code
	statusOnlyErr := &memory.Error{
		StatusCode: 404,
	}

	expectedStatusMsg := "API error 404"
	if statusOnlyErr.Error() != expectedStatusMsg {
		t.Errorf("Expected error message %s, got %s", expectedStatusMsg, statusOnlyErr.Error())
	}

	// Test field-specific errors without status code
	fieldErrNoStatus := &memory.Error{
		Errors: map[string][]string{
			"content": {"is invalid"},
		},
	}

	errorMsgNoStatus := fieldErrNoStatus.Error()
	expectedNoStatus := "API error: content is invalid"
	if errorMsgNoStatus != expectedNoStatus {
		t.Errorf("Expected error message %s, got %s", expectedNoStatus, errorMsgNoStatus)
	}
}

// TestMemoryCreatePromptWithFieldErrors tests creating a prompt with field validation errors.
func TestMemoryCreatePromptWithFieldErrors(t *testing.T) {
	// Test API response with field validation errors
	server := createMockServerForMemory(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/memory/spaces/space-123/prompts" {
			t.Errorf("Expected path /provision/memory/spaces/space-123/prompts, got %s", r.URL.Path)
		}

		// Return field validation errors
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := map[string]interface{}{
			"errors": map[string][]string{
				"name":    {"is required"},
				"content": {"is too short", "must be at least 10 characters"},
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	createReq := memory.CreatePromptRequest{
		Prompt: memory.PromptRequestData{
			Name:    "test-prompt",               // Valid name to bypass client validation
			Content: "Valid content for testing", // Valid content to bypass client validation
			Role:    "system",                    // Valid role to bypass client validation
		},
	}

	_, err := client.Memory.CreatePrompt("space-123", createReq)
	if err == nil {
		t.Fatal("Expected error for invalid prompt data")
	}

	// Check that the error contains field-specific messages
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "name is required") {
		t.Errorf("Expected error to contain 'name is required', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "content is too short") {
		t.Errorf("Expected error to contain 'content is too short', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "content must be at least 10 characters") {
		t.Errorf("Expected error to contain 'content must be at least 10 characters', got %s", errorMsg)
	}
}
