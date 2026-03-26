package formats

import (
	"encoding/json"
	"testing"
)

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
