package validate

import "strings"

// CLDRPluralForms maps locale codes to their required CLDR plural categories.
// Source: Unicode CLDR plural rules for the most common languages.
var CLDRPluralForms = map[string][]string{
	// Germanic
	"en":    {"one", "other"},
	"en-US": {"one", "other"},
	"en-GB": {"one", "other"},
	"de":    {"one", "other"},
	"nl":    {"one", "other"},
	"sv":    {"one", "other"},
	"da":    {"one", "other"},
	"nb":    {"one", "other"},
	"fi":    {"one", "other"},

	// Romance
	"fr":    {"one", "many", "other"},
	"fr-CA": {"one", "many", "other"},
	"es":    {"one", "many", "other"},
	"es-MX": {"one", "many", "other"},
	"pt":    {"one", "many", "other"},
	"pt-BR": {"one", "many", "other"},
	"it":    {"one", "many", "other"},
	"ro":    {"one", "few", "other"},

	// Slavic
	"ru": {"one", "few", "many", "other"},
	"uk": {"one", "few", "many", "other"},
	"pl": {"one", "few", "many", "other"},
	"cs": {"one", "few", "many", "other"},

	// Semitic
	"ar": {"zero", "one", "two", "few", "many", "other"},
	"he": {"one", "two", "other"},

	// Indic
	"hi": {"one", "other"},
	"bn": {"one", "other"},
	"pa": {"one", "other"},
	"te": {"one", "other"},

	// CJK (no plural forms — only "other")
	"ja":    {"other"},
	"ko":    {"other"},
	"zh":    {"other"},
	"zh-CN": {"other"},
	"zh-TW": {"other"},
	"yue":   {"other"},

	// Other
	"tr": {"one", "other"},
	"id": {"other"},
	"vi": {"other"},
	"th": {"other"},
	"el": {"one", "other"},
	"hu": {"one", "other"},
}

// PluralFormsFor returns the CLDR plural categories for a locale.
// Falls back to ["one", "other"] if the locale is not in the table.
func PluralFormsFor(locale string) []string {
	if forms, ok := KnownPluralFormsFor(locale); ok {
		return forms
	}
	return []string{"one", "other"}
}

// KnownPluralFormsFor returns configured CLDR categories without guessing for
// unknown locales. Strict validation uses this form to avoid false failures.
func KnownPluralFormsFor(locale string) ([]string, bool) {
	if forms, ok := CLDRPluralForms[locale]; ok {
		return forms, true
	}
	for configuredLocale, forms := range CLDRPluralForms {
		if strings.EqualFold(configuredLocale, locale) {
			return forms, true
		}
	}
	if separator := strings.IndexAny(locale, "-_"); separator > 0 {
		forms, ok := CLDRPluralForms[strings.ToLower(locale[:separator])]
		return forms, ok
	}
	return nil, false
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
