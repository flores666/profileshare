package content

import (
	"errors"
	"slices"
)

var contentTypes = []string{"photo", "video"}

func validateCreate(request CreateContentRequest) error {
	if request.UserId == "" {
		return errors.New("userId is required")
	}

	if request.DisplayName == "" {
		return errors.New("displayName is required")
	}

	if request.FolderId == "" {
		return errors.New("folderId is required")
	}

	if request.Type == "" || !slices.Contains(contentTypes, request.Type) {
		return errors.New("type is required")
	}

	return nil
}

func validateFilter(filter Filter) error {
	if filter.FolderId == "" {
		return errors.New("folder id is required")
	}

	if filter.UserId == "" {
		return errors.New("user id is required")
	}

	return nil
}

func validateUpdate(request UpdateContentRequest) error {
	if request.DisplayName != nil && *request.DisplayName == "" {
		return errors.New("displayName is required")
	}

	if request.Id == "" {
		return errors.New("id is required")
	}

	return nil
}
