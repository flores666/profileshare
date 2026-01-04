package content

import (
	"net/http"

	"github.com/flores666/profileshare-lib/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	basePath      = "/api/content"
	errValidation = "Ошибка проверки данных"
	errMissingId  = "Отсутствует id"
)

type Handler struct {
	service Service
}

func NewContentHandler(service Service) *Handler {
	handler := &Handler{
		service: service,
	}

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
		respond(w, r, http.StatusBadRequest, api.NewError(errMissingId, nil))
		return
	}

	response := h.service.GetById(r.Context(), id)
	if !response.Ok() {
		respond(w, r, http.StatusInternalServerError, response)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func (h *Handler) getByFilter(w http.ResponseWriter, r *http.Request) {
	filter := getFilter(r)
	err := &api.ValidationErrors{}

	if filter.FolderId == "" {
		err.Add("folderId", "поле обязательное")
	}
	if filter.UserId == "" {
		err.Add("userId", "поле обязательное")
	}
	if !err.Ok() {
		respond(w, r, http.StatusBadRequest, api.NewError(errValidation, err))
		return
	}

	response := h.service.GetByFilter(r.Context(), filter)
	if !response.Ok() {
		respond(w, r, http.StatusInternalServerError, response)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request CreateContentRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		respond(w, r, http.StatusBadRequest, api.NewError(errValidation, nil))
		return
	}

	// todo: авторизация и подстановка userId текущего авторизованного пользователя
	response := h.service.Create(r.Context(), request)
	if !response.Ok() {
		respond(w, r, http.StatusInternalServerError, response)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var request UpdateContentRequest
	if err := api.GetBodyWithValidation(r, &request); err != nil {
		respond(w, r, http.StatusBadRequest, api.NewError(errValidation, nil))
		return
	}

	// todo: авторизация и проверка на владение сущностью
	response := h.service.Update(r.Context(), request)
	if !response.Ok() {
		respond(w, r, http.StatusInternalServerError, response)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respond(w, r, http.StatusBadRequest, api.NewError(errMissingId, nil))
		return
	}

	// todo: авторизация и проверка на владение сущностью
	response := h.service.SafeDelete(r.Context(), id)
	if !response.Ok() {
		respond(w, r, http.StatusInternalServerError, response)
		return
	}

	respond(w, r, http.StatusOK, response)
}

func getFilter(r *http.Request) Filter {
	return Filter{
		UserId:   r.URL.Query().Get("userId"),
		Search:   r.URL.Query().Get("search"),
		FolderId: r.URL.Query().Get("folderId"),
	}
}

func respond(w http.ResponseWriter, r *http.Request, status int, response api.AppResponse) {
	render.Status(r, status)
	render.JSON(w, r, response)
}
