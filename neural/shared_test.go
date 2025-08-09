package neural_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	tama "github.com/upmaru/tama-go"
)

// CreateMockServer creates a test HTTP server with the given handler.
func CreateMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// createTestClient is a helper to create a test client with a given base URL.
func createTestClient(baseURL string) *tama.Client {
	config := tama.Config{
		BaseURL: baseURL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}
	return tama.NewClient(config)
}
