package review

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/policy"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

func TestApproveVerifiesCurrentArtifactAndPreservesOrigin(t *testing.T) {
	cfg := reviewFixture(t, "Enregistrer")
	reviewedAt := time.Date(2026, time.September, 4, 16, 0, 0, 0, time.UTC)

	approved, err := Approve(cfg, Filter{Locale: "fr", Bundle: "default", Keys: []string{"save"}}, reviewedAt)
	if err != nil {
		t.Fatal(err)
	}
	if len(approved) != 1 || approved[0].Origin != "provider" || approved[0].ReviewStatus != state.ReviewApproved {
		t.Fatalf("approved = %#v", approved)
	}

	loaded, err := state.Load(cfg.ManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := loaded.Get("default", "save", "fr")
	if !ok || entry.ReviewedAt == nil || !entry.ReviewedAt.Equal(reviewedAt) {
		t.Fatalf("saved approval = %#v, found = %v", entry, ok)
	}
}

func TestApproveRejectsModifiedTarget(t *testing.T) {
	cfg := reviewFixture(t, "Enregistrer")
	targetPath := filepath.Join(filepath.Dir(cfg.SourcePath), "fr.json")
	if err := os.WriteFile(targetPath, []byte(`{"save":"Sauvegarder"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Approve(cfg, Filter{Locale: "fr", All: true}, time.Now())
	if err == nil || !strings.Contains(err.Error(), "target changed") {
		t.Fatalf("Approve error = %v, want target-changed rejection", err)
	}
}

func TestListFiltersAndSortsReviewState(t *testing.T) {
	manifest := state.New()
	manifest.Set(state.Entry{Bundle: "b", Key: "z", Locale: "fr", ReviewStatus: state.ReviewNeedsReview})
	manifest.Set(state.Entry{Bundle: "a", Key: "a", Locale: "fr", ReviewStatus: state.ReviewNeedsReview})
	manifest.Set(state.Entry{Bundle: "a", Key: "a", Locale: "de", ReviewStatus: state.ReviewNeedsReview})

	entries, err := List(manifest, Filter{Locale: "fr", Status: state.ReviewNeedsReview})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || entries[0].Bundle != "a" || entries[1].Bundle != "b" {
		t.Fatalf("entries = %#v", entries)
	}
}

func reviewFixture(t *testing.T, target string) *config.Config {
	t.Helper()
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"save":"Save"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"save":"`+target+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		SourceLocale:  "en",
		TargetLocales: []string{"fr"},
		SourcePath:    sourcePath,
		ManifestPath:  filepath.Join(dir, "manifest.json"),
		LLM:           config.LLM{Provider: "gemini", Model: "test-model"},
	}
	cfg.ApplyDefaults()
	resolved, err := policy.Resolve(cfg, "fr", "json", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	manifest := state.New()
	manifest.Set(state.Entry{
		Bundle:       "default",
		Key:          "save",
		Locale:       "fr",
		SourceHash:   state.SourceHash("json", "Save"),
		PolicyHash:   resolved.Hash,
		TargetHash:   state.TargetHash(target),
		Origin:       "provider",
		ReviewStatus: state.ReviewNeedsReview,
		UpdatedAt:    time.Now().UTC(),
	})
	if err := manifest.Save(cfg.ManifestPath); err != nil {
		t.Fatal(err)
	}
	return cfg
}
