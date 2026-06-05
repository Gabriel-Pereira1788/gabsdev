package handlers

import (
	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/views/pages"
	"net/http"
)

func About(w http.ResponseWriter, r *http.Request) {
	postsCount := len(services.GetPosts())
	pages.AboutPage(r.URL.Path, postsCount).Render(r.Context(), w)
}
