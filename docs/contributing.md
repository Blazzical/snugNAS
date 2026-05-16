# Contributing to snugNAS

Thank you for considering a contribution.

## Dev setup

You'll need:

- **Go 1.22+** — https://go.dev/dl/
- **Docker Desktop** — https://www.docker.com/products/docker-desktop/
- **Git**

Then:

```powershell
git clone https://github.com/Blazzical/snugNAS.git
cd snugNAS
go build ./cmd/snugnas
.\snugnas.exe --help
```

## Project layout

See [`docs/architecture.md`](architecture.md).

## Conventions

- **No upstream forks.** We ship Immich, Jellyfin, and FileBrowser as official Docker images. If something needs changing upstream, send the patch upstream — don't fork.
- **Web wizard, not native GUI.** Keep new user-facing flows in `internal/landing` / `internal/wizard`.
- **Caddy is the only host-exposed service.** Don't add `ports:` mappings in the compose template for individual app containers.
- **Embed everything.** Dashboard assets go under `web/` and are embedded via `embed.FS`. No external static-asset directories at runtime.

## Filing issues

Please include:

- snugNAS version (`snugnas --version`)
- Windows version
- Docker Desktop version (`docker version`)
- `snugnas status` output
- The relevant lines from `%APPDATA%\snugNAS\snugnas.log` (planned)

## Code of conduct

Be kind. We're all here because we don't want to pay Google.
