// Package installer detects and helps the user install Docker Desktop + WSL2
// on Windows. It deliberately does not silently install — see docs/architecture.md
// for the rationale (Docker Desktop's licensing requires explicit acceptance).
package installer

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// Status is the result of probing the local machine for snugNAS prerequisites.
type Status struct {
	DockerInstalled bool
	DockerRunning   bool
	WSL2Default     bool
	Notes           []string
}

// OK reports whether the host is ready to run snugnas up.
func (s Status) OK() bool {
	return s.DockerInstalled && s.DockerRunning && s.WSL2Default
}

// Check runs all detection probes and returns a Status. Never returns an error
// — failures show up as false fields with human-readable Notes.
func Check(ctx context.Context) Status {
	var s Status

	if _, err := exec.LookPath("docker"); err == nil {
		s.DockerInstalled = true
	} else {
		s.Notes = append(s.Notes,
			"Docker is not on PATH. Install Docker Desktop from https://www.docker.com/products/docker-desktop/")
	}

	if s.DockerInstalled {
		out, err := exec.CommandContext(ctx, "docker", "version", "--format", "{{.Server.Version}}").Output()
		if err == nil && len(bytes.TrimSpace(out)) > 0 {
			s.DockerRunning = true
		} else {
			s.Notes = append(s.Notes,
				"Docker is installed but the daemon is not reachable. Launch Docker Desktop and wait for the whale icon to stop animating.")
		}
	}

	if out, err := exec.CommandContext(ctx, "wsl", "--status").CombinedOutput(); err == nil {
		// `wsl --status` writes UTF-16 LE on Windows 10. Strip null bytes for a
		// rough-and-ready ASCII view that's good enough for substring checks.
		text := strings.ReplaceAll(string(out), "\x00", "")
		if strings.Contains(text, "Default Version: 2") {
			s.WSL2Default = true
		}
	}
	if !s.WSL2Default {
		s.Notes = append(s.Notes,
			"WSL2 doesn't appear to be the default. Run `wsl --update --web-download` from an admin PowerShell, then `wsl --set-default-version 2`.")
	}

	return s
}
