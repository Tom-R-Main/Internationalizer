> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

సాఫ్ట్‌వేర్ ప్రాజెక్ట్‌ల కోసం AI-నేటివ్ ఇంటర్నేషనలైజేషన్ పైప్‌లైన్. LLMలను ఉపయోగించి i18n ఫైల్‌లను అనువదించండి, ధృవీకరించండి మరియు నిర్వహించండి.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Internationalizer ఎందుకు?

చాలా i18n టూల్స్ రన్‌టైమ్ లైబ్రరీలు (i18next, react-intl) లేదా కీ-మేనేజ్‌మెంట్ SaaS ప్లాట్‌ఫారమ్‌లు (Crowdin, Lokalise). వాటిలో ఏవీ వాస్తవ అనువాద సమస్యను సమర్థవంతంగా పరిష్కరించలేవు:

- **మాన్యువల్ అనువాదం** కొన్ని భాషల కంటే ఎక్కువ విస్తరించలేదు
- **మెషిన్ ట్రాన్స్‌లేషన్ APIలు** (Google Translate, DeepL) మీ పరిభాష, టోన్ మరియు UI కన్వెన్షన్‌లను పరిగణనలోకి తీసుకోవు
- **సాధారణ LLM అనువాదం** మెరుగ్గానే పనిచేస్తుంది, కానీ గ్లాసరీలు మరియు స్టైల్ గైడ్‌లు లేకపోతే అస్థిరమైన ఫలితాలు వస్తాయి

Internationalizer భిన్నమైనది. ఇది LLM అనువాదాన్ని వీటితో అనుసంధానించే ఒక **CLI పైప్‌లైన్**:

- **భాషల వారీగా గ్లాసరీలు** — మీ యాప్ అంతటా స్థిరమైన సాంకేతిక పరిభాషను అమలు చేయండి
- **భాషల వారీగా స్టైల్ గైడ్‌లు** — టోన్, ఫార్మాలిటీ, బహువచనాలు మరియు టైపోగ్రఫీని నియంత్రించండి
- **ట్రాన్స్‌లేషన్ మెమరీ** — మార్పులేని స్ట్రింగ్‌లను దాటవేయండి, తద్వారా API కాల్ ఖర్చులను ఆదా చేయండి
- **డిటర్మినిస్టిక్ వాలిడేషన్** — కోడ్ విడుదల కావడానికి ముందే తప్పిపోయిన లేదా అదనపు కీలు, ప్రొటెక్టెడ్-స్ట్రక్చర్ మార్పులు, గ్లాసరీ సమస్యలు మరియు బహువచన లేదా ICU లోపాలను గుర్తించండి

<!-- internationalizer:unit markdown:installation -->
## ఇన్‌స్టాలేషన్

npm నుండి ఇన్‌స్టాల్ చేయండి:

```bash
npm install -g internationalizer
```

లేదా గ్లోబల్ ఇన్‌స్టాలేషన్ లేకుండా నేరుగా రన్ చేయండి:

```bash
npx internationalizer --help
```

npm ప్యాకేజీ ప్లాట్‌ఫారమ్-నిర్దిష్ట ఆప్షనల్ డిపెండెన్సీల ద్వారా npm నుండి సరిపోలే ప్రీబిల్ట్ బైనరీని ఇన్‌స్టాల్ చేస్తుంది.

Go తో ఇన్‌స్టాల్ చేయండి:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

లేదా సోర్స్ కోడ్ నుండి బిల్డ్ చేయండి:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

<!-- internationalizer:unit markdown:npm-packages -->
## npm ప్యాకేజీలు

- Git ట్యాగ్‌లు మరియు npm ప్యాకేజీ వెర్షన్‌లు ఖచ్చితంగా సరిపోలాలి, ఉదాహరణకు `v0.1.0` మరియు `0.1.0`
- రూట్ `internationalizer` ప్యాకేజీ `internationalizer-darwin-arm64` వంటి ప్లాట్‌ఫారమ్ ప్యాకేజీలపై ఆధారపడి ఉంటుంది
- మద్దతు ఉన్న npm టార్గెట్‌లు: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI పబ్లిషింగ్ కోసం `NPM_TOKEN` పేరుతో ఒక GitHub సీక్రెట్ అవసరం

<!-- internationalizer:unit markdown:quick-start -->
## త్వరిత ప్రారంభం

1. ప్రాజెక్ట్ రూట్ డైరెక్టరీలో కాన్ఫిగరేషన్ ఫైల్‌ను సృష్టించండి:

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

5. అన్ని లొకేల్‌లను ధృవీకరించండి:

```bash
internationalizer validate
```

<!-- internationalizer:unit markdown:commands -->
## ఆదేశాలు

### `translate`

తప్పిపోయిన లేదా పాతబడిన కీలను గుర్తించి, LLM ద్వారా వాటిని అనువదించండి.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

అనువాద స్థితి మిస్సింగ్, సోర్స్-స్టేల్, పాలసీ-స్టేల్, కరెంట్ మరియు మాన్యువల్‌గా ఎడిట్ చేయబడిన స్థితులను వేర్వేరుగా నివేదిస్తుంది; కనుక మాన్యువల్ ఎడిట్ సోర్స్ లేదా పాలసీ మార్పును దాచలేదు. పాలసీ-స్టేల్ విలువలు నివేదించబడతాయి, కానీ కేవలం `--refresh-policy` ఉపయోగించినప్పుడు మాత్రమే తిరిగి అనువదించబడతాయి. మాన్యువల్‌గా ఎడిట్ చేయబడిన విలువలు ఎప్పటికీ ఆటోమేటిక్‌గా ఓవర్‌రైట్ కావు. సమీక్షించిన అనువాదాలకు మేనిఫెస్ట్‌ను పరిచయం చేస్తున్నప్పుడు లేదా సమీక్షించిన మాన్యువల్ ఎడిట్‌ను కొత్త బేస్‌లైన్‌గా స్పష్టంగా ఆమోదించేటప్పుడు `--adopt-existing` ను ఉపయోగించండి.

### `validate`

సోర్స్ బండిల్స్‌తో పోల్చి అన్ని లొకేల్ ఫైల్‌లను తనిఖీ చేయండి. డిఫాల్ట్ ధృవీకరణ స్ట్రక్చరల్ కవరేజీని (అవసరమైన టార్గెట్ కీలు ఎంత శాతం ఉన్నాయో) తనిఖీ చేస్తుంది, అదనపు కీలను హెచ్చరికలుగా నివేదిస్తుంది మరియు తప్పిపోయిన కీలు, ఇంటర్‌పోలేషన్ అసమతుల్యతలు లేదా చెల్లని ICU MessageFormat నిర్మాణం ఉన్నప్పుడు విఫలమవుతుంది.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` అనువదించబడిన కవరేజీని కూడా నివేదిస్తుంది. సోర్స్‌తో సమానంగా ఉండే భాషాపరమైన విలువ అనువదించబడనట్లే పరిగణించబడుతుంది; అయితే మొత్తం విలువకు ఖచ్చితమైన అదే-సోర్స్, అదే-టార్గెట్ ఎంట్రీ గ్లాసరీలో స్పష్టంగా ఉంటే తప్ప మినహాయింపు ఉండదు; `ignore_case` గౌరవించబడుతుంది, కానీ పెద్ద విలువలో భాగమైన ఒక గ్లాసరీ పదానికి మినహాయింపు వర్తించదు. అదనపు కీలు, సోర్స్‌తో సమానమైన విలువలు, మారిన ఇంటర్‌పోలేషన్/HTML/కోడ్/మార్క్‌డౌన్-లింక్ నిర్మాణం, గ్లాసరీ ఉల్లంఘనలు మరియు కాన్ఫిగర్ చేసిన బహువచన రూపాల లోపాలపై స్ట్రిక్ట్ మోడ్ విఫలమవుతుంది.

`--require-state` ప్రతి టార్గెట్‌ను `.internationalizer.lock` ఆధారంగా ధృవీకరిస్తుంది. ఒక కీ అన్‌ట్రాక్ చేయబడినా, లేదా దాని రికార్డ్ చేసిన సోర్స్, అనువాద విధానం (పాలసీ) లేదా టార్గెట్ హాష్ పాతబడినా (స్టేల్ అయినా) ఇది విఫలమవుతుంది. దీనిని `--strict` తో కలిపి కూడా ఉపయోగించవచ్చు.

మానవ రీడబుల్ మరియు JSON నివేదికలు స్థిరమైన ఫైండింగ్ కోడ్‌లను ఉపయోగిస్తాయి:

| కోడ్ | అర్థం |
| --- | --- |
| `missing_key` / `extra_key` | సోర్స్ మరియు టార్గెట్ కీ సెట్‌లు భిన్నంగా ఉన్నాయి |
| `blank_translation` | ఖాళీ లేని సోర్స్‌కు స్ట్రిక్ట్-మోడ్‌లో ఖాళీ టార్గెట్ ఉంది |
| `source_identical` | స్ట్రిక్ట్-మోడ్‌లో భాషాపరమైన విలువ అనువదించబడకుండా మిగిలిపోయింది |
| `protected_structure_mismatch` | ఇంటర్‌పోలేషన్, HTML, కోడ్ లేదా లింక్ నిర్మాణం మారింది |
| `glossary_violation` | ఆమోదించబడిన టార్గెట్ పదం లేదా వేరియంట్ కనుగొనబడలేదు |
| `plural_form_missing` | కాన్ఫిగర్ చేసిన లొకేల్ బహువచన రూపం లేదు |
| `icu_message_syntax` | సోర్స్ లేదా టార్గెట్ ICU సందేశం సరిగ్గా లేదు |
| `icu_argument_mismatch` | ICU ఆర్గ్యుమెంట్ పేర్లు, రకాలు లేదా ఫార్మాటర్ శైలులు భిన్నంగా ఉన్నాయి |
| `icu_selector_mismatch` | సెలెక్టర్లు భిన్నంగా ఉన్నాయి లేదా టార్గెట్ లొకేల్‌కు బహువచన కేటగిరీ చెల్లదు |
| `untracked` | టార్గెట్ కోసం ఎలాంటి మేనిఫెస్ట్ రికార్డ్ లేదు |
| `source_stale` | రికార్డ్ చేసిన అనువాదం తర్వాత సోర్స్ కంటెంట్ మారింది |
| `policy_stale` | రూపొందించిన ప్రాంప్ట్ లేదా మోడల్ సెట్టింగ్‌లు మారాయి |
| `target_modified` | టార్గెట్ కంటెంట్ మేనిఫెస్ట్ రికార్డుతో సరిపోలడం లేదు |

### `detect`

i18n ఫ్రేమ్‌వర్క్‌ను ఆటోమేటిక్‌గా గుర్తించి, అనుకూలమైన కాన్ఫిగరేషన్‌ను సూచించండి.

```bash
internationalizer detect
```

మద్దతు ఉన్నవి: react-i18next, next-intl, vue-i18n, సాధారణ JSON, మార్క్‌డౌన్ డాక్యుమెంట్‌లు.

### `glossary`

అనువాద సమయంలో అమలు చేయబడే భాషల వారీగా గ్లాసరీ నిబంధనలను నిర్వహించండి.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ట్రాన్స్‌లేషన్ మెమరీని నిర్వహించండి (గతంలో అనువదించబడిన స్ట్రింగ్‌ల JSONL కాష్).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## కాన్ఫిగరేషన్ రిఫరెన్స్

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

లొకేల్ ఐడెంటిఫైయర్‌లు `fr`, `pt-BR`, లేదా `sr-Latn-RS` వంటి సరైన రూపంలో ఉన్న BCP 47 ట్యాగ్‌లై ఉండాలి. ప్రామాణిక సమానమైన (canonical-equivalent) టార్గెట్ లొకేల్‌లు డూప్లికేట్‌లుగా తిరస్కరించబడతాయి, అలాగే లొకేల్-నిర్దిష్ట ప్రొవైడర్ ఓవర్‌రైడ్‌లు ప్రామాణిక సమానమైన స్పెల్లింగ్‌తో సరిపోలుతాయి. పై ఉదాహరణలో, జపనీస్‌తో సహా ఎలాంటి ఓవర్‌రైడ్ లేని లొకేల్‌లు గ్లోబల్ Gemini కాన్ఫిగరేషన్‌ను పొందుతాయి.

ICU MessageFormat విలువలు నిర్మాణాత్మకంగా అన్వయించబడతాయి. సాధారణ ఆర్గ్యుమెంట్‌లు, `select`, `plural`, `selectordinal`, `number`, `date` మరియు `time` లకు మద్దతు ఉంది; వీటిలో నెస్టెడ్ సందేశాలు, బహువచన ఆఫ్‌సెట్‌లు, ఖచ్చితమైన సంఖ్య సెలెక్టర్లు మరియు `#` లు కూడా ఉంటాయి. సింటాక్స్, ఆర్గ్యుమెంట్ రకాలు మరియు ఫార్మాటర్ శైలులు, బహువచన ఆఫ్‌సెట్‌లు, సెలెక్ట్ బ్రాంచ్ గుర్తింపు మరియు టార్గెట్ లొకేల్ CLDR బహువచన కేటగిరీలను వాలిడేషన్ తనిఖీ చేస్తుంది. ఈ నియమాలను ఉల్లంఘించే ప్రొవైడర్ అవుట్‌పుట్, లొకేల్ ఫైల్ లేదా ట్రాన్స్‌లేషన్ మెమరీ రికార్డ్ రాయబడక ముందే తిరస్కరించబడుతుంది.

`i18next-v4` తో, గుర్తించబడిన సోర్స్ బహువచన కుటుంబాలు అనువాద సమయంలో టార్గెట్ లొకేల్ యొక్క CLDR కేటగిరీలకు విస్తరించబడతాయి. టార్గెట్‌కు మాత్రమే చెందిన కేటగిరీ సోర్స్ కుటుంబం యొక్క `_other` విలువను తన అనువాద టెంప్లేట్‌గా ఉపయోగిస్తుంది. స్ట్రిక్ట్ వాలిడేషన్‌కు ఆ టార్గెట్ కేటగిరీలు తప్పనిసరిగా అవసరం; వాటిని ఉపయోగించని టార్గెట్ లొకేల్‌లకు సోర్స్‌కు మాత్రమే చెందిన కేటగిరీలు ఆప్షనల్.

<!-- internationalizer:unit markdown:style-guides -->
## స్టైల్ గైడ్‌లు

స్టైల్ గైడ్‌లు అనేవి LLM అనువాద ప్రాంప్ట్‌లోకి చేర్చబడే మార్క్‌డౌన్ ఫైల్‌లు. ఇవి టోన్, ఫార్మాలిటీ, టైపోగ్రఫీ మరియు ఇతర భాషా-నిర్దిష్ట నియమాలను నియంత్రిస్తాయి.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### ఉమ్మడి నియమాలు (`_conventions.md`)

అన్ని భాషలకు వర్తించే నియమాలను నిర్వచించండి: ఇంటర్‌పోలేషన్ సింటాక్స్, HTML పరిరక్షణ, స్ట్రింగ్ రకం కన్వెన్షన్‌లు (బటన్‌లు vs. లేబుల్‌లు vs. లోపాలు) మొదలైనవి.

### భాషల వారీ గైడ్‌లు (`{locale}.md`)

భాషా-నిర్దిష్ట నియమాలను నిర్వచించండి: మర్యాదపూర్వక సంబోధన (tu vs. vous), విరామచిహ్నాలు (guillemets, ఇన్వర్టెడ్ ప్రశ్నార్థకాలు), బహువచన రూపాలు, తేదీ/సంఖ్య ఫార్మాటింగ్ మరియు పరిభాష గ్లాసరీ.

స్టైల్ గైడ్‌లు అనేవి శాశ్వతమైన పాలసీ ఇన్‌పుట్‌లు, రూపొందించబడిన అవుట్‌పుట్ కాదు. Internationalizer వాటిని చదువుతుంది కానీ ఎప్పటికీ తిరిగి రాయదు. వాటి కంటెంట్ గ్లాసరీ మరియు ప్రాంప్ట్ కాంట్రాక్ట్ నుండి విడిగా హాష్ చేయబడుతుంది, కనుక అప్లికేషన్ కోడ్ మార్పు అనువాదాన్ని పాతబడేలా (స్టేల్ అయ్యేలా) చేయదు. ఒక గైడ్‌ను ఎడిట్ చేసినప్పుడు ఆ లొకేల్ ఉద్దేశపూర్వకంగా పాలసీ సమీక్ష కోసం గుర్తించబడుతుంది; ప్రాంప్ట్ కాంట్రాక్ట్ వెర్షన్ మారితే తప్ప అంతర్గత ప్రాంప్ట్ పదాల మార్పు దానిని ప్రభావితం చేయదు.

పనిచేసే ఉదాహరణ కోసం [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) చూడండి.

<!-- internationalizer:unit markdown:glossary-format -->
## గ్లాసరీ ఫార్మాట్

గ్లాసరీ ఫైల్‌లు `{glossary_dir}/{locale}.json` లో భద్రపరచబడే JSON శ్రేణులు (arrays):

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

`variants` అనేది ఆమోదించబడిన ఇతర టార్గెట్ రూపాలను సూచిస్తుంది. `enforcement` అనేది `error`, `warning` కావచ్చు లేదా డిఫాల్ట్ ఎర్రర్ ప్రవర్తన కోసం వదిలివేయబడవచ్చు. పదాలు LLM ప్రాంప్ట్‌లోకి టెర్మినాలజీ పట్టికగా ఇంజెక్ట్ చేయబడతాయి, ఇది మీ అప్లికేషన్ అంతటా స్థిరమైన అనువాదాన్ని నిర్ధారిస్తుంది. `{"source":"API","target":"API"}` వంటి ఖచ్చితమైన ఎంట్రీ, పూర్తిగా సోర్స్‌తో సమానమైన ఆ విలువను స్ట్రిక్ట్ అన్‌ట్రాన్స్‌లేటెడ్-విలువ ఫైండింగ్‌ల నుండి మినహాయిస్తుంది; అయితే కేవలం `API` ఉన్న పెద్ద విలువకు ఈ మినహాయింపు వర్తించదు.

<!-- internationalizer:unit markdown:translation-memory -->
## ట్రాన్స్‌లేషన్ మెమరీ

ట్రాన్స్‌లేషన్ మెమరీ ఒక JSONL ఫైల్‌గా భద్రపరచబడుతుంది (ఒక్కో లైన్‌కు ఒక JSON రికార్డు). ప్రతి రికార్డు వీటిని కలిగి ఉంటుంది:

- బండిల్, కీ, సోర్స్ విలువ, అనువదించబడిన విలువ మరియు ప్రామాణిక టార్గెట్ లొకేల్
- సోర్స్, స్టైల్ గైడ్, గ్లాసరీ, ప్రాంప్ట్-కాంట్రాక్ట్ మరియు మిశ్రమ పాలసీ హాష్‌లు
- అనువాదాన్ని రూపొందించిన ప్రొవైడర్ మరియు మోడల్
- ఒక టైమ్‌స్టాంప్

తదుపరి రన్‌లలో, అదే సోర్స్ మరియు పాలసీ హాష్‌లను కలిగి ఉన్న స్ట్రింగ్‌లు LLMకి కాల్ చేయకుండానే కాష్ నుండి అందించబడతాయి. డిఫాల్ట్ పాత్ ఇగ్నోర్ చేయబడిన `.internationalizer/` డైరెక్టరీలో ఉంటుంది, కాబట్టి ఇది లోకల్ కాష్‌గానే ఉంటుంది. మీ ప్రాజెక్ట్ ఉద్దేశపూర్వకంగా ట్రాన్స్‌లేషన్ మెమరీని షేర్ చేసుకోవాలనుకుంటే `tm_path` ను ట్రాక్ చేయబడే ప్రదేశానికి సెట్ చేయండి. సమీక్షించదగిన `.internationalizer.lock` మేనిఫెస్ట్ విడిగా వెర్షన్ చేయబడుతుంది.

<!-- internationalizer:unit markdown:supported-formats -->
## మద్దతు ఉన్న ఫార్మాట్‌లు

| ఫార్మాట్ | ఎక్స్‌టెన్షన్‌లు | మోడ్ |
|--------|-----------|------|
| JSON | `.json` | కీ-విలువ (నెస్టెడ్, డాట్-నొటేషన్ ఫ్లాటెన్డ్) |
| YAML | `.yml`, `.yaml` | కీ-విలువ (కామెంట్‌లు మరియు క్రమాన్ని పరిరక్షిస్తుంది) |
| Markdown | `.md`, `.mdx` | పీఠిక మరియు H2-స్థాయి విభాగాలు |

మార్క్‌డౌన్ టార్గెట్‌లు H2 విభాగాల కంటే ముందు అదృశ్య `internationalizer:unit` కామెంట్‌లను కలిగి ఉంటాయి. ఈ స్థిరమైన మార్కర్‌లు ఇతర విభాగాలను తిరిగి అనువదించాల్సిన అవసరం లేకుండానే ఏదైనా ఒక సోర్స్ విభాగాన్ని జోడించడానికి, తరలించడానికి లేదా సవరించడానికి Internationalizer కు సహాయపడతాయి. మార్కర్‌లు లేని ఇప్పటికే ఉన్న డాక్యుమెంట్‌లు వాటి తదుపరి విజయవంతమైన అప్‌డేట్‌లో ఈ మార్కర్‌లను పొందుతాయి.

<!-- internationalizer:unit markdown:project-type-detection -->
## ప్రాజెక్ట్ రకం గుర్తింపు

`internationalizer detect` వీటిని తనిఖీ చేయడం ద్వారా మీ i18n సెటప్‌ను గుర్తిస్తుంది:

- react-i18next, next-intl, లేదా vue-i18n కోసం `package.json` డిపెండెన్సీలు
- సాధారణ లొకేల్ ప్యాటర్న్‌లకు సరిపోలే డైరెక్టరీ నిర్మాణాలు
- ఫైల్ ఎక్స్‌టెన్షన్‌లు మరియు నేమింగ్ కన్వెన్షన్‌లు

<!-- internationalizer:unit markdown:architecture -->
## ఆర్కిటెక్చర్

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
## ప్రత్యామ్నాయాలతో పోలిక

| ఫీచర్ | Internationalizer | i18next | Crowdin | సాధారణ LLM |
|---------|------------------|---------|---------|-------------|
| LLM-ఆధారిత అనువాదం | అవును | కాదు | పాక్షికం | అవును |
| భాషల వారీగా స్టైల్ గైడ్‌లు | అవును | కాదు | కాదు | కాదు |
| గ్లాసరీ అమలు | అవును | కాదు | అవును | కాదు |
| ట్రాన్స్‌లేషన్ మెమరీ | అవును | కాదు | అవును | కాదు |
| CLI / లోకల్ ఎగ్జిక్యూషన్ | అవును | వర్తించదు | కాదు | మాన్యువల్ |
| Git-ఫ్రెండ్లీ ఫైల్‌లు | అవును | అవును | పాక్షికం | మాన్యువల్ |
| SaaS డిపెండెన్సీ లేదు | అవును | అవును | కాదు | మారుతుంది |
| ఓపెన్ సోర్స్ (AGPL-3.0) | అవును | అవును | కాదు | మారుతుంది |

<!-- internationalizer:unit markdown:license -->
## లైసెన్స్

[AGPL-3.0](../../LICENSE)

డిపెండెన్సీ నోటీసుల కోసం [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) చూడండి.

<!-- internationalizer:unit markdown:contributing -->
## సహకారం

డెవలప్‌మెంట్ సెటప్ మరియు మార్గదర్శకాల కోసం [CONTRIBUTING.md](../../CONTRIBUTING.md) చూడండి. అన్ని కాంట్రిబ్యూషన్‌లకు DCO సైన్-ఆఫ్ అవసరం.
