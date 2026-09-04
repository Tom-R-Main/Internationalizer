package onboarding

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
)

func retargetFixture(t *testing.T, explicit bool) string {
	t.Helper()
	root := t.TempDir()
	body := "# keep header\nsource_locale: en\ntarget_locales: [fr]\nmessage_syntax: plain\nllm:\n  provider: openai\n  locale_overrides:\n    fr:\n      model: custom\nglossary_dir: custom/glossary\nfuture_setting: retained\n"
	if explicit {
		body += "bundles:\n  - id: marketing\n    source: en.json\n    target: linked/{locale}.json # keep target comment\n    format: json\n    message_syntax: plain\n    future_bundle_setting: retained\n"
	} else {
		body += "source_path: en.json # keep source comment\n"
	}
	planWrite(t, root, ".internationalizer.yml", body)
	planWrite(t, root, "en.json", `{"hello":"Hello"}`)
	return root
}

func TestRetargetPreservesExistingSettings(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		t.Run(map[bool]string{false: "legacy", true: "explicit"}[explicit], func(t *testing.T) {
			root := retargetFixture(t, explicit)
			id := "default"
			if explicit {
				id = "marketing"
			}
			p, err := BuildPlan(root, "", PlanOptions{UpdateTargets: map[string]string{id: "catalogs/{locale}.json"}})
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{"# keep header", "source_locale: en", "provider: openai", "model: custom", "glossary_dir: custom/glossary", "future_setting: retained", "id: " + id, "source: en.json", "target: catalogs/{locale}.json", "message_syntax: plain"} {
				if !strings.Contains(p.ProposedYAML, want) {
					t.Errorf("lost %q: %s", want, p.ProposedYAML)
				}
			}
			if explicit && (!strings.Contains(p.ProposedYAML, "# keep target comment") || !strings.Contains(p.ProposedYAML, "future_bundle_setting: retained")) {
				t.Fatal("lost bundle metadata")
			}
			if !explicit && !strings.Contains(p.ProposedYAML, "# keep source comment") {
				t.Fatal("lost legacy comment")
			}
			receipt, err := ApplyPlan(p)
			if err != nil || receipt.Status != "applied" {
				t.Fatalf("apply: %+v %v", receipt, err)
			}
			receipt, err = ApplyPlan(p)
			if err != nil || receipt.Status != "already_applied" {
				t.Fatalf("replay: %+v %v", receipt, err)
			}
		})
	}
}

func TestRetargetRejectsInvalidDecisions(t *testing.T) {
	for name, options := range map[string]PlanOptions{
		"unknown":        {UpdateTargets: map[string]string{"unknown": "catalogs/{locale}.json"}},
		"empty_id":       {UpdateTargets: map[string]string{"": "catalogs/{locale}.json"}},
		"empty_target":   {UpdateTargets: map[string]string{"default": ""}},
		"add_existing":   {AddBundles: []config.Bundle{{ID: "default", Source: "en.json", Target: "catalogs/{locale}.json"}}},
		"add_and_update": {AddBundles: []config.Bundle{{ID: "new", Source: "en.json", Target: "catalogs/{locale}.json"}}, UpdateTargets: map[string]string{"new": "catalogs/{locale}.json"}},
	} {
		t.Run(name, func(t *testing.T) {
			root := retargetFixture(t, false)
			_, err := BuildPlan(root, "", options)
			assertPlanCode(t, err, "invalid_decision")
		})
	}
}

func TestRetargetRepairsUnsafeOldLinksWithoutReadingThem(t *testing.T) {
	for _, kind := range []string{"external", "dangling", "cyclic", "ancestor"} {
		t.Run(kind, func(t *testing.T) {
			root := retargetFixture(t, true)
			linkPath := filepath.Join(root, "linked", "fr.json")
			if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
				t.Fatal(err)
			}
			var target string
			switch kind {
			case "external":
				target = filepath.Join(t.TempDir(), "never-read.json")
			case "dangling":
				target = "missing.json"
			case "cyclic":
				target = "fr.json"
			case "ancestor":
				if err := os.Remove(filepath.Dir(linkPath)); err != nil {
					t.Fatal(err)
				}
				linkPath = filepath.Dir(linkPath)
				target = "missing-dir"
			}
			if err := os.Symlink(target, linkPath); err != nil {
				t.Skipf("symlink unavailable: %v", err)
			}
			before, err := os.ReadFile(filepath.Join(root, ".internationalizer.yml"))
			if err != nil {
				t.Fatal(err)
			}
			_, err = BuildPlan(root, "", PlanOptions{})
			assertPlanCode(t, err, "unsafe_path")
			if !strings.Contains(err.Error(), `bundle "marketing" target locale "fr"`) || !strings.Contains(err.Error(), "linked/fr.json") {
				t.Fatalf("missing path context: %v", err)
			}
			p, err := BuildPlan(root, "", PlanOptions{UpdateTargets: map[string]string{"marketing": "safe/{locale}.json"}})
			if err != nil {
				t.Fatal(err)
			}
			current, _ := os.ReadFile(filepath.Join(root, ".internationalizer.yml"))
			if !bytes.Equal(current, before) {
				t.Fatal("planning changed config")
			}
			// A caller supplied the replacement independently: old link identity
			// is not plan evidence and changing it must not require following it.
			if err := os.Remove(linkPath); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("another-missing-path", linkPath); err != nil {
				t.Fatal(err)
			}
			if _, err := ApplyPlan(p); err != nil {
				t.Fatal(err)
			}
			if after, err := os.Readlink(linkPath); err != nil || after != "another-missing-path" {
				t.Fatalf("old link mutated: %q %v", after, err)
			}
		})
	}
}

func TestRetargetRejectsUnsafeReplacementAndLateLink(t *testing.T) {
	root := retargetFixture(t, false)
	_, err := BuildPlan(root, "", PlanOptions{UpdateTargets: map[string]string{"default": "../outside/{locale}.json"}})
	assertPlanCode(t, err, "unsafe_path")
	p, err := BuildPlan(root, "", PlanOptions{UpdateTargets: map[string]string{"default": "safe/{locale}.json"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(root, "safe")); err != nil {
		t.Skip(err)
	}
	_, err = ApplyPlan(p)
	assertPlanCode(t, err, "unsafe_path")
}

func TestSymlinkDiagnosticsBoundedAndInRoot(t *testing.T) {
	for _, kind := range []string{"regular", "ancestor", "external", "dangling", "cyclic", "sensitive", "directory"} {
		t.Run(kind, func(t *testing.T) {
			root := retargetFixture(t, false)
			planWrite(t, root, "real/fr.json", `not-json-and-never-opened`)
			link, target, want := "fr.json", "real/fr.json", `review in-root destination "real/fr.json"`
			switch kind {
			case "ancestor":
				link, target = "linked", "real"
			case "external":
				target, want = filepath.Join(t.TempDir(), "secret.json"), "outside the project root"
			case "dangling":
				target, want = "missing.json", "dangling"
			case "cyclic":
				target, want = "fr.json", "cyclic"
			case "sensitive":
				target, want = ".env", "credential-shaped"
			case "directory":
				target, want = "real/.", "not a regular catalog file"
			}
			if err := os.Symlink(target, filepath.Join(root, link)); err != nil {
				t.Skip(err)
			}
			path := link
			if kind == "ancestor" {
				path += "/fr.json"
			}
			_, err := safePlanPath(root, path)
			assertPlanCode(t, err, "unsafe_path")
			if !strings.Contains(err.Error(), want) || !strings.Contains(err.Error(), `link component "`+link+`"`) {
				t.Fatalf("missing %q: %v", want, err)
			}
			if kind == "external" && strings.Contains(err.Error(), target) {
				t.Fatalf("leaked external destination: %v", err)
			}
		})
	}
}

func TestSymlinkDiagnosticResolvesParentsAfterLinks(t *testing.T) {
	for _, absolute := range []bool{false, true} {
		t.Run(map[bool]string{false: "relative", true: "absolute"}[absolute], func(t *testing.T) {
			root := retargetFixture(t, false)
			// Canonicalize platform aliases before constructing an absolute link.
			root, err := canonicalPlanRoot(root)
			if err != nil {
				t.Fatal(err)
			}
			planWrite(t, root, "other/child/unused.json", "{}")
			planWrite(t, root, "other/real/fr.json", "{}")
			if err := os.Symlink(filepath.Join("other", "child"), filepath.Join(root, "branch")); err != nil {
				t.Skip(err)
			}
			target := "branch" + string(filepath.Separator) + ".." + string(filepath.Separator) + "real" + string(filepath.Separator) + "fr.json"
			if absolute {
				target = root + string(filepath.Separator) + target
			}
			if err := os.Symlink(target, filepath.Join(root, "fr.json")); err != nil {
				t.Fatal(err)
			}
			_, err = safePlanPath(root, "fr.json")
			assertPlanCode(t, err, "unsafe_path")
			if !strings.Contains(err.Error(), `review in-root destination "other/real/fr.json"`) {
				t.Fatalf("wrong parent resolution: %v", err)
			}
		})
	}
}
