package onboarding

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
)

func planFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	planWrite(t, root, ".internationalizer.yml", "# retain me\nsource_locale: en\ntarget_locales: [fr]\n# marketing owner\nsource_path: tmp/en.json # staging source\nllm:\n  provider: openai\n  api_key_env: TEST_KEY\n  locale_overrides:\n    fr:\n      model: custom\nglossary_dir: custom/glossary\nfuture_setting: retained\n")
	planWrite(t, root, "tmp/en.json", `{"code":"<code>{.sift,.agents}</code>"}`)
	planWrite(t, root, "web/locales/en.json", `{"hello":"Hello {{name}}"}`)
	return root
}

func planWrite(t *testing.T, root, path, data string) {
	t.Helper()
	p := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
}

func planOptions() PlanOptions {
	return PlanOptions{AddBundles: []config.Bundle{{ID: "web", Source: "web/locales/en.json", Target: "web/locales/{locale}.json", MessageSyntax: message.I18next}}, Syntax: map[string]message.Syntax{"default": message.Plain}, ConfirmSources: []string{"tmp/en.json"}}
}

func TestPlanPreservesConfigAndReplay(t *testing.T) {
	root := planFixture(t)
	p, err := BuildPlan(root, "", planOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(p.RequiredDecisions) != 0 {
		t.Fatalf("decisions: %+v", p.RequiredDecisions)
	}
	p2, err := BuildPlan(root, "", planOptions())
	if err != nil || p.ID != p2.ID {
		t.Fatalf("plan not deterministic: %v", err)
	}
	for _, retained := range []string{"# retain me", "# marketing owner", "# staging source", "future_setting: retained", "glossary_dir: custom/glossary", "model: custom", "id: default", "id: web", "message_syntax: plain", "message_syntax: i18next"} {
		if !strings.Contains(p.ProposedYAML, retained) {
			t.Errorf("missing %q in %s", retained, p.ProposedYAML)
		}
	}
	before, _ := os.ReadFile(filepath.Join(root, ".internationalizer.yml"))
	if strings.Contains(string(before), "id: web") {
		t.Fatal("planning wrote config")
	}
	r, err := ApplyPlan(p)
	if err != nil || r.Status != "applied" || len(r.ChangedPaths) != 1 {
		t.Fatalf("apply: %+v %v", r, err)
	}
	r, err = ApplyPlan(p)
	if err != nil || r.Status != "already_applied" || len(r.ChangedPaths) != 0 {
		t.Fatalf("replay: %+v %v", r, err)
	}
}

func TestPlanRequiresDecisions(t *testing.T) {
	root := planFixture(t)
	p, err := BuildPlan(root, "", PlanOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.RequiredDecisions) == 0 {
		t.Fatal("missing source/syntax decisions")
	}
	_, err = ApplyPlan(p)
	assertPlanCode(t, err, "decisions_required")
}

func TestApplyRejectsDriftTamperingAndLocks(t *testing.T) {
	for _, kind := range []string{"source", "config", "tamper", "lock", "symlink"} {
		t.Run(kind, func(t *testing.T) {
			root := planFixture(t)
			p, err := BuildPlan(root, "", planOptions())
			if err != nil {
				t.Fatal(err)
			}
			code := "stale_plan"
			switch kind {
			case "source":
				planWrite(t, root, "tmp/en.json", `{"new":"changed"}`)
			case "config":
				planWrite(t, root, ".internationalizer.yml", "changed: true\n")
			case "tamper":
				p.ProposedYAML += "tampered: true\n"
				code = "invalid_plan"
			case "lock":
				planWrite(t, root, ".internationalizer.yml.apply-lock", "owned")
				code = "apply_locked"
			case "symlink":
				if err := os.Remove(filepath.Join(root, "tmp/en.json")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(root, "web/locales/en.json"), filepath.Join(root, "tmp/en.json")); err != nil {
					t.Skip(err)
				}
				code = "unsafe_path"
			}
			_, err = ApplyPlan(p)
			assertPlanCode(t, err, code)
		})
	}
}

func TestPlanRejectsUnsafePathsAndSecrets(t *testing.T) {
	root := planFixture(t)
	_, err := BuildPlan(root, "../outside.yml", planOptions())
	assertPlanCode(t, err, "unsafe_path")
	planWrite(t, root, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr]\nsource_path: tmp/en.json\napi_key: never-serialize-this\n")
	p, err := BuildPlan(root, "", planOptions())
	assertPlanCode(t, err, "inline_secret")
	if p != nil || strings.Contains(err.Error(), "never-serialize-this") {
		t.Fatal("secret exposed")
	}
}

func TestPlanCreatesConfigAndObservesAbsentRuntimeEvidence(t *testing.T) {
	root := t.TempDir()
	planWrite(t, root, "locales/en.json", `{"hello":"Hello {{name}}"}`)
	options := PlanOptions{TargetLocales: []string{"fr"}, AddBundles: []config.Bundle{{ID: "web", Source: "locales/en.json", Target: "locales/{locale}.json", MessageSyntax: message.I18next}}}
	plan, err := BuildPlan(root, "", options)
	if err != nil || plan.BeforeExists || len(plan.RequiredDecisions) > 0 {
		t.Fatalf("fresh plan: %+v %v", plan, err)
	}
	planWrite(t, root, "package.json", `{"dependencies":{"i18next-icu":"1.0.0"}}`)
	_, err = ApplyPlan(plan)
	assertPlanCode(t, err, "stale_plan")
	plan, err = BuildPlan(root, "", options)
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := ApplyPlan(plan)
	if err != nil || receipt.Status != "applied" {
		t.Fatalf("fresh apply: %+v %v", receipt, err)
	}
	info, err := os.Stat(filepath.Join(root, ".internationalizer.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Fatalf("unexpected permissions %v", info.Mode().Perm())
	}
}

func TestPlanPreservesExistingBundleSettings(t *testing.T) {
	root := t.TempDir()
	planWrite(t, root, ".internationalizer.yaml", "source_locale: en\ntarget_locales: [fr]\nmessage_syntax: auto\nbundles:\n  - id: marketing\n    source: locales/en.json\n    target: locales/{locale}.json\n    message_syntax: auto # profile note\n    future_bundle_option: keep\nllm:\n  provider: openai\n  locale_overrides:\n    fr:\n      provider: gemini\n      api_key_env: FRENCH_KEY\n")
	planWrite(t, root, "locales/en.json", `{"hello":"Hello"}`)
	p, err := BuildPlan(root, "", PlanOptions{Syntax: map[string]message.Syntax{"marketing": message.Plain}})
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"id: marketing", "future_bundle_option: keep", "# profile note", "provider: gemini", "api_key_env: FRENCH_KEY"} {
		if !strings.Contains(p.ProposedYAML, fragment) {
			t.Errorf("lost %q", fragment)
		}
	}
	if !strings.HasSuffix(p.ConfigPath, ".internationalizer.yaml") {
		t.Fatalf("wrong config path %s", p.ConfigPath)
	}
	if _, err = ApplyPlan(p); err != nil {
		t.Fatal(err)
	}
}

func TestPlanMissingChoicesAndInvalidDecisions(t *testing.T) {
	t.Run("empty project", func(t *testing.T) {
		p, err := BuildPlan(t.TempDir(), "", PlanOptions{})
		if err != nil || len(p.RequiredDecisions) != 2 {
			t.Fatalf("empty plan: %+v %v", p, err)
		}
	})
	t.Run("missing syntax", func(t *testing.T) {
		root := t.TempDir()
		planWrite(t, root, "en.json", `{"hello":"hello"}`)
		p, err := BuildPlan(root, "", PlanOptions{TargetLocales: []string{"fr"}, AddBundles: []config.Bundle{{ID: "web", Source: "en.json", Target: "{locale}.json"}}})
		if err != nil || len(p.RequiredDecisions) != 1 || p.RequiredDecisions[0].Code != "SYNTAX_SELECTION_REQUIRED" {
			t.Fatalf("missing syntax: %+v %v", p, err)
		}
	})
	t.Run("missing source", func(t *testing.T) {
		p, err := BuildPlan(t.TempDir(), "", PlanOptions{TargetLocales: []string{"fr"}, AddBundles: []config.Bundle{{ID: "web", Source: "en.json", Target: "{locale}.json", MessageSyntax: message.Plain}}})
		if err != nil || len(p.RequiredDecisions) != 1 || p.RequiredDecisions[0].Code != "SOURCE_NOT_FOUND" {
			t.Fatalf("missing source: %+v %v", p, err)
		}
	})
	for _, options := range []PlanOptions{
		{Syntax: map[string]message.Syntax{"absent": message.Plain}},
		{Syntax: map[string]message.Syntax{"default": "unsupported"}},
		{AddBundles: []config.Bundle{{ID: "default", Source: "tmp/en.json", Target: "tmp/{locale}.json"}}},
		{AddBundles: []config.Bundle{{ID: "other", Source: "tmp/en.json"}}},
	} {
		root := planFixture(t)
		_, err := BuildPlan(root, "", options)
		assertPlanCode(t, err, "invalid_decision")
	}
}

func TestPlanRejectsYAMLAliasesDuplicateKeysAndCredentialURLs(t *testing.T) {
	for _, tc := range []struct{ extra, code string }{
		{"future: &alias {one: two}\nother: *alias\n", "unsupported_config"},
		{"future: one\nfuture: two\n", "invalid_config"},
		{"llm:\n  base_url: https://example.com/v1?key=secret\n", "inline_secret"},
	} {
		root := t.TempDir()
		planWrite(t, root, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr]\nsource_path: en.json\n"+tc.extra)
		_, err := BuildPlan(root, "", PlanOptions{})
		assertPlanCode(t, err, tc.code)
	}
}

func TestApplyReplayReportsConfigMatchWithoutRevalidatingSource(t *testing.T) {
	root := planFixture(t)
	p, err := BuildPlan(root, "", planOptions())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ApplyPlan(p); err != nil {
		t.Fatal(err)
	}
	planWrite(t, root, "tmp/en.json", `{"new":"new"}`)
	receipt, err := ApplyPlan(p)
	if err != nil || receipt.Status != "already_applied" || receipt.ObservationsRevalidated {
		t.Fatalf("replay: %+v %v", receipt, err)
	}
}

func TestApplyRejectsNewDiscoveryEvidence(t *testing.T) {
	root := planFixture(t)
	planWrite(t, root, "web/package.json", `{"dependencies":{"i18next":"1.0.0"}}`)
	p, err := BuildPlan(root, "", planOptions())
	if err != nil {
		t.Fatal(err)
	}
	planWrite(t, root, "web/i18n.ts", `import ICU from "i18next-icu";`)
	_, err = ApplyPlan(p)
	assertPlanCode(t, err, "stale_plan")
}

func TestPlanDoesNotSilentlyShadowHomeConfiguration(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	planWrite(t, home, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr]\nsource_path: en.json\n")
	root := t.TempDir()
	_, err := BuildPlan(root, "", PlanOptions{})
	assertPlanCode(t, err, "external_config_scope")
	p, err := BuildPlan(root, ".internationalizer.yml", PlanOptions{})
	if err != nil || p.BeforeExists {
		t.Fatalf("explicit local config: %+v %v", p, err)
	}
}

func TestPlanRejectsCredentialShapedFiles(t *testing.T) {
	for _, name := range []string{".env", ".env.production", "private.pem", "private.key", "certificate.p12"} {
		root := t.TempDir()
		planWrite(t, root, name, `{"hello":"private material"}`)
		_, err := BuildPlan(root, ".internationalizer.yml", PlanOptions{TargetLocales: []string{"fr"}, AddBundles: []config.Bundle{{ID: "web", Source: name, Target: "{locale}.json", Format: "json", MessageSyntax: message.Plain}}})
		assertPlanCode(t, err, "unsafe_path")
	}
}

func TestPlanAllowsIntrinsicFluentGrammar(t *testing.T) {
	root := t.TempDir()
	planWrite(t, root, "en.ftl", "hello = Hello\n")
	p, err := BuildPlan(root, ".internationalizer.yml", PlanOptions{TargetLocales: []string{"fr"}, AddBundles: []config.Bundle{{ID: "fluent", Source: "en.ftl", Target: "{locale}.ftl", MessageSyntax: message.Auto}}})
	if err != nil || len(p.RequiredDecisions) > 0 {
		t.Fatalf("Fluent plan: %+v %v", p, err)
	}
	if _, err = ApplyPlan(p); err != nil {
		t.Fatal(err)
	}
}

func TestApplyFailurePreservesConfigAndReleasesOwnedLock(t *testing.T) {
	for _, drift := range []string{"source", "config", "runtime"} {
		t.Run(drift, func(t *testing.T) {
			root := planFixture(t)
			p, err := BuildPlan(root, "", planOptions())
			if err != nil {
				t.Fatal(err)
			}
			switch drift {
			case "source":
				planWrite(t, root, "tmp/en.json", `{"hello":"Changed"}`)
			case "config":
				planWrite(t, root, ".internationalizer.yml", "# concurrent edit\nsource_locale: en\ntarget_locales: [fr]\nsource_path: tmp/en.json\n")
			case "runtime":
				planWrite(t, root, "web/package.json", `{"dependencies":{"i18next-icu":"1"}}`)
			}
			configPath := filepath.Join(root, ".internationalizer.yml")
			before, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatal(err)
			}
			receipt, err := ApplyPlan(p)
			assertPlanCode(t, err, "stale_plan")
			if receipt != nil {
				t.Fatalf("failed application returned a success receipt: %+v", receipt)
			}
			after, err := os.ReadFile(configPath)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("failed application changed existing config: %v", err)
			}
			if _, err := os.Lstat(configPath + ".apply-lock"); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("owned lock remains after failure: %v", err)
			}
			artifacts, err := filepath.Glob(filepath.Join(root, ".internationalizer-apply-*"))
			if err != nil || len(artifacts) != 0 {
				t.Fatalf("temporary replacements remain: %v, %v", artifacts, err)
			}
		})
	}
}

func TestApplyPreservesPermissionsAndAttestsExactReadback(t *testing.T) {
	root := planFixture(t)
	configPath := filepath.Join(root, ".internationalizer.yml")
	if err := os.Chmod(configPath, 0640); err != nil {
		t.Fatal(err)
	}
	p, err := BuildPlan(root, "", planOptions())
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := ApplyPlan(p)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != p.ProposedYAML || planHash(data) != receipt.ConfigSHA256 || receipt.PlanID != p.ID || !receipt.ObservationsRevalidated {
		t.Fatalf("receipt does not attest committed bytes and evidence: %+v", receipt)
	}
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0640 {
		t.Fatalf("existing config permissions changed: %v", info.Mode().Perm())
	}
	if _, err := os.Lstat(configPath + ".apply-lock"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("owned lock remains after success: %v", err)
	}
}

func assertPlanCode(t *testing.T, err error, want string) {
	t.Helper()
	var e *PlanError
	if !errors.As(err, &e) || e.Code != want {
		t.Fatalf("got %v, want code %s", err, want)
	}
}
