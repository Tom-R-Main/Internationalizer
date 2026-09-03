package validate

import (
	"reflect"
	"testing"
)

func TestKnownPluralFormsForFallsBackToBaseLocaleWithoutGuessing(t *testing.T) {
	forms, ok := KnownPluralFormsFor("ru-RU")
	if !ok || len(forms) != 4 || forms[1] != "few" {
		t.Fatalf("Russian forms = %v, known = %v", forms, ok)
	}
	if forms, ok := KnownPluralFormsFor("xx-ZZ"); ok || forms != nil {
		t.Fatalf("unknown locale forms = %v, known = %v", forms, ok)
	}
}

func TestExpandI18nextV4SourceUsesLocalePluralFamilies(t *testing.T) {
	expanded, required, optional := ExpandI18nextV4Source(
		map[string]string{"items_other": "{{count}} 個"},
		"ja",
		"ar",
	)
	wantKeys := []string{"items_few", "items_many", "items_one", "items_other", "items_two", "items_zero"}
	var gotKeys []string
	for _, key := range wantKeys {
		if expanded[key] != "{{count}} 個" {
			t.Fatalf("expanded[%q] = %q", key, expanded[key])
		}
		if _, ok := required[key]; !ok {
			t.Fatalf("required lacks %q: %#v", key, required)
		}
		gotKeys = append(gotKeys, key)
	}
	if !reflect.DeepEqual(gotKeys, wantKeys) || len(optional) != 0 || len(expanded) != len(wantKeys) {
		t.Fatalf("expanded = %#v, required = %#v, optional = %#v", expanded, required, optional)
	}
}

func TestExpandI18nextV4SourceMarksSourceOnlyFormsOptional(t *testing.T) {
	expanded, required, optional := ExpandI18nextV4Source(
		map[string]string{"items_one": "{{count}} item", "items_other": "{{count}} items"},
		"en",
		"ja",
	)
	if len(expanded) != 2 || len(required) != 1 {
		t.Fatalf("expanded = %#v, required = %#v", expanded, required)
	}
	if _, ok := optional["items_one"]; !ok {
		t.Fatalf("optional = %#v, want items_one", optional)
	}
}

func TestHebrewPluralFormsMatchI18nextV4CardinalCategories(t *testing.T) {
	forms, ok := KnownPluralFormsFor("he")
	if !ok || !reflect.DeepEqual(forms, []string{"one", "two", "other"}) {
		t.Fatalf("Hebrew forms = %v, known = %t", forms, ok)
	}
}

func TestKnownPluralFormsForAcceptsCaseInsensitiveLocaleTags(t *testing.T) {
	for _, locale := range []string{"FR", "PT-br", "ja-JP"} {
		if _, ok := KnownPluralFormsFor(locale); !ok {
			t.Fatalf("KnownPluralFormsFor(%q) did not normalize locale casing", locale)
		}
	}
}
