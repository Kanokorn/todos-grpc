package storage

import (
	"context"

	"github.com/Kanokorn/todos-grpc/internal/todos"

	"github.com/google/uuid"
)

type InMemory struct {
	storage map[string]*todos.Todo
}

func NewInMemory() *InMemory {
	return &InMemory{
		storage: map[string]*todos.Todo{},
	}
}

func (im *InMemory) Add(ctx context.Context, todo *todos.Todo) (*todos.Todo, error) {
	todo.ID = uuid.NewString()
	im.storage[todo.ID] = todo
	return todo, nil
}

func (im *InMemory) ChangeStatus(ctx context.Context, id string) (*todos.Todo, error) {
	todo := im.storage[id]
	todo.Completed = !todo.Completed
	return todo, nil
}

func (im *InMemory) List(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error) {
	var result []*todos.Todo

	switch option {
	case todos.All:
		for _, todo := range im.storage {
			result = append(result, todo)
		}
	case todos.Completed:
		for _, todo := range im.storage {
			if todo.Completed {
				result = append(result, todo)
			}
		}
	case todos.Incompleted:
		for _, todo := range im.storage {
			if !todo.Completed {
				result = append(result, todo)
			}
		}
	}

	return result, nil
}

func (im *InMemory) Remove(ctx context.Context, id string) error {
	delete(im.storage, id)
	return nil
}
