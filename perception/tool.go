package perception

import (
	"errors"
	"fmt"
)

// Tool represents a perception tool resource.
type Tool struct {
	ID             string `json:"id,omitempty"`
	ThoughtID      string `json:"thought_id"`
	ActionID       string `json:"action_id"`
	ProvisionState string `json:"provision_state"`
}

// ToolResponse represents the API response for tool operations.
type ToolResponse struct {
	Data Tool `json:"data"`
}

// CreateToolRequest represents the request payload for creating a tool.
type CreateToolRequest struct {
	Tool CreateToolData `json:"tool"`
}

// CreateToolData represents the tool data in the request.
type CreateToolData struct {
	ActionID string `json:"action_id"`
}

// UpdateToolRequest represents the request payload for updating a tool.
type UpdateToolRequest struct {
	Tool UpdateToolData `json:"tool"`
}

// UpdateToolData represents the tool update data.
type UpdateToolData struct {
	ActionID string `json:"action_id,omitempty"`
}

// GetTool retrieves a specific tool by ID.
// GET /provision/perception/tools/:id.
func (s *Service) GetTool(id string) (*Tool, error) {
	if id == "" {
		return nil, errors.New("tool ID is required")
	}

	var toolResp ToolResponse
	resp, err := s.client.R().
		SetResult(&toolResp).
		Get(fmt.Sprintf("/provision/perception/tools/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &toolResp.Data, nil
}

// CreateTool creates a new tool within a thought.
// POST /provision/perception/thoughts/:thought_id/tools.
func (s *Service) CreateTool(thoughtID string, req CreateToolRequest) (*Tool, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if req.Tool.ActionID == "" {
		return nil, errors.New("action ID is required")
	}

	var toolResp ToolResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&toolResp).
		Post(fmt.Sprintf("/provision/perception/thoughts/%s/tools", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to create tool: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &toolResp.Data, nil
}

// UpdateTool updates an existing tool using PATCH.
// PATCH /provision/perception/tools/:id.
func (s *Service) UpdateTool(id string, req UpdateToolRequest) (*Tool, error) {
	if id == "" {
		return nil, errors.New("tool ID is required")
	}

	var toolResp ToolResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&toolResp).
		Patch(fmt.Sprintf("/provision/perception/tools/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update tool: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &toolResp.Data, nil
}

// ReplaceTool replaces an existing tool using PUT.
// PUT /provision/perception/tools/:id.
func (s *Service) ReplaceTool(id string, req UpdateToolRequest) (*Tool, error) {
	if id == "" {
		return nil, errors.New("tool ID is required")
	}

	var toolResp ToolResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&toolResp).
		Put(fmt.Sprintf("/provision/perception/tools/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace tool: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &toolResp.Data, nil
}

// DeleteTool deletes a tool by ID.
// DELETE /provision/perception/tools/:id.
func (s *Service) DeleteTool(id string) error {
	if id == "" {
		return errors.New("tool ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/tools/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete tool: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
