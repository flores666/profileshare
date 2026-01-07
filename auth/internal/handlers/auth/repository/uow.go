package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork interface {
	Users() UsersRepository
	Tokens() TokensRepository
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type sqlTxKey struct{}

type unitOfWork struct {
	db               *sqlx.DB
	usersRepository  UsersRepository
	tokensRepository TokensRepository
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &unitOfWork{
		db:               db,
		usersRepository:  NewUsersRepository(db),
		tokensRepository: NewTokensRepository(db),
	}
}

func (u *unitOfWork) Users() UsersRepository {
	if u.usersRepository == nil {
		u.usersRepository = NewUsersRepository(u.db)
	}

	return u.usersRepository
}

func (u *unitOfWork) Tokens() TokensRepository {
	if u.tokensRepository == nil {
		u.tokensRepository = NewTokensRepository(u.db)
	}

	return u.tokensRepository
}

func (u *unitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	ctx = context.WithValue(ctx, sqlTxKey{}, tx)
	if err = fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func getExecutor(ctx context.Context, db *sqlx.DB) Executor {
	if tx, ok := ctx.Value(sqlTxKey{}).(*sqlx.Tx); ok {
		return tx
	}

	return db
}
