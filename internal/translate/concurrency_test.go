package translate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

type outOfOrderProvider struct {
	frStarted chan struct{}
	releaseFR chan struct{}
}

func (p *outOfOrderProvider) Name() string { return "ordered-test" }

func (p *outOfOrderProvider) Translate(ctx context.Context, req llm.TranslateRequest) (*llm.TranslateResponse, error) {
	switch req.TargetLocale {
	case "fr":
		close(p.frStarted)
		select {
		case <-p.releaseFR:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	case "de":
		select {
		case <-p.frStarted:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		close(p.releaseFR)
	}

	translations := make(map[string]string, len(req.Entries))
	for _, entry := range req.Entries {
		translations[entry.Key] = req.TargetLocale + ":" + entry.Value
	}
	return &llm.TranslateResponse{Translations: translations}, nil
}

func TestRunKeepsConfiguredResultOrderWhenConcurrentJobsFinishOutOfOrder(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"Save","b":"Cancel"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	cfg.TargetLocales = []string{"fr", "de"}
	cfg.Concurrency = 2
	provider := &outOfOrderProvider{
		frStarted: make(chan struct{}),
		releaseFR: make(chan struct{}),
	}

	results, err := Run(t.Context(), cfg, provider, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[0].Locale != "fr" || results[1].Locale != "de" {
		t.Fatalf("result order = %#v, want fr then de", results)
	}
	for _, locale := range cfg.TargetLocales {
		data, err := os.ReadFile(filepath.Join(dir, locale+".json"))
		if err != nil {
			t.Fatal(err)
		}
		if want := `"a": "` + locale + `:Save"`; !bytes.Contains(data, []byte(want)) {
			t.Fatalf("%s target = %s, want %s", locale, data, want)
		}
	}
	manifest, err := state.Load(cfg.ManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Translations) != 4 {
		t.Fatalf("manifest entries = %d, want 4", len(manifest.Translations))
	}
}

func TestRunHonorsCancellationBeforeWritingTargets(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"Save"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(dir, sourcePath)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	provider := &fakeProvider{}
	results, err := Run(ctx, cfg, provider, Options{})
	if err != context.Canceled {
		t.Fatalf("Run error = %v, want context.Canceled", err)
	}
	if len(results) != 0 {
		t.Fatalf("canceled run results = %#v, want none", results)
	}
	if provider.calls != 0 {
		t.Fatalf("provider calls = %d, want none", provider.calls)
	}
	if _, err := os.Stat(filepath.Join(dir, "fr.json")); !os.IsNotExist(err) {
		t.Fatalf("canceled run created target: %v", err)
	}
	if _, err := os.Stat(cfg.ManifestPath); !os.IsNotExist(err) {
		t.Fatalf("canceled run created manifest: %v", err)
	}
}
