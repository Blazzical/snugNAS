package compose

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
)

// ServiceStatus is the trimmed-down view of `docker compose ps --format json`
// that the dashboard cares about.
type ServiceStatus struct {
	Service string `json:"service"`
	State   string `json:"state"`            // running, exited, restarting, ...
	Health  string `json:"health,omitempty"` // "", starting, healthy, unhealthy
}

// PS calls `docker compose ps --format json --all` and returns one
// ServiceStatus per container. The output of docker compose is newline-
// delimited JSON; one object per line.
func PS(ctx context.Context) ([]ServiceStatus, error) {
	p, err := ComposePath()
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", p, "ps", "--format", "json", "--all")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var statuses []ServiceStatus
	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Buffer(make([]byte, 1<<20), 1<<20) // 1MB per line (Labels can be long)
	for scanner.Scan() {
		var raw struct {
			Service string
			State   string
			Health  string
		}
		if err := json.Unmarshal(scanner.Bytes(), &raw); err != nil {
			// Skip malformed lines rather than failing the whole call.
			continue
		}
		statuses = append(statuses, ServiceStatus{
			Service: raw.Service,
			State:   raw.State,
			Health:  raw.Health,
		})
	}
	if err := scanner.Err(); err != nil {
		return statuses, err
	}
	return statuses, nil
}
