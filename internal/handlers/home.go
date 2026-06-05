package handlers

import (
	"net/http"

	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/views/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts := services.GetPosts()
	tags := services.GetUniqueTags()

	pages.HomePage(r.URL.Path, posts, tags).Render(r.Context(), w)
}
