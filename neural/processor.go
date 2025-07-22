package neural

import (
	"errors"
	"fmt"
)

// GetProcessor retrieves a specific processor by space ID and type.
// GET /provision/neural/spaces/:space_id/types/:type/processor.
func (s *Service) GetProcessor(spaceID, processorType string) (*Processor, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if processorType == "" {
		return nil, errors.New("processor type is required")
	}

	var processorResp ProcessorResponse
	resp, err := s.client.R().
		SetResult(&processorResp).
		Get(fmt.Sprintf("/provision/neural/spaces/%s/types/%s/processor", spaceID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to get processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// CreateProcessor creates a new processor.
// POST /provision/neural/spaces/:space_id/types/:type/processor.
func (s *Service) CreateProcessor(
	spaceID, processorType string, req CreateProcessorRequest,
) (*Processor, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
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
		Post(fmt.Sprintf("/provision/neural/spaces/%s/types/%s/processor", spaceID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to create processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// UpdateProcessor updates an existing processor using PATCH.
// PATCH /provision/neural/spaces/:space_id/types/:type/processor.
func (s *Service) UpdateProcessor(
	spaceID, processorType string, req UpdateProcessorRequest,
) (*Processor, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if processorType == "" {
		return nil, errors.New("processor type is required")
	}

	var processorResp ProcessorResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&processorResp).
		Patch(fmt.Sprintf("/provision/neural/spaces/%s/types/%s/processor", spaceID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to update processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// ReplaceProcessor replaces an existing processor using PUT.
// PUT /provision/neural/spaces/:space_id/types/:type/processor.
func (s *Service) ReplaceProcessor(
	spaceID, processorType string, req UpdateProcessorRequest,
) (*Processor, error) {
	if spaceID == "" {
		return nil, errors.New("space ID is required")
	}
	if processorType == "" {
		return nil, errors.New("processor type is required")
	}

	var processorResp ProcessorResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&processorResp).
		Put(fmt.Sprintf("/provision/neural/spaces/%s/types/%s/processor", spaceID, processorType))

	if err != nil {
		return nil, fmt.Errorf("failed to replace processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &processorResp.Data, nil
}

// DeleteProcessor deletes a processor by space ID and type.
// DELETE /provision/neural/spaces/:space_id/types/:type/processor.
func (s *Service) DeleteProcessor(spaceID, processorType string) error {
	if spaceID == "" {
		return errors.New("space ID is required")
	}
	if processorType == "" {
		return errors.New("processor type is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/neural/spaces/%s/types/%s/processor", spaceID, processorType))

	if err != nil {
		return fmt.Errorf("failed to delete processor: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
