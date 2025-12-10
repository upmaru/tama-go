package perception_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/upmaru/tama-go/perception"
)

// createMockServer creates a test HTTP server with the given handler.
func createMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// ValidateThoughtResponse validates that actual thought matches expected thought.
func ValidateThoughtResponse(t *testing.T, actual, expected perception.Thought) {
	if actual.ID != expected.ID {
		t.Errorf("Expected thought ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ChainID != expected.ChainID {
		t.Errorf("Expected thought chain_id %s, got %s", expected.ChainID, actual.ChainID)
	}

	if actual.OutputClassID != expected.OutputClassID {
		t.Errorf(
			"Expected thought output_class_id %s, got %s",
			expected.OutputClassID,
			actual.OutputClassID,
		)
	}

	// Handle Module pointer comparison
	if (actual.Module == nil) != (expected.Module == nil) {
		t.Errorf(
			"Expected module nil status %v, got %v",
			expected.Module == nil,
			actual.Module == nil,
		)
	}
	if actual.Module != nil && expected.Module != nil {
		if actual.Module.ID != expected.Module.ID {
			t.Errorf("Expected module ID %s, got %s", expected.Module.ID, actual.Module.ID)
		}
		if actual.Module.Reference != expected.Module.Reference {
			t.Errorf(
				"Expected module reference %s, got %s",
				expected.Module.Reference,
				actual.Module.Reference,
			)
		}
	}

	// Handle Delegation pointer comparison
	if (actual.Delegation == nil) != (expected.Delegation == nil) {
		t.Errorf(
			"Expected delegation nil status %v, got %v",
			expected.Delegation == nil,
			actual.Delegation == nil,
		)
	}
	if actual.Delegation != nil && expected.Delegation != nil {
		if actual.Delegation.TargetThoughtID != expected.Delegation.TargetThoughtID {
			t.Errorf(
				"Expected delegation target_thought_id %s, got %s",
				expected.Delegation.TargetThoughtID,
				actual.Delegation.TargetThoughtID,
			)
		}
	}

	// Handle Faculty pointer comparison
	if (actual.Faculty == nil) != (expected.Faculty == nil) {
		t.Errorf(
			"Expected faculty nil status %v, got %v",
			expected.Faculty == nil,
			actual.Faculty == nil,
		)
	}
	if actual.Faculty != nil && expected.Faculty != nil {
		if actual.Faculty.QueueID != expected.Faculty.QueueID {
			t.Errorf(
				"Expected faculty queue_id %s, got %s",
				expected.Faculty.QueueID,
				actual.Faculty.QueueID,
			)
		}
		if actual.Faculty.Priority != expected.Faculty.Priority {
			t.Errorf(
				"Expected faculty priority %d, got %d",
				expected.Faculty.Priority,
				actual.Faculty.Priority,
			)
		}
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf(
			"Expected thought provision_state %s, got %s",
			expected.ProvisionState,
			actual.ProvisionState,
		)
	}

	if actual.Relation != expected.Relation {
		t.Errorf("Expected thought relation %s, got %s", expected.Relation, actual.Relation)
	}

	if actual.Index != expected.Index {
		t.Errorf("Expected thought index %d, got %d", expected.Index, actual.Index)
	}
}

// ValidatePathResponse validates that actual path matches expected path.
func ValidatePathResponse(t *testing.T, actual, expected perception.Path) {
	if actual.ID != expected.ID {
		t.Errorf("Expected path ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ThoughtID != expected.ThoughtID {
		t.Errorf("Expected path thought_id %s, got %s", expected.ThoughtID, actual.ThoughtID)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf(
			"Expected path provision_state %s, got %s",
			expected.ProvisionState,
			actual.ProvisionState,
		)
	}

	if actual.TargetClassID != expected.TargetClassID {
		t.Errorf(
			"Expected path target_class_id %s, got %s",
			expected.TargetClassID,
			actual.TargetClassID,
		)
	}
}

// ValidateContextResponse validates that actual context matches expected context.
func ValidateContextResponse(t *testing.T, actual, expected perception.Context) {
	if actual.ID != expected.ID {
		t.Errorf("Expected context ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ThoughtID != expected.ThoughtID {
		t.Errorf("Expected context thought_id %s, got %s", expected.ThoughtID, actual.ThoughtID)
	}

	if actual.PromptID != expected.PromptID {
		t.Errorf("Expected context prompt_id %s, got %s", expected.PromptID, actual.PromptID)
	}

	if actual.Layer != expected.Layer {
		t.Errorf("Expected context layer %d, got %d", expected.Layer, actual.Layer)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf(
			"Expected context provision_state %s, got %s",
			expected.ProvisionState,
			actual.ProvisionState,
		)
	}
}

// ValidateRequestIndex is a helper function to validate request index based on request count.
func ValidateRequestIndex(t *testing.T, receivedRequest *perception.CreateThoughtRequest, requestCount int) {
	switch requestCount {
	case 1:
		// First request should have index 0
		expectedIndex := 0
		if receivedRequest.Thought.Index == nil {
			t.Errorf("Expected thought index %d, got nil", expectedIndex)
			return
		}
		if *receivedRequest.Thought.Index != expectedIndex {
			t.Errorf("Expected thought index %d, got %d", expectedIndex, *receivedRequest.Thought.Index)
		}
	case 2:
		// Second request should have nil index
		if receivedRequest.Thought.Index != nil {
			t.Errorf("Expected nil index for second request, got %d", *receivedRequest.Thought.Index)
		}
	}
}

// ValidateInitializerResponse validates that actual initializer matches expected initializer.
func ValidateInitializerResponse(t *testing.T, actual, expected perception.Initializer) {
	if actual.ID != expected.ID {
		t.Errorf("Expected initializer ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ThoughtID != expected.ThoughtID {
		t.Errorf("Expected initializer thought_id %s, got %s", expected.ThoughtID, actual.ThoughtID)
	}

	if actual.ClassID != expected.ClassID {
		t.Errorf("Expected initializer class_id %s, got %s", expected.ClassID, actual.ClassID)
	}

	if actual.Reference != expected.Reference {
		t.Errorf("Expected initializer reference %s, got %s", expected.Reference, actual.Reference)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf(
			"Expected initializer provision_state %s, got %s",
			expected.ProvisionState,
			actual.ProvisionState,
		)
	}

	// Handle Index pointer comparison
	if (actual.Index == nil) != (expected.Index == nil) {
		t.Errorf(
			"Expected index nil status %v, got %v",
			expected.Index == nil,
			actual.Index == nil,
		)
	}
	if actual.Index != nil && expected.Index != nil {
		if *actual.Index != *expected.Index {
			t.Errorf("Expected index %d, got %d", *expected.Index, *actual.Index)
		}
	}
}

// ValidateActivationResponse validates that actual activation matches expected activation.
func ValidateActivationResponse(t *testing.T, actual, expected perception.Activation) {
	if actual.ID != expected.ID {
		t.Errorf("Expected activation ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.ThoughtPathID != expected.ThoughtPathID {
		t.Errorf("Expected activation thought_path_id %s, got %s", expected.ThoughtPathID, actual.ThoughtPathID)
	}

	if actual.ChainID != expected.ChainID {
		t.Errorf("Expected activation chain_id %s, got %s", expected.ChainID, actual.ChainID)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf(
			"Expected activation provision_state %s, got %s",
			expected.ProvisionState,
			actual.ProvisionState,
		)
	}
}

// GetKeys returns the keys from a map for debugging purposes.
func GetKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
