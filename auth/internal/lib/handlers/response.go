package handlers

import (
	"github.com/flores666/profileshare-lib/api"
	"net/http"

	"github.com/go-chi/render"
)

func Respond(w http.ResponseWriter, r *http.Request, status int, payload any) {
	render.Status(r, status)
	if payload != nil {
		render.JSON(w, r, api.NewOk(payload))
	}
}

func Error(w http.ResponseWriter, r *http.Request, status int, err *api.ValidationErrors) {
	render.Status(r, status)
	render.JSON(w, r, api.NewError(err.Message, err))
}
