package users

import (
	"auth/internal/lib/handlers"
	"net/http"

	"github.com/flores666/profileshare-lib/api"

	"github.com/go-chi/chi/v5"
)

const BaseRoutePath = "/api/users"

type Handler struct {
	service Service
}

func NewUsersHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Put(BaseRoutePath, h.update)
	r.Get(BaseRoutePath+"/{id}", h.getById)
	r.Get(BaseRoutePath, h.getByFilter)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var request UpdateUserRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		handlers.Respond(w, r, http.StatusBadRequest, api.NewError(err.Error(), nil))
		return
	}

	result := h.service.Update(r.Context(), request)
	if !result.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, result)
		return
	}

	handlers.Respond(w, r, http.StatusOK, result)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		handlers.Respond(w, r, http.StatusBadRequest, api.NewError("Отсутствует id пользователя", nil))
		return
	}

	response := h.service.GetById(r.Context(), id)
	if !response.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, response)
		return
	}

	handlers.Respond(w, r, http.StatusOK, response)
}

func (h *Handler) getByFilter(w http.ResponseWriter, r *http.Request) {
	filter := getFilter(r)
	response := h.service.GetByFilter(r.Context(), filter)
	if !response.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, response)
		return
	}

	handlers.Respond(w, r, http.StatusOK, response)
}

func getFilter(r *http.Request) QueryFilter {
	return QueryFilter{
		Search: r.URL.Query().Get("search"),
	}
}
