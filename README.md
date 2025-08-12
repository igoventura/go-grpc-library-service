# Go gRPC Library Service

A production-ready gRPC service built with Go for managing a library of books. This project demonstrates modern Go microservice development patterns including clean architecture, database integration, and comprehensive testing.

## ✨ Features

- 📚 Complete CRUD operations for books (Create, Read, Update, Delete, List)
- 🗄️ CockroachDB integration with repository pattern
- 🔒 Thread-safe concurrent request handling
- 🎯 Professional gRPC error handling with status codes
- ✅ Comprehensive test suite with mock repository
- 🛠️ Production-ready server with reflection enabled
- 🔧 Environment-based configuration
## 🏗️ Project Structure

```
library-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   └── book.go             # Domain models and DTOs
│   ├── repository/
│   │   ├── book.go             # Repository interface
│   │   └── cockroach/
│   │       └── book.go         # CockroachDB implementation
│   ├── server/
│   │   └── server.go           # Server factory
│   └── service/
│       ├── library.go          # Business logic
│       └── library_test.go     # Service tests
├── migrations/
│   └── 0001_create_books_table.sql  # Database migrations
├── proto/
│   ├── book_model.proto        # Book data model
│   └── library_service.proto   # Service definitions
├── pkg/
│   └── pb/                     # Generated protobuf code
├── .env.example                # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── TUTORIAL.md                 # Complete learning guide
```

### 📁 Architecture Overview

- **`cmd/`** - Application executables and entry points
- **`internal/domain/`** - Domain models and business entities
- **`internal/repository/`** - Data access layer with interface and implementations
- **`internal/service/`** - Business logic and gRPC service handlers
- **`internal/server/`** - Server configuration and dependency injection
- **`migrations/`** - Database schema migrations
- **`proto/`** - Protocol buffer schemas (API contracts)
- **`pkg/`** - Generated code and public libraries

### 🎯 Why This Structure?

- **Clean Architecture**: Clear separation between domain, service, and data layers
- **Dependency Injection**: Easy testing with mock repositories
- **Scalability**: Supports future microservice expansion
- **Maintainability**: Follows Go community standards
- **Testability**: Each layer can be tested in isolation

## 🚀 Getting Started

### Prerequisites

1. **Go 1.21+**
2. **Protocol Buffers Compiler**
3. **CockroachDB** (local or cloud)
```bash
# Install protoc compiler
brew install protobuf  # macOS
# or apt-get install protobuf-compiler  # Ubuntu

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 📋 Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-grpc-library-service
```

2. **Install dependencies**
```bash
make install-deps
# or manually:
go mod tidy
```

3. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your database URL
```

4. **Setup database**
```bash
# Run migrations (you'll need to implement a migration runner or run manually)
# For now, execute the SQL in migrations/0001_create_books_table.sql
```

5. **Generate Protocol Buffer code**
```bash
make generate
```

## 🔧 Available Commands

### Development
```bash
make run              # Start the server
make test             # Run all tests
make test-coverage    # Run tests with coverage
make format           # Format code
make generate         # Generate protobuf code
```

### Production
```bash
make build            # Build binary
./bin/server          # Run built binary
```

### Database Operations
```bash
# Example database URL for CockroachDB:
# postgresql://username:password@localhost:26257/library?sslmode=require
```

## 🗄️ Database Schema

```sql
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title STRING NOT NULL,
    author STRING NOT NULL,
    edition INT NOT NULL,
    isbn STRING NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
```

## 🧪 Testing

The project includes comprehensive tests:

- **Unit Tests**: Service layer with mock repository
- **Integration Tests**: Can be added for database operations
- **Race Detection**: Automatic with `make test-coverage`

```bash
# Run tests
make test

# Run with race detection and coverage
make test-coverage

# Test specific package
go test ./internal/service/...
```

## 📡 API Usage

### Using grpcurl

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Create a book
grpcurl -plaintext -d '{
  "title": "The Go Programming Language",
  "author": "Alan Donovan",
  "edition": 1,
  "isbn": "978-0134190440"
}' localhost:50051 library.v1.LibraryService/CreateBook

# Get a book
grpcurl -plaintext -d '{"id": "book-uuid-here"}' \
  localhost:50051 library.v1.LibraryService/GetBook

# List all books
grpcurl -plaintext -d '{}' localhost:50051 library.v1.LibraryService/ListBooks
```

## 🌍 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | CockroachDB connection string | `postgresql://user:pass@localhost:26257/library?sslmode=require` |

## 🏆 Production Considerations

### What's Included
- ✅ Structured logging ready
- ✅ gRPC reflection for debugging
- ✅ Clean error handling
- ✅ Environment-based configuration
- ✅ Database connection management

### What You Might Add
- 🔄 Database connection pooling
- 📊 Metrics and monitoring (Prometheus)
- 🔍 Distributed tracing (Jaeger)
- 🔐 Authentication and authorization
- 🛡️ Rate limiting and circuit breakers
- 📝 Structured logging (zerolog/logrus)
- 🐳 Docker containerization
- ☸️ Kubernetes deployment manifests

## 📚 Learning Resources

- **📖 [TUTORIAL.md](TUTORIAL.md)** - Complete step-by-step learning guide
- **🎥 [gRPC Docs](https://grpc.io/docs/languages/go/)** - Official gRPC Go documentation
- **📘 [Go Standards](https://github.com/golang-standards/project-layout)** - Project layout standards

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [gRPC](https://grpc.io/) and [Protocol Buffers](https://developers.google.com/protocol-buffers)
- Database powered by [CockroachDB](https://www.cockroachlabs.com/)
- Inspired by Go community best practices

---

*Happy coding! 🚀*
