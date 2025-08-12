package memory_test

import (
	"net/http"
	"net/http/httptest"
)

// CreateMockServer creates a test HTTP server with the given handler.
func CreateMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}
