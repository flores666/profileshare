package content

import (
	"content/internal/lib/api"
)

type CreateRequest struct {
	UserId      string `json:"userId" validate:"required"`
	DisplayName string `json:"displayName" validate:"required"`
	Text        string `json:"text,omitempty"`
	MediaUrl    string `json:"mediaUrl" validate:"required"`
	Type        string `json:"type" validate:"required"`
}

type Response struct {
	api.HttpResponse
	Id string `json:"id"`
}
