# Makefile for Go project
# This Makefile provides a set of commands to manage the Go project, including building, testing, and generating database models.
PKG := ticket-reservation
GOLINT ?= golangci-lint
GO_FILES = $(shell go list ./... | grep -v -e /mocks -e /example)
GO_BIN = $(shell go env GOPATH)/bin


# install target installs the necessary Go tools.
install:
	@echo "Installing Go tools... 🚀"
	@test -e $(GO_BIN)/mockgen || go install github.com/golang/mock/mockgen@v1.7.0-rc.1
	@test -e $(GO_BIN)/swag || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Go tools installed successfully. ✅"

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