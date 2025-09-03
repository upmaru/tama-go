package neural

import (
	"errors"
	"fmt"
)

// Listener represents a neural listener resource.
// All properties are strings per API contract.
type Listener struct {
	ID             string `json:"id,omitempty"`
	SpaceID        string `json:"space_id"`
	Endpoint       string `json:"endpoint"`
	ProvisionState string `json:"provision_state"`
}

// ListenerResponse represents the API response for listener operations.
type ListenerResponse struct {
	Data Listener `json:"data"`
}

// CreateListenerRequest represents the request payload for creating a listener.
type CreateListenerRequest struct {
	Listener ListenerRequestData `json:"listener"`
}

// ListenerRequestData represents the listener data in the create request.
type ListenerRequestData struct {
	Endpoint string `json:"endpoint"`
}

// UpdateListenerRequest represents the request payload for updating or replacing a listener.
type UpdateListenerRequest struct {
	Listener UpdateListenerData `json:"listener"`
}

// UpdateListenerData represents the listener update data.
type UpdateListenerData struct {
	Endpoint string `json:"endpoint,omitempty"`
}

// GetListener retrieves a specific listener by ID.
// GET /provision/neural/listeners/:id.
func (s *Service) GetListener(id string) (*Listener, error) {
	if id == "" {
		return nil, errors.New("listener ID is required")
	}

	var listenerResp ListenerResponse
	resp, err := s.client.R().
		SetResult(&listenerResp).
		Get(fmt.Sprintf("/provision/neural/listeners/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get listener: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &listenerResp.Data, nil
}

// CreateListener creates a new listener within a space.
// POST /provision/neural/spaces/:space_id/listeners.
func (s *Service) CreateListener(spaceID string, req CreateListenerRequest) (*Listener, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if req.Listener.Endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	var listenerResp ListenerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&listenerResp).
		Post(fmt.Sprintf("/provision/neural/spaces/%s/listeners", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &listenerResp.Data, nil
}

// UpdateListener updates an existing listener using PATCH.
// PATCH /provision/neural/listeners/:id.
func (s *Service) UpdateListener(id string, req UpdateListenerRequest) (*Listener, error) {
	if id == "" {
		return nil, errors.New("listener ID is required")
	}
	if req.Listener.Endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	var listenerResp ListenerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&listenerResp).
		Patch(fmt.Sprintf("/provision/neural/listeners/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update listener: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &listenerResp.Data, nil
}

// ReplaceListener replaces an existing listener using PUT.
// PUT /provision/neural/listeners/:id.
func (s *Service) ReplaceListener(id string, req UpdateListenerRequest) (*Listener, error) {
	if id == "" {
		return nil, errors.New("listener ID is required")
	}
	if req.Listener.Endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	var listenerResp ListenerResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&listenerResp).
		Put(fmt.Sprintf("/provision/neural/listeners/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace listener: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &listenerResp.Data, nil
}

// DeleteListener deletes a listener by ID.
// DELETE /provision/neural/listeners/:id.
func (s *Service) DeleteListener(id string) error {
	if id == "" {
		return errors.New("listener ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/neural/listeners/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete listener: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

