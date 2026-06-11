package handlers

import (
	"net/http"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/views/pages"
)

func About(lang i18n.Lang) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postsCount := len(services.GetPosts(lang))
		pages.AboutPage(lang, r.URL.Path, r.URL.RequestURI(), postsCount).Render(r.Context(), w)
	}
}
