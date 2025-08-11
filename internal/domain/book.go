package domain

import (
	"time"

	v1 "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
)

type Book struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Author    string    `db:"author"`
	Edition   int       `db:"edition"`
	ISBN      string    `db:"isbn"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func BookToDto(book *Book) *v1.Book {
	return &v1.Book{
		Id:      book.ID,
		Title:   book.Title,
		Author:  book.Author,
		Edition: int32(book.Edition),
		Isbn:    book.ISBN,
	}
}
