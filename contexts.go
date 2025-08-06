package tama

import (
	"github.com/upmaru/tama-go/contexts"
)

// ContextsService handles all contexts-related API operations
//
// The contexts service operations are organized in a separate package:
// - contexts/input.go: Input operations (contexts inputs with types, corpus IDs, and CRUD operations).
type ContextsService struct {
	*contexts.Service
}

// newContextsService creates a new contexts service instance.
func newContextsService(client *Client) *ContextsService {
	return &ContextsService{
		Service: contexts.NewService(client.httpClient),
	}
}
