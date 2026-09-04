package onboarding

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/message"
)

func discoveryFile(t *testing.T, root, path, content string) {
	t.Helper()
	path = filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestScanMonorepo(t *testing.T) {
	root := t.TempDir()
	t.Setenv("ONBOARDING_TEST_KEY", "sensitive-value-must-not-appear")
	discoveryFile(t, root, ".internationalizer.yml", `source_locale: en
target_locales: [fr, ja]
source_path: tmp/english-keys.json
llm:
  provider: openai
  api_key_env: ONBOARDING_TEST_KEY
  locale_overrides:
    ja:
      provider: gemini
      api_key_env: ONBOARDING_TEST_MISSING_KEY
`)
	discoveryFile(t, root, "tmp/english-keys.json", `{"docs.tui.skills.desc":"<code>&lt;root&gt;/{.sift,.claude,.codex,.agents}/skills</code>"}`)
	discoveryFile(t, root, "web/package.json", `{"dependencies":{"i18next":"1","react-i18next":"1"}}`)
	discoveryFile(t, root, "web/src/i18n/locales/en.json", `{"welcome":"Hello {{name}}"}`)
	discoveryFile(t, root, "node_modules/fake/en.json", `{"bad":"must not be found"}`)
	inspection, err := Scan(root, ".internationalizer.yml")
	if err != nil {
		t.Fatal(err)
	}
	if !inspection.ConfigExists || len(inspection.Candidates) != 2 || len(inspection.Bundles) != 1 {
		t.Fatalf("unexpected inspection: %+v", inspection)
	}
	bundle := inspection.Bundles[0]
	if bundle.ID != "default" || bundle.MessageSyntax != message.Auto || bundle.Provenance["source"] != "source_path" {
		t.Fatalf("bundle = %+v", bundle)
	}
	if !filepath.IsAbs(bundle.Source) || !strings.HasSuffix(bundle.Targets["ja"], filepath.Join("tmp", "ja.json")) {
		t.Fatalf("paths = %+v", bundle)
	}
	var web Candidate
	for _, candidate := range inspection.Candidates {
		if strings.HasPrefix(candidate.Source, "web/") {
			web = candidate
		}
	}
	if web.Framework != "i18next" || web.SuggestedSyntax != message.I18next || len(web.ConfiguredBundles) != 0 || len(web.Evidence) < 2 {
		t.Fatalf("web = %+v", web)
	}
	counts := map[string]int{}
	for _, diagnostic := range inspection.Diagnostics {
		counts[diagnostic.Code]++
	}
	for _, code := range []string{"SOURCE_CONFIRMATION_REQUIRED", "AUTO_SYNTAX_AMBIGUOUS", "UNCOVERED_CATALOG"} {
		if counts[code] != 1 {
			t.Fatalf("%s count = %d; diagnostics = %+v", code, counts[code], inspection.Diagnostics)
		}
	}
	if !inspection.Credentials[0].Present || inspection.Credentials[1].Present || inspection.Credentials[0].ProviderVerified {
		t.Fatalf("credentials = %+v", inspection.Credentials)
	}
	encoded, err := json.Marshal(inspection)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "sensitive-value") {
		t.Fatal("credential value leaked")
	}
}

func TestScanICUIntegrationAndNamespace(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr]\nsource_path: app/locales/en/common.json\nmessage_syntax: icu\n")
	discoveryFile(t, root, "app/package.json", `{"dependencies":{"i18next":"1","i18next-icu":"1"}}`)
	discoveryFile(t, root, "app/locales/en/common.json", `{"count":"{count, plural, one {One} other {Many}}"}`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(inspection.Candidates) != 1 {
		t.Fatalf("candidates = %+v", inspection.Candidates)
	}
	candidate := inspection.Candidates[0]
	if candidate.Target != "app/locales/{locale}/common.json" || candidate.SuggestedSyntax != "" || !strings.Contains(candidate.Uncertainty, "plugin registration") {
		t.Fatalf("candidate = %+v", candidate)
	}
	if inspection.Bundles[0].MessageSyntax != message.ICU {
		t.Fatal("explicit ICU changed")
	}
}

func TestScanNoConfigAndRuntimeReference(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, "package.json", `{"dependencies":{"i18next":"1"}}`)
	discoveryFile(t, root, "src/i18n.ts", `import ICU from 'i18next-icu';`)
	discoveryFile(t, root, "locales/en.json", `{"hello":"Hello"}`)
	inspection, err := Scan(root, "absent.yml")
	if err != nil {
		t.Fatal(err)
	}
	if inspection.ConfigExists || inspection.SourceLocale != "en" || len(inspection.Candidates) != 1 {
		t.Fatalf("inspection = %+v", inspection)
	}
	if inspection.Candidates[0].SuggestedSyntax != "" {
		t.Fatal("runtime ICU reference ignored")
	}
}

func TestScanInvalidConfigDoesNotEchoValues(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, "broken.yml", "batch_size: super-secret-invalid-value\n")
	inspection, err := Scan(root, "broken.yml")
	if err == nil || strings.Contains(err.Error(), "super-secret") {
		t.Fatalf("error = %v", err)
	}
	if len(inspection.Diagnostics) != 1 || inspection.Diagnostics[0].Code != "CONFIG_INVALID" {
		t.Fatalf("inspection = %+v", inspection)
	}
}

func TestScanSymlinksExcluded(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr]\nsource_path: linked/en.json\n")
	other := t.TempDir()
	discoveryFile(t, other, "en.json", `{"hello":"Hello"}`)
	if err := os.Symlink(other, filepath.Join(root, "linked")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	var unreadable bool
	for _, diagnostic := range inspection.Diagnostics {
		unreadable = unreadable || diagnostic.Code == "SOURCE_UNREADABLE"
	}
	if !unreadable {
		t.Fatal("configured symlink should not be read")
	}
}

func TestScanExplicitICUStrict(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "source_locale: en\ntarget_locales: [fr, ja]\nsource_path: en.json\nmessage_syntax: icu\n")
	discoveryFile(t, root, "en.json", `{"bad":"{oops, invalid}"}`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, diagnostic := range inspection.Diagnostics {
		if diagnostic.Code == "SOURCE_SYNTAX_INVALID" {
			count++
		}
		if diagnostic.Code == "AUTO_SYNTAX_AMBIGUOUS" {
			t.Fatal("explicit ICU downgraded")
		}
	}
	if count != 1 {
		t.Fatalf("diagnostics = %+v", inspection.Diagnostics)
	}
}

func TestScanAutoCodeBracesWarnEvenWhenICUParses(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "target_locales: [fr]\nsource_path: en.json\n")
	discoveryFile(t, root, "en.json", `{"code":"Run <code>{command}</code>"}`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range inspection.Diagnostics {
		if diagnostic.Code == "AUTO_SYNTAX_AMBIGUOUS" {
			if diagnostic.Severity != "warning" || diagnostic.Key != "code" {
				t.Fatalf("diagnostic = %+v", diagnostic)
			}
			return
		}
	}
	t.Fatalf("missing ambiguity warning: %+v", inspection.Diagnostics)
}

func TestScanUsesRootNotConfigDirectory(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, "configuration/project.yml", "target_locales: [fr]\nmessage_syntax: plain\nbundles:\n  - id: web\n    source: catalogs/en.json\n    target: catalogs/{locale}.json\n    message_syntax: i18next\n")
	discoveryFile(t, root, "catalogs/en.json", `{"hello":"Hello {{name}}"}`)
	inspection, err := Scan(root, "configuration/project.yml")
	if err != nil {
		t.Fatal(err)
	}
	if len(inspection.Bundles) != 1 {
		t.Fatalf("bundles = %+v", inspection.Bundles)
	}
	bundle := inspection.Bundles[0]
	if bundle.Source != filepath.Join(inspection.Root, "catalogs", "en.json") || bundle.Provenance["message_syntax"] != "bundle.message_syntax" || bundle.MessageSyntax != message.I18next {
		t.Fatalf("bundle = %+v", bundle)
	}
	for _, diagnostic := range inspection.Diagnostics {
		if diagnostic.Severity == "error" {
			t.Fatalf("unexpected diagnostic = %+v", diagnostic)
		}
	}
}

func TestScanRejectsSecretShapedSource(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "target_locales: [fr]\nbundles:\n  - id: bad\n    source: .env.json\n    target: fr.json\n    format: json\n")
	discoveryFile(t, root, ".env.json", `{"TOKEN":"must-not-be-read"}`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range inspection.Diagnostics {
		if diagnostic.Code == "SOURCE_UNREADABLE" {
			return
		}
	}
	t.Fatal("secret-shaped configured source was not rejected")
}

func TestScanInvalidRoot(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, "file", "not a directory")
	if _, err := Scan(filepath.Join(root, "file"), ""); err == nil {
		t.Fatal("expected non-directory error")
	}
	if _, err := Scan(filepath.Join(root, "missing"), ""); err == nil {
		t.Fatal("expected missing root error")
	}
}

func TestScanExcludesConfiguredSupportDirectoriesButKeepsExplicitBundles(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "target_locales: [fr]\nstyle_guides_dir: docs/i18n-style-guides\nglossary_dir: terms\nbundles:\n  - id: intentional-guide\n    source: docs/i18n-style-guides/published/en.md\n    target: docs/i18n-style-guides/published/{locale}.md\n")
	discoveryFile(t, root, "docs/i18n-style-guides/en.md", "# English style guide\n")
	discoveryFile(t, root, "docs/i18n-style-guides/microsoft/en.md", "# Microsoft English guide\n")
	discoveryFile(t, root, "docs/i18n-style-guides/published/en.md", "# Intentionally translated guide\n")
	discoveryFile(t, root, "terms/en.json", `{"term":"Definition"}`)
	discoveryFile(t, root, "terms-extra/en.json", `{"hello":"Hello"}`)
	discoveryFile(t, root, "app/locales/en.json", `{"hello":"Hello"}`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(inspection.Candidates) != 3 {
		t.Fatalf("candidates = %+v", inspection.Candidates)
	}
	expected := map[string]bool{"docs/i18n-style-guides/published/en.md": true, "terms-extra/en.json": true, "app/locales/en.json": true}
	for _, candidate := range inspection.Candidates {
		if !expected[candidate.Source] {
			t.Fatalf("unexpected inferred support catalog = %+v", candidate)
		}
		if candidate.Source == "docs/i18n-style-guides/published/en.md" && (len(candidate.ConfiguredBundles) != 1 || candidate.ConfiguredBundles[0] != "intentional-guide") {
			t.Fatalf("explicit bundle missing = %+v", candidate)
		}
	}
}

func TestScanConfigurationDecisionsHaveSafeRecovery(t *testing.T) {
	root := t.TempDir()
	discoveryFile(t, root, ".internationalizer.yml", "target_locales: [fr]\nsource_path: tmp/english-keys.json\n")
	discoveryFile(t, root, "tmp/english-keys.json", `{"code":"<code>{a,b,c}</code>"}`)
	discoveryFile(t, root, "web/locales/en.json", `{"hello":"Hello"}`)
	inspection, err := Scan(root, "")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, diagnostic := range inspection.Diagnostics {
		if len(diagnostic.RequiredDecisions) == 0 {
			continue
		}
		count++
		if len(diagnostic.Recovery) != 1 {
			t.Fatalf("missing recovery = %+v", diagnostic)
		}
		recovery := diagnostic.Recovery[0]
		if strings.Join(recovery.Argv, " ") != "internationalizer config plan --help" || len(recovery.SideEffects) != 0 || strings.Join(recovery.RequiredDecisions, ",") != strings.Join(diagnostic.RequiredDecisions, ",") {
			t.Fatalf("unsafe or incomplete recovery = %+v", diagnostic)
		}
	}
	if count != 3 {
		t.Fatalf("expected 3 decisions, got %d: %+v", count, inspection.Diagnostics)
	}
}

func TestScanOtherFrameworksAndMixedRuntimeEvidence(t *testing.T) {
	for _, tc := range []struct {
		name, dependencies, framework, uncertainty string
		syntax                                     message.Syntax
	}{
		{name: "next-intl", dependencies: `"next-intl":"1"`, framework: "next-intl", syntax: message.ICU, uncertainty: "static"},
		{name: "vue-i18n", dependencies: `"vue-i18n":"1"`, framework: "vue-i18n", uncertainty: "compatibility"},
		{name: "mixed-next-i18next", dependencies: `"next-intl":"1","i18next":"1"`, framework: "multiple", uncertainty: "multiple frameworks"},
		{name: "mixed-vue-next", dependencies: `"vue-i18n":"1","next-intl":"1"`, framework: "multiple", uncertainty: "multiple frameworks"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			discoveryFile(t, root, "app/package.json", `{"dependencies":{`+tc.dependencies+`}}`)
			discoveryFile(t, root, "app/messages/en.json", `{"hello":"Hello"}`)
			inspection, err := Scan(root, "absent.yml")
			if err != nil {
				t.Fatal(err)
			}
			if len(inspection.Candidates) != 1 {
				t.Fatalf("candidates = %+v", inspection.Candidates)
			}
			candidate := inspection.Candidates[0]
			if candidate.Framework != tc.framework || candidate.SuggestedSyntax != tc.syntax || !strings.Contains(candidate.Uncertainty, tc.uncertainty) {
				t.Fatalf("candidate = %+v", candidate)
			}
		})
	}
}
