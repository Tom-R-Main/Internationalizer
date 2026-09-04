package acceptance_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

func TestExecuFunctionSyntaxLifecycle(t *testing.T) {
	dir := copyFixture(t, "execufunction-syntax")
	targetPath := filepath.Join(dir, "web", "fr.json")
	original := mustReadFile(t, targetPath)
	runCLI(t, dir, nil, "validate", "--strict", "--json").requireSuccess(t)
	dry := runCLI(t, dir, nil, "translate", "--dry-run")
	dry.requireSuccess(t)
	if !strings.Contains(dry.stdout, "Would translate 0 keys") {
		t.Fatal(dry.stdout)
	}
	if !bytes.Equal(original, mustReadFile(t, targetPath)) {
		t.Fatal("dry-run modified target")
	}
	if _, err := os.Stat(filepath.Join(dir, ".internationalizer.lock")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote state: %v", err)
	}
	for _, strategy := range []string{"accented", "bidi"} {
		result := runCLI(t, dir, nil, "pseudo", "--strategy", strategy, "--dry-run")
		result.requireSuccess(t)
		runCLI(t, dir, nil, "pseudo", "--strategy", strategy).requireSuccess(t)
		locale := "en-XA"
		if strategy == "bidi" {
			locale = "ar-XB"
		}
		var generated map[string]string
		if err := json.Unmarshal(mustReadFile(t, filepath.Join(dir, "web", locale+".json")), &generated); err != nil {
			t.Fatal(err)
		}
		for _, literal := range []string{`<code>&lt;root&gt;/{.sift,.claude,.codex,.agents}/skills</code>`, "{{user.name}}"} {
			if !strings.Contains(generated["docs.tui.skills.desc"], literal) {
				t.Fatalf("%s changed %s: %s", strategy, literal, generated["docs.tui.skills.desc"])
			}
		}
		if !strings.Contains(generated["price"], "{{amount, number}}") || !strings.Contains(generated["price"], "{{- customer.name}}") {
			t.Fatalf("%s damaged interpolation: %s", strategy, generated["price"])
		}
	}
	runCLI(t, dir, nil, "translate", "--adopt-existing").requireSuccess(t)
	manifest, err := state.Load(filepath.Join(dir, ".internationalizer.lock"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range manifest.Translations {
		if entry.ReviewStatus != state.ReviewNeedsReview {
			t.Fatal("adoption/pseudo implied approval")
		}
	}
	runCLI(t, dir, nil, "review", "approve", "--locale", "fr", "--all").requireSuccess(t)
	runCLI(t, dir, nil, "validate", "--require-approved").requireSuccess(t)
	configPath := filepath.Join(dir, ".internationalizer.yml")
	config := string(mustReadFile(t, configPath))
	mustWriteFile(t, configPath, strings.Replace(config, "message_syntax: i18next", "message_syntax: plain", 1))
	stale := runCLI(t, dir, nil, "validate", "--require-approved", "--json")
	if stale.exitCode == 0 || !strings.Contains(stale.stdout, "policy_stale") {
		t.Fatalf("policy change retained approval: %+v", stale)
	}
	if result := runCLI(t, dir, nil, "review", "approve", "--locale", "fr", "--all"); result.exitCode == 0 {
		t.Fatal("approved stale syntax policy")
	}
}

func TestSourceErrorsReportedOnceAcrossLocales(t *testing.T) {
	dir := copyFixture(t, "execufunction-syntax")
	mustWriteFile(t, filepath.Join(dir, "icu", "en.json"), `{"items":"{count, plural, one {One}}"}`)
	configPath := filepath.Join(dir, ".internationalizer.yml")
	mustWriteFile(t, configPath, strings.Replace(string(mustReadFile(t, configPath)), "[fr]", "[fr, de, es]", 1))
	for _, args := range [][]string{{"validate", "--json"}, {"translate", "--dry-run"}} {
		result := runCLI(t, dir, nil, args...)
		if result.exitCode == 0 || strings.Count(result.stdout, "source ICU message:") != 1 {
			t.Fatalf("%v did not report source once: %+v", args, result)
		}
		if !strings.Contains(result.stdout, "icu/en.json") && !strings.Contains(result.stdout, `icu\\en.json`) && !strings.Contains(result.stdout, `icu\en.json`) {
			t.Fatalf("missing source path: %s", result.stdout)
		}
	}
}
