package motor

import (
	"errors"
	"fmt"
)

// Action represents a motor action.
type Action struct {
	ID              string `json:"id"`
	Identifier      string `json:"identifier"`
	Path            string `json:"path"`
	Method          string `json:"method"`
	SpecificationID string `json:"specification_id"`
}

// ActionResponse wraps the Action data in API responses.
type ActionResponse struct {
	Data Action `json:"data"`
}

// GetAction retrieves a specific action by specification ID and action ID.
func (s *Service) GetAction(specID, id string) (*Action, error) {
	if specID == "" || id == "" {
		return nil, errors.New("specification ID and action ID are required")
	}

	var respWrapper ActionResponse
	resp, err := s.client.R().
		SetResult(&respWrapper).
		Get(fmt.Sprintf("/provision/motor/specifications/%s/actions/%s", specID, id))

	if err != nil {
		return nil, err
	}

	if handleErr := s.handleAPIError(resp); handleErr != nil {
		return nil, handleErr
	}

	return &respWrapper.Data, nil
}

// Execute runs the motor action.
func (a *Action) Execute() error {
	return nil
}
