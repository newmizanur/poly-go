package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/newmizanur/poly-go/internal/transpile"
)

func main() {
	repoRoot, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	langDir := filepath.Join(repoRoot, "lang")
	locales, err := listLocales(langDir)
	if err != nil {
		fatal(err)
	}
	if len(locales) == 0 {
		fatal(fmt.Errorf("no locales found in %s", langDir))
	}

	templatePath := filepath.Join(repoRoot, "tools", "gen_features", "template.pgo.txt")
	input, err := os.ReadFile(templatePath)
	if err != nil {
		fatal(err)
	}

	featuresRoot := filepath.Join(repoRoot, "testdata", "features")
	for _, locale := range locales {
		mapPath := filepath.Join(langDir, locale+".json")
		mapData, err := os.ReadFile(mapPath)
		if err != nil {
			fatal(err)
		}
		maps, err := transpile.LoadKeywordMapData(mapData, false)
		if err != nil {
			fatal(err)
		}

		localized, err := transpile.TranspileFile(templatePath, input, maps)
		if err != nil {
			fatal(err)
		}
		localized = bytes.ReplaceAll(localized, []byte("\r\n"), []byte("\n"))

		localeDir := filepath.Join(featuresRoot, locale)
		expectedDir := filepath.Join(localeDir, ".expected")
		if err := os.MkdirAll(expectedDir, 0o755); err != nil {
			fatal(err)
		}
		if err := os.WriteFile(filepath.Join(localeDir, "main.p.go"), localized, 0o644); err != nil {
			fatal(err)
		}
		if err := os.WriteFile(filepath.Join(expectedDir, "main_p.go"), normalize(input), 0o644); err != nil {
			fatal(err)
		}
	}
}

func listLocales(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	locales := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		locales = append(locales, strings.TrimSuffix(name, ".json"))
	}
	return locales, nil
}

func normalize(src []byte) []byte {
	return bytes.ReplaceAll(src, []byte("\r\n"), []byte("\n"))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
