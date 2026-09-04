<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-native internationalization pipeline for software projects. Translate, validate, and manage i18n files using LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/cs.md">Čeština</a> · <a href="docs/i18n/da.md">Dansk</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/el.md">Ελληνικά</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fi.md">Suomi</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/he.md">עברית</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/it.md">Italiano</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/ko.md">한국어</a> · <a href="docs/i18n/ms.md">Bahasa Melayu</a><br><a href="docs/i18n/nl.md">Nederlands</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pl.md">Polski</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ro.md">Română</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/sv.md">Svenska</a> · <a href="docs/i18n/te.md">తెలుగు</a> · <a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/tr.md">Türkçe</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/vi.md">Tiếng Việt</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
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
- **Deterministic validation** — catch missing or extra keys, protected-structure drift, glossary issues, and plural or ICU errors before they ship
- **Explicit approval** — keep provider or adopted provenance separate from human review
- **Fluent and pseudolocales** — preserve translator context and exercise accented or bidirectional layouts without an API call

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
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    format: json

llm:
  provider: gemini
  model: gemini-3.8-flash
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

5. Validate and approve the exact generated artifacts:

```bash
internationalizer validate
internationalizer review list --status needs_review
internationalizer review approve --locale fr --all
internationalizer validate --require-approved
```

## Commands

### `translate`

Find missing or stale keys and translate them via an LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Translation state independently reports missing, source-stale, policy-stale,
current, and manually edited conditions, so a manual edit cannot conceal a
source or policy change. Policy-stale values are reported but only retranslated
with `--refresh-policy`. Manually edited values are never overwritten
automatically. Use `--adopt-existing` when introducing the manifest to existing
translations or when explicitly accepting a manual edit as the new provenance
baseline. Adoption does not imply human approval; use `review approve` for that
separate decision.

### `pseudo`

Generate deterministic test locales without a provider or translation-memory
lookup. Accented output defaults to `en-XA`; bidirectional output defaults to
`ar-XB`. ICU and Fluent runtime syntax, code, links, and markup are preserved.

```bash
internationalizer pseudo                         # accented en-XA
internationalizer pseudo --strategy bidi         # bidi ar-XB
internationalizer pseudo --dry-run                # show planned artifacts only
internationalizer pseudo --locale qps-ploc        # choose another valid locale tag
```

The generator refreshes only artifacts it previously recorded as pseudo
output. Use `--force` to replace any other existing target deliberately.

### `review`

Inspect and approve the exact target content currently bound to its source and
translation policy. Generated, cached, and adopted values all begin in
`needs_review`; approval is invalidated by later source, policy, or target
changes. Pseudolocales are tracked separately as test artifacts.

```bash
internationalizer review list --status needs_review
internationalizer review approve --locale fr --bundle app --key common.save
internationalizer review approve --locale fr --all
```

### `validate`

Check all locale files against their source bundles. Default validation checks
structural coverage (the percentage of required target keys present), reports
extra keys as warnings, and fails for missing keys, interpolation mismatches,
or invalid ICU MessageFormat structure.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
internationalizer validate --require-approved   # also require explicit approval
```

`--strict` also reports translated coverage. A linguistic value identical to
its source is untranslated unless the glossary explicitly contains an exact
same-source, same-target entry for the complete value; `ignore_case` is honored,
but a glossary term embedded in a longer value is not an exemption. Strict mode
fails on extra keys, source-identical values, changed interpolation/HTML/code/
Markdown-link structure, glossary violations, and configured plural forms.

`--require-state` verifies each target against `.internationalizer.lock`. It
fails when a key is untracked, or when its recorded source, translation policy,
or target hash is stale. `--require-approved` implies `--require-state` and also
fails if the exact current artifact has not been approved. Both can be combined
with `--strict`.

Human and JSON reports use stable finding codes:

| Code | Meaning |
| --- | --- |
| `missing_key` / `extra_key` | Source and target key sets differ |
| `blank_translation` | A non-empty source has an empty strict-mode target |
| `source_identical` | A strict-mode linguistic value remains untranslated |
| `protected_structure_mismatch` | Interpolation, HTML, code, or link structure changed |
| `glossary_violation` | No approved target term or variant was found |
| `plural_form_missing` | A configured locale plural form is absent |
| `icu_message_syntax` | A source or target ICU message is malformed |
| `icu_argument_mismatch` | ICU argument names, types, or formatter styles differ |
| `icu_selector_mismatch` | Selectors differ or a plural category is invalid for the target locale |
| `untracked` | No manifest record exists for the target |
| `source_stale` | Source content changed after the recorded translation |
| `policy_stale` | The generated prompt or model settings changed |
| `target_modified` | Target content differs from the manifest record |
| `needs_review` | Current provenance exists, but the exact target is not approved |

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
target_locales: [fr, de, es, ja, yue, zh-CN, zh-TW, ar]

# One or more source-to-target mappings (required).
# {locale} is replaced with each configured target locale.
bundles:
  - id: app
    source: locales/en.json
    target: locales/{locale}.json
    format: json
  - id: docs
    source: README.md
    target: docs/i18n/{locale}.md
    format: markdown
  - id: browser
    source: browser/locales/en-US/browser.ftl
    target: browser/locales/{locale}/browser.ftl
    format: fluent

# Backward compatibility: source_path still maps targets to sibling files
# such as locales/fr.json. Prefer bundles for new projects.
# source_path: locales/en.json

# LLM provider settings
llm:
  # Provider: "anthropic", "openai", "gemini", or "openrouter" (default: gemini)
  provider: gemini

  # Model name defaults by provider:
  #   anthropic:  claude-opus-5
  #   openai:     gpt-5.6-luna (reasoning effort defaults to max)
  #   gemini:     gemini-3.8-flash
  #   openrouter: deepseek/deepseek-v4-pro-0813
  model: gemini-3.8-flash

  # Environment variable containing the API key
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Base URL for OpenAI-compatible endpoints (optional)
  # base_url: https://api.openai.com

  # OpenAI GPT-5-series Responses API reasoning effort
  # (default: max for the OpenAI provider)
  reasoning_effort: max

  # Optional LLM settings for individual target locales. An override using the
  # global provider inherits unspecified global settings. A different provider
  # uses that provider's defaults for unspecified settings.
  locale_overrides:
    yue:
      provider: openrouter
      model: deepseek/deepseek-v4-flash-0731
      api_key_env: OPENROUTER_API_KEY
    zh-CN:
      provider: openrouter
      model: deepseek/deepseek-v4-flash-0731
      api_key_env: OPENROUTER_API_KEY
    zh-TW:
      provider: openrouter
      model: deepseek/deepseek-v4-flash-0731
      api_key_env: OPENROUTER_API_KEY

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

# Versioned source, policy, target, and provenance state
# (default: .internationalizer.lock; commit this file)
manifest_path: .internationalizer.lock

# Optional translation and strict-validation rules
validation:
  plural_style: i18next-v4 # generate and validate target-locale plural forms
```

Locale identifiers must be well-formed BCP 47 tags such as `fr`, `pt-BR`, or
`sr-Latn-RS`. Canonical-equivalent target locales are rejected as duplicates,
and locale-specific provider overrides match canonical-equivalent spelling.
In the example above, locales without an override—including Japanese—inherit
the global Gemini configuration.

ICU MessageFormat values are parsed structurally. Simple arguments, `select`,
`plural`, `selectordinal`, `number`, `date`, and `time` are supported, including
nested messages, plural offsets, exact-number selectors, and `#`. Validation
checks syntax, argument types and formatter styles, plural offsets, select
branch identity, and target-locale CLDR plural categories. Provider output that
breaks these invariants is rejected before a locale file or translation-memory
record is written.

Fluent (`.ftl`) resources are handled as semantic source documents rather than
flattened maps. Message values, terms, and attributes become independent units;
comments are passed to the provider as developer context and included in source
provenance. Serialization preserves resource comments and ordering. Validation
protects variables, references, functions, selector defaults and branches, and
`data-l10n-name` markup slots while allowing target-locale selector variants and
natural reordering of named rich-text elements.

With `i18next-v4`, recognized source plural families are expanded during
translation to the target locale's CLDR categories. A target-only category uses
the source family's `_other` value as its translation template. Strict
validation requires those target categories; source-only categories are
optional for target locales that do not use them.

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
    "variants": ["Panneau de contrôle"],
    "enforcement": "error",
    "ignore_case": false,
    "whole_word": true
  }
]
```

`variants` lists other approved target forms. `enforcement` may be `error`,
`warning`, or omitted for the default error behavior. Terms are injected into
the LLM prompt as a terminology table, ensuring consistent translation across
your application. An exact entry such as `{"source":"API","target":"API"}`
also exempts that complete source-identical value from strict untranslated-value
findings; it does not exempt a longer value merely containing `API`.

## Translation Memory

Translation memory is stored as a JSONL file (one JSON record per line). Each record contains:

- The bundle, key, source value, translated value, and canonical target locale
- Source and translation-policy hashes
- The provider and model that produced the translation
- A timestamp

On subsequent runs, strings with the same source and policy hashes are served
from the cache without calling the LLM. The default path is under the ignored
`.internationalizer/` directory, so it remains a local cache. Set `tm_path` to a
tracked location if your project intentionally shares translation memory. The
reviewable `.internationalizer.lock` manifest is versioned separately.
Manifest schema v2 records provenance origin and review status independently,
so “generated successfully” and “approved by a person” cannot be conflated.

## Supported Formats

| Format | Extensions | Mode |
|--------|-----------|------|
| JSON | `.json` | Key-value (nested, dot-notation flattened) |
| YAML | `.yml`, `.yaml` | Key-value (preserves comments and ordering) |
| Markdown | `.md`, `.mdx` | Whole-document translation |
| Fluent | `.ftl` | Semantic messages, terms, attributes, comments, and selectors |

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
  fluentpattern/           Fluent pattern validation and safe text transforms
  formats/                 Format adapters (JSON, YAML, Markdown, Fluent)
  glossary/                Per-locale glossary management
  llm/                     LLM provider interface + implementations
    anthropic.go           Anthropic Claude backend
    openai.go              OpenAI / compatible backend
    gemini.go              Google Gemini via AI Studio backend
                           OpenRouter uses openai.go with custom base_url
  locale/                  BCP 47 identity and CLDR plural categories
  message/                 ICU MessageFormat parser and structural comparison
  policy/                  Stable translation-policy hashing
  pseudo/                  Provider-free accented and bidi test locales
  review/                  Explicit artifact approval workflow
  state/                   Versioned translation manifest
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

See [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) for dependency notices.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines. All contributions require DCO sign-off.
