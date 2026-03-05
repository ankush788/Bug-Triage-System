.PHONY: help build run test docker-up docker-down migrate lint fmt clean

help:
	@echo "AI Bug Triage System - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make server              Run API server"
	@echo "  make worker              Run bug analyzer worker"
	@echo "  make build               Build server and worker binaries"
	@echo "  make test                Run tests"
	@echo "  make lint                Run linter"
	@echo "  make fmt                 Format code"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up           Start services with Docker Compose"
	@echo "  make docker-down         Stop Docker Compose services"
	@echo "  make docker-build        Build Docker images"
	@echo ""
	@echo "Database:"
	@echo "  make migrate             Run database migrations"
	@echo "  make db-create           Create development database"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean               Remove build artifacts"
	@echo "  make clean-docker        Remove Docker volumes"

# Development servers
server:
	@echo "Starting API server..."
	go run ./cmd/server

worker:
	@echo "Starting bug analyzer worker..."
	go run ./cmd/worker

# Build binaries
build:
	@echo "Building server..."
	go build -o bin/server ./cmd/server
	@echo "Building worker..."
	go build -o bin/worker ./cmd/worker

# Testing and linting
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

vet:
	@echo "Running go vet..."
	go vet ./...

# Docker operations
docker-build:
	@echo "Building Docker image..."
	docker-compose build

docker-up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d
	@echo "Services started!"
	@echo "API: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "Kafka: localhost:9092"

docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-logs:
	@echo "Showing Docker Compose logs..."
	docker-compose logs -f

# Database operations
db-create:
	@echo "Creating development database..."
	createdb bug_triage 2>/dev/null || echo "Database already exists"

migrate:
	@echo "Running migrations..."
	psql bug_triage < migrations/001_initial_schema.sql

db-reset: db-drop db-create migrate
	@echo "Database reset complete"

db-drop:
	@echo "Dropping database..."
	dropdb bug_triage 2>/dev/null || echo "Database already dropped"

# Cleanup
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

clean-docker:
	@echo "Removing Docker volumes..."
	docker-compose down -v
	@echo "Docker resources cleaned"

# Dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Generate API documentation (optional)
gen-docs:
	@echo "Generating API documentation..."
	swag init -g cmd/server/main.go

# Development utilities
.env:
	@echo "Creating .env file..."
	cp .env.example .env
	@echo "Please edit .env with your configuration"

setup: .env deps db-create migrate
	@echo "Development setup complete!"
	@echo "Run 'make server' to start the API server"

run-local: setup server
