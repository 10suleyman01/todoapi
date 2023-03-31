package todo

import "context"

type Repository interface {
	GetAllTodoByUserId(ctx context.Context, userId string) (t []Todo, err error)
	GetTodoById(ctx context.Context, id string) (todo Todo, err error)
	Create(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, todo Todo) error
	Delete(ctx context.Context, id string, userId string) error
}
