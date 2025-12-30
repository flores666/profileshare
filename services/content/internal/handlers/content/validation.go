package content

import (
	"slices"

	"api"
)

var contentTypes = []string{"photo", "video"}

func validateCreate(request CreateContentRequest) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if request.UserId == "" {
		errs.Add("userId", "is required")
	}

	if request.DisplayName == "" || len([]rune(request.DisplayName)) <= 2 {
		errs.Add("displayName", "must be at least 2 characters")
	}

	if request.FolderId == "" {
		errs.Add("folderId", "is required")
	}

	if request.Type == "" {
		errs.Add("type", "is required")
	}

	if !slices.Contains(contentTypes, request.Type) {
		errs.Add("type", "invalid input")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}

func validateFilter(filter Filter) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if filter.FolderId == "" {
		errs.Add("folderId", "is required")
	}

	if filter.UserId == "" {
		errs.Add("userId", "is required")
	}

	if len([]rune(filter.Search)) < 2 {
		errs.Add("search", "must contain at least 2 characters")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}

func validateUpdate(request UpdateContentRequest) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if request.DisplayName != nil && len([]rune(*request.DisplayName)) <= 2 {
		errs.Add("displayName", "must be at least 2 characters")
	}

	if request.Id == "" {
		errs.Add("id", "is required")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}

func validateId(id string) *api.ValidationErrors {
	errs := &api.ValidationErrors{}
	if id == "" {
		errs.Add("id", "is required")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}
