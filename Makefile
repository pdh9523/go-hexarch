# Makefile for go-hexarch

.PHONY: build run test coverage coverage-profile coverage-html coverage-func clean run deps-update wire wire-clean install-tools health help migrate-up migrate-down migrate-create migrate-status migrate-force

build:
	@echo "Build application..."
	go work sync
	go build -o bin/web-api ./applications/web-api

run:
	@echo "Running server based on web-api..."
	go run ./applications/web-api

test:
	@echo "Running all tests..."
	go work sync
	make test-security

test-security:
	@echo "Running security test..."
	cd shared/security && go test ./...

coverage:
	@echo "Running tests with coverage..."
	go work sync

coverage-profile:
	@echo "Generating coverage profiles..."
	go work sync
	@mkdir -p coverage

coverage-html: coverage-profile
	@echo "Generating HTML coverage reports..."
	@mkdir -p coverage/html

coverage-html: coverage-profile
	@echo "Generating HTML coverage reports..."
	@mkdir -p coverage/html

coverage-func: coverage-profile
	@echo "Function coverage report:"

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf coverage/
	go clean -cache

deps-update:
	@echo "Updating dependencies..."
	go work sync
	go get -u ./...
	go mod tidy

wire:
	@echo "Running wire for dependency injection..."
	@if ! command -v wire &> /dev/null; then \
  		echo "Wire not found. Installing..." \
  		go install github.com/google/wire/cmd/wire@latest; \
  	fi
	cd applications/web-api/config && wire
	@echo "Wire code generation completed."

wire-clean:
	@echo "Cleaning wire generated files..."
	find . -name "wire_gen.go" -type f -delete
	@echo "Wire generated files cleaned."

install-tools:
	@echo "Installing development tools..."
	go install github.com/google/wire/cmd/wire@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

health:
	@echo "Checking health..."
	curl -f http://localhost:8080/health || echo "Server is not running"

migrate-up:
	@echo "Running database migrations..."
	@if ! command -v migrate &> /dev/null; then \
		echo "migrate not found. Installing..."; \
		make install-tools; \
	fi
	migrate -path migrations -database "$${DATABASE_URL:-postgres://admin:1234@localhost:5432/go-hexarch?sslmode=disable}" up
	@echo "Migrations completed."

migrate-down:
	@echo "Rolling back last migration..."
	migrate -path migrations -database "$${DATABASE_URL:-postgres://admin:1234@localhost:5432/go-hexarch?sslmode=disable}" down 1
	@echo "Rollback completed."

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(name)..."
	@if ! command -v migrate &> /dev/null; then \
		echo "migrate not found. Installing..."; \
		make install-tools; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)
	@echo "Migration files created."

migrate-status:
	@echo "Checking migration status..."
	migrate -path migrations -database "$${DATABASE_URL:-postgres://admin:1234@localhost:5432/go-hexarch?sslmode=disable}" version

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: Version is required. Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	@echo "Forcing migration version to $(version)..."
	migrate -path migrations -database "$${DATABASE_URL:-postgres://admin:1234@localhost:5432/go-hexarch?sslmode=disable}" force $(version)

help:
	@echo "Available commands:"
	@echo "  build                        - Build applications"
	@echo "  run-web                      - Run web-api server"
	@echo "  test                         - Run all tests"
	@echo "  test-{domain}                - Run specific domain tests"
	@echo "  coverage                     - Run tests with coverage"
	@echo "  coverage-profile             - Generate coverage profiles"
	@echo "  coverage-html                - Generate HTML coverage reports"
	@echo "  coverage-func                - Show function coverage report"
	@echo "  clean                        - Clean build artifacts"
	@echo "  deps-update                  - Update dependencies"
	@echo "  generate                     - Generate code (wire)"
	@echo "  wire                         - Run wire for dependency injection"
	@echo "  wire-clean                   - Clean wire generated files"
	@echo "  install-tools                - Install development tools"
	@echo "  health                       - Check server health"
	@echo "  help                         - Show this help message"
	@echo ""
	@echo "Database Migration commands:"
	@echo "  migrate-up     - Run all pending migrations"
	@echo "  migrate-down   - Rollback last migration"
	@echo "  migrate-create name=<name> - Create new migration files"
	@echo "  migrate-status - Show current migration version"
	@echo "  migrate-force version=<n> - Force set migration version"