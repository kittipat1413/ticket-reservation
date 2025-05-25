# 🎟️ Ticket Reservation System
A simple ticket reservation system for managing concert tickets, including seat reservations and payment processing.
This project is designed to be a starting point for building a more complex ticket reservation system.
> **Note:** The high-level system design, API design details, and architectural diagrams are available in the [`/docs`](./docs) directory.

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
│   ├── generate_sql_builder.go   # Command to generate SQL builder files
│   └── ...                       # Other commands
├── db/                           # Database-related files
│   ├── migrations/               # Database migration files
│   │   ├── 202504221350_init.up.sql
│   │   ├── 202504221350_init.down.sql
│   │   └── ...
|── docs/                         # Documentation files
│   └── swagger.yaml              # API documentation in Swagger format
├── internal/                     # Internal application code
│   ├── api/                      # Delivery layer (HTTP handlers, routes, middleware)
│   │   ├── http/
│   │   │   ├── handler/          # HTTP handlers
│   │   │   │   ├── healthcheck/  # Health check handlers
│   │   │   │   └── ...
│   │   │   ├── middleware/       # HTTP middleware
│   │   │   └── route/            # HTTP route definitions
│   ├── config/                   # Configuration structs and loaders
│   ├── domain/                   # Domain entities and interfaces
│   │   ├── entity/               # Domain models (e.g., Concert, Reservation)
│   │   ├── repository/           # Repository interfaces
│   │   │   └── mocks/            # Mock implementations for repositories
│   ├── infra/                    # Infrastructure implementations
│   │   ├── db/                   # Database logic
│   │   │   ├── connection.go     # Database connection setup
│   │   │   ├── transactor.go     # Database transactor
│   │   │   ├── sql_execer.go     # Interface for executing SQL queries
│   │   │   ├── model_gen/        # Generated database models and table definitions
│   │   │   ├── healthcheck/      # Health check repository implementation
│   │   │   └── ...               # Other repository implementation
│   ├── server/                   # Server setup and initialization
│   ├── usecase/                  # Application business logic
│   │   ├── healthcheck/          # Health check use case
│   │   └── ...                   # Other use cases
│   └── util/                     # Utility functions
│       ├── httpresponse/         # HTTP response helpers
│       └── ...                   # Other utility functions
├── pkg/                          # Shared helper packages
├── env.yaml                      # Environment variables configuration (for local dev)
├── go.mod                        # Go module definition
├── go.sum                        # Go module dependencies
├── main.go                       # Main application entry point
├── Makefile                      # Makefile for common tasks
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

## 🛠️ Available Commands

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

Here’s a brief overview of the key targets in the provided Makefile:
- `install`: Installs all required Go tools (swag, mockgen, etc.).
- `gen-swag`: Generates Swagger documentation from code annotations.
- `gen-db`: Generates database models using go-jet.
- `gen-mock`: Generates mock files from code annotations.
- `vet`: Runs Go’s vet tool for static analysis.
- `fmt`: Formats the Go codebase.

## 🚀 Getting Started
### ⚡ Quick Start
Spin up the full local dev environment using Docker Compose:
```bash
docker-compose up --build
```