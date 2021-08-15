#!/usr/bin/env make

IMAGE := todos

.PHONY: build clean docker-build docker-run help proto run run-grpc-ui test

build: clean
	go build -a -o todos ./cmd/grpc

clean:
	rm -rf todos

docker-build:
	docker build -t $(IMAGE) --target $(IMAGE) .
	docker image prune -f
	@echo $(IMAGE)

docker-run:
	docker run --rm -p 50051:50051 $(IMAGE)

help:
	@echo "   Avaliable commands:"
	@echo "    build          Remove existing executable file, compile main package, and build a new executable file."
	@echo "    clean          Remove executable file."
	@echo "    docker-build   Build a docker image."
	@echo "    docker-run     Run a docker image as a container."
	@echo "    run            Compile and run main package."
	@echo "    run-grpc-ui    Run GRPC UI client."
	@echo "    proto          Generate protobuf."
	@echo "    test           Run all packages tests."
	@echo

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/todos.proto

run:
	go run ./cmd/grpc

run-grpc-ui:
	grpcui -import-path ./proto -proto todos.proto -plaintext 0.0.0.0:50051

test:
	go test -v -cover ./...
