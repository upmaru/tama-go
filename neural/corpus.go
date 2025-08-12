package neural

import (
	"errors"
	"fmt"
)

// Corpus represents a neural corpus resource.
type Corpus struct {
	ID             string `json:"id,omitempty"`
	Main           bool   `json:"main"`
	Name           string `json:"name"`
	Slug           string `json:"slug,omitempty"`
	Template       string `json:"template"`
	ProvisionState string `json:"provision_state,omitempty"`
}

// CorpusResponse represents the API response for corpus operations.
type CorpusResponse struct {
	Data Corpus `json:"data"`
}

// CreateCorpusRequest represents the request payload for creating a corpus.
type CreateCorpusRequest struct {
	Corpus CorpusRequestData `json:"corpus"`
}

// CorpusRequestData represents the corpus data in the request.
type CorpusRequestData struct {
	Main     bool   `json:"main"`
	Name     string `json:"name"`
	Template string `json:"template"`
}

// UpdateCorpusRequest represents the request payload for updating a corpus.
type UpdateCorpusRequest struct {
	Corpus UpdateCorpusData `json:"corpus"`
}

// UpdateCorpusData represents the corpus update data.
type UpdateCorpusData struct {
	Main     *bool  `json:"main,omitempty"`
	Name     string `json:"name,omitempty"`
	Template string `json:"template,omitempty"`
}

// GetCorpus retrieves a specific corpus by ID.
// GET /provision/neural/corpora/:id.
func (s *Service) GetCorpus(id string) (*Corpus, error) {
	if id == "" {
		return nil, errors.New("corpus ID is required")
	}

	var corpusResp CorpusResponse
	resp, err := s.client.R().
		SetResult(&corpusResp).
		Get(fmt.Sprintf("/provision/neural/corpora/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get corpus: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &corpusResp.Data, nil
}

// CreateCorpus creates a new corpus within a class.
// POST /provision/neural/classes/:class_id/corpora.
func (s *Service) CreateCorpus(classID string, req CreateCorpusRequest) (*Corpus, error) {
	if classID == "" {
		return nil, errors.New("class ID is required")
	}
	if req.Corpus.Name == "" {
		return nil, errors.New("corpus name is required")
	}
	if req.Corpus.Template == "" {
		return nil, errors.New("corpus template is required")
	}

	var corpusResp CorpusResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&corpusResp).
		Post(fmt.Sprintf("/provision/neural/classes/%s/corpora", classID))

	if err != nil {
		return nil, fmt.Errorf("failed to create corpus: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &corpusResp.Data, nil
}

// UpdateCorpus updates an existing corpus using PATCH.
// PATCH /provision/neural/corpora/:id.
func (s *Service) UpdateCorpus(id string, req UpdateCorpusRequest) (*Corpus, error) {
	if id == "" {
		return nil, errors.New("corpus ID is required")
	}

	var corpusResp CorpusResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&corpusResp).
		Patch(fmt.Sprintf("/provision/neural/corpora/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update corpus: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &corpusResp.Data, nil
}

// ReplaceCorpus replaces an existing corpus using PUT.
// PUT /provision/neural/corpora/:id.
func (s *Service) ReplaceCorpus(id string, req UpdateCorpusRequest) (*Corpus, error) {
	if id == "" {
		return nil, errors.New("corpus ID is required")
	}

	var corpusResp CorpusResponse
	resp, err := s.client.R().
		SetBody(req).
		SetResult(&corpusResp).
		Put(fmt.Sprintf("/provision/neural/corpora/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to replace corpus: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &corpusResp.Data, nil
}

// GetCorpusByClassAndSlug retrieves a specific corpus by class ID and slug.
// GET /provision/neural/classes/:class_id/corpora/:id.
func (s *Service) GetCorpusByClassAndSlug(classID string, slug string) (*Corpus, error) {
	if classID == "" {
		return nil, errors.New("class ID is required")
	}
	if slug == "" {
		return nil, errors.New("corpus slug is required")
	}

	var corpusResp CorpusResponse
	resp, err := s.client.R().
		SetResult(&corpusResp).
		Get(fmt.Sprintf("/provision/neural/classes/%s/corpora/%s", classID, slug))

	if err != nil {
		return nil, fmt.Errorf("failed to get corpus: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return nil, apiErr
	}

	return &corpusResp.Data, nil
}

// DeleteCorpus deletes a corpus by ID.
// DELETE /provision/neural/corpora/:id.
func (s *Service) DeleteCorpus(id string) error {
	if id == "" {
		return errors.New("corpus ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf("/provision/neural/corpora/%s", id))

	if err != nil {
		return fmt.Errorf("failed to delete corpus: %w", err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
