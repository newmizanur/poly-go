package transpile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoldenTranspile(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	testdataRoot := filepath.Join(repoRoot, "testdata")
	localeSet := mustLocaleSet(filepath.Join(repoRoot, "lang"))

	err := filepath.WalkDir(testdataRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".p.go") {
			return nil
		}

		t.Run(relForTestName(repoRoot, path), func(t *testing.T) {
			locale := localeForPath(testdataRoot, path, localeSet)
			mapPath := filepath.Join(repoRoot, "lang", "bn.json")
			if locale != "" {
				mapPath = filepath.Join(repoRoot, "lang", locale+".json")
			}
			maps, err := LoadKeywordMapData(mustRead(mapPath), false)
			if err != nil {
				t.Fatalf("LoadKeywordMap(%s): %v", mapPath, err)
			}
			src, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read input %s: %v", path, err)
			}

			got, err := TranspileFile(path, src, maps)
			if err != nil {
				t.Fatalf("TranspileFile(%s): %v", path, err)
			}

			expectedPath := expectedForInput(path)
			want, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("read expected %s: %v", expectedPath, err)
			}

			if normalize(string(got)) != normalize(string(want)) {
				t.Fatalf(
					"golden mismatch\ninput: %s\nexpected: %s\n\n--- GOT ---\n%s\n\n--- WANT ---\n%s",
					path, expectedPath, got, want,
				)
			}
		})

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

func mustRead(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}

func expectedForInput(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	outBase := strings.Replace(base, ".p.go", "_p.go", 1)
	return filepath.Join(dir, ".expected", outBase)
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.Join(lines, "\n")
}

func relForTestName(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return rel
	}
	return path
}

func mustLocaleSet(langDir string) map[string]struct{} {
	entries, err := os.ReadDir(langDir)
	if err != nil {
		panic(err)
	}
	set := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		set[strings.TrimSuffix(name, ".json")] = struct{}{}
	}
	return set
}

func localeForPath(root, path string, locales map[string]struct{}) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		if _, ok := locales[part]; ok {
			return part
		}
	}
	return ""
}
