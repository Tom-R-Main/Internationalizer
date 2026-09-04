> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

সফটওয়্যার প্রজেক্টের জন্য AI-নেটিভ আন্তর্জাতিকীকরণ (internationalization) পাইপলাইন। LLM ব্যবহার করে i18n ফাইল অনুবাদ করুন, যাচাই (validate) করুন এবং পরিচালনা করুন।

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Internationalizer কেন?

অধিকাংশ i18n টুল হয় রানটাইম লাইব্রেরি (i18next, react-intl) নয়তো কি-ম্যানেজমেন্ট SaaS প্ল্যাটফর্ম (Crowdin, Lokalise)। এগুলোর কোনোটিই অনুবাদের আসল সমস্যার কার্যকর সমাধান করে না:

- **ম্যানুয়াল অনুবাদ** কয়েকটি ভাষার বেশি স্কেল করতে পারে না
- **মেশিন ট্রান্সলেশন API** (Google Translate, DeepL) আপনার পরিভাষা, টোন এবং UI কনভেনশন উপেক্ষা করে
- **সাধারণ LLM অনুবাদ** তুলনামূলক ভালো হলেও গ্লোসারি এবং স্টাইল গাইড ছাড়া অসামঞ্জস্যপূর্ণ ফলাফল দেয়

Internationalizer সম্পূর্ণ আলাদা। এটি একটি **CLI পাইপলাইন**, যা LLM অনুবাদের সাথে নিচের সুবিধাগুলো যুক্ত করে:

- **ভাষাভিত্তিক গ্লোসারি** — আপনার পুরো অ্যাপ্লিকেশনে সুসংগত পরিভাষা বজায় রাখে
- **ভাষাভিত্তিক স্টাইল গাইড** — টোন, ফর্মালিটি, বহুবচন রূপ এবং টাইপোগ্রাফি নিয়ন্ত্রণ করে
- **ট্রান্সলেশন মেমরি** — অপরিবর্তিত স্ট্রিংগুলো এড়িয়ে গিয়ে API কলের খরচ বাঁচায়
- **ডিটারমিনিস্টিক ভ্যালিডেশন** — বাদ পড়া বা অতিরিক্ত কি, সুরক্ষিত কাঠামোর বিচ্যুতি, গ্লোসারি সংক্রান্ত অসঙ্গতি এবং বহুবচন বা ICU ত্রুটি রিলিজের আগেই শনাক্ত করে
<!-- internationalizer:unit markdown:installation -->
## ইনস্টলেশন

npm থেকে ইনস্টল করুন:

```bash
npm install -g internationalizer
```

অথবা গ্লোবাল ইনস্টলেশন ছাড়াই রান করুন:

```bash
npx internationalizer --help
```

npm প্যাকেজটি প্ল্যাটফর্ম-নির্দিষ্ট অপশনাল ডিপেন্ডেন্সির মাধ্যমে npm থেকে উপযুক্ত প্রি-বিল্ট বাইনারি ইনস্টল করে।

Go ব্যবহার করে ইনস্টল করুন:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

অথবা সোর্স থেকে বিল্ড করুন:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## npm প্যাকেজসমূহ

- Git ট্যাগ এবং npm প্যাকেজের সংস্করণ অবশ্যই এক হতে হবে, যেমন `v0.1.0` এবং `0.1.0`
- রুট `internationalizer` প্যাকেজটি প্ল্যাটফর্ম প্যাকেজের ওপর নির্ভরশীল, যেমন `internationalizer-darwin-arm64`
- সমর্থিত npm প্ল্যাটফর্ম: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI প্রকাশের জন্য `NPM_TOKEN` নামের একটি GitHub সিক্রেট প্রয়োজন
<!-- internationalizer:unit markdown:quick-start -->
## কুইক স্টার্ট

১. আপনার প্রজেক্টের রুটে একটি কনফিগারেশন ফাইল তৈরি করুন:

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

২. আপনার API কী সেট করুন:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

৩. কী অনুবাদ হতে চলেছে তার প্রিভিউ দেখুন:

```bash
internationalizer translate --dry-run
```

৪. অনুবাদ প্রক্রিয়াটি রান করুন:

```bash
internationalizer translate
```

৫. সমস্ত লোকেল ভ্যালিডেট করুন:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## কমান্ডসমূহ

### `translate`

অনুপস্থিত বা পুরানো কি (stale keys) খুঁজে বের করুন এবং LLM-এর সাহায্যে অনুবাদ করুন।

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

ট্রান্সলেশন স্টেট স্বতন্ত্রভাবে অনুপস্থিত (missing), সোর্স পরিবর্তনের কারণে পুরানো (source-stale), পলিসি পরিবর্তনের কারণে পুরানো (policy-stale), বর্তমান (current) এবং ম্যানুয়ালি এডিট করা (manually edited) অবস্থা রিপোর্ট করে; ফলে কোনো ম্যানুয়াল এডিট কখনোই সোর্স বা পলিসির পরিবর্তন আড়াল করতে পারে না। পলিসি পরিবর্তনের কারণে পুরানো মানগুলো চিহ্নিত করা হয়, তবে শুধুমাত্র `--refresh-policy` ফ্ল্যাগ দিলে পুনরায় অনুবাদ করা হয়। ম্যানুয়ালি এডিট করা মান কখনোই স্বয়ংক্রিয়ভাবে ওভাররাইট হয় না। পূর্বে পর্যালোচিত কোনো অনুবাদকে ম্যানিফেস্টে যুক্ত করতে অথবা পর্যালোচিত ম্যানুয়াল এডিটকে নতুন বেসলাইন হিসেবে সুস্পষ্টভাবে গ্রহণ করতে `--adopt-existing` ব্যবহার করুন।

### `validate`

সোর্স বান্ডেলের বিপরীতে সমস্ত লোকেল ফাইল যাচাই করুন। ডিফল্ট ভ্যালিডেশন কাঠামোগত কভারেজ (প্রয়োজনীয় টার্গেট কি-র উপস্থিতির শতাংশ) পরীক্ষা করে, অতিরিক্ত কি-গুলোকে ওয়ার্নিং হিসেবে রিপোর্ট করে এবং কোনো কি অনুপস্থিত থাকলে, ইন্টারপোলেশন অমিল হলে বা অবৈধ ICU MessageFormat কাঠামো থাকলে ব্যর্থ (fail) হয়।

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` ফ্ল্যাগ অনুবাদের কভারেজও রিপোর্ট করে। সোর্সের সাথে হুবহু অভিন্ন কোনো ভাষাগত মানকে অনুদিত ধরা হয় না, যদি না গ্লোসারিতে সম্পূর্ণ মানটির জন্য একই সোর্স ও একই টার্গেট সংবলিত সুস্পষ্ট এন্ট্রি থাকে; এতে `ignore_case` মান্য করা হয়, তবে কোনো দীর্ঘ মানের অংশ হিসেবে উপস্থিত গ্লোসারি টার্ম এর আওতাভুক্ত নয়। স্ট্রিক্ট মোড অতিরিক্ত কি, সোর্স-অভিন্ন মান, পরিবর্তিত ইন্টারপোলেশন/HTML/কোড/Markdown লিংকের কাঠামো, গ্লোসারি লঙ্ঘন এবং কনফিগার করা বহুবচন রূপের অনুপস্থিতির ক্ষেত্রে ব্যর্থ হয়।

`--require-state` প্রতিটি টার্গেটকে `.internationalizer.lock`-এর বিপরীতে যাচাই করে। কোনো কি ট্র্যাক না করা থাকলে, অথবা নথিভুক্ত সোর্স, ট্রান্সলেশন পলিসি বা টার্গেট হ্যাশ পুরানো হলে এটি ব্যর্থ হয়। এটি `--strict`-এর সাথেও ব্যবহার করা যায়।

হিউম্যান এবং JSON রিপোর্টে নিচের অপরিবর্তনশীল ফাইন্ডিং কোডগুলো ব্যবহৃত হয়:

| কোড | অর্থ |
| --- | --- |
| `missing_key` / `extra_key` | সোর্স এবং টার্গেটের কি সেটের মধ্যে অমিল রয়েছে |
| `blank_translation` | একটি অ-শূন্য সোর্সের জন্য স্ট্রিক্ট মোডে টার্গেট শূন্য পাওয়া গেছে |
| `source_identical` | স্ট্রিক্ট মোডে একটি ভাষাগত মান অনুবাদহীন রয়ে গেছে |
| `protected_structure_mismatch` | ইন্টারপোলেশন, HTML, কোড বা লিংকের কাঠামো পরিবর্তিত হয়েছে |
| `glossary_violation` | কোনো অনুমোদিত টার্গেট টার্ম বা ভ্যারিয়েন্ট পাওয়া যায়নি |
| `plural_form_missing` | লোকেলের জন্য কনফিগার করা বহুবচন রূপ অনুপস্থিত |
| `icu_message_syntax` | সোর্স বা টার্গেটের ICU মেসেজের গঠন ত্রুটিপূর্ণ |
| `icu_argument_mismatch` | ICU আর্গুমেন্টের নাম, ধরন বা ফরম্যাটার স্টাইল আলাদা |
| `icu_selector_mismatch` | সিলেক্টরে অমিল রয়েছে অথবা টার্গেট লোকেলের জন্য প্লুরাল ক্যাটাগরি অবৈধ |
| `untracked` | টার্গেটের জন্য ম্যানিফেস্টে কোনো রেকর্ড নেই |
| `source_stale` | নথিভুক্ত অনুবাদের পর সোর্সের বিষয়বস্তু পরিবর্তিত হয়েছে |
| `policy_stale` | জেনারেট করা প্রম্পট বা মডেল সেটিংস পরিবর্তিত হয়েছে |
| `target_modified` | টার্গেটের বিষয়বস্তু ম্যানিফেস্ট রেকর্ডের সাথে মিলছে না |

### `detect`

i18n ফ্রেমওয়ার্ক স্বয়ংক্রিয়ভাবে শনাক্ত করুন এবং উপযুক্ত কনফিগারেশনের প্রস্তাব পান।

```bash
internationalizer detect
```

সমর্থিত ফ্রেমওয়ার্ক: react-i18next, next-intl, vue-i18n, ভ্যানিলা JSON, Markdown ডক্স।

### `glossary`

ভাষাভিত্তিক গ্লোসারি টার্ম পরিচালনা করুন, যা অনুবাদের সময় কঠোরভাবে প্রয়োগ করা হয়।

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ট্রান্সলেশন মেমরি (পূর্বে অনুবাদ করা স্ট্রিংয়ের JSONL ক্যাশ) পরিচালনা করুন।

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## কনফিগারেশন রেফারেন্স

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

লোকেল শনাক্তকারীগুলো অবশ্যই BCP 47 ফরম্যাট মেনে গঠিত হতে হবে, যেমন `fr`, `pt-BR`, বা `sr-Latn-RS`। ক্যানোনিকাল-সমতুল্য টার্গেট লোকেলগুলো ডুপ্লিকেট হিসেবে বাতিল করা হয় এবং লোকেল-নির্দিষ্ট প্রোভাইডার ওভাররাইড ক্যানোনিকাল-সমতুল্য বানানের সাথে মেলানো হয়। উপরের উদাহরণে, ওভাররাইড ছাড়া অন্যান্য লোকেল—যেমন জাপানি ভাষা—স্বয়ংক্রিয়ভাবে গ্লোবাল Gemini কনফিগারেশন গ্রহণ করে।

ICU MessageFormat মানগুলো কাঠামোগতভাবে পার্স করা হয়। সাধারণ আর্গুমেন্ট, `select`, `plural`, `selectordinal`, `number`, `date`, এবং `time` সমর্থিত, যার মধ্যে নেস্টেড মেসেজ, প্লুরাল অফসেট, নির্দিষ্ট সংখ্যার সিলেক্টর এবং `#` অন্তর্ভুক্ত। ভ্যালিডেশন প্রক্রিয়ায় সিনট্যাক্স, আর্গুমেন্টের ধরন ও ফরম্যাটার স্টাইল, প্লুরাল অফসেট, সিলেক্ট ব্রাঞ্চের সমতা এবং টার্গেট লোকেলের CLDR প্লুরাল ক্যাটাগরি পরীক্ষা করা হয়। প্রোভাইডারের আউটপুট যদি এই নিয়মগুলো ভঙ্গ করে, তবে লোকেল ফাইল বা ট্রান্সলেশন মেমরিতে রেকর্ড লেখার আগেই তা বাতিল করা হয়।

`i18next-v4`-এর ক্ষেত্রে, অনুবাদের সময় স্বীকৃত সোর্স প্লুরাল ফ্যামিলিগুলোকে টার্গেট লোকেলের CLDR ক্যাটাগরি অনুযায়ী প্রসারিত করা হয়। শুধুমাত্র টার্গেটে থাকা কোনো ক্যাটাগরির ক্ষেত্রে সোর্স ফ্যামিলির `_other` মানকে অনুবাদের টেমপ্লেট হিসেবে ব্যবহার করা হয়। স্ট্রিক্ট ভ্যালিডেশনে এই টার্গেট ক্যাটাগরিগুলো থাকা বাধ্যতামূলক; সোর্সে থাকা কিন্তু টার্গেট লোকেলে অপ্রয়োজনীয় ক্যাটাগরিগুলো ঐচ্ছিক।
<!-- internationalizer:unit markdown:style-guides -->
## স্টাইল গাইড

স্টাইল গাইড হলো Markdown ফাইল, যা অনুবাদের সময় LLM প্রম্পটের মধ্যে যুক্ত করা হয়। এগুলো টোন, ফর্মালিটি, টাইপোগ্রাফি এবং ভাষা-নির্দিষ্ট অন্যান্য রীতিনীতি নিয়ন্ত্রণ করে।

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### শেয়ার্ড কনভেনশন (`_conventions.md`)

সমস্ত ভাষার ক্ষেত্রে প্রযোজ্য নিয়মগুলো সংজ্ঞায়িত করুন: ইন্টারপোলেশন সিনট্যাক্স, HTML অক্ষুণ্ণ রাখা, স্ট্রিংয়ের ধরনের নিয়ম (বাটন বনাম লেবেল বনাম এরর) ইত্যাদি।

### ভাষাভিত্তিক গাইড (`{locale}.md`)

ভাষা-নির্দিষ্ট নিয়মগুলো নির্ধারণ করুন: ফর্মালিটির মাত্রা (tu বনাম vous), বিরামচিহ্ন (guillemets, উল্টানো প্রশ্নবোধক চিহ্ন), বহুবচনের রূপ, তারিখ/সংখ্যার ফরম্যাটিং এবং পরিভাষার গ্লোসারি।

স্টাইল গাইডগুলো স্থায়ী পলিসি ইনপুট, কোনো জেনারেট করা আউটপুট নয়। Internationalizer এগুলো পড়ে কিন্তু কখনোই পরিবর্তন করে না। এদের বিষয়বস্তু গ্লোসারি এবং প্রম্পট চুক্তি থেকে আলাদাভাবে হ্যাশ করা হয়, ফলে অ্যাপ্লিকেশনের কোড পরিবর্তিত হলেও অনুবাদ পুরানো (stale) হয়ে যায় না। একটি গাইড ইচ্ছাকৃতভাবে এডিট করলে সেই লোকেলটি পলিসি পর্যালোচনার জন্য চিহ্নিত হয়; তবে অভ্যন্তরীণ প্রম্পটের ভাষা বদলালে এমনটা ঘটে না, যদি না প্রম্পট চুক্তির সংস্করণও পরিবর্তিত হয়।

কার্যকর উদাহরণের জন্য [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) দেখুন।
<!-- internationalizer:unit markdown:glossary-format -->
## গ্লোসারি ফরম্যাট

গ্লোসারি ফাইলগুলো হলো JSON অ্যারে, যা `{glossary_dir}/{locale}.json`-এ সংরক্ষণ করা হয়:

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

`variants` অনুমোদিত অন্যান্য টার্গেট রূপগুলোর তালিকা রাখে। `enforcement`-এর মান `error` বা `warning` হতে পারে, অথবা ডিফল্ট error আচরণের জন্য এটি বাদ দেওয়া যেতে পারে। পরিভাষার একটি টেবিল হিসেবে টার্মগুলো LLM প্রম্পটে যুক্ত হয়, যা আপনার পুরো অ্যাপ্লিকেশনে সামঞ্জস্যপূর্ণ অনুবাদ নিশ্চিত করে। `{"source":"API","target":"API"}`-এর মতো একটি অবিকল এন্ট্রি সম্পূর্ণ সোর্স-অভিন্ন মানটিকে স্ট্রিক্ট অনুবাদের ফাইন্ডিং থেকে অব্যাহতি দেয়; তবে এটি শুধুমাত্র `API` ধারণকারী কোনো দীর্ঘ বাক্যাংশকে অব্যাহতি দেয় না।
<!-- internationalizer:unit markdown:translation-memory -->
## ট্রান্সলেশন মেমরি

ট্রান্সলেশন মেমরি একটি JSONL ফাইল (প্রতি লাইনে একটি করে JSON রেকর্ড) হিসেবে সংরক্ষিত হয়। প্রতিটি রেকর্ডে থাকে:

- বান্ডেল, কি, সোর্স মান, অনুদিত মান এবং ক্যানোনিকাল টার্গেট লোকেল
- সোর্স, স্টাইল গাইড, গ্লোসারি, প্রম্পট চুক্তি এবং সমন্বিত পলিসি হ্যাশ
- যে প্রোভাইডার এবং মডেল অনুবাদটি তৈরি করেছে
- একটি টাইমস্ট্যাম্প

পরবর্তী রানগুলোতে, একই সোর্স ও পলিসি হ্যাশযুক্ত স্ট্রিংগুলো পুনরায় LLM-কে কল না করে ক্যাশ থেকে পরিবেশিত হয়। ডিফল্ট পাথ হলো গিট দ্বারা উপেক্ষিত `.internationalizer/` ডিরেক্টরি, যাতে এটি লোকাল ক্যাশ হিসেবে থাকে। আপনার প্রজেক্ট যদি দলগতভাবে ট্রান্সলেশন মেমরি শেয়ার করতে চায়, তবে `tm_path`-কে একটি ট্র্যাক করা লোকেশনে সেট করুন। পর্যালোচনাযোগ্য `.internationalizer.lock` ম্যানিফেস্টটি পৃথকভাবে সংস্করণ নিয়ন্ত্রিত হয়।
<!-- internationalizer:unit markdown:supported-formats -->
## সাপোর্টেড ফরম্যাট

| ফরম্যাট | এক্সটেনশন | মোড |
|--------|-----------|------|
| JSON | `.json` | কি-ভ্যালু (নেস্টেড, ডট-নোটেশন দ্বারা ফ্ল্যাটেন করা) |
| YAML | `.yml`, `.yaml` | কি-ভ্যালু (মন্তব্য এবং ক্রম অক্ষুণ্ণ রাখে) |
| Markdown | `.md`, `.mdx` | প্রস্তাবনা (preamble) এবং H2-স্তরের বিভাগসমূহ |

Markdown টার্গেটে H2 বিভাগের পূর্বে অদৃশ্য `internationalizer:unit` কমেন্ট থাকে। এই অপরিবর্তনশীল মার্কারগুলোর সাহায্যে Internationalizer অন্য কোনো অসম্পর্কিত অংশ পুনরায় অনুবাদ না করেই সোর্সের নির্দিষ্ট একটি বিভাগ যোগ, সরানো বা এডিট করতে পারে। বিদ্যমান মার্কারবিহীন ডকুমেন্টগুলোতে পরবর্তী সফল আপডেটের সময় মার্কার যোগ করা হয়।
<!-- internationalizer:unit markdown:project-type-detection -->
## প্রজেক্ট টাইপ ডিটেকশন

`internationalizer detect` নিচের বিষয়গুলো পরীক্ষা করে আপনার i18n সেটআপ শনাক্ত করে:

- react-i18next, next-intl, অথবা vue-i18n-এর জন্য `package.json`-এর ডিপেন্ডেন্সি
- প্রচলিত লোকেল কাঠামোর সাথে মিলে যাওয়া ডিরেক্টরি প্যাটার্ন
- ফাইলের এক্সটেনশন এবং নামকরণের নিয়ম
<!-- internationalizer:unit markdown:architecture -->
## আর্কিটেকচার

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
## বিকল্পগুলোর সাথে তুলনা

| ফিচার | Internationalizer | i18next | Crowdin | সাধারণ LLM |
|---------|------------------|---------|---------|-------------|
| LLM-চালিত অনুবাদ | হ্যাঁ | না | আংশিক | হ্যাঁ |
| ভাষাভিত্তিক স্টাইল গাইড | হ্যাঁ | না | না | না |
| গ্লোসারি প্রয়োগ | হ্যাঁ | না | হ্যাঁ | না |
| ট্রান্সলেশন মেমরি | হ্যাঁ | না | হ্যাঁ | না |
| CLI / লোকাল রান | হ্যাঁ | প্রযোজ্য নয় | না | ম্যানুয়াল |
| Git-বান্ধব ফাইল | হ্যাঁ | হ্যাঁ | আংশিক | ম্যানুয়াল |
| কোনো SaaS নির্ভরতা নেই | হ্যাঁ | হ্যাঁ | না | ভিন্ন হয় |
| ওপেন সোর্স (AGPL-3.0) | হ্যাঁ | হ্যাঁ | না | ভিন্ন হয় |
<!-- internationalizer:unit markdown:license -->
## লাইসেন্স

[AGPL-3.0](../../LICENSE)

ডিপেন্ডেন্সি সংক্রান্ত বিজ্ঞপ্তির জন্য [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) দেখুন।
<!-- internationalizer:unit markdown:contributing -->
## কন্ট্রিবিউটিং

ডেভেলপমেন্ট সেটআপ এবং নির্দেশিকার জন্য [CONTRIBUTING.md](../../CONTRIBUTING.md) দেখুন। সমস্ত অবদানের ক্ষেত্রে DCO সাইন-অফ বাধ্যতামূলক।
