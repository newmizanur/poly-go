package transpile

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type KeywordMap struct {
	Keywords    map[string]string `json:"keywords"`
	Predeclared map[string]string `json:"predeclared"`
}

type Maps struct {
	LocalToGo        map[string]string
	LocalPredeclared map[string]string
	GoToLocal        map[string]string
	GoPredeclared    map[string]string
	LocalAll         map[string]struct{}
	AllowGoKeywords  bool
}

func LoadKeywordMap(path string) (Maps, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Maps{}, err
	}
	return LoadKeywordMapData(data, false)
}

func LoadKeywordMapData(data []byte, allowGoKeywords bool) (Maps, error) {
	var km KeywordMap
	if err := json.Unmarshal(data, &km); err != nil {
		return Maps{}, err
	}
	maps := Maps{
		LocalToGo:        make(map[string]string),
		LocalPredeclared: make(map[string]string),
		GoToLocal:        make(map[string]string),
		GoPredeclared:    make(map[string]string),
		LocalAll:         make(map[string]struct{}),
		AllowGoKeywords:  allowGoKeywords,
	}
	for k, v := range km.Keywords {
		maps.LocalToGo[k] = v
		maps.GoToLocal[v] = k
		maps.LocalAll[k] = struct{}{}
	}
	for k, v := range km.Predeclared {
		maps.LocalPredeclared[k] = v
		maps.GoPredeclared[v] = k
		maps.LocalAll[k] = struct{}{}
	}
	return maps, nil
}

// ContainsLocalizedKeywords reports whether src includes any localized keyword/predeclared identifiers.
// It ignores strings/comments and respects the same identifier rules as the transpiler.
func ContainsLocalizedKeywords(src []byte, maps Maps) bool {
	found := false
	scanIdentifiers(src, func(ident string) {
		if _, ok := maps.GoToLocal[ident]; ok {
			return
		}
		if _, ok := maps.GoPredeclared[ident]; ok {
			return
		}
		if _, ok := maps.LocalAll[ident]; ok {
			found = true
		}
	})
	return found
}

type Direction int

const (
	LocalToGo Direction = iota
	GoToLocal
)

func TranspileFile(srcPath string, src []byte, maps Maps) ([]byte, error) {
	prefixLen := buildTagPrefixLen(src)
	prefix := src[:prefixLen]
	body := src[prefixLen:]

	direction, err := detectDirection(body, maps)
	if err != nil {
		return nil, err
	}

	transpiledBody, err := transpileBody(body, maps, direction)
	if err != nil {
		return nil, err
	}
	if direction == GoToLocal {
		transpiledBody = compactCommasInBraces(transpiledBody)
	}

	out := make([]byte, 0, len(prefix)+len(transpiledBody))
	out = append(out, prefix...)
	out = append(out, transpiledBody...)
	return out, nil
}

func TranspileFileLocalizedToGo(srcPath string, src []byte, maps Maps) ([]byte, error) {
	prefixLen := buildTagPrefixLen(src)
	prefix := src[:prefixLen]
	body := src[prefixLen:]

	transpiledBody, err := transpileBodyLocalizedToGo(body, maps)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 0, len(prefix)+len(transpiledBody))
	out = append(out, prefix...)
	out = append(out, transpiledBody...)
	return out, nil
}

func transpileBodyLocalizedToGo(body []byte, maps Maps) ([]byte, error) {
	out := make([]byte, 0, len(body))
	last := 0
	idx := 0
	escapedNames := make(map[string]struct{})

	for idx < len(body) {
		r, size := utf8.DecodeRune(body[idx:])
		if r == utf8.RuneError && size == 1 {
			return nil, fmt.Errorf("invalid utf-8 at %d", idx)
		}

		if r == '/' && idx+1 < len(body) {
			next := body[idx+1]
			if next == '/' {
				idx = skipLineComment(body, idx)
				continue
			}
			if next == '*' {
				idx = skipBlockComment(body, idx)
				continue
			}
		}

		if r == '"' {
			idx = skipInterpretedString(body, idx)
			continue
		}
		if r == '\'' {
			idx = skipRuneLiteral(body, idx)
			continue
		}
		if r == '`' {
			idx = skipRawString(body, idx)
			continue
		}

		if r == '@' {
			_, nsize := utf8.DecodeRune(body[idx+size:])
			if nsize > 0 {
				nr, _ := utf8.DecodeRune(body[idx+size:])
				if isIdentStart(nr) {
					identStart := idx + size
					identEnd := readIdent(body, identStart)
					out = append(out, body[last:idx]...)
					ident := string(body[identStart:identEnd])
					escapedNames[ident] = struct{}{}
					replacement := translateIdentLocalizedToGo(ident, maps, true, escapedNames)
					out = append(out, []byte(replacement)...)
					last = identEnd
					idx = identEnd
					continue
				}
			}
		}

		if isIdentStart(r) {
			identStart := idx
			identEnd := readIdent(body, identStart)
			out = append(out, body[last:identStart]...)
			ident := string(body[identStart:identEnd])
			if _, escaped := escapedNames[ident]; !escaped {
				if !maps.AllowGoKeywords {
					if _, ok := maps.LocalToGo[ident]; !ok {
						if _, ok := maps.LocalPredeclared[ident]; !ok {
							if _, ok := maps.GoToLocal[ident]; ok {
								return nil, fmt.Errorf("go keyword %q is not allowed in .p.go; use localized keyword", ident)
							}
							if _, ok := maps.GoPredeclared[ident]; ok {
								return nil, fmt.Errorf("go predeclared %q is not allowed in .p.go; use localized keyword", ident)
							}
						}
					}
				}
			}
			if ident == "চ্যানেল" && shouldDropChanKeyword(body, identEnd) {
				last = identEnd
				idx = identEnd
				continue
			}
			replacement := translateIdentLocalizedToGo(ident, maps, false, escapedNames)
			out = append(out, []byte(replacement)...)
			last = identEnd
			idx = identEnd
			continue
		}

		idx += size
	}

	out = append(out, body[last:]...)
	return out, nil
}

func translateIdentLocalizedToGo(ident string, maps Maps, escaped bool, escapedNames map[string]struct{}) string {
	if escaped {
		if !isValidGoIdent(ident) {
			return mangleIdent(ident)
		}
		return ident
	}
	if _, ok := escapedNames[ident]; ok {
		if !isValidGoIdent(ident) {
			return mangleIdent(ident)
		}
		return ident
	}
	if mapped, ok := maps.LocalToGo[ident]; ok {
		return mapped
	}
	if mapped, ok := maps.LocalPredeclared[ident]; ok {
		return mapped
	}
	if !isValidGoIdent(ident) {
		return mangleIdent(ident)
	}
	return ident
}

func shouldDropChanKeyword(src []byte, identEnd int) bool {
	idx := identEnd
	for idx < len(src) {
		if src[idx] == '\n' {
			return false
		}
		if isSpace(src[idx]) {
			idx++
			continue
		}
		if src[idx] == '/' && idx+1 < len(src) {
			if src[idx+1] == '/' {
				idx = skipLineComment(src, idx)
				return false
			}
			if src[idx+1] == '*' {
				idx = skipBlockComment(src, idx)
				continue
			}
		}
		break
	}
	if idx >= len(src) {
		return false
	}
	r, size := utf8.DecodeRune(src[idx:])
	if r == utf8.RuneError && size == 1 {
		return false
	}
	if !isIdentStart(r) {
		return false
	}
	idx = readIdent(src, idx)
	for idx < len(src) {
		if src[idx] == '\n' {
			return false
		}
		if isSpace(src[idx]) {
			idx++
			continue
		}
		if src[idx] == '/' && idx+1 < len(src) {
			if src[idx+1] == '/' {
				idx = skipLineComment(src, idx)
				return false
			}
			if src[idx+1] == '*' {
				idx = skipBlockComment(src, idx)
				continue
			}
		}
		break
	}
	return idx+1 < len(src) && src[idx] == ':' && src[idx+1] == '='
}

func isValidGoIdent(ident string) bool {
	if ident == "" {
		return false
	}
	i := 0
	for _, r := range ident {
		if i == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
		} else {
			if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return false
			}
		}
		i++
	}
	return true
}

func mangleIdent(ident string) string {
	sum := sha1.Sum([]byte(ident))
	// Use exported prefix so mangled package-level identifiers are accessible cross-package.
	return "Bgo_" + hex.EncodeToString(sum[:8])
}

func detectDirection(body []byte, maps Maps) (Direction, error) {
	seenLocal := false
	seenGo := false

	scanIdentifiers(body, func(ident string) {
		if _, ok := maps.LocalToGo[ident]; ok {
			seenLocal = true
		}
		if _, ok := maps.LocalPredeclared[ident]; ok {
			seenLocal = true
		}
		if _, ok := maps.GoToLocal[ident]; ok {
			seenGo = true
		}
		if _, ok := maps.GoPredeclared[ident]; ok {
			seenGo = true
		}
	})

	if seenGo {
		return GoToLocal, nil
	}
	if seenLocal {
		return LocalToGo, nil
	}
	return LocalToGo, nil
}

func transpileBody(body []byte, maps Maps, direction Direction) ([]byte, error) {
	out := make([]byte, 0, len(body))
	last := 0
	idx := 0

	for idx < len(body) {
		r, size := utf8.DecodeRune(body[idx:])
		if r == utf8.RuneError && size == 1 {
			return nil, fmt.Errorf("invalid utf-8 at %d", idx)
		}

		if r == '/' && idx+1 < len(body) {
			next := body[idx+1]
			if next == '/' {
				idx = skipLineComment(body, idx)
				continue
			}
			if next == '*' {
				idx = skipBlockComment(body, idx)
				continue
			}
		}

		if r == '"' {
			idx = skipInterpretedString(body, idx)
			continue
		}
		if r == '\'' {
			idx = skipRuneLiteral(body, idx)
			continue
		}
		if r == '`' {
			idx = skipRawString(body, idx)
			continue
		}

		if r == '@' {
			_, nsize := utf8.DecodeRune(body[idx+size:])
			if nsize > 0 {
				nr, _ := utf8.DecodeRune(body[idx+size:])
				if isIdentStart(nr) {
					identStart := idx + size
					identEnd := readIdent(body, identStart)
					out = append(out, body[last:idx]...)
					ident := string(body[identStart:identEnd])
					replacement := translateIdent(ident, maps, direction, true, false)
					out = append(out, []byte(replacement)...)
					last = identEnd
					idx = identEnd
					continue
				}
			}
		}

		if isIdentStart(r) {
			identStart := idx
			identEnd := readIdent(body, identStart)
			out = append(out, body[last:identStart]...)
			ident := string(body[identStart:identEnd])
			escapeNeeded := false
			if direction == GoToLocal {
				if _, ok := maps.LocalAll[ident]; ok {
					escapeNeeded = shouldEscapeGoIdent(body, identEnd)
				}
			}
			replacement := translateIdent(ident, maps, direction, false, escapeNeeded)
			out = append(out, []byte(replacement)...)
			last = identEnd
			idx = identEnd
			continue
		}

		idx += size
	}

	out = append(out, body[last:]...)
	return out, nil
}

func translateIdent(ident string, maps Maps, direction Direction, escaped bool, escapeNeeded bool) string {
	if escaped {
		return ident
	}
	if direction == LocalToGo {
		if mapped, ok := maps.LocalToGo[ident]; ok {
			return mapped
		}
		if mapped, ok := maps.LocalPredeclared[ident]; ok {
			return mapped
		}
		return ident
	}

	if mapped, ok := maps.GoPredeclared[ident]; ok {
		return mapped
	}
	if mapped, ok := maps.GoToLocal[ident]; ok {
		return mapped
	}
	if escapeNeeded {
		return "@" + ident
	}
	return ident
}

func shouldEscapeGoIdent(src []byte, identEnd int) bool {
	idx := identEnd
	for idx < len(src) {
		if isSpace(src[idx]) {
			idx++
			continue
		}
		if src[idx] == '/' && idx+1 < len(src) {
			if src[idx+1] == '/' {
				idx = skipLineComment(src, idx)
				continue
			}
			if src[idx+1] == '*' {
				idx = skipBlockComment(src, idx)
				continue
			}
		}
		if src[idx] == ':' && idx+1 < len(src) && src[idx+1] == '=' {
			return true
		}
		return false
	}
	return false
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func scanIdentifiers(src []byte, fn func(ident string)) {
	idx := 0
	for idx < len(src) {
		r, size := utf8.DecodeRune(src[idx:])
		if r == utf8.RuneError && size == 1 {
			return
		}

		if r == '/' && idx+1 < len(src) {
			next := src[idx+1]
			if next == '/' {
				idx = skipLineComment(src, idx)
				continue
			}
			if next == '*' {
				idx = skipBlockComment(src, idx)
				continue
			}
		}

		if r == '"' {
			idx = skipInterpretedString(src, idx)
			continue
		}
		if r == '\'' {
			idx = skipRuneLiteral(src, idx)
			continue
		}
		if r == '`' {
			idx = skipRawString(src, idx)
			continue
		}

		if r == '@' {
			_, nsize := utf8.DecodeRune(src[idx+size:])
			if nsize > 0 {
				nr, _ := utf8.DecodeRune(src[idx+size:])
				if isIdentStart(nr) {
					idx = readIdent(src, idx+size)
					continue
				}
			}
		}

		if isIdentStart(r) {
			end := readIdent(src, idx)
			fn(string(src[idx:end]))
			idx = end
			continue
		}

		idx += size
	}
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentPart(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.Is(unicode.Mark, r)
}

func readIdent(src []byte, start int) int {
	idx := start
	first := true
	for idx < len(src) {
		r, size := utf8.DecodeRune(src[idx:])
		if r == utf8.RuneError && size == 1 {
			return idx
		}
		if first {
			if !isIdentStart(r) {
				return idx
			}
			first = false
			idx += size
			continue
		}
		if !isIdentPart(r) {
			return idx
		}
		idx += size
	}
	return idx
}

func skipLineComment(src []byte, start int) int {
	idx := start
	for idx < len(src) {
		if src[idx] == '\n' {
			return idx + 1
		}
		idx++
	}
	return idx
}

func skipBlockComment(src []byte, start int) int {
	idx := start + 2
	for idx+1 < len(src) {
		if src[idx] == '*' && src[idx+1] == '/' {
			return idx + 2
		}
		idx++
	}
	return len(src)
}

func skipInterpretedString(src []byte, start int) int {
	idx := start + 1
	for idx < len(src) {
		r, size := utf8.DecodeRune(src[idx:])
		if r == '\\' {
			idx += size
			if idx < len(src) {
				_, s2 := utf8.DecodeRune(src[idx:])
				idx += s2
			}
			continue
		}
		if r == '"' {
			return idx + size
		}
		idx += size
	}
	return idx
}

func skipRuneLiteral(src []byte, start int) int {
	idx := start + 1
	for idx < len(src) {
		r, size := utf8.DecodeRune(src[idx:])
		if r == '\\' {
			idx += size
			if idx < len(src) {
				_, s2 := utf8.DecodeRune(src[idx:])
				idx += s2
			}
			continue
		}
		if r == '\'' {
			return idx + size
		}
		idx += size
	}
	return idx
}

func skipRawString(src []byte, start int) int {
	idx := start + 1
	for idx < len(src) {
		if src[idx] == '`' {
			return idx + 1
		}
		idx++
	}
	return idx
}

func buildTagPrefixLen(src []byte) int {
	reader := bufio.NewReader(strings.NewReader(string(src)))
	pos := 0
	lineNum := 0
	for {
		line, err := reader.ReadString('\n')
		lineNum++
		if err != nil && len(line) == 0 {
			break
		}
		trimmed := strings.TrimSpace(strings.TrimRight(line, "\r\n"))
		if lineNum == 1 && trimmed == "" {
			return 0
		}
		if strings.HasPrefix(trimmed, "//go:build") || strings.HasPrefix(trimmed, "// +build") {
			pos += len(line)
			if err != nil {
				break
			}
			continue
		}
		break
	}
	return pos
}

func compactCommasInBraces(src []byte) []byte {
	out := make([]byte, 0, len(src))
	idx := 0
	type tokKind int
	const (
		tokNone tokKind = iota
		tokIdent
		tokRBracket
		tokOther
	)
	prevTok := tokNone
	lastTok := tokNone
	var braceStack []bool
	for idx < len(src) {
		if src[idx] == '/' && idx+1 < len(src) {
			if src[idx+1] == '/' {
				end := skipLineComment(src, idx)
				out = append(out, src[idx:end]...)
				idx = end
				continue
			}
			if src[idx+1] == '*' {
				end := skipBlockComment(src, idx)
				out = append(out, src[idx:end]...)
				idx = end
				continue
			}
		}
		if src[idx] == '"' {
			end := skipInterpretedString(src, idx)
			out = append(out, src[idx:end]...)
			idx = end
			continue
		}
		if src[idx] == '\'' {
			end := skipRuneLiteral(src, idx)
			out = append(out, src[idx:end]...)
			idx = end
			continue
		}
		if src[idx] == '`' {
			end := skipRawString(src, idx)
			out = append(out, src[idx:end]...)
			idx = end
			continue
		}

		r, size := utf8.DecodeRune(src[idx:])
		if r == utf8.RuneError && size == 1 {
			out = append(out, src[idx])
			prevTok = lastTok
			lastTok = tokOther
			idx++
			continue
		}
		if isIdentStart(r) {
			end := readIdent(src, idx)
			out = append(out, src[idx:end]...)
			prevTok = lastTok
			lastTok = tokIdent
			idx = end
			continue
		}

		switch src[idx] {
		case ']':
			out = append(out, src[idx])
			prevTok = lastTok
			lastTok = tokRBracket
			idx++
			continue
		case '{':
			isComposite := prevTok == tokRBracket && lastTok == tokIdent
			braceStack = append(braceStack, isComposite)
			out = append(out, src[idx])
			prevTok = lastTok
			lastTok = tokOther
			idx++
			continue
		case '}':
			if len(braceStack) > 0 {
				braceStack = braceStack[:len(braceStack)-1]
			}
			out = append(out, src[idx])
			prevTok = lastTok
			lastTok = tokOther
			idx++
			continue
		case ',':
			if len(braceStack) > 0 && braceStack[len(braceStack)-1] {
				out = append(out, src[idx])
				idx++
				for idx < len(src) && (src[idx] == ' ' || src[idx] == '\t') {
					idx++
				}
				prevTok = lastTok
				lastTok = tokOther
				continue
			}
		}

		out = append(out, src[idx])
		prevTok = lastTok
		lastTok = tokOther
		idx++
	}
	return out
}
