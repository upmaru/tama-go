package sensory_test

import (
	"strings"
	"testing"

	"github.com/upmaru/tama-go/sensory"
)

func TestSensoryFieldSpecificErrors(t *testing.T) {
	// Test sensory field-specific errors
	fieldErr := &sensory.Error{
		StatusCode: 422,
		Errors: map[string][]string{
			"source_id": {"has already been taken"},
			"name":      {"is required", "must be at least 3 characters"},
		},
	}

	errorMsg := fieldErr.Error()
	// Check that all field errors are included
	if !strings.Contains(errorMsg, "source_id has already been taken") {
		t.Errorf("Expected error message to contain 'source_id has already been taken', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "name is required") {
		t.Errorf("Expected error message to contain 'name is required', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "name must be at least 3 characters") {
		t.Errorf("Expected error message to contain 'name must be at least 3 characters', got %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "API error 422:") {
		t.Errorf("Expected error message to contain status code, got %s", errorMsg)
	}

	// Test error with only status code
	statusOnlyErr := &sensory.Error{
		StatusCode: 404,
	}

	expectedStatusMsg := "API error 404"
	if statusOnlyErr.Error() != expectedStatusMsg {
		t.Errorf("Expected error message %s, got %s", expectedStatusMsg, statusOnlyErr.Error())
	}

	// Test field-specific errors without status code
	fieldErrNoStatus := &sensory.Error{
		Errors: map[string][]string{
			"endpoint": {"is invalid URL"},
		},
	}

	errorMsgNoStatus := fieldErrNoStatus.Error()
	expectedNoStatus := "API error: endpoint is invalid URL"
	if errorMsgNoStatus != expectedNoStatus {
		t.Errorf("Expected error message %s, got %s", expectedNoStatus, errorMsgNoStatus)
	}
}
