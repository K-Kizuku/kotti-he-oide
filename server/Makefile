.PHONY: run build test clean proto migrate-up migrate-down migrate-create migrate-up-rds migrate-down-rds

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

# Run database migrations on RDS
# Required environment variables:
#   RDS_ENDPOINT: RDS endpoint (e.g., mydb.abc123.ap-northeast-1.rds.amazonaws.com)
#   RDS_PORT: RDS port (default: 3306)
#   RDS_DATABASE: Database name
#   RDS_USER: Database username
#   RDS_PASSWORD: Database password
migrate-up-rds:
	@if [ -z "$(RDS_ENDPOINT)" ]; then \
		echo "Error: RDS_ENDPOINT environment variable is not set"; \
		exit 1; \
	fi
	@if [ -z "$(RDS_DATABASE)" ]; then \
		echo "Error: RDS_DATABASE environment variable is not set"; \
		exit 1; \
	fi
	@if [ -z "$(RDS_USER)" ]; then \
		echo "Error: RDS_USER environment variable is not set"; \
		exit 1; \
	fi
	@if [ -z "$(RDS_PASSWORD)" ]; then \
		echo "Error: RDS_PASSWORD environment variable is not set"; \
		exit 1; \
	fi
	@RDS_PORT=$${RDS_PORT:-3306}; \
	echo "Connecting to RDS: $(RDS_USER)@$(RDS_ENDPOINT):$$RDS_PORT/$(RDS_DATABASE)"; \
	migrate -path migrations -database "mysql://$(RDS_USER):$(RDS_PASSWORD)@tcp($(RDS_ENDPOINT):$$RDS_PORT)/$(RDS_DATABASE)?charset=utf8mb4&parseTime=True&loc=Local" up

# Rollback database migrations on RDS
migrate-down-rds:
	@if [ -z "$(RDS_ENDPOINT)" ]; then \
		echo "Error: RDS_ENDPOINT environment variable is not set"; \
		exit 1; \
	fi
	@if [ -z "$(RDS_DATABASE)" ]; then \
		echo "Error: RDS_DATABASE environment variable is not set"; \
		exit 1; \
	fi
	@if [ -z "$(RDS_USER)" ]; then \
		echo "Error: RDS_USER environment variable is not set"; \
		exit 1; \
	fi
	@if [ -z "$(RDS_PASSWORD)" ]; then \
		echo "Error: RDS_PASSWORD environment variable is not set"; \
		exit 1; \
	fi
	@RDS_PORT=$${RDS_PORT:-3306}; \
	echo "Connecting to RDS: $(RDS_USER)@$(RDS_ENDPOINT):$$RDS_PORT/$(RDS_DATABASE)"; \
	migrate -path migrations -database "mysql://$(RDS_USER):$(RDS_PASSWORD)@tcp($(RDS_ENDPOINT):$$RDS_PORT)/$(RDS_DATABASE)?charset=utf8mb4&parseTime=True&loc=Local" down
