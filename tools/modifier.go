package tools

import (
	"errors"
	"fmt"
)

const (
	// ModifierMissingPolicyError fails the tool call when the configured value is unavailable.
	ModifierMissingPolicyError = "error"
	// ModifierMissingPolicySkip leaves the request unchanged when the configured value is unavailable.
	ModifierMissingPolicySkip = "skip"

	// ModifierSourceTypeMetadata resolves modifier values from trusted runtime metadata.
	ModifierSourceTypeMetadata = "metadata"

	// ModifierSourcePathActorIdentifier resolves the current actor identifier.
	ModifierSourcePathActorIdentifier = "actor_identifier"
	// ModifierSourcePathOriginEntityIdentifier resolves the originating entity identifier.
	ModifierSourcePathOriginEntityIdentifier = "origin_entity_identifier"
	// ModifierSourcePathCurrentTimestamp resolves the current timestamp.
	ModifierSourcePathCurrentTimestamp = "current_timestamp"
)

// ModifierSource identifies a trusted runtime value used by a tool modifier.
type ModifierSource struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

// Modifier represents a trusted thought-tool modifier resource.
type Modifier struct {
	ID              string         `json:"id,omitempty"`
	Index           int            `json:"index"`
	Target          string         `json:"target"`
	OnMissingParent string         `json:"on_missing_parent"`
	OnMissingSource string         `json:"on_missing_source"`
	Source          ModifierSource `json:"source"`
	ThoughtToolID   string         `json:"thought_tool_id,omitempty"`
	ProvisionState  string         `json:"provision_state"`
}

// ModifierResponse represents the API response for modifier operations.
type ModifierResponse struct {
	Data Modifier `json:"data"`
}

// CreateModifierRequest represents the request payload for creating a modifier.
type CreateModifierRequest struct {
	Modifier ModifierRequestData `json:"modifier"`
}

// ModifierRequestData represents the complete modifier data in a create request.
type ModifierRequestData struct {
	Index           int            `json:"index"`
	Target          string         `json:"target"`
	OnMissingParent string         `json:"on_missing_parent"`
	OnMissingSource string         `json:"on_missing_source"`
	Source          ModifierSource `json:"source"`
}

// UpdateModifierRequest represents the request payload for updating a modifier.
type UpdateModifierRequest struct {
	Modifier UpdateModifierData `json:"modifier"`
}

// UpdateModifierData represents the mutable fields in a modifier update.
type UpdateModifierData struct {
	Target          string          `json:"target,omitempty"`
	OnMissingParent string          `json:"on_missing_parent,omitempty"`
	OnMissingSource string          `json:"on_missing_source,omitempty"`
	Source          *ModifierSource `json:"source,omitempty"`
}

// GetModifier retrieves a specific active modifier by ID.
// GET /provision/tools/modifiers/:id.
func (s *Service) GetModifier(id string) (*Modifier, error) {
	if id == "" {
		return nil, errors.New("modifier ID is required")
	}

	var modifierResp ModifierResponse
	resp, err := s.client.R().
		SetResult(&modifierResp).
		Get(fmt.Sprintf("/provision/tools/modifiers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get modifier: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &modifierResp.Data, nil
}

// CreateModifier creates or reactivates a modifier within a thought tool.
// POST /provision/tools/:thought_tool_id/modifiers.
func (s *Service) CreateModifier(
	thoughtToolID string,
	req CreateModifierRequest,
) (*Modifier, error) {
	if thoughtToolID == "" {
		return nil, errors.New("thought tool ID is required")
	}
	if err := validateCreateModifier(req.Modifier); err != nil {
		return nil, err
	}

	var modifierResp ModifierResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&modifierResp).
		Post(fmt.Sprintf("/provision/tools/%s/modifiers", thoughtToolID))

	if err != nil {
		return nil, fmt.Errorf("failed to create modifier: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &modifierResp.Data, nil
}

// UpdateModifier partially updates an active or inactive modifier using PATCH.
// PATCH /provision/tools/modifiers/:id.
func (s *Service) UpdateModifier(id string, req UpdateModifierRequest) (*Modifier, error) {
	if id == "" {
		return nil, errors.New("modifier ID is required")
	}
	if err := validateUpdateModifier(req.Modifier); err != nil {
		return nil, err
	}

	var modifierResp ModifierResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&modifierResp).
		Patch(fmt.Sprintf("/provision/tools/modifiers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update modifier: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &modifierResp.Data, nil
}

// ReplaceModifier sends the mutable modifier fields using PUT.
// PUT /provision/tools/modifiers/:id.
func (s *Service) ReplaceModifier(id string, req UpdateModifierRequest) (*Modifier, error) {
	if id == "" {
		return nil, errors.New("modifier ID is required")
	}
	if err := validateUpdateModifier(req.Modifier); err != nil {
		return nil, err
	}

	var modifierResp ModifierResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&modifierResp).
		Put(fmt.Sprintf("/provision/tools/modifiers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace modifier: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &modifierResp.Data, nil
}

// DeleteModifier deactivates an active modifier by ID.
// DELETE /provision/tools/modifiers/:id.
func (s *Service) DeleteModifier(id string) error {
	if id == "" {
		return errors.New("modifier ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/tools/modifiers/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete modifier: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

func validateCreateModifier(modifier ModifierRequestData) error {
	if modifier.Index < 0 {
		return errors.New("modifier index must be non-negative")
	}
	if modifier.Target == "" {
		return errors.New("modifier target is required")
	}
	if modifier.OnMissingParent == "" {
		return errors.New("modifier on missing parent policy is required")
	}
	if modifier.OnMissingSource == "" {
		return errors.New("modifier on missing source policy is required")
	}

	return validateModifierSource(modifier.Source)
}

func validateModifierSource(source ModifierSource) error {
	if source.Type == "" {
		return errors.New("modifier source type is required")
	}
	if source.Path == "" {
		return errors.New("modifier source path is required")
	}

	return nil
}

func validateUpdateModifier(modifier UpdateModifierData) error {
	if modifier.Source == nil {
		return nil
	}

	return validateModifierSource(*modifier.Source)
}
