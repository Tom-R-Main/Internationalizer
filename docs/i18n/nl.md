> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-native internationalisatiepipeline voor softwareprojecten. Vertaal, valideer en beheer i18n-bestanden met behulp van LLM's.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Waarom Internationalizer?

De meeste i18n-tools zijn ofwel runtimebibliotheken (i18next, react-intl) of SaaS-platforms voor sleutelbeheer (Crowdin, Lokalise). Geen van beide lost het daadwerkelijke vertaalprobleem goed op:

- **Handmatige vertaling** schaalt niet voorbij een paar talen
- **Machinevertaling-API's** (Google Translate, DeepL) negeren je terminologie, toon en UI-conventies
- **Generieke LLM-vertaling** werkt beter, maar zonder woordenlijsten en stijlgidsen krijg je inconsistente resultaten

Internationalizer pakt het anders aan. Het is een **CLI-pipeline** die LLM-vertaling combineert met:

- **Woordenlijsten per taal** — dwing consistente terminologie af in je hele app
- **Stijlgidsen per taal** — stuur toon, formaliteit, meervoudsvormen en typografie aan
- **Vertaalgeheugen** — sla ongewijzigde strings over en bespaar kosten op API-aanroepen
- **Deterministische validatie** — spoor ontbrekende of overtollige sleutels, afwijkingen in beschermde structuren, terminologieproblemen en meervouds- of ICU-fouten op voordat ze in productie gaan

<!-- internationalizer:unit markdown:installation -->
## Installatie

Installeer via npm:

```bash
npm install -g internationalizer
```

Of voer uit zonder globale installatie:

```bash
npx internationalizer --help
```

Het npm-pakket installeert de bijbehorende voorgebouwde binary vanaf npm via platformspecifieke optionele afhankelijkheden.

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

<!-- internationalizer:unit markdown:npm-packages -->
## npm-pakketten

- Git-tags en npm-pakketversies moeten exact overeenkomen, bijvoorbeeld `v0.1.0` en `0.1.0`
- Het rootpakket `internationalizer` is afhankelijk van platformspecifieke pakketten zoals `internationalizer-darwin-arm64`
- Ondersteunde npm-platforms: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publicatie via CI vereist een GitHub-secret met de naam `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Snel aan de slag

1. Maak een configuratiebestand aan in de hoofdmap van je project:

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

2. Stel je API-sleutel in:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Bekijk vooraf wat er vertaald gaat worden:

```bash
internationalizer translate --dry-run
```

4. Voer de vertaling uit:

```bash
internationalizer translate
```

5. Valideer alle doeltalen:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## Commando's

### `translate`

Zoek ontbrekende of verouderde sleutels en vertaal ze via een LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

De vertaalstatus rapporteert onafhankelijk over ontbrekende, brongewijzigde (source-stale), beleidsgewijzigde (policy-stale), actuele en handmatig bewerkte situaties. Hierdoor kan een handmatige aanpassing nooit een bron- of beleidswijziging maskeren. Beleidsgewijzigde waarden worden wel gemeld, maar uitsluitend opnieuw vertaald met `--refresh-policy`. Handmatig bewerkte waarden worden nooit automatisch overschreven. Gebruik `--adopt-existing` wanneer je het manifest introduceert voor reeds gecontroleerde vertalingen of wanneer je een gecontroleerde handmatige bewerking expliciet als nieuwe uitgangssituatie accepteert.

### `validate`

Controleer alle doelbestanden tegen hun bronbundels. De standaardvalidatie verifieert de structurele dekking (het percentage aanwezige vereiste doelsleutels), meldt overtollige sleutels als waarschuwingen en mislukt bij ontbrekende sleutels, niet-overeenkomende interpolatievariabelen of een ongeldige ICU MessageFormat-structuur.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` rapporteert daarnaast de vertaalde dekking. Een taalkundige waarde die identiek is aan de bron geldt als onvertaald, tenzij de woordenlijst expliciet een exacte vermelding bevat met dezelfde bron en hetzelfde doel voor de volledige waarde; hierbij wordt rekening gehouden met `ignore_case`, maar een woordenlijstterm die deel uitmaakt van een langere string geldt niet als uitzondering. De strikte modus faalt op overtollige sleutels, met de bron identieke waarden, gewijzigde interpolatie-, HTML-, code- of Markdown-linkstructuren, woordenlijstovertredingen en geconfigureerde meervoudsvormen.

`--require-state` verifieert elk doelbestand tegen `.internationalizer.lock`. Dit commando faalt wanneer een sleutel niet wordt bijgehouden, of wanneer de vastgelegde bron, het vertaalbeleid of de doelhash verouderd is. Dit kan worden gecombineerd met `--strict`.

De rapportages voor mensen en JSON gebruiken vaste bevindingcodes:

| Code | Betekenis |
| --- | --- |
| `missing_key` / `extra_key` | De sleutelsets van bron en doel wijken af |
| `blank_translation` | Een niet-lege bronwaarde heeft een leeg doel in strikte modus |
| `source_identical` | Een taalkundige waarde in strikte modus is onvertaald gebleven |
| `protected_structure_mismatch` | De opzet van interpolaties, HTML, code of links is gewijzigd |
| `glossary_violation` | Er is geen goedgekeurde doelterm of goedgekeurde variant gevonden |
| `plural_form_missing` | Een geconfigureerde meervoudsvorm voor de doeltaal ontbreekt |
| `icu_message_syntax` | Een ICU-bericht in de bron of het doel bevat syntaxfouten |
| `icu_argument_mismatch` | ICU-argumentnamen, typen of formatteerstijlen komen niet overeen |
| `icu_selector_mismatch` | Selectors wijken af of een meervoudscategorie is ongeldig voor de doeltaal |
| `untracked` | Er bestaat geen manifestvermelding voor het doel |
| `source_stale` | De broninhoud is gewijzigd na de vastgelegde vertaling |
| `policy_stale` | De gegenereerde prompt of modelinstellingen zijn gewijzigd |
| `target_modified` | De inhoud van het doel wijkt af van de vermelding in het manifest |

### `detect`

Detecteer automatisch het i18n-framework en stel een configuratie voor.

```bash
internationalizer detect
```

Ondersteunt: react-i18next, next-intl, vue-i18n, standaard JSON, Markdown-documentatie.

### `glossary`

Beheer woordenlijsttermen per taal die tijdens de vertaling worden afgedwongen.

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

<!-- internationalizer:unit markdown:configuration-reference -->
## Configuratiereferentie

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

Taal-ID's moeten correct opgemaakte BCP 47-tags zijn zoals `fr`, `pt-BR` of `sr-Latn-RS`. Canoniek gelijkwaardige doeltalen worden geweigerd als duplicaten, en taalspecifieke provider-overrides matchen op canoniek gelijkwaardige schrijfwijzen. In het bovenstaande voorbeeld nemen doeltalen zonder override — inclusief Japans — de globale Gemini-configuratie over.

ICU MessageFormat-waarden worden structureel geanalyseerd. Eenvoudige argumenten, `select`, `plural`, `selectordinal`, `number`, `date` en `time` worden ondersteund, inclusief geneste berichten, meervoudsoffsets, selectors met exacte getallen en `#`. Validatie controleert syntax, argumenttypen, formatteerstijlen, meervoudsoffsets, de identiteit van select-takken en de CLDR-meervoudscategorieën van de doeltaal. Uitvoer van providers die deze invarianten schendt, wordt afgewezen voordat een doeltaalbestand of vertaalgeheugenrecord wordt weggeschreven.

Met `i18next-v4` worden herkende meervoudsfamilies in de bron tijdens de vertaling uitgebreid naar de CLDR-categorieën van de doeltaal. Een categorie die alleen in de doeltaal voorkomt, gebruikt de `_other`-waarde van de bronfamilie als vertaalsjabloon. Strikte validatie vereist deze doelcategorieën; categorieën die alleen in de bron voorkomen, zijn optioneel voor doeltalen die deze niet hanteren.

<!-- internationalizer:unit markdown:style-guides -->
## Stijlgidsen

Stijlgidsen zijn Markdown-bestanden die worden toegevoegd aan de LLM-vertaalprompt. Ze sturen de toon, formaliteit, typografie en andere taalspecifieke conventies aan.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Gedeelde conventies (`_conventions.md`)

Definieer regels die gelden voor alle talen: interpolatiesyntaxis, behoud van HTML, conventies voor stringtypen (knoppen vs. labels vs. foutmeldingen), enzovoort.

### Gidsen per taal (`{locale}.md`)

Definieer taalspecifieke regels: formaliteitsregister (je/jouw vs. u/uw), interpunctie (aanhalingstekens, gedachtestreepjes), meervoudsvormen, opmaak van datums en getallen, en een terminologiewoordenlijst.

Stijlgidsen zijn duurzame beleidsinvoer, geen gegenereerde uitvoer. Internationalizer leest ze uit, maar herschrijft ze nooit. De inhoud ervan wordt afzonderlijk gehasht van de woordenlijst en het promptcontract, waardoor een codewijziging in de applicatie een vertaling niet zomaar als verouderd markeert. Het bewerken van een gids markeert die taal bewust voor beleidsbeoordeling; interne aanpassingen aan promptformuleringen doen dat niet, tenzij ook de versie van het promptcontract wijzigt.

Zie [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) voor een werkend voorbeeld.

<!-- internationalizer:unit markdown:glossary-format -->
## Woordenlijstindeling

Woordenlijstbestanden zijn JSON-arrays die worden opgeslagen in `{glossary_dir}/{locale}.json`:

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

`variants` geeft andere goedgekeurde doelvarianten weer. `enforcement` kan worden ingesteld op `error`, `warning`, of worden weggelaten voor het standaardfoutgedrag. Termen worden als terminologietabel aan de LLM-prompt toegevoegd, wat zorgt voor een consistente vertaling binnen je hele applicatie. Een exacte vermelding zoals `{"source":"API","target":"API"}` stelt die volledige, met de bron identieke waarde bovendien vrij van meldingen over onvertaalde waarden in de strikte modus; een langere waarde die toevallig `API` bevat, wordt hiermee niet vrijgesteld.

<!-- internationalizer:unit markdown:translation-memory -->
## Vertaalgeheugen

Het vertaalgeheugen wordt opgeslagen als een JSONL-bestand (één JSON-record per regel). Elk record bevat:

- De bundel, sleutel, bronwaarde, vertaalde waarde en canonieke doeltaal
- Hashes van de bron, de stijlgids, de woordenlijst, het promptcontract en het gecombineerde beleid
- De provider en het model waarmee de vertaling is gegenereerd
- Een tijdstempel

Bij volgende runs worden strings met dezelfde bron- en beleidshashes direct vanuit de cache geleverd zonder de LLM aan te roepen. Het standaardpad bevindt zich in de genegeerde map `.internationalizer/`, waardoor het een lokale cache blijft. Stel `tm_path` in op een getrackte locatie als je project het vertaalgeheugen bewust wil delen via versiebeheer. Het controleerbare `.internationalizer.lock`-manifest wordt apart geversioneerd.

<!-- internationalizer:unit markdown:supported-formats -->
## Ondersteunde bestandsindelingen

| Formaat | Extensies | Modus |
|---|---|---|
| JSON | `.json` | Sleutel-waarde (genest, afgevlakt met puntnotatie) |
| YAML | `.yml`, `.yaml` | Sleutel-waarde (behoudt commentaren en volgorde) |
| Markdown | `.md`, `.mdx` | Preambule en secties op H2-niveau |

Markdown-doelbestanden bevatten onzichtbare `internationalizer:unit`-commentaren vóór H2-secties. Dankzij deze stabiele markeringen kan Internationalizer één bronsectie toevoegen, verplaatsen of bewerken zonder niet-gerelateerde secties opnieuw te vertalen. Bestaande niet-gemarkeerde documenten krijgen deze markeringen bij de eerstvolgende succesvolle bijwerking.

<!-- internationalizer:unit markdown:project-type-detection -->
## Projecttypedetectie

`internationalizer detect` herkent je i18n-configuratie door te controleren op:

- Afhankelijkheden in `package.json` voor react-i18next, next-intl of vue-i18n
- Mappenstructuren die overeenkomen met gangbare taalpatronen
- Bestandsextensies en naamgevingsconventies

<!-- internationalizer:unit markdown:architecture -->
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
## Vergelijking met alternatieven

| Functie | Internationalizer | i18next | Crowdin | Generieke LLM |
|---|---|---|---|---|
| LLM-gestuurde vertaling | Ja | Nee | Gedeeltelijk | Ja |
| Stijlgidsen per taal | Ja | Nee | Nee | Nee |
| Woordenlijsthandhaving | Ja | Nee | Ja | Nee |
| Vertaalgeheugen | Ja | Nee | Ja | Nee |
| CLI / lokale uitvoering | Ja | N.v.t. | Nee | Handmatig |
| Git-vriendelijke bestanden | Ja | Ja | Gedeeltelijk | Handmatig |
| Geen SaaS-afhankelijkheid | Ja | Ja | Nee | Varieert |
| Open source (AGPL-3.0) | Ja | Ja | Nee | Varieert |

<!-- internationalizer:unit markdown:license -->
## Licentie

[AGPL-3.0](../../LICENSE)

Zie [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) voor vermeldingen van afhankelijkheden.

<!-- internationalizer:unit markdown:contributing -->
## Bijdragen

Zie [CONTRIBUTING.md](../../CONTRIBUTING.md) voor instructies over het opzetten van de ontwikkelomgeving en richtlijnen. Alle bijdragen vereisen een DCO-ondertekening.
