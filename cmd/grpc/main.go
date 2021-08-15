package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"todos/internal/server"
	"todos/internal/storage"
	"todos/proto"

	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := os.Getenv("TODO_GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	inMemory := storage.NewInMemory()
	proto.RegisterTodoServiceServer(grpcServer, server.NewServer(inMemory))

	go func() {
		ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		<-ctx.Done()
		grpcServer.GracefulStop()
		log.Println("shutting down server...")
	}()

	log.Printf("server started at port: %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	log.Println("server exited properly")
	return nil
}
