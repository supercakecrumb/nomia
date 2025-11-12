.PHONY: help build test clean docker-build docker-up docker-down migrate-up migrate-down changelog-new changelog-batch changelog-merge release

# Default target
help:
	@echo "Available targets:"
	@echo "  build              - Build all binaries"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  docker-build       - Build Docker images"
	@echo "  docker-up          - Start Docker containers"
	@echo "  docker-down        - Stop Docker containers"
	@echo "  migrate-up         - Run database migrations"
	@echo "  migrate-down       - Rollback database migrations"
	@echo "  changelog-new      - Create new changelog entry"
	@echo "  changelog-batch    - Batch unreleased changes"
	@echo "  changelog-merge    - Merge changes into CHANGELOG.md"
	@echo "  release            - Create new release (VERSION=patch|minor|major)"

# Build targets
build:
	@echo "Building binaries..."
	@go build -o bin/api ./cmd/api
	@go build -o bin/worker ./cmd/worker
	@go build -o bin/migrate ./cmd/migrate
	@echo "Build complete!"

# Test targets
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Tests complete! Coverage report: coverage.html"

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Docker targets
docker-build:
	@echo "Building Docker images..."
	@docker-compose build
	@echo "Docker build complete!"

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Docker containers started!"

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Docker containers stopped!"

# Migration targets
migrate-up:
	@echo "Running database migrations..."
	@./bin/migrate up
	@echo "Migrations complete!"

migrate-down:
	@echo "Rolling back database migrations..."
	@./bin/migrate down
	@echo "Rollback complete!"

# Changelog targets
changelog-new:
	@echo "Creating new changelog entry..."
	@if ! command -v changie &> /dev/null; then \
		echo "Error: changie is not installed. Install it with:"; \
		echo "  go install github.com/miniscruff/changie@latest"; \
		exit 1; \
	fi
	@changie new

changelog-batch:
	@echo "Batching unreleased changes..."
	@if ! command -v changie &> /dev/null; then \
		echo "Error: changie is not installed. Install it with:"; \
		echo "  go install github.com/miniscruff/changie@latest"; \
		exit 1; \
	fi
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make changelog-batch VERSION=0.2.0"; \
		exit 1; \
	fi
	@changie batch $(VERSION)
	@echo "Changes batched for version $(VERSION)"

changelog-merge:
	@echo "Merging changes into CHANGELOG.md..."
	@if ! command -v changie &> /dev/null; then \
		echo "Error: changie is not installed. Install it with:"; \
		echo "  go install github.com/miniscruff/changie@latest"; \
		exit 1; \
	fi
	@changie merge
	@echo "Changes merged into CHANGELOG.md"

# Release target
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=patch|minor|major"; \
		exit 1; \
	fi
	@if [ ! -f "scripts/release.sh" ]; then \
		echo "Error: scripts/release.sh not found"; \
		exit 1; \
	fi
	@bash scripts/release.sh $(VERSION)

# Install development dependencies
install-deps:
	@echo "Installing development dependencies..."
	@go install github.com/miniscruff/changie@latest
	@echo "Dependencies installed!"

# Lint code
lint:
	@echo "Running linters..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

# Run all checks
check: fmt lint test
	@echo "All checks passed!"