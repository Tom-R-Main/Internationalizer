// Package fluentpattern validates and compares runtime-significant Fluent
// pattern structure without owning resource-file serialization.
package fluentpattern

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Analysis is the protected structure of one Fluent pattern.
type Analysis struct {
	Expressions []string
	Selectors   []Selector
}

// Selector is one Fluent select expression and its locale variants.
type Selector struct {
	Head       string
	DefaultKey string
	Variants   map[string]Analysis
}

// LooksLike reports whether text contains syntax specific to Fluent rather
// than a plain ICU-style argument.
func LooksLike(value string) bool {
	for offset := 0; offset < len(value); {
		start := strings.IndexByte(value[offset:], '{')
		if start < 0 {
			return false
		}
		start += offset
		end, err := matchingBrace(value, start)
		if err != nil {
			return true
		}
		expression := strings.TrimSpace(value[start+1 : end])
		if strings.Contains(expression, "->") || strings.ContainsAny(expression, "$.") || strings.HasPrefix(expression, "-") || strings.HasPrefix(expression, "\"") || strings.HasPrefix(expression, "'") {
			return true
		}
		offset = end + 1
	}
	return false
}

// Analyze validates braces and select-expression structure and returns a
// canonical representation of all non-linguistic expressions.
func Analyze(value string) (Analysis, error) {
	analysis := Analysis{}
	for offset := 0; offset < len(value); {
		relative := strings.IndexAny(value[offset:], "{}")
		if relative < 0 {
			break
		}
		start := offset + relative
		if value[start] == '}' {
			return Analysis{}, fmt.Errorf("unexpected closing brace at byte %d", start)
		}
		end, err := matchingBrace(value, start)
		if err != nil {
			return Analysis{}, err
		}
		expression := strings.TrimSpace(value[start+1 : end])
		if expression == "" {
			return Analysis{}, fmt.Errorf("empty Fluent placeable at byte %d", start)
		}
		arrow := topLevelArrow(expression)
		if arrow < 0 {
			analysis.Expressions = append(analysis.Expressions, canonicalExpression(expression))
		} else {
			selector, err := parseSelector(expression[:arrow], expression[arrow+2:])
			if err != nil {
				return Analysis{}, fmt.Errorf("fluent selector at byte %d: %w", start, err)
			}
			analysis.Selectors = append(analysis.Selectors, selector)
		}
		offset = end + 1
	}
	sort.Strings(analysis.Expressions)
	sort.Slice(analysis.Selectors, func(left, right int) bool {
		return selectorSignature(analysis.Selectors[left]) < selectorSignature(analysis.Selectors[right])
	})
	return analysis, nil
}

// Compare checks that target preserves source expressions and selectors.
// Target-only locale variants are allowed when they preserve the source
// selector's default-branch structure.
func Compare(source, target string) (expected, actual []string, preserved bool, err error) {
	sourceAnalysis, err := Analyze(source)
	if err != nil {
		return nil, nil, false, fmt.Errorf("source pattern: %w", err)
	}
	targetAnalysis, err := Analyze(target)
	if err != nil {
		return sourceAnalysis.Signatures(), nil, false, fmt.Errorf("target pattern: %w", err)
	}
	return sourceAnalysis.Signatures(), targetAnalysis.Signatures(), compareAnalysis(sourceAnalysis, targetAnalysis), nil
}

// TransformText applies transform to linguistic text while retaining Fluent
// placeables, selector heads, variant keys, and default markers byte-for-byte.
func TransformText(value string, transform func(string) string) (string, error) {
	var result strings.Builder
	offset := 0
	for offset < len(value) {
		relative := strings.IndexAny(value[offset:], "{}")
		if relative < 0 {
			result.WriteString(transform(value[offset:]))
			break
		}
		start := offset + relative
		if value[start] == '}' {
			return "", fmt.Errorf("unexpected closing brace at byte %d", start)
		}
		result.WriteString(transform(value[offset:start]))
		end, err := matchingBrace(value, start)
		if err != nil {
			return "", err
		}
		expression := value[start+1 : end]
		arrow := topLevelArrow(expression)
		if arrow < 0 {
			result.WriteString(value[start : end+1])
		} else {
			body := expression[arrow+2:]
			headers, err := variantHeaders(body)
			if err != nil || len(headers) == 0 {
				if err == nil {
					err = fmt.Errorf("selector has no variants")
				}
				return "", err
			}
			result.WriteByte('{')
			result.WriteString(expression[:arrow+2])
			result.WriteString(body[:headers[0].end])
			for index, header := range headers {
				branchEnd := len(body)
				if index+1 < len(headers) {
					branchEnd = headers[index+1].start
				}
				transformed, err := TransformText(body[header.end:branchEnd], transform)
				if err != nil {
					return "", err
				}
				result.WriteString(transformed)
				if index+1 < len(headers) {
					result.WriteString(body[headers[index+1].start:headers[index+1].end])
				}
			}
			result.WriteByte('}')
		}
		offset = end + 1
	}
	return result.String(), nil
}

// Signatures returns a deterministic diagnostic view of an analysis.
func (analysis Analysis) Signatures() []string {
	result := make([]string, 0, len(analysis.Expressions)+len(analysis.Selectors))
	for _, expression := range analysis.Expressions {
		result = append(result, "expression:"+expression)
	}
	for _, selector := range analysis.Selectors {
		result = append(result, selectorSignature(selector))
	}
	sort.Strings(result)
	return result
}

func compareAnalysis(source, target Analysis) bool {
	if !equalList(source.Expressions, target.Expressions) || len(source.Selectors) != len(target.Selectors) {
		return false
	}
	used := make([]bool, len(target.Selectors))
	for _, sourceSelector := range source.Selectors {
		matched := false
		for index, targetSelector := range target.Selectors {
			if used[index] || sourceSelector.Head != targetSelector.Head || !compareSelector(sourceSelector, targetSelector) {
				continue
			}
			used[index] = true
			matched = true
			break
		}
		if !matched {
			return false
		}
	}
	return true
}

func compareSelector(source, target Selector) bool {
	if source.DefaultKey != target.DefaultKey {
		return false
	}
	fallback, ok := source.Variants[source.DefaultKey]
	if !ok {
		return false
	}
	for key, sourceBranch := range source.Variants {
		targetBranch, exists := target.Variants[key]
		if !exists || !compareAnalysis(sourceBranch, targetBranch) {
			return false
		}
	}
	for key, targetBranch := range target.Variants {
		if _, exists := source.Variants[key]; exists {
			continue
		}
		if !compareAnalysis(fallback, targetBranch) {
			return false
		}
	}
	return true
}

func parseSelector(head, body string) (Selector, error) {
	head = canonicalExpression(strings.TrimSpace(head))
	if head == "" {
		return Selector{}, fmt.Errorf("selector expression is empty")
	}
	headers, err := variantHeaders(body)
	if err != nil {
		return Selector{}, err
	}
	if len(headers) == 0 {
		return Selector{}, fmt.Errorf("selector has no variants")
	}
	selector := Selector{Head: head, Variants: make(map[string]Analysis)}
	for index, header := range headers {
		end := len(body)
		if index+1 < len(headers) {
			end = headers[index+1].start
		}
		branch, err := Analyze(body[header.end:end])
		if err != nil {
			return Selector{}, fmt.Errorf("variant %q: %w", header.key, err)
		}
		if _, duplicate := selector.Variants[header.key]; duplicate {
			return Selector{}, fmt.Errorf("duplicate variant %q", header.key)
		}
		selector.Variants[header.key] = branch
		if header.defaultVariant {
			if selector.DefaultKey != "" {
				return Selector{}, fmt.Errorf("selector has multiple default variants")
			}
			selector.DefaultKey = header.key
		}
	}
	if selector.DefaultKey == "" {
		return Selector{}, fmt.Errorf("selector has no default variant")
	}
	return selector, nil
}

type variantHeader struct {
	start          int
	end            int
	key            string
	defaultVariant bool
}

func variantHeaders(body string) ([]variantHeader, error) {
	var headers []variantHeader
	depth := 0
	var quote byte
	escaped := false
	lineStart := true
	for index := 0; index < len(body); index++ {
		character := body[index]
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			switch character {
			case '\\':
				escaped = true
			case quote:
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote = character
			continue
		}
		switch character {
		case '{':
			depth++
		case '}':
			if depth == 0 {
				return nil, fmt.Errorf("unexpected closing brace in selector")
			}
			depth--
		case '\n', '\r':
			lineStart = true
			continue
		}
		if depth != 0 {
			continue
		}
		if lineStart && (character == ' ' || character == '\t') {
			continue
		}
		if !lineStart {
			continue
		}
		defaultVariant := false
		start := index
		if character == '*' {
			defaultVariant = true
			index++
			if index >= len(body) {
				return nil, fmt.Errorf("incomplete default variant")
			}
			character = body[index]
		}
		if character != '[' {
			lineStart = false
			continue
		}
		close := strings.IndexByte(body[index+1:], ']')
		if close < 0 {
			return nil, fmt.Errorf("unterminated variant key")
		}
		close += index + 1
		key := strings.TrimSpace(body[index+1 : close])
		if key == "" || strings.ContainsAny(key, "[]{}") {
			return nil, fmt.Errorf("invalid variant key %q", key)
		}
		headers = append(headers, variantHeader{start: start, end: close + 1, key: key, defaultVariant: defaultVariant})
		index = close
		lineStart = false
	}
	if depth != 0 || quote != 0 {
		return nil, fmt.Errorf("unterminated selector expression")
	}
	if len(headers) > 0 && strings.TrimSpace(body[:headers[0].start]) != "" {
		return nil, fmt.Errorf("content appears before the first variant")
	}
	return headers, nil
}

func matchingBrace(value string, start int) (int, error) {
	depth := 0
	var quote byte
	escaped := false
	for index := start; index < len(value); index++ {
		character := value[index]
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			switch character {
			case '\\':
				escaped = true
			case quote:
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote = character
			continue
		}
		switch character {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return index, nil
			}
		}
	}
	return 0, fmt.Errorf("unterminated Fluent placeable at byte %d", start)
}

func topLevelArrow(expression string) int {
	depth := 0
	var quote byte
	escaped := false
	for index := 0; index+1 < len(expression); index++ {
		character := expression[index]
		if quote != 0 {
			if escaped {
				escaped = false
			} else if character == '\\' {
				escaped = true
			} else if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote = character
			continue
		}
		switch character {
		case '{', '(':
			depth++
		case '}', ')':
			if depth > 0 {
				depth--
			}
		case '-':
			if depth == 0 && expression[index+1] == '>' {
				return index
			}
		}
	}
	return -1
}

func canonicalExpression(expression string) string {
	var result strings.Builder
	var quote rune
	escaped := false
	for _, character := range strings.TrimSpace(expression) {
		if quote != 0 {
			result.WriteRune(character)
			if escaped {
				escaped = false
			} else if character == '\\' {
				escaped = true
			} else if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote = character
			result.WriteRune(character)
			continue
		}
		if !unicode.IsSpace(character) {
			result.WriteRune(character)
		}
	}
	return result.String()
}

func selectorSignature(selector Selector) string {
	keys := make([]string, 0, len(selector.Variants))
	for key := range selector.Variants {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var result strings.Builder
	result.WriteString("selector:")
	result.WriteString(selector.Head)
	result.WriteString(":default=")
	result.WriteString(selector.DefaultKey)
	for _, key := range keys {
		result.WriteString(":variant=")
		result.WriteString(key)
		result.WriteByte('[')
		result.WriteString(strings.Join(selector.Variants[key].Signatures(), ","))
		result.WriteByte(']')
	}
	return result.String()
}

func equalList(left, right []string) bool {
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
