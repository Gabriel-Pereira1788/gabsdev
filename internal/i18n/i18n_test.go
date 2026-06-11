package i18n_test

import (
	"testing"

	"gabsdev-go/internal/i18n"
)

func TestT_KnownKeys(t *testing.T) {
	ptVal := i18n.T(i18n.LangPT, "nav.home")
	if ptVal != "~/home" {
		t.Errorf("expected ~/home, got %q", ptVal)
	}

	enVal := i18n.T(i18n.LangEN, "home.view_all")
	if enVal != "view all" {
		t.Errorf("expected 'view all', got %q", enVal)
	}
}

func TestT_FallbackToPT(t *testing.T) {
	got := i18n.T(i18n.LangEN, "test.pt_only")
	if got != "pt-fallback-value" {
		t.Errorf("expected PT fallback value %q, got %q", "pt-fallback-value", got)
	}
}

func TestT_MissingKeyFallsBackToKey(t *testing.T) {
	result := i18n.T(i18n.LangEN, "key.that.does.not.exist")
	if result != "key.that.does.not.exist" {
		t.Errorf("expected raw key fallback, got %q", result)
	}
}

func TestToggleURL_PTToEN(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/about", "/en/about"},
		{"/", "/en"},
		{"/articles/next-para-go-templ-htmx", "/en/articles/next-para-go-templ-htmx"},
		{"/articles?tag=go", "/en/articles?tag=go"},
	}
	for _, tc := range cases {
		got := i18n.ToggleURL(i18n.LangPT, tc.path)
		if got != tc.want {
			t.Errorf("ToggleURL(PT, %q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestToggleURL_ENToPT(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/en/about", "/about"},
		{"/en", "/"},
		{"/en/articles/next-para-go-templ-htmx", "/articles/next-para-go-templ-htmx"},
		{"/en/articles?tag=go", "/articles?tag=go"},
	}
	for _, tc := range cases {
		got := i18n.ToggleURL(i18n.LangEN, tc.path)
		if got != tc.want {
			t.Errorf("ToggleURL(EN, %q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestPath(t *testing.T) {
	if got := i18n.Path(i18n.LangEN, "/articles/x"); got != "/en/articles/x" {
		t.Errorf("Path(EN, /articles/x) = %q, want /en/articles/x", got)
	}
	if got := i18n.Path(i18n.LangPT, "/articles/x"); got != "/articles/x" {
		t.Errorf("Path(PT, /articles/x) = %q, want /articles/x", got)
	}
}
