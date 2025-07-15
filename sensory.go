package tama

import (
	"github.com/upmaru/tama-go/sensory"
)

// SensoryService handles all sensory-related API operations
//
// The sensory service operations are organized in a separate package:
// - sensory/source.go: Source operations (data sources with endpoints and credentials)
// - sensory/model.go: Model operations (machine learning models with identifiers and paths)
// - sensory/limit.go: Limit operations (rate limits and restrictions)
type SensoryService struct {
	*sensory.Service
}

// newSensoryService creates a new sensory service instance
func newSensoryService(client *Client) *SensoryService {
	return &SensoryService{
		Service: sensory.NewService(client.httpClient),
	}
}
