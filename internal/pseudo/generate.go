package pseudo

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
	localeid "github.com/Tom-R-Main/Internationalizer/internal/locale"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
)

// GenerateOptions configures pseudolocale generation.
type GenerateOptions struct {
	Strategy Strategy
	Locale   string
	Force    bool
	DryRun   bool
}

// GenerateResult describes one generated bundle artifact.
type GenerateResult struct {
	Bundle     string
	Locale     string
	TargetPath string
	Units      int
	Written    bool
}

// Generate creates deterministic pseudolocale artifacts without an LLM or TM.
func Generate(cfg *config.Config, opts GenerateOptions) ([]GenerateResult, error) {
	effectiveConfig := *cfg
	effectiveConfig.ApplyDefaults()
	cfg = &effectiveConfig
	if err := cfg.ValidateProject(); err != nil {
		return nil, err
	}
	if _, err := DefaultLocale(opts.Strategy); err != nil {
		return nil, err
	}
	locale := opts.Locale
	if locale == "" {
		var err error
		locale, err = DefaultLocale(opts.Strategy)
		if err != nil {
			return nil, err
		}
	}
	canonicalLocale, err := localeid.Canonical(locale)
	if err != nil {
		return nil, fmt.Errorf("invalid pseudo locale %q: %w", locale, err)
	}

	manifest, err := state.Load(cfg.ManifestPath)
	if err != nil {
		return nil, err
	}
	policyHash, err := state.HashValue(struct {
		Version  int      `json:"version"`
		Strategy Strategy `json:"strategy"`
	}{Version: 1, Strategy: opts.Strategy})
	if err != nil {
		return nil, err
	}

	results := make([]GenerateResult, 0, len(cfg.EffectiveBundles()))
	now := time.Now().UTC()
	for _, bundle := range cfg.EffectiveBundles() {
		format, err := formatForBundle(bundle)
		if err != nil {
			return nil, fmt.Errorf("bundle %q format: %w", bundle.ID, err)
		}
		sourceData, err := os.ReadFile(bundle.Source)
		if err != nil {
			return nil, fmt.Errorf("reading bundle %q source %s: %w", bundle.ID, bundle.Source, err)
		}
		sourceUnits, err := formats.ParseUnits(format, sourceData)
		if err != nil {
			return nil, fmt.Errorf("parsing bundle %q source: %w", bundle.ID, err)
		}
		targetPath, err := bundle.TargetPath(locale)
		if err != nil {
			return nil, err
		}
		if existing, readErr := os.ReadFile(targetPath); readErr == nil {
			if !opts.Force && !pseudoOwnsArtifact(manifest, bundle.ID, canonicalLocale, format, sourceData, existing) {
				return nil, fmt.Errorf("refusing to overwrite %s without --force because it is not a tracked pseudo artifact", targetPath)
			}
		} else if !os.IsNotExist(readErr) {
			return nil, fmt.Errorf("reading pseudo target %s: %w", targetPath, readErr)
		}

		pseudoUnits := make([]formats.Unit, len(sourceUnits))
		for index, unit := range sourceUnits {
			transformed, err := transformUnit(format, unit, opts.Strategy)
			if err != nil {
				return nil, fmt.Errorf("pseudolocalizing bundle %q unit %q: %w", bundle.ID, unit.ID, err)
			}
			unit.Value = transformed
			pseudoUnits[index] = unit
		}
		var output []byte
		if paired, ok := format.(formats.PairedFormat); ok {
			output, err = paired.SerializeTarget(formats.UnitValues(pseudoUnits), sourceData, sourceData)
		} else {
			output, err = formats.SerializeUnits(format, pseudoUnits, sourceData)
		}
		if err != nil {
			return nil, fmt.Errorf("serializing pseudo target %s: %w", targetPath, err)
		}
		output = appendOneNewline(output)
		serializedValues, err := parseTargetValues(format, sourceData, output)
		if err != nil {
			return nil, fmt.Errorf("validating pseudo target %s: %w", targetPath, err)
		}
		result := GenerateResult{Bundle: bundle.ID, Locale: canonicalLocale, TargetPath: targetPath, Units: len(pseudoUnits), Written: !opts.DryRun}
		results = append(results, result)
		if opts.DryRun {
			continue
		}
		if err := state.WriteFileAtomic(targetPath, output, 0o644); err != nil {
			return nil, err
		}
		for _, sourceUnit := range sourceUnits {
			targetValue, ok := serializedValues[sourceUnit.ID]
			if !ok {
				return nil, fmt.Errorf("validating pseudo target %s: unit %q is missing after serialization", targetPath, sourceUnit.ID)
			}
			manifest.Set(state.Entry{
				Bundle:       bundle.ID,
				Key:          sourceUnit.ID,
				Locale:       canonicalLocale,
				SourceHash:   state.SourceUnitHash(format.Name(), sourceUnit.Value, sourceUnit.Context, sourceUnit.Structure),
				PolicyHash:   policyHash,
				TargetHash:   state.TargetHash(targetValue),
				Origin:       "pseudo",
				ReviewStatus: state.ReviewNeedsReview,
				UpdatedAt:    now,
			})
		}
	}
	if !opts.DryRun {
		if err := manifest.Save(cfg.ManifestPath); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func formatForBundle(bundle config.Bundle) (formats.Format, error) {
	if bundle.Format != "" {
		return formats.FormatByName(bundle.Format)
	}
	return formats.FormatForFile(bundle.Source)
}

func transformUnit(format formats.Format, unit formats.Unit, strategy Strategy) (string, error) {
	if format.Name() != "markdown" {
		return Transform(unit.Value, strategy)
	}
	lineEnd := strings.IndexByte(unit.Value, '\n')
	if lineEnd < 0 {
		lineEnd = len(unit.Value)
	}
	line := strings.TrimSuffix(unit.Value[:lineEnd], "\r")
	prefixEnd := 0
	for prefixEnd < len(line) && (line[prefixEnd] == ' ' || line[prefixEnd] == '\t') {
		prefixEnd++
	}
	headingStart := prefixEnd
	for prefixEnd < len(line) && line[prefixEnd] == '#' && prefixEnd-headingStart < 6 {
		prefixEnd++
	}
	if prefixEnd == headingStart || prefixEnd >= len(line) || (line[prefixEnd] != ' ' && line[prefixEnd] != '\t') {
		return Transform(unit.Value, strategy)
	}
	for prefixEnd < len(line) && (line[prefixEnd] == ' ' || line[prefixEnd] == '\t') {
		prefixEnd++
	}
	heading, err := Transform(line[prefixEnd:], strategy)
	if err != nil {
		return "", err
	}
	newline := ""
	body := ""
	if lineEnd < len(unit.Value) {
		newline = "\n"
		if lineEnd > 0 && unit.Value[lineEnd-1] == '\r' {
			newline = "\r\n"
		}
		body = unit.Value[lineEnd+1:]
	}
	transformedBody, err := transformPadded(body, strategy)
	if err != nil {
		return "", err
	}
	return line[:prefixEnd] + heading + newline + transformedBody, nil
}

func transformPadded(value string, strategy Strategy) (string, error) {
	start := 0
	for start < len(value) && strings.ContainsRune(" \t\r\n", rune(value[start])) {
		start++
	}
	end := len(value)
	for end > start && strings.ContainsRune(" \t\r\n", rune(value[end-1])) {
		end--
	}
	if start == end {
		return value, nil
	}
	transformed, err := Transform(value[start:end], strategy)
	if err != nil {
		return "", err
	}
	return value[:start] + transformed + value[end:], nil
}

func pseudoOwnsArtifact(manifest *state.Manifest, bundle, locale string, format formats.Format, source, target []byte) bool {
	sourceUnits, err := formats.ParseUnits(format, source)
	if err != nil || len(sourceUnits) == 0 {
		return false
	}
	targetValues, err := parseTargetValues(format, source, target)
	if err != nil || len(targetValues) != len(sourceUnits) {
		return false
	}
	for _, sourceUnit := range sourceUnits {
		targetValue, ok := targetValues[sourceUnit.ID]
		if !ok {
			return false
		}
		entry, ok := manifest.Get(bundle, sourceUnit.ID, locale)
		if !ok || entry.Origin != "pseudo" || entry.TargetHash != state.TargetHash(targetValue) {
			return false
		}
	}
	return true
}

func parseTargetValues(format formats.Format, source, target []byte) (map[string]string, error) {
	if paired, ok := format.(formats.PairedFormat); ok {
		return paired.ParseTarget(source, target)
	}
	targetUnits, err := formats.ParseUnits(format, target)
	if err != nil {
		return nil, err
	}
	return formats.UnitValues(targetUnits), nil
}

func appendOneNewline(data []byte) []byte {
	return append(bytes.TrimRight(data, "\n"), '\n')
}
