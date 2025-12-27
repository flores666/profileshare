package content

import (
	"content/internal/handlers/content/entity"
	"content/internal/lib/api"
	"content/internal/lib/utils"
	"database/sql"
	"log/slog"
	"time"
)

type Service interface {
	CreateContent(request CreateContentRequest) CreateContentResponse
	GetContentById(id string) (*Content, error)
}

type service struct {
	repository Repository
	logger     *slog.Logger
}

type Content struct {
	entity.Content
}

func NewService(repository Repository, logger *slog.Logger) Service {
	srv := service{
		repository: repository,
		logger:     logger,
	}

	srv.logger = srv.logger.With(slog.String("caller", "handlers.content.service"))

	return srv
}

func (s service) CreateContent(request CreateContentRequest) CreateContentResponse {
	id := utils.NewGuid()
	now := time.Now()

	model := Content{
		Content: entity.Content{
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
		},
	}

	err := s.repository.Create(model.Content)

	if err != nil {
		s.logger.Error("could not create content, error = ", err.Error())
		return CreateContentResponse{
			HttpResponse: api.NewError("could not create content"),
		}
	}

	return CreateContentResponse{
		Content:      model,
		HttpResponse: api.NewOk(),
	}
}

func (s service) GetContentById(id string) (*Content, error) {
	item, err := s.repository.GetContentById(id)
	if err != nil {
		return nil, err
	}

	return &Content{
		Content: *item,
	}, nil
}
