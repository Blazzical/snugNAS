package wizard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// validateStorage returns nil if storageRoot is usable: the path is non-empty,
// reasonably-formed, and either exists writable or can be created.
//
// We don't validate hostname format here — Caddy and mDNS are lenient enough
// that anything non-empty works for v0.8.
func validateStorage(storageRoot string) error {
	if storageRoot == "" {
		return fmt.Errorf("storage folder is required")
	}
	if strings.TrimSpace(storageRoot) != storageRoot {
		return fmt.Errorf("storage folder has surrounding whitespace")
	}

	// If the path exists, make sure it's a directory and we can write a probe file.
	if info, err := os.Stat(storageRoot); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%s exists but isn't a directory", storageRoot)
		}
		return writeProbe(storageRoot)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("can't read %s: %w", storageRoot, err)
	}

	// Path doesn't exist — verify we can create it.
	if err := os.MkdirAll(storageRoot, 0o755); err != nil {
		return fmt.Errorf("can't create %s: %w", storageRoot, err)
	}
	return writeProbe(storageRoot)
}

// writeProbe creates and removes a small file under dir to confirm we have
// write access. The file lives a few milliseconds at most.
func writeProbe(dir string) error {
	f, err := os.CreateTemp(dir, ".snugnas-write-probe-*")
	if err != nil {
		return fmt.Errorf("can't write to %s: %w", dir, err)
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	_ = filepath.Clean(name) // silence unused on some platforms
	return nil
}
