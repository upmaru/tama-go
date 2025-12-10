package tama

import (
	"github.com/upmaru/tama-go/system"
)

// SystemService handles all system-related API operations.
//
// The system service operations are organized in a separate package:
// - system/queue.go: Queue operations (queue role/name definitions with concurrency limits).
type SystemService struct {
	*system.Service
}

// newSystemService creates a new system service instance.
func newSystemService(client *Client) *SystemService {
	return &SystemService{
		Service: system.NewService(client.GetHTTPClient()),
	}
}
