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
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Endpoint       string `json:"endpoint"`
	SpaceID        string `json:"space_id"`
	Type           string `json:"type"`
	ProvisionState string `json:"provision_state"`
}

// Model represents a sensory model resource.
type Model struct {
	ID             string         `json:"id,omitempty"`
	Identifier     string         `json:"identifier"`
	Path           string         `json:"path"`
	Parameters     map[string]any `json:"parameters,omitempty"`
	ProvisionState string         `json:"provision_state"`
}

// Limit represents a sensory limit resource.
type Limit struct {
	ID             string `json:"id,omitempty"`
	SourceID       string `json:"source_id"`
	Count          int    `json:"count"`
	ScaleUnit      string `json:"scale_unit"`
	ScaleCount     int    `json:"scale_count"`
	ProvisionState string `json:"provision_state"`
}

// Specification represents a sensory specification resource.
type Specification struct {
	ID             string         `json:"id,omitempty"`
	SpaceID        string         `json:"space_id"`
	Schema         map[string]any `json:"schema"`
	Version        string         `json:"version"`
	Endpoint       string         `json:"endpoint"`
	CurrentState   string         `json:"current_state"`
	ProvisionState string         `json:"provision_state"`
}

// Validation represents the validation configuration for an identity.
type Validation struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Codes  []int  `json:"codes"`
}

// Identity represents a sensory identity resource.
type Identity struct {
	ID              string     `json:"id,omitempty"`
	SpecificationID string     `json:"specification_id"`
	ProvisionState  string     `json:"provision_state"`
	CurrentState    string     `json:"current_state"`
	Identifier      string     `json:"identifier"`
	Validation      Validation `json:"validation"`
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

// SpecificationResponse represents the API response for specification operations.
type SpecificationResponse struct {
	Data Specification `json:"data"`
}

// IdentityResponse represents the API response for identity operations.
type IdentityResponse struct {
	Data Identity `json:"data"`
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
	ScaleUnit      string `json:"scale_unit,omitempty"`
	ScaleCount     int    `json:"scale_count,omitempty"`
	Count          int    `json:"count,omitempty"`
	ProvisionState string `json:"provision_state,omitempty"`
}

// CreateSpecificationRequest represents the request payload for creating a specification.
type CreateSpecificationRequest struct {
	Specification SpecificationRequestData `json:"specification"`
}

// SpecificationRequestData represents the specification data in the request.
type SpecificationRequestData struct {
	Schema   map[string]any `json:"schema"`
	Version  string         `json:"version"`
	Endpoint string         `json:"endpoint"`
}

// UpdateSpecificationRequest represents the request payload for updating a specification.
type UpdateSpecificationRequest struct {
	Specification UpdateSpecificationData `json:"specification"`
}

// UpdateSpecificationData represents the specification update data.
type UpdateSpecificationData struct {
	Schema   map[string]any `json:"schema,omitempty"`
	Version  string         `json:"version,omitempty"`
	Endpoint string         `json:"endpoint,omitempty"`
}

// CreateIdentityRequest represents the request payload for creating an identity.
type CreateIdentityRequest struct {
	Identity IdentityRequestData `json:"identity"`
}

// IdentityRequestData represents the identity data in the request.
type IdentityRequestData struct {
	APIKey     string     `json:"api_key"`
	Validation Validation `json:"validation"`
}

// UpdateIdentityRequest represents the request payload for updating an identity.
type UpdateIdentityRequest struct {
	Identity UpdateIdentityData `json:"identity"`
}

// UpdateIdentityData represents the identity update data.
type UpdateIdentityData struct {
	APIKey     string      `json:"api_key,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

// handleAPIError processes API error responses.
func (s *Service) handleAPIError(resp any) error {
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
func (s *Service) extractErrorResponse(resp any) (errorResponse, bool) {
	type errorResponse interface {
		IsError() bool
		Error() any
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
	// Try to parse as nested error format (e.g., module.reference)
	if err := s.parseNestedError(body, statusCode); err != nil {
		return err
	}

	// Try to parse as map[string][]string (array format)
	if err := s.parseArrayError(body, statusCode); err != nil {
		return err
	}

	// Try to parse as map[string]string (single string format)
	return s.parseStringError(body, statusCode)
}

// parseNestedError parses errors with arbitrary nesting depth (e.g., module.config.database.connection).
func (s *Service) parseNestedError(body []byte, statusCode int) error {
	var rawNestedError struct {
		Errors interface{} `json:"errors"`
	}

	if err := json.Unmarshal(body, &rawNestedError); err == nil && rawNestedError.Errors != nil {
		flatErrors := make(map[string][]string)
		s.flattenErrors(rawNestedError.Errors, "", flatErrors)

		if len(flatErrors) > 0 {
			return &Error{
				StatusCode: statusCode,
				Errors:     flatErrors,
			}
		}
	}
	return nil
}

// flattenErrors recursively flattens nested error structures into dot-notation keys.
func (s *Service) flattenErrors(data interface{}, prefix string, result map[string][]string) {
	switch v := data.(type) {
	case map[string]interface{}:
		// Handle nested objects
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = fmt.Sprintf("%s.%s", prefix, key)
			}
			s.flattenErrors(value, newPrefix, result)
		}

	case []interface{}:
		// Handle arrays - convert to []string if possible
		var stringSlice []string
		for _, item := range v {
			if str, ok := item.(string); ok {
				stringSlice = append(stringSlice, str)
			} else {
				// If array contains non-strings, convert to string representation
				stringSlice = append(stringSlice, fmt.Sprintf("%v", item))
			}
		}
		if prefix != "" {
			result[prefix] = stringSlice
		}

	case []string:
		// Handle string arrays directly
		if prefix != "" {
			result[prefix] = v
		}

	case string:
		// Handle single string values
		if prefix != "" {
			result[prefix] = []string{v}
		}

	default:
		// Handle other types by converting to string
		if prefix != "" {
			result[prefix] = []string{fmt.Sprintf("%v", v)}
		}
	}
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
	Error() any
	StatusCode() int
	Status() string
	Body() []byte
}
