> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

תהליך עבודה (pipeline) מבוסס AI לבינאום (internationalization) של פרויקטי תוכנה. תרגום, אימות וניהול של קובצי i18n באמצעות מודלי LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br>
<a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## למה לבחור ב-Internationalizer?

רוב כלי ה-i18n הם ספריות זמן ריצה (runtime) כמו i18next ו-react-intl, או פלטפורמות SaaS לניהול מפתחות (Crowdin, Lokalise). אף אחד מהם לא פותר היטב את בעיית התרגום עצמה:

- **תרגום ידני** אינו סקיילבילי מעבר למספר שפות בודדות
- **ממשקי API לתרגום מכונה** (Google Translate, DeepL) מתעלמים מהמונחים, מהטון וממוסכמות ממשק המשתמש (UI) שלכם
- **תרגום LLM גנרי** עובד טוב יותר, אך ללא מילוני מונחים ומדריכי סגנון, מתקבלות תוצאות חסרות עקביות

Internationalizer הוא שונה. זהו **CLI pipeline** המשלב תרגום LLM עם:

- **מילוני מונחים לכל שפה** — אכיפת טרמינולוגיה עקבית ברחבי האפליקציה
- **מדריכי סגנון לכל שפה** — שליטה בטון, ברמת הרשמיות, בצורות ריבוי ובטיפוגרפיה
- **זיכרון תרגום (Translation memory)** — דילוג על מחרוזות שלא השתנו וחיסכון בעלויות קריאות API
- **אימות מפתחות** — איתור תרגומים חסרים וחוסר התאמה במשתני אינטרפולציה לפני השחרור ללקוחות

## התקנה

התקנה מ-npm:

```bash
npm install -g internationalizer
```

או הרצה ללא התקנה גלובלית:

```bash
npx internationalizer --help
```

חבילת ה-npm מתקינה את הקובץ הבינארי המקומפל מראש התואם למערכת ההפעלה מ-npm דרך תלויות אופציונליות ספציפיות לפלטפורמה.

התקנה באמצעות Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

או בנייה מקוד המקור:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## חבילות npm

- תגיות Git וגרסאות חבילת npm חייבות להיות תואמות, לדוגמה `v0.1.0` ו-`0.1.0`
- חבילת השורש `internationalizer` תלויה בחבילות פלטפורמה כגון `internationalizer-darwin-arm64`
- יעדי npm נתמכים: macOS arm64/x64, Linux arm64/x64, Windows x64
- פרסום דרך CI דורש סוד (secret) ב-GitHub בשם `NPM_TOKEN`

## התחלה מהירה

1. יצירת קובץ הגדרות בשורש הפרויקט:

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

2. הגדרת מפתח ה-API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. תצוגה מקדימה של מה שיתורגם:

```bash
internationalizer translate --dry-run
```

4. הרצת התרגום:

```bash
internationalizer translate
```

5. אימות כל השפות (locales):

```bash
internationalizer validate
```

## פקודות

### `translate`

איתור מפתחות חסרים ותרגומם באמצעות LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

בדיקת כל קובצי השפות לאיתור מפתחות חסרים, מפתחות מיותרים וחוסר התאמה באינטרפולציה.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

זיהוי אוטומטי של תשתית ה-i18n והצעת תצורה (configuration).

```bash
internationalizer detect
```

תומך ב: react-i18next, next-intl, vue-i18n, vanilla JSON, מסמכי markdown.

### `glossary`

ניהול מונחי מילון לכל שפה הנאכפים במהלך התרגום.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ניהול זיכרון תרגום (מטמון JSONL של מחרוזות שתורגמו בעבר).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## מדריך הגדרות (Configuration)

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

## מדריכי סגנון

מדריכי סגנון הם קובצי Markdown המוזרקים לפרומפט התרגום של ה-LLM. הם שולטים בטון, ברמת הרשמיות, בטיפוגרפיה ובמוסכמות נוספות הספציפיות לשפה.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### מוסכמות משותפות (`_conventions.md`)

הגדרת כללים החלים על כל השפות: תחביר אינטרפולציה, שימור HTML, מוסכמות לסוגי מחרוזות (כפתורים לעומת תוויות לעומת שגיאות) וכו'.

### מדריכים לכל שפה (`{locale}.md`)

הגדרת כללים ספציפיים לשפה: משלב רשמיות (tu לעומת vous), פיסוק (מרכאות כפולות, סימני שאלה הפוכים), צורות ריבוי, עיצוב תאריכים/מספרים ומילון מונחים.

ראו את [`examples/react-app/style-guides/`](examples/react-app/style-guides/) לדוגמה עובדת.

## פורמט מילון מונחים

קובצי מילון מונחים הם מערכי JSON המאוחסנים בנתיב `{glossary_dir}/{locale}.json`:

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

המונחים מוזרקים לפרומפט ה-LLM כטבלת טרמינולוגיה, מה שמבטיח תרגום עקבי של מונחי מפתח ברחבי האפליקציה.

## זיכרון תרגום (Translation Memory)

זיכרון התרגום מאוחסן כקובץ JSONL (רשומת JSON אחת בכל שורה). כל רשומה מכילה:

- את מפתח המקור והערך
- את הערך המתורגם
- גיבוב (hash) מסוג SHA-256 של ערך המקור
- חותמת זמן (timestamp)

בהרצות עוקבות, מחרוזות שלא השתנו מוגשות ממטמון ה-TM ללא קריאה ל-LLM, מה שחוסך זמן ועלויות API. קובץ ה-TM ידידותי ל-Git וניתן לדחוף אותו (commit) יחד עם קובצי השפות.

## פורמטים נתמכים

| פורמט | סיומות | מצב | 
|--------|-----------|------|
| JSON | `.json` | מפתח-ערך (מקונן, שטוח בשיטת נקודה) |
| YAML | `.yml`, `.yaml` | מפתח-ערך (משמר הערות וסדר) |
| Markdown | `.md`, `.mdx` | תרגום מסמך שלם |

## זיהוי סוג פרויקט

הפקודה `internationalizer detect` מזהה את תצורת ה-i18n שלכם על ידי בדיקת:

- תלויות ב-`package.json` עבור react-i18next, next-intl או vue-i18n
- מבני ספריות התואמים לתבניות שפה (locale) נפוצות
- סיומות קבצים ומוסכמות שמות

## ארכיטקטורה

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

## השוואה לחלופות

| תכונה | Internationalizer | i18next | Crowdin | LLM גנרי |
|---------|------------------|---------|---------|-------------|
| תרגום מבוסס LLM | כן | לא | חלקי | כן |
| מדריכי סגנון לכל שפה | כן | לא | לא | לא |
| אכיפת מילון מונחים | כן | לא | כן | לא |
| זיכרון תרגום | כן | לא | כן | לא |
| CLI / הרצה מקומית | כן | לא רלוונטי | לא | ידני |
| קבצים ידידותיים ל-Git | כן | כן | חלקי | ידני |
| ללא תלות ב-SaaS | כן | כן | לא | משתנה |
| קוד פתוח (AGPL-3.0) | כן | כן | לא | משתנה |

## רישיון

[AGPL-3.0](LICENSE)

## תרומה לקוד (Contributing)

ראו את [CONTRIBUTING.md](CONTRIBUTING.md) להנחיות והגדרות פיתוח. כל התרומות דורשות אישור DCO.

