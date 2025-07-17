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
	Type string `json:"type"` // "root" or "component"
}

// UpdateSpaceRequest represents the request payload for updating a space.
type UpdateSpaceRequest struct {
	Space UpdateSpaceData `json:"space"`
}

// UpdateSpaceData represents the space update data.
type UpdateSpaceData struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"` // "root" or "component"
}

// handleAPIError processes API error responses.
func (s *Service) handleAPIError(resp interface{}) error {
	type errorResponse interface {
		IsError() bool
		Error() interface{}
		StatusCode() int
		Status() string
		Body() []byte
	}

	if errResp, ok := resp.(errorResponse); ok && errResp.IsError() {
		// Always try to parse the raw JSON body first for better error handling
		if body := errResp.Body(); len(body) > 0 {
			// First try to parse as map[string][]string (array format)
			var rawArrayError struct {
				Errors map[string][]string `json:"errors"`
			}

			if err := json.Unmarshal(body, &rawArrayError); err == nil && rawArrayError.Errors != nil {
				return &Error{
					StatusCode: errResp.StatusCode(),
					Errors:     rawArrayError.Errors,
				}
			}

			// If that fails, try to parse as map[string]string (single string format)
			var rawStringError struct {
				Errors map[string]string `json:"errors"`
			}

			if err := json.Unmarshal(body, &rawStringError); err == nil && rawStringError.Errors != nil {
				// Convert string values to string arrays
				convertedErrors := make(map[string][]string)
				for field, message := range rawStringError.Errors {
					convertedErrors[field] = []string{message}
				}
				return &Error{
					StatusCode: errResp.StatusCode(),
					Errors:     convertedErrors,
				}
			}

			// Fallback: try to parse as a general error response
			var generalError Error
			if err := json.Unmarshal(body, &generalError); err == nil {
				generalError.StatusCode = errResp.StatusCode()
				return &generalError
			}
		}

		// Fallback: try to get the structured error from resty
		if apiErrorResp, isError := errResp.Error().(*Error); isError {
			apiErrorResp.StatusCode = errResp.StatusCode()
			return apiErrorResp
		}

		return fmt.Errorf("API error: %s", errResp.Status())
	}
	return nil
}
