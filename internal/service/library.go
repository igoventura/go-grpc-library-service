package service

import (
	"context"
	"errors"

	"github.com/igoventura/go-grpc-library-service/internal/domain"
	"github.com/igoventura/go-grpc-library-service/internal/repository"
	v1 "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type LibraryServiceServerImpl struct {
	v1.UnimplementedLibraryServiceServer

	repo repository.BookRepository
}

func New(bookRepo repository.BookRepository) *LibraryServiceServerImpl {
	return &LibraryServiceServerImpl{
		repo: bookRepo,
	}
}

func (s *LibraryServiceServerImpl) CreateBook(ctx context.Context, req *v1.CreateBookRequest) (*v1.Book, error) {
	domainBook := &domain.Book{
		Title:   req.Title,
		Author:  req.Author,
		Edition: int(req.Edition),
		ISBN:    req.Isbn,
	}

	createdBook, err := s.repo.CreateBook(ctx, domainBook)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create book: %v", err)
	}

	responseDto := domain.BookToDto(createdBook)

	return responseDto, nil
}

func (s *LibraryServiceServerImpl) GetBook(ctx context.Context, req *v1.GetBookRequest) (*v1.Book, error) {
	book, err := s.repo.GetBookByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "book not found: %s", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to get book: %v", err)
	}

	responseDto := domain.BookToDto(book)

	return responseDto, nil
}

func (s *LibraryServiceServerImpl) UpdateBook(ctx context.Context, req *v1.UpdateBookRequest) (*v1.Book, error) {
	domainBook := &domain.Book{
		ID:      req.Id,
		Title:   req.Title,
		Author:  req.Author,
		Edition: int(req.Edition),
		ISBN:    req.Isbn,
	}

	updatedBook, err := s.repo.UpdateBook(ctx, domainBook)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "book not found: %s", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to update book: %v", err)
	}

	responseDto := domain.BookToDto(updatedBook)
	return responseDto, nil
}

func (s *LibraryServiceServerImpl) DeleteBook(ctx context.Context, req *v1.DeleteBookRequest) (*emptypb.Empty, error) {
	err := s.repo.DeleteBook(ctx, req.Id)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "book not found: %s", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete book: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *LibraryServiceServerImpl) ListBooks(ctx context.Context, req *v1.ListBooksRequest) (*v1.ListBooksResponse, error) {
	response := &v1.ListBooksResponse{}
	books, err := s.repo.ListBooks(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list books: %v", err)
	}

	for _, book := range books {
		if book.ISBN == "" || book.Title == "" {
			continue
		}
		bookDto := domain.BookToDto(book)
		response.Books = append(response.Books, bookDto)
	}

	return response, nil
}
