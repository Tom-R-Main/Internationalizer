package validate

import (
	"strings"

	localeid "github.com/Tom-R-Main/Internationalizer/internal/locale"
)

// PluralFormsFor returns the CLDR plural categories for a locale.
// It retains the historical ["one", "other"] fallback for an invalid tag.
func PluralFormsFor(locale string) []string {
	if forms, ok := KnownPluralFormsFor(locale); ok {
		return forms
	}
	return []string{"one", "other"}
}

// KnownPluralFormsFor returns CLDR-backed cardinal categories for a valid
// locale. Invalid or unrecognized tags are not guessed.
func KnownPluralFormsFor(locale string) ([]string, bool) {
	forms, err := localeid.CardinalCategories(locale)
	return forms, err == nil
}

// ExpandI18nextV4Source returns the source key set plus any target-only plural
// categories needed by targetLocale. A missing target category uses the source
// locale's "other" value as its translation template. optional contains source
// plural keys that the target locale does not require.
func ExpandI18nextV4Source(sourceKeys map[string]string, sourceLocale, targetLocale string) (expanded map[string]string, required, optional map[string]struct{}) {
	expanded = cloneStrings(sourceKeys)
	required = make(map[string]struct{})
	optional = make(map[string]struct{})

	sourceForms, sourceKnown := KnownPluralFormsFor(sourceLocale)
	targetForms, targetKnown := KnownPluralFormsFor(targetLocale)
	if !sourceKnown || !targetKnown {
		return expanded, required, optional
	}

	for base, values := range pluralFamilies(sourceKeys, sourceForms) {
		template := values["other"]
		if template == "" {
			for _, form := range sourceForms {
				if value, ok := values[form]; ok {
					template = value
					break
				}
			}
		}
		for _, form := range targetForms {
			key := base + "_" + form
			required[key] = struct{}{}
			if _, ok := expanded[key]; !ok {
				expanded[key] = template
			}
		}
		for _, form := range sourceForms {
			key := base + "_" + form
			if _, targetRequires := required[key]; !targetRequires {
				optional[key] = struct{}{}
			}
		}
	}
	return expanded, required, optional
}

func pluralFamilies(sourceKeys map[string]string, sourceForms []string) map[string]map[string]string {
	forms := make(map[string]struct{}, len(sourceForms))
	for _, form := range sourceForms {
		forms[form] = struct{}{}
	}
	candidates := make(map[string]map[string]string)
	for key, value := range sourceKeys {
		base, form, ok := splitI18nextPluralKey(key)
		if !ok {
			continue
		}
		if _, sourceForm := forms[form]; !sourceForm {
			continue
		}
		if candidates[base] == nil {
			candidates[base] = make(map[string]string)
		}
		candidates[base][form] = value
	}
	for base, values := range candidates {
		for _, form := range sourceForms {
			if _, ok := values[form]; !ok {
				delete(candidates, base)
				break
			}
		}
	}
	return candidates
}

func cloneStrings(values map[string]string) map[string]string {
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func splitI18nextPluralKey(key string) (string, string, bool) {
	separator := strings.LastIndexByte(key, '_')
	if separator <= 0 || separator == len(key)-1 {
		return "", "", false
	}
	category := key[separator+1:]
	for _, known := range []string{"zero", "one", "two", "few", "many", "other"} {
		if category == known {
			return key[:separator], category, true
		}
	}
	return "", "", false
}
