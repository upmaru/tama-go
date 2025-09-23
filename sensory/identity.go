package sensory

import (
	"errors"
	"fmt"
)

// This file contains all Identity-related operations for the SensoryService.
// Identities represent provisioned instances of specifications with validation endpoints.

// Validation represents the validation configuration for an identity.
type Validation struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Codes  []int  `json:"codes"`
}

// Identity represents a sensory identity resource.
type Identity struct {
	ID              string     `json:"id,omitempty"`
	SpecificationID string     `json:"specification_id"`
	ProvisionState  string     `json:"provision_state"`
	CurrentState    string     `json:"current_state"`
	Identifier      string     `json:"identifier"`
	Validation      Validation `json:"validation"`
}

// IdentityResponse represents the API response for identity operations.
type IdentityResponse struct {
	Data Identity `json:"data"`
}

// CreateIdentityRequest represents the request payload for creating an identity.
type CreateIdentityRequest struct {
	Identity IdentityRequestData `json:"identity"`
}

// IdentityRequestData represents the identity data in the request.
type IdentityRequestData struct {
	APIKey       string     `json:"api_key,omitempty"`
	ClientID     string     `json:"client_id,omitempty"`
	ClientSecret string     `json:"client_secret,omitempty"`
	Validation   Validation `json:"validation"`
}

// UpdateIdentityRequest represents the request payload for updating an identity.
type UpdateIdentityRequest struct {
	Identity UpdateIdentityData `json:"identity"`
}

// UpdateIdentityData represents the identity update data.
type UpdateIdentityData struct {
	APIKey       string      `json:"api_key,omitempty"`
	ClientID     string      `json:"client_id,omitempty"`
	ClientSecret string      `json:"client_secret,omitempty"`
	Validation   *Validation `json:"validation,omitempty"`
}

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
	// Validate authentication: either API key OR client credentials must be provided
	hasAPIKey := req.Identity.APIKey != ""
	hasClientCredentials := req.Identity.ClientID != "" && req.Identity.ClientSecret != ""

	if !hasAPIKey && !hasClientCredentials {
		return nil, errors.New("either API key or client credentials (client_id and client_secret) are required")
	}

	if hasAPIKey && hasClientCredentials {
		return nil, errors.New("provide either API key or client credentials, not both")
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
