package formats

import "testing"

func TestYAMLParseExcludesNonStringLeaves(t *testing.T) {
	f := &YAMLFormat{}
	entries, err := f.Parse([]byte("label: Save\nenabled: true\nlimit: 3\nempty: null\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries["label"] != "Save" {
		t.Fatalf("entries = %#v, want only the translatable string", entries)
	}
}
