package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/message"
)

var (
	htmlTagRe    = regexp.MustCompile(`(?s)<!--.*?-->|</?[A-Za-z][^>]*>`)
	inlineCodeRe = regexp.MustCompile("`+[^`\\n]*`+")
)

// ProtectedFindings compares source and target structures that translations
// must preserve exactly. Multiple damaged structures may yield multiple
// findings with the same stable code and distinct messages.
func ProtectedFindings(key, source, target, targetLocale string) []Finding {
	var findings []Finding
	icu := usesICUValidation(source, target)
	if !icu {
		if mismatch := InterpolationMismatch(key, source, target); mismatch != nil {
			findings = append(findings, protectedFinding(key, "interpolation variables", mismatch.SourceVars, mismatch.TargetVars))
		}
	}
	checks := []struct {
		name    string
		extract func(string) []string
	}{
		{"HTML structure", extractHTMLTags},
		{"fenced code", extractFencedCode},
		{"inline code", extractInlineCode},
		{"markdown link destinations", extractLinkDestinations},
	}
	tokenizers := make([]func(string) []string, len(checks))
	for index, check := range checks {
		tokenizers[index] = check.extract
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
		sourceTokens := check.extract(source)
		targetTokens := check.extract(target)
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

func extractHTMLTags(input string) []string {
	return htmlTagRe.FindAllString(input, -1)
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
