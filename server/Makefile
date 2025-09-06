.PHONY: run build test clean

# Run the server in development mode
run:
	go run main.go

# Run the server from cmd directory
run-cmd:
	go run cmd/server/main.go

# Build the server binary
build:
	go build -o bin/server main.go

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