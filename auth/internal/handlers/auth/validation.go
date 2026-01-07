package auth

import (
	"github.com/flores666/profileshare-lib/api"
)

func validateRegister(request RegisterUserRequest) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if len([]rune(request.Nickname)) < 2 {
		errs.Add("nickname", "must contain at least 2 characters")
	}

	if len([]rune(request.Email)) < 2 {
		errs.Add("email", "must contain at least 2 characters")
	}

	if len([]rune(request.Password)) < 8 {
		errs.Add("password", "must contain at least 8 characters")
	}

	if request.ReturnUrl == "" {
		errs.Add("returnUrl", "is required")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}

func validateLogin(request LoginUserRequest) *api.ValidationErrors {
	errs := &api.ValidationErrors{}

	if len([]rune(request.Email)) < 2 {
		errs.Add("email", "must contain at least 2 characters")
	}

	if len([]rune(request.Password)) < 8 {
		errs.Add("password", "must contain at least 8 characters")
	}

	if errs.Ok() {
		return nil
	}

	return errs
}
