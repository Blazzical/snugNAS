package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/Blazzical/snugNAS/internal/compose"
	"github.com/Blazzical/snugNAS/internal/config"
	"github.com/spf13/cobra"
)

func wizardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wizard",
		Short: "Write a starter config and render the compose stack",
		Long: `Writes a default config.toml if none exists, then renders the
docker-compose.yml and Caddyfile into the runtime directory. The web-based
interactive wizard ships in v0.3 — for now, edit config.toml by hand to set
storage_root and admin_email.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				if !config.IsNotConfigured(err) {
					return err
				}
				cfg = config.Default()
				if err := cfg.Save(); err != nil {
					return err
				}
				p, _ := config.Path()
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote default config to %s\n", p)
			}

			if cfg.StorageRoot == "" {
				p, _ := config.Path()
				return errors.New("set storage_root in " + p + " before running `snugnas up`")
			}

			if err := os.MkdirAll(cfg.StorageRoot, 0o755); err != nil {
				return fmt.Errorf("create storage root: %w", err)
			}
			if err := compose.Render(cfg); err != nil {
				return err
			}
			rt, _ := config.RuntimeDir()
			fmt.Fprintf(cmd.OutOrStdout(), "Rendered compose stack to %s\n", rt)
			fmt.Fprintln(cmd.OutOrStdout(), "Run `snugnas up` to start the stack.")
			return nil
		},
	}
}
