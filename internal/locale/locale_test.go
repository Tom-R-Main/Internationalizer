package locale

import (
	"slices"
	"testing"
)

func TestCanonical(t *testing.T) {
	tests := map[string]string{
		"EN-us":      "en-US",
		"sr-latn-rs": "sr-Latn-RS",
		"zh-hant-tw": "zh-Hant-TW",
	}
	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			got, err := Canonical(input)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("Canonical(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestCanonicalRejectsMalformedTags(t *testing.T) {
	for _, input := range []string{"", "en_US", "../fr", "en--US"} {
		t.Run(input, func(t *testing.T) {
			if _, err := Canonical(input); err == nil {
				t.Fatalf("Canonical(%q) succeeded", input)
			}
		})
	}
}

func TestCardinalCategoriesUseCLDRRules(t *testing.T) {
	tests := map[string][]string{
		"ar": {"zero", "one", "two", "few", "many", "other"},
		"ja": {"other"},
		"ru": {"one", "few", "many", "other"},
		"te": {"one", "other"},
	}
	for tag, want := range tests {
		t.Run(tag, func(t *testing.T) {
			got, err := CardinalCategories(tag)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, want) {
				t.Fatalf("CardinalCategories(%q) = %v, want %v", tag, got, want)
			}
		})
	}
}

func TestOrdinalCategoriesUseCLDRRules(t *testing.T) {
	got, err := OrdinalCategories("en-US")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"one", "two", "few", "other"}
	if !slices.Equal(got, want) {
		t.Fatalf("OrdinalCategories(en-US) = %v, want %v", got, want)
	}
}

func TestSupportedLocalesMatchCLDR49CategorySets(t *testing.T) {
	tests := map[string]struct {
		cardinal []string
		ordinal  []string
	}{
		"ar":    {[]string{"zero", "one", "two", "few", "many", "other"}, []string{"other"}},
		"bn":    {[]string{"one", "other"}, []string{"one", "two", "few", "many", "other"}},
		"cs":    {[]string{"one", "few", "many", "other"}, []string{"other"}},
		"da":    {[]string{"one", "other"}, []string{"other"}},
		"de":    {[]string{"one", "other"}, []string{"other"}},
		"el":    {[]string{"one", "other"}, []string{"other"}},
		"es":    {[]string{"one", "many", "other"}, []string{"other"}},
		"fi":    {[]string{"one", "other"}, []string{"other"}},
		"fr":    {[]string{"one", "many", "other"}, []string{"one", "other"}},
		"he":    {[]string{"one", "two", "other"}, []string{"other"}},
		"hi":    {[]string{"one", "other"}, []string{"one", "two", "few", "many", "other"}},
		"id":    {[]string{"other"}, []string{"other"}},
		"it":    {[]string{"one", "many", "other"}, []string{"many", "other"}},
		"ja":    {[]string{"other"}, []string{"other"}},
		"ko":    {[]string{"other"}, []string{"other"}},
		"kw":    {[]string{"zero", "one", "two", "few", "many", "other"}, []string{"one", "many", "other"}},
		"lld":   {[]string{"one", "many", "other"}, []string{"many", "other"}},
		"ms":    {[]string{"other"}, []string{"one", "other"}},
		"mt":    {[]string{"one", "two", "few", "many", "other"}, []string{"other"}},
		"nl":    {[]string{"one", "other"}, []string{"other"}},
		"pa":    {[]string{"one", "other"}, []string{"other"}},
		"pl":    {[]string{"one", "few", "many", "other"}, []string{"other"}},
		"pt-BR": {[]string{"one", "many", "other"}, []string{"other"}},
		"ro":    {[]string{"one", "few", "other"}, []string{"one", "other"}},
		"ru":    {[]string{"one", "few", "many", "other"}, []string{"other"}},
		"scn":   {[]string{"one", "many", "other"}, []string{"many", "other"}},
		"sv":    {[]string{"one", "other"}, []string{"one", "other"}},
		"te":    {[]string{"one", "other"}, []string{"other"}},
		"th":    {[]string{"other"}, []string{"other"}},
		"tr":    {[]string{"one", "other"}, []string{"other"}},
		"uk":    {[]string{"one", "few", "many", "other"}, []string{"few", "other"}},
		"vec":   {[]string{"one", "many", "other"}, []string{"many", "other"}},
		"vi":    {[]string{"other"}, []string{"one", "other"}},
		"yue":   {[]string{"other"}, []string{"other"}},
		"zh-CN": {[]string{"other"}, []string{"other"}},
		"zh-TW": {[]string{"other"}, []string{"other"}},
	}

	for tag, want := range tests {
		t.Run(tag, func(t *testing.T) {
			cardinal, err := CardinalCategories(tag)
			if err != nil {
				t.Fatal(err)
			}
			ordinal, err := OrdinalCategories(tag)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(cardinal, want.cardinal) || !slices.Equal(ordinal, want.ordinal) {
				t.Fatalf("%s categories: cardinal=%v ordinal=%v, want cardinal=%v ordinal=%v", tag, cardinal, ordinal, want.cardinal, want.ordinal)
			}
		})
	}
}
