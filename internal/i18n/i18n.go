package i18n

import (
	"embed"
	"encoding/json"
	"strings"
)

//go:embed locales/*.json
var localesFS embed.FS

type Lang string

const (
	LangPT      Lang = "pt"
	LangEN      Lang = "en"
	DefaultLang Lang = LangPT
)

var translations map[Lang]map[string]string

func init() {
	translations = make(map[Lang]map[string]string)
	for _, lang := range []Lang{LangPT, LangEN} {
		data, err := localesFS.ReadFile("locales/" + string(lang) + ".json")
		if err != nil {
			translations[lang] = map[string]string{}
			continue
		}
		var m map[string]string
		if err := json.Unmarshal(data, &m); err != nil {
			translations[lang] = map[string]string{}
			continue
		}
		translations[lang] = m
	}
}

func T(lang Lang, key string) string {
	if m, ok := translations[lang]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	if lang != LangPT {
		if m, ok := translations[LangPT]; ok {
			if v, ok := m[key]; ok {
				return v
			}
		}
	}
	return key
}

func Prefix(lang Lang) string {
	if lang == LangEN {
		return "/en"
	}
	return ""
}

func Path(lang Lang, path string) string {
	if lang == LangEN {
		return "/en" + path
	}
	return path
}

func ToggleURL(lang Lang, fullPath string) string {
	path := fullPath
	query := ""
	if idx := strings.Index(fullPath, "?"); idx >= 0 {
		path = fullPath[:idx]
		query = fullPath[idx:]
	}

	var newPath string
	if lang == LangPT {
		if path == "/" {
			newPath = "/en"
		} else {
			newPath = "/en" + path
		}
	} else {
		if path == "/en" || path == "/en/" {
			newPath = "/"
		} else if strings.HasPrefix(path, "/en/") {
			newPath = path[3:]
		} else {
			newPath = path
		}
	}

	return newPath + query
}

func HTMLLang(lang Lang) string {
	if lang == LangEN {
		return "en"
	}
	return "pt-BR"
}
