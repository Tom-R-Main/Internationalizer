package validate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/fluentpattern"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"github.com/Tom-R-Main/Internationalizer/internal/protectedtext"
)

var (
	htmlTagRe      = regexp.MustCompile(`(?s)<!--.*?-->|</?[A-Za-z][^>]*>`)
	htmlPathAttrRe = regexp.MustCompile(`(?i)\b(href|src)\s*=\s*("[^"]*"|'[^']*')`)
	inlineCodeRe   = regexp.MustCompile("`+[^`\\n]*`+")
	uriSchemeRe    = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)
)

// ProtectedFindings compares source and target structures that translations
// must preserve exactly. Multiple damaged structures may yield multiple
// findings with the same stable code and distinct messages.
func ProtectedFindings(key, source, target, targetLocale string) []Finding {
	return ProtectedSyntaxFindings(key, source, target, targetLocale, message.ResolveSyntax("", message.Auto, source))
}

func ProtectedSyntaxFindings(key, source, target, targetLocale string, syntax message.Syntax) []Finding {
	return protectedFindings(key, source, target, targetLocale, "", "", syntax)
}

// ProtectedDocumentFindings compares Markdown document structure while
// resolving relative links from each document's own directory. A localized
// document may add one link back to its source.
func ProtectedDocumentFindings(key, source, target, targetLocale, sourcePath, targetPath string, syntaxes ...message.Syntax) []Finding {
	syntax := message.Legacy
	if len(syntaxes) > 0 {
		syntax = syntaxes[0]
	}
	return protectedFindings(key, source, target, targetLocale, sourcePath, targetPath, syntax)
}

func protectedFindings(key, source, target, targetLocale, sourcePath, targetPath string, syntax message.Syntax) []Finding {
	var findings []Finding
	document := sourcePath != "" && targetPath != ""
	icu := false
	if !document {
		icu = syntax == message.ICU && usesICUValidation(source, target)
	}
	if !icu {
		if mismatch := SyntaxInterpolationMismatch(key, source, target, syntax); mismatch != nil {
			findings = append(findings, protectedFinding(key, "interpolation variables", mismatch.SourceVars, mismatch.TargetVars))
		}
		if syntax == message.Fluent {
			expected, actual, preserved, err := fluentpattern.Compare(source, target)
			if err != nil {
				actual = []string{err.Error()}
			}
			if err != nil || !preserved {
				findings = append(findings, protectedFinding(key, "Fluent pattern", expected, actual))
			}
		}
	}
	checks := []struct {
		name          string
		extractSource func(string) []string
		extractTarget func(string) []string
	}{
		{"HTML structure", extractHTMLTags, extractHTMLTags},
		{"HTML code", protectedtext.HTMLCode, protectedtext.HTMLCode},
		{"fenced code", extractFencedCode, extractFencedCode},
		{"inline code", extractInlineCode, extractInlineCode},
		{"markdown link destinations", extractLinkDestinations, extractLinkDestinations},
	}
	if document {
		checks[0].extractSource = func(input string) []string { return extractDocumentHTMLTags(input, sourcePath) }
		checks[0].extractTarget = func(input string) []string { return extractDocumentHTMLTags(input, targetPath) }
		checks[4].extractSource = func(input string) []string { return extractDocumentLinkDestinations(input, sourcePath) }
		checks[4].extractTarget = func(input string) []string {
			tokens := extractDocumentLinkDestinations(input, targetPath)
			return removeSourceBacklink(tokens, sourcePath)
		}
	}
	tokenizers := make([]func(string) []string, len(checks))
	for index, check := range checks {
		tokenizers[index] = check.extractSource
	}
	preservedKinds := make([]bool, len(checks))
	if icu {
		var err error
		preservedKinds, err = message.PreservedTextTokenKinds(source, target, targetLocale, tokenizers)
		if err != nil {
			icu = false
		}
	}
	for index, check := range checks {
		sourceTokens := check.extractSource(source)
		targetTokens := check.extractTarget(target)
		preserved := equalStrings(sourceTokens, targetTokens)
		if icu {
			preserved = preservedKinds[index]
		}
		if !preserved {
			findings = append(findings, protectedFinding(key, check.name, sourceTokens, targetTokens))
		}
	}
	return findings
}

func extractDocumentHTMLTags(input, documentPath string) []string {
	tags := htmlTagRe.FindAllString(input, -1)
	for index, tag := range tags {
		tags[index] = htmlPathAttrRe.ReplaceAllStringFunc(tag, func(attribute string) string {
			match := htmlPathAttrRe.FindStringSubmatch(attribute)
			if len(match) != 3 {
				return attribute
			}
			value := strings.Trim(match[2], `"'`)
			return strings.ToLower(match[1]) + `="` + resolveDocumentDestination(value, documentPath) + `"`
		})
	}
	return tags
}

func extractDocumentLinkDestinations(input, documentPath string) []string {
	destinations := extractLinkDestinations(input)
	for index, destination := range destinations {
		destinations[index] = resolveDocumentDestination(strings.Trim(destination, "<>"), documentPath)
	}
	return destinations
}

func removeSourceBacklink(destinations []string, sourcePath string) []string {
	want, err := filepath.Abs(sourcePath)
	if err != nil {
		want = filepath.Clean(sourcePath)
	}
	want = filepath.ToSlash(want)
	for index, destination := range destinations {
		if destination == want {
			return append(destinations[:index:index], destinations[index+1:]...)
		}
	}
	return destinations
}

func resolveDocumentDestination(destination, documentPath string) string {
	if destination == "" || strings.HasPrefix(destination, "#") || strings.HasPrefix(destination, "/") || uriSchemeRe.MatchString(destination) {
		return destination
	}
	pathPart := destination
	suffix := ""
	if index := strings.IndexAny(pathPart, "?#"); index >= 0 {
		pathPart, suffix = pathPart[:index], pathPart[index:]
	}
	if pathPart == "" {
		return destination
	}
	base, err := filepath.Abs(filepath.Dir(documentPath))
	if err != nil {
		base = filepath.Clean(filepath.Dir(documentPath))
	}
	return filepath.ToSlash(filepath.Clean(filepath.Join(base, pathPart))) + suffix
}

func extractHTMLTags(input string) []string {
	return extractHTMLStructure(input)
}

func extractInlineCode(input string) []string {
	return inlineCodeRe.FindAllString(input, -1)
}

func usesICUValidation(source, target string) bool {
	if !message.LooksLike(source) || !message.LooksLike(target) {
		return false
	}
	if _, err := message.Parse(source); err != nil {
		return false
	}
	_, err := message.Parse(target)
	return err == nil
}

func protectedFinding(key, structure string, expected, actual []string) Finding {
	return Finding{
		Code:     CodeProtectedStructureMismatch,
		Severity: SeverityError,
		Key:      key,
		Message:  fmt.Sprintf("protected %s mismatch", structure),
		Expected: expected,
		Actual:   actual,
	}
}

func extractFencedCode(input string) []string {
	lines := strings.SplitAfter(input, "\n")
	var blocks []string
	var block strings.Builder
	var marker byte
	var markerLength int
	inFence := false
	for _, line := range lines {
		trimmed := strings.TrimLeft(strings.TrimSuffix(line, "\n"), " \t")
		if !inFence {
			candidate, length, ok := fenceMarker(trimmed)
			if !ok {
				continue
			}
			inFence = true
			marker = candidate
			markerLength = length
			block.Reset()
			block.WriteString(line)
			continue
		}

		block.WriteString(line)
		if isClosingFence(trimmed, marker, markerLength) {
			blocks = append(blocks, block.String())
			inFence = false
		}
	}
	if inFence {
		blocks = append(blocks, block.String())
	}
	return blocks
}

func fenceMarker(line string) (byte, int, bool) {
	if len(line) < 3 || (line[0] != '`' && line[0] != '~') {
		return 0, 0, false
	}
	marker := line[0]
	length := 0
	for length < len(line) && line[length] == marker {
		length++
	}
	return marker, length, length >= 3
}

func isClosingFence(line string, marker byte, minimum int) bool {
	if len(line) < minimum || line[0] != marker {
		return false
	}
	length := 0
	for length < len(line) && line[length] == marker {
		length++
	}
	return length >= minimum && strings.TrimSpace(line[length:]) == ""
}

func extractLinkDestinations(input string) []string {
	var destinations []string
	for offset := 0; offset+1 < len(input); {
		relative := strings.Index(input[offset:], "](")
		if relative < 0 {
			break
		}
		start := offset + relative + 2
		depth := 1
		escaped := false
		end := start
		for ; end < len(input); end++ {
			character := input[end]
			if escaped {
				escaped = false
				continue
			}
			if character == '\\' {
				escaped = true
				continue
			}
			switch character {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					inside := strings.TrimSpace(input[start:end])
					destinations = append(destinations, firstLinkToken(inside))
				}
			}
			if depth == 0 {
				break
			}
		}
		if depth != 0 {
			break
		}
		offset = end + 1
	}
	return destinations
}

func firstLinkToken(inside string) string {
	if strings.HasPrefix(inside, "<") {
		if end := strings.IndexByte(inside, '>'); end >= 0 {
			return inside[:end+1]
		}
	}
	if end := strings.IndexAny(inside, " \t\n"); end >= 0 {
		return inside[:end]
	}
	return inside
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
