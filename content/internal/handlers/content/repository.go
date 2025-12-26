package content

import (
	"content/internal/handlers/content/entity"
	"database/sql"
)

type Repository interface {
	Create(content entity.Content) error
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
