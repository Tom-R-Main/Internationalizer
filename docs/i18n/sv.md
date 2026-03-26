> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-baserad pipeline för internationalisering av programvaruprojekt. Översätt, validera och hantera i18n-filer med LLM:er.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Varför Internationalizer?

De flesta i18n-verktyg är antingen runtime-bibliotek (i18next, react-intl) eller SaaS-plattformar för nyckelhantering (Crowdin, Lokalise). Inget av dem löser det faktiska översättningsproblemet på ett bra sätt:

- **Manuell översättning** skalar inte förbi ett fåtal språk
- **API:er för maskinöversättning** (Google Translate, DeepL) ignorerar din terminologi, ton och dina gränssnittskonventioner
- **Generisk LLM-översättning** fungerar bättre, men utan ordlistor och stilguider får du inkonsekventa resultat

Internationalizer är annorlunda. Det är en **CLI-pipeline** som kombinerar LLM-översättning med:

- **Språkspecifika ordlistor** — säkerställ konsekvent terminologi i hela din app
- **Språkspecifika stilguider** — kontrollera ton, formalitet, pluralisering och typografi
- **Översättningsminne** — hoppa över oförändrade strängar och spara pengar på API-anrop
- **Nyckelvalidering** — fånga upp saknade översättningar och interpoleringsfel innan de släpps

## Installation

Installera från npm:

```bash
npm install -g internationalizer
```

Eller kör utan global installation:

```bash
npx internationalizer --help
```

npm-paketet installerar den matchande förbyggda binären från npm via plattformsspecifika valfria beroenden.

Installera med Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Eller bygg från källkod:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm-paket

- Git-taggar och npm-paketversioner måste matcha, till exempel `v0.1.0` och `0.1.0`
- Rotpaketet `internationalizer` är beroende av plattformspaket som `internationalizer-darwin-arm64`
- npm-mål som stöds: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI-publicering kräver en GitHub-hemlighet med namnet `NPM_TOKEN`

## Snabbstart

1. Skapa en konfigurationsfil i din projektrot:

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

2. Ange din API-nyckel:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Förhandsgranska vad som kommer att översättas:

```bash
internationalizer translate --dry-run
```

4. Kör översättningen:

```bash
internationalizer translate
```

5. Validera alla språk:

```bash
internationalizer validate
```

## Kommandon

### `translate`

Hitta saknade nycklar och översätt dem via en LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Kontrollera alla språkfiler efter saknade nycklar, extra nycklar och interpoleringsfel.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Identifiera i18n-ramverket automatiskt och föreslå en konfiguration.

```bash
internationalizer detect
```

Stöder: react-i18next, next-intl, vue-i18n, vanlig JSON, markdown-dokument.

### `glossary`

Hantera språkspecifika ordlistetermer som upprätthålls under översättningen.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Hantera översättningsminne (JSONL-cache med tidigare översatta strängar).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Konfigurationsreferens

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

Stilguider är Markdown-filer som injiceras i LLM-översättningsprompten. De styr ton, formalitet, typografi och andra språkspecifika konventioner.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Delade konventioner (`_conventions.md`)

Definiera regler som gäller för alla språk: interpoleringssyntax, bevarande av HTML, konventioner för strängtyper (knappar vs. etiketter vs. fel) etc.

### Språkspecifika guider (`{locale}.md`)

Definiera språkspecifika regler: formalitetsregister (tu vs. vous), interpunktion (citattecken, omvända frågetecken), pluralformer, datum-/nummerformatering och en terminologiordlista.

Se [`examples/react-app/style-guides/`](examples/react-app/style-guides/) för ett fungerande exempel.

## Ordlisteformat

Ordlistefiler är JSON-vektorer som lagras i `{glossary_dir}/{locale}.json`:

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

Termer injiceras i LLM-prompten som en terminologitabell, vilket säkerställer konsekvent översättning av nyckeltermer i hela din applikation.

## Översättningsminne

Översättningsminnet lagras som en JSONL-fil (en JSON-post per rad). Varje post innehåller:

- Källnyckeln och värdet
- Det översatta värdet
- En SHA-256-hash av källvärdet
- En tidsstämpel

Vid efterföljande körningar hämtas oförändrade strängar från TM-cachen utan att anropa LLM:en, vilket sparar både tid och API-kostnader. TM-filen är git-vänlig och kan checkas in tillsammans med dina språkfiler.

## Format som stöds

| Format | Filändelser | Läge |
|--------|-----------|------|
| JSON | `.json` | Nyckel-värde (nästlad, tillplattad med punktnotation) |
| YAML | `.yml`, `.yaml` | Nyckel-värde (bevarar kommentarer och ordning) |
| Markdown | `.md`, `.mdx` | Översättning av hela dokument |

## Identifiering av projekttyp

`internationalizer detect` identifierar din i18n-uppsättning genom att kontrollera:

- Beroenden i `package.json` för react-i18next, next-intl eller vue-i18n
- Katalogstrukturer som matchar vanliga språkmönster
- Filändelser och namngivningskonventioner

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

## Jämförelse med alternativ

| Funktion | Internationalizer | i18next | Crowdin | Generisk LLM |
|---------|------------------|---------|---------|-------------|
| LLM-driven översättning | Ja | Nej | Delvis | Ja |
| Språkspecifika stilguider | Ja | Nej | Nej | Nej |
| Upprätthållande av ordlista | Ja | Nej | Ja | Nej |
| Översättningsminne | Ja | Nej | Ja | Nej |
| CLI / lokal körning | Ja | Ej tillämpligt | Nej | Manuellt |
| Git-vänliga filer | Ja | Ja | Delvis | Manuellt |
| Inget SaaS-beroende | Ja | Ja | Nej | Varierar |
| Öppen källkod (AGPL-3.0) | Ja | Ja | Nej | Varierar |

## Licens

[AGPL-3.0](LICENSE)

## Bidra

Se [CONTRIBUTING.md](CONTRIBUTING.md) för utvecklingsmiljö och riktlinjer. Alla bidrag kräver DCO-godkännande.

