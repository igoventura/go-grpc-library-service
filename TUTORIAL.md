# Building Your First Go gRPC Service: A Complete Journey 🚀

*From zero to production-ready microservice in Go with database integration*

---

## 🎯 What You'll Build

Welcome to the world of Go and gRPC! In this hands-on guide, we'll build a complete **Library Management Service** that demonstrates the core principles of modern microservice development. By the end, you'll have:

- ✅ A fully functional gRPC API with CRUD operations
- ✅ CockroachDB integration with repository pattern
- ✅ Clean architecture with domain models and DTOs
- ✅ Professional error handling with gRPC status codes
- ✅ Comprehensive test suite with mock repositories
- ✅ Production-ready server with database connections
- ✅ Environment-based configuration management

**Why this matters:** This isn't just a tutorial—it's a blueprint for building scalable, maintainable microservices with real database persistence that you'll encounter in production Go development.

---

## 📋 Table of Contents

1. [API-First Design with Protocol Buffers](#step-1-api-first-design)
2. [Clean Architecture with Repository Pattern](#step-2-clean-architecture)
3. [Database Integration with CockroachDB](#step-3-database-integration)
4. [Professional Error Handling](#step-4-error-handling)
5. [Test-Driven Development with Mocks](#step-5-comprehensive-testing)
6. [Production-Ready Server](#step-6-production-server)
7. [Key Takeaways](#chapter-summary)

---

## Step 1: API-First Design with Protocol Buffers 📝

### The Philosophy

Every successful gRPC service begins with a **contract-first approach**. Your `.proto` file isn't just documentation—it's the single source of truth that defines your API's interface, data structures, and behavior.

### What We Built

Our `library_service.proto` defines a complete CRUD API:

```proto
service LibraryService {
    rpc CreateBook(CreateBookRequest) returns (Book);
    rpc GetBook(GetBookRequest) returns (Book);
    rpc UpdateBook(UpdateBookRequest) returns (Book);
    rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty);
    rpc ListBooks(ListBooksRequest) returns (ListBooksResponse);
}
```

And our `book_model.proto` defines the data structure:

```proto
message Book {
    string id = 1;
    string title = 2;
    string author = 3;
    int32 edition = 4;
    string isbn = 5;
}
```

### Pro Tips You Learned

- **Code Generation**: Protocol Buffers automatically generate type-safe Go structs and interfaces
- **Version Compatibility**: Proto3 syntax ensures forward/backward compatibility
- **Language Agnostic**: The same `.proto` file can generate clients in Python, Java, Node.js, etc.
- **Separation of Concerns**: Separate model and service definitions for better organization

### Real-World Impact

This approach scales beautifully—teams can work on different services simultaneously, knowing exactly what to expect from each API.
---

## Step 2: Clean Architecture with Repository Pattern 🏗️

### The Challenge

As applications grow, tight coupling between business logic and data storage becomes a nightmare. Changes to the database affect business logic, making testing difficult and maintenance expensive.

### The Solution: Repository Pattern

We implemented a clean architecture with clear separation of concerns:
```go
// Domain model - pure business entity
type Book struct {
    ID        string    `db:"id"`
    Title     string    `db:"title"`
    Author    string    `db:"author"`
    Edition   int       `db:"edition"`
    ISBN      string    `db:"isbn"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// Repository interface - defines data access contract
type BookRepository interface {
    CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
    GetBookByID(ctx context.Context, id string) (*domain.Book, error)
    UpdateBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
    DeleteBook(ctx context.Context, id string) error
    ListBooks(ctx context.Context) ([]*domain.Book, error)
}
```

### DTO Pattern for API Conversion

We created conversion functions to transform between domain models and API DTOs:
```go
func BookToDto(book *Book) *v1.Book {
    return &v1.Book{
        Id:      book.ID,
        Title:   book.Title,
        Author:  book.Author,
        Edition: int32(book.Edition),
        Isbn:    book.ISBN,
    }
}
```

### Architecture Benefits

- **Testability**: Easy to mock repositories for unit testing
- **Flexibility**: Can swap database implementations without changing business logic
- **Maintainability**: Clear separation makes code easier to understand and modify
- **Scalability**: Each layer can be optimized independently
---

## Step 3: Database Integration with CockroachDB 🗄️

### From In-Memory to Production Database

We evolved from a simple in-memory map to a robust CockroachDB implementation:

### Database Schema Design

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

### Repository Implementation
```go
func (r *BookRepository) CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    stmt := `INSERT INTO books (title, author, edition, isbn)
             VALUES ($1, $2, $3, $4)
             RETURNING id, created_at, updated_at`

    row := tx.QueryRowContext(ctx, stmt, book.Title, book.Author, book.Edition, book.ISBN)
    err = row.Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)

    if err != nil {
        return nil, err
    }

    return book, tx.Commit()
}
```

### Production Database Features

- **ACID Transactions**: Ensure data consistency
- **UUID Primary Keys**: Distributed-system friendly identifiers
- **Timestamps**: Track creation and modification times
- **Connection Pooling**: Efficient resource management
- **Context Propagation**: Proper cancellation and timeout handling

### Environment Configuration
```go
// Load environment variables
connStr := os.Getenv("DATABASE_URL")
if connStr == "" {
    log.Fatal("DATABASE_URL environment variable is not set")
}

db, err := sql.Open("postgres", connStr)
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
```

---

## Step 4: Professional Error Handling 🎯

### Beyond Basic Errors

Generic Go errors (`fmt.Errorf`) don't cut it in distributed systems. We implemented a sophisticated error handling strategy that maps database errors to appropriate gRPC status codes.

### Repository-Level Error Handling

```go
var ErrNotFound = errors.New("not found")

func (r *BookRepository) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
    // ... database query ...
    err := row.Scan(&book.ID, &book.Title, /* ... */)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, repository.ErrNotFound  // Convert DB error to domain error
        }
        return nil, err
    }
    return book, nil
}
```

### Service-Level Error Translation

```go
func (s *LibraryServiceServerImpl) GetBook(ctx context.Context, req *v1.GetBookRequest) (*v1.Book, error) {
    book, err := s.repo.GetBookByID(ctx, req.Id)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            return nil, status.Errorf(codes.NotFound, "book not found: %s", req.Id)
        }
        return nil, status.Errorf(codes.Internal, "failed to get book: %v", err)
    }

    return domain.BookToDto(book), nil
}
```

### Error Mapping Strategy

- `repository.ErrNotFound` → `codes.NotFound`
- Database constraint violations → `codes.InvalidArgument`
- Connection/timeout errors → `codes.Internal`
- Context cancellation → `codes.Cancelled`

This creates a **consistent error experience** across your entire microservice ecosystem.
---

## Step 5: Test-Driven Development with Mocks 🧪

### Testing Philosophy

"Code without tests is broken by design" - But with database integration, we need a smarter testing strategy than hitting a real database for every test.

### Mock Repository Pattern

We created a comprehensive mock repository for fast, reliable unit tests:

```go
type MockBookRepository struct {
    books   map[string]*domain.Book
    counter int
}

func (m *MockBookRepository) CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
    m.counter++
    book.ID = fmt.Sprintf("test-id-%d", m.counter)
    m.books[book.ID] = book
    return book, nil
}

func (m *MockBookRepository) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
    book, exists := m.books[id]
    if !exists {
        return nil, repository.ErrNotFound
    }
    return book, nil
}
```

### Testing Strategy Layers

1. **Unit Tests**: Service layer with mock repository (fast, isolated)
2. **Integration Tests**: Repository layer with test database (real DB interactions)
3. **End-to-End Tests**: Full gRPC calls with test database

### Sample Test Structure

```go
func TestLibraryServiceServerImpl_GetBook_NotFound(t *testing.T) {
    mockRepo := NewMockBookRepository()
    service := New(mockRepo)
    ctx := context.Background()

    _, err := service.GetBook(ctx, &v1.GetBookRequest{Id: "nonexistent"})
    
    // Verify it's the RIGHT kind of error
    st, ok := status.FromError(err)
    if !ok {
        t.Fatal("Expected gRPC status error")
    }
    if st.Code() != codes.NotFound {
        t.Errorf("Expected NotFound, got %v", st.Code())
    }
}
```

### Testing Benefits

- **Speed**: Mock tests run in milliseconds
- **Reliability**: No external dependencies
- **Coverage**: Easy to test error scenarios
- **Isolation**: Each test starts with clean state

---

## Step 6: Production-Ready Server 🚀

### Database Connection Management

Our production server handles database lifecycle properly:

```go
func main() {
    // Load environment configuration
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, reading from environment")
    }

    // Database connection
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        log.Fatal("DATABASE_URL environment variable is not set")
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Verify connection
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    // Dependency injection
    bookRepo := cockroach.NewBookRepository(db)
    libraryServer := server.NewLibraryServer(bookRepo)

    // gRPC server setup
    grpcServer := grpc.NewServer()
    pb.RegisterLibraryServiceServer(grpcServer, libraryServer)
    reflection.Register(grpcServer)

    // Start serving
    grpcServer.Serve(lis)
}
```

### Production Features

- **Environment Variables**: Configuration through `.env` files
- **Connection Health Checks**: Database ping on startup
- **Dependency Injection**: Clean separation of concerns
- **gRPC Reflection**: Development and debugging support
- **Graceful Resource Cleanup**: Proper `defer` usage

### Development Tools Integration

**gRPC Reflection** enables powerful debugging:
```bash
# Discover available services
grpcurl -plaintext localhost:50051 list

# Call methods with real data
grpcurl -plaintext -d '{
  "title": "Clean Architecture",
  "author": "Robert Martin",
  "edition": 1,
  "isbn": "978-0134494166"
}' localhost:50051 library.v1.LibraryService/CreateBook
```

---

## 🎯 Chapter Summary: Your Production Journey Complete

### What You've Accomplished

Congratulations! You've built a **production-ready gRPC microservice** that demonstrates:

🏗️ **Clean Architecture**: Repository pattern with domain models and DTOs
🗄️ **Database Integration**: CockroachDB with proper transaction handling
🎯 **Error Handling**: Multi-layer error mapping with proper gRPC status codes
🧪 **Testing Excellence**: Mock repositories for fast, reliable unit tests
⚡ **Production Readiness**: Environment configuration and connection management
🔧 **Developer Experience**: Reflection-enabled debugging and tooling

### The Evolution Journey

You've seen how a service evolves from simple to sophisticated:

1. **Phase 1**: In-memory storage with mutexes
2. **Phase 2**: Clean architecture with repository pattern
3. **Phase 3**: Database integration with proper error handling
4. **Phase 4**: Comprehensive testing with mocks
5. **Phase 5**: Production-ready configuration management

### Real-World Applications

This pattern applies to virtually any microservice:

- **E-commerce**: Product catalogs, order management, inventory tracking
- **Finance**: Account management, transaction processing, audit trails
- **Healthcare**: Patient records, appointment scheduling, medical history
- **IoT**: Device management, sensor data, telemetry processing
- **Social Media**: User profiles, content management, activity feeds

### Next Level Features

Ready to go beyond? Consider adding:

- **Database Migrations**: Automated schema evolution with tools like `golang-migrate`
- **Connection Pooling**: Advanced database performance with `pgxpool`
- **Caching Layer**: Redis integration for frequently accessed data
- **Observability**: Metrics with Prometheus, tracing with Jaeger
- **Security**: JWT authentication, rate limiting, input validation
- **Deployment**: Docker containers, Kubernetes manifests, CI/CD pipelines

### Resources & Community

- 📚 **This Repository**: [github.com/igoventura/go-grpc-library-service](https://github.com/igoventura/go-grpc-library-service)
- 🏗️ **Clean Architecture**: [Uncle Bob's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- 🗄️ **CockroachDB Docs**: [CockroachDB Go Tutorial](https://www.cockroachlabs.com/docs/stable/build-a-go-app-with-cockroachdb.html)
- 🧪 **Testing in Go**: [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- 💬 **Discussion**: Open GitHub issues for questions or improvements

---

### Final Thoughts

You've just built something remarkable. This isn't just a library service—it's a **template for enterprise-grade microservices**. The patterns, practices, and architecture you've learned here will serve you well as you build larger, more complex distributed systems.

The Go community values simplicity, reliability, and performance. Your service embodies all three. You've written code that's not just functional, but maintainable, testable, and scalable.

*Welcome to the world of production Go development. You're ready to build the future! 🎉*

---
*Happy coding, and welcome to the Go community! 🚀*

