> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

सॉफ़्टवेयर प्रोजेक्ट्स के लिए AI-नेटिव इंटरनेशनलाइज़ेशन पाइपलाइन। LLMs का उपयोग करके i18n फ़ाइलों का अनुवाद, सत्यापन और प्रबंधन करें।

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Internationalizer क्यों?

अधिकांश i18n टूल्स या तो रनटाइम लाइब्रेरीज़ (i18next, react-intl) होते हैं या की-मैनेजमेंट SaaS प्लेटफ़ॉर्म (Crowdin, Lokalise)। इनमें से कोई भी अनुवाद की वास्तविक समस्या का प्रभावी समाधान नहीं करता है:

- **मैन्युअल अनुवाद** कुछ भाषाओं से आगे स्केल नहीं हो पाता
- **मशीन ट्रांसलेशन APIs** (Google Translate, DeepL) आपकी शब्दावली, टोन और UI नियमों की अनदेखी करते हैं
- **सामान्य LLM अनुवाद** बेहतर काम करता है, लेकिन ग्लॉसरी और स्टाइल गाइड के बिना असंगत परिणाम देता है

Internationalizer अलग है। यह एक **CLI पाइपलाइन** है जो LLM अनुवाद को इनके साथ जोड़ती है:

- **प्रति-भाषा ग्लॉसरी** — आपके पूरे ऐप में एक समान शब्दावली लागू करें
- **प्रति-भाषा स्टाइल गाइड्स** — टोन, औपचारिकता, बहुवचन (pluralization) और टाइपोग्राफी को नियंत्रित करें
- **ट्रांसलेशन मेमोरी** — अपरिवर्तित स्ट्रिंग्स को छोड़ें और API कॉल्स के खर्च बचाएँ
- **डिटरमिनिस्टिक वैलिडेशन** — शिप करने से पहले अनुपलब्ध या अतिरिक्त कीज़, सुरक्षित संरचना के विचलन, ग्लॉसरी संबंधी समस्याओं और बहुवचन या ICU त्रुटियों को पकड़ें

<!-- internationalizer:unit markdown:installation -->
## इंस्टॉलेशन

npm से इंस्टॉल करें:

```bash
npm install -g internationalizer
```

या बिना ग्लोबल इंस्टॉल के रन करें:

```bash
npx internationalizer --help
```

npm पैकेज प्लेटफ़ॉर्म-विशिष्ट वैकल्पिक निर्भरताओं (optional dependencies) के ज़रिए npm से संगत प्रीबिल्ट बाइनरी इंस्टॉल करता है।

Go के साथ इंस्टॉल करें:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

या सोर्स से बिल्ड करें:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm पैकेजेस

- Git टैग्स और npm पैकेज वर्ज़न मेल खाने चाहिए, उदाहरण के लिए `v0.1.0` और `0.1.0`
- रूट `internationalizer` पैकेज `internationalizer-darwin-arm64` जैसे प्लेटफ़ॉर्म पैकेजेस पर निर्भर करता है
- समर्थित npm टार्गेट्स: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI पब्लिशिंग के लिए `NPM_TOKEN` नाम के GitHub सीक्रेट की आवश्यकता होती है

<!-- internationalizer:unit markdown:quick-start -->
## क्विक स्टार्ट

1. अपने प्रोजेक्ट रूट में एक कॉन्फ़िगरेशन फ़ाइल बनाएँ:

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

2. अपनी API की (key) सेट करें:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. अनुवाद होने वाली सामग्री का प्रीव्यू देखें:

```bash
internationalizer translate --dry-run
```

4. अनुवाद रन करें:

```bash
internationalizer translate
```

5. सभी लोकेल्स को वैलिडेट करें:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## कमांड्स

### `translate`

अनुपलब्ध या पुरानी (stale) कीज़ ढूँढें और LLM के माध्यम से उनका अनुवाद करें।

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

ट्रांसलेशन स्टेट अनुपलब्ध (missing), स्रोत-पुरानी (source-stale), नीति-पुरानी (policy-stale), वर्तमान (current) और मैन्युअल रूप से संपादित (manually edited) स्थितियों की स्वतंत्र रूप से रिपोर्ट करती है, ताकि कोई मैन्युअल संपादन किसी स्रोत या नीति परिवर्तन को छिपा न सके। नीति-पुरानी वैल्यूज़ रिपोर्ट की जाती हैं, लेकिन केवल `--refresh-policy` के साथ ही उनका पुनः अनुवाद किया जाता है। मैन्युअल रूप से संपादित वैल्यूज़ कभी भी स्वचालित रूप से अधिलेखित (overwrite) नहीं की जाती हैं। समीक्षित अनुवादों में मेनिफ़ेस्ट को शामिल करते समय या किसी समीक्षित मैन्युअल संपादन को स्पष्ट रूप से नई आधार रेखा (baseline) के रूप में स्वीकार करते समय `--adopt-existing` का उपयोग करें।

### `validate`

सभी लोकेल फ़ाइलों की उनके स्रोत बंडलों से जाँच करें। डिफ़ॉल्ट वैलिडेशन संरचनात्मक कवरेज (उपस्थित आवश्यक टार्गेट कीज़ का प्रतिशत) की जाँच करता है, अतिरिक्त कीज़ को चेतावनियों के रूप में रिपोर्ट करता है, और अनुपलब्ध कीज़, इंटरपोलेशन मिसमैच या अमान्य ICU MessageFormat संरचना पर विफल हो जाता है।

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` अनुवादित कवरेज की भी रिपोर्ट करता है। अपने स्रोत के समान भाषाई मान तब तक अनुवादित नहीं माना जाता जब तक कि ग्लॉसरी में संपूर्ण मान के लिए स्पष्ट रूप से समान-स्रोत, समान-टार्गेट प्रविष्टि मौजूद न हो; `ignore_case` का पालन किया जाता है, लेकिन किसी लंबे मान में अंतर्निहित ग्लॉसरी शब्द छूट के दायरे में नहीं आता। स्ट्रिक्ट मोड अतिरिक्त कीज़, स्रोत-समान मानों, परिवर्तित इंटरपोलेशन/HTML/कोड/Markdown-लिंक संरचना, ग्लॉसरी उल्लंघनों और कॉन्फ़िगर किए गए बहुवचन रूपों पर विफल होता है।

`--require-state` प्रत्येक टार्गेट को `.internationalizer.lock` के सापेक्ष सत्यापित करता है। जब कोई की अनट्रैक्ड होती है, या उसका रिकॉर्ड किया गया स्रोत, अनुवाद नीति, या टार्गेट हैश पुराना होता है, तो यह विफल हो जाता है। इसे `--strict` के साथ जोड़ा जा सकता है।

मानव-पठनीय और JSON रिपोर्ट स्थिर कोड्स का उपयोग करती हैं:

| कोड | अर्थ |
| --- | --- |
| `missing_key` / `extra_key` | स्रोत और टार्गेट कीज़ के सेट भिन्न हैं |
| `blank_translation` | गैर-रिक्त स्रोत का स्ट्रिक्ट-मोड टार्गेट रिक्त है |
| `source_identical` | स्ट्रिक्ट-मोड भाषाई मान अनुवादित नहीं हुआ है |
| `protected_structure_mismatch` | इंटरपोलेशन, HTML, कोड, या लिंक संरचना में परिवर्तन हुआ है |
| `glossary_violation` | कोई स्वीकृत टार्गेट शब्द या रूप नहीं मिला |
| `plural_form_missing` | कॉन्फ़िगर किया गया लोकेल बहुवचन रूप अनुपस्थित है |
| `icu_message_syntax` | स्रोत या टार्गेट ICU संदेश विकृत है |
| `icu_argument_mismatch` | ICU आर्गुमेंट के नाम, प्रकार या फ़ॉर्मेटर शैलियाँ भिन्न हैं |
| `icu_selector_mismatch` | चयनकर्ता भिन्न हैं या टार्गेट लोकेल के लिए बहुवचन श्रेणी अमान्य है |
| `untracked` | टार्गेट के लिए कोई मेनिफ़ेस्ट रिकॉर्ड मौजूद नहीं है |
| `source_stale` | रिकॉर्ड किए गए अनुवाद के बाद स्रोत सामग्री बदल गई |
| `policy_stale` | जनरेट किया गया प्रॉम्प्ट या मॉडल सेटिंग्स बदल गईं |
| `target_modified` | टार्गेट सामग्री मेनिफ़ेस्ट रिकॉर्ड से भिन्न है |

### `detect`

i18n फ़्रेमवर्क का स्वतः पता लगाएँ और कॉन्फ़िगरेशन का सुझाव दें।

```bash
internationalizer detect
```

समर्थन: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs।

### `glossary`

प्रति-भाषा ग्लॉसरी शब्दों का प्रबंधन करें जिन्हें अनुवाद के दौरान लागू किया जाता है।

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ट्रांसलेशन मेमोरी (पहले अनुवादित स्ट्रिंग्स का JSONL कैश) प्रबंधित करें।

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## कॉन्फ़िगरेशन संदर्भ

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

लोकेल आइडेंटिफ़ायर सही रूप से गठित BCP 47 टैग होने चाहिए जैसे `fr`, `pt-BR`, या `sr-Latn-RS`। विहित-समतुल्य (canonical-equivalent) टार्गेट लोकेल्स को डुप्लिकेट के रूप में अस्वीकार कर दिया जाता है, और लोकेल-विशिष्ट प्रोवाइडर ओवरराइड्स विहित-समतुल्य वर्तनी से मेल खाते हैं। ऊपर दिए गए उदाहरण में, जापानी सहित बिना ओवरराइड वाले लोकेल्स ग्लोबल Gemini कॉन्फ़िगरेशन इनहेरिट करते हैं।

ICU MessageFormat मानों को संरचनात्मक रूप से पार्स किया जाता है। सामान्य आर्गुमेंट्स, `select`, `plural`, `selectordinal`, `number`, `date`, और `time` समर्थित हैं, जिनमें नेस्टेड संदेश, प्लुरल ऑफ़सेट्स, सटीक-संख्या चयनकर्ता और `#` शामिल हैं। वैलिडेशन सिंटैक्स, आर्गुमेंट प्रकारों और फ़ॉर्मेटर शैलियों, प्लुरल ऑफ़सेट्स, सेलेक्ट शाखा पहचान, और टार्गेट लोकेल की CLDR बहुवचन श्रेणियों की जाँच करता है। इन अपरिवर्तनीयों (invariants) को तोड़ने वाले प्रोवाइडर आउटपुट को लोकेल फ़ाइल या ट्रांसलेशन-मेमोरी रिकॉर्ड लिखे जाने से पहले ही अस्वीकार कर दिया जाता है।

`i18next-v4` के साथ, पहचानी गई स्रोत बहुवचन फ़ैमिलीज़ का अनुवाद के दौरान टार्गेट लोकेल की CLDR श्रेणियों में विस्तार किया जाता है। केवल-टार्गेट श्रेणी स्रोत फ़ैमिली के `_other` मान को अपने अनुवाद टेम्पलेट के रूप में उपयोग करती है। स्ट्रिक्ट वैलिडेशन के लिए उन टार्गेट श्रेणियों की आवश्यकता होती है; केवल-स्रोत श्रेणियाँ उन टार्गेट लोकेल्स के लिए वैकल्पिक हैं जो उनका उपयोग नहीं करते हैं।

<!-- internationalizer:unit markdown:style-guides -->
## स्टाइल गाइड्स

स्टाइल गाइड्स मार्कडाउन फ़ाइलें होती हैं जिन्हें LLM अनुवाद प्रॉम्प्ट में इंजेक्ट किया जाता है। वे टोन, औपचारिकता, टाइपोग्राफी और भाषा-विशिष्ट अन्य नियमों को नियंत्रित करती हैं।

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### साझा परिपाटियाँ (`_conventions.md`)

सभी भाषाओं पर लागू होने वाले नियमों को परिभाषित करें: इंटरपोलेशन सिंटैक्स, HTML संरक्षण, स्ट्रिंग प्रकार के नियम (बटन बनाम लेबल बनाम त्रुटियाँ), आदि।

### प्रति-भाषा गाइड्स (`{locale}.md`)

भाषा-विशिष्ट नियमों को परिभाषित करें: औपचारिकता रजिस्टर (tu बनाम vous), विराम चिह्न (guillemets, उल्टे प्रश्न चिह्न), बहुवचन रूप, दिनांक/संख्या फ़ॉर्मेटिंग और शब्दावली ग्लॉसरी।

स्टाइल गाइड्स स्थायी नीति इनपुट हैं, न कि जनरेट किया गया आउटपुट। Internationalizer उन्हें पढ़ता है लेकिन कभी दोबारा नहीं लिखता। उनकी सामग्री को ग्लॉसरी और प्रॉम्प्ट अनुबंध से अलग हैश किया जाता है, ताकि किसी एप्लिकेशन कोड परिवर्तन से कोई अनुवाद पुराना न हो। किसी गाइड को संपादित करने से जानबूझकर वह लोकेल नीति समीक्षा के लिए चिह्नित हो जाता है; आंतरिक प्रॉम्प्ट शब्दावली बदलने से ऐसा तब तक नहीं होता, जब तक कि प्रॉम्प्ट अनुबंध संस्करण भी न बदले।

कार्यशील उदाहरण के लिए [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) देखें।

<!-- internationalizer:unit markdown:glossary-format -->
## ग्लॉसरी फ़ॉर्मेट

ग्लॉसरी फ़ाइलें `{glossary_dir}/{locale}.json` में संग्रहीत JSON ऐरे होती हैं:

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

`variants` अन्य अनुमोदित टार्गेट रूपों को सूचीबद्ध करता है। `enforcement` का मान `error`, `warning` हो सकता है, या डिफ़ॉल्ट error व्यवहार के लिए इसे छोड़ा जा सकता है। शब्दों को LLM प्रॉम्प्ट में शब्दावली तालिका के रूप में इंजेक्ट किया जाता है, जिससे आपके एप्लिकेशन में सुसंगत अनुवाद सुनिश्चित होता है। `{"source":"API","target":"API"}` जैसी सटीक प्रविष्टि उस संपूर्ण स्रोत-समान मान को स्ट्रिक्ट गैर-अनुवादित निष्कर्षों से छूट भी प्रदान करती है; यह केवल `API` युक्त किसी लंबे मान को छूट नहीं देती।

<!-- internationalizer:unit markdown:translation-memory -->
## ट्रांसलेशन मेमोरी

ट्रांसलेशन मेमोरी JSONL फ़ाइल (प्रति पंक्ति एक JSON रिकॉर्ड) के रूप में संग्रहीत होती है। प्रत्येक रिकॉर्ड में शामिल होते हैं:

- बंडल, की, स्रोत मान, अनुवादित मान, और विहित टार्गेट लोकेल
- स्रोत, स्टाइल-गाइड, ग्लॉसरी, प्रॉम्प्ट-अनुबंध और संयुक्त नीति हैश
- अनुवाद उत्पन्न करने वाला प्रोवाइडर और मॉडल
- एक टाइमस्टैम्प

बाद के रन में, समान स्रोत और नीति हैश वाले स्ट्रिंग्स LLM को कॉल किए बिना कैश से प्रदान किए जाते हैं। डिफ़ॉल्ट पाथ अनदेखी की जाने वाली `.internationalizer/` डायरेक्टरी के अंतर्गत होता है, इसलिए यह एक स्थानीय कैश बना रहता है। यदि आपका प्रोजेक्ट जानबूझकर ट्रांसलेशन मेमोरी साझा करता है, तो `tm_path` को किसी ट्रैक्ड स्थान पर सेट करें। समीक्षा योग्य `.internationalizer.lock` मेनिफ़ेस्ट को अलग से वर्ज़न किया जाता है।

<!-- internationalizer:unit markdown:supported-formats -->
## समर्थित फ़ॉर्मेट

| फ़ॉर्मेट | एक्सटेंशन | मोड |
|--------|-----------|------|
| JSON | `.json` | की-वैल्यू (नेस्टेड, डॉट-नोटेशन फ़्लैटन्ड) |
| YAML | `.yml`, `.yaml` | की-वैल्यू (टिप्पणियों और क्रम को सुरक्षित रखता है) |
| Markdown | `.md`, `.mdx` | प्रस्तावना और H2-स्तरीय अनुभाग |

Markdown टार्गेट्स में H2 अनुभागों से पहले अदृश्य `internationalizer:unit` टिप्पणियाँ होती हैं। ये स्थिर मार्कर Internationalizer को असंबंधित अनुभागों का पुनः अनुवाद किए बिना किसी एक स्रोत अनुभाग को जोड़ने, स्थानांतरित करने या संपादित करने की सुविधा देते हैं। मौजूदा बिना मार्कर वाले दस्तावेज़ों को उनके अगले सफल अपडेट पर मार्कर प्राप्त होते हैं।

<!-- internationalizer:unit markdown:project-type-detection -->
## प्रोजेक्ट प्रकार की पहचान

`internationalizer detect` इनकी जाँच करके आपके i18n सेटअप की पहचान करता है:

- react-i18next, next-intl, या vue-i18n के लिए `package.json` डिपेंडेंसीज़
- सामान्य लोकेल पैटर्न से मेल खाने वाली डायरेक्टरी संरचनाएँ
- फ़ाइल एक्सटेंशन और नामकरण परिपाटियाँ

<!-- internationalizer:unit markdown:architecture -->
## आर्किटेक्चर

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
## विकल्पों से तुलना

| सुविधा | Internationalizer | i18next | Crowdin | सामान्य LLM |
|---------|------------------|---------|---------|-------------|
| LLM-संचालित अनुवाद | हाँ | नहीं | आंशिक | हाँ |
| प्रति-भाषा स्टाइल गाइड्स | हाँ | नहीं | नहीं | नहीं |
| ग्लॉसरी प्रवर्तन | हाँ | नहीं | हाँ | नहीं |
| ट्रांसलेशन मेमोरी | हाँ | नहीं | हाँ | नहीं |
| CLI / स्थानीय निष्पादन | हाँ | लागू नहीं | नहीं | मैन्युअल |
| Git-अनुकूल फ़ाइलें | हाँ | हाँ | आंशिक | मैन्युअल |
| कोई SaaS निर्भरता नहीं | हाँ | हाँ | नहीं | भिन्न होता है |
| ओपन सोर्स (AGPL-3.0) | हाँ | हाँ | नहीं | भिन्न होता है |

<!-- internationalizer:unit markdown:license -->
## लाइसेंस

[AGPL-3.0](../../LICENSE)

डिपेंडेंसी सूचनाओं के लिए [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) देखें।

<!-- internationalizer:unit markdown:contributing -->
## योगदान

डेवलपमेंट सेटअप और दिशानिर्देशों के लिए [CONTRIBUTING.md](../../CONTRIBUTING.md) देखें। सभी योगदानों के लिए DCO साइन-ऑफ़ आवश्यक है।
