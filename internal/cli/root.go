package cli

import (
	"github.com/Blazzical/snugNAS/internal/buildinfo"
	"github.com/spf13/cobra"
)

// version is sourced from internal/buildinfo so the ldflags override has a
// single canonical target (see internal/buildinfo/buildinfo.go).
var version = buildinfo.Version

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
		restartCmd(),
		statusCmd(),
		logsCmd(),
		openCmd(),
		dashboardCmd(),
		trayCmd(),
		autostartCmd(),
		resetCmd(),
		updateCmd(),
	)
	return root
}
