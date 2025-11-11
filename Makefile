.PHONY: help docker-up docker-down docker-logs docker-clean migrate-up migrate-down migrate-force migrate-version migrate-create dev build run test

# Default target
help:
	@echo "Available commands:"
	@echo "  make docker-up         - Start PostgreSQL container"
	@echo "  make docker-down       - Stop PostgreSQL container"
	@echo "  make docker-logs       - View PostgreSQL logs"
	@echo "  make docker-clean      - Stop and remove container + volume (deletes all data)"
	@echo "  make migrate-up        - Run all pending migrations"
	@echo "  make migrate-down      - Rollback last migration"
	@echo "  make migrate-force V=n - Force migration version (use with caution)"
	@echo "  make migrate-version   - Show current migration version"
	@echo "  make migrate-create NAME=name - Create new migration files"
	@echo "  make dev               - Start development server with Air"
	@echo "  make build             - Build the application"
	@echo "  make run               - Run the application"
	@echo "  make test              - Run tests"

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f postgres

docker-clean:
	docker-compose down -v

# Migration commands (require migrate CLI: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/event_booking?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/event_booking?sslmode=disable" down 1

migrate-force:
	@if [ -z "$(V)" ]; then echo "Usage: make migrate-force V=<version>"; exit 1; fi
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/event_booking?sslmode=disable" force $(V)

migrate-version:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/event_booking?sslmode=disable" version

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=<migration_name>"; exit 1; fi
	migrate create -ext sql -dir db/migrations -seq $(NAME)

# Development commands
dev:
	air

build:
	go build -o bin/api main.go

run:
	go run main.go

test:
	go test -v ./...
