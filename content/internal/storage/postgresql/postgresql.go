package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(driverName, connectionString string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var version string
	if err := db.QueryRow("SELECT version()").Scan(&version); err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
