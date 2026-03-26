package translate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/formats"
	"github.com/Tom-R-Main/Internationalizer/internal/glossary"
	"github.com/Tom-R-Main/Internationalizer/internal/llm"
	"github.com/Tom-R-Main/Internationalizer/internal/styleguide"
	"github.com/Tom-R-Main/Internationalizer/internal/tm"
)

// Options configures a translation run.
type Options struct {
	DryRun      bool
	Locales     []string // filter to specific locales; empty = all from config
	BatchSize   int      // override config; 0 = use config
	Concurrency int      // override config; 0 = use config
}

// Result holds the outcome of a translation run.
type Result struct {
	Locale         string
	KeysTotal      int
	KeysCached     int
	KeysTranslated int
	KeysSkipped    int
	Batches        int
	TokensIn       int
	TokensOut      int
	Errors         []string
}

// Run executes the translation pipeline.
func Run(ctx context.Context, cfg *config.Config, provider llm.Provider, opts Options) ([]Result, error) {
	batchSize := cfg.BatchSize
	if opts.BatchSize > 0 {
		batchSize = opts.BatchSize
	}
	concurrency := cfg.Concurrency
	if opts.Concurrency > 0 {
		concurrency = opts.Concurrency
	}

	// Determine target locales.
	locales := cfg.TargetLocales
	if len(opts.Locales) > 0 {
		locales = opts.Locales
	}

	// Detect format.
	format, err := formats.FormatForFile(cfg.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("detecting format: %w", err)
	}

	// Load source.
	sourceData, err := os.ReadFile(cfg.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("reading source %s: %w", cfg.SourcePath, err)
	}
	sourceKeys, err := format.Parse(sourceData)
	if err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}

	// Load translation memory.
	memory, err := tm.Load(cfg.TMPath)
	if err != nil {
		return nil, fmt.Errorf("loading TM: %w", err)
	}

	// Process locales with concurrency limit.
	sem := make(chan struct{}, concurrency)
	var mu sync.Mutex
	var results []Result

	var wg sync.WaitGroup
	for _, locale := range locales {
		wg.Add(1)
		go func(locale string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := translateLocale(ctx, cfg, format, provider, memory, sourceKeys, sourceData, locale, batchSize, opts.DryRun)
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(locale)
	}
	wg.Wait()

	return results, nil
}

func translateLocale(
	ctx context.Context,
	cfg *config.Config,
	format formats.Format,
	provider llm.Provider,
	memory *tm.TM,
	sourceKeys map[string]string,
	sourceData []byte,
	locale string,
	batchSize int,
	dryRun bool,
) Result {
	result := Result{Locale: locale}
	sourceDir := filepath.Dir(cfg.SourcePath)
	ext := filepath.Ext(cfg.SourcePath)
	targetPath := filepath.Join(sourceDir, locale+ext)

	// Load existing target translations.
	targetKeys := make(map[string]string)
	var targetData []byte
	if data, err := os.ReadFile(targetPath); err == nil {
		targetData = data
		if parsed, err := format.Parse(data); err == nil {
			targetKeys = parsed
		}
	}

	// Find missing keys.
	var missing []llm.Entry
	cached := make(map[string]string)
	for key, sourceVal := range sourceKeys {
		if _, exists := targetKeys[key]; exists {
			continue
		}
		result.KeysTotal++

		// Check TM.
		hash := tm.HashSource(sourceVal)
		if translation, ok := memory.Lookup(locale, key, hash); ok {
			cached[key] = translation
			result.KeysCached++
			continue
		}

		missing = append(missing, llm.Entry{Key: key, Value: sourceVal})
	}

	result.KeysTotal += result.KeysCached
	result.Batches = (len(missing) + batchSize - 1) / batchSize
	if batchSize == 0 {
		result.Batches = 0
	}

	fmt.Fprintf(os.Stderr, "[%s] %d missing, %d cached, %d to translate in %d batches\n",
		locale, len(missing)+result.KeysCached, result.KeysCached, len(missing), result.Batches)

	if dryRun {
		result.KeysSkipped = len(missing)
		return result
	}

	// Load glossary and style guide for this locale.
	terms, _ := glossary.Load(cfg.GlossaryDir, locale)
	guide, _ := styleguide.Load(cfg.StyleGuidesDir, locale)
	systemPrompt := llm.BuildSystemPrompt(cfg.SourceLocale, locale, guide, terms)

	// Translate batches.
	var tmRecords []tm.Record
	for i := 0; i < len(missing); i += batchSize {
		end := i + batchSize
		if end > len(missing) {
			end = len(missing)
		}
		batch := missing[i:end]

		resp, err := provider.Translate(ctx, llm.TranslateRequest{
			SourceLocale: cfg.SourceLocale,
			TargetLocale: locale,
			Entries:      batch,
			SystemPrompt: systemPrompt,
		})
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("batch %d: %v", i/batchSize+1, err))
			continue
		}

		result.TokensIn += resp.Usage.InputTokens
		result.TokensOut += resp.Usage.OutputTokens

		for _, entry := range batch {
			if translation, ok := resp.Translations[entry.Key]; ok {
				targetKeys[entry.Key] = translation
				result.KeysTranslated++

				tmRecords = append(tmRecords, tm.Record{
					Key:       entry.Key,
					Source:    entry.Value,
					Target:    translation,
					Locale:    locale,
					Hash:      tm.HashSource(entry.Value),
					Timestamp: time.Now(),
				})
			}
		}
	}

	// Apply cached translations.
	for key, val := range cached {
		targetKeys[key] = val
	}

	// Write updated target file.
	output, err := format.Serialize(targetKeys, targetData)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("serialize: %v", err))
		return result
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("mkdir: %v", err))
		return result
	}
	if err := os.WriteFile(targetPath, append(output, '\n'), 0o644); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("write: %v", err))
		return result
	}

	// Update TM.
	if len(tmRecords) > 0 {
		if err := memory.AddBatch(tmRecords); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("tm update: %v", err))
		}
	}

	return result
}

// FormatResults returns a human-readable summary of translation results.
func FormatResults(results []Result, elapsed time.Duration) string {
	var totalTranslated, totalCached, totalIn, totalOut int
	var hasErrors bool

	for _, r := range results {
		totalTranslated += r.KeysTranslated
		totalCached += r.KeysCached
		totalIn += r.TokensIn
		totalOut += r.TokensOut
		if len(r.Errors) > 0 {
			hasErrors = true
		}
	}

	summary := fmt.Sprintf("\nTranslated %d keys across %d locales (%d from cache) in %s\n",
		totalTranslated, len(results), totalCached, elapsed.Round(time.Millisecond))
	if totalIn > 0 || totalOut > 0 {
		summary += fmt.Sprintf("Tokens: %d input, %d output\n", totalIn, totalOut)
	}

	if hasErrors {
		summary += "\nErrors:\n"
		for _, r := range results {
			for _, e := range r.Errors {
				summary += fmt.Sprintf("  [%s] %s\n", r.Locale, e)
			}
		}
	}
	return summary
}
