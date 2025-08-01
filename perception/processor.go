package perception

import (
	"errors"
	"fmt"
)

// Processor represents a perception processor resource.
type Processor struct {
	ID             string         `json:"id,omitempty"`
	ThoughtID      string         `json:"thought_id,omitempty"`
	ModelID        string         `json:"model_id,omitempty"`
	Configuration  map[string]any `json:"configuration"`
	ProvisionState string         `json:"provision_state"`
	Type           string         `json:"type"`
}

// ProcessorResponse represents the API response for processor operations.
type ProcessorResponse struct {
	Data Processor `json:"data"`
}

// CreateProcessorRequest represents the request payload for creating a processor.
type CreateProcessorRequest struct {
	Processor ProcessorRequestData `json:"processor"`
}

// ProcessorRequestData represents the processor data in the request.
type ProcessorRequestData struct {
	ModelID       string         `json:"model_id"`
	Configuration map[string]any `json:"configuration"`
}

// UpdateProcessorRequest represents the request payload for updating a processor.
type UpdateProcessorRequest struct {
	Processor UpdateProcessorData `json:"processor"`
}

// UpdateProcessorData represents the processor update data.
type UpdateProcessorData struct {
	ModelID       string         `json:"model_id,omitempty"`
	Configuration map[string]any `json:"configuration,omitempty"`
}

// GetProcessor retrieves a specific processor by thought ID and type.
// GET /provision/perception/thoughts/:thought_id/types/:type/processor.
func (s *Service) GetProcessor(thoughtID, processorType string) (*Processor, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if processorType == "" {
		return nil, errors.New("processor type is required")
	}

	var processorResp ProcessorResponse
	resp, err := s.client.R().
		SetResult(&processorResp).
		Get(fmt.Sprintf("/provision/perception/thoughts/%s/types/%s/processor", thoughtID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to get processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// CreateProcessor creates a new processor.
// POST /provision/perception/thoughts/:thought_id/types/:type/processor.
func (s *Service) CreateProcessor(
	thoughtID, processorType string, req CreateProcessorRequest,
) (*Processor, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if processorType == "" {
		return nil, errors.New("processor type is required")
	}
	if req.Processor.ModelID == "" {
		return nil, errors.New("model ID is required")
	}

	var processorResp ProcessorResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&processorResp).
		Post(fmt.Sprintf("/provision/perception/thoughts/%s/types/%s/processor", thoughtID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to create processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// UpdateProcessor updates an existing processor using PATCH.
// PATCH /provision/perception/thoughts/:thought_id/types/:type/processor.
func (s *Service) UpdateProcessor(
	thoughtID, processorType string, req UpdateProcessorRequest,
) (*Processor, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if processorType == "" {
		return nil, errors.New("processor type is required")
	}

	var processorResp ProcessorResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&processorResp).
		Patch(fmt.Sprintf("/provision/perception/thoughts/%s/types/%s/processor", thoughtID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to update processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// ReplaceProcessor replaces an existing processor using PUT.
// PUT /provision/perception/thoughts/:thought_id/types/:type/processor.
func (s *Service) ReplaceProcessor(
	thoughtID, processorType string, req UpdateProcessorRequest,
) (*Processor, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if processorType == "" {
		return nil, errors.New("processor type is required")
	}

	var processorResp ProcessorResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&processorResp).
		Put(fmt.Sprintf("/provision/perception/thoughts/%s/types/%s/processor", thoughtID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to replace processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// DeleteProcessor deletes a processor by thought ID and type.
// DELETE /provision/perception/thoughts/:thought_id/types/:type/processor.
func (s *Service) DeleteProcessor(thoughtID, processorType string) error {
	if thoughtID == "" {
		return errors.New("thought ID is required")
	}
	if processorType == "" {
		return errors.New("processor type is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/thoughts/%s/types/%s/processor", thoughtID, processorType))

	if err != nil {
		return fmt.Errorf("failed to delete processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
