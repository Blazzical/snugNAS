# snugNAS roadmap

## v0.1 — Scaffolding (done)

- [x] Repo layout (cmd/, internal/, compose/, web/, docs/)
- [x] go.mod, README, LICENSE, .gitignore
- [x] CLI skeleton with cobra (install/wizard/up/down/status/open/tray)
- [x] docker-compose + Caddyfile templates
- [x] Static dashboard HTML/CSS

## v0.2 — Buildable (done)

- [x] `internal/config` reads/writes `%APPDATA%\snugNAS\config.toml` (TOML, BOM-tolerant)
- [x] `internal/compose` renders templates with config values
- [x] `internal/landing` serves embedded `web/` via embed.FS, generates QR codes
- [x] `snugnas up` shells out to `docker compose up -d`
- [x] `snugnas down`, `snugnas status`, `snugnas dashboard`, `snugnas open` wired up

## v0.3 — Wizard (done)

- [x] `snugnas wizard` serves a single-page web flow on :7777
- [x] Collects storage path, hostname, admin email (password collection deferred — services own admin signup)
- [x] `POST /api/commit` writes config, kicks off provisioning in a goroutine
- [x] Page polls `GET /api/state` for live log + phase
- [x] On `phase=done`, page redirects to the dashboard at http://localhost:8080/
- [x] CLI auto-opens the browser when `snugnas wizard` starts

## v0.4 — Installer & mDNS (done)

- [x] `snugnas install` checks docker on PATH, daemon reachable, WSL2 default
- [x] Prints a checklist with `[ok]` / `[--]` marks and one-line remediation per failure
- [x] mDNS broadcast for `<hostname>.local` via grandcat/zeroconf
- [x] mDNS started by `snugnas wizard` and `snugnas dashboard`; stopped on Ctrl+C

## v0.5 — Jellyfin Base URL via subpath routing (done)

- [x] Wizard writes/patches `<StorageRoot>/jellyfin/config/config/network.xml`
      with `<BaseUrl>/jellyfin</BaseUrl>` before first start
- [x] Direct port 8096:8096 mapping removed from compose
- [x] Caddy `/jellyfin/*` routes work end-to-end (web UI assets included)
- [x] Dashboard tile + QR code updated to `http://<hostname>:8080/jellyfin/`

## Deferred to later

- Subdomain routing in Caddy (`immich.snugnas.local`, etc.) — depends on
  reliable mDNS multi-label `.local` resolution; mixed support across Android
- Caddy `tls internal` for LAN-only HTTPS — needs trust-store install flow
- Immich behind Caddy subpath — Immich has no Base URL config; needs subdomains
- Option in `snugnas install` to free port 80 from IIS (W3SVC)

## v0.6 — Windows polish

- [ ] Tray app: start/stop/open dashboard/quit
- [ ] Auto-start at login (Windows service or scheduled task)
- [ ] NSIS-based `snugnas-setup.exe` (matches Jellyfin's installer style)
- [ ] Signed binary

## v1.0 — Ship

- [ ] Auto-update for the wrapper
- [ ] Compose-image pinning + version channel (stable/beta)
- [ ] Backup/restore for the storage root
- [ ] Documentation site
- [ ] Public release announcement

## Maybe later

- Optional SSO (Authelia + the Jellyfin Authelia plugin + Immich OAuth + filebrowser proxy auth)
- Plug-in app catalog (Vaultwarden, Navidrome, Paperless-ngx, etc.) — "CasaOS for Windows" in the fullest sense
- macOS and Linux host support (Go binary is already portable; mostly a packaging problem)
- Custom branded mobile app — only if user demand justifies it; Jellyfin and Immich apps are excellent already
