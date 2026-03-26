<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

مسار عمل أقلمة مدعوم بالذكاء الاصطناعي لمشاريع البرمجيات. ترجم ملفات i18n وتحقق من صحتها وأدرها باستخدام النماذج اللغوية الكبيرة (LLMs).

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## لماذا Internationalizer؟

معظم أدوات i18n هي إما مكتبات وقت التشغيل (مثل i18next، react-intl) أو منصات برمجيات كخدمة (SaaS) لإدارة المفاتيح (مثل Crowdin، Lokalise). لا تحل أي منها مشكلة الترجمة الفعلية بشكل جيد:

- **الترجمة اليدوية** لا تتوسع بشكل جيد عند تجاوز بضع لغات
- **واجهات برمجة تطبيقات الترجمة الآلية** (Google Translate، DeepL) تتجاهل مصطلحاتك ونبرتك وأعراف واجهة المستخدم الخاصة بك
- **الترجمة العامة باستخدام LLM** تعمل بشكل أفضل، ولكن بدون مسارد وأدلة أسلوب، ستحصل على نتائج غير متسقة

أداة Internationalizer مختلفة. إنها **مسار عمل عبر واجهة سطر الأوامر (CLI)** يجمع بين ترجمة LLM مع:

- **مسارد لكل لغة** — لفرض مصطلحات متسقة عبر تطبيقك
- **أدلة أسلوب لكل لغة** — للتحكم في النبرة، ومستوى الرسمية، وصيغ الجمع، والطباعة
- **ذاكرة الترجمة** — لتخطي السلاسل النصية غير المتغيرة، وتوفير تكاليف استدعاءات واجهة برمجة التطبيقات (API)
- **التحقق من صحة المفاتيح** — لاكتشاف الترجمات المفقودة وعدم تطابق المتغيرات (interpolation) قبل الإصدار

## التثبيت

التثبيت من npm:

```bash
npm install -g internationalizer
```

أو التشغيل بدون تثبيت عام:

```bash
npx internationalizer --help
```

تقوم حزمة npm بتثبيت الملف الثنائي (binary) المبني مسبقًا والمطابق من npm عبر التبعيات الاختيارية الخاصة بنظام التشغيل.

التثبيت باستخدام Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

أو البناء من المصدر:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## حزم npm

- يجب أن تتطابق علامات Git (tags) وإصدارات حزمة npm، على سبيل المثال `v0.1.0` و `0.1.0`
- تعتمد حزمة `internationalizer` الجذرية على حزم أنظمة التشغيل مثل `internationalizer-darwin-arm64`
- أهداف npm المدعومة: macOS arm64/x64، Linux arm64/x64، Windows x64
- يتطلب النشر عبر CI وجود سر (secret) في GitHub باسم `NPM_TOKEN`

## البدء السريع

1. أنشئ ملف تكوين في جذر مشروعك:

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

2. قم بتعيين مفتاح API الخاص بك:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. عاين ما سيتم ترجمته:

```bash
internationalizer translate --dry-run
```

4. قم بتشغيل الترجمة:

```bash
internationalizer translate
```

5. تحقق من صحة جميع اللغات:

```bash
internationalizer validate
```

## الأوامر

### `translate`

ابحث عن المفاتيح المفقودة وترجمها عبر LLM.

```bash
internationalizer translate                    # ترجمة جميع اللغات
internationalizer translate -l fr              # ترجمة الفرنسية فقط
internationalizer translate --dry-run          # معاينة بدون استدعاءات API
internationalizer translate --batch-size 20    # دفعات أصغر
internationalizer translate --concurrency 2    # استدعاءات متوازية أقل
```

### `validate`

تحقق من جميع ملفات اللغات بحثًا عن المفاتيح المفقودة، والمفاتيح الزائدة، وعدم تطابق المتغيرات (interpolation).

```bash
internationalizer validate                     # مخرجات قابلة للقراءة بشريًا
internationalizer validate --json              # مخرجات JSON قابلة للقراءة آليًا
internationalizer validate -q                  # رمز الخروج (exit code) فقط
```

### `detect`

اكتشاف إطار عمل i18n تلقائيًا واقتراح تكوين.

```bash
internationalizer detect
```

يدعم: react-i18next، next-intl، vue-i18n، vanilla JSON، مستندات markdown.

### `glossary`

إدارة مصطلحات المسرد لكل لغة والتي يتم فرضها أثناء الترجمة.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

إدارة ذاكرة الترجمة (ذاكرة تخزين مؤقت JSONL للسلاسل النصية المترجمة مسبقًا).

```bash
internationalizer tm stats                     # عرض عدد السجلات
internationalizer tm export                    # تصدير كـ JSON
internationalizer tm clear --force             # حذف جميع السجلات
```

## مرجع التكوين

```yaml
# .internationalizer.yml

# لغة المصدر (الافتراضي: en)
source_locale: en

# اللغات المراد الترجمة إليها (مطلوب)
target_locales: [fr, de, es, ja, zh-CN, ar]

# مسار ملف لغة المصدر (مطلوب)
source_path: locales/en.json

# إعدادات مزود LLM
llm:
  # المزود: "anthropic" أو "openai" أو "gemini" أو "openrouter" (الافتراضي: gemini)
  provider: gemini

  # الأسماء الافتراضية للنماذج حسب المزود:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # متغير البيئة الذي يحتوي على مفتاح API
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # عنوان URL الأساسي لنقاط النهاية المتوافقة مع OpenAI (اختياري)
  # base_url: https://api.openai.com

# عدد المفاتيح لكل استدعاء LLM (الافتراضي: 40)
batch_size: 40

# استدعاءات LLM المتوازية (الافتراضي: 4)
concurrency: 4

# الدليل الذي يحتوي على ملفات Markdown لأدلة الأسلوب لكل لغة (الافتراضي: style-guides)
style_guides_dir: style-guides

# الدليل الذي يحتوي على ملفات JSON للمسارد لكل لغة (الافتراضي: glossary)
glossary_dir: glossary

# مسار ملف ذاكرة الترجمة (الافتراضي: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## أدلة الأسلوب

أدلة الأسلوب هي ملفات Markdown يتم حقنها في موجه (prompt) ترجمة LLM. إنها تتحكم في النبرة، ومستوى الرسمية، والطباعة، والأعراف الأخرى الخاصة باللغة.

```
style-guides/
  _conventions.md    # قواعد مشتركة لجميع اللغات
  fr.md              # قواعد خاصة باللغة الفرنسية
  ja.md              # قواعد خاصة باللغة اليابانية
  ar.md              # قواعد خاصة باللغة العربية
```

### الأعراف المشتركة (`_conventions.md`)

تحديد القواعد التي تنطبق على جميع اللغات: بناء جملة المتغيرات (interpolation)، والحفاظ على HTML، وأعراف أنواع السلاسل النصية (الأزرار مقابل التسميات مقابل الأخطاء)، إلخ.

### أدلة لكل لغة (`{locale}.md`)

تحديد القواعد الخاصة باللغة: مستوى الرسمية (tu مقابل vous)، وعلامات الترقيم (علامات الاقتباس المزدوجة، وعلامات الاستفهام المقلوبة)، وصيغ الجمع، وتنسيق التواريخ/الأرقام، ومسرد المصطلحات.

راجع [`examples/react-app/style-guides/`](examples/react-app/style-guides/) للحصول على مثال عملي.

## تنسيق المسرد

ملفات المسرد عبارة عن مصفوفات JSON مخزنة في `{glossary_dir}/{locale}.json`:

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

يتم حقن المصطلحات في موجه LLM كجدول مصطلحات، مما يضمن ترجمة متسقة للمصطلحات الأساسية عبر تطبيقك.

## ذاكرة الترجمة

يتم تخزين ذاكرة الترجمة كملف JSONL (سجل JSON واحد لكل سطر). يحتوي كل سجل على:

- مفتاح وقيمة المصدر
- القيمة المترجمة
- تجزئة SHA-256 لقيمة المصدر
- طابع زمني

في عمليات التشغيل اللاحقة، يتم تقديم السلاسل النصية غير المتغيرة من ذاكرة التخزين المؤقت لـ TM دون استدعاء LLM، مما يوفر الوقت وتكاليف API. ملف TM متوافق مع git ويمكن إيداعه (commit) جنبًا إلى جنب مع ملفات اللغات الخاصة بك.

## التنسيقات المدعومة

| التنسيق | الامتدادات | الوضع |
|--------|-----------|------|
| JSON | `.json` | مفتاح-قيمة (متداخل، مسطح بترميز النقطة) |
| YAML | `.yml`, `.yaml` | مفتاح-قيمة (يحافظ على التعليقات والترتيب) |
| Markdown | `.md`, `.mdx` | ترجمة المستند بالكامل |

## اكتشاف نوع المشروع

يحدد `internationalizer detect` إعداد i18n الخاص بك عن طريق التحقق من:

- تبعيات `package.json` لـ react-i18next أو next-intl أو vue-i18n
- هياكل الأدلة التي تطابق أنماط اللغات الشائعة
- امتدادات الملفات وأعراف التسمية

## البنية

```
cmd/internationalizer/     نقطة دخول CLI وتعريفات الأوامر
internal/
  config/                  تحميل تكوين YAML مع الإعدادات الافتراضية
  detect/                  اكتشاف نوع المشروع تلقائيًا
  formats/                 محللات التنسيق (JSON، YAML، Markdown)
  glossary/                إدارة المسرد لكل لغة
  llm/                     واجهة مزود LLM + التطبيقات
    anthropic.go           الواجهة الخلفية لـ Anthropic Claude
    openai.go              الواجهة الخلفية لـ OpenAI / المتوافقة
    gemini.go              الواجهة الخلفية لـ Google Gemini عبر AI Studio
                           يستخدم OpenRouter ملف openai.go مع base_url مخصص
  styleguide/              محمل دليل الأسلوب
  tm/                      ذاكرة ترجمة JSONL
  translate/               منسق الترجمة
  validate/                التحقق من صحة اللغات ومقارنتها
```

## مقارنة مع البدائل

| الميزة | Internationalizer | i18next | Crowdin | LLM عام |
|---------|------------------|---------|---------|-------------|
| ترجمة مدعومة بـ LLM | نعم | لا | جزئي | نعم |
| أدلة أسلوب لكل لغة | نعم | لا | لا | لا |
| فرض المسرد | نعم | لا | نعم | لا |
| ذاكرة الترجمة | نعم | لا | نعم | لا |
| CLI / تنفيذ محلي | نعم | غير متوفر | لا | يدوي |
| ملفات متوافقة مع Git | نعم | نعم | جزئي | يدوي |
| بدون الاعتماد على SaaS | نعم | نعم | لا | يختلف |
| مفتوح المصدر (AGPL-3.0) | نعم | نعم | لا | يختلف |

## الترخيص

[AGPL-3.0](LICENSE)

## المساهمة

راجع [CONTRIBUTING.md](CONTRIBUTING.md) لإعداد التطوير والإرشادات. تتطلب جميع المساهمات توقيع DCO.

