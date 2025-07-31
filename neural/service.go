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
	if err := s.parseStringError(body, statusCode); err != nil {
		return err
	}

	// Try to parse as a general error response
	return s.parseGeneralError(body, statusCode)
}

// parseNestedError parses errors with arbitrary nesting depth (e.g., module.config.database.connection).
func (s *Service) parseNestedError(body []byte, statusCode int) error {
	var rawNestedError struct {
		Errors any `json:"errors"`
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
func (s *Service) flattenErrors(data any, prefix string, result map[string][]string) {
	switch v := data.(type) {
	case map[string]any:
		// Handle nested objects
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = fmt.Sprintf("%s.%s", prefix, key)
			}
			s.flattenErrors(value, newPrefix, result)
		}

	case []any:
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
