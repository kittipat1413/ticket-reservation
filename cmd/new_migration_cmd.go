package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
)

var newMigrationCmd = &cobra.Command{
	Use:   "new-migration [name]",
	Short: "Creates a new migration file with timestamp and the given name",
	Args:  cobra.ExactArgs(1),
	RunE:  runNewMigrationCmd,
}

func runNewMigrationCmd(cmd *cobra.Command, args []string) error {
	name := flect.Underscore(args[0])
	timestamp := time.Now().Format("200601021504")
	filename := fmt.Sprintf("%s_%s", timestamp, name)

	upname := filename + ".up.sql"
	downname := filename + ".down.sql"

	migrationDir := "./db/migrations"
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations dir: %w", err)
	}

	uppath := filepath.Join(migrationDir, upname)
	downpath := filepath.Join(migrationDir, downname)

	fmt.Printf("Creating migration files:\n- %s\n- %s\n", uppath, downpath)

	if err := os.WriteFile(uppath, []byte("-- "+upname+"\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(downpath, []byte("-- "+downname+"\n"), 0644); err != nil {
		return err
	}

	return nil
}
