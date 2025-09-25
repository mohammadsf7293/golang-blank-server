.PHONY: all build test clean sqlc run docker-up docker-down mocks

# Default target
all: build

# Generate mocks
mocks:
	go generate ./...

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Generate SQL code
sqlc:
	sqlc generate

# Run the application
run:
	go run cmd/api/main.go

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Start the application with all dependencies
start: docker-up sqlc run

# Stop all services
stop: docker-down
	pkill -f "go run cmd/api/main.go" || true

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Install dependencies
deps:
	go mod tidy
	go mod verify

# Help target
help:
	@echo "Available targets:"
	@echo "  make          : Build the application"
	@echo "  make build    : Build the application"
	@echo "  make test     : Run tests"
	@echo "  make clean    : Clean build artifacts"
	@echo "  make sqlc     : Generate SQL code"
	@echo "  make run      : Run the application"
	@echo "  make docker-up: Start Docker services"
	@echo "  make docker-down: Stop Docker services"
	@echo "  make start    : Start application with dependencies"
	@echo "  make stop     : Stop all services"
	@echo "  make fmt      : Format code"
	@echo "  make lint     : Run linter"
	@echo "  make deps     : Install dependencies"