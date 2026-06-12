package services_test

import (
	"slices"
	"testing"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
)

func TestGetPostsByTag_EmptyTagReturnsAll(t *testing.T) {
	all := services.GetPosts(i18n.LangEN)
	got := services.GetPostsByTag(i18n.LangEN, "")
	if len(got) != len(all) {
		t.Errorf("empty tag should return all %d posts, got %d", len(all), len(got))
	}
}

func TestGetPostsByTag_FiltersByTag(t *testing.T) {
	posts := services.GetPosts(i18n.LangEN)

	var tag string
	for _, p := range posts {
		if len(p.Tags) > 0 {
			tag = p.Tags[0]
			break
		}
	}
	if tag == "" {
		t.Skip("no tagged post available")
	}

	got := services.GetPostsByTag(i18n.LangEN, tag)
	if len(got) == 0 {
		t.Fatalf("expected at least one post for tag %q", tag)
	}
	for _, p := range got {
		if !slices.Contains(p.Tags, tag) {
			t.Errorf("post %q returned for tag %q but lacks it", p.Slug, tag)
		}
	}
}
