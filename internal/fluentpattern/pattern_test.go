package fluentpattern

import (
	"strings"
	"testing"
)

func TestCompareAllowsLocaleSpecificVariantsWithDefaultStructure(t *testing.T) {
	source := `{ $count ->
    [one] One item for { $name }
   *[other] { $count } items for { $name }
}`
	target := `{ $count ->
    [one] Один элемент для { $name }
    [few] { $count } элемента для { $name }
    [many] { $count } элементов для { $name }
   *[other] { $count } элемента для { $name }
}`
	_, _, preserved, err := Compare(source, target)
	if err != nil || !preserved {
		t.Fatalf("Compare() = preserved %v, error %v", preserved, err)
	}
}

func TestCompareRejectsDamagedFluentRuntimeStructure(t *testing.T) {
	source := `{ PLATFORM() ->
    [windows] Open { -brand-short-name } for { $user }
   *[other] Open { -brand-short-name } for { $user }
}`
	tests := map[string]string{
		"selector": `{ $platform ->
    [windows] Ouvrir { -brand-short-name } pour { $user }
   *[other] Ouvrir { -brand-short-name } pour { $user }
}`,
		"reference": `{ PLATFORM() ->
    [windows] Ouvrir { -brand-full-name } pour { $user }
   *[other] Ouvrir { -brand-short-name } pour { $user }
}`,
		"variable": `{ PLATFORM() ->
    [windows] Ouvrir { -brand-short-name } pour { $account }
   *[other] Ouvrir { -brand-short-name } pour { $user }
}`,
		"default": `{ PLATFORM() ->
   *[windows] Ouvrir { -brand-short-name } pour { $user }
    [other] Ouvrir { -brand-short-name } pour { $user }
}`,
	}
	for name, target := range tests {
		t.Run(name, func(t *testing.T) {
			_, _, preserved, err := Compare(source, target)
			if err == nil && preserved {
				t.Fatal("damaged Fluent structure was accepted")
			}
		})
	}
}

func TestAnalyzeRejectsMalformedPatterns(t *testing.T) {
	for _, pattern := range []string{
		`Hello { $name`,
		`Hello }`,
		`{ $count -> [one] One }`,
		"{ $count ->\n   *[one] One\n   *[other] Other\n}",
	} {
		if _, err := Analyze(pattern); err == nil {
			t.Fatalf("Analyze(%q) accepted malformed pattern", pattern)
		}
	}
}

func TestTransformTextPreservesFluentSyntaxInsideSelectors(t *testing.T) {
	input := `{ $count ->
    [one] One item for { -brand-short-name }
   *[other] { $count } items for <strong>{ $user }</strong>
}`
	output, err := TransformText(input, strings.ToUpper)
	if err != nil {
		t.Fatal(err)
	}
	for _, syntax := range []string{"{ $count ->", "[one]", "*[other]", "{ -brand-short-name }", "{ $count }", "{ $user }"} {
		if !strings.Contains(output, syntax) {
			t.Fatalf("transformed pattern lost %q: %s", syntax, output)
		}
	}
	if !strings.Contains(output, "ONE ITEM FOR") || !strings.Contains(output, "ITEMS FOR <STRONG>") {
		t.Fatalf("literal text was not transformed: %s", output)
	}
}
