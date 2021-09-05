package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kanokorn/todos-grpc/internal/server"
	"github.com/Kanokorn/todos-grpc/internal/storage"
	"github.com/Kanokorn/todos-grpc/proto"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
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

	db, err := sqlx.Connect("postgres", "postgres://postgres:password@127.0.0.1:5432/todos?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	inPG := storage.NewPostgres(db)
	proto.RegisterTodoServiceServer(grpcServer, server.NewServer(inPG))

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
