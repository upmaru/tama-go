package sensory

import (
	"fmt"
)

// This file contains all Model-related operations for the SensoryService.
// Models represent machine learning models with identifiers and paths.

// Model operations

// GetModel retrieves a specific model by ID
// GET /provision/sensory/models/:id
func (s *Service) GetModel(id string) (*Model, error) {
	if id == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	var modelResp ModelResponse
	resp, err := s.client.R().
		SetResult(&modelResp).
		SetError(&Error{}).
		Get(fmt.Sprintf("/provision/sensory/models/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &modelResp.Data, nil
}

// CreateModel creates a new model for a specific source
// POST /provision/sensory/sources/:source_id/models
func (s *Service) CreateModel(sourceID string, req CreateModelRequest) (*Model, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID is required")
	}
	if req.Model.Identifier == "" {
		return nil, fmt.Errorf("model identifier is required")
	}
	if req.Model.Path == "" {
		return nil, fmt.Errorf("model path is required")
	}

	var modelResp ModelResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&modelResp).
		SetError(&Error{}).
		Post(fmt.Sprintf("/provision/sensory/sources/%s/models", sourceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &modelResp.Data, nil
}

// UpdateModel updates an existing model using PATCH
// PATCH /provision/sensory/models/:id
func (s *Service) UpdateModel(id string, req UpdateModelRequest) (*Model, error) {
	if id == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	var modelResp ModelResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&modelResp).
		SetError(&Error{}).
		Patch(fmt.Sprintf("/provision/sensory/models/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update model: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &modelResp.Data, nil
}

// ReplaceModel replaces an existing model using PUT
// PUT /provision/sensory/models/:id
func (s *Service) ReplaceModel(id string, req UpdateModelRequest) (*Model, error) {
	if id == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	var modelResp ModelResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&modelResp).
		SetError(&Error{}).
		Put(fmt.Sprintf("/provision/sensory/models/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace model: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &modelResp.Data, nil
}

// DeleteModel deletes a model by ID
// DELETE /provision/sensory/models/:id
func (s *Service) DeleteModel(id string) error {
	if id == "" {
		return fmt.Errorf("model ID is required")
	}

	resp, err := s.client.R().
		SetError(&Error{}).
		Delete(fmt.Sprintf("/provision/sensory/models/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return errorResp
		}
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}
