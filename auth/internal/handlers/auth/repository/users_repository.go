package repository

import (
	"auth/internal/storage"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *storage.User) error
	GetUserByEmail(ctx context.Context, email string) (*storage.User, error)
	GetUserById(ctx context.Context, id string) (*storage.User, error)
	Update(ctx context.Context, userId string, code string, codeRequestedAt time.Time, isConfirmed bool) error
}

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) CreateUser(ctx context.Context, user *storage.User) error {
	query := `
		INSERT INTO authorization_service.users (
			id,
			nickname,
			email,
			password_hash,
		    code,
		    code_requested_at,
			created_at
		) VALUES (
			:id,
			:nickname,
			:email,
			:password_hash,
		    :code,
		    :code_requested_at,
			:created_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *usersRepository) GetUserByEmail(ctx context.Context, email string) (*storage.User, error) {
	query := `
		SELECT id, 
		       nickname, 
		       email, 
		       password_hash, 
		       is_confirmed,
		       code, 
		       COALESCE(code_requested_at, make_timestamptz(1,1,1,0,0,0)) AS code_requested_at
		FROM authorization_service.users
		WHERE LOWER(email) = LOWER($1)
	`

	var user storage.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *usersRepository) GetUserById(ctx context.Context, id string) (*storage.User, error) {
	query := `
		SELECT id, 
		       nickname, 
		       email, 
		       password_hash, 
		       is_confirmed,
		       code, 
		       COALESCE(code_requested_at, make_timestamptz(1,1,1,0,0,0)) AS code_requested_at
		FROM authorization_service.users
		WHERE id = $1
	`

	var user storage.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *usersRepository) Update(ctx context.Context, userId string, code string, codeRequestedAt time.Time, isConfirmed bool) error {
	executor := getExecutor(ctx, r.db)

	query := `UPDATE authorization_service.users SET code = $1, code_requested_at = $2, is_confirmed = $3 WHERE id = $4`

	updateTime := &codeRequestedAt
	if codeRequestedAt.IsZero() {
		updateTime = nil
	}

	_, err := executor.ExecContext(ctx, query, code, updateTime, isConfirmed, userId)
	return err
}
