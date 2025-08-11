package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	server "github.com/igoventura/go-grpc-library-service/internal/server"
	pb "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
)

func main() {
	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Starting gRPC server on port 50051...")

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create and register the library server
	libraryServer := server.NewLibraryServer()
	pb.RegisterLibraryServiceServer(grpcServer, libraryServer)

	// Enable gRPC reflection for debugging tools like grpcurl
	reflection.Register(grpcServer)

	log.Println("Library gRPC service registered")
	log.Println("Server is ready to accept connections...")

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
