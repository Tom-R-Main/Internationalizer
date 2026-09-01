package validate

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
)

// Report holds validation results for a single locale.
type Report struct {
	Bundle     string     `json:"bundle"`
	Locale     string     `json:"locale"`
	TargetPath string     `json:"target_path"`
	Missing    []string   `json:"missing"`
	Extra      []string   `json:"extra"`
	Mismatches []Mismatch `json:"mismatches,omitempty"`
	Errors     []string   `json:"errors,omitempty"`
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
	if err := cfg.ValidateProject(); err != nil {
		return nil, err
	}
	var reports []Report
	for _, bundle := range cfg.EffectiveBundles() {
		format, err := formatForBundle(bundle)
		if err != nil {
			return nil, fmt.Errorf("bundle %q format: %w", bundle.ID, err)
		}
		sourceData, err := os.ReadFile(bundle.Source)
		if err != nil {
			return nil, fmt.Errorf("reading bundle %q source %s: %w", bundle.ID, bundle.Source, err)
		}
		sourceKeys, err := format.Parse(sourceData)
		if err != nil {
			return nil, fmt.Errorf("parsing bundle %q source: %w", bundle.ID, err)
		}
		for _, locale := range cfg.TargetLocales {
			targetPath, err := bundle.TargetPath(locale)
			if err != nil {
				return nil, err
			}
			report := validateLocale(bundle.ID, locale, sourceKeys, targetPath, format)
			reports = append(reports, report)
		}
	}
	return reports, nil
}

func formatForBundle(bundle config.Bundle) (formats.Format, error) {
	if bundle.Format != "" {
		return formats.FormatByName(bundle.Format)
	}
	return formats.FormatForFile(bundle.Source)
}

func validateLocale(bundle, locale string, sourceKeys map[string]string, targetPath string, format formats.Format) Report {
	report := Report{Bundle: bundle, Locale: locale, TargetPath: targetPath}

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
		report.Errors = append(report.Errors, fmt.Sprintf("parsing target: %v", err))
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
		if mismatch := InterpolationMismatch(key, sourceVal, targetVal); mismatch != nil {
			report.Mismatches = append(report.Mismatches, *mismatch)
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

// InterpolationMismatch compares the placeholder multiset in two values.
func InterpolationMismatch(key, source, target string) *Mismatch {
	sourceVars := extractVars(source)
	targetVars := extractVars(target)
	if sameVars(sourceVars, targetVars) {
		return nil
	}
	return &Mismatch{Key: key, SourceVars: sourceVars, TargetVars: targetVars}
}

// HasFailures reports whether validation should return a non-zero exit status.
func HasFailures(reports []Report) bool {
	for _, report := range reports {
		if len(report.Missing) > 0 || len(report.Mismatches) > 0 || len(report.Errors) > 0 {
			return true
		}
	}
	return false
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
		if len(r.Missing) > 0 || len(r.Mismatches) > 0 || len(r.Errors) > 0 {
			status = "FAIL"
			hasErrors = true
		}

		_, _ = fmt.Fprintf(&b, "[%s/%s] %s — %.1f%% coverage", r.Bundle, r.Locale, status, r.Coverage)

		if len(r.Missing) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d missing", len(r.Missing))
		}
		if len(r.Extra) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d extra", len(r.Extra))
		}
		if len(r.Mismatches) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d interpolation mismatches", len(r.Mismatches))
		}
		if len(r.Errors) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d errors", len(r.Errors))
		}
		b.WriteString("\n")

		// Show details for failures.
		if len(r.Missing) > 0 && len(r.Missing) <= 20 {
			for _, key := range r.Missing {
				_, _ = fmt.Fprintf(&b, "  - missing: %s\n", key)
			}
		}
		for _, m := range r.Mismatches {
			_, _ = fmt.Fprintf(&b, "  - mismatch: %s (source: %v, target: %v)\n",
				m.Key, m.SourceVars, m.TargetVars)
		}
		for _, reportErr := range r.Errors {
			_, _ = fmt.Fprintf(&b, "  - error: %s\n", reportErr)
		}
	}

	if !hasErrors {
		b.WriteString("\nAll locales valid.\n")
	}
	return b.String()
}
