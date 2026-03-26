<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

KI-native Internationalisierungs-Pipeline für Softwareprojekte. Übersetzen, validieren und verwalten Sie i18n-Dateien mithilfe von LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Warum Internationalizer?

Die meisten i18n-Tools sind entweder Laufzeitbibliotheken (i18next, react-intl) oder SaaS-Plattformen zur Schlüsselverwaltung (Crowdin, Lokalise). Keines davon löst das eigentliche Übersetzungsproblem gut:

- **Manuelle Übersetzung** skaliert nicht über ein paar Sprachen hinaus
- **Maschinelle Übersetzungs-APIs** (Google Translate, DeepL) ignorieren Ihre Terminologie, Ihren Tonfall und Ihre UI-Konventionen
- **Generische LLM-Übersetzung** funktioniert besser, aber ohne Glossare und Styleguides erhalten Sie inkonsistente Ergebnisse

Internationalizer ist anders. Es ist eine **CLI-Pipeline**, die LLM-Übersetzung kombiniert mit:

- **Sprachspezifischen Glossaren** — erzwingen eine konsistente Terminologie in Ihrer gesamten App
- **Sprachspezifischen Styleguides** — steuern Tonfall, Formalität, Pluralbildung und Typografie
- **Translation Memory** — überspringt unveränderte Zeichenfolgen und spart Kosten für API-Aufrufe
- **Schlüsselvalidierung** — erkennt fehlende Übersetzungen und Interpolationsfehler vor der Auslieferung

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

## npm-Pakete

- Git-Tags und npm-Paketversionen müssen übereinstimmen, zum Beispiel `v0.1.0` und `0.1.0`
- Das Hauptpaket `internationalizer` hängt von Plattformpaketen wie `internationalizer-darwin-arm64` ab
- Unterstützte npm-Ziele: macOS arm64/x64, Linux arm64/x64, Windows x64
- Die CI-Veröffentlichung erfordert ein GitHub-Secret namens `NPM_TOKEN`

## Schnellstart

1. Erstellen Sie eine Konfigurationsdatei im Stammverzeichnis Ihres Projekts:

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

2. Legen Sie Ihren API-Schlüssel fest:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Vorschau der zu übersetzenden Inhalte anzeigen:

```bash
internationalizer translate --dry-run
```

4. Übersetzung ausführen:

```bash
internationalizer translate
```

5. Alle Gebietsschemas validieren:

```bash
internationalizer validate
```

## Befehle

### `translate`

Fehlende Schlüssel finden und über ein LLM übersetzen.

```bash
internationalizer translate                    # Alle Gebietsschemas übersetzen
internationalizer translate -l fr              # Nur Französisch übersetzen
internationalizer translate --dry-run          # Vorschau ohne API-Aufrufe
internationalizer translate --batch-size 20    # Kleinere Batches
internationalizer translate --concurrency 2    # Weniger parallele Aufrufe
```

### `validate`

Alle Gebietsschema-Dateien auf fehlende Schlüssel, zusätzliche Schlüssel und Interpolationsfehler prüfen.

```bash
internationalizer validate                     # Menschenlesbare Ausgabe
internationalizer validate --json              # Maschinenlesbares JSON
internationalizer validate -q                  # Nur Exit-Code
```

### `detect`

Das i18n-Framework automatisch erkennen und eine Konfiguration vorschlagen.

```bash
internationalizer detect
```

Unterstützt: react-i18next, next-intl, vue-i18n, Vanilla JSON, Markdown-Dokumente.

### `glossary`

Sprachspezifische Glossarbegriffe verwalten, die während der Übersetzung erzwungen werden.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Translation Memory verwalten (JSONL-Cache von zuvor übersetzten Zeichenfolgen).

```bash
internationalizer tm stats                     # Datensatzanzahl anzeigen
internationalizer tm export                    # Als JSON exportieren
internationalizer tm clear --force             # Alle Datensätze löschen
```

## Konfigurationsreferenz

```yaml
# .internationalizer.yml

# Quellsprache (Standard: en)
source_locale: en

# Zielsprachen für die Übersetzung (erforderlich)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Pfad zur Quell-Gebietsschema-Datei (erforderlich)
source_path: locales/en.json

# LLM-Anbieter-Einstellungen
llm:
  # Anbieter: "anthropic", "openai", "gemini" oder "openrouter" (Standard: gemini)
  provider: gemini

  # Standard-Modellnamen nach Anbieter:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Umgebungsvariable, die den API-Schlüssel enthält
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Basis-URL für OpenAI-kompatible Endpunkte (optional)
  # base_url: https://api.openai.com

# Schlüssel pro LLM-Aufruf (Standard: 40)
batch_size: 40

# Parallele LLM-Aufrufe (Standard: 4)
concurrency: 4

# Verzeichnis mit sprachspezifischen Styleguide-Markdown-Dateien (Standard: style-guides)
style_guides_dir: style-guides

# Verzeichnis mit sprachspezifischen Glossar-JSON-Dateien (Standard: glossary)
glossary_dir: glossary

# Pfad zur Translation-Memory-Datei (Standard: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## Styleguides

Styleguides sind Markdown-Dateien, die in den LLM-Übersetzungs-Prompt eingefügt werden. Sie steuern Tonfall, Formalität, Typografie und andere sprachspezifische Konventionen.

```
style-guides/
  _conventions.md    # Gemeinsame Regeln für alle Sprachen
  fr.md              # Französisch-spezifische Regeln
  ja.md              # Japanisch-spezifische Regeln
  ar.md              # Arabisch-spezifische Regeln
```

### Gemeinsame Konventionen (`_conventions.md`)

Definieren Sie Regeln, die für alle Sprachen gelten: Interpolationssyntax, Beibehaltung von HTML, Konventionen für Zeichenfolgentypen (Schaltflächen vs. Labels vs. Fehler) usw.

### Sprachspezifische Guides (`{locale}.md`)

Definieren Sie sprachspezifische Regeln: Formalitätsregister (Du vs. Sie), Interpunktion (Guillemets, umgekehrte Fragezeichen), Pluralformen, Datums-/Zahlenformatierung und ein Terminologie-Glossar.

Ein funktionierendes Beispiel finden Sie unter [`examples/react-app/style-guides/`](examples/react-app/style-guides/).

## Glossar-Format

Glossardateien sind JSON-Arrays, die in `{glossary_dir}/{locale}.json` gespeichert werden:

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

Begriffe werden als Terminologietabelle in den LLM-Prompt eingefügt, um eine konsistente Übersetzung von Schlüsselbegriffen in Ihrer gesamten Anwendung sicherzustellen.

## Translation Memory

Das Translation Memory wird als JSONL-Datei gespeichert (ein JSON-Datensatz pro Zeile). Jeder Datensatz enthält:

- Den Quellschlüssel und -wert
- Den übersetzten Wert
- Einen SHA-256-Hash des Quellwerts
- Einen Zeitstempel

Bei nachfolgenden Ausführungen werden unveränderte Zeichenfolgen aus dem TM-Cache bereitgestellt, ohne das LLM aufzurufen, was sowohl Zeit als auch API-Kosten spart. Die TM-Datei ist Git-freundlich und kann zusammen mit Ihren Gebietsschema-Dateien committet werden.

## Unterstützte Formate

| Format | Erweiterungen | Modus |
|--------|-----------|------|
| JSON | `.json` | Schlüssel-Wert (verschachtelt, Punktnotation abgeflacht) |
| YAML | `.yml`, `.yaml` | Schlüssel-Wert (behält Kommentare und Reihenfolge bei) |
| Markdown | `.md`, `.mdx` | Übersetzung des gesamten Dokuments |

## Erkennung des Projekttyps

`internationalizer detect` identifiziert Ihr i18n-Setup durch Überprüfung von:

- `package.json`-Abhängigkeiten für react-i18next, next-intl oder vue-i18n
- Verzeichnisstrukturen, die gängigen Gebietsschema-Mustern entsprechen
- Dateierweiterungen und Namenskonventionen

## Architektur

```
cmd/internationalizer/     CLI-Einstiegspunkt und Befehlsdefinitionen
internal/
  config/                  Laden der YAML-Konfiguration mit Standardwerten
  detect/                  Automatische Erkennung des Projekttyps
  formats/                 Format-Parser (JSON, YAML, Markdown)
  glossary/                Sprachspezifische Glossarverwaltung
  llm/                     LLM-Anbieter-Schnittstelle + Implementierungen
    anthropic.go           Anthropic Claude-Backend
    openai.go              OpenAI / kompatibles Backend
    gemini.go              Google Gemini über AI Studio-Backend
                           OpenRouter verwendet openai.go mit benutzerdefinierter base_url
  styleguide/              Styleguide-Lader
  tm/                      JSONL Translation Memory
  translate/               Übersetzungs-Orchestrator
  validate/                Gebietsschema-Validierung und Diffing
```

## Vergleich mit Alternativen

| Funktion | Internationalizer | i18next | Crowdin | Generisches LLM |
|---------|------------------|---------|---------|-------------|
| LLM-gestützte Übersetzung | Ja | Nein | Teilweise | Ja |
| Sprachspezifische Styleguides | Ja | Nein | Nein | Nein |
| Glossar-Erzwingung | Ja | Nein | Ja | Nein |
| Translation Memory | Ja | Nein | Ja | Nein |
| CLI / lokale Ausführung | Ja | N/A | Nein | Manuell |
| Git-freundliche Dateien | Ja | Ja | Teilweise | Manuell |
| Keine SaaS-Abhängigkeit | Ja | Ja | Nein | Variiert |
| Open Source (AGPL-3.0) | Ja | Ja | Nein | Variiert |

## Lizenz

[AGPL-3.0](LICENSE)

## Mitwirken

Siehe [CONTRIBUTING.md](CONTRIBUTING.md) für Entwicklungs-Setup und Richtlinien. Alle Beiträge erfordern einen DCO-Sign-off.

