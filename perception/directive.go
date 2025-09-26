package perception

import "errors"

// Directive represents a perception directive resource.
type Directive struct {
	ID              string `json:"id,omitempty"`
	ThoughtPathID   string `json:"thought_path_id,omitempty"`
	PromptID        string `json:"prompt_id,omitempty"`
	TargetThoughtID string `json:"target_thought_id,omitempty"`
	ProvisionState  string `json:"provision_state"`
}

// DirectiveResponse represents the API response for directive operations.
type DirectiveResponse struct {
	Data Directive `json:"data"`
}

// CreateDirectiveRequest represents the request payload for creating a directive.
type CreateDirectiveRequest struct {
	Directive DirectiveRequestData `json:"directive"`
}

// DirectiveRequestData represents the directive data in the request.
type DirectiveRequestData struct {
	PromptID        string `json:"prompt_id"`
	TargetThoughtID string `json:"target_thought_id"`
}

// UpdateDirectiveRequest represents the request payload for updating a directive.
type UpdateDirectiveRequest struct {
	Directive UpdateDirectiveData `json:"directive"`
}

// UpdateDirectiveData represents the directive update data.
type UpdateDirectiveData struct {
	PromptID        string `json:"prompt_id,omitempty"`
	TargetThoughtID string `json:"target_thought_id,omitempty"`
}

// GetDirective retrieves a specific directive by ID.
// GET /provision/perception/directives/:id.
func (s *Service) GetDirective(id string) (*Directive, error) {
	var directiveResp DirectiveResponse
	if err := genericGet(s, id, "directive", "/provision/perception/directives/%s", &directiveResp); err != nil {
		return nil, err
	}
	return &directiveResp.Data, nil
}

// CreateDirective creates a new directive within a path.
// POST /provision/perception/paths/:path_id/directives.
func (s *Service) CreateDirective(pathID string, req CreateDirectiveRequest) (*Directive, error) {
	if req.Directive.PromptID == "" {
		return nil, errors.New("prompt ID is required")
	}
	if req.Directive.TargetThoughtID == "" {
		return nil, errors.New("target thought ID is required")
	}

	var directiveResp DirectiveResponse
	if err := genericCreate(s, pathID, req, "directive", "path", "/provision/perception/paths/%s/directives", &directiveResp); err != nil {
		return nil, err
	}
	return &directiveResp.Data, nil
}

// UpdateDirective updates an existing directive using PATCH.
// PATCH /provision/perception/directives/:id.
func (s *Service) UpdateDirective(id string, req UpdateDirectiveRequest) (*Directive, error) {
	var directiveResp DirectiveResponse
	if err := genericUpdate(s, id, req, "directive", "/provision/perception/directives/%s", &directiveResp); err != nil {
		return nil, err
	}
	return &directiveResp.Data, nil
}

// ReplaceDirective replaces an existing directive using PUT.
// PUT /provision/perception/directives/:id.
func (s *Service) ReplaceDirective(id string, req UpdateDirectiveRequest) (*Directive, error) {
	var directiveResp DirectiveResponse
	if err := genericReplace(s, id, req, "directive", "/provision/perception/directives/%s", &directiveResp); err != nil {
		return nil, err
	}
	return &directiveResp.Data, nil
}

// DeleteDirective deletes a directive by ID.
// DELETE /provision/perception/directives/:id.
func (s *Service) DeleteDirective(id string) error {
	return genericDelete(s, id, "directive", "/provision/perception/directives/%s")
}
