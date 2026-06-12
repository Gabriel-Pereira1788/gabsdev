package services_test

import (
	"strings"
	"testing"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
)

func TestSearchPosts_EmptyQueryReturnsTagResults(t *testing.T) {
	got := services.SearchPosts(i18n.LangEN, "", "")
	all := services.GetPostsByTag(i18n.LangEN, "")
	if len(got) != len(all) {
		t.Errorf("empty query should delegate to GetPostsByTag (%d), got %d", len(all), len(got))
	}
}

func TestSearchPosts_MatchesTitleCaseInsensitive(t *testing.T) {
	posts := services.GetPosts(i18n.LangEN)
	if len(posts) == 0 {
		t.Skip("no posts available")
	}
	title := posts[0].Title
	if len(title) < 3 {
		t.Skip("title too short to derive query")
	}
	query := strings.ToUpper(title[:3])

	got := services.SearchPosts(i18n.LangEN, query, "")
	found := false
	for _, p := range got {
		if p.Slug == posts[0].Slug {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected post %q in results for query %q", posts[0].Slug, query)
	}
}

func TestSearchPosts_NoMatchReturnsEmpty(t *testing.T) {
	got := services.SearchPosts(i18n.LangEN, "zzz-no-such-content-zzz", "")
	if len(got) != 0 {
		t.Errorf("expected empty result for non-matching query, got %d", len(got))
	}
}
