package transpile

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

const DefaultLocale = "bn"

//go:embed lang/*.json
var embeddedKeywordMaps embed.FS

func EmbeddedKeywordMap(locale string) ([]byte, bool) {
	if locale == "" {
		locale = DefaultLocale
	}
	path := filepath.ToSlash(filepath.Join("lang", locale+".json"))
	data, err := embeddedKeywordMaps.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return data, true
}

func EmbeddedLocales() []string {
	matches, err := fs.Glob(embeddedKeywordMaps, "lang/*.json")
	if err != nil || len(matches) == 0 {
		return nil
	}
	locales := make([]string, 0, len(matches))
	for _, match := range matches {
		base := filepath.Base(match)
		locales = append(locales, strings.TrimSuffix(base, filepath.Ext(base)))
	}
	return locales
}
