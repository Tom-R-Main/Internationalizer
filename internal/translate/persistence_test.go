package translate

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/llm"
)

type persistenceProvider struct {
	beforeResponse func()
}

func (*persistenceProvider) Name() string { return "persistence-test" }

func (p *persistenceProvider) Translate(_ context.Context, _ llm.TranslateRequest) (*llm.TranslateResponse, error) {
	if p.beforeResponse != nil {
		p.beforeResponse()
	}
	return &llm.TranslateResponse{Translations: map[string]string{"a": "Un"}}, nil
}

func TestRunPersistenceFlagsFollowActualWriteBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name            string
		catalogWritten  bool
		manifestUpdated bool
	}{
		{name: "success", catalogWritten: true, manifestUpdated: true},
		{name: "catalog failure"},
		{name: "TM failure", catalogWritten: true, manifestUpdated: true},
		{name: "manifest failure", catalogWritten: true},
		{name: "cancellation after catalog", catalogWritten: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			source := filepath.Join(dir, "en.json")
			if err := os.WriteFile(source, []byte(`{"a":"A"}`), 0o600); err != nil {
				t.Fatal(err)
			}
			cfg := testConfig(dir, source)
			cfg.TMPath = filepath.Join(dir, "tm.jsonl")
			cfg.ManifestPath = filepath.Join(dir, "manifest.json")
			target := filepath.Join(dir, "fr.json")
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			provider := &persistenceProvider{beforeResponse: func() {
				var blocked string
				switch tc.name {
				case "catalog failure":
					blocked = target
				case "TM failure":
					blocked = cfg.TMPath
				case "manifest failure":
					blocked = cfg.ManifestPath
				case "cancellation after catalog":
					// The provider completed despite cancellation. Its result is
					// committed before Run observes cancellation and skips Save.
					cancel()
				}
				if blocked != "" {
					if err := os.Mkdir(blocked, 0o700); err != nil {
						t.Error(err)
					}
				}
			}}
			results, err := Run(ctx, cfg, provider, Options{})
			if (err != nil) != (tc.name != "success") || len(results) != 1 {
				t.Fatalf("results=%+v error=%v", results, err)
			}
			if tc.name == "cancellation after catalog" && !errors.Is(err, context.Canceled) {
				t.Fatalf("expected cancellation, got %v", err)
			}
			result := results[0]
			if result.CatalogWritten != tc.catalogWritten || result.ManifestUpdated != tc.manifestUpdated {
				t.Fatalf("persistence flags=%+v, want catalog=%t manifest=%t", result, tc.catalogWritten, tc.manifestUpdated)
			}
			if result.KeysTranslated != 1 {
				t.Fatalf("provider-generated count should not depend on persistence: %+v", result)
			}
			if tc.catalogWritten {
				if info, err := os.Stat(target); err != nil || !info.Mode().IsRegular() {
					t.Fatalf("catalog flag has no retained file: %v", err)
				}
			}
			if tc.manifestUpdated {
				if info, err := os.Stat(cfg.ManifestPath); err != nil || !info.Mode().IsRegular() {
					t.Fatalf("manifest flag has no retained file: %v", err)
				}
			}
		})
	}
}
