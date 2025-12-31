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
	ErrFailedSave         = "failed to save data"
	ErrAlreadyRegistered  = "user already registered"
	ErrFailedQuery        = "failed to query data"
	ErrCodeRequestTimeout = "code request timeout"
	CodeRequestTimeout    = time.Minute * 5
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

func (s *service) Register(ctx context.Context, request RegisterUserRequest) (*storage.User, *api.ValidationErrors) {
	if err := validateRegister(request); err != nil {
		return nil, err
	}

	existingUser, err := s.repository.GetUser(ctx, request.Email)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	if existingUser != nil {
		if existingUser.IsConfirmed {
			return nil, api.NewValidationErrors(ErrAlreadyRegistered)
		}

		if existingUser.CodeRequestedAt.Add(CodeRequestTimeout).After(time.Now()) {
			return nil, api.NewValidationErrors(ErrCodeRequestTimeout)
		} else {
			// todo send email again
		}

		return existingUser, nil
	}

	now := time.Now()
	id := utils.NewGuid()

	model := &storage.User{
		Id:              id,
		Nickname:        request.Nickname,
		Email:           request.Email,
		PasswordHash:    password.Hash(request.Password),
		CodeRequestedAt: now,
		CreatedAt:       now,
	}

	//todo email confirm
	if err := s.repository.CreateUser(ctx, model); err != nil {
		s.logger.Error("could not create user", slog.String("error", err.Error()))

		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	model.PasswordHash = ""

	return model, nil
}
