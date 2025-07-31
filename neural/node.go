package neural

import (
	"errors"
	"fmt"
)

// Node represents a neural node resource.
type Node struct {
	ID             string `json:"id,omitempty"`
	On             string `json:"on,omitempty"`
	Type           string `json:"type"`
	SpaceID        string `json:"space_id,omitempty"`
	ClassID        string `json:"class_id"`
	ChainID        string `json:"chain_id"`
	ProvisionState string `json:"provision_state,omitempty"`
}

// NodeResponse represents the API response for node operations.
type NodeResponse struct {
	Data Node `json:"data"`
}

// CreateNodeRequest represents the request payload for creating a node.
type CreateNodeRequest struct {
	Node NodeRequestData `json:"node"`
}

// NodeRequestData represents the node data in the request.
type NodeRequestData struct {
	On      string `json:"on,omitempty"`
	Type    string `json:"type"`
	ClassID string `json:"class_id"`
	ChainID string `json:"chain_id"`
}

// UpdateNodeRequest represents the request payload for updating a node.
type UpdateNodeRequest struct {
	Node UpdateNodeData `json:"node"`
}

// UpdateNodeData represents the node update data.
type UpdateNodeData struct {
	On   string `json:"on,omitempty"`
	Type string `json:"type,omitempty"`
}

// GetNode retrieves a specific node by ID.
// GET /provision/neural/nodes/:id.
func (s *Service) GetNode(id string) (*Node, error) {
	if id == "" {
		return nil, errors.New("node ID is required")
	}

	var nodeResp NodeResponse
	resp, err := s.client.R().
		SetResult(&nodeResp).
		Get(fmt.Sprintf("/provision/neural/nodes/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &nodeResp.Data, nil
}

// CreateNode creates a new node.
// POST /provision/neural/spaces/:space_id/nodes.
func (s *Service) CreateNode(spaceID string, req CreateNodeRequest) (*Node, error) {
	if spaceID == "" || req.Node.ClassID == "" || req.Node.ChainID == "" || req.Node.Type == "" {
		return nil, errors.New("required parameters are missing")
	}

	var nodeResp NodeResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&nodeResp).
		Post(fmt.Sprintf("/provision/neural/spaces/%s/nodes", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &nodeResp.Data, nil
}

// UpdateNode updates an existing node using PATCH.
// PATCH /provision/neural/nodes/:id.
func (s *Service) UpdateNode(id string, req UpdateNodeRequest) (*Node, error) {
	if id == "" {
		return nil, errors.New("node ID is required")
	}

	var nodeResp NodeResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&nodeResp).
		Patch(fmt.Sprintf("/provision/neural/nodes/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &nodeResp.Data, nil
}

// ReplaceNode replaces an existing node using PUT.
// PUT /provision/neural/nodes/:id.
func (s *Service) ReplaceNode(id string, req UpdateNodeRequest) (*Node, error) {
	if id == "" {
		return nil, errors.New("node ID is required")
	}

	var nodeResp NodeResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&nodeResp).
		Put(fmt.Sprintf("/provision/neural/nodes/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace node: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &nodeResp.Data, nil
}

// DeleteNode deletes a node by ID.
// DELETE /provision/neural/nodes/:id.
func (s *Service) DeleteNode(id string) error {
	if id == "" {
		return errors.New("node ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/neural/nodes/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
