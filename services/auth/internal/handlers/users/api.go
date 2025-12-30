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

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post(BaseRoutePath, h.create)
	r.Put(BaseRoutePath, h.update)
	r.Get(BaseRoutePath+"/{id}", h.getById)
	r.Get(BaseRoutePath, h.getByFilter)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		respondError(w, r, http.StatusBadRequest, api.NewValidationErrors(err.Error()))
		return
	}

	result, err := h.service.Create(r.Context(), request)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, result)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var request UpdateUserRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		respondError(w, r, http.StatusBadRequest, api.NewValidationErrors(err.Error()))
		return
	}

	err := h.service.Update(r.Context(), request)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		err := &api.ValidationErrors{}
		err.Add("id", "is required")
		respondError(w, r, http.StatusBadRequest, err)

		return
	}

	response, err := h.service.GetById(r.Context(), id)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func (h *Handler) getByFilter(w http.ResponseWriter, r *http.Request) {
	filter := getFilter(r)
	response, err := h.service.GetByFilter(r.Context(), filter)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func getFilter(r *http.Request) QueryFilter {
	return QueryFilter{
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
