package tools

import (
	"errors"
	"fmt"
)

// Initializer represents a tools initializer resource.
type Initializer struct {
	ID             string         `json:"id,omitempty"`
	Reference      string         `json:"reference"`
	Index          string         `json:"index"`
	Parameters     map[string]any `json:"parameters"`
	ThoughtToolID  string         `json:"thought_tool_id,omitempty"`
	ProvisionState string         `json:"provision_state"`
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
	Reference  string         `json:"reference"`
	Index      string         `json:"index"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// UpdateInitializerRequest represents the request payload for updating an initializer.
type UpdateInitializerRequest struct {
	Initializer UpdateInitializerData `json:"initializer"`
}

// UpdateInitializerData represents the initializer update data.
type UpdateInitializerData struct {
	Reference  string         `json:"reference,omitempty"`
	Index      string         `json:"index,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// GetInitializer retrieves a specific initializer by ID.
// GET /provision/tools/initializers/:id.
func (s *Service) GetInitializer(id string) (*Initializer, error) {
	if id == "" {
		return nil, errors.New("initializer ID is required")
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetResult(&initializerResp).
		Get(fmt.Sprintf("/provision/tools/initializers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// CreateInitializer creates a new initializer within a thought tool.
// POST /provision/tools/:thought_tool_id/initializers.
func (s *Service) CreateInitializer(thoughtToolID string, req CreateInitializerRequest) (*Initializer, error) {
	if thoughtToolID == "" {
		return nil, errors.New("thought tool ID is required")
	}
	if req.Initializer.Reference == "" {
		return nil, errors.New("initializer reference is required")
	}
	if req.Initializer.Index == "" {
		return nil, errors.New("initializer index is required")
	}

	// Initialize parameters map if nil
	if req.Initializer.Parameters == nil {
		req.Initializer.Parameters = make(map[string]any)
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&initializerResp).
		Post(fmt.Sprintf("/provision/tools/%s/initializers", thoughtToolID))

	if err != nil {
		return nil, fmt.Errorf("failed to create initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// UpdateInitializer updates an existing initializer using PATCH.
// PATCH /provision/tools/initializers/:id.
func (s *Service) UpdateInitializer(id string, req UpdateInitializerRequest) (*Initializer, error) {
	if id == "" {
		return nil, errors.New("initializer ID is required")
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&initializerResp).
		Patch(fmt.Sprintf("/provision/tools/initializers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// ReplaceInitializer replaces an existing initializer using PUT.
// PUT /provision/tools/initializers/:id.
func (s *Service) ReplaceInitializer(id string, req UpdateInitializerRequest) (*Initializer, error) {
	if id == "" {
		return nil, errors.New("initializer ID is required")
	}

	var initializerResp InitializerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&initializerResp).
		Put(fmt.Sprintf("/provision/tools/initializers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &initializerResp.Data, nil
}

// DeleteInitializer deletes an initializer by ID.
// DELETE /provision/tools/initializers/:id.
func (s *Service) DeleteInitializer(id string) error {
	if id == "" {
		return errors.New("initializer ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/tools/initializers/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete initializer: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
