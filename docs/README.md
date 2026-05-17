# snugNAS docs

| File | What's inside |
|---|---|
| [`architecture.md`](architecture.md) | Process and package layout, why each design decision was made, storage layout under `<StorageRoot>` |
| [`roadmap.md`](roadmap.md) | What shipped in each version, what's deferred, what's planned for v1.0 |
| [`build.md`](build.md) | Building the Go binary, cross-compiling, building the NSIS installer, ldflags version stamping |
| [`contributing.md`](contributing.md) | Conventions: no upstream forks, embed assets, Caddy is the only host-exposed ingress |

## Quick reference: what each `snugnas` command does

| Command | Purpose |
|---|---|
| `snugnas install` | Check Docker Desktop and WSL2 are ready; print remediation if not |
| `snugnas wizard` | Run the first-run web wizard — opens browser, walks through storage/hostname, provisions stack |
| `snugnas up` | Start the stack (`docker compose up -d`) |
| `snugnas down` | Stop the stack |
| `snugnas restart [service]` | Bounce all services or one specific service |
| `snugnas status` | Show `docker compose ps` output |
| `snugnas logs [service] [--follow] [--tail=N]` | Stream service logs |
| `snugnas doctor` | Combined prereq + config + service health report |
| `snugnas open` | Open the dashboard in the default browser |
| `snugnas dashboard` | Run the dashboard HTTP server in the foreground |
| `snugnas tray` | Windows tray app — dashboard + mDNS in one always-on process |
| `snugnas autostart {enable\|disable\|status}` | Manage HKCU\Run autostart entry |
| `snugnas reset [--storage] [--yes]` | Wipe config + runtime (and optionally storage). Dry-run by default. |
| `snugnas update` | Check GitHub Releases for a newer snugNAS version |
| `snugnas --version` | Print the version (overridable at build time via ldflags) |
