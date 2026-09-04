> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

مسار عمل تدويل أصيل بالذكاء الاصطناعي لمشاريع البرمجيات. ترجمة ملفات i18n والتحقق من صحتها وإدارتها باستخدام LLMs.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## لماذا Internationalizer؟

معظم أدوات i18n إما أن تكون مكتبات لوقت التشغيل (مثل i18next وreact-intl) أو منصات برمجيات كخدمة (SaaS) لإدارة المفاتيح (مثل Crowdin وLokalise). ولا تنجح أي منها في معالجة جوهر مشكلة الترجمة بكفاءة:

- **الترجمة اليدوية** يتعذر توسيع نطاقها عند تجاوز بضع لغات
- **واجهات برمجة تطبيقات (API) الترجمة الآلية** (مثل Google Translate وDeepL) تتجاهل المصطلحات المعتمدة والنبرة وأعراف واجهة المستخدم في تطبيقك
- **الترجمة العامة عبر LLM** تقدم أداءً أفضل، ولكن دون مسارد وأدلة أسلوب ستكون النتائج متباينة وغير متسقة

تتميز أداة Internationalizer بنهج مختلف؛ إذ تمثل **مسار عمل عبر واجهة سطر الأوامر (CLI)** يدمج ترجمة LLM مع العناصر التالية:

- **مسارد مخصصة لكل لغة** — لفرض مصطلحات متسقة في أرجاء التطبيق
- **أدلة أسلوب مخصصة لكل لغة** — للتحكم في النبرة ومستوى الرسمية وصيغ الجمع وعلامات الترقيم
- **ذاكرة ترجمة** — لتخطي السلاسل النصية غير المتغيرة وتوفير تكاليف استدعاءات API
- **تحقق قطعي من الصحة** — لاكتشاف المفاتيح المفقودة أو الزائدة، والانحراف في البنى المحمية، ومخالفات المسرد، وأخطاء صيغ الجمع أو رسائل ICU قبل النشر الفعلي
<!-- internationalizer:unit markdown:installation -->
## التثبيت

التثبيت عبر npm:

```bash
npm install -g internationalizer
```

أو التشغيل المباشر دون تثبيت عام:

```bash
npx internationalizer --help
```

تتولى حزمة npm تثبيت الملف الثنائي الجاهز والمطابق من npm عبر التبعيات الاختيارية الخاصة بكل منصة تشغيل.

التثبيت باستخدام Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

أو البناء من المصدر البرمجي:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## حزم npm

- يجب تطابق وسوم Git وإصدارات حزمة npm، على سبيل المثال `v0.1.0` و`0.1.0`
- تعتمد الحزمة الجذرية `internationalizer` على حزم المنصات المحددة مثل `internationalizer-darwin-arm64`
- منصات npm المدعومة: macOS arm64/x64 وLinux arm64/x64 وWindows x64
- يتطلب النشر عبر مسار CI وجود سر في GitHub باسم `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## البدء السريع

1. إنشاء ملف تكوين في جذر المشروع:

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

2. ضبط مفتاح API في متغيرات البيئة:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. معاينة ما ستتم ترجمته:

```bash
internationalizer translate --dry-run
```

4. تنفيذ الترجمة:

```bash
internationalizer translate
```

5. التحقق من صحة جميع ملفات اللغات:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## الأوامر

### `translate`

استكشاف المفاتيح المفقودة أو القديمة وترجمتها عبر LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

تُبلغ حالة الترجمة بصورة مستقلة عن حالات الفقد، وقدم المصدر، وقدم السياسة، والحالة الراهنة، والتعديل اليدوي؛ ومن ثم لا يمكن لتعديل يدوي إخفاء أي تغيير يطرأ على المصدر أو السياسة. ويجري الإبلاغ عن القيم ذات السياسات القديمة دون إعادة ترجمتها تلقائيًا إلا بإضافة `--refresh-policy`. أما القيم المعدلة يدويًا فلا يتم استبدالها تلقائيًا على الإطلاق. يُستخدم الخيار `--adopt-existing` عند إدخال سجل التتبع إلى ترجمات تمت مراجعتها مسبقًا، أو عند اعتماد تعديل يدوي تمت مراجعته كخط أساس جديد صراحةً.

### `validate`

فحص كافة ملفات اللغات ومقارنتها بحزم المصدر. يتحقق وضع الفحص الافتراضي من التغطية البنيوية (نسبة المفاتيح المطلوبة المتوفرة في ملفات الهدف)، ويُبلغ عن المفاتيح الزائدة كتحذيرات، بينما يفشل الفحص عند وجود مفاتيح مفقودة أو عدم تطابق في بنية المتغيرات (interpolation) أو وجود بنية غير صالحة في صياغة ICU MessageFormat.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

يُبلغ الخيار `--strict` أيضًا عن نسبة تغطية الترجمة. وتُعد القيمة اللغوية المطابقة لنصها المصدري غير مترجمة ما لم يتضمن المسرد صراحةً مدخلاً متطابق المصدر والهدف للقيمة بأكملها؛ مع مراعاة خيار `ignore_case`، علمًا بأن ورود مصطلح المسرد ضمن قيمة أطول لا يمنحها استثناءً. ويفشل الوضع الصارم عند وجود مفاتيح زائدة، أو قيم مطابقة للمصدر، أو حدوث تغيير في بنية المتغيرات أو HTML أو الأكواد أو روابط Markdown، أو انتهاك المسرد، أو غياب صيغ الجمع المحددة في الإعدادات.

يتحقق الخيار `--require-state` من كل هدف مقابل ملف التتبع `.internationalizer.lock`. ويفشل الفحص إذا كان المفتاح غير متتبع، أو إذا كان تجزؤ المصدر أو سياسة الترجمة أو تجزؤ الهدف المسجل قديمًا. ويمكن دمجه مع الخيار `--strict`.

تعتمد التقارير النصية وتقارير JSON رموز نتائج ثابتة:

| الرمز | المعنى |
| --- | --- |
| `missing_key` / `extra_key` | اختلاف بين مجموعتي مفاتيح المصدر والهدف |
| `blank_translation` | نص مصدري غير فارغ يقابله هدف فارغ في الوضع الصارم |
| `source_identical` | بقاء قيمة لغوية غير مترجمة في الوضع الصارم |
| `protected_structure_mismatch` | حدوث تغيير في بنية المتغيرات أو HTML أو التعليمات البرمجية أو الروابط |
| `glossary_violation` | تعذر العثور على المصطلح المعتمد أو أحد متغيراته في لغة الهدف |
| `plural_form_missing` | غياب إحدى صيغ الجمع المكونة للغة الهدف |
| `icu_message_syntax` | بنية غير صحيحة لرسالة ICU في المصدر أو الهدف |
| `icu_argument_mismatch` | اختلاف أسماء معاملات ICU أو أنواعها أو أنماط منسقاتها |
| `icu_selector_mismatch` | اختلاف المحددات أو عدم صلاحية فئة الجمع للغة الهدف |
| `untracked` | عدم وجود سجل في ملف التتبع للهدف |
| `source_stale` | تغير محتوى المصدر بعد عملية الترجمة المسجلة |
| `policy_stale` | تغير الموجه التوليدي (prompt) أو إعدادات النموذج |
| `target_modified` | اختلاف محتوى الهدف عن السجل المحفوظ في ملف التتبع |

### `detect`

اكتشاف إطار عمل i18n تلقائيًا واقتراح تكوين ملائم.

```bash
internationalizer detect
```

يدعم: react-i18next وnext-intl وvue-i18n وملفات JSON القياسية ووثائق markdown.

### `glossary`

إدارة مصطلحات المسرد لكل لغة وفرضها أثناء الترجمة.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

إدارة ذاكرة الترجمة (ذاكرة تخزين مؤقت بصيغة JSONL للسلاسل النصية المترجمة سابقًا).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## مرجع التكوين

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

يجب أن تكون معرفات اللغات وسوم BCP 47 سليمة الصياغة مثل `fr` أو `pt-BR` أو `sr-Latn-RS`. وتُرفض اللغات الهدف المكافئة قياسيًا بصفتها تكرارًا، وتطابق التجاوزات المحددة للغات الهدف التهجئة المكافئة قياسيًا. وفي المثال الموضح أعلاه، ترث اللغات التي لم يُحدد لها تخصيص تجاوُزي — بما فيها اليابانية — تكوين Gemini العام.

يجري تحليل قيم ICU MessageFormat بنيويًا. ويدعم النظام المعاملات البسيطة، و`select`، و`plural`، و`selectordinal`، و`number`، و`date`، و`time`، بما في ذلك الرسائل المتداخلة، وإزاحات الجمع (plural offsets)، ومحددات الأرقام الدقيقة، ورمز `#`. ويتحقق الفحص من الصحة النحوية، وأنواع المعاملات، وأنماط التنسيق، وإزاحات الجمع، وتطابق فروع الاختيار، وفئات جمع CLDR الخاصة باللغة الهدف. وتُرفض مخرجات المزود التي تخل بهذه الثوابت قبل حفظها في ملف اللغة أو في سجل ذاكرة الترجمة.

عند تفعيل `i18next-v4`، يجري توسيع عائلات صيغ الجمع المكتشفة في لغة المصدر أثناء الترجمة إلى فئات CLDR الخاصة باللغة الهدف. وتستخدم الفئة المقتصرة على لغة الهدف قيمة `_other` في عائلة المصدر كقالب للترجمة. ويشترط التحقق الصارم وجود فئات الهدف تلك؛ في حين تُعد الفئات المقتصرة على لغة المصدر اختيارية للغات الهدف التي لا تستخدمها.
<!-- internationalizer:unit markdown:style-guides -->
## أدلة الأسلوب

أدلة الأسلوب هي ملفات Markdown تُحقن ضمن موجه (prompt) ترجمة LLM؛ للتحكم في النبرة، ومستوى الرسمية، وعلامات الترقيم، وسائر أعراف اللغة الخاصة.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### الاتفاقيات المشتركة (`_conventions.md`)

تحدد القواعد المطبقة على جميع اللغات: صياغة المتغيرات، والحفاظ على وسوم HTML، وأعراف أنواع السلاسل النصية (الأزرار مقابل التسميات مقابل رسائل الخطأ)، وما إلى ذلك.

### الأدلة المخصصة لكل لغة (`{locale}.md`)

تحدد القواعد الخاصة بكل لغة: درجة الرسمية (مثل tu مقابل vous)، وعلامات الترقيم (كالأقواس المزدوجة وعلامات الاستفهام المعكوسة)، وصيغ الجمع، وتنسيق التواريخ والأرقام، ومسرد المصطلحات.

تُعد أدلة الأسلوب مدخلات سياسة ثابتة وليست مخرجات مولدة؛ إذ تقرؤها أداة Internationalizer دون تعديل محتواها مطلقًا. ويجري تجزؤ محتواها بمعزل عن المسرد ومواصفات الموجه البرمجي، ومن ثم لا يؤدي تعديل شيفرة التطبيق إلى تقادم الترجمة. يؤدي تعديل الدليل عمدًا إلى تعليم تلك اللغة للمراجعة وفق السياسة؛ بينما لا يؤدي تعديل صياغة الموجه الداخلية إلى ذلك، ما لم يتغير إصدار مواصفات الموجه نفسه.

راجع [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) للاطلاع على مثال عملي مكتمل.
<!-- internationalizer:unit markdown:glossary-format -->
## صيغة المسرد

ملفات المسرد عبارة عن مصفوفات JSON مخزنة في المسار `{glossary_dir}/{locale}.json`:

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

يحدد الحقل `variants` الصيغ المعتمدة الأخرى في لغة الهدف. ويأخذ الحقل `enforcement` إما القيمة `error` أو `warning`، أو يُترك فارغًا لاعتماد سلوك الخطأ الافتراضي. وتُحقن هذه المصطلحات في موجه LLM في هيئة جدول مصطلحات لضمان اتساق الترجمة عبر التطبيق بأكمله. كما يمنح وجود مدخل مطابق تمامًا مثل `{"source":"API","target":"API"}` استثناءً للقيمة المطابقة بالكامل لمصدرها من نتائج التحقق الصارم للقيم غير المترجمة؛ ولا يمنح هذا الاستثناء لقيمة أطول تشتمل على لفظ `API` فحسب.
<!-- internationalizer:unit markdown:translation-memory -->
## ذاكرة الترجمة

تُخزن ذاكرة الترجمة في ملف بتنسيق JSONL (سجل JSON واحد في كل سطر). ويتضمن كل سجل ما يلي:

- الحزمة، والمفتاح، وقيمة المصدر، والقيمة المترجمة، والاسم القياسي للغة الهدف
- تجزؤ المصدر، ودليل الأسلوب، والمسرد، ومواصفات الموجه، وتجزؤ السياسة المجمعة
- المزود والنموذج اللذين أنتجا الترجمة
- طابعًا زمنيًا

في عمليات التشغيل اللاحقة، تُسترجع السلاسل النصية ذات تجزؤ المصدر والسياسة المتطابقين مباشرة من التخزين المؤقت دون استدعاء LLM. ويكون المسار الافتراضي تحت المجلد المتجاهل `.internationalizer/` ليبقى تخزينًا مؤقتًا محليًا. ويمكن ضبط `tm_path` على مسار خاضع للتتبع إذا كان المشروع يتشارك ذاكرة الترجمة عن قصد. علمًا بأن ملف التتبع الخاضع للمراجعة `.internationalizer.lock` تتم إدارته بنسخ مستقلة.
<!-- internationalizer:unit markdown:supported-formats -->
## التنسيقات المدعومة

| التنسيق | الامتدادات | الوضع |
|--------|-----------|------|
| JSON | `.json` | مفتاح-قيمة (متداخل، ومسطح بترميز النقطة) |
| YAML | `.yml`, `.yaml` | مفتاح-قيمة (يحافظ على التعليقات والترتيب) |
| Markdown | `.md`, `.mdx` | المقدمة الاستهلالية والأقسام على مستوى عناوين H2 |

تتضمن ملفات Markdown الهدف تعليقات `internationalizer:unit` غير مرئية قبل عناوين H2. وتتيح هذه العلامات الثابتة لأداة Internationalizer إضافة قسم واحد في المصدر أو نقله أو تعديله دون إعادة ترجمة الأقسام الأخرى التي لم تتغير. وستتلقى المستندات الحالية غير المميزة هذه العلامات عند أول تحديث ناجح لاحق.
<!-- internationalizer:unit markdown:project-type-detection -->
## اكتشاف نوع المشروع

يتعرف أمر `internationalizer detect` على إعدادات i18n في مشروعك من خلال فحص ما يلي:

- تبعيات `package.json` الخاصة بـ react-i18next أو next-intl أو vue-i18n
- هياكل الأدلة المتطابقة مع أنماط اللغات الشائعة
- امتدادات الملفات وأعراف التسمية
<!-- internationalizer:unit markdown:architecture -->
## البنية المعمارية

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
## مقارنة مع البدائل

| الميزة | Internationalizer | i18next | Crowdin | LLM عام |
|---------|------------------|---------|---------|-------------|
| ترجمة مدعومة بـ LLM | نعم | لا | جزئي | نعم |
| أدلة أسلوب مخصصة لكل لغة | نعم | لا | لا | لا |
| فرض المسرد | نعم | لا | نعم | لا |
| ذاكرة الترجمة | نعم | لا | نعم | لا |
| واجهة سطر أوامر / تنفيذ محلي | نعم | لا ينطبق | لا | يدوي |
| ملفات ملائمة لنظام Git | نعم | نعم | جزئي | يدوي |
| خلو من الاعتماد على SaaS | نعم | نعم | لا | يتفاوت |
| مصدر مفتوح (AGPL-3.0) | نعم | نعم | لا | يتفاوت |
<!-- internationalizer:unit markdown:license -->
## الترخيص

[AGPL-3.0](../../LICENSE)

راجع [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) للاطلاع على إشعارات تبعيات الطرف الثالث.
<!-- internationalizer:unit markdown:contributing -->
## المساهمة

راجع [CONTRIBUTING.md](../../CONTRIBUTING.md) لمعرفة إرشادات الإعداد والمساهمة في التطوير. تتطلب كافة المساهمات التوقيع وفق شهادة أصل المطور (DCO).
