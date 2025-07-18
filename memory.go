package tama

import (
	"github.com/upmaru/tama-go/memory"
)

// MemoryService handles all memory-related API operations
//
// The memory service operations are organized in a separate package:
// - memory/prompt.go: Prompt operations (memory prompts with content, roles, and associated spaces).
type MemoryService struct {
	*memory.Service
}

// newMemoryService creates a new memory service instance.
func newMemoryService(client *Client) *MemoryService {
	return &MemoryService{
		Service: memory.NewService(client.httpClient),
	}
}
