package tama_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// CreateMockServer creates a test HTTP server with the given handler.
// This helper is exported for use across all test files.
func CreateMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// Contains checks if a string contains a substring.
// This helper is exported for use across all test files.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
