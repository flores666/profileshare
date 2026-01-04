package auth

import (
	"auth/internal/lib/handlers"
	"net/http"

	"github.com/flores666/profileshare-lib/api"

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
		handlers.Respond(w, r, http.StatusBadRequest, api.NewError(err.Error(), nil))
		return
	}

	result := h.service.Register(r.Context(), request)
	if !result.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, result)
		return
	}

	handlers.Respond(w, r, http.StatusOK, result)
}
