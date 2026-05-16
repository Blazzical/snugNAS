package cli

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Blazzical/snugNAS/internal/config"
	"github.com/Blazzical/snugNAS/internal/wizard"
	"github.com/spf13/cobra"
)

func wizardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wizard",
		Short: "Run the first-run web wizard (storage path, hostname, admin email)",
		Long: `Starts an HTTP server on 127.0.0.1:7777 and opens a browser at it.
The page collects storage path, hostname, and admin email; submission writes
config.toml and runs ` + "`docker compose up -d`" + ` in the background.
The page polls /api/state for progress and redirects to the dashboard on success.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				if !config.IsNotConfigured(err) {
					return err
				}
				cfg = config.Default()
			}
			port := cfg.DashboardPort
			url := fmt.Sprintf("http://127.0.0.1:%d/", port)

			fmt.Fprintln(cmd.OutOrStdout(), "Wizard listening at", url)

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			go func() {
				time.Sleep(500 * time.Millisecond)
				_ = openBrowser(url)
			}()

			return wizard.Serve(ctx, port)
		},
	}
}
