package translate

import (
	"fmt"
	"strings"

	validation "github.com/Tom-R-Main/Internationalizer/internal/validate"
)

// validateTranslationValue enforces deterministic invariants at both provider
// and explicit-adoption boundaries. Instructions in a prompt are not proof
// that protected source structure survived translation.
func validateTranslationValue(key, source, target, targetLocale string) error {
	return validateTranslationValueWithContext(key, source, target, targetLocale, valueValidationContext{})
}

type valueValidationContext struct {
	document   bool
	sourcePath string
	targetPath string
}

func validateTranslationValueWithContext(key, source, target, targetLocale string, validationContext valueValidationContext) error {
	if strings.TrimSpace(source) != "" && strings.TrimSpace(target) == "" {
		return fmt.Errorf("blank translation for %q", key)
	}
	protected := validation.ProtectedFindings(key, source, target, targetLocale)
	if validationContext.document {
		protected = validation.ProtectedDocumentFindings(key, source, target, targetLocale, validationContext.sourcePath, validationContext.targetPath)
	}
	if len(protected) > 0 {
		return fmt.Errorf("%s for %q", protected[0].Message, key)
	}
	if validationContext.document {
		return nil
	}
	if findings := validation.ICUFindings(key, source, target, targetLocale); len(findings) > 0 {
		return fmt.Errorf("%s for %q", findings[0].Message, key)
	}
	return nil
}
