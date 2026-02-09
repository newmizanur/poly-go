package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/newmizanur/poly-go/internal/transpile"
	"github.com/newmizanur/poly-go/internal/workspace"
)

const version = "0.0.1"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	if cmd == "-h" || cmd == "--help" || cmd == "help" {
		usage()
		return
	}
	if cmd == "-v" || cmd == "--version" {
		fmt.Println("pgo", version)
		return
	}
	switch cmd {
	case "version":
		fmt.Println("pgo", version)
		return
	case "clean":
		if err := runClean(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	case "gen":
		lang, mapPath, allowGo, _, err := parseFlags(cmd, os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := runGen(lang, mapPath, allowGo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	case "build", "run", "test":
		lang, mapPath, allowGo, goArgs, err := parseFlags(cmd, os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := runGo(cmd, goArgs, lang, mapPath, allowGo); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	case "set":
		lang, _, _, rest, err := parseFlags(cmd, os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if lang == "" && len(rest) > 0 {
			lang = rest[0]
		}
		if lang == "" {
			fmt.Fprintln(os.Stderr, "usage: pgo set <lang>")
			os.Exit(1)
		}
		if err := setDefaultLang(lang); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: pgo <gen|build|run|test|clean|version|set> [--lang=<locale>] [--map=<path>] [--allow-go] [args...]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  gen       generate .pgo_gen")
	fmt.Fprintln(os.Stderr, "  build     build module via .pgo_gen")
	fmt.Fprintln(os.Stderr, "  run       run module or files via .pgo_gen")
	fmt.Fprintln(os.Stderr, "  test      test module via .pgo_gen")
	fmt.Fprintln(os.Stderr, "  clean     remove .pgo_gen")
	fmt.Fprintln(os.Stderr, "  set       set default locale in .pgo_lang")
	fmt.Fprintln(os.Stderr, "  version   print version")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "flags:")
	fmt.Fprintln(os.Stderr, "  --lang     locale override (e.g. bn, es, jp, zh)")
	fmt.Fprintln(os.Stderr, "  --map      custom keyword map path")
	fmt.Fprintln(os.Stderr, "  --allow-go allow Go keywords in .p.go")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "examples:")
	fmt.Fprintln(os.Stderr, "  pgo run --lang=bn ./examples/bn.p.go")
	fmt.Fprintln(os.Stderr, "  pgo set jp")
}

func runClean() error {
	moduleRoot, err := workspace.FindModuleRoot(mustGetwd())
	if err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(moduleRoot, workspace.GeneratedDirName))
}

func runGen(lang string, mapPath string, allowGo bool) error {
	moduleRoot, err := workspace.FindModuleRoot(mustGetwd())
	if err != nil {
		return err
	}
	maps, resolvedLang, err := loadMaps(moduleRoot, lang, mapPath, allowGo)
	if err != nil {
		return err
	}
	return workspace.Generate(moduleRoot, maps, resolvedLang)
}

func runGo(subcmd string, args []string, lang string, mapPath string, allowGo bool) error {
	moduleRoot, err := workspace.FindModuleRoot(mustGetwd())
	if err != nil {
		return err
	}
	maps, resolvedLang, err := loadMaps(moduleRoot, lang, mapPath, allowGo)
	if err != nil {
		return err
	}
	if err := workspace.Generate(moduleRoot, maps, resolvedLang); err != nil {
		return err
	}

	goArgs := mapArgsForGenerated(args)
	cmd := exec.Command("go", append([]string{subcmd}, goArgs...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = filepath.Join(moduleRoot, workspace.GeneratedDirName)
	cmd.Env = os.Environ()
	return cmd.Run()
}

func mapArgsForGenerated(args []string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasSuffix(arg, ".p.go") {
			out = append(out, strings.TrimSuffix(arg, ".p.go")+"_p.go")
			continue
		}
		out = append(out, arg)
	}
	return out
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return wd
}

func parseFlags(_ string, args []string) (string, string, bool, []string, error) {
	var lang string
	var mapPath string
	var allowGo bool
	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--lang=") {
			lang = strings.TrimPrefix(arg, "--lang=")
			continue
		}
		if arg == "--lang" {
			if i+1 >= len(args) {
				return "", "", false, nil, fmt.Errorf("missing value for --lang")
			}
			lang = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(arg, "--map=") {
			mapPath = strings.TrimPrefix(arg, "--map=")
			continue
		}
		if arg == "--map" {
			if i+1 >= len(args) {
				return "", "", false, nil, fmt.Errorf("missing value for --map")
			}
			mapPath = args[i+1]
			i++
			continue
		}
		if arg == "--allow-go" {
			allowGo = true
			continue
		}
		rest = append(rest, arg)
	}
	return lang, mapPath, allowGo, rest, nil
}

func loadMaps(moduleRoot string, langFlag string, mapPath string, allowGo bool) (transpile.Maps, string, error) {
	resolvedLang, err := resolveLang(moduleRoot, langFlag)
	if err != nil {
		return transpile.Maps{}, "", err
	}
	if mapPath != "" {
		data, err := os.ReadFile(mapPath)
		if err != nil {
			return transpile.Maps{}, "", err
		}
		maps, err := transpile.LoadKeywordMapData(data, allowGo)
		return maps, resolvedLang, err
	}
	data, err := keywordMapData(moduleRoot, resolvedLang)
	if err != nil {
		return transpile.Maps{}, "", err
	}
	maps, err := transpile.LoadKeywordMapData(data, allowGo)
	return maps, resolvedLang, err
}

func resolveLang(moduleRoot, langFlag string) (string, error) {
	if langFlag != "" {
		return langFlag, nil
	}
	if env := strings.TrimSpace(os.Getenv("PGO_LANG")); env != "" {
		return env, nil
	}
	if env := strings.TrimSpace(os.Getenv("POLYGO_LANG")); env != "" {
		return env, nil
	}
	if env := strings.TrimSpace(os.Getenv("BGO_LANG")); env != "" {
		return env, nil
	}
	path := filepath.Join(moduleRoot, ".pgo_lang")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if lang := strings.TrimSpace(string(data)); lang != "" {
		return lang, nil
	}
	return "", nil
}

func keywordMapData(moduleRoot, lang string) ([]byte, error) {
	if lang != "" {
		if data, err := os.ReadFile(filepath.Join(moduleRoot, "lang", lang+".json")); err == nil {
			return data, nil
		} else if !os.IsNotExist(err) {
			return nil, err
		}
		if data, ok := transpile.EmbeddedKeywordMap(lang); ok {
			return data, nil
		}
		return nil, fmt.Errorf("keyword map for lang %q not found (moduleRoot/lang or embedded)", lang)
	}
	if data, err := os.ReadFile(filepath.Join(moduleRoot, "keywords.json")); err == nil {
		return data, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if data, ok := transpile.EmbeddedKeywordMap(transpile.DefaultLocale); ok {
		return data, nil
	}
	return nil, fmt.Errorf("no embedded keyword maps found")
}

func setDefaultLang(lang string) error {
	moduleRoot, err := workspace.FindModuleRoot(mustGetwd())
	if err != nil {
		return err
	}
	if !localeAvailable(moduleRoot, lang) {
		return fmt.Errorf("unknown locale %q (no moduleRoot/lang or embedded map)", lang)
	}
	path := filepath.Join(moduleRoot, ".pgo_lang")
	return os.WriteFile(path, []byte(strings.TrimSpace(lang)+"\n"), 0o644)
}

func localeAvailable(moduleRoot, lang string) bool {
	if lang == "" {
		return false
	}
	if _, err := os.Stat(filepath.Join(moduleRoot, "lang", lang+".json")); err == nil {
		return true
	}
	_, ok := transpile.EmbeddedKeywordMap(lang)
	return ok
}
