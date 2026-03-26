package glossary

import (
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
