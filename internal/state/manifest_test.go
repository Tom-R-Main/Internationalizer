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
		Bundle:       "app",
		Key:          "common.save",
		Locale:       "fr",
		SourceHash:   SourceHash("json", "Save"),
		PolicyHash:   mustHashValue(t, "policy"),
		TargetHash:   TargetHash("Enregistrer"),
		Origin:       "provider",
		Provider:     "openai",
		Model:        "model",
		ReviewStatus: ReviewNeedsReview,
		UpdatedAt:    time.Now().UTC().Truncate(time.Second),
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

func TestLoadMigratesLegacyLocaleIdentityInMemory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	data := []byte(`{"schema_version":1,"translations":{"legacy-raw-key":{"bundle":"app","key":"save","locale":"pt-br","source_hash":"source","policy_hash":"policy","target_hash":"target","updated_at":"2026-01-01T00:00:00Z"}}}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	manifest, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := manifest.Get("app", "save", "pt-BR")
	if !ok || entry.Locale != "pt-BR" {
		t.Fatalf("legacy entry = %#v, %v; want canonicalized lookup", entry, ok)
	}
	if manifest.SchemaVersion != SchemaVersion || entry.ReviewStatus != ReviewNeedsReview || entry.ReviewedAt != nil {
		t.Fatalf("legacy review state = %#v, schema = %d; want needs_review schema v2", entry, manifest.SchemaVersion)
	}
}

func TestApproveRecordsExplicitReviewSeparatelyFromOrigin(t *testing.T) {
	manifest := New()
	manifest.Set(Entry{Bundle: "app", Key: "save", Locale: "fr", Origin: "provider"})
	reviewedAt := time.Date(2026, time.September, 4, 12, 0, 0, 0, time.FixedZone("EDT", -4*60*60))

	approved, err := manifest.Approve("app", "save", "fr", reviewedAt)
	if err != nil {
		t.Fatal(err)
	}
	if approved.Origin != "provider" || approved.ReviewStatus != ReviewApproved || approved.ReviewedAt == nil || !approved.ReviewedAt.Equal(reviewedAt.UTC()) {
		t.Fatalf("approved entry = %#v", approved)
	}
}

func TestLoadRejectsInvalidApprovedState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	data := []byte(`{"schema_version":2,"translations":{"entry":{"bundle":"app","key":"save","locale":"fr","review_status":"approved","updated_at":"2026-01-01T00:00:00Z"}}}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load accepted approved state without reviewed_at")
	}
}

func TestLoadRejectsConflictingCanonicalLocaleIdentities(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	data := []byte(`{"schema_version":1,"translations":{"first":{"bundle":"app","key":"save","locale":"pt-br","source_hash":"one","updated_at":"2026-01-01T00:00:00Z"},"second":{"bundle":"app","key":"save","locale":"pt-BR","source_hash":"two","updated_at":"2026-01-01T00:00:00Z"}}}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load accepted conflicting canonical-equivalent identities")
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

func TestSourceUnitHashTracksTranslatorContextWithoutChangingFlatIdentity(t *testing.T) {
	legacy := SourceHash("json", "Save")
	if got := SourceUnitHash("json", "Save", "", ""); got != legacy {
		t.Fatalf("flat unit hash = %q, want legacy %q", got, legacy)
	}
	first := SourceUnitHash("fluent", "Open", "Verb used on a button.", "fluent-pattern-v1")
	second := SourceUnitHash("fluent", "Open", "Noun shown in a menu.", "fluent-pattern-v1")
	if first == second {
		t.Fatal("translator context did not affect source unit hash")
	}
}

func TestManifestIdentityCanonicalizesLocale(t *testing.T) {
	manifest := New()
	entry := Entry{Bundle: "app", Key: "save", Locale: "pt-br"}
	manifest.Set(entry)

	got, ok := manifest.Get("app", "save", "pt-BR")
	if !ok {
		t.Fatal("canonical-equivalent locale did not find the manifest entry")
	}
	if got.Locale != "pt-BR" {
		t.Fatalf("stored locale = %q, want %q", got.Locale, "pt-BR")
	}
	if Identity("app", "save", "pt-br") != Identity("app", "save", "pt-BR") {
		t.Fatal("canonical-equivalent locales produced different identities")
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
