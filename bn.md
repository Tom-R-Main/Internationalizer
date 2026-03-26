<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

সফটওয়্যার প্রজেক্টের জন্য এআই-নেটিভ (AI-native) ইন্টারন্যাশনালাইজেশন পাইপলাইন। LLM ব্যবহার করে i18n ফাইলগুলো অনুবাদ, যাচাই এবং পরিচালনা করুন।

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Internationalizer কেন?

বেশিরভাগ i18n টুল হয় রানটাইম লাইব্রেরি (i18next, react-intl) অথবা কি-ম্যানেজমেন্ট SaaS প্ল্যাটফর্ম (Crowdin, Lokalise)। এগুলোর কোনোটিই অনুবাদের মূল সমস্যাটি ভালোভাবে সমাধান করতে পারে না:

- **ম্যানুয়াল অনুবাদ** কয়েকটি ভাষার বেশি স্কেল করা যায় না
- **মেশিন ট্রান্সলেশন API** (Google Translate, DeepL) আপনার পরিভাষা, টোন এবং UI কনভেনশনগুলোকে উপেক্ষা করে
- **সাধারণ LLM অনুবাদ** তুলনামূলক ভালো কাজ করে, কিন্তু গ্লোসারি এবং স্টাইল গাইড ছাড়া আপনি অসামঞ্জস্যপূর্ণ ফলাফল পাবেন

Internationalizer আলাদা। এটি একটি **CLI পাইপলাইন** যা LLM অনুবাদের সাথে নিচের বিষয়গুলোকে যুক্ত করে:

- **ভাষাভিত্তিক গ্লোসারি** — আপনার অ্যাপজুড়ে সামঞ্জস্যপূর্ণ পরিভাষা নিশ্চিত করে
- **ভাষাভিত্তিক স্টাইল গাইড** — টোন, ফর্মালিটি, বহুবচন এবং টাইপোগ্রাফি নিয়ন্ত্রণ করে
- **ট্রান্সলেশন মেমরি** — অপরিবর্তিত স্ট্রিংগুলো এড়িয়ে যায়, API কলের খরচ বাঁচায়
- **কি (Key) ভ্যালিডেশন** — রিলিজের আগেই বাদ পড়া অনুবাদ এবং ইন্টারপোলেশনের অমিলগুলো ধরে ফেলে

## ইনস্টলেশন

npm থেকে ইনস্টল করুন:

```bash
npm install -g internationalizer
```

অথবা গ্লোবাল ইনস্টল ছাড়াই রান করুন:

```bash
npx internationalizer --help
```

npm প্যাকেজটি প্ল্যাটফর্ম-নির্দিষ্ট অপশনাল ডিপেন্ডেন্সির মাধ্যমে npm থেকে মানানসই প্রি-বিল্ট বাইনারি ইনস্টল করে।

Go দিয়ে ইনস্টল করুন:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

অথবা সোর্স থেকে বিল্ড করুন:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm প্যাকেজসমূহ

- Git ট্যাগ এবং npm প্যাকেজ ভার্সন অবশ্যই মিলতে হবে, উদাহরণস্বরূপ `v0.1.0` এবং `0.1.0`
- রুট `internationalizer` প্যাকেজটি প্ল্যাটফর্ম প্যাকেজগুলোর ওপর নির্ভর করে, যেমন `internationalizer-darwin-arm64`
- সাপোর্টেড npm টার্গেট: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI পাবলিশিংয়ের জন্য `NPM_TOKEN` নামের একটি GitHub সিক্রেট প্রয়োজন

## কুইক স্টার্ট

১. আপনার প্রজেক্ট রুটে একটি কনফিগ ফাইল তৈরি করুন:

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

২. আপনার API কি (key) সেট করুন:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

৩. কী অনুবাদ করা হবে তার প্রিভিউ দেখুন:

```bash
internationalizer translate --dry-run
```

৪. অনুবাদ শুরু করুন:

```bash
internationalizer translate
```

৫. সব লোকেল (locales) ভ্যালিডেট করুন:

```bash
internationalizer validate
```

## কমান্ডসমূহ

### `translate`

বাদ পড়া কি (keys) খুঁজে বের করুন এবং একটি LLM-এর মাধ্যমে সেগুলো অনুবাদ করুন।

```bash
internationalizer translate                    # সব লোকেল অনুবাদ করুন
internationalizer translate -l fr              # শুধু ফ্রেঞ্চ অনুবাদ করুন
internationalizer translate --dry-run          # API কল ছাড়াই প্রিভিউ দেখুন
internationalizer translate --batch-size 20    # ছোট ব্যাচ
internationalizer translate --concurrency 2    # কম প্যারালাল কল
```

### `validate`

বাদ পড়া কি, অতিরিক্ত কি এবং ইন্টারপোলেশনের অমিল খুঁজতে সব লোকেল ফাইল চেক করুন।

```bash
internationalizer validate                     # মানুষের পড়ার যোগ্য আউটপুট
internationalizer validate --json              # মেশিন-রিডেবল JSON
internationalizer validate -q                  # শুধু এক্সিট কোড
```

### `detect`

i18n ফ্রেমওয়ার্ক স্বয়ংক্রিয়ভাবে শনাক্ত করুন এবং একটি কনফিগারেশনের পরামর্শ দিন।

```bash
internationalizer detect
```

সাপোর্ট করে: react-i18next, next-intl, vue-i18n, ভ্যানিলা JSON, মার্কডাউন ডক্স।

### `glossary`

ভাষাভিত্তিক গ্লোসারি টার্মগুলো পরিচালনা করুন যা অনুবাদের সময় প্রয়োগ করা হয়।

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ট্রান্সলেশন মেমরি পরিচালনা করুন (আগে অনুবাদ করা স্ট্রিংগুলোর JSONL ক্যাশ)।

```bash
internationalizer tm stats                     # রেকর্ডের সংখ্যা দেখান
internationalizer tm export                    # JSON হিসেবে ডাম্প করুন
internationalizer tm clear --force             # সব রেকর্ড মুছে ফেলুন
```

## কনফিগারেশন রেফারেন্স

```yaml
# .internationalizer.yml

# সোর্স ভাষা (ডিফল্ট: en)
source_locale: en

# যেসব ভাষায় অনুবাদ করতে হবে (আবশ্যক)
target_locales: [fr, de, es, ja, zh-CN, ar]

# সোর্স লোকেল ফাইলের পাথ (আবশ্যক)
source_path: locales/en.json

# LLM প্রোভাইডার সেটিংস
llm:
  # প্রোভাইডার: "anthropic", "openai", "gemini", অথবা "openrouter" (ডিফল্ট: gemini)
  provider: gemini

  # প্রোভাইডার অনুযায়ী ডিফল্ট মডেলের নাম:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # API কি ধারণকারী এনভায়রনমেন্ট ভেরিয়েবল
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # OpenAI-সামঞ্জস্যপূর্ণ এন্ডপয়েন্টগুলোর জন্য বেস URL (ঐচ্ছিক)
  # base_url: https://api.openai.com

# প্রতি LLM কলে কি-এর সংখ্যা (ডিফল্ট: 40)
batch_size: 40

# প্যারালাল LLM কল (ডিফল্ট: 4)
concurrency: 4

# ভাষাভিত্তিক স্টাইল গাইড মার্কডাউন ফাইল ধারণকারী ডিরেক্টরি (ডিফল্ট: style-guides)
style_guides_dir: style-guides

# ভাষাভিত্তিক গ্লোসারি JSON ফাইল ধারণকারী ডিরেক্টরি (ডিফল্ট: glossary)
glossary_dir: glossary

# ট্রান্সলেশন মেমরি ফাইলের পাথ (ডিফল্ট: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## স্টাইল গাইড

স্টাইল গাইড হলো মার্কডাউন ফাইল যা LLM ট্রান্সলেশন প্রম্পটে ইনজেক্ট করা হয়। এগুলো টোন, ফর্মালিটি, টাইপোগ্রাফি এবং অন্যান্য ভাষাগত কনভেনশন নিয়ন্ত্রণ করে।

```
style-guides/
  _conventions.md    # সব ভাষার জন্য সাধারণ নিয়ম
  fr.md              # ফ্রেঞ্চ-নির্দিষ্ট নিয়ম
  ja.md              # জাপানি-নির্দিষ্ট নিয়ম
  ar.md              # আরবি-নির্দিষ্ট নিয়ম
```

### সাধারণ কনভেনশন (`_conventions.md`)

সব ভাষার ক্ষেত্রে প্রযোজ্য নিয়মগুলো সংজ্ঞায়িত করুন: ইন্টারপোলেশন সিনট্যাক্স, HTML সংরক্ষণ, স্ট্রিং টাইপ কনভেনশন (বাটন বনাম লেবেল বনাম এরর) ইত্যাদি।

### ভাষাভিত্তিক গাইড (`{locale}.md`)

ভাষা-নির্দিষ্ট নিয়মগুলো সংজ্ঞায়িত করুন: ফর্মালিটি রেজিস্টার (tu বনাম vous), বিরামচিহ্ন (guillemets, উল্টানো প্রশ্নবোধক চিহ্ন), বহুবচন রূপ, তারিখ/সংখ্যার ফরম্যাটিং এবং পরিভাষার গ্লোসারি।

একটি কার্যকরী উদাহরণের জন্য [`examples/react-app/style-guides/`](examples/react-app/style-guides/) দেখুন।

## গ্লোসারি ফরম্যাট

গ্লোসারি ফাইলগুলো হলো JSON অ্যারে যা `{glossary_dir}/{locale}.json`-এ সংরক্ষিত থাকে:

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

টার্মগুলো একটি পরিভাষা টেবিল হিসেবে LLM প্রম্পটে ইনজেক্ট করা হয়, যা আপনার অ্যাপ্লিকেশনের মূল টার্মগুলোর সামঞ্জস্যপূর্ণ অনুবাদ নিশ্চিত করে।

## ট্রান্সলেশন মেমরি

ট্রান্সলেশন মেমরি একটি JSONL ফাইল হিসেবে সংরক্ষিত থাকে (প্রতি লাইনে একটি JSON রেকর্ড)। প্রতিটি রেকর্ডে থাকে:

- সোর্স কি এবং ভ্যালু
- অনুবাদ করা ভ্যালু
- সোর্স ভ্যালুর একটি SHA-256 হ্যাশ
- একটি টাইমস্ট্যাম্প

পরবর্তী রানগুলোতে, অপরিবর্তিত স্ট্রিংগুলো LLM কল না করেই TM ক্যাশ থেকে সার্ভ করা হয়, যা সময় এবং API খরচ উভয়ই বাঁচায়। TM ফাইলটি গিট-বান্ধব (git-friendly) এবং আপনার লোকেল ফাইলগুলোর পাশাপাশি কমিট করা যেতে পারে।

## সাপোর্টেড ফরম্যাট

| ফরম্যাট | এক্সটেনশন | মোড |
|--------|-----------|------|
| JSON | `.json` | কি-ভ্যালু (নেস্টেড, ডট-নোটেশন ফ্ল্যাটেনড) |
| YAML | `.yml`, `.yaml` | কি-ভ্যালু (কমেন্ট এবং ক্রম সংরক্ষণ করে) |
| Markdown | `.md`, `.mdx` | সম্পূর্ণ-ডকুমেন্ট অনুবাদ |

## প্রজেক্ট টাইপ ডিটেকশন

`internationalizer detect` নিচের বিষয়গুলো চেক করে আপনার i18n সেটআপ শনাক্ত করে:

- react-i18next, next-intl, অথবা vue-i18n-এর জন্য `package.json` ডিপেন্ডেন্সি
- সাধারণ লোকেল প্যাটার্নের সাথে মিলে যাওয়া ডিরেক্টরি স্ট্রাকচার
- ফাইল এক্সটেনশন এবং নেমিং কনভেনশন

## আর্কিটেকচার

```
cmd/internationalizer/     CLI এন্ট্রি পয়েন্ট এবং কমান্ড ডেফিনিশন
internal/
  config/                  ডিফল্টসহ YAML কনফিগ লোডিং
  detect/                  প্রজেক্ট টাইপ অটো-ডিটেকশন
  formats/                 ফরম্যাট পার্সার (JSON, YAML, Markdown)
  glossary/                ভাষাভিত্তিক গ্লোসারি ম্যানেজমেন্ট
  llm/                     LLM প্রোভাইডার ইন্টারফেস + ইমপ্লিমেন্টেশন
    anthropic.go           Anthropic Claude ব্যাকএন্ড
    openai.go              OpenAI / সামঞ্জস্যপূর্ণ ব্যাকএন্ড
    gemini.go              AI Studio-এর মাধ্যমে Google Gemini ব্যাকএন্ড
                           OpenRouter কাস্টম base_url-এর সাথে openai.go ব্যবহার করে
  styleguide/              স্টাইল গাইড লোডার
  tm/                      JSONL ট্রান্সলেশন মেমরি
  translate/               ট্রান্সলেশন অর্কেস্ট্রেটর
  validate/                লোকেল ভ্যালিডেশন এবং ডিফিং
```

## বিকল্পগুলোর সাথে তুলনা

| ফিচার | Internationalizer | i18next | Crowdin | সাধারণ LLM |
|---------|------------------|---------|---------|-------------|
| LLM-চালিত অনুবাদ | হ্যাঁ | না | আংশিক | হ্যাঁ |
| ভাষাভিত্তিক স্টাইল গাইড | হ্যাঁ | না | না | না |
| গ্লোসারি প্রয়োগ | হ্যাঁ | না | হ্যাঁ | না |
| ট্রান্সলেশন মেমরি | হ্যাঁ | না | হ্যাঁ | না |
| CLI / লোকাল এক্সিকিউশন | হ্যাঁ | প্রযোজ্য নয় | না | ম্যানুয়াল |
| গিট-বান্ধব ফাইল | হ্যাঁ | হ্যাঁ | আংশিক | ম্যানুয়াল |
| কোনো SaaS নির্ভরতা নেই | হ্যাঁ | হ্যাঁ | না | ভিন্ন হয় |
| ওপেন সোর্স (AGPL-3.0) | হ্যাঁ | হ্যাঁ | না | ভিন্ন হয় |

## লাইসেন্স

[AGPL-3.0](LICENSE)

## কন্ট্রিবিউটিং

ডেভেলপমেন্ট সেটআপ এবং নির্দেশিকার জন্য [CONTRIBUTING.md](CONTRIBUTING.md) দেখুন। সব কন্ট্রিবিউশনের জন্য DCO সাইন-অফ প্রয়োজন।

