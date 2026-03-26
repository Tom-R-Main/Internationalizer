<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-native internationalization pipeline for software projects. Translate, validate, and manage i18n files using LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Why Internationalizer?

Most i18n tools are either runtime libraries (i18next, react-intl) or key-management SaaS platforms (Crowdin, Lokalise). None of them solve the actual translation problem well:

- **Manual translation** doesn't scale past a few languages
- **Machine translation APIs** (Google Translate, DeepL) ignore your terminology, tone, and UI conventions
- **Generic LLM translation** works better, but without glossaries and style guides, you get inconsistent results

Internationalizer is different. It's a **CLI pipeline** that combines LLM translation with:

- **Per-language glossaries** — enforce consistent terminology across your app
- **Per-language style guides** — control tone, formality, pluralization, and typography
- **Translation memory** — skip unchanged strings, save money on API calls
- **Key validation** — catch missing translations and interpolation mismatches before they ship

## Installation

Install from npm:

```bash
npm install -g internationalizer
```

Or run without a global install:

```bash
npx internationalizer --help
```

The npm package installs the matching prebuilt binary from npm via platform-specific optional dependencies.

Install with Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Or build from source:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm Packages

- Git tags and npm package versions must match, for example `v0.1.0` and `0.1.0`
- The root `internationalizer` package depends on platform packages such as `internationalizer-darwin-arm64`
- Supported npm targets: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI publishing requires a GitHub secret named `NPM_TOKEN`

## Quick Start

1. Create a config file in your project root:

```yaml
# .internationalizer.yml
source_locale: en
target_locales: [fr, de, es, ja]
source_path: locales/en.json

llm:
  provider: gemini
  model: gemini-3.1-pro-preview
  api_key_env: GOOGLE_AI_STUDIO_API_KEY
```

2. Set your API key:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Preview what will be translated:

```bash
internationalizer translate --dry-run
```

4. Run the translation:

```bash
internationalizer translate
```

5. Validate all locales:

```bash
internationalizer validate
```

## Commands

### `translate`

Find missing keys and translate them via an LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Check all locale files for missing keys, extra keys, and interpolation mismatches.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Auto-detect the i18n framework and suggest a configuration.

```bash
internationalizer detect
```

Supports: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs.

### `glossary`

Manage per-language glossary terms that are enforced during translation.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Manage translation memory (JSONL cache of previously translated strings).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Configuration Reference

```yaml
# .internationalizer.yml

# Source language (default: en)
source_locale: en

# Languages to translate into (required)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Path to the source locale file (required)
source_path: locales/en.json

# LLM provider settings
llm:
  # Provider: "anthropic", "openai", "gemini", or "openrouter" (default: gemini)
  provider: gemini

  # Model name defaults by provider:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Environment variable containing the API key
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Base URL for OpenAI-compatible endpoints (optional)
  # base_url: https://api.openai.com

# Keys per LLM call (default: 40)
batch_size: 40

# Parallel LLM calls (default: 4)
concurrency: 4

# Directory containing per-locale style guide Markdown files (default: style-guides)
style_guides_dir: style-guides

# Directory containing per-locale glossary JSON files (default: glossary)
glossary_dir: glossary

# Path to translation memory file (default: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Style Guides

Style guides are Markdown files that get injected into the LLM translation prompt. They control tone, formality, typography, and other language-specific conventions.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Shared conventions (`_conventions.md`)

Define rules that apply to all languages: interpolation syntax, HTML preservation, string type conventions (buttons vs. labels vs. errors), etc.

### Per-language guides (`{locale}.md`)

Define language-specific rules: formality register (tu vs. vous), punctuation (guillemets, inverted question marks), plural forms, date/number formatting, and a terminology glossary.

See [`examples/react-app/style-guides/`](examples/react-app/style-guides/) for a working example.

## Glossary Format

Glossary files are JSON arrays stored in `{glossary_dir}/{locale}.json`:

```json
[
  {
    "source": "Dashboard",
    "target": "Tableau de bord",
    "ignore_case": false,
    "whole_word": true
  }
]
```

Terms are injected into the LLM prompt as a terminology table, ensuring consistent translation of key terms across your application.

## Translation Memory

Translation memory is stored as a JSONL file (one JSON record per line). Each record contains:

- The source key and value
- The translated value
- A SHA-256 hash of the source value
- A timestamp

On subsequent runs, unchanged strings are served from the TM cache without calling the LLM, saving both time and API costs. The TM file is git-friendly and can be committed alongside your locale files.

## Supported Formats

| Format | Extensions | Mode |
|--------|-----------|------|
| JSON | `.json` | Key-value (nested, dot-notation flattened) |
| YAML | `.yml`, `.yaml` | Key-value (preserves comments and ordering) |
| Markdown | `.md`, `.mdx` | Whole-document translation |

## Project Type Detection

`internationalizer detect` identifies your i18n setup by checking:

- `package.json` dependencies for react-i18next, next-intl, or vue-i18n
- Directory structures matching common locale patterns
- File extensions and naming conventions

## Architecture

```
cmd/internationalizer/     CLI entry point and command definitions
internal/
  config/                  YAML config loading with defaults
  detect/                  Project type auto-detection
  formats/                 Format parsers (JSON, YAML, Markdown)
  glossary/                Per-locale glossary management
  llm/                     LLM provider interface + implementations
    anthropic.go           Anthropic Claude backend
    openai.go              OpenAI / compatible backend
    gemini.go              Google Gemini via AI Studio backend
                           OpenRouter uses openai.go with custom base_url
  styleguide/              Style guide loader
  tm/                      JSONL translation memory
  translate/               Translation orchestrator
  validate/                Locale validation and diffing
```

## Comparison to Alternatives

| Feature | Internationalizer | i18next | Crowdin | Generic LLM |
|---------|------------------|---------|---------|-------------|
| LLM-powered translation | Yes | No | Partial | Yes |
| Per-language style guides | Yes | No | No | No |
| Glossary enforcement | Yes | No | Yes | No |
| Translation memory | Yes | No | Yes | No |
| CLI / local execution | Yes | N/A | No | Manual |
| Git-friendly files | Yes | Yes | Partial | Manual |
| No SaaS dependency | Yes | Yes | No | Varies |
| Open source (AGPL-3.0) | Yes | Yes | No | Varies |

## License

[AGPL-3.0](LICENSE)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines. All contributions require DCO sign-off.
