> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline di internazionalizzazione nativa per l'AI per progetti software. Traduci, convalida e gestisci i file i18n usando gli LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Perché Internationalizer?

La maggior parte degli strumenti i18n sono librerie di runtime (i18next, react-intl) o piattaforme SaaS per la gestione delle chiavi (Crowdin, Lokalise). Nessuno di essi risolve davvero a fondo il problema della traduzione:

- **La traduzione manuale** non scala oltre un numero limitato di lingue
- **Le API di traduzione automatica** (Google Translate, DeepL) ignorano la tua terminologia, il tono e le convenzioni dell'interfaccia utente
- **La traduzione con LLM generici** funziona meglio, ma senza glossari e guide di stile produce risultati incoerenti

Internationalizer è diverso. È una **pipeline CLI** che combina la traduzione con LLM con:

- **Glossari per lingua** — garantiscono una terminologia coerente in tutta l'applicazione
- **Guide di stile per lingua** — definiscono tono, registro di formalità, pluralizzazione e tipografia
- **Memoria di traduzione** — salta le stringhe invariate, risparmiando sui costi delle chiamate API
- **Convalida deterministica** — individua chiavi mancanti o in eccesso, derive della struttura protetta, problemi del glossario ed errori nei plurali o nei messaggi ICU prima del rilascio

<!-- internationalizer:unit markdown:installation -->
## Installazione

Installa da npm:

```bash
npm install -g internationalizer
```

Oppure esegui senza un'installazione globale:

```bash
npx internationalizer --help
```

Il pacchetto npm installa il file binario precompilato corrispondente da npm tramite dipendenze opzionali specifiche per la piattaforma.

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

<!-- internationalizer:unit markdown:npm-packages -->
## Pacchetti npm

- I tag Git e le versioni dei pacchetti npm devono corrispondere, ad esempio `v0.1.0` e `0.1.0`
- Il pacchetto radice `internationalizer` dipende da pacchetti specifici per piattaforma come `internationalizer-darwin-arm64`
- Target npm supportati: macOS arm64/x64, Linux arm64/x64, Windows x64
- La pubblicazione tramite CI richiede un secret di GitHub denominato `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Avvio rapido

1. Crea un file di configurazione nella radice del progetto:

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

2. Imposta la chiave API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Visualizza in anteprima gli elementi che verranno tradotti:

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

<!-- internationalizer:unit markdown:commands -->
## Comandi

### `translate`

Individua le chiavi mancanti o obsolete e traducile tramite un LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Lo stato della traduzione segnala in modo indipendente le condizioni di elementi mancanti, con sorgente obsoleta, con criteri obsoleti, aggiornati e modificati manualmente; in questo modo una modifica manuale non può nascondere una modifica ai sorgenti o ai criteri. I valori con criteri obsoleti vengono segnalati ma ritradotti solo specificando `--refresh-policy`. I valori modificati manualmente non vengono mai sovrascritti automaticamente. Usa `--adopt-existing` quando introduci il manifest su traduzioni già revisionate o quando intendi accettare esplicitamente una modifica manuale revisionata come nuova baseline.

### `validate`

Verifica tutti i file di lingua rispetto ai relativi bundle di origine. La convalida predefinita verifica la copertura strutturale (la percentuale di chiavi di destinazione richieste presenti), segnala le chiavi in eccesso come avvisi e restituisce un errore in presenza di chiavi mancanti, discrepanze di interpolazione o strutture ICU MessageFormat non valide.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` segnala inoltre la percentuale di copertura tradotta. Un valore linguistico identico alla sorgente è considerato non tradotto, a meno che il glossario non contenga una voce identica per sorgente e destinazione per l'intero valore; l'opzione `ignore_case` viene rispettata, ma un termine di glossario presente all'interno di un valore più lungo non costituisce un'esenzione. La modalità strict genera un errore in caso di chiavi in eccesso, valori identici alla sorgente, alterazioni nella struttura di interpolazione/HTML/codice/link Markdown, violazioni del glossario e forme plurali configurate non conformi.

`--require-state` verifica ogni destinazione rispetto al file `.internationalizer.lock`. Restituisce un errore se una chiave non è tracciata o se il valore registrato per sorgente, criteri di traduzione o hash di destinazione risulta obsoleto. Può essere combinato con `--strict`.

I report per l'utente e in formato JSON usano codici di riscontro stabili:

| Codice | Significato |
| --- | --- |
| `missing_key` / `extra_key` | I set di chiavi di origine e destinazione differiscono |
| `blank_translation` | Una sorgente non vuota presenta una destinazione vuota in modalità strict |
| `source_identical` | Un valore linguistico in modalità strict non è stato tradotto |
| `protected_structure_mismatch` | La struttura di interpolazione, HTML, codice o link è cambiata |
| `glossary_violation` | Non è stato trovato alcun termine di destinazione o variante approvata |
| `plural_form_missing` | Una forma plurale configurata per la lingua è assente |
| `icu_message_syntax` | Un messaggio ICU di origine o destinazione presenta errori di sintassi |
| `icu_argument_mismatch` | I nomi, i tipi o gli stili di formattazione degli argomenti ICU differiscono |
| `icu_selector_mismatch` | I selettori differiscono oppure una categoria plurale non è valida per la lingua di destinazione |
| `untracked` | Non è presente alcun record nel manifest per la destinazione |
| `source_stale` | Il contenuto di origine è stato modificato dopo la traduzione registrata |
| `policy_stale` | Il prompt generato o le impostazioni del modello sono cambiati |
| `target_modified` | Il contenuto di destinazione differisce dal record presente nel manifest |

### `detect`

Rileva automaticamente il framework i18n in uso e propone una configurazione.

```bash
internationalizer detect
```

Supporta: react-i18next, next-intl, vue-i18n, JSON standard, documentazione Markdown.

### `glossary`

Gestisci i termini del glossario per lingua applicati durante la traduzione.

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

<!-- internationalizer:unit markdown:configuration-reference -->
## Riferimento per la configurazione

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

Gli identificatori delle impostazioni internazionali devono essere tag BCP 47 corretti come `fr`, `pt-BR` o `sr-Latn-RS`. Le impostazioni internazionali di destinazione con equivalenza canonica vengono rifiutate come duplicati e gli override del provider specifici per lingua corrispondono alla grafia canonicamente equivalente. Nell'esempio precedente, le impostazioni internazionali senza override (incluso il giapponese) ereditano la configurazione globale di Gemini.

I valori in formato ICU MessageFormat vengono analizzati in modo strutturale. Sono supportati argomenti semplici, `select`, `plural`, `selectordinal`, `number`, `date` e `time`, inclusi messaggi annidati, offset dei plurali, selettori di numeri esatti e `#`. La convalida verifica la sintassi, i tipi di argomenti e gli stili di formattazione, gli offset dei plurali, l'identità dei rami del costrutto select e le categorie dei plurali CLDR della lingua di destinazione. L'output del provider che viola queste invarianti viene rifiutato prima che venga scritto un file di lingua o un record della memoria di traduzione.

Con `i18next-v4`, le famiglie di plurali di origine riconosciute vengono espanse durante la traduzione nelle categorie CLDR della lingua di destinazione. Una categoria presente solo nella destinazione usa il valore `_other` della famiglia di origine come modello di traduzione. La convalida strict richiede la presenza di tali categorie di destinazione; le categorie presenti solo nella sorgente sono facoltative per le lingue di destinazione che non le prevedono.

<!-- internationalizer:unit markdown:style-guides -->
## Guide di stile

Le guide di stile sono file Markdown inseriti nel prompt di traduzione dell'LLM. Definiscono tono, registro di formalità, tipografia e altre convenzioni specifiche per ogni lingua.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Convenzioni condivise (`_conventions.md`)

Definiscono le regole valide per tutte le lingue: sintassi di interpolazione, conservazione dell'HTML, convenzioni sui tipi di stringa (pulsanti, etichette, errori) e altro ancora.

### Guide per lingua (`{locale}.md`)

Definiscono regole specifiche per la lingua: registro di formalità (tu vs. Lei/voi), punteggiatura (virgolette caporali, punti interrogativi invertiti), forme plurali, formattazione di date/numeri e un glossario terminologico.

Le guide di stile rappresentano criteri durevoli, non file generati. Internationalizer le legge ma non le riscrive mai. Il loro contenuto viene sottoposto a hash separatamente dal glossario e dal contratto del prompt, in modo che una modifica al codice dell'applicazione non renda obsoleta una traduzione. La modifica di una guida contrassegna intenzionalmente tale lingua per la revisione dei criteri; la sola modifica del testo interno dei prompt non produce questo effetto, a meno che non cambi anche la versione del contratto del prompt.

Vedi [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) per un esempio pratico.

<!-- internationalizer:unit markdown:glossary-format -->
## Formato del glossario

I file di glossario sono array JSON archiviati nel percorso `{glossary_dir}/{locale}.json`:

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

`variants` elenca le altre forme approvate per la destinazione. `enforcement` può essere impostato su `error`, `warning` oppure omesso per applicare il comportamento predefinito di errore. I termini vengono inseriti nel prompt dell'LLM sotto forma di tabella terminologica, garantendo coerenza nelle traduzioni in tutta l'applicazione. Una voce esatta come `{"source":"API","target":"API"}` esenta inoltre il valore completo identico all'originale dalle segnalazioni di mancata traduzione in modalità strict; la presenza del termine `API` all'interno di un valore più lungo non costituisce invece un'esenzione.

<!-- internationalizer:unit markdown:translation-memory -->
## Memoria di traduzione

La memoria di traduzione è memorizzata in un file JSONL (un record JSON per riga). Ciascun record include:

- Bundle, chiave, valore di origine, valore tradotto e impostazione internazionale canonica di destinazione
- Hash di origine, guida di stile, glossario, contratto del prompt e hash combinato dei criteri
- Provider e modello che hanno generato la traduzione
- Un timestamp

Nelle esecuzioni successive, le stringhe aventi gli stessi hash di sorgente e criteri vengono recuperate direttamente dalla cache senza chiamare l'LLM. Il percorso predefinito si trova all'interno della directory ignorata `.internationalizer/`, mantenendo la cache in locale. Imposta `tm_path` su una directory tracciata se il progetto prevede la condivisione esplicita della memoria di traduzione. Il file manifest verificabile `.internationalizer.lock` viene invece gestito con un proprio controllo di versione separato.

<!-- internationalizer:unit markdown:supported-formats -->
## Formati supportati

| Formato | Estensioni | Modalità |
|---|---|---|
| JSON | `.json` | Chiave-valore (annidato, appiattito con notazione puntata) |
| YAML | `.yml`, `.yaml` | Chiave-valore (preserva commenti e ordinamento) |
| Markdown | `.md`, `.mdx` | Preambolo e sezioni di livello H2 |

I file Markdown di destinazione includono commenti invisibili `internationalizer:unit` prima delle sezioni H2. Questi marcatori stabili consentono a Internationalizer di aggiungere, spostare o modificare una sezione di origine senza dover ritradurre le sezioni non correlate. I documenti esistenti privi di marcatori li riceveranno in occasione del successivo aggiornamento eseguito con successo.

<!-- internationalizer:unit markdown:project-type-detection -->
## Rilevamento del tipo di progetto

`internationalizer detect` analizza l'impostazione i18n del progetto controllando:

- Le dipendenze in `package.json` relative a react-i18next, next-intl o vue-i18n
- Le strutture delle directory corrispondenti a schemi di localizzazione diffusi
- Le estensioni dei file e le convenzioni di denominazione

<!-- internationalizer:unit markdown:architecture -->
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
## Confronto con le alternative

| Funzionalità | Internationalizer | i18next | Crowdin | LLM generico |
|---|---|---|---|---|
| Traduzione basata su LLM | Sì | No | Parziale | Sì |
| Guide di stile per lingua | Sì | No | No | No |
| Applicazione del glossario | Sì | No | Sì | No |
| Memoria di traduzione | Sì | No | Sì | No |
| Esecuzione da CLI / in locale | Sì | N/D | No | Manuale |
| File compatibili con Git | Sì | Sì | Parziale | Manuale |
| Nessuna dipendenza SaaS | Sì | Sì | No | Variabile |
| Open source (AGPL-3.0) | Sì | Sì | No | Variabile |

<!-- internationalizer:unit markdown:license -->
## Licenza

[AGPL-3.0](../../LICENSE)

Consulta [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) per le comunicazioni sulle dipendenze.

<!-- internationalizer:unit markdown:contributing -->
## Contribuire

Consulta [CONTRIBUTING.md](../../CONTRIBUTING.md) per le linee guida e la configurazione dell'ambiente di sviluppo. Tutti i contributi richiedono la firma del DCO (Developer Certificate of Origin).
