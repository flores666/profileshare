package postgresql

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// NewStorage creates a connection pool to a PostgreSQL database.
// Make sure to close the connection by calling db.Close() when done.
func NewStorage(driverName, connectionString string) (*sqlx.DB, error) {
	const op = "storage.postgresql.NewStorage"

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conStr := stdlib.RegisterConnConfig(config)

	db, err := sqlx.Open(driverName, conStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var version string
	if err = db.QueryRow("SELECT version()").Scan(&version); err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return db, nil
}
