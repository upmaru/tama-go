package sensory_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/sensory"
)

func TestSensoryGetLimit(t *testing.T) {
	expectedLimit := sensory.Limit{
		ID:             "limit-123",
		SourceID:       "source-456",
		Count:          32,
		ScaleUnit:      "seconds",
		ScaleCount:     1,
		ProvisionState: "active",
	}

	expectedResponse := sensory.LimitResponse{
		Data: expectedLimit,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/limits/limit-123" {
			t.Errorf("Expected path /provision/sensory/limits/limit-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	client := tama.NewClient(tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	limit, err := client.Sensory.GetLimit("limit-123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if limit.ID != expectedLimit.ID {
		t.Errorf("Expected limit ID %s, got %s", expectedLimit.ID, limit.ID)
	}

	if limit.Count != expectedLimit.Count {
		t.Errorf("Expected count %d, got %d", expectedLimit.Count, limit.Count)
	}

	if limit.ScaleUnit != expectedLimit.ScaleUnit {
		t.Errorf("Expected scale unit %s, got %s", expectedLimit.ScaleUnit, limit.ScaleUnit)
	}

	if limit.ScaleCount != expectedLimit.ScaleCount {
		t.Errorf("Expected scale count %d, got %d", expectedLimit.ScaleCount, limit.ScaleCount)
	}
}

func TestSensoryCreateLimit(t *testing.T) {
	expectedLimit := sensory.Limit{
		ID:             "limit-789",
		SourceID:       "source-123",
		Count:          64,
		ScaleUnit:      "minutes",
		ScaleCount:     5,
		ProvisionState: "active",
	}

	expectedResponse := sensory.LimitResponse{
		Data: expectedLimit,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/sensory/sources/source-123/limits" {
			t.Errorf("Expected path /provision/sensory/sources/source-123/limits, got %s", r.URL.Path)
		}

		var req sensory.CreateLimitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Limit.Count != 64 {
			t.Errorf("Expected count 64, got %d", req.Limit.Count)
		}

		if req.Limit.ScaleUnit != "minutes" {
			t.Errorf("Expected request scale unit 'minutes', got %s", req.Limit.ScaleUnit)
		}

		if req.Limit.ScaleCount != 5 {
			t.Errorf("Expected request scale count 5, got %d", req.Limit.ScaleCount)
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

	createReq := sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      64,
			ScaleUnit:  "minutes",
			ScaleCount: 5,
		},
	}

	limit, err := client.Sensory.CreateLimit("source-123", createReq)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if limit.ID != expectedLimit.ID {
		t.Errorf("Expected limit ID %s, got %s", expectedLimit.ID, limit.ID)
	}
}

func TestSensoryCreateLimitValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	})

	// Test empty source ID validation
	_, err := client.Sensory.CreateLimit("", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      32,
			ScaleUnit:  "seconds",
			ScaleCount: 1,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty source ID")
	}

	// Test empty scale_unit validation
	_, err = client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      32,
			ScaleCount: 1,
		},
	})
	if err == nil {
		t.Error("Expected validation error for empty scale_unit")
	}

	// Test invalid scale_count validation
	_, err = client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      32,
			ScaleUnit:  "seconds",
			ScaleCount: 0,
		},
	})
	if err == nil {
		t.Error("Expected validation error for zero scale_count")
	}

	// Test invalid count value validation
	_, err = client.Sensory.CreateLimit("source-123", sensory.CreateLimitRequest{
		Limit: sensory.LimitRequestData{
			Count:      0,
			ScaleUnit:  "seconds",
			ScaleCount: 1,
		},
	})
	if err == nil {
		t.Error("Expected validation error for zero count value")
	}
}

func TestSensoryGetLimit_EmptyIDValidation(t *testing.T) {
	client := tama.NewClient(tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	})

	_, err := client.Sensory.GetLimit("")
	if err == nil {
		t.Error("Expected validation error for empty limit ID in GetLimit")
	}
}
