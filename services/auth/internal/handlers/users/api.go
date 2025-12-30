package users

import (
	"api"
	"auth/internal/lib/handlers"
	"net/http"

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
		handlers.Error(w, r, http.StatusBadRequest, api.NewValidationErrors(err.Error()))
		return
	}

	err := h.service.Update(r.Context(), request)
	if err != nil {
		handlers.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	handlers.Respond(w, r, http.StatusOK, nil)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		err := &api.ValidationErrors{}
		err.Add("id", "is required")
		handlers.Error(w, r, http.StatusBadRequest, err)

		return
	}

	resp, err := h.service.GetById(r.Context(), id)
	if err != nil {
		handlers.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	handlers.Respond(w, r, http.StatusOK, resp)
}

func (h *Handler) getByFilter(w http.ResponseWriter, r *http.Request) {
	filter := getFilter(r)
	resp, err := h.service.GetByFilter(r.Context(), filter)
	if err != nil {
		handlers.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	handlers.Respond(w, r, http.StatusOK, resp)
}

func getFilter(r *http.Request) QueryFilter {
	return QueryFilter{
		Search: r.URL.Query().Get("search"),
	}
}
