package utils

func IconSVG(name string) string {
	switch name {
	case "terminal":
		return `<rect x="1" y="2" width="14" height="12"/>` +
			`<rect x="2" y="4" width="12" height="9" fill="var(--bg-1)"/>` +
			`<rect x="3" y="6" width="2" height="2"/>` +
			`<rect x="5" y="8" width="2" height="2"/>` +
			`<rect x="3" y="10" width="4" height="1"/>`
	case "tag":
		return `<rect x="2" y="2" width="8" height="8"/>` +
			`<rect x="9" y="9" width="5" height="5"/>` +
			`<rect x="4" y="4" width="2" height="2" fill="var(--bg-1)"/>` +
			`<rect x="7" y="7" width="6" height="6"/>` +
			`<rect x="9" y="9" width="2" height="2" fill="var(--bg-1)"/>`
	case "clock":
		return `<rect x="5" y="1" width="6" height="1"/>` +
			`<rect x="3" y="2" width="2" height="1"/>` +
			`<rect x="11" y="2" width="2" height="1"/>` +
			`<rect x="2" y="3" width="1" height="2"/>` +
			`<rect x="13" y="3" width="1" height="2"/>` +
			`<rect x="1" y="5" width="1" height="6"/>` +
			`<rect x="14" y="5" width="1" height="6"/>` +
			`<rect x="2" y="11" width="1" height="2"/>` +
			`<rect x="13" y="11" width="1" height="2"/>` +
			`<rect x="3" y="13" width="2" height="1"/>` +
			`<rect x="11" y="13" width="2" height="1"/>` +
			`<rect x="5" y="14" width="6" height="1"/>` +
			`<rect x="7" y="4" width="2" height="5"/>` +
			`<rect x="8" y="7" width="3" height="2"/>`
	case "moon":
		return `<rect x="6" y="2" width="6" height="2"/>` +
			`<rect x="4" y="3" width="2" height="2"/>` +
			`<rect x="3" y="5" width="2" height="6"/>` +
			`<rect x="4" y="11" width="2" height="2"/>` +
			`<rect x="6" y="12" width="6" height="2"/>` +
			`<rect x="9" y="4" width="4" height="8" fill="var(--bg-1)"/>`
	case "sun":
		return `<rect x="7" y="1" width="2" height="2"/>` +
			`<rect x="7" y="13" width="2" height="2"/>` +
			`<rect x="1" y="7" width="2" height="2"/>` +
			`<rect x="13" y="7" width="2" height="2"/>` +
			`<rect x="3" y="3" width="2" height="2"/>` +
			`<rect x="11" y="3" width="2" height="2"/>` +
			`<rect x="3" y="11" width="2" height="2"/>` +
			`<rect x="11" y="11" width="2" height="2"/>` +
			`<rect x="5" y="5" width="6" height="6"/>` +
			`<rect x="7" y="7" width="2" height="2" fill="var(--bg-1)"/>`
	case "arrow":
		return `<rect x="2" y="7" width="9" height="2"/>` +
			`<rect x="8" y="4" width="2" height="2"/>` +
			`<rect x="10" y="6" width="2" height="2"/>` +
			`<rect x="10" y="8" width="2" height="2"/>` +
			`<rect x="8" y="10" width="2" height="2"/>`
	case "back":
		return `<rect x="5" y="7" width="9" height="2"/>` +
			`<rect x="6" y="4" width="2" height="2"/>` +
			`<rect x="4" y="6" width="2" height="2"/>` +
			`<rect x="4" y="8" width="2" height="2"/>` +
			`<rect x="6" y="10" width="2" height="2"/>`
	case "rss":
		return `<rect x="2" y="11" width="3" height="3"/>` +
			`<rect x="2" y="7" width="2" height="2"/>` +
			`<rect x="4" y="7" width="2" height="2"/>` +
			`<rect x="6" y="9" width="2" height="2"/>` +
			`<rect x="6" y="11" width="2" height="3"/>` +
			`<rect x="2" y="3" width="2" height="2"/>` +
			`<rect x="4" y="3" width="2" height="2"/>` +
			`<rect x="6" y="4" width="2" height="2"/>` +
			`<rect x="8" y="6" width="2" height="2"/>` +
			`<rect x="10" y="8" width="2" height="2"/>` +
			`<rect x="10" y="11" width="2" height="3"/>`
	case "github":
		return `<rect x="4" y="2" width="8" height="2"/>` +
			`<rect x="2" y="4" width="2" height="6"/>` +
			`<rect x="12" y="4" width="2" height="6"/>` +
			`<rect x="4" y="10" width="2" height="4"/>` +
			`<rect x="10" y="10" width="2" height="4"/>` +
			`<rect x="6" y="12" width="4" height="2"/>` +
			`<rect x="5" y="5" width="2" height="2" fill="var(--bg-1)"/>` +
			`<rect x="9" y="5" width="2" height="2" fill="var(--bg-1)"/>`
	case "folder":
		return `<rect x="1" y="3" width="6" height="2"/>` +
			`<rect x="1" y="5" width="14" height="9"/>` +
			`<rect x="2" y="7" width="12" height="6" fill="var(--bg-1)"/>`
	case "doc":
		return `<rect x="3" y="1" width="7" height="14"/>` +
			`<rect x="10" y="4" width="3" height="11"/>` +
			`<rect x="10" y="1" width="1" height="3"/>` +
			`<rect x="11" y="2" width="1" height="2"/>` +
			`<rect x="12" y="3" width="1" height="1"/>` +
			`<rect x="5" y="4" width="4" height="1" fill="var(--bg-1)"/>` +
			`<rect x="5" y="6" width="6" height="1" fill="var(--bg-1)"/>` +
			`<rect x="5" y="8" width="6" height="1" fill="var(--bg-1)"/>` +
			`<rect x="5" y="10" width="4" height="1" fill="var(--bg-1)"/>`
	case "hash":
		return `<rect x="4" y="2" width="2" height="12"/>` +
			`<rect x="9" y="2" width="2" height="12"/>` +
			`<rect x="2" y="5" width="12" height="2"/>` +
			`<rect x="2" y="9" width="12" height="2"/>`
	case "heart":
		return `<rect x="2" y="3" width="4" height="2"/>` +
			`<rect x="10" y="3" width="4" height="2"/>` +
			`<rect x="1" y="5" width="6" height="4"/>` +
			`<rect x="9" y="5" width="6" height="4"/>` +
			`<rect x="2" y="9" width="12" height="2"/>` +
			`<rect x="3" y="11" width="10" height="2"/>` +
			`<rect x="5" y="13" width="6" height="1"/>` +
			`<rect x="6" y="14" width="4" height="1"/>`
	case "linkedin":
		return `<rect x="2" y="2" width="3" height="3"/>` +
			`<rect x="2" y="6" width="3" height="8"/>` +
			`<rect x="7" y="6" width="3" height="8"/>` +
			`<rect x="7" y="6" width="7" height="3"/>` +
			`<rect x="11" y="6" width="3" height="8"/>`
	}
	return ""
}
