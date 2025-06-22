# 🎟️ Ticket Reservation System
A simple ticket reservation system for managing concert tickets, including seat reservations and payment processing.
This project is designed to be a starting point for building a more complex ticket reservation system.
> **Note:** The high-level system design, API design details, and architectural diagrams are available in the [`/docs`](./docs) directory.

> **🧩 Powered by go-common:**
This project integrates [go-common](https://github.com/kittipat1413/go-common) — our standardized backend framework — across all layers:
Structured Logging, Error Handling, Config Loader (Viper), Validation & Middleware Utilities, Tracing & Retry Helpers and more...

## 📝 Overview

The project leverages a Clean Architecture structure to cleanly separate concerns:

- **Delivery Layer:** HTTP API handlers, routes, and middleware under `internal/api/http/`
- **Application (Usecase) Layer:** Business logic under `internal/usecase/`
- **Domain Layer:** Core domain models (entities) and repository interfaces under `internal/domain/`
- **Infrastructure Layer:** Database connections, migrations, and external integrations under `internal/infra/`
- **Shared Utilities:** Helpers and configurations under `pkg/` and `internal/util/`

The project also includes several CLI commands (in the [`cmd/`](./cmd) directory) for tasks such as serving the API, running migrations, generating new migrations, printing configuration, and auto-generating SQL builder files.


## 📂 Project Structure
```
ticket-reservation/
├── cmd/                          # Application entry points and CLI commands
│   ├── root.go                   # Root command for the CLI
│   ├── serve_cmd.go              # Command to start the server
│   ├── migrate_cmd.go            # Command to run database migrations
│   ├── new_migration_cmd.go      # Command to create new migration files
│   ├── print_config_cmd.go       # Command to print the current configuration
│   └── generate_sql_builder.go   # Command to generate SQL builder files
│   └── ...                       # Other commands
├── db/                           # Database-related files
│   └── migrations/               # Database migration files
├── docs/                         # Documentation files
│   ├── swagger.yaml              # API documentation in Swagger format
│   └── ...                       # Other design and architecture documents
├── internal/                     # Internal application code
│   ├── api/                      # Delivery layer
│   │   └── http/                 # HTTP API components
│   │       ├── handler/          # HTTP handlers by domain
│   │       │   ├── healthcheck/  # Health check handlers
│   │       │   └── ...           # Other domain handlers
│   │       ├── middleware/       # HTTP middleware
│   │       └── route/            # HTTP route definitions
│   ├── config/                   # Application configuration
│   ├── domain/                   # Domain layer
│   │   ├── cache/                # Cache-related interfaces
│   │   ├── entity/               # Domain models (e.g., Concert, Reservation)
│   │   ├── errs/                 # Domain-specific errors
│   │   └── repository/           # Repository interfaces
│   ├── infra/                    # Infrastructure layer
│   │   ├── db/                   # Database implementations
│   │   │   ├── connection.go     # Database connection setup
│   │   │   ├── sql_execer.go     # SQL execution interface
│   │   │   ├── transactor.go     # Transaction management
│   │   │   ├── mocks/            # Mock implementations
│   │   │   ├── model_gen/        # Generated models from DB schema
│   │   │   └── repository/       # Repository implementations
│   │   │       ├── healthcheck/  # Health check repository
│   │   │       └── ...           # Other repositories
│   │   └── redis/                # Redis implementations
│   │       ├── client.go         # Redis client interface
│   │       ├── connection.go     # Redis connection setup
│   │       ├── mocks/            # Mock implementations
│   │       └── repository/       # Repository implementations
│   │           ├── seat/         # Seat locking and cache implementations
│   │           └── ...           # Other repositories
│   ├── server/                   # Server setup and initialization
│   │   ├── dependency.go         # Dependency injection
│   │   ├── middleware.go         # Server middleware
│   │   └── server.go             # HTTP server setup
│   ├── usecase/                  # Application business logic
│   │   ├── healthcheck/          # Health check use case
│   │   └── ...                   # Other use cases
│   └── util/                     # Utility functions
│       ├── httpresponse/         # HTTP response helpers
│       └── ...                   # Other utility functions
├── pkg/                          # Shared helper packages
├── docker-compose.yaml           # Docker Compose for local development
├── Dockerfile                    # Docker build definition
├── env.yaml                      # Environment variables configuration
├── go.mod                        # Go module definition
├── go.sum                        # Go module dependencies
├── main.go                       # Main application entry point
├── Makefile                      # Makefile for common tasks
├── otel-collector-config.yaml    # OpenTelemetry collector configuration
└── README.md                     # Project documentation
```

## 🚀 Getting Started

### Installation & Setup
1.	**Install the required Go tools**:
Use the provided `Makefile` target to install tools like `swag` for Swagger docs and `mockgen` for generating mocks:
```bash
make install
```
2. **Generate Documentation and Database Models**:
- _**Generate Swagger Docs**_: The project uses `swaggo/swag` to generate API documentation from annotated code. To generate the Swagger files, run:
	```bash
	make gen-swag
	```
	> (The Swagger YAML file will be created in the `/docs` directory.)
- _**Generate Database Models**_: The project uses `go-jet` for generating type-safe SQL builder files. Generate these files with:
	```bash
	make gen-db
	```
	> (The generated files will be placed in the `/internal/infra/db/model_gen` directory.)
- _**Generate Mocks**_: Generate mock implementations for any interfaces that are used across the project for testing purposes:
	```bash
	make gen-mock
	```
	> These mocks are generated using go generate directives (e.g., //go:generate mockgen) and placed in appropriate mocks/ directories next to the interfaces they mock (e.g., internal/domain/repository/mocks/, internal/usecase/mocks/, etc.).
3. **Starting the Application**:
You can start the server with the CLI command provided:
```bash
go run main.go serve
```

## 🧪 Testing Strategy

The project follows a comprehensive testing strategy focusing on testability through interfaces:

- **Interface-Driven Design:** All major components are defined through interfaces, allowing easy mock substitution for testing
- **Mock Implementations:** Mock implementations are provided for all interfaces in the `mocks/` directories throughout the codebase
- **Unit Tests:** Unit tests are implemented for key business logic, particularly in the usecase layer
- **Test Examples:** See `internal/usecase/seat/seat_usecase_test.go` for an example of how to use mocks to test the seat reservation flow

To run the tests:
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/usecase/seat/...
```

## 🛠️ Available Cobra Commands

The project leverages a CLI powered by Cobra with the following key commands:
- `serve`: Start the application server.
- `migrate`: Run database migrations (up or down).
- `new-migration`: Create new migration files with a timestamped filename.
- `print-config`: Print the current effective configuration.
- `generate-db`: Generate SQL builder and model files using go-jet.

For a full list of commands, run:
```bash
go run main.go --help
```

## 🔍 Design & Documentation
- System Design & Architecture:
The architectural decisions, design patterns, and system diagrams are documented in the [`/docs`](./docs) directory.
- API Documentation:
The generated Swagger documentation (see [`/docs/swagger.yaml`](./docs/swagger.yaml)) provides a detailed API specification, including endpoints, request/response schemas, and usage examples.

## ⚙️ Makefile Overview
The project includes a `Makefile` to simplify common tasks and commands. Here are the available targets:
```
Available commands ⚙️
  make install                  - Install necessary Go tools for the project
  make gen-all                  - Generate all necessary files (Swagger, DB models, mocks)
  make gen-swag                 - Generate Swagger documentation
  make gen-db                   - Generate database models
  make gen-mock                 - Generate mock files
  make precommit                - Run linters, go vet, and go fmt
  make lint                     - Run linters
  make vet                      - Run go vet
  make fmt                      - Format Go code
  make test                     - Run tests
  make test-coverage            - Run tests with coverage and generate a coverage report
  make open-coverage-report     - Open the coverage report in a web browser
  make new-migration            - Create a new migration file
  make migrate-up               - Apply all pending migrations
  make migrate-down             - Roll back the last migration
```

## 🚀 Getting Started
### ⚡ Quick Start
Spin up the full local dev environment using Docker Compose:
```bash
docker-compose up --build
```