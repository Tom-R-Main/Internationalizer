package translate

import (
	"fmt"
	"regexp"
	"strings"

	validation "github.com/Tom-R-Main/Internationalizer/internal/validate"
)

var (
	htmlTagRe    = regexp.MustCompile(`(?s)<!--.*?-->|</?[A-Za-z][^>]*>`)
	inlineCodeRe = regexp.MustCompile("`+[^`\\n]*`+")
)

// validateTranslationValue enforces deterministic invariants at both provider
// and explicit-adoption boundaries. Instructions in a prompt are not proof
// that protected source structure survived translation.
func validateTranslationValue(key, source, target string) error {
	if strings.TrimSpace(source) != "" && strings.TrimSpace(target) == "" {
		return fmt.Errorf("blank translation for %q", key)
	}
	if mismatch := validationMismatch(key, source, target); mismatch != nil {
		return mismatch
	}
	if !equalStrings(htmlTagRe.FindAllString(source, -1), htmlTagRe.FindAllString(target, -1)) {
		return fmt.Errorf("protected HTML structure mismatch for %q", key)
	}
	if !equalStrings(extractFencedCode(source), extractFencedCode(target)) {
		return fmt.Errorf("fenced code mismatch for %q", key)
	}
	if !equalStrings(inlineCodeRe.FindAllString(source, -1), inlineCodeRe.FindAllString(target, -1)) {
		return fmt.Errorf("inline code mismatch for %q", key)
	}
	if !equalStrings(extractLinkDestinations(source), extractLinkDestinations(target)) {
		return fmt.Errorf("markdown link destination mismatch for %q", key)
	}
	return nil
}

func validationMismatch(key, source, target string) error {
	mismatch := validation.InterpolationMismatch(key, source, target)
	if mismatch == nil {
		return nil
	}
	return fmt.Errorf("interpolation mismatch for %q (source: %v, target: %v)", key, mismatch.SourceVars, mismatch.TargetVars)
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
