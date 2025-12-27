package entity

import (
	"database/sql"
	"time"
)

// Content represents content entity in database
type Content struct {
	Id          string
	UserId      string
	DisplayName string
	Text        string
	MediaUrl    string
	Type        string
	FolderId    string
	CreatedAt   time.Time
	DeletedAt   sql.NullTime
}
