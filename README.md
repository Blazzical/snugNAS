# snugNAS

[![CI](https://github.com/Blazzical/snugNAS/actions/workflows/ci.yml/badge.svg)](https://github.com/Blazzical/snugNAS/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**CasaOS for Windows.** A one-installer NAS stack that runs on top of your existing Windows machine — keep Windows, get a Synology-quality experience.

snugNAS bundles best-in-class open-source services behind a single installer, a non-technical first-run wizard, and a web dashboard with QR codes that point your phone and smart TV at the right app.

## What you get

| Service | What it's for | Phone / TV |
|---|---|---|
| [Immich](https://immich.app/) | Photo & video backup (Google Photos replacement) | Native iOS / Android apps with background auto-upload |
| [Jellyfin](https://jellyfin.org/) | Movies, TV, music | Android, iOS, AndroidTV, AppleTV, Roku, Fire TV, Chromecast, Kodi, DLNA |
| [FileBrowser](https://filebrowser.org/) | Everything else (docs, downloads, archives) | Mobile-friendly web UI |
| [Caddy](https://caddyserver.com/) | Reverse proxy, single URL, auto TLS | (internal) |

All four services point at the same storage root and coexist safely.

## Status

**v0.7 — works end to end on a developer machine.** Full stack pulls and starts via the wizard, dashboard shows live service health, mDNS exposes `snugnas.local` to phones on the LAN. NSIS installer scaffolded but not yet release-built. See [`docs/roadmap.md`](docs/roadmap.md) for what each version shipped.

## Quick start (today, from a Git clone)

You need [Go 1.22+](https://go.dev/dl/), [Docker Desktop](https://www.docker.com/products/docker-desktop/), and Git. On Windows 10, also run `wsl --update --web-download` from an admin PowerShell once.

```powershell
git clone https://github.com/Blazzical/snugNAS.git
cd snugNAS
go build .\cmd\snugnas

# 1. Verify prerequisites
.\snugnas.exe install

# 2. First-run wizard (opens browser at http://127.0.0.1:7777)
.\snugnas.exe wizard

# After provisioning completes, the wizard redirects to the dashboard.
# Optional: run as a tray app so it survives the terminal closing
.\snugnas.exe tray

# Optional: launch automatically at login
.\snugnas.exe autostart enable
```

The dashboard lives at `http://localhost:8080/` (locally) or `http://snugnas.local:8080/` (from a phone on the same Wi-Fi — needs mDNS, which iOS / modern Android / macOS all do natively).

## How it works

```
snugNAS (Windows host)
 ├─ Docker Desktop / WSL2          ← runtime
 │   ├─ Immich (+ postgres, redis, ML)
 │   ├─ Jellyfin (BaseUrl=/jellyfin)
 │   ├─ FileBrowser
 │   └─ Caddy                       ← single ingress at :8080
 └─ snugnas.exe (Go binary)
     ├─ install       → checks Docker Desktop + WSL2
     ├─ wizard        → web-based first-run setup (storage / hostname / admin)
     ├─ up / down     → wraps `docker compose up -d` / down
     ├─ restart       → bounces one service or all
     ├─ status        → docker compose ps
     ├─ logs          → docker compose logs --tail=...
     ├─ open          → opens dashboard in default browser
     ├─ dashboard     → runs the dashboard HTTP server (foreground)
     ├─ tray          → dashboard + mDNS + Windows tray icon
     ├─ autostart     → enable / disable / status — HKCU\Run entry
     └─ reset         → wipe config + runtime (--storage for data too)
```

We do **not** fork Immich, Jellyfin, or FileBrowser. We ship their official Docker images and add value via the installer, wizard, dashboard, and tray app.

## Installer (planned end-user flow)

Build with NSIS (see [`docs/build.md`](docs/build.md)):

```powershell
.\scripts\build-release.ps1 -Version 0.7.4
```

Then a typical end-user install will be:

1. Download `snugnas-setup.exe` from a GitHub release.
2. Run it. If Docker Desktop isn't already installed, the wizard's `install` check will guide that first.
3. Open the wizard from the Start menu — it'll ask for storage path and admin email, generate the Jellyfin BaseUrl config, render the docker-compose, pull all the images, and start everything.
4. Dashboard opens at `http://snugnas.local:8080/`. Scan a QR code per service to set up the matching app on your phone.

## Develop

See [`docs/architecture.md`](docs/architecture.md) for package layout, [`docs/contributing.md`](docs/contributing.md) for conventions, and [`docs/build.md`](docs/build.md) for cross-compile and installer build steps.

```powershell
# Run all tests
go test ./...

# Run with race detector + clean cache
go test -race -count=1 ./...

# Cross-compile to Linux (e.g. for CI testing)
$env:GOOS = "linux"; go build .\cmd\snugnas; Remove-Item Env:GOOS
```

## Licenses

snugNAS itself: **MIT** (see [`LICENSE`](LICENSE)).

Bundled services keep their own licenses — they are not modified or redistributed as source, only pulled as Docker images at install time:

- Immich — AGPL-3.0
- Jellyfin — GPL-2.0
- FileBrowser — Apache-2.0
- Caddy — Apache-2.0

## Credits

Built on the shoulders of giants. Huge thanks to the Immich, Jellyfin, FileBrowser, and Caddy maintainers.
