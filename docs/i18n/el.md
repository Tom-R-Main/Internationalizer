> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

AI-native pipeline διεθνοποίησης για έργα λογισμικού. Μεταφράστε, επικυρώστε και διαχειριστείτε αρχεία i18n χρησιμοποιώντας LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Γιατί το Internationalizer;

Τα περισσότερα εργαλεία i18n είναι είτε βιβλιοθήκες χρόνου εκτέλεσης (i18next, react-intl) είτε πλατφόρμες SaaS διαχείρισης κλειδιών (Crowdin, Lokalise). Κανένα από αυτά δεν επιλύει ικανοποιητικά το πραγματικό πρόβλημα της μετάφρασης:

- Η **χειροκίνητη μετάφραση** δεν κλιμακώνεται πέρα από λίγες γλώσσες
- Τα **API μηχανικής μετάφρασης** (Google Translate, DeepL) αγνοούν την ορολογία, το ύφος και τις συμβάσεις του UI σας
- Η **γενική μετάφραση μέσω LLM** λειτουργεί καλύτερα, αλλά χωρίς γλωσσάρια και οδηγούς στυλ, τα αποτελέσματα είναι ασυνεπή

Το Internationalizer διαφέρει. Είναι ένα **CLI pipeline** που συνδυάζει τη μετάφραση μέσω LLM με:

- **Γλωσσάρια ανά γλώσσα** — επιβολή συνεκτικής ορολογίας σε ολόκληρη την εφαρμογή σας
- **Οδηγούς στυλ ανά γλώσσα** — έλεγχος ύφους, βαθμού επισημότητας, σχηματισμού πληθυντικού και τυπογραφίας
- **Μεταφραστική μνήμη** — παράλειψη αμετάβλητων συμβολοσειρών και εξοικονόμηση πόρων σε κλήσεις API
- **Ντετερμινιστική επικύρωση** — εντοπισμός κλειδιών που λείπουν ή πλεονάζουν, απόκλισης προστατευμένης δομής, ζητημάτων γλωσσαρίου και σφαλμάτων πληθυντικού ή ICU πριν από τη διάθεση του κώδικα στην παραγωγή

<!-- internationalizer:unit markdown:installation -->
## Εγκατάσταση

Εγκαταστήστε από το npm:

```bash
npm install -g internationalizer
```

Ή εκτελέστε χωρίς καθολική εγκατάσταση:

```bash
npx internationalizer --help
```

Το πακέτο npm εγκαθιστά το αντίστοιχο προμεταγλωττισμένο εκτελέσιμο από το npm μέσω προαιρετικών εξαρτήσεων για τη συγκεκριμένη πλατφόρμα.

Εγκαταστήστε με το Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Ή δημιουργήστε το εκτελέσιμο από τον πηγαίο κώδικα:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## Πακέτα npm

- Οι ετικέτες Git και οι εκδόσεις των πακέτων npm πρέπει να συμφωνούν απόλυτα, για παράδειγμα `v0.1.0` και `0.1.0`
- Το ριζικό πακέτο `internationalizer` εξαρτάται από πακέτα πλατφόρμας όπως το `internationalizer-darwin-arm64`
- Υποστηριζόμενοι στόχοι npm: macOS arm64/x64, Linux arm64/x64, Windows x64
- Η δημοσίευση μέσω CI απαιτεί ένα μυστικό GitHub με το όνομα `NPM_TOKEN`

<!-- internationalizer:unit markdown:quick-start -->
## Γρήγορη εκκίνηση

1. Δημιουργήστε ένα αρχείο ρυθμίσεων στον ριζικό κατάλογο του έργου σας:

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

2. Ορίστε το κλειδί API σας:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Προεπισκοπήστε τι πρόκειται να μεταφραστεί:

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

<!-- internationalizer:unit markdown:commands -->
## Εντολές

### `translate`

Εντοπίστε κλειδιά που λείπουν ή είναι παρωχημένα και μεταφράστε τα μέσω ενός LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Η κατάσταση μετάφρασης αναφέρει ανεξάρτητα συνθήκες κλειδιών που λείπουν, παρωχημένης πηγής (source-stale), παρωχημένης πολιτικής (policy-stale), ενημερωμένων (current) και μη αυτόματα επεξεργασμένων (manually edited), ώστε μια χειροκίνητη τροποποίηση να μην μπορεί να αποκρύψει μια αλλαγή στην πηγή ή στην πολιτική. Οι τιμές με παρωχημένη πολιτική αναφέρονται, αλλά επαναμεταφράζονται μόνο με τη σημαία `--refresh-policy`. Οι τιμές που έχουν υποστεί χειροκίνητη επεξεργασία δεν αντικαθίστανται ποτέ αυτόματα. Χρησιμοποιήστε τη σημαία `--adopt-existing` κατά την εισαγωγή του manifest σε ελεγμένες μεταφράσεις ή όταν αποδέχεστε ρητά μια ελεγμένη μη αυτόματη τροποποίηση ως τη νέα γραμμή βάσης.

### `validate`

Ελέγξτε όλα τα αρχεία locale σε σχέση με τα πηγαία bundle τους. Η προεπιλεγμένη επικύρωση ελέγχει τη δομική κάλυψη (το ποσοστό των απαιτούμενων κλειδιών προορισμού που υπάρχουν), αναφέρει τυχόν πλεονάζοντα κλειδιά ως προειδοποιήσεις και αποτυγχάνει όταν υπάρχουν κλειδιά που λείπουν, αναντιστοιχίες παρεμβολής ή μη έγκυρη δομή ICU MessageFormat.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

Η σημαία `--strict` αναφέρει επιπλέον τη μεταφρασμένη κάλυψη. Μια γλωσσική τιμή πανομοιότυπη με την πηγή της θεωρείται αμετάφραστη, εκτός εάν το γλωσσάρι περιέχει ρητά μια ακριβή εγγραφή με ίδια πηγή και ίδιο προορισμό για ολόκληρη την τιμή. Η ρύθμιση `ignore_case` τηρείται, αλλά ένας όρος γλωσσαρίου ενσωματωμένος μέσα σε μεγαλύτερη τιμή δεν αποτελεί εξαίρεση. Η αυστηρή λειτουργία (strict mode) αποτυγχάνει σε πλεονάζοντα κλειδιά, τιμές πανομοιότυπες με την πηγή, τροποποιημένη δομή παρεμβολής/HTML/κώδικα/συνδέσμων Markdown, παραβιάσεις γλωσσαρίου και μη έγκυρες ρυθμισμένες μορφές πληθυντικού.

Η σημαία `--require-state` επαληθεύει κάθε αρχείο προορισμού έναντι του αρχείου `.internationalizer.lock`. Αποτυγχάνει όταν ένα κλειδί δεν παρακολουθείται (untracked) ή όταν η καταγεγραμμένη πηγή του, η πολιτική μετάφρασης ή το hash προορισμού είναι παρωχημένα. Μπορεί να συνδυαστεί με το `--strict`.

Οι αναφορές για ανθρώπινη ανάγνωση και σε μορφή JSON χρησιμοποιούν σταθερούς κωδικούς ευρημάτων:

| Κωδικός | Ερμηνεία |
| --- | --- |
| `missing_key` / `extra_key` | Τα σύνολα κλειδιών πηγής και προορισμού διαφέρουν |
| `blank_translation` | Μια μη κενή πηγή έχει κενό προορισμό σε αυστηρή λειτουργία (strict mode) |
| `source_identical` | Μια γλωσσική τιμή σε αυστηρή λειτουργία παραμένει αμετάφραστη |
| `protected_structure_mismatch` | Η δομή παρεμβολής, HTML, κώδικα ή συνδέσμου έχει τροποποιηθεί |
| `glossary_violation` | Δεν βρέθηκε εγκεκριμένος όρος προορισμού ή παραλλαγή |
| `plural_form_missing` | Απουσιάζει μια ρυθμισμένη μορφή πληθυντικού για το locale |
| `icu_message_syntax` | Ένα πηγαίο μήνυμα ή μήνυμα προορισμού ICU είναι κακοσχηματισμένο |
| `icu_argument_mismatch` | Τα ονόματα, οι τύποι ή τα στυλ μορφοποιητή των ορισμάτων ICU διαφέρουν |
| `icu_selector_mismatch` | Οι επιλογείς (selectors) διαφέρουν ή μια κατηγορία πληθυντικού δεν είναι έγκυρη για το locale προορισμού |
| `untracked` | Δεν υπάρχει εγγραφή manifest για τον προορισμό |
| `source_stale` | Το πηγαίο περιεχόμενο τροποποιήθηκε μετά την καταγεγραμμένη μετάφραση |
| `policy_stale` | Το παραγόμενο prompt ή οι ρυθμίσεις μοντέλου άλλαξαν |
| `target_modified` | Το περιεχόμενο προορισμού διαφέρει από την εγγραφή του manifest |

### `detect`

Εντοπίστε αυτόματα το πλαίσιο i18n και λάβετε προτεινόμενες ρυθμίσεις παραμέτρων.

```bash
internationalizer detect
```

Υποστηρίζει: react-i18next, next-intl, vue-i18n, vanilla JSON, έγγραφα markdown.

### `glossary`

Διαχειριστείτε όρους γλωσσαρίου ανά γλώσσα που επιβάλλονται υποχρεωτικά κατά τη μετάφραση.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Διαχειριστείτε τη μεταφραστική μνήμη (προσωρινή μνήμη JSONL προηγουμένως μεταφρασμένων συμβολοσειρών).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## Αναφορά ρυθμίσεων

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

Τα αναγνωριστικά locale πρέπει να είναι καλοσχηματισμένες ετικέτες κατά BCP 47, όπως `fr`, `pt-BR` ή `sr-Latn-RS`. Τα κανονικά ισοδύναμα locale προορισμού απορρίπτονται ως διπλότυπα και οι ειδικές παρακάμψεις παρόχου ανά locale αντιστοιχίζονται με την κανονικά ισοδύναμη γραφή. Στο παραπάνω παράδειγμα, τα locale χωρίς παράκαμψη —συμπεριλαμβανομένων των Ιαπωνικών— κληρονομούν τις καθολικές ρυθμίσεις του Gemini.

Οι τιμές ICU MessageFormat αναλύονται δομικά. Υποστηρίζονται απλά ορίσματα, `select`, `plural`, `selectordinal`, `number`, `date` και `time`, συμπεριλαμβανομένων ένθετων μηνυμάτων, μετατοπίσεων πληθυντικού (plural offsets), επιλογέων ακριβούς αριθμού και του συμβόλου `#`. Η επικύρωση ελέγχει τη σύνταξη, τους τύπους ορισμάτων και τα στυλ μορφοποιητή, τις μετατοπίσεις πληθυντικού, την ταυτότητα των κλάδων select και τις κατηγορίες πληθυντικού CLDR του locale προορισμού. Τα αποτελέσματα του παρόχου που παραβιάζουν αυτές τις αναλλοίωτες συνθήκες απορρίπτονται προτού εγγραφεί ένα αρχείο locale ή μια εγγραφή μεταφραστικής μνήμης.

Με το `i18next-v4`, οι αναγνωρισμένες οικογένειες πληθυντικού της πηγής επεκτείνονται κατά τη μετάφραση στις κατηγορίες CLDR του locale προορισμού. Μια κατηγορία που υπάρχει μόνο στον προορισμό χρησιμοποιεί την τιμή `_other` της οικογένειας πηγής ως πρότυπο μετάφρασης. Η αυστηρή επικύρωση απαιτεί αυτές τις κατηγορίες προορισμού. Οι κατηγορίες που υπάρχουν μόνο στην πηγή είναι προαιρετικές για τα locale προορισμού που δεν τις χρησιμοποιούν.

<!-- internationalizer:unit markdown:style-guides -->
## Οδηγοί στυλ

Οι οδηγοί στυλ είναι αρχεία Markdown που εισάγονται στο prompt μετάφρασης του LLM. Ελέγχουν το ύφος, την επισημότητα, την τυπογραφία και άλλες συμβάσεις ειδικές για κάθε γλώσσα.

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

Ορίστε κανόνες για τη συγκεκριμένη γλώσσα: επίπεδο επισημότητας (ενικός έναντι πληθυντικού ευγενείας), στίξη (εισαγωγικά, ερωτηματικά), μορφές πληθυντικού, μορφοποίηση ημερομηνιών/αριθμών και ένα γλωσσάρι ορολογίας.

Οι οδηγοί στυλ αποτελούν ανθεκτικές εισόδους πολιτικής και όχι παραγόμενα αποτελέσματα. Το Internationalizer τους διαβάζει, αλλά δεν τους επανεγγράφει ποτέ. Το περιεχόμενό τους κατακερματίζεται (hashed) ξεχωριστά από το γλωσσάρι και το συμβόλαιο prompt, ώστε μια αλλαγή στον κώδικα της εφαρμογής να μην καθιστά μια μετάφραση παρωχημένη. Η επεξεργασία ενός οδηγού επισημαίνει σκόπιμα το συγκεκριμένο locale για επανέλεγχο πολιτικής. Αντίθετα, η αλλαγή της εσωτερικής διατύπωσης του prompt δεν προκαλεί επανέλεγχο, εκτός εάν αλλάξει και η έκδοση του συμβολαίου prompt.

Δείτε το [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) για ένα λειτουργικό παράδειγμα.

<!-- internationalizer:unit markdown:glossary-format -->
## Μορφή γλωσσαρίου

Τα αρχεία γλωσσαρίου είναι πίνακες JSON που αποθηκεύονται στη διαδρομή `{glossary_dir}/{locale}.json`:

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

Το πεδίο `variants` παραθέτει άλλες εγκεκριμένες μορφές προορισμού. Το πεδίο `enforcement` μπορεί να οριστεί σε `error`, `warning` ή να παραλειφθεί για την προεπιλεγμένη συμπεριφορά σφάλματος. Οι όροι εισάγονται στο prompt του LLM ως πίνακας ορολογίας, διασφαλίζοντας συνεπή μετάφραση σε ολόκληρη την εφαρμογή σας. Μια ακριβής εγγραφή όπως το `{"source":"API","target":"API"}` εξαιρεί επίσης αυτή την πλήρη, ταυτόσημη με την πηγή τιμή από τα αυστηρά ευρήματα αμετάφραστων τιμών. Δεν εξαιρεί, ωστόσο, μια μεγαλύτερη τιμή που απλώς περιέχει τη λέξη `API`.

<!-- internationalizer:unit markdown:translation-memory -->
## Μεταφραστική μνήμη

Η μεταφραστική μνήμη αποθηκεύεται ως αρχείο JSONL (μία εγγραφή JSON ανά γραμμή). Κάθε εγγραφή περιέχει:

- Το bundle, το κλειδί, την τιμή πηγής, τη μεταφρασμένη τιμή και το κανονικό locale προορισμού
- Τα hash πηγής, οδηγού στυλ, γλωσσαρίου, συμβολαίου prompt και συνδυασμένης πολιτικής
- Τον πάροχο και το μοντέλο που παρήγαγαν τη μετάφραση
- Μια χρονική σήμανση (timestamp)

Σε επόμενες εκτελέσεις, οι συμβολοσειρές με τα ίδια hash πηγής και πολιτικής εξυπηρετούνται απευθείας από την προσωρινή μνήμη χωρίς κλήση στο LLM. Η προεπιλεγμένη διαδρομή βρίσκεται στον αγνοούμενο κατάλογο `.internationalizer/`, ώστε να παραμένει τοπική μνήμη cache. Ορίστε το `tm_path` σε μια παρακολουθούμενη τοποθεσία εάν το έργο σας μοιράζεται σκόπιμα τη μεταφραστική μνήμη. Το αναθεωρήσιμο manifest `.internationalizer.lock` παρακολουθείται σε ξεχωριστή έκδοση.

<!-- internationalizer:unit markdown:supported-formats -->
## Υποστηριζόμενες μορφές

| Μορφή | Επεκτάσεις | Λειτουργία |
|---|---|---|
| JSON | `.json` | Key-value (ένθετο, επίπεδο με dot-notation) |
| YAML | `.yml`, `.yaml` | Key-value (διατηρεί σχόλια και σειρά) |
| Markdown | `.md`, `.mdx` | Προοίμιο και ενότητες επιπέδου H2 |

Οι προορισμοί Markdown περιέχουν αόρατα σχόλια `internationalizer:unit` πριν από τις ενότητες H2. Αυτοί οι σταθεροί δείκτες επιτρέπουν στο Internationalizer να προσθέτει, να μετακινεί ή να επεξεργάζεται μία πηγαία ενότητα χωρίς να επαναμεταφράζει μη σχετικές ενότητες. Τα υπάρχοντα έγγραφα χωρίς σήμανση λαμβάνουν δείκτες κατά την επόμενη επιτυχή ενημέρωσή τους.

<!-- internationalizer:unit markdown:project-type-detection -->
## Εντοπισμός τύπου έργου

Η εντολή `internationalizer detect` αναγνωρίζει τις ρυθμίσεις i18n σας ελέγχοντας:

- Τις εξαρτήσεις του `package.json` για react-i18next, next-intl ή vue-i18n
- Τις δομές καταλόγων που αντιστοιχούν σε κοινά μοτίβα locale
- Τις επεκτάσεις αρχείων και τις συμβάσεις ονομασίας

<!-- internationalizer:unit markdown:architecture -->
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
## Σύγκριση με εναλλακτικές

| Δυνατότητα | Internationalizer | i18next | Crowdin | Γενικό LLM |
|---|---|---|---|---|
| Μετάφραση μέσω LLM | Ναι | Όχι | Μερικώς | Ναι |
| Οδηγοί στυλ ανά γλώσσα | Ναι | Όχι | Όχι | Όχι |
| Επιβολή γλωσσαρίου | Ναι | Όχι | Ναι | Όχι |
| Μεταφραστική μνήμη | Ναι | Όχι | Ναι | Όχι |
| CLI / τοπική εκτέλεση | Ναι | Δ/Υ | Όχι | Χειροκίνητα |
| Αρχεία φιλικά προς το Git | Ναι | Ναι | Μερικώς | Χειροκίνητα |
| Χωρίς εξάρτηση από SaaS | Ναι | Ναι | Όχι | Διαφέρει |
| Ανοιχτού κώδικα (AGPL-3.0) | Ναι | Ναι | Όχι | Διαφέρει |

<!-- internationalizer:unit markdown:license -->
## Άδεια χρήσης

[AGPL-3.0](../../LICENSE)

Δείτε το αρχείο [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) για ειδοποιήσεις σχετικά με εξαρτήσεις τρίτων.

<!-- internationalizer:unit markdown:contributing -->
## Συνεισφορά

Δείτε το αρχείο [CONTRIBUTING.md](../../CONTRIBUTING.md) για τις ρυθμίσεις περιβάλλοντος ανάπτυξης και τις σχετικές οδηγίες. Όλες οι συνεισφορές απαιτούν υπογραφή DCO.
