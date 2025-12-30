package users

import (
	"api"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const BaseRoutePath = "/api/users"

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewUsersHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h Handler) RegisterRoutes(r chi.Router) {
	r.Get(BaseRoutePath, h.test)
}

func (h Handler) test(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "Hello World")
}

func getFilter(r *http.Request) Filter {
	return Filter{
		Search: r.URL.Query().Get("search"),
	}
}

func respond(w http.ResponseWriter, r *http.Request, status int, payload any) {
	render.Status(r, status)
	if payload != nil {
		render.JSON(w, r, api.NewOk(payload))
	}
}

func respondError(w http.ResponseWriter, r *http.Request, status int, err *api.ValidationErrors) {
	render.Status(r, status)
	render.JSON(w, r, api.NewError(err.Message, err))
}
