package neural

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Service handles all neural-related API operations.
type Service struct {
	client *resty.Client
}

// NewService creates a new neural service instance.
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

// Space represents a neural space resource.
type Space struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	Slug         string `json:"slug,omitempty"`
	Type         string `json:"type"`
	CurrentState string `json:"current_state"`
}

// SpaceResponse represents the API response for space operations.
type SpaceResponse struct {
	Data Space `json:"data"`
}

// CreateSpaceRequest represents the request payload for creating a space.
type CreateSpaceRequest struct {
	Space SpaceRequestData `json:"space"`
}

// SpaceRequestData represents the space data in the request.
type SpaceRequestData struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// UpdateSpaceRequest represents the request payload for updating a space.
type UpdateSpaceRequest struct {
	Space UpdateSpaceData `json:"space"`
}

// UpdateSpaceData represents the space update data.
type UpdateSpaceData struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// Processor represents a neural processor resource.
type Processor struct {
	ID            string         `json:"id,omitempty"`
	SpaceID       string         `json:"space_id,omitempty"`
	ModelID       string         `json:"model_id,omitempty"`
	Configuration map[string]any `json:"configuration"`
	CurrentState  string         `json:"current_state"`
	Type          string         `json:"type"`
}

// ProcessorResponse represents the API response for processor operations.
type ProcessorResponse struct {
	Data Processor `json:"data"`
}

// CreateProcessorRequest represents the request payload for creating a processor.
type CreateProcessorRequest struct {
	Processor ProcessorRequestData `json:"processor"`
}

// ProcessorRequestData represents the processor data in the request.
type ProcessorRequestData struct {
	ModelID       string         `json:"model_id"`
	Configuration map[string]any `json:"configuration"`
}

// UpdateProcessorRequest represents the request payload for updating a processor.
type UpdateProcessorRequest struct {
	Processor UpdateProcessorData `json:"processor"`
}

// UpdateProcessorData represents the processor update data.
type UpdateProcessorData struct {
	ModelID       string         `json:"model_id,omitempty"`
	Configuration map[string]any `json:"configuration,omitempty"`
}

// Class represents a neural class resource.
type Class struct {
	ID           string         `json:"id,omitempty"`
	SpaceID      string         `json:"space_id,omitempty"`
	CurrentState string         `json:"current_state"`
	Schema       map[string]any `json:"schema"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
}

// ClassResponse represents the API response for class operations.
type ClassResponse struct {
	Data Class `json:"data"`
}

// CreateClassRequest represents the request payload for creating a class.
type CreateClassRequest struct {
	Class ClassRequestData `json:"class"`
}

// ClassRequestData represents the class data in the request.
type ClassRequestData struct {
	Schema map[string]any `json:"schema"`
}

// UpdateClassRequest represents the request payload for updating a class.
type UpdateClassRequest struct {
	Class UpdateClassData `json:"class"`
}

// UpdateClassData represents the class update data.
type UpdateClassData struct {
	Schema map[string]any `json:"schema,omitempty"`
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
	// Try to parse as map[string][]string (array format)
	if err := s.parseArrayError(body, statusCode); err != nil {
		return err
	}

	// Try to parse as map[string]string (single string format)
	if err := s.parseStringError(body, statusCode); err != nil {
		return err
	}

	// Try to parse as a general error response
	return s.parseGeneralError(body, statusCode)
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

// parseGeneralError parses general error response format.
func (s *Service) parseGeneralError(body []byte, statusCode int) error {
	var generalError Error
	if err := json.Unmarshal(body, &generalError); err == nil {
		generalError.StatusCode = statusCode
		return &generalError
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
