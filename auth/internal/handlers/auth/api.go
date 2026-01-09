package auth

import (
	"auth/internal/handlers/auth/security"
	"auth/internal/lib/handlers"
	"net/http"
	"strings"

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
	r.Post(BaseRoutePath+"/login", h.login)
	r.Post(BaseRoutePath+"/logout", h.logout)
	r.Post(BaseRoutePath+"/refresh", h.refresh)
	r.Post(BaseRoutePath+"/confirm", h.confirm)
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

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var request LoginUserRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		handlers.Respond(w, r, http.StatusBadRequest, api.NewError(err.Error(), nil))
		return
	}

	result := h.service.Login(r.Context(), request)
	if !result.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, result)
		return
	}

	if tokens, ok := result.Data.(*security.TokenPair); ok {
		createTokenCookie(w, tokens.RefreshToken)
	}

	handlers.Respond(w, r, http.StatusOK, result)
}

func (h *Handler) confirm(w http.ResponseWriter, r *http.Request) {
	var request ConfirmUserRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		handlers.Respond(w, r, http.StatusBadRequest, api.NewError(err.Error(), nil))
		return
	}

	result := h.service.Confirm(r.Context(), request)
	if !result.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, result)
		return
	}

	handlers.Respond(w, r, http.StatusOK, result)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	var request LogoutRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		cookie, cerr := r.Cookie("rt")

		if cerr != nil {
			handlers.Respond(w, r, http.StatusBadRequest, api.NewError(err.Error(), nil))
			return
		}

		request.RefreshToken = cookie.Value
	}

	request.AccessToken = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	result := h.service.Logout(r.Context(), request)
	if !result.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, result)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "rt",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	handlers.Respond(w, r, http.StatusOK, result)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var request RefreshTokenRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		cookie, cerr := r.Cookie("rt")

		if cerr != nil {
			handlers.Respond(w, r, http.StatusBadRequest, api.NewError(err.Error(), nil))
			return
		}

		request.RefreshToken = cookie.Value
	}

	result := h.service.RefreshTokens(r.Context(), request)
	if !result.Ok() {
		handlers.Respond(w, r, http.StatusInternalServerError, result)
		return
	}

	if tokens, ok := result.Data.(*security.TokenPair); ok {
		createTokenCookie(w, tokens.RefreshToken)
	}

	handlers.Respond(w, r, http.StatusOK, result)
}

func createTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "rt",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})
}
