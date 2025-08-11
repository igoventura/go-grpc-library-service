# Go gRPC Library Service

A simple gRPC service built with Go for learning purposes.

## Project Structure

```
library-service/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── server/
│   │   └── server.go        # gRPC server implementation
│   ├── service/
│   │   └── library.go       # Business logic
│   └── repository/
│       └── repository.go    # Data access layer
├── proto/
│   └── library.proto        # Protocol buffer definitions
├── pkg/
│   └── pb/                  # Generated protobuf code
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Key folder purposes:

- `cmd/` - Application executables
- `internal/` - Private application code (not importable by others)
- `proto/` - Protocol buffer schemas
- `pkg/` - Public library code (importable)

### Why this structure?

- Follows Go community standards
- Clear separation of concerns
- Supports future microservice expansion

## Prerequisites

Install protobuf compiler and Go plugins:

```bash
# Install protoc compiler
brew install protobuf  # macOS
# or apt-get install protobuf-compiler  # Ubuntu

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Commands

### Generate/Update Protocol Buffer Files
```bash
make generate
# or manually:
protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb proto/*.proto
```

### Install Dependencies
```bash
go get google.golang.org/protobuf
go get google.golang.org/grpc
go mod tidy
```

### Run Server
```bash
go run cmd/server/main.go
```

### Run Tests
```bash
go test ./...
go test -race -cover ./...
```

### Build
```bash
go build -o bin/server cmd/server/main.go
```

### Format Code
```bash
go fmt ./...
```

## Development Workflow

1. Define/modify `.proto` files
2. Run `make generate` to update Go code
3. Implement service handlers
4. Run tests
5. Start server