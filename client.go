package tama

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultTimeout is the default timeout for API requests.
	DefaultTimeout = 30 * time.Second
)

// Client represents the main Tama API client.
type Client struct {
	httpClient *resty.Client
	baseURL    string
	apiKey     string
	Neural     *NeuralService
	Sensory    *SensoryService
}

// Config holds configuration options for the client.
type Config struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// NewClient creates a new Tama API client.
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	httpClient := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	if config.APIKey != "" {
		httpClient.SetAuthToken(config.APIKey)
	}

	client := &Client{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
	}

	// Initialize services
	client.Neural = newNeuralService(client)
	client.Sensory = newSensoryService(client)

	return client
}

// SetAPIKey sets the API key for authentication.
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
	c.httpClient.SetAuthToken(apiKey)
}

// SetDebug enables or disables debug mode for HTTP requests.
func (c *Client) SetDebug(debug bool) {
	c.httpClient.SetDebug(debug)
}

// SetHeader sets a header on the HTTP client.
func (c *Client) SetHeader(header, value string) {
	c.httpClient.SetHeader(header, value)
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
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}
