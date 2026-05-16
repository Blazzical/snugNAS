package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func wizardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wizard",
		Short: "Run the first-run web wizard (opens http://localhost:7777)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO(wizard): start internal/wizard HTTP server, open browser,
			// collect storage dir + admin creds + hostname, render compose
			// from compose/docker-compose.yml.tmpl, write .env, run `docker compose up -d`.
			fmt.Println("wizard: not yet implemented")
			return nil
		},
	}
}
