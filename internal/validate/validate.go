package validate

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/Tom-R-Main/Internationalizer/internal/jsonintegrity"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"github.com/Tom-R-Main/Internationalizer/internal/policy"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
	"github.com/Tom-R-Main/Internationalizer/internal/styleguide"
)

// Options enables validation contracts that are intentionally opt-in during
// their compatibility period.
type Options struct {
	Strict          bool
	RequireState    bool
	RequireApproved bool
}

// Report holds validation results for a single locale. Coverage remains a
// compatibility alias for StructuralCoverage.
type Report struct {
	SourcePath         string     `json:"source_path,omitempty"`
	BlockedBySource    bool       `json:"blocked_by_source,omitempty"`
	BlockedLocales     []string   `json:"blocked_locales,omitempty"`
	Bundle             string     `json:"bundle"`
	Locale             string     `json:"locale"`
	TargetPath         string     `json:"target_path"`
	Missing            []string   `json:"missing"`
	Extra              []string   `json:"extra"`
	Mismatches         []Mismatch `json:"mismatches,omitempty"`
	Errors             []string   `json:"errors,omitempty"`
	Coverage           float64    `json:"coverage"`
	StructuralCoverage float64    `json:"structural_coverage"`
	TranslatedCoverage *float64   `json:"translated_coverage,omitempty"`
	Findings           []Finding  `json:"findings,omitempty"`
}

// Mismatch indicates an interpolation variable difference between source and target.
type Mismatch struct {
	Key        string   `json:"key"`
	SourceVars []string `json:"source_vars"`
	TargetVars []string `json:"target_vars"`
}

var interpolationRe = regexp.MustCompile(`(?:\{\{(\w+)\}\}|\{(\w+)\}|%\{(\w+)\})`)

// Validate preserves the original validation behavior.
func Validate(cfg *config.Config) ([]Report, error) {
	return ValidateWithOptions(cfg, Options{})
}

// ValidateWithOptions checks all target locales against the source locale.
func ValidateWithOptions(cfg *config.Config, opts Options) ([]Report, error) {
	if opts.RequireApproved {
		opts.RequireState = true
	}
	effectiveConfig := *cfg
	effectiveConfig.ApplyDefaults()
	cfg = &effectiveConfig
	if err := cfg.ValidateProject(); err != nil {
		return nil, err
	}

	var manifest *state.Manifest
	if opts.RequireState {
		var err error
		manifest, err = state.Load(cfg.ManifestPath)
		if err != nil {
			return nil, err
		}
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
		sourceUnits, err := formats.ParseSourceUnits(format, sourceData, bundle.MessageSyntax)
		if err != nil {
			return nil, fmt.Errorf("parsing bundle %q source %s: %w", bundle.ID, bundle.Source, err)
		}
		sourceKeys := formats.UnitValues(sourceUnits)
		sourceFindings := make(map[string][]Finding)
		for _, unit := range sourceUnits {
			if findings := SyntaxSourceFindings(unit.ID, unit.Value, cfg.SourceLocale, unit.Syntax, bundle.MessageSyntax); len(findings) > 0 {
				sourceFindings[unit.ID] = findings
			}
		}
		for localeIndex, locale := range cfg.TargetLocales {
			targetPath, err := bundle.TargetPath(locale)
			if err != nil {
				return nil, err
			}

			var terms []glossary.Term
			if opts.Strict || opts.RequireState {
				terms, err = glossary.Load(cfg.GlossaryDir, locale)
				if err != nil {
					return nil, err
				}
			}

			policyHash := ""
			if opts.RequireState {
				guide, loadErr := styleguide.Load(cfg.StyleGuidesDir, locale)
				if loadErr != nil {
					return nil, loadErr
				}
				resolved, resolveErr := policy.Resolve(cfg, locale, format.Name(), guide, terms, bundle.MessageSyntax)
				if resolveErr != nil {
					return nil, fmt.Errorf("hashing translation policy: %w", resolveErr)
				}
				policyHash = resolved.Hash
			}

			pluralStyle := bundle.PluralStyle(cfg.Validation.PluralStyle)
			report := validateLocale(bundle.ID, cfg.SourceLocale, locale, bundle.Source, sourceData, sourceUnits, sourceKeys, targetPath, format, terms, manifest, policyHash, pluralStyle, opts, bundle.MessageSyntax, sourceFindings)
			if len(sourceFindings) > 0 {
				report.SourcePath = bundle.Source
				report.BlockedBySource = true
				if localeIndex == 0 {
					report.BlockedLocales = cfg.TargetLocales
					for _, findings := range sourceFindings {
						report.Findings = append(report.Findings, findings...)
					}
					sortFindings(report.Findings)
				}
			}
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

func validateLocale(bundle, sourceLocale, locale, sourcePath string, sourceData []byte, sourceUnits []formats.Unit, sourceKeys map[string]string, targetPath string, format formats.Format, terms []glossary.Term, manifest *state.Manifest, policyHash, pluralStyle string, opts Options, syntaxPolicy message.Syntax, sourceFindings map[string][]Finding) Report {
	report := Report{Bundle: bundle, Locale: locale, TargetPath: targetPath}
	if opts.Strict {
		translated := 0.0
		report.TranslatedCoverage = &translated
	}

	validationKeys := sourceKeys
	requiredPluralKeys := make(map[string]struct{})
	optionalPluralKeys := make(map[string]struct{})
	if (opts.Strict || opts.RequireState || syntaxPolicy == message.I18next) && pluralStyle == "i18next-v4" {
		validationKeys, requiredPluralKeys, optionalPluralKeys = ExpandI18nextV4Source(sourceKeys, sourceLocale, locale)
	}
	sourceICUValid := make(map[string]bool, len(validationKeys))
	syntaxes := make(map[string]message.Syntax, len(validationKeys))
	for _, unit := range sourceUnits {
		syntaxes[unit.ID] = unit.Syntax
	}
	document := format.Name() == "markdown"
	for _, key := range allKeys(validationKeys) {
		sourceValue := validationKeys[key]
		if syntaxes[key] == "" {
			syntaxes[key] = message.ResolveSyntax(format.Name(), syntaxPolicy, sourceValue)
		}
		findings := sourceFindings[key]
		sourceICUValid[key] = len(findings) == 0
	}

	targetData, err := os.ReadFile(targetPath)
	if err != nil {
		for _, key := range allKeys(validationKeys) {
			if _, optional := optionalPluralKeys[key]; optional {
				continue
			}
			report.Missing = append(report.Missing, key)
			report.Findings = append(report.Findings, missingFinding(key, requiredPluralKeys))
		}
		sortFindings(report.Findings)
		return report
	}

	targetKeys, err := parseTarget(format, sourceData, targetData)
	if err != nil {
		report.Missing = allKeys(validationKeys)
		report.Errors = append(report.Errors, fmt.Sprintf("parsing target %s: %v", targetPath, err))
		var integrity *jsonintegrity.Error
		if errors.As(err, &integrity) {
			report.Findings = append(report.Findings, Finding{Code: FindingCode(integrity.JSONCode()), Severity: SeverityError, Message: err.Error(), Key: integrity.Key, Path: integrity.Path, OtherPath: integrity.OtherPath})
		}
		sortFindings(report.Findings)
		return report
	}
	sourceUnitsByID := make(map[string]formats.Unit, len(sourceUnits))
	for _, unit := range sourceUnits {
		sourceUnitsByID[unit.ID] = unit
	}

	present := 0
	translated := 0
	for key, sourceValue := range validationKeys {
		_, coverageKey := optionalPluralKeys[key]
		coverageKey = !coverageKey
		targetValue, ok := targetKeys[key]
		if !ok {
			if _, optional := optionalPluralKeys[key]; optional {
				if coverageKey {
					present++
					translated++
				}
				continue
			}
			report.Missing = append(report.Missing, key)
			report.Findings = append(report.Findings, missingFinding(key, requiredPluralKeys))
			continue
		}
		if coverageKey {
			present++
		}

		blank := strings.TrimSpace(sourceValue) != "" && strings.TrimSpace(targetValue) == ""
		identical := sourceValue == targetValue && !isNonLinguistic(sourceValue)
		exempt := identical && glossary.SourceIdenticalExempt(terms, sourceValue, targetValue)
		if coverageKey && !blank && (!identical || exempt) {
			translated++
		}
		if opts.Strict && blank {
			report.Findings = append(report.Findings, Finding{Code: CodeBlankTranslation, Severity: SeverityError, Key: key, Message: "non-empty source has a blank target"})
		}
		if opts.Strict && identical && !exempt {
			report.Findings = append(report.Findings, Finding{Code: CodeSourceIdentical, Severity: SeverityError, Key: key, Message: "target is identical to the source without an exact glossary exemption"})
		}

		if mismatch := SyntaxInterpolationMismatch(key, sourceValue, targetValue, syntaxes[key]); mismatch != nil {
			report.Mismatches = append(report.Mismatches, *mismatch)
		}
		if sourceICUValid[key] && !document {
			report.Findings = append(report.Findings, SyntaxFindings(key, sourceValue, targetValue, locale, syntaxes[key])...)
		}
		if opts.Strict || opts.RequireState || (syntaxPolicy != "" && syntaxPolicy != message.Auto) {
			if document {
				report.Findings = append(report.Findings, ProtectedDocumentFindings(key, sourceValue, targetValue, locale, sourcePath, targetPath, syntaxes[key])...)
			} else {
				report.Findings = append(report.Findings, ProtectedSyntaxFindings(key, sourceValue, targetValue, locale, syntaxes[key])...)
			}
		}
		if opts.Strict {
			report.Findings = append(report.Findings, glossaryFindings(key, sourceValue, targetValue, terms)...)
		}
		if opts.RequireState {
			sourceUnit := sourceUnitsByID[key]
			sourceUnit.Value = sourceValue
			report.Findings = append(report.Findings, provenanceFindings(bundle, key, locale, format.Name(), sourceUnit, targetValue, policyHash, manifest, opts.RequireApproved)...)
		}
	}

	for key := range targetKeys {
		if _, ok := validationKeys[key]; ok {
			continue
		}
		report.Extra = append(report.Extra, key)
		severity := SeverityWarning
		if opts.Strict || opts.RequireState {
			severity = SeverityError
		}
		report.Findings = append(report.Findings, Finding{Code: CodeExtraKey, Severity: severity, Key: key, Message: "target contains a key absent from the source"})
	}

	sort.Strings(report.Missing)
	sort.Strings(report.Extra)
	sort.Slice(report.Mismatches, func(i, j int) bool { return report.Mismatches[i].Key < report.Mismatches[j].Key })
	sortFindings(report.Findings)

	total := len(validationKeys) - len(optionalPluralKeys)
	if total > 0 {
		report.StructuralCoverage = float64(present) / float64(total) * 100
		report.Coverage = report.StructuralCoverage
		if report.TranslatedCoverage != nil {
			value := float64(translated) / float64(total) * 100
			report.TranslatedCoverage = &value
		}
	}
	return report
}

func parseTarget(format formats.Format, source, target []byte) (map[string]string, error) {
	if paired, ok := format.(formats.PairedFormat); ok {
		return paired.ParseTarget(source, target)
	}
	return format.Parse(target)
}

// ICUFindings compares ICU MessageFormat structure independently of linguistic
// strictness. Message syntax and branch identity are runtime correctness.
func ICUFindings(key, source, target, targetLocale string) []Finding {
	return messageFindings(key, message.Compare(source, target, targetLocale))
}

// ICUSourceFindings validates source syntax and source-locale plural selectors
// even when a target file does not exist yet.
func ICUSourceFindings(key, source, sourceLocale string) []Finding {
	return messageFindings(key, message.Compare(source, source, sourceLocale))
}

func messageFindings(key string, issues []message.Issue) []Finding {
	findings := make([]Finding, 0, len(issues))
	for _, issue := range issues {
		code := CodeICUArgumentMismatch
		switch issue.Code {
		case message.CodeSyntax:
			code = CodeICUMessageSyntax
		case message.CodeSelectorMismatch, message.CodeInvalidPluralCategory:
			code = CodeICUSelectorMismatch
		}
		findings = append(findings, Finding{
			Code:     code,
			Severity: SeverityError,
			Key:      key,
			Message:  issue.Message,
		})
	}
	return findings
}

func missingFinding(key string, requiredPluralKeys map[string]struct{}) Finding {
	if _, required := requiredPluralKeys[key]; required {
		return Finding{Code: CodePluralFormMissing, Severity: SeverityError, Key: key, Message: "target locale requires this plural form"}
	}
	return Finding{Code: CodeMissingKey, Severity: SeverityError, Key: key, Message: "target key is missing"}
}

func provenanceFindings(bundle, key, locale, format string, source formats.Unit, target, policyHash string, manifest *state.Manifest, requireApproved bool) []Finding {
	recorded, ok := manifest.Get(bundle, key, locale)
	if !ok {
		return []Finding{{Code: CodeUntracked, Severity: SeverityError, Key: key, Message: "target has no manifest provenance"}}
	}
	var findings []Finding
	if recorded.SourceHash != state.SourceUnitHash(format, source.Value, source.Context, source.Structure) {
		findings = append(findings, Finding{Code: CodeSourceStale, Severity: SeverityError, Key: key, Message: "source changed after the recorded translation"})
	}
	if recorded.PolicyHash != policyHash {
		findings = append(findings, Finding{Code: CodePolicyStale, Severity: SeverityError, Key: key, Message: "translation policy changed after the recorded translation"})
	}
	if recorded.TargetHash != state.TargetHash(target) {
		findings = append(findings, Finding{Code: CodeTargetModified, Severity: SeverityError, Key: key, Message: "target changed after the recorded translation"})
	}
	if requireApproved && recorded.ReviewStatus != state.ReviewApproved {
		findings = append(findings, Finding{Code: CodeNeedsReview, Severity: SeverityError, Key: key, Message: "translation has not been explicitly approved"})
	}
	return findings
}

func glossaryFindings(key, source, target string, terms []glossary.Term) []Finding {
	var findings []Finding
	for _, term := range terms {
		if !containsGlossaryValue(source, term.Source, term.IgnoreCase, term.WholeWord) {
			continue
		}
		approved := glossary.ApprovedTargets(term)
		matched := false
		for _, candidate := range approved {
			if containsGlossaryValue(target, candidate, term.IgnoreCase, term.WholeWord) {
				matched = true
				break
			}
		}
		if matched {
			continue
		}
		severity := SeverityError
		if term.Enforcement == glossary.EnforcementWarning {
			severity = SeverityWarning
		}
		findings = append(findings, Finding{Code: CodeGlossaryViolation, Severity: severity, Key: key, Message: fmt.Sprintf("source term %q is missing an approved target form", term.Source), Expected: approved})
	}
	return findings
}

func containsGlossaryValue(value, term string, ignoreCase, wholeWord bool) bool {
	if term == "" {
		return false
	}
	if ignoreCase {
		matcher := regexp.MustCompile(`(?i:` + regexp.QuoteMeta(term) + `)`)
		for _, location := range matcher.FindAllStringIndex(value, -1) {
			if !wholeWord || (wordBoundary(value, location[0]-1) && wordBoundary(value, location[1])) {
				return true
			}
		}
		return false
	}
	for offset := 0; ; {
		index := strings.Index(value[offset:], term)
		if index < 0 {
			return false
		}
		start := offset + index
		end := start + len(term)
		if !wholeWord || (wordBoundary(value, start-1) && wordBoundary(value, end)) {
			return true
		}
		offset = end
		if offset >= len(value) {
			return false
		}
	}
}

func wordBoundary(value string, byteIndex int) bool {
	if byteIndex < 0 || byteIndex >= len(value) {
		return true
	}
	for byteIndex > 0 && value[byteIndex]&0xC0 == 0x80 {
		byteIndex--
	}
	for _, r := range value[byteIndex:] {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_'
	}
	return true
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
		if report.BlockedBySource || len(report.Missing) > 0 || len(report.Mismatches) > 0 || len(report.Errors) > 0 {
			return true
		}
		for _, finding := range report.Findings {
			if finding.Severity == SeverityError {
				return true
			}
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

func isNonLinguistic(value string) bool {
	for _, r := range value {
		if unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// FormatHuman returns a human-readable summary of validation reports.
func FormatHuman(reports []Report) string {
	var b strings.Builder
	hasErrors := false
	for _, report := range reports {
		status := "OK"
		if HasFailures([]Report{report}) {
			status = "FAIL"
			hasErrors = true
		}
		if report.TranslatedCoverage == nil {
			_, _ = fmt.Fprintf(&b, "[%s/%s] %s — %.1f%% coverage", report.Bundle, report.Locale, status, report.Coverage)
		} else {
			_, _ = fmt.Fprintf(&b, "[%s/%s] %s — %.1f%% structural coverage", report.Bundle, report.Locale, status, report.StructuralCoverage)
			_, _ = fmt.Fprintf(&b, ", %.1f%% translated coverage", *report.TranslatedCoverage)
		}
		if len(report.Missing) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d missing", len(report.Missing))
		}
		if len(report.Extra) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d extra", len(report.Extra))
		}
		if len(report.Mismatches) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d interpolation mismatches", len(report.Mismatches))
		}
		if len(report.Errors) > 0 {
			_, _ = fmt.Fprintf(&b, ", %d errors", len(report.Errors))
		}
		b.WriteByte('\n')
		if len(report.BlockedLocales) > 0 {
			_, _ = fmt.Fprintf(&b, "  Source %s blocks %d locales: %s\n", report.SourcePath, len(report.BlockedLocales), strings.Join(report.BlockedLocales, ", "))
		} else if report.BlockedBySource {
			_, _ = fmt.Fprintf(&b, "  Blocked by source errors in %s (reported above)\n", report.SourcePath)
		}
		for _, finding := range report.Findings {
			_, _ = fmt.Fprintf(&b, "  - %s: %s", finding.Code, finding.Message)
			if finding.Key != "" {
				_, _ = fmt.Fprintf(&b, " (%s)", finding.Key)
			}
			b.WriteByte('\n')
		}
		for _, mismatch := range report.Mismatches {
			if hasProtectedInterpolationFinding(report.Findings, mismatch.Key) {
				continue
			}
			_, _ = fmt.Fprintf(&b, "  - mismatch: %s (source: %v, target: %v)\n", mismatch.Key, mismatch.SourceVars, mismatch.TargetVars)
		}
		for _, reportErr := range report.Errors {
			_, _ = fmt.Fprintf(&b, "  - error: %s\n", reportErr)
		}
	}
	if !hasErrors {
		b.WriteString("\nAll locales valid.\n")
	}
	return b.String()
}

func hasProtectedInterpolationFinding(findings []Finding, key string) bool {
	for _, finding := range findings {
		if finding.Code == CodeProtectedStructureMismatch && finding.Key == key && finding.Message == "protected interpolation variables mismatch" {
			return true
		}
	}
	return false
}
