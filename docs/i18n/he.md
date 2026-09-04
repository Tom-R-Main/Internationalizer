> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

צינור עבודה (pipeline) מבוסס AI לבינאום פרויקטי תוכנה. תרגום, אימות וניהול קובצי i18n באמצעות מודלי LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## מדוע Internationalizer?

רוב כלי ה-i18n הם ספריות זמן ריצה (כגון i18next או react-intl) או פלטפורמות SaaS לניהול מפתחות (כגון Crowdin או Lokalise). אף אחד מהם אינו פותר היטב את בעיית התרגום עצמה:

- **תרגום ידני** אינו מדרגי מעבר למספר שפות בודדות
- **ממשקי API לתרגום מכונה** (Google Translate, DeepL) מתעלמים מהמונחים, מהטון וממוסכמות ממשק המשתמש (UI)
- **תרגום באמצעות LLM כללי** מספק תוצאות טובות יותר, אך ללא מילוני מונחים ומדריכי סגנון, מתקבלות תוצאות חסרות עקביות

Internationalizer פועל אחרת. זהו **צינור עבודה מבוסס CLI** המשלב תרגום מבוסס LLM עם:

- **מילוני מונחים ייעודיים לכל שפה** — אכיפת טרמינולוגיה עקבית ברחבי היישום
- **מדריכי סגנון ייעודיים לכל שפה** — שליטה בטון, ברמת הרשמיות, בריבוי ובטיפוגרפיה
- **זיכרון תרגום (Translation memory)** — דילוג על מחרוזות שלא השתנו וחיסכון בעלויות של קריאות API
- **אימות דטרמיניסטי** — איתור מפתחות חסרים או עודפים, סטיות במבנים מוגנים, חריגות ממילון המונחים ושגיאות ריבוי או ICU לפני השחרור לגרסת הייצור
<!-- internationalizer:unit markdown:installation -->
## התקנה

התקנה דרך npm:

```bash
npm install -g internationalizer
```

או הרצה ללא התקנה גלובלית:

```bash
npx internationalizer --help
```

חבילת ה-npm מתקינה מ-npm את הקובץ הבינארי התואם שנבנה מראש, באמצעות תלויות רשות התלויות בפלטפורמה.

התקנה באמצעות Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

או הידור מקוד המקור:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## חבילות npm

- תגיות Git וגרסאות חבילת ה-npm חייבות להיות תואמות, לדוגמה `v0.1.0` ו-`0.1.0`
- חבילת השורש `internationalizer` תלויה בחבילות ייעודיות לפלטפורמות השונות, כגון `internationalizer-darwin-arm64`
- יעדי npm נתמכים: macOS arm64/x64, Linux arm64/x64, Windows x64
- פרסום ב-CI דורש סוד ב-GitHub בשם `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## תחילת עבודה מהירה

1. יצירת קובץ הגדרות בשורש הפרויקט:

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

2. הגדרת מפתח ה-API:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. הצגת תצוגה מקדימה של המחרוזות שיתורגמו:

```bash
internationalizer translate --dry-run
```

4. הרצת התרגום:

```bash
internationalizer translate
```

5. אימות כל קובצי השפות:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## פקודות

### `translate`

איתור מפתחות חסרים או לא מעודכנים ותרגומם באמצעות LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

מצב התרגום מדווח בנפרד על מפתחות חסרים, חוסר עדכון במקור (source-stale), חוסר עדכון במדיניות (policy-stale), ערכים עדכניים וערכים שנערכו ידנית. באופן זה, עריכה ידנית אינה יכולה להסוות שינוי שבוצע במקור או במדיניות. ערכים שהמדיניות שלהם אינה מעודכנת מדווחים, אך מתורגמים מחדש רק בעת שימוש בדגל `--refresh-policy`. ערכים שנערכו ידנית לעולם אינם נדרסים באופן אוטומטי. יש להשתמש בדגל `--adopt-existing` בעת הטמעת המניפסט בפרויקט עם תרגומים קיימים שעברו בדיקה, או בעת אישור מפורש של עריכה ידנית כבסיס ההשוואה החדש.

### `validate`

בדיקת כל קובצי השפות מול מאגרי המקור (source bundles). אימות ברירת המחדל בודק כיסוי מבני (אחוז מפתחות היעד הנדרשים שנמצאו), מדווח על מפתחות עודפים כהתראות, ונכשל כאשר נמצאים מפתחות חסרים, חוסר התאמה באינטרפולציה או מבנה שגוי של ICU MessageFormat.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

הדגל `--strict` מדווח גם על כיסוי התרגום. ערך טקסטואלי הזהה למקור נחשב כלא מתורגם, אלא אם מילון המונחים כולל במפורש רשומה זהה מקור-יעד עבור הערך המלא; ההגדרה `ignore_case` נשמרת, אך מונח מילון המשובץ בתוך ערך ארוך יותר אינו מקנה פטור. מצב מחמיר (Strict mode) נכשל בעת הימצאות מפתחות עודפים, ערכים הזהים למקור, שינויים במבנה האינטרפולציה/HTML/קוד/קישורי Markdown, חריגות ממילון המונחים או חוסר בצורות ריבוי שהוגדרו.

הדגל `--require-state` מאמת כל קובץ יעד מול `.internationalizer.lock`. האימות נכשל כאשר מפתח אינו מנוהל (untracked), או כאשר רשומת המקור, מדיניות התרגום או גיבוב ערך היעד אינם מעודכנים. ניתן לשלב דגל זה עם `--strict`.

דוחות הפלט לקריאה אנושית ודוחות ה-JSON משתמשים בקודי ממצאים קבועים:

| קוד | משמעות |
| --- | --- |
| `missing_key` / `extra_key` | קבוצות המפתחות במקור וביעד אינן תואמות |
| `blank_translation` | מקור שאינו ריק מכיל ערך יעד ריק במצב מחמיר |
| `source_identical` | ערך טקסטואלי במצב מחמיר נותר ללא תרגום |
| `protected_structure_mismatch` | מבנה האינטרפולציה, ה-HTML, הקוד או הקישורים השתנה |
| `glossary_violation` | לא נמצא מונח יעד מאושר או וריאנט מאושר |
| `plural_form_missing` | צורת ריבוי שהוגדרה עבור השפה אינה קיימת |
| `icu_message_syntax` | הודעת ICU במקור או ביעד אינה תקינה מבנית |
| `icu_argument_mismatch` | שמות, סוגים או סגנונות מעצבים של ארגומנטי ICU אינם תואמים |
| `icu_selector_mismatch` | בוררים (selectors) אינם תואמים או שקטגוריית הריבוי אינה חוקית בשפת היעד |
| `untracked` | לא קיימת רשומת מניפסט עבור ערך היעד |
| `source_stale` | תוכן המקור השתנה לאחר רישום התרגום |
| `policy_stale` | הפרומפט שנוצר או הגדרות המודל השתנו |
| `target_modified` | תוכן היעד שונה מרשומת המניפסט |

### `detect`

זיהוי אוטומטי של מסגרת ה-i18n בפרויקט והצעת תצורה מתאימה.

```bash
internationalizer detect
```

תמיכה במסגרות: react-i18next, next-intl, vue-i18n, vanilla JSON, מסמכי markdown.

### `glossary`

ניהול מונחי מילון ייעודיים לכל שפה הנאכפים במהלך התרגום.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ניהול זיכרון התרגום (מטמון JSONL של מחרוזות שתורגמו בעבר).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## מדריך הגדרות

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

מזהי שפה חייבים להיות תגיות BCP 47 בעלות מבנה תקין, כגון `fr`, `pt-BR` או `sr-Latn-RS`. שפות יעד בעלות שקילות קנונית נדחות ככפולות, ודריסות של ספקי מודל לפי שפה תואמות לאיות השקול קנונית. בדוגמה שלעיל, שפות ללא דריסה ספציפית — כולל יפנית — יורשות את תצורת Gemini הגלובלית.

ערכי ICU MessageFormat מנותחים באופן מבני. נתמכים ארגומנטים פשוטים, `select`, `plural`, `selectordinal`, `number`, `date` ו-`time`, כולל הודעות מקוננות, היסטים של ריבוי (plural offsets), בוררי מספרים מדויקים ו-`#`. תהליך האימות בודק תחביר, סוגי ארגומנטים וסגנונות מעצבים, היסטים של ריבוי, זהות ענפי בחירה (select branch identity) וקטגוריות ריבוי של CLDR בשפת היעד. פלט מהספק שמפר כללים מבניים אלה נדחה לפני כתיבתו לקובץ שפה או לרשומת זיכרון תרגום.

בעת שימוש ב-`i18next-v4`, משפחות ריבוי מזוהות במקור מורחבות במהלך התרגום לקטגוריות ה-CLDR של שפת היעד. קטגוריה הקיימת ביעד בלבד משתמשת בערך `_other` ממשפחת המקור כתבנית לתרגום. אימות מחמיר מחייב את קיומן של אותן קטגוריות יעד; קטגוריות הקיימות במקור בלבד הן אופציונליות בשפות יעד שאינן משתמשות בהן.
<!-- internationalizer:unit markdown:style-guides -->
## מדריכי סגנון

מדריכי סגנון הם קובצי Markdown המוזרקים לפרומפט התרגום של ה-LLM. הם קובעים את הטון, רמת הרשמיות, הטיפוגרפיה ומוסכמות ייחודיות נוספות של השפה.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### מוסכמות משותפות (`_conventions.md`)

הגדרת כללים החלים על כל השפות: תחביר אינטרפולציה, שימור תגיות HTML, מוסכמות עבור סוגי מחרוזות (לחצנים לעומת תוויות לעומת שגיאות) וכדומה.

### מדריכים לפי שפה (`{locale}.md`)

הגדרת כללים ספציפיים לשפה: משלב רשמיות (פנייה בגוף ראשון או שני), סימני פיסוק (מרכאות, סימני שאלה הפוכים), צורות ריבוי, עיצוב תאריכים ומספרים, ומילון מונחים.

מדריכי סגנון הם קלט מדיניות קבוע, ולא פלט מיוצר. Internationalizer קורא אותם אך לעולם אינו משכתב אותם. התוכן שלהם מגובב (hashed) בנפרד ממילון המונחים ומחוזה הפרומפט, כך ששינוי בקוד היישום אינו הופך את התרגום ללא-מעודכן. עריכה מכוונת של מדריך סגנון מסמנת את אותה שפה לבדיקת מדיניות; שינוי ניסוח פנימי של הפרומפט אינו גורם לכך, אלא אם גרסת חוזה הפרומפט השתנתה אף היא.

דוגמה מעשית זמינה בנתיב [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/).
<!-- internationalizer:unit markdown:glossary-format -->
## מבנה מילון המונחים

קובצי מילון מונחים הם מערכי JSON המאוחסנים בנתיב `{glossary_dir}/{locale}.json`:

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

השדה `variants` מפרט צורות יעד מאושרות נוספות. השדה `enforcement` יכול להכיל את הערכים `error`, `warning`, או להישאר ריק להתנהגות ברירת מחדל מסוג שגיאה. המונחים מוזרקים לפרומפט ה-LLM כטבלת מונחים, מה שמבטיח תרגום עקבי ברחבי היישום שלכם. רשומה מדויקת כגון `{"source":"API","target":"API"}` מעניקה פטור לאותו ערך שלם הזהה למקור מממצאי ערך לא-מתורגם במצב מחמיר; היא אינה מעניקה פטור לערך ארוך יותר שרק מכיל בתוכו את המונח `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## זיכרון תרגום

זיכרון התרגום נשמר כקובץ JSONL (רשומת JSON אחת בכל שורה). כל רשומה כוללת:

- את המאגר (bundle), המפתח, ערך המקור, הערך המתורגם ושפת היעד הקנונית
- גיבובי מקור, מדריך סגנון, מילון מונחים, חוזה פרומפט ומדיניות משולבת
- את הספק והמודל שהפיקו את התרגום
- חותמת זמן

בהרצות עוקבות, מחרוזות בעלות אותם גיבובי מקור ומדיניות מוגשות מתוך המטמון ללא קריאה ל-LLM. נתיב ברירת המחדל נמצא תחת התיקייה המוחרגת `.internationalizer/`, כך שהוא נשמר כמטמון מקומי. ניתן להגדיר את `tm_path` לנתיב המנוהל בגיט אם הפרויקט משתף במכוון את זיכרון התרגום. מניפסט המעקב `.internationalizer.lock` מנוהל בגרסאות בנפרד.
<!-- internationalizer:unit markdown:supported-formats -->
## פורמטים נתמכים

| פורמט | סיומות | מצב |
|---|---|---|
| JSON | `.json` | מפתח-ערך (מקונן, משוטח בייצוג נקודה) |
| YAML | `.yml`, `.yaml` | מפתח-ערך (משמר הערות וסדר) |
| Markdown | `.md`, `.mdx` | פתיח ומקטעים ברמת H2 |

קובצי יעד מסוג Markdown מכילים הערות `internationalizer:unit` נסתרות לפני מקטעי H2. סמנים יציבים אלה מאפשרים ל-Internationalizer להוסיף, להזיז או לערוך מקטע מקור יחיד מבלי לתרגם מחדש מקטעים שאינם קשורים. מסמכים קיימים ללא סמנים יקבלו סמנים אלה בעדכון המוצלח הבא שלהם.
<!-- internationalizer:unit markdown:project-type-detection -->
## זיהוי סוג הפרויקט

הפקודה `internationalizer detect` מזהה את הגדרות ה-i18n שלכם על ידי בדיקת:

- תלויות ב-`package.json` עבור react-i18next, next-intl או vue-i18n
- מבני תיקיות התואמים לדפוסי שפות (locales) נפוצים
- סיומות קבצים ומוסכמות שמות
<!-- internationalizer:unit markdown:architecture -->
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
## השוואה לחלופות

| תכונה | Internationalizer | i18next | Crowdin | LLM כללי |
|---|---|---|---|---|
| תרגום מבוסס LLM | כן | לא | חלקי | כן |
| מדריכי סגנון ייעודיים לכל שפה | כן | לא | לא | לא |
| אכיפת מילון מונחים | כן | לא | כן | לא |
| זיכרון תרגום | כן | לא | כן | לא |
| CLI / הרצה מקומית | כן | לא רלוונטי | לא | ידני |
| קבצים מותאמים ל-Git | כן | כן | חלקי | ידני |
| ללא תלות ב-SaaS | כן | כן | לא | משתנה |
| קוד פתוח (AGPL-3.0) | כן | כן | לא | משתנה |
<!-- internationalizer:unit markdown:license -->
## רישיון

[AGPL-3.0](../../LICENSE)

למידע על תלויות, ראו את [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md).
<!-- internationalizer:unit markdown:contributing -->
## תרומה לפרויקט

להנחיות ולהגדרת סביבת הפיתוח, ראו את [CONTRIBUTING.md](../../CONTRIBUTING.md). כל התרומות דורשות אישור DCO.
