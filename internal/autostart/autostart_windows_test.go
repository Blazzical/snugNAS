//go:build windows

package autostart

import (
	"os"
	"testing"
)

// TestRoundTrip touches the user's actual HKCU\Run key. To avoid surprising
// users running tests, set SNUGNAS_TEST_REGISTRY=1 to opt in.
func TestRoundTrip(t *testing.T) {
	if os.Getenv("SNUGNAS_TEST_REGISTRY") != "1" {
		t.Skip("set SNUGNAS_TEST_REGISTRY=1 to enable (touches HKCU\\Run)")
	}

	// Save current state so we restore it cleanly.
	wasEnabled, originalCmd, _ := Status()
	t.Cleanup(func() {
		_ = Disable()
		if wasEnabled {
			// Re-enable, but with the originally registered command.
			// Best-effort — Enable() will overwrite with our exe path.
			_, _ = Enable()
			_ = originalCmd // documented for the reader
		}
	})

	_ = Disable()
	if e, _, _ := Status(); e {
		t.Fatal("Status reports enabled immediately after Disable")
	}

	line, err := Enable()
	if err != nil {
		t.Fatalf("Enable: %v", err)
	}
	if line == "" {
		t.Error("Enable returned empty command")
	}

	e, got, err := Status()
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !e {
		t.Error("Status reports disabled after Enable")
	}
	if got != line {
		t.Errorf("Status command = %q, want %q", got, line)
	}
}
