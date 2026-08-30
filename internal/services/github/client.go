// Package github fetches the personal GitHub contribution heatmap and
// recent push activity for the home page, cached in memory so no page load
// ever waits on a live GitHub API round trip (and traffic never comes close
// to rate limits).
package github

import (
	"net/http"
	"os"
	"sync"
	"time"

	"gabsdev-go/internal/domain"
)

const (
	graphqlURL  = "https://api.github.com/graphql"
	restBase    = "https://api.github.com"
	cacheTTL    = 15 * time.Minute
	httpTimeout = 8 * time.Second
)

var httpClient = &http.Client{Timeout: httpTimeout}

func token() string {
	return os.Getenv("GITHUB_TOKEN")
}

func username() string {
	if u := os.Getenv("GITHUB_USERNAME"); u != "" {
		return u
	}
	return "Gabriel-Pereira1788"
}

var (
	cacheMu      sync.RWMutex
	cached       domain.GitHubActivity
	cachedAt     time.Time
	hasCache     bool
	refreshing   bool
	refreshingMu sync.Mutex
)

// GetActivity returns the cached GitHub activity and never blocks: a cold
// cache (nothing fetched yet) returns a zero-value GitHubActivity — the
// section renders its empty state — while a background fetch fills it in
// for the next request. A stale-but-populated cache is served immediately
// and refreshed in the background. A slow or failing GitHub API can only
// ever delay freshness, never a page load.
func GetActivity() domain.GitHubActivity {
	cacheMu.RLock()
	snapshot := cached
	has := hasCache
	stale := time.Since(cachedAt) > cacheTTL
	cacheMu.RUnlock()

	if !has || stale {
		go backgroundRefresh()
	}
	return snapshot
}

// WarmCache triggers an initial fetch in the background, meant to be called
// once at process startup so the cache is already populated (or at least
// in flight) before the first real visitor ever hits GetActivity.
func WarmCache() {
	go backgroundRefresh()
}

func backgroundRefresh() {
	refreshingMu.Lock()
	if refreshing {
		refreshingMu.Unlock()
		return
	}
	refreshing = true
	refreshingMu.Unlock()

	defer func() {
		refreshingMu.Lock()
		refreshing = false
		refreshingMu.Unlock()
	}()

	fetchAndStore()
}

func fetchAndStore() domain.GitHubActivity {
	weeks, total, err := fetchContributions()
	if err != nil {
		cacheMu.RLock()
		snapshot, has := cached, hasCache
		cacheMu.RUnlock()
		if has {
			return snapshot
		}
		return domain.GitHubActivity{}
	}

	activity := domain.GitHubActivity{
		Weeks:              weeks,
		TotalContributions: total,
		Commits:            fetchRecentCommits(),
	}

	cacheMu.Lock()
	cached = activity
	cachedAt = time.Now()
	hasCache = true
	cacheMu.Unlock()

	return activity
}
