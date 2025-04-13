# Makefile for Notification Service

# Variables
BINARY_NAME=notification-service
MAIN_PACKAGE=.
DOCKER_COMPOSE=docker-compose

# Go commands
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build

# Docker commands
DOCKER=docker
DOCKER_BUILD=$(DOCKER) build

.PHONY: all build clean test coverage run docker-build docker-run docker-stop help

all: test build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

docker-build:
	@echo "Building Docker image..."
	$(DOCKER_BUILD) -t $(BINARY_NAME) .

docker-run:
	@echo "Starting services with Docker Compose..."
	$(DOCKER_COMPOSE) up -d

docker-stop:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make coverage    - Generate test coverage report"
	@echo "  make run         - Build and run the application"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Start services with Docker Compose"
	@echo "  make docker-stop - Stop services"
	@echo "  make help        - Show this help message"
