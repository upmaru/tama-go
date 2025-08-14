package perception

import (
	"errors"
	"fmt"
)

// Initializer represents a perception initializer resource.
type Initializer struct {
	ID             string         `json:"id,omitempty"`
	Parameters     map[string]any `json:"parameters,omitempty"`
	Index          *int           `json:"index,omitempty"`
	ProvisionState string         `json:"provision_state"`
	ThoughtID      string         `json:"thought_id,omitempty"`
	ClassID        string         `json:"class_id"`
	Reference      string         `json:"reference"`
}

// InitializerResponse represents the API response for initializer operations.
type InitializerResponse struct {
	Data Initializer `json:"data"`
}

// CreateInitializerRequest represents the request payload for creating an initializer.
type CreateInitializerRequest struct {
	Initializer InitializerRequestData `json:"initializer"`
}

// InitializerRequestData represents the initializer data in the request.
type InitializerRequestData struct {
	Parameters map[string]any `json:"parameters,omitempty"`
	Index      *int           `json:"index,omitempty"`
	ClassID    string         `json:"class_id"`
	Reference  string         `json:"reference"`
}

// UpdateInitializerRequest represents the request payload for updating an initializer.
type UpdateInitializerRequest struct {
	Initializer UpdateInitializerData `json:"initializer"`
}

// UpdateInitializerData represents the initializer update data.
type UpdateInitializerData struct {
	Parameters map[string]any `json:"parameters,omitempty"`
	Index      *int           `json:"index,omitempty"`
	ClassID    string         `json:"class_id,omitempty"`
	Reference  string         `json:"reference,omitempty"`
}

// GetInitializer retrieves a specific initializer by ID.
// GET /provision/perception/initializers/:id.
func (s *Service) GetInitializer(id string) (*Initializer, error) {
	if id == "" {
		return nil, errors.New("initializer ID is required")
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetResult(&initializerResp).
		Get(fmt.Sprintf("/provision/perception/initializers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// CreateInitializer creates a new initializer within a thought.
// POST /provision/perception/thoughts/:thought_id/initializers.
func (s *Service) CreateInitializer(thoughtID string, req CreateInitializerRequest) (*Initializer, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if req.Initializer.ClassID == "" {
		return nil, errors.New("initializer class ID is required")
	}
	if req.Initializer.Reference == "" {
		return nil, errors.New("initializer reference is required")
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&initializerResp).
		Post(fmt.Sprintf("/provision/perception/thoughts/%s/initializers", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to create initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// UpdateInitializer updates an existing initializer using PATCH.
// PATCH /provision/perception/initializers/:id.
func (s *Service) UpdateInitializer(id string, req UpdateInitializerRequest) (*Initializer, error) {
	if id == "" {
		return nil, errors.New("initializer ID is required")
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&initializerResp).
		Patch(fmt.Sprintf("/provision/perception/initializers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// ReplaceInitializer replaces an existing initializer using PUT.
// PUT /provision/perception/initializers/:id.
func (s *Service) ReplaceInitializer(id string, req UpdateInitializerRequest) (*Initializer, error) {
	if id == "" {
		return nil, errors.New("initializer ID is required")
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&initializerResp).
		Put(fmt.Sprintf("/provision/perception/initializers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// DeleteInitializer deletes an initializer by ID.
// DELETE /provision/perception/initializers/:id.
func (s *Service) DeleteInitializer(id string) error {
	if id == "" {
		return errors.New("initializer ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/initializers/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
