> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-baseret internationaliserings-pipeline til softwareprojekter. Oversæt, valider og administrer i18n-filer ved hjælp af LLM'er.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Hvorfor Internationalizer?

De fleste i18n-værktøjer er enten runtime-biblioteker (i18next, react-intl) eller SaaS-platforme til nøgleadministration (Crowdin, Lokalise). Ingen af dem løser det egentlige oversættelsesproblem særlig godt:

- **Manuel oversættelse** skalerer ikke ud over et par sprog
- **Maskinoversættelses-API'er** (Google Translate, DeepL) ignorerer din terminologi, tone og UI-konventioner
- **Generisk LLM-oversættelse** fungerer bedre, men uden ordlister og stilguider får du inkonsistente resultater

Internationalizer er anderledes. Det er en **CLI-pipeline**, der kombinerer LLM-oversættelse med:

- **Sprogspecifikke ordlister** – håndhæv ensartet terminologi på tværs af din app
- **Sprogspecifikke stilguider** – styr tone, formalitet, flertalsbøjning og typografi
- **Oversættelseshukommelse** – spring uændrede strenge over, og spar penge på API-kald
- **Deterministisk validering** – fang manglende eller ekstra nøgler, ændringer i beskyttet struktur, ordlisteproblemer samt flertals- eller ICU-fejl, før de udgives
<!-- internationalizer:unit markdown:installation -->
## Installation

Installer fra npm:

```bash
npm install -g internationalizer
```

Eller kør uden en global installation:

```bash
npx internationalizer --help
```

npm-pakken installerer den matchende forudbyggede binære fil fra npm via platformsspecifikke valgfrie afhængigheder.

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
<!-- internationalizer:unit markdown:npm-packages -->
## npm-pakker

- Git-tags og npm-pakkeversioner skal stemme overens, for eksempel `v0.1.0` og `0.1.0`
- Rodpakken `internationalizer` afhænger af platformspakker såsom `internationalizer-darwin-arm64`
- Understøttede npm-mål: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI-udgivelse kræver en GitHub-hemmelighed ved navn `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Hurtig start

1. Opret en konfigurationsfil i din projektrod:

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
<!-- internationalizer:unit markdown:commands -->
## Kommandoer

### `translate`

Find manglende eller forældede nøgler, og oversæt dem via en LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Oversættelsesstatus rapporterer uafhængigt om manglende (missing), kildeforældede (source-stale), politikforældede (policy-stale), aktuelle (current) og manuelt redigerede tilstande, så en manuel redigering ikke kan skjule en kilde- eller politikændring. Politikforældede værdier rapporteres, men genoversættes kun med `--refresh-policy`. Manuelt redigerede værdier overskrives aldrig automatisk. Brug `--adopt-existing`, når du introducerer manifestet til gennemgåede oversættelser, eller når du udtrykkeligt accepterer en gennemgået manuel redigering som det nye udgangspunkt.

### `validate`

Kontroller alle sprogfiler mod deres kildebundter. Standardvalidering kontrollerer strukturel dækning (procentdelen af påkrævede målnøgler, der er til stede), rapporterer ekstra nøgler som advarsler og mislykkes ved manglende nøgler, uoverensstemmelser i interpolation eller ugyldig ICU MessageFormat-struktur.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` rapporterer også oversat dækning. En sproglig værdi, der er identisk med sin kilde, betragtes som uoversat, medmindre ordlisten udtrykkeligt indeholder en nøjagtig post med samme kilde og samme mål for den fulde værdi; der tages højde for `ignore_case`, men et ordlistebegreb, der er indlejret i en længere værdi, udgør ikke en undtagelse. Strikt tilstand mislykkes ved ekstra nøgler, kildeidentiske værdier, ændret struktur for interpolation, HTML, kode eller Markdown-links, overtrædelser af ordlisten samt konfigurerede flertalsformer.

`--require-state` verificerer hvert mål mod `.internationalizer.lock`. Den mislykkes, hvis en nøgle ikke spores, eller hvis dens registrerede kilde-, oversættelsespolitik- eller mål-hash er forældet. Den kan kombineres med `--strict`.

Menneske- og JSON-rapporter bruger stabile resultatkoder:

| Kode | Betydning |
| --- | --- |
| `missing_key` / `extra_key` | Kilde- og målnøglesæt er forskellige |
| `blank_translation` | En ikketom kilde har et tomt mål i strikt tilstand |
| `source_identical` | En sproglig værdi forbliver uoversat i strikt tilstand |
| `protected_structure_mismatch` | Struktur for interpolation, HTML, kode eller links er ændret |
| `glossary_violation` | Der blev ikke fundet noget godkendt målbegreb eller variant |
| `plural_form_missing` | En konfigureret sprogspecifik flertalsform mangler |
| `icu_message_syntax` | En kilde- eller mål-ICU-meddelelse er misdannet |
| `icu_argument_mismatch` | ICU-argumentnavne, -typer eller formateringsstile er forskellige |
| `icu_selector_mismatch` | Vælgere er forskellige, eller en flertalskategori er ugyldig for målsproget |
| `untracked` | Der findes ingen manifestpost for målet |
| `source_stale` | Kildeindholdet blev ændret efter den registrerede oversættelse |
| `policy_stale` | Den genererede prompt eller modelindstillingerne blev ændret |
| `target_modified` | Målindholdet afviger fra manifestposten |

### `detect`

Registrer automatisk i18n-frameworket, og foreslå en konfiguration.

```bash
internationalizer detect
```

Understøtter: react-i18next, next-intl, vue-i18n, ren JSON, Markdown-dokumenter.

### `glossary`

Administrer sprogspecifikke ordlistebegreber, der håndhæves under oversættelse.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Administrer oversættelseshukommelse (JSONL-cache med tidligere oversatte strenge).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Konfigurationsreference

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

Sprog-id'er skal være velformede BCP 47-tags såsom `fr`, `pt-BR` eller `sr-Latn-RS`. Kanonisk ækvivalente målsprog afvises som dubletter, og sprogspecifikke tilsidesættelser for udbydere matcher kanonisk ækvivalent stavning. I eksemplet ovenfor arver sprog uden en tilsidesættelse – herunder japansk – den globale Gemini-konfiguration.

ICU MessageFormat-værdier parses strukturelt. Simple argumenter, `select`, `plural`, `selectordinal`, `number`, `date` og `time` understøttes, herunder indlejrede meddelelser, flertalsforskydninger, eksakte talvælgere og `#`. Valideringen kontrollerer syntaks, argumenttyper og formateringsstile, flertalsforskydninger, select-forgreningsidentitet og målsprogets CLDR-flertalskategorier. Udbyderoutput, der bryder disse invarianter, afvises, før der skrives en sprogfil eller en post i oversættelseshukommelsen.

Med `i18next-v4` udvides genkendte kildeflertalsfamilier under oversættelsen til målsprogets CLDR-kategorier. En kategori, der kun findes i målet, bruger kildefamiliens `_other`-værdi som sin oversættelsesskabelon. Strikt validering kræver disse målkategorier; kategorier, der kun findes i kilden, er valgfrie for målsprog, der ikke anvender dem.
<!-- internationalizer:unit markdown:style-guides -->
## Stilguider

Stilguider er Markdown-filer, der injiceres i LLM-oversættelsesprompten. De styrer tone, formalitet, typografi og andre sprogspecifikke konventioner.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Fælles konventioner (`_conventions.md`)

Definer regler, der gælder for alle sprog: interpolationssyntaks, bevarelse af HTML, konventioner for strengtyper (knapper vs. etiketter vs. fejl) osv.

### Sprogspecifikke guider (`{locale}.md`)

Definer sprogspecifikke regler: formalitetsregister (du vs. De), tegnsætning (anførselstegn, omvendte spørgsmålstegn), flertalsformer, dato-/talformatering og en terminologiordliste.

Stilguider er varige politiske input, ikke genereret output. Internationalizer læser dem, men omskriver dem aldrig. Deres indhold hashes separat fra ordlisten og promptkontrakten, så en ændring i applikationskoden ikke gør en oversættelse forældet. Redigering af en guide markerer bevidst det pågældende sprog til politikgennemgang; ændring af intern promptordlyd gør ikke, medmindre versionen af promptkontrakten også ændres.

Se [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) for et fungerende eksempel.
<!-- internationalizer:unit markdown:glossary-format -->
## Ordlisteformat

Ordlistefiler er JSON-arrays, der gemmes i `{glossary_dir}/{locale}.json`:

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

`variants` angiver andre godkendte målformer. `enforcement` kan være `error`, `warning` eller udelades for standardadfærden med fejl. Begreber injiceres i LLM-prompten som en terminologitabel, hvilket sikrer ensartet oversættelse på tværs af din applikation. En nøjagtig post som f.eks. `{"source":"API","target":"API"}` fritager også den komplette kildeidentiske værdi for fund af uoversatte værdier i strikt tilstand; den fritager ikke en længere værdi, der blot indeholder `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## Oversættelseshukommelse

Oversættelseshukommelse gemmes som en JSONL-fil (én JSON-post pr. linje). Hver post indeholder:

- Bundtet, nøglen, kildeværdien, den oversatte værdi og det kanoniske målsprog
- Hashes for kilde, stilguide, ordliste, promptkontrakt og kombineret politik
- Den udbyder og model, der producerede oversættelsen
- Et tidsstempel

Ved efterfølgende kørsler leveres strenge med de samme kilde- og politik-hashes fra cachen uden at kalde LLM'en. Standardstien er under den ignorerede mappe `.internationalizer/`, så den forbliver en lokal cache. Angiv `tm_path` til en sporet placering, hvis dit projekt bevidst deler oversættelseshukommelse. Det gennemgåelige `.internationalizer.lock`-manifest versioneres separat.
<!-- internationalizer:unit markdown:supported-formats -->
## Understøttede formater

| Format | Udvidelser | Tilstand |
| --- | --- | --- |
| JSON | `.json` | Nøgle-værdi (indlejret, fladgjort med punktumnotation) |
| YAML | `.yml`, `.yaml` | Nøgle-værdi (bevarer kommentarer og rækkefølge) |
| Markdown | `.md`, `.mdx` | Præambel og afsnit på H2-niveau |

Markdown-mål indeholder usynlige `internationalizer:unit`-kommentarer før H2-afsnit. Disse stabile markører gør det muligt for Internationalizer at tilføje, flytte eller redigere et enkelt kildeafsnit uden at genoversætte urelaterede afsnit. Eksisterende uafmærkede dokumenter modtager markører ved deres næste vellykkede opdatering.
<!-- internationalizer:unit markdown:project-type-detection -->
## Registrering af projekttype

`internationalizer detect` identificerer din i18n-opsætning ved at kontrollere:

- `package.json`-afhængigheder for react-i18next, next-intl eller vue-i18n
- Mappestrukturer, der matcher almindelige sprogmønstre
- Filudvidelser og navngivningskonventioner
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
## Sammenligning med alternativer

| Funktion | Internationalizer | i18next | Crowdin | Generisk LLM |
| --- | --- | --- | --- | --- |
| LLM-drevet oversættelse | Ja | Nej | Delvist | Ja |
| Sprogspecifikke stilguider | Ja | Nej | Nej | Nej |
| Håndhævelse af ordliste | Ja | Nej | Ja | Nej |
| Oversættelseshukommelse | Ja | Nej | Ja | Nej |
| CLI / lokal udførelse | Ja | Ikke relevant | Nej | Manuelt |
| Git-venlige filer | Ja | Ja | Delvist | Manuelt |
| Ingen SaaS-afhængighed | Ja | Ja | Nej | Varierer |
| Open source (AGPL-3.0) | Ja | Ja | Nej | Varierer |
<!-- internationalizer:unit markdown:license -->
## Licens

[AGPL-3.0](../../LICENSE)

Se [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) for meddelelser om afhængigheder.
<!-- internationalizer:unit markdown:contributing -->
## Bidrag

Se [CONTRIBUTING.md](../../CONTRIBUTING.md) for opsætning af udviklingsmiljø og retningslinjer. Alle bidrag kræver DCO-godkendelse.
