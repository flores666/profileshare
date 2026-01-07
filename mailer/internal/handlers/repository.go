package handlers

import (
	"context"
	"database/sql"
	"errors"
	"mailer/internal/storage"

	"github.com/jmoiron/sqlx"
)

type EmailsRepository interface {
	GetByIdempotencyKey(ctx context.Context, key string) (*storage.Email, error)
	Save(ctx context.Context, email *storage.Email) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) EmailsRepository {
	return &repository{db: db}
}

func (r *repository) GetByIdempotencyKey(ctx context.Context, key string) (*storage.Email, error) {
	query := `SELECT m.id, m.recipient, m.text, m.subject, m.status, m.idempotency_key, m.created_at FROM mailer.mails m WHERE idempotency_key = $1`

	var item storage.Email
	err := r.db.GetContext(ctx, &item, query, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &item, nil
}

func (r *repository) Save(ctx context.Context, email *storage.Email) error {
	query := `
		INSERT INTO mailer.mails (
			id,
			recipient,
			text,
			subject,
		    status,
		    idempotency_key,
			created_at
		) VALUES (
			:id,
			:recipient,
			:text,
			:subject,
		    :status,
		    :idempotency_key,
			:created_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, email)
	return err
}
