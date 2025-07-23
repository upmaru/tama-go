package perception

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Service handles all perception-related API operations.
type Service struct {
	client *resty.Client
}

// NewService creates a new perception service instance.
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

// Chain represents a perception chain resource.
type Chain struct {
	ID           string `json:"id,omitempty"`
	SpaceID      string `json:"space_id,omitempty"`
	Name         string `json:"name"`
	Slug         string `json:"slug,omitempty"`
	CurrentState string `json:"current_state"`
}

// ChainResponse represents the API response for chain operations.
type ChainResponse struct {
	Data Chain `json:"data"`
}

// CreateChainRequest represents the request payload for creating a chain.
type CreateChainRequest struct {
	Chain ChainRequestData `json:"chain"`
}

// ChainRequestData represents the chain data in the request.
type ChainRequestData struct {
	Name string `json:"name"`
}

// UpdateChainRequest represents the request payload for updating a chain.
type UpdateChainRequest struct {
	Chain UpdateChainData `json:"chain"`
}

// UpdateChainData represents the chain update data.
type UpdateChainData struct {
	Name string `json:"name,omitempty"`
}

// Module represents a thought module configuration.
type Module struct {
	ID         string         `json:"id,omitempty"`
	Reference  string         `json:"reference"`
	Parameters map[string]any `json:"parameters"`
}

// Thought represents a perception thought resource.
type Thought struct {
	ID            string `json:"id,omitempty"`
	ChainID       string `json:"chain_id,omitempty"`
	OutputClassID string `json:"output_class_id,omitempty"`
	Module        Module `json:"module"`
	CurrentState  string `json:"current_state"`
	Relation      string `json:"relation"`
	Index         int    `json:"index"`
}

// ThoughtResponse represents the API response for thought operations.
type ThoughtResponse struct {
	Data Thought `json:"data"`
}

// CreateThoughtRequest represents the request payload for creating a thought.
type CreateThoughtRequest struct {
	Thought ThoughtRequestData `json:"thought"`
}

// ThoughtRequestData represents the thought data in the request.
type ThoughtRequestData struct {
	Relation      string `json:"relation"`
	OutputClassID string `json:"output_class_id,omitempty"`
	Module        Module `json:"module"`
}

// UpdateThoughtRequest represents the request payload for updating a thought.
type UpdateThoughtRequest struct {
	Thought UpdateThoughtData `json:"thought"`
}

// UpdateThoughtData represents the thought update data.
type UpdateThoughtData struct {
	Relation      string `json:"relation,omitempty"`
	OutputClassID string `json:"output_class_id,omitempty"`
	Module        Module `json:"module,omitempty"`
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
