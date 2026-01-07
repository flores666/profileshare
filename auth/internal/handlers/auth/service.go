package auth

import (
	"auth/internal/handlers/auth/repository"
	"auth/internal/handlers/auth/security"
	"auth/internal/lib/mapper"
	"auth/internal/lib/masking"
	"auth/internal/lib/password"
	"auth/internal/storage"
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/flores666/profileshare-lib/api"
	"github.com/flores666/profileshare-lib/eventBus"
	"github.com/flores666/profileshare-lib/utils"
)

type Service interface {
	Register(ctx context.Context, request RegisterUserRequest) api.AppResponse
	Confirm(ctx context.Context, request ConfirmUserRequest) api.AppResponse
	Login(ctx context.Context, request LoginUserRequest) api.AppResponse
}

const (
	ErrFailedSave         = "Не удалось сохранить данные"
	ErrAlreadyRegistered  = "Пользователь уже зарегистрирован"
	ErrCodeRequestTimeout = "Повторите попытку через 5 минут"
	ErrInternal           = "Внутрення ошибка"
	ErrInvalidCredentials = "Неверные логин или пароль"
	CodeRequestTimeout    = time.Minute * 2
	AccConfirmTimeout     = time.Minute * 10
	Success               = "Успешно"
)

type service struct {
	unitOfWork repository.UnitOfWork
	logger     *slog.Logger
	producer   eventBus.Producer
	jwtService *security.JWTService
}

func NewService(
	unitOfWork repository.UnitOfWork,
	jwtService *security.JWTService,
	logger *slog.Logger,
	producer eventBus.Producer,
) Service {
	return &service{
		logger:     logger,
		producer:   producer,
		jwtService: jwtService,
		unitOfWork: unitOfWork,
	}
}

func (s *service) Register(ctx context.Context, request RegisterUserRequest) api.AppResponse {
	if err := validateRegister(request); err != nil {
		return api.NewError("Ошибка проверки данных", err)
	}

	existingUser, err := s.unitOfWork.Users().GetUserByEmail(ctx, request.Email)
	if err != nil {
		s.logger.Error("failed to get user", slog.String("error", err.Error()))
		return api.NewError(ErrInternal, nil)
	}

	if existingUser != nil {
		return s.handleExistingUser(ctx, existingUser, request.ReturnUrl)
	}

	return s.createUser(ctx, request)
}

func (s *service) Login(ctx context.Context, request LoginUserRequest) api.AppResponse {
	if err := validateLogin(request); err != nil {
		return api.NewError("Ошибка проверки данных", err)
	}

	user, err := s.unitOfWork.Users().GetUserByEmail(ctx, request.Email)
	if err != nil {
		s.logger.Error("failed to get user", slog.String("error", err.Error()))
		return api.NewError(ErrInternal, nil)
	}

	ok, err := password.Verify(request.Password, user.PasswordHash)
	if err != nil {
		s.logger.Error("failed to verify password", slog.String("error", err.Error()))
		return api.NewError(ErrInvalidCredentials, nil)
	}

	if !ok {
		return api.NewError(ErrInvalidCredentials, nil)
	}

	tokens, err := s.issueTokens(ctx, user.Id)
	if err != nil {
		s.logger.Error("failed to issue tokens", slog.String("error", err.Error()))
		return api.NewError(ErrInternal, nil)
	}

	return api.NewOk(Success, tokens)
}

func (s *service) Confirm(ctx context.Context, request ConfirmUserRequest) api.AppResponse {
	user, err := s.unitOfWork.Users().GetUserById(ctx, request.UserId)
	if err != nil {
		s.logger.Error("failed to get user", slog.String("error", err.Error()))
		return api.NewError("Внутренняя ошибка", nil)
	}

	if user == nil {
		return api.NewError("Неверная ссылка или пользователь не найден", nil)
	}

	if user.IsConfirmed {
		return api.NewError(ErrAlreadyRegistered, nil)
	}

	if user.Code != request.Code {
		return api.NewError("Неверный код подтверждения", nil)
	}

	if user.CodeRequestedAt.Add(AccConfirmTimeout).Before(time.Now().UTC()) {
		return api.NewError("Ссылка устарела, запросите новый код подтверждения", nil)
	}

	var response api.AppResponse
	err = s.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if uowError := s.unitOfWork.Users().Update(ctx, user.Id, "", time.Time{}, true); uowError != nil {
			s.logger.Error("failed to confirm user", slog.String("error", uowError.Error()))
			return uowError
		}

		tokens, uowError := s.issueTokens(ctx, user.Id)
		if uowError != nil {
			s.logger.Error("failed to issue tokens after confirmation", slog.String("error", uowError.Error()))
			return uowError
		}

		response = api.NewOk("Аккаунт успешно подтверждён", tokens)

		return nil
	})

	if err != nil {
		s.logger.Error("failed to confirm user", slog.String("error", err.Error()))
		return api.NewError("Не удалось подтвердить пользователя", nil)
	}

	return response
}

func (s *service) issueTokens(ctx context.Context, userId string) (*security.TokenPair, error) {
	tokens, err := s.jwtService.GenerateTokens(userId)
	if err != nil {
		return nil, err
	}

	err = s.unitOfWork.Tokens().SaveToken(ctx, &storage.Token{
		Id:           utils.NewGuid(),
		UserId:       userId,
		ProviderName: security.ProviderLumo,
		Token:        tokens.RefreshToken,
		ExpiresAt:    time.Now().UTC().Add(s.jwtService.RefreshTTL),
		CreatedAt:    time.Now().UTC(),
	})

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *service) handleExistingUser(ctx context.Context, user *storage.User, redirectUrl string) api.AppResponse {
	if user.IsConfirmed {
		return api.NewError(ErrAlreadyRegistered, nil)
	}

	if user.CodeRequestedAt.Add(CodeRequestTimeout).After(time.Now().UTC()) {
		return api.NewError(ErrCodeRequestTimeout, nil)
	}

	user.Code = masking.RandStringBytesMask(10)
	user.CodeRequestedAt = time.Now().UTC()

	if err := s.unitOfWork.Users().Update(ctx, user.Id, user.Code, user.CodeRequestedAt, false); err != nil {
		s.logger.Error("could not update user code", slog.String("error", err.Error()))
		return api.NewError(ErrFailedSave, nil)
	}

	go s.publishUser(user, redirectUrl)

	return api.NewOk("Сообщение с новым кодом подтверждения отправлено на вашу почту", mapper.MapUserToDto(user))
}

func (s *service) createUser(ctx context.Context, request RegisterUserRequest) api.AppResponse {
	now := time.Now().UTC()
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

	if err := s.unitOfWork.Users().CreateUser(ctx, model); err != nil {
		s.logger.Error("could not create user", slog.String("error", err.Error()))
		return api.NewError(ErrFailedSave, nil)
	}

	go s.publishUser(model, request.ReturnUrl)

	return api.NewOk(Success, mapper.MapUserToDto(model))
}

func (s *service) publishUser(user *storage.User, redirectUrl string) {
	r, err := addQueryParam(redirectUrl, "code", user.Code)
	if err != nil {
		s.logger.Error("could not add query param", slog.String("error", err.Error()))
		r = redirectUrl
	}

	event := &UserRegisteredMessage{
		UserId:         user.Id,
		Email:          user.Email,
		ReturnUrl:      r,
		IdempotencyKey: user.Id + ";" + user.Code,
	}

	if err := s.producer.Produce(context.Background(), UserCreatedTopic, event); err != nil {
		s.logger.Error("failed to produce event", slog.String("error", err.Error()))
	}
}

func addQueryParam(rawUrl, key, value string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add(key, value)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
