package cmd

import (
	"fmt"
	"log"

	"ticket-reservation/internal/config"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	jetPostgres "github.com/go-jet/jet/v2/postgres"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var generateSQLBuilderCmd = &cobra.Command{
	Use:   "generate-db",
	Short: "Generate SQL builder and model files using go-jet.",
	Long: `Connects to the database and generates type-safe SQL builder and model files using go-jet.

This command will:
- Connect to the database schema.
- Generate table builders, model structs, and view builders (if any).
- Apply custom field mappings (e.g., map payments.amount to decimal.Decimal).

Example:
	generate-db
	generate-db --schema public --dir ./internal/infra/db/model_gen

Documentation:
https://github.com/go-jet/jet/wiki/Generator
`,

	RunE: runGenerateSQLBuilderCmd,
}

func runGenerateSQLBuilderCmd(cmd *cobra.Command, args []string) error {
	cfg := config.MustConfigure()
	schema, _ := cmd.Flags().GetString("schema")
	dir, _ := cmd.Flags().GetString("dir")

	log.Println("Starting SQL builder code generation...")

	if err := generateSQLBuilder(cfg.DB.URL, schema, dir); err != nil {
		return fmt.Errorf("failed to generate SQL builder files: %w", err)
	}

	log.Println("SQL builder generation completed successfully.")
	return nil
}

func generateSQLBuilder(databaseURL string, schema string, dir string) error {
	// Generate SQL builder files using go-jet
	// documentation: https://github.com/go-jet/jet/wiki/Generator
	err := postgres.GenerateDSN(
		databaseURL, // database connection
		schema,      // schema name
		dir,         // output directory
		template.Default(jetPostgres.Dialect).
			UseSchema(func(schema metadata.Schema) template.Schema {
				return template.DefaultSchema(schema).
					UseModel(template.DefaultModel().
						UseTable(func(table metadata.Table) template.TableModel {
							return template.DefaultTableModel(table).
								UseField(func(column metadata.Column) template.TableModelField {
									field := template.DefaultTableModelField(column)

									// Customize the field type for specific columns
									// For example, map payments.amount to *decimal.Decimal
									if schema.Name == "public" && table.Name == "payments" && column.Name == "amount" {
										field.Type = template.NewType(&decimal.Decimal{})
									}

									// Set the `db` struct tag to match Jet's fully-qualified column alias (e.g., `concerts.name`).
									// This is required for correct mapping when using sqlx, which matches result columns to struct tags.
									// For example, Jet generates: SELECT concerts.name AS "concerts.name"
									// So we use: `db:"concerts.name"`
									//
									// Alternatively, you can remove the table alias in the generated SQL by calling .As("") on the column:
									// e.g., concerts.Name.As("") will generate: SELECT name AS "name"
									// In that case, the struct tag can simply be `db:"name"`
									//
									// Here, we keep the alias for clarity and consistency across queries.
									field = field.UseTags(
										fmt.Sprintf(`db:"%s.%s"`, table.Name, column.Name),
									)

									return field
								})
						}),
					)
			}),
	)
	return err
}

func init() {
	generateSQLBuilderCmd.Flags().String("schema", "public", "Schema name to generate files for")
	generateSQLBuilderCmd.Flags().String("dir", "./internal/infra/db/model_gen", "Directory to save generated files")
}
