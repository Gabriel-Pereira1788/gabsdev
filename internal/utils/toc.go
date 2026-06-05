package utils

import (
	"regexp"
	"strings"
)

type TOCItem struct {
	ID    string
	Title string
}

var h2Regex = regexp.MustCompile(`(?m)^## (.+)$`)

func ExtractTOC(content string) []TOCItem {
	matches := h2Regex.FindAllStringSubmatch(content, -1)
	items := make([]TOCItem, 0, len(matches))
	for _, m := range matches {
		title := strings.TrimSpace(m[1])
		id := slugifyHeading(title)
		items = append(items, TOCItem{ID: id, Title: title})
	}
	return items
}

var nonAlphanumRegex = regexp.MustCompile(`[^a-z0-9\s-]`)
var whitespaceRegex = regexp.MustCompile(`\s+`)

func slugifyHeading(s string) string {
	s = strings.ToLower(s)
	s = nonAlphanumRegex.ReplaceAllString(s, "")
	s = whitespaceRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
