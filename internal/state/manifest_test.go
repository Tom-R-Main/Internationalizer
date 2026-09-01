package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestManifestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "manifest.json")
	manifest := New()
	entry := Entry{
		Bundle:     "app",
		Key:        "common.save",
		Locale:     "fr",
		SourceHash: SourceHash("json", "Save"),
		PolicyHash: mustHashValue(t, "policy"),
		TargetHash: TargetHash("Enregistrer"),
		Origin:     "provider",
		Provider:   "openai",
		Model:      "model",
		UpdatedAt:  time.Now().UTC().Truncate(time.Second),
	}
	manifest.Set(entry)
	if err := manifest.Save(path); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := loaded.Get("app", "common.save", "fr")
	if !ok {
		t.Fatal("saved entry was not loaded")
	}
	if got != entry {
		t.Fatalf("loaded entry = %#v, want %#v", got, entry)
	}
}

func TestLoadRejectsMalformedManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load accepted malformed provenance")
	}
}

func TestHashesSeparateSourceFormatAndPolicy(t *testing.T) {
	if SourceHash("json", "Save") == SourceHash("markdown", "Save") {
		t.Fatal("source hash ignored format")
	}
	if mustHashValue(t, struct{ Prompt string }{"one"}) == mustHashValue(t, struct{ Prompt string }{"two"}) {
		t.Fatal("policy hash ignored prompt")
	}
}

func TestHashValueRejectsUnsupportedValue(t *testing.T) {
	if _, err := HashValue(func() {}); err == nil {
		t.Fatal("HashValue accepted an unsupported value")
	}
}

func mustHashValue(t *testing.T, value interface{}) string {
	t.Helper()
	hash, err := HashValue(value)
	if err != nil {
		t.Fatal(err)
	}
	return hash
}
