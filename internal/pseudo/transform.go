// Package pseudo generates deterministic accented and bidirectional test text.
package pseudo

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tom-R-Main/Internationalizer/internal/fluentpattern"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
)

// Strategy selects a pseudolocalization transformation.
type Strategy string

const (
	Accented Strategy = "accented"
	Bidi     Strategy = "bidi"
)

var protectedTokenRe = regexp.MustCompile("(?s)<!--.*?-->|```.*?```|~~~.*?~~~|`+[^`\\n]*`+|</?[A-Za-z][^>]*>|\\{\\{[A-Za-z0-9_.-]+\\}\\}|%\\{[A-Za-z0-9_.-]+\\}|\\{[A-Za-z0-9_.-]+\\}")

type span struct {
	start int
	end   int
}

// Transform transforms linguistic text while preserving runtime syntax.
func Transform(input string, strategy Strategy) (string, error) {
	if strategy != Accented && strategy != Bidi {
		return "", fmt.Errorf("unsupported pseudolocalization strategy %q", strategy)
	}
	transformLiteral := func(value string) string {
		return transformPreserving(value, strategy)
	}
	output := ""
	var err error
	if fluentpattern.LooksLike(input) {
		output, err = fluentpattern.TransformText(input, transformLiteral)
		if err != nil {
			return "", fmt.Errorf("pseudolocalizing Fluent pattern: %w", err)
		}
	} else if message.LooksLike(input) {
		output, err = message.TransformText(input, transformLiteral)
		if err != nil {
			return "", fmt.Errorf("pseudolocalizing ICU message: %w", err)
		}
	} else {
		output = transformLiteral(input)
	}
	if strategy == Accented {
		return "[!! " + output + " !!]", nil
	}
	return "\u2067" + output + "\u2069", nil
}

// DefaultLocale returns the conventional Unicode pseudo-locale for strategy.
func DefaultLocale(strategy Strategy) (string, error) {
	switch strategy {
	case Accented:
		return "en-XA", nil
	case Bidi:
		return "ar-XB", nil
	default:
		return "", fmt.Errorf("unsupported pseudolocalization strategy %q", strategy)
	}
}

func transformPreserving(input string, strategy Strategy) string {
	spans := protectedSpans(input)
	var output strings.Builder
	position := 0
	for _, protected := range spans {
		if protected.start > position {
			output.WriteString(transformText(input[position:protected.start], strategy))
		}
		output.WriteString(input[protected.start:protected.end])
		position = protected.end
	}
	output.WriteString(transformText(input[position:], strategy))
	return output.String()
}

func protectedSpans(input string) []span {
	spans := make([]span, 0)
	for _, location := range protectedTokenRe.FindAllStringIndex(input, -1) {
		spans = append(spans, span{start: location[0], end: location[1]})
	}
	spans = append(spans, markdownDestinationSpans(input)...)
	sort.Slice(spans, func(i, j int) bool {
		if spans[i].start != spans[j].start {
			return spans[i].start < spans[j].start
		}
		return spans[i].end > spans[j].end
	})
	merged := spans[:0]
	for _, candidate := range spans {
		if len(merged) == 0 || candidate.start > merged[len(merged)-1].end {
			merged = append(merged, candidate)
			continue
		}
		if candidate.end > merged[len(merged)-1].end {
			merged[len(merged)-1].end = candidate.end
		}
	}
	return merged
}

func markdownDestinationSpans(input string) []span {
	var spans []span
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
			}
			if depth == 0 {
				break
			}
		}
		if depth != 0 {
			break
		}
		spans = append(spans, span{start: start, end: end})
		offset = end + 1
	}
	return spans
}

func transformText(input string, strategy Strategy) string {
	if strategy == Bidi {
		return reverseRunes(mapRunes(input, bidiRune))
	}
	mapped := mapRunes(input, accentedRune)
	var expanded strings.Builder
	for _, character := range mapped {
		expanded.WriteRune(character)
		if isVowel(character) {
			expanded.WriteRune(character)
		}
	}
	return expanded.String()
}

func mapRunes(input string, transform func(rune) rune) string {
	return strings.Map(transform, input)
}

func reverseRunes(input string) string {
	characters := []rune(input)
	for left, right := 0, len(characters)-1; left < right; left, right = left+1, right-1 {
		characters[left], characters[right] = characters[right], characters[left]
	}
	return string(characters)
}

func accentedRune(character rune) rune {
	if replacement, ok := accentedCharacters[character]; ok {
		return replacement
	}
	return character
}

func bidiRune(character rune) rune {
	if replacement, ok := bidiCharacters[character]; ok {
		return replacement
	}
	return character
}

func isVowel(character rune) bool {
	decomposed := unicode.ToLower(character)
	return strings.ContainsRune("aàáâãäåæeèéêëiìíîïoòóôõöøuùúûüyýÿȧḗī", decomposed) && utf8.RuneLen(character) > 0
}

var accentedCharacters = map[rune]rune{
	'a': 'ȧ', 'b': 'ƀ', 'c': 'ƈ', 'd': 'ḓ', 'e': 'ḗ', 'f': 'ƒ', 'g': 'ɠ', 'h': 'ħ', 'i': 'ī', 'j': 'ĵ', 'k': 'ķ', 'l': 'ŀ', 'm': 'ḿ',
	'n': 'ƞ', 'o': 'ǿ', 'p': 'ƥ', 'q': 'ɋ', 'r': 'ř', 's': 'ş', 't': 'ŧ', 'u': 'ŭ', 'v': 'ṽ', 'w': 'ẇ', 'x': 'ẋ', 'y': 'ẏ', 'z': 'ž',
	'A': 'Ȧ', 'B': 'Ɓ', 'C': 'Ƈ', 'D': 'Ḓ', 'E': 'Ḗ', 'F': 'Ƒ', 'G': 'Ɠ', 'H': 'Ħ', 'I': 'Ī', 'J': 'Ĵ', 'K': 'Ķ', 'L': 'Ŀ', 'M': 'Ḿ',
	'N': 'Ƞ', 'O': 'Ø', 'P': 'Ƥ', 'Q': 'Ɋ', 'R': 'Ř', 'S': 'Ş', 'T': 'Ŧ', 'U': 'Ŭ', 'V': 'Ṽ', 'W': 'Ẇ', 'X': 'Ẋ', 'Y': 'Ẏ', 'Z': 'Ž',
}

var bidiCharacters = map[rune]rune{
	'a': 'ɐ', 'b': 'q', 'c': 'ɔ', 'd': 'p', 'e': 'ǝ', 'f': 'ɟ', 'g': 'ƃ', 'h': 'ɥ', 'i': 'ı', 'j': 'ɾ', 'k': 'ʞ', 'l': 'ן', 'm': 'ɯ',
	'n': 'u', 'o': 'o', 'p': 'd', 'q': 'b', 'r': 'ɹ', 's': 's', 't': 'ʇ', 'u': 'n', 'v': 'ʌ', 'w': 'ʍ', 'x': 'x', 'y': 'ʎ', 'z': 'z',
	'A': '∀', 'B': 'ꓭ', 'C': 'Ɔ', 'D': '◖', 'E': 'Ǝ', 'F': 'Ⅎ', 'G': '⅁', 'H': 'H', 'I': 'I', 'J': 'ſ', 'K': 'ꓘ', 'L': '˥', 'M': 'W',
	'N': 'N', 'O': 'O', 'P': 'Ԁ', 'Q': 'Ό', 'R': 'ꓤ', 'S': 'S', 'T': '⊥', 'U': '∩', 'V': 'Λ', 'W': 'M', 'X': 'X', 'Y': '⅄', 'Z': 'Z',
}
