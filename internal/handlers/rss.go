package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gabsdev-go/internal/i18n"
	services "gabsdev-go/internal/services/posts"
)

func RSS(w http.ResponseWriter, r *http.Request) {
	siteURL := os.Getenv("SITE_URL")
	if siteURL == "" {
		siteURL = "https://gabsdev.dev"
	}

	posts := services.GetPosts(i18n.LangPT)

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")

	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fmt.Fprintf(&sb, "<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n")
	sb.WriteString("<channel>\n")
	sb.WriteString("<title>gabsdev</title>\n")
	fmt.Fprintf(&sb, "<link>%s</link>\n", siteURL)
	sb.WriteString("<description>Notas de engenharia, anotacoes de terminal e artigos tecnicos condensados.</description>\n")
	sb.WriteString("<language>pt-BR</language>\n")
	fmt.Fprintf(&sb, "<atom:link href=\"%s/rss.xml\" rel=\"self\" type=\"application/rss+xml\"/>\n", siteURL)

	for _, post := range posts {
		t, err := time.Parse("2006-01-02", post.Date)
		pubDate := ""
		if err == nil {
			pubDate = t.Format(time.RFC1123Z)
		}

		sb.WriteString("<item>\n")
		sb.WriteString("<title>")
		xml.EscapeText(&sb, []byte(post.Title))
		sb.WriteString("</title>\n")
		fmt.Fprintf(&sb, "<link>%s/articles/%s</link>\n", siteURL, post.Slug)
		fmt.Fprintf(&sb, "<guid isPermaLink=\"true\">%s/articles/%s</guid>\n", siteURL, post.Slug)
		sb.WriteString("<description>")
		xml.EscapeText(&sb, []byte(post.Excerpt))
		sb.WriteString("</description>\n")
		if pubDate != "" {
			fmt.Fprintf(&sb, "<pubDate>%s</pubDate>\n", pubDate)
		}
		for _, tag := range post.Tags {
			sb.WriteString("<category>")
			xml.EscapeText(&sb, []byte(tag))
			sb.WriteString("</category>\n")
		}
		sb.WriteString("</item>\n")
	}

	sb.WriteString("</channel>\n")
	sb.WriteString("</rss>")

	fmt.Fprint(w, sb.String())
}
