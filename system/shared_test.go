package system_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/upmaru/tama-go/system"
)

// createSystemMockServer creates a test HTTP server with the given handler.
func createSystemMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// validateQueueResponse validates that actual queue matches expected queue.
func validateQueueResponse(t *testing.T, actual, expected system.Queue) {
	if actual.ID != expected.ID {
		t.Errorf("Expected queue ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.Role != expected.Role {
		t.Errorf("Expected queue role %s, got %s", expected.Role, actual.Role)
	}

	if actual.Name != expected.Name {
		t.Errorf("Expected queue name %s, got %s", expected.Name, actual.Name)
	}

	if actual.Concurrency != expected.Concurrency {
		t.Errorf("Expected queue concurrency %d, got %d", expected.Concurrency, actual.Concurrency)
	}
}
