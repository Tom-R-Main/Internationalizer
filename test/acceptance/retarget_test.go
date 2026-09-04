package acceptance_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOnboardingRetargetTwentyOneMarketingSymlinks(t *testing.T) {
	root := t.TempDir()
	locales := []string{"fr", "ja", "de", "es", "it", "pt", "nl", "sv", "da", "fi", "pl", "cs", "sk", "hu", "ro", "tr", "ru", "uk", "ko", "zh", "ar"}
	configPath := filepath.Join(root, ".internationalizer.yml")
	original := []byte(fmt.Sprintf(`# keep marketing ownership
source_locale: en
target_locales: [%s]
source_path: tmp/english-keys.json
message_syntax: plain
llm:
  provider: openai
  api_key_env: RETARGET_TEST_KEY
  locale_overrides:
    ja:
      provider: gemini
      model: preserved-model
      api_key_env: RETARGET_JA_TEST_KEY
glossary_dir: custom/glossary
future_setting: retained
`, strings.Join(locales, ", ")))
	if err := os.WriteFile(configPath, original, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "tmp"), 0755); err != nil {
		t.Fatal(err)
	}
	catalog := []byte(`{"hello":"Hello","code":"<code>{.sift,.agents}</code>"}`)
	if err := os.WriteFile(filepath.Join(root, "tmp", "english-keys.json"), catalog, 0600); err != nil {
		t.Fatal(err)
	}
	for _, locale := range locales {
		target := "translations-" + locale + ".json"
		if err := os.WriteFile(filepath.Join(root, "tmp", target), catalog, 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(root, "tmp", locale+".json")); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
	}
	blocked := runCLI(t, root, nil, "config", "plan", "--confirm-source", "tmp/english-keys.json", "--json")
	if blocked.exitCode == 0 {
		t.Fatal("unrepaired symlink targets accepted")
	}
	decodeOnboardingJSON(t, blocked)
	for _, want := range []string{`"code": "unsafe_path"`, "tmp/fr.json", "tmp/translations-fr.json", `locale \"fr\"`} {
		if !strings.Contains(blocked.stdout, want) {
			t.Fatalf("missing %q: %s", want, blocked.stdout)
		}
	}
	planned := runCLI(t, root, nil, "config", "plan", "--update-bundle", "default", "--target", "default=tmp/translations-{locale}.json", "--confirm-source", "tmp/english-keys.json", "--out", "repair-plan.json", "--json")
	planned.requireSuccess(t)
	if decodeOnboardingJSON(t, planned)["status"] != "planned" {
		t.Fatalf("repair requires unexpected decisions: %s", planned.stdout)
	}
	if !bytes.Equal(original, mustReadFile(t, configPath)) {
		t.Fatal("plan changed config")
	}
	applied := runCLI(t, root, nil, "config", "apply", "--plan", "repair-plan.json", "--no-input", "--json")
	applied.requireSuccess(t)
	configured := string(mustReadFile(t, configPath))
	for _, want := range []string{"# keep marketing ownership", "id: default", "source: tmp/english-keys.json", "target: tmp/translations-{locale}.json", "preserved-model", "RETARGET_JA_TEST_KEY", "custom/glossary", "future_setting: retained"} {
		if !strings.Contains(configured, want) {
			t.Fatalf("lost %q: %s", want, configured)
		}
	}
	replay := runCLI(t, root, nil, "config", "apply", "--plan", "repair-plan.json", "--no-input", "--json")
	replay.requireSuccess(t)
	if decodeOnboardingJSON(t, replay)["status"] != "already_applied" {
		t.Fatalf("replay: %s", replay.stdout)
	}
	dry := runCLI(t, root, nil, "translate", "--dry-run", "--json")
	dry.requireSuccess(t)
	decodeOnboardingJSON(t, dry)
	if !strings.Contains(dry.stdout, `"provider_called": false`) {
		t.Fatalf("dry-run called provider: %s", dry.stdout)
	}
	for _, locale := range locales {
		target := "translations-" + locale + ".json"
		if link, err := os.Readlink(filepath.Join(root, "tmp", locale+".json")); err != nil || link != target {
			t.Fatalf("link changed: %q %v", link, err)
		}
		if !bytes.Equal(catalog, mustReadFile(t, filepath.Join(root, "tmp", target))) {
			t.Fatalf("catalog %s changed", locale)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".internationalizer.lock")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote state: %v", err)
	}
}
