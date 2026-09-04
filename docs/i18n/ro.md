> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Flux de internaționalizare nativ AI pentru proiecte software. Traduceți, validați și gestionați fișierele i18n folosind LLM-uri.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## De ce Internationalizer?

Majoritatea instrumentelor i18n sunt fie biblioteci pentru runtime (i18next, react-intl), fie platforme SaaS pentru gestionarea cheilor (Crowdin, Lokalise). Niciunul dintre ele nu rezolvă cu adevărat problema traducerii:

- **Traducerea manuală** nu este scalabilă dincolo de câteva limbi
- **API-urile de traducere automată** (Google Translate, DeepL) vă ignoră terminologia, tonul și convențiile interfeței utilizator
- **Traducerea generică prin LLM** funcționează mai bine, însă fără glosare și ghiduri de stil obțineți rezultate inconsecvente

Internationalizer este diferit. Este un **pipeline CLI** care combină traducerea prin LLM cu:

- **Glosare per limbă** — impun o terminologie consecventă în întreaga aplicație
- **Ghiduri de stil per limbă** — controlează tonul, nivelul de formalitate, pluralul și tipografia
- **Memorie de traducere** — omite șirurile neschimbate, economisind costurile apelurilor API
- **Validare deterministă** — detectează cheile lipsă sau suplimentare, abaterile de structură protejată, problemele de glosar și erorile de plural sau ICU înainte de livrare

<!-- internationalizer:unit markdown:installation -->
## Instalare

Instalați din npm:

```bash
npm install -g internationalizer
```

Sau rulați fără instalare globală:

```bash
npx internationalizer --help
```

Pachetul npm instalează binarul precompilat corespunzător din npm prin dependențe opționale specifice fiecărei platforme.

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

<!-- internationalizer:unit markdown:npm-packages -->
## Pachete npm

- Etichetele Git și versiunile pachetelor npm trebuie să coincidă, de exemplu `v0.1.0` și `0.1.0`
- Pachetul rădăcină `internationalizer` depinde de pachete de platformă precum `internationalizer-darwin-arm64`
- Ținte npm acceptate: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publicarea automată în CI necesită un secret GitHub denumit `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Pornire rapidă

1. Creați un fișier de configurare în rădăcina proiectului:

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

2. Setați cheia API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Previzualizați ce va fi tradus:

```bash
internationalizer translate --dry-run
```

4. Rulați traducerea:

```bash
internationalizer translate
```

5. Validați toate pachetele de limbă:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## Comenzi

### `translate`

Găsiți cheile lipsă sau învechite și traduceți-le printr-un LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Starea traducerii raportează în mod independent stările lipsă, sursă învechită, politică învechită, actuală și editată manual, astfel încât o editare manuală nu poate masca o modificare a sursei sau a politicilor. Valorile cu politică învechită sunt raportate, dar sunt retraduse doar cu `--refresh-policy`. Valorile editate manual nu sunt niciodată suprascrise automat. Folosiți `--adopt-existing` atunci când asociați manifestul unor traduceri revizuite sau când acceptați în mod explicit o editare manuală revizuită drept nou nivel de referință.

### `validate`

Verificați toate fișierele de limbă în raport cu pachetele sursă. Validarea implicită verifică acoperirea structurală (procentul de chei țintă obligatorii prezente), raportează cheile suplimentare ca avertismente și eșuează în caz de chei lipsă, nepotriviri de interpolare sau structură ICU MessageFormat nevalidă.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` raportează și acoperirea tradusă. O valoare lingvistică identică cu sursa sa este considerată netradusă, cu excepția cazului în care glosarul conține explicit o intrare exactă cu aceeași sursă și aceeași țintă pentru valoarea completă; `ignore_case` este respectat, dar un termen de glosar inclus într-o valoare mai lungă nu constituie o scutire. Modul strict eșuează la chei suplimentare, valori identice cu sursa, structură de interpolare/HTML/cod/linkuri Markdown modificată, încălcări ale glosarului și lipsa formelor de plural configurate.

`--require-state` verifică fiecare țintă în raport cu `.internationalizer.lock`. Eșuează când o cheie este neurmărită sau când hash-ul înregistrat pentru sursă, politica de traducere ori țintă este învechit. Poate fi combinat cu `--strict`.

Rapoartele pentru utilizator și JSON folosesc coduri de constatare stabile:

| Cod | Semnificație |
| --- | --- |
| `missing_key` / `extra_key` | Seturile de chei sursă și țintă diferă |
| `blank_translation` | O sursă nevidă are o țintă goală în modul strict |
| `source_identical` | O valoare lingvistică în modul strict rămâne netradusă |
| `protected_structure_mismatch` | Structura de interpolare, HTML, cod sau linkuri a fost modificată |
| `glossary_violation` | Nu a fost găsit niciun termen țintă sau nicio variantă aprobată |
| `plural_form_missing` | O formă de plural configurată pentru limbă lipsește |
| `icu_message_syntax` | Un mesaj ICU sursă sau țintă este formatat incorect |
| `icu_argument_mismatch` | Numele argumentelor ICU, tipurile sau stilurile de formatare diferă |
| `icu_selector_mismatch` | Selectorii diferă sau o categorie de plural este nevalidă pentru limba țintă |
| `untracked` | Nu există nicio înregistrare în manifest pentru țintă |
| `source_stale` | Conținutul sursă s-a modificat după înregistrarea traducerii |
| `policy_stale` | Promptul generat sau setările modelului s-au modificat |
| `target_modified` | Conținutul țintă diferă de înregistrarea din manifest |

### `detect`

Detectați automat framework-ul i18n și sugerați o configurație.

```bash
internationalizer detect
```

Acceptă: react-i18next, next-intl, vue-i18n, JSON simplu, documentații Markdown.

### `glossary`

Gestionați termenii de glosar per limbă a căror aplicare este impusă în timpul traducerii.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Gestionați memoria de traducere (cache JSONL cu șiruri traduse anterior).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## Referință de configurare

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

Identificatorii de limbă trebuie să fie etichete BCP 47 valide, precum `fr`, `pt-BR` sau `sr-Latn-RS`. Limbile țintă echivalente canonic sunt respinse ca duplicate, iar suprascrierile de furnizor specifice unei limbi se potrivesc cu scrierea canonică echivalentă. În exemplul de mai sus, limbile fără o configurare dedicată — inclusiv japoneza — moștenesc configurația globală Gemini.

Valorile ICU MessageFormat sunt analizate structural. Argumentele simple, `select`, `plural`, `selectordinal`, `number`, `date` și `time` sunt acceptate, inclusiv mesajele imbricate, decalajele de plural, selectorii numerici exacți și `#`. Validarea verifică sintaxa, tipurile de argumente și stilurile de formatare, decalajele de plural, identitatea ramurilor select și categoriile de plural CLDR pentru limba țintă. Răspunsul furnizorului care încalcă aceste condiții invariante este respins înainte de a fi scris un fișier de limbă sau o înregistrare în memoria de traducere.

Cu `i18next-v4`, familiile sursă de plural recunoscute sunt extinse în timpul traducerii la categoriile CLDR ale limbii țintă. O categorie existentă doar în limba țintă folosește valoarea `_other` a familiei sursă drept șablon de traducere. Validarea strictă impune acele categorii țintă; categoriile existente doar în sursă sunt opționale pentru limbile țintă care nu le folosesc.

<!-- internationalizer:unit markdown:style-guides -->
## Ghiduri de stil

Ghidurile de stil sunt fișiere Markdown transmise în promptul de traducere trimis LLM-ului. Ele controlează tonul, nivelul de formalitate, tipografia și alte convenții specifice fiecărei limbi.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Convenții comune (`_conventions.md`)

Definiți reguli aplicabile tuturor limbilor: sintaxa de interpolare, păstrarea HTML, convențiile privind tipurile de șiruri (butoane vs. etichete vs. erori) etc.

### Ghiduri per limbă (`{locale}.md`)

Definiți reguli specifice limbii: registrul de adresare (tu vs. dumneavoastră), punctuația (ghilimele, semne de întrebare inversate), formele de plural, formatarea datei și a numerelor și un glosar terminologic.

Ghidurile de stil sunt date de intrare durabile pentru politici, nu conținut generat. Internationalizer le citește, dar nu le suprascrie niciodată. Conținutul lor este asociat unui hash calculat separat de glosar și contractul de prompt, astfel încât o modificare a codului aplicației nu face o traducere să devină învechită. Editarea unui ghid marchează intenționat respectiva limbă pentru revizuirea politicii; ajustarea formulării interne a promptului nu face acest lucru, cu excepția cazului în care se modifică și versiunea contractului de prompt.

Consultați [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) pentru un exemplu funcțional.

<!-- internationalizer:unit markdown:glossary-format -->
## Formatul glosarelor

Fișierele de glosar sunt matrice JSON stocate în `{glossary_dir}/{locale}.json`:

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

`variants` listează alte forme aprobate în limba țintă. `enforcement` poate fi `error`, `warning` sau poate fi omisă pentru comportamentul implicit de eroare. Termenii sunt transmiși în promptul LLM sub forma unui tabel de terminologie, asigurând o traducere consecventă în întreaga aplicație. O intrare exactă precum `{"source":"API","target":"API"}` scutește, de asemenea, acea valoare completă identică cu sursa de la constatările de valori netraduse din modul strict; nu scutește însă o valoare mai lungă care conține doar `API`.

<!-- internationalizer:unit markdown:translation-memory -->
## Memorie de traducere

Memoria de traducere este stocată sub forma unui fișier JSONL (o înregistrare JSON pe fiecare linie). Fiecare înregistrare conține:

- Pachetul, cheia, valoarea sursă, valoarea tradusă și codul canonic al limbii țintă
- Hash-urile pentru sursă, ghidul de stil, glosar, contractul de prompt și hash-ul combinat de politică
- Furnizorul și modelul care au produs traducerea
- Un marcaj temporal

La rulările ulterioare, șirurile cu aceleași hash-uri de sursă și politică sunt servite direct din cache, fără a apela LLM-ul. Calea implicită se află în directorul ignorat `.internationalizer/`, astfel încât rămâne un cache local. Setați `tm_path` către o locație urmărită prin controlul versiunilor dacă proiectul partajează intenționat memoria de traducere. Manifestul revizuibil `.internationalizer.lock` este versionat separat.

<!-- internationalizer:unit markdown:supported-formats -->
## Formate acceptate

| Format | Extensii | Mod |
|--------|-----------|------|
| JSON | `.json` | Cheie-valoare (imbricate, aplatizate prin notație cu punct) |
| YAML | `.yml`, `.yaml` | Cheie-valoare (păstrează comentariile și ordonarea) |
| Markdown | `.md`, `.mdx` | Preambul și secțiuni la nivel H2 |

Fișierele Markdown țintă conțin comentarii invizibile `internationalizer:unit` înainte de secțiunile H2. Acești markeri stabili îi permit instrumentului Internationalizer să adauge, să mute sau să editeze o secțiune sursă fără a retraduce secțiunile neafectate. Documentele existente nemarcate primesc markeri la următoarea actualizare reușită.

<!-- internationalizer:unit markdown:project-type-detection -->
## Detectarea tipului de proiect

`internationalizer detect` identifică configurarea dumneavoastră i18n verificând:

- Dependențele din `package.json` pentru react-i18next, next-intl sau vue-i18n
- Structurile de directoare care corespund modelelor comune de configurare a limbilor
- Extensiile de fișiere și convențiile de denumire

<!-- internationalizer:unit markdown:architecture -->
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
## Comparație cu alternativele

| Funcționalitate | Internationalizer | i18next | Crowdin | LLM generic |
|---------|------------------|---------|---------|-------------|
| Traducere bazată pe LLM | Da | Nu | Parțial | Da |
| Ghiduri de stil per limbă | Da | Nu | Nu | Nu |
| Impunerea terminologiei din glosar | Da | Nu | Da | Nu |
| Memorie de traducere | Da | Nu | Da | Nu |
| Rulare CLI / locală | Da | N/A | Nu | Manual |
| Fișiere adaptate pentru Git | Da | Da | Parțial | Manual |
| Fără dependență de SaaS | Da | Da | Nu | Variază |
| Open source (AGPL-3.0) | Da | Da | Nu | Variază |

<!-- internationalizer:unit markdown:license -->
## Licență

[AGPL-3.0](../../LICENSE)

Consultați [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) pentru notificările privind dependențele.

<!-- internationalizer:unit markdown:contributing -->
## Contribuții

Consultați [CONTRIBUTING.md](../../CONTRIBUTING.md) pentru configurarea mediului de dezvoltare și ghiduri. Toate contribuțiile necesită asumarea DCO.
