package onboarding

import (
	"strings"
	"testing"
)

func TestScanRetainsMalformedCatalogEvidence(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr]\nsource_path: configured/en.json\nmessage_syntax: plain\n")
	discoveryFile(t, root, "configured/en.json", `{"hello":"Hello"}`)
	discoveryFile(t, root, "app/locales/en.json", `{"hello":`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range inspection.Diagnostics {
		if d.Code == "CATALOG_PARSE_FAILED" && d.Severity == "error" && strings.Contains(d.Message, "app/locales/en.json") {
			found = true
		}
	}
	if !found {
		t.Fatalf("malformed catalog silently disappeared: %+v", inspection)
	}
}

func TestScanReportsJSONIntegrityForConfiguredTargets(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr]\nsource_path: locales/en.json\nmessage_syntax: plain\n")
	discoveryFile(t, root, "locales/en.json", `{"hello":"Hello"}`)
	discoveryFile(t, root, "locales/fr.json", `{"hello":"PRIVATE FIRST VALUE","hello":"PRIVATE LAST VALUE"}`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range inspection.Diagnostics {
		if strings.Contains(d.Message, "PRIVATE") {
			t.Fatal("catalog value leaked")
		}
		if d.Code == "JSON_DUPLICATE_MEMBER" && d.Bundle == "default" && strings.Contains(d.Message, "locales/fr.json") {
			found = true
		}
	}
	if !found {
		t.Fatalf("target integrity omitted: %+v", inspection.Diagnostics)
	}
}
