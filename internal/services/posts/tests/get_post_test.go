package services_test

import (
	"testing"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
)

func TestGetPost_ExistingSlug(t *testing.T) {
	slug := "next-para-go-templ-htmx"

	post, ok := services.GetPost(i18n.LangEN, slug)
	if !ok {
		t.Fatalf("expected EN post for slug %q to be found", slug)
	}
	if post.Slug != slug {
		t.Errorf("expected slug %q, got %q", slug, post.Slug)
	}
}

func TestGetPost_MissingSlugReturnsNotFound(t *testing.T) {
	_, ok := services.GetPost(i18n.LangEN, "does-not-exist")
	if ok {
		t.Error("expected ok=false for non-existent slug")
	}
}

func TestGetAdjacentPosts_MissingSlugReturnsNil(t *testing.T) {
	prev, next := services.GetAdjacentPosts(i18n.LangEN, "does-not-exist")
	if prev != nil || next != nil {
		t.Errorf("expected nil prev/next for missing slug, got prev=%v next=%v", prev, next)
	}
}

func TestGetAdjacentPosts_MiddlePost(t *testing.T) {
	posts := services.GetPosts(i18n.LangEN)
	if len(posts) < 3 {
		t.Skip("need at least three posts to assert middle adjacency")
	}
	mid := len(posts) / 2
	prev, next := services.GetAdjacentPosts(i18n.LangEN, posts[mid].Slug)

	if prev == nil || prev.Slug != posts[mid-1].Slug {
		t.Errorf("expected prev %q, got %v", posts[mid-1].Slug, prev)
	}
	if next == nil || next.Slug != posts[mid+1].Slug {
		t.Errorf("expected next %q, got %v", posts[mid+1].Slug, next)
	}
}
