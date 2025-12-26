package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// GetBodyWithValidation reads the request body, decodes JSON into out,
// and validates the resulting structure.
// It returns an error if decoding or validation fails.
func GetBodyWithValidation(r *http.Request, out interface{}) error {
	err := render.DecodeJSON(r.Body, &out)

	if err != nil {
		const message = "failed to decode body"

		return errors.New(message)
	}

	if err = validator.New().Struct(out); err != nil {
		const message = "body validation error"

		return errors.New(message)
	}

	return nil
}
