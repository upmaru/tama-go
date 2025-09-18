//nolint:dupl // Similar CRUD shape as listener; intentional duplication for consistent service pattern
package tools

import (
	"errors"
	"fmt"
)

// Option represents a tools option resource.
//
// Fields:
// - ID: server-generated identifier
// - ThoughtToolOutputID: parent thought tool output ID
// - ActionModifierID: references a motor action modifier
// - ProvisionState: server-side provision state
//
// Example JSON shape:
//
//	{
//	  "id": "option-123",
//	  "thought_tool_output_id": "output-456",
//	  "action_modifier_id": "modifier-789",
//	  "provision_state": "active"
//	}
type Option struct {
	ID                  string `json:"id,omitempty"`
	ThoughtToolOutputID string `json:"thought_tool_output_id,omitempty"`
	ActionModifierID    string `json:"action_modifier_id"`
	ProvisionState      string `json:"provision_state"`
}

// OptionResponse represents the API response for option operations.
type OptionResponse struct {
	Data Option `json:"data"`
}

// CreateOptionRequest represents the request payload for creating an option.
type CreateOptionRequest struct {
	Option OptionRequestData `json:"option"`
}

// OptionRequestData represents the option data in the request.
type OptionRequestData struct {
	ActionModifierID string `json:"action_modifier_id"`
}

// UpdateOptionRequest represents the request payload for updating an option.
type UpdateOptionRequest struct {
	Option UpdateOptionData `json:"option"`
}

// UpdateOptionData represents the option update data.
type UpdateOptionData struct {
	ActionModifierID string `json:"action_modifier_id,omitempty"`
}

// GetOption retrieves a specific option by ID.
// GET /provision/tools/options/:id.
func (s *Service) GetOption(id string) (*Option, error) {
	if id == "" {
		return nil, errors.New("option ID is required")
	}

	var optionResp OptionResponse
	resp, err := s.client.R().
		SetResult(&optionResp).
		Get(fmt.Sprintf("/provision/tools/options/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get option: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &optionResp.Data, nil
}

// CreateOption creates a new option within a thought tool output.
// POST /provision/tools/outputs/:output_id/options.
func (s *Service) CreateOption(outputID string, req CreateOptionRequest) (*Option, error) {
	if outputID == "" {
		return nil, errors.New("output ID is required")
	}
	if req.Option.ActionModifierID == "" {
		return nil, errors.New("option action modifier ID is required")
	}

	var optionResp OptionResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&optionResp).
		Post(fmt.Sprintf("/provision/tools/outputs/%s/options", outputID))

	if err != nil {
		return nil, fmt.Errorf("failed to create option: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &optionResp.Data, nil
}

// UpdateOption updates an existing option using PATCH.
// PATCH /provision/tools/options/:id.
func (s *Service) UpdateOption(id string, req UpdateOptionRequest) (*Option, error) {
	if id == "" {
		return nil, errors.New("option ID is required")
	}

	var optionResp OptionResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&optionResp).
		Patch(fmt.Sprintf("/provision/tools/options/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update option: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &optionResp.Data, nil
}

// ReplaceOption replaces an existing option using PUT.
// PUT /provision/tools/options/:id.
func (s *Service) ReplaceOption(id string, req UpdateOptionRequest) (*Option, error) {
	if id == "" {
		return nil, errors.New("option ID is required")
	}

	var optionResp OptionResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&optionResp).
		Put(fmt.Sprintf("/provision/tools/options/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace option: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &optionResp.Data, nil
}

// DeleteOption deletes an option by ID.
// DELETE /provision/tools/options/:id.
func (s *Service) DeleteOption(id string) error {
	if id == "" {
		return errors.New("option ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/tools/options/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete option: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
