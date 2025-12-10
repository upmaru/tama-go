package system

import (
	"errors"
	"fmt"
)

// Queue represents a system queue resource.
type Queue struct {
	ID          string `json:"id,omitempty"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	Concurrency int    `json:"concurrency"`
}

// QueueResponse represents the API response for queue operations.
type QueueResponse struct {
	Data Queue `json:"data"`
}

// CreateQueueRequest represents the request payload for creating a queue.
type CreateQueueRequest struct {
	Queue QueueRequestData `json:"queue"`
}

// QueueRequestData represents the queue data in the create request.
type QueueRequestData struct {
	Role        string `json:"role"`
	Name        string `json:"name"`
	Concurrency int    `json:"concurrency"`
}

// UpdateQueueRequest represents the request payload for updating a queue.
type UpdateQueueRequest struct {
	Queue UpdateQueueData `json:"queue"`
}

// UpdateQueueData represents the queue update data.
type UpdateQueueData struct {
	Role        string `json:"role,omitempty"`
	Name        string `json:"name,omitempty"`
	Concurrency *int   `json:"concurrency,omitempty"`
}

// GetQueue retrieves a specific queue by ID.
// GET /provision/system/queues/:id.
func (s *Service) GetQueue(id string) (*Queue, error) {
	if id == "" {
		return nil, errors.New("queue ID is required")
	}

	var queueResp QueueResponse
	resp, err := s.client.R().
		SetResult(&queueResp).
		Get(fmt.Sprintf("/provision/system/queues/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &queueResp.Data, nil
}

// CreateQueue creates a new queue.
// POST /provision/system/queues.
func (s *Service) CreateQueue(req CreateQueueRequest) (*Queue, error) {
	if req.Queue.Role == "" {
		return nil, errors.New("queue role is required")
	}
	if req.Queue.Name == "" {
		return nil, errors.New("queue name is required")
	}
	if req.Queue.Concurrency <= 0 {
		return nil, errors.New("queue concurrency must be greater than 0")
	}

	var queueResp QueueResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&queueResp).
		Post("/provision/system/queues")

	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &queueResp.Data, nil
}

// UpdateQueue updates an existing queue using PATCH.
// PATCH /provision/system/queues/:id.
func (s *Service) UpdateQueue(id string, req UpdateQueueRequest) (*Queue, error) {
	if id == "" {
		return nil, errors.New("queue ID is required")
	}
	if req.Queue.Concurrency != nil && *req.Queue.Concurrency <= 0 {
		return nil, errors.New("queue concurrency must be greater than 0")
	}

	var queueResp QueueResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&queueResp).
		Patch(fmt.Sprintf("/provision/system/queues/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update queue: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &queueResp.Data, nil
}

// DeleteQueue deletes a queue by ID.
// DELETE /provision/system/queues/:id.
func (s *Service) DeleteQueue(id string) error {
	if id == "" {
		return errors.New("queue ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/system/queues/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete queue: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
