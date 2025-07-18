package sensory

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Service handles all sensory-related API operations.
type Service struct {
	client *resty.Client
}

// NewService creates a new sensory service instance.
func NewService(client *resty.Client) *Service {
	return &Service{
		client: client,
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

// SourceCredential represents the credential structure for sources.
type SourceCredential struct {
	APIKey string `json:"api_key"`
}

// Source represents a sensory source resource.
type Source struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	Endpoint     string `json:"endpoint"`
	SpaceID      string `json:"space_id"`
	CurrentState string `json:"current_state"`
}

// Model represents a sensory model resource.
type Model struct {
	ID           string         `json:"id,omitempty"`
	Identifier   string         `json:"identifier"`
	Path         string         `json:"path"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	CurrentState string         `json:"current_state"`
}

// Limit represents a sensory limit resource.
type Limit struct {
	ID           string `json:"id,omitempty"`
	SourceID     string `json:"source_id"`
	Count        int    `json:"count"`
	ScaleUnit    string `json:"scale_unit"`
	ScaleCount   int    `json:"scale_count"`
	CurrentState string `json:"current_state"`
}

// SourceResponse represents the API response for source operations.
type SourceResponse struct {
	Data Source `json:"data"`
}

// ModelResponse represents the API response for model operations.
type ModelResponse struct {
	Data Model `json:"data"`
}

// LimitResponse represents the API response for limit operations.
type LimitResponse struct {
	Data Limit `json:"data"`
}

// CreateSourceRequest represents the request payload for creating a source.
type CreateSourceRequest struct {
	Source SourceRequestData `json:"source"`
}

// SourceRequestData represents the source data in the request.
type SourceRequestData struct {
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Endpoint   string           `json:"endpoint"`
	Credential SourceCredential `json:"credential"`
}

// UpdateSourceRequest represents the request payload for updating a source.
type UpdateSourceRequest struct {
	Source UpdateSourceData `json:"source"`
}

// UpdateSourceData represents the source update data.
type UpdateSourceData struct {
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Endpoint   string            `json:"endpoint,omitempty"`
	Credential *SourceCredential `json:"credential,omitempty"`
}

// CreateModelRequest represents the request payload for creating a model.
type CreateModelRequest struct {
	Model ModelRequestData `json:"model"`
}

// ModelRequestData represents the model data in the request.
type ModelRequestData struct {
	Identifier string         `json:"identifier"`
	Path       string         `json:"path"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// UpdateModelRequest represents the request payload for updating a model.
type UpdateModelRequest struct {
	Model UpdateModelData `json:"model"`
}

// UpdateModelData represents the model update data.
type UpdateModelData struct {
	Identifier string         `json:"identifier,omitempty"`
	Path       string         `json:"path,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// CreateLimitRequest represents the request payload for creating a limit.
type CreateLimitRequest struct {
	Limit LimitRequestData `json:"limit"`
}

// LimitRequestData represents the limit data in the request.
type LimitRequestData struct {
	ScaleUnit  string `json:"scale_unit"`
	ScaleCount int    `json:"scale_count"`
	Count      int    `json:"count"`
}

// UpdateLimitRequest represents the request payload for updating a limit.
type UpdateLimitRequest struct {
	Limit UpdateLimitData `json:"limit"`
}

// UpdateLimitData represents the limit update data.
type UpdateLimitData struct {
	ScaleUnit    string `json:"scale_unit,omitempty"`
	ScaleCount   int    `json:"scale_count,omitempty"`
	Count        int    `json:"count,omitempty"`
	CurrentState string `json:"current_state,omitempty"`
}

// handleAPIError processes API error responses.
func (s *Service) handleAPIError(resp interface{}) error {
	errResp, ok := s.extractErrorResponse(resp)
	if !ok {
		return nil
	}

	if body := errResp.Body(); len(body) > 0 {
		if err := s.parseErrorFromBody(body, errResp.StatusCode()); err != nil {
			return err
		}
	}

	return s.fallbackError(errResp)
}

// extractErrorResponse extracts error response interface from resp.
func (s *Service) extractErrorResponse(resp interface{}) (errorResponse, bool) {
	type errorResponse interface {
		IsError() bool
		Error() interface{}
		StatusCode() int
		Status() string
		Body() []byte
	}

	if errResp, ok := resp.(errorResponse); ok && errResp.IsError() {
		return errResp, true
	}
	return nil, false
}

// parseErrorFromBody attempts to parse error from response body.
func (s *Service) parseErrorFromBody(body []byte, statusCode int) error {
	// Try to parse as map[string][]string (array format)
	if err := s.parseArrayError(body, statusCode); err != nil {
		return err
	}

	// Try to parse as map[string]string (single string format)
	return s.parseStringError(body, statusCode)
}

// parseArrayError parses errors in array format.
func (s *Service) parseArrayError(body []byte, statusCode int) error {
	var rawArrayError struct {
		Errors map[string][]string `json:"errors"`
	}

	if err := json.Unmarshal(body, &rawArrayError); err == nil && rawArrayError.Errors != nil {
		return &Error{
			StatusCode: statusCode,
			Errors:     rawArrayError.Errors,
		}
	}
	return nil
}

// parseStringError parses errors in string format and converts to array format.
func (s *Service) parseStringError(body []byte, statusCode int) error {
	var rawStringError struct {
		Errors map[string]string `json:"errors"`
	}

	if err := json.Unmarshal(body, &rawStringError); err == nil && rawStringError.Errors != nil {
		convertedErrors := make(map[string][]string)
		for field, message := range rawStringError.Errors {
			convertedErrors[field] = []string{message}
		}
		return &Error{
			StatusCode: statusCode,
			Errors:     convertedErrors,
		}
	}
	return nil
}

// fallbackError handles fallback error cases.
func (s *Service) fallbackError(errResp errorResponse) error {
	if apiErrorResp, isError := errResp.Error().(*Error); isError {
		apiErrorResp.StatusCode = errResp.StatusCode()
		return apiErrorResp
	}
	return fmt.Errorf("API error: %s", errResp.Status())
}

// errorResponse interface for type assertion.
type errorResponse interface {
	IsError() bool
	Error() interface{}
	StatusCode() int
	Status() string
	Body() []byte
}
