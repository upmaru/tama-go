package memory_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/memory"
)

func TestMemoryGetTopic(t *testing.T) {
	expected := memory.Topic{
		ID:             "topic-123",
		ListenerID:     "listener-123",
		ClassID:        "class-123",
		ProvisionState: "active",
	}
	expectedResp := memory.TopicResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/memory/topics/topic-123" {
			t.Errorf("Expected path /provision/memory/topics/topic-123, got %s", r.URL.Path)
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

	topic, err := client.Memory.GetTopic("topic-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if topic.ID != expected.ID {
		t.Errorf("Expected topic ID %s, got %s", expected.ID, topic.ID)
	}
	if topic.ListenerID != expected.ListenerID {
		t.Errorf("Expected listener ID %s, got %s", expected.ListenerID, topic.ListenerID)
	}
	if topic.ClassID != expected.ClassID {
		t.Errorf("Expected class ID %s, got %s", expected.ClassID, topic.ClassID)
	}
	if topic.ProvisionState != expected.ProvisionState {
		t.Errorf("Expected provision state %s, got %s", expected.ProvisionState, topic.ProvisionState)
	}
}

func TestMemoryGetTopicError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := memory.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"topic": {"not found"},
			},
		}
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

	_, err = client.Memory.GetTopic("missing")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	var merr *memory.Error
	if errors.As(err, &merr) {
		if merr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", merr.StatusCode)
		}
		if merr.Errors == nil || len(merr.Errors["topic"]) == 0 || merr.Errors["topic"][0] != "not found" {
			t.Errorf("Expected 'topic not found', got %v", merr.Errors)
		}
	} else {
		t.Errorf("Expected memory.Error, got %T", err)
	}
}

func TestMemoryCreateTopic(t *testing.T) {
	expected := memory.Topic{
		ID:             "topic-789",
		ListenerID:     "listener-123",
		ClassID:        "class-456",
		ProvisionState: "pending",
	}
	expectedResp := memory.TopicResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/memory/listeners/listener-123/topics" {
			t.Errorf("Expected path /provision/memory/listeners/listener-123/topics, got %s", r.URL.Path)
		}
		var req memory.CreateTopicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Topic.ClassID != "class-456" {
			t.Errorf("Expected class_id 'class-456', got %s", req.Topic.ClassID)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
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

	createReq := memory.CreateTopicRequest{Topic: memory.TopicRequestData{ClassID: "class-456"}}
	topic, err := client.Memory.CreateTopic("listener-123", createReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if topic.ID != expected.ID {
		t.Errorf("Expected topic ID %s, got %s", expected.ID, topic.ID)
	}
	if topic.ClassID != expected.ClassID {
		t.Errorf("Expected class ID %s, got %s", expected.ClassID, topic.ClassID)
	}
}

func TestMemoryCreateTopicValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Empty listener ID
	_, err = client.Memory.CreateTopic(
		"",
		memory.CreateTopicRequest{
			Topic: memory.TopicRequestData{
				ClassID: "class-123",
			},
		},
	)
	if err == nil {
		t.Error("Expected validation error for empty listener ID")
	}
	// Empty class ID
	_, err = client.Memory.CreateTopic("listener-123", memory.CreateTopicRequest{Topic: memory.TopicRequestData{}})
	if err == nil {
		t.Error("Expected validation error for empty class ID")
	}
}

func TestMemoryUpdateTopic(t *testing.T) {
	expected := memory.Topic{
		ID:             "topic-123",
		ListenerID:     "listener-123",
		ClassID:        "class-999",
		ProvisionState: "active",
	}
	expectedResp := memory.TopicResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/memory/topics/topic-123" {
			t.Errorf("Expected path /provision/memory/topics/topic-123, got %s", r.URL.Path)
		}
		var req memory.UpdateTopicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Topic.ClassID != "class-999" {
			t.Errorf("Expected class_id 'class-999', got %s", req.Topic.ClassID)
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

	updateReq := memory.UpdateTopicRequest{Topic: memory.UpdateTopicData{ClassID: "class-999"}}
	topic, err := client.Memory.UpdateTopic("topic-123", updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if topic.ClassID != expected.ClassID {
		t.Errorf("Expected class ID %s, got %s", expected.ClassID, topic.ClassID)
	}
}

func TestMemoryUpdateTopicValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Empty topic ID
	_, err = client.Memory.UpdateTopic("", memory.UpdateTopicRequest{Topic: memory.UpdateTopicData{ClassID: "x"}})
	if err == nil {
		t.Error("Expected validation error for empty topic ID")
	}
	// Empty class ID
	_, err = client.Memory.UpdateTopic("topic-123", memory.UpdateTopicRequest{Topic: memory.UpdateTopicData{}})
	if err == nil {
		t.Error("Expected validation error for empty class ID")
	}
}

func TestMemoryReplaceTopic(t *testing.T) {
	expected := memory.Topic{
		ID:             "topic-123",
		ListenerID:     "listener-123",
		ClassID:        "class-555",
		ProvisionState: "active",
	}
	expectedResp := memory.TopicResponse{Data: expected}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/memory/topics/topic-123" {
			t.Errorf("Expected path /provision/memory/topics/topic-123, got %s", r.URL.Path)
		}
		var req memory.UpdateTopicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Topic.ClassID != "class-555" {
			t.Errorf("Expected class_id 'class-555', got %s", req.Topic.ClassID)
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

	replaceReq := memory.UpdateTopicRequest{Topic: memory.UpdateTopicData{ClassID: "class-555"}}
	topic, err := client.Memory.ReplaceTopic("topic-123", replaceReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if topic.ClassID != expected.ClassID {
		t.Errorf("Expected class ID %s, got %s", expected.ClassID, topic.ClassID)
	}
}

func TestMemoryReplaceTopicValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	// Empty topic ID
	_, err = client.Memory.ReplaceTopic("", memory.UpdateTopicRequest{Topic: memory.UpdateTopicData{ClassID: "x"}})
	if err == nil {
		t.Error("Expected validation error for empty topic ID")
	}
	// Empty class ID
	_, err = client.Memory.ReplaceTopic("topic-123", memory.UpdateTopicRequest{Topic: memory.UpdateTopicData{}})
	if err == nil {
		t.Error("Expected validation error for empty class ID")
	}
}

func TestMemoryDeleteTopic(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/provision/memory/topics/topic-123" {
			t.Errorf("Expected path /provision/memory/topics/topic-123, got %s", r.URL.Path)
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

	err = client.Memory.DeleteTopic("topic-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMemoryDeleteTopicValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		Timeout:        10 * time.Second,
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}
	err = client.Memory.DeleteTopic("")
	if err == nil {
		t.Error("Expected validation error for empty topic ID")
	}
}
