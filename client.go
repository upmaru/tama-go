package tama

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents the main Tama API client
type Client struct {
	httpClient *resty.Client
	baseURL    string
	apiKey     string
	Neural     *NeuralService
	Sensory    *SensoryService
}

// Config holds configuration options for the client
type Config struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// NewClient creates a new Tama API client
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
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

// SetAPIKey sets the API key for authentication
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
	c.httpClient.SetAuthToken(apiKey)
}

// SetDebug enables debug mode for HTTP requests
func (c *Client) SetDebug(debug bool) {
	c.httpClient.SetDebug(debug)
}

// Error represents an API error response
type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s - %s", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Response represents a standard API response wrapper
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}
