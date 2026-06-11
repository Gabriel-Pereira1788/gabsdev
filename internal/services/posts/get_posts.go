package services

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gabsdev-go/internal/domain"
	"gabsdev-go/internal/i18n"

	"gopkg.in/yaml.v3"
)

func postsDir(lang i18n.Lang) string {
	return filepath.Join("content/posts", string(lang))
}

func GetPosts(lang i18n.Lang) []domain.Post {
	dir := postsDir(lang)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []domain.Post{}
	}

	var posts []domain.Post
	for _, entry := range entries {
		name := entry.Name()

		if !strings.HasSuffix(name, ".md") && !strings.HasSuffix(name, ".mdx") {
			continue
		}
		post, err := parsePost(lang, name)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	return posts
}

func parsePost(lang i18n.Lang, filename string) (domain.Post, error) {
	fullPath := filepath.Join(postsDir(lang), filename)

	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return domain.Post{}, err
	}

	parts := strings.SplitN(string(raw), "---", 3)

	var meta domain.PostMeta
	if len(parts) >= 3 {
		yaml.Unmarshal([]byte(parts[1]), &meta)
	}

	content := ""
	if len(parts) == 3 {
		content = strings.TrimSpace(parts[2])
	}

	if meta.Slug == "" {
		meta.Slug = strings.TrimSuffix(strings.TrimSuffix(filename, ".mdx"), ".md")
	}

	if meta.ReadMin == 0 {
		words := len(strings.Fields(content))
		meta.ReadMin = words / 200
		if meta.ReadMin < 1 {
			meta.ReadMin = 1
		}
	}

	return domain.Post{PostMeta: meta, Content: content}, nil
}
