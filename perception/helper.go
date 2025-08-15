package perception

import (
	"errors"
	"fmt"
)

// genericGet is a helper function for GET operations.
func genericGet[T any](
	s *Service,
	id string,
	resourceName string,
	pathFormat string,
	result T,
) error {
	if id == "" {
		return fmt.Errorf("%s ID is required", resourceName)
	}

	resp, err := s.client.R().
		SetResult(result).
		Get(fmt.Sprintf(pathFormat, id))

	if err != nil {
		return fmt.Errorf("failed to get %s: %w", resourceName, err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

// genericCreate is a helper function for POST operations.
func genericCreate[T any](
	s *Service,
	parentID string,
	req any,
	resourceName string,
	parentIDName string,
	pathFormat string,
	result T,
) error {
	if parentID == "" {
		return fmt.Errorf("%s ID is required", parentIDName)
	}

	resp, err := s.client.R().
		SetBody(req).
		SetResult(result).
		Post(fmt.Sprintf(pathFormat, parentID))

	if err != nil {
		return fmt.Errorf("failed to create %s: %w", resourceName, err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

// genericUpdate is a helper function for PATCH operations.
func genericUpdate[T any](
	s *Service,
	id string,
	req any,
	resourceName string,
	pathFormat string,
	result T,
) error {
	if id == "" {
		return fmt.Errorf("%s ID is required", resourceName)
	}

	resp, err := s.client.R().
		SetBody(req).
		SetResult(result).
		Patch(fmt.Sprintf(pathFormat, id))

	if err != nil {
		return fmt.Errorf("failed to update %s: %w", resourceName, err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

// genericReplace is a helper function for PUT operations.
func genericReplace[T any](
	s *Service,
	id string,
	req any,
	resourceName string,
	pathFormat string,
	result T,
) error {
	if id == "" {
		return fmt.Errorf("%s ID is required", resourceName)
	}

	resp, err := s.client.R().
		SetBody(req).
		SetResult(result).
		Put(fmt.Sprintf(pathFormat, id))

	if err != nil {
		return fmt.Errorf("failed to replace %s: %w", resourceName, err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}

// genericDelete is a helper function for DELETE operations.
func genericDelete(s *Service, id string, resourceName string, pathFormat string) error {
	if id == "" {
		return errors.New(resourceName + " ID is required")
	}

	resp, err := s.client.R().
		Delete(fmt.Sprintf(pathFormat, id))

	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", resourceName, err)
	}

	if apiErr := s.handleAPIError(resp); apiErr != nil {
		return apiErr
	}

	return nil
}
