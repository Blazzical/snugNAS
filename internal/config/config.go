// Package config reads and writes the snugNAS config file (typically at
// %APPDATA%\snugNAS\config.toml on Windows). It is the source of truth for
// storage paths, hostname, ports, and admin identity.
package config

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// Config is the on-disk representation of the user's snugNAS install.
// Field names match the template variables in internal/compose/templates.
type Config struct {
	StorageRoot      string `toml:"storage_root"`
	Hostname         string `toml:"hostname"`
	DashboardPort    int    `toml:"dashboard_port"`
	AdminEmail       string `toml:"admin_email"`
	PostgresPassword string `toml:"postgres_password"`
}

// Default returns a Config with sensible defaults but no storage or admin set.
// The wizard fills in the empty fields.
func Default() *Config {
	pw, _ := randomHex(24)
	return &Config{
		Hostname:         "snugnas.local",
		DashboardPort:    7777,
		PostgresPassword: pw,
	}
}

// Dir is the directory holding config.toml and the rendered runtime files.
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "snugNAS"), nil
}

// Path returns the absolute path to config.toml.
func Path() (string, error) {
	d, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "config.toml"), nil
}

// RuntimeDir is where rendered docker-compose.yml and Caddyfile are written.
func RuntimeDir() (string, error) {
	d, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "runtime"), nil
}

// Load reads config.toml. Returns os.ErrNotExist (wrapped) if no config exists yet.
func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimPrefix(data, utf8BOM)
	var c Config
	if err := toml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", p, err)
	}
	return &c, nil
}

// Save writes the config back to disk, creating the directory if needed.
func (c *Config) Save() error {
	d, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	p, err := Path()
	if err != nil {
		return err
	}
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

// IsNotConfigured reports whether err means "no config has been written yet".
func IsNotConfigured(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

// PosixStorageRoot returns StorageRoot with backslashes converted to forward
// slashes. Used by templates so YAML doesn't choke on Windows paths and Docker
// accepts the bind-mount source consistently.
//
// We use strings.ReplaceAll (not filepath.ToSlash) on purpose: on Linux,
// filepath.ToSlash treats `\` as a literal character, not a separator, so it
// would leave Windows paths intact. We always want backslashes converted
// regardless of the host OS, because the generated docker-compose.yml is
// consumed by Docker (which wants posix paths) not the host filesystem.
func (c *Config) PosixStorageRoot() string {
	return strings.ReplaceAll(c.StorageRoot, `\`, "/")
}

func randomHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
