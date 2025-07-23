package tama

import (
	"github.com/upmaru/tama-go/perception"
)

// PerceptionService handles all perception-related API operations
//
// The perception service operations are organized in a separate package:
// - perception/chain.go: Chain operations (perception chains with names, states, and CRUD operations).
type PerceptionService struct {
	*perception.Service
}

// newPerceptionService creates a new perception service instance.
func newPerceptionService(client *Client) *PerceptionService {
	return &PerceptionService{
		Service: perception.NewService(client.httpClient),
	}
}
