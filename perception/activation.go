//nolint:dupl // activation and tool follow similar CRUD patterns
package perception

import "errors"

// Activation represents a perception activation resource.
type Activation struct {
	ID             string `json:"id,omitempty"`
	ThoughtPathID  string `json:"thought_path_id,omitempty"`
	ChainID        string `json:"chain_id,omitempty"`
	ProvisionState string `json:"provision_state"`
}

// ActivationResponse represents the API response for activation operations.
type ActivationResponse struct {
	Data Activation `json:"data"`
}

// CreateActivationRequest represents the request payload for creating an activation.
type CreateActivationRequest struct {
	Activation ActivationRequestData `json:"activation"`
}

// ActivationRequestData represents the activation data in the request.
type ActivationRequestData struct {
	ChainID string `json:"chain_id"`
}

// UpdateActivationRequest represents the request payload for updating an activation.
type UpdateActivationRequest struct {
	Activation UpdateActivationData `json:"activation"`
}

// UpdateActivationData represents the activation update data.
type UpdateActivationData struct {
	ChainID string `json:"chain_id,omitempty"`
}

// GetActivation retrieves a specific activation by ID.
// GET /provision/perception/activations/:id.
func (s *Service) GetActivation(id string) (*Activation, error) {
	var activationResp ActivationResponse
	if err := genericGet(s, id, "activation", "/provision/perception/activations/%s", &activationResp); err != nil {
		return nil, err
	}
	return &activationResp.Data, nil
}

// CreateActivation creates a new activation within a path.
// POST /provision/perception/paths/:path_id/activations.
func (s *Service) CreateActivation(pathID string, req CreateActivationRequest) (*Activation, error) {
	if req.Activation.ChainID == "" {
		return nil, errors.New("chain ID is required")
	}

	var activationResp ActivationResponse
	if err := genericCreate(
		s,
		pathID,
		req,
		"activation",
		"path",
		"/provision/perception/paths/%s/activations",
		&activationResp,
	); err != nil {
		return nil, err
	}
	return &activationResp.Data, nil
}

// UpdateActivation updates an existing activation using PATCH.
// PATCH /provision/perception/activations/:id.
func (s *Service) UpdateActivation(id string, req UpdateActivationRequest) (*Activation, error) {
	var activationResp ActivationResponse
	if err := genericUpdate(
		s,
		id,
		req,
		"activation",
		"/provision/perception/activations/%s",
		&activationResp,
	); err != nil {
		return nil, err
	}
	return &activationResp.Data, nil
}

// ReplaceActivation replaces an existing activation using PUT.
// PUT /provision/perception/activations/:id.
func (s *Service) ReplaceActivation(id string, req UpdateActivationRequest) (*Activation, error) {
	var activationResp ActivationResponse
	if err := genericReplace(
		s,
		id,
		req,
		"activation",
		"/provision/perception/activations/%s",
		&activationResp,
	); err != nil {
		return nil, err
	}
	return &activationResp.Data, nil
}

// DeleteActivation deletes an activation by ID.
// DELETE /provision/perception/activations/:id.
func (s *Service) DeleteActivation(id string) error {
	return genericDelete(s, id, "activation", "/provision/perception/activations/%s")
}
