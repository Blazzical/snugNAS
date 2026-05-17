# snugNAS roadmap

## v0.1 — Scaffolding (done)

- [x] Repo layout (cmd/, internal/, compose/, web/, docs/)
- [x] go.mod, README, LICENSE, .gitignore
- [x] CLI skeleton with cobra
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
- [x] On `phase=done`, page redirects to the dashboard
- [x] CLI auto-opens the browser when `snugnas wizard` starts
- [x] Same server transparently swaps to dashboard handlers after `phase=done`
      so the post-up redirect just works (v0.3.1)

## v0.4 — Installer prereq check + mDNS (done)

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

## v0.6 — Windows polish (done)

- [x] **v0.6.0** Tray app via getlantern/systray — dashboard + mDNS in one
      always-on process. Menu: Open dashboard / Open on LAN / Quit. Hand-rolled
      16×16 ICO icon (pure Go, no Cgo).
- [x] **v0.6.1** `snugnas autostart enable|disable|status` — HKCU\Run entry,
      build-tag split so non-Windows builds still compile (stub returns
      ErrUnsupported).
- [x] **v0.6.2** NSIS installer script + MUI2 pages + Start-menu shortcuts +
      Add/Remove Programs registration + uninstaller (data preserved).
      `scripts\build-release.ps1` automates the full build.
- [x] **v0.6.3** `snugnas logs [service]` + `snugnas restart [service]` wrappers
      around docker compose.

## v0.7 — Reset, tests, CI, live dashboard (done)

- [x] **v0.7.0** `snugnas reset` — wipes config + runtime; opt-in `--storage`
      to also delete data. Dry-run by default; requires `--yes`.
- [x] **v0.7.1** Unit tests for config, mdns, installer, tray, autostart.
      `t.Setenv` isolates tests from real user state.
- [x] **v0.7.2** GitHub Actions CI on ubuntu-latest + windows-latest: vet,
      test -race, build, plus a cross-compile to windows-amd64 from linux.
      Tray split into `tray_windows.go` / `tray_other.go` so Linux can build.
- [x] **v0.7.3** Dashboard live service health — `/api/health` returns docker
      compose ps JSON, tiles show a coloured status dot polled every 5s.
- [x] **v0.7.4** Cross-platform fix: `PosixStorageRoot` now uses
      `strings.ReplaceAll` instead of `filepath.ToSlash` so backslash → slash
      conversion is unconditional regardless of the build host.

## v0.8 — Polish + release machinery (done)

- [x] **v0.8.0** README polish — accurate status, full CLI reference, quickstart
- [x] **v0.8.1** `snugnas update` — checks GitHub Releases API for newer
      versions, prints assets without auto-downloading
- [x] **v0.8.2** `gofmt` enforced in CI (fails the build if any file is unformatted)
- [x] **v0.8.3** Unit tests for the updater using httptest
- [x] **v0.8.4** `/api/info` endpoint + dashboard footer showing version /
      hostname / uptime
- [x] **v0.8.5** GoReleaser config (`.goreleaser.yml`) + release workflow.
      Cross-builds win/linux/darwin × amd64/arm64, zip on Windows + tar.gz
      elsewhere, SHA256 checksums, **draft releases** (never auto-publishes).
- [x] **v0.8.6** `snugnas doctor` — combined prereq + config + service health
      report; landing HTTP handler tests (httptest covers /, /api/qr/info)
- [x] **v0.8.7** Wizard storage validation — path must exist or be creatable
      and writable, surfaced as a 400 on /api/commit. Tested.

## Deferred (still on the wishlist)

- Subdomain routing in Caddy (`immich.snugnas.local`, etc.) — depends on
  reliable mDNS multi-label `.local` resolution; mixed support across Android
- Caddy `tls internal` for LAN-only HTTPS — needs trust-store install flow
- Immich behind Caddy subpath — Immich has no Base URL config; needs subdomains
- Option in `snugnas install` / `snugnas doctor --fix` to free port 80 from IIS

## v0.9 / v1.0 — Ship

- [ ] Tag a v0.7+ release so `snugnas update` finds something (the GoReleaser
      workflow will create a *draft* release that the maintainer publishes
      manually)
- [ ] Signed binary (Windows code signing cert)
- [ ] Compose-image pinning (currently `:latest` and `:release` floats)
- [ ] Backup / restore for the storage root
- [ ] Documentation site (mkdocs or similar) — most content already in docs/
- [ ] Public release announcement

## Maybe later

- Optional SSO (Authelia + the Jellyfin Authelia plugin + Immich OAuth + filebrowser proxy auth)
- Plug-in app catalog (Vaultwarden, Navidrome, Paperless-ngx, etc.) — "CasaOS for Windows" in the fullest sense
- macOS and Linux host support (Go binary is already portable; mostly a packaging problem)
- Custom branded mobile app — only if user demand justifies it; Jellyfin and Immich apps are excellent already
