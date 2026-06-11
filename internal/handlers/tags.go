package handlers

import (
	"net/http"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/views/pages"
)

func Tags(lang i18n.Lang) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := services.GetUniqueTags(lang)
		pages.TagsPage(lang, r.URL.Path, r.URL.RequestURI(), tags).Render(r.Context(), w)
	}
}
