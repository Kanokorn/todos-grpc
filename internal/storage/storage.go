package storage

import (
	"context"
	"fmt"

	"github.com/Kanokorn/todos-grpc/internal/todos"
)

type TodoRepository interface {
	Add(ctx context.Context, todo *todos.Todo) (*todos.Todo, error)
	ChangeStatus(ctx context.Context, id string) (*todos.Todo, error)
	List(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error)
	Remove(ctx context.Context, id string) error
}

type ErrNotFound struct {
	ID string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("not found todo id: %s", e.ID)
}

type ErrUpdate struct {
	ID string
}

func (e *ErrUpdate) Error() string {
	return fmt.Sprintf("failed to update todo id: %s", e.ID)
}
