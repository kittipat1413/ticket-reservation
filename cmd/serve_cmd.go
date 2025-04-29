package cmd

import (
	"ticket-reservation/internal/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Starts application server",
	RunE:    runServeCmd,
	Aliases: []string{"s", "server"},
}

func runServeCmd(cmd *cobra.Command, args []string) error {
	return server.New().Start()
}
