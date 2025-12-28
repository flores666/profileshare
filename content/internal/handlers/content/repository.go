package content

import (
	"content/internal/handlers/content/entity"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Repository interface {
	Create(content entity.Content) error
	GetById(id string) (*entity.Content, error)
	Query(filter Filter) ([]*entity.Content, error)
	Update(model entity.UpdateContent) error
	SafeDelete(id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r repository) Create(content entity.Content) (err error) {
	useTransaction := content.FolderId != ""

	return r.exec(useTransaction, func(exec func(query string, args ...any) (sql.Result, error)) error {
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

func (r repository) GetById(id string) (*entity.Content, error) {
	query := `SELECT         
    	id,
        user_id,
        display_name,
        text,
        media_url,
        type,
        deleted_at,
        created_at
    FROM content.content WHERE id = $1`

	var content entity.Content
	err := r.db.QueryRow(query, id).Scan(
		&content.Id,
		&content.UserId,
		&content.DisplayName,
		&content.Text,
		&content.MediaUrl,
		&content.Type,
		&content.DeletedAt,
		&content.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}

		return nil, err
	}

	return &content, nil
}

func (r repository) Query(filter Filter) ([]*entity.Content, error) {
	query := `
		SELECT
			c.id, c.user_id, c.display_name, c.text,
			c.media_url, c.type, c.deleted_at, c.created_at
		FROM content.content c`

	args := []any{filter.UserId}
	argID := 2

	if filter.FolderId != "" {
		query += fmt.Sprintf(" JOIN content.folders_contents f on f.content_id = c.id")
	}

	query += " WHERE c.user_id = $1"

	if filter.FolderId != "" {
		query += fmt.Sprintf(" AND f.folder_id = $%d", argID)
		args = append(args, filter.FolderId)
		argID++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (c.text ILIKE $%d OR c.display_name ILIKE $%d)", argID, argID)
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	//todo: cursor pagination
	query += " ORDER BY c.created_at DESC LIMIT 20"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]*entity.Content, 0)

	for rows.Next() {
		var content entity.Content
		if err = rows.Scan(
			&content.Id,
			&content.UserId,
			&content.DisplayName,
			&content.Text,
			&content.MediaUrl,
			&content.Type,
			&content.DeletedAt,
			&content.CreatedAt); err != nil {
			return nil, err
		}

		list = append(list, &content)
	}

	return list, nil
}

func (r repository) Update(model entity.UpdateContent) error {
	if model.Id == "" {
		return errors.New("id is required")
	}

	query := "UPDATE content.content SET "
	var args []any
	argsId := 1

	if model.DisplayName != nil {
		query += fmt.Sprintf("display_name = $%d, ", argsId)
		args = append(args, *model.DisplayName)
		argsId++
	}
	if model.Text != nil {
		query += fmt.Sprintf("text = $%d, ", argsId)
		args = append(args, *model.Text)
		argsId++
	}
	if model.MediaUrl != nil {
		query += fmt.Sprintf("media_url = $%d, ", argsId)
		args = append(args, *model.MediaUrl)
		argsId++
	}

	query = strings.TrimSuffix(query, ", ")

	query += fmt.Sprintf(" WHERE id = $%d", argsId)
	args = append(args, model.Id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r repository) SafeDelete(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	now := time.Now()

	query := "UPDATE content.content SET deleted_at = $1 WHERE id = $2"
	_, err := r.db.Exec(query, now, id)

	return err
}

func (r repository) exec(useTransaction bool, fn func(exec func(query string, args ...any) (sql.Result, error)) error) (err error) {
	if useTransaction {
		tran, tranErr := r.db.Begin()
		if tranErr != nil {
			return tranErr
		}

		defer func() {
			if rec := recover(); rec != nil {
				tran.Rollback()
				panic(rec)
			} else if err != nil {
				tran.Rollback()
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
