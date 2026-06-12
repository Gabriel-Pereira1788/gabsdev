package services_test

import (
	"testing"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
)

func TestGetUniqueTags_NoDuplicates(t *testing.T) {
	tags := services.GetUniqueTags(i18n.LangEN)
	seen := map[string]bool{}
	for _, tag := range tags {
		if seen[tag] {
			t.Errorf("duplicate tag %q in result", tag)
		}
		seen[tag] = true
	}
}

func TestGetUniqueTags_SupersetOfSinglePostTags(t *testing.T) {
	tags := services.GetUniqueTags(i18n.LangEN)
	set := map[string]bool{}
	for _, tag := range tags {
		set[tag] = true
	}

	for _, p := range services.GetPosts(i18n.LangEN) {
		for _, tag := range p.Tags {
			if !set[tag] {
				t.Errorf("tag %q from post %q missing in GetUniqueTags", tag, p.Slug)
			}
		}
	}
}
