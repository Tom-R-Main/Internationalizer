package formats

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestJSONRejectsIntegrityLoss(t *testing.T) {
	cases := []struct{ name, input, code string }{
		{"duplicate", `{"key":"{{name}}","key":"Hello"}`, "json_duplicate_member"},
		{"escaped duplicate", `{"key":"one","\u006bey":"two"}`, "json_duplicate_member"},
		{"nested duplicate", `{"a":[{"key":"one","key":"two"}]}`, "json_duplicate_member"},
		{"flattened collision", `{"a.b":"one","a":{"b":"two"}}`, "json_flattened_key_collision"},
		{"array collision", `{"a.0":"one","a":["two"]}`, "json_flattened_key_collision"},
		{"metadata collision", `{"a.b":false,"a":{"b":"two"}}`, "json_flattened_key_collision"},
		{"empty key collision", `{"":{"a":"one"},"a":"two"}`, "json_flattened_key_collision"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &JSONFormat{}
			operations := []func() error{
				func() error { _, err := f.Parse([]byte(tc.input)); return err },
				func() error { _, err := f.Serialize(map[string]string{}, []byte(tc.input)); return err },
				func() error { _, err := f.RemoveEntries([]byte(tc.input), map[string]struct{}{}); return err },
			}
			for i, operation := range operations {
				var last string
				for range 100 {
					err := operation()
					var coded interface{ JSONCode() string }
					if !errors.As(err, &coded) || coded.JSONCode() != tc.code {
						t.Fatalf("operation %d: got %v, want %s", i, err, tc.code)
					}
					if last != "" && last != err.Error() {
						t.Fatalf("nondeterministic error: %q != %q", last, err.Error())
					}
					last = err.Error()
				}
			}
		})
	}
}

func TestJSONRejectsTrailingContent(t *testing.T) {
	f := &JSONFormat{}
	for _, input := range []string{`{"key":"one"} {"key":"two"}`, `{"key":"one"} trailing`} {
		if _, err := f.Serialize(nil, []byte(input)); err == nil {
			t.Fatal("Serialize accepted trailing JSON")
		}
		if _, err := f.RemoveEntries([]byte(input), nil); err == nil {
			t.Fatal("RemoveEntries accepted trailing JSON")
		}
	}
}

func TestJSONRejectsUnboundedNesting(t *testing.T) {
	if _, err := (&JSONFormat{}).Parse([]byte(strings.Repeat("[", 300) + `"leaf"` + strings.Repeat("]", 300))); err == nil {
		t.Fatal("accepted excessive nesting")
	}
}

func TestJSONSerializeRejectsDestructiveInsertion(t *testing.T) {
	cases := []struct {
		original string
		entries  map[string]string
	}{
		{`{"a":{"b":"Keep"}}`, map[string]string{"a": "Lose child"}},
		{`{"a":null}`, map[string]string{"a": "Lose metadata"}},
		{`{"a":null}`, map[string]string{"a.b": "Lose metadata"}},
		{`{"a":42}`, map[string]string{"a": "Lose metadata"}},
		{`{"a":[null]}`, map[string]string{"a.0": "Lose metadata"}},
		{`{"a":[null]}`, map[string]string{"a.0.b": "Lose metadata"}},
		{`null`, map[string]string{"a": "Lose root metadata"}},
		{`{"a":["Keep"]}`, map[string]string{"a.-1": "Invalid index"}},
		{`{"a":["Keep"]}`, map[string]string{"a.999999999": "Invalid index"}},
		{`{"a":["Keep"]}`, map[string]string{"a.00": "Alias"}},
		{`{}`, map[string]string{".a": "Lost identity"}},
		{"", map[string]string{"a": "one", "a.b": "two"}},
	}
	for _, tc := range cases {
		if _, err := (&JSONFormat{}).Serialize(tc.entries, []byte(tc.original)); err == nil {
			t.Fatalf("accepted destructive insertion into %s with keys %v", tc.original, tc.entries)
		}
	}
}

func TestJSONValidDottedContainersRoundTripAndRemoval(t *testing.T) {
	f := &JSONFormat{}
	for _, original := range []string{`{"a.b":{"c":"X"},"a":{"b":{"d":"Y"}}}`, `{"":{"":"X"}}`, `{"a.b":"X","a":{"b":{"d":"Y"}}}`} {
		entries, err := f.Parse([]byte(original))
		if err != nil {
			t.Fatal(err)
		}
		output, err := f.Serialize(entries, []byte(original))
		if err != nil {
			t.Fatal(err)
		}
		reparsed, err := f.Parse(output)
		if err != nil || !reflect.DeepEqual(entries, reparsed) {
			t.Fatalf("roundtrip: %v", err)
		}
	}
	output, err := f.RemoveEntries([]byte(`{"a.b":"X","a":{"b":{"d":"Y"}},"metadata":true}`), map[string]struct{}{"a.b": {}, "metadata": {}})
	if err != nil {
		t.Fatal(err)
	}
	entries, err := f.Parse(output)
	if err != nil || !reflect.DeepEqual(entries, map[string]string{"a.b.d": "Y"}) {
		t.Fatalf("removed unrelated content: %s %v", output, err)
	}
	var raw map[string]any
	if err := json.Unmarshal(output, &raw); err != nil {
		t.Fatal(err)
	}
	if raw["metadata"] != true {
		t.Fatal("removed nonstring metadata")
	}
}

func FuzzJSONCatalogRoundTrip(f *testing.F) {
	for _, input := range []string{`{"a.b":"dotted","a":{"c":"nested"}}`, `{"n":9007199254740993,"ok":true,"items":["a",null]}`, `{"":"empty"}`, `"root"`, `["root array"]`, `{"a":"one","a":"two"}`} {
		f.Add(input)
	}
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 65536 {
			t.Skip()
		}
		format := &JSONFormat{}
		entries, err := format.Parse([]byte(input))
		if err != nil {
			return
		}
		output, err := format.Serialize(entries, []byte(input))
		if err != nil {
			t.Fatalf("cannot serialize valid input: %v", err)
		}
		reparsed, err := format.Parse(output)
		if err != nil || !reflect.DeepEqual(entries, reparsed) {
			t.Fatalf("catalog identities changed: %v", err)
		}
		// Compare with json.Number decoding to preserve integers beyond float64.
		decode := func(data []byte) any {
			decoder := json.NewDecoder(strings.NewReader(string(data)))
			decoder.UseNumber()
			var value any
			if err := decoder.Decode(&value); err != nil {
				t.Fatal(err)
			}
			return value
		}
		if !reflect.DeepEqual(decode([]byte(input)), decode(output)) {
			t.Fatal("roundtrip changed non-string metadata or structure")
		}
	})
}

func TestJSONParse(t *testing.T) {
	input := `{
		"common": {
			"save": "Save",
			"cancel": "Cancel"
		},
		"dashboard": {
			"title": "Dashboard",
			"welcome": "Hello, {{name}}!"
		}
	}`

	f := &JSONFormat{}
	result, err := f.Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"common.save":       "Save",
		"common.cancel":     "Cancel",
		"dashboard.title":   "Dashboard",
		"dashboard.welcome": "Hello, {{name}}!",
	}

	for key, want := range expected {
		got, ok := result[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", key, got, want)
		}
	}

	if len(result) != len(expected) {
		t.Errorf("got %d keys, want %d", len(result), len(expected))
	}
}

func TestJSONParseExcludesNonStringLeaves(t *testing.T) {
	f := &JSONFormat{}
	entries, err := f.Parse([]byte(`{"label":"Save","enabled":true,"limit":3,"empty":null}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries["label"] != "Save" {
		t.Fatalf("entries = %#v, want only the translatable string", entries)
	}
	output, err := f.Serialize(map[string]string{"label": "Enregistrer"}, []byte(`{"label":"Save","enabled":true,"limit":3,"empty":null}`))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(output, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["enabled"] != true || decoded["limit"] != float64(3) || decoded["empty"] != nil {
		t.Fatalf("non-string values changed: %#v", decoded)
	}
}

func TestJSONRoundTripsRootString(t *testing.T) {
	f := &JSONFormat{}
	entries, err := f.Parse([]byte(`"Save"`))
	if err != nil {
		t.Fatal(err)
	}
	entries[""] = "Enregistrer"
	output, err := f.Serialize(entries, []byte(`"Save"`))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != `"Enregistrer"` {
		t.Fatalf("Serialize() = %s, want translated root string", output)
	}
}

func TestJSONSerializePreservesOrder(t *testing.T) {
	original := `{
  "b": "B",
  "a": "A",
  "nested": {
    "z": "Z",
    "y": "Y"
  }
}`

	f := &JSONFormat{}
	entries, err := f.Parse([]byte(original))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries["b"] = "B-translated"
	entries["nested.z"] = "Z-translated"

	output, err := f.Serialize(entries, []byte(original))
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// Verify it's valid JSON.
	var check map[string]interface{}
	if err := json.Unmarshal(output, &check); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Verify values replaced.
	reparsed, _ := f.Parse(output)
	if reparsed["b"] != "B-translated" {
		t.Errorf("b not replaced: got %q", reparsed["b"])
	}
	if reparsed["nested.z"] != "Z-translated" {
		t.Errorf("nested.z not replaced: got %q", reparsed["nested.z"])
	}
	if reparsed["a"] != "A" {
		t.Errorf("a should be unchanged: got %q", reparsed["a"])
	}
}

func TestJSONRoundTrip(t *testing.T) {
	original := `{
  "greetings": {
    "hello": "Hello",
    "goodbye": "Goodbye"
  },
  "count_one": "{{count}} item",
  "count_other": "{{count}} items"
}`

	f := &JSONFormat{}
	entries, err := f.Parse([]byte(original))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	output, err := f.Serialize(entries, []byte(original))
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}

	reparsed, err := f.Parse(output)
	if err != nil {
		t.Fatalf("Re-parse: %v", err)
	}

	for key, want := range entries {
		if got := reparsed[key]; got != want {
			t.Errorf("round-trip key %q: got %q, want %q", key, got, want)
		}
	}
}

func TestJSONSerializeAddsMissingKeysToExistingFile(t *testing.T) {
	original := `{
  "common": {
    "save": "Enregistrer"
  },
  "dashboard": {
    "title": "Tableau de bord"
  }
}`

	f := &JSONFormat{}
	entries := map[string]string{
		"common.save":       "Enregistrer",
		"common.cancel":     "Annuler",
		"dashboard.title":   "Tableau de bord",
		"dashboard.welcome": "Bon retour, {{name}} !",
	}

	output, err := f.Serialize(entries, []byte(original))
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	reparsed, err := f.Parse(output)
	if err != nil {
		t.Fatalf("Re-parse failed: %v", err)
	}

	for key, want := range entries {
		if got := reparsed[key]; got != want {
			t.Errorf("key %q: got %q, want %q", key, got, want)
		}
	}
}

func TestJSONSerializeWithSourceStructurePreservesArrays(t *testing.T) {
	f := &JSONFormat{}
	output, err := f.Serialize(map[string]string{
		"steps.0.label": "Premier",
		"steps.1.label": "Deuxième",
	}, []byte(`{"steps":[{"label":"First"},{"label":"Second"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Steps []struct {
			Label string `json:"label"`
		} `json:"steps"`
	}
	if err := json.Unmarshal(output, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Steps) != 2 || decoded.Steps[1].Label != "Deuxième" {
		t.Fatalf("decoded = %#v", decoded)
	}
}

func TestJSONSerializePreservesNumericObjectKeys(t *testing.T) {
	f := &JSONFormat{}
	original := []byte(`{"http":{"404":"Not found"}}`)
	output, err := f.Serialize(map[string]string{"http.404": "Introuvable"}, original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded struct {
		HTTP map[string]string `json:"http"`
	}
	if err := json.Unmarshal(output, &decoded); err != nil {
		t.Fatal(err)
	}
	if got := decoded.HTTP["404"]; got != "Introuvable" {
		t.Fatalf("http.404 = %q, want %q; output = %s", got, "Introuvable", output)
	}
}

func TestJSONSerializeRejectsAmbiguousNumericPathWithoutStructure(t *testing.T) {
	f := &JSONFormat{}
	if _, err := f.Serialize(map[string]string{"steps.0": "Premier"}, nil); err == nil {
		t.Fatal("Serialize guessed an array shape without source structure")
	}
}
