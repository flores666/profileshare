package users

import (
	"auth/internal/storage"
	"context"
	"github.com/flores666/profileshare-lib/api"
	"log/slog"
)

type Service interface {
	GetById(ctx context.Context, id string) (*storage.User, *api.ValidationErrors)
	GetByFilter(ctx context.Context, filter QueryFilter) ([]*storage.User, *api.ValidationErrors)
	Update(ctx context.Context, request UpdateUserRequest) *api.ValidationErrors
}

const (
	ErrFailedSave  = "failed to save data"
	ErrFailedQuery = "failed to query data"
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

func (s *service) GetById(ctx context.Context, id string) (*storage.User, *api.ValidationErrors) {
	if err := validateId(id); err != nil {
		return nil, err
	}

	model, err := s.repository.GetById(ctx, id)
	if err != nil {
		s.logger.Error("could not get user, error = ", err, "id = ", id)
		return nil, api.NewValidationErrors(ErrFailedQuery)
	}

	return model, nil
}

func (s *service) GetByFilter(ctx context.Context, filter QueryFilter) ([]*storage.User, *api.ValidationErrors) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}

	list, err := s.repository.Query(ctx, filter)

	if err != nil {
		s.logger.Error("could not get users by filter, error = ", err, "filter = ", filter)
		return nil, api.NewValidationErrors(ErrFailedQuery)
	}

	return list, nil
}

func (s *service) Update(ctx context.Context, request UpdateUserRequest) *api.ValidationErrors {
	if err := validateUpdate(request); err != nil {
		return err
	}

	model := storage.UpdateUser{
		Id:       request.Id,
		Nickname: request.Nickname,
		Email:    request.Email,
	}

	if err := s.repository.Update(ctx, model); err != nil {
		s.logger.Error("could not update user, error = ", err, "id = ", request.Id)
		return api.NewValidationErrors(ErrFailedSave)
	}

	return nil
}
