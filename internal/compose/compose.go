// Package compose renders the docker-compose.yml and Caddyfile from templates
// using a *config.Config, and provides up/down/status helpers around the
// `docker compose` CLI.
package compose

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/Blazzical/snugNAS/internal/config"
)

//go:embed templates
var templates embed.FS

// Render writes docker-compose.yml and Caddyfile into the user's runtime dir.
func Render(cfg *config.Config) error {
	rt, err := config.RuntimeDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(rt, 0o755); err != nil {
		return err
	}
	if err := renderTemplate(cfg, "templates/docker-compose.yml.tmpl", filepath.Join(rt, "docker-compose.yml")); err != nil {
		return fmt.Errorf("render compose: %w", err)
	}
	if err := renderTemplate(cfg, "templates/Caddyfile.tmpl", filepath.Join(rt, "Caddyfile")); err != nil {
		return fmt.Errorf("render Caddyfile: %w", err)
	}
	return nil
}

func renderTemplate(cfg *config.Config, src, dst string) error {
	raw, err := templates.ReadFile(src)
	if err != nil {
		return err
	}
	t, err := template.New(filepath.Base(src)).Parse(string(raw))
	if err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, cfg)
}

// ComposePath returns the absolute path to the rendered docker-compose.yml.
func ComposePath() (string, error) {
	rt, err := config.RuntimeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(rt, "docker-compose.yml"), nil
}

// Up runs `docker compose up -d`.
func Up(ctx context.Context, stdout, stderr io.Writer) error {
	return run(ctx, stdout, stderr, "up", "-d")
}

// Down runs `docker compose down`.
func Down(ctx context.Context, stdout, stderr io.Writer) error {
	return run(ctx, stdout, stderr, "down")
}

// Status runs `docker compose ps`.
func Status(ctx context.Context, stdout, stderr io.Writer) error {
	return run(ctx, stdout, stderr, "ps")
}

func run(ctx context.Context, stdout, stderr io.Writer, args ...string) error {
	p, err := ComposePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(p); err != nil {
		return fmt.Errorf("no rendered compose file at %s — run `snugnas wizard` first", p)
	}
	full := append([]string{"compose", "-f", p}, args...)
	cmd := exec.CommandContext(ctx, "docker", full...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
