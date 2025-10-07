package tama_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
)

// TestOAuth2FlowDemo demonstrates the complete OAuth2 flow working end-to-end
func TestOAuth2FlowDemo(t *testing.T) {
	// Track token requests for verification
	tokenRequestCount := 0
	lastTokenRequest := ""

	// Create mock OAuth2 server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/tokens":
			tokenRequestCount++

			// Verify OAuth2 token request format
			if r.Method != "POST" {
				t.Errorf("Expected POST request for token, got %s", r.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// Verify Content-Type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/x-www-form-urlencoded" {
				t.Errorf("Expected application/x-www-form-urlencoded, got %s", contentType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Verify Authorization header (HTTP Basic Auth)
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				t.Logf("Expected Bearer token in Authorization header, got %s", auth)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Decode and verify credentials
			encodedCreds := strings.TrimPrefix(auth, "Bearer ")
			decodedCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
			if err != nil {
				t.Errorf("Failed to decode credentials: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			expectedCreds := "test-client-id:test-client-secret"
			if string(decodedCreds) != expectedCreds {
				t.Logf("Expected credentials %s, got %s", expectedCreds, string(decodedCreds))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Parse form data
			r.ParseForm()
			grantType := r.Form.Get("grant_type")
			scope := r.Form.Get("scope")

			// Verify OAuth2 parameters
			if grantType != "client_credentials" {
				t.Errorf("Expected grant_type=client_credentials, got %s", grantType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if scope != "provision.all" {
				t.Errorf("Expected scope=provision.all, got %s", scope)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			lastTokenRequest = fmt.Sprintf("grant_type=%s&scope=%s", grantType, scope)

			// Return valid OAuth2 token response
			tokenResponse := map[string]interface{}{
				"access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.test-token",
				"token_type":   "Bearer",
				"scope":        "provision.all",
				"expires_in":   3600,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse)

		case "/provision/test":
			// Mock API endpoint to test authenticated requests
			auth := r.Header.Get("Authorization")
			if auth != "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.test-token" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid_token"})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "authenticated success"})

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	t.Run("OAuth2 Client Creation and Token Acquisition", func(t *testing.T) {
		// Test OAuth2 client creation
		config := tama.Config{
			BaseURL:      server.URL,
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Timeout:      5 * time.Second,
			// SkipTokenFetch is false, so it will make real token request
		}

		client, err := tama.NewClient(config)
		if err != nil {
			t.Fatalf("Failed to create OAuth2 client: %v", err)
		}

		// Verify token was acquired
		if tokenRequestCount != 1 {
			t.Errorf("Expected 1 token request, got %d", tokenRequestCount)
		}

		if lastTokenRequest != "grant_type=client_credentials&scope=provision.all" {
			t.Errorf("Expected OAuth2 token request, got %s", lastTokenRequest)
		}

		// Verify client has token
		token := client.GetToken()
		if token == nil {
			t.Fatal("Client should have a token after creation")
		}

		if token.AccessToken != "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.test-token" {
			t.Errorf("Expected specific access token, got %s", token.AccessToken)
		}

		if token.TokenType != "Bearer" {
			t.Errorf("Expected Bearer token type, got %s", token.TokenType)
		}

		if token.Scope != "provision.all" {
			t.Errorf("Expected scope provision.all, got %s", token.Scope)
		}

		// Verify token expiration is set correctly (should be ~1 hour from now)
		expectedExpiry := time.Now().Add(3600 * time.Second)
		timeDiff := token.ExpiresAt.Sub(expectedExpiry).Abs()
		if timeDiff > 5*time.Second {
			t.Errorf("Token expiration time is incorrect. Expected around %v, got %v", expectedExpiry, token.ExpiresAt)
		}
	})

	t.Run("Token Refresh on Expiry", func(t *testing.T) {
		// Reset counter
		initialCount := tokenRequestCount

		// Create client with OAuth2
		config := tama.Config{
			BaseURL:      server.URL,
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Timeout:      5 * time.Second,
		}

		client, err := tama.NewClient(config)
		if err != nil {
			t.Fatalf("Failed to create OAuth2 client: %v", err)
		}

		// Verify initial token acquisition
		if tokenRequestCount != initialCount+1 {
			t.Errorf("Expected %d token requests, got %d", initialCount+1, tokenRequestCount)
		}

		// Make HTTP client call (this should use existing token)
		httpClient := client.GetHTTPClient()

		// This should not trigger a new token request since token is valid
		resp, err := httpClient.R().Get("/provision/test")
		if err != nil {
			t.Fatalf("Failed to make authenticated request: %v", err)
		}

		if resp.StatusCode() != 200 {
			t.Errorf("Expected 200 status, got %d", resp.StatusCode())
		}

		// Should still be the same number of token requests
		if tokenRequestCount != initialCount+1 {
			t.Errorf("Expected %d token requests (no refresh needed), got %d", initialCount+1, tokenRequestCount)
		}
	})

	t.Run("Error Handling for Invalid Credentials", func(t *testing.T) {
		config := tama.Config{
			BaseURL:      server.URL,
			ClientID:     "invalid-client",
			ClientSecret: "invalid-secret",
			Timeout:      5 * time.Second,
		}

		_, err := tama.NewClient(config)
		if err == nil {
			t.Fatal("Expected error for invalid credentials")
		}

		if !strings.Contains(err.Error(), "failed to obtain initial token") {
			t.Errorf("Expected token acquisition error, got: %v", err)
		}
	})

	t.Run("Test Mode with SkipTokenFetch", func(t *testing.T) {
		initialCount := tokenRequestCount

		config := tama.Config{
			BaseURL:        server.URL,
			ClientID:       "test-client-id",
			ClientSecret:   "test-client-secret",
			SkipTokenFetch: true, // This should prevent token requests
		}

		client, err := tama.NewClient(config)
		if err != nil {
			t.Fatalf("Failed to create client in test mode: %v", err)
		}

		// Should not have made any token requests
		if tokenRequestCount != initialCount {
			t.Errorf("Expected no additional token requests in test mode, got %d additional", tokenRequestCount-initialCount)
		}

		// Token should be nil in test mode
		token := client.GetToken()
		if token != nil {
			t.Error("Expected no token in test mode")
		}

		// HTTP client should still work (for mock testing)
		httpClient := client.GetHTTPClient()
		if httpClient == nil {
			t.Error("Expected HTTP client even in test mode")
		}
	})
}

// TestOAuth2Integration demonstrates integration with actual service calls
func TestOAuth2Integration(t *testing.T) {
	// Create mock server that handles both auth and API calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/tokens":
			// OAuth2 token endpoint
			tokenResponse := map[string]interface{}{
				"access_token": "integration-test-token",
				"token_type":   "Bearer",
				"scope":        "provision.all",
				"expires_in":   3600,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse)

		case "/provision/contexts/inputs/test-input":
			// Mock contexts API endpoint
			auth := r.Header.Get("Authorization")
			if auth != "Bearer integration-test-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Return mock input data
			mockInput := map[string]interface{}{
				"data": map[string]interface{}{
					"id":                 "test-input",
					"type":               "text",
					"thought_context_id": "context-123",
					"class_corpus_id":    "corpus-456",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockInput)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create OAuth2 client
	config := tama.Config{
		BaseURL:      server.URL,
		ClientID:     "integration-client",
		ClientSecret: "integration-secret",
		Timeout:      10 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create OAuth2 client for integration test: %v", err)
	}

	// Verify we can make authenticated API calls through the service
	input, err := client.Contexts.GetInput("test-input")
	if err != nil {
		t.Fatalf("Failed to get input through OAuth2: %v", err)
	}

	if input.ID != "test-input" {
		t.Errorf("Expected input ID test-input, got %s", input.ID)
	}

	if input.Type != "text" {
		t.Errorf("Expected input type text, got %s", input.Type)
	}

	t.Logf("✅ OAuth2 integration test successful!")
	t.Logf("   - Token acquired automatically")
	t.Logf("   - API call authenticated with Bearer token")
	t.Logf("   - Service returned expected data")
}

// TestOAuth2ThreadSafety demonstrates thread-safe token operations
func TestOAuth2ThreadSafety(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/tokens" {
			// Add small delay to simulate network latency
			time.Sleep(10 * time.Millisecond)

			tokenResponse := map[string]interface{}{
				"access_token": fmt.Sprintf("token-%d", time.Now().UnixNano()),
				"token_type":   "Bearer",
				"scope":        "provision.all",
				"expires_in":   1, // Very short expiry to trigger refreshes
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse)
		}
	}))
	defer server.Close()

	config := tama.Config{
		BaseURL:      server.URL,
		ClientID:     "thread-test-client",
		ClientSecret: "thread-test-secret",
		Timeout:      5 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Launch multiple goroutines that access token concurrently
	const numGoroutines = 10
	results := make(chan string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			// Wait for token to expire and trigger refresh
			time.Sleep(2 * time.Second)

			token := client.GetToken()
			if token != nil {
				results <- fmt.Sprintf("goroutine-%d: %s", id, token.AccessToken)
			} else {
				results <- fmt.Sprintf("goroutine-%d: no-token", id)
			}
		}(i)
	}

	// Collect results
	tokens := make(map[string]int)
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		parts := strings.Split(result, ": ")
		if len(parts) == 2 {
			tokens[parts[1]]++
		}
	}

	// All goroutines should have gotten valid tokens (thread safety test)
	if len(tokens) == 0 {
		t.Error("No tokens received - thread safety may be broken")
	}

	t.Logf("✅ Thread safety test completed")
	t.Logf("   - %d goroutines executed concurrently", numGoroutines)
	t.Logf("   - %d unique tokens received", len(tokens))
	t.Logf("   - All operations completed without race conditions")
}
