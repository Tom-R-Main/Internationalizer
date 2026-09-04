package acceptance_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

func jsonIntegrityFixture(t *testing.T, providerURL string) string {
	t.Helper()
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, ".internationalizer.yml"), fmt.Sprintf(`source_locale: en
target_locales: [fr]
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    message_syntax: i18next
llm:
  provider: openai
  model: gpt-5.6-luna
  api_key_env: JSON_INTEGRITY_TEST_KEY
  base_url: %s/api.openai.com
manifest_path: .internationalizer.lock
tm_path: .internationalizer/tm.jsonl
`, providerURL))
	mustWriteFile(t, filepath.Join(root, "locales/en.json"), `{"hello":"Hello {{name}}"}`)
	mustWriteFile(t, filepath.Join(root, "locales/fr.json"), `{"hello":"Bonjour {{name}}"}`)
	return root
}

func requireJSONIntegrityFailure(t *testing.T, result cliResult, code string) {
	t.Helper()
	if result.exitCode == 0 {
		t.Fatalf("ambiguous JSON passed: stdout=%s stderr=%s", result.stdout, result.stderr)
	}
	decodeOnboardingJSON(t, result)
	if !strings.Contains(result.stdout, `"`+code+`"`) {
		t.Fatalf("missing integrity code %q: stdout=%s stderr=%s", code, result.stdout, result.stderr)
	}
}

func TestJSONIntegrityCollisionAlwaysFailsValidation(t *testing.T) {
	root := jsonIntegrityFixture(t, "http://127.0.0.1:1")
	mustWriteFile(t, filepath.Join(root, "locales/en.json"), `{"a.b":"Hello {{name}}","a":{"b":"Hello"}}`)
	mustWriteFile(t, filepath.Join(root, "locales/fr.json"), `{"a.b":"Bonjour"}`)
	for _, strict := range []bool{false, true} {
		t.Run(fmt.Sprintf("strict=%t", strict), func(t *testing.T) {
			args := []string{"validate", "--json"}
			if strict {
				args = append(args, "--strict")
			}
			var first string
			for run := range 100 {
				result := runCLI(t, root, nil, args...)
				requireJSONIntegrityFailure(t, result, "json_flattened_key_collision")
				if run == 0 {
					first = result.stdout
				} else if result.stdout != first {
					t.Fatalf("diagnostic changed on run %d:\nfirst=%s\nnow=%s", run+1, first, result.stdout)
				}
			}
		})
	}
}

func TestJSONIntegrityValidationRejectsSourceAndTargetAmbiguity(t *testing.T) {
	cases := []struct{ name, content, code string }{
		{"duplicate", `{"hello":"Hello {{name}}","hello":"Hello"}`, "json_duplicate_member"},
		{"equal-duplicate", `{"hello":"Hello","hello":"Hello"}`, "json_duplicate_member"},
		{"escaped-duplicate", `{"hello":"Hello {{name}}","\u0068ello":"Hello"}`, "json_duplicate_member"},
		{"nested-duplicate", `{"section":{"hello":"Hello {{name}}","hello":"Hello"}}`, "json_duplicate_member"},
		{"nonstring-duplicate", `{"hello":"Hello {{name}}","hello":false}`, "json_duplicate_member"},
		{"flattened", `{"a.b":"Hello {{name}}","a":{"b":"Hello"}}`, "json_flattened_key_collision"},
		{"array-alias", `{"a.0":"Hello {{name}}","a":["Hello"]}`, "json_flattened_key_collision"},
		{"nonstring-alias", `{"a.b":"Hello {{name}}","a":{"b":false}}`, "json_flattened_key_collision"},
	}
	for _, tc := range cases {
		for _, locale := range []string{"en", "fr"} {
			t.Run(tc.name+"/"+locale, func(t *testing.T) {
				root := jsonIntegrityFixture(t, "http://127.0.0.1:1")
				mustWriteFile(t, filepath.Join(root, "locales", locale+".json"), tc.content)
				for _, strict := range []bool{false, true} {
					args := []string{"validate", "--json"}
					if strict {
						args = append(args, "--strict")
					}
					requireJSONIntegrityFailure(t, runCLI(t, root, nil, args...), tc.code)
				}
				filtered := runCLI(t, root, nil, "validate", "--json", "--limit", "1", "--finding-code", "missing")
				requireJSONIntegrityFailure(t, filtered, tc.code)
			})
		}
	}
}

func TestJSONIntegrityDiscoveryRetainsMalformedCatalog(t *testing.T) {
	for _, configured := range []bool{false, true} {
		t.Run(fmt.Sprintf("configured=%t", configured), func(t *testing.T) {
			root := jsonIntegrityFixture(t, "http://127.0.0.1:1")
			catalog := "web/locales/en.json"
			if configured {
				catalog = "locales/en.json"
			}
			mustWriteFile(t, filepath.Join(root, catalog), `{"hello":"Hello {{name}}","hello":"Hello"}`)
			before := jsonIntegritySnapshot(t, root)
			for _, args := range [][]string{{"detect", "--json"}, {"config", "check", "--json"}} {
				result := runCLI(t, root, nil, args...)
				decodeOnboardingJSON(t, result)
				if args[0] == "config" && result.exitCode == 0 {
					t.Fatalf("config check passed malformed catalog: %s", result.stdout)
				}
				if !strings.Contains(result.stdout, "JSON_DUPLICATE_MEMBER") {
					t.Fatalf("malformed catalog silently disappeared: %s", result.stdout)
				}
				var decoded struct {
					Data struct {
						Inspection struct {
							Candidates []struct {
								Source         string `json:"source"`
								ParseErrorCode string `json:"parse_error_code"`
							} `json:"candidates"`
						} `json:"inspection"`
					} `json:"data"`
				}
				if err := json.Unmarshal([]byte(result.stdout), &decoded); err != nil {
					t.Fatal(err)
				}
				found := false
				for _, candidate := range decoded.Data.Inspection.Candidates {
					if filepath.ToSlash(candidate.Source) == catalog && candidate.ParseErrorCode == "json_duplicate_member" {
						found = true
					}
				}
				if !found {
					t.Fatalf("malformed candidate %s was not retained with its parse code: %s", catalog, result.stdout)
				}
			}
			if after := jsonIntegritySnapshot(t, root); !reflect.DeepEqual(before, after) {
				t.Fatal("discovery changed project files")
			}
		})
	}
}

func TestJSONIntegrityMutatingCommandsRejectWithoutSideEffects(t *testing.T) {
	for _, command := range []struct {
		name string
		args []string
	}{
		{"translate", []string{"translate", "--json"}},
		{"dry-run", []string{"translate", "--dry-run", "--json"}},
		{"adopt", []string{"translate", "--adopt-existing", "--json"}},
		{"pseudo", []string{"pseudo", "--locale", "en-XA", "--force"}},
		{"approve", []string{"review", "approve", "--locale", "fr", "--all"}},
	} {
		for _, input := range []struct{ name, content, diagnostic string }{
			{"duplicate", `{"hello":"Hello {{name}}","hello":"Hello"}`, "duplicate"},
			{"collision", `{"a.b":"Hello {{name}}","a":{"b":"Hello"}}`, "collid"},
		} {
			for _, source := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%s/source=%t", command.name, input.name, source), func(t *testing.T) {
					var calls atomic.Int32
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						calls.Add(1)
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"output":[{"type":"message","content":[{"type":"output_text","text":"{\"hello\":\"Bonjour\"}"}]}]}`))
					}))
					defer server.Close()
					root := jsonIntegrityFixture(t, server.URL)
					if command.name == "approve" {
						runCLI(t, root, nil, "translate", "--adopt-existing").requireSuccess(t)
					}
					locale := "fr"
					if command.name == "pseudo" {
						locale = "en-XA"
					}
					if source {
						locale = "en"
					}
					mustWriteFile(t, filepath.Join(root, "locales", locale+".json"), input.content)
					before := jsonIntegritySnapshot(t, root)
					result := runCLI(t, root, []string{"JSON_INTEGRITY_TEST_KEY=synthetic-only"}, command.args...)
					if result.exitCode == 0 {
						t.Errorf("command accepted malformed catalog: %s", result.stdout)
					}
					if !strings.Contains(strings.ToLower(result.stdout+result.stderr), input.diagnostic) {
						t.Errorf("failure did not explain %s: stdout=%s stderr=%s", input.diagnostic, result.stdout, result.stderr)
					}
					if command.args[0] == "translate" {
						var output struct {
							Errors []struct {
								Code string `json:"code"`
							} `json:"errors"`
							Data struct {
								Jobs []struct {
									InputError *struct {
										Code      string `json:"code"`
										Path      string `json:"path"`
										OtherPath string `json:"other_path"`
									} `json:"input_error"`
								} `json:"jobs"`
							} `json:"data"`
						}
						if err := json.Unmarshal([]byte(result.stdout), &output); err != nil {
							t.Fatal(err)
						}
						code := "json_duplicate_member"
						if input.name == "collision" {
							code = "json_flattened_key_collision"
						}
						if source {
							if len(output.Errors) != 1 || output.Errors[0].Code != code {
								t.Errorf("source integrity code was hidden by generic translation error: %s", result.stdout)
							}
						} else if len(output.Data.Jobs) != 1 || output.Data.Jobs[0].InputError == nil {
							t.Errorf("target job omitted structured input_error: %s", result.stdout)
						} else if detail := output.Data.Jobs[0].InputError; detail.Code != code || detail.Path == "" || detail.OtherPath == "" {
							t.Errorf("target input_error omitted integrity code or original paths: %s", result.stdout)
						}
					}
					if input.name == "collision" {
						for _, pointer := range []string{"/a.b", "/a/b"} {
							if !strings.Contains(result.stdout+result.stderr, pointer) {
								t.Errorf("collision diagnostic omitted original path %q: stdout=%s stderr=%s", pointer, result.stdout, result.stderr)
							}
						}
					}
					if calls.Load() != 0 {
						t.Errorf("malformed catalog triggered %d provider requests", calls.Load())
					}
					if after := jsonIntegritySnapshot(t, root); !reflect.DeepEqual(before, after) {
						t.Errorf("command changed catalog, config, or state after malformed input: before=%v after=%v", before, after)
					}
				})
			}
		}
	}
}

func jsonIntegritySnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files[relative] = string(mustReadFile(t, path))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return files
}
