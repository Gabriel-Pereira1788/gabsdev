package services

import "gabsdev-go/internal/domain"

func GetPostsByTag(tag string) []domain.Post {
	posts := GetPosts()
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
