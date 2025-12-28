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
		const message = "missing id"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.NewError(message))
		h.logger.Warn(message)

		return
	}

	response := h.service.GetById(id)
	if !response.Ok() {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, api.NewError(response.Error))
		h.logger.Info(response.Error)

		return
	}

	render.JSON(w, r, response)
}

func (h *Handler) getByFilter(w http.ResponseWriter, r *http.Request) {
	filter := getFilter(r)
	if filter.FolderId == "" || filter.UserId == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.NewError("Missing parameters: folderId, userId"))
		return
	}

	response := h.service.GetByFilter(filter)
	if !response.Ok() {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.NewError(response.Error))
		h.logger.Warn(response.Error)
		return
	}

	render.JSON(w, r, response)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request CreateContentRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.NewError(err.Error()))
		h.logger.Warn(err.Error())

		return
	}

	// todo: авторизация и подстановка userId текущего авторизованного пользователя
	response := h.service.Create(request)

	if !response.Ok() {
		h.logger.Warn(response.Error)
		render.Status(r, http.StatusInternalServerError)
	}

	render.JSON(w, r, response)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var request UpdateContentRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.NewError(err.Error()))
		h.logger.Warn(err.Error())

		return
	}

	// todo: авторизация и проверка на владение сущностью
	response := h.service.Update(request)

	if !response.Ok() {
		h.logger.Warn(response.Error)
		render.Status(r, http.StatusInternalServerError)
	}

	render.JSON(w, r, response)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		const message = "missing id"
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.NewError(message))
		h.logger.Warn(message)

		return
	}

	// todo: авторизация и проверка на владение сущностью
	response := h.service.SafeDelete(id)

	if !response.Ok() {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.NewError(response.Error))
		h.logger.Info(response.Error)

		return
	}

	render.JSON(w, r, response)
}

func getFilter(r *http.Request) Filter {
	return Filter{
		UserId:   r.URL.Query().Get("userId"),
		Search:   r.URL.Query().Get("search"),
		FolderId: r.URL.Query().Get("folderId"),
	}
}
