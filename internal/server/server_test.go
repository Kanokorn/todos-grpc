package server_test

import (
	"context"
	"fmt"
	"testing"
	"todos"
	"todos/internal/server"
	"todos/proto"

	"github.com/stretchr/testify/require"
)

type mockTodoManager struct {
	AddFn          func(ctx context.Context, todo *todos.Todo) (*todos.Todo, error)
	ChangeStatusFn func(ctx context.Context, id string) (*todos.Todo, error)
	ListFn         func(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error)
	RemoveFn       func(ctx context.Context, id string) error
}

func (m *mockTodoManager) Add(ctx context.Context, todo *todos.Todo) (*todos.Todo, error) {
	return m.AddFn(ctx, todo)
}

func (m *mockTodoManager) ChangeStatus(ctx context.Context, id string) (*todos.Todo, error) {
	return m.ChangeStatusFn(ctx, id)
}

func (m *mockTodoManager) List(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error) {
	return m.ListFn(ctx, option)
}

func (m *mockTodoManager) Remove(ctx context.Context, id string) error {
	return m.RemoveFn(ctx, id)
}

func TestAddTodoSuccess(t *testing.T) {
	mock := &mockTodoManager{}
	mock.AddFn = func(ctx context.Context, todo *todos.Todo) (*todos.Todo, error) {
		todo.ID = "1"
		return todo, nil
	}

	srv := server.NewServer(mock)
	resp, err := srv.Add(context.Background(), &proto.AddRequest{
		Label: "Hello",
	})

	require.NoError(t, err)
	require.Equal(t, "1", resp.Id)
	require.Equal(t, "Hello", resp.Label)
	require.False(t, resp.Completed)
}

func TestAddTodoError(t *testing.T) {
	mock := &mockTodoManager{}
	mock.AddFn = func(ctx context.Context, todo *todos.Todo) (*todos.Todo, error) {
		return nil, fmt.Errorf("boom!")
	}

	srv := server.NewServer(mock)
	_, err := srv.Add(context.Background(), &proto.AddRequest{
		Label: "Hello",
	})

	require.Error(t, err)
}

func TestChangeStatusSuccess(t *testing.T) {
	mock := &mockTodoManager{}
	mock.ChangeStatusFn = func(ctx context.Context, id string) (*todos.Todo, error) {
		return &todos.Todo{
			ID:        id,
			Label:     "Hello",
			Completed: true,
		}, nil
	}

	srv := server.NewServer(mock)
	todo, err := srv.ChangeStatus(context.Background(), &proto.ChangeStatusRequest{
		Id: "1",
	})

	require.NoError(t, err)
	require.Equal(t, "1", todo.Id)
	require.Equal(t, "Hello", todo.Label)
	require.True(t, todo.Completed)
}

func TestChangeStatusError(t *testing.T) {
	mock := &mockTodoManager{}
	mock.ChangeStatusFn = func(ctx context.Context, id string) (*todos.Todo, error) {
		return nil, fmt.Errorf("boom!")
	}

	srv := server.NewServer(mock)
	_, err := srv.ChangeStatus(context.Background(), &proto.ChangeStatusRequest{
		Id: "1",
	})

	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = boom!", err.Error())
}

func TestListAllTodoSuccess(t *testing.T) {
	mock := &mockTodoManager{}
	mock.ListFn = func(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error) {
		var result []*todos.Todo

		t1 := &todos.Todo{
			ID:        "1",
			Label:     "Hello",
			Completed: false,
		}

		t2 := &todos.Todo{
			ID:        "2",
			Label:     "World",
			Completed: true,
		}

		result = append(result, t1, t2)
		return result, nil
	}

	srv := server.NewServer(mock)
	todos, err := srv.ListAll(context.Background(), &proto.ListAllRequest{
		Option: proto.ListAllRequest_ALL,
	})

	require.NoError(t, err)

	first := todos.Todos[0]
	require.Equal(t, "1", first.Id)
	require.Equal(t, "Hello", first.Label)
	require.False(t, first.Completed)

	second := todos.Todos[1]
	require.Equal(t, "2", second.Id)
	require.Equal(t, "World", second.Label)
	require.True(t, second.Completed)
}

func TestListAllTodoError(t *testing.T) {
	mock := &mockTodoManager{}
	mock.ListFn = func(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error) {
		return nil, fmt.Errorf("boom!")
	}

	srv := server.NewServer(mock)
	_, err := srv.ListAll(context.Background(), &proto.ListAllRequest{
		Option: proto.ListAllRequest_ALL,
	})

	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = boom!", err.Error())
}

func TestListAllCompletedTodo(t *testing.T) {
	mock := &mockTodoManager{}
	mock.ListFn = func(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error) {
		var result []*todos.Todo

		t1 := &todos.Todo{
			ID:        "1",
			Label:     "Hello",
			Completed: true,
		}
		result = append(result, t1)

		return result, nil
	}

	srv := server.NewServer(mock)
	todos, err := srv.ListAll(context.Background(), &proto.ListAllRequest{
		Option: proto.ListAllRequest_COMPLETED,
	})

	require.NoError(t, err)

	first := todos.Todos[0]
	require.Equal(t, "1", first.Id)
	require.Equal(t, "Hello", first.Label)
	require.True(t, first.Completed)
}

func TestListAllInCompletedTodo(t *testing.T) {
	mock := &mockTodoManager{}
	mock.ListFn = func(ctx context.Context, option todos.ListOption) ([]*todos.Todo, error) {
		var result []*todos.Todo

		t1 := &todos.Todo{
			ID:        "1",
			Label:     "Hello",
			Completed: false,
		}
		result = append(result, t1)

		return result, nil
	}

	srv := server.NewServer(mock)
	todos, err := srv.ListAll(context.Background(), &proto.ListAllRequest{
		Option: proto.ListAllRequest_INCOMPLETED,
	})

	require.NoError(t, err)

	first := todos.Todos[0]
	require.Equal(t, "1", first.Id)
	require.Equal(t, "Hello", first.Label)
	require.False(t, first.Completed)
}

func TestRemoveTodoSuccess(t *testing.T) {
	mock := &mockTodoManager{}
	mock.RemoveFn = func(ctx context.Context, id string) error {
		return nil
	}

	srv := server.NewServer(mock)
	_, err := srv.Remove(context.Background(), &proto.RemoveRequest{
		Id: "1",
	})

	require.NoError(t, err)
}

func TestRemoveTodoError(t *testing.T) {
	mock := &mockTodoManager{}
	mock.RemoveFn = func(ctx context.Context, id string) error {
		return fmt.Errorf("boom!")
	}

	srv := server.NewServer(mock)
	_, err := srv.Remove(context.Background(), &proto.RemoveRequest{
		Id: "1",
	})

	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = boom!", err.Error())
}
