# PostEaze Backend Makefile
# Provides convenient commands for testing, building, and development

.PHONY: help test test-unit test-integration test-coverage test-verbose test-race test-bench clean build run dev lint fmt vet deps

# Default target
help: ## Show this help message
	@echo "PostEaze Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make test                    # Run all tests"
	@echo "  make test-unit               # Run unit tests only"
	@echo "  make test-coverage           # Run tests with coverage"
	@echo "  make test-integration        # Run integration tests"
	@echo "  make test-verbose            # Run tests with verbose output"

# Test commands
test: ## Run all tests
	@echo "Running all tests..."
	@./scripts/test.sh --type all

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@./scripts/test.sh --type unit

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@./scripts/test.sh --type integration

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@./scripts/test.sh --type all --coverage

test-coverage-unit: ## Run unit tests with coverage
	@echo "Running unit tests with coverage..."
	@./scripts/test.sh --type unit --coverage

test-coverage-integration: ## Run integration tests with coverage
	@echo "Running integration tests with coverage..."
	@./scripts/test.sh --type integration --coverage

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@./scripts/test.sh --type all --verbose

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@./scripts/test.sh --type all --race

test-bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	@./scripts/test.sh --benchmark

test-short: ## Run tests in short mode
	@echo "Running tests in short mode..."
	@./scripts/test.sh --short

# Specific test suites
test-api: ## Run API handler tests
	@echo "Running API handler tests..."
	@./scripts/test.sh --type unit --suite api

test-business: ## Run business logic tests
	@echo "Running business logic tests..."
	@./scripts/test.sh --type unit --suite business

test-models: ## Run model tests
	@echo "Running model tests..."
	@./scripts/test.sh --type unit --suite models

test-utils: ## Run utility tests
	@echo "Running utility tests..."
	@./scripts/test.sh --type unit --suite utils

test-auth: ## Run authentication integration tests
	@echo "Running authentication integration tests..."
	@./scripts/test.sh --type integration --suite auth

test-log: ## Run log integration tests
	@echo "Running log integration tests..."
	@./scripts/test.sh --type integration --suite log

test-database: ## Run database integration tests
	@echo "Running database integration tests..."
	@./scripts/test.sh --type integration --suite database

test-e2e: ## Run end-to-end integration tests
	@echo "Running end-to-end integration tests..."
	@./scripts/test.sh --type integration --suite e2e

# Coverage commands with thresholds
test-coverage-80: ## Run tests with 80% coverage threshold
	@echo "Running tests with 80% coverage threshold..."
	@./scripts/test.sh --coverage --threshold 80

test-coverage-90: ## Run tests with 90% coverage threshold
	@echo "Running tests with 90% coverage threshold..."
	@./scripts/test.sh --coverage --threshold 90

# Parallel test execution
test-parallel: ## Run tests in parallel
	@echo "Running tests in parallel..."
	@./scripts/test.sh --parallel

test-parallel-coverage: ## Run tests in parallel with coverage
	@echo "Running tests in parallel with coverage..."
	@./scripts/test.sh --parallel --coverage

# Database-specific tests
test-sqlite: ## Run integration tests with SQLite
	@echo "Running integration tests with SQLite..."
	@./scripts/test.sh --type integration --database sqlite

test-postgres: ## Run integration tests with PostgreSQL
	@echo "Running integration tests with PostgreSQL..."
	@./scripts/test.sh --type integration --database postgres

# Output format tests
test-json: ## Run tests with JSON output
	@echo "Running tests with JSON output..."
	@./scripts/test.sh --format json

test-junit: ## Run tests with JUnit output
	@echo "Running tests with JUnit output..."
	@./scripts/test.sh --format junit

# CI/CD specific commands
test-ci: ## Run tests for CI environment (coverage + race detection)
	@echo "Running CI tests..."
	@./scripts/test.sh --coverage --race --parallel --threshold 80

test-ci-unit: ## Run unit tests for CI
	@echo "Running CI unit tests..."
	@./scripts/test.sh --type unit --coverage --race --parallel

test-ci-integration: ## Run integration tests for CI
	@echo "Running CI integration tests..."
	@./scripts/test.sh --type integration --coverage --parallel

# Development commands
build: ## Build the application
	@echo "Building application..."
	@go build -o bin/postease ./main.go

run: ## Run the application
	@echo "Running application..."
	@go run main.go

dev: ## Run application in development mode
	@echo "Running application in development mode..."
	@MODE=dev go run main.go

# Code quality commands
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go mod tidy

deps-vendor: ## Vendor dependencies
	@echo "Vendoring dependencies..."
	@go mod vendor

# Clean commands
clean: ## Clean build artifacts and test files
	@echo "Cleaning build artifacts..."
	@rm -f bin/postease
	@rm -f coverage.out coverage.html
	@rm -f integration_coverage.out integration_coverage.html
	@rm -f *.test
	@rm -rf vendor/

clean-test: ## Clean test artifacts only
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@rm -f integration_coverage.out integration_coverage.html
	@rm -f *.test

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t postease-backend .

docker-test: ## Run tests in Docker container
	@echo "Running tests in Docker container..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.23 make test

# Database commands
db-setup: ## Setup test database
	@echo "Setting up test database..."
	@if command -v psql >/dev/null 2>&1; then \
		createdb postease_test || echo "Database may already exist"; \
	else \
		echo "PostgreSQL not found. Using SQLite for tests."; \
	fi

db-clean: ## Clean test database
	@echo "Cleaning test database..."
	@if command -v psql >/dev/null 2>&1; then \
		dropdb postease_test || echo "Database may not exist"; \
	fi

# Performance testing
perf-test: ## Run performance tests
	@echo "Running performance tests..."
	@./scripts/test.sh --benchmark --race

perf-profile: ## Run tests with CPU profiling
	@echo "Running tests with CPU profiling..."
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...

# Security testing
security-test: ## Run security tests
	@echo "Running security tests..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Test data management
test-fixtures: ## Generate test fixtures
	@echo "Generating test fixtures..."
	@go run tests/testutils/fixtures.go

# Continuous testing
test-watch: ## Run tests in watch mode (requires entr)
	@echo "Running tests in watch mode..."
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -c make test; \
	else \
		echo "entr not found. Install it to enable watch mode."; \
	fi

# Report generation
test-report: ## Generate comprehensive test report
	@echo "Generating comprehensive test report..."
	@./scripts/test.sh --coverage --benchmark --format junit
	@echo "Test report generated with coverage and benchmark data"

# Quick commands for development
quick-test: ## Quick test run (short mode, no coverage)
	@echo "Running quick tests..."
	@./scripts/test.sh --short

full-test: ## Full test suite (all tests, coverage, race detection)
	@echo "Running full test suite..."
	@./scripts/test.sh --coverage --race --benchmark --parallel

# Help for Windows users
help-windows: ## Show Windows-specific commands
	@echo "Windows PowerShell Commands:"
	@echo "  .\scripts\test.ps1 -Type all"
	@echo "  .\scripts\test.ps1 -Coverage"
	@echo "  .\scripts\test.ps1 -Type unit -Suite api"
	@echo "  .\scripts\test.ps1 -Help"