package utils

import (
	"fmt"
	"time"

	"gabsdev-go/internal/i18n"
)

// RelativeTime formats a past timestamp as a short relative string, matching
// the density GitHub's own activity feed uses ("2h" rather than "2 hours"),
// localized per lang since the grammar (not just the words) differs.
func RelativeTime(t time.Time, lang i18n.Lang) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		if lang == i18n.LangPT {
			return "agora"
		}
		return "now"
	case d < time.Hour:
		return relValue(lang, int(d.Minutes()), "min")
	case d < 24*time.Hour:
		return relValue(lang, int(d.Hours()), "h")
	case d < 30*24*time.Hour:
		return relValue(lang, int(d.Hours()/24), "d")
	default:
		return relValue(lang, int(d.Hours()/24/30), "mo")
	}
}

func relValue(lang i18n.Lang, n int, unit string) string {
	if lang == i18n.LangPT {
		return fmt.Sprintf("há %d%s", n, unit)
	}
	return fmt.Sprintf("%d%s ago", n, unit)
}
