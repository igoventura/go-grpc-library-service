package cockroach

import (
	"context"
	"database/sql"

	"github.com/igoventura/go-grpc-library-service/internal/domain"
	"github.com/igoventura/go-grpc-library-service/internal/repository"
)

type BookRepository struct {
	repository.BookRepository

	db *sql.DB
}

func NewBookRepository(db *sql.DB) repository.BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) CreateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmt := `INSERT INTO books (title, author, edition, isbn) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	row := tx.QueryRowContext(ctx, stmt, book.Title, book.Author, book.Edition, book.ISBN)

	err = row.Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return book, nil
}
