package system_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/system"
)

func TestSystemGetQueue(t *testing.T) {
	expectedQueue := system.Queue{
		ID:          "queue-123",
		Role:        "translator",
		Name:        "primary",
		Concurrency: 5,
	}

	expectedResponse := system.QueueResponse{
		Data: expectedQueue,
	}

	server := createSystemMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/system/queues/queue-123" {
			t.Errorf("Expected path /provision/system/queues/queue-123, got %s", r.URL.Path)
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

	queue, err := client.System.GetQueue("queue-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateQueueResponse(t, *queue, expectedQueue)
}

func TestSystemCreateQueue(t *testing.T) {
	request := system.CreateQueueRequest{
		Queue: system.QueueRequestData{
			Role:        "planner",
			Name:        "analysis",
			Concurrency: 3,
		},
	}

	expectedQueue := system.Queue{
		ID:          "queue-456",
		Role:        "planner",
		Name:        "analysis",
		Concurrency: 3,
	}

	expectedResponse := system.QueueResponse{
		Data: expectedQueue,
	}

	server := createSystemMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/system/queues" {
			t.Errorf("Expected path /provision/system/queues, got %s", r.URL.Path)
		}

		var receivedRequest system.CreateQueueRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Queue.Role != request.Queue.Role {
			t.Errorf("Expected queue role %s, got %s", request.Queue.Role, receivedRequest.Queue.Role)
		}
		if receivedRequest.Queue.Name != request.Queue.Name {
			t.Errorf("Expected queue name %s, got %s", request.Queue.Name, receivedRequest.Queue.Name)
		}
		if receivedRequest.Queue.Concurrency != request.Queue.Concurrency {
			t.Errorf(
				"Expected queue concurrency %d, got %d",
				request.Queue.Concurrency,
				receivedRequest.Queue.Concurrency,
			)
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

	queue, err := client.System.CreateQueue(request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateQueueResponse(t, *queue, expectedQueue)
}

func TestSystemCreateQueueValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.System.CreateQueue(system.CreateQueueRequest{
		Queue: system.QueueRequestData{
			Name:        "analysis",
			Concurrency: 2,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty queue role")
	}

	_, err = client.System.CreateQueue(system.CreateQueueRequest{
		Queue: system.QueueRequestData{
			Role:        "planner",
			Concurrency: 2,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty queue name")
	}

	_, err = client.System.CreateQueue(system.CreateQueueRequest{
		Queue: system.QueueRequestData{
			Role:        "planner",
			Name:        "analysis",
			Concurrency: 0,
		},
	})
	if err == nil {
		t.Error("Expected validation error for invalid queue concurrency")
	}
}

func TestSystemUpdateQueue(t *testing.T) {
	newConcurrency := 10
	request := system.UpdateQueueRequest{
		Queue: system.UpdateQueueData{
			Role:        "planner",
			Name:        "critical",
			Concurrency: &newConcurrency,
		},
	}

	expectedQueue := system.Queue{
		ID:          "queue-123",
		Role:        "planner",
		Name:        "critical",
		Concurrency: 10,
	}

	expectedResponse := system.QueueResponse{
		Data: expectedQueue,
	}

	server := createSystemMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/system/queues/queue-123" {
			t.Errorf("Expected path /provision/system/queues/queue-123, got %s", r.URL.Path)
		}

		var receivedRequest system.UpdateQueueRequest
		if err := json.NewDecoder(r.Body).Decode(&receivedRequest); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedRequest.Queue.Role != request.Queue.Role {
			t.Errorf("Expected queue role %s, got %s", request.Queue.Role, receivedRequest.Queue.Role)
		}
		if receivedRequest.Queue.Name != request.Queue.Name {
			t.Errorf("Expected queue name %s, got %s", request.Queue.Name, receivedRequest.Queue.Name)
		}
		if receivedRequest.Queue.Concurrency == nil {
			t.Fatal("Expected queue concurrency to be present")
		}
		if *receivedRequest.Queue.Concurrency != *request.Queue.Concurrency {
			t.Errorf(
				"Expected queue concurrency %d, got %d",
				*request.Queue.Concurrency,
				*receivedRequest.Queue.Concurrency,
			)
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

	queue, err := client.System.UpdateQueue("queue-123", request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateQueueResponse(t, *queue, expectedQueue)
}

func TestSystemUpdateQueueValidation(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	_, err = client.System.UpdateQueue("", system.UpdateQueueRequest{
		Queue: system.UpdateQueueData{
			Role: "planner",
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty queue ID")
	}

	invalidConcurrency := 0
	_, err = client.System.UpdateQueue("queue-123", system.UpdateQueueRequest{
		Queue: system.UpdateQueueData{
			Concurrency: &invalidConcurrency,
		},
	})
	if err == nil {
		t.Error("Expected validation error for invalid queue concurrency")
	}
}

func TestSystemDeleteQueue(t *testing.T) {
	server := createSystemMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/system/queues/queue-123" {
			t.Errorf("Expected path /provision/system/queues/queue-123, got %s", r.URL.Path)
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

	err = client.System.DeleteQueue("queue-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSystemDeleteQueueEmptyID(t *testing.T) {
	client, err := tama.NewClient(tama.Config{
		BaseURL:        "https://api.example.com",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
	})
	if err != nil {
		t.Skipf("Skipping test due to client creation failure: %v", err)
	}

	err = client.System.DeleteQueue("")
	if err == nil {
		t.Error("Expected validation error for empty queue ID")
	}
}
