package validate

import (
	"fmt"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"github.com/Tom-R-Main/Internationalizer/internal/protectedtext"
)

// SyntaxSourceFindings checks a source with its already resolved runtime grammar.
func SyntaxSourceFindings(key, source, sourceLocale string, syntax message.Syntax, requested ...message.Syntax) []Finding {
	if syntax == message.ICU {
		findings := messageFindings(key, message.CompareICU(source, source, sourceLocale))
		if len(requested) > 0 && (requested[0] == message.Auto || requested[0] == "") {
			for i := range findings {
				if findings[i].Code != CodeICUMessageSyntax {
					continue
				}
				context := "contains ambiguous brace syntax"
				for _, code := range protectedtext.HTMLCode(source) {
					if strings.Contains(code, "{") {
						context = "contains brace syntax inside HTML code"
						break
					}
				}
				findings[i].Message = fmt.Sprintf("%s %s. With message_syntax: auto, it was interpreted as ICU. Select the bundle's runtime syntax: plain, i18next, or icu. %s", key, context, findings[i].Message)
			}
		}
		return findings
	}
	return nil
}

func SyntaxFindings(key, source, target, targetLocale string, syntax message.Syntax) []Finding {
	if syntax == message.ICU {
		return messageFindings(key, message.CompareICU(source, target, targetLocale))
	}
	return nil
}

func SyntaxInterpolationMismatch(key, source, target string, syntax message.Syntax) *Mismatch {
	switch syntax {
	case message.Plain, message.ICU, message.Fluent:
		return nil
	case message.I18next:
		expected, actual := message.I18nextTokens(source), message.I18nextTokens(target)
		if !sameVars(expected, actual) {
			return &Mismatch{Key: key, SourceVars: expected, TargetVars: actual}
		}
		return nil
	default:
		return InterpolationMismatch(key, source, target)
	}
}
