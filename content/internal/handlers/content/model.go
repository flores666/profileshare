package content

import (
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

type CreateContentResponse struct {
	api.HttpResponse
	Content
}
