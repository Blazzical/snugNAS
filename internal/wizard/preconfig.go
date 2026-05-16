package wizard

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Blazzical/snugNAS/internal/config"
)

// ensureJellyfinBaseURL writes (or patches) <StorageRoot>/jellyfin/config/network.xml
// so Jellyfin starts with BaseUrl=/jellyfin set. Lets Caddy route /jellyfin/* to
// Jellyfin without breaking its absolute-path asset URLs.
//
// First run: writes a minimal network.xml. Jellyfin's first start merges in
// defaults for the other fields. Re-run after Jellyfin has been started:
// replaces an empty <BaseUrl/> or <BaseUrl></BaseUrl> with our value. Leaves
// a non-empty existing BaseUrl alone (don't stomp user customisation).
func ensureJellyfinBaseURL(cfg *config.Config) error {
	// The bind mount target is /config in the container; Jellyfin nests its
	// XML config files under /config/config/ inside that. That's why this path
	// has "config" twice.
	dir := filepath.Join(cfg.StorageRoot, "jellyfin", "config", "config")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "network.xml")

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		minimal := `<?xml version="1.0" encoding="utf-8"?>
<NetworkConfiguration>
  <BaseUrl>/jellyfin</BaseUrl>
</NetworkConfiguration>
`
		return os.WriteFile(path, []byte(minimal), 0o644)
	}
	if err != nil {
		return err
	}

	emptyBaseURL := regexp.MustCompile(`<BaseUrl\s*/>|<BaseUrl></BaseUrl>`)
	if emptyBaseURL.Match(data) {
		data = emptyBaseURL.ReplaceAll(data, []byte("<BaseUrl>/jellyfin</BaseUrl>"))
		return os.WriteFile(path, data, 0o644)
	}
	// BaseUrl already non-empty — leave alone.
	return nil
}
