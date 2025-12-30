package content

import (
	"content/internal/lib/api"
	"content/internal/lib/utils"
	"context"
	"log/slog"
	"time"
)

type Service interface {
	Create(ctx context.Context, request CreateContentRequest) (*Content, *api.ValidationErrors)
	Update(ctx context.Context, request UpdateContentRequest) *api.ValidationErrors
	GetById(ctx context.Context, id string) (*Content, *api.ValidationErrors)
	GetByFilter(ctx context.Context, filter Filter) ([]*Content, *api.ValidationErrors)
	SafeDelete(ctx context.Context, id string) *api.ValidationErrors
}

type service struct {
	repository Repository
	logger     *slog.Logger
}

const (
	ErrFailedSave  = "failed to save data"
	ErrFailedQuery = "failed to query data"
)

func NewService(repository Repository, logger *slog.Logger) Service {
	srv := service{
		repository: repository,
		logger:     logger,
	}

	srv.logger = srv.logger.With(slog.String("caller", "handlers.content.service"))

	return srv
}

func (s service) Create(ctx context.Context, request CreateContentRequest) (*Content, *api.ValidationErrors) {
	if err := validateCreate(request); err != nil {
		return nil, err
	}

	id := utils.NewGuid()
	now := time.Now()

	model := Content{
		Id:          id,
		UserId:      request.UserId,
		DisplayName: request.DisplayName,
		Text:        request.Text,
		MediaUrl:    request.MediaUrl,
		Type:        request.Type,
		FolderId:    request.FolderId,
		CreatedAt:   now,
	}

	if repoErr := s.repository.Create(ctx, model); repoErr != nil {
		s.logger.Error("could not create content, error = ", repoErr.Error())
		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	return &model, nil
}

func (s service) GetById(ctx context.Context, id string) (*Content, *api.ValidationErrors) {
	if err := validateId(id); err != nil {
		return nil, err
	}

	item, err := s.repository.GetById(ctx, id)
	if err != nil {
		s.logger.Error("could not get content, error = ", err, "id = ", id)
		return nil, api.NewValidationErrors(ErrFailedQuery)
	}

	return item, nil
}

func (s service) GetByFilter(ctx context.Context, filter Filter) ([]*Content, *api.ValidationErrors) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}

	list, err := s.repository.Query(ctx, filter)

	if err != nil {
		s.logger.Error("could not get content by filter, error = ", err, "filter = ", filter)
		return nil, api.NewValidationErrors(ErrFailedQuery)
	}

	return list, nil
}

func (s service) Update(ctx context.Context, request UpdateContentRequest) *api.ValidationErrors {
	if err := validateUpdate(request); err != nil {
		return err
	}

	model := UpdateContent{
		Id:          request.Id,
		DisplayName: request.DisplayName,
		Text:        request.Text,
		MediaUrl:    request.MediaUrl,
	}

	if err := s.repository.Update(ctx, model); err != nil {
		s.logger.Error("could not update content, error = ", err, "id = ", request.Id)
		return api.NewValidationErrors(ErrFailedSave)
	}

	return nil
}

func (s service) SafeDelete(ctx context.Context, id string) *api.ValidationErrors {
	err := s.repository.SafeDelete(ctx, id)
	if err != nil {
		s.logger.Error("could not safe delete content, error = ", err, "id = ", id)
		return api.NewValidationErrors(ErrFailedSave)
	}

	return nil
}
