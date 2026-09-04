package acceptance_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This trajectory deliberately starts with a valid marketing-only config and
// a nested app. No network, real credentials, or repository catalogs are used.
func onboardingFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		".internationalizer.yml": `# Marketing pipeline; preserve this comment.
source_locale: en
target_locales: [fr, ja]
source_path: tmp/english-keys.json
llm:
  provider: openai
  api_key_env: ONBOARDING_TEST_KEY
  locale_overrides:
    ja:
      provider: gemini
      model: custom-japanese-model
      api_key_env: ONBOARDING_JA_KEY
glossary_dir: custom/glossary
style_guides_dir: custom/guides
future_setting: preserve-me
`,
		"tmp/english-keys.json":                `{"docs.tui.skills.desc":"Use <code>{.sift,.claude,.codex,.agents}/skills</code>","hello":"Hello {{name}}"}`,
		"tmp/fr.json":                          `{"docs.tui.skills.desc":"Utiliser <code>{.sift,.claude,.codex,.agents}/skills</code>","hello":"Bonjour {{name}}"}`,
		"tmp/ja.json":                          `{"docs.tui.skills.desc":"使用 <code>{.sift,.claude,.codex,.agents}/skills</code>","hello":"こんにちは {{name}}"}`,
		"exf-app/web/package.json":             `{"dependencies":{"i18next":"25.0.0","react-i18next":"16.0.0"}}`,
		"exf-app/web/src/i18n/index.ts":        `import i18next from 'i18next'; import en from './locales/en.json'; i18next.init({resources:{en}});`,
		"exf-app/web/src/i18n/locales/en.json": `{"hello":"Hello {{name}}","count_one":"{{count}} item","count_other":"{{count}} items"}`,
	}
	for name, contents := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func decodeOnboardingJSON(t *testing.T, result cliResult) map[string]any {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal([]byte(result.stdout), &doc); err != nil {
		t.Fatalf("invalid JSON: %v\nstdout=%s\nstderr=%s", err, result.stdout, result.stderr)
	}
	if doc["schema_version"] != float64(1) {
		t.Fatalf("unversioned output: %s", result.stdout)
	}
	return doc
}

func TestOnboardingDiscoverPlanApplyVerify(t *testing.T) {
	root := onboardingFixture(t)
	original := mustReadFile(t, filepath.Join(root, ".internationalizer.yml"))
	discovery := runCLI(t, root, nil, "detect", "--json")
	discovery.requireSuccess(t)
	decodeOnboardingJSON(t, discovery)
	for _, want := range []string{"UNCOVERED_CATALOG", "SOURCE_CONFIRMATION_REQUIRED", "AUTO_SYNTAX_AMBIGUOUS", "exf-app/web/src/i18n/locales/en.json", `"suggested_syntax": "i18next"`} {
		if !strings.Contains(discovery.stdout, want) {
			t.Fatalf("missing %q: %s", want, discovery.stdout)
		}
	}
	check := runCLI(t, root, nil, "config", "check", "--json")
	decodeOnboardingJSON(t, check)
	if !strings.Contains(check.stdout, `"provider_verified": false`) {
		t.Fatal("offline readiness missing")
	}
	blocked := runCLI(t, root, nil, "config", "plan", "--json")
	blocked.requireSuccess(t)
	if decodeOnboardingJSON(t, blocked)["status"] != "needs_decision" {
		t.Fatalf("should need decisions: %s", blocked.stdout)
	}
	if !bytes.Equal(original, mustReadFile(t, filepath.Join(root, ".internationalizer.yml"))) {
		t.Fatal("inspection or planning modified config")
	}
	planArgs := []string{"config", "plan", "--json", "--add-bundle", "web=exf-app/web/src/i18n/locales/en.json", "--syntax", "web=i18next", "--syntax", "default=i18next", "--confirm-source", "tmp/english-keys.json", "--out", "onboarding-plan.json"}
	plan := runCLI(t, root, nil, planArgs...)
	plan.requireSuccess(t)
	if decodeOnboardingJSON(t, plan)["status"] != "planned" {
		t.Fatalf("explicit choices should resolve plan: %s", plan.stdout)
	}
	apply := runCLI(t, root, nil, "config", "apply", "--plan", "onboarding-plan.json", "--no-input", "--json")
	apply.requireSuccess(t)
	decodeOnboardingJSON(t, apply)
	configured := string(mustReadFile(t, filepath.Join(root, ".internationalizer.yml")))
	for _, want := range []string{"# Marketing pipeline", "custom-japanese-model", "ONBOARDING_JA_KEY", "custom/glossary", "custom/guides", "future_setting: preserve-me", "id: default", "id: web"} {
		if !strings.Contains(configured, want) {
			t.Fatalf("lost %q: %s", want, configured)
		}
	}
	replay := runCLI(t, root, nil, "config", "apply", "--plan", "onboarding-plan.json", "--no-input", "--json")
	replay.requireSuccess(t)
	if !strings.Contains(replay.stdout, "already_applied") {
		t.Fatalf("retry not recognized: %s", replay.stdout)
	}
	dry := runCLI(t, root, nil, "translate", "--dry-run", "--json")
	dry.requireSuccess(t)
	if decodeOnboardingJSON(t, dry)["status"] != "planned" {
		t.Fatalf("dry run did not plan: %s", dry.stdout)
	}
	for _, want := range []string{`"provider_called": false`, `"generated_keys": 0`} {
		if !strings.Contains(dry.stdout, want) {
			t.Fatalf("missing %q: %s", want, dry.stdout)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".internationalizer.lock")); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote state: %v", err)
	}
	// A configured runtime is not a waiver for catalog damage.
	target := filepath.Join(root, "tmp/fr.json")
	if err := os.WriteFile(target, []byte(`{"docs.tui.skills.desc":"<code>different-command</code>","hello":"Bonjour {{other}}"}`), 0600); err != nil {
		t.Fatal(err)
	}
	invalid := runCLI(t, root, nil, "validate", "--bundle", "default", "--locale", "fr", "--json")
	if invalid.exitCode == 0 {
		t.Fatal("damaged placeholders and code passed validation")
	}
	decodeOnboardingJSON(t, invalid)
	if !strings.Contains(invalid.stdout, "HTML code mismatch") || !strings.Contains(invalid.stdout, "interpolation") {
		t.Fatalf("expected code and interpolation findings: %s", invalid.stdout)
	}
}

func TestOnboardingJSONFailuresAndPlanDrift(t *testing.T) {
	root := onboardingFixture(t)
	invalid := runCLI(t, root, nil, "config", "plan", "--unknown-option", "--json")
	if invalid.exitCode == 0 {
		t.Fatal("unknown option accepted")
	}
	decodeOnboardingJSON(t, invalid)
	planned := runCLI(t, root, nil, "config", "plan", "--syntax", "default=i18next", "--confirm-source", "tmp/english-keys.json", "--out", "plan.json", "--json")
	planned.requireSuccess(t)
	source := filepath.Join(root, "tmp/english-keys.json")
	if err := os.WriteFile(source, []byte(`{"changed":"new source"}`), 0600); err != nil {
		t.Fatal(err)
	}
	applied := runCLI(t, root, nil, "config", "apply", "--plan", "plan.json", "--no-input", "--json")
	if applied.exitCode == 0 {
		t.Fatal("stale plan applied")
	}
	decodeOnboardingJSON(t, applied)
	if !strings.Contains(applied.stdout, `"code": "stale_plan"`) {
		t.Fatalf("missing drift code: %s", applied.stdout)
	}
}
