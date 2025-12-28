package content

import (
	"content/internal/lib/api"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const basePath = "/api/content"

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewContentHandler(service Service, logger *slog.Logger) *Handler {
	handler := &Handler{
		logger:  logger,
		service: service,
	}

	handler.logger = logger.With(
		slog.String("caller", "handlers.content.api"),
	)

	return handler
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get(basePath+"/{id}", h.getById)
	r.Get(basePath, h.getByFilter)
	r.Post(basePath, h.create)
	r.Put(basePath, h.update)
	r.Delete(basePath+"/{id}", h.delete)
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
	err := &api.ValidationErrors{}

	if filter.FolderId == "" {
		err.Add("folderId", "is required")
	}
	if filter.UserId == "" {
		err.Add("userId", "is required")
	}
	if !err.Ok() {
		respondError(w, r, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.GetByFilter(r.Context(), filter)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request CreateContentRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		respondError(w, r, http.StatusBadRequest, api.NewValidationErrors(err.Error()))
		return
	}

	// todo: авторизация и подстановка userId текущего авторизованного пользователя
	response, err := h.service.Create(r.Context(), request)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
	}

	respond(w, r, http.StatusOK, response)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var request UpdateContentRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		respondError(w, r, http.StatusBadRequest, api.NewValidationErrors(err.Error()))
		return
	}

	// todo: авторизация и проверка на владение сущностью
	err := h.service.Update(r.Context(), request)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		err := &api.ValidationErrors{}
		err.Add("id", "is required")
		respondError(w, r, http.StatusBadRequest, err)

		return
	}

	// todo: авторизация и проверка на владение сущностью
	err := h.service.SafeDelete(r.Context(), id)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, err)
}

func getFilter(r *http.Request) Filter {
	return Filter{
		UserId:   r.URL.Query().Get("userId"),
		Search:   r.URL.Query().Get("search"),
		FolderId: r.URL.Query().Get("folderId"),
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
