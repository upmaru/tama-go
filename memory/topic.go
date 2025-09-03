package memory

import (
	"errors"
	"fmt"
)

// Topic represents a memory topic resource.
// All properties are strings per API contract.
type Topic struct {
	ID             string `json:"id,omitempty"`
	ListenerID     string `json:"listener_id"`
	ClassID        string `json:"class_id"`
	ProvisionState string `json:"provision_state"`
}

// TopicResponse represents the API response for topic operations.
type TopicResponse struct {
	Data Topic `json:"data"`
}

// CreateTopicRequest represents the request payload for creating a topic.
type CreateTopicRequest struct {
	Topic TopicRequestData `json:"topic"`
}

// TopicRequestData represents the topic data in the create request.
type TopicRequestData struct {
	ClassID string `json:"class_id"`
}

// UpdateTopicRequest represents the request payload for updating or replacing a topic.
type UpdateTopicRequest struct {
	Topic UpdateTopicData `json:"topic"`
}

// UpdateTopicData represents the topic update data.
type UpdateTopicData struct {
	ClassID string `json:"class_id,omitempty"`
}

// GetTopic retrieves a specific topic by ID.
// GET /provision/memory/topics/:id.
func (s *Service) GetTopic(id string) (*Topic, error) {
	if id == "" {
		return nil, errors.New("topic ID is required")
	}

	var topicResp TopicResponse
	resp, err := s.client.R().
		SetResult(&topicResp).
		Get(fmt.Sprintf("/provision/memory/topics/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &topicResp.Data, nil
}

// CreateTopic creates a new topic within a listener.
// POST /provision/memory/listeners/:listener_id/topics.
func (s *Service) CreateTopic(listenerID string, req CreateTopicRequest) (*Topic, error) {
	if listenerID == "" {
		return nil, errors.New("listener ID is required")
	}
	if req.Topic.ClassID == "" {
		return nil, errors.New("class ID is required")
	}

	var topicResp TopicResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&topicResp).
		Post(fmt.Sprintf("/provision/memory/listeners/%s/topics", listenerID))

	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &topicResp.Data, nil
}

// UpdateTopic updates an existing topic using PATCH.
// PATCH /provision/memory/topics/:id.
func (s *Service) UpdateTopic(id string, req UpdateTopicRequest) (*Topic, error) {
	if id == "" {
		return nil, errors.New("topic ID is required")
	}
	if req.Topic.ClassID == "" {
		return nil, errors.New("class ID is required")
	}

	var topicResp TopicResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&topicResp).
		Patch(fmt.Sprintf("/provision/memory/topics/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &topicResp.Data, nil
}

// ReplaceTopic replaces an existing topic using PUT.
// PUT /provision/memory/topics/:id.
func (s *Service) ReplaceTopic(id string, req UpdateTopicRequest) (*Topic, error) {
	if id == "" {
		return nil, errors.New("topic ID is required")
	}
	if req.Topic.ClassID == "" {
		return nil, errors.New("class ID is required")
	}

	var topicResp TopicResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&topicResp).
		Put(fmt.Sprintf("/provision/memory/topics/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace topic: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &topicResp.Data, nil
}

// DeleteTopic deletes a topic by ID.
// DELETE /provision/memory/topics/:id.
func (s *Service) DeleteTopic(id string) error {
	if id == "" {
		return errors.New("topic ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/memory/topics/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

