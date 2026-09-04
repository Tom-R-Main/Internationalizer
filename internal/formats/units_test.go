package formats

import (
	"reflect"
	"testing"
)

func TestParseUnitsAdaptsLegacyFormatsDeterministically(t *testing.T) {
	units, err := ParseUnits(&JSONFormat{}, []byte(`{"z":"Last","a":"First"}`))
	if err != nil {
		t.Fatal(err)
	}
	want := []Unit{
		{ID: "a", Value: "First", Kind: UnitMessage},
		{ID: "z", Value: "Last", Kind: UnitMessage},
	}
	if !reflect.DeepEqual(units, want) {
		t.Fatalf("units = %#v, want %#v", units, want)
	}
}

func TestMarkdownUsesDocumentUnit(t *testing.T) {
	units, err := ParseUnits(&MarkdownFormat{}, []byte("# Hello\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != 1 || units[0].ID != "_content" || units[0].Kind != UnitDocument {
		t.Fatalf("units = %#v", units)
	}
}

func TestMergeUnitValuesPreservesMetadataAndSourceUnits(t *testing.T) {
	baseline := []Unit{{ID: "title", Value: "Ancien", Kind: UnitAttribute, Context: "window title"}}
	source := []Unit{
		{ID: "title", Value: "Title", Kind: UnitAttribute, Context: "window title"},
		{ID: "body", Value: "Body", Kind: UnitMessage, Structure: "text"},
	}
	merged := MergeUnitValues(baseline, source, map[string]string{"title": "Nouveau", "body": "Corps"})
	if len(merged) != 2 || merged[0].Value != "Nouveau" || merged[0].Context != "window title" || merged[1].ID != "body" || merged[1].Structure != "text" {
		t.Fatalf("merged = %#v", merged)
	}
}

func TestValidateUnitsRejectsDuplicateIdentity(t *testing.T) {
	if err := ValidateUnits([]Unit{{ID: "same"}, {ID: "same"}}); err == nil {
		t.Fatal("ValidateUnits accepted duplicate identity")
	}
}
