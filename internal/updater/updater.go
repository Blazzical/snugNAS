// Package updater talks to the GitHub Releases API to check whether a newer
// snugNAS version is available. It does not download or apply updates — it
// just surfaces information so the CLI / dashboard can nudge the user.
package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// DefaultRepo is the GitHub repo we check for releases of.
const DefaultRepo = "Blazzical/snugNAS"

// Release is a trimmed-down view of the GitHub release payload.
type Release struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	HTMLURL     string  `json:"html_url"`
	PublishedAt string  `json:"published_at"`
	Assets      []Asset `json:"assets"`
}

// Asset is one downloadable file attached to a release.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// ErrNoReleases is returned when GitHub responds 404 — typically because no
// releases have been published for the repo yet.
var ErrNoReleases = errors.New("no releases published yet")

// LatestRelease fetches the latest release of the default repo.
func LatestRelease(ctx context.Context) (*Release, error) {
	return LatestReleaseFor(ctx, DefaultRepo)
}

// LatestReleaseFor fetches the latest release of a specific GitHub repo.
func LatestReleaseFor(ctx context.Context, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	return fetchFrom(ctx, url)
}

// fetchFrom is the inner implementation, factored out so tests can point at
// an httptest.Server instead of the real GitHub API.
func fetchFrom(ctx context.Context, url string) (*Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "snugnas")

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrNoReleases
	case http.StatusOK:
		var r Release
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return nil, fmt.Errorf("decode release: %w", err)
		}
		return &r, nil
	default:
		return nil, fmt.Errorf("github api: HTTP %d", resp.StatusCode)
	}
}
