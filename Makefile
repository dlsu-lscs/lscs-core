# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Database migrations (requires goose: go install github.com/pressly/goose/v3/cmd/goose@latest)
GOOSE_DRIVER=mysql
GOOSE_DBSTRING=$(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_DATABASE)?parseTime=true&multiStatements=true
GOOSE_MIGRATION_DIR=./migrations

migrate-up:
	@echo "Running migrations..."
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up

migrate-down:
	@echo "Rolling back last migration..."
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down

migrate-status:
	@echo "Migration status..."
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" status

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(GOOSE_MIGRATION_DIR) create $$name sql

# NOTE: goose doesn't have a builtin "mark migrations as applied" without overwriting the database (for existing dbs with data)
# so we just workaround by manually adding the table for goose
migrate-baseline:
	@echo "Marking baseline migration as applied..."
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" version
	@echo "INSERT INTO goose_db_version (version_id, is_applied, tstamp) VALUES (1, true, NOW()) ON DUPLICATE KEY UPDATE is_applied=true;" | mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USERNAME) -p$(DB_PASSWORD) $(DB_DATABASE)

# Swagger documentation (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

.PHONY: all build run test clean watch docker-run docker-down itest migrate-up migrate-down migrate-status migrate-create migrate-baseline swagger
