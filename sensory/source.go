package sensory

import (
	"errors"
	"fmt"
)

// This file contains all Source-related operations for the SensoryService.
// Sources represent data sources with endpoints and credentials.

// SourceCredential represents the credential structure for sources.
type SourceCredential struct {
	APIKey string `json:"api_key"`
}

// Source represents a sensory source resource.
type Source struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Endpoint       string `json:"endpoint"`
	SpaceID        string `json:"space_id"`
	Type           string `json:"type"`
	ProvisionState string `json:"provision_state"`
}

// SourceResponse represents the API response for source operations.
type SourceResponse struct {
	Data Source `json:"data"`
}

// CreateSourceRequest represents the request payload for creating a source.
type CreateSourceRequest struct {
	Source SourceRequestData `json:"source"`
}

// SourceRequestData represents the source data in the request.
type SourceRequestData struct {
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Endpoint   string           `json:"endpoint"`
	Credential SourceCredential `json:"credential"`
}

// UpdateSourceRequest represents the request payload for updating a source.
type UpdateSourceRequest struct {
	Source UpdateSourceData `json:"source"`
}

// UpdateSourceData represents the source update data.
type UpdateSourceData struct {
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Endpoint   string            `json:"endpoint,omitempty"`
	Credential *SourceCredential `json:"credential,omitempty"`
}

// Source operations

// GetSource retrieves a specific source by ID.
// GET /provision/sensory/sources/:id.
func (s *Service) GetSource(id string) (*Source, error) {
	if id == "" {
		return nil, errors.New("source ID is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetResult(&sourceResp).
		Get(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &sourceResp.Data, nil
}

// GetBySpecificationAndSlug retrieves a source by specification ID and source slug.
// GET /provision/sensory/specifications/:specification_id/sources/:id.
func (s *Service) GetSourceBySpecificationAndSlug(specificationID string, slug string) (*Source, error) {
	if specificationID == "" {
		return nil, errors.New("specification ID is required")
	}
	if slug == "" {
		return nil, errors.New("source slug is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetResult(&sourceResp).
		Get(fmt.Sprintf("/provision/sensory/specifications/%s/sources/%s", specificationID, slug))

	if err != nil {
		return nil, fmt.Errorf("failed to get source by specification and slug: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &sourceResp.Data, nil
}

// CreateSource creates a new source in a specific space.
// POST /provision/sensory/spaces/:space_id/sources.
func (s *Service) CreateSource(spaceID string, req CreateSourceRequest) (*Source, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if req.Source.Name == "" {
		return nil, errors.New("source name is required")
	}
	if req.Source.Type == "" {
		return nil, errors.New("source type is required")
	}
	if req.Source.Endpoint == "" {
		return nil, errors.New("source endpoint is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&sourceResp).
		Post(fmt.Sprintf("/provision/sensory/spaces/%s/sources", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &sourceResp.Data, nil
}

// UpdateSource updates an existing source using PATCH.
// PATCH /provision/sensory/sources/:id.
func (s *Service) UpdateSource(id string, req UpdateSourceRequest) (*Source, error) {
	if id == "" {
		return nil, errors.New("source ID is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&sourceResp).
		Patch(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update source: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &sourceResp.Data, nil
}

// ReplaceSource replaces an existing source using PUT.
// PUT /provision/sensory/sources/:id.
func (s *Service) ReplaceSource(id string, req UpdateSourceRequest) (*Source, error) {
	if id == "" {
		return nil, errors.New("source ID is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&sourceResp).
		Put(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace source: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &sourceResp.Data, nil
}

// DeleteSource deletes a source by ID.
// DELETE /provision/sensory/sources/:id.
func (s *Service) DeleteSource(id string) error {
	if id == "" {
		return errors.New("source ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
