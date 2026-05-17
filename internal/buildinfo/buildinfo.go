// Package buildinfo holds compile-time metadata for snugNAS — currently just
// the version string. Overridable at build time via:
//
//	go build -ldflags "-X github.com/Blazzical/snugNAS/internal/buildinfo.Version=0.7.5" ./cmd/snugnas
//
// Kept in its own tiny package so cli, landing, updater, and tests can all
// reference the same value without importing each other.
package buildinfo

// Version is the snugNAS release version. The default value ships in dev
// builds; release builds override it via -ldflags.
var Version = "0.0.1-dev"
