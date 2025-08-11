# Continue.dev Context - Go gRPC Library Service

## 1. CRITICAL RESPONSE RULES - READ FIRST
**MINIMAL TOKENS**: Keep all responses short and concise.
**NO UNSOLICITED CODE**: Only provide code when explicitly asked or for autocomplete.
**LEARNING FOCUSED**: Guide understanding, don't solve problems for the user.
**BE BRIEF**: Use bullet points, short sentences, and essential info only.

---

## 2. Project Summary & Overview

This document summarizes the development of a simple library management gRPC service built in Go. The project served as a hands-on exercise covering the end-to-end lifecycle of a microservice, from API design with Protocol Buffers to a fully tested, runnable server.

The final service provides basic CRUD (Create, Read, Update, Delete) and list functionality for a collection of books, managed in-memory. The primary goal was to help a developer learn Go while building a practical application.

---

## 3. Guiding Principles & Learning Objectives

This project was guided by the following principles to focus on learning and best practices.

* **Learning Objectives**:
    * Master Go syntax, idioms, and conventions.
    * Understand Go's type system, interfaces, and concurrency.
    * Learn gRPC service development in Go.
    * Apply Go best practices for project structure and error handling.

* **Common Go Pitfalls to Avoid**:
    * Nil pointer dereferences.
    * Race conditions (which we solved with mutexes).
    * Goroutine leaks and improper resource cleanup (`defer` is key).
    * Forgetting to handle `context` propagation.

* **Testing Philosophy**:
    * Test both the "happy path" and specific error cases.
    * Use `go test -race` to automatically detect race conditions.
    * For future work, consider **table-driven tests** to make test cases more concise and easier to expand.

---

## 4. Key Go Concepts Applied

This project provided practical experience with several core Go features:

* **Concurrency (`sync.RWMutex`)**: We used a Read-Write mutex to safely handle concurrent requests to our in-memory book map.
    * **Write Locks (`Lock`/`Unlock`)**: Applied to `CreateBook`, `UpdateBook`, and `DeleteBook` for exclusive access.
    * **Read Locks (`RLock`/`RUnlock`)**: Applied to `GetBook` and `ListBooks` for improved performance by allowing concurrent reads.

* **Pointers vs. Values**: We discussed how modifying a struct retrieved from a map of pointers (`map[string]*Book`) directly alters the stored value.

* **Built-in Testing (`testing` package)**: We created a comprehensive test suite in `library_test.go` to validate our service's logic.

* **Error Handling**: We moved from generic Go errors to structured gRPC errors.
    * **`status.Errorf`**: Used to return specific gRPC error codes to the client (e.g., `codes.NotFound`).
    * **`status.FromError`**: Used in tests to inspect the gRPC status code of a returned error.

* **Packages & Project Structure**: We followed Go's standard project layout conventions:
    * `cmd/server/main.go`: The executable entry point for running the server.
    * `internal/service`: The location for our core business logic.
    * `pkg/pb`: The output directory for our generated Protobuf Go code.
    * `proto`: The source for our `.proto` API definitions.

---

## 5. gRPC Development Workflow Followed

We successfully executed the following development workflow:

1.  **Define Protobuf Schema**: Created `book_model.proto` and `library_service.proto`.
2.  **Generate Go Code**: Used `protoc` with the `go` and `go-grpc` plugins.
3.  **Implement Service Interface**: Created a `Server` struct in `internal/service/library.go`.
4.  **Write Tests**: Wrote unit tests to validate all logic and error handling.
5.  **Create Server Executable**: Wrote the `main.go` file to listen for TCP connections.
6.  **Enable Reflection**: Added `reflection.Register` for easy debugging with `grpcurl`.

---

## 6. Key Code Patterns & Snippets

The following patterns were established as best practices during development.

**Pattern: Using a Read Lock for Read-Only Operations**
```go
func (s *Server) GetBook(ctx context.Context, req *v1.GetBookRequest) (*v1.Book, error) {
  s.mu.RLock()
  defer s.mu.RUnlock()

  book, ok := s.books[req.Id]
  if !ok {
    // ... return NotFound error
  }
  return book, nil
}
```

**Pattern: Returning gRPC-Specific Status Errors**
```go
if !ok {
    err := status.Errorf(codes.NotFound, "book not found: %s", id)
    return nil, err
}
```

**Pattern: Verifying gRPC Status Codes in Tests**
```go
func TestGetBook_NotFound(t *testing.T) {
    // ... setup
    _, err := service.GetBook(ctx, getReq)

    st, ok := status.FromError(err)
    if !ok {
        t.Fatal("Expected gRPC status error")
    }
    if st.Code() != codes.NotFound {
        t.Errorf("Expected NotFound error code, got %v", st.Code())
    }
}
```

---

## 7. Comprehensive Tooling & Commands

* **Dependency Management**: `go mod init`, `go mod tidy`
* **Code Quality**: `go fmt ./...`, `go vet ./...`
* **Testing**: `go test ./...`, `go test -race ./...`, `go test -cover ./...`
* **Building & Running**: `go build ./cmd/server`, `go run cmd/server/main.go`
* **Protocol Buffers**: `protoc --go_out=. --go-grpc_out=. proto/**/*.proto`
* **API Client**: `grpcurl -plaintext localhost:50051 list`
