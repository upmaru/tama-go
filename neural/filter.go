//nolint:dupl // Similar CRUD patterns to bridge.go are intentional by design.
package neural

import (
	"errors"
	"fmt"
)

// Filter represents a neural filter resource.
// All properties are strings per API contract.
// Structure:
//
//	id, listener_id, chain_id, provision_state
//
// All fields are strings.
type Filter struct {
	ID             string `json:"id,omitempty"`
	ListenerID     string `json:"listener_id"`
	ChainID        string `json:"chain_id"`
	ProvisionState string `json:"provision_state"`
}

// FilterResponse represents the API response for filter operations.
type FilterResponse struct {
	Data Filter `json:"data"`
}

// CreateFilterRequest represents the request payload for creating a filter.
type CreateFilterRequest struct {
	Filter FilterRequestData `json:"filter"`
}

// FilterRequestData represents the filter data in the create request.
type FilterRequestData struct {
	ChainID string `json:"chain_id"`
}

// UpdateFilterRequest represents the request payload for updating or replacing a filter.
type UpdateFilterRequest struct {
	Filter UpdateFilterData `json:"filter"`
}

// UpdateFilterData represents the filter update data.
// Currently the only mutable field is the chain association.
type UpdateFilterData struct {
	ChainID string `json:"chain_id,omitempty"`
}

// GetFilter retrieves a specific filter by ID.
// GET /provision/neural/filters/:id.
func (s *Service) GetFilter(id string) (*Filter, error) {
	if id == "" {
		return nil, errors.New("filter ID is required")
	}

	var filterResp FilterResponse
	resp, err := s.client.R().
		SetResult(&filterResp).
		Get(fmt.Sprintf("/provision/neural/filters/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get filter: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &filterResp.Data, nil
}

// CreateFilter creates a new filter for a listener.
// POST /provision/neural/listeners/:listener_id/filters.
func (s *Service) CreateFilter(listenerID string, req CreateFilterRequest) (*Filter, error) {
	if listenerID == "" {
		return nil, errors.New("listener ID is required")
	}
	if req.Filter.ChainID == "" {
		return nil, errors.New("chain ID is required")
	}

	var filterResp FilterResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&filterResp).
		Post(fmt.Sprintf("/provision/neural/listeners/%s/filters", listenerID))

	if err != nil {
		return nil, fmt.Errorf("failed to create filter: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &filterResp.Data, nil
}

// UpdateFilter updates an existing filter using PATCH.
// PATCH /provision/neural/filters/:id.
func (s *Service) UpdateFilter(id string, req UpdateFilterRequest) (*Filter, error) {
	if id == "" {
		return nil, errors.New("filter ID is required")
	}
	// Require at least one field; currently ChainID is the only updatable field
	if req.Filter.ChainID == "" {
		return nil, errors.New("chain ID is required")
	}

	var filterResp FilterResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&filterResp).
		Patch(fmt.Sprintf("/provision/neural/filters/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update filter: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &filterResp.Data, nil
}

// ReplaceFilter replaces an existing filter using PUT.
// PUT /provision/neural/filters/:id.
func (s *Service) ReplaceFilter(id string, req UpdateFilterRequest) (*Filter, error) {
	if id == "" {
		return nil, errors.New("filter ID is required")
	}
	if req.Filter.ChainID == "" {
		return nil, errors.New("chain ID is required")
	}

	var filterResp FilterResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&filterResp).
		Put(fmt.Sprintf("/provision/neural/filters/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace filter: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &filterResp.Data, nil
}

// DeleteFilter deletes a filter by ID.
// DELETE /provision/neural/filters/:id.
func (s *Service) DeleteFilter(id string) error {
	if id == "" {
		return errors.New("filter ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/neural/filters/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete filter: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
