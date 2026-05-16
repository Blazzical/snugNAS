// Package autostart enables snugNAS to launch automatically when the user
// logs in. On Windows this is done by adding an entry under
// HKCU\Software\Microsoft\Windows\CurrentVersion\Run that runs `snugnas tray`.
// On other platforms the functions return ErrUnsupported.
package autostart

import "errors"

// ErrUnsupported is returned by autostart functions on platforms where the
// feature isn't implemented.
var ErrUnsupported = errors.New("autostart is only supported on Windows")

// valueName is the registry value name used for the autostart entry. Exported
// so docs/tests can reference the canonical name.
const ValueName = "snugNAS"

// runKey is the subkey under HKCU where Run entries live.
const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`
