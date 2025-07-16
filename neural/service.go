package neural

import (
	"fmt"

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
	}

	if errResp, ok := resp.(errorResponse); ok && errResp.IsError() {
		if apiErrorResp, isError := errResp.Error().(*Error); isError {
			apiErrorResp.StatusCode = errResp.StatusCode()
			return apiErrorResp
		}
		return fmt.Errorf("API error: %s", errResp.Status())
	}
	return nil
}
