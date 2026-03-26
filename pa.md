<p align="center">
  <img src="assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

ਸਾਫਟਵੇਅਰ ਪ੍ਰੋਜੈਕਟਾਂ ਲਈ AI-native internationalization ਪਾਈਪਲਾਈਨ। LLMs ਦੀ ਵਰਤੋਂ ਕਰਕੇ i18n ਫਾਈਲਾਂ ਦਾ ਅਨੁਵਾਦ ਕਰੋ, ਪ੍ਰਮਾਣਿਤ ਕਰੋ ਅਤੇ ਪ੍ਰਬੰਧਿਤ ਕਰੋ।

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

## Internationalizer ਕਿਉਂ?

ਜ਼ਿਆਦਾਤਰ i18n ਟੂਲ ਜਾਂ ਤਾਂ ਰਨਟਾਈਮ ਲਾਇਬ੍ਰੇਰੀਆਂ (i18next, react-intl) ਹਨ ਜਾਂ key-management SaaS ਪਲੇਟਫਾਰਮ (Crowdin, Lokalise) ਹਨ। ਇਹਨਾਂ ਵਿੱਚੋਂ ਕੋਈ ਵੀ ਅਸਲ ਅਨੁਵਾਦ ਸਮੱਸਿਆ ਨੂੰ ਚੰਗੀ ਤਰ੍ਹਾਂ ਹੱਲ ਨਹੀਂ ਕਰਦਾ:

- **ਮੈਨੁਅਲ ਅਨੁਵਾਦ** ਕੁਝ ਭਾਸ਼ਾਵਾਂ ਤੋਂ ਅੱਗੇ ਸਕੇਲ ਨਹੀਂ ਹੁੰਦਾ
- **ਮਸ਼ੀਨ ਅਨੁਵਾਦ APIs** (Google Translate, DeepL) ਤੁਹਾਡੀ ਸ਼ਬਦਾਵਲੀ, ਟੋਨ ਅਤੇ UI ਪਰੰਪਰਾਵਾਂ ਨੂੰ ਨਜ਼ਰਅੰਦਾਜ਼ ਕਰਦੇ ਹਨ
- **ਆਮ LLM ਅਨੁਵਾਦ** ਬਿਹਤਰ ਕੰਮ ਕਰਦਾ ਹੈ, ਪਰ ਗਲੋਸਰੀਆਂ ਅਤੇ ਸਟਾਈਲ ਗਾਈਡਾਂ ਤੋਂ ਬਿਨਾਂ, ਤੁਹਾਨੂੰ ਅਸੰਗਤ ਨਤੀਜੇ ਮਿਲਦੇ ਹਨ

Internationalizer ਵੱਖਰਾ ਹੈ। ਇਹ ਇੱਕ **CLI ਪਾਈਪਲਾਈਨ** ਹੈ ਜੋ LLM ਅਨੁਵਾਦ ਨੂੰ ਇਹਨਾਂ ਨਾਲ ਜੋੜਦੀ ਹੈ:

- **ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਗਲੋਸਰੀਆਂ** — ਤੁਹਾਡੀ ਐਪ ਵਿੱਚ ਇਕਸਾਰ ਸ਼ਬਦਾਵਲੀ ਲਾਗੂ ਕਰੋ
- **ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਸਟਾਈਲ ਗਾਈਡਾਂ** — ਟੋਨ, ਰਸਮੀਪਣ, ਬਹੁਵਚਨ ਅਤੇ ਟਾਈਪੋਗ੍ਰਾਫੀ ਨੂੰ ਕੰਟਰੋਲ ਕਰੋ
- **ਅਨੁਵਾਦ ਮੈਮੋਰੀ** — ਨਾ ਬਦਲੀਆਂ ਸਟ੍ਰਿੰਗਾਂ ਨੂੰ ਛੱਡੋ, API ਕਾਲਾਂ 'ਤੇ ਪੈਸੇ ਬਚਾਓ
- **ਕੁੰਜੀ ਪ੍ਰਮਾਣਿਕਤਾ** — ਸ਼ਿਪ ਕਰਨ ਤੋਂ ਪਹਿਲਾਂ ਗੁੰਮ ਹੋਏ ਅਨੁਵਾਦਾਂ ਅਤੇ ਇੰਟਰਪੋਲੇਸ਼ਨ ਬੇਮੇਲਾਂ ਨੂੰ ਫੜੋ

## ਇੰਸਟਾਲੇਸ਼ਨ

npm ਤੋਂ ਇੰਸਟਾਲ ਕਰੋ:

```bash
npm install -g internationalizer
```

ਜਾਂ ਗਲੋਬਲ ਇੰਸਟਾਲ ਤੋਂ ਬਿਨਾਂ ਚਲਾਓ:

```bash
npx internationalizer --help
```

npm ਪੈਕੇਜ ਪਲੇਟਫਾਰਮ-ਵਿਸ਼ੇਸ਼ ਵਿਕਲਪਿਕ ਨਿਰਭਰਤਾਵਾਂ ਰਾਹੀਂ npm ਤੋਂ ਮੇਲ ਖਾਂਦੀ ਪ੍ਰੀਬਿਲਟ ਬਾਈਨਰੀ ਨੂੰ ਇੰਸਟਾਲ ਕਰਦਾ ਹੈ।

Go ਨਾਲ ਇੰਸਟਾਲ ਕਰੋ:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

ਜਾਂ ਸਰੋਤ ਤੋਂ ਬਿਲਡ ਕਰੋ:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm ਪੈਕੇਜ

- Git ਟੈਗ ਅਤੇ npm ਪੈਕੇਜ ਵਰਜਨ ਮੇਲ ਖਾਂਦੇ ਹੋਣੇ ਚਾਹੀਦੇ ਹਨ, ਉਦਾਹਰਨ ਲਈ `v0.1.0` ਅਤੇ `0.1.0`
- ਰੂਟ `internationalizer` ਪੈਕੇਜ ਪਲੇਟਫਾਰਮ ਪੈਕੇਜਾਂ 'ਤੇ ਨਿਰਭਰ ਕਰਦਾ ਹੈ ਜਿਵੇਂ ਕਿ `internationalizer-darwin-arm64`
- ਸਮਰਥਿਤ npm ਟੀਚੇ: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI ਪਬਲਿਸ਼ਿੰਗ ਲਈ `NPM_TOKEN` ਨਾਮਕ GitHub ਸੀਕ੍ਰੇਟ ਦੀ ਲੋੜ ਹੁੰਦੀ ਹੈ

## ਤੁਰੰਤ ਸ਼ੁਰੂਆਤ

1. ਆਪਣੇ ਪ੍ਰੋਜੈਕਟ ਰੂਟ ਵਿੱਚ ਇੱਕ ਕੌਂਫਿਗ ਫਾਈਲ ਬਣਾਓ:

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

2. ਆਪਣੀ API ਕੁੰਜੀ ਸੈੱਟ ਕਰੋ:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. ਪੂਰਵ-ਦਰਸ਼ਨ ਕਰੋ ਕਿ ਕੀ ਅਨੁਵਾਦ ਕੀਤਾ ਜਾਵੇਗਾ:

```bash
internationalizer translate --dry-run
```

4. ਅਨੁਵਾਦ ਚਲਾਓ:

```bash
internationalizer translate
```

5. ਸਾਰੇ ਲੋਕੇਲ ਪ੍ਰਮਾਣਿਤ ਕਰੋ:

```bash
internationalizer validate
```

## ਕਮਾਂਡਾਂ

### `translate`

ਗੁੰਮ ਹੋਈਆਂ ਕੁੰਜੀਆਂ ਲੱਭੋ ਅਤੇ ਉਹਨਾਂ ਦਾ LLM ਰਾਹੀਂ ਅਨੁਵਾਦ ਕਰੋ।

```bash
internationalizer translate                    # ਸਾਰੇ ਲੋਕੇਲਾਂ ਦਾ ਅਨੁਵਾਦ ਕਰੋ
internationalizer translate -l fr              # ਸਿਰਫ਼ ਫ੍ਰੈਂਚ ਦਾ ਅਨੁਵਾਦ ਕਰੋ
internationalizer translate --dry-run          # API ਕਾਲਾਂ ਤੋਂ ਬਿਨਾਂ ਪੂਰਵ-ਦਰਸ਼ਨ ਕਰੋ
internationalizer translate --batch-size 20    # ਛੋਟੇ ਬੈਚ
internationalizer translate --concurrency 2    # ਘੱਟ ਸਮਾਨਾਂਤਰ ਕਾਲਾਂ
```

### `validate`

ਗੁੰਮ ਹੋਈਆਂ ਕੁੰਜੀਆਂ, ਵਾਧੂ ਕੁੰਜੀਆਂ, ਅਤੇ ਇੰਟਰਪੋਲੇਸ਼ਨ ਬੇਮੇਲਾਂ ਲਈ ਸਾਰੀਆਂ ਲੋਕੇਲ ਫਾਈਲਾਂ ਦੀ ਜਾਂਚ ਕਰੋ।

```bash
internationalizer validate                     # ਮਨੁੱਖੀ-ਪੜ੍ਹਨਯੋਗ ਆਉਟਪੁੱਟ
internationalizer validate --json              # ਮਸ਼ੀਨ-ਪੜ੍ਹਨਯੋਗ JSON
internationalizer validate -q                  # ਸਿਰਫ਼ ਐਗਜ਼ਿਟ ਕੋਡ
```

### `detect`

i18n ਫਰੇਮਵਰਕ ਦਾ ਸਵੈ-ਪਤਾ ਲਗਾਓ ਅਤੇ ਇੱਕ ਕੌਂਫਿਗਰੇਸ਼ਨ ਦਾ ਸੁਝਾਅ ਦਿਓ।

```bash
internationalizer detect
```

ਸਮਰਥਨ ਕਰਦਾ ਹੈ: react-i18next, next-intl, vue-i18n, vanilla JSON, markdown docs।

### `glossary`

ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਗਲੋਸਰੀ ਸ਼ਬਦਾਂ ਦਾ ਪ੍ਰਬੰਧਨ ਕਰੋ ਜੋ ਅਨੁਵਾਦ ਦੌਰਾਨ ਲਾਗੂ ਕੀਤੇ ਜਾਂਦੇ ਹਨ।

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

ਅਨੁਵਾਦ ਮੈਮੋਰੀ (ਪਹਿਲਾਂ ਅਨੁਵਾਦ ਕੀਤੀਆਂ ਸਟ੍ਰਿੰਗਾਂ ਦਾ JSONL ਕੈਸ਼) ਦਾ ਪ੍ਰਬੰਧਨ ਕਰੋ।

```bash
internationalizer tm stats                     # ਰਿਕਾਰਡ ਗਿਣਤੀ ਦਿਖਾਓ
internationalizer tm export                    # JSON ਵਜੋਂ ਡੰਪ ਕਰੋ
internationalizer tm clear --force             # ਸਾਰੇ ਰਿਕਾਰਡ ਮਿਟਾਓ
```

## ਕੌਂਫਿਗਰੇਸ਼ਨ ਹਵਾਲਾ

```yaml
# .internationalizer.yml

# ਸਰੋਤ ਭਾਸ਼ਾ (ਡਿਫੌਲਟ: en)
source_locale: en

# ਅਨੁਵਾਦ ਕਰਨ ਲਈ ਭਾਸ਼ਾਵਾਂ (ਲੋੜੀਂਦਾ)
target_locales: [fr, de, es, ja, zh-CN, ar]

# ਸਰੋਤ ਲੋਕੇਲ ਫਾਈਲ ਦਾ ਮਾਰਗ (ਲੋੜੀਂਦਾ)
source_path: locales/en.json

# LLM ਪ੍ਰਦਾਤਾ ਸੈਟਿੰਗਾਂ
llm:
  # ਪ੍ਰਦਾਤਾ: "anthropic", "openai", "gemini", ਜਾਂ "openrouter" (ਡਿਫੌਲਟ: gemini)
  provider: gemini

  # ਮਾਡਲ ਨਾਮ ਪ੍ਰਦਾਤਾ ਦੁਆਰਾ ਡਿਫੌਲਟ:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # API ਕੁੰਜੀ ਵਾਲਾ ਵਾਤਾਵਰਣ ਵੇਰੀਏਬਲ
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # OpenAI-ਅਨੁਕੂਲ ਐਂਡਪੁਆਇੰਟਾਂ ਲਈ ਬੇਸ URL (ਵਿਕਲਪਿਕ)
  # base_url: https://api.openai.com

# ਪ੍ਰਤੀ LLM ਕਾਲ ਕੁੰਜੀਆਂ (ਡਿਫੌਲਟ: 40)
batch_size: 40

# ਸਮਾਨਾਂਤਰ LLM ਕਾਲਾਂ (ਡਿਫੌਲਟ: 4)
concurrency: 4

# ਪ੍ਰਤੀ-ਲੋਕੇਲ ਸਟਾਈਲ ਗਾਈਡ ਮਾਰਕਡਾਊਨ ਫਾਈਲਾਂ ਵਾਲੀ ਡਾਇਰੈਕਟਰੀ (ਡਿਫੌਲਟ: style-guides)
style_guides_dir: style-guides

# ਪ੍ਰਤੀ-ਲੋਕੇਲ ਗਲੋਸਰੀ JSON ਫਾਈਲਾਂ ਵਾਲੀ ਡਾਇਰੈਕਟਰੀ (ਡਿਫੌਲਟ: glossary)
glossary_dir: glossary

# ਅਨੁਵਾਦ ਮੈਮੋਰੀ ਫਾਈਲ ਦਾ ਮਾਰਗ (ਡਿਫੌਲਟ: .internationalizer/tm.jsonl)
tm_path: .internationalizer/tm.jsonl
```

## ਸਟਾਈਲ ਗਾਈਡਾਂ

ਸਟਾਈਲ ਗਾਈਡਾਂ ਮਾਰਕਡਾਊਨ ਫਾਈਲਾਂ ਹਨ ਜੋ LLM ਅਨੁਵਾਦ ਪ੍ਰੋਂਪਟ ਵਿੱਚ ਇੰਜੈਕਟ ਕੀਤੀਆਂ ਜਾਂਦੀਆਂ ਹਨ। ਉਹ ਟੋਨ, ਰਸਮੀਪਣ, ਟਾਈਪੋਗ੍ਰਾਫੀ, ਅਤੇ ਹੋਰ ਭਾਸ਼ਾ-ਵਿਸ਼ੇਸ਼ ਪਰੰਪਰਾਵਾਂ ਨੂੰ ਕੰਟਰੋਲ ਕਰਦੀਆਂ ਹਨ।

```
style-guides/
  _conventions.md    # ਸਾਰੀਆਂ ਭਾਸ਼ਾਵਾਂ ਲਈ ਸਾਂਝੇ ਨਿਯਮ
  fr.md              # ਫ੍ਰੈਂਚ-ਵਿਸ਼ੇਸ਼ ਨਿਯਮ
  ja.md              # ਜਾਪਾਨੀ-ਵਿਸ਼ੇਸ਼ ਨਿਯਮ
  ar.md              # ਅਰਬੀ-ਵਿਸ਼ੇਸ਼ ਨਿਯਮ
```

### ਸਾਂਝੀਆਂ ਪਰੰਪਰਾਵਾਂ (`_conventions.md`)

ਉਹ ਨਿਯਮ ਪਰਿਭਾਸ਼ਿਤ ਕਰੋ ਜੋ ਸਾਰੀਆਂ ਭਾਸ਼ਾਵਾਂ 'ਤੇ ਲਾਗੂ ਹੁੰਦੇ ਹਨ: ਇੰਟਰਪੋਲੇਸ਼ਨ ਸਿੰਟੈਕਸ, HTML ਸੰਭਾਲ, ਸਟ੍ਰਿੰਗ ਕਿਸਮ ਦੀਆਂ ਪਰੰਪਰਾਵਾਂ (ਬਟਨ ਬਨਾਮ ਲੇਬਲ ਬਨਾਮ ਤਰੁੱਟੀਆਂ), ਆਦਿ।

### ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਗਾਈਡਾਂ (`{locale}.md`)

ਭਾਸ਼ਾ-ਵਿਸ਼ੇਸ਼ ਨਿਯਮ ਪਰਿਭਾਸ਼ਿਤ ਕਰੋ: ਰਸਮੀਪਣ ਰਜਿਸਟਰ (tu ਬਨਾਮ vous), ਵਿਸ਼ਰਾਮ ਚਿੰਨ੍ਹ (guillemets, ਉਲਟੇ ਪ੍ਰਸ਼ਨ ਚਿੰਨ੍ਹ), ਬਹੁਵਚਨ ਰੂਪ, ਮਿਤੀ/ਨੰਬਰ ਫਾਰਮੈਟਿੰਗ, ਅਤੇ ਇੱਕ ਸ਼ਬਦਾਵਲੀ ਗਲੋਸਰੀ।

ਕੰਮ ਕਰਨ ਵਾਲੀ ਉਦਾਹਰਨ ਲਈ [`examples/react-app/style-guides/`](examples/react-app/style-guides/) ਦੇਖੋ।

## ਗਲੋਸਰੀ ਫਾਰਮੈਟ

ਗਲੋਸਰੀ ਫਾਈਲਾਂ `{glossary_dir}/{locale}.json` ਵਿੱਚ ਸਟੋਰ ਕੀਤੀਆਂ JSON ਐਰੇ ਹਨ:

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

ਸ਼ਬਦਾਂ ਨੂੰ LLM ਪ੍ਰੋਂਪਟ ਵਿੱਚ ਇੱਕ ਸ਼ਬਦਾਵਲੀ ਸਾਰਣੀ ਵਜੋਂ ਇੰਜੈਕਟ ਕੀਤਾ ਜਾਂਦਾ ਹੈ, ਜੋ ਤੁਹਾਡੀ ਐਪਲੀਕੇਸ਼ਨ ਵਿੱਚ ਮੁੱਖ ਸ਼ਬਦਾਂ ਦੇ ਇਕਸਾਰ ਅਨੁਵਾਦ ਨੂੰ ਯਕੀਨੀ ਬਣਾਉਂਦਾ ਹੈ।

## ਅਨੁਵਾਦ ਮੈਮੋਰੀ

ਅਨੁਵਾਦ ਮੈਮੋਰੀ ਨੂੰ ਇੱਕ JSONL ਫਾਈਲ (ਪ੍ਰਤੀ ਲਾਈਨ ਇੱਕ JSON ਰਿਕਾਰਡ) ਵਜੋਂ ਸਟੋਰ ਕੀਤਾ ਜਾਂਦਾ ਹੈ। ਹਰੇਕ ਰਿਕਾਰਡ ਵਿੱਚ ਸ਼ਾਮਲ ਹਨ:

- ਸਰੋਤ ਕੁੰਜੀ ਅਤੇ ਮੁੱਲ
- ਅਨੁਵਾਦ ਕੀਤਾ ਮੁੱਲ
- ਸਰੋਤ ਮੁੱਲ ਦਾ ਇੱਕ SHA-256 ਹੈਸ਼
- ਇੱਕ ਟਾਈਮਸਟੈਂਪ

ਅਗਲੀਆਂ ਦੌੜਾਂ 'ਤੇ, ਨਾ ਬਦਲੀਆਂ ਸਟ੍ਰਿੰਗਾਂ ਨੂੰ LLM ਨੂੰ ਕਾਲ ਕੀਤੇ ਬਿਨਾਂ TM ਕੈਸ਼ ਤੋਂ ਸਰਵ ਕੀਤਾ ਜਾਂਦਾ ਹੈ, ਜਿਸ ਨਾਲ ਸਮਾਂ ਅਤੇ API ਲਾਗਤਾਂ ਦੋਵਾਂ ਦੀ ਬਚਤ ਹੁੰਦੀ ਹੈ। TM ਫਾਈਲ git-ਅਨੁਕੂਲ ਹੈ ਅਤੇ ਤੁਹਾਡੀਆਂ ਲੋਕੇਲ ਫਾਈਲਾਂ ਦੇ ਨਾਲ ਕਮਿਟ ਕੀਤੀ ਜਾ ਸਕਦੀ ਹੈ।

## ਸਮਰਥਿਤ ਫਾਰਮੈਟ

| ਫਾਰਮੈਟ | ਐਕਸਟੈਂਸ਼ਨਾਂ | ਮੋਡ |
|--------|-----------|------|
| JSON | `.json` | Key-value (nested, dot-notation flattened) |
| YAML | `.yml`, `.yaml` | Key-value (ਟਿੱਪਣੀਆਂ ਅਤੇ ਕ੍ਰਮ ਨੂੰ ਸੁਰੱਖਿਅਤ ਰੱਖਦਾ ਹੈ) |
| Markdown | `.md`, `.mdx` | ਪੂਰੇ-ਦਸਤਾਵੇਜ਼ ਦਾ ਅਨੁਵਾਦ |

## ਪ੍ਰੋਜੈਕਟ ਕਿਸਮ ਦੀ ਪਛਾਣ

`internationalizer detect` ਇਹਨਾਂ ਦੀ ਜਾਂਚ ਕਰਕੇ ਤੁਹਾਡੇ i18n ਸੈੱਟਅੱਪ ਦੀ ਪਛਾਣ ਕਰਦਾ ਹੈ:

- react-i18next, next-intl, ਜਾਂ vue-i18n ਲਈ `package.json` ਨਿਰਭਰਤਾਵਾਂ
- ਆਮ ਲੋਕੇਲ ਪੈਟਰਨਾਂ ਨਾਲ ਮੇਲ ਖਾਂਦੀਆਂ ਡਾਇਰੈਕਟਰੀ ਬਣਤਰਾਂ
- ਫਾਈਲ ਐਕਸਟੈਂਸ਼ਨਾਂ ਅਤੇ ਨਾਮਕਰਨ ਪਰੰਪਰਾਵਾਂ

## ਆਰਕੀਟੈਕਚਰ

```
cmd/internationalizer/     CLI ਐਂਟਰੀ ਪੁਆਇੰਟ ਅਤੇ ਕਮਾਂਡ ਪਰਿਭਾਸ਼ਾਵਾਂ
internal/
  config/                  ਡਿਫੌਲਟਾਂ ਨਾਲ YAML ਕੌਂਫਿਗ ਲੋਡਿੰਗ
  detect/                  ਪ੍ਰੋਜੈਕਟ ਕਿਸਮ ਦਾ ਸਵੈ-ਪਤਾ ਲਗਾਉਣਾ
  formats/                 ਫਾਰਮੈਟ ਪਾਰਸਰ (JSON, YAML, Markdown)
  glossary/                ਪ੍ਰਤੀ-ਲੋਕੇਲ ਗਲੋਸਰੀ ਪ੍ਰਬੰਧਨ
  llm/                     LLM ਪ੍ਰਦਾਤਾ ਇੰਟਰਫੇਸ + ਲਾਗੂਕਰਨ
    anthropic.go           Anthropic Claude ਬੈਕਐਂਡ
    openai.go              OpenAI / ਅਨੁਕੂਲ ਬੈਕਐਂਡ
    gemini.go              AI Studio ਬੈਕਐਂਡ ਰਾਹੀਂ Google Gemini
                           OpenRouter ਕਸਟਮ base_url ਨਾਲ openai.go ਦੀ ਵਰਤੋਂ ਕਰਦਾ ਹੈ
  styleguide/              ਸਟਾਈਲ ਗਾਈਡ ਲੋਡਰ
  tm/                      JSONL ਅਨੁਵਾਦ ਮੈਮੋਰੀ
  translate/               ਅਨੁਵਾਦ ਆਰਕੈਸਟ੍ਰੇਟਰ
  validate/                ਲੋਕੇਲ ਪ੍ਰਮਾਣਿਕਤਾ ਅਤੇ ਡਿਫਿੰਗ
```

## ਵਿਕਲਪਾਂ ਨਾਲ ਤੁਲਨਾ

| ਵਿਸ਼ੇਸ਼ਤਾ | Internationalizer | i18next | Crowdin | ਆਮ LLM |
|---------|------------------|---------|---------|-------------|
| LLM-ਸੰਚਾਲਿਤ ਅਨੁਵਾਦ | ਹਾਂ | ਨਹੀਂ | ਅੰਸ਼ਕ | ਹਾਂ |
| ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਸਟਾਈਲ ਗਾਈਡਾਂ | ਹਾਂ | ਨਹੀਂ | ਨਹੀਂ | ਨਹੀਂ |
| ਗਲੋਸਰੀ ਲਾਗੂਕਰਨ | ਹਾਂ | ਨਹੀਂ | ਹਾਂ | ਨਹੀਂ |
| ਅਨੁਵਾਦ ਮੈਮੋਰੀ | ਹਾਂ | ਨਹੀਂ | ਹਾਂ | ਨਹੀਂ |
| CLI / ਸਥਾਨਕ ਐਗਜ਼ੀਕਿਊਸ਼ਨ | ਹਾਂ | ਲਾਗੂ ਨਹੀਂ | ਨਹੀਂ | ਮੈਨੁਅਲ |
| Git-ਅਨੁਕੂਲ ਫਾਈਲਾਂ | ਹਾਂ | ਹਾਂ | ਅੰਸ਼ਕ | ਮੈਨੁਅਲ |
| ਕੋਈ SaaS ਨਿਰਭਰਤਾ ਨਹੀਂ | ਹਾਂ | ਹਾਂ | ਨਹੀਂ | ਬਦਲਦਾ ਹੈ |
| ਓਪਨ ਸੋਰਸ (AGPL-3.0) | ਹਾਂ | ਹਾਂ | ਨਹੀਂ | ਬਦਲਦਾ ਹੈ |

## ਲਾਇਸੰਸ

[AGPL-3.0](LICENSE)

## ਯੋਗਦਾਨ ਪਾਉਣਾ

ਵਿਕਾਸ ਸੈੱਟਅੱਪ ਅਤੇ ਦਿਸ਼ਾ-ਨਿਰਦੇਸ਼ਾਂ ਲਈ [CONTRIBUTING.md](CONTRIBUTING.md) ਦੇਖੋ। ਸਾਰੇ ਯੋਗਦਾਨਾਂ ਲਈ DCO ਸਾਈਨ-ਆਫ ਦੀ ਲੋੜ ਹੁੰਦੀ ਹੈ。

