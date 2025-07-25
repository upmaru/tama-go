package sensory

import (
	"errors"
	"fmt"
)

// This file contains all Identity-related operations for the SensoryService.
// Identities represent provisioned instances of specifications with validation endpoints.

// Identity operations

// GetIdentity retrieves a specific identity by ID.
// GET /provision/sensory/identities/:id.
func (s *Service) GetIdentity(id string) (*Identity, error) {
	if id == "" {
		return nil, errors.New("identity ID is required")
	}

	var identityResp IdentityResponse
	resp, err := s.client.R().
		SetResult(&identityResp).
		Get(fmt.Sprintf("/provision/sensory/identities/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &identityResp.Data, nil
}

// CreateIdentity creates a new identity for a specific specification and identifier.
// POST /provision/sensory/specifications/:specification_id/identifiers/:identifier/identities.
func (s *Service) CreateIdentity(specificationID, identifier string, req CreateIdentityRequest) (*Identity, error) {
	if specificationID == "" {
		return nil, errors.New("specification ID is required")
	}
	if identifier == "" {
		return nil, errors.New("identifier is required")
	}
	if req.Identity.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	if req.Identity.Validation.Path == "" {
		return nil, errors.New("validation path is required")
	}
	if req.Identity.Validation.Method == "" {
		return nil, errors.New("validation method is required")
	}
	if len(req.Identity.Validation.Codes) == 0 {
		return nil, errors.New("validation codes are required")
	}

	var identityResp IdentityResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&identityResp).
		Post(fmt.Sprintf("/provision/sensory/specifications/%s/identifiers/%s/identities", specificationID, identifier))

	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &identityResp.Data, nil
}

// UpdateIdentity updates an existing identity using PATCH.
// PATCH /provision/sensory/identities/:id.
func (s *Service) UpdateIdentity(id string, req UpdateIdentityRequest) (*Identity, error) {
	if id == "" {
		return nil, errors.New("identity ID is required")
	}

	var identityResp IdentityResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&identityResp).
		Patch(fmt.Sprintf("/provision/sensory/identities/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update identity: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &identityResp.Data, nil
}

// ReplaceIdentity replaces an existing identity using PUT.
// PUT /provision/sensory/identities/:id.
func (s *Service) ReplaceIdentity(id string, req UpdateIdentityRequest) (*Identity, error) {
	if id == "" {
		return nil, errors.New("identity ID is required")
	}

	var identityResp IdentityResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&identityResp).
		Put(fmt.Sprintf("/provision/sensory/identities/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace identity: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &identityResp.Data, nil
}

// DeleteIdentity deletes an identity by ID.
// DELETE /provision/sensory/identities/:id.
func (s *Service) DeleteIdentity(id string) error {
	if id == "" {
		return errors.New("identity ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/sensory/identities/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete identity: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
