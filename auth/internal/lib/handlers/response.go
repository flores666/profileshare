package handlers

import (
	"net/http"

	"github.com/flores666/profileshare-lib/api"

	"github.com/go-chi/render"
)

func Respond(w http.ResponseWriter, r *http.Request, status int, response api.AppResponse) {
	render.Status(r, status)
	render.JSON(w, r, response)
}
