# Makefile for Go API project

.PHONY: test test-verbose test-cover test-bench clean build run help

# Default target
help:
	@echo "Available commands:"
	@echo "  test        - Run all tests"
	@echo "  test-verbose - Run all tests with verbose output"
	@echo "  test-cover  - Run tests with coverage report"
	@echo "  test-bench  - Run benchmark tests"

	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  clean       - Clean build artifacts"

	@echo "  migrate         - Run all migrations up"
	@echo "  migrate-down    - Rollback last migration"
	@echo "  migrate-fresh   - Fresh the migration and delete entire data in database"
	@echo "  migrate-status  - Show migration status"
	@echo "  migrate-create  - Create new migration file (usage: make migrate-create name=your_migration_name)"

# Run all tests
test:
	go test ./...

# Run all tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -cover ./...

# Generate detailed coverage report
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmark tests
test-bench:
	go test -bench=. ./...

# Build the application
build:
	go build -o bin/main main.go

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests and build (CI pipeline)
ci: test build

# Run specific test domain
test-domain:
	go test -v ./domain/...

# Run specific test pkg
test-pkg:
	go test -v -cover ./pkg/...

# Run migration test
test-migrations:
	go test -v -cover -count=1 ./pkg/migrations/...

# Format code
fmt:
	go fmt ./...

# Migration commands (using goose)
migrate:
	goose up

migrate-down:
	goose down

migrate-down-all:
	goose down-to 0

migrate-status:
	goose status

migrate-fresh:
	goose down-to 0 && goose up

migrate-create:
	goose create $(name) sql