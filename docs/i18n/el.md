> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-native ροή εργασιών διεθνοποίησης (internationalization) για έργα λογισμικού. Μεταφράστε, επικυρώστε και διαχειριστείτε αρχεία i18n χρησιμοποιώντας LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br>
<a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Γιατί το Internationalizer;

Τα περισσότερα εργαλεία i18n είναι είτε βιβλιοθήκες χρόνου εκτέλεσης (i18next, react-intl) είτε πλατφόρμες SaaS διαχείρισης κλειδιών (Crowdin, Lokalise). Κανένα από αυτά δεν λύνει καλά το πραγματικό πρόβλημα της μετάφρασης:

- Η **χειροκίνητη μετάφραση** δεν κλιμακώνεται πέρα από μερικές γλώσσες
- Τα **API μηχανικής μετάφρασης** (Google Translate, DeepL) αγνοούν την ορολογία, το ύφος και τις συμβάσεις του περιβάλλοντος χρήστη (UI)
- Η **γενική μετάφραση μέσω LLM** λειτουργεί καλύτερα, αλλά χωρίς γλωσσάρια και οδηγούς στυλ, έχετε ασυνεπή αποτελέσματα

Το Internationalizer είναι διαφορετικό. Είναι ένα **CLI pipeline** που συνδυάζει τη μετάφραση μέσω LLM με:

- **Γλωσσάρια ανά γλώσσα** — επιβάλλουν συνεπή ορολογία σε όλη την εφαρμογή σας
- **Οδηγούς στυλ ανά γλώσσα** — ελέγχουν το ύφος, την επισημότητα, τον πληθυντικό και την τυπογραφία
- **Μεταφραστική μνήμη** — παραλείπει τις αμετάβλητες συμβολοσειρές, εξοικονομώντας χρήματα από κλήσεις API
- **Επικύρωση κλειδιών** — εντοπίζει τις ελλιπείς μεταφράσεις και τις αναντιστοιχίες παρεμβολής (interpolation) πριν την κυκλοφορία

## Εγκατάσταση

Εγκατάσταση από το npm:

```bash
npm install -g internationalizer
```

Ή εκτέλεση χωρίς καθολική εγκατάσταση:

```bash
npx internationalizer --help
```

Το πακέτο npm εγκαθιστά το αντίστοιχο προμεταγλωττισμένο εκτελέσιμο από το npm μέσω προαιρετικών εξαρτήσεων για τη συγκεκριμένη πλατφόρμα.

Εγκατάσταση με Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Ή δημιουργία από τον πηγαίο κώδικα:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Πακέτα npm

- Οι ετικέτες Git και οι εκδόσεις των πακέτων npm πρέπει να ταιριάζουν, για παράδειγμα `v0.1.0` και `0.1.0`
- Το βασικό πακέτο `internationalizer` εξαρτάται από πακέτα πλατφόρμας όπως το `internationalizer-darwin-arm64`
- Υποστηριζόμενοι στόχοι npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Η δημοσίευση μέσω CI απαιτεί ένα GitHub secret με το όνομα `NPM_TOKEN`

## Γρήγορη εκκίνηση

1. Δημιουργήστε ένα αρχείο ρυθμίσεων στον ριζικό κατάλογο του έργου σας:

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

2. Ορίστε το κλειδί API σας:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Δείτε μια προεπισκόπηση του τι θα μεταφραστεί:

```bash
internationalizer translate --dry-run
```

4. Εκτελέστε τη μετάφραση:

```bash
internationalizer translate
```

5. Επικυρώστε όλα τα locale:

```bash
internationalizer validate
```

## Εντολές

### `translate`

Βρείτε τα κλειδιά που λείπουν και μεταφράστε τα μέσω ενός LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Ελέγξτε όλα τα αρχεία locale για κλειδιά που λείπουν, πλεονάζοντα κλειδιά και αναντιστοιχίες παρεμβολής.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Εντοπίστε αυτόματα το πλαίσιο i18n και λάβετε προτάσεις για τις ρυθμίσεις.

```bash
internationalizer detect
```

Υποστηρίζει: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs.

### `glossary`

Διαχειριστείτε τους όρους του γλωσσαρίου ανά γλώσσα που επιβάλλονται κατά τη μετάφραση.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Διαχειριστείτε τη μεταφραστική μνήμη (προσωρινή μνήμη JSONL των ήδη μεταφρασμένων συμβολοσειρών).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Αναφορά ρυθμίσεων

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

## Οδηγοί στυλ

Οι οδηγοί στυλ είναι αρχεία Markdown που εισάγονται στο prompt μετάφρασης του LLM. Ελέγχουν το ύφος, την επισημότητα, την τυπογραφία και άλλες συμβάσεις που αφορούν τη συγκεκριμένη γλώσσα.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Κοινές συμβάσεις (`_conventions.md`)

Ορίστε κανόνες που ισχύουν για όλες τις γλώσσες: σύνταξη παρεμβολής, διατήρηση HTML, συμβάσεις τύπων συμβολοσειρών (κουμπιά έναντι ετικετών έναντι σφαλμάτων) κ.λπ.

### Οδηγοί ανά γλώσσα (`{locale}.md`)

Ορίστε κανόνες για τη συγκεκριμένη γλώσσα: επίπεδο επισημότητας (π.χ. ενικός/πληθυντικός ευγενείας), στίξη (εισαγωγικά, ανεστραμμένα ερωτηματικά), μορφές πληθυντικού, μορφοποίηση ημερομηνίας/αριθμών και ένα γλωσσάρι ορολογίας.

Δείτε το [`examples/react-app/style-guides/`](examples/react-app/style-guides/) για ένα παράδειγμα σε λειτουργία.

## Μορφή γλωσσαρίου

Τα αρχεία γλωσσαρίου είναι πίνακες JSON που αποθηκεύονται στο `{glossary_dir}/{locale}.json`:

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

Οι όροι εισάγονται στο prompt του LLM ως πίνακας ορολογίας, διασφαλίζοντας τη συνεπή μετάφραση βασικών όρων σε όλη την εφαρμογή σας.

## Μεταφραστική μνήμη

Η μεταφραστική μνήμη αποθηκεύεται ως αρχείο JSONL (μία εγγραφή JSON ανά γραμμή). Κάθε εγγραφή περιέχει:

- Το κλειδί και την τιμή προέλευσης
- Τη μεταφρασμένη τιμή
- Ένα κατακερματισμό (hash) SHA-256 της τιμής προέλευσης
- Μια χρονική σήμανση (timestamp)

Σε επόμενες εκτελέσεις, οι αμετάβλητες συμβολοσειρές εξυπηρετούνται από την προσωρινή μνήμη TM χωρίς κλήση στο LLM, εξοικονομώντας χρόνο και κόστος API. Το αρχείο TM είναι φιλικό προς το git και μπορεί να υποβληθεί (commit) μαζί με τα αρχεία locale σας.

## Υποστηριζόμενες μορφές

| Μορφή | Επεκτάσεις | Λειτουργία |
|--------|-----------|------|
| JSON | `.json` | Key-value (ένθετο, επίπεδο με dot-notation) |
| YAML | `.yml`, `.yaml` | Key-value (διατηρεί σχόλια και σειρά) |
| Markdown | `.md`, `.mdx` | Μετάφραση ολόκληρου εγγράφου |

## Εντοπισμός τύπου έργου

Το `internationalizer detect` αναγνωρίζει τις ρυθμίσεις i18n σας ελέγχοντας:

- Τις εξαρτήσεις του `package.json` για react-i18next, next-intl ή vue-i18n
- Τις δομές καταλόγων που ταιριάζουν με κοινά μοτίβα locale
- Τις επεκτάσεις αρχείων και τις συμβάσεις ονομασίας

## Αρχιτεκτονική

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

## Σύγκριση με εναλλακτικές

| Δυνατότητα | Internationalizer | i18next | Crowdin | Γενικό LLM |
|---------|------------------|---------|---------|-------------|
| Μετάφραση μέσω LLM | Ναι | Όχι | Μερικώς | Ναι |
| Οδηγοί στυλ ανά γλώσσα | Ναι | Όχι | Όχι | Όχι |
| Επιβολή γλωσσαρίου | Ναι | Όχι | Ναι | Όχι |
| Μεταφραστική μνήμη | Ναι | Όχι | Ναι | Όχι |
| CLI / τοπική εκτέλεση | Ναι | Δ/Υ | Όχι | Χειροκίνητα |
| Αρχεία φιλικά προς το Git | Ναι | Ναι | Μερικώς | Χειροκίνητα |
| Χωρίς εξάρτηση από SaaS | Ναι | Ναι | Όχι | Διαφέρει |
| Ανοιχτού κώδικα (AGPL-3.0) | Ναι | Ναι | Όχι | Διαφέρει |

## Άδεια χρήσης

[AGPL-3.0](LICENSE)

## Συνεισφορά

Δείτε το [CONTRIBUTING.md](CONTRIBUTING.md) για τις ρυθμίσεις ανάπτυξης και τις οδηγίες. Όλες οι συνεισφορές απαιτούν υπογραφή DCO.

