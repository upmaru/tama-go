package perception

import (
	"errors"
	"fmt"
)

// Delegation represents a perception delegation resource.
type Delegation struct {
	ID              string `json:"id,omitempty"`
	ThoughtID       string `json:"thought_id,omitempty"`
	TargetThoughtID string `json:"target_thought_id"`
	ProvisionState  string `json:"provision_state"`
}

// DelegationResponse represents the API response for delegation operations.
type DelegationResponse struct {
	Data Delegation `json:"data"`
}

// CreateDelegationRequest represents the request payload for creating a delegation.
type CreateDelegationRequest struct {
	Delegation DelegationRequestData `json:"delegation"`
}

// DelegationRequestData represents the delegation data in the request.
type DelegationRequestData struct {
	TargetThoughtID string `json:"target_thought_id"`
}

// UpdateDelegationRequest represents the request payload for updating a delegation.
type UpdateDelegationRequest struct {
	Delegation UpdateDelegationData `json:"delegation"`
}

// UpdateDelegationData represents the delegation update data.
type UpdateDelegationData struct {
	TargetThoughtID string `json:"target_thought_id,omitempty"`
}

// GetDelegation retrieves the delegation for a specific thought.
// GET /provision/perception/thoughts/:thought_id/delegation.
func (s *Service) GetDelegation(thoughtID string) (*Delegation, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}

	var delegationResp DelegationResponse
	resp, err := s.client.R().
		SetResult(&delegationResp).
		Get(fmt.Sprintf("/provision/perception/thoughts/%s/delegation", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to get delegation: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &delegationResp.Data, nil
}

// CreateDelegation creates a new delegation within a thought.
// POST /provision/perception/thoughts/:thought_id/delegation.
func (s *Service) CreateDelegation(thoughtID string, req CreateDelegationRequest) (*Delegation, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if req.Delegation.TargetThoughtID == "" {
		return nil, errors.New("target thought ID is required")
	}

	var delegationResp DelegationResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&delegationResp).
		Post(fmt.Sprintf("/provision/perception/thoughts/%s/delegation", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to create delegation: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &delegationResp.Data, nil
}

// UpdateDelegation updates an existing delegation using PATCH.
// PATCH /provision/perception/thoughts/:thought_id/delegation.
func (s *Service) UpdateDelegation(thoughtID string, req UpdateDelegationRequest) (*Delegation, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}

	var delegationResp DelegationResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&delegationResp).
		Patch(fmt.Sprintf("/provision/perception/thoughts/%s/delegation", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to update delegation: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &delegationResp.Data, nil
}

// ReplaceDelegation replaces an existing delegation using PUT.
// PUT /provision/perception/thoughts/:thought_id/delegation.
func (s *Service) ReplaceDelegation(thoughtID string, req UpdateDelegationRequest) (*Delegation, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}

	var delegationResp DelegationResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&delegationResp).
		Put(fmt.Sprintf("/provision/perception/thoughts/%s/delegation", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to replace delegation: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &delegationResp.Data, nil
}

// DeleteDelegation deletes a delegation by thought ID.
// DELETE /provision/perception/thoughts/:thought_id/delegation.
func (s *Service) DeleteDelegation(thoughtID string) error {
	if thoughtID == "" {
		return errors.New("thought ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/thoughts/%s/delegation", thoughtID))

	if err != nil {
		return fmt.Errorf("failed to delete delegation: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
