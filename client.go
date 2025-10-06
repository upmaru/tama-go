package tama

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultTimeout is the default timeout for API requests.
	DefaultTimeout = 30 * time.Second
)

// TokenResponse represents the OAuth2 token response from /auth/tokens
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Token represents an OAuth2 access token with expiration tracking
type Token struct {
	AccessToken string
	TokenType   string
	Scope       string
	ExpiresAt   time.Time
}

// IsExpired checks if the token is expired or will expire within 30 seconds
func (t *Token) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return true
	}
	// Consider token expired if it expires within 30 seconds to account for request time
	return time.Now().Add(30 * time.Second).After(t.ExpiresAt)
}

// Client represents the main Tama API client.
type Client struct {
	httpClient     *resty.Client
	baseURL        string
	clientID       string
	clientSecret   string
	token          *Token
	tokenMutex     sync.RWMutex
	skipTokenFetch bool
	Neural         *NeuralService
	Sensory        *SensoryService
	Memory         *MemoryService
	Perception     *PerceptionService
	Motor          *MotorService
	Contexts       *ContextsService
}

// Config holds configuration options for the client.
type Config struct {
	BaseURL        string
	ClientID       string
	ClientSecret   string
	Timeout        time.Duration
	SkipTokenFetch bool // For testing - skips initial token fetch
}

// NewClient creates a new Tama API client with OAuth2 authentication.
func NewClient(config Config) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	if config.ClientID == "" {
		return nil, fmt.Errorf("client_id is required")
	}

	if config.ClientSecret == "" {
		return nil, fmt.Errorf("client_secret is required")
	}

	httpClient := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	client := &Client{
		httpClient:     httpClient,
		baseURL:        config.BaseURL,
		clientID:       config.ClientID,
		clientSecret:   config.ClientSecret,
		skipTokenFetch: config.SkipTokenFetch,
	}

	// Get initial OAuth2 token (unless skipped for testing)
	if !config.SkipTokenFetch {
		if err := client.refreshToken(); err != nil {
			return nil, fmt.Errorf("failed to obtain initial token: %w", err)
		}
	}

	// Set up request interceptor to ensure token validity on each request
	httpClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Skip token validation if SkipTokenFetch is enabled (for testing)
		if client.skipTokenFetch {
			return nil
		}
		// Ensure valid token before each request
		if err := client.ensureValidToken(); err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		return nil
	})

	// Initialize services
	client.Neural = newNeuralService(client)
	client.Sensory = newSensoryService(client)
	client.Memory = newMemoryService(client)
	client.Perception = newPerceptionService(client)
	client.Motor = newMotorService(client)
	client.Contexts = newContextsService(client)

	return client, nil
}

// refreshToken obtains a new OAuth2 access token using client credentials flow
func (c *Client) refreshToken() error {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Create Basic Auth credentials: base64(client_id:client_secret)
	credentials := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))

	// Create a separate HTTP client for token requests to avoid circular dependency
	tokenClient := resty.New().
		SetBaseURL(c.baseURL).
		SetTimeout(c.httpClient.GetClient().Timeout).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+credentials)

	var tokenResponse TokenResponse
	resp, err := tokenClient.R().
		SetFormData(map[string]string{
			"grant_type": "client_credentials",
			"scope":      "provision.all",
		}).
		SetResult(&tokenResponse).
		Post("/auth/tokens")

	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	// Store the new token
	c.token = &Token{
		AccessToken: tokenResponse.AccessToken,
		TokenType:   tokenResponse.TokenType,
		Scope:       tokenResponse.Scope,
		ExpiresAt:   expiresAt,
	}

	// Update the HTTP client with the new token
	c.httpClient.SetAuthToken(c.token.AccessToken)

	return nil
}

// ensureValidToken ensures that we have a valid, non-expired token
func (c *Client) ensureValidToken() error {
	// Skip token validation if SkipTokenFetch is enabled (for testing)
	if c.skipTokenFetch {
		return nil
	}

	c.tokenMutex.RLock()
	if c.token != nil && !c.token.IsExpired() {
		c.tokenMutex.RUnlock()
		return nil
	}
	c.tokenMutex.RUnlock()

	// Token is expired or missing, refresh it
	return c.refreshToken()
}

// SetDebug enables or disables debug mode for HTTP requests.
func (c *Client) SetDebug(debug bool) {
	c.httpClient.SetDebug(debug)
}

// SetHeader sets a header on the HTTP client.
func (c *Client) SetHeader(header, value string) {
	c.httpClient.SetHeader(header, value)
}

// GetHTTPClient returns the underlying HTTP client for use by services.
func (c *Client) GetHTTPClient() *resty.Client {
	// Ensure we have a valid token before returning the client
	// If token refresh fails, we'll let the individual API calls handle the error
	c.ensureValidToken()
	return c.httpClient
}

// TestOAuth2Authentication tests actual OAuth2 authentication flow
// This test requires a real server with valid client credentials
// Run with: go test -run TestOAuth2Authentication -tags=oauth2
func TestOAuth2Authentication(clientID, clientSecret, baseURL string) error {
	config := Config{
		BaseURL:      baseURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Timeout:      30 * time.Second,
		// SkipTokenFetch is false by default, so it will attempt real authentication
	}

	client, err := NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create OAuth2 client: %w", err)
	}

	// Test that we have a valid token
	token := client.GetToken()
	if token == nil {
		return fmt.Errorf("no token available after authentication")
	}

	if token.AccessToken == "" {
		return fmt.Errorf("empty access token")
	}

	if token.TokenType != "Bearer" {
		return fmt.Errorf("expected Bearer token type, got %s", token.TokenType)
	}

	return nil
}

// GetToken returns the current access token information
func (c *Client) GetToken() *Token {
	c.tokenMutex.RLock()
	defer c.tokenMutex.RUnlock()

	if c.token == nil {
		return nil
	}

	// Return a copy to prevent external modification
	return &Token{
		AccessToken: c.token.AccessToken,
		TokenType:   c.token.TokenType,
		Scope:       c.token.Scope,
		ExpiresAt:   c.token.ExpiresAt,
	}
}

// Error represents an API error response.
type Error struct {
	StatusCode int                 `json:"status_code"`
	Errors     map[string][]string `json:"errors,omitempty"`
}

func (e *Error) Error() string {
	if len(e.Errors) > 0 {
		var errorParts []string
		for field, messages := range e.Errors {
			for _, message := range messages {
				errorParts = append(errorParts, fmt.Sprintf("%s %s", field, message))
			}
		}
		if e.StatusCode > 0 {
			return fmt.Sprintf("API error %d: %s", e.StatusCode, strings.Join(errorParts, ", "))
		}
		return fmt.Sprintf("API error: %s", strings.Join(errorParts, ", "))
	}

	if e.StatusCode > 0 {
		return fmt.Sprintf("API error %d", e.StatusCode)
	}
	return "API error"
}

// Response represents a standard API response wrapper.
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}
