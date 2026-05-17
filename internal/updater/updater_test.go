package updater

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLatestReleaseSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"tag_name": "v1.2.3",
			"name": "v1.2.3 — release name",
			"html_url": "https://example.com/release/1",
			"published_at": "2026-05-17T00:00:00Z",
			"assets": [
				{"name": "snugnas.exe", "browser_download_url": "https://example.com/dl/snugnas.exe"}
			]
		}`))
	}))
	defer srv.Close()

	r, err := fetchFrom(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("fetchFrom: %v", err)
	}
	if r.TagName != "v1.2.3" {
		t.Errorf("TagName = %q, want v1.2.3", r.TagName)
	}
	if len(r.Assets) != 1 || r.Assets[0].Name != "snugnas.exe" {
		t.Errorf("Assets[0] = %+v", r.Assets)
	}
}

func TestLatestReleaseNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := fetchFrom(context.Background(), srv.URL)
	if !errors.Is(err, ErrNoReleases) {
		t.Errorf("err = %v, want ErrNoReleases", err)
	}
}

func TestLatestReleaseServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := fetchFrom(context.Background(), srv.URL)
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Errorf("err = %v, want HTTP 500", err)
	}
}
