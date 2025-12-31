package auth

import (
	"auth/internal/storage"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(ctx context.Context, user *storage.User) error
	GetUser(ctx context.Context, email string) (*storage.User, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) CreateUser(ctx context.Context, user *storage.User) error {
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

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r repository) GetUser(ctx context.Context, email string) (*storage.User, error) {
	query := `SELECT id, nickname, email, password_hash, code_requested_at FROM authorization_service.users WHERE LOWER(email) = LOWER(:email)`

	var user *storage.User
	err := r.db.GetContext(ctx, user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}
