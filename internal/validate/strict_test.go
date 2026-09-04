package validate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/Tom-R-Main/Internationalizer/internal/policy"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

func TestStrictValidationSeparatesStructuralAndTranslatedCoverage(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"translated": "Save", "seeded": "Cancel"}, map[string]string{"translated": "Enregistrer", "seeded": "Cancel"})

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	report := reports[0]
	if report.Coverage != 100 || report.StructuralCoverage != 100 || report.TranslatedCoverage == nil || *report.TranslatedCoverage != 50 {
		t.Fatalf("coverage = legacy %.1f, structural %.1f, translated %v", report.Coverage, report.StructuralCoverage, report.TranslatedCoverage)
	}
	assertFindingCodes(t, report, CodeSourceIdentical)
	if !HasFailures(reports) {
		t.Fatal("strict source-identical translation did not fail")
	}
}

func TestStrictValidationAllowsExactSameSourceGlossaryTerm(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"brand": "Lens"}, map[string]string{"brand": "Lens"})
	if err := glossary.Save(cfg.GlossaryDir, "fr", []glossary.Term{{Source: "Lens", Target: "Lens", WholeWord: true}}); err != nil {
		t.Fatal(err)
	}

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if HasFailures(reports) || reports[0].TranslatedCoverage == nil || *reports[0].TranslatedCoverage != 100 {
		t.Fatalf("exact glossary exemption failed: %#v", reports[0])
	}
}

func TestStrictValidationAllowsNonLinguisticValues(t *testing.T) {
	cfg := validationConfig(t,
		map[string]string{"count": "123", "separator": "---"},
		map[string]string{"count": "123", "separator": "---"},
	)

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if HasFailures(reports) || reports[0].TranslatedCoverage == nil || *reports[0].TranslatedCoverage != 100 {
		t.Fatalf("non-linguistic values failed strict validation: %#v", reports[0])
	}
}

func TestExtraKeyWarnsByDefaultAndFailsStrictOrRequiredState(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"save": "Save"}, map[string]string{"save": "Enregistrer", "extra": "Supplémentaire"})

	legacy, err := Validate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if HasFailures(legacy) {
		t.Fatalf("legacy extra key became a failure: %#v", legacy[0])
	}
	if findingByCode(legacy[0], CodeExtraKey).Severity != SeverityWarning {
		t.Fatalf("legacy extra finding = %#v", findingByCode(legacy[0], CodeExtraKey))
	}

	strict, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if !HasFailures(strict) || findingByCode(strict[0], CodeExtraKey).Severity != SeverityError {
		t.Fatalf("strict extra key did not fail: %#v", strict[0])
	}

	requiredState, err := ValidateWithOptions(cfg, Options{RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	if !HasFailures(requiredState) || findingByCode(requiredState[0], CodeExtraKey).Severity != SeverityError {
		t.Fatalf("required-state extra key did not fail: %#v", requiredState[0])
	}
}

func TestStrictValidationReportsProtectedAndGlossaryFailures(t *testing.T) {
	cfg := validationConfig(t,
		map[string]string{
			"html":        "Click <strong>here</strong>",
			"placeholder": "Hello, {{name}}",
			"term":        "Open the Dashboard",
		},
		map[string]string{
			"html":        "Cliquez ici",
			"placeholder": "Bonjour",
			"term":        "Ouvrez le tableau",
		},
	)
	if err := glossary.Save(cfg.GlossaryDir, "fr", []glossary.Term{{Source: "Dashboard", Target: "Tableau de bord", WholeWord: true}}); err != nil {
		t.Fatal(err)
	}

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	assertFindingCodes(t, reports[0], CodeGlossaryViolation, CodeProtectedStructureMismatch)
	if len(reports[0].Mismatches) != 1 || reports[0].Mismatches[0].Key != "placeholder" {
		t.Fatalf("legacy mismatch field changed: %#v", reports[0].Mismatches)
	}
}

func TestWarningGlossaryFindingDoesNotFailOtherwiseValidStrictReport(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"term": "Open Dashboard"}, map[string]string{"term": "Ouvrir tableau"})
	if err := glossary.Save(cfg.GlossaryDir, "fr", []glossary.Term{{Source: "Dashboard", Target: "Tableau de bord", Enforcement: glossary.EnforcementWarning}}); err != nil {
		t.Fatal(err)
	}

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	finding := findingByCode(reports[0], CodeGlossaryViolation)
	if finding.Severity != SeverityWarning || HasFailures(reports) {
		t.Fatalf("warning glossary result = %#v, report = %#v", finding, reports[0])
	}
}

func TestGlossaryIgnoreCaseUsesUnicodeCaseFolding(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"term": "ΟΣ"}, map[string]string{"term": "ος"})
	if err := glossary.Save(cfg.GlossaryDir, "fr", []glossary.Term{{Source: "οσ", Target: "οσ", IgnoreCase: true, WholeWord: true}}); err != nil {
		t.Fatal(err)
	}

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if findingByCode(reports[0], CodeGlossaryViolation).Code != "" {
		t.Fatalf("Unicode case-folded glossary term failed: %#v", reports[0])
	}
}

func TestStrictValidationRejectsBlankTranslation(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"save": "Save"}, map[string]string{"save": ""})

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	assertFindingCodes(t, reports[0], CodeBlankTranslation)
	if !HasFailures(reports) || reports[0].TranslatedCoverage == nil || *reports[0].TranslatedCoverage != 0 {
		t.Fatalf("blank translation passed strict validation: %#v", reports[0])
	}
}

func TestI18nextPluralValidationAllowsRequiredTargetForms(t *testing.T) {
	cfg := validationConfig(t,
		map[string]string{"items_one": "{{count}} item", "items_other": "{{count}} items"},
		map[string]string{"items_one": "{{count}} article", "items_many": "{{count}} articles", "items_other": "{{count}} articles"},
	)
	cfg.Validation.PluralStyle = "i18next-v4"

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if HasFailures(reports) || len(reports[0].Extra) != 0 {
		t.Fatalf("valid French plural forms failed: %#v", reports[0])
	}

	writeJSON(t, targetPath(cfg), map[string]string{"items_one": "{{count}} article", "items_other": "{{count}} articles"})
	reports, err = ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	assertFindingCodes(t, reports[0], CodePluralFormMissing)
	if findingByCode(reports[0], CodePluralFormMissing).Key != "items_many" {
		t.Fatalf("plural finding = %#v", findingByCode(reports[0], CodePluralFormMissing))
	}
}

func TestTargetOnlyPluralFormsReceiveContentAndStateValidation(t *testing.T) {
	cfg := validationConfig(t,
		map[string]string{"items_one": "{{count}} item", "items_other": "{{count}} items"},
		map[string]string{"items_one": "{{count}} article", "items_many": "", "items_other": "{{count}} articles"},
	)
	cfg.Validation.PluralStyle = "i18next-v4"

	reports, err := ValidateWithOptions(cfg, Options{Strict: true, RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	var pluralFindings []FindingCode
	for _, finding := range reports[0].Findings {
		if finding.Key == "items_many" {
			pluralFindings = append(pluralFindings, finding.Code)
		}
	}
	for _, want := range []FindingCode{CodeBlankTranslation, CodeProtectedStructureMismatch, CodeUntracked} {
		found := false
		for _, got := range pluralFindings {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("items_many findings = %v, missing %q", pluralFindings, want)
		}
	}
}

func TestPrivateUseLocaleUsesRootPluralRequirements(t *testing.T) {
	cfg := validationConfig(t,
		map[string]string{"items_one": "{{count}} item", "items_other": "{{count}} items"},
		map[string]string{"items_one": "{{count}} x", "items_other": "{{count}} xs"},
	)
	oldTarget := targetPath(cfg)
	cfg.TargetLocales = []string{"qaa-ZZ"}
	newTarget := filepath.Join(filepath.Dir(oldTarget), "qaa-ZZ.json")
	if err := os.Rename(oldTarget, newTarget); err != nil {
		t.Fatal(err)
	}
	cfg.Validation.PluralStyle = "i18next-v4"

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if findingByCode(reports[0], CodePluralFormMissing).Code != "" {
		t.Fatalf("private-use locale received an unsupported plural finding: %#v", reports[0])
	}

	writeJSON(t, newTarget, map[string]string{"items_other": "{{count}} xs"})
	reports, err = ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if HasFailures(reports) {
		t.Fatalf("root plural fallback rejected an other-only target: %#v", reports[0])
	}
}

func TestI18nextPluralCoverageUsesTargetLocaleRequirements(t *testing.T) {
	cfg := validationConfig(t,
		map[string]string{"items_one": "{{count}} item", "items_other": "{{count}} items"},
		map[string]string{"items_other": "{{count}} 個"},
	)
	oldTarget := targetPath(cfg)
	cfg.TargetLocales = []string{"ja"}
	newTarget := filepath.Join(filepath.Dir(oldTarget), "ja.json")
	if err := os.Rename(oldTarget, newTarget); err != nil {
		t.Fatal(err)
	}
	cfg.Validation.PluralStyle = "i18next-v4"

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if HasFailures(reports) || reports[0].StructuralCoverage != 100 || reports[0].TranslatedCoverage == nil || *reports[0].TranslatedCoverage != 100 {
		t.Fatalf("valid Japanese plural coverage = %#v", reports[0])
	}
}

func TestRequireStateReportsAndClearsBoundProvenanceFindings(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"save": "Save"}, map[string]string{"save": "Enregistrer"})
	reports, err := ValidateWithOptions(cfg, Options{RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	assertFindingCodes(t, reports[0], CodeUntracked)

	resolved, err := policy.Resolve(cfg, "fr", "json", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	manifest := state.New()
	manifest.Set(state.Entry{
		Bundle:     "default",
		Key:        "save",
		Locale:     "fr",
		SourceHash: state.SourceHash("json", "Save"),
		PolicyHash: resolved.Hash,
		TargetHash: state.TargetHash("Enregistrer"),
		UpdatedAt:  time.Now().UTC(),
	})
	if err := manifest.Save(cfg.ManifestPath); err != nil {
		t.Fatal(err)
	}

	reports, err = ValidateWithOptions(cfg, Options{RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	if HasFailures(reports) {
		t.Fatalf("matching state failed: %#v", reports[0])
	}

	writeJSON(t, cfg.SourcePath, map[string]string{"save": "Save now"})
	writeJSON(t, targetPath(cfg), map[string]string{"save": "Sauvegarder"})
	cfg.LLM.Model = "test-v2"
	reports, err = ValidateWithOptions(cfg, Options{RequireState: true})
	if err != nil {
		t.Fatal(err)
	}
	assertFindingCodes(t, reports[0], CodePolicyStale, CodeSourceStale, CodeTargetModified)
}

func TestRequireStateRejectsMalformedManifest(t *testing.T) {
	cfg := validationConfig(t, map[string]string{"save": "Save"}, map[string]string{"save": "Enregistrer"})
	if err := os.WriteFile(cfg.ManifestPath, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateWithOptions(cfg, Options{RequireState: true}); err == nil {
		t.Fatal("malformed manifest was accepted")
	}
}

func TestProtectedFindingsUseStableCodeForEveryStructure(t *testing.T) {
	tests := map[string][2]string{
		"interpolation": {"Hello {{name}}", "Bonjour"},
		"html":          {"<strong>Save</strong>", "Enregistrer"},
		"inline code":   {"Run `go test`", "Lancez `go test ./...`"},
		"fenced code":   {"```go\nfmt.Println(1)\n```\n", "```go\nfmt.Println(2)\n```\n"},
		"link":          {"[Guide](https://example.com/a)", "[Guide](https://example.com/b)"},
	}
	for name, values := range tests {
		t.Run(name, func(t *testing.T) {
			findings := ProtectedFindings("key", values[0], values[1], "fr")
			if len(findings) == 0 || findings[0].Code != CodeProtectedStructureMismatch {
				t.Fatalf("findings = %#v", findings)
			}
		})
	}
}

func TestProtectedDocumentFindingsResolveRelativeLinksAndAllowSourceBacklink(t *testing.T) {
	source := `<img src="assets/logo.svg"> [License](LICENSE) [Guide](docs/guide.md)`
	target := `[English](../../README.md) <img src="../../assets/logo.svg"> [Licence](../../LICENSE) [Guide](../guide.md)`
	findings := ProtectedDocumentFindings(
		"markdown:preamble",
		source,
		target,
		"fr",
		"/project/README.md",
		"/project/docs/i18n/fr.md",
	)
	if len(findings) != 0 {
		t.Fatalf("equivalent document links were rejected: %#v", findings)
	}
}

func TestProtectedDocumentFindingsRejectDifferentResolvedLink(t *testing.T) {
	source := `[Guide](docs/guide.md)`
	target := `[Guide](../other.md)`
	findings := ProtectedDocumentFindings(
		"markdown:preamble",
		source,
		target,
		"fr",
		"/project/README.md",
		"/project/docs/i18n/fr.md",
	)
	if len(findings) != 1 || findings[0].Code != CodeProtectedStructureMismatch {
		t.Fatalf("different document link was accepted: %#v", findings)
	}
}

func TestProtectedFindingsCompareICUStructuresPerLocaleBranch(t *testing.T) {
	tests := map[string]string{
		"html":        "<strong>Save</strong>",
		"inline code": "`go test`",
		"fenced code": "```go\ngo test\n```\n",
		"link":        "[Guide](https://example.com/guide)",
	}
	for name, protected := range tests {
		t.Run(name, func(t *testing.T) {
			source := "{n, plural, one {" + protected + "} other {" + protected + "}}"
			target := "{n, plural, one {" + protected + "} few {" + protected + "} many {" + protected + "} other {" + protected + "}}"
			if findings := ProtectedFindings("key", source, target, "ru"); len(findings) != 0 {
				t.Fatalf("valid target-only ICU branches produced findings: %#v", findings)
			}

			damaged := "{n, plural, one {" + protected + "} few {damaged} many {" + protected + "} other {" + protected + "}}"
			if findings := ProtectedFindings("key", source, damaged, "ru"); len(findings) == 0 {
				t.Fatal("damaged target-only ICU branch was accepted")
			}
		})
	}
}

func TestProtectedFindingsCarryContextIntoNestedICUBranches(t *testing.T) {
	argument := "{os, select, win {dir} other {ls}}"
	damagedArgument := "{os, select, win {répertoire} other {ls}}"
	tests := map[string][2]string{
		"html attribute": {`<span data-command="` + argument + `">Run</span>`, `<span data-command="` + damagedArgument + `">Exécuter</span>`},
		"inline code":    {"`" + argument + "`", "`" + damagedArgument + "`"},
		"fenced code":    {"```sh\n" + argument + "\n```\n", "```sh\n" + damagedArgument + "\n```\n"},
		"link":           {"[Run](" + argument + ")", "[Exécuter](" + damagedArgument + ")"},
	}
	for name, values := range tests {
		t.Run(name, func(t *testing.T) {
			if findings := ProtectedFindings("key", values[0], values[0], "fr"); len(findings) != 0 {
				t.Fatalf("unchanged nested ICU protection produced findings: %#v", findings)
			}
			if findings := ProtectedFindings("key", values[0], values[1], "fr"); len(findings) == 0 {
				t.Fatal("damaged nested ICU protected content was accepted")
			}
		})
	}
}

func TestEvaluationCorpusSchemaAndIDsAreStable(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "test", "evaluation", "v1", "cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		SchemaVersion int `json:"schema_version"`
		Cases         []struct {
			ID   string `json:"id"`
			Kind string `json:"kind"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	if corpus.SchemaVersion != 1 || len(corpus.Cases) < 10 {
		t.Fatalf("corpus metadata = %#v", corpus)
	}
	seen := make(map[string]struct{}, len(corpus.Cases))
	for _, testCase := range corpus.Cases {
		if testCase.ID == "" || testCase.Kind == "" {
			t.Fatalf("invalid case metadata: %#v", testCase)
		}
		if _, exists := seen[testCase.ID]; exists {
			t.Fatalf("duplicate evaluation case id %q", testCase.ID)
		}
		seen[testCase.ID] = struct{}{}
	}
}

func TestStrictEvaluationProjectExposesBaselineBlindSpots(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	t.Chdir(repositoryRoot)
	cfg, err := config.Load(filepath.Join("test", "evaluation", "v1", "project", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}

	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 {
		t.Fatalf("reports = %#v", reports)
	}
	report := reports[0]
	if report.StructuralCoverage < 72.7 || report.StructuralCoverage > 72.8 || report.TranslatedCoverage == nil || *report.TranslatedCoverage < 63.6 || *report.TranslatedCoverage > 63.7 {
		t.Fatalf("evaluation coverage = structural %.1f, translated %v", report.StructuralCoverage, report.TranslatedCoverage)
	}
	for _, code := range []FindingCode{
		CodeMissingKey,
		CodeExtraKey,
		CodeSourceIdentical,
		CodeProtectedStructureMismatch,
		CodeGlossaryViolation,
		CodePluralFormMissing,
	} {
		if findingByCode(report, code).Code == "" {
			t.Fatalf("evaluation findings = %#v, missing %q", report.Findings, code)
		}
	}
}

func TestEvaluationValidationCasesMatchStrictFindings(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	t.Chdir(repositoryRoot)
	data, err := os.ReadFile(filepath.Join("test", "evaluation", "v1", "cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Cases []struct {
			ID                   string        `json:"id"`
			Kind                 string        `json:"kind"`
			ExpectedFindingCodes []FindingCode `json:"expected_finding_codes"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(filepath.Join("test", "evaluation", "v1", "project", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	reports, err := ValidateWithOptions(cfg, Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	keyByCase := map[string]string{
		"missing-key":           "missing",
		"unexpected-target-key": "extra",
		"english-seed":          "seeded",
		"allowed-product-name":  "brand",
		"placeholder-damage":    "placeholder",
		"html-damage":           "markup",
		"glossary-drift":        "term",
		"target-plural-gap":     "items_many",
	}
	for _, testCase := range corpus.Cases {
		if testCase.Kind != "validation" {
			continue
		}
		key, ok := keyByCase[testCase.ID]
		if !ok {
			t.Fatalf("validation case %q has no executable project mapping", testCase.ID)
		}
		var got []FindingCode
		for _, finding := range reports[0].Findings {
			if finding.Key == key {
				got = append(got, finding.Code)
			}
		}
		sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
		want := append([]FindingCode(nil), testCase.ExpectedFindingCodes...)
		sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("case %q findings = %v, want %v", testCase.ID, got, want)
		}
	}
}

func validationConfig(t *testing.T, source, target map[string]string) *config.Config {
	t.Helper()
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	writeJSON(t, sourcePath, source)
	writeJSON(t, filepath.Join(dir, "fr.json"), target)
	return &config.Config{
		SourceLocale:   "en",
		TargetLocales:  []string{"fr"},
		SourcePath:     sourcePath,
		GlossaryDir:    filepath.Join(dir, "glossary"),
		StyleGuidesDir: filepath.Join(dir, "style-guides"),
		ManifestPath:   filepath.Join(dir, "manifest.json"),
		LLM:            config.LLM{Provider: "test", Model: "test-v1"},
	}
}

func targetPath(cfg *config.Config) string {
	path, _ := cfg.EffectiveBundles()[0].TargetPath(cfg.TargetLocales[0])
	return path
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func findingByCode(report Report, code FindingCode) Finding {
	for _, finding := range report.Findings {
		if finding.Code == code {
			return finding
		}
	}
	return Finding{}
}

func assertFindingCodes(t *testing.T, report Report, expected ...FindingCode) {
	t.Helper()
	gotSet := make(map[FindingCode]struct{})
	for _, finding := range report.Findings {
		gotSet[finding.Code] = struct{}{}
	}
	for _, code := range expected {
		if _, ok := gotSet[code]; !ok {
			t.Fatalf("findings = %#v, missing code %q", report.Findings, code)
		}
	}
}

func TestFindingsSortDeterministicallyByKeyThenCode(t *testing.T) {
	findings := []Finding{
		{Key: "z", Code: CodeSourceIdentical, Message: "z"},
		{Key: "a", Code: CodeSourceIdentical, Message: "b"},
		{Key: "a", Code: CodeExtraKey, Message: "a"},
	}
	sortFindings(findings)
	want := []Finding{
		{Key: "a", Code: CodeExtraKey, Message: "a"},
		{Key: "a", Code: CodeSourceIdentical, Message: "b"},
		{Key: "z", Code: CodeSourceIdentical, Message: "z"},
	}
	if !reflect.DeepEqual(findings, want) {
		t.Fatalf("findings = %#v, want %#v", findings, want)
	}
}
