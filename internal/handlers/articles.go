package handlers

import (
	"net/http"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/utils"
	"gabsdev-go/internal/views/components"
	"gabsdev-go/internal/views/pages"
)

func Articles(lang i18n.Lang) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activeTag := r.URL.Query().Get("tag")
		posts := services.GetPostsByTag(lang, activeTag)
		tags := services.GetUniqueTags(lang)

		pages.ArticlesPage(lang, r.URL.Path, r.URL.RequestURI(), posts, tags, activeTag).Render(r.Context(), w)
	}
}

func ArticlesSearch(lang i18n.Lang) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		tag := r.URL.Query().Get("tag")

		posts := services.SearchPosts(lang, query, tag)

		components.PostList(lang, posts, tag).Render(r.Context(), w)
	}
}

func ArticleDetail(lang i18n.Lang) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		post, ok := services.GetPost(lang, slug)
		if !ok {
			http.NotFound(w, r)
			return
		}

		renderedContent := utils.RenderMarkdown(post.Content)
		toc := utils.ExtractTOC(post.Content)
		prev, next := services.GetAdjacentPosts(lang, slug)

		pages.ArticlesDetailPage(lang, r.URL.Path, r.URL.RequestURI(), post, renderedContent, toc, prev, next).Render(r.Context(), w)
	}
}
