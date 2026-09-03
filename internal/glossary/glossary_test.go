package glossary

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func mustAdd(t *testing.T, dir, locale, source, target string) {
	t.Helper()
	if err := Add(dir, locale, source, target); err != nil {
		t.Fatalf("Add(%q, %q, %q): %v", locale, source, target, err)
	}
}

func TestAddAndLoad(t *testing.T) {
	dir := t.TempDir()

	mustAdd(t, dir, "fr", "Dashboard", "Tableau de bord")
	mustAdd(t, dir, "fr", "Settings", "Paramètres")

	terms, err := Load(dir, "fr")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(terms) != 2 {
		t.Fatalf("got %d terms, want 2", len(terms))
	}
	if terms[0].Source != "Dashboard" || terms[0].Target != "Tableau de bord" {
		t.Errorf("term 0: got %+v", terms[0])
	}
}

func TestAddUpdatesExisting(t *testing.T) {
	dir := t.TempDir()

	mustAdd(t, dir, "fr", "Dashboard", "Tableau de bord")
	mustAdd(t, dir, "fr", "Dashboard", "Panneau de contrôle")

	terms, err := Load(dir, "fr")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(terms) != 1 {
		t.Fatalf("got %d terms, want 1 (should update, not duplicate)", len(terms))
	}
	if terms[0].Target != "Panneau de contrôle" {
		t.Errorf("target not updated: got %q", terms[0].Target)
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()

	mustAdd(t, dir, "fr", "Dashboard", "Tableau de bord")
	mustAdd(t, dir, "fr", "Settings", "Paramètres")

	if err := Remove(dir, "fr", "Dashboard"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	terms, err := Load(dir, "fr")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(terms) != 1 {
		t.Fatalf("got %d terms, want 1", len(terms))
	}
	if terms[0].Source != "Settings" {
		t.Errorf("wrong remaining term: %q", terms[0].Source)
	}
}

func TestRemoveNotFound(t *testing.T) {
	dir := t.TempDir()
	mustAdd(t, dir, "fr", "Dashboard", "Tableau de bord")

	err := Remove(dir, "fr", "Nonexistent")
	if err == nil {
		t.Error("expected error for removing nonexistent term")
	}
}

func TestLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	terms, err := Load(dir, "fr")
	if err != nil {
		t.Fatalf("Load empty: %v", err)
	}
	if terms != nil {
		t.Errorf("expected nil for missing glossary, got %v", terms)
	}
}

func TestLoadOldGlossaryWithoutValidationMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(path, []byte(`[{"source":"Dashboard","target":"Tableau de bord","ignore_case":true}]`), 0o644); err != nil {
		t.Fatal(err)
	}

	terms, err := Load(dir, "fr")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(terms) != 1 {
		t.Fatalf("Load returned %d terms, want 1", len(terms))
	}
	if terms[0].Enforcement != "" || terms[0].Variants != nil {
		t.Fatalf("old glossary metadata = %#v, want zero values", terms[0])
	}
}

func TestSaveAndLoadVariantsAndEnforcement(t *testing.T) {
	dir := t.TempDir()
	want := []Term{{
		Source:      "Sign in",
		Target:      "Connexion",
		Variants:    []string{"Se connecter"},
		Enforcement: EnforcementWarning,
	}}
	if err := Save(dir, "fr", want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load(dir, "fr")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Load = %#v, want %#v", got, want)
	}
}

func TestLoadRejectsInvalidEnforcement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(path, []byte(`[{"source":"Save","target":"Enregistrer","enforcement":"block"}]`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir, "fr")
	if err == nil {
		t.Fatal("Load accepted invalid enforcement")
	}
	if !strings.Contains(err.Error(), `invalid enforcement "block"`) {
		t.Fatalf("Load error = %q, want invalid enforcement detail", err)
	}
}

func TestSaveRejectsInvalidEnforcementBeforeWriting(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "glossary")
	err := Save(dir, "fr", []Term{{Source: "Save", Target: "Enregistrer", Enforcement: "block"}})
	if err == nil {
		t.Fatal("Save accepted invalid enforcement")
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatalf("Save created glossary directory for invalid terms: %v", statErr)
	}
}

func TestApprovedTargetsIncludesPrimaryAndVariants(t *testing.T) {
	term := Term{Target: "Connexion", Variants: []string{"Se connecter", "Identification"}}
	want := []string{"Connexion", "Se connecter", "Identification"}
	if got := ApprovedTargets(term); !reflect.DeepEqual(got, want) {
		t.Fatalf("ApprovedTargets = %#v, want %#v", got, want)
	}
}

func TestSourceIdenticalExemptRequiresCompleteExplicitMatch(t *testing.T) {
	terms := []Term{
		{Source: "API", Target: "API"},
		{Source: "GitHub", Target: "GitHub", IgnoreCase: true},
	}
	tests := []struct {
		name   string
		source string
		target string
		want   bool
	}{
		{name: "exact", source: "API", target: "API", want: true},
		{name: "ignore case", source: "github", target: "GITHUB", want: true},
		{name: "partial", source: "API access", target: "API access", want: false},
		{name: "not source identical", source: "API", target: "Apis", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := SourceIdenticalExempt(terms, test.source, test.target); got != test.want {
				t.Fatalf("SourceIdenticalExempt(%q, %q) = %t, want %t", test.source, test.target, got, test.want)
			}
		})
	}
}

func TestFormatForPrompt(t *testing.T) {
	terms := []Term{
		{Source: "Dashboard", Target: "Tableau de bord", WholeWord: true},
		{Source: "Save", Target: "Enregistrer", IgnoreCase: true},
	}

	output := FormatForPrompt(terms)
	if output == "" {
		t.Fatal("empty output")
	}
	if !containsStr(output, "Source") || !containsStr(output, "Translation") {
		t.Error("missing table headers")
	}
	if !containsStr(output, "Tableau de bord") {
		t.Error("missing term translation")
	}
}

func TestFormatForPromptEmpty(t *testing.T) {
	if output := FormatForPrompt(nil); output != "" {
		t.Errorf("expected empty for nil terms, got %q", output)
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
