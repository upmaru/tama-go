package motor

import (
	"errors"
	"fmt"
)

// Modifier represents a motor modifier resource.
//
// Fields:
// - ID: server-generated identifier
// - Name: human-readable name of the modifier
// - ActionID: the parent action this modifier belongs to
// - Schema: arbitrary configuration for the modifier
// - ProvisionState: server-side state tracking
//
// Example JSON shape:
//
//	{
//	  "id": "modifier-123",
//	  "name": "sanitize",
//	  "action_id": "action-456",
//	  "schema": {"key": "value"},
//	  "provision_state": "active"
//	}
type Modifier struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name"`
	ActionID       string         `json:"action_id,omitempty"`
	Schema         map[string]any `json:"schema"`
	ProvisionState string         `json:"provision_state,omitempty"`
}

// ModifierResponse wraps the modifier data in API responses.
type ModifierResponse struct {
	Data Modifier `json:"data"`
}

// CreateModifierRequest represents the payload for creating a modifier under an action.
type CreateModifierRequest struct {
	Modifier ModifierRequestData `json:"modifier"`
}

// ModifierRequestData holds the fields required to create a modifier.
type ModifierRequestData struct {
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"`
}

// UpdateModifierRequest represents the payload for updating/replacing a modifier.
type UpdateModifierRequest struct {
	Modifier UpdateModifierData `json:"modifier"`
}

// UpdateModifierData holds the fields that can be updated on a modifier.
type UpdateModifierData struct {
	Name   string         `json:"name,omitempty"`
	Schema map[string]any `json:"schema,omitempty"`
}

// GetModifier retrieves a specific modifier by ID.
// GET /provision/motor/modifiers/:id.
func (s *Service) GetModifier(id string) (*Modifier, error) {
	if id == "" {
		return nil, errors.New("modifier ID is required")
	}

	var respWrapper ModifierResponse
	resp, err := s.client.R().
		SetResult(&respWrapper).
		Get(fmt.Sprintf("/provision/motor/modifiers/%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to get modifier: %w", err)
	}
	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}
	return &respWrapper.Data, nil
}

// CreateModifier creates a new modifier under the specified action.
// POST /provision/motor/actions/:action_id/modifiers.
func (s *Service) CreateModifier(actionID string, req CreateModifierRequest) (*Modifier, error) {
	if actionID == "" {
		return nil, errors.New("action ID is required")
	}
	if req.Modifier.Name == "" {
		return nil, errors.New("modifier name is required")
	}
	if req.Modifier.Schema == nil {
		return nil, errors.New("modifier schema is required")
	}

	var respWrapper ModifierResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&respWrapper).
		Post(fmt.Sprintf("/provision/motor/actions/%s/modifiers", actionID))
	if err != nil {
		return nil, fmt.Errorf("failed to create modifier: %w", err)
	}
	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}
	return &respWrapper.Data, nil
}

// UpdateModifier updates an existing modifier using PATCH.
// PATCH /provision/motor/modifiers/:id.
func (s *Service) UpdateModifier(id string, req UpdateModifierRequest) (*Modifier, error) {
	if id == "" {
		return nil, errors.New("modifier ID is required")
	}

	var respWrapper ModifierResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&respWrapper).
		Patch(fmt.Sprintf("/provision/motor/modifiers/%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to update modifier: %w", err)
	}
	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}
	return &respWrapper.Data, nil
}

// ReplaceModifier replaces an existing modifier using PUT.
// PUT /provision/motor/modifiers/:id.
func (s *Service) ReplaceModifier(id string, req UpdateModifierRequest) (*Modifier, error) {
	if id == "" {
		return nil, errors.New("modifier ID is required")
	}

	var respWrapper ModifierResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&respWrapper).
		Put(fmt.Sprintf("/provision/motor/modifiers/%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to replace modifier: %w", err)
	}
	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}
	return &respWrapper.Data, nil
}

// DeleteModifier deletes a modifier by ID.
// DELETE /provision/motor/modifiers/:id.
func (s *Service) DeleteModifier(id string) error {
	if id == "" {
		return errors.New("modifier ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/motor/modifiers/%s", id))
	if err != nil {
		return fmt.Errorf("failed to delete modifier: %w", err)
	}
	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}
	return nil
}
