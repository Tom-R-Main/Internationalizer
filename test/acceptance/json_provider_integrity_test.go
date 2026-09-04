package acceptance_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/jsonintegrity"
)

func TestJSONIntegrityProviderResponseFailsWithoutPersistence(t *testing.T) {
	for _, tc := range []struct {
		name, response, code, path, otherPath string
	}{
		{"duplicate", `{"a.b":"PRIVATE FIRST VALUE","a.b":"PRIVATE LAST VALUE"}`, "json_duplicate_member", "/a.b", "/a.b"},
		{"collision", `{"a.b":"PRIVATE FIRST VALUE","a":{"b":"PRIVATE LAST VALUE"}}`, "json_flattened_key_collision", "/a/b", "/a.b"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.URL.Path != "/api.openai.com/v1/responses" {
					t.Errorf("unexpected provider path: %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{
					"output": []any{map[string]any{
						"type":    "message",
						"content": []any{map[string]any{"type": "output_text", "text": tc.response}},
					}},
				}); err != nil {
					t.Errorf("writing mock response: %v", err)
				}
			}))
			defer server.Close()
			root := jsonIntegrityFixture(t, server.URL)
			// Seed real provenance without a provider, then leave a single new
			// source key pending. Existing target, manifest, and TM must survive.
			runCLI(t, root, nil, "translate", "--adopt-existing").requireSuccess(t)
			mustReadFile(t, filepath.Join(root, ".internationalizer.lock"))
			mustWriteFile(t, filepath.Join(root, ".internationalizer", "tm.jsonl"), "{\"bundle\":\"unrelated\",\"key\":\"retained\",\"source\":\"Seed\",\"target\":\"Graine\",\"locale\":\"fr\",\"hash\":\"existing\",\"policy_hash\":\"existing\"}\n")
			mustWriteFile(t, filepath.Join(root, "locales", "en.json"), `{"hello":"Hello {{name}}","a.b":"New message"}`)
			before := jsonIntegritySnapshot(t, root)
			result := runCLI(t, root, []string{"JSON_INTEGRITY_TEST_KEY=synthetic-only"}, "translate", "--json")
			requireJSONIntegrityFailure(t, result, tc.code)
			var output struct {
				Status string `json:"status"`
				Data   struct {
					ProviderCalled bool `json:"provider_called"`
					Summary        struct {
						BlockedJobs   int `json:"blocked_jobs"`
						GeneratedKeys int `json:"generated_keys"`
						PersistedJobs int `json:"persisted_jobs"`
					} `json:"summary"`
					Jobs []struct {
						InputError      *jsonintegrity.Error `json:"input_error"`
						ProviderCalls   int                  `json:"provider_calls"`
						CatalogWritten  bool                 `json:"catalog_written"`
						ManifestUpdated bool                 `json:"manifest_updated"`
					} `json:"jobs"`
				} `json:"data"`
			}
			if err := json.Unmarshal([]byte(result.stdout), &output); err != nil {
				t.Fatal(err)
			}
			if calls.Load() != 1 || !output.Data.ProviderCalled || len(output.Data.Jobs) != 1 {
				t.Fatalf("expected exactly one attempted provider job, calls=%d: %s", calls.Load(), result.stdout)
			}
			job := output.Data.Jobs[0]
			if job.InputError == nil || job.InputError.Code != tc.code || job.InputError.Key != "a.b" || job.InputError.Path != tc.path || job.InputError.OtherPath != tc.otherPath {
				t.Fatalf("provider integrity details missing: %s", result.stdout)
			}
			if job.ProviderCalls != 1 || job.CatalogWritten || job.ManifestUpdated || output.Status != "blocked" || output.Data.Summary.BlockedJobs != 1 || output.Data.Summary.GeneratedKeys != 0 || output.Data.Summary.PersistedJobs != 0 {
				t.Fatalf("failure reported generation or persistence: %s", result.stdout)
			}
			if strings.Contains(result.stdout+result.stderr, "PRIVATE") {
				t.Fatal("provider values leaked through integrity diagnostics")
			}
			if after := jsonIntegritySnapshot(t, root); !reflect.DeepEqual(before, after) {
				t.Fatal("failed provider response changed target, manifest, TM, or other project files")
			}
		})
	}
}
