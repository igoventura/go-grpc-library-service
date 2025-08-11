package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/igoventura/go-grpc-library-service/internal/repository/cockroach"
	server "github.com/igoventura/go-grpc-library-service/internal/server"
	pb "github.com/igoventura/go-grpc-library-service/pkg/pb/library/v1"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, reading from environment")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close() // Ensure the connection is closed when main exits.

	// It's good practice to ping the database to verify the connection.
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Starting gRPC server on port 50051...")

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	bookRepo := cockroach.NewBookRepository(db)

	// Create and register the library server
	libraryServer := server.NewLibraryServer(bookRepo)
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
