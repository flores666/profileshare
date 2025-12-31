package auth

import (
	"api"
	"auth/internal/lib/masking"
	"auth/internal/lib/password"
	"auth/internal/storage"
	"context"
	"eventBus"
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
	producer   eventBus.Producer
}

func NewService(repository Repository, logger *slog.Logger, producer eventBus.Producer) Service {
	return &service{
		repository: repository,
		logger:     logger,
		producer:   producer,
	}
}

func (s *service) Register(ctx context.Context, request RegisterUserRequest) (*storage.User, *api.ValidationErrors) {
	if err := validateRegister(request); err != nil {
		return nil, err
	}

	existingUser, err := s.repository.GetUser(ctx, request.Email)
	if err != nil {
		s.logger.Error("failed to get user", slog.String("error", err.Error()))
		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	if existingUser != nil {
		return s.handleExistingUser(ctx, existingUser)
	}

	return s.createUser(ctx, request)
}

func (s *service) handleExistingUser(ctx context.Context, user *storage.User) (*storage.User, *api.ValidationErrors) {
	if user.IsConfirmed {
		return nil, api.NewValidationErrors(ErrAlreadyRegistered)
	}

	if user.CodeRequestedAt.Add(CodeRequestTimeout).After(time.Now()) {
		return nil, api.NewValidationErrors(ErrCodeRequestTimeout)
	}

	user.Code = masking.RandStringBytesMask(10)
	user.CodeRequestedAt = time.Now()

	if err := s.repository.UpdateCode(ctx, user.Id, user.Code, user.CodeRequestedAt); err != nil {
		s.logger.Error("could not update user code", slog.String("error", err.Error()))
		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	s.publishUserRegistered(ctx, user)

	return user, nil
}

func (s *service) createUser(ctx context.Context, request RegisterUserRequest) (*storage.User, *api.ValidationErrors) {
	now := time.Now()
	id := utils.NewGuid()

	model := &storage.User{
		Id:              id,
		Nickname:        request.Nickname,
		Email:           request.Email,
		PasswordHash:    password.Hash(request.Password),
		Code:            masking.RandStringBytesMask(10),
		CodeRequestedAt: now,
		CreatedAt:       now,
	}

	if err := s.repository.CreateUser(ctx, model); err != nil {
		s.logger.Error("could not create user", slog.String("error", err.Error()))

		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	s.publishUserRegistered(ctx, model)

	model.PasswordHash = ""
	model.Code = ""

	return model, nil
}

func (s *service) publishUserRegistered(ctx context.Context, user *storage.User) {
	event := &UserRegisteredEvent{
		UserId: user.Id,
		Email:  user.Email,
		Code:   user.Code,
	}

	if err := s.producer.Produce(ctx, UserCreatedTopic, event); err != nil {
		s.logger.Error("failed to produce event", slog.String("error", err.Error()))
	}
}
