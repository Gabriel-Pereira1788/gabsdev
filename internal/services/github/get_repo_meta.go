package github

import (
	"encoding/json"
	"net/http"
	"sync"
)

// langClasses maps a repo's primary language (as reported by the GitHub API)
// to one of the fixed CSS class suffixes defined in
// static/css/global.css (.gh-dot.lang-*). This is a closed, hand-picked set
// — never raw API text — so it's safe to interpolate directly into a class
// attribute in the template.
var langClasses = map[string]string{
	"Go":         "go",
	"TypeScript": "ts",
	"JavaScript": "js",
	"Swift":      "swift",
	"Kotlin":     "kotlin",
	"Python":     "py",
	"HTML":       "html",
	"CSS":        "css",
	"Shell":      "shell",
	"Dart":       "dart",
	"Java":       "java",
	"C++":        "cpp",
	"Ruby":       "ruby",
}

const defaultLangClass = "other"

var (
	repoLangMu    sync.Mutex
	repoLangCache = map[string]string{}
)

// repoLangClass resolves (and caches, for the lifetime of the process — repo
// primary language rarely changes) a repo's language dot class. Any failure
// falls back to the neutral "other" class instead of breaking the row.
func repoLangClass(fullName string) string {
	repoLangMu.Lock()
	if c, ok := repoLangCache[fullName]; ok {
		repoLangMu.Unlock()
		return c
	}
	repoLangMu.Unlock()

	class := defaultLangClass
	if lang := fetchRepoPrimaryLanguage(fullName); lang != "" {
		if c, ok := langClasses[lang]; ok {
			class = c
		}
	}

	repoLangMu.Lock()
	repoLangCache[fullName] = class
	repoLangMu.Unlock()

	return class
}

func fetchRepoPrimaryLanguage(fullName string) string {
	req, err := http.NewRequest(http.MethodGet, restBase+"/repos/"+fullName, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if tok := token(); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var parsed struct {
		Language string `json:"language"`
	}
	if json.NewDecoder(resp.Body).Decode(&parsed) != nil {
		return ""
	}
	return parsed.Language
}

// fetchCommitMessage fetches a single commit's message by SHA — used as a
// fallback when a PushEvent's payload.commits array is empty (GitHub trims
// commit details from the events feed for some pushes, e.g. merges, but
// still includes the head SHA). Returns "" on any failure, which the caller
// treats as "skip this event".
func fetchCommitMessage(fullName, sha string) string {
	req, err := http.NewRequest(http.MethodGet, restBase+"/repos/"+fullName+"/commits/"+sha, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if tok := token(); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var parsed struct {
		Commit struct {
			Message string `json:"message"`
		} `json:"commit"`
	}
	if json.NewDecoder(resp.Body).Decode(&parsed) != nil {
		return ""
	}
	return firstLine(parsed.Commit.Message)
}
