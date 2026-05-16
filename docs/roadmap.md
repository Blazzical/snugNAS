# snugNAS roadmap

## v0.1 — Scaffolding (done)

- [x] Repo layout (cmd/, internal/, compose/, web/, docs/)
- [x] go.mod, README, LICENSE, .gitignore
- [x] CLI skeleton with cobra (install/wizard/up/down/status/open/tray)
- [x] docker-compose + Caddyfile templates
- [x] Static dashboard HTML/CSS

## v0.2 — Buildable

- [ ] `internal/config` reads/writes `%APPDATA%\snugNAS\config.toml`
- [ ] `internal/compose` renders templates with config values
- [ ] `internal/landing` serves `web/` via embed.FS, generates QR codes
- [ ] `snugnas up` actually runs `docker compose up -d`
- [ ] `snugnas down` and `snugnas status` work end-to-end

## v0.3 — Wizard

- [ ] `snugnas wizard` serves a 3-step web flow on :7777
  - Step 1: pick storage directory
  - Step 2: admin email + password + hostname
  - Step 3: pull images, start stack, show progress
- [ ] On success, browser auto-redirects to the dashboard

## v0.4 — Installer

- [ ] `snugnas install` detects Docker Desktop + WSL2
- [ ] If missing, opens download page and waits for the user to finish install
- [ ] mDNS broadcast for `snugnas.local` (Bonjour service publisher)

## v0.5 — Windows polish

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
