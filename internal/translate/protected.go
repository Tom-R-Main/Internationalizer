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
	if strings.TrimSpace(source) != "" && strings.TrimSpace(target) == "" {
		return fmt.Errorf("blank translation for %q", key)
	}
	if findings := validation.ProtectedFindings(key, source, target); len(findings) > 0 {
		return fmt.Errorf("%s for %q", findings[0].Message, key)
	}
	if findings := validation.ICUFindings(key, source, target, targetLocale); len(findings) > 0 {
		return fmt.Errorf("%s for %q", findings[0].Message, key)
	}
	return nil
}
