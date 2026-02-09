# PolyGo System Design

## Overview

PolyGo is a small Go CLI that lets developers write Go code using localized keywords in `.p.go` files. It transpiles localized source into standard Go, then executes the official Go toolchain inside a generated workspace.

Key goals:
- Keep Go semantics intact (no fork).
- Support multiple languages through data‑driven keyword maps.
- Keep tooling minimal and predictable.

## Architecture

### 1) Keyword maps
- Each locale is defined by `lang/<locale>.json` with:
  - `keywords`: localized tokens → Go keywords
  - `predeclared`: localized tokens → predeclared identifiers
- Maps are embedded into the binary (from `internal/transpile/lang/*.json`).
- `pgo set <locale>` writes `.pgo_lang` to select a default locale.

### 2) Transpiler
- A fast, single‑pass tokenizer scans identifiers.
- Identifiers are replaced if they match locale keywords/predeclared entries.
- Strings/comments are preserved.
- Escape prefix `@` allows using localized keywords as identifiers.

### 3) Workspace generation
- `.pgo_gen` is created at module root.
- Copies `go.mod`/`go.sum` and mirrors the directory tree.
- Transpiles `.p.go` → `_p.go` (to avoid name collisions).
- Normal `.go` files are copied as‑is.

### 4) CLI flow
Commands: `gen`, `build`, `run`, `test`, `clean`, `version`.

Flow:
1. Resolve locale and keyword map.
2. Generate `.pgo_gen`.
3. Run `go <cmd>` inside `.pgo_gen`.

## Locale resolution
Order of precedence:
1. `--lang=<locale>` flag
2. `PGO_LANG` / `POLYGO_LANG` / `BGO_LANG`
3. `.pgo_lang`
4. Embedded default locale

## Key design tradeoffs
- **Generated workspace** keeps Go tooling untouched.
- **Data‑driven maps** make new languages trivial to add.
- **Escape prefix** keeps keyword rules strict without blocking common names.

## Extensibility
- Add new locale: `lang/<locale>.json` + `examples/<locale>.p.go`.
- VS Code extension reads `.pgo_lang` for completions/highlighting.

## Non‑goals
- No custom compiler or runtime.
- No semantic linting beyond Go’s toolchain.

