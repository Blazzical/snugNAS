# snugNAS

[![CI](https://github.com/Blazzical/snugNAS/actions/workflows/ci.yml/badge.svg)](https://github.com/Blazzical/snugNAS/actions/workflows/ci.yml)

**CasaOS for Windows.** A one-installer NAS stack that runs on top of your existing Windows machine — keep Windows, get a Synology-quality experience.

snugNAS bundles best-in-class open-source services behind a single installer, a non-technical first-run wizard, and a web dashboard with QR codes that point your phone and smart TV at the right app.

## What you get

| Service | What it's for | Phone/TV access |
|---|---|---|
| [Immich](https://immich.app/) | Photo & video backup (Google Photos replacement) | Native iOS / Android apps, background auto-upload |
| [Jellyfin](https://jellyfin.org/) | Movies, TV, music | Android, iOS, AndroidTV, AppleTV, Roku, Fire TV, Chromecast, Kodi, DLNA |
| [FileBrowser](https://filebrowser.org/) | Everything else (docs, downloads, archives) | Mobile-friendly web UI |
| [Caddy](https://caddyserver.com/) | Reverse proxy + automatic HTTPS | n/a |

All four services point at the same storage root and coexist safely.

## Status

**Pre-alpha.** Project scaffolded 2026-05-16. Nothing is buildable yet — see `docs/roadmap.md`.

## How it works

```
snugNAS (Windows host)
 ├─ Docker Desktop / WSL2          ← runtime, installed by our installer
 │   ├─ Immich (+ postgres, redis, ML)
 │   ├─ Jellyfin
 │   ├─ FileBrowser
 │   └─ Caddy                       ← single URL, auto TLS
 └─ snugnas.exe (Go binary)
     ├─ install      → detects/guides Docker Desktop install
     ├─ wizard       → web-based first-run setup (storage dir, admin, hostname)
     ├─ up / down    → starts/stops the compose stack
     ├─ status       → service health
     └─ tray         → Windows tray app + landing page on http://snugnas.local
```

We do **not** fork Immich, Jellyfin, or FileBrowser. We ship their official Docker images and add value via the installer, wizard, dashboard, and tray app.

## Install (planned end-user flow)

1. Download `snugnas-setup.exe` from Releases.
2. Run it. If Docker Desktop isn't installed, it walks you through that first.
3. The wizard opens in your browser:
   - Pick a storage folder (e.g. `D:\snugNAS-data`).
   - Set an admin email + password.
   - Choose a hostname (default: `snugnas.local`).
4. snugNAS pulls the Docker images and starts everything.
5. The dashboard opens at `http://snugnas.local`. Scan a QR code per service to set up the matching app on your phone.

## Develop

Prerequisites: Go 1.22+, git, Docker Desktop.

```powershell
git clone https://github.com/Blazzical/snugNAS.git
cd snugNAS
go build ./cmd/snugnas
.\snugnas.exe wizard
```

See `docs/architecture.md` for layout and `docs/contributing.md` for contribution guidelines (TBD).

## Licenses

snugNAS itself: **MIT** (see `LICENSE`).

Bundled services keep their own licenses — they are not modified or redistributed as source, only pulled as Docker images at install time:

- Immich — AGPL-3.0
- Jellyfin — GPL-2.0
- FileBrowser — Apache-2.0
- Caddy — Apache-2.0

## Credits

Built on the shoulders of giants. Huge thanks to the Immich, Jellyfin, FileBrowser, and Caddy maintainers.
