package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"

	v1 "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
)

type LibraryServiceServerImpl struct {
	v1.UnimplementedLibraryServiceServer

	mu    sync.RWMutex
	books map[string]*v1.Book
}

func New() *LibraryServiceServerImpl {
	return &LibraryServiceServerImpl{
		books: make(map[string]*v1.Book),
	}
}

// CreateBook handles the RPC call to create a new book.
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
