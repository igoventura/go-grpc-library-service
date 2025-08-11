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

func (r *BookRepository) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
	stmt := `SELECT id, title, author, edition, isbn, created_at, updated_at FROM books WHERE id = $1`
	row := r.db.QueryRowContext(ctx, stmt, id)
	book := &domain.Book{}
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Edition, &book.ISBN, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `UPDATE books SET title = $1, author = $2, edition = $3, isbn = $4, updated_at = now() WHERE id = $5 RETURNING updated_at`
	err = tx.QueryRowContext(ctx, stmt, book.Title, book.Author, book.Edition, book.ISBN, book.ID).Scan(&book.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return book, nil
}

func (r *BookRepository) DeleteBook(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt := `DELETE FROM books WHERE id = $1`
	res, err := tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *BookRepository) ListBooks(ctx context.Context) ([]*domain.Book, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, author, edition, isbn, created_at, updated_at FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []*domain.Book
	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Edition, &book.ISBN, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
