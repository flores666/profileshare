package repository

import (
	"auth/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type TokensRepository interface {
	SaveToken(ctx context.Context, token *storage.Token) error
	GetByRefresh(ctx context.Context, token string) (*storage.Token, error)
	Revoke(ctx context.Context, token string) error
}

type tokensRepository struct {
	db *sqlx.DB
}

func NewTokensRepository(db *sqlx.DB) TokensRepository {
	return &tokensRepository{db: db}
}

func (t *tokensRepository) SaveToken(ctx context.Context, token *storage.Token) error {
	executor := getExecutor(ctx, t.db)

	columns := []string{
		"id",
		"user_id",
		"provider_name",
		"token",
		"expires_at",
		"created_at",
	}

	values := []string{
		":id",
		":user_id",
		":provider_name",
		":token",
		":expires_at",
		":created_at",
	}

	if token.ReplacedByToken != "" {
		columns = append(columns, "replaced_by_token")
		values = append(values, ":replaced_by_token")
	}

	if token.RevokedByIp != "" {
		columns = append(columns, "revoked_by_ip")
		values = append(values, ":revoked_by_ip")
	}

	if !token.RevokedAt.IsZero() {
		columns = append(columns, "revoked_at")
		values = append(values, ":revoked_at")
	}

	query := fmt.Sprintf(`
		INSERT INTO authorization_service.tokens (%s)
		VALUES (%s)
	`, strings.Join(columns, ", "), strings.Join(values, ", "))

	_, err := executor.NamedExecContext(ctx, query, token)
	return err
}

func (t *tokensRepository) GetByRefresh(ctx context.Context, token string) (*storage.Token, error) {
	query := `SELECT * FROM authorization_service.tokens WHERE token = $1`

	var item storage.Token
	err := t.db.GetContext(ctx, &item, query, token)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &item, err
}

func (t *tokensRepository) Revoke(ctx context.Context, token string) error {
	query := `
		UPDATE authorization_service.tokens
		SET
			revoked_at = $1,
			expires_at = $2,
		WHERE token = $3 AND revoked_at IS NULL
	`

	now := time.Now().UTC()

	_, err := t.db.ExecContext(ctx, query, token, now, now)
	return err
}
