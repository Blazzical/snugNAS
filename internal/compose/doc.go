// Package compose renders compose/docker-compose.yml.tmpl with the values
// collected by the wizard, writes the result to runtime/docker-compose.yml,
// and provides up/down/status helpers around `docker compose`.
package compose
