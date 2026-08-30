package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	baseURLOnce sync.Once
	baseURL     string
)

func BaseURL() string {
	baseURLOnce.Do(func() {
		baseURL = os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		baseURL = strings.TrimRight(baseURL, "/")
	})
	return baseURL
}

var (
	assetVersionOnce sync.Once
	assetVersion     string
)

// AssetVersion returns a cache-busting token appended to static asset URLs
// (?v=...) so browsers fetch fresh CSS/JS/images after every deploy instead
// of silently serving a stale cached copy — Go's http.FileServer sets
// Last-Modified but no Cache-Control, which lets browsers cache static
// files far more aggressively (via heuristic freshness) than intended.
// Derived from the running binary's own mtime, so it changes on every
// rebuild without needing a build-time injected version string.
func AssetVersion() string {
	assetVersionOnce.Do(func() {
		exe, err := os.Executable()
		if err == nil {
			if info, err := os.Stat(exe); err == nil {
				assetVersion = strconv.FormatInt(info.ModTime().Unix(), 10)
				return
			}
		}
		assetVersion = strconv.FormatInt(time.Now().Unix(), 10)
	})
	return assetVersion
}
