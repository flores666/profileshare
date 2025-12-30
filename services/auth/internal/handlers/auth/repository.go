package auth

import (
	"auth/internal/storage"
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, user storage.User) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Create(ctx context.Context, user storage.User) error {
	query := `
		INSERT INTO authorization_service.users (
			id,
			nickname,
			email,
			password_hash,
		    code_requested_at,
			created_at
		) VALUES (
			:id,
			:nickname,
			:email,
			:password_hash,
		    :code_requested_at,
			:created_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, &user)
	return err
}
