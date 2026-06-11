package services

import (
	"gabsdev-go/internal/domain"
	"gabsdev-go/internal/i18n"
)

func GetPost(lang i18n.Lang, slug string) (domain.Post, bool) {
	for _, post := range GetPosts(lang) {
		if post.Slug == slug {
			return post, true
		}
	}
	return domain.Post{}, false
}

func GetAdjacentPosts(lang i18n.Lang, slug string) (prev, next *domain.Post) {
	posts := GetPosts(lang)
	for i, p := range posts {
		if p.Slug == slug {
			if i > 0 {
				prev = &posts[i-1]
			}
			if i < len(posts)-1 {
				next = &posts[i+1]
			}
			return
		}
	}
	return
}
