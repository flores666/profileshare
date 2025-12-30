package users

import (
	"api"
	"auth/internal/handlers/users/repository"
	"context"
	"log/slog"
	"time"
	"utils"
)

type Service interface {
	Create(ctx context.Context, request CreateUserRequest) (*repository.User, *api.ValidationErrors)
	GetById(ctx context.Context, id string) (*repository.User, *api.ValidationErrors)
	GetByFilter(ctx context.Context, filter QueryFilter) ([]*repository.User, *api.ValidationErrors)
	Update(ctx context.Context, request UpdateUserRequest) *api.ValidationErrors
}

const (
	ErrFailedSave  = "failed to save data"
	ErrFailedQuery = "failed to query data"
)

type service struct {
	repository repository.UsersRepository
	logger     *slog.Logger
}

func NewService(repository repository.UsersRepository, logger *slog.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s service) Create(ctx context.Context, request CreateUserRequest) (*repository.User, *api.ValidationErrors) {
	if err := validateCreate(request); err != nil {
		return nil, err
	}

	id := utils.NewGuid()
	now := time.Now()

	model := repository.User{
		Id:        id,
		Nickname:  request.Nickname,
		Email:     request.Email,
		CreatedAt: now,
	}

	if repoErr := s.repository.Create(ctx, model); repoErr != nil {
		s.logger.Error("could not create content, error = ", repoErr.Error())
		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	return &model, nil
}

func (s service) GetById(ctx context.Context, id string) (*repository.User, *api.ValidationErrors) {
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

func (s service) GetByFilter(ctx context.Context, filter QueryFilter) ([]*repository.User, *api.ValidationErrors) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}

	list, err := s.repository.Query(ctx, getRepoFilter(filter))

	if err != nil {
		s.logger.Error("could not get users by filter, error = ", err, "filter = ", filter)
		return nil, api.NewValidationErrors(ErrFailedQuery)
	}

	return list, nil
}

func (s service) Update(ctx context.Context, request UpdateUserRequest) *api.ValidationErrors {
	if err := validateUpdate(request); err != nil {
		return err
	}

	model := repository.UpdateUser{
		Id:           request.Id,
		Nickname:     request.Nickname,
		Email:        request.Email,
		RoleId:       request.RoleId,
		BannedBefore: time.Time{},
	}

	if err := s.repository.Update(ctx, model); err != nil {
		s.logger.Error("could not update user, error = ", err, "id = ", request.Id)
		return api.NewValidationErrors(ErrFailedSave)
	}

	return nil
}

func getRepoFilter(filter QueryFilter) repository.Filter {
	return repository.Filter{
		Search: filter.Search,
	}
}
