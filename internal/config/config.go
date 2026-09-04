package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	localeid "github.com/Tom-R-Main/Internationalizer/internal/locale"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"gopkg.in/yaml.v3"
)

// Provider model defaults track the production translation choices in
// ExecuFunction while remaining overrideable in project configuration.
const (
	DefaultAnthropicModel  = "claude-opus-5"
	DefaultOpenAIModel     = "gpt-5.6-luna"
	DefaultOpenAIReasoning = "max"
	DefaultGeminiModel     = "gemini-3.8-flash"
	DefaultOpenRouterModel = "deepseek/deepseek-v4-pro-0813"
)

type Config struct {
	MessageSyntax  message.Syntax `yaml:"message_syntax"`
	SourceLocale   string         `yaml:"source_locale"`
	TargetLocales  []string       `yaml:"target_locales"`
	SourcePath     string         `yaml:"source_path"`
	Bundles        []Bundle       `yaml:"bundles"`
	LLM            LLM            `yaml:"llm"`
	BatchSize      int            `yaml:"batch_size"`
	Concurrency    int            `yaml:"concurrency"`
	StyleGuidesDir string         `yaml:"style_guides_dir"`
	GlossaryDir    string         `yaml:"glossary_dir"`
	TMPath         string         `yaml:"tm_path"`
	ManifestPath   string         `yaml:"manifest_path"`
	Formats        []string       `yaml:"formats"`
	Validation     Validation     `yaml:"validation"`
}

// Validation configures optional project-specific validation rules.
type Validation struct {
	PluralStyle string `yaml:"plural_style"`
}

// Bundle maps one source file to a locale-specific target path.
// Target must contain the literal {locale} placeholder.
type Bundle struct {
	MessageSyntax message.Syntax `yaml:"message_syntax"`
	ID            string         `yaml:"id"` // required stable identity for explicit bundles
	Source        string         `yaml:"source"`
	Target        string         `yaml:"target"`
	Format        string         `yaml:"format"`
}

type LLM struct {
	Provider        string                 `yaml:"provider"`
	Model           string                 `yaml:"model"`
	APIKeyEnv       string                 `yaml:"api_key_env"`
	BaseURL         string                 `yaml:"base_url"`
	ReasoningEffort string                 `yaml:"reasoning_effort"`
	LocaleOverrides map[string]LLMOverride `yaml:"locale_overrides"`
}

// LLMOverride replaces LLM settings for one target locale. Fields omitted
// while keeping the global provider inherit their global values.
type LLMOverride struct {
	Provider        string `yaml:"provider"`
	Model           string `yaml:"model"`
	APIKeyEnv       string `yaml:"api_key_env"`
	BaseURL         string `yaml:"base_url"`
	ReasoningEffort string `yaml:"reasoning_effort"`
}

func (c *Config) ApplyDefaults() {
	if c.MessageSyntax == "" {
		c.MessageSyntax = message.Auto
	}
	if c.SourceLocale == "" {
		c.SourceLocale = "en"
	}
	if c.BatchSize <= 0 {
		c.BatchSize = 40
	}
	if c.Concurrency <= 0 {
		c.Concurrency = 4
	}
	if c.TMPath == "" {
		c.TMPath = ".internationalizer/tm.jsonl"
	}
	if c.ManifestPath == "" {
		c.ManifestPath = ".internationalizer.lock"
	}
	if c.StyleGuidesDir == "" {
		c.StyleGuidesDir = "style-guides"
	}
	if c.GlossaryDir == "" {
		c.GlossaryDir = "glossary"
	}
	applyLLMDefaults(&c.LLM)
}

func applyLLMDefaults(settings *LLM) {
	if settings.Provider == "" {
		settings.Provider = "gemini"
	}
	if settings.Model == "" {
		switch settings.Provider {
		case "anthropic":
			settings.Model = DefaultAnthropicModel
		case "openai":
			settings.Model = DefaultOpenAIModel
		case "gemini":
			settings.Model = DefaultGeminiModel
		case "openrouter":
			settings.Model = DefaultOpenRouterModel
		}
	}
	if settings.APIKeyEnv == "" {
		switch settings.Provider {
		case "anthropic":
			settings.APIKeyEnv = "ANTHROPIC_API_KEY"
		case "openai":
			settings.APIKeyEnv = "OPENAI_API_KEY"
		case "gemini":
			settings.APIKeyEnv = "GOOGLE_AI_STUDIO_API_KEY"
		case "openrouter":
			settings.APIKeyEnv = "OPENROUTER_API_KEY"
		}
	}
	if settings.Provider == "openai" && settings.ReasoningEffort == "" {
		settings.ReasoningEffort = DefaultOpenAIReasoning
	}
}

// LLMForLocale returns the fully resolved LLM settings for a target locale.
func (c *Config) LLMForLocale(locale string) LLM {
	global := c.LLM
	global.LocaleOverrides = nil
	applyLLMDefaults(&global)

	override, ok := localeOverride(c.LLM.LocaleOverrides, locale)
	if !ok {
		return global
	}
	effective := LLM{
		Provider:        override.Provider,
		Model:           override.Model,
		APIKeyEnv:       override.APIKeyEnv,
		BaseURL:         override.BaseURL,
		ReasoningEffort: override.ReasoningEffort,
	}
	if effective.Provider == "" {
		effective.Provider = global.Provider
	}
	if effective.Provider == global.Provider {
		if effective.Model == "" {
			effective.Model = global.Model
		}
		if effective.APIKeyEnv == "" {
			effective.APIKeyEnv = global.APIKeyEnv
		}
		if effective.BaseURL == "" {
			effective.BaseURL = global.BaseURL
		}
		if effective.ReasoningEffort == "" {
			effective.ReasoningEffort = global.ReasoningEffort
		}
	}
	applyLLMDefaults(&effective)
	return effective
}

// ConfiguredTargetLocale returns the configured spelling of a target locale
// that is canonically equivalent to requested.
func (c *Config) ConfiguredTargetLocale(requested string) (string, bool) {
	canonical, err := localeid.Canonical(requested)
	if err != nil {
		return "", false
	}
	for _, candidate := range c.TargetLocales {
		candidateCanonical, candidateErr := localeid.Canonical(candidate)
		if candidateErr == nil && candidateCanonical == canonical {
			return candidate, true
		}
	}
	return "", false
}

// HasLLMOverrideForLocale reports whether a canonical-equivalent locale has
// explicit provider settings.
func (c *Config) HasLLMOverrideForLocale(locale string) bool {
	_, ok := localeOverride(c.LLM.LocaleOverrides, locale)
	return ok
}

func (c *Config) Validate() error {
	if err := c.ValidateProject(); err != nil {
		return err
	}
	return c.ValidateCredentials()
}

// ValidateProject checks configuration that is required even for dry runs.
func (c *Config) ValidateProject() error {
	if c.Validation.PluralStyle != "" && c.Validation.PluralStyle != "i18next-v4" {
		return fmt.Errorf("unsupported validation.plural_style %q", c.Validation.PluralStyle)
	}
	sourceLocale := c.SourceLocale
	if sourceLocale == "" {
		sourceLocale = "en"
	}
	if _, err := localeid.Canonical(sourceLocale); err != nil {
		return fmt.Errorf("invalid source locale %q: %w", c.SourceLocale, err)
	}
	if len(c.TargetLocales) == 0 {
		return fmt.Errorf("target_locales must not be empty")
	}
	localeSeen := make(map[string]struct{}, len(c.TargetLocales))
	for _, locale := range c.TargetLocales {
		canonical, err := localeid.Canonical(locale)
		if err != nil {
			return fmt.Errorf("invalid target locale %q: %w", locale, err)
		}
		if _, ok := localeSeen[canonical]; ok {
			return fmt.Errorf("duplicate target locale %q", locale)
		}
		localeSeen[canonical] = struct{}{}
	}
	overrideSeen := make(map[string]struct{}, len(c.LLM.LocaleOverrides))
	for locale := range c.LLM.LocaleOverrides {
		canonical, err := localeid.Canonical(locale)
		if err != nil {
			return fmt.Errorf("invalid llm locale override %q: %w", locale, err)
		}
		if _, duplicate := overrideSeen[canonical]; duplicate {
			return fmt.Errorf("duplicate llm locale override %q", locale)
		}
		overrideSeen[canonical] = struct{}{}
		if _, ok := localeSeen[canonical]; !ok {
			return fmt.Errorf("llm locale override %q is not in target_locales", locale)
		}
	}
	bundles := c.EffectiveBundles()
	if err := message.ValidateSyntax(c.MessageSyntax); err != nil {
		return err
	}
	if len(bundles) == 0 {
		return fmt.Errorf("source_path or bundles is required")
	}
	seen := make(map[string]struct{}, len(bundles))
	targets := make(map[string]string, len(bundles)*len(c.TargetLocales))
	for _, bundle := range bundles {
		if err := message.ValidateSyntax(bundle.MessageSyntax); err != nil {
			return fmt.Errorf("bundle %q: %w", bundle.ID, err)
		}
		format := strings.ToLower(bundle.Format)
		extension := strings.ToLower(filepath.Ext(bundle.Source))
		if format == "" && extension == ".ftl" {
			format = "fluent"
		}
		if format == "fluent" && bundle.MessageSyntax != "" && bundle.MessageSyntax != message.Auto {
			return fmt.Errorf("bundle %q: Fluent resources require message_syntax: auto", bundle.ID)
		}
		if (format == "markdown" || (format == "" && (extension == ".md" || extension == ".mdx"))) && bundle.MessageSyntax == message.ICU {
			return fmt.Errorf("bundle %q: Markdown documents cannot use message_syntax: icu", bundle.ID)
		}
		if bundle.ID == "" {
			return fmt.Errorf("explicit bundle id is required")
		}
		if bundle.Source == "" {
			return fmt.Errorf("bundle %q source is required", bundle.ID)
		}
		if bundle.Target == "" || !strings.Contains(bundle.Target, "{locale}") {
			return fmt.Errorf("bundle %q target must contain {locale}", bundle.ID)
		}
		if _, ok := seen[bundle.ID]; ok {
			return fmt.Errorf("duplicate bundle id %q", bundle.ID)
		}
		seen[bundle.ID] = struct{}{}
		for _, locale := range c.TargetLocales {
			target, err := bundle.TargetPath(locale)
			if err != nil {
				return err
			}
			if filepath.Clean(bundle.Source) == target {
				return fmt.Errorf("bundle %q target for %s resolves to its source path", bundle.ID, locale)
			}
			if owner, ok := targets[target]; ok {
				return fmt.Errorf("bundles %q and %q both target %s", owner, bundle.ID, target)
			}
			targets[target] = bundle.ID
		}
	}
	return nil
}

// ValidateCredentials checks provider credentials for a live translation run.
func (c *Config) ValidateCredentials() error {
	return c.ValidateCredentialsForLocales(nil)
}

// ValidateCredentialsForLocales checks only providers used by the requested
// locales. An empty locale list checks all configured targets.
func (c *Config) ValidateCredentialsForLocales(locales []string) error {
	if len(locales) == 0 {
		locales = c.TargetLocales
	}
	checked := make(map[string]struct{})
	if len(locales) == 0 {
		settings := c.LLM
		settings.LocaleOverrides = nil
		applyLLMDefaults(&settings)
		return validateAPIKeyEnv(settings.APIKeyEnv, checked)
	}
	for _, locale := range locales {
		apiKeyEnv := c.LLMForLocale(locale).APIKeyEnv
		if err := validateAPIKeyEnv(apiKeyEnv, checked); err != nil {
			return err
		}
	}
	return nil
}

func validateAPIKeyEnv(apiKeyEnv string, checked map[string]struct{}) error {
	if apiKeyEnv == "" {
		return nil
	}
	if _, ok := checked[apiKeyEnv]; ok {
		return nil
	}
	checked[apiKeyEnv] = struct{}{}
	if os.Getenv(apiKeyEnv) == "" {
		return fmt.Errorf("environment variable %s is not set", apiKeyEnv)
	}
	return nil
}

// EffectiveBundles returns explicit bundles, or a backward-compatible bundle
// derived from source_path and the historical sibling target convention.
func (c *Config) EffectiveBundles() []Bundle {
	if len(c.Bundles) > 0 {
		bundles := make([]Bundle, len(c.Bundles))
		copy(bundles, c.Bundles)
		for i := range bundles {
			if bundles[i].MessageSyntax == "" {
				bundles[i].MessageSyntax = c.MessageSyntax
			}
		}
		return bundles
	}
	if c.SourcePath == "" {
		return nil
	}
	return []Bundle{{
		MessageSyntax: c.MessageSyntax,
		ID:            "default",
		Source:        c.SourcePath,
		Target:        filepath.Join(filepath.Dir(c.SourcePath), "{locale}"+filepath.Ext(c.SourcePath)),
	}}
}

// TargetPath resolves a bundle target for a validated locale.
func (b Bundle) TargetPath(locale string) (string, error) {
	if _, err := localeid.Canonical(locale); err != nil {
		return "", fmt.Errorf("invalid target locale %q: %w", locale, err)
	}
	if !strings.Contains(b.Target, "{locale}") {
		return "", fmt.Errorf("bundle %q target must contain {locale}", b.ID)
	}
	return filepath.Clean(strings.ReplaceAll(b.Target, "{locale}", locale)), nil
}

// PluralStyle uses v4 key families for i18next; explicit non-i18next grammars
// must not reinterpret keys just because they end in a plural suffix.
func (b Bundle) PluralStyle(legacyStyle string) string {
	switch b.MessageSyntax {
	case message.I18next:
		return "i18next-v4"
	case "", message.Auto:
		return legacyStyle
	default:
		return ""
	}
}

func localeOverride(overrides map[string]LLMOverride, requested string) (LLMOverride, bool) {
	if override, ok := overrides[requested]; ok {
		return override, true
	}
	canonical, err := localeid.Canonical(requested)
	if err != nil {
		return LLMOverride{}, false
	}
	keys := make([]string, 0, len(overrides))
	for key := range overrides {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if candidate, candidateErr := localeid.Canonical(key); candidateErr == nil && candidate == canonical {
			return overrides[key], true
		}
	}
	return LLMOverride{}, false
}

// APIKey returns the resolved API key from the environment.
func (c *Config) APIKey() string {
	return os.Getenv(c.LLM.APIKeyEnv)
}

// APIKeyForLocale resolves the configured credential for a target locale.
func (c *Config) APIKeyForLocale(locale string) string {
	return os.Getenv(c.LLMForLocale(locale).APIKeyEnv)
}

// Load reads the config from the given path, or searches default locations.
func Load(path string) (*Config, error) {
	if path != "" {
		return loadFile(path)
	}
	candidates := []string{
		".internationalizer.yml",
		".internationalizer.yaml",
	}
	home, err := os.UserHomeDir()
	if err == nil {
		candidates = append(candidates,
			filepath.Join(home, ".internationalizer.yml"),
			filepath.Join(home, ".internationalizer.yaml"),
		)
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return loadFile(p)
		}
	}
	return nil, fmt.Errorf("no config file found; create .internationalizer.yml or pass --config")
}

func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	cfg.ApplyDefaults()
	return &cfg, nil
}
