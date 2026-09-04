package validate

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

const localizationNameAttribute = "data-l10n-name"

type markupToken struct {
	raw         string
	name        string
	attributes  []markupAttribute
	closing     bool
	selfClosing bool
	comment     bool
	valid       bool
}

type markupAttribute struct {
	name     string
	value    string
	hasValue bool
}

type markupFrame struct {
	tag   string
	slot  string
	owner string
}

// extractHTMLStructure preserves ordinary HTML structure in order while
// comparing data-l10n-name elements by semantic identity. This permits a
// translator to move named rich-text slots without changing their element,
// protected attributes, nesting, or contained markup.
func extractHTMLStructure(input string) []string {
	tokens := scanMarkup(input)
	root := make([]string, 0, len(tokens))
	named := make([]string, 0, len(tokens))
	bodies := make(map[string][]string)
	stack := make([]markupFrame, 0)

	appendOwned := func(owner, value string) {
		if owner == "" {
			root = append(root, value)
			return
		}
		bodies[owner] = append(bodies[owner], value)
	}
	currentOwner := func() string {
		for index := len(stack) - 1; index >= 0; index-- {
			if stack[index].slot != "" {
				return stack[index].slot
			}
		}
		return ""
	}

	for _, token := range tokens {
		owner := currentOwner()
		if token.comment {
			appendOwned(owner, "comment:"+token.raw)
			continue
		}
		if !token.valid {
			appendOwned(owner, "invalid:"+token.raw)
			continue
		}
		if token.closing {
			if len(stack) == 0 || stack[len(stack)-1].tag != token.name {
				appendOwned(owner, "unmatched-close:"+token.canonical())
				continue
			}
			frame := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if frame.slot != "" {
				named = append(named, fmt.Sprintf("slot:%s:parent:%s:close:%s", frame.slot, frame.owner, token.name))
			} else {
				appendOwned(owner, "close:"+token.name)
			}
			continue
		}

		slot := token.attribute(localizationNameAttribute)
		if slot != "" {
			named = append(named, fmt.Sprintf("slot:%s:parent:%s:open:%s", slot, owner, token.canonical()))
			if owner != "" {
				appendOwned(owner, "child-slot:"+slot)
			}
			if !token.selfClosing && !isVoidElement(token.name) {
				stack = append(stack, markupFrame{tag: token.name, slot: slot, owner: owner})
			}
			continue
		}

		appendOwned(owner, "open:"+token.canonical())
		if !token.selfClosing && !isVoidElement(token.name) {
			stack = append(stack, markupFrame{tag: token.name, owner: owner})
		}
	}

	for len(stack) > 0 {
		frame := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		owner := frame.owner
		if frame.slot != "" {
			owner = frame.slot
		}
		appendOwned(owner, "unclosed:"+frame.tag)
	}

	slotNames := make([]string, 0, len(bodies))
	for slot := range bodies {
		slotNames = append(slotNames, slot)
	}
	sort.Strings(slotNames)
	structure := append([]string(nil), root...)
	for _, slot := range slotNames {
		for index, value := range bodies[slot] {
			structure = append(structure, fmt.Sprintf("slot-body:%s:%08d:%s", slot, index, value))
		}
	}
	sort.Strings(named)
	return append(structure, named...)
}

func scanMarkup(input string) []markupToken {
	var tokens []markupToken
	for offset := 0; offset < len(input); {
		relative := strings.IndexByte(input[offset:], '<')
		if relative < 0 {
			break
		}
		start := offset + relative
		if strings.HasPrefix(input[start:], "<!--") {
			end := strings.Index(input[start+4:], "-->")
			if end < 0 {
				break
			}
			end += start + 7
			tokens = append(tokens, markupToken{raw: input[start:end], comment: true, valid: true})
			offset = end
			continue
		}
		cursor := start + 1
		if cursor < len(input) && input[cursor] == '/' {
			cursor++
		}
		if cursor >= len(input) || !isMarkupNameStart(rune(input[cursor])) {
			offset = start + 1
			continue
		}
		end := markupTagEnd(input, cursor)
		if end < 0 {
			break
		}
		raw := input[start : end+1]
		token := parseMarkupToken(raw)
		token.raw = raw
		tokens = append(tokens, token)
		offset = end + 1
	}
	return tokens
}

func markupTagEnd(input string, start int) int {
	var quote byte
	for index := start; index < len(input); index++ {
		character := input[index]
		if quote != 0 {
			if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' {
			quote = character
			continue
		}
		if character == '>' {
			return index
		}
	}
	return -1
}

func parseMarkupToken(raw string) markupToken {
	inside := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(raw, "<"), ">"))
	token := markupToken{valid: true}
	if strings.HasPrefix(inside, "/") {
		token.closing = true
		inside = strings.TrimSpace(inside[1:])
	}
	if strings.HasSuffix(inside, "/") {
		token.selfClosing = true
		inside = strings.TrimSpace(inside[:len(inside)-1])
	}
	nameEnd := 0
	for nameEnd < len(inside) && isMarkupNameRune(rune(inside[nameEnd])) {
		nameEnd++
	}
	if nameEnd == 0 {
		token.valid = false
		return token
	}
	token.name = strings.ToLower(inside[:nameEnd])
	remainder := inside[nameEnd:]
	if token.closing {
		token.valid = strings.TrimSpace(remainder) == "" && !token.selfClosing
		return token
	}

	seen := make(map[string]struct{})
	for position := 0; ; {
		position = skipMarkupSpace(remainder, position)
		if position >= len(remainder) {
			break
		}
		start := position
		for position < len(remainder) && isMarkupAttributeRune(rune(remainder[position])) {
			position++
		}
		if position == start {
			token.valid = false
			return token
		}
		attribute := markupAttribute{name: strings.ToLower(remainder[start:position])}
		position = skipMarkupSpace(remainder, position)
		if position < len(remainder) && remainder[position] == '=' {
			attribute.hasValue = true
			position++
			position = skipMarkupSpace(remainder, position)
			if position >= len(remainder) {
				token.valid = false
				return token
			}
			if remainder[position] == '\'' || remainder[position] == '"' {
				quote := remainder[position]
				position++
				valueStart := position
				for position < len(remainder) && remainder[position] != quote {
					position++
				}
				if position >= len(remainder) {
					token.valid = false
					return token
				}
				attribute.value = remainder[valueStart:position]
				position++
			} else {
				valueStart := position
				for position < len(remainder) && !unicode.IsSpace(rune(remainder[position])) {
					position++
				}
				attribute.value = remainder[valueStart:position]
			}
		}
		if _, duplicate := seen[attribute.name]; duplicate {
			token.valid = false
			return token
		}
		seen[attribute.name] = struct{}{}
		token.attributes = append(token.attributes, attribute)
	}
	return token
}

func (token markupToken) attribute(name string) string {
	for _, attribute := range token.attributes {
		if attribute.name == name && attribute.hasValue {
			return attribute.value
		}
	}
	return ""
}

func (token markupToken) canonical() string {
	if token.closing {
		return "/" + token.name
	}
	attributes := append([]markupAttribute(nil), token.attributes...)
	sort.Slice(attributes, func(left, right int) bool { return attributes[left].name < attributes[right].name })
	var result strings.Builder
	result.WriteString(token.name)
	for _, attribute := range attributes {
		result.WriteByte('|')
		result.WriteString(attribute.name)
		if attribute.hasValue {
			fmt.Fprintf(&result, "=%d:%s", len(attribute.value), attribute.value)
		}
	}
	if token.selfClosing {
		result.WriteString("|/")
	}
	return result.String()
}

func skipMarkupSpace(input string, position int) int {
	for position < len(input) && unicode.IsSpace(rune(input[position])) {
		position++
	}
	return position
}

func isMarkupNameStart(character rune) bool {
	return unicode.IsLetter(character)
}

func isMarkupNameRune(character rune) bool {
	return unicode.IsLetter(character) || unicode.IsDigit(character) || character == ':' || character == '-'
}

func isMarkupAttributeRune(character rune) bool {
	return isMarkupNameRune(character) || character == '_' || character == '.'
}

func isVoidElement(name string) bool {
	switch name {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr":
		return true
	default:
		return false
	}
}
