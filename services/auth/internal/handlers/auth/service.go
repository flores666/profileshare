package auth

import (
	"api"
	"auth/internal/lib/password"
	"auth/internal/storage"
	"context"
	"log/slog"
	"time"
	"utils"
)

type Service interface {
	Register(ctx context.Context, request RegisterUserRequest) (*storage.User, *api.ValidationErrors)
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

func (s service) Register(ctx context.Context, request RegisterUserRequest) (*storage.User, *api.ValidationErrors) {
	if err := validateRegister(request); err != nil {
		return nil, err
	}

	// todo check already registered with email and nickname

	now := time.Now()
	id := utils.NewGuid()

	model := storage.User{
		Id:              id,
		Nickname:        request.Nickname,
		Email:           request.Email,
		PasswordHash:    password.Hash(request.Password),
		CodeRequestedAt: now,
		CreatedAt:       now,
	}

	//todo email confirm
	if err := s.repository.Create(ctx, model); err != nil {
		s.logger.Error("could not create user", slog.String("error", err.Error()))

		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	model.PasswordHash = ""

	return &model, nil
}
