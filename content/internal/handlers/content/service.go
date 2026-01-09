package content

import (
	"context"
	"log/slog"
	"time"

	"github.com/flores666/profileshare-lib/api"
	"github.com/flores666/profileshare-lib/utils"
)

type Service interface {
	Create(ctx context.Context, request CreateContentRequest, userId string) api.AppResponse
	Update(ctx context.Context, request UpdateContentRequest, userId string) api.AppResponse
	GetById(ctx context.Context, id string) api.AppResponse
	GetByFilter(ctx context.Context, filter Filter) api.AppResponse
	SafeDelete(ctx context.Context, id string, userId string) api.AppResponse
}

type service struct {
	repository Repository
	logger     *slog.Logger
}

const (
	ErrFailedSave  = "Не удалось сохранить данные"
	ErrFailedQuery = "Не удалось выполнить запрос"
	ErrValidation  = "Ошибка проверки данных"
	ErrForbidden   = "Запись вам не принадлежит"
	Success        = "Успешно"
)

func NewService(repository Repository, logger *slog.Logger) Service {
	srv := &service{
		repository: repository,
		logger:     logger,
	}

	srv.logger = srv.logger.With(slog.String("caller", "handlers.content.service"))

	return srv
}

func (s *service) Create(ctx context.Context, request CreateContentRequest, userId string) api.AppResponse {
	if err := validateCreate(request); err != nil {
		return api.NewError(ErrValidation, err)
	}

	id := utils.NewGuid()
	now := time.Now().UTC()

	model := Content{
		Id:          id,
		UserId:      userId,
		DisplayName: request.DisplayName,
		Text:        request.Text,
		MediaUrl:    request.MediaUrl,
		Type:        request.Type,
		FolderId:    request.FolderId,
		CreatedAt:   now,
	}

	if repoErr := s.repository.Create(ctx, model); repoErr != nil {
		s.logger.Error("could not create content, error = ", repoErr.Error())
		return api.NewError(ErrFailedSave, nil)
	}

	return api.NewOk(Success, model)
}

func (s *service) GetById(ctx context.Context, id string) api.AppResponse {
	if err := validateId(id); err != nil {
		return api.NewError(ErrValidation, err)
	}

	item, err := s.repository.GetById(ctx, id)
	if err != nil {
		s.logger.Error("could not get content, error = ", err, "id = ", id)
		return api.NewError(ErrFailedQuery, nil)
	}

	return api.NewOk(Success, MapContentToDto(item))
}

func (s *service) GetByFilter(ctx context.Context, filter Filter) api.AppResponse {
	if err := validateFilter(filter); err != nil {
		return api.NewError(ErrValidation, err)
	}

	list, err := s.repository.Query(ctx, filter)

	if err != nil {
		s.logger.Error("could not get content by filter, error = ", err, "filter = ", filter)
		return api.NewError(ErrFailedQuery, nil)
	}

	return api.NewOk(Success, MapContentSliceToDto(list))
}

func (s *service) Update(ctx context.Context, request UpdateContentRequest, userId string) api.AppResponse {
	if err := validateUpdate(request); err != nil {
		return api.NewError(ErrValidation, err)
	}

	content, err := s.repository.GetById(ctx, request.Id)
	if err != nil {
		s.logger.Error("could not get content, error = ", err, "id = ", request.Id)
		return api.NewError(ErrFailedQuery, nil)
	}
	if content.UserId != userId {
		return api.NewError(ErrForbidden, nil)
	}

	model := UpdateContent{
		Id:          request.Id,
		DisplayName: request.DisplayName,
		Text:        request.Text,
		MediaUrl:    request.MediaUrl,
	}

	if err := s.repository.Update(ctx, model); err != nil {
		s.logger.Error("could not update content, error = ", err, "id = ", request.Id)
		return api.NewError(ErrFailedSave, nil)
	}

	return api.NewOk(Success, MapContentToDto(content))
}

func (s *service) SafeDelete(ctx context.Context, id string, userId string) api.AppResponse {
	content, err := s.repository.GetById(ctx, id)
	if err != nil {
		s.logger.Error("could not get content, error = ", err, "id = ", id)
		return api.NewError(ErrFailedQuery, nil)
	}
	if content.UserId != userId {
		return api.NewError(ErrForbidden, nil)
	}

	err = s.repository.SafeDelete(ctx, id)
	if err != nil {
		s.logger.Error("could not safe delete content, error = ", err, "id = ", id)
		return api.NewError(ErrFailedSave, nil)
	}

	return api.NewOk(Success, nil)
}
