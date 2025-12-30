package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UsersRepository interface {
	Create(ctx context.Context, user User) error
	GetById(ctx context.Context, id string) (*User, error)
	Query(ctx context.Context, filter Filter) ([]*User, error)
	Update(ctx context.Context, model UpdateUser) error
}

type usersRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) UsersRepository {
	return usersRepository{db: db}
}

func (r usersRepository) Create(ctx context.Context, user User) error {
	//TODO implement me
	panic("implement me")
}

func (r usersRepository) GetById(ctx context.Context, id string) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (r usersRepository) Query(ctx context.Context, filter Filter) ([]*User, error) {
	//TODO implement me
	panic("implement me")
}

func (r usersRepository) Update(ctx context.Context, model UpdateUser) error {
	//TODO implement me
	panic("implement me")
}
