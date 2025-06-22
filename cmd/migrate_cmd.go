package cmd

import (
	"fmt"
	"log"
	"ticket-reservation/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run DB migrations.",
	Long: `Run database schema migrations using the golang-migrate library.
	
This command supports both 'up' and 'down' directions:
- 'up': applies all available migration scripts (e.g., .up.sql)
- 'down': rolls back all applied migrations (e.g., .down.sql)

Migration files must follow the naming convention:
	<version>_<description>.up.sql   - for applying schema changes
	<version>_<description>.down.sql - for rolling them back

Example:
	migrate --action up
	migrate --action down
	migrate --action up --dir ./custom/path/to/migrations

The default migration directory is './data/migrations'.
Database connection is read from config.DatabaseURL().
`,

	RunE: runMigrateCmd,
}

func runMigrateCmd(cmd *cobra.Command, args []string) error {
	cfg := config.MustConfigure()

	dir, _ := cmd.Flags().GetString("dir")
	action, _ := cmd.Flags().GetString("action")

	m, err := migrate.New(
		"file://"+dir,
		cfg.DB.URL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	switch action {
	case "up":
		log.Println("Running migration: UP")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run up migration: %w", err)
		}
		log.Println("All migrations applied successfully.")
	case "down":
		log.Println("Running migration: DOWN")
		if err := m.Steps(-1); err != nil {
			return fmt.Errorf("failed to run down migration: %w", err)
		}
		log.Println("Rolled back one migration successfully.")
	default:
		return fmt.Errorf("invalid action: %s (use up or down)", action)
	}
	return nil
}

func init() {
	migrateCmd.Flags().String("dir", "./db/migrations", "Path to the migration files")
	migrateCmd.Flags().String("action", "up", "Migration action: up or down")
}
