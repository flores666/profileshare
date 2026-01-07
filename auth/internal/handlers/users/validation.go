package users

import (
	"github.com/flores666/profileshare-lib/api"
	"time"
)

func validateId(id string) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if id == "" {
		errs.Add("Id", "is required")
	}

	if errs.Ok() {
		return nil
	}

	return nil
}

func validateFilter(filter QueryFilter) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if filter.Search != "" && len([]rune(filter.Search)) < 2 {
		errs.Add("search", "must contain at least 2 characters")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}

func validateUpdate(request UpdateUserRequest) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if request.Email != nil && len([]rune(*request.Email)) < 2 {
		errs.Add("email", "must contain at least 2 characters")
	}

	if request.Nickname != nil && len([]rune(*request.Nickname)) < 2 {
		errs.Add("nickname", "must contain at least 2 characters")
	}

	if request.Id == "" {
		errs.Add("id", "is required")
	}

	if request.BannedBefore != nil && (*request.BannedBefore).Before(time.Now().UTC()) {
		errs.Add("bannedBefore", "is invalid")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}
