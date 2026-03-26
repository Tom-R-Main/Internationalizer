<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

सॉफ़्टवेयर प्रोजेक्ट्स के लिए AI-नेटिव इंटरनेशनलाइज़ेशन पाइपलाइन। LLMs का उपयोग करके i18n फ़ाइलों का अनुवाद, सत्यापन और प्रबंधन करें।

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Internationalizer क्यों?

ज़्यादातर i18n टूल्स या तो रनटाइम लाइब्रेरीज़ (i18next, react-intl) हैं या की-मैनेजमेंट SaaS प्लेटफ़ॉर्म (Crowdin, Lokalise) हैं। इनमें से कोई भी अनुवाद की वास्तविक समस्या को अच्छी तरह से हल नहीं करता है:

- **मैन्युअल अनुवाद** कुछ भाषाओं से आगे स्केल नहीं हो पाता है
- **मशीन ट्रांसलेशन APIs** (Google Translate, DeepL) आपकी शब्दावली, टोन और UI नियमों को नज़रअंदाज़ कर देते हैं
- **सामान्य LLM अनुवाद** बेहतर काम करता है, लेकिन ग्लॉसरी और स्टाइल गाइड के बिना, आपको असंगत परिणाम मिलते हैं

Internationalizer अलग है। यह एक **CLI पाइपलाइन** है जो LLM अनुवाद को इनके साथ जोड़ती है:

- **प्रति-भाषा ग्लॉसरी (Per-language glossaries)** — आपके पूरे ऐप में एक समान शब्दावली लागू करती है
- **प्रति-भाषा स्टाइल गाइड** — टोन, औपचारिकता, बहुवचन (pluralization) और टाइपोग्राफी को नियंत्रित करती है
- **ट्रांसलेशन मेमोरी** — बिना बदले गए स्ट्रिंग्स को छोड़ें, API कॉल्स पर पैसे बचाएं
- **की वैलिडेशन (Key validation)** — शिप करने से पहले छूटे हुए अनुवादों और इंटरपोलेशन मिसमैच को पकड़ें

## इंस्टॉलेशन (Installation)

npm से इंस्टॉल करें:

```bash
npm install -g internationalizer
```

या ग्लोबल इंस्टॉल के बिना रन करें:

```bash
npx internationalizer --help
```

npm पैकेज प्लेटफ़ॉर्म-विशिष्ट वैकल्पिक निर्भरताओं (optional dependencies) के माध्यम से npm से मैचिंग प्रीबिल्ट बाइनरी इंस्टॉल करता है।

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

## npm पैकेजेस

- Git टैग्स और npm पैकेज वर्ज़न मैच होने चाहिए, उदाहरण के लिए `v0.1.0` और `0.1.0`
- रूट `internationalizer` पैकेज प्लेटफ़ॉर्म पैकेजेस पर निर्भर करता है जैसे कि `internationalizer-darwin-arm64`
- समर्थित npm टार्गेट्स: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI पब्लिशिंग के लिए `NPM_TOKEN` नाम के GitHub सीक्रेट की आवश्यकता होती है

## क्विक स्टार्ट

1. अपने प्रोजेक्ट रूट में एक कॉन्फ़िग फ़ाइल बनाएँ:

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

2. अपनी API की (key) सेट करें:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. प्रीव्यू करें कि क्या अनुवाद किया जाएगा:

```bash
internationalizer translate --dry-run
```

4. अनुवाद रन करें:

```bash
internationalizer translate
```

5. सभी लोकेल्स (locales) को वैलिडेट करें:

```bash
internationalizer validate
```

## कमांड्स

### `translate`

मिसिंग कीज़ (keys) खोजें और LLM के माध्यम से उनका अनुवाद करें।

```bash
internationalizer translate                    # सभी लोकेल्स का अनुवाद करें
internationalizer translate -l fr              # केवल फ्रेंच का अनुवाद करें
internationalizer translate --dry-run          # API कॉल्स के बिना प्रीव्यू करें
internationalizer translate --batch-size 20    # छोटे बैचेस
internationalizer translate --concurrency 2    # कम पैरेलल कॉल्स
```

### `validate`

मिसिंग कीज़, अतिरिक्त कीज़ और इंटरपोलेशन मिसमैच के लिए सभी लोकेल फ़ाइलों की जाँच करें।

```bash
internationalizer validate                     # ह्यूमन-रीडेबल आउटपुट
internationalizer validate --json              # मशीन-रीडेबल JSON
internationalizer validate -q                  # केवल एग्ज़िट कोड
```

### `detect`

i18n फ़्रेमवर्क का स्वतः पता लगाएँ (Auto-detect) और कॉन्फ़िगरेशन का सुझाव दें।

```bash
internationalizer detect
```

समर्थित: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs।

### `glossary`

प्रति-भाषा ग्लॉसरी शब्दों को प्रबंधित करें जिन्हें अनुवाद के दौरान लागू किया जाता है।

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ट्रांसलेशन मेमोरी (पहले अनुवादित स्ट्रिंग्स का JSONL कैश) प्रबंधित करें।

```bash
internationalizer tm stats                     # रिकॉर्ड काउंट्स दिखाएं
internationalizer tm export                    # JSON के रूप में डंप करें
internationalizer tm clear --force             # सभी रिकॉर्ड्स हटाएं
```

## कॉन्फ़िगरेशन संदर्भ (Configuration Reference)

```yaml
# .internationalizer.yml

# स्रोत भाषा (डिफ़ॉल्ट: en)
source_locale: en

# जिन भाषाओं में अनुवाद करना है (आवश्यक)
target_locales: [fr, de, es, ja, zh-CN, ar]

# स्रोत लोकेल फ़ाइल का पाथ (आवश्यक)
source_path: locales/en.json

# LLM प्रोवाइडर सेटिंग्स
llm:
  # प्रोवाइडर: "anthropic", "openai", "gemini", या "openrouter" (डिफ़ॉल्ट: gemini)
  provider: gemini

  # प्रोवाइडर के अनुसार डिफ़ॉल्ट मॉडल नाम:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # API की (key) वाला एनवायरनमेंट वेरिएबल
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # OpenAI-कम्पैटिबल एंडपॉइंट्स के लिए बेस URL (वैकल्पिक)
  # base_url: https://api.openai.com

# प्रति LLM कॉल कीज़ (डिफ़ॉल्ट: 40)
batch_size: 40

# पैरेलल LLM कॉल्स (डिफ़ॉल्ट: 4)
concurrency: 4

# प्रति-लोकेल स्टाइल गाइड मार्कडाउन फ़ाइलों वाली डायरेक्टरी (डिफ़ॉल्ट: style-guides)
style_guides_dir: style-guides

# प्रति-लोकेल ग्लॉसरी JSON फ़ाइलों वाली डायरेक्टरी (डिफ़ॉल्ट: glossary)
glossary_dir: glossary

# ट्रांसलेशन मेमोरी फ़ाइल का पाथ (डिफ़ॉल्ट: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## स्टाइल गाइड्स

स्टाइल गाइड्स मार्कडाउन फ़ाइलें हैं जिन्हें LLM ट्रांसलेशन प्रॉम्प्ट में इंजेक्ट किया जाता है। वे टोन, औपचारिकता, टाइपोग्राफी और अन्य भाषा-विशिष्ट नियमों को नियंत्रित करते हैं।

```
style-guides/
  _conventions.md    # सभी भाषाओं के लिए साझा नियम
  fr.md              # फ्रेंच-विशिष्ट नियम
  ja.md              # जापानी-विशिष्ट नियम
  ar.md              # अरबी-विशिष्ट नियम
```

### साझा नियम (`_conventions.md`)

उन नियमों को परिभाषित करें जो सभी भाषाओं पर लागू होते हैं: इंटरपोलेशन सिंटैक्स, HTML संरक्षण, स्ट्रिंग प्रकार के नियम (बटन बनाम लेबल बनाम त्रुटियां), आदि।

### प्रति-भाषा गाइड्स (`{locale}.md`)

भाषा-विशिष्ट नियमों को परिभाषित करें: औपचारिकता रजिस्टर (tu बनाम vous), विराम चिह्न (guillemets, उल्टे प्रश्न चिह्न), बहुवचन रूप, दिनांक/संख्या स्वरूपण, और एक शब्दावली ग्लॉसरी।

काम करने वाले उदाहरण के लिए [`examples/react-app/style-guides/`](examples/react-app/style-guides/) देखें。

## ग्लॉसरी फ़ॉर्मेट

ग्लॉसरी फ़ाइलें `{glossary_dir}/{locale}.json` में स्टोर किए गए JSON ऐरे (arrays) हैं:

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

शब्दों को LLM प्रॉम्प्ट में एक शब्दावली तालिका (terminology table) के रूप में इंजेक्ट किया जाता है, जो आपके एप्लिकेशन में प्रमुख शब्दों का एक समान अनुवाद सुनिश्चित करता है।

## ट्रांसलेशन मेमोरी

ट्रांसलेशन मेमोरी को JSONL फ़ाइल (प्रति पंक्ति एक JSON रिकॉर्ड) के रूप में स्टोर किया जाता है। प्रत्येक रिकॉर्ड में शामिल हैं:

- स्रोत की (key) और वैल्यू
- अनुवादित वैल्यू
- स्रोत वैल्यू का SHA-256 हैश
- एक टाइमस्टैम्प

बाद के रन पर, बिना बदले गए स्ट्रिंग्स को LLM को कॉल किए बिना TM कैश से सर्व किया जाता है, जिससे समय और API लागत दोनों की बचत होती है। TM फ़ाइल git-फ़्रेंडली है और इसे आपकी लोकेल फ़ाइलों के साथ कमिट किया जा सकता है।

## समर्थित फ़ॉर्मेट्स

| फ़ॉर्मेट | एक्सटेंशन्स | मोड |
|--------|-----------|------|
| JSON | `.json` | की-वैल्यू (नेस्टेड, डॉट-नोटेशन फ़्लैटन्ड) |
| YAML | `.yml`, `.yaml` | की-वैल्यू (टिप्पणियों और क्रम को सुरक्षित रखता है) |
| Markdown | `.md`, `.mdx` | संपूर्ण-दस्तावेज़ अनुवाद |

## प्रोजेक्ट टाइप डिटेक्शन

`internationalizer detect` इनकी जाँच करके आपके i18n सेटअप की पहचान करता है:

- react-i18next, next-intl, या vue-i18n के लिए `package.json` डिपेंडेंसीज़
- सामान्य लोकेल पैटर्न से मेल खाने वाले डायरेक्टरी स्ट्रक्चर्स
- फ़ाइल एक्सटेंशन और नेमिंग कन्वेंशन्स

## आर्किटेक्चर

```
cmd/internationalizer/     CLI एंट्री पॉइंट और कमांड डेफिनिशन्स
internal/
  config/                  डिफ़ॉल्ट्स के साथ YAML कॉन्फ़िग लोडिंग
  detect/                  प्रोजेक्ट टाइप ऑटो-डिटेक्शन
  formats/                 फ़ॉर्मेट पार्सर्स (JSON, YAML, Markdown)
  glossary/                प्रति-लोकेल ग्लॉसरी मैनेजमेंट
  llm/                     LLM प्रोवाइडर इंटरफ़ेस + इम्प्लीमेंटेशन्स
    anthropic.go           Anthropic Claude बैकएंड
    openai.go              OpenAI / कम्पैटिबल बैकएंड
    gemini.go              AI Studio बैकएंड के माध्यम से Google Gemini
                           OpenRouter कस्टम base_url के साथ openai.go का उपयोग करता है
  styleguide/              स्टाइल गाइड लोडर
  tm/                      JSONL ट्रांसलेशन मेमोरी
  translate/               ट्रांसलेशन ऑर्केस्ट्रेटर
  validate/                लोकेल वैलिडेशन और डिफ़िंग
```

## विकल्पों से तुलना

| फ़ीचर | Internationalizer | i18next | Crowdin | Generic LLM |
|---------|------------------|---------|---------|-------------|
| LLM-पावर्ड अनुवाद | हाँ | नहीं | आंशिक | हाँ |
| प्रति-भाषा स्टाइल गाइड्स | हाँ | नहीं | नहीं | नहीं |
| ग्लॉसरी एन्फ़ोर्समेंट | हाँ | नहीं | हाँ | नहीं |
| ट्रांसलेशन मेमोरी | हाँ | नहीं | हाँ | नहीं |
| CLI / लोकल एक्ज़ीक्यूशन | हाँ | लागू नहीं | नहीं | मैन्युअल |
| Git-फ़्रेंडली फ़ाइलें | हाँ | हाँ | आंशिक | मैन्युअल |
| कोई SaaS डिपेंडेंसी नहीं | हाँ | हाँ | नहीं | भिन्न होता है |
| ओपन सोर्स (AGPL-3.0) | हाँ | हाँ | नहीं | भिन्न होता है |

## लाइसेंस

[AGPL-3.0](LICENSE)

## योगदान (Contributing)

डेवलपमेंट सेटअप और दिशानिर्देशों के लिए [CONTRIBUTING.md](CONTRIBUTING.md) देखें। सभी योगदानों के लिए DCO साइन-ऑफ़ आवश्यक है。

