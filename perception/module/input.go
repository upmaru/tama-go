package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Input represents a module input resource.
type Input struct {
	ID              string `json:"id,omitempty"`
	Type            string `json:"type"`
	ThoughtModuleID string `json:"thought_module_id,omitempty"`
	ThoughtID       string `json:"thought_id,omitempty"`
	ClassCorpusID   string `json:"class_corpus_id"`
	ProvisionState  string `json:"provision_state,omitempty"`
}

// InputResponse represents the API response for input operations.
type InputResponse struct {
	Data Input `json:"data"`
}

// CreateInputRequest represents the request payload for creating a module input.
type CreateInputRequest struct {
	Input CreateInputData `json:"input"`
}

// CreateInputData represents the input data in the create request.
type CreateInputData struct {
	Type          string `json:"type"`
	ClassCorpusID string `json:"class_corpus_id"`
}

// UpdateInputRequest represents the request payload for updating a module input.
type UpdateInputRequest struct {
	Input UpdateInputData `json:"input"`
}

// UpdateInputData represents the input update data.
type UpdateInputData struct {
	Type          string `json:"type,omitempty"`
	ClassCorpusID string `json:"class_corpus_id,omitempty"`
}

// Service handles all module input related API operations.
type Service struct {
	client *resty.Client
}

// NewService creates a new module input service instance.
func NewService(client *resty.Client) *Service {
	return &Service{
		client: client,
	}
}

// Error represents an API error response.
type Error struct {
	StatusCode int                 `json:"status_code"`
	Errors     map[string][]string `json:"errors,omitempty"`
}

func (e *Error) Error() string {
	if len(e.Errors) > 0 {
		var errorParts []string
		for field, messages := range e.Errors {
			for _, message := range messages {
				errorParts = append(errorParts, fmt.Sprintf("%s %s", field, message))
			}
		}
		if e.StatusCode > 0 {
			return fmt.Sprintf("API error %d: %s", e.StatusCode, strings.Join(errorParts, ", "))
		}
		return fmt.Sprintf("API error: %s", strings.Join(errorParts, ", "))
	}

	if e.StatusCode > 0 {
		return fmt.Sprintf("API error %d", e.StatusCode)
	}
	return "API error"
}

// GetInput retrieves a specific input by ID.
// GET /provision/perception/module/inputs/:id.
func (s *Service) GetInput(id string) (*Input, error) {
	if id == "" {
		return nil, errors.New("input ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetResult(&inputResp).
		Get(fmt.Sprintf("/provision/perception/module/inputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// CreateInput creates a new module input within a thought.
// POST /provision/perception/thoughts/:thought_id/module/inputs.
func (s *Service) CreateInput(thoughtID string, req CreateInputRequest) (*Input, error) {
	if thoughtID == "" {
		return nil, errors.New("thought ID is required")
	}
	if req.Input.Type == "" {
		return nil, errors.New("input type is required")
	}
	if req.Input.ClassCorpusID == "" {
		return nil, errors.New("class corpus ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&inputResp).
		Post(fmt.Sprintf("/provision/perception/thoughts/%s/module/inputs", thoughtID))

	if err != nil {
		return nil, fmt.Errorf("failed to create input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// UpdateInput updates an existing input using PATCH.
// PATCH /provision/perception/module/inputs/:id.
func (s *Service) UpdateInput(id string, req UpdateInputRequest) (*Input, error) {
	if id == "" {
		return nil, errors.New("input ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&inputResp).
		Patch(fmt.Sprintf("/provision/perception/module/inputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// ReplaceInput replaces an existing input using PUT.
// PUT /provision/perception/module/inputs/:id.
func (s *Service) ReplaceInput(id string, req UpdateInputRequest) (*Input, error) {
	if id == "" {
		return nil, errors.New("input ID is required")
	}

	var inputResp InputResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&inputResp).
		Put(fmt.Sprintf("/provision/perception/module/inputs/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &inputResp.Data, nil
}

// DeleteInput deletes an input by ID.
// DELETE /provision/perception/module/inputs/:id.
func (s *Service) DeleteInput(id string) error {
	if id == "" {
		return errors.New("input ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/perception/module/inputs/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete input: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

// handleAPIError processes API error responses.
func (s *Service) handleAPIError(resp any) error {
	type errorResponse interface {
		IsError() bool
		Error() any
		StatusCode() int
		Status() string
		Body() []byte
	}

	if errResp, ok := resp.(errorResponse); ok && errResp.IsError() {
		if body := errResp.Body(); len(body) > 0 {
			var errorResp struct {
				Errors map[string][]string `json:"errors"`
			}

			if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Errors != nil {
				return &Error{
					StatusCode: errResp.StatusCode(),
					Errors:     errorResp.Errors,
				}
			}
		}

		return fmt.Errorf("API error: %s", errResp.Status())
	}

	return nil
}
