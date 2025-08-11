package service

import (
	"context"
	"testing"

	v1 "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLibraryServiceServerImpl_CreateBook(t *testing.T) {
	// Create a new service instance for testing
	service := New()
	ctx := context.Background()

	// Test creating a book
	req := &v1.CreateBookRequest{
		Title:   "The Go Programming Language",
		Author:  "Alan Donovan",
		Edition: 1,
		Isbn:    "978-0134190440",
	}

	book, err := service.CreateBook(ctx, req)
	if err != nil {
		t.Fatalf("CreateBook failed: %v", err)
	}

	// Verify the book was created correctly
	if book.Title != req.Title {
		t.Errorf("Expected title %q, got %q", req.Title, book.Title)
	}
	if book.Author != req.Author {
		t.Errorf("Expected author %q, got %q", req.Author, book.Author)
	}
	if book.Edition != req.Edition {
		t.Errorf("Expected edition %q, got %q", req.Edition, book.Edition)
	}
	if book.Isbn != req.Isbn {
		t.Errorf("Expected ISBN %q, got %q", req.Isbn, book.Isbn)
	}
	if book.Id == "" {
		t.Error("Expected book ID to be generated, got empty string")
	}
}

func TestLibraryServiceServerImpl_GetBook(t *testing.T) {
	service := New()
	ctx := context.Background()

	// First create a book
	createReq := &v1.CreateBookRequest{
		Title:   "Clean Code",
		Author:  "Robert Martin",
		Edition: 1,
		Isbn:    "978-0132350884",
	}
	createdBook, err := service.CreateBook(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateBook failed: %v", err)
	}

	// Test getting the book
	getReq := &v1.GetBookRequest{Id: createdBook.Id}
	book, err := service.GetBook(ctx, getReq)
	if err != nil {
		t.Fatalf("GetBook failed: %v", err)
	}

	if book.Id != createdBook.Id {
		t.Errorf("Expected ID %q, got %q", createdBook.Id, book.Id)
	}
	if book.Title != createdBook.Title {
		t.Errorf("Expected title %q, got %q", createdBook.Title, book.Title)
	}
}

func TestLibraryServiceServerImpl_GetBook_NotFound(t *testing.T) {
	service := New()
	ctx := context.Background()

	// Test getting a non-existent book
	getReq := &v1.GetBookRequest{Id: "non-existent-id"}
	_, err := service.GetBook(ctx, getReq)

	if err == nil {
		t.Fatal("Expected error for non-existent book, got nil")
	}

	// Check if it's the correct gRPC error
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}
	if st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound error code, got %v", st.Code())
	}
}

func TestLibraryServiceServerImpl_UpdateBook(t *testing.T) {
	service := New()
	ctx := context.Background()

	// First create a book
	createReq := &v1.CreateBookRequest{
		Title:   "Original Title",
		Author:  "Original Author",
		Edition: 1,
		Isbn:    "978-0000000000",
	}
	createdBook, err := service.CreateBook(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateBook failed: %v", err)
	}

	// Update the book
	updateReq := &v1.UpdateBookRequest{
		Id:      createdBook.Id,
		Title:   "Updated Title",
		Author:  "Updated Author",
		Edition: 2,
		Isbn:    "978-1111111111",
	}
	updatedBook, err := service.UpdateBook(ctx, updateReq)
	if err != nil {
		t.Fatalf("UpdateBook failed: %v", err)
	}

	if updatedBook.Title != updateReq.Title {
		t.Errorf("Expected title %q, got %q", updateReq.Title, updatedBook.Title)
	}
	if updatedBook.Author != updateReq.Author {
		t.Errorf("Expected author %q, got %q", updateReq.Author, updatedBook.Author)
	}
}

func TestLibraryServiceServerImpl_UpdateBook_NotFound(t *testing.T) {
	service := New()
	ctx := context.Background()

	// Test updating a non-existent book
	updateReq := &v1.UpdateBookRequest{
		Id:      "non-existent-id",
		Title:   "Some Title",
		Author:  "Some Author",
		Edition: 1,
		Isbn:    "978-0000000000",
	}
	_, err := service.UpdateBook(ctx, updateReq)

	if err == nil {
		t.Fatal("Expected error for non-existent book, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}
	if st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound error code, got %v", st.Code())
	}
}

func TestLibraryServiceServerImpl_DeleteBook(t *testing.T) {
	service := New()
	ctx := context.Background()

	// First create a book
	createReq := &v1.CreateBookRequest{
		Title:   "Book to Delete",
		Author:  "Some Author",
		Edition: 1,
		Isbn:    "978-0000000000",
	}
	createdBook, err := service.CreateBook(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateBook failed: %v", err)
	}

	// Delete the book
	deleteReq := &v1.DeleteBookRequest{Id: createdBook.Id}
	_, err = service.DeleteBook(ctx, deleteReq)
	if err != nil {
		t.Fatalf("DeleteBook failed: %v", err)
	}

	// Verify the book is deleted
	getReq := &v1.GetBookRequest{Id: createdBook.Id}
	_, err = service.GetBook(ctx, getReq)
	if err == nil {
		t.Fatal("Expected error when getting deleted book, got nil")
	}
}

func TestLibraryServiceServerImpl_ListBooks(t *testing.T) {
	service := New()
	ctx := context.Background()

	// Create some test books
	books := []*v1.CreateBookRequest{
		{Title: "Book 1", Author: "Author 1", Edition: 1, Isbn: "978-0000000001"},
		{Title: "Book 2", Author: "Author 2", Edition: 1, Isbn: "978-0000000002"},
		{Title: "", Author: "Author 3", Edition: 1, Isbn: ""}, // This should be filtered out
	}

	for _, book := range books {
		_, err := service.CreateBook(ctx, book)
		if err != nil {
			t.Fatalf("CreateBook failed: %v", err)
		}
	}

	// List books
	listReq := &v1.ListBooksRequest{}
	response, err := service.ListBooks(ctx, listReq)
	if err != nil {
		t.Fatalf("ListBooks failed: %v", err)
	}

	// Should only return 2 books (the third one has empty title and ISBN)
	if len(response.Books) != 2 {
		t.Errorf("Expected 2 books, got %d", len(response.Books))
	}
}
