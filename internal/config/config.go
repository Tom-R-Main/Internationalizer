package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Provider model defaults track the production translation choices in
// ExecuFunction while remaining overrideable in project configuration.
const (
	DefaultAnthropicModel  = "claude-opus-5"
	DefaultOpenAIModel     = "gpt-5.6-luna"
	DefaultOpenAIReasoning = "max"
	DefaultGeminiModel     = "gemini-3.7-flash"
	DefaultOpenRouterModel = "deepseek/deepseek-v4-pro-0813"
)

type Config struct {
	SourceLocale   string   `yaml:"source_locale"`
	TargetLocales  []string `yaml:"target_locales"`
	SourcePath     string   `yaml:"source_path"`
	Bundles        []Bundle `yaml:"bundles"`
	LLM            LLM      `yaml:"llm"`
	BatchSize      int      `yaml:"batch_size"`
	Concurrency    int      `yaml:"concurrency"`
	StyleGuidesDir string   `yaml:"style_guides_dir"`
	GlossaryDir    string   `yaml:"glossary_dir"`
	TMPath         string   `yaml:"tm_path"`
	ManifestPath   string   `yaml:"manifest_path"`
	Formats        []string `yaml:"formats"`
}

// Bundle maps one source file to a locale-specific target path.
// Target must contain the literal {locale} placeholder.
type Bundle struct {
	ID     string `yaml:"id"` // required stable identity for explicit bundles
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Format string `yaml:"format"`
}

type LLM struct {
	Provider        string `yaml:"provider"`
	Model           string `yaml:"model"`
	APIKeyEnv       string `yaml:"api_key_env"`
	BaseURL         string `yaml:"base_url"`
	ReasoningEffort string `yaml:"reasoning_effort"`
}

func (c *Config) ApplyDefaults() {
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
	if c.LLM.Provider == "" {
		c.LLM.Provider = "gemini"
	}
	if c.LLM.Model == "" {
		switch c.LLM.Provider {
		case "anthropic":
			c.LLM.Model = DefaultAnthropicModel
		case "openai":
			c.LLM.Model = DefaultOpenAIModel
		case "gemini":
			c.LLM.Model = DefaultGeminiModel
		case "openrouter":
			c.LLM.Model = DefaultOpenRouterModel
		}
	}
	if c.LLM.APIKeyEnv == "" {
		switch c.LLM.Provider {
		case "anthropic":
			c.LLM.APIKeyEnv = "ANTHROPIC_API_KEY"
		case "openai":
			c.LLM.APIKeyEnv = "OPENAI_API_KEY"
		case "gemini":
			c.LLM.APIKeyEnv = "GOOGLE_AI_STUDIO_API_KEY"
		case "openrouter":
			c.LLM.APIKeyEnv = "OPENROUTER_API_KEY"
		}
	}
	if c.LLM.Provider == "openai" && c.LLM.ReasoningEffort == "" {
		c.LLM.ReasoningEffort = DefaultOpenAIReasoning
	}
}

func (c *Config) Validate() error {
	if err := c.ValidateProject(); err != nil {
		return err
	}
	return c.ValidateCredentials()
}

// ValidateProject checks configuration that is required even for dry runs.
func (c *Config) ValidateProject() error {
	if len(c.TargetLocales) == 0 {
		return fmt.Errorf("target_locales must not be empty")
	}
	localeSeen := make(map[string]struct{}, len(c.TargetLocales))
	for _, locale := range c.TargetLocales {
		if !validLocale.MatchString(locale) {
			return fmt.Errorf("invalid target locale %q", locale)
		}
		if _, ok := localeSeen[locale]; ok {
			return fmt.Errorf("duplicate target locale %q", locale)
		}
		localeSeen[locale] = struct{}{}
	}
	bundles := c.EffectiveBundles()
	if len(bundles) == 0 {
		return fmt.Errorf("source_path or bundles is required")
	}
	seen := make(map[string]struct{}, len(bundles))
	targets := make(map[string]string, len(bundles)*len(c.TargetLocales))
	for _, bundle := range bundles {
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
	if c.LLM.APIKeyEnv != "" && os.Getenv(c.LLM.APIKeyEnv) == "" {
		return fmt.Errorf("environment variable %s is not set", c.LLM.APIKeyEnv)
	}
	return nil
}

var validLocale = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

// EffectiveBundles returns explicit bundles, or a backward-compatible bundle
// derived from source_path and the historical sibling target convention.
func (c *Config) EffectiveBundles() []Bundle {
	if len(c.Bundles) > 0 {
		bundles := make([]Bundle, len(c.Bundles))
		copy(bundles, c.Bundles)
		return bundles
	}
	if c.SourcePath == "" {
		return nil
	}
	return []Bundle{{
		ID:     "default",
		Source: c.SourcePath,
		Target: filepath.Join(filepath.Dir(c.SourcePath), "{locale}"+filepath.Ext(c.SourcePath)),
	}}
}

// TargetPath resolves a bundle target for a validated locale.
func (b Bundle) TargetPath(locale string) (string, error) {
	if !validLocale.MatchString(locale) {
		return "", fmt.Errorf("invalid target locale %q", locale)
	}
	if !strings.Contains(b.Target, "{locale}") {
		return "", fmt.Errorf("bundle %q target must contain {locale}", b.ID)
	}
	return filepath.Clean(strings.ReplaceAll(b.Target, "{locale}", locale)), nil
}

// APIKey returns the resolved API key from the environment.
func (c *Config) APIKey() string {
	return os.Getenv(c.LLM.APIKeyEnv)
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
