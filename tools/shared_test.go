package tools_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/upmaru/tama-go/tools"
)

// createMockServer creates a test HTTP server with the given handler.
func createMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// ValidateInputResponse validates that actual input matches expected input.
func ValidateInputResponse(t *testing.T, actual, expected tools.Input) {
	if actual.ID != expected.ID {
		t.Errorf("Expected input ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.Type != expected.Type {
		t.Errorf("Expected input type %s, got %s", expected.Type, actual.Type)
	}

	if actual.ThoughtToolID != expected.ThoughtToolID {
		t.Errorf("Expected input thought_tool_id %s, got %s", expected.ThoughtToolID, actual.ThoughtToolID)
	}

	if actual.ClassCorpusID != expected.ClassCorpusID {
		t.Errorf("Expected input class_corpus_id %s, got %s", expected.ClassCorpusID, actual.ClassCorpusID)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf(
			"Expected input provision_state %s, got %s",
			expected.ProvisionState,
			actual.ProvisionState,
		)
	}
}

// ValidateInitializerResponse validates that actual initializer matches expected initializer.
func ValidateInitializerResponse(t *testing.T, actual, expected tools.Initializer) {
	if actual.ID != expected.ID {
		t.Errorf("Expected initializer ID %s, got %s", expected.ID, actual.ID)
	}

	if actual.Reference != expected.Reference {
		t.Errorf("Expected initializer reference %s, got %s", expected.Reference, actual.Reference)
	}

	if actual.Index != expected.Index {
		t.Errorf("Expected initializer index %d, got %d", expected.Index, actual.Index)
	}

	if actual.ThoughtToolID != expected.ThoughtToolID {
		t.Errorf(
			"Expected initializer thought_tool_id %s, got %s",
			expected.ThoughtToolID,
			actual.ThoughtToolID,
		)
	}

	if actual.ProvisionState != expected.ProvisionState {
		t.Errorf(
			"Expected initializer provision_state %s, got %s",
			expected.ProvisionState,
			actual.ProvisionState,
		)
	}

	// Validate parameters map
	if len(actual.Parameters) != len(expected.Parameters) {
		t.Errorf(
			"Expected %d parameters, got %d",
			len(expected.Parameters),
			len(actual.Parameters),
		)
		return
	}

	for key, expectedValue := range expected.Parameters {
		actualValue, exists := actual.Parameters[key]
		if !exists {
			t.Errorf("Expected parameter %s not found in actual parameters", key)
			continue
		}
		// Convert both values to string for comparison to handle JSON number conversions
		if fmt.Sprintf("%v", actualValue) != fmt.Sprintf("%v", expectedValue) {
			t.Errorf(
				"Expected parameter %s to be %v, got %v",
				key,
				expectedValue,
				actualValue,
			)
		}
	}
}
