// Package protectedtext locates literal code shared by validation and pseudo.
package protectedtext

import (
	"regexp"
	"strings"
)

type Span struct{ Start, End int }

// Quoted attributes may contain >; comments must not introduce code elements.
var htmlToken = regexp.MustCompile(`(?s)<!--.*?-->|</?[A-Za-z][A-Za-z0-9:-]*(?:[^<>"']|"[^"]*"|'[^']*')*>`)

// HTMLCodeSpans includes the full code element, with nested markup intact.
// An unclosed code element protects through EOF rather than rewriting code.
func HTMLCodeSpans(input string) []Span {
	var spans []Span
	start, depth := 0, 0
	for _, loc := range htmlToken.FindAllStringIndex(input, -1) {
		token := input[loc[0]:loc[1]]
		if strings.HasPrefix(token, "<!--") {
			continue
		}
		closing := strings.HasPrefix(token, "</")
		name := strings.TrimLeft(token[1:], "/")
		end := strings.IndexAny(name, " \t\r\n/>")
		if end < 0 || !strings.EqualFold(name[:end], "code") {
			continue
		}
		if closing {
			if depth > 0 {
				depth--
				if depth == 0 {
					spans = append(spans, Span{start, loc[1]})
				}
			}
		} else if !strings.HasSuffix(token, "/>") {
			if depth == 0 {
				start = loc[0]
			}
			depth++
		}
	}
	if depth > 0 {
		spans = append(spans, Span{start, len(input)})
	}
	return spans
}

func HTMLCode(input string) []string {
	var tokens []string
	for _, span := range HTMLCodeSpans(input) {
		tokens = append(tokens, input[span.Start:span.End])
	}
	return tokens
}
