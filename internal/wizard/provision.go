package wizard

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Blazzical/snugNAS/internal/compose"
)

// provision runs the long-form work after the user clicks "Set Up":
// save config, render templates, ensure storage dir exists, `docker compose up -d`.
// All output from docker is streamed into s.logs for the polling page.
func (s *Server) provision() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	s.append("Saving config...")
	if err := s.cfg.Save(); err != nil {
		s.fail(fmt.Errorf("save config: %w", err))
		return
	}

	s.append("Creating storage directory at " + s.cfg.StorageRoot + "...")
	if err := os.MkdirAll(s.cfg.StorageRoot, 0o755); err != nil {
		s.fail(fmt.Errorf("create storage root: %w", err))
		return
	}

	s.append("Rendering compose stack...")
	if err := compose.Render(s.cfg); err != nil {
		s.fail(fmt.Errorf("render compose: %w", err))
		return
	}

	s.append("Pulling images and starting containers (this can take a while on first run)...")
	w := &lineWriter{mu: &s.mu, dst: &s.logs}
	if err := compose.Up(ctx, w, w); err != nil {
		s.fail(fmt.Errorf("docker compose up: %w", err))
		return
	}

	s.append("All services started. Opening dashboard.")
	s.mu.Lock()
	s.phase = PhaseDone
	s.mu.Unlock()
}

func (s *Server) append(line string) {
	s.mu.Lock()
	s.logs = append(s.logs, line)
	s.mu.Unlock()
}

func (s *Server) fail(err error) {
	s.mu.Lock()
	s.phase = PhaseFailed
	s.err = err
	s.logs = append(s.logs, "ERROR: "+err.Error())
	s.mu.Unlock()
}

// lineWriter is an io.Writer that splits incoming bytes by newline and pushes
// each non-empty line into a shared slice. Used to stream docker's progress
// output into the wizard's state.
type lineWriter struct {
	mu  *sync.Mutex
	dst *[]string
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, line := range strings.Split(string(p), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		*w.dst = append(*w.dst, line)
	}
	return len(p), nil
}
