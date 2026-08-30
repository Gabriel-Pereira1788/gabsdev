package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"gabsdev-go/internal/domain"
)

// maxRecentCommits caps how many push events render in the activity list —
// enough to feel alive without turning the section into a full feed.
const maxRecentCommits = 6

type ghEvent struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Repo      struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Head    string `json:"head"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
	} `json:"payload"`
}

// fetchRecentCommits is best-effort: any failure returns nil so the
// contribution heatmap can still render without the commit list. It renders
// one row per push (GitHub's own activity feed does the same, not one row
// per individual commit within a push).
func fetchRecentCommits() []domain.Commit {
	url := fmt.Sprintf("%s/users/%s/events/public?per_page=30", restBase, username())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if tok := token(); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var events []ghEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil
	}

	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt.After(events[j].CreatedAt) })

	commits := make([]domain.Commit, 0, maxRecentCommits)
	for _, e := range events {
		if e.Type != "PushEvent" {
			continue
		}

		message := ""
		switch {
		case len(e.Payload.Commits) > 0:
			message = firstLine(e.Payload.Commits[len(e.Payload.Commits)-1].Message)
		case e.Payload.Head != "":
			// GitHub's events API often omits the commits array even for
			// real pushes (privacy trimming) but still includes the head
			// SHA — fetch that single commit's message directly.
			message = fetchCommitMessage(e.Repo.Name, e.Payload.Head)
		}
		if message == "" {
			continue
		}

		commits = append(commits, domain.Commit{
			Repo:      shortRepoName(e.Repo.Name),
			RepoURL:   "https://github.com/" + e.Repo.Name,
			Message:   message,
			Timestamp: e.CreatedAt,
			LangClass: repoLangClass(e.Repo.Name),
		})
		if len(commits) >= maxRecentCommits {
			break
		}
	}

	return commits
}

func shortRepoName(full string) string {
	for i := len(full) - 1; i >= 0; i-- {
		if full[i] == '/' {
			return full[i+1:]
		}
	}
	return full
}

func firstLine(msg string) string {
	for i, r := range msg {
		if r == '\n' {
			return msg[:i]
		}
	}
	return msg
}
