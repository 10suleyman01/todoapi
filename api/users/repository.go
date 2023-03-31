package users

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByNameAndPassword(ctx context.Context, name, password string) (User, error)
	GetById(ctx context.Context, id string) (User, error)
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}
