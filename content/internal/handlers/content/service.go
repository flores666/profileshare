package content

import (
	"content/internal/handlers/content/entity"
	"content/internal/lib/api"
	"content/internal/lib/utils"
	"database/sql"
	"log/slog"
	"time"
)

type Service interface {
	Create(request CreateContentRequest) Response
	Update(request UpdateContentRequest) Response
	GetById(id string) Response
	GetByFilter(filter Filter) QueryResponse
	SafeDelete(id string) Response
}

type service struct {
	repository Repository
	logger     *slog.Logger
}

func NewService(repository Repository, logger *slog.Logger) Service {
	srv := service{
		repository: repository,
		logger:     logger,
	}

	srv.logger = srv.logger.With(slog.String("caller", "handlers.content.service"))

	return srv
}

func (s service) Create(request CreateContentRequest) Response {
	if err := validateCreate(request); err != nil {
		return Response{
			HttpResponse: api.NewError(err.Error()),
		}
	}

	id := utils.NewGuid()
	now := time.Now()

	model := entity.Content{
		Id:          id,
		UserId:      request.UserId,
		DisplayName: request.DisplayName,
		Text:        request.Text,
		MediaUrl:    request.MediaUrl,
		Type:        request.Type,
		FolderId:    request.FolderId,
		CreatedAt:   now,
		DeletedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	err := s.repository.Create(model)

	if err != nil {
		s.logger.Error("could not create content, error = ", err.Error())
		return Response{
			HttpResponse: api.NewError("could not create content"),
		}
	}

	return Response{
		Data:         &model,
		HttpResponse: api.NewOk(),
	}
}

func (s service) GetById(id string) Response {
	if id == "" {
		return Response{
			HttpResponse: api.NewError("id is required"),
		}
	}

	item, err := s.repository.GetById(id)
	if err != nil {
		return Response{
			HttpResponse: api.NewError(err.Error()),
		}
	}

	return Response{
		HttpResponse: api.NewOk(),
		Data:         item,
	}
}

func (s service) GetByFilter(filter Filter) QueryResponse {
	if err := validateFilter(filter); err != nil {
		return QueryResponse{
			HttpResponse: api.NewError(err.Error()),
		}
	}

	list, err := s.repository.Query(filter)

	if err != nil {
		return QueryResponse{
			HttpResponse: api.NewError(err.Error()),
			Data:         []*entity.Content{},
		}
	}

	return QueryResponse{
		HttpResponse: api.NewOk(),
		Data:         list,
	}
}

func (s service) Update(request UpdateContentRequest) Response {
	if err := validateUpdate(request); err != nil {
		return Response{
			HttpResponse: api.NewError(err.Error()),
		}
	}

	model := entity.UpdateContent{
		Id:          request.Id,
		DisplayName: request.DisplayName,
		Text:        request.Text,
		MediaUrl:    request.MediaUrl,
	}

	err := s.repository.Update(model)
	if err != nil {
		return Response{
			HttpResponse: api.NewError(err.Error()),
		}
	}

	return Response{
		HttpResponse: api.NewOk(),
	}
}

func (s service) SafeDelete(id string) Response {
	err := s.repository.SafeDelete(id)
	if err != nil {
		return Response{
			HttpResponse: api.NewError(err.Error()),
		}
	}

	return Response{
		HttpResponse: api.NewOk(),
	}
}
