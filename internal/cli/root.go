package cli

import (
	"github.com/spf13/cobra"
)

const version = "0.0.1-dev"

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "snugnas",
		Short: "snugNAS — CasaOS for Windows",
		Long: `snugNAS bundles Immich, Jellyfin, and FileBrowser behind a single
installer, first-run wizard, and dashboard. It runs on top of Docker Desktop
on Windows.`,
		Version: version,
	}

	root.AddCommand(
		installCmd(),
		wizardCmd(),
		upCmd(),
		downCmd(),
		statusCmd(),
		openCmd(),
		dashboardCmd(),
		trayCmd(),
		autostartCmd(),
	)
	return root
}
