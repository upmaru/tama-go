package sensory

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// Service handles all sensory-related API operations
type Service struct {
	client *resty.Client
}

// NewService creates a new sensory service instance
func NewService(client *resty.Client) *Service {
	return &Service{
		client: client,
	}
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

// SourceCredential represents the credential structure for sources
type SourceCredential struct {
	ApiKey string `json:"api_key"`
}

// Source represents a sensory source resource
type Source struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// Model represents a sensory model resource
type Model struct {
	ID         string `json:"id,omitempty"`
	Identifier string `json:"identifier"`
}

// Limit represents a sensory limit resource
type Limit struct {
	ID           string `json:"id,omitempty"`
	Count        int    `json:"count"`
	ScaleUnit    string `json:"scale_unit"`
	ScaleCount   int    `json:"scale_count"`
	CurrentState string `json:"current_state"`
}

// SourceResponse represents the API response for source operations
type SourceResponse struct {
	Data Source `json:"data"`
}

// ModelResponse represents the API response for model operations
type ModelResponse struct {
	Data Model `json:"data"`
}

// LimitResponse represents the API response for limit operations
type LimitResponse struct {
	Data Limit `json:"data"`
}

// CreateSourceRequest represents the request payload for creating a source
type CreateSourceRequest struct {
	Source SourceRequestData `json:"source"`
}

// SourceRequestData represents the source data in the request
type SourceRequestData struct {
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Endpoint   string           `json:"endpoint"`
	Credential SourceCredential `json:"credential"`
}

// UpdateSourceRequest represents the request payload for updating a source
type UpdateSourceRequest struct {
	Source UpdateSourceData `json:"source"`
}

// UpdateSourceData represents the source update data
type UpdateSourceData struct {
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Endpoint   string            `json:"endpoint,omitempty"`
	Credential *SourceCredential `json:"credential,omitempty"`
}

// CreateModelRequest represents the request payload for creating a model
type CreateModelRequest struct {
	Model ModelRequestData `json:"model"`
}

// ModelRequestData represents the model data in the request
type ModelRequestData struct {
	Identifier string `json:"identifier"`
	Path       string `json:"path"`
}

// UpdateModelRequest represents the request payload for updating a model
type UpdateModelRequest struct {
	Model UpdateModelData `json:"model"`
}

// UpdateModelData represents the model update data
type UpdateModelData struct {
	Identifier string `json:"identifier,omitempty"`
	Path       string `json:"path,omitempty"`
}

// CreateLimitRequest represents the request payload for creating a limit
type CreateLimitRequest struct {
	Limit LimitRequestData `json:"limit"`
}

// LimitRequestData represents the limit data in the request
type LimitRequestData struct {
	ScaleUnit  string `json:"scale_unit"`
	ScaleCount int    `json:"scale_count"`
	Count      int    `json:"count"`
}

// UpdateLimitRequest represents the request payload for updating a limit
type UpdateLimitRequest struct {
	Limit UpdateLimitData `json:"limit"`
}

// UpdateLimitData represents the limit update data
type UpdateLimitData struct {
	ScaleUnit    string `json:"scale_unit,omitempty"`
	ScaleCount   int    `json:"scale_count,omitempty"`
	Count        int    `json:"count,omitempty"`
	CurrentState string `json:"current_state,omitempty"`
}
