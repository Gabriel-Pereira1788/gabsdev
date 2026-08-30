package handlers

import (
	"net/http"
	"strings"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/views/pages"
)

func Home(lang i18n.Lang) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		prefix := i18n.Prefix(lang)
		if prefix != "" && strings.HasPrefix(path, prefix) {
			path = path[len(prefix):]
			if path == "" {
				path = "/"
			}
		}

		if path != "/" {
			http.NotFound(w, r)
			return
		}

		posts := services.GetPosts(lang)

		pages.HomePage(lang, r.URL.Path, r.URL.RequestURI(), posts).Render(r.Context(), w)
	}
}
