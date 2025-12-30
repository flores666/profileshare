package users

import (
	"api"
	"context"
	"log/slog"
)

type Service interface {
	Create(ctx context.Context, user CreateUserRequest) (*User, *api.ValidationErrors)
	GetById(ctx context.Context, id string) (*User, *api.ValidationErrors)
	Query(ctx context.Context, filter Filter) ([]*User, *api.ValidationErrors)
	Update(ctx context.Context, model UpdateUser) *api.ValidationErrors
}

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

func (s service) Create(ctx context.Context, user CreateUserRequest) (*User, *api.ValidationErrors) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetById(ctx context.Context, id string) (*User, *api.ValidationErrors) {
	//TODO implement me
	panic("implement me")
}

func (s service) Query(ctx context.Context, filter Filter) ([]*User, *api.ValidationErrors) {
	//TODO implement me
	panic("implement me")
}

func (s service) Update(ctx context.Context, model UpdateUser) *api.ValidationErrors {
	//TODO implement me
	panic("implement me")
}
