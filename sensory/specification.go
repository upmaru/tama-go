package sensory

import (
	"errors"
	"fmt"
)

// This file contains all Specification-related operations for the SensoryService.
// Specifications represent configuration schemas with endpoints and versioning.

// Specification represents a sensory specification resource.
type Specification struct {
	ID             string         `json:"id,omitempty"`
	SpaceID        string         `json:"space_id"`
	Schema         map[string]any `json:"schema"`
	Version        string         `json:"version"`
	Endpoint       string         `json:"endpoint"`
	CurrentState   string         `json:"current_state"`
	ProvisionState string         `json:"provision_state"`
}

// SpecificationResponse represents the API response for specification operations.
type SpecificationResponse struct {
	Data Specification `json:"data"`
}

// CreateSpecificationRequest represents the request payload for creating a specification.
type CreateSpecificationRequest struct {
	Specification SpecificationRequestData `json:"specification"`
}

// SpecificationRequestData represents the specification data in the request.
type SpecificationRequestData struct {
	Schema   map[string]any `json:"schema"`
	Version  string         `json:"version"`
	Endpoint string         `json:"endpoint"`
}

// UpdateSpecificationRequest represents the request payload for updating a specification.
type UpdateSpecificationRequest struct {
	Specification UpdateSpecificationData `json:"specification"`
}

// UpdateSpecificationData represents the specification update data.
type UpdateSpecificationData struct {
	Schema   map[string]any `json:"schema,omitempty"`
	Version  string         `json:"version,omitempty"`
	Endpoint string         `json:"endpoint,omitempty"`
}

// Specification operations

// GetSpecification retrieves a specific specification by ID.
// GET /provision/sensory/specifications/:id.
func (s *Service) GetSpecification(id string) (*Specification, error) {
	if id == "" {
		return nil, errors.New("specification ID is required")
	}

	var specResp SpecificationResponse
	resp, err := s.client.R().
		SetResult(&specResp).
		Get(fmt.Sprintf("/provision/sensory/specifications/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get specification: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &specResp.Data, nil
}

// CreateSpecification creates a new specification in a specific space.
// POST /provision/sensory/spaces/:space_id/specifications.
func (s *Service) CreateSpecification(spaceID string, req CreateSpecificationRequest) (*Specification, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if req.Specification.Schema == nil {
		return nil, errors.New("specification schema is required")
	}
	if req.Specification.Version == "" {
		return nil, errors.New("specification version is required")
	}
	if req.Specification.Endpoint == "" {
		return nil, errors.New("specification endpoint is required")
	}

	var specResp SpecificationResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&specResp).
		Post(fmt.Sprintf("/provision/sensory/spaces/%s/specifications", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create specification: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &specResp.Data, nil
}

// UpdateSpecification updates an existing specification using PATCH.
// PATCH /provision/sensory/specifications/:id.
func (s *Service) UpdateSpecification(id string, req UpdateSpecificationRequest) (*Specification, error) {
	if id == "" {
		return nil, errors.New("specification ID is required")
	}

	var specResp SpecificationResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&specResp).
		Patch(fmt.Sprintf("/provision/sensory/specifications/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update specification: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &specResp.Data, nil
}

// ReplaceSpecification replaces an existing specification using PUT.
// PUT /provision/sensory/specifications/:id.
func (s *Service) ReplaceSpecification(id string, req UpdateSpecificationRequest) (*Specification, error) {
	if id == "" {
		return nil, errors.New("specification ID is required")
	}

	var specResp SpecificationResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&specResp).
		Put(fmt.Sprintf("/provision/sensory/specifications/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace specification: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &specResp.Data, nil
}

// DeleteSpecification deletes a specification by ID.
// DELETE /provision/sensory/specifications/:id.
func (s *Service) DeleteSpecification(id string) error {
	if id == "" {
		return errors.New("specification ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/sensory/specifications/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete specification: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
