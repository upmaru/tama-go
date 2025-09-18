//nolint:dupl // Similar CRUD shape as listener; intentional duplication for consistent service pattern
package tools

import (
	"errors"
	"fmt"
)

// Output represents a tools output resource.
//
// Fields:
// - ID: server-generated identifier
// - ThoughtToolID: parent thought tool ID
// - ClassCorpusID: target class corpus ID
// - ProvisionState: server-side provision state
//
// Example JSON shape:
//
//	{
//	  "id": "output-123",
//	  "thought_tool_id": "tool-456",
//	  "class_corpus_id": "corpus-789",
//	  "provision_state": "active"
//	}
type Output struct {
	ID             string `json:"id,omitempty"`
	ThoughtToolID  string `json:"thought_tool_id,omitempty"`
	ClassCorpusID  string `json:"class_corpus_id"`
	ProvisionState string `json:"provision_state"`
}

// OutputResponse represents the API response for output operations.
type OutputResponse struct {
	Data Output `json:"data"`
}

// CreateOutputRequest represents the request payload for creating an output.
type CreateOutputRequest struct {
	Output OutputRequestData `json:"output"`
}

// OutputRequestData represents the output data in the request.
type OutputRequestData struct {
	ClassCorpusID string `json:"class_corpus_id"`
}

// UpdateOutputRequest represents the request payload for updating an output.
type UpdateOutputRequest struct {
	Output UpdateOutputData `json:"output"`
}

// UpdateOutputData represents the output update data.
type UpdateOutputData struct {
	ClassCorpusID string `json:"class_corpus_id,omitempty"`
}

// GetOutput retrieves a specific output by ID.
// GET /provision/tools/outputs/:id.
func (s *Service) GetOutput(id string) (*Output, error) {
	if id == "" {
		return nil, errors.New("output ID is required")
	}

	var outputResp OutputResponse
	resp, err := s.client.R().
		SetResult(&outputResp).
		Get(fmt.Sprintf("/provision/tools/outputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get output: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &outputResp.Data, nil
}

// CreateOutput creates a new output within a thought tool.
// POST /provision/tools/:thought_tool_id/outputs.
func (s *Service) CreateOutput(thoughtToolID string, req CreateOutputRequest) (*Output, error) {
	if thoughtToolID == "" {
		return nil, errors.New("thought tool ID is required")
	}
	if req.Output.ClassCorpusID == "" {
		return nil, errors.New("output class corpus ID is required")
	}

	var outputResp OutputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&outputResp).
		Post(fmt.Sprintf("/provision/tools/%s/outputs", thoughtToolID))

	if err != nil {
		return nil, fmt.Errorf("failed to create output: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &outputResp.Data, nil
}

// UpdateOutput updates an existing output using PATCH.
// PATCH /provision/tools/outputs/:id.
func (s *Service) UpdateOutput(id string, req UpdateOutputRequest) (*Output, error) {
	if id == "" {
		return nil, errors.New("output ID is required")
	}

	var outputResp OutputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&outputResp).
		Patch(fmt.Sprintf("/provision/tools/outputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update output: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &outputResp.Data, nil
}

// ReplaceOutput replaces an existing output using PUT.
// PUT /provision/tools/outputs/:id.
func (s *Service) ReplaceOutput(id string, req UpdateOutputRequest) (*Output, error) {
	if id == "" {
		return nil, errors.New("output ID is required")
	}

	var outputResp OutputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&outputResp).
		Put(fmt.Sprintf("/provision/tools/outputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace output: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &outputResp.Data, nil
}

// DeleteOutput deletes an output by ID.
// DELETE /provision/tools/outputs/:id.
func (s *Service) DeleteOutput(id string) error {
	if id == "" {
		return errors.New("output ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/tools/outputs/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete output: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
