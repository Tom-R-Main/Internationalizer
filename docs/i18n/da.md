> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-baseret internationaliserings-pipeline til softwareprojekter. Oversæt, valider og administrer i18n-filer ved hjælp af LLM'er.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Hvorfor Internationalizer?

De fleste i18n-værktøjer er enten runtime-biblioteker (i18next, react-intl) eller SaaS-platforme til nøglehåndtering (Crowdin, Lokalise). Ingen af dem løser det egentlige oversættelsesproblem særlig godt:

- **Manuel oversættelse** skalerer ikke ud over et par sprog
- **Maskinoversættelses-API'er** (Google Translate, DeepL) ignorerer din terminologi, tone og UI-konventioner
- **Generisk LLM-oversættelse** fungerer bedre, men uden ordlister og stilguider får du inkonsistente resultater

Internationalizer er anderledes. Det er en **CLI-pipeline**, der kombinerer LLM-oversættelse med:

- **Sprogspecifikke ordlister** — håndhæver konsistent terminologi på tværs af din app
- **Sprogspecifikke stilguider** — styrer tone, formalitet, flertalsbøjning og typografi
- **Oversættelseshukommelse** — springer uændrede strenge over og sparer penge på API-kald
- **Nøglevalidering** — fanger manglende oversættelser og interpolationsfejl, før de udgives

## Installation

Installer fra npm:

```bash
npm install -g internationalizer
```

Eller kør uden en global installation:

```bash
npx internationalizer --help
```

npm-pakken installerer den matchende forudkompilerede binære fil fra npm via platformsspecifikke valgfrie afhængigheder.

Installer med Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Eller byg fra kildekoden:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm-pakker

- Git-tags og npm-pakkeversioner skal matche, for eksempel `v0.1.0` og `0.1.0`
- Rodpakken `internationalizer` afhænger af platformspakker som `internationalizer-darwin-arm64`
- Understøttede npm-mål: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI-udgivelse kræver en GitHub-hemmelighed med navnet `NPM_TOKEN`

## Kom hurtigt i gang

1. Opret en konfigurationsfil i din projektrod:

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

2. Angiv din API-nøgle:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Se en forhåndsvisning af, hvad der vil blive oversat:

```bash
internationalizer translate --dry-run
```

4. Kør oversættelsen:

```bash
internationalizer translate
```

5. Valider alle sprog:

```bash
internationalizer validate
```

## Kommandoer

### `translate`

Find manglende nøgler, og oversæt dem via en LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Tjek alle sprogfiler for manglende nøgler, ekstra nøgler og interpolationsfejl.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Registrer automatisk i18n-frameworket, og foreslå en konfiguration.

```bash
internationalizer detect
```

Understøtter: react-i18next, next-intl, vue-i18n, ren JSON, markdown-dokumenter.

### `glossary`

Administrer sprogspecifikke ordlistetermer, der håndhæves under oversættelsen.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Administrer oversættelseshukommelse (JSONL-cache af tidligere oversatte strenge).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Konfigurationsreference

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

## Stilguider

Stilguider er Markdown-filer, der injiceres i LLM-oversættelsesprompten. De styrer tone, formalitet, typografi og andre sprogspecifikke konventioner.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Delte konventioner (`_conventions.md`)

Definer regler, der gælder for alle sprog: interpolationssyntaks, bevarelse af HTML, konventioner for strengtyper (knapper vs. etiketter vs. fejl) osv.

### Sprogspecifikke guider (`{locale}.md`)

Definer sprogspecifikke regler: formalitetsregister (tu vs. vous), tegnsætning (vinkelanførselstegn, omvendte spørgsmålstegn), flertalsformer, dato-/talformatering og en terminologiordliste.

Se [`examples/react-app/style-guides/`](examples/react-app/style-guides/) for et fungerende eksempel.

## Ordlisteformat

Ordlistefiler er JSON-arrays, der gemmes i `{glossary_dir}/{locale}.json`:

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

Termer injiceres i LLM-prompten som en terminologitabel, hvilket sikrer konsistent oversættelse af nøgletermer på tværs af din applikation.

## Oversættelseshukommelse

Oversættelseshukommelse gemmes som en JSONL-fil (én JSON-post pr. linje). Hver post indeholder:

- Kildenøglen og -værdien
- Den oversatte værdi
- En SHA-256-hash af kildeværdien
- Et tidsstempel

Ved efterfølgende kørsler leveres uændrede strenge fra TM-cachen uden at kalde LLM'en, hvilket sparer både tid og API-omkostninger. TM-filen er Git-venlig og kan committes sammen med dine sprogfiler.

## Understøttede formater

| Format | Udvidelser | Tilstand |
|--------|-----------|------|
| JSON | `.json` | Nøgle-værdi (indlejret, fladgjort med punktumnotation) |
| YAML | `.yml`, `.yaml` | Nøgle-værdi (bevarer kommentarer og rækkefølge) |
| Markdown | `.md`, `.mdx` | Oversættelse af hele dokumentet |

## Registrering af projekttype

`internationalizer detect` identificerer din i18n-opsætning ved at tjekke:

- `package.json`-afhængigheder for react-i18next, next-intl eller vue-i18n
- Mappestrukturer, der matcher almindelige sprogmønstre
- Filudvidelser og navngivningskonventioner

## Arkitektur

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

## Sammenligning med alternativer

| Funktion | Internationalizer | i18next | Crowdin | Generisk LLM |
|---------|------------------|---------|---------|-------------|
| LLM-drevet oversættelse | Ja | Nej | Delvist | Ja |
| Sprogspecifikke stilguider | Ja | Nej | Nej | Nej |
| Håndhævelse af ordliste | Ja | Nej | Ja | Nej |
| Oversættelseshukommelse | Ja | Nej | Ja | Nej |
| CLI / lokal kørsel | Ja | N/A | Nej | Manuelt |
| Git-venlige filer | Ja | Ja | Delvist | Manuelt |
| Ingen SaaS-afhængighed | Ja | Ja | Nej | Varierer |
| Open source (AGPL-3.0) | Ja | Ja | Nej | Varierer |

## Licens

[AGPL-3.0](LICENSE)

## Bidrag

Se [CONTRIBUTING.md](CONTRIBUTING.md) for opsætning af udviklingsmiljø og retningslinjer. Alle bidrag kræver DCO-godkendelse.

