.PHONY: help test test-cover test-race run build clean lint fmt vet

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run all tests
	go test -v ./...

test-cover: ## Run tests with coverage report
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	go test -race ./...

test-watch: ## Run tests in watch mode (requires entr)
	find . -name '*.go' | entr -c go test ./...

run: ## Run the server
	go run cmd/server/main.go

build: ## Build the server binary
	go build -o bin/server cmd/server/main.go
	@echo "Binary built: bin/server"

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	gofmt -s -w .
	go mod tidy

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)

dev: ## Start development server
	@echo "Starting server in development mode (fixture mode)..."
	FIXTURE_MODE=true go run cmd/server/main.go

prod: ## Start production server
	@echo "Starting server in production mode (database mode)..."
	FIXTURE_MODE=false go run cmd/server/main.go

deps: ## Download dependencies
	go mod download
	go mod verify

tidy: ## Tidy dependencies
	go mod tidy

install-tools: ## Install development tools
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

db-up: ## Start database
	cd .. && docker-compose up -d

db-down: ## Stop database
	cd .. && docker-compose down

db-logs: ## View database logs
	cd .. && docker-compose logs -f postgres

db-reset: ## Reset database (removes all data)
	cd .. && docker-compose down -v
	cd .. && docker-compose up -d

import-data: ## Import US name data
	bash scripts/import-us-data.sh

.DEFAULT_GOAL := help