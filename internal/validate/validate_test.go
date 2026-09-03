package validate

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
)

func TestValidateUsesBundleTargetTemplate(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "README.md")
	targetPath := filepath.Join(dir, "docs", "i18n", "fr.md")
	if err := os.WriteFile(sourcePath, []byte("# Hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte("# Bonjour\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		SourceLocale:  "en",
		TargetLocales: []string{"fr"},
		Bundles: []config.Bundle{{
			ID:     "docs",
			Source: sourcePath,
			Target: filepath.Join(dir, "docs", "i18n", "{locale}.md"),
		}},
	}

	reports, err := Validate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 || reports[0].Bundle != "docs" || reports[0].Coverage != 100 {
		t.Fatalf("reports = %#v", reports)
	}
	if HasFailures(reports) {
		t.Fatalf("valid translated document reported as failure: %#v", reports)
	}
}

func TestValidateReportsMalformedTargetAsFailure(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, "fr.json")
	if err := os.WriteFile(sourcePath, []byte(`{"a":"A"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		TargetLocales: []string{"fr"},
		SourcePath:    sourcePath,
	}
	reports, err := Validate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 || len(reports[0].Errors) != 1 || !HasFailures(reports) {
		t.Fatalf("reports = %#v, want malformed target failure", reports)
	}
}

func TestFormatHumanPreservesInterpolationMismatchDetails(t *testing.T) {
	reports := []Report{{
		Bundle: "app",
		Locale: "fr",
		Mismatches: []Mismatch{{
			Key:        "welcome",
			SourceVars: []string{"name"},
			TargetVars: nil,
		}},
	}}

	got := FormatHuman(reports)
	if !strings.Contains(got, "mismatch: welcome (source: [name], target: [])") {
		t.Fatalf("FormatHuman omitted mismatch details: %q", got)
	}
}

func TestValidateReportsICUStructureFailures(t *testing.T) {
	tests := []struct {
		name   string
		locale string
		source string
		target string
		code   FindingCode
	}{
		{
			name:   "malformed target",
			locale: "fr",
			source: `{count, plural, one {One item} other {# items}}`,
			target: `{count, plural, one {Un article} other {# articles}`,
			code:   CodeICUMessageSyntax,
		},
		{
			name:   "argument type",
			locale: "fr",
			source: `{count, plural, one {One item} other {# items}}`,
			target: `{count, select, one {Un article} other {Des articles}}`,
			code:   CodeICUArgumentMismatch,
		},
		{
			name:   "locale category",
			locale: "te",
			source: `{count, plural, one {One item} other {# items}}`,
			target: `{count, plural, one {ఒక అంశం} few {కొన్ని అంశాలు} other {# అంశాలు}}`,
			code:   CodeICUSelectorMismatch,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			report := validateMessagePair(t, test.locale, test.source, test.target)
			if !reportHasFinding(report, test.code) || !HasFailures([]Report{report}) {
				t.Fatalf("report = %#v, want failing finding %q", report, test.code)
			}
		})
	}
}

func TestValidateAllowsLocalePluralBranchesAndExactSelectors(t *testing.T) {
	report := validateMessagePair(t, "ru",
		`{count, plural, one {One item} other {# items}}`,
		`{count, plural, =0 {Нет элементов} one {Один элемент} few {# элемента} other {# элементов}}`,
	)
	if HasFailures([]Report{report}) {
		t.Fatalf("valid Russian ICU message reported as failure: %#v", report)
	}
}

func TestValidateReportsMalformedICUSourceWithoutTarget(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	if err := os.WriteFile(sourcePath, []byte(`{"items":"{count, plural, one {One}}"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	reports, err := Validate(&config.Config{SourceLocale: "en", TargetLocales: []string{"fr"}, SourcePath: sourcePath})
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 || !reportHasFinding(reports[0], CodeICUMessageSyntax) {
		t.Fatalf("reports = %#v, want source syntax finding", reports)
	}
}

func validateMessagePair(t *testing.T, locale, source, target string) Report {
	t.Helper()
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "en.json")
	targetPath := filepath.Join(dir, locale+".json")
	if err := os.WriteFile(sourcePath, []byte(`{"items":`+strconv.Quote(source)+`}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"items":`+strconv.Quote(target)+`}`), 0o644); err != nil {
		t.Fatal(err)
	}
	reports, err := Validate(&config.Config{SourceLocale: "en", TargetLocales: []string{locale}, SourcePath: sourcePath})
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 {
		t.Fatalf("Validate returned %d reports, want 1", len(reports))
	}
	return reports[0]
}

func reportHasFinding(report Report, code FindingCode) bool {
	for _, finding := range report.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}
