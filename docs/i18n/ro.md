> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Flux de internaționalizare nativ AI pentru proiecte software. Traduceți, validați și gestionați fișierele i18n folosind LLM-uri.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br>
<a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## De ce Internationalizer?

Majoritatea instrumentelor i18n sunt fie biblioteci de execuție (i18next, react-intl), fie platforme SaaS pentru gestionarea cheilor (Crowdin, Lokalise). Niciuna dintre ele nu rezolvă bine problema reală a traducerii:

- **Traducerea manuală** nu este scalabilă dincolo de câteva limbi
- **API-urile de traducere automată** (Google Translate, DeepL) vă ignoră terminologia, tonul și convențiile interfeței cu utilizatorul
- **Traducerea LLM generică** funcționează mai bine, dar fără glosare și ghiduri de stil, obțineți rezultate inconsecvente

Internationalizer este diferit. Este un **flux CLI** care combină traducerea LLM cu:

- **Glosare specifice fiecărei limbi** — impun o terminologie consecventă în întreaga aplicație
- **Ghiduri de stil specifice fiecărei limbi** — controlează tonul, nivelul de formalitate, pluralizarea și tipografia
- **Memorie de traducere** — omite șirurile neschimbate, economisind bani la apelurile API
- **Validarea cheilor** — detectează traducerile lipsă și nepotrivirile de interpolare înainte de lansare

## Instalare

Instalați din npm:

```bash
npm install -g internationalizer
```

Sau rulați fără o instalare globală:

```bash
npx internationalizer --help
```

Pachetul npm instalează binarul precompilat corespunzător din npm prin intermediul dependențelor opționale specifice platformei.

Instalați cu Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Sau compilați din sursă:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Pachete npm

- Etichetele Git și versiunile pachetelor npm trebuie să se potrivească, de exemplu `v0.1.0` și `0.1.0`
- Pachetul rădăcină `internationalizer` depinde de pachete de platformă precum `internationalizer-darwin-arm64`
- Ținte npm acceptate: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publicarea CI necesită un secret GitHub numit `NPM_TOKEN`

## Pornire rapidă

1. Creați un fișier de configurare în rădăcina proiectului:

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

2. Setați cheia API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Previzualizați ce se va traduce:

```bash
internationalizer translate --dry-run
```

4. Rulați traducerea:

```bash
internationalizer translate
```

5. Validați toate limbile:

```bash
internationalizer validate
```

## Comenzi

### `translate`

Găsiți cheile lipsă și traduceți-le prin intermediul unui LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Verificați toate fișierele de limbă pentru chei lipsă, chei suplimentare și nepotriviri de interpolare.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Detectați automat framework-ul i18n și generați o sugestie de configurare.

```bash
internationalizer detect
```

Acceptă: react-i18next, next-intl, vue-i18n, JSON simplu, documente Markdown.

### `glossary`

Gestionați termenii de glosar specifici fiecărei limbi, care sunt impuși în timpul traducerii.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gestionați memoria de traducere (cache JSONL al șirurilor traduse anterior).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Referință de configurare

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

## Ghiduri de stil

Ghidurile de stil sunt fișiere Markdown care sunt injectate în promptul de traducere LLM. Acestea controlează tonul, nivelul de formalitate, tipografia și alte convenții specifice limbii.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Convenții partajate (`_conventions.md`)

Definiți regulile care se aplică tuturor limbilor: sintaxa de interpolare, păstrarea HTML, convențiile pentru tipurile de șiruri (butoane vs. etichete vs. erori) etc.

### Ghiduri specifice fiecărei limbi (`{locale}.md`)

Definiți regulile specifice limbii: registrul de formalitate (tu vs. vous), punctuația (ghilimele, semne de întrebare inversate), formele de plural, formatarea datelor/numerelor și un glosar terminologic.

Consultați [`examples/react-app/style-guides/`](examples/react-app/style-guides/) pentru un exemplu funcțional.

## Formatul glosarului

Fișierele de glosar sunt matrice JSON stocate în `{glossary_dir}/{locale}.json`:

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

Termenii sunt injectați în promptul LLM sub forma unui tabel terminologic, asigurând o traducere consecventă a termenilor cheie în întreaga aplicație.

## Memorie de traducere

Memoria de traducere este stocată ca un fișier JSONL (o înregistrare JSON pe linie). Fiecare înregistrare conține:

- Cheia și valoarea sursă
- Valoarea tradusă
- Un hash SHA-256 al valorii sursă
- Un marcaj de timp

La rulările ulterioare, șirurile neschimbate sunt servite din memoria cache TM fără a apela LLM-ul, economisind atât timp, cât și costuri API. Fișierul TM este compatibil cu Git și poate fi inclus în commit-uri alături de fișierele de limbă.

## Formate acceptate

| Format | Extensii | Mod |
|--------|-----------|------|
| JSON | `.json` | Cheie-valoare (imbricate, aplatizate prin notație cu punct) |
| YAML | `.yml`, `.yaml` | Cheie-valoare (păstrează comentariile și ordinea) |
| Markdown | `.md`, `.mdx` | Traducerea întregului document |

## Detectarea tipului de proiect

`internationalizer detect` identifică configurația i18n verificând:

- Dependențele din `package.json` pentru react-i18next, next-intl sau vue-i18n
- Structurile de directoare care se potrivesc cu modelele comune de limbă
- Extensiile de fișiere și convențiile de denumire

## Arhitectură

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

## Comparație cu alternativele

| Funcționalitate | Internationalizer | i18next | Crowdin | LLM generic |
|---------|------------------|---------|---------|-------------|
| Traducere bazată pe LLM | Da | Nu | Parțial | Da |
| Ghiduri de stil specifice limbii | Da | Nu | Nu | Nu |
| Impunerea glosarului | Da | Nu | Da | Nu |
| Memorie de traducere | Da | Nu | Da | Nu |
| CLI / execuție locală | Da | N/A | Nu | Manual |
| Fișiere compatibile cu Git | Da | Da | Parțial | Manual |
| Fără dependență de SaaS | Da | Da | Nu | Variază |
| Open source (AGPL-3.0) | Da | Da | Nu | Variază |

## Licență

[AGPL-3.0](LICENSE)

## Contribuții

Consultați [CONTRIBUTING.md](CONTRIBUTING.md) pentru configurarea mediului de dezvoltare și instrucțiuni. Toate contribuțiile necesită aprobarea DCO.

