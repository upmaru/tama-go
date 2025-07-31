//nolint:dupl // CRUD operations follow the same pattern as path.go
package perception

import (
	"errors"
	"fmt"
)

// Context represents a perception context resource.
type Context struct {
	ID             string `json:"id,omitempty"`
	ThoughtID      string `json:"thought_id,omitempty"`
	PromptID       string `json:"prompt_id"`
	Layer          int    `json:"layer"`
	ProvisionState string `json:"provision_state"`
}

// ContextResponse represents the API response for context operations.
type ContextResponse struct {
	Data Context `json:"data"`
}

// CreateContextRequest represents the request payload for creating a context.
type CreateContextRequest struct {
	Context ContextRequestData `json:"context"`
}

// ContextRequestData represents the context data in the request.
type ContextRequestData struct {
	PromptID string `json:"prompt_id"`
	Layer    int    `json:"layer"`
}

// UpdateContextRequest represents the request payload for updating a context.
type UpdateContextRequest struct {
	Context UpdateContextData `json:"context"`
}

// UpdateContextData represents the context update data.
type UpdateContextData struct {
	PromptID string `json:"prompt_id,omitempty"`
	Layer    int    `json:"layer,omitempty"`
}

// GetContext retrieves a specific context by ID.
// GET /provision/perception/contexts/:id.
func (s *Service) GetContext(id string) (*Context, error) {
	if id == "" {
		return nil, errors.New("context ID is required")
	}

	var contextResp ContextResponse
	resp, err := s.client.R().
		SetResult(&contextResp).
		Get(fmt.Sprintf("/provision/perception/contexts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get context: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &contextResp.Data, nil
}

// CreateContext creates a new context within a thought.
// POST /provision/perception/thoughts/:thought_id/contexts.
func (s *Service) CreateContext(thoughtID string, req CreateContextRequest) (*Context, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if req.Context.PromptID == "" {
		return nil, errors.New("context prompt ID is required")
	}

	var contextResp ContextResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&contextResp).
		Post(fmt.Sprintf("/provision/perception/thoughts/%s/contexts", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &contextResp.Data, nil
}

// UpdateContext updates an existing context using PATCH.
// PATCH /provision/perception/contexts/:id.
func (s *Service) UpdateContext(id string, req UpdateContextRequest) (*Context, error) {
	if id == "" {
		return nil, errors.New("context ID is required")
	}

	var contextResp ContextResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&contextResp).
		Patch(fmt.Sprintf("/provision/perception/contexts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update context: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &contextResp.Data, nil
}

// ReplaceContext replaces an existing context using PUT.
// PUT /provision/perception/contexts/:id.
func (s *Service) ReplaceContext(id string, req UpdateContextRequest) (*Context, error) {
	if id == "" {
		return nil, errors.New("context ID is required")
	}

	var contextResp ContextResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&contextResp).
		Put(fmt.Sprintf("/provision/perception/contexts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace context: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &contextResp.Data, nil
}

// DeleteContext deletes a context by ID.
// DELETE /provision/perception/contexts/:id.
func (s *Service) DeleteContext(id string) error {
	if id == "" {
		return errors.New("context ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/contexts/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete context: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
