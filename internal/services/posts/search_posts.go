package services

import (
	"gabsdev-go/internal/domain"
	"strings"
)

func SearchPosts(query, tag string) []domain.Post {
	posts := GetPostsByTag(tag)
	if query == "" {
		return posts
	}
	q := strings.ToLower(query)
	var result []domain.Post
	for _, p := range posts {
		if strings.Contains(strings.ToLower(p.Title), q) ||
			strings.Contains(strings.ToLower(p.Excerpt), q) ||
			tagsContain(p.Tags, q) {
			result = append(result, p)
		}
	}
	return result
}

func tagsContain(tags []string, q string) bool {
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), q) {
			return true
		}
	}
	return false
}
