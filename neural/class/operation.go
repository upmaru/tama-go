package class

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Operation represents a neural class operation resource.
type Operation struct {
	ID           string   `json:"id,omitempty"`
	CurrentState string   `json:"current_state"`
	ClassID      string   `json:"class_id"`
	NodeIDs      []string `json:"node_ids"`
}

// OperationResponse represents the API response for operation operations.
type OperationResponse struct {
	Data Operation `json:"data"`
}

// CreateOperationRequest represents the request payload for creating a neural class operation.
type CreateOperationRequest struct {
	Operation CreateOperationData `json:"operation"`
}

// CreateOperationData represents the operation data in the create request.
type CreateOperationData struct {
	ChainIDs []string `json:"chain_ids"`
	NodeType *string  `json:"node_type,omitempty"`
}

// Service handles all neural class operation related API operations.
type Service struct {
	client *resty.Client
}

// NewService creates a new neural class operation service instance.
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

// GetOperation retrieves a specific operation by ID.
// GET /provision/neural/classes/:class_id/operations/:id.
func (s *Service) GetOperation(classID string, operationID string) (*Operation, error) {
	if classID == "" {
		return nil, errors.New("class ID is required")
	}
	if operationID == "" {
		return nil, errors.New("operation ID is required")
	}

	var operationResp OperationResponse
	resp, err := s.client.R().
		SetResult(&operationResp).
		Get(fmt.Sprintf("/provision/neural/classes/%s/operations/%s", classID, operationID))

	if err != nil {
		return nil, fmt.Errorf("failed to get operation: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &operationResp.Data, nil
}

// CreateOperation creates a new operation for a neural class.
// POST /provision/neural/classes/:class_id/operations.
func (s *Service) CreateOperation(classID string, req CreateOperationRequest) (*Operation, error) {
	if classID == "" {
		return nil, errors.New("class ID is required")
	}

	if len(req.Operation.ChainIDs) == 0 {
		return nil, errors.New("chain IDs are required")
	}

	var operationResp OperationResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&operationResp).
		Post(fmt.Sprintf("/provision/neural/classes/%s/operations", classID))

	if err != nil {
		return nil, fmt.Errorf("failed to create operation: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &operationResp.Data, nil
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
