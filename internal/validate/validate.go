package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
)

// Report holds validation results for a single locale.
type Report struct {
	Locale     string     `json:"locale"`
	Missing    []string   `json:"missing"`
	Extra      []string   `json:"extra"`
	Mismatches []Mismatch `json:"mismatches,omitempty"`
	Coverage   float64    `json:"coverage"`
}

// Mismatch indicates an interpolation variable difference between source and target.
type Mismatch struct {
	Key        string   `json:"key"`
	SourceVars []string `json:"source_vars"`
	TargetVars []string `json:"target_vars"`
}

// interpolation patterns: {{var}}, {var}, %{var}
var interpolationRe = regexp.MustCompile(`(?:\{\{(\w+)\}\}|\{(\w+)\}|%\{(\w+)\})`)

// Validate checks all target locales against the source locale.
func Validate(cfg *config.Config) ([]Report, error) {
	sourceDir := filepath.Dir(cfg.SourcePath)
	sourceFile := filepath.Base(cfg.SourcePath)

	format, err := formats.FormatForFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("detecting format: %w", err)
	}

	sourceData, err := os.ReadFile(cfg.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("reading source %s: %w", cfg.SourcePath, err)
	}

	sourceKeys, err := format.Parse(sourceData)
	if err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}

	var reports []Report
	for _, locale := range cfg.TargetLocales {
		targetPath := filepath.Join(sourceDir, locale+filepath.Ext(sourceFile))
		report := validateLocale(locale, sourceKeys, targetPath, format)
		reports = append(reports, report)
	}
	return reports, nil
}

func validateLocale(locale string, sourceKeys map[string]string, targetPath string, format formats.Format) Report {
	report := Report{Locale: locale}

	targetData, err := os.ReadFile(targetPath)
	if err != nil {
		// Target file doesn't exist — all keys are missing.
		for key := range sourceKeys {
			report.Missing = append(report.Missing, key)
		}
		sort.Strings(report.Missing)
		report.Coverage = 0
		return report
	}

	targetKeys, err := format.Parse(targetData)
	if err != nil {
		report.Missing = allKeys(sourceKeys)
		report.Coverage = 0
		return report
	}

	// Find missing and extra keys.
	for key := range sourceKeys {
		if _, ok := targetKeys[key]; !ok {
			report.Missing = append(report.Missing, key)
		}
	}
	for key := range targetKeys {
		if _, ok := sourceKeys[key]; !ok {
			report.Extra = append(report.Extra, key)
		}
	}

	// Check interpolation mismatches on shared keys.
	for key, sourceVal := range sourceKeys {
		targetVal, ok := targetKeys[key]
		if !ok {
			continue
		}
		srcVars := extractVars(sourceVal)
		tgtVars := extractVars(targetVal)
		if !sameVars(srcVars, tgtVars) {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Key:        key,
				SourceVars: srcVars,
				TargetVars: tgtVars,
			})
		}
	}

	sort.Strings(report.Missing)
	sort.Strings(report.Extra)

	total := len(sourceKeys)
	if total > 0 {
		report.Coverage = float64(total-len(report.Missing)) / float64(total) * 100
	}
	return report
}

func extractVars(s string) []string {
	matches := interpolationRe.FindAllStringSubmatch(s, -1)
	var vars []string
	for _, m := range matches {
		for _, g := range m[1:] {
			if g != "" {
				vars = append(vars, g)
			}
		}
	}
	sort.Strings(vars)
	return vars
}

func sameVars(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func allKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FormatHuman returns a human-readable summary of validation reports.
func FormatHuman(reports []Report) string {
	var b strings.Builder
	hasErrors := false

	for _, r := range reports {
		status := "OK"
		if len(r.Missing) > 0 || len(r.Mismatches) > 0 {
			status = "FAIL"
			hasErrors = true
		}

		b.WriteString(fmt.Sprintf("[%s] %s — %.1f%% coverage", r.Locale, status, r.Coverage))

		if len(r.Missing) > 0 {
			b.WriteString(fmt.Sprintf(", %d missing", len(r.Missing)))
		}
		if len(r.Extra) > 0 {
			b.WriteString(fmt.Sprintf(", %d extra", len(r.Extra)))
		}
		if len(r.Mismatches) > 0 {
			b.WriteString(fmt.Sprintf(", %d interpolation mismatches", len(r.Mismatches)))
		}
		b.WriteString("\n")

		// Show details for failures.
		if len(r.Missing) > 0 && len(r.Missing) <= 20 {
			for _, key := range r.Missing {
				b.WriteString(fmt.Sprintf("  - missing: %s\n", key))
			}
		}
		for _, m := range r.Mismatches {
			b.WriteString(fmt.Sprintf("  - mismatch: %s (source: %v, target: %v)\n",
				m.Key, m.SourceVars, m.TargetVars))
		}
	}

	if !hasErrors {
		b.WriteString("\nAll locales valid.\n")
	}
	return b.String()
}
