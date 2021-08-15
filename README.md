# Todo gRPC

Todo APIs with grpc.

Features:
- Add todo
- Change todo status (completed/incompleted)
- ListAll todos and you can filter only "completed" or "incompleted" status
- Remove todo

## Installation

This project is written in Go. If this is your first time encountering Go, follow instructions to install it on your computer.

After installing Go, run the following command to download the project:

```bash
git clone https://github.com/Kanokorn/todos-grpc.git
```

## Configuration

Configuration required enviroment variables.

```bash
set -x TODO_GRPC_PORT ${TODO_GRPC_PORT}
```

## Run

Run the following command:

```bash
make run
```

It will run gRPC server on port 50051.

If your machine has Docker installed, you can build and run application with:

```bash
make docker-build
make docker-run
```

## Test

Tests are written in each package directory. To run all packages tests:

```bash
make test
```

## Other make file commands can be seen by

```bash
make help
```
