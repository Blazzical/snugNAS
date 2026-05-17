package landing

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Blazzical/snugNAS/internal/config"
)

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	cfg := &config.Config{
		Hostname:      "test.local",
		DashboardPort: 7777,
	}
	return Handlers(cfg)
}

func TestIndexServesHTML(t *testing.T) {
	srv := httptest.NewServer(newTestHandler(t))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html...", ct)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "<title>snugNAS</title>") {
		t.Error("body missing <title>snugNAS</title>")
	}
}

func TestQR(t *testing.T) {
	srv := httptest.NewServer(newTestHandler(t))
	defer srv.Close()

	for _, svc := range []string{"immich", "jellyfin", "filebrowser"} {
		t.Run(svc, func(t *testing.T) {
			resp, err := http.Get(srv.URL + "/api/qr?service=" + svc)
			if err != nil {
				t.Fatalf("GET /api/qr: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want 200", resp.StatusCode)
			}
			if ct := resp.Header.Get("Content-Type"); ct != "image/png" {
				t.Errorf("Content-Type = %q, want image/png", ct)
			}
			body, _ := io.ReadAll(resp.Body)
			// PNG magic header: 89 50 4E 47 0D 0A 1A 0A
			if len(body) < 8 || string(body[1:4]) != "PNG" {
				t.Errorf("body doesn't start with PNG magic: %x", body[:min(8, len(body))])
			}
		})
	}
}

func TestQRUnknownService(t *testing.T) {
	srv := httptest.NewServer(newTestHandler(t))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/qr?service=banana")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestInfo(t *testing.T) {
	srv := httptest.NewServer(newTestHandler(t))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/info")
	if err != nil {
		t.Fatalf("GET /api/info: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var info struct {
		Version  string `json:"version"`
		Hostname string `json:"hostname"`
		Uptime   string `json:"uptime"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if info.Hostname != "test.local" {
		t.Errorf("Hostname = %q, want test.local", info.Hostname)
	}
	if info.Version == "" {
		t.Error("Version is empty")
	}
	if info.Uptime == "" {
		t.Error("Uptime is empty")
	}
}

func TestStaticAssets(t *testing.T) {
	srv := httptest.NewServer(newTestHandler(t))
	defer srv.Close()

	for _, path := range []string{"/static/style.css", "/static/dashboard.js"} {
		resp, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s status = %d, want 200", path, resp.StatusCode)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
