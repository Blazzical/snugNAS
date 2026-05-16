# Building snugNAS

## The Go binary

You need Go 1.22 or newer.

```powershell
go build -o snugnas.exe .\cmd\snugnas
```

Outputs `snugnas.exe` at the repo root. Optional: pass a release version with
ldflags:

```powershell
go build -ldflags "-s -w -X github.com/Blazzical/snugNAS/internal/cli.version=0.6.2" -o snugnas.exe .\cmd\snugnas
```

`-s -w` strips DWARF and symbol table to shave a few MB; `-X ...version=...`
stamps the version that `snugnas --version` reports.

## Cross-compiling

snugNAS is Windows-first but the wrapper compiles for Linux and macOS too:

```powershell
$env:GOOS = "linux";   $env:GOARCH = "amd64"; go build -o snugnas-linux-amd64    .\cmd\snugnas
$env:GOOS = "darwin";  $env:GOARCH = "amd64"; go build -o snugnas-darwin-amd64   .\cmd\snugnas
$env:GOOS = "darwin";  $env:GOARCH = "arm64"; go build -o snugnas-darwin-arm64   .\cmd\snugnas
Remove-Item Env:GOOS, Env:GOARCH
```

Note: the tray app and autostart command only have real implementations on
Windows. On other platforms `snugnas tray` returns a "platform not supported"
error and `snugnas autostart enable` returns `ErrUnsupported`.

## Windows installer (NSIS)

The Windows `snugnas-setup.exe` is built with [NSIS](https://nsis.sourceforge.io/)
3.x. Install NSIS (winget will do):

```powershell
winget install --id NSIS.NSIS --silent
```

Then, from the repo root:

```powershell
# 1. Build the binary
go build -o snugnas.exe .\cmd\snugnas

# 2. Compile the installer (output: installer\snugnas-setup.exe)
makensis installer\snugnas.nsi
```

The script bundles `snugnas.exe`, `LICENSE`, and `README.md`. It writes a
Start-menu entry and an Add/Remove Programs record. Uninstalling keeps the
user's config and storage in place — that data is preserved deliberately.

The installer requests admin elevation. It does **not** enable autostart
from inside the installer because Run keys are per-user (HKCU) and an admin
installer runs in a different user context. Tell users to run
`snugnas autostart enable` from an ordinary PowerShell to opt in.

## Releasing

For now: tag the commit (`git tag v0.6.2 && git push --tags`), build the
binary with the matching version, build the installer, and attach the
installer plus the raw `snugnas.exe` to a GitHub release.

CI-based release builds aren't wired up yet — that's on the "v1.0 — Ship"
roadmap item.
