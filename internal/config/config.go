package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SourceLocale   string   `yaml:"source_locale"`
	TargetLocales  []string `yaml:"target_locales"`
	SourcePath     string   `yaml:"source_path"`
	LLM            LLM      `yaml:"llm"`
	BatchSize      int      `yaml:"batch_size"`
	Concurrency    int      `yaml:"concurrency"`
	StyleGuidesDir string   `yaml:"style_guides_dir"`
	GlossaryDir    string   `yaml:"glossary_dir"`
	TMPath         string   `yaml:"tm_path"`
	Formats        []string `yaml:"formats"`
}

type LLM struct {
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"api_key_env"`
	BaseURL   string `yaml:"base_url"`
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
			c.LLM.Model = "claude-sonnet-4-6"
		case "openai":
			c.LLM.Model = "gpt-5.4"
		case "gemini":
			c.LLM.Model = "gemini-3.1-pro-preview"
		case "openrouter":
			c.LLM.Model = "google/gemini-3-flash-preview"
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
}

func (c *Config) Validate() error {
	if len(c.TargetLocales) == 0 {
		return fmt.Errorf("target_locales must not be empty")
	}
	if c.SourcePath == "" {
		return fmt.Errorf("source_path is required")
	}
	if c.LLM.APIKeyEnv != "" && os.Getenv(c.LLM.APIKeyEnv) == "" {
		return fmt.Errorf("environment variable %s is not set", c.LLM.APIKeyEnv)
	}
	return nil
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
