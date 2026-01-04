package auth

import (
	"auth/internal/lib/mapper"
	"auth/internal/lib/masking"
	"auth/internal/lib/password"
	"auth/internal/storage"
	"context"
	"log/slog"
	"time"

	"github.com/flores666/profileshare-lib/api"
	"github.com/flores666/profileshare-lib/eventBus"
	"github.com/flores666/profileshare-lib/utils"
)

type Service interface {
	Register(ctx context.Context, request RegisterUserRequest) api.AppResponse
}

const (
	ErrFailedSave         = "Не удалось сохранить данные"
	ErrAlreadyRegistered  = "Пользователь уже зарегистрирован"
	ErrCodeRequestTimeout = "Повторите попытку через 5 минут"
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

func (s *service) Register(ctx context.Context, request RegisterUserRequest) api.AppResponse {
	if err := validateRegister(request); err != nil {
		return api.NewError("Ошибка проверки данных", err)
	}

	existingUser, err := s.repository.GetUser(ctx, request.Email)
	if err != nil {
		s.logger.Error("failed to get user", slog.String("error", err.Error()))
		return api.NewError("Внутрення ошибка", nil)
	}

	if existingUser != nil {
		return s.handleExistingUser(ctx, existingUser)
	}

	return s.createUser(ctx, request)
}

func (s *service) handleExistingUser(ctx context.Context, user *storage.User) api.AppResponse {
	if user.IsConfirmed {
		return api.NewError(ErrAlreadyRegistered, nil)
	}

	if user.CodeRequestedAt.Add(CodeRequestTimeout).After(time.Now()) {
		return api.NewError(ErrCodeRequestTimeout, nil)
	}

	user.Code = masking.RandStringBytesMask(10)
	user.CodeRequestedAt = time.Now()

	if err := s.repository.UpdateCode(ctx, user.Id, user.Code, user.CodeRequestedAt); err != nil {
		s.logger.Error("could not update user code", slog.String("error", err.Error()))
		return api.NewError(ErrFailedSave, nil)
	}

	s.publishUserRegistered(user)

	return api.NewOk("Сообщение с новым кодом подтверждения отправлено на вашу почту", mapper.MapUserToDto(user))
}

func (s *service) createUser(ctx context.Context, request RegisterUserRequest) api.AppResponse {
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
		return api.NewError(ErrFailedSave, nil)
	}

	go s.publishUserRegistered(model)

	return api.NewOk("Успешно", mapper.MapUserToDto(model))
}

func (s *service) publishUserRegistered(user *storage.User) {
	event := &UserRegisteredMessage{
		UserId:         user.Id,
		Email:          user.Email,
		Code:           user.Code,
		IdempotencyKey: user.Id + ";" + user.Code,
	}

	if err := s.producer.Produce(context.Background(), UserCreatedTopic, event); err != nil {
		s.logger.Error("failed to produce event", slog.String("error", err.Error()))
	}
}
