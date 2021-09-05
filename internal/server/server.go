package server

import (
	"context"
	"errors"

	"github.com/Kanokorn/todos-grpc/internal/storage"
	"github.com/Kanokorn/todos-grpc/internal/todos"
	"github.com/Kanokorn/todos-grpc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoServer struct {
	Storage storage.Service

	proto.UnimplementedTodoServiceServer
}

func NewServer(s storage.Service) *TodoServer {
	return &TodoServer{
		Storage: s,
	}
}

func (s *TodoServer) Add(ctx context.Context, r *proto.AddRequest) (*proto.Todo, error) {
	todo, err := s.Storage.Add(ctx, &todos.Todo{
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
	todo, err := s.Storage.ChangeStatus(ctx, r.Id)
	if err != nil {
		switch err := err.(type) {
		case *storage.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case *storage.ErrUpdate:
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var errNotFound *storage.ErrNotFound
	var errUpdate *storage.ErrUpdate

	if errors.As(err, &errNotFound) {
		return nil, status.Error(codes.NotFound, errNotFound.Error())
	} else if errors.As(err, &errUpdate) {
		return nil, status.Error(codes.Internal, errUpdate.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.Todo{
		Id:        todo.ID,
		Label:     todo.Label,
		Completed: todo.Completed,
	}, nil
}

func (s *TodoServer) ListAll(ctx context.Context, r *proto.ListAllRequest) (*proto.Todos, error) {
	todos, err := s.Storage.List(ctx, todos.ListOption(r.Option))
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
	err := s.Storage.Remove(ctx, r.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.RemoveResponse{}, nil
}
