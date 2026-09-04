package acceptance_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

var binaryPath string

func TestMain(m *testing.M) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "locating acceptance test source")
		os.Exit(1)
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	buildDir, err := os.MkdirTemp("", "internationalizer-acceptance-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating acceptance build directory: %v\n", err)
		os.Exit(1)
	}
	binaryName := "internationalizer"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath = filepath.Join(buildDir, binaryName)
	build := exec.Command("go", "build", "-o", binaryPath, "./cmd/internationalizer")
	build.Dir = repositoryRoot
	build.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	if output, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "building acceptance binary: %v\n%s", err, output)
		os.Exit(1)
	}

	exitCode := m.Run()
	if err := os.RemoveAll(buildDir); err != nil && exitCode == 0 {
		fmt.Fprintf(os.Stderr, "removing acceptance build directory: %v\n", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}

func TestDetectUsesProjectFixture(t *testing.T) {
	projectDir := copyFixture(t, "react-i18next")

	result := runCLI(t, projectDir, nil, "detect")
	result.requireSuccess(t)
	normalizedStdout := strings.ReplaceAll(result.stdout, `\`, "/")
	for _, want := range []string{
		"Catalog public/locales/en.json: framework=i18next",
		"suggested syntax=i18next",
		"UNCOVERED_CATALOG",
	} {
		if !strings.Contains(normalizedStdout, want) {
			t.Fatalf("detect stdout does not contain %q:\n%s", want, result.stdout)
		}
	}
}

func TestMultiFormatFixtureAdoptsAndValidatesAsOneProject(t *testing.T) {
	projectDir := copyFixture(t, "multi-format")

	adoption := runCLI(t, projectDir, nil, "translate", "--adopt-existing")
	adoption.requireSuccess(t)
	manifest, err := state.Load(filepath.Join(projectDir, ".internationalizer.lock"))
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Translations) != 6 {
		t.Fatalf("multi-format manifest entries = %d, want 6", len(manifest.Translations))
	}

	validation := runCLI(t, projectDir, nil, "validate", "--json")
	validation.requireSuccess(t)
	if got := strings.Count(validation.stdout, `"coverage": 100`); got != 3 {
		t.Fatalf("full-coverage reports = %d, want 3:\n%s", got, validation.stdout)
	}

	dryRun := runCLI(t, projectDir, nil, "translate", "--dry-run")
	dryRun.requireSuccess(t)
	if !strings.Contains(dryRun.stdout, "0 missing, 0 source-stale, 0 policy-stale, 0 manual edits, 0 untracked") {
		t.Fatalf("adopted multi-format project was not current:\n%s", dryRun.stdout)
	}
}

func TestAdoptionLifecycleIsExplicitIdempotentAndReadOnlyWhenPlanned(t *testing.T) {
	projectDir := copyFixture(t, "lifecycle")
	targetPath := filepath.Join(projectDir, "locales", "fr.json")
	manifestPath := filepath.Join(projectDir, ".internationalizer.lock")
	originalTarget := mustReadFile(t, targetPath)

	dryRun := runCLI(t, projectDir, nil, "translate", "--dry-run")
	dryRun.requireSuccess(t)
	if !strings.Contains(dryRun.stdout, "2 untracked") {
		t.Fatalf("initial dry-run did not report untracked translations:\n%s", dryRun.stdout)
	}
	if got := mustReadFile(t, targetPath); !bytes.Equal(got, originalTarget) {
		t.Fatalf("dry-run changed target:\n%s", got)
	}
	if _, err := os.Stat(manifestPath); !os.IsNotExist(err) {
		t.Fatalf("dry-run created manifest: %v", err)
	}

	adopt := runCLI(t, projectDir, nil, "translate", "--adopt-existing")
	adopt.requireSuccess(t)
	manifest, err := state.Load(manifestPath)
	if err != nil {
		t.Fatalf("loading adopted manifest: %v", err)
	}
	if len(manifest.Translations) != 2 {
		t.Fatalf("manifest entries = %d, want 2", len(manifest.Translations))
	}
	for identity, entry := range manifest.Translations {
		if entry.Origin != "adopted" || entry.Provider != "" || entry.Model != "" {
			t.Fatalf("manifest entry %s provenance = %#v, want adopted without provider", identity, entry)
		}
	}
	manifestAfterAdoption := mustReadFile(t, manifestPath)

	repeatedAdoption := runCLI(t, projectDir, nil, "translate", "--adopt-existing")
	repeatedAdoption.requireSuccess(t)
	if got := mustReadFile(t, manifestPath); !bytes.Equal(got, manifestAfterAdoption) {
		t.Fatalf("repeated adoption rewrote current manifest:\n%s", got)
	}

	validation := runCLI(t, projectDir, nil, "validate", "--json")
	validation.requireSuccess(t)
	if !strings.Contains(validation.stdout, `"coverage": 100`) {
		t.Fatalf("successful validation did not report full coverage:\n%s", validation.stdout)
	}

	sourcePath := filepath.Join(projectDir, "locales", "en.json")
	if err := os.WriteFile(sourcePath, []byte("{\n  \"save\": \"Save changes\",\n  \"welcome\": \"Hello, {name}\"\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	staleDryRun := runCLI(t, projectDir, nil, "translate", "--dry-run")
	staleDryRun.requireSuccess(t)
	if !strings.Contains(staleDryRun.stdout, "1 source-stale") {
		t.Fatalf("source edit was not reported as stale:\n%s", staleDryRun.stdout)
	}
	if got := mustReadFile(t, manifestPath); !bytes.Equal(got, manifestAfterAdoption) {
		t.Fatalf("stale dry-run changed manifest:\n%s", got)
	}

	if err := os.WriteFile(targetPath, []byte("{\n  \"save\": \"Enregistrer\"\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	failedValidation := runCLI(t, projectDir, nil, "validate", "--json")
	if failedValidation.exitCode != 1 {
		t.Fatalf("validation exit code = %d, want 1; stderr=%s", failedValidation.exitCode, failedValidation.stderr)
	}
	if !strings.Contains(failedValidation.stdout, `"welcome"`) || !strings.Contains(failedValidation.stdout, `"missing"`) {
		t.Fatalf("failed validation did not report missing key:\n%s", failedValidation.stdout)
	}
}

func TestStrictValidationDistinguishesSeededContentFromCurrentState(t *testing.T) {
	projectDir := t.TempDir()
	mustWriteFile(t, filepath.Join(projectDir, "locales", "en.json"), "{\n  \"brand\": \"Lens\",\n  \"seeded\": \"Save\"\n}\n")
	mustWriteFile(t, filepath.Join(projectDir, "locales", "fr.json"), "{\n  \"brand\": \"Lens\",\n  \"seeded\": \"Save\"\n}\n")
	mustWriteFile(t, filepath.Join(projectDir, "glossary", "fr.json"), "[{\"source\":\"Lens\",\"target\":\"Lens\",\"whole_word\":true}]\n")
	mustWriteFile(t, filepath.Join(projectDir, ".internationalizer.yml"), `source_locale: en
target_locales: [fr]
source_path: locales/en.json
glossary_dir: glossary
manifest_path: .internationalizer.lock
`)

	legacy := runCLI(t, projectDir, nil, "validate", "--json")
	legacy.requireSuccess(t)

	strict := runCLI(t, projectDir, nil, "validate", "--strict", "--json")
	if strict.exitCode != 1 {
		t.Fatalf("strict exit code = %d, want 1; stderr=%s", strict.exitCode, strict.stderr)
	}
	for _, want := range []string{`"structural_coverage": 100`, `"translated_coverage": 50`, `"code": "source_identical"`} {
		if !strings.Contains(strict.stdout, want) {
			t.Fatalf("strict report does not contain %q:\n%s", want, strict.stdout)
		}
	}

	adoption := runCLI(t, projectDir, nil, "translate", "--adopt-existing")
	adoption.requireSuccess(t)
	requireState := runCLI(t, projectDir, nil, "validate", "--require-state", "--json")
	requireState.requireSuccess(t)
	if strings.Contains(requireState.stdout, `"code": "untracked"`) {
		t.Fatalf("adopted state remained untracked:\n%s", requireState.stdout)
	}
}

func TestFluentPseudoAndExplicitReviewLifecycle(t *testing.T) {
	projectDir := t.TempDir()
	mustWriteFile(t, filepath.Join(projectDir, "locales", "en.ftl"), `# Link shown below the item count.
items =
    { $count ->
        [one] One item
       *[other] { $count } items
    }
learn-more = See <a data-l10n-name="docs">the documentation</a>.
save-button =
    .label = Save
`)
	mustWriteFile(t, filepath.Join(projectDir, "locales", "fr.ftl"), `# Link shown below the item count.
items =
    { $count ->
        [one] Un élément
       *[other] { $count } éléments
    }
learn-more = Consultez <a data-l10n-name="docs">la documentation</a>.
save-button =
    .label = Enregistrer
`)
	mustWriteFile(t, filepath.Join(projectDir, ".internationalizer.yml"), `source_locale: en
target_locales: [fr]
bundles:
  - id: browser
    source: locales/en.ftl
    target: locales/{locale}.ftl
    format: fluent
manifest_path: .internationalizer.lock
`)

	adopt := runCLI(t, projectDir, nil, "translate", "--adopt-existing")
	adopt.requireSuccess(t)
	listed := runCLI(t, projectDir, nil, "review", "list", "--status", "needs_review")
	listed.requireSuccess(t)
	if got := strings.Count(listed.stdout, "needs_review"); got != 3 {
		t.Fatalf("needs-review entries = %d, want 3:\n%s", got, listed.stdout)
	}

	required := runCLI(t, projectDir, nil, "validate", "--require-approved", "--json")
	if required.exitCode != 1 || !strings.Contains(required.stdout, `"code": "needs_review"`) {
		t.Fatalf("unapproved validation = exit %d:\n%s\n%s", required.exitCode, required.stdout, required.stderr)
	}
	approve := runCLI(t, projectDir, nil, "review", "approve", "--locale", "fr", "--all")
	approve.requireSuccess(t)
	approved := runCLI(t, projectDir, nil, "validate", "--require-approved", "--json")
	approved.requireSuccess(t)

	pseudo := runCLI(t, projectDir, nil, "pseudo", "--strategy", "accented")
	pseudo.requireSuccess(t)
	pseudoOutput := string(mustReadFile(t, filepath.Join(projectDir, "locales", "en-XA.ftl")))
	for _, protected := range []string{"{ $count ->", "*[other]", "{ $count }", `data-l10n-name="docs"`, ".label ="} {
		if !strings.Contains(pseudoOutput, protected) {
			t.Fatalf("pseudo Fluent output lost %q:\n%s", protected, pseudoOutput)
		}
	}
	manifest, err := state.Load(filepath.Join(projectDir, ".internationalizer.lock"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range manifest.Translations {
		switch entry.Locale {
		case "fr":
			if entry.ReviewStatus != state.ReviewApproved {
				t.Fatalf("French entry was not approved: %#v", entry)
			}
		case "en-XA":
			if entry.Origin != "pseudo" || entry.ReviewStatus != state.ReviewNeedsReview {
				t.Fatalf("pseudo entry provenance = %#v", entry)
			}
		}
	}
}

func TestOpenAITranslationExercisesResponsesContractAndPreservesSourceShape(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.URL.Path != "/api.openai.com/v1/responses" {
			t.Errorf("request path = %q, want Responses API path", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer acceptance-key" {
			t.Errorf("authorization = %q", got)
		}
		bodyData, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request: %v", err)
			return
		}
		var body map[string]any
		if err := json.Unmarshal(bodyData, &body); err != nil {
			t.Errorf("decoding request: %v", err)
			return
		}
		if body["model"] != "gpt-5.6-luna" {
			t.Errorf("model = %#v, want gpt-5.6-luna", body["model"])
		}
		reasoning, ok := body["reasoning"].(map[string]any)
		if !ok || reasoning["effort"] != "max" {
			t.Errorf("reasoning = %#v, want max", body["reasoning"])
		}
		if _, ok := body["temperature"]; ok {
			t.Error("GPT-5.6 Responses request included temperature")
		}
		input, ok := body["input"].(string)
		if !ok {
			t.Errorf("input = %#v, want encoded JSON string", body["input"])
		} else {
			var entries map[string]string
			if err := json.Unmarshal([]byte(input), &entries); err != nil {
				t.Errorf("decoding input entries: %v", err)
			} else if len(entries) != 5 {
				t.Errorf("input entries = %d, want 5: %#v", len(entries), entries)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"output":[{"type":"message","content":[{"type":"output_text","text":"{\"account\":{\"save\":\"Enregistrer\",\"welcome\":\"Bonjour, {name}\"},\"steps\":{\"0\":\"Un\",\"1\":\"Deux\"},\"limits\":{\"0\":\"Zéro\"}}"}]}],"usage":{"input_tokens":41,"output_tokens":19}}`)
	}))
	defer server.Close()

	projectDir := t.TempDir()
	mustWriteFile(t, filepath.Join(projectDir, "locales", "en.json"), `{
  "account": {"save": "Save", "welcome": "Hello, {name}"},
  "steps": ["One", "Two"],
  "limits": {"0": "Zero"},
  "enabled": true,
  "retries": 3,
  "nullable": null
}
`)
	configData := fmt.Sprintf(`source_locale: en
target_locales: [fr]
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
llm:
  provider: openai
  model: gpt-5.6-luna
  api_key_env: INTERNATIONALIZER_ACCEPTANCE_KEY
  # Keep the official-host marker while routing the request to the local server.
  base_url: %s/api.openai.com
manifest_path: .internationalizer.lock
tm_path: .internationalizer/tm.jsonl
`, server.URL)
	mustWriteFile(t, filepath.Join(projectDir, ".internationalizer.yml"), configData)

	translation := runCLI(t, projectDir, []string{"INTERNATIONALIZER_ACCEPTANCE_KEY=acceptance-key"}, "translate")
	translation.requireSuccess(t)
	if !strings.Contains(translation.stdout, "Translated 5 keys") || !strings.Contains(translation.stdout, "Tokens: 41 input, 19 output") {
		t.Fatalf("translation summary missing provider results:\n%s", translation.stdout)
	}
	if requests.Load() != 1 {
		t.Fatalf("provider requests = %d, want 1", requests.Load())
	}

	var target map[string]any
	if err := json.Unmarshal(mustReadFile(t, filepath.Join(projectDir, "locales", "fr.json")), &target); err != nil {
		t.Fatalf("decoding target: %v", err)
	}
	if target["enabled"] != true || target["retries"] != float64(3) || target["nullable"] != nil {
		t.Fatalf("non-string source shape was not preserved: %#v", target)
	}
	steps, ok := target["steps"].([]any)
	if !ok || len(steps) != 2 || steps[0] != "Un" || steps[1] != "Deux" {
		t.Fatalf("translated array = %#v", target["steps"])
	}
	limits, ok := target["limits"].(map[string]any)
	if !ok || limits["0"] != "Zéro" {
		t.Fatalf("numeric object key was not preserved: %#v", target["limits"])
	}

	manifest, err := state.Load(filepath.Join(projectDir, ".internationalizer.lock"))
	if err != nil {
		t.Fatalf("loading provider manifest: %v", err)
	}
	if len(manifest.Translations) != 5 {
		t.Fatalf("provider manifest entries = %d, want 5", len(manifest.Translations))
	}
	for identity, entry := range manifest.Translations {
		if entry.Origin != "provider" || entry.Provider != "openai" || entry.Model != "gpt-5.6-luna" {
			t.Fatalf("manifest entry %s provenance = %#v", identity, entry)
		}
	}
}

type cliResult struct {
	stdout   string
	stderr   string
	exitCode int
}

func (r cliResult) requireSuccess(t *testing.T) {
	t.Helper()
	if r.exitCode != 0 {
		t.Fatalf("CLI exit code = %d, want 0\nstdout:\n%s\nstderr:\n%s", r.exitCode, r.stdout, r.stderr)
	}
	if r.stderr != "" {
		t.Fatalf("successful CLI wrote stderr:\n%s", r.stderr)
	}
}

func runCLI(t *testing.T, dir string, environment []string, args ...string) cliResult {
	t.Helper()
	command := exec.Command(binaryPath, args...)
	command.Dir = dir
	command.Env = append(os.Environ(), environment...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err == nil {
		return cliResult{stdout: stdout.String(), stderr: stderr.String()}
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("running CLI: %v", err)
	}
	return cliResult{stdout: stdout.String(), stderr: stderr.String(), exitCode: exitError.ExitCode()}
}

func copyFixture(t *testing.T, name string) string {
	t.Helper()
	projectDir := t.TempDir()
	fixtureRoot := filepath.Join("testdata", name)
	if err := os.CopyFS(projectDir, os.DirFS(fixtureRoot)); err != nil {
		t.Fatalf("copying fixture %s: %v", name, err)
	}
	return projectDir
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
