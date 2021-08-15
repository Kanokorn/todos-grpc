FROM golang:1.16-buster as builder

RUN apt update && apt -y --no-install-recommends install \
    ca-certificates \
    git

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todos ./cmd/grpc

FROM builder as tester
WORKDIR /app
RUN go test -v ./...

FROM scratch as todos
ENV TODO_GRPC_PORT=50051

WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/todos .

EXPOSE ${TODO_GRPC_PORT}
CMD ["/todos"]
