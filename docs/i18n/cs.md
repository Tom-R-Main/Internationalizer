> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI nativní internacionalizační pipeline pro softwarové projekty. Překládejte, ověřujte a spravujte soubory i18n pomocí LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Proč Internationalizer?

Většina nástrojů pro i18n představuje buď runtime knihovny (i18next, react-intl), nebo SaaS platformy pro správu klíčů (Crowdin, Lokalise). Žádný z nich však skutečný problém s překladem neřeší uspokojivě:

- **Ruční překlad** nelze škálovat pro více než několik jazyků
- **Rozhraní API pro strojový překlad** (Google Translate, DeepL) ignorují vaši terminologii, tón i konvence uživatelského rozhraní
- **Běžný překlad pomocí LLM** funguje lépe, ale bez glosářů a stylistických příruček vede k nekonzistentním výsledkům

Internationalizer se liší. Představuje **CLI pipeline**, která spojuje překlad pomocí LLM s těmito prvky:

- **Glosáře pro jednotlivé jazyky** – vynucují jednotnou terminologii v celé aplikaci
- **Stylistické příručky pro jednotlivé jazyky** – řídí tón, formálnost, pravidla plurálu a typografii
- **Překladová paměť** – přeskakuje nezměněné řetězce a šetří náklady na volání API
- **Deterministické ověřování** – zachytí chybějící i přebývající klíče, porušení chráněných struktur, odchylky od glosáře a chyby v plurálech či syntaxi ICU ještě před nasazením
<!-- internationalizer:unit markdown:installation -->
## Instalace

Instalace přes npm:

```bash
npm install -g internationalizer
```

Nebo spuštění bez globální instalace:

```bash
npx internationalizer --help
```

Balíček npm nainstaluje odpovídající předkompilovanou binárku z npm prostřednictvím volitelných závislostí pro konkrétní platformu.

Instalace v Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Nebo sestavení ze zdrojových kódů:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## Balíčky npm

- Značky v Gitu a verze balíčků npm se musí přesně shodovat, například `v0.1.0` a `0.1.0`
- Kořenový balíček `internationalizer` závisí na platformních balíčcích, jako je `internationalizer-darwin-arm64`
- Podporované cíle npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publikování přes CI vyžaduje secret v GitHubu s názvem `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Rychlý start

1. V kořenovém adresáři projektu vytvořte konfigurační soubor:

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

2. Nastavte klíč rozhraní API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Zobrazte si náhled chystaného překladu:

```bash
internationalizer translate --dry-run
```

4. Spusťte překlad:

```bash
internationalizer translate
```

5. Zkontrolujte všechny lokalizace:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## Příkazy

### `translate`

Vyhledá chybějící nebo zastaralé klíče a přeloží je prostřednictvím LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Stav překladu nezávisle zaznamenává chybějící řetězce, zastaralý zdroj, zastaralá pravidla, aktuální stav a ručně upravené záznamy, takže ruční úprava nemůže zakrýt změnu zdroje ani pravidel. Hodnoty se zastaralými pravidly se hlásí, ale znovu překládají pouze s přepínačem `--refresh-policy`. Ručně upravené hodnoty se nikdy automaticky nepřepíší. Přepínač `--adopt-existing` použijte při zavádění manifestu pro již zkontrolované překlady nebo při výslovném přijetí zkontrolované ruční úpravy jako nového výchozího stavu.

### `validate`

Zkontroluje všechny soubory lokalizace oproti zdrojovým balíčkům. Výchozí ověření kontroluje strukturální pokrytí (procento přítomných požadovaných cílových klíčů), hlásí přebývající klíče jako varování a selže při chybějících klíčích, neshodách v interpolaci nebo neplatné struktuře ICU MessageFormat.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

Přepínač `--strict` navíc sleduje pokrytí překladu. Jazyková hodnota shodná se zdrojem se považuje za nepřeloženou, pokud glosář výslovně neobsahuje záznam se shodným zdrojem i cílem pro celou tuto hodnotu; volba `ignore_case` se respektuje, ale výskyt termínu z glosáře v delším textu výjimku netvoří. Striktní režim selže při přebývajících klíčích, hodnotách identických se zdrojem, změněné struktuře interpolace, HTML, kódu nebo odkazů v Markdownu, při porušení glosáře a chybějících konfigurovaných formách plurálu.

Volba `--require-state` ověřuje každý cíl proti souboru `.internationalizer.lock`. Selže, pokud klíč není sledován nebo pokud je zastaralý záznam o zdroji, pravidlech překladu či hashi cíle. Lze ji kombinovat s volbou `--strict`.

Textové i formátované JSON výstupy používají stálé kódy nálezů:

| Kód | Význam |
| --- | --- |
| `missing_key` / `extra_key` | Sady zdrojových a cílových klíčů se liší |
| `blank_translation` | Neprázdný zdroj má ve striktním režimu prázdný cíl |
| `source_identical` | Jazyková hodnota zůstala ve striktním režimu nepřeložená |
| `protected_structure_mismatch` | Změnila se struktura interpolace, kódů, značek HTML nebo odkazů |
| `glossary_violation` | Nebyl nalezen žádný schválený cílový termín ani varianta |
| `plural_form_missing` | Chybí nastavený tvar množného čísla pro danou lokalizaci |
| `icu_message_syntax` | Zdrojová nebo cílová zpráva ICU má neplatnou syntaxi |
| `icu_argument_mismatch` | Liší se názvy, typy nebo styly formátovačů argumentů ICU |
| `icu_selector_mismatch` | Selektory se liší nebo je kategorie plurálu neplatná pro cílovou lokalizaci |
| `untracked` | Pro cíl neexistuje v manifestu žádný záznam |
| `source_stale` | Zdrojový obsah se od zaznamenaného překladu změnil |
| `policy_stale` | Změnil se vygenerovaný prompt nebo nastavení modelu |
| `target_modified` | Cílový obsah se neshoduje se záznamem v manifestu |

### `detect`

Automaticky rozpozná framework pro i18n a navrhne konfiguraci.

```bash
internationalizer detect
```

Podporuje: react-i18next, next-intl, vue-i18n, čistý JSON, dokumentaci v Markdownu.

### `glossary`

Spravuje termíny glosáře pro jednotlivé jazyky, které jsou vynucovány během překladu.

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
<!-- internationalizer:unit markdown:configuration-reference -->
## Referenční příručka konfigurace

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

Identifikátory lokalizací musí být správně utvořené značky BCP 47, jako například `fr`, `pt-BR` nebo `sr-Latn-RS`. Kanonicky ekvivalentní cílové lokalizace jsou odmítnuty jako duplicity a specifická přepsání poskytovatelů se porovnávají podle kanonicky ekvivalentního zápisu. Ve výše uvedeném příkladu lokalizace bez přepsání – včetně japonštiny – dědí globální konfiguraci Gemini.

Hodnoty ICU MessageFormat se analyzují strukturálně. Podporovány jsou jednoduché argumenty, `select`, `plural`, `selectordinal`, `number`, `date` a `time`, včetně vnořených zpráv, posunů plurálu (offsets), přesných číselných selektorů i znaku `#`. Validace kontroluje syntaxi, typy argumentů a styly formátovačů, posuny plurálu, identitu větví select a kategorie plurálu podle CLDR pro cílovou lokalizaci. Výstupy poskytovatelů porušující tyto invarianty jsou odmítnuty ještě před zápisem do souboru lokalizace nebo záznamu překladové paměti.

Při použití `i18next-v4` se rozpoznané zdrojové rodiny plurálu při překladu rozšiřují na kategorie CLDR cílové lokalizace. Kategorie existující pouze v cíli používá jako šablonu překladu zdrojovou hodnotu `_other`. Striktní validace tyto cílové kategorie vyžaduje; kategorie přítomné pouze ve zdroji jsou pro cílové lokalizace, které je nepoužívají, volitelné.
<!-- internationalizer:unit markdown:style-guides -->
## Stylistické příručky

Stylistické příručky jsou soubory v Markdownu, které se vkládají do promptu pro překlad pomocí LLM. Určují tón, formálnost, typografii a další konvence specifické pro daný jazyk.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Sdílené konvence (`_conventions.md`)

Definují pravidla platná pro všechny jazyky: syntaxi interpolace, zachování HTML, konvence pro typy řetězců (tlačítka vs. popisky vs. chyby) atd.

### Příručky pro jednotlivé jazyky (`{locale}.md`)

Definují pravidla specifická pro daný jazyk: rovinu formálnosti (tykání vs. vykání), interpunkci (české uvozovky, obrácené otazníky), tvary množného čísla, formátování data a čísel a terminologický glosář.

Stylistické příručky představují trvalé konfigurační vstupy, nikoli generovaný výstup. Internationalizer z nich čte, ale nikdy je nepřepisuje. Jejich obsah se hashuje odděleně od glosáře a kontraktu promptu, takže změna kódu aplikace nezpůsobí zastarání překladu. Úprava příručky záměrně označí danou lokalizaci k revizi pravidel; změna interní formulace promptu k tomu nevede, pokud se zároveň nezmění verze kontraktu promptu.

Funkční příklad najdete v adresáři [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/).
<!-- internationalizer:unit markdown:glossary-format -->
## Formát glosáře

Soubory glosáře jsou pole JSON uložená v souboru `{glossary_dir}/{locale}.json`:

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

Položka `variants` uvádí další schválené tvary v cílovém jazyce. `enforcement` může nabývat hodnot `error`, `warning`, nebo může být vynechána pro výchozí chování error. Termíny jsou vloženy do systémového promptu LLM jako tabulka terminologie, což zajišťuje konzistentní překlad v celé aplikaci. Přesný záznam, jako například `{"source":"API","target":"API"}`, navíc osvobozuje celou tuto se zdrojem shodnou hodnotu od striktních nálezů o nepřeloženém obsahu; neosvobozuje však delší hodnotu, která řetězec `API` pouze obsahuje.
<!-- internationalizer:unit markdown:translation-memory -->
## Překladová paměť

Překladová paměť je uložena v souboru JSONL (jeden záznam JSON na řádek). Každý záznam obsahuje:

- Balíček, klíč, zdrojovou hodnotu, přeloženou hodnotu a kanonickou cílovou lokalizaci
- Hashe zdroje, stylistické příručky, glosáře, kontraktu promptu a kombinovaných pravidel
- Poskytovatele a model, které překlad vygenerovaly
- Časové razítko

Při dalších spuštěních se řetězce se stejným hashem zdroje i pravidel obsluhují z mezipaměti bez volání LLM. Výchozí cesta směřuje do ignorovaného adresáře `.internationalizer/`, takže zůstává lokální mezipamětí. Pokud projekt překladovou paměť záměrně sdílí, nastavte `tm_path` na sledované umístění. Revidovatelný manifest `.internationalizer.lock` se verzuje samostatně.
<!-- internationalizer:unit markdown:supported-formats -->
## Podporované formáty

| Formát | Přípony | Režim |
|--------|-----------|------|
| JSON | `.json` | Klíč-hodnota (vnořené, zploštělé tečkovou notací) |
| YAML | `.yml`, `.yaml` | Klíč-hodnota (zachovává komentáře a řazení) |
| Markdown | `.md`, `.mdx` | Úvodní část a oddíly na úrovni nadpisů H2 |

Cíle v Markdownu obsahují před oddíly H2 neviditelné komentáře `internationalizer:unit`. Tyto stabilní značky umožňují nástroji Internationalizer přidat, přesunout nebo upravit jeden zdrojový oddíl bez nutnosti znovu překládat nesouvisející oddíly. Stávající neoznačené dokumenty tyto značky obdrží při příští úspěšné aktualizaci.
<!-- internationalizer:unit markdown:project-type-detection -->
## Detekce typu projektu

Příkaz `internationalizer detect` rozpozná nastavení vaší internacionalizace kontrolou:

- Závislostí v `package.json` pro react-i18next, next-intl nebo vue-i18n
- Adresářových struktur odpovídajících běžným vzorům lokalizace
- Přípon souborů a konvencí pro pojmenování
<!-- internationalizer:unit markdown:architecture -->
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
  locale/                  BCP 47 identity and CLDR plural categories
  message/                 ICU MessageFormat parser and structural comparison
  policy/                  Stable translation-policy hashing
  state/                   Versioned translation manifest
  styleguide/              Style guide loader
  tm/                      JSONL translation memory
  translate/               Translation orchestrator
  validate/                Locale validation and diffing
```
<!-- internationalizer:unit markdown:comparison-to-alternatives -->
## Srovnání s alternativami

| Funkce | Internationalizer | i18next | Crowdin | Běžné LLM |
|---------|------------------|---------|---------|-------------|
| Překlad pomocí LLM | Ano | Ne | Částečně | Ano |
| Stylistické příručky pro jednotlivé jazyky | Ano | Ne | Ne | Ne |
| Vynucování glosáře | Ano | Ne | Ano | Ne |
| Překladová paměť | Ano | Ne | Ano | Ne |
| CLI / lokální spuštění | Ano | N/A | Ne | Ručně |
| Soubory vhodné pro Git | Ano | Ano | Částečně | Ručně |
| Žádná závislost na SaaS | Ano | Ano | Ne | Různé |
| Open source (AGPL-3.0) | Ano | Ano | Ne | Různé |
<!-- internationalizer:unit markdown:license -->
## Licence

[AGPL-3.0](../../LICENSE)

Oznámení o závislostech třetích stran najdete v souboru [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md).
<!-- internationalizer:unit markdown:contributing -->
## Přispívání

Pokyny pro nastavení vývojového prostředí a pravidla pro přispívání najdete v souboru [CONTRIBUTING.md](../../CONTRIBUTING.md). Všechny příspěvky vyžadují potvrzení DCO (Developer Certificate of Origin).
