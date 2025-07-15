package neural

import (
	"fmt"
)

// GetSpace retrieves a specific space by ID
// GET /provision/neural/spaces/:id
func (s *Service) GetSpace(id string) (*Space, error) {
	if id == "" {
		return nil, fmt.Errorf("space ID is required")
	}

	var spaceResp SpaceResponse
	resp, err := s.client.R().
		SetResult(&spaceResp).
		SetError(&Error{}).
		Get(fmt.Sprintf("/provision/neural/spaces/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get space: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &spaceResp.Data, nil
}

// CreateSpace creates a new space
// POST /provision/neural/spaces
func (s *Service) CreateSpace(req CreateSpaceRequest) (*Space, error) {
	if req.Space.Name == "" {
		return nil, fmt.Errorf("space name is required")
	}
	if req.Space.Type == "" {
		return nil, fmt.Errorf("space type is required")
	}
	if req.Space.Type != "root" && req.Space.Type != "component" {
		return nil, fmt.Errorf("space type must be 'root' or 'component'")
	}

	var spaceResp SpaceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&spaceResp).
		SetError(&Error{}).
		Post("/provision/neural/spaces")

	if err != nil {
		return nil, fmt.Errorf("failed to create space: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &spaceResp.Data, nil
}

// UpdateSpace updates an existing space using PATCH
// PATCH /provision/neural/spaces/:id
func (s *Service) UpdateSpace(id string, req UpdateSpaceRequest) (*Space, error) {
	if id == "" {
		return nil, fmt.Errorf("space ID is required")
	}

	var spaceResp SpaceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&spaceResp).
		SetError(&Error{}).
		Patch(fmt.Sprintf("/provision/neural/spaces/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update space: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &spaceResp.Data, nil
}

// ReplaceSpace replaces an existing space using PUT
// PUT /provision/neural/spaces/:id
func (s *Service) ReplaceSpace(id string, req UpdateSpaceRequest) (*Space, error) {
	if id == "" {
		return nil, fmt.Errorf("space ID is required")
	}

	var spaceResp SpaceResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&spaceResp).
		SetError(&Error{}).
		Put(fmt.Sprintf("/provision/neural/spaces/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace space: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &spaceResp.Data, nil
}

// DeleteSpace deletes a space by ID
// DELETE /provision/neural/spaces/:id
func (s *Service) DeleteSpace(id string) error {
	if id == "" {
		return fmt.Errorf("space ID is required")
	}

	resp, err := s.client.R().
		SetError(&Error{}).
		Delete(fmt.Sprintf("/provision/neural/spaces/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete space: %w", err)
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
