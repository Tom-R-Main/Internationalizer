package formats

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const fluentFixture = `## Account panel

# Shown when the user has unread messages.
inbox-status =
    { $count ->
        [one] One unread message
       *[other] { $count } unread messages
    }
    .aria-label = Inbox status

save-button =
    .label = Save
    .accesskey = S

-brand-short-name = Acme
welcome = Welcome to { -brand-short-name }.
`

func TestFluentParseUnitsPreservesSemanticKindsAndContext(t *testing.T) {
	units, err := ParseUnits(&FluentFormat{}, []byte(fluentFixture))
	if err != nil {
		t.Fatal(err)
	}
	wantIDs := []string{"inbox-status", "inbox-status.aria-label", "save-button.label", "save-button.accesskey", "-brand-short-name", "welcome"}
	if len(units) != len(wantIDs) {
		t.Fatalf("units = %#v", units)
	}
	for index, want := range wantIDs {
		if units[index].ID != want || units[index].Structure != "fluent-pattern-v1" {
			t.Fatalf("unit %d = %#v, want ID %q", index, units[index], want)
		}
	}
	if units[0].Kind != UnitMessage || units[1].Kind != UnitAttribute || units[4].Kind != UnitTerm {
		t.Fatalf("unexpected semantic kinds: %#v", units)
	}
	if units[0].Context != "# Shown when the user has unread messages." {
		t.Fatalf("context = %q", units[0].Context)
	}
	if !strings.Contains(units[0].Value, "{ $count ->") || !strings.Contains(units[0].Value, "*[other]") {
		t.Fatalf("selector pattern = %q", units[0].Value)
	}
	if units[2].ID != "save-button.label" {
		t.Fatalf("attribute-only message produced unexpected units: %#v", units)
	}
}

func TestFluentUnchangedSerializationIsLossless(t *testing.T) {
	format := &FluentFormat{}
	units, err := format.ParseUnits([]byte(fluentFixture))
	if err != nil {
		t.Fatal(err)
	}
	output, err := format.SerializeUnits(units, []byte(fluentFixture))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != fluentFixture {
		t.Fatalf("unchanged round trip drifted:\n%s", output)
	}
}

func TestFluentSerializationChangesValuesWithoutLosingResourceStructure(t *testing.T) {
	format := &FluentFormat{}
	units, err := format.ParseUnits([]byte(fluentFixture))
	if err != nil {
		t.Fatal(err)
	}
	for index := range units {
		switch units[index].ID {
		case "save-button.label":
			units[index].Value = "Enregistrer"
		case "inbox-status":
			units[index].Value = "{ $count ->\n    [one] Un message non lu\n   *[other] { $count } messages non lus\n}"
		}
	}
	output, err := format.SerializeUnits(units, []byte(fluentFixture))
	if err != nil {
		t.Fatal(err)
	}
	text := string(output)
	for _, preserved := range []string{"## Account panel", "# Shown when", ".accesskey = S", "-brand-short-name = Acme"} {
		if !strings.Contains(text, preserved) {
			t.Fatalf("serialized resource lost %q:\n%s", preserved, text)
		}
	}
	reparsed, err := format.ParseUnits(output)
	if err != nil {
		t.Fatal(err)
	}
	values := UnitValues(reparsed)
	if values["save-button.label"] != "Enregistrer" || !strings.Contains(values["inbox-status"], "Un message non lu") {
		t.Fatalf("reparsed values = %#v", values)
	}
}

func TestFluentSerializationAppendsNewSourceEntriesInSourceOrder(t *testing.T) {
	format := &FluentFormat{}
	target := []byte("existing = Existant\n")
	baseline, err := format.ParseUnits(target)
	if err != nil {
		t.Fatal(err)
	}
	source, err := format.ParseUnits([]byte("existing = Existing\n\n# First note\nfirst = First\nsecond =\n    .label = Second\n"))
	if err != nil {
		t.Fatal(err)
	}
	values := map[string]string{"existing": "Existant", "first": "Premier", "second.label": "Deuxième"}
	merged := MergeUnitValues(baseline, source, values)
	if got := []string{merged[0].ID, merged[1].ID, merged[2].ID}; !reflect.DeepEqual(got, []string{"existing", "first", "second.label"}) {
		t.Fatalf("merged order = %v", got)
	}
	output, err := format.SerializeUnits(merged, target)
	if err != nil {
		t.Fatal(err)
	}
	text := string(output)
	if strings.Index(text, "first = Premier") > strings.Index(text, "second =") || !strings.Contains(text, "# First note") {
		t.Fatalf("new entries were not appended in source order with context:\n%s", text)
	}
	parsed, err := format.ParseUnits(output)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(UnitValues(parsed), values) {
		t.Fatalf("parsed appended values = %#v, want %#v", UnitValues(parsed), values)
	}
}

func TestFluentRejectsDuplicateUnitIdentity(t *testing.T) {
	_, err := (&FluentFormat{}).ParseUnits([]byte("same = One\nsame = Two\n"))
	if err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("ParseUnits() error = %v", err)
	}
}

func TestFluentFormatIsRegistered(t *testing.T) {
	format, err := FormatForFile("browser.ftl")
	if err != nil || format.Name() != "fluent" {
		t.Fatalf("FormatForFile() = %v, %v", format, err)
	}
}

func TestFluentOptionalCorpusRoundTrip(t *testing.T) {
	root := os.Getenv("INTERNATIONALIZER_FLUENT_CORPUS")
	if root == "" {
		t.Skip("set INTERNATIONALIZER_FLUENT_CORPUS to a directory of .ftl resources")
	}
	format := &FluentFormat{}
	files := 0
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".ftl" {
			return nil
		}
		clean := filepath.ToSlash(path)
		if strings.Contains(clean, "/test/") || strings.Contains(clean, "/tests/") {
			return nil
		}
		files++
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		units, err := format.ParseUnits(data)
		if err != nil {
			return &fluentCorpusError{path: path, err: err}
		}
		output, err := format.SerializeUnits(units, data)
		if err != nil {
			return &fluentCorpusError{path: path, err: err}
		}
		if !reflect.DeepEqual(output, data) {
			return &fluentCorpusError{path: path, err: errFluentRoundTripDrift}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if files == 0 {
		t.Fatalf("no .ftl files found under %s", root)
	}
}

type fluentCorpusError struct {
	path string
	err  error
}

func (err *fluentCorpusError) Error() string { return err.path + ": " + err.err.Error() }

var errFluentRoundTripDrift = &fluentRoundTripError{}

type fluentRoundTripError struct{}

func (*fluentRoundTripError) Error() string { return "unchanged serialization drifted" }
