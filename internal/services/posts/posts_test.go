package services_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root (go.mod)")
		}
		dir = parent
	}
	if err := os.Chdir(dir); err != nil {
		panic("chdir to project root: " + err.Error())
	}
}

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
