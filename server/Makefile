.PHONY: run build test clean proto migrate-up migrate-down migrate-create

# Run the server in development mode
run:
	go run cmd/server/main.go

# Run the server from cmd directory
run-cmd:
	go run cmd/server/main.go

# Build the server binary
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Development server with auto-reload (requires air)
dev:
	air

# Generate gRPC stubs (Go/Python) via Buf v2
proto:
	cd .. && buf generate

# Run database migrations
migrate-up:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable is not set"; \
		exit 1; \
	fi
	migrate -path migrations -database "mysql://$(DATABASE_URL)" up

# Rollback database migrations
migrate-down:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable is not set"; \
		exit 1; \
	fi
	migrate -path migrations -database "mysql://$(DATABASE_URL)" down

# Create a new migration file
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME variable is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(NAME)
