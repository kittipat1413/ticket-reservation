package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Starts the application",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute initializes the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(
		printConfigCmd,
		generateSQLBuilderCmd,
		newMigrationCmd,
		migrateCmd,
		serveCmd,
	)
}
