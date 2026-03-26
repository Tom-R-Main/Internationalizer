> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline di internazionalizzazione nativa per l'AI per progetti software. Traduci, convalida e gestisci i file i18n usando gli LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br>
<a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Perché Internationalizer?

La maggior parte degli strumenti i18n sono librerie di runtime (i18next, react-intl) o piattaforme SaaS per la gestione delle chiavi (Crowdin, Lokalise). Nessuno di questi risolve bene il vero problema della traduzione:

- **La traduzione manuale** non scala oltre poche lingue
- **Le API di traduzione automatica** (Google Translate, DeepL) ignorano la tua terminologia, il tono e le convenzioni della UI
- **La traduzione LLM generica** funziona meglio, ma senza glossari e guide di stile ottieni risultati incoerenti

Internationalizer è diverso. È una **pipeline CLI** che combina la traduzione LLM con:

- **Glossari per lingua** — applicano una terminologia coerente in tutta l'app
- **Guide di stile per lingua** — controllano tono, formalità, pluralizzazione e tipografia
- **Memoria di traduzione** — salta le stringhe non modificate, risparmiando sulle chiamate API
- **Convalida delle chiavi** — individua traduzioni mancanti e mancate corrispondenze di interpolazione prima del rilascio

## Installazione

Installa da npm:

```bash
npm install -g internationalizer
```

Oppure esegui senza un'installazione globale:

```bash
npx internationalizer --help
```

Il pacchetto npm installa il binario precompilato corrispondente da npm tramite dipendenze opzionali specifiche per la piattaforma.

Installa con Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Oppure compila dal codice sorgente:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Pacchetti npm

- I tag Git e le versioni dei pacchetti npm devono corrispondere, ad esempio `v0.1.0` e `0.1.0`
- Il pacchetto radice `internationalizer` dipende da pacchetti di piattaforma come `internationalizer-darwin-arm64`
- Target npm supportati: macOS arm64/x64, Linux arm64/x64, Windows x64
- La pubblicazione CI richiede un secret di GitHub chiamato `NPM_TOKEN`

## Avvio rapido

1. Crea un file di configurazione nella radice del tuo progetto:

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

2. Imposta la tua chiave API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Visualizza in anteprima cosa verrà tradotto:

```bash
internationalizer translate --dry-run
```

4. Esegui la traduzione:

```bash
internationalizer translate
```

5. Convalida tutte le lingue:

```bash
internationalizer validate
```

## Comandi

### `translate`

Trova le chiavi mancanti e traducile tramite un LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Controlla tutti i file delle lingue per chiavi mancanti, chiavi extra e mancate corrispondenze di interpolazione.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Rileva automaticamente il framework i18n e suggerisce una configurazione.

```bash
internationalizer detect
```

Supporta: react-i18next, next-intl, vue-i18n, JSON vanilla, documenti markdown.

### `glossary`

Gestisci i termini del glossario per lingua che vengono applicati durante la traduzione.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gestisci la memoria di traduzione (cache JSONL delle stringhe tradotte in precedenza).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Riferimento per la configurazione

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

## Guide di stile

Le guide di stile sono file Markdown che vengono iniettati nel prompt di traduzione dell'LLM. Controllano il tono, la formalità, la tipografia e altre convenzioni specifiche della lingua.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Convenzioni condivise (`_conventions.md`)

Definisci le regole che si applicano a tutte le lingue: sintassi di interpolazione, conservazione dell'HTML, convenzioni sui tipi di stringa (pulsanti vs. etichette vs. errori), ecc.

### Guide per lingua (`{locale}.md`)

Definisci le regole specifiche della lingua: registro di formalità (tu vs. vous), punteggiatura (virgolette, punti interrogativi invertiti), forme plurali, formattazione di date/numeri e un glossario terminologico.

Vedi [`examples/react-app/style-guides/`](examples/react-app/style-guides/) per un esempio funzionante.

## Formato del glossario

I file del glossario sono array JSON memorizzati in `{glossary_dir}/{locale}.json`:

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

I termini vengono iniettati nel prompt dell'LLM come tabella terminologica, garantendo una traduzione coerente dei termini chiave in tutta l'applicazione.

## Memoria di traduzione

La memoria di traduzione è memorizzata come file JSONL (un record JSON per riga). Ogni record contiene:

- La chiave e il valore di origine
- Il valore tradotto
- Un hash SHA-256 del valore di origine
- Un timestamp

Nelle esecuzioni successive, le stringhe non modificate vengono fornite dalla cache della TM senza chiamare l'LLM, risparmiando sia tempo che costi delle API. Il file della TM è compatibile con git e può essere committato insieme ai file delle lingue.

## Formati supportati

| Formato | Estensioni | Modalità |
|--------|-----------|------|
| JSON | `.json` | Chiave-valore (nidificato, appiattito con notazione a punti) |
| YAML | `.yml`, `.yaml` | Chiave-valore (preserva commenti e ordinamento) |
| Markdown | `.md`, `.mdx` | Traduzione dell'intero documento |

## Rilevamento del tipo di progetto

`internationalizer detect` identifica la tua configurazione i18n controllando:

- Le dipendenze in `package.json` per react-i18next, next-intl o vue-i18n
- Le strutture delle directory che corrispondono a pattern comuni per le lingue
- Le estensioni dei file e le convenzioni di denominazione

## Architettura

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

## Confronto con le alternative

| Funzionalità | Internationalizer | i18next | Crowdin | LLM generico |
|---------|------------------|---------|---------|-------------|
| Traduzione basata su LLM | Sì | No | Parziale | Sì |
| Guide di stile per lingua | Sì | No | No | No |
| Applicazione del glossario | Sì | No | Sì | No |
| Memoria di traduzione | Sì | No | Sì | No |
| CLI / esecuzione locale | Sì | N/D | No | Manuale |
| File compatibili con Git | Sì | Sì | Parziale | Manuale |
| Nessuna dipendenza SaaS | Sì | Sì | No | Varia |
| Open source (AGPL-3.0) | Sì | Sì | No | Varia |

## Licenza

[AGPL-3.0](LICENSE)

## Contribuire

Vedi [CONTRIBUTING.md](CONTRIBUTING.md) per la configurazione dello sviluppo e le linee guida. Tutti i contributi richiedono l'approvazione DCO.

