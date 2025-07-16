package sensory

import (
	"fmt"
)

// This file contains all Limit-related operations for the SensoryService.
// Limits represent rate limits and restrictions with scale units and counts.

// Limit operations

// GetLimit retrieves a specific limit by ID
// GET /provision/sensory/limits/:id
func (s *Service) GetLimit(id string) (*Limit, error) {
	if id == "" {
		return nil, fmt.Errorf("limit ID is required")
	}

	var limitResp LimitResponse
	resp, err := s.client.R().
		SetResult(&limitResp).
		SetError(&Error{}).
		Get(fmt.Sprintf("/provision/sensory/limits/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get limit: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &limitResp.Data, nil
}

// CreateLimit creates a new limit for a specific source
// POST /provision/sensory/sources/:source_id/limits
func (s *Service) CreateLimit(sourceID string, req CreateLimitRequest) (*Limit, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID is required")
	}
	if req.Limit.ScaleUnit == "" {
		return nil, fmt.Errorf("limit scale_unit is required")
	}
	if req.Limit.ScaleCount <= 0 {
		return nil, fmt.Errorf("limit scale_count must be greater than 0")
	}
	if req.Limit.Count <= 0 {
		return nil, fmt.Errorf("count value must be greater than 0")
	}

	var limitResp LimitResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&limitResp).
		SetError(&Error{}).
		Post(fmt.Sprintf("/provision/sensory/sources/%s/limits", sourceID))

	if err != nil {
		return nil, fmt.Errorf("failed to create limit: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &limitResp.Data, nil
}

// UpdateLimit updates an existing limit using PATCH
// PATCH /provision/sensory/limits/:id
func (s *Service) UpdateLimit(id string, req UpdateLimitRequest) (*Limit, error) {
	if id == "" {
		return nil, fmt.Errorf("limit ID is required")
	}

	var limitResp LimitResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&limitResp).
		SetError(&Error{}).
		Patch(fmt.Sprintf("/provision/sensory/limits/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update limit: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &limitResp.Data, nil
}

// ReplaceLimit replaces an existing limit using PUT
// PUT /provision/sensory/limits/:id
func (s *Service) ReplaceLimit(id string, req UpdateLimitRequest) (*Limit, error) {
	if id == "" {
		return nil, fmt.Errorf("limit ID is required")
	}

	var limitResp LimitResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&limitResp).
		SetError(&Error{}).
		Put(fmt.Sprintf("/provision/sensory/limits/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace limit: %w", err)
	}

	if resp.IsError() {
		if errorResp, ok := resp.Error().(*Error); ok {
			errorResp.StatusCode = resp.StatusCode()
			return nil, errorResp
		}
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &limitResp.Data, nil
}

// DeleteLimit deletes a limit by ID
// DELETE /provision/sensory/limits/:id
func (s *Service) DeleteLimit(id string) error {
	if id == "" {
		return fmt.Errorf("limit ID is required")
	}

	resp, err := s.client.R().
		SetError(&Error{}).
		Delete(fmt.Sprintf("/provision/sensory/limits/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete limit: %w", err)
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
