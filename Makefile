.PHONY: build test run clean help backend

# Variables
BINARY_NAME=go-balancer
BACKEND_BINARY=backend-server
GO=go
GOFLAGS=-v

# Default target
all: test build

# Build the load balancer
build:
	@echo "Building load balancer..."
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) ./examples/simple

# Build backend server
backend:
	@echo "Building backend server..."
	$(GO) build $(GOFLAGS) -o bin/$(BACKEND_BINARY) ./examples/backend-server

# Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run the load balancer
run: build
	@echo "Starting load balancer..."
	./bin/$(BINARY_NAME)

# Run backend servers for testing (requires separate terminals)
run-backends: backend
	@echo "Start backend servers in separate terminals:"
	@echo "Terminal 1: ./bin/$(BACKEND_BINARY) -port 8081"
	@echo "Terminal 2: ./bin/$(BACKEND_BINARY) -port 8082"
	@echo "Terminal 3: ./bin/$(BACKEND_BINARY) -port 8083"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  all            - Run tests and build"
	@echo "  build          - Build the load balancer"
	@echo "  backend        - Build the backend server"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  run            - Build and run the load balancer"
	@echo "  run-backends   - Show commands to run backend servers"
	@echo "  clean          - Remove build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  help           - Show this help message"
