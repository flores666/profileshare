package content

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, content Content) error
	GetById(ctx context.Context, id string) (*Content, error)
	Query(ctx context.Context, filter Filter) ([]*Content, error)
	Update(ctx context.Context, model UpdateContent) error
	SafeDelete(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, content Content) (err error) {
	useTransaction := content.FolderId != ""

	return r.exec(ctx, useTransaction, func(exec func(query string, args ...any) (sql.Result, error)) error {
		insertContentQuery := `INSERT INTO content.content (id, user_id, display_name, text, media_url, type, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`

		_, err = exec(insertContentQuery, content.Id, content.UserId, content.DisplayName, content.Text, content.MediaUrl, content.Type, content.CreatedAt)
		if err == nil && content.FolderId != "" {
			insertLinkQuery := `INSERT INTO content.folders_contents (folder_id, content_id, created_at) VALUES ($1,$2,$3)`

			_, err = exec(insertLinkQuery,
				content.FolderId,
				content.Id,
				content.CreatedAt)
		}

		return err
	})
}

func (r *repository) GetById(ctx context.Context, id string) (*Content, error) {
	query := `SELECT         
    	id,
        user_id,
        display_name,
        text,
        media_url,
        type,
        COALESCE(deleted_at, '0001-01-01 00:00:00+00') as deleted_at,
        created_at
    FROM content.content WHERE id = $1`

	var content Content
	err := r.db.GetContext(ctx, &content, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &content, nil
}

func (r *repository) Query(ctx context.Context, filter Filter) ([]*Content, error) {
	query := `
		SELECT
			c.id, c.user_id, c.display_name, c.text,
			c.media_url, c.type, COALESCE(c.deleted_at, make_timestamptz(1,1,1,0,0,0)) as deleted_at, c.created_at
		FROM content.content c`

	params := map[string]any{
		"user_id": filter.UserId,
	}

	if filter.FolderId != "" {
		query += " JOIN content.folders_contents f on f.content_id = c.id"
	}

	query += " WHERE c.user_id = :user_id"

	if filter.FolderId != "" {
		query += " AND f.folder_id = :folder_id"
		params["folder_id"] = filter.FolderId
	}

	if filter.Search != "" {
		query += " AND (c.text ILIKE :search OR c.display_name ILIKE :search)"
		params["search"] = "%" + filter.Search + "%"
	}

	//todo: cursor pagination
	query += " ORDER BY c.created_at DESC LIMIT 20"

	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, err
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var result []*Content

	err = r.db.SelectContext(ctx, &result, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) Update(ctx context.Context, model UpdateContent) error {
	if model.Id == "" {
		return errors.New("id is required")
	}

	query := "UPDATE content.content SET "
	params := map[string]any{
		"id": model.Id,
	}

	var sets []string

	if model.DisplayName != nil {
		sets = append(sets, "display_name = :display_name")
		params["display_name"] = *model.DisplayName
	}
	if model.Text != nil {
		sets = append(sets, "text = :text")
		params["text"] = *model.Text
	}
	if model.MediaUrl != nil {
		sets = append(sets, "media_url = :media_url")
		params["media_url"] = *model.MediaUrl
	}

	if len(sets) == 0 {
		return errors.New("nothing to update")
	}

	query += strings.Join(sets, ", ")
	query += " WHERE id = :id"

	_, err := r.db.NamedExecContext(ctx, query, params)
	return err
}

func (r *repository) SafeDelete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	now := time.Now().UTC()

	query := "UPDATE content.content SET deleted_at = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, now, id)

	return err
}

func (r *repository) exec(ctx context.Context, useTransaction bool, fn func(exec func(query string, args ...any) (sql.Result, error)) error) (err error) {
	if useTransaction {
		tran, tranErr := r.db.BeginTx(ctx, &sql.TxOptions{})
		if tranErr != nil {
			return tranErr
		}

		defer func() {
			if rec := recover(); rec != nil {
				_ = tran.Rollback()
				panic(rec)
			} else if err != nil {
				_ = tran.Rollback()
			} else {
				err = tran.Commit()
			}
		}()

		err = fn(tran.Exec)
	} else {
		err = fn(r.db.Exec)
	}

	return
}
