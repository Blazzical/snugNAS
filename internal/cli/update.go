package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Blazzical/snugNAS/internal/updater"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Check GitHub for a newer snugNAS release",
		Long: `Queries https://api.github.com/repos/Blazzical/snugNAS/releases/latest
and reports whether a newer version is available. Does not download or
install anything — prints the release URL and assets for you to grab manually.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := updater.LatestRelease(cmd.Context())
			if err != nil {
				if errors.Is(err, updater.ErrNoReleases) {
					fmt.Fprintln(cmd.OutOrStdout(), "No releases published yet. You're on", version)
					return nil
				}
				return err
			}

			current := strings.TrimPrefix(version, "v")
			latest := strings.TrimPrefix(r.TagName, "v")

			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "Current version:", current)
			fmt.Fprintln(out, "Latest release: ", latest, "—", r.Name)
			fmt.Fprintln(out, "Released:       ", r.PublishedAt)
			fmt.Fprintln(out, "URL:            ", r.HTMLURL)

			if current == latest {
				fmt.Fprintln(out, "\nYou're up to date.")
				return nil
			}

			fmt.Fprintln(out, "\nUpdate available. Download:")
			if len(r.Assets) == 0 {
				fmt.Fprintln(out, "  (no assets — source-only release)")
				return nil
			}
			for _, a := range r.Assets {
				fmt.Fprintf(out, "  - %s\n    %s\n", a.Name, a.BrowserDownloadURL)
			}
			return nil
		},
	}
}
