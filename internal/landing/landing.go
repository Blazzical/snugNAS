// Package landing serves the snugNAS dashboard at the root of snugnas.local —
// service tiles, QR codes for mobile app onboarding, and a status panel.
// Assets are embedded via embed.FS so the binary remains self-contained.
package landing

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Blazzical/snugNAS/internal/compose"
	"github.com/Blazzical/snugNAS/internal/config"
	qrcode "github.com/skip2/go-qrcode"
)

//go:embed web
var webFS embed.FS

// Handlers returns an http.Handler that serves the dashboard at /, its static
// assets at /static/*, and PNG QR codes at /api/qr?service=<name>. Useful for
// mounting the dashboard inside a larger server (e.g. the wizard, after
// provisioning completes).
func Handlers(cfg *config.Config) http.Handler {
	sub, _ := fs.Sub(webFS, "web")
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/api/qr", qrHandler(cfg))
	mux.HandleFunc("/api/health", healthHandler())
	mux.HandleFunc("/", indexHandler(sub))
	return mux
}

// healthHandler returns the per-service docker compose ps status as JSON.
// Used by the dashboard to colour the status dot on each tile.
func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		statuses, err := compose.PS(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if statuses == nil {
			statuses = []compose.ServiceStatus{} // avoid `null` in the JSON
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(statuses)
	}
}

// Serve starts the dashboard HTTP server on localhost:<DashboardPort> and
// blocks until ctx is canceled.
func Serve(ctx context.Context, cfg *config.Config) error {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.DashboardPort))
	srv := &http.Server{
		Addr:         addr,
		Handler:      Handlers(cfg),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func indexHandler(sub fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := fs.ReadFile(sub, "index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	}
}

// qrHandler renders a PNG QR code that encodes the public URL for the
// requested service. Used by the dashboard tiles.
func qrHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := r.URL.Query().Get("service")
		target, err := serviceURL(cfg, service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		png, err := qrcode.Encode(target, qrcode.Medium, 256)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write(png)
	}
}

// serviceURL returns the URL a phone or browser should use for a service.
// v0.2 uses direct host:port for each service; v0.4 will switch to
// https://hostname/<service>/ via Caddy once Base URLs and TLS are sorted.
func serviceURL(cfg *config.Config, service string) (string, error) {
	switch service {
	case "immich":
		return fmt.Sprintf("http://%s:2283/", cfg.Hostname), nil
	case "jellyfin":
		return fmt.Sprintf("http://%s:8080/jellyfin/", cfg.Hostname), nil
	case "filebrowser":
		return fmt.Sprintf("http://%s:8080/files/", cfg.Hostname), nil
	default:
		return "", fmt.Errorf("unknown service %q", service)
	}
}
