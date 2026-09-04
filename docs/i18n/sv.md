> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-baserad pipeline för internationalisering av programvaruprojekt. Översätt, validera och hantera i18n-filer med LLM:er.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Varför Internationalizer?

De flesta i18n-verktyg är antingen runtime-bibliotek (i18next, react-intl) eller SaaS-plattformar för nyckelhantering (Crowdin, Lokalise). Inget av dem löser det faktiska översättningsproblemet på ett bra sätt:

- **Manuell översättning** skalar inte förbi ett fåtal språk
- **API:er för maskinöversättning** (Google Translate, DeepL) ignorerar din terminologi, ton och dina gränssnittskonventioner
- **Generisk LLM-översättning** fungerar bättre, men utan ordlistor och stilguider får du inkonsekventa resultat

Internationalizer skiljer sig från mängden. Det är en **CLI-pipeline** som kombinerar LLM-översättning med:

- **Språkspecifika ordlistor** – säkerställ konsekvent terminologi i hela programmet
- **Språkspecifika stilguider** – styr ton, formalitet, pluralisering och typografi
- **Översättningsminne** – hoppa över oförändrade strängar och spara pengar på API-anrop
- **Deterministisk validering** – fånga upp saknade eller extra nycklar, avvikelser i skyddad struktur, ordlisteproblem samt plural- eller ICU-fel innan de når produktion
<!-- internationalizer:unit markdown:installation -->
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
<!-- internationalizer:unit markdown:npm-packages -->
## npm-paket

- Git-taggar och npm-paketversioner måste matcha, till exempel `v0.1.0` och `0.1.0`
- Rotpaketet `internationalizer` är beroende av plattformspaket som `internationalizer-darwin-arm64`
- npm-mål som stöds: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publicering via CI kräver en GitHub-hemlighet med namnet `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Snabbstart

1. Skapa en konfigurationsfil i projektroten:

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

5. Validera alla språkmål:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## Kommandon

### `translate`

Hitta saknade eller inaktuella nycklar och översätt dem via en LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Översättningstillståndet rapporterar oberoende av varandra tillstånden saknad (missing), inaktuell källa (source-stale), inaktuell policy (policy-stale), aktuell (current) och manuellt redigerad (manually edited), så att en manuell redigering inte kan dölja en ändring i källan eller policyn. Värden med inaktuell policy rapporteras men översätts endast om med `--refresh-policy`. Manuellt redigerade värden skrivs aldrig över automatiskt. Använd `--adopt-existing` när manifestet introduceras för granskade översättningar eller när du uttryckligen vill acceptera en granskad manuell redigering som ny baslinje.

### `validate`

Kontrollera alla språkfiler mot källknippena. Standardvalideringen kontrollerar strukturell täckning (procentandelen obligatoriska målnycklar som finns), rapporterar extra nycklar som varningar och misslyckas vid saknade nycklar, bristande överensstämmelse i interpolering eller ogiltig struktur enligt ICU MessageFormat.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` rapporterar även översättningstäckning. Ett språkligt värde som är identiskt med källan räknas som oöversatt, såvida inte ordlistan uttryckligen innehåller en post med exakt samma källa och mål för hela värdet. Flaggan `ignore_case` respekteras, men en ordlisteterm som ingår i ett längre värde utgör inget undantag. Strikt läge ger fel vid extra nycklar, källidentiska värden, ändrad struktur för interpolering/HTML/kod/Markdown-länkar, överträdelser mot ordlistan och konfigurerade pluralformer.

`--require-state` verifierar varje målfil mot `.internationalizer.lock`. Kommandot misslyckas om en nyckel inte spåras, eller om dess registrerade källa, översättningspolicy eller målhash är inaktuell. Det kan kombineras med `--strict`.

Rapporter i text- och JSON-format använder stabila identifieringskoder:

| Kod | Betydelse |
| --- | --- |
| `missing_key` / `extra_key` | Nyckeluppsättningarna i källan och målet skiljer sig åt |
| `blank_translation` | En icke-tom källa har ett tomt målvärde i strikt läge |
| `source_identical` | Ett språkligt värde i strikt läge är fortfarande oöversatt |
| `protected_structure_mismatch` | Strukturen för interpolering, HTML, kod eller länk har ändrats |
| `glossary_violation` | Ingen godkänd målterm eller variant hittades |
| `plural_form_missing` | En konfigurerad pluralform saknas för målspråket |
| `icu_message_syntax` | Ett ICU-meddelande i källan eller målet är felaktigt utformat |
| `icu_argument_mismatch` | ICU-argumentnamn, typer eller formateringsstilar skiljer sig åt |
| `icu_selector_mismatch` | Selektorer skiljer sig åt eller en pluralkategori är ogiltig för målspråket |
| `untracked` | Det finns ingen manifestpost för målet |
| `source_stale` | Källinnehållet ändrades efter den registrerade översättningen |
| `policy_stale` | Den genererade prompten eller modellinställningarna har ändrats |
| `target_modified` | Målinnehållet skiljer sig från manifestposten |

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

Hantera översättningsminnet (JSONL-cache med tidigare översatta strängar).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Konfigurationsreferens

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

Språkidentifierare måste vara välformade BCP 47-taggar som `fr`, `pt-BR` eller `sr-Latn-RS`. Kanoniskt likvärdiga målspråk avvisas som dubbletter, och språkspecifika åsidosättningar av leverantör matchas mot kanoniskt likvärdig stavning. I exemplet ovan ärver språkmål utan en åsidosättning – däribland japanska – den globala Gemini-konfigurationen.

Värden i ICU MessageFormat tolkas strukturellt. Enkla argument, `select`, `plural`, `selectordinal`, `number`, `date` och `time` stöds, inklusive nästlade meddelanden, pluralförskjutningar, selektorer med exakta tal samt `#`. Valideringen kontrollerar syntax, argumenttyper och formateringsstilar, pluralförskjutningar, identitet för select-grenar och CLDR-pluralkategorier för målspråket. Leverantörssvar som bryter mot dessa invarianter avvisas innan en språkfil eller en post i översättningsminnet skrivs.

Med `i18next-v4` expanderas identifierade pluralfamiljer i källan under översättningen till målspråkets CLDR-kategorier. En kategori som bara finns i målet använder källfamiljens `_other`-värde som översättningsmall. Strikt validering kräver dessa målkategorier; kategorier som bara finns i källan är valfria för målspråk som inte använder dem.
<!-- internationalizer:unit markdown:style-guides -->
## Stilguider

Stilguider är Markdown-filer som injiceras i prompten för LLM-översättning. De styr ton, formalitet, typografi och andra språkspecifika konventioner.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Delade konventioner (`_conventions.md`)

Definiera regler som gäller för alla språk: interpoleringssyntax, bevarande av HTML, konventioner för strängtyper (knappar jämfört med etiketter och fel) osv.

### Språkspecifika guider (`{locale}.md`)

Definiera språkspecifika regler: formalitetsregister (tu jämfört med vous), interpunktion (citattecken, omvända frågetecken), pluralformer, datum- och talformatering samt en terminologiordlista.

Stilguider är beständiga policyunderlag, inte genererade utdata. Internationalizer läser dem men skriver aldrig om dem. Deras innehåll hashas separat från ordlistan och promptkontraktet, så att en kodändring i programmet inte gör en översättning inaktuell. Om du redigerar en guide markeras det språket avsiktligt för policygranskning; ändringar i interna promptformuleringar gör det inte, såvida inte promptkontraktets version också ändras.

Se [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) för ett fungerande exempel.
<!-- internationalizer:unit markdown:glossary-format -->
## Ordlisteformat

Ordlistefiler är JSON-matriser som lagras i `{glossary_dir}/{locale}.json`:

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

`variants` listar andra godkända målformer. `enforcement` kan vara `error`, `warning` eller utelämnas för standardbeteendet fel. Termerna injiceras i LLM-prompten som en terminologitabell, vilket säkerställer konsekvent översättning i hela programmet. En exakt post som `{"source":"API","target":"API"}` undantar även det fullständiga källidentiska värdet från anmärkningar om oöversatta värden i strikt läge; den undantar inte ett längre värde som bara innehåller `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## Översättningsminne

Översättningsminnet lagras som en JSONL-fil (en JSON-post per rad). Varje post innehåller:

- Knippe, nyckel, källvärde, översatt värde och kanoniskt målspråk
- Hashar för källa, stilguide, ordlista, promptkontrakt och kombinerad policy
- Leverantören och modellen som genererade översättningen
- En tidsstämpel

Vid efterföljande körningar hämtas strängar med samma käll- och policyhashar från cachen utan att anropa LLM:en. Standardsökvägen ligger under den ignorerade katalogen `.internationalizer/`, så att den förblir en lokal cache. Sätt `tm_path` till en spårad plats om projektet avsiktligt delar översättningsminne. Det granskningsbara manifestet `.internationalizer.lock` versionshanteras separat.
<!-- internationalizer:unit markdown:supported-formats -->
## Format som stöds

| Format | Filnamnstillägg | Läge |
|--------|-----------|------|
| JSON | `.json` | Nyckel-värde (nästlad, tillplattad med punktnotation) |
| YAML | `.yml`, `.yaml` | Nyckel-värde (bevarar kommentarer och ordning) |
| Markdown | `.md`, `.mdx` | Ingress och avsnitt på H2-nivå |

Markdown-målfiler innehåller osynliga `internationalizer:unit`-kommentarer före H2-avsnitt. Dessa stabila markörer gör att Internationalizer kan lägga till, flytta eller redigera ett enskilt källavsnitt utan att översätta om orelaterade avsnitt. Befintliga omarkerade dokument får markörer vid nästa lyckade uppdatering.
<!-- internationalizer:unit markdown:project-type-detection -->
## Identifiering av projekttyp

`internationalizer detect` identifierar din i18n-konfiguration genom att kontrollera:

- Beroenden i `package.json` för react-i18next, next-intl eller vue-i18n
- Katalogstrukturer som matchar vanliga språkmönster
- Filnamnstillägg och namngivningskonventioner
<!-- internationalizer:unit markdown:architecture -->
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
## Jämförelse med alternativ

| Funktion | Internationalizer | i18next | Crowdin | Generisk LLM |
|---------|------------------|---------|---------|-------------|
| LLM-driven översättning | Ja | Nej | Delvis | Ja |
| Språkspecifika stilguider | Ja | Nej | Nej | Nej |
| Upprätthållande av ordlista | Ja | Nej | Ja | Nej |
| Översättningsminne | Ja | Nej | Ja | Nej |
| CLI / lokal körning | Ja | Ej tillämpligt | Nej | Manuell |
| Git-vänliga filer | Ja | Ja | Delvis | Manuell |
| Inget SaaS-beroende | Ja | Ja | Nej | Varierar |
| Öppen källkod (AGPL-3.0) | Ja | Ja | Nej | Varierar |
<!-- internationalizer:unit markdown:license -->
## Licens

[AGPL-3.0](../../LICENSE)

Se [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) för information om beroendens licenser.
<!-- internationalizer:unit markdown:contributing -->
## Bidra

Se [CONTRIBUTING.md](../../CONTRIBUTING.md) för utvecklingsmiljö och riktlinjer. Alla bidrag kräver DCO-godkännande.
