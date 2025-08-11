package server

import (
	"github.com/igoventura/go-grpc-library-service/internal/repository/cockroach"
	"github.com/igoventura/go-grpc-library-service/internal/service"
	v1 "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
)

// NewLibraryServer creates a new instance of the library service server
func NewLibraryServer() v1.LibraryServiceServer {
	return service.New(cockroach.NewBookRepository(nil))
}
