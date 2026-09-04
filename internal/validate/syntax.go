package validate

import "github.com/Tom-R-Main/Internationalizer/internal/message"

// SyntaxSourceFindings checks a source with its already resolved runtime grammar.
func SyntaxSourceFindings(key, source, sourceLocale string, syntax message.Syntax) []Finding {
	if syntax == message.ICU {
		return messageFindings(key, message.CompareICU(source, source, sourceLocale))
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
