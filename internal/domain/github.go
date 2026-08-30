package domain

import "time"

// ContributionDay is a single day in the GitHub contribution heatmap.
// Level is a 0-4 bucket (0 = no activity, 4 = busiest day in the window),
// computed relative to the window's own max so the heatmap stays readable
// regardless of how active a given period was.
type ContributionDay struct {
	Date  string
	Count int
	Level int
}

type ContributionWeek struct {
	Days []ContributionDay
}

// Commit is one push event from the GitHub activity feed, condensed to a
// single row (GitHub's own feed does the same: one entry per push, not per
// individual commit inside it).
type Commit struct {
	Repo      string
	RepoURL   string
	Message   string
	Timestamp time.Time
	// LangClass is a small fixed set of CSS class suffixes (see
	// static/css/global.css .gh-dot.lang-*), never raw user/API data, so it
	// is safe to interpolate directly into a class attribute.
	LangClass string
}

// GitHubActivity is the full payload rendered by the GitHub activity
// section on the home page. A zero value (no weeks, no commits) means the
// data could not be fetched and the section renders its empty-state copy
// instead of breaking the page.
type GitHubActivity struct {
	Weeks              []ContributionWeek
	TotalContributions int
	Commits            []Commit
}
