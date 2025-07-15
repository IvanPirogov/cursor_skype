.PHONY: help build run test clean docker-up docker-down migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start Docker containers"
	@echo "  docker-down - Stop Docker containers"
	@echo "  migrate     - Run database migrations"
	@echo "  mod         - Download Go modules"
	@echo "  fmt         - Format Go code"
	@echo "  lint        - Run linter"

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf uploads/
	go clean

# Start Docker containers
docker-up:
	docker-compose up -d

# Stop Docker containers
docker-down:
	docker-compose down

# Build and run with Docker
docker-build:
	docker-compose build

# Run database migrations
migrate:
	go run cmd/server/main.go migrate

# Download Go modules
mod:
	go mod download
	go mod tidy

# Format Go code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install development dependencies
install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o ./docs

# Setup development environment
setup-dev:
	cp .env.example .env
	mkdir -p uploads
	make docker-up
	sleep 5
	make migrate

# Development server with hot reload
dev:
	go install github.com/cosmtrek/air@latest
	air

# Production build
prod-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go