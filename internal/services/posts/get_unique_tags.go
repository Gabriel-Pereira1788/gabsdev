package services

import "gabsdev-go/internal/i18n"

func GetUniqueTags(lang i18n.Lang) []string {
	posts := GetPosts(lang)
	seen := map[string]bool{}
	var tags []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	return tags
}
