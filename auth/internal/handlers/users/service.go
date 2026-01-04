package users

import (
	"auth/internal/lib/mapper"
	"auth/internal/storage"
	"context"
	"log/slog"

	"github.com/flores666/profileshare-lib/api"
)

type Service interface {
	GetById(ctx context.Context, id string) api.AppResponse
	GetByFilter(ctx context.Context, filter QueryFilter) api.AppResponse
	Update(ctx context.Context, request UpdateUserRequest) api.AppResponse
}

const (
	ErrFailedQuery = "Не удалось выполнить запрос"
	ErrFailedSave  = "Не удалось сохранить данные"
)

type service struct {
	repository Repository
	logger     *slog.Logger
}

func NewService(repository Repository, logger *slog.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) GetById(ctx context.Context, id string) api.AppResponse {
	if err := validateId(id); err != nil {
		return api.NewError("Ошибка проверки данных", err)
	}

	model, err := s.repository.GetById(ctx, id)
	if err != nil {
		s.logger.Error("could not get user, error = ", err, "id = ", id)
		return api.NewError(ErrFailedQuery, nil)
	}

	return api.NewOk("Успешно", mapper.MapUserToDto(model))
}

func (s *service) GetByFilter(ctx context.Context, filter QueryFilter) api.AppResponse {
	if err := validateFilter(filter); err != nil {
		return api.NewError("Ошибка проверки данных", err)
	}

	list, err := s.repository.Query(ctx, filter)

	if err != nil {
		s.logger.Error("could not get users by filter, error = ", err, "filter = ", filter)
		return api.NewError(ErrFailedQuery, nil)
	}

	return api.NewOk("Успешно", mapper.MapUserSliceToDto(list))
}

func (s *service) Update(ctx context.Context, request UpdateUserRequest) api.AppResponse {
	if err := validateUpdate(request); err != nil {
		return api.NewError("Ошибка проверки данных", err)
	}

	model := storage.UpdateUser{
		Id:       request.Id,
		Nickname: request.Nickname,
		Email:    request.Email,
	}

	if err := s.repository.Update(ctx, model); err != nil {
		s.logger.Error("could not update user, error = ", err, "id = ", request.Id)
		return api.NewError(ErrFailedSave, nil)
	}

	response := api.NewOk("Успешно", nil)

	user, err := s.repository.GetById(ctx, request.Id)
	if err != nil {
		s.logger.Error("could not get user, error = ", err, "id = ", request.Id)
		return response
	}

	response.Data = mapper.MapUserToDto(user)
	return response
}
