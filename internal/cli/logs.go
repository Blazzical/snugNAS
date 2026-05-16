package cli

import (
	"github.com/Blazzical/snugNAS/internal/compose"
	"github.com/spf13/cobra"
)

func logsCmd() *cobra.Command {
	var follow bool
	var tail int
	cmd := &cobra.Command{
		Use:   "logs [service]",
		Short: "Show logs from snugNAS services (wrapper around docker compose logs)",
		Long: `Without a service argument, shows merged logs for every service.
With a service name (caddy, immich-server, jellyfin, filebrowser, ...) shows
logs for just that one.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := loadConfig(); err != nil {
				return err
			}
			service := ""
			if len(args) > 0 {
				service = args[0]
			}
			return compose.Logs(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), follow, tail, service)
		},
	}
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Tail logs continuously (Ctrl+C to stop)")
	cmd.Flags().IntVar(&tail, "tail", 200, "Show only the last N lines per service")
	return cmd
}

func restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart [service]",
		Short: "Restart the snugNAS stack (or one specific service)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := loadConfig(); err != nil {
				return err
			}
			service := ""
			if len(args) > 0 {
				service = args[0]
			}
			return compose.Restart(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), service)
		},
	}
}
