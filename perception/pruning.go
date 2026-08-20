package perception

// Pruning represents a perception pruning resource.
type Pruning struct {
	ID                    string `json:"id"`
	ThoughtID             string `json:"thought_id"`
	PreviousVersionsCount int    `json:"previous_versions_count"`
}

// PruningResponse represents the API response for pruning operations.
type PruningResponse struct {
	Data Pruning `json:"data"`
}

// CreatePruningRequest represents the request payload for creating a pruning.
type CreatePruningRequest struct {
	Pruning CreatePruningData `json:"pruning"`
}

// CreatePruningData represents the pruning data in the request.
type CreatePruningData struct {
	PreviousVersionsCount int `json:"previous_versions_count"`
}

// UpdatePruningRequest represents the request payload for updating a pruning.
type UpdatePruningRequest struct {
	Pruning UpdatePruningData `json:"pruning"`
}

// UpdatePruningData represents the pruning update data.
type UpdatePruningData struct {
	PreviousVersionsCount *int `json:"previous_versions_count,omitempty"`
}

// GetPruning retrieves a specific pruning by ID.
// GET /provision/perception/prunings/:id.
func (s *Service) GetPruning(id string) (*Pruning, error) {
	var pruningResp PruningResponse
	if err := genericGet(s, id, "pruning", "/provision/perception/prunings/%s", &pruningResp); err != nil {
		return nil, err
	}
	return &pruningResp.Data, nil
}

// CreatePruning creates a new pruning within a thought.
// POST /provision/perception/thoughts/:thought_id/prunings.
func (s *Service) CreatePruning(thoughtID string, req CreatePruningRequest) (*Pruning, error) {
	var pruningResp PruningResponse
	if err := genericCreate(
		s,
		thoughtID,
		req,
		"pruning",
		"thought",
		"/provision/perception/thoughts/%s/prunings",
		&pruningResp,
	); err != nil {
		return nil, err
	}
	return &pruningResp.Data, nil
}

// UpdatePruning updates an existing pruning using PATCH.
// PATCH /provision/perception/prunings/:id.
func (s *Service) UpdatePruning(id string, req UpdatePruningRequest) (*Pruning, error) {
	var pruningResp PruningResponse
	if err := genericUpdate(s, id, req, "pruning", "/provision/perception/prunings/%s", &pruningResp); err != nil {
		return nil, err
	}
	return &pruningResp.Data, nil
}

// ReplacePruning replaces an existing pruning using PUT.
// PUT /provision/perception/prunings/:id.
func (s *Service) ReplacePruning(id string, req UpdatePruningRequest) (*Pruning, error) {
	var pruningResp PruningResponse
	if err := genericReplace(s, id, req, "pruning", "/provision/perception/prunings/%s", &pruningResp); err != nil {
		return nil, err
	}
	return &pruningResp.Data, nil
}

// DeletePruning deletes a pruning by ID.
// DELETE /provision/perception/prunings/:id.
func (s *Service) DeletePruning(id string) error {
	return genericDelete(s, id, "pruning", "/provision/perception/prunings/%s")
}
