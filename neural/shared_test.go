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
		BaseURL:        baseURL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true,
		Timeout:        10 * time.Second,
	}
	client, err := tama.NewClient(config)
	if err != nil {
		return nil
	}
	return client
}
