package content

import (
	"content/internal/lib/api"
	"content/internal/lib/logger/sl"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const basePath = "/api/content"
const caller = "handlers.content.api"

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewContentHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post(basePath, h.createContent)
	r.Get(basePath+"/{id}", h.getContent)
}

func (h *Handler) createContent(w http.ResponseWriter, r *http.Request) {
	var request CreateRequest

	if err := h.getBody(w, r, &request); err != nil {
		return
	}

	fmt.Println(request)
}

func (h *Handler) getContent(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getBody(w http.ResponseWriter, r *http.Request, out interface{}) error {
	h.logger = h.logger.With(
		slog.String("caller", caller),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	err := render.DecodeJSON(r.Body, &out)

	if err != nil {
		const message = "failed to decode body"

		h.logger.Error(message, sl.Error(err))
		render.JSON(w, r, api.NewError(message))

		return err
	}

	h.logger.Debug("body decoded", slog.Any("request", out))

	if err := validator.New().Struct(out); err != nil {
		h.logger.Warn("validation error", sl.Error(err))

		render.JSON(w, r, api.NewError("body validation error"))

		return err
	}

	return nil
}
