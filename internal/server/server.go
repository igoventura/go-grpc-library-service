package server

import (
	"github.com/igoventura/go-grpc-library-service/internal/repository"
	"github.com/igoventura/go-grpc-library-service/internal/service"
	v1 "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
)

func NewLibraryServer(bookRepo repository.BookRepository) v1.LibraryServiceServer {
	return service.New(bookRepo)
}
