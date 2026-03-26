package validate

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
	"he": {"one", "two", "many", "other"},

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
	if forms, ok := CLDRPluralForms[locale]; ok {
		return forms
	}
	return []string{"one", "other"}
}
