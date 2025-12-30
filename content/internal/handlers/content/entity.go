package content

import (
	"time"
)

// Content represents content entity in database
type Content struct {
	Id          string    `db:"id"`
	UserId      string    `db:"user_id"`
	DisplayName string    `db:"display_name"`
	Text        string    `db:"text"`
	MediaUrl    string    `db:"media_url"`
	Type        string    `db:"type"`
	FolderId    string    `db:"folder_id"`
	CreatedAt   time.Time `db:"created_at"`
	DeletedAt   time.Time `db:"deleted_at"`
}

type UpdateContent struct {
	Id          string  `db:"id"`
	DisplayName *string `db:"display_name"`
	Text        *string `db:"text"`
	MediaUrl    *string `db:"media_url"`
}
