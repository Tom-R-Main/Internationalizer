> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

KI-native Internationalisierungs-Pipeline für Softwareprojekte. Übersetze, validiere und verwalte i18n-Dateien mithilfe von LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Warum Internationalizer?

Die meisten i18n-Tools sind entweder Laufzeitbibliotheken (i18next, react-intl) oder SaaS-Plattformen zur Schlüsselverwaltung (Crowdin, Lokalise). Keines davon löst das eigentliche Übersetzungsproblem zufriedenstellend:

- **Manuelle Übersetzung** skaliert über wenige Sprachen hinaus nicht mehr
- **APIs für maschinelle Übersetzung** (Google Translate, DeepL) ignorieren deine Terminologie, deinen Tonfall und deine UI-Konventionen
- **Generische LLM-Übersetzungen** funktionieren besser, führen aber ohne Glossare und Styleguides zu inkonsistenten Ergebnissen

Internationalizer geht einen anderen Weg. Es ist eine **CLI-Pipeline**, die LLM-Übersetzung mit folgenden Bausteinen kombiniert:

- **Sprachspezifische Glossare** – setzen eine konsistente Terminologie in deiner gesamten App durch
- **Sprachspezifische Styleguides** – steuern Tonfall, Formalitätsgrad, Pluralisierung und Typografie
- **Translation Memory** – überspringt unveränderte Strings und spart Kosten bei API-Aufrufen
- **Deterministische Validierung** – fängt fehlende oder überflüssige Schlüssel, Abweichungen geschützter Strukturen, Glossarprobleme sowie Plural- oder ICU-Fehler ab, bevor sie ausgeliefert werden

<!-- internationalizer:unit markdown:installation -->
## Installation

Über npm installieren:

```bash
npm install -g internationalizer
```

Oder ohne globale Installation ausführen:

```bash
npx internationalizer --help
```

Das npm-Paket installiert das passende vorkompilierte Binary von npm über plattformspezifische optionale Abhängigkeiten.

Mit Go installieren:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Oder aus dem Quellcode kompilieren:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm-Pakete

- Git-Tags und npm-Paketversionen müssen übereinstimmen, beispielsweise `v0.1.0` und `0.1.0`
- Das Root-Paket `internationalizer` hängt von Plattformpaketen wie `internationalizer-darwin-arm64` ab
- Unterstützte npm-Zielplattformen: macOS arm64/x64, Linux arm64/x64, Windows x64
- Das Veröffentlichen über CI erfordert ein GitHub-Secret namens `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Schnellstart

1. Erstelle eine Konfigurationsdatei im Root-Verzeichnis deines Projekts:

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

2. Setze deinen API-Schlüssel:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Zeige eine Vorschau der anstehenden Übersetzungen an:

```bash
internationalizer translate --dry-run
```

4. Führe die Übersetzung aus:

```bash
internationalizer translate
```

5. Validiere alle Locales:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## Befehle

### `translate`

Fehlende oder veraltete Schlüssel ermitteln und über ein LLM übersetzen.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Der Übersetzungsstatus meldet fehlende, quellseitig veraltete (source-stale), richtlinienseitig veraltete (policy-stale), aktuelle und manuell bearbeitete Zustände unabhängig voneinander, sodass eine manuelle Bearbeitung eine Quell- oder Richtlinienänderung nicht verdecken kann. Richtlinienseitig veraltete Werte werden gemeldet, aber nur mit `--refresh-policy` neu übersetzt. Manuell bearbeitete Werte werden niemals automatisch überschrieben. Verwende `--adopt-existing`, wenn du das Manifest für bereits geprüfte Übersetzungen einführst oder eine geprüfte manuelle Änderung ausdrücklich als neue Baseline übernimmst.

### `validate`

Alle Locale-Dateien mit ihren Quell-Bundles abgleichen. Die Standardvalidierung prüft die strukturelle Abdeckung (den Prozentsatz der vorhandenen erforderlichen Zielschlüssel), meldet überzählige Schlüssel als Warnungen und schlägt bei fehlenden Schlüsseln, Interpolationsabweichungen oder ungültiger ICU-MessageFormat-Struktur fehl.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` gibt zusätzlich die übersetzte Abdeckung aus. Ein sprachlicher Wert, der mit seiner Quelle identisch ist, gilt als unübersetzt, es sei denn, das Glossar enthält ausdrücklich einen exakten Eintrag mit identischer Quelle und identischem Ziel für den vollständigen Wert; `ignore_case` wird beachtet, aber ein in einen längeren Wert eingebetteter Glossarbegriff gilt nicht als Ausnahme. Der Strict-Modus schlägt bei überzähligen Schlüsseln, quellidentischen Werten, veränderten Interpolations-, HTML-, Code- oder Markdown-Link-Strukturen, Glossarverstößen und konfigurierten Pluralformen fehl.

`--require-state` gleicht jedes Ziel mit `.internationalizer.lock` ab. Die Prüfung schlägt fehl, wenn ein Schlüssel nicht erfasst ist (untracked) oder wenn der hinterlegte Hash für Quelle, Übersetzungsrichtlinie oder Ziel veraltet ist. Die Option kann mit `--strict` kombiniert werden.

Ausgaben für Menschen und maschinenlesbare JSON-Berichte verwenden stabile Finding-Codes:

| Code | Bedeutung |
| --- | --- |
| `missing_key` / `extra_key` | Quell- und Zielschlüsselsätze weichen voneinander ab |
| `blank_translation` | Eine nicht leere Quelle hat im Strict-Modus ein leeres Ziel |
| `source_identical` | Ein sprachlicher Wert bleibt im Strict-Modus unübersetzt |
| `protected_structure_mismatch` | Struktur von Interpolationen, HTML, Code oder Links wurde verändert |
| `glossary_violation` | Kein zulässiger Zielbegriff oder keine zulässige Variante gefunden |
| `plural_form_missing` | Eine konfigurierte Pluralform des Locales fehlt |
| `icu_message_syntax` | Eine ICU-Nachricht in der Quelle oder im Ziel ist syntaktisch fehlerhaft |
| `icu_argument_mismatch` | ICU-Argumentnamen, -Typen oder -Formatierungsstile weichen ab |
| `icu_selector_mismatch` | Selektoren weichen ab oder eine Pluralkategorie ist für das Ziel-Locale ungültig |
| `untracked` | Für das Ziel existiert kein Eintrag im Manifest |
| `source_stale` | Quellinhalt wurde nach der erfassten Übersetzung geändert |
| `policy_stale` | Der generierte Prompt oder die Modelleinstellungen wurden geändert |
| `target_modified` | Zielinhalt weicht vom Manifest-Eintrag ab |

### `detect`

Das i18n-Framework automatisch erkennen und eine Konfiguration vorschlagen.

```bash
internationalizer detect
```

Unterstützt: react-i18next, next-intl, vue-i18n, Standard-JSON, Markdown-Dokumentation.

### `glossary`

Sprachspezifische Glossarbegriffe verwalten, die während der Übersetzung durchgesetzt werden.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Das Translation Memory verwalten (JSONL-Cache zuvor übersetzter Strings).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## Konfigurationsreferenz

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

Locale-IDs müssen wohlgeformte BCP-47-Tags wie `fr`, `pt-BR` oder `sr-Latn-RS` sein. Kanonisch äquivalente Ziel-Locales werden als Duplikate abgelehnt, und locale-spezifische Provider-Overrides berücksichtigen die kanonisch äquivalente Schreibweise. Im obigen Beispiel erben Locales ohne Override – einschließlich Japanisch – die globale Gemini-Konfiguration.

Werte im ICU MessageFormat werden strukturell geparst. Einfache Argumente, `select`, `plural`, `selectordinal`, `number`, `date` und `time` werden unterstützt, einschließlich verschachtelter Nachrichten, Plural-Offsets, Selektoren für exakte Zahlen und `#`. Die Validierung prüft Syntax, Argumenttypen und Formatierungsstile, Plural-Offsets, Übereinstimmung der Select-Zweige sowie CLDR-Pluralkategorien des Ziel-Locales. Provider-Ausgaben, die diese Invarianten verletzen, werden verworfen, bevor eine Locale-Datei oder ein Datensatz im Translation Memory geschrieben wird.

Mit `i18next-v4` werden erkannte Pluralfamilien der Quelle bei der Übersetzung in die CLDR-Kategorien des Ziel-Locales expandiert. Eine Kategorie, die nur im Ziel existiert, verwendet den `_other`-Wert der Quellfamilie als Übersetzungsvorlage. Die strikte Validierung verlangt diese Zielkategorien; reine Quellkategorien sind für Ziel-Locales optional, die diese nicht verwenden.

<!-- internationalizer:unit markdown:style-guides -->
## Styleguides

Styleguides sind Markdown-Dateien, die in den LLM-Übersetzungs-Prompt eingefügt werden. Sie steuern Tonfall, Formalitätsgrad, Typografie und weitere sprachspezifische Konventionen.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Gemeinsame Konventionen (`_conventions.md`)

Definiert Regeln, die für alle Sprachen gelten: Interpolationssyntax, Beibehaltung von HTML, Konventionen für String-Typen (Schaltflächen vs. Labels vs. Fehlermeldungen) usw.

### Sprachspezifische Leitfäden (`{locale}.md`)

Definiert sprachspezifische Regeln: Höflichkeitsformen (du vs. Sie), Zeichensetzung (Guillemets, umgekehrte Fragezeichen), Pluralformen, Datums- und Zahlenformate sowie ein Terminologieglossar.

Styleguides sind dauerhafte Richtlinien-Eingaben (Policy Inputs), keine generierten Ausgaben. Internationalizer liest sie ein, verändert sie jedoch niemals. Ihr Inhalt wird getrennt vom Glossar und dem Prompt-Contract gehasht, sodass Änderungen am Anwendungscode eine Übersetzung nicht als veraltet einstufen. Das Bearbeiten eines Guides markiert dieses Locale bewusst für eine Richtlinienüberprüfung; eine Änderung des internen Prompt-Wortlauts tut dies nicht, es sei denn, die Version des Prompt-Contracts ändert sich ebenfalls.

Ein funktionierendes Beispiel findest du unter [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/).

<!-- internationalizer:unit markdown:glossary-format -->
## Glossarformat

Glossardateien sind JSON-Arrays, die in `{glossary_dir}/{locale}.json` gespeichert werden:

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

`variants` listet weitere zulässige Zielformen auf. `enforcement` kann den Wert `error` oder `warning` annehmen oder für das standardmäßige Error-Verhalten weggelassen werden. Begriffe werden als Terminologietabelle in den LLM-Prompt übergeben, um eine konsistente Übersetzung in deiner gesamten Anwendung sicherzustellen. Ein exakter Eintrag wie `{"source":"API","target":"API"}` nimmt diesen vollständigen quellidentischen Wert zudem von strikten Befunden zu unübersetzten Werten aus; ein längerer Wert, der `API` lediglich enthält, wird dadurch jedoch nicht freigestellt.

<!-- internationalizer:unit markdown:translation-memory -->
## Translation Memory

Das Translation Memory wird als JSONL-Datei gespeichert (ein JSON-Datensatz pro Zeile). Jeder Datensatz enthält:

- Das Bundle, den Schlüssel, den Quellwert, den übersetzten Wert und das kanonische Ziel-Locale
- Hashes für Quelle, Styleguide, Glossar, Prompt-Contract und die kombinierte Richtlinie
- Den Provider und das Modell, die die Übersetzung erzeugt haben
- Einen Zeitstempel

Bei nachfolgenden Durchläufen werden Strings mit denselben Quell- und Richtlinien-Hashes direkt aus dem Cache bedient, ohne das LLM aufzurufen. Der Standardpfad liegt im ignorierten Verzeichnis `.internationalizer/`, bleibt also ein lokaler Cache. Setze `tm_path` auf einen versionierten Pfad, falls dein Projekt das Translation Memory gezielt teilen soll. Das prüfbare Manifest `.internationalizer.lock` wird separat versioniert.

<!-- internationalizer:unit markdown:supported-formats -->
## Unterstützte Formate

| Format | Erweiterungen | Modus |
|--------|--------------|-------|
| JSON | `.json` | Schlüssel-Wert (verschachtelt, mit Punktnotation verflacht) |
| YAML | `.yml`, `.yaml` | Schlüssel-Wert (behält Kommentare und Reihenfolge bei) |
| Markdown | `.md`, `.mdx` | Präambel und Abschnitte auf H2-Ebene |

Markdown-Zieldateien enthalten unsichtbare `internationalizer:unit`-Kommentare vor H2-Abschnitten. Diese stabilen Markierungen ermöglichen es Internationalizer, einzelne Quellabschnitte hinzuzufügen, zu verschieben oder zu bearbeiten, ohne nicht betroffene Abschnitte erneut zu übersetzen. Bestehende Dokumente ohne Markierungen erhalten diese bei der nächsten erfolgreichen Aktualisierung.

<!-- internationalizer:unit markdown:project-type-detection -->
## Erkennung des Projekttyps

`internationalizer detect` ermittelt dein i18n-Setup anhand folgender Kriterien:

- Abhängigkeiten in `package.json` nach react-i18next, next-intl oder vue-i18n
- Verzeichnisstrukturen, die gängigen Locale-Mustern entsprechen
- Dateierweiterungen und Benennungskonventionen

<!-- internationalizer:unit markdown:architecture -->
## Architektur

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
## Vergleich mit Alternativen

| Funktion | Internationalizer | i18next | Crowdin | Generisches LLM |
|----------|-------------------|---------|---------|-----------------|
| LLM-gestützte Übersetzung | Ja | Nein | Teilweise | Ja |
| Sprachspezifische Styleguides | Ja | Nein | Nein | Nein |
| Glossardurchsetzung | Ja | Nein | Ja | Nein |
| Translation Memory | Ja | Nein | Ja | Nein |
| CLI- / lokale Ausführung | Ja | k. A. | Nein | Manuell |
| Git-freundliche Dateien | Ja | Ja | Teilweise | Manuell |
| Keine SaaS-Abhängigkeit | Ja | Ja | Nein | Variiert |
| Open Source (AGPL-3.0) | Ja | Ja | Nein | Variiert |

<!-- internationalizer:unit markdown:license -->
## Lizenz

[AGPL-3.0](../../LICENSE)

Hinweise zu Drittanbieter-Lizenzen findest du in [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md).

<!-- internationalizer:unit markdown:contributing -->
## Mitwirken

Informationen zum Einrichten der Entwicklungsumgebung und Leitfäden findest du in [CONTRIBUTING.md](../../CONTRIBUTING.md). Alle Beiträge erfordern ein DCO-Sign-off.
