package message

import (
	"slices"
	"strings"
	"testing"
)

func TestParseAndPrintNestedICUMessage(t *testing.T) {
	input := "{gender, select, female {{count, plural, offset:1 =0 {None} one {She and one other} other {She and # others}}} other {{count, number, ::compact-short} items}}"
	parsed, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	printed := parsed.String()
	reparsed, err := Parse(printed)
	if err != nil {
		t.Fatalf("printed message is invalid: %v\n%s", err, printed)
	}
	if reparsed.String() != printed {
		t.Fatalf("printer is unstable:\nfirst:  %s\nsecond: %s", printed, reparsed.String())
	}
}

func TestParseRejectsMalformedMessages(t *testing.T) {
	tests := []string{
		"{count, plural, one {One}",
		"{count, plural, one {One}}",
		"{count, plural other {Items}}",
		"{kind, select, other {One} other {Two}}",
		"{count, plural, offset:nope other {Items}}",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, err := Parse(input); err == nil {
				t.Fatalf("Parse(%q) succeeded", input)
			}
		})
	}
}

func TestParseRejectsExcessiveNesting(t *testing.T) {
	input := strings.Repeat("{value, select, other {", maxNestingDepth+1) + "value" + strings.Repeat("}}", maxNestingDepth+1)
	if _, err := Parse(input); err == nil || !strings.Contains(err.Error(), "nesting exceeds") {
		t.Fatalf("Parse error = %v, want nesting limit error", err)
	}
}

func TestCompareChecksArgumentsAndSelectors(t *testing.T) {
	issues := Compare(
		"{gender, select, female {{count, plural, one {One} other {#}}} other {They}}",
		"{gender, select, female {{count, select, one {Un} other {Plusieurs}}} male {Il} other {Iel}}",
		"fr",
	)
	got := make([]Code, len(issues))
	for index, issue := range issues {
		got[index] = issue.Code
	}
	if !slices.Contains(got, CodeSelectorMismatch) || !slices.Contains(got, CodeArgumentTypeMismatch) {
		t.Fatalf("Compare issues = %#v", issues)
	}
}

func TestCompareRejectsLocaleInapplicablePluralCategory(t *testing.T) {
	issues := Compare(
		"{count, plural, one {One} other {# items}}",
		"{count, plural, one {ఒక అంశం} few {కొన్ని అంశాలు} other {# అంశాలు}}",
		"te",
	)
	if !hasCode(issues, CodeInvalidPluralCategory) {
		t.Fatalf("Compare issues = %#v, want %q", issues, CodeInvalidPluralCategory)
	}
}

func TestCompareAllowsLocaleCategoryAndExactSelector(t *testing.T) {
	issues := Compare(
		"{count, plural, one {One} other {# items}}",
		"{count, plural, =0 {Нет элементов} one {Один элемент} few {# элемента} other {# элементов}}",
		"ru",
	)
	if len(issues) != 0 {
		t.Fatalf("Compare returned issues for valid Russian message: %#v", issues)
	}
}

func TestCompareRejectsArgumentIntroducedByTargetOnlyExactSelector(t *testing.T) {
	issues := Compare(
		"{n, plural, other {# items}}",
		"{n, plural, =0 {{price, number}} other {# articles}}",
		"fr",
	)
	if !hasCode(issues, CodeArgumentMismatch) {
		t.Fatalf("Compare issues = %#v, want target-only exact argument mismatch", issues)
	}
}

func TestParseRejectsNonFiniteAndMalformedExactSelectors(t *testing.T) {
	for _, selector := range []string{"NaN", "Inf", "Infinity", "1.", "0x1p2"} {
		t.Run(selector, func(t *testing.T) {
			if _, err := Parse("{n, plural, =" + selector + " {bad} other {ok}}"); err == nil {
				t.Fatalf("Parse accepted invalid exact selector %q", selector)
			}
		})
	}
	for _, selector := range []string{"0", "-1", "1.5"} {
		t.Run("valid_"+selector, func(t *testing.T) {
			if _, err := Parse("{n, plural, =" + selector + " {ok} other {other}}"); err != nil {
				t.Fatalf("Parse rejected valid exact selector %q: %v", selector, err)
			}
		})
	}
}

func TestCompareChecksTargetOnlyLocaleCategoryAgainstSourceOther(t *testing.T) {
	source := "{count, plural, one {{name} has one item} other {{name} has # items}}"
	valid := "{count, plural, one {{name} имеет один предмет} few {{name} имеет # предмета} other {{name} имеет # предметов}}"
	if issues := Compare(source, valid, "ru"); len(issues) != 0 {
		t.Fatalf("Compare rejected structurally valid target-only category: %#v", issues)
	}

	invalid := "{count, plural, one {{name} имеет один предмет} few {несколько предметов} other {{name} имеет # предметов}}"
	issues := Compare(source, invalid, "ru")
	if !hasCode(issues, CodeArgumentMismatch) || !hasCode(issues, CodeSelectorMismatch) {
		t.Fatalf("Compare issues = %#v, want missing argument and pound placeholder", issues)
	}
}

func TestComparePairsRepeatedPluralsWithTargetOnlyCategories(t *testing.T) {
	source := "{n, plural, one {{a}} other {{a} #}} then {n, plural, one {{z}} other {{z} #}}"
	target := "{n, plural, one {{z}} few {{z} #} other {{z} #}} puis {n, plural, one {{a}} other {{a} #}}"
	if issues := Compare(source, target, "ru"); len(issues) != 0 {
		t.Fatalf("Compare rejected valid repeated plural reordering: %#v", issues)
	}
}

func TestComparePreservesPluralOffsetsAndFormatterStyles(t *testing.T) {
	tests := []struct {
		name   string
		source string
		target string
		code   Code
	}{
		{
			name:   "plural offset",
			source: "{count, plural, offset:1 one {One} other {#}}",
			target: "{count, plural, one {Un} other {#}}",
			code:   CodeSelectorMismatch,
		},
		{
			name:   "number style",
			source: "Total: {amount, number, ::currency/USD}",
			target: "Total : {amount, number, ::percent}",
			code:   CodeArgumentStyleMismatch,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			issues := Compare(test.source, test.target, "fr")
			if !hasCode(issues, test.code) {
				t.Fatalf("Compare issues = %#v, want %q", issues, test.code)
			}
		})
	}
}

func TestComparePreservesArgumentsWithinCorrespondingBranches(t *testing.T) {
	source := "{gender, select, female {{count, plural, offset:1 one {One} other {#}}} other {{count, plural, one {One} other {#}}}}"
	target := "{gender, select, female {{count, plural, one {Une} other {#}}} other {{count, plural, offset:1 one {Une} other {#}}}}"
	issues := Compare(source, target, "fr")
	if !hasCode(issues, CodeSelectorMismatch) {
		t.Fatalf("Compare issues = %#v, want branch-specific offset mismatch", issues)
	}
}

func TestCompareAllowsRepeatedArgumentReorderingWithinBranch(t *testing.T) {
	source := "Created {date, date, short}; archived {date, date, long}"
	target := "Archivé {date, date, long}; créé {date, date, short}"
	if issues := Compare(source, target, "fr"); len(issues) != 0 {
		t.Fatalf("Compare rejected valid repeated-argument reordering: %#v", issues)
	}
}

func TestCompareAllowsRepeatedNestedArgumentReorderingWithinBranch(t *testing.T) {
	source := "{choice, select, first {{first}} other {A}} then {choice, select, first {{last}} other {B}}"
	target := "{choice, select, first {{last}} other {B traduit}} puis {choice, select, first {{first}} other {A traduit}}"
	if issues := Compare(source, target, "fr"); len(issues) != 0 {
		t.Fatalf("Compare rejected valid repeated nested-argument reordering: %#v", issues)
	}
}

func TestComparePreservesPoundPlaceholders(t *testing.T) {
	issues := Compare(
		"{count, plural, one {One item} other {# items}}",
		"{count, plural, one {Un article} other {articles}}",
		"fr",
	)
	if !hasCode(issues, CodeSelectorMismatch) {
		t.Fatalf("Compare issues = %#v, want missing pound placeholder", issues)
	}
}

func TestComparePreservesPluralPoundThroughNestedSelect(t *testing.T) {
	source := "{n, plural, other {{gender, select, male {# men} other {# people}}}}"
	target := "{n, plural, other {{gender, select, male {hommes} other {# personnes}}}}"
	issues := Compare(source, target, "fr")
	if !hasCode(issues, CodeSelectorMismatch) {
		t.Fatalf("Compare issues = %#v, want nested missing pound placeholder", issues)
	}
}

func TestCompareSupportsUnicodeArgumentNames(t *testing.T) {
	source := "Bienvenue, {имя}"
	target := "Bienvenue"
	if !LooksLike(source) {
		t.Fatalf("LooksLike(%q) = false", source)
	}
	issues := Compare(source, target, "fr")
	if !hasCode(issues, CodeArgumentMismatch) {
		t.Fatalf("Compare issues = %#v, want missing Unicode argument", issues)
	}

	parsed, err := Parse("{имя, select, admin {Привет, {имя}} other {Здравствуйте, {имя}}}")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(parsed.String()); err != nil {
		t.Fatalf("Parse(String()) failed: %v", err)
	}
}

func TestCompareReportsMalformedSourceAndTarget(t *testing.T) {
	for name, values := range map[string][2]string{
		"source": {"{count, plural, one {One}}", "{count, plural, one {Un} other {Autres}}"},
		"target": {"{count, plural, one {One} other {Other}}", "{count, plural, one {Un} other {Autres}"},
	} {
		t.Run(name, func(t *testing.T) {
			issues := Compare(values[0], values[1], "fr")
			if !hasCode(issues, CodeSyntax) {
				t.Fatalf("Compare issues = %#v, want syntax issue", issues)
			}
		})
	}
}

func hasCode(issues []Issue, code Code) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func FuzzParsePrint(f *testing.F) {
	for _, seed := range []string{
		"plain text",
		"Hello {name}",
		"{gender, select, female {She} other {They}}",
		"{count, plural, offset:1 =0 {None} one {One} other {# items}}",
		"{place, selectordinal, one {#st} two {#nd} few {#rd} other {#th}}",
		"This '{' is literal and this apostrophe isn''t ambiguous",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := Parse(input)
		if err != nil {
			return
		}
		printed := parsed.String()
		reparsed, err := Parse(printed)
		if err != nil {
			t.Fatalf("Parse(String(Parse(%q))) failed: %v; printed %q", input, err, printed)
		}
		if reparsed.String() != printed {
			t.Fatalf("unstable print for %q: first %q, second %q", input, printed, reparsed.String())
		}
	})
}
