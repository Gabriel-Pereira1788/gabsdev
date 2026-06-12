package services_test

import (
	"testing"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
)

func TestGetPosts_ReturnsByLanguage(t *testing.T) {
	ptPosts := services.GetPosts(i18n.LangPT)
	enPosts := services.GetPosts(i18n.LangEN)

	if len(ptPosts) == 0 {
		t.Fatal("expected at least one PT post")
	}
	if len(enPosts) == 0 {
		t.Fatal("expected at least one EN post")
	}

	if ptPosts[0].Title == enPosts[0].Title {
		t.Errorf("expected different titles for PT/EN; both have %q", ptPosts[0].Title)
	}
}

func TestGetPosts_SortedByDateDesc(t *testing.T) {
	posts := services.GetPosts(i18n.LangEN)
	if len(posts) < 2 {
		t.Skip("need at least two posts to assert ordering")
	}
	for i := 1; i < len(posts); i++ {
		if posts[i-1].Date < posts[i].Date {
			t.Errorf("posts not sorted by date desc: %q < %q at index %d", posts[i-1].Date, posts[i].Date, i)
		}
	}
}

func TestGetPosts_UnknownLangReturnsEmpty(t *testing.T) {
	posts := services.GetPosts(i18n.Lang("zz"))
	if len(posts) != 0 {
		t.Errorf("expected empty slice for unknown lang, got %d posts", len(posts))
	}
}
