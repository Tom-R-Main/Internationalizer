package message

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/fluentpattern"
)

// Syntax describes the runtime grammar, independently of the resource format.
type Syntax string

const (
	Auto    Syntax = "auto"
	ICU     Syntax = "icu"
	I18next Syntax = "i18next"
	Plain   Syntax = "plain"
	Fluent  Syntax = "fluent"
	// Legacy preserves the non-ICU interpolation checks used by auto mode.
	Legacy Syntax = "legacy"
)

func ValidateSyntax(syntax Syntax) error {
	switch syntax {
	case "", Auto, ICU, I18next, Plain:
		return nil
	default:
		return fmt.Errorf("unsupported message_syntax %q (expected auto, icu, i18next, or plain)", syntax)
	}
}

// ResolveSyntax uses the source only: target text must never select its grammar.
// Fluent resources own their grammar; Markdown documents are not ICU messages.
func ResolveSyntax(format string, policy Syntax, source string) Syntax {
	if format == "fluent" {
		return Fluent
	}
	if policy != "" && policy != Auto {
		return policy
	}
	if format == "markdown" {
		return Legacy
	}
	for _, location := range typedICUArgument.FindAllStringIndex(source, -1) {
		// A formatter inside {{value, number}} belongs to i18next, not ICU.
		if location[0] == 0 || source[location[0]-1] != '{' {
			return ICU
		}
	}
	if fluentpattern.LooksLike(source) {
		return Fluent
	}
	if LooksLike(source) {
		return ICU
	}
	return Legacy
}

var typedICUArgument = regexp.MustCompile(`\{\s*[\pL\pN_.-]+\s*,`)

// I18nextTokens includes escaping and formatting modifiers as part of the
// contract, as well as nested variable paths. Repeated tokens remain repeated.
var I18nextToken = regexp.MustCompile(`\{\{[^{}\r\n]+\}\}`)

func I18nextTokens(value string) []string {
	tokens := I18nextToken.FindAllString(value, -1)
	for i, token := range tokens {
		tokens[i] = strings.TrimSpace(token[2 : len(token)-2])
	}
	sort.Strings(tokens)
	return tokens
}
