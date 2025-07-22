package neural

import (
	"errors"
	"fmt"
)

// GetClass retrieves a specific class by ID.
// GET /provision/neural/classes/:id.
func (s *Service) GetClass(id string) (*Class, error) {
	if id == "" {
		return nil, errors.New("class ID is required")
	}

	var classResp ClassResponse
	resp, err := s.client.R().
		SetResult(&classResp).
		Get(fmt.Sprintf("/provision/neural/classes/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get class: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &classResp.Data, nil
}

// CreateClass creates a new class within a space.
// POST /provision/neural/spaces/:space_id/classes.
func (s *Service) CreateClass(spaceID string, req CreateClassRequest) (*Class, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if req.Class.Schema == nil {
		return nil, errors.New("class schema is required")
	}

	var classResp ClassResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&classResp).
		Post(fmt.Sprintf("/provision/neural/spaces/%s/classes", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create class: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &classResp.Data, nil
}

// UpdateClass updates an existing class using PATCH.
// PATCH /provision/neural/classes/:id.
func (s *Service) UpdateClass(id string, req UpdateClassRequest) (*Class, error) {
	if id == "" {
		return nil, errors.New("class ID is required")
	}

	var classResp ClassResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&classResp).
		Patch(fmt.Sprintf("/provision/neural/classes/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update class: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &classResp.Data, nil
}

// ReplaceClass replaces an existing class using PUT.
// PUT /provision/neural/classes/:id.
func (s *Service) ReplaceClass(id string, req UpdateClassRequest) (*Class, error) {
	if id == "" {
		return nil, errors.New("class ID is required")
	}

	var classResp ClassResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&classResp).
		Put(fmt.Sprintf("/provision/neural/classes/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace class: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &classResp.Data, nil
}

// DeleteClass deletes a class by ID.
// DELETE /provision/neural/classes/:id.
func (s *Service) DeleteClass(id string) error {
	if id == "" {
		return errors.New("class ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/neural/classes/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
