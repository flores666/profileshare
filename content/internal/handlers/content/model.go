package content

import (
	"content/internal/handlers/content/entity"
	"content/internal/lib/api"
)

type CreateContentRequest struct {
	UserId      string `json:"userId" validate:"required"`
	DisplayName string `json:"displayName" validate:"required"`
	Text        string `json:"text,omitempty"`
	MediaUrl    string `json:"mediaUrl" validate:"required"`
	Type        string `json:"type" validate:"required"`
	FolderId    string `json:"folderId" validate:"required"`
}

type UpdateContentRequest struct {
	Id          string  `json:"id" validate:"required"`
	DisplayName *string `json:"displayName,omitempty"`
	Text        *string `json:"text,omitempty"`
	MediaUrl    *string `json:"mediaUrl,omitempty"`
}

type Response struct {
	api.HttpResponse
	Data *entity.Content `json:"data"`
}

type Filter struct {
	UserId   string
	Search   string
	FolderId string
}

type QueryResponse struct {
	api.HttpResponse
	Data []*entity.Content `json:"data"`
}
