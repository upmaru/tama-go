package sensory

import (
	"fmt"
)

// This file contains all Source-related operations for the SensoryService.
// Sources represent data sources with endpoints and credentials.

// Source operations

// GetSource retrieves a specific source by ID
// GET /provision/sensory/sources/:id
func (s *Service) GetSource(id string) (*Source, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetResult(&sourceResp).
		SetError(&Error{}).
		Get(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &sourceResp.Data, nil
}

// CreateSource creates a new source in a specific space
// POST /provision/sensory/spaces/:space_id/sources
func (s *Service) CreateSource(spaceID string, req CreateSourceRequest) (*Source, error) {
	if spaceID == "" {
		return nil, fmt.Errorf("space ID is required")
	}
	if req.Source.Name == "" {
		return nil, fmt.Errorf("source name is required")
	}
	if req.Source.Type == "" {
		return nil, fmt.Errorf("source type is required")
	}
	if req.Source.Endpoint == "" {
		return nil, fmt.Errorf("source endpoint is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&sourceResp).
		SetError(&Error{}).
		Post(fmt.Sprintf("/provision/sensory/spaces/%s/sources", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &sourceResp.Data, nil
}

// UpdateSource updates an existing source using PATCH
// PATCH /provision/sensory/sources/:id
func (s *Service) UpdateSource(id string, req UpdateSourceRequest) (*Source, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&sourceResp).
		SetError(&Error{}).
		Patch(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update source: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &sourceResp.Data, nil
}

// ReplaceSource replaces an existing source using PUT
// PUT /provision/sensory/sources/:id
func (s *Service) ReplaceSource(id string, req UpdateSourceRequest) (*Source, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID is required")
	}

	var sourceResp SourceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&sourceResp).
		SetError(&Error{}).
		Put(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace source: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &sourceResp.Data, nil
}

// DeleteSource deletes a source by ID
// DELETE /provision/sensory/sources/:id
func (s *Service) DeleteSource(id string) error {
	if id == "" {
		return fmt.Errorf("source ID is required")
	}

	resp, err := s.client.R().
		SetError(&Error{}).
		Delete(fmt.Sprintf("/provision/sensory/sources/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return errorResp
		}
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}
