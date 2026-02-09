package workspace

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/newmizanur/poly-go/internal/transpile"
)

const GeneratedDirName = ".pgo_gen"

func FindModuleRoot(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func Generate(moduleRoot string, maps transpile.Maps, locale string) error {
	genDir := filepath.Join(moduleRoot, GeneratedDirName)
	if err := os.RemoveAll(genDir); err != nil {
		return err
	}
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return err
	}

	if err := copyIfExists(filepath.Join(moduleRoot, "go.mod"), filepath.Join(genDir, "go.mod")); err != nil {
		return err
	}
	if err := copyIfExists(filepath.Join(moduleRoot, "go.sum"), filepath.Join(genDir, "go.sum")); err != nil {
		return err
	}

	return filepath.WalkDir(moduleRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == moduleRoot {
			return nil
		}

		name := d.Name()
		if d.IsDir() {
			if name == GeneratedDirName || name == ".git" || name == "vendor" {
				return fs.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(moduleRoot, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(genDir, rel)

		if filepath.Base(path) == "go.mod" || filepath.Base(path) == "go.sum" {
			return nil
		}

		if strings.HasSuffix(path, ".p.go") {
			if !shouldIncludeLocalized(rel, locale) {
				return nil
			}
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			out, err := transpile.TranspileFileLocalizedToGo(rel, src, maps)
			if err != nil {
				return err
			}
			outPath := filepath.Join(genDir, bgToGoOutput(rel))
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}
			return os.WriteFile(outPath, out, 0o644)
		}

		if filepath.Ext(path) == ".go" {
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if transpile.ContainsLocalizedKeywords(src, maps) {
				return fmt.Errorf("localized keywords found in %s; rename file to *.p.go so it can be transpiled", rel)
			}
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return err
			}
			return copyFile(path, dest)
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		return copyFile(path, dest)
	})
}

func shouldIncludeLocalized(rel string, locale string) bool {
	if locale == "" {
		return true
	}
	parts := strings.Split(rel, string(filepath.Separator))
	inExamples := false
	inTestdata := false
	for _, part := range parts {
		if part == "examples" {
			inExamples = true
		} else if part == "testdata" {
			inTestdata = true
		}
		if part == locale {
			return true
		}
	}
	if !inExamples && !inTestdata {
		return true
	}
	base := filepath.Base(rel)
	if strings.HasPrefix(base, locale+".") || strings.HasPrefix(base, locale+"_") {
		return true
	}
	return false
}

func bgToGoOutput(rel string) string {
	if !strings.HasSuffix(rel, ".p.go") {
		return rel
	}
	base := strings.TrimSuffix(rel, ".p.go")
	return base + "_p.go"
}

func copyIfExists(src, dest string) error {
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return copyFile(src, dest)
}

func copyFile(src, dest string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	info, err := s.Stat()
	if err != nil {
		return err
	}

	d, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
