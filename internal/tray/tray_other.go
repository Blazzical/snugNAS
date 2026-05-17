//go:build !windows

package tray

import (
	"errors"

	"github.com/Blazzical/snugNAS/internal/config"
)

// Run returns an error on non-Windows platforms — the tray uses getlantern/systray
// which needs Cgo + platform UI libraries (GTK on Linux, Cocoa on macOS) and we
// haven't wired those up. snugnas dashboard provides the equivalent foreground
// experience on those platforms.
func Run(cfg *config.Config) error {
	return errors.New("snugnas tray is Windows-only — use `snugnas dashboard` instead")
}
