package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewStorage creates a connection pool to a PostgreSQL database.
// Make sure to close the connection by calling db.Close() when done.
func NewStorage(driverName, connectionString string) (*sql.DB, error) {
	const op = "storage.postgresql.NewStorage"

	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var version string
	if err := db.QueryRow("SELECT version()").Scan(&version); err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return db, nil
}
