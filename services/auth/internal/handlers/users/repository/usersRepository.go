package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UsersRepository interface {
	Create(ctx context.Context, user User) error
	GetById(ctx context.Context, id string) (*User, error)
	Query(ctx context.Context, filter Filter) ([]*User, error)
	Update(ctx context.Context, model UpdateUser) error
}

type usersRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) UsersRepository {
	return usersRepository{db: db}
}

func (r usersRepository) Create(ctx context.Context, user User) error {
	query := `
		INSERT INTO users.users (
			id,
			nickname,
			email,
			role_id,
			created_at
		) VALUES (
			:id,
			:nickname,
			:email,
			:role_id,
			:created_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r usersRepository) GetById(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT
			id,
			nickname,
			email,
			role_id,
			COALESCE(banned_before, '0001-01-01 00:00:00+00') AS banned_before,
			created_at
		FROM users.users
		WHERE id = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r usersRepository) Query(ctx context.Context, filter Filter) ([]*User, error) {
	query := `
		SELECT
			id,
			nickname,
			email,
			role_id,
			COALESCE(banned_before, '0001-01-01 00:00:00+00') AS banned_before,
			created_at
		FROM users.users
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

	var users []*User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r usersRepository) Update(ctx context.Context, model UpdateUser) error {
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

	if model.RoleId != nil {
		sets = append(sets, "role_id = :role_id")
		params["role_id"] = *model.RoleId
	}

	if !model.BannedBefore.IsZero() {
		sets = append(sets, "banned_before = :banned_before")
		params["banned_before"] = model.BannedBefore
	}

	if len(sets) == 0 {
		return errors.New("nothing to update")
	}

	query += strings.Join(sets, ", ")
	query += " WHERE id = :id"

	_, err := r.db.NamedExecContext(ctx, query, params)
	return err
}
