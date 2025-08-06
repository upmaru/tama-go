package contexts

import (
	"errors"
	"fmt"
)

// Input represents a contexts input resource.
type Input struct {
	ID               string `json:"id,omitempty"`
	Type             string `json:"type"`
	ThoughtContextID string `json:"thought_context_id,omitempty"`
	ClassCorpusID    string `json:"class_corpus_id"`
	ProvisionState   string `json:"provision_state,omitempty"`
}

// InputResponse represents the API response for input operations.
type InputResponse struct {
	Data Input `json:"data"`
}

// CreateInputRequest represents the request payload for creating a contexts input.
type CreateInputRequest struct {
	Input CreateInputData `json:"input"`
}

// CreateInputData represents the input data in the create request.
type CreateInputData struct {
	Type          string `json:"type"`
	ClassCorpusID string `json:"class_corpus_id"`
}

// UpdateInputRequest represents the request payload for updating a contexts input.
type UpdateInputRequest struct {
	Input UpdateInputData `json:"input"`
}

// UpdateInputData represents the input update data.
type UpdateInputData struct {
	Type          string `json:"type,omitempty"`
	ClassCorpusID string `json:"class_corpus_id,omitempty"`
}


// GetInput retrieves a specific input by ID.
// GET /provision/contexts/inputs/:id
func (s *Service) GetInput(id string) (*Input, error) {
	if id == "" {
		return nil, errors.New("input ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetResult(&inputResp).
		Get(fmt.Sprintf("/provision/contexts/inputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// CreateInput creates a new contexts input within a thought context.
// POST /provision/contexts/:thought_context_id/inputs
func (s *Service) CreateInput(thoughtContextID string, req CreateInputRequest) (*Input, error) {
	if thoughtContextID == "" {
		return nil, errors.New("thought context ID is required")
	}
	if req.Input.Type == "" {
		return nil, errors.New("input type is required")
	}
	if req.Input.ClassCorpusID == "" {
		return nil, errors.New("class corpus ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&inputResp).
		Post(fmt.Sprintf("/provision/contexts/%s/inputs", thoughtContextID))

	if err != nil {
		return nil, fmt.Errorf("failed to create input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// UpdateInput updates an existing input using PATCH.
// PATCH /provision/contexts/inputs/:id
func (s *Service) UpdateInput(id string, req UpdateInputRequest) (*Input, error) {
	if id == "" {
		return nil, errors.New("input ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&inputResp).
		Patch(fmt.Sprintf("/provision/contexts/inputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// ReplaceInput replaces an existing input using PUT.
// PUT /provision/contexts/inputs/:id
func (s *Service) ReplaceInput(id string, req UpdateInputRequest) (*Input, error) {
	if id == "" {
		return nil, errors.New("input ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&inputResp).
		Put(fmt.Sprintf("/provision/contexts/inputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// DeleteInput deletes an input by ID.
// DELETE /provision/contexts/inputs/:id
func (s *Service) DeleteInput(id string) error {
	if id == "" {
		return errors.New("input ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/contexts/inputs/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

