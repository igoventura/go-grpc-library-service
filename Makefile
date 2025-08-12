.PHONY: generate install-deps run test test-coverage build format clean help

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=server
BUILD_DIR=bin
PROTO_DIR=proto
PROTO_OUT=./pkg/pb

ifeq (,$(wildcard .env))
		$(warning .env file not found. Environment variables might be missing.)
else
		include .env
		export
endif

# Generate/Update Protocol Buffer Files
generate:
	@echo "Generating protocol buffer files..."
	protoc --go_out=${PROTO_OUT} --go-grpc_out=${PROTO_OUT} $(PROTO_DIR)/*.proto

# Install Dependencies
install-deps:
	@echo "Installing dependencies..."
	go get google.golang.org/protobuf
	go get google.golang.org/grpc
	go mod tidy

# Run Server
run:
	@echo "Starting server..."
	go run cmd/server/main.go

# Run Tests
test:
	@echo "Running tests..."
	go test ./...

# Run Tests with race detection and coverage
test-coverage:
	@echo "Running tests with race detection and coverage..."
	go test -race -cover ./...

# Build application
build:
	@echo "Building application..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

# Install protobuf tools
install-tools:
	@echo "Installing protobuf tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Development setup (install tools and dependencies)
setup: install-tools install-deps
	@echo "Development setup complete!"

# Full development workflow
dev: generate format test build
	@echo "Development workflow complete!"

# Database Migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path ./migrations -database ${DATABASE_URL} up

# Help target
help:
	@echo "Available targets:"
	@echo "  generate      - Generate/Update Protocol Buffer Files"
	@echo "  install-deps  - Install Go dependencies"
	@echo "  install-tools - Install protobuf compiler tools"
	@echo "  setup         - Complete development setup (tools + deps)"
	@echo "  run           - Run the server"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with race detection and coverage"
	@echo "  build         - Build the application"
	@echo "  format        - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  dev           - Full development workflow (generate + format + test + build)"
	@echo "  migrate-up    - Run database migrations"
	@echo "  help          - Show this help message"