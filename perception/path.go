//nolint:dupl // CRUD operations follow the same pattern as chain.go
package perception

import (
	"errors"
	"fmt"
)

// Path represents a perception path resource.
type Path struct {
	ID             string         `json:"id,omitempty"`
	ThoughtID      string         `json:"thought_id,omitempty"`
	ProvisionState string         `json:"provision_state"`
	TargetClassID  string         `json:"target_class_id"`
	Parameters     map[string]any `json:"parameters,omitempty"`
}

// PathResponse represents the API response for path operations.
type PathResponse struct {
	Data Path `json:"data"`
}

// CreatePathRequest represents the request payload for creating a path.
type CreatePathRequest struct {
	Path PathRequestData `json:"path"`
}

// PathRequestData represents the path data in the request.
type PathRequestData struct {
	TargetClassID string         `json:"target_class_id"`
	Parameters    map[string]any `json:"parameters,omitempty"`
}

// UpdatePathRequest represents the request payload for updating a path.
type UpdatePathRequest struct {
	Path UpdatePathData `json:"path"`
}

// UpdatePathData represents the path update data.
type UpdatePathData struct {
	TargetClassID string         `json:"target_class_id,omitempty"`
	Parameters    map[string]any `json:"parameters,omitempty"`
}

// GetPath retrieves a specific path by ID.
// GET /provision/perception/paths/:id.
func (s *Service) GetPath(id string) (*Path, error) {
	if id == "" {
		return nil, errors.New("path ID is required")
	}

	var pathResp PathResponse
	resp, err := s.client.R().
		SetResult(&pathResp).
		Get(fmt.Sprintf("/provision/perception/paths/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get path: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &pathResp.Data, nil
}

// CreatePath creates a new path within a thought.
// POST /provision/perception/thoughts/:thought_id/paths.
func (s *Service) CreatePath(thoughtID string, req CreatePathRequest) (*Path, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if req.Path.TargetClassID == "" {
		return nil, errors.New("path target class ID is required")
	}

	var pathResp PathResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&pathResp).
		Post(fmt.Sprintf("/provision/perception/thoughts/%s/paths", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to create path: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &pathResp.Data, nil
}

// UpdatePath updates an existing path using PATCH.
// PATCH /provision/perception/paths/:id.
func (s *Service) UpdatePath(id string, req UpdatePathRequest) (*Path, error) {
	if id == "" {
		return nil, errors.New("path ID is required")
	}

	var pathResp PathResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&pathResp).
		Patch(fmt.Sprintf("/provision/perception/paths/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update path: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &pathResp.Data, nil
}

// ReplacePath replaces an existing path using PUT.
// PUT /provision/perception/paths/:id.
func (s *Service) ReplacePath(id string, req UpdatePathRequest) (*Path, error) {
	if id == "" {
		return nil, errors.New("path ID is required")
	}

	var pathResp PathResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&pathResp).
		Put(fmt.Sprintf("/provision/perception/paths/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace path: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &pathResp.Data, nil
}

// DeletePath deletes a path by ID.
// DELETE /provision/perception/paths/:id.
func (s *Service) DeletePath(id string) error {
	if id == "" {
		return errors.New("path ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/paths/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete path: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
