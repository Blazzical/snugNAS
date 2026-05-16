package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/Blazzical/snugNAS/internal/compose"
	"github.com/Blazzical/snugNAS/internal/config"
	"github.com/spf13/cobra"
)

func resetCmd() *cobra.Command {
	var yes, alsoStorage bool
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Wipe snugNAS config + rendered runtime (dry-run by default)",
		Long: `Removes the snugNAS config.toml and the rendered runtime files.
With --storage, also removes the storage root that holds your photos,
videos, and service databases.

Always stops the running stack before deleting anything. Without --yes
this command is a dry-run that lists what would be removed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil && !config.IsNotConfigured(err) {
				return err
			}

			cpath, _ := config.Path()
			rdir, _ := config.RuntimeDir()

			targets := []string{}
			if _, statErr := os.Stat(cpath); statErr == nil {
				targets = append(targets, cpath)
			}
			if _, statErr := os.Stat(rdir); statErr == nil {
				targets = append(targets, rdir)
			}
			if alsoStorage && cfg != nil && cfg.StorageRoot != "" {
				if _, statErr := os.Stat(cfg.StorageRoot); statErr == nil {
					targets = append(targets, cfg.StorageRoot+"  (ALL DATA)")
				}
			}

			out := cmd.OutOrStdout()
			if len(targets) == 0 {
				fmt.Fprintln(out, "Nothing to remove. snugNAS doesn't appear to be configured.")
				return nil
			}

			fmt.Fprintln(out, "Will remove:")
			for _, t := range targets {
				fmt.Fprintln(out, "  -", t)
			}

			if !yes {
				fmt.Fprintln(out, "\nDry-run. Pass --yes to actually delete.")
				if !alsoStorage && cfg != nil && cfg.StorageRoot != "" {
					fmt.Fprintf(out, "Note: --storage is NOT set, so %s and its data will be kept.\n", cfg.StorageRoot)
				}
				return nil
			}

			// Try to stop the stack first so we don't leave orphaned containers
			// pointing at a now-deleted compose file. Best-effort: ignore errors.
			if cfg != nil {
				fmt.Fprintln(out, "Stopping stack...")
				if err := compose.Down(cmd.Context(), io.Discard, cmd.ErrOrStderr()); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), "warning: docker compose down:", err)
				}
			}

			if _, err := os.Stat(cpath); err == nil {
				if err := os.Remove(cpath); err != nil {
					return fmt.Errorf("remove config: %w", err)
				}
			}
			if _, err := os.Stat(rdir); err == nil {
				if err := os.RemoveAll(rdir); err != nil {
					return fmt.Errorf("remove runtime: %w", err)
				}
			}
			if alsoStorage && cfg != nil && cfg.StorageRoot != "" {
				if _, err := os.Stat(cfg.StorageRoot); err == nil {
					if err := os.RemoveAll(cfg.StorageRoot); err != nil {
						return fmt.Errorf("remove storage: %w", err)
					}
				}
			}

			fmt.Fprintln(out, "\nDone. Run `snugnas wizard` to set up again.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Actually delete (default is dry-run)")
	cmd.Flags().BoolVar(&alsoStorage, "storage", false, "Also remove the storage root (irreversible — wipes photos, videos, DBs)")
	return cmd
}
