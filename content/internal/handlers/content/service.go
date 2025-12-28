package content

import (
	"content/internal/lib/api"
	"content/internal/lib/utils"
	"database/sql"
	"log/slog"
	"time"
)

type Service interface {
	Create(request CreateContentRequest) (*Content, *api.ValidationErrors)
	Update(request UpdateContentRequest) *api.ValidationErrors
	GetById(id string) (*Content, *api.ValidationErrors)
	GetByFilter(filter Filter) ([]*Content, *api.ValidationErrors)
	SafeDelete(id string) *api.ValidationErrors
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

func (s service) Create(request CreateContentRequest) (*Content, *api.ValidationErrors) {
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
		DeletedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	if repoErr := s.repository.Create(model); repoErr != nil {
		s.logger.Error("could not create content, error = ", repoErr.Error())
		return nil, api.NewValidationErrors(ErrFailedSave)
	}

	return &model, nil
}

func (s service) GetById(id string) (*Content, *api.ValidationErrors) {
	if err := validateId(id); err != nil {
		return nil, err
	}

	item, err := s.repository.GetById(id)
	if err != nil {
		return nil, api.NewValidationErrors(ErrFailedQuery)
	}

	return item, nil
}

func (s service) GetByFilter(filter Filter) ([]*Content, *api.ValidationErrors) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}

	list, err := s.repository.Query(filter)

	if err != nil {
		return nil, api.NewValidationErrors(ErrFailedQuery)
	}

	return list, nil
}

func (s service) Update(request UpdateContentRequest) *api.ValidationErrors {
	if err := validateUpdate(request); err != nil {
		return err
	}

	model := UpdateContent{
		Id:          request.Id,
		DisplayName: request.DisplayName,
		Text:        request.Text,
		MediaUrl:    request.MediaUrl,
	}

	err := s.repository.Update(model)
	if err != nil {
		return api.NewValidationErrors(ErrFailedSave)
	}

	return nil
}

func (s service) SafeDelete(id string) *api.ValidationErrors {
	err := s.repository.SafeDelete(id)
	if err != nil {
		return api.NewValidationErrors(ErrFailedSave)
	}

	return nil
}
