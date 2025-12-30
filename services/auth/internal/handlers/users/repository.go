package users

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, user User) error
	GetById(ctx context.Context, id string) (*User, error)
	Query(ctx context.Context, filter Filter) ([]*User, error)
	Update(ctx context.Context, model UpdateUser) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Create(ctx context.Context, user User) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetById(ctx context.Context, id string) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Query(ctx context.Context, filter Filter) ([]*User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(ctx context.Context, model UpdateUser) error {
	//TODO implement me
	panic("implement me")
}
