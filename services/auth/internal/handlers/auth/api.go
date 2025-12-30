package auth

import (
	"api"
	"auth/internal/lib/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const BaseRoutePath = "/api/auth"

type Handler struct {
	service Service
}

func NewAuthHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post(BaseRoutePath+"/register", h.register)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var request RegisterUserRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		handlers.Error(w, r, http.StatusBadRequest, api.NewValidationErrors(err.Error()))
		return
	}

	result, err := h.service.Register(r.Context(), request)
	if err != nil {
		handlers.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	handlers.Respond(w, r, http.StatusOK, result)
}
