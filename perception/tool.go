//nolint:dupl // activation and tool follow similar CRUD patterns
package perception

import "errors"

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
	var toolResp ToolResponse
	if err := genericGet(s, id, "tool", "/provision/perception/tools/%s", &toolResp); err != nil {
		return nil, err
	}
	return &toolResp.Data, nil
}

// CreateTool creates a new tool within a thought.
// POST /provision/perception/thoughts/:thought_id/tools.
func (s *Service) CreateTool(thoughtID string, req CreateToolRequest) (*Tool, error) {
	if req.Tool.ActionID == "" {
		return nil, errors.New("action ID is required")
	}

	var toolResp ToolResponse
	if err := genericCreate(
		s,
		thoughtID,
		req,
		"tool",
		"thought",
		"/provision/perception/thoughts/%s/tools",
		&toolResp,
	); err != nil {
		return nil, err
	}
	return &toolResp.Data, nil
}

// UpdateTool updates an existing tool using PATCH.
// PATCH /provision/perception/tools/:id.
func (s *Service) UpdateTool(id string, req UpdateToolRequest) (*Tool, error) {
	var toolResp ToolResponse
	if err := genericUpdate(s, id, req, "tool", "/provision/perception/tools/%s", &toolResp); err != nil {
		return nil, err
	}
	return &toolResp.Data, nil
}

// ReplaceTool replaces an existing tool using PUT.
// PUT /provision/perception/tools/:id.
func (s *Service) ReplaceTool(id string, req UpdateToolRequest) (*Tool, error) {
	var toolResp ToolResponse
	if err := genericReplace(s, id, req, "tool", "/provision/perception/tools/%s", &toolResp); err != nil {
		return nil, err
	}
	return &toolResp.Data, nil
}

// DeleteTool deletes a tool by ID.
// DELETE /provision/perception/tools/:id.
func (s *Service) DeleteTool(id string) error {
	return genericDelete(s, id, "tool", "/provision/perception/tools/%s")
}
