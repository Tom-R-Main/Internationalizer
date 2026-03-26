> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI nativní internacionalizační pipeline pro softwarové projekty. Překládejte, ověřujte a spravujte soubory i18n pomocí LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br>
<a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Proč Internationalizer?

Většina nástrojů pro i18n jsou buď runtime knihovny (i18next, react-intl), nebo SaaS platformy pro správu klíčů (Crowdin, Lokalise). Žádný z nich však dobře neřeší samotný problém s překladem:

- **Manuální překlad** nelze efektivně škálovat pro více než několik jazyků.
- **API pro strojový překlad** (Google Translate, DeepL) ignorují vaši terminologii, tón a konvence uživatelského rozhraní.
- **Obecný překlad pomocí LLM** funguje lépe, ale bez glosářů a stylistických příruček získáte nekonzistentní výsledky.

Internationalizer je jiný. Je to **CLI pipeline**, která kombinuje překlad pomocí LLM s následujícími funkcemi:

- **Glosáře pro jednotlivé jazyky** — vynucují konzistentní terminologii v celé aplikaci.
- **Stylistické příručky pro jednotlivé jazyky** — řídí tón, formálnost, pluralizaci a typografii.
- **Překladová paměť** — přeskakuje nezměněné řetězce a šetří peníze za volání API.
- **Ověřování klíčů** — odhalí chybějící překlady a neshody v interpolaci ještě před vydáním.

## Instalace

Instalace přes npm:

```bash
npm install -g internationalizer
```

Nebo spuštění bez globální instalace:

```bash
npx internationalizer --help
```

Balíček npm nainstaluje odpovídající předkompilovanou binární verzi z npm prostřednictvím volitelných závislostí specifických pro danou platformu.

Instalace pomocí Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Nebo sestavení ze zdrojových kódů:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Balíčky npm

- Značky (tags) v Gitu a verze balíčků npm se musí shodovat, například `v0.1.0` a `0.1.0`.
- Kořenový balíček `internationalizer` závisí na balíčcích pro konkrétní platformy, jako je `internationalizer-darwin-arm64`.
- Podporované cíle npm: macOS arm64/x64, Linux arm64/x64, Windows x64.
- Publikování přes CI vyžaduje GitHub secret s názvem `NPM_TOKEN`.

## Rychlý start

1. Vytvořte konfigurační soubor v kořenovém adresáři projektu:

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

2. Nastavte svůj klíč API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Zobrazte si náhled toho, co se bude překládat:

```bash
internationalizer translate --dry-run
```

4. Spusťte překlad:

```bash
internationalizer translate
```

5. Ověřte všechny lokalizace:

```bash
internationalizer validate
```

## Příkazy

### `translate`

Najde chybějící klíče a přeloží je pomocí LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Zkontroluje všechny soubory lokalizace na chybějící klíče, přebytečné klíče a neshody v interpolaci.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Automaticky detekuje framework pro i18n a navrhne konfiguraci.

```bash
internationalizer detect
```

Podporuje: react-i18next, next-intl, vue-i18n, čistý JSON, dokumentaci v Markdownu.

### `glossary`

Spravuje termíny v glosáři pro jednotlivé jazyky, které jsou vynucovány během překladu.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Spravuje překladovou paměť (JSONL mezipaměť dříve přeložených řetězců).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Referenční příručka konfigurace

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

## Stylistické příručky

Stylistické příručky jsou soubory Markdown, které se vkládají do promptu pro překlad pomocí LLM. Řídí tón, formálnost, typografii a další konvence specifické pro daný jazyk.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Sdílené konvence (`_conventions.md`)

Definují pravidla, která platí pro všechny jazyky: syntaxi interpolace, zachování HTML, konvence pro typy řetězců (tlačítka vs. štítky vs. chyby) atd.

### Příručky pro jednotlivé jazyky (`{locale}.md`)

Definují pravidla specifická pro daný jazyk: úroveň formálnosti (tykání vs. vykání), interpunkci (uvozovky, obrácené otazníky), tvary množného čísla, formátování data/čísel a terminologický glosář.

Funkční příklad najdete v [`examples/react-app/style-guides/`](examples/react-app/style-guides/).

## Formát glosáře

Soubory glosáře jsou pole ve formátu JSON uložená v `{glossary_dir}/{locale}.json`:

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

Termíny jsou vloženy do promptu pro LLM jako terminologická tabulka, což zajišťuje konzistentní překlad klíčových termínů v celé aplikaci.

## Překladová paměť

Překladová paměť je uložena jako soubor JSONL (jeden záznam JSON na řádek). Každý záznam obsahuje:

- Zdrojový klíč a hodnotu.
- Přeloženou hodnotu.
- Hash SHA-256 zdrojové hodnoty.
- Časové razítko.

Při dalších spuštěních jsou nezměněné řetězce načteny z mezipaměti překladové paměti (TM) bez volání LLM, což šetří čas i náklady na API. Soubor TM je vhodný pro verzování v Gitu a lze jej commitovat společně se soubory lokalizace.

## Podporované formáty

| Formát | Přípony | Režim |
|--------|-----------|------|
| JSON | `.json` | Klíč-hodnota (vnořené, zploštělé pomocí tečkové notace) |
| YAML | `.yml`, `.yaml` | Klíč-hodnota (zachovává komentáře a pořadí) |
| Markdown | `.md`, `.mdx` | Překlad celého dokumentu |

## Detekce typu projektu

Příkaz `internationalizer detect` identifikuje vaše nastavení i18n kontrolou:

- Závislostí v `package.json` pro react-i18next, next-intl nebo vue-i18n.
- Struktur adresářů odpovídajících běžným vzorům lokalizace.
- Přípon souborů a konvencí pojmenování.

## Architektura

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

## Srovnání s alternativami

| Funkce | Internationalizer | i18next | Crowdin | Obecné LLM |
|---------|------------------|---------|---------|-------------|
| Překlad pomocí LLM | Ano | Ne | Částečně | Ano |
| Stylistické příručky pro jednotlivé jazyky | Ano | Ne | Ne | Ne |
| Vynucování glosáře | Ano | Ne | Ano | Ne |
| Překladová paměť | Ano | Ne | Ano | Ne |
| CLI / lokální spuštění | Ano | N/A | Ne | Manuálně |
| Soubory vhodné pro Git | Ano | Ano | Částečně | Manuálně |
| Bez závislosti na SaaS | Ano | Ano | Ne | Různé |
| Open source (AGPL-3.0) | Ano | Ano | Ne | Různé |

## Licence

[AGPL-3.0](LICENSE)

## Přispívání

Pokyny pro nastavení vývojového prostředí a pravidla najdete v [CONTRIBUTING.md](CONTRIBUTING.md). Všechny příspěvky vyžadují schválení DCO (Developer Certificate of Origin).

