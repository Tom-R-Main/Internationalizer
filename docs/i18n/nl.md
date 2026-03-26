> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-native internationalisatie-pipeline voor softwareprojecten. Vertaal, valideer en beheer i18n-bestanden met behulp van LLM's.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Waarom Internationalizer?

De meeste i18n-tools zijn ofwel runtime-bibliotheken (i18next, react-intl) of SaaS-platforms voor sleutelbeheer (Crowdin, Lokalise). Geen van deze lost het daadwerkelijke vertaalprobleem goed op:

- **Handmatige vertaling** is niet schaalbaar voorbij een paar talen
- **Machinevertaling-API's** (Google Translate, DeepL) negeren je terminologie, toon en UI-conventies
- **Generieke LLM-vertaling** werkt beter, maar zonder woordenlijsten en stijlgidsen krijg je inconsistente resultaten

Internationalizer is anders. Het is een **CLI-pipeline** die LLM-vertaling combineert met:

- **Woordenlijsten per taal** — dwing consistente terminologie af in je hele app
- **Stijlgidsen per taal** — beheer toon, formaliteit, meervoudsvormen en typografie
- **Vertaalgeheugen** — sla ongewijzigde strings over en bespaar geld op API-aanroepen
- **Sleutelvalidatie** — spoor ontbrekende vertalingen en interpolatiefouten op voordat ze live gaan

## Installatie

Installeer via npm:

```bash
npm install -g internationalizer
```

Of voer uit zonder globale installatie:

```bash
npx internationalizer --help
```

Het npm-pakket installeert de bijbehorende voorgebouwde binary via platformspecifieke optionele afhankelijkheden.

Installeer met Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Of bouw vanuit de broncode:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm-pakketten

- Git-tags en npm-pakketversies moeten overeenkomen, bijvoorbeeld `v0.1.0` en `0.1.0`
- Het hoofd-`internationalizer`-pakket is afhankelijk van platformpakketten zoals `internationalizer-darwin-arm64`
- Ondersteunde npm-targets: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI-publicatie vereist een GitHub-secret genaamd `NPM_TOKEN`

## Snel aan de slag

1. Maak een configuratiebestand aan in de hoofdmap van je project:

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

2. Stel je API-sleutel in:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Bekijk een voorbeeld van wat er vertaald gaat worden:

```bash
internationalizer translate --dry-run
```

4. Voer de vertaling uit:

```bash
internationalizer translate
```

5. Valideer alle talen:

```bash
internationalizer validate
```

## Commando's

### `translate`

Vind ontbrekende sleutels en vertaal ze via een LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Controleer alle taalbestanden op ontbrekende sleutels, extra sleutels en interpolatiefouten.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Detecteer automatisch het i18n-framework en stel een configuratie voor.

```bash
internationalizer detect
```

Ondersteunt: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown-documenten.

### `glossary`

Beheer woordenlijsttermen per taal die worden afgedwongen tijdens de vertaling.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Beheer het vertaalgeheugen (JSONL-cache van eerder vertaalde strings).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Configuratiereferentie

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

## Stijlgidsen

Stijlgidsen zijn Markdown-bestanden die worden geïnjecteerd in de LLM-vertaalprompt. Ze beheren de toon, formaliteit, typografie en andere taalspecifieke conventies.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Gedeelde conventies (`_conventions.md`)

Definieer regels die van toepassing zijn op alle talen: interpolatiesyntaxis, behoud van HTML, conventies voor stringtypen (knoppen vs. labels vs. foutmeldingen), enz.

### Gidsen per taal (`{locale}.md`)

Definieer taalspecifieke regels: formaliteitsregister (tu vs. vous), interpunctie (guillemets, omgekeerde vraagtekens), meervoudsvormen, datum-/getalnotatie en een terminologiewoordenlijst.

Zie [`examples/react-app/style-guides/`](examples/react-app/style-guides/) voor een werkend voorbeeld.

## Woordenlijstformaat

Woordenlijstbestanden zijn JSON-arrays die worden opgeslagen in `{glossary_dir}/{locale}.json`:

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

Termen worden in de LLM-prompt geïnjecteerd als een terminologietabel, wat zorgt voor een consistente vertaling van belangrijke termen in je hele applicatie.

## Vertaalgeheugen

Het vertaalgeheugen wordt opgeslagen als een JSONL-bestand (één JSON-record per regel). Elk record bevat:

- De bronsleutel en -waarde
- De vertaalde waarde
- Een SHA-256-hash van de bronwaarde
- Een tijdstempel

Bij volgende uitvoeringen worden ongewijzigde strings uit de TM-cache gehaald zonder de LLM aan te roepen, wat zowel tijd als API-kosten bespaart. Het TM-bestand is git-vriendelijk en kan samen met je taalbestanden worden gecommit.

## Ondersteunde formaten

| Formaat | Extensies | Modus |
|--------|-----------|------|
| JSON | `.json` | Key-value (genest, afgevlakt met puntnotatie) |
| YAML | `.yml`, `.yaml` | Key-value (behoudt opmerkingen en volgorde) |
| Markdown | `.md`, `.mdx` | Vertaling van het hele document |

## Detectie van projecttype

`internationalizer detect` identificeert je i18n-setup door te controleren op:

- `package.json`-afhankelijkheden voor react-i18next, next-intl of vue-i18n
- Mappenstructuren die overeenkomen met veelvoorkomende taalpatronen
- Bestandsextensies en naamgevingsconventies

## Architectuur

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

## Vergelijking met alternatieven

| Functie | Internationalizer | i18next | Crowdin | Generieke LLM |
|---------|------------------|---------|---------|-------------|
| LLM-aangedreven vertaling | Ja | Nee | Gedeeltelijk | Ja |
| Stijlgidsen per taal | Ja | Nee | Nee | Nee |
| Handhaving van woordenlijst | Ja | Nee | Ja | Nee |
| Vertaalgeheugen | Ja | Nee | Ja | Nee |
| CLI / lokale uitvoering | Ja | N.v.t. | Nee | Handmatig |
| Git-vriendelijke bestanden | Ja | Ja | Gedeeltelijk | Handmatig |
| Geen SaaS-afhankelijkheid | Ja | Ja | Nee | Varieert |
| Open source (AGPL-3.0) | Ja | Ja | Nee | Varieert |

## Licentie

[AGPL-3.0](LICENSE)

## Bijdragen

Zie [CONTRIBUTING.md](CONTRIBUTING.md) voor de ontwikkelingssetup en richtlijnen. Alle bijdragen vereisen een DCO-aftekening.

