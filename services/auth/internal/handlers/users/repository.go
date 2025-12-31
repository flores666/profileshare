package users

import (
	"auth/internal/storage"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*storage.User, error)
	Query(ctx context.Context, filter QueryFilter) ([]*storage.User, error)
	Update(ctx context.Context, model storage.UpdateUser) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetById(ctx context.Context, id string) (*storage.User, error) {
	query := `
		SELECT
			id,
			nickname,
			email,
			COALESCE(role_id, '00000000-0000-0000-0000-000000000000') AS role_id,
			COALESCE(code_requested_at, make_timestamptz(1,1,1,0,0,0)) AS code_requested_at,
			is_confirmed,
			COALESCE(banned_before, make_timestamptz(1,1,1,0,0,0)) AS banned_before,
			created_at
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

func (r *repository) Query(ctx context.Context, filter QueryFilter) ([]*storage.User, error) {
	query := `
		SELECT
			id,
			nickname,
			email,
			COALESCE(role_id, '00000000-0000-0000-0000-000000000000') AS role_id,
			is_confirmed,
			COALESCE(banned_before, make_timestamptz(1,1,1,0,0,0)) AS banned_before,
			created_at
		FROM authorization_service.users
	`

	params := map[string]any{}

	if filter.Search != "" {
		query += `
			WHERE nickname ILIKE :search
			   OR email ILIKE :search
		`
		params["search"] = "%" + filter.Search + "%"
	}

	query += " ORDER BY created_at DESC LIMIT 20"

	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, err
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var users []*storage.User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) Update(ctx context.Context, model storage.UpdateUser) error {
	query := "UPDATE users.users SET "
	params := map[string]any{
		"id": model.Id,
	}

	var sets []string

	if model.Nickname != nil {
		sets = append(sets, "nickname = :nickname")
		params["nickname"] = *model.Nickname
	}

	if model.Email != nil {
		sets = append(sets, "email = :email")
		params["email"] = *model.Email
	}

	if len(sets) == 0 {
		return errors.New("nothing to update")
	}

	query += strings.Join(sets, ", ")
	query += " WHERE id = :id"

	_, err := r.db.NamedExecContext(ctx, query, params)
	return err
}
