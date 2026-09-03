package formats

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

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

func TestYAMLSerializePreservingAddsMissingEntries(t *testing.T) {
	f := &YAMLFormat{}
	output, err := f.Serialize(map[string]string{
		"items_one":  "{{count}} article",
		"items_many": "{{count}} articles",
	}, []byte("# translations\nitems_one: '{{count}} item'\nenabled: true\n"))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := f.Parse(output)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["items_one"] != "{{count}} article" || parsed["items_many"] != "{{count}} articles" {
		t.Fatalf("serialized entries = %#v\n%s", parsed, output)
	}
	if !strings.Contains(string(output), "# translations") || !strings.Contains(string(output), "enabled: true") {
		t.Fatalf("serialization lost source shape:\n%s", output)
	}
}

func TestYAMLSerializePreservingAddsMissingEntriesInsideSequence(t *testing.T) {
	f := &YAMLFormat{}
	output, err := f.Serialize(map[string]string{
		"screens.0.items_one":  "{{count}} article",
		"screens.0.items_many": "{{count}} articles",
	}, []byte("screens:\n  - items_one: '{{count}} item'\n"))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := f.Parse(output)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["screens.0.items_many"] != "{{count}} articles" {
		t.Fatalf("serialized entries = %#v\n%s", parsed, output)
	}
}

func TestYAMLSerializePreservingCreatesMissingSequenceBranch(t *testing.T) {
	f := &YAMLFormat{}
	output, err := f.Serialize(map[string]string{
		"title":           "Accueil",
		"screens.0.title": "Paramètres",
	}, []byte("title: Home\n"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Screens []map[string]string `yaml:"screens"`
	}
	if err := yaml.Unmarshal(output, &decoded); err != nil {
		t.Fatalf("missing sequence branch changed shape: %v\n%s", err, output)
	}
	if len(decoded.Screens) != 1 || decoded.Screens[0]["title"] != "Paramètres" {
		t.Fatalf("missing sequence branch = %#v\n%s", decoded.Screens, output)
	}
}

func TestYAMLSerializePreservingKeepsNumericMappingKeys(t *testing.T) {
	f := &YAMLFormat{}
	output, err := f.Serialize(map[string]string{
		"http.200": "OK",
		"http.404": "Introuvable",
	}, []byte("http:\n  '200': OK\n  '500': Error\n"))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := f.Parse(output)
	if err != nil {
		t.Fatal(err)
	}
	if parsed["http.200"] != "OK" || parsed["http.404"] != "Introuvable" || parsed["http.500"] != "Error" {
		t.Fatalf("numeric mapping changed shape: %#v\n%s", parsed, output)
	}
}
