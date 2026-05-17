package wizard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateStorage(t *testing.T) {
	tmp := t.TempDir()

	// Existing writable dir.
	if err := validateStorage(tmp); err != nil {
		t.Errorf("existing writable dir: %v", err)
	}

	// Non-existing subdir that we can create.
	sub := filepath.Join(tmp, "new-subdir")
	if err := validateStorage(sub); err != nil {
		t.Errorf("creatable subdir: %v", err)
	}
	if _, err := os.Stat(sub); err != nil {
		t.Errorf("validate didn't create the subdir: %v", err)
	}

	// Empty path.
	if err := validateStorage(""); err == nil {
		t.Error("expected error for empty path")
	}

	// Leading whitespace.
	if err := validateStorage("  /tmp/foo"); err == nil {
		t.Error("expected error for whitespace-padded path")
	}

	// Existing file (not a directory).
	filePath := filepath.Join(tmp, "not-a-dir")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	err := validateStorage(filePath)
	if err == nil || !strings.Contains(err.Error(), "isn't a directory") {
		t.Errorf("expected 'isn't a directory' error, got %v", err)
	}
}
