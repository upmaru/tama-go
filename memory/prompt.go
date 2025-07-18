package memory

import (
	"errors"
	"fmt"
)

// This file contains all Prompt-related operations for the MemoryService.
// Prompts represent memory prompts with content, roles, and associated spaces.

// Prompt operations

// GetPrompt retrieves a specific prompt by ID.
// GET /provision/memory/prompts/:id.
func (s *Service) GetPrompt(id string) (*Prompt, error) {
	if id == "" {
		return nil, errors.New("prompt ID is required")
	}

	var promptResp PromptResponse
	resp, err := s.client.R().
		SetResult(&promptResp).
		Get(fmt.Sprintf("/provision/memory/prompts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &promptResp.Data, nil
}

// CreatePrompt creates a new prompt in a specific space.
// POST /provision/memory/spaces/:space_id/prompts.
func (s *Service) CreatePrompt(spaceID string, req CreatePromptRequest) (*Prompt, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if req.Prompt.Name == "" {
		return nil, errors.New("prompt name is required")
	}
	if req.Prompt.Content == "" {
		return nil, errors.New("prompt content is required")
	}
	if req.Prompt.Role == "" {
		return nil, errors.New("prompt role is required")
	}

	var promptResp PromptResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&promptResp).
		Post(fmt.Sprintf("/provision/memory/spaces/%s/prompts", spaceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create prompt: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &promptResp.Data, nil
}

// UpdatePrompt updates an existing prompt using PATCH.
// PATCH /provision/memory/prompts/:id.
func (s *Service) UpdatePrompt(id string, req UpdatePromptRequest) (*Prompt, error) {
	if id == "" {
		return nil, errors.New("prompt ID is required")
	}

	var promptResp PromptResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&promptResp).
		Patch(fmt.Sprintf("/provision/memory/prompts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update prompt: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &promptResp.Data, nil
}

// ReplacePrompt replaces an existing prompt using PUT.
// PUT /provision/memory/prompts/:id.
func (s *Service) ReplacePrompt(id string, req UpdatePromptRequest) (*Prompt, error) {
	if id == "" {
		return nil, errors.New("prompt ID is required")
	}

	var promptResp PromptResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&promptResp).
		Put(fmt.Sprintf("/provision/memory/prompts/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace prompt: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &promptResp.Data, nil
}

// DeletePrompt deletes a prompt by ID.
// DELETE /provision/memory/prompts/:id.
func (s *Service) DeletePrompt(id string) error {
	if id == "" {
		return errors.New("prompt ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/memory/prompts/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete prompt: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
