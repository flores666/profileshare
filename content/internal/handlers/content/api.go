package content

import (
	"content/internal/lib/api"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const basePath = "/api/content"
const caller = "handlers.content.api"

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewContentHandler(logger *slog.Logger) *Handler {
	handler := &Handler{
		logger: logger,
	}

	handler.logger = logger.With(
		slog.String("caller", caller),
	)

	return handler
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post(basePath, h.createContent)
	r.Get(basePath+"/{id}", h.getContent)
}

func (h *Handler) createContent(w http.ResponseWriter, r *http.Request) {
	var request CreateRequest

	if err := api.GetBodyWithValidation(r, &request); err != nil {
		render.JSON(w, r, api.NewError(err.Error()))
		h.logger.Warn(err.Error())

		return
	}

	fmt.Println(request)
}

func (h *Handler) getContent(w http.ResponseWriter, r *http.Request) {

}
