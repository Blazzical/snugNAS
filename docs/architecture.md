# snugNAS architecture

## One-paragraph summary

snugNAS is a Go-based wrapper that turns a Windows machine into an easy-to-use NAS by orchestrating Docker containers. It does not fork or modify Immich, Jellyfin, or FileBrowser — it ships their official images, generates a docker-compose file from a template, runs a first-run wizard, and exposes a dashboard with QR codes for phone-app onboarding. Caddy fronts the stack so users see a single hostname (`snugnas.local` by default).

## Process layout

```
Windows host
├── snugnas.exe (this repo)              ← long-running tray app
│   ├── HTTP server on :7777
│   │   ├── /            → dashboard (web/index.html via embed.FS)
│   │   ├── /api/qr      → QR codes for each service
│   │   └── /wizard/*    → first-run wizard (only when no config exists)
│   └── systray icon     → start/stop/open/quit
└── Docker Desktop
    └── compose project "snugnas"
        ├── caddy            (port 80, 443) ← only host-exposed service
        ├── immich-server    + postgres + redis + machine-learning
        ├── jellyfin
        └── filebrowser
```

Caddy reverse-proxies the dashboard at `/` to `host.docker.internal:7777` (the Go binary on the Windows host). Other paths route to the respective service containers.

## Why this shape

- **Docker required.** Immich is Docker-only and the maintainers refuse to support bare-metal. We can't ship Immich without Docker, so we accept Docker as the runtime for everything.
- **One ingress.** A single hostname/port makes mDNS (`snugnas.local`) work cleanly and means QR codes have one URL pattern.
- **Wrapper on the host, not in a container.** We need to launch the user's browser, manage Windows tray, and detect Docker Desktop — all host-side concerns.
- **Embedded assets.** The dashboard HTML/CSS/JS lives in `web/` and is embedded into the Go binary via `embed.FS`. Single-binary install story is core to the project's pitch.

## Storage layout

The wizard asks for one `StorageRoot` (e.g. `D:\snugNAS-data`). Under that:

```
{StorageRoot}/
├── photos/              ← Immich uploads
├── immich-db/           ← Immich postgres data
├── media/               ← Jellyfin reads (read-only mount)
├── jellyfin/config/
├── jellyfin/cache/
├── files/               ← FileBrowser serves this
└── filebrowser/         ← FileBrowser config + database
```

All three apps can read the same files safely (Immich documents this explicitly).

## Module map

| Package | Purpose |
|---|---|
| `cmd/snugnas` | CLI entrypoint, just calls into `internal/cli` |
| `internal/cli` | Cobra commands: install, wizard, up, down, status, open, tray |
| `internal/installer` | Detects Docker Desktop + WSL2, guides install if missing |
| `internal/compose` | Renders `compose/docker-compose.yml.tmpl` and `compose/Caddyfile.tmpl`, drives `docker compose` |
| `internal/wizard` | First-run web wizard (storage path, admin creds, hostname) |
| `internal/landing` | Dashboard HTTP server (tiles, QR codes, status) |
| `internal/tray` | Windows tray app via getlantern/systray |
| `internal/config` | Read/write `%APPDATA%\snugNAS\config.toml` |
| `internal/compose/templates/` | Embedded Compose + Caddyfile templates |
| `internal/landing/web/` | Embedded dashboard assets |

## Decisions deliberately made

- **No fork of Immich/Jellyfin/FileBrowser.** Adds maintenance burden with no clear win.
- **Web wizard, not native GUI.** We already need a web stack for the dashboard.
- **Docker Desktop over Podman as primary.** Podman supported as a documented alternative.
- **No bundled OIDC / unified auth in v1.** Each service has its own admin user. v2 may add an SSO layer if there's demand.
