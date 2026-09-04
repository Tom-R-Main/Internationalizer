package formats

import (
	"reflect"
	"strings"
	"testing"
)

func TestMarkdownParseKeepsUnchangedSectionsStable(t *testing.T) {
	format := &MarkdownFormat{}
	before := []byte("# Product\n\nIntro.\n\n## Install\n\nRun this.\n\n## Usage\n\nUse it.\n")
	after := []byte("# Product\n\nIntro.\n\n## Install\n\nRun this command.\n\n## Usage\n\nUse it.\n")

	beforeEntries, err := format.Parse(before)
	if err != nil {
		t.Fatal(err)
	}
	afterEntries, err := format.Parse(after)
	if err != nil {
		t.Fatal(err)
	}

	if len(beforeEntries) != 3 {
		t.Fatalf("entry count = %d, want 3", len(beforeEntries))
	}
	if _, ok := beforeEntries["markdown:install"]; !ok {
		t.Fatalf("entries = %#v, want markdown:install", beforeEntries)
	}
	if beforeEntries["markdown:preamble"] != afterEntries["markdown:preamble"] {
		t.Fatal("preamble changed when only the install section changed")
	}
	if beforeEntries["markdown:usage"] != afterEntries["markdown:usage"] {
		t.Fatal("usage changed when only the install section changed")
	}
	if beforeEntries["markdown:install"] == afterEntries["markdown:install"] {
		t.Fatal("install section did not change")
	}
}

func TestMarkdownPairedSerializationUsesStableMarkers(t *testing.T) {
	format := &MarkdownFormat{}
	source := []byte("# Product\n\nIntro.\n\n## Install\n\nRun this.\n\n## Usage\n\nUse it.\n")
	target := []byte("# Produit\n\nIntroduction.\n\n## Installation\n\nExécutez ceci.\n\n## Utilisation\n\nUtilisez-le.\n")

	entries, err := format.ParseTarget(source, target)
	if err != nil {
		t.Fatal(err)
	}
	if got := entries["markdown:install"]; !strings.Contains(got, "Installation") {
		t.Fatalf("install target = %q", got)
	}

	output, err := format.SerializeTarget(entries, source, target)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "<!-- internationalizer:unit markdown:install -->\n## Installation") {
		t.Fatalf("output lacks stable install marker:\n%s", output)
	}

	reparsed, err := format.ParseTarget(source, output)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(reparsed, entries) {
		t.Fatalf("reparsed = %#v, want %#v", reparsed, entries)
	}
}

func TestMarkdownMarkersSurviveSourceSectionInsertion(t *testing.T) {
	format := &MarkdownFormat{}
	before := []byte("# Product\n\n## Install\n\nRun this.\n\n## Usage\n\nUse it.\n")
	target := []byte("# Produit\n\n## Installation\n\nExécutez ceci.\n\n## Utilisation\n\nUtilisez-le.\n")
	entries, err := format.ParseTarget(before, target)
	if err != nil {
		t.Fatal(err)
	}
	marked, err := format.SerializeTarget(entries, before, target)
	if err != nil {
		t.Fatal(err)
	}

	after := []byte("# Product\n\n## Install\n\nRun this.\n\n## Configuration\n\nConfigure it.\n\n## Usage\n\nUse it.\n")
	afterEntries, err := format.ParseTarget(after, marked)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := afterEntries["markdown:configuration"]; ok {
		t.Fatal("new source section unexpectedly has a target")
	}
	if got := afterEntries["markdown:usage"]; !strings.Contains(got, "Utilisation") {
		t.Fatalf("usage target lost after insertion: %q", got)
	}
}
