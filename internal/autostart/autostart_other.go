//go:build !windows

package autostart

// Enable returns ErrUnsupported on non-Windows platforms.
func Enable() (string, error) { return "", ErrUnsupported }

// Disable returns ErrUnsupported on non-Windows platforms.
func Disable() error { return ErrUnsupported }

// Status returns (false, "", nil) on non-Windows platforms.
func Status() (bool, string, error) { return false, "", nil }
