package server

import (
	"context"
	"todos"
	"todos/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoServer struct {
	TodoManager todos.TodoManager

	proto.UnimplementedTodoServiceServer
}

func NewServer(todoManager todos.TodoManager) *TodoServer {
	return &TodoServer{
		TodoManager: todoManager,
	}
}

func (s *TodoServer) Add(ctx context.Context, r *proto.AddRequest) (*proto.Todo, error) {
	todo, err := s.TodoManager.Add(ctx, &todos.Todo{
		Label: r.GetLabel(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Todo{
		Id:        todo.ID,
		Label:     todo.Label,
		Completed: todo.Completed,
	}, nil
}

func (s *TodoServer) ChangeStatus(ctx context.Context, r *proto.ChangeStatusRequest) (*proto.Todo, error) {
	todo, err := s.TodoManager.ChangeStatus(ctx, r.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.Todo{
		Id:        todo.ID,
		Label:     todo.Label,
		Completed: todo.Completed,
	}, nil
}

func (s *TodoServer) ListAll(ctx context.Context, r *proto.ListAllRequest) (*proto.Todos, error) {
	todos, err := s.TodoManager.List(ctx, todos.ListOption(r.Option))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var ptodos []*proto.Todo

	for _, todo := range todos {
		ptodos = append(ptodos, &proto.Todo{
			Id:        todo.ID,
			Label:     todo.Label,
			Completed: todo.Completed,
		})
	}

	return &proto.Todos{
		Todos: ptodos,
	}, nil
}

func (s *TodoServer) Remove(ctx context.Context, r *proto.RemoveRequest) (*proto.RemoveResponse, error) {
	err := s.TodoManager.Remove(ctx, r.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.RemoveResponse{}, nil
}
