# Hunter-Seeker Job Tracker Makefile

.PHONY: help build run dev clean test docker-build docker-run sample-data deps

# Default target
help:
	@echo "Hunter-Seeker Job Tracker - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Run in development mode"
	@echo "  make run          - Build and run the application"
	@echo "  make build        - Build the application binary"
	@echo "  make deps         - Download Go dependencies"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker containers"
	@echo "  make docker-clean - Clean Docker containers and images"
	@echo ""
	@echo "Database:"
	@echo "  make sample-data  - Add sample job applications"
	@echo "  make clear-data   - Clear all job application data (with backup)"
	@echo "  make clean-db     - Remove database file"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make clean-root   - Clean accidentally created root binaries"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format Go code"
	@echo "  make lint         - Run linter"

# Go binary name
BINARY_NAME=hunter-seeker
BUILD_DIR=./bin

# Development - run with live reload if air is available
dev:
	@if command -v air >/dev/null 2>&1; then \
		echo "Running with air (live reload)..."; \
		air; \
	else \
		echo "Air not found. Running normally (install air with: go install github.com/cosmtrek/air@latest)"; \
		go run cmd/server/main.go; \
	fi

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run: build
	@echo "Starting Hunter-Seeker..."
	@mkdir -p ./data
	./$(BUILD_DIR)/$(BINARY_NAME)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Add sample data
sample-data:
	@echo "Adding sample job applications..."
	@mkdir -p ./data
	go run cmd/sample-data/main.go

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t hunter-seeker .

docker-run:
	@echo "Starting with Docker Compose..."
	docker-compose up -d
	@echo "Application running at http://localhost:8080"

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-clean: docker-stop
	@echo "Cleaning Docker containers and images..."
	docker-compose down -v
	docker rmi hunter-seeker 2>/dev/null || true

# Testing
test:
	@echo "Running tests..."
	go test -v ./...

# Code formatting
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	goimports -w . 2>/dev/null || echo "goimports not found, skipping"

# Linting
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

# Clean accidentally created root binaries
clean-root:
	@echo "Cleaning accidentally created root binaries..."
	@rm -f server sample-data debug hunter-seeker main
	@echo "Root binaries cleaned"

# Clear all data with backup
clear-data:
	@echo "Clearing job application data..."
	./scripts/clear_data.sh

# Clean database
clean-db:
	@echo "Removing database file..."
	rm -f ./data/jobs.db

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Quick start for new developers
quickstart: deps install-tools clean-root sample-data
	@echo ""
	@echo "🎯 Hunter-Seeker Quick Start Complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Run 'make dev' to start development server"
	@echo "2. Open http://localhost:8080 in your browser"
	@echo "3. Start tracking your job applications!"
	@echo ""

# Production build
build-prod:
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/server
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/server
	@echo "Production binaries built in $(BUILD_DIR)/"

# Show current status
status:
	@echo "Hunter-Seeker Status:"
	@echo "Port: $(shell cat .env 2>/dev/null | grep PORT | cut -d= -f2 || echo '8080')"
	@echo "Database: $(shell ls -la ./data/jobs.db 2>/dev/null || echo 'Not found')"
	@echo "Docker: $(shell docker-compose ps 2>/dev/null | grep hunter-seeker || echo 'Not running')"

# Backup database
backup:
	@if [ -f ./data/jobs.db ]; then \
		BACKUP_NAME="backup-$$(date +%Y%m%d-%H%M%S).db"; \
		cp ./data/jobs.db "./data/$$BACKUP_NAME"; \
		echo "Database backed up to ./data/$$BACKUP_NAME"; \
	else \
		echo "No database file found to backup"; \
	fi

# Restore database from backup
restore:
	@echo "Available backups:"
	@ls -la ./data/backup-*.db 2>/dev/null || echo "No backups found"
	@echo ""
	@echo "To restore, run: cp ./data/backup-YYYYMMDD-HHMMSS.db ./data/jobs.db"
