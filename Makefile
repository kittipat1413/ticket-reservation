.PHONY: help install gen-all gen-swag gen-db gen-mock precommit lint vet fmt new-migration migrate-up migrate-down test test-coverage open-coverage-report

# Makefile for Go project
# This Makefile provides a set of commands to manage the Go project, including building, testing, and generating database models.
PKG := ticket-reservation
GOLINT ?= golangci-lint
GO_FILES = $(shell go list ./... | grep -v -e /mocks -e /example)
GO_BIN = $(shell go env GOPATH)/bin

# Default target is help, which lists available commands.
help:
	@echo "Available commands ⚙️"
	@echo "  make install                  - Install necessary Go tools for the project"
	@echo "  make gen-all                  - Generate all necessary files (Swagger, DB models, mocks)"
	@echo "  make gen-swag                 - Generate Swagger documentation"
	@echo "  make gen-db                   - Generate database models"
	@echo "  make gen-mock                 - Generate mock files"
	@echo "  make precommit                - Run linters, go vet, and go fmt"
	@echo "  make lint                     - Run linters"
	@echo "  make vet                      - Run go vet"
	@echo "  make fmt                      - Format Go code"
	@echo "  make test                     - Run tests"
	@echo "  make test-coverage            - Run tests with coverage and generate a coverage report"
	@echo "  make open-coverage-report     - Open the coverage report in a web browser"
	@echo "  make new-migration            - Create a new migration file"
	@echo "  make migrate-up               - Apply all pending migrations"
	@echo "  make migrate-down             - Roll back the last migration"

# install target installs the necessary Go tools.
install:
	@echo "Installing Go tools... 🚀"
	@test -e $(GO_BIN)/mockgen || go install github.com/golang/mock/mockgen@v1.7.0-rc.1
	@test -e $(GO_BIN)/swag || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Go tools installed successfully. ✅"

################################ Generation Targets ################################

# gen-all target generates all necessary files.
# It runs gen-swag, gen-db, and gen-mock targets.
gen-all: gen-swag gen-db gen-mock

# gen-swag target generates Swagger documentation using swag.
# docs: https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format
gen-swag:
	@echo "Generating Swagger documentation... 📜"
	@swag fmt
	@swag init \
		--generalInfo ./internal/api/http/route/routes.go \
		--output ./docs \
		--outputTypes yml
	@echo "Swagger documentation generated successfully. ✅"

# gen-db target generates database models using go-jet.
# docs: https://github.com/go-jet/jet/wiki/Generator
gen-db:
	@echo "Generating database models... 🚧"
	@go run main.go generate-db
	@echo "Database models generated successfully. ✅"

# gen-mock target generates mock files using mockgen.
gen-mock:
	@echo "Generating mock files... 🚧"
	@go generate ./...
	@echo "Mock files generated successfully. ✅"

################################ Precommit ################################

# precommit target runs the linters, go vet, and go fmt.
precommit: lint vet fmt test

# lint target runs the linters using golangci-lint.
lint:
	@echo "Running linters... 🧹"
	@$(GOLINT) run
	@echo "Linters completed successfully. ✅"

# vet target runs go vet on the project.
vet:
	@echo "Running go vet... 🔍"
	@go vet $(GO_FILES)
	@echo "go vet completed successfully. ✅"

# fmt target formats the Go code.
fmt:
	@echo "Formatting Go code... 📝"
	@go fmt $(GO_FILES)
	@echo "Go code formatted successfully. ✅"

# test target runs the tests for the project.
test:
	@echo "Running tests... 🧪"
	@go test $(GO_FILES)/... -cover --race

# test-coverage target runs the tests with coverage and generates a coverage report.
test-coverage:
	@echo "Running tests with coverage... 🧪"
	@go test $(GO_FILES)/... -race -covermode=atomic -coverprofile coverage.out
	@go tool cover -func=coverage.out -o=coverage_summary.out
	@cat coverage_summary.out | grep total | awk '{print "Total coverage: " $$3}'

# open-coverage-report target opens the coverage report in a web browser.
open-coverage-report:
	@echo "Opening coverage report... 📊"
	@go tool cover -html coverage.out -o coverage.html;
	@open coverage.html

################################ Migrations ################################

# new-migration target creates a new migration file.
new-migration:
	@echo "Creating new migration... 🛠️"
	@read -p "Enter migration name: " MIGRATION_NAME; \
	if [ -z "$$MIGRATION_NAME" ]; then \
		echo "Migration name cannot be empty. ❌"; \
		exit 1; \
	fi; \
	go run main.go new-migration "$$MIGRATION_NAME"; \
	echo "Migration created successfully. ✅"

# migrate-up target applies all pending migrations.
migrate-up:
	@echo "Applying migrations... ⬆️"
	@go run main.go migrate --action up
	@echo "Migrations applied successfully. ✅"

# migrate-down target rolls back the last migration.
migrate-down:
	@echo "Rolling back to previous migration... ⬇️"
	@go run main.go migrate --action down
	@echo "Rolled back one migration successfully. ✅"