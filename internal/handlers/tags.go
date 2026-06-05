package handlers

import (
	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/views/pages"
	"log"
	"net/http"
)

func Tags(w http.ResponseWriter, r *http.Request) {
	tags := services.GetUniqueTags()
	log.Println("TAGS", tags)
	pages.TagsPage(r.URL.Path, tags).Render(r.Context(), w)
}
