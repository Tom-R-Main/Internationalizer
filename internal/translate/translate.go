package translate

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/policy"
	"github.com/Tom-R-Main/Internationalizer/internal/state"
	"github.com/Tom-R-Main/Internationalizer/internal/styleguide"
	"github.com/Tom-R-Main/Internationalizer/internal/tm"
	validation "github.com/Tom-R-Main/Internationalizer/internal/validate"
)

// Options configures a translation run.
type Options struct {
	DryRun          bool
	AdoptExisting   bool
	RefreshPolicy   bool
	Locales         []string                // filter to specific locales; empty = all from config
	BatchSize       int                     // override config; 0 = use config
	Concurrency     int                     // override config; 0 = use config
	LocaleProviders map[string]llm.Provider // provider overrides keyed by target locale
}

// Result holds the outcome of one bundle and locale. Lifecycle counters are
// pre-run observations; manual, source-stale, and policy-stale may overlap.
type Result struct {
	Bundle          string
	Locale          string
	TargetPath      string
	KeysTotal       int
	KeysMissing     int
	KeysSourceStale int
	KeysPolicyStale int
	KeysManualEdit  int
	KeysUntracked   int
	KeysCurrent     int
	KeysCached      int
	KeysTranslated  int
	KeysSkipped     int
	Batches         int
	TokensIn        int
	TokensOut       int
	Errors          []string
}

// RunError reports that one or more locale jobs failed. Results remain
// available to the caller for human- or machine-readable reporting.
type RunError struct {
	Failed int
}

func (e *RunError) Error() string {
	return fmt.Sprintf("translation failed for %d bundle/locale job(s)", e.Failed)
}

type preparedBundle struct {
	bundle      config.Bundle
	format      formats.Format
	sourceUnits []formats.Unit
	sourceKeys  map[string]string
	sourceData  []byte
}

type job struct {
	index  int
	bundle preparedBundle
	locale string
}

type jobOutput struct {
	result  Result
	updates []state.Entry
}

// Run executes the translation pipeline.
func Run(ctx context.Context, cfg *config.Config, provider llm.Provider, opts Options) ([]Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	effectiveConfig := *cfg
	effectiveConfig.ApplyDefaults()
	cfg = &effectiveConfig
	if err := cfg.ValidateProject(); err != nil {
		return nil, err
	}

	batchSize := cfg.BatchSize
	if opts.BatchSize > 0 {
		batchSize = opts.BatchSize
	}
	if batchSize <= 0 {
		batchSize = 40
	}
	concurrency := cfg.Concurrency
	if opts.Concurrency > 0 {
		concurrency = opts.Concurrency
	}
	if concurrency <= 0 {
		concurrency = 1
	}

	locales := cfg.TargetLocales
	if len(opts.Locales) > 0 {
		seen := make(map[string]struct{}, len(opts.Locales))
		locales = make([]string, 0, len(opts.Locales))
		for _, requested := range opts.Locales {
			locale, ok := cfg.ConfiguredTargetLocale(requested)
			if !ok {
				return nil, fmt.Errorf("locale %q is not in target_locales", requested)
			}
			if _, ok := seen[locale]; ok {
				continue
			}
			seen[locale] = struct{}{}
			locales = append(locales, locale)
		}
	}

	bundles, err := prepareBundles(cfg.EffectiveBundles())
	if err != nil {
		return nil, err
	}
	memory, err := tm.Load(cfg.TMPath)
	if err != nil {
		return nil, fmt.Errorf("loading TM: %w", err)
	}
	manifest, err := state.Load(cfg.ManifestPath)
	if err != nil {
		return nil, err
	}

	jobCount := len(bundles) * len(locales)
	outputs := make([]jobOutput, jobCount)
	jobs := make(chan job)
	workerCount := concurrency
	if workerCount > jobCount {
		workerCount = jobCount
	}

	var workers sync.WaitGroup
	workers.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer workers.Done()
			for current := range jobs {
				currentProvider := provider
				if localeProvider, ok := opts.LocaleProviders[current.locale]; ok {
					currentProvider = localeProvider
				} else if cfg.HasLLMOverrideForLocale(current.locale) {
					currentProvider = nil
				}
				outputs[current.index] = translateLocale(ctx, cfg, current.bundle, currentProvider, memory, manifest, current.locale, batchSize, opts)
			}
		}()
	}

	next := 0
enqueue:
	for _, bundle := range bundles {
		for _, locale := range locales {
			select {
			case jobs <- job{index: next, bundle: bundle, locale: locale}:
				next++
			case <-ctx.Done():
				break enqueue
			}
		}
	}
	close(jobs)
	workers.Wait()

	results := make([]Result, 0, next)
	for i := 0; i < next; i++ {
		results = append(results, outputs[i].result)
		if !opts.DryRun {
			for _, update := range outputs[i].updates {
				manifest.Set(update)
			}
		}
	}
	if ctx.Err() != nil {
		return results, ctx.Err()
	}

	if !opts.DryRun && hasUpdates(outputs[:next]) {
		if err := manifest.Save(cfg.ManifestPath); err != nil {
			return results, err
		}
	}

	failed := 0
	for _, result := range results {
		if len(result.Errors) > 0 {
			failed++
		}
	}
	if failed > 0 {
		return results, &RunError{Failed: failed}
	}
	return results, nil
}

func prepareBundles(bundles []config.Bundle) ([]preparedBundle, error) {
	prepared := make([]preparedBundle, 0, len(bundles))
	for _, bundle := range bundles {
		var format formats.Format
		var err error
		if bundle.Format != "" {
			format, err = formats.FormatByName(bundle.Format)
		} else {
			format, err = formats.FormatForFile(bundle.Source)
		}
		if err != nil {
			return nil, fmt.Errorf("bundle %q format: %w", bundle.ID, err)
		}
		data, err := os.ReadFile(bundle.Source)
		if err != nil {
			return nil, fmt.Errorf("reading bundle %q source %s: %w", bundle.ID, bundle.Source, err)
		}
		units, err := formats.ParseUnits(format, data)
		if err != nil {
			return nil, fmt.Errorf("parsing bundle %q source: %w", bundle.ID, err)
		}
		prepared = append(prepared, preparedBundle{bundle: bundle, format: format, sourceUnits: units, sourceKeys: formats.UnitValues(units), sourceData: data})
	}
	return prepared, nil
}

func translateLocale(
	ctx context.Context,
	cfg *config.Config,
	bundle preparedBundle,
	provider llm.Provider,
	memory *tm.TM,
	manifest *state.Manifest,
	locale string,
	batchSize int,
	opts Options,
) jobOutput {
	sourceKeys := bundle.sourceKeys
	var optionalPluralKeys map[string]struct{}
	if cfg.Validation.PluralStyle == "i18next-v4" {
		sourceKeys, _, optionalPluralKeys = validation.ExpandI18nextV4Source(bundle.sourceKeys, cfg.SourceLocale, locale)
		for key := range optionalPluralKeys {
			delete(sourceKeys, key)
		}
	}
	result := Result{Bundle: bundle.bundle.ID, Locale: locale, KeysTotal: len(sourceKeys)}
	for _, key := range sortedKeys(sourceKeys) {
		if findings := validation.ICUSourceFindings(key, sourceKeys[key], cfg.SourceLocale); len(findings) > 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("source %q: %s", key, findings[0].Message))
		}
	}
	if len(result.Errors) > 0 {
		return jobOutput{result: result}
	}
	targetPath, err := bundle.bundle.TargetPath(locale)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return jobOutput{result: result}
	}
	result.TargetPath = targetPath

	terms, err := glossary.Load(cfg.GlossaryDir, locale)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("glossary: %v", err))
		return jobOutput{result: result}
	}
	guide, err := styleguide.Load(cfg.StyleGuidesDir, locale)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("style guide: %v", err))
		return jobOutput{result: result}
	}
	translationPolicy, err := policy.Resolve(cfg, locale, bundle.format.Name(), guide, terms)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("hashing translation policy: %v", err))
		return jobOutput{result: result}
	}
	localeLLM := translationPolicy.LLM
	prompt := translationPolicy.Prompt
	policyHash := translationPolicy.Hash

	targetKeys := make(map[string]string)
	var targetUnits []formats.Unit
	var targetData []byte
	targetExists := false
	data, err := os.ReadFile(targetPath)
	switch {
	case err == nil:
		targetExists = true
		targetData = data
		targetUnits, err = formats.ParseUnits(bundle.format, data)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("parsing target %s: %v", targetPath, err))
			return jobOutput{result: result}
		}
		targetKeys = formats.UnitValues(targetUnits)
	case os.IsNotExist(err):
	case err != nil:
		result.Errors = append(result.Errors, fmt.Sprintf("reading target %s: %v", targetPath, err))
		return jobOutput{result: result}
	}

	keys := sortedKeys(sourceKeys)
	sourceUnits := make(map[string]formats.Unit, len(bundle.sourceUnits))
	for _, unit := range bundle.sourceUnits {
		sourceUnits[unit.ID] = unit
	}
	plans := make([]plannedEntry, 0, len(keys))
	for _, key := range keys {
		sourceValue := sourceKeys[key]
		sourceUnit := sourceUnits[key]
		sourceHash := state.SourceUnitHash(bundle.format.Name(), sourceValue, sourceUnit.Context, sourceUnit.Structure)
		targetValue, exists := targetKeys[key]
		recorded, recordedOK := manifest.Get(bundle.bundle.ID, key, locale)
		entryState := classify(exists, targetValue, sourceHash, policyHash, recorded, recordedOK)
		plans = append(plans, plannedEntry{key: key, source: sourceValue, context: sourceUnit.Context, sourceHash: sourceHash, state: entryState})
		result.addState(entryState)
	}
	if opts.AdoptExisting {
		for _, plan := range plans {
			if !plan.state.adoptable() {
				continue
			}
			if err := validateTranslationValue(plan.key, plan.source, targetKeys[plan.key], locale); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("cannot adopt %q: %v", plan.key, err))
			}
		}
		if len(result.Errors) > 0 {
			return jobOutput{result: result}
		}
	}

	candidates := make([]plannedEntry, 0)
	for _, plan := range plans {
		if plan.state.missing || (!plan.state.manualEdit && (plan.state.sourceStale || (plan.state.policyStale && opts.RefreshPolicy))) {
			candidates = append(candidates, plan)
		}
	}
	result.Batches = batchCount(len(candidates), batchSize)
	if opts.DryRun {
		result.KeysSkipped = len(candidates)
		return jobOutput{result: result}
	}
	if opts.AdoptExisting {
		result.KeysSkipped = len(candidates)
		candidates = nil
		result.Batches = 0
	}
	if len(candidates) > 0 && provider == nil {
		result.Errors = append(result.Errors, "translation provider is required")
		return jobOutput{result: result}
	}

	staged := cloneMap(targetKeys)
	origins := make(map[string]translationOrigin)
	toTranslate := make([]plannedEntry, 0, len(candidates))
	for _, plan := range candidates {
		if record, ok := memory.Lookup(locale, bundle.bundle.ID, plan.key, plan.sourceHash, policyHash); ok {
			if err := validateTranslationValue(plan.key, plan.source, record.Target, locale); err == nil {
				staged[plan.key] = record.Target
				origins[plan.key] = translationOrigin{kind: "tm", provider: record.Provider, model: record.Model}
				result.KeysCached++
				continue
			}
		}
		toTranslate = append(toTranslate, plan)
	}
	result.Batches = batchCount(len(toTranslate), batchSize)

	var records []tm.Record
	for i := 0; i < len(toTranslate); i += batchSize {
		end := i + batchSize
		if end > len(toTranslate) {
			end = len(toTranslate)
		}
		batchPlans := toTranslate[i:end]
		entries := make([]llm.Entry, len(batchPlans))
		for j, plan := range batchPlans {
			entries[j] = llm.Entry{Key: plan.key, Value: plan.source, Context: plan.context}
		}

		response, err := provider.Translate(ctx, llm.TranslateRequest{
			SourceLocale: cfg.SourceLocale,
			TargetLocale: locale,
			Entries:      entries,
			SystemPrompt: prompt,
		})
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("batch %d: %v", i/batchSize+1, err))
			return jobOutput{result: result}
		}
		if response == nil {
			result.Errors = append(result.Errors, fmt.Sprintf("batch %d: provider returned no response", i/batchSize+1))
			return jobOutput{result: result}
		}
		if err := validateBatch(entries, response.Translations, locale); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("batch %d: %v", i/batchSize+1, err))
			return jobOutput{result: result}
		}

		result.TokensIn += response.Usage.InputTokens
		result.TokensOut += response.Usage.OutputTokens
		for _, plan := range batchPlans {
			translation := response.Translations[plan.key]
			staged[plan.key] = translation
			origins[plan.key] = translationOrigin{kind: "provider", provider: provider.Name(), model: localeLLM.Model}
			result.KeysTranslated++
			records = append(records, tm.Record{
				Bundle:     bundle.bundle.ID,
				Key:        plan.key,
				Source:     plan.source,
				Target:     translation,
				Locale:     locale,
				Hash:       plan.sourceHash,
				PolicyHash: policyHash,
				Provider:   provider.Name(),
				Model:      localeLLM.Model,
				Timestamp:  time.Now().UTC(),
			})
		}
	}

	changed := len(candidates) > 0
	if changed {
		serializationBaseline := targetData
		if !targetExists {
			serializationBaseline = bundle.sourceData
			if len(optionalPluralKeys) > 0 {
				remover, ok := bundle.format.(formats.EntryRemover)
				if !ok {
					result.Errors = append(result.Errors, fmt.Sprintf("format %q cannot omit source-only plural forms", bundle.format.Name()))
					return jobOutput{result: result}
				}
				serializationBaseline, err = remover.RemoveEntries(serializationBaseline, optionalPluralKeys)
				if err != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("preparing target structure %s: %v", targetPath, err))
					return jobOutput{result: result}
				}
			}
		}
		unitBaseline := targetUnits
		if !targetExists {
			unitBaseline = bundle.sourceUnits
		}
		stagedUnits := formats.MergeUnitValues(unitBaseline, bundle.sourceUnits, staged)
		output, err := formats.SerializeUnits(bundle.format, stagedUnits, serializationBaseline)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("serializing target %s: %v", targetPath, err))
			return jobOutput{result: result}
		}
		output = appendOneNewline(output)
		parsedUnits, err := formats.ParseUnits(bundle.format, output)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("validating staged target %s: %v", targetPath, err))
			return jobOutput{result: result}
		}
		parsed := formats.UnitValues(parsedUnits)
		for _, plan := range candidates {
			if parsed[plan.key] != staged[plan.key] {
				result.Errors = append(result.Errors, fmt.Sprintf("validating staged target %s: key %q changed during serialization", targetPath, plan.key))
				return jobOutput{result: result}
			}
		}
		if err := state.WriteFileAtomic(targetPath, output, 0o644); err != nil {
			result.Errors = append(result.Errors, err.Error())
			return jobOutput{result: result}
		}
		targetExists = true
	}

	if len(records) > 0 {
		if err := memory.AddBatch(records); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("TM update: %v", err))
		}
	}

	now := time.Now().UTC()
	updates := make([]state.Entry, 0)
	if targetExists {
		for _, plan := range plans {
			targetValue, exists := staged[plan.key]
			if !exists {
				continue
			}
			origin := origins[plan.key]
			if plan.state.adoptable() && opts.AdoptExisting {
				origin = translationOrigin{kind: "adopted"}
			}
			if origin.kind == "" {
				continue
			}
			updates = append(updates, state.Entry{
				Bundle:       bundle.bundle.ID,
				Key:          plan.key,
				Locale:       locale,
				SourceHash:   plan.sourceHash,
				PolicyHash:   policyHash,
				TargetHash:   state.TargetHash(targetValue),
				Origin:       origin.kind,
				Provider:     origin.provider,
				Model:        origin.model,
				ReviewStatus: state.ReviewNeedsReview,
				UpdatedAt:    now,
			})
		}
	}
	return jobOutput{result: result, updates: updates}
}

type entryState struct {
	missing     bool
	sourceStale bool
	policyStale bool
	manualEdit  bool
	untracked   bool
	current     bool
}

func (s entryState) adoptable() bool {
	return !s.missing && !s.current && (s.untracked || s.manualEdit || s.sourceStale || s.policyStale)
}

type translationOrigin struct {
	kind     string
	provider string
	model    string
}

type plannedEntry struct {
	key        string
	source     string
	context    string
	sourceHash string
	state      entryState
}

func classify(targetExists bool, target string, sourceHash string, policyHash string, recorded state.Entry, recordedOK bool) entryState {
	if !targetExists {
		return entryState{missing: true}
	}
	if !recordedOK {
		return entryState{untracked: true}
	}
	classified := entryState{
		manualEdit:  recorded.TargetHash != state.TargetHash(target),
		sourceStale: recorded.SourceHash != sourceHash,
		policyStale: recorded.PolicyHash != policyHash,
	}
	classified.current = !classified.manualEdit && !classified.sourceStale && !classified.policyStale
	return classified
}

func (r *Result) addState(state entryState) {
	if state.missing {
		r.KeysMissing++
	}
	if state.sourceStale {
		r.KeysSourceStale++
	}
	if state.policyStale {
		r.KeysPolicyStale++
	}
	if state.manualEdit {
		r.KeysManualEdit++
	}
	if state.untracked {
		r.KeysUntracked++
	}
	if state.current {
		r.KeysCurrent++
	}
}

func validateBatch(entries []llm.Entry, translations map[string]string, targetLocale string) error {
	expected := make(map[string]struct{}, len(entries))
	var missing []string
	for _, entry := range entries {
		expected[entry.Key] = struct{}{}
		if _, ok := translations[entry.Key]; !ok {
			missing = append(missing, entry.Key)
		}
	}
	var extra []string
	for key := range translations {
		if _, ok := expected[key]; !ok {
			extra = append(extra, key)
		}
	}
	if len(missing) == 0 && len(extra) == 0 {
		for _, entry := range entries {
			if err := validateTranslationValue(entry.Key, entry.Value, translations[entry.Key], targetLocale); err != nil {
				return err
			}
		}
		return nil
	}
	sort.Strings(missing)
	sort.Strings(extra)
	return fmt.Errorf("provider response key set mismatch (missing: %v, extra: %v)", missing, extra)
}

func sortedKeys(entries map[string]string) []string {
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func cloneMap(input map[string]string) map[string]string {
	output := make(map[string]string, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func batchCount(entries, batchSize int) int {
	if entries == 0 || batchSize <= 0 {
		return 0
	}
	return (entries + batchSize - 1) / batchSize
}

func appendOneNewline(data []byte) []byte {
	return append([]byte(strings.TrimRight(string(data), "\n")), '\n')
}

func hasUpdates(outputs []jobOutput) bool {
	for _, output := range outputs {
		if len(output.updates) > 0 {
			return true
		}
	}
	return false
}

// FormatResults returns a human-readable summary of translation results.
func FormatResults(results []Result, elapsed time.Duration) string {
	var translated, cached, missing, sourceStale, policyStale, manual, untracked, tokensIn, tokensOut int
	var hasErrors bool
	for _, result := range results {
		translated += result.KeysTranslated
		cached += result.KeysCached
		missing += result.KeysMissing
		sourceStale += result.KeysSourceStale
		policyStale += result.KeysPolicyStale
		manual += result.KeysManualEdit
		untracked += result.KeysUntracked
		tokensIn += result.TokensIn
		tokensOut += result.TokensOut
		hasErrors = hasErrors || len(result.Errors) > 0
	}

	summary := fmt.Sprintf("\nTranslated %d keys across %d bundle/locale jobs (%d from cache) in %s\n", translated, len(results), cached, elapsed.Round(time.Millisecond))
	summary += fmt.Sprintf("Observed before run: %d missing, %d source-stale, %d policy-stale, %d manual edits, %d untracked\n", missing, sourceStale, policyStale, manual, untracked)
	if tokensIn > 0 || tokensOut > 0 {
		summary += fmt.Sprintf("Tokens: %d input, %d output\n", tokensIn, tokensOut)
	}
	if missing+sourceStale+policyStale+manual+untracked > 0 {
		summary += "\nPre-run observations:\n"
		for _, result := range results {
			if result.KeysMissing+result.KeysSourceStale+result.KeysPolicyStale+result.KeysManualEdit+result.KeysUntracked == 0 {
				continue
			}
			summary += fmt.Sprintf("  [%s/%s] %s: %d missing, %d source-stale, %d policy-stale, %d manual edits, %d untracked\n",
				result.Bundle, result.Locale, result.TargetPath, result.KeysMissing, result.KeysSourceStale, result.KeysPolicyStale, result.KeysManualEdit, result.KeysUntracked)
		}
	}
	if hasErrors {
		summary += "\nErrors:\n"
		for _, result := range results {
			for _, runErr := range result.Errors {
				summary += fmt.Sprintf("  [%s/%s] %s\n", result.Bundle, result.Locale, runErr)
			}
		}
	}
	return summary
}
