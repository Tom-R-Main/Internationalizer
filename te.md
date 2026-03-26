<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

సాఫ్ట్‌వేర్ ప్రాజెక్ట్‌ల కోసం AI-నేటివ్ ఇంటర్నేషనలైజేషన్ పైప్‌లైన్. LLMలను ఉపయోగించి i18n ఫైల్‌లను అనువదించండి, ధృవీకరించండి మరియు నిర్వహించండి.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Internationalizer ఎందుకు?

చాలా i18n టూల్స్ రన్‌టైమ్ లైబ్రరీలు (i18next, react-intl) లేదా కీ-మేనేజ్‌మెంట్ SaaS ప్లాట్‌ఫారమ్‌లు (Crowdin, Lokalise). వాటిలో ఏవీ అసలు అనువాద సమస్యను సరిగ్గా పరిష్కరించలేవు:

- **మాన్యువల్ అనువాదం** (Manual translation) కొన్ని భాషల కంటే ఎక్కువ స్కేల్ చేయబడదు
- **మెషిన్ ట్రాన్స్‌లేషన్ APIలు** (Google Translate, DeepL) మీ పరిభాష, టోన్ మరియు UI కన్వెన్షన్‌లను విస్మరిస్తాయి
- **సాధారణ LLM అనువాదం** మెరుగ్గా పనిచేస్తుంది, కానీ గ్లాసరీలు మరియు స్టైల్ గైడ్‌లు లేకుండా, మీకు అస్థిరమైన ఫలితాలు వస్తాయి

Internationalizer భిన్నమైనది. ఇది LLM అనువాదాన్ని వీటితో కలిపే ఒక **CLI పైప్‌లైన్**:

- **భాషల వారీగా గ్లాసరీలు** — మీ యాప్ అంతటా స్థిరమైన పరిభాషను అమలు చేస్తుంది
- **భాషల వారీగా స్టైల్ గైడ్‌లు** — టోన్, ఫార్మాలిటీ, బహువచనాలు మరియు టైపోగ్రఫీని నియంత్రిస్తుంది
- **ట్రాన్స్‌లేషన్ మెమరీ** — మార్చబడని స్ట్రింగ్‌లను దాటవేస్తుంది, API కాల్‌లలో డబ్బును ఆదా చేస్తుంది
- **కీ వాలిడేషన్** — షిప్ చేయడానికి ముందే తప్పిపోయిన అనువాదాలు మరియు ఇంటర్‌పోలేషన్ మిస్‌మ్యాచ్‌లను పట్టుకుంటుంది

## ఇన్‌స్టాలేషన్

npm నుండి ఇన్‌స్టాల్ చేయండి:

```bash
npm install -g internationalizer
```

లేదా గ్లోబల్ ఇన్‌స్టాల్ లేకుండా రన్ చేయండి:

```bash
npx internationalizer --help
```

npm ప్యాకేజీ ప్లాట్‌ఫారమ్-నిర్దిష్ట ఆప్షనల్ డిపెండెన్సీల ద్వారా npm నుండి సరిపోలే ప్రీబిల్ట్ బైనరీని ఇన్‌స్టాల్ చేస్తుంది.

Go తో ఇన్‌స్టాల్ చేయండి:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

లేదా సోర్స్ నుండి బిల్డ్ చేయండి:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm ప్యాకేజీలు

- Git ట్యాగ్‌లు మరియు npm ప్యాకేజీ వెర్షన్‌లు తప్పనిసరిగా సరిపోలాలి, ఉదాహరణకు `v0.1.0` మరియు `0.1.0`
- రూట్ `internationalizer` ప్యాకేజీ `internationalizer-darwin-arm64` వంటి ప్లాట్‌ఫారమ్ ప్యాకేజీలపై ఆధారపడి ఉంటుంది
- సపోర్ట్ చేసే npm టార్గెట్‌లు: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI పబ్లిషింగ్ కోసం `NPM_TOKEN` పేరుతో ఒక GitHub సీక్రెట్ అవసరం

## క్విక్ స్టార్ట్

1. మీ ప్రాజెక్ట్ రూట్‌లో కాన్ఫిగ్ ఫైల్‌ను సృష్టించండి:

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

2. మీ API కీని సెట్ చేయండి:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. ఏవి అనువదించబడతాయో ప్రివ్యూ చూడండి:

```bash
internationalizer translate --dry-run
```

4. అనువాదాన్ని రన్ చేయండి:

```bash
internationalizer translate
```

5. అన్ని లొకేల్‌లను వాలిడేట్ చేయండి:

```bash
internationalizer validate
```

## కమాండ్‌లు

### `translate`

తప్పిపోయిన కీలను కనుగొని, వాటిని LLM ద్వారా అనువదించండి.

```bash
internationalizer translate                    # అన్ని లొకేల్‌లను అనువదించండి
internationalizer translate -l fr              # ఫ్రెంచ్ మాత్రమే అనువదించండి
internationalizer translate --dry-run          # API కాల్స్ లేకుండా ప్రివ్యూ చూడండి
internationalizer translate --batch-size 20    # చిన్న బ్యాచ్‌లు
internationalizer translate --concurrency 2    # తక్కువ ప్యారలల్ కాల్స్
```

### `validate`

తప్పిపోయిన కీలు, అదనపు కీలు మరియు ఇంటర్‌పోలేషన్ మిస్‌మ్యాచ్‌ల కోసం అన్ని లొకేల్ ఫైల్‌లను తనిఖీ చేయండి.

```bash
internationalizer validate                     # చదవగలిగే అవుట్‌పుట్
internationalizer validate --json              # మెషిన్ చదవగలిగే JSON
internationalizer validate -q                  # ఎగ్జిట్ కోడ్ మాత్రమే
```

### `detect`

i18n ఫ్రేమ్‌వర్క్‌ను ఆటో-డిటెక్ట్ చేసి, కాన్ఫిగరేషన్‌ను సూచించండి.

```bash
internationalizer detect
```

సపోర్ట్ చేసేవి: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs.

### `glossary`

అనువాద సమయంలో అమలు చేయబడే భాషల వారీగా గ్లాసరీ పదాలను నిర్వహించండి.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ట్రాన్స్‌లేషన్ మెమరీని నిర్వహించండి (గతంలో అనువదించబడిన స్ట్రింగ్‌ల JSONL కాష్).

```bash
internationalizer tm stats                     # రికార్డ్ కౌంట్‌లను చూపించు
internationalizer tm export                    # JSON గా డంప్ చేయి
internationalizer tm clear --force             # అన్ని రికార్డులను తొలగించు
```

## కాన్ఫిగరేషన్ రిఫరెన్స్

```yaml
# .internationalizer.yml

# సోర్స్ భాష (డిఫాల్ట్: en)
source_locale: en

# అనువదించాల్సిన భాషలు (తప్పనిసరి)
target_locales: [fr, de, es, ja, zh-CN, ar]

# సోర్స్ లొకేల్ ఫైల్ పాత్ (తప్పనిసరి)
source_path: locales/en.json

# LLM ప్రొవైడర్ సెట్టింగ్‌లు
llm:
  # ప్రొవైడర్: "anthropic", "openai", "gemini", లేదా "openrouter" (డిఫాల్ట్: gemini)
  provider: gemini

  # ప్రొవైడర్ ద్వారా మోడల్ పేరు డిఫాల్ట్‌లు:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # API కీని కలిగి ఉన్న ఎన్విరాన్‌మెంట్ వేరియబుల్
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # OpenAI-అనుకూల ఎండ్‌పాయింట్‌ల కోసం బేస్ URL (ఆప్షనల్)
  # base_url: https://api.openai.com

# ప్రతి LLM కాల్‌కు కీలు (డిఫాల్ట్: 40)
batch_size: 40

# ప్యారలల్ LLM కాల్స్ (డిఫాల్ట్: 4)
concurrency: 4

# లొకేల్ వారీగా స్టైల్ గైడ్ మార్క్‌డౌన్ ఫైల్‌లను కలిగి ఉన్న డైరెక్టరీ (డిఫాల్ట్: style-guides)
style_guides_dir: style-guides

# లొకేల్ వారీగా గ్లాసరీ JSON ఫైల్‌లను కలిగి ఉన్న డైరెక్టరీ (డిఫాల్ట్: glossary)
glossary_dir: glossary

# ట్రాన్స్‌లేషన్ మెమరీ ఫైల్ పాత్ (డిఫాల్ట్: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## స్టైల్ గైడ్‌లు

స్టైల్ గైడ్‌లు అనేవి LLM అనువాద ప్రాంప్ట్‌లోకి ఇంజెక్ట్ చేయబడే మార్క్‌డౌన్ ఫైల్‌లు. ఇవి టోన్, ఫార్మాలిటీ, టైపోగ్రఫీ మరియు ఇతర భాషా-నిర్దిష్ట కన్వెన్షన్‌లను నియంత్రిస్తాయి.

```
style-guides/
  _conventions.md    # అన్ని భాషల కోసం షేర్ చేయబడిన నియమాలు
  fr.md              # ఫ్రెంచ్-నిర్దిష్ట నియమాలు
  ja.md              # జపనీస్-నిర్దిష్ట నియమాలు
  ar.md              # అరబిక్-నిర్దిష్ట నియమాలు
```

### షేర్ చేయబడిన కన్వెన్షన్‌లు (`_conventions.md`)

అన్ని భాషలకు వర్తించే నియమాలను నిర్వచించండి: ఇంటర్‌పోలేషన్ సింటాక్స్, HTML ప్రిజర్వేషన్, స్ట్రింగ్ టైప్ కన్వెన్షన్‌లు (బటన్‌లు vs. లేబుల్‌లు vs. ఎర్రర్‌లు) మొదలైనవి.

### భాషల వారీగా గైడ్‌లు (`{locale}.md`)

భాషా-నిర్దిష్ట నియమాలను నిర్వచించండి: ఫార్మాలిటీ రిజిస్టర్ (tu vs. vous), విరామచిహ్నాలు (guillemets, inverted question marks), బహువచన రూపాలు, తేదీ/సంఖ్య ఫార్మాటింగ్ మరియు పరిభాష గ్లాసరీ.

పనిచేసే ఉదాహరణ కోసం [`examples/react-app/style-guides/`](examples/react-app/style-guides/) చూడండి.

## గ్లాసరీ ఫార్మాట్

గ్లాసరీ ఫైల్‌లు `{glossary_dir}/{locale}.json` లో నిల్వ చేయబడిన JSON శ్రేణులు:

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

పదాలు LLM ప్రాంప్ట్‌లోకి టెర్మినాలజీ టేబుల్‌గా ఇంజెక్ట్ చేయబడతాయి, ఇది మీ అప్లికేషన్ అంతటా కీలక పదాల స్థిరమైన అనువాదాన్ని నిర్ధారిస్తుంది.

## ట్రాన్స్‌లేషన్ మెమరీ

ట్రాన్స్‌లేషన్ మెమరీ JSONL ఫైల్‌గా నిల్వ చేయబడుతుంది (లైన్‌కు ఒక JSON రికార్డ్). ప్రతి రికార్డ్ వీటిని కలిగి ఉంటుంది:

- సోర్స్ కీ మరియు వాల్యూ
- అనువదించబడిన వాల్యూ
- సోర్స్ వాల్యూ యొక్క SHA-256 హాష్
- ఒక టైమ్‌స్టాంప్

తదుపరి రన్‌లలో, మార్చబడని స్ట్రింగ్‌లు LLMని కాల్ చేయకుండానే TM కాష్ నుండి అందించబడతాయి, దీని వలన సమయం మరియు API ఖర్చులు రెండూ ఆదా అవుతాయి. TM ఫైల్ git-ఫ్రెండ్లీ మరియు మీ లొకేల్ ఫైల్‌లతో పాటు కమిట్ చేయవచ్చు.

## సపోర్ట్ చేసే ఫార్మాట్‌లు

| ఫార్మాట్ | ఎక్స్‌టెన్షన్‌లు | మోడ్ |
|--------|-----------|------|
| JSON | `.json` | కీ-వాల్యూ (నెస్టెడ్, డాట్-నొటేషన్ ఫ్లాటెన్డ్) |
| YAML | `.yml`, `.yaml` | కీ-వాల్యూ (కామెంట్‌లు మరియు ఆర్డరింగ్‌ను సంరక్షిస్తుంది) |
| Markdown | `.md`, `.mdx` | పూర్తి-డాక్యుమెంట్ అనువాదం |

## ప్రాజెక్ట్ టైప్ డిటెక్షన్

`internationalizer detect` వీటిని తనిఖీ చేయడం ద్వారా మీ i18n సెటప్‌ను గుర్తిస్తుంది:

- react-i18next, next-intl, లేదా vue-i18n కోసం `package.json` డిపెండెన్సీలు
- సాధారణ లొకేల్ ప్యాటర్న్‌లకు సరిపోలే డైరెక్టరీ స్ట్రక్చర్‌లు
- ఫైల్ ఎక్స్‌టెన్షన్‌లు మరియు నేమింగ్ కన్వెన్షన్‌లు

## ఆర్కిటెక్చర్

```
cmd/internationalizer/     CLI ఎంట్రీ పాయింట్ మరియు కమాండ్ డెఫినిషన్స్
internal/
  config/                  డిఫాల్ట్‌లతో YAML కాన్ఫిగ్ లోడింగ్
  detect/                  ప్రాజెక్ట్ టైప్ ఆటో-డిటెక్షన్
  formats/                 ఫార్మాట్ పార్సర్‌లు (JSON, YAML, Markdown)
  glossary/                లొకేల్ వారీగా గ్లాసరీ మేనేజ్‌మెంట్
  llm/                     LLM ప్రొవైడర్ ఇంటర్‌ఫేస్ + ఇంప్లిమెంటేషన్స్
    anthropic.go           Anthropic Claude బ్యాకెండ్
    openai.go              OpenAI / అనుకూల బ్యాకెండ్
    gemini.go              AI Studio బ్యాకెండ్ ద్వారా Google Gemini
                           OpenRouter కస్టమ్ base_url తో openai.go ని ఉపయోగిస్తుంది
  styleguide/              స్టైల్ గైడ్ లోడర్
  tm/                      JSONL ట్రాన్స్‌లేషన్ మెమరీ
  translate/               ట్రాన్స్‌లేషన్ ఆర్కెస్ట్రేటర్
  validate/                లొకేల్ వాలిడేషన్ మరియు డిఫింగ్
```

## ప్రత్యామ్నాయాలతో పోలిక

| ఫీచర్ | Internationalizer | i18next | Crowdin | సాధారణ LLM |
|---------|------------------|---------|---------|-------------|
| LLM-ఆధారిత అనువాదం | అవును | కాదు | పాక్షికం | అవును |
| భాషల వారీగా స్టైల్ గైడ్‌లు | అవును | కాదు | కాదు | కాదు |
| గ్లాసరీ ఎన్‌ఫోర్స్‌మెంట్ | అవును | కాదు | అవును | కాదు |
| ట్రాన్స్‌లేషన్ మెమరీ | అవును | కాదు | అవును | కాదు |
| CLI / లోకల్ ఎగ్జిక్యూషన్ | అవును | వర్తించదు | కాదు | మాన్యువల్ |
| Git-ఫ్రెండ్లీ ఫైల్‌లు | అవును | అవును | పాక్షికం | మాన్యువల్ |
| SaaS డిపెండెన్సీ లేదు | అవును | అవును | కాదు | మారుతుంది |
| ఓపెన్ సోర్స్ (AGPL-3.0) | అవును | అవును | కాదు | మారుతుంది |

## లైసెన్స్

[AGPL-3.0](LICENSE)

## కాంట్రిబ్యూటింగ్

డెవలప్‌మెంట్ సెటప్ మరియు మార్గదర్శకాల కోసం [CONTRIBUTING.md](CONTRIBUTING.md) చూడండి. అన్ని కాంట్రిబ్యూషన్‌లకు DCO సైన్-ఆఫ్ అవసరం.

