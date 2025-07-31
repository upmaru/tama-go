package perception

import (
	"errors"
	"fmt"
)

// Chain represents a perception chain resource.
type Chain struct {
	ID             string `json:"id,omitempty"`
	SpaceID        string `json:"space_id,omitempty"`
	Name           string `json:"name"`
	Slug           string `json:"slug,omitempty"`
	ProvisionState string `json:"provision_state"`
}

// ChainResponse represents the API response for chain operations.
type ChainResponse struct {
	Data Chain `json:"data"`
}

// CreateChainRequest represents the request payload for creating a chain.
type CreateChainRequest struct {
	Chain ChainRequestData `json:"chain"`
}

// ChainRequestData represents the chain data in the request.
type ChainRequestData struct {
	Name string `json:"name"`
}

// UpdateChainRequest represents the request payload for updating a chain.
type UpdateChainRequest struct {
	Chain UpdateChainData `json:"chain"`
}

// UpdateChainData represents the chain update data.
type UpdateChainData struct {
	Name string `json:"name,omitempty"`
}

// GetChain retrieves a specific chain by ID.
// GET /provision/perception/chains/:id.
func (s *Service) GetChain(id string) (*Chain, error) {
	if id == "" {
		return nil, errors.New("chain ID is required")
	}

	var chainResp ChainResponse
	resp, err := s.client.R().
		SetResult(&chainResp).
		Get(fmt.Sprintf("/provision/perception/chains/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get chain: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &chainResp.Data, nil
}

// CreateChain creates a new chain within a space.
// POST /provision/perception/spaces/:space_id/chains.
func (s *Service) CreateChain(spaceID string, req CreateChainRequest) (*Chain, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if req.Chain.Name == "" {
		return nil, errors.New("chain name is required")
	}

	var chainResp ChainResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&chainResp).
		Post(fmt.Sprintf("/provision/perception/spaces/%s/chains", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &chainResp.Data, nil
}

// UpdateChain updates an existing chain using PATCH.
// PATCH /provision/perception/chains/:id.
func (s *Service) UpdateChain(id string, req UpdateChainRequest) (*Chain, error) {
	if id == "" {
		return nil, errors.New("chain ID is required")
	}

	var chainResp ChainResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&chainResp).
		Patch(fmt.Sprintf("/provision/perception/chains/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update chain: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &chainResp.Data, nil
}

// ReplaceChain replaces an existing chain using PUT.
// PUT /provision/perception/chains/:id.
func (s *Service) ReplaceChain(id string, req UpdateChainRequest) (*Chain, error) {
	if id == "" {
		return nil, errors.New("chain ID is required")
	}

	var chainResp ChainResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&chainResp).
		Put(fmt.Sprintf("/provision/perception/chains/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace chain: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &chainResp.Data, nil
}

// DeleteChain deletes a chain by ID.
// DELETE /provision/perception/chains/:id.
func (s *Service) DeleteChain(id string) error {
	if id == "" {
		return errors.New("chain ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/chains/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete chain: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
