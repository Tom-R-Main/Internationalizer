package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"gopkg.in/yaml.v3"
)

func persistenceCLIProject(t *testing.T, source, target, endpoint string, locales []string) (string, *config.Config) {
	t.Helper()
	dir := t.TempDir()
	cfg := &config.Config{
		SourceLocale: "en", TargetLocales: locales, MessageSyntax: message.Plain,
		SourcePath: filepath.Join(dir, "en.json"), BatchSize: 1, Concurrency: 1,
		TMPath: filepath.Join(dir, "tm.jsonl"), ManifestPath: filepath.Join(dir, "manifest.json"),
		StyleGuidesDir: filepath.Join(dir, "guides"), GlossaryDir: filepath.Join(dir, "glossary"),
		LLM: config.LLM{Provider: "openai", Model: "test-model", BaseURL: endpoint, APIKeyEnv: "PERSISTENCE_TEST_KEY"},
	}
	t.Setenv("PERSISTENCE_TEST_KEY", "synthetic-test-credential")
	for name, content := range map[string]string{"en.json": source, "fr.json": target} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path, cfg
}

func persistenceResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"a\":\"Un\"}"}}]}`))
}

func TestTranslationJSONLaterBatchFailureDoesNotClaimRetainedWork(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) == 1 {
			persistenceResponse(w)
			return
		}
		http.Error(w, "synthetic batch failure", http.StatusBadRequest)
	}))
	defer server.Close()
	path, cfg := persistenceCLIProject(t, `{"a":"A","b":"B"}`, `{}`, server.URL, []string{"fr"})
	result, err := runJSONCommand(t, "translate", "--config", path, "--json")
	if err == nil || result.Status != "blocked" {
		t.Fatalf("unpersisted batches must be blocked, not partial success: status=%s error=%v", result.Status, err)
	}
	data, _ := json.Marshal(result.Data)
	var run translationJSON
	if err := json.Unmarshal(data, &run); err != nil {
		t.Fatal(err)
	}
	if run.Summary.GeneratedKeys != 1 || !run.ProviderCalled || run.Summary.PersistedJobs != 0 {
		t.Fatalf("generation must remain distinct from persistence: %+v", run)
	}
	if len(run.Jobs) != 1 || run.Jobs[0].CatalogWritten || run.Jobs[0].ManifestUpdated {
		t.Fatalf("uncommitted batch claimed persistence: %+v", run.Jobs)
	}
	target, err := os.ReadFile(filepath.Join(filepath.Dir(path), "fr.json"))
	if err != nil || string(target) != `{}` {
		t.Fatalf("failed locale wrote its staged batch: %s %v", target, err)
	}
	if _, err := os.Stat(cfg.ManifestPath); !os.IsNotExist(err) {
		t.Fatalf("unexpected manifest after unpersisted failure: %v", err)
	}
}

func TestTranslationJSONRetainedCachedAndAdoptedJobsArePartialFailures(t *testing.T) {
	for _, mode := range []string{"cached", "adopted"} {
		t.Run(mode, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				persistenceResponse(w)
			}))
			defer server.Close()
			path, cfg := persistenceCLIProject(t, `{"a":"A"}`, `{"a":"Un"}`, server.URL, []string{"fr", "de"})
			args := []string{"translate", "--config", path, "--json"}
			if mode == "cached" {
				if err := os.WriteFile(filepath.Join(filepath.Dir(path), "fr.json"), []byte(`{}`), 0o600); err != nil {
					t.Fatal(err)
				}
				if _, err := runJSONCommand(t, "translate", "--config", path, "--locale", "fr", "--json"); err != nil {
					t.Fatal(err)
				}
				for _, stale := range []string{filepath.Join(filepath.Dir(path), "fr.json"), cfg.ManifestPath} {
					if err := os.Remove(stale); err != nil {
						t.Fatal(err)
					}
				}
				calls.Store(0)
			} else {
				args = append(args, "--adopt-existing")
			}
			if err := os.WriteFile(filepath.Join(filepath.Dir(path), "de.json"), []byte(`invalid target`), 0o600); err != nil {
				t.Fatal(err)
			}
			result, err := runJSONCommand(t, args...)
			if err == nil || result.Status != "partial_failure" {
				t.Fatalf("retained %s job must report partial_failure: status=%s error=%v", mode, result.Status, err)
			}
			data, _ := json.Marshal(result.Data)
			var run translationJSON
			if err := json.Unmarshal(data, &run); err != nil {
				t.Fatal(err)
			}
			if run.Summary.PersistedJobs != 1 || run.Summary.GeneratedKeys != 0 || run.ProviderCalled || len(run.Jobs) != 2 {
				t.Fatalf("incorrect retained-work summary: %+v", run)
			}
			if !run.Jobs[0].ManifestUpdated || run.Jobs[0].CatalogWritten != (mode == "cached") || run.Jobs[1].ManifestUpdated || run.Jobs[1].CatalogWritten {
				t.Fatalf("persistence was attributed to the wrong jobs: %+v", run.Jobs)
			}
			if calls.Load() != 0 {
				t.Fatalf("%s execution unexpectedly called provider", mode)
			}
			if _, err := os.Stat(cfg.ManifestPath); err != nil {
				t.Fatalf("successful %s job was not retained: %v", mode, err)
			}
		})
	}
}
