package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func upCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start the snugNAS stack (docker compose up -d)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: shell out to `docker compose -f runtime/docker-compose.yml up -d`.
			fmt.Println("up: not yet implemented")
			return nil
		},
	}
}

func downCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop the snugNAS stack (docker compose down)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: shell out to `docker compose -f runtime/docker-compose.yml down`.
			fmt.Println("down: not yet implemented")
			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the running status of snugNAS services",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: parse `docker compose ps --format json` and render a table.
			fmt.Println("status: not yet implemented")
			return nil
		},
	}
}

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open the snugNAS dashboard in the default browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: cross-platform browser open of http://snugnas.local
			fmt.Println("open: not yet implemented")
			return nil
		},
	}
}

func trayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tray",
		Short: "Run the Windows tray app (background, foreground on Linux/macOS)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: getlantern/systray loop. See internal/tray.
			fmt.Println("tray: not yet implemented")
			return nil
		},
	}
}
