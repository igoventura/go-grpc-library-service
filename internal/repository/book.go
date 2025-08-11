package repository

import (
	"context"

	"github.com/igoventura/go-grpc-library-service/internal/domain"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
	GetBookByID(ctx context.Context, id string) (*domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
	DeleteBook(ctx context.Context, id string) error
	ListBooks(ctx context.Context) ([]*domain.Book, error)
}
