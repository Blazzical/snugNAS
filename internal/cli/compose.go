package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"github.com/Blazzical/snugNAS/internal/compose"
	"github.com/Blazzical/snugNAS/internal/config"
	"github.com/Blazzical/snugNAS/internal/landing"
	"github.com/Blazzical/snugNAS/internal/mdns"
	"github.com/spf13/cobra"
)

func upCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start the snugNAS stack (docker compose up -d)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if err := compose.Render(cfg); err != nil {
				return err
			}
			return compose.Up(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

func downCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop the snugNAS stack (docker compose down)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := loadConfig(); err != nil {
				return err
			}
			return compose.Down(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the running status of snugNAS services",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := loadConfig(); err != nil {
				return err
			}
			return compose.Status(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open the snugNAS dashboard in the default browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			url := fmt.Sprintf("http://127.0.0.1:%d/", cfg.DashboardPort)
			fmt.Fprintln(cmd.OutOrStdout(), "Opening", url)
			return openBrowser(url)
		},
	}
}

func dashboardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dashboard",
		Short: "Run the dashboard HTTP server in the foreground",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()
			fmt.Fprintf(cmd.OutOrStdout(), "Dashboard listening at http://127.0.0.1:%d/\n", cfg.DashboardPort)

			pub, err := mdns.Publish(cfg.Hostname, cfg.DashboardPort)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning: mDNS publish failed:", err)
			} else {
				defer pub.Shutdown()
				fmt.Fprintf(cmd.OutOrStdout(), "Advertising %s on the LAN (IPs: %v)\n", pub.Hostname, pub.IPs)
			}

			return landing.Serve(ctx, cfg)
		},
	}
}

func trayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tray",
		Short: "Run the Windows tray app (planned for v0.5)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("tray app not yet implemented (planned for v0.5)")
		},
	}
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		if config.IsNotConfigured(err) {
			p, _ := config.Path()
			return nil, fmt.Errorf("no config found at %s — run `snugnas wizard` first", p)
		}
		return nil, err
	}
	if cfg.StorageRoot == "" {
		p, _ := config.Path()
		return nil, fmt.Errorf("storage_root is empty in %s — edit it and try again", p)
	}
	return cfg, nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

