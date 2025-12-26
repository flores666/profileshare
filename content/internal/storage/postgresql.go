package storage

import "database/sql"

type Storage struct {
	db *sql.DB
}

func NewStorage(connectionString string) *Storage {
	const op = "storage.postgresql.New"

	//db, err := sql.Open()
	return nil
}
