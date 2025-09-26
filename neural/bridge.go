//nolint:dupl // Similar CRUD patterns to filter.go are intentional by design.
package neural

import (
	"errors"
	"fmt"
)

// Bridge represents a neural bridge resource.
type Bridge struct {
	ID             string `json:"id,omitempty"`
	SpaceID        string `json:"space_id"`
	TargetSpaceID  string `json:"target_space_id"`
	ProvisionState string `json:"provision_state"`
}

// BridgeResponse represents the API response for bridge operations.
type BridgeResponse struct {
	Data Bridge `json:"data"`
}

// CreateBridgeRequest represents the request payload for creating a bridge.
type CreateBridgeRequest struct {
	Bridge BridgeRequestData `json:"bridge"`
}

// BridgeRequestData represents the bridge data in the request.
type BridgeRequestData struct {
	TargetSpaceID string `json:"target_space_id"`
}

// UpdateBridgeRequest represents the request payload for updating a bridge.
type UpdateBridgeRequest struct {
	Bridge UpdateBridgeData `json:"bridge"`
}

// UpdateBridgeData represents the bridge update data.
type UpdateBridgeData struct {
	TargetSpaceID string `json:"target_space_id,omitempty"`
}

// GetBridge retrieves a specific bridge by ID.
// GET /provision/neural/bridges/:id.
func (s *Service) GetBridge(id string) (*Bridge, error) {
	if id == "" {
		return nil, errors.New("bridge ID is required")
	}

	var bridgeResp BridgeResponse
	resp, err := s.client.R().
		SetResult(&bridgeResp).
		Get(fmt.Sprintf("/provision/neural/bridges/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get bridge: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &bridgeResp.Data, nil
}

// CreateBridge creates a new bridge within a space.
// POST /provision/neural/spaces/:space_id/bridges.
func (s *Service) CreateBridge(spaceID string, req CreateBridgeRequest) (*Bridge, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if req.Bridge.TargetSpaceID == "" {
		return nil, errors.New("target space ID is required")
	}

	var bridgeResp BridgeResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&bridgeResp).
		Post(fmt.Sprintf("/provision/neural/spaces/%s/bridges", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create bridge: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &bridgeResp.Data, nil
}

// UpdateBridge updates an existing bridge using PATCH.
// PATCH /provision/neural/bridges/:id.
func (s *Service) UpdateBridge(id string, req UpdateBridgeRequest) (*Bridge, error) {
	if id == "" {
		return nil, errors.New("bridge ID is required")
	}
	if req.Bridge.TargetSpaceID == "" {
		return nil, errors.New("target space ID is required")
	}

	var bridgeResp BridgeResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&bridgeResp).
		Patch(fmt.Sprintf("/provision/neural/bridges/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update bridge: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &bridgeResp.Data, nil
}

// ReplaceBridge replaces an existing bridge using PUT.
// PUT /provision/neural/bridges/:id.
func (s *Service) ReplaceBridge(id string, req UpdateBridgeRequest) (*Bridge, error) {
	if id == "" {
		return nil, errors.New("bridge ID is required")
	}
	if req.Bridge.TargetSpaceID == "" {
		return nil, errors.New("target space ID is required")
	}

	var bridgeResp BridgeResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&bridgeResp).
		Put(fmt.Sprintf("/provision/neural/bridges/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace bridge: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &bridgeResp.Data, nil
}

// DeleteBridge deletes a bridge by ID.
// DELETE /provision/neural/bridges/:id.
func (s *Service) DeleteBridge(id string) error {
	if id == "" {
		return errors.New("bridge ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/neural/bridges/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete bridge: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
