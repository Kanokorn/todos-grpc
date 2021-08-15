package storage

import (
	"context"
	"testing"
	"todos"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	db := NewInMemory()
	require.Len(t, db.storage, 0)

	todo, err := db.Add(context.Background(), &todos.Todo{
		Label: "Hello",
	})

	require.NoError(t, err)
	require.NotEmpty(t, todo.ID)
	require.Equal(t, "Hello", todo.Label)
	require.False(t, todo.Completed)
	require.Len(t, db.storage, 1)
}

func TestChangeStatusFromIncompletedToCompleted(t *testing.T) {
	db := NewInMemory()
	todo, err := db.Add(context.Background(), &todos.Todo{
		Label:     "Hello",
		Completed: false,
	})
	require.NoError(t, err)
	require.False(t, todo.Completed)

	result, err := db.ChangeStatus(context.Background(), todo.ID)

	require.NoError(t, err)
	require.True(t, result.Completed)
}

func TestChangeStatusFromCompletedToIncompleted(t *testing.T) {
	db := NewInMemory()
	todo, err := db.Add(context.Background(), &todos.Todo{
		Label:     "Hello",
		Completed: true,
	})
	require.NoError(t, err)
	require.True(t, todo.Completed)

	result, err := db.ChangeStatus(context.Background(), todo.ID)

	require.NoError(t, err)
	require.False(t, result.Completed)
}

func TestListAllTodos(t *testing.T) {
	db := NewInMemory()
	db.Add(context.Background(), &todos.Todo{
		Label:     "Hello",
		Completed: true,
	})
	db.Add(context.Background(), &todos.Todo{
		Label:     "World",
		Completed: false,
	})
	require.Len(t, db.storage, 2)

	result, err := db.List(context.Background(), todos.All)

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestListAllCompleted(t *testing.T) {
	db := NewInMemory()
	db.Add(context.Background(), &todos.Todo{
		Label:     "Hello",
		Completed: true,
	})
	db.Add(context.Background(), &todos.Todo{
		Label:     "World",
		Completed: false,
	})
	require.Len(t, db.storage, 2)

	result, err := db.List(context.Background(), todos.Completed)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "Hello", result[0].Label)
	require.True(t, result[0].Completed)
}

func TestListAllInCompleted(t *testing.T) {
	db := NewInMemory()
	db.Add(context.Background(), &todos.Todo{
		Label:     "Hello",
		Completed: true,
	})
	db.Add(context.Background(), &todos.Todo{
		Label:     "World",
		Completed: false,
	})
	require.Len(t, db.storage, 2)

	result, err := db.List(context.Background(), todos.Incompleted)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "World", result[0].Label)
	require.False(t, result[0].Completed)
}

func TestRemove(t *testing.T) {
	db := NewInMemory()
	todo, _ := db.Add(context.Background(), &todos.Todo{
		Label:     "Hello",
		Completed: true,
	})
	require.Len(t, db.storage, 1)

	err := db.Remove(context.Background(), todo.ID)

	require.NoError(t, err)
	require.Len(t, db.storage, 0)
}
