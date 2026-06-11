package services

import (
	"gabsdev-go/internal/domain"
	"gabsdev-go/internal/i18n"
)

func GetPostsByTag(lang i18n.Lang, tag string) []domain.Post {
	posts := GetPosts(lang)
	if tag == "" {
		return posts
	}
	var filteredPosts []domain.Post
	for _, post := range posts {
		for _, postTag := range post.Tags {
			if postTag == tag {
				filteredPosts = append(filteredPosts, post)
			}
		}
	}
	return filteredPosts
}
