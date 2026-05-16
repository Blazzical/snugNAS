//go:build windows

package autostart

import (
	"errors"
	"os"

	"golang.org/x/sys/windows/registry"
)

// Enable adds (or overwrites) the autostart entry. Returns the command string
// that was written so callers can show the user what will run on login.
func Enable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	cmd := `"` + exe + `" tray`

	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	if err := k.SetStringValue(ValueName, cmd); err != nil {
		return "", err
	}
	return cmd, nil
}

// Disable removes the autostart entry. Idempotent — returns nil if no entry
// exists.
func Disable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return nil
		}
		return err
	}
	defer k.Close()
	err = k.DeleteValue(ValueName)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	return err
}

// Status reports whether autostart is enabled and the registered command.
func Status() (enabled bool, command string, err error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, "", nil
		}
		return false, "", err
	}
	defer k.Close()
	cmd, _, err := k.GetStringValue(ValueName)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, "", nil
		}
		return false, "", err
	}
	return true, cmd, nil
}
