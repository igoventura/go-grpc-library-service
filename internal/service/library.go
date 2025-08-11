package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/igoventura/go-grpc-library-service/internal/repository"
	v1 "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type LibraryServiceServerImpl struct {
	v1.UnimplementedLibraryServiceServer

	mu    sync.RWMutex
	books map[string]*v1.Book
}

func New(bookRepo repository.BookRepository) *LibraryServiceServerImpl {
	return &LibraryServiceServerImpl{
		books: make(map[string]*v1.Book),
	}
}

func (s *LibraryServiceServerImpl) CreateBook(ctx context.Context, req *v1.CreateBookRequest) (*v1.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id := hex.EncodeToString(idBytes)

	book := &v1.Book{
		Id:      id,
		Title:   req.Title,
		Author:  req.Author,
		Edition: req.Edition,
		Isbn:    req.Isbn,
	}

	s.books[id] = book

	return book, nil
}

func (s *LibraryServiceServerImpl) GetBook(ctx context.Context, req *v1.GetBookRequest) (*v1.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[req.Id]

	if !ok {
		err := status.Errorf(codes.NotFound, "book not found: %s", req.Id)
		return nil, err
	}

	return book, nil
}

func (s *LibraryServiceServerImpl) UpdateBook(ctx context.Context, req *v1.UpdateBookRequest) (*v1.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[req.Id]

	if !ok {
		err := status.Errorf(codes.NotFound, "book not found: %s", req.Id)
		return nil, err
	}

	book.Title = req.Title
	book.Author = req.Author
	book.Edition = req.Edition
	book.Isbn = req.Isbn

	s.books[req.Id] = book
	return book, nil
}

func (s *LibraryServiceServerImpl) DeleteBook(ctx context.Context, req *v1.DeleteBookRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.books[req.Id]
	if !ok {
		err := status.Errorf(codes.NotFound, "book not found: %s", req.Id)
		return nil, err
	}

	delete(s.books, req.Id)

	return &emptypb.Empty{}, nil
}

func (s *LibraryServiceServerImpl) ListBooks(ctx context.Context, req *v1.ListBooksRequest) (*v1.ListBooksResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	response := &v1.ListBooksResponse{}
	for _, book := range s.books {
		if book.Isbn == "" || book.Title == "" {
			continue
		}
		response.Books = append(response.Books, book)
	}

	return response, nil
}
