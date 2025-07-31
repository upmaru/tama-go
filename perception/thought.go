package perception

import (
	"errors"
	"fmt"
)

// Module represents a thought module configuration.
type Module struct {
	ID         string         `json:"id,omitempty"`
	Reference  string         `json:"reference"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// Thought represents a perception thought resource.
type Thought struct {
	ID             string `json:"id,omitempty"`
	ChainID        string `json:"chain_id,omitempty"`
	OutputClassID  string `json:"output_class_id,omitempty"`
	Module         Module `json:"module"`
	ProvisionState string `json:"provision_state"`
	Relation       string `json:"relation"`
	Index          int    `json:"index"`
}

// ThoughtResponse represents the API response for thought operations.
type ThoughtResponse struct {
	Data Thought `json:"data"`
}

// CreateThoughtRequest represents the request payload for creating a thought.
type CreateThoughtRequest struct {
	Thought ThoughtRequestData `json:"thought"`
}

// ThoughtRequestData represents the thought data in the request.
type ThoughtRequestData struct {
	Relation      string `json:"relation"`
	OutputClassID string `json:"output_class_id,omitempty"`
	Index         int    `json:"index,omitempty"`
	Module        Module `json:"module"`
}

// UpdateThoughtRequest represents the request payload for updating a thought.
type UpdateThoughtRequest struct {
	Thought UpdateThoughtData `json:"thought"`
}

// UpdateThoughtData represents the thought update data.
type UpdateThoughtData struct {
	Relation      string `json:"relation,omitempty"`
	OutputClassID string `json:"output_class_id,omitempty"`
	Index         int    `json:"index,omitempty"`
	Module        Module `json:"module,omitempty"`
}

// GetThought retrieves a specific thought by ID.
// GET /provision/perception/thoughts/:id.
func (s *Service) GetThought(id string) (*Thought, error) {
	if id == "" {
		return nil, errors.New("thought ID is required")
	}

	var thoughtResp ThoughtResponse
	resp, err := s.client.R().
		SetResult(&thoughtResp).
		Get(fmt.Sprintf("/provision/perception/thoughts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get thought: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &thoughtResp.Data, nil
}

// CreateThought creates a new thought within a chain.
// POST /provision/perception/chains/:chain_id/thoughts.
func (s *Service) CreateThought(chainID string, req CreateThoughtRequest) (*Thought, error) {
	if chainID == "" {
		return nil, errors.New("chain ID is required")
	}
	if req.Thought.Relation == "" {
		return nil, errors.New("thought relation is required")
	}
	if req.Thought.Module.Reference == "" {
		return nil, errors.New("thought module reference is required")
	}

	var thoughtResp ThoughtResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&thoughtResp).
		Post(fmt.Sprintf("/provision/perception/chains/%s/thoughts", chainID))

	if err != nil {
		return nil, fmt.Errorf("failed to create thought: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &thoughtResp.Data, nil
}

// UpdateThought updates an existing thought using PATCH.
// PATCH /provision/perception/thoughts/:id.
func (s *Service) UpdateThought(id string, req UpdateThoughtRequest) (*Thought, error) {
	if id == "" {
		return nil, errors.New("thought ID is required")
	}

	var thoughtResp ThoughtResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&thoughtResp).
		Patch(fmt.Sprintf("/provision/perception/thoughts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update thought: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &thoughtResp.Data, nil
}

// DeleteThought deletes a thought by ID.
// DELETE /provision/perception/thoughts/:id.
func (s *Service) DeleteThought(id string) error {
	if id == "" {
		return errors.New("thought ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/thoughts/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete thought: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
