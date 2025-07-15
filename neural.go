package tama

import (
	"github.com/upmaru/tama-go/neural"
)

// NeuralService handles all neural-related API operations
//
// The neural service operations are organized in a separate package:
// - neural/space.go: Space operations (neural spaces with names, types, and CRUD operations)
type NeuralService struct {
	*neural.Service
}

// newNeuralService creates a new neural service instance
func newNeuralService(client *Client) *NeuralService {
	return &NeuralService{
		Service: neural.NewService(client.httpClient),
	}
}
