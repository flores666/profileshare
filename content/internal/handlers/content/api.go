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
	r.Post(basePath, h.create)
	r.Get(basePath+"/{id}", h.getById)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request CreateContentRequest

	if err := api.GetBodyWithValidation(r, &request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.NewError(err.Error()))
		h.logger.Warn(err.Error())

		return
	}

	response := h.service.CreateContent(request)

	if !response.IsOk() {
		h.logger.Warn(response.Error)
		render.Status(r, http.StatusInternalServerError)
	}

	render.JSON(w, r, response)
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

	item, err := h.service.GetContentById(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, api.NewError(err.Error()))
		h.logger.Info(err.Error())

		return
	}

	render.JSON(w, r, item)
}
