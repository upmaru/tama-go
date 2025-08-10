package neural

import (
	"errors"
	"fmt"
)

// Class represents a neural class resource.
type Class struct {
	ID             string         `json:"id,omitempty"`
	SpaceID        string         `json:"space_id,omitempty"`
	ProvisionState string         `json:"provision_state"`
	Schema         map[string]any `json:"schema"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
}

// ClassResponse represents the API response for class operations.
type ClassResponse struct {
	Data Class `json:"data"`
}

// CreateClassRequest represents the request payload for creating a class.
type CreateClassRequest struct {
	Class ClassRequestData `json:"class"`
}

// ClassRequestData represents the class data in the request.
type ClassRequestData struct {
	Schema map[string]any `json:"schema"`
}

// UpdateClassRequest represents the request payload for updating a class.
type UpdateClassRequest struct {
	Class UpdateClassData `json:"class"`
}

// UpdateClassData represents the class update data.
type UpdateClassData struct {
	Schema map[string]any `json:"schema,omitempty"`
}

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

// GetClassBySpecificationAndName retrieves a specific class by specification ID and name.
// GET /provision/neural/specifications/:specification_id/classes/:name.
func (s *Service) GetClassBySpecificationAndName(specificationID string, name string) (*Class, error) {
	if specificationID == "" {
		return nil, errors.New("specification ID is required")
	}
	if name == "" {
		return nil, errors.New("class name is required")
	}

	var classResp ClassResponse
	resp, err := s.client.R().
		SetResult(&classResp).
		Get(fmt.Sprintf("/provision/neural/specifications/%s/classes/%s", specificationID, name))

	if err != nil {
		return nil, fmt.Errorf("failed to get class by specification and name: %w", err)
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
