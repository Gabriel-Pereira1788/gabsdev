package handlers

import (
	services "gabsdev-go/internal/services/posts"
	"gabsdev-go/internal/utils"
	"gabsdev-go/internal/views/components"
	"gabsdev-go/internal/views/pages"
	"net/http"
)

func Articles(w http.ResponseWriter, r *http.Request) {
	activeTag := r.URL.Query().Get("tag")
	posts := services.GetPostsByTag(activeTag)
	tags := services.GetUniqueTags()

	pages.ArticlesPage(r.URL.Path, posts, tags, activeTag).Render(r.Context(), w)
}

func ArticlesSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tag := r.URL.Query().Get("tag")

	posts := services.SearchPosts(query, tag)

	components.PostList(posts, tag).Render(r.Context(), w)
}

func ArticleDetail(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, ok := services.GetPost(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	renderedContent := utils.RenderMarkdown(post.Content)
	toc := utils.ExtractTOC(post.Content)
	prev, next := services.GetAdjacentPosts(slug)

	pages.ArticlesDetailPage(r.URL.Path, post, renderedContent, toc, prev, next).Render(r.Context(), w)
}
