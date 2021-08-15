package main

import (
	"log"
	"net"
	"os"
	"todos/internal/server"
	"todos/internal/storage"
	"todos/proto"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("TODO_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	inMemory := storage.NewInMemory()
	proto.RegisterTodoServiceServer(grpcServer, server.NewServer(inMemory))

	grpcServer.Serve(lis)
}
