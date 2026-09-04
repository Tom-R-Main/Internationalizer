package translate

import (
	"fmt"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/message"
	validation "github.com/Tom-R-Main/Internationalizer/internal/validate"
)

// validateTranslationValue enforces deterministic invariants at both provider
// and explicit-adoption boundaries. Instructions in a prompt are not proof
// that protected source structure survived translation.
func validateTranslationValue(key, source, target, targetLocale string) error {
	return validateTranslationValueWithContext(key, source, target, targetLocale, valueValidationContext{})
}

type valueValidationContext struct {
	syntaxes   map[string]message.Syntax
	document   bool
	sourcePath string
	targetPath string
}

func validateTranslationValueWithContext(key, source, target, targetLocale string, validationContext valueValidationContext) error {
	if strings.TrimSpace(source) != "" && strings.TrimSpace(target) == "" {
		return fmt.Errorf("blank translation for %q", key)
	}
	syntax := validationContext.syntaxes[key]
	if syntax == "" {
		format := ""
		if validationContext.document {
			format = "markdown"
		}
		syntax = message.ResolveSyntax(format, message.Auto, source)
	}
	protected := validation.ProtectedSyntaxFindings(key, source, target, targetLocale, syntax)
	if validationContext.document {
		protected = validation.ProtectedDocumentFindings(key, source, target, targetLocale, validationContext.sourcePath, validationContext.targetPath, syntax)
	}
	if len(protected) > 0 {
		return fmt.Errorf("%s for %q", protected[0].Message, key)
	}
	if validationContext.document {
		return nil
	}
	if findings := validation.SyntaxFindings(key, source, target, targetLocale, syntax); len(findings) > 0 {
		return fmt.Errorf("%s for %q", findings[0].Message, key)
	}
	return nil
}
