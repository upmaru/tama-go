package neural_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
)

func TestNeuralGetCorpus(t *testing.T) {
	expectedCorpus := neural.Corpus{
		ID:             "corpus-123",
		Main:           true,
		Name:           "test-corpus",
		Slug:           "test-corpus-slug",
		Template:       "test-template",
		ProvisionState: "active",
	}

	expectedResponse := neural.CorpusResponse{
		Data: expectedCorpus,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/corpora/corpus-123" {
			t.Errorf("Expected path /provision/neural/corpora/corpus-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	corpus, err := client.Neural.GetCorpus("corpus-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if corpus.ID != expectedCorpus.ID {
		t.Errorf("Expected corpus ID %s, got %s", expectedCorpus.ID, corpus.ID)
	}

	if corpus.Main != expectedCorpus.Main {
		t.Errorf("Expected corpus main %v, got %v", expectedCorpus.Main, corpus.Main)
	}

	if corpus.Name != expectedCorpus.Name {
		t.Errorf("Expected corpus name %s, got %s", expectedCorpus.Name, corpus.Name)
	}

	if corpus.Template != expectedCorpus.Template {
		t.Errorf("Expected corpus template %s, got %s", expectedCorpus.Template, corpus.Template)
	}
}

func TestNeuralGetCorpusError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"corpus": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	_, err := client.Neural.GetCorpus("nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var neuralErr *neural.Error
	if errors.As(err, &neuralErr) {
		if neuralErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", neuralErr.StatusCode)
		}
		if neuralErr.Errors == nil || len(neuralErr.Errors["corpus"]) == 0 ||
			neuralErr.Errors["corpus"][0] != "not found" {
			t.Errorf("Expected error 'corpus not found', got %v", neuralErr.Errors)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralCreateCorpus(t *testing.T) {
	expectedCorpus := neural.Corpus{
		ID:             "corpus-789",
		Main:           false,
		Name:           "new-corpus",
		Slug:           "new-corpus-slug",
		Template:       "new-template",
		ProvisionState: "active",
	}

	expectedResponse := neural.CorpusResponse{
		Data: expectedCorpus,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-123/corpora" {
			t.Errorf("Expected path /provision/neural/classes/class-123/corpora, got %s", r.URL.Path)
		}

		var req neural.CreateCorpusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Corpus.Name != "new-corpus" {
			t.Errorf("Expected request name 'new-corpus', got %s", req.Corpus.Name)
		}

		if req.Corpus.Template != "new-template" {
			t.Errorf("Expected request template 'new-template', got %s", req.Corpus.Template)
		}

		if req.Corpus.Main != false {
			t.Errorf("Expected request main false, got %v", req.Corpus.Main)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	createReq := neural.CreateCorpusRequest{
		Corpus: neural.CorpusRequestData{
			Main:     false,
			Name:     "new-corpus",
			Template: "new-template",
		},
	}

	corpus, err := client.Neural.CreateCorpus("class-123", createReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if corpus.ID != expectedCorpus.ID {
		t.Errorf("Expected corpus ID %s, got %s", expectedCorpus.ID, corpus.ID)
	}

	if corpus.Name != expectedCorpus.Name {
		t.Errorf("Expected corpus name %s, got %s", expectedCorpus.Name, corpus.Name)
	}

	if corpus.Template != expectedCorpus.Template {
		t.Errorf("Expected corpus template %s, got %s", expectedCorpus.Template, corpus.Template)
	}
}

func TestNeuralCreateCorpusValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	// Test empty class ID validation
	_, err := client.Neural.CreateCorpus("", neural.CreateCorpusRequest{
		Corpus: neural.CorpusRequestData{
			Main:     true,
			Name:     "test-corpus",
			Template: "test-template",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty class ID")
	}

	// Test empty name validation
	_, err = client.Neural.CreateCorpus("class-123", neural.CreateCorpusRequest{
		Corpus: neural.CorpusRequestData{
			Main:     true,
			Template: "test-template",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	// Test empty template validation
	_, err = client.Neural.CreateCorpus("class-123", neural.CreateCorpusRequest{
		Corpus: neural.CorpusRequestData{
			Main: true,
			Name: "test-corpus",
		},
	})

	if err == nil {
		t.Error("Expected validation error for empty template")
	}
}

func TestNeuralUpdateCorpus(t *testing.T) {
	expectedCorpus := neural.Corpus{
		ID:             "corpus-123",
		Main:           true,
		Name:           "updated-corpus",
		Slug:           "updated-corpus-slug",
		Template:       "updated-template",
		ProvisionState: "active",
	}

	expectedResponse := neural.CorpusResponse{
		Data: expectedCorpus,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/corpora/corpus-123" {
			t.Errorf("Expected path /provision/neural/corpora/corpus-123, got %s", r.URL.Path)
		}

		var req neural.UpdateCorpusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	main := true
	updateReq := neural.UpdateCorpusRequest{
		Corpus: neural.UpdateCorpusData{
			Main:     &main,
			Name:     "updated-corpus",
			Template: "updated-template",
		},
	}

	corpus, err := client.Neural.UpdateCorpus("corpus-123", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if corpus.Name != expectedCorpus.Name {
		t.Errorf("Expected corpus name %s, got %s", expectedCorpus.Name, corpus.Name)
	}

	if corpus.Template != expectedCorpus.Template {
		t.Errorf("Expected corpus template %s, got %s", expectedCorpus.Template, corpus.Template)
	}
}

func TestNeuralReplaceCorpus(t *testing.T) {
	expectedCorpus := neural.Corpus{
		ID:             "corpus-123",
		Main:           false,
		Name:           "replaced-corpus",
		Slug:           "replaced-corpus-slug",
		Template:       "replaced-template",
		ProvisionState: "active",
	}

	expectedResponse := neural.CorpusResponse{
		Data: expectedCorpus,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/corpora/corpus-123" {
			t.Errorf("Expected path /provision/neural/corpora/corpus-123, got %s", r.URL.Path)
		}

		var req neural.UpdateCorpusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	main := false
	replaceReq := neural.UpdateCorpusRequest{
		Corpus: neural.UpdateCorpusData{
			Main:     &main,
			Name:     "replaced-corpus",
			Template: "replaced-template",
		},
	}

	corpus, err := client.Neural.ReplaceCorpus("corpus-123", replaceReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if corpus.Name != expectedCorpus.Name {
		t.Errorf("Expected corpus name %s, got %s", expectedCorpus.Name, corpus.Name)
	}

	if corpus.Template != expectedCorpus.Template {
		t.Errorf("Expected corpus template %s, got %s", expectedCorpus.Template, corpus.Template)
	}

	if corpus.Main != expectedCorpus.Main {
		t.Errorf("Expected corpus main %v, got %v", expectedCorpus.Main, corpus.Main)
	}
}

func TestNeuralGetCorpusByClassAndSlug(t *testing.T) {
	expectedCorpus := neural.Corpus{
		ID:             "corpus-456",
		Main:           false,
		Name:           "test-corpus-by-slug",
		Slug:           "test-slug",
		Template:       "test-template",
		ProvisionState: "active",
	}

	expectedResponse := neural.CorpusResponse{
		Data: expectedCorpus,
	}

	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/classes/class-789/corpora/test-slug" {
			t.Errorf("Expected path /provision/neural/classes/class-789/corpora/test-slug, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	corpus, err := client.Neural.GetCorpusByClassAndSlug("class-789", "test-slug")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if corpus.ID != expectedCorpus.ID {
		t.Errorf("Expected corpus ID %s, got %s", expectedCorpus.ID, corpus.ID)
	}

	if corpus.Main != expectedCorpus.Main {
		t.Errorf("Expected corpus main %v, got %v", expectedCorpus.Main, corpus.Main)
	}

	if corpus.Name != expectedCorpus.Name {
		t.Errorf("Expected corpus name %s, got %s", expectedCorpus.Name, corpus.Name)
	}

	if corpus.Slug != expectedCorpus.Slug {
		t.Errorf("Expected corpus slug %s, got %s", expectedCorpus.Slug, corpus.Slug)
	}

	if corpus.Template != expectedCorpus.Template {
		t.Errorf("Expected corpus template %s, got %s", expectedCorpus.Template, corpus.Template)
	}
}

func TestNeuralGetCorpusByClassAndSlugValidation(t *testing.T) {
	config := tama.Config{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	// Test empty class ID validation
	_, err := client.Neural.GetCorpusByClassAndSlug("", "test-slug")
	if err == nil {
		t.Error("Expected validation error for empty class ID")
	}
	if err.Error() != "class ID is required" {
		t.Errorf("Expected error 'class ID is required', got %s", err.Error())
	}

	// Test empty slug validation
	_, err = client.Neural.GetCorpusByClassAndSlug("class-789", "")
	if err == nil {
		t.Error("Expected validation error for empty slug")
	}
	if err.Error() != "corpus slug is required" {
		t.Errorf("Expected error 'corpus slug is required', got %s", err.Error())
	}
}

func TestNeuralGetCorpusByClassAndSlugError(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResp := neural.Error{
			StatusCode: 404,
			Errors: map[string][]string{
				"corpus": {"not found"},
			},
		}
		json.NewEncoder(w).Encode(errorResp)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)
	_, err := client.Neural.GetCorpusByClassAndSlug("class-789", "nonexistent-slug")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var neuralErr *neural.Error
	if errors.As(err, &neuralErr) {
		if neuralErr.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", neuralErr.StatusCode)
		}
		if neuralErr.Errors == nil || len(neuralErr.Errors["corpus"]) == 0 ||
			neuralErr.Errors["corpus"][0] != "not found" {
			t.Errorf("Expected error 'corpus not found', got %v", neuralErr.Errors)
		}
	} else {
		t.Errorf("Expected neural.Error, got %T", err)
	}
}

func TestNeuralDeleteCorpus(t *testing.T) {
	server := CreateMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/provision/neural/corpora/corpus-123" {
			t.Errorf("Expected path /provision/neural/corpora/corpus-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	config := tama.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := tama.NewClient(config)

	err := client.Neural.DeleteCorpus("corpus-123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
