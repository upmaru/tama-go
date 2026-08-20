package tama

import (
	"github.com/upmaru/tama-go/tools"
)

// ToolsService handles all tools-related API operations.
//
// The tools service operations are organized in a separate package:
// - tools/input.go: Input operations (tools inputs with type, class corpus ID, and CRUD operations)
// - tools/initializer.go: Initializer operations (reference, index, parameters, and CRUD operations)
// - tools/output.go: Output operations (class corpus ID and CRUD operations)
// - tools/option.go: Output option operations (action modifier association and CRUD operations)
// - tools/modifier.go: Trusted thought-tool modifier operations (runtime metadata projection and CRUD operations).
type ToolsService struct {
	*tools.Service
}

// newToolsService creates a new tools service instance.
func newToolsService(client *Client) *ToolsService {
	return &ToolsService{
		Service: tools.NewService(client.httpClient),
	}
}
