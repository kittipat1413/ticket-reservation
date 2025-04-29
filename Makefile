# Makefile for Go project
# This Makefile provides a set of commands to manage the Go project, including building, testing, and generating database models.


# gen-db target generates database models using go-jet.
gen-db:
	@echo "Generating database models... 🚧"
	@go run main.go generate-db
	@echo "Database models generated successfully. ✅"

# gen-mock target generates mock files using mockgen.
gen-mock:
	@echo "Generating mock files... 🚧"
	@go generate ./...
	@echo "Mock files generated successfully. ✅"