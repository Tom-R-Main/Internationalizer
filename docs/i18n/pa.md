> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

ਸਾਫਟਵੇਅਰ ਪ੍ਰੋਜੈਕਟਾਂ ਲਈ AI-native internationalization ਪਾਈਪਲਾਈਨ। LLMs ਦੀ ਵਰਤੋਂ ਕਰਕੇ i18n ਫਾਈਲਾਂ ਦਾ ਅਨੁਵਾਦ ਕਰੋ, ਪ੍ਰਮਾਣਿਤ ਕਰੋ ਅਤੇ ਪ੍ਰਬੰਧਿਤ ਕਰੋ।

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---

<!-- internationalizer:unit markdown:why-internationalizer -->
## Internationalizer ਕਿਉਂ?

ਜ਼ਿਆਦਾਤਰ i18n ਟੂਲ ਜਾਂ ਤਾਂ ਰਨਟਾਈਮ ਲਾਇਬ੍ਰੇਰੀਆਂ (i18next, react-intl) ਹਨ ਜਾਂ key-management SaaS ਪਲੇਟਫਾਰਮ (Crowdin, Lokalise) ਹਨ। ਇਹਨਾਂ ਵਿੱਚੋਂ ਕੋਈ ਵੀ ਅਸਲ ਅਨੁਵਾਦ ਸਮੱਸਿਆ ਨੂੰ ਚੰਗੀ ਤਰ੍ਹਾਂ ਹੱਲ ਨਹੀਂ ਕਰਦਾ:

- **ਮੈਨੁਅਲ ਅਨੁਵਾਦ** ਕੁਝ ਭਾਸ਼ਾਵਾਂ ਤੋਂ ਅੱਗੇ ਸਕੇਲ ਨਹੀਂ ਹੁੰਦਾ
- **ਮਸ਼ੀਨ ਅਨੁਵਾਦ APIs** (Google Translate, DeepL) ਤੁਹਾਡੀ ਸ਼ਬਦਾਵਲੀ, ਟੋਨ ਅਤੇ UI ਪਰੰਪਰਾਵਾਂ ਨੂੰ ਨਜ਼ਰਅੰਦਾਜ਼ ਕਰਦੇ ਹਨ
- **ਆਮ LLM ਅਨੁਵਾਦ** ਬਿਹਤਰ ਕੰਮ ਕਰਦਾ ਹੈ, ਪਰ ਗਲੋਸਰੀਆਂ ਅਤੇ ਸਟਾਈਲ ਗਾਈਡਾਂ ਤੋਂ ਬਿਨਾਂ, ਤੁਹਾਨੂੰ ਅਸੰਗਤ ਨਤੀਜੇ ਮਿਲਦੇ ਹਨ

Internationalizer ਵੱਖਰਾ ਹੈ। ਇਹ ਇੱਕ **CLI ਪਾਈਪਲਾਈਨ** ਹੈ ਜੋ LLM ਅਨੁਵਾਦ ਨੂੰ ਇਹਨਾਂ ਨਾਲ ਜੋੜਦੀ ਹੈ:

- **ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਗਲੋਸਰੀਆਂ** — ਤੁਹਾਡੀ ਐਪ ਵਿੱਚ ਇਕਸਾਰ ਸ਼ਬਦਾਵਲੀ ਲਾਗੂ ਕਰੋ
- **ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਸਟਾਈਲ ਗਾਈਡਾਂ** — ਟੋਨ, ਰਸਮੀਪਣ, ਬਹੁਵਚਨ ਰੂਪਾਂ ਅਤੇ ਟਾਈਪੋਗ੍ਰਾਫੀ ਨੂੰ ਕੰਟਰੋਲ ਕਰੋ
- **ਅਨੁਵਾਦ ਮੈਮੋਰੀ** — ਨਾ ਬਦਲੀਆਂ ਸਟ੍ਰਿੰਗਾਂ ਨੂੰ ਛੱਡੋ, API ਕਾਲਾਂ 'ਤੇ ਪੈਸੇ ਬਚਾਓ
- **ਨਿਰਧਾਰਿਤ ਪ੍ਰਮਾਣੀਕਰਨ** — ਸ਼ਿਪ ਕਰਨ ਤੋਂ ਪਹਿਲਾਂ ਗੁੰਮ ਜਾਂ ਵਾਧੂ ਕੁੰਜੀਆਂ, ਸੁਰੱਖਿਅਤ ਬਣਤਰ ਬਦਲਾਵਾਂ, ਗਲੋਸਰੀ ਸੰਬੰਧੀ ਮੁੱਦਿਆਂ ਅਤੇ ਬਹੁਵਚਨ ਜਾਂ ICU ਗਲਤੀਆਂ ਨੂੰ ਫੜੋ

<!-- internationalizer:unit markdown:installation -->
## ਇੰਸਟਾਲੇਸ਼ਨ

npm ਤੋਂ ਇੰਸਟਾਲ ਕਰੋ:

```bash
npm install -g internationalizer
```

ਜਾਂ ਗਲੋਬਲ ਇੰਸਟਾਲ ਤੋਂ ਬਿਨਾਂ ਚਲਾਓ:

```bash
npx internationalizer --help
```

npm ਪੈਕੇਜ ਪਲੇਟਫਾਰਮ-ਵਿਸ਼ੇਸ਼ ਵਿਕਲਪਿਕ ਨਿਰਭਰਤਾਵਾਂ ਰਾਹੀਂ npm ਤੋਂ ਮੇਲ ਖਾਂਦੀ ਪਹਿਲਾਂ ਤੋਂ ਬਣੀ ਬਾਈਨਰੀ ਨੂੰ ਇੰਸਟਾਲ ਕਰਦਾ ਹੈ।

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

<!-- internationalizer:unit markdown:npm-packages -->
## npm ਪੈਕੇਜ

- Git ਟੈਗ ਅਤੇ npm ਪੈਕੇਜ ਵਰਜਨ ਮੇਲ ਖਾਂਦੇ ਹੋਣੇ ਚਾਹੀਦੇ ਹਨ, ਉਦਾਹਰਨ ਲਈ `v0.1.0` ਅਤੇ `0.1.0`
- ਰੂਟ `internationalizer` ਪੈਕੇਜ ਪਲੇਟਫਾਰਮ ਪੈਕੇਜਾਂ 'ਤੇ ਨਿਰਭਰ ਕਰਦਾ ਹੈ ਜਿਵੇਂ ਕਿ `internationalizer-darwin-arm64`
- ਸਮਰਥਿਤ npm ਟੀਚੇ: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI ਪਬਲਿਸ਼ਿੰਗ ਲਈ `NPM_TOKEN` ਨਾਮਕ GitHub ਸੀਕ੍ਰੇਟ ਦੀ ਲੋੜ ਹੁੰਦੀ ਹੈ

<!-- internationalizer:unit markdown:quick-start -->
## ਤੁਰੰਤ ਸ਼ੁਰੂਆਤ

1. ਆਪਣੇ ਪ੍ਰੋਜੈਕਟ ਰੂਟ ਵਿੱਚ ਇੱਕ ਕੌਂਫਿਗ ਫਾਈਲ ਬਣਾਓ:

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

<!-- internationalizer:unit markdown:commands -->
## ਕਮਾਂਡਾਂ

### `translate`

ਗੁੰਮ ਹੋਈਆਂ ਜਾਂ ਪੁਰਾਣੀਆਂ ਕੁੰਜੀਆਂ ਲੱਭੋ ਅਤੇ ਉਹਨਾਂ ਦਾ LLM ਰਾਹੀਂ ਅਨੁਵਾਦ ਕਰੋ।

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

ਅਨੁਵਾਦ ਸਥਿਤੀ ਸੁਤੰਤਰ ਤੌਰ 'ਤੇ ਗੁੰਮ, ਸਰੋਤ-ਪੁਰਾਣੀਆਂ, ਨੀਤੀ-ਪੁਰਾਣੀਆਂ, ਮੌਜੂਦਾ, ਅਤੇ ਮੈਨੁਅਲ ਤੌਰ 'ਤੇ ਸੰਪਾਦਿਤ ਸਥਿਤੀਆਂ ਦੀ ਰਿਪੋਰਟ ਕਰਦੀ ਹੈ, ਇਸ ਲਈ ਕੋਈ ਮੈਨੁਅਲ ਸੰਪਾਦਨ ਕਿਸੇ ਸਰੋਤ ਜਾਂ ਨੀਤੀ ਤਬਦੀਲੀ ਨੂੰ ਲੁਕਾ ਨਹੀਂ ਸਕਦਾ। ਨੀਤੀ-ਪੁਰਾਣੇ ਮੁੱਲਾਂ ਦੀ ਰਿਪੋਰਟ ਕੀਤੀ ਜਾਂਦੀ ਹੈ ਪਰ ਸਿਰਫ਼ `--refresh-policy` ਨਾਲ ਹੀ ਮੁੜ-ਅਨੁਵਾਦ ਕੀਤਾ ਜਾਂਦਾ ਹੈ। ਮੈਨੁਅਲ ਤੌਰ 'ਤੇ ਸੰਪਾਦਿਤ ਕੀਤੇ ਮੁੱਲ ਕਦੇ ਵੀ ਸਵੈਚਲਿਤ ਤੌਰ 'ਤੇ ਓਵਰਰਾਈਟ ਨਹੀਂ ਹੁੰਦੇ। ਸਮੀਖਿਆ ਕੀਤੇ ਅਨੁਵਾਦਾਂ ਵਿੱਚ ਮੈਨੀਫੈਸਟ ਸ਼ਾਮਲ ਕਰਦੇ ਸਮੇਂ ਜਾਂ ਸਮੀਖਿਆ ਕੀਤੇ ਮੈਨੁਅਲ ਸੰਪਾਦਨ ਨੂੰ ਸਪਸ਼ਟ ਤੌਰ 'ਤੇ ਨਵੇਂ ਬੇਸਲਾਈਨ ਵਜੋਂ ਸਵੀਕਾਰ ਕਰਦੇ ਸਮੇਂ `--adopt-existing` ਦੀ ਵਰਤੋਂ ਕਰੋ।

### `validate`

ਸਾਰੀਆਂ ਲੋਕੇਲ ਫਾਈਲਾਂ ਦੀ ਉਹਨਾਂ ਦੇ ਸਰੋਤ ਬੰਡਲਾਂ ਨਾਲ ਜਾਂਚ ਕਰੋ। ਡਿਫੌਲਟ ਪ੍ਰਮਾਣੀਕਰਨ ਬਣਤਰ ਕਵਰੇਜ ਦੀ ਜਾਂਚ ਕਰਦਾ ਹੈ (ਮੌਜੂਦ ਲੋੜੀਂਦੀਆਂ ਟਾਰਗੇਟ ਕੁੰਜੀਆਂ ਦੀ ਪ੍ਰਤੀਸ਼ਤਤਾ), ਵਾਧੂ ਕੁੰਜੀਆਂ ਨੂੰ ਚੇਤਾਵਨੀਆਂ ਵਜੋਂ ਰਿਪੋਰਟ ਕਰਦਾ ਹੈ, ਅਤੇ ਗੁੰਮ ਕੁੰਜੀਆਂ, ਇੰਟਰਪੋਲੇਸ਼ਨ ਬੇਮੇਲਾਂ, ਜਾਂ ਅਵੈਧ ICU MessageFormat ਬਣਤਰ ਲਈ ਅਸਫਲ ਹੁੰਦਾ ਹੈ।

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` ਅਨੁਵਾਦਿਤ ਕਵਰੇਜ ਦੀ ਵੀ ਰਿਪੋਰਟ ਕਰਦਾ ਹੈ। ਆਪਣੇ ਸਰੋਤ ਦੇ ਸਮਾਨ ਇੱਕ ਭਾਸ਼ਾਈ ਮੁੱਲ ਉਦੋਂ ਤੱਕ ਅਣ-ਅਨੁਵਾਦਿਤ ਰਹਿੰਦਾ ਹੈ ਜਦੋਂ ਤੱਕ ਗਲੋਸਰੀ ਵਿੱਚ ਪੂਰੇ ਮੁੱਲ ਲਈ ਸਪਸ਼ਟ ਤੌਰ 'ਤੇ ਇੱਕੋ ਸਰੋਤ, ਇੱਕੋ ਟਾਰਗੇਟ ਐਂਟਰੀ ਮੌਜੂਦ ਨਾ ਹੋਵੇ; `ignore_case` ਦਾ ਸਤਿਕਾਰ ਕੀਤਾ ਜਾਂਦਾ ਹੈ, ਪਰ ਲੰਬੇ ਮੁੱਲ ਵਿੱਚ ਸ਼ਾਮਲ ਗਲੋਸਰੀ ਸ਼ਬਦ ਛੋਟ ਨਹੀਂ ਹੈ। ਸਖ਼ਤ (Strict) ਮੋਡ ਵਾਧੂ ਕੁੰਜੀਆਂ, ਸਰੋਤ-ਸਮਾਨ ਮੁੱਲਾਂ, ਬਦਲੇ ਹੋਏ ਇੰਟਰਪੋਲੇਸ਼ਨ/HTML/ਕੋਡ/ਮਾਰਕਡਾਊਨ-ਲਿੰਕ ਬਣਤਰ, ਗਲੋਸਰੀ ਉਲੰਘਣਾਵਾਂ, ਅਤੇ ਕੌਂਫਿਗਰ ਕੀਤੇ ਬਹੁਵਚਨ ਰੂਪਾਂ 'ਤੇ ਅਸਫਲ ਹੁੰਦਾ ਹੈ।

`--require-state` `.internationalizer.lock` ਦੇ ਵਿਰੁੱਧ ਹਰੇਕ ਟਾਰਗੇਟ ਦੀ ਪੁਸ਼ਟੀ ਕਰਦਾ ਹੈ। ਇਹ ਉਦੋਂ ਅਸਫਲ ਹੁੰਦਾ ਹੈ ਜਦੋਂ ਕੋਈ ਕੁੰਜੀ ਅਣ-ਟਰੈਕ ਹੁੰਦੀ ਹੈ, ਜਾਂ ਜਦੋਂ ਇਸਦਾ ਰਿਕਾਰਡ ਕੀਤਾ ਸਰੋਤ, ਅਨੁਵਾਦ ਨੀਤੀ, ਜਾਂ ਟਾਰਗੇਟ ਹੈਸ਼ ਪੁਰਾਣਾ ਹੁੰਦਾ ਹੈ। ਇਸਨੂੰ `--strict` ਨਾਲ ਜੋੜਿਆ ਜਾ ਸਕਦਾ ਹੈ।

ਮਨੁੱਖੀ ਅਤੇ JSON ਰਿਪੋਰਟਾਂ ਸਥਿਰ ਖੋਜ ਕੋਡਾਂ ਦੀ ਵਰਤੋਂ ਕਰਦੀਆਂ ਹਨ:

| ਕੋਡ | ਅਰਥ |
| --- | --- |
| `missing_key` / `extra_key` | ਸਰੋਤ ਅਤੇ ਟਾਰਗੇਟ ਕੁੰਜੀ ਸੈੱਟ ਵੱਖਰੇ ਹਨ |
| `blank_translation` | ਇੱਕ ਗੈਰ-ਖਾਲੀ ਸਰੋਤ ਦਾ ਇੱਕ ਖਾਲੀ ਸਖ਼ਤ-ਮੋਡ ਟਾਰਗੇਟ ਹੁੰਦਾ ਹੈ |
| `source_identical` | ਇੱਕ ਸਖ਼ਤ-ਮੋਡ ਭਾਸ਼ਾਈ ਮੁੱਲ ਅਣ-ਅਨੁਵਾਦਿਤ ਰਹਿੰਦਾ ਹੈ |
| `protected_structure_mismatch` | ਇੰਟਰਪੋਲੇਸ਼ਨ, HTML, ਕੋਡ, ਜਾਂ ਲਿੰਕ ਬਣਤਰ ਬਦਲ ਗਈ |
| `glossary_violation` | ਕੋਈ ਪ੍ਰਵਾਨਿਤ ਟਾਰਗੇਟ ਸ਼ਬਦ ਜਾਂ ਰੂਪ ਨਹੀਂ ਮਿਲਿਆ |
| `plural_form_missing` | ਇੱਕ ਕੌਂਫਿਗਰ ਕੀਤਾ ਲੋਕੇਲ ਬਹੁਵਚਨ ਰੂਪ ਗੈਰ-ਹਾਜ਼ਰ ਹੈ |
| `icu_message_syntax` | ਇੱਕ ਸਰੋਤ ਜਾਂ ਟਾਰਗੇਟ ICU ਸੁਨੇਹਾ ਖਰਾਬ ਹੈ |
| `icu_argument_mismatch` | ICU ਆਰਗੂਮੈਂਟ ਦੇ ਨਾਮ, ਕਿਸਮਾਂ, ਜਾਂ ਫਾਰਮੈਟਰ ਸਟਾਈਲ ਵੱਖਰੇ ਹਨ |
| `icu_selector_mismatch` | ਚੋਣਕਾਰ ਵੱਖਰੇ ਹਨ ਜਾਂ ਇੱਕ ਬਹੁਵਚਨ ਸ਼੍ਰੇਣੀ ਟਾਰਗੇਟ ਲੋਕੇਲ ਲਈ ਅਵੈਧ ਹੈ |
| `untracked` | ਟਾਰਗੇਟ ਲਈ ਕੋਈ ਮੈਨੀਫੈਸਟ ਰਿਕਾਰਡ ਮੌਜੂਦ ਨਹੀਂ ਹੈ |
| `source_stale` | ਰਿਕਾਰਡ ਕੀਤੇ ਅਨੁਵਾਦ ਤੋਂ ਬਾਅਦ ਸਰੋਤ ਸਮੱਗਰੀ ਬਦਲ ਗਈ |
| `policy_stale` | ਤਿਆਰ ਕੀਤਾ ਪ੍ਰੋਂਪਟ ਜਾਂ ਮਾਡਲ ਸੈਟਿੰਗਾਂ ਬਦਲ ਗਈਆਂ |
| `target_modified` | ਟਾਰਗੇਟ ਸਮੱਗਰੀ ਮੈਨੀਫੈਸਟ ਰਿਕਾਰਡ ਤੋਂ ਵੱਖਰੀ ਹੈ |

### `detect`

i18n ਫਰੇਮਵਰਕ ਦਾ ਸਵੈ-ਪਤਾ ਲਗਾਓ ਅਤੇ ਇੱਕ ਕੌਂਫਿਗਰੇਸ਼ਨ ਦਾ ਸੁਝਾਅ ਦਿਓ।

```bash
internationalizer detect
```

ਸਮਰਥਨ ਕਰਦਾ ਹੈ: react-i18next, next-intl, vue-i18n, vanilla JSON, ਮਾਰਕਡਾਊਨ ਦਸਤਾਵੇਜ਼।

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
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

<!-- internationalizer:unit markdown:configuration-reference -->
## ਕੌਂਫਿਗਰੇਸ਼ਨ ਹਵਾਲਾ

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

ਲੋਕੇਲ ਪਛਾਣਕਰਤਾ ਚੰਗੀ ਤਰ੍ਹਾਂ ਬਣੇ BCP 47 ਟੈਗ ਹੋਣੇ ਚਾਹੀਦੇ ਹਨ ਜਿਵੇਂ ਕਿ `fr`, `pt-BR`, ਜਾਂ `sr-Latn-RS`। ਪ੍ਰਮਾਣਿਕ-ਬਰਾਬਰ ਟਾਰਗੇਟ ਲੋਕੇਲ ਡੁਪਲੀਕੇਟ ਵਜੋਂ ਰੱਦ ਕੀਤੇ ਜਾਂਦੇ ਹਨ, ਅਤੇ ਲੋਕੇਲ-ਵਿਸ਼ੇਸ਼ ਪ੍ਰਦਾਤਾ ਓਵਰਰਾਈਡ ਪ੍ਰਮਾਣਿਕ-ਬਰਾਬਰ ਸਪੈਲਿੰਗ ਨਾਲ ਮੇਲ ਖਾਂਦੇ ਹਨ। ਉਪਰੋਕਤ ਉਦਾਹਰਨ ਵਿੱਚ, ਬਿਨਾਂ ਕਿਸੇ ਓਵਰਰਾਈਡ ਵਾਲੇ ਲੋਕੇਲ—ਜਾਪਾਨੀ ਸਮੇਤ—ਗਲੋਬਲ Gemini ਕੌਂਫਿਗਰੇਸ਼ਨ ਨੂੰ ਵਿਰਾਸਤ ਵਿੱਚ ਪ੍ਰਾਪਤ ਕਰਦੇ ਹਨ।

ICU MessageFormat ਮੁੱਲਾਂ ਨੂੰ ਬਣਤਰ ਦੇ ਅਧਾਰ 'ਤੇ ਪਾਰਸ ਕੀਤਾ ਜਾਂਦਾ ਹੈ। ਸਰਲ ਆਰਗੂਮੈਂਟ, `select`, `plural`, `selectordinal`, `number`, `date`, ਅਤੇ `time` ਸਮਰਥਿਤ ਹਨ, ਜਿਸ ਵਿੱਚ ਨੇਸਟਡ ਸੁਨੇਹੇ, ਬਹੁਵਚਨ ਆਫਸੈੱਟ, ਸਹੀ-ਸੰਖਿਆ ਚੋਣਕਾਰ, ਅਤੇ `#` ਸ਼ਾਮਲ ਹਨ। ਪ੍ਰਮਾਣੀਕਰਨ ਸਿੰਟੈਕਸ, ਆਰਗੂਮੈਂਟ ਕਿਸਮਾਂ ਅਤੇ ਫਾਰਮੈਟਰ ਸਟਾਈਲ, ਬਹੁਵਚਨ ਆਫਸੈੱਟ, ਚੋਣ ਸ਼ਾਖਾ ਪਛਾਣ, ਅਤੇ ਟਾਰਗੇਟ-ਲੋਕੇਲ CLDR ਬਹੁਵਚਨ ਸ਼੍ਰੇਣੀਆਂ ਦੀ ਜਾਂਚ ਕਰਦਾ ਹੈ। ਪ੍ਰਦਾਤਾ ਆਉਟਪੁੱਟ ਜੋ ਇਹਨਾਂ ਨਿਯਮਾਂ ਨੂੰ ਤੋੜਦੀ ਹੈ, ਇੱਕ ਲੋਕੇਲ ਫਾਈਲ ਜਾਂ ਅਨੁਵਾਦ-ਮੈਮੋਰੀ ਰਿਕਾਰਡ ਲਿਖੇ ਜਾਣ ਤੋਂ ਪਹਿਲਾਂ ਰੱਦ ਕਰ ਦਿੱਤੀ ਜਾਂਦੀ ਹੈ।

`i18next-v4` ਦੇ ਨਾਲ, ਮਾਨਤਾ ਪ੍ਰਾਪਤ ਸਰੋਤ ਬਹੁਵਚਨ ਪਰਿਵਾਰਾਂ ਦਾ ਅਨੁਵਾਦ ਦੌਰਾਨ ਟਾਰਗੇਟ ਲੋਕੇਲ ਦੀਆਂ CLDR ਸ਼੍ਰੇਣੀਆਂ ਵਿੱਚ ਵਿਸਤਾਰ ਕੀਤਾ ਜਾਂਦਾ ਹੈ। ਇੱਕ ਸਿਰਫ਼-ਟਾਰਗੇਟ ਸ਼੍ਰੇਣੀ ਸਰੋਤ ਪਰਿਵਾਰ ਦੇ `_other` ਮੁੱਲ ਨੂੰ ਆਪਣੇ ਅਨੁਵਾਦ ਟੈਂਪਲੇਟ ਵਜੋਂ ਵਰਤਦੀ ਹੈ। ਸਖ਼ਤ ਪ੍ਰਮਾਣੀਕਰਨ ਲਈ ਉਹਨਾਂ ਟਾਰਗੇਟ ਸ਼੍ਰੇਣੀਆਂ ਦੀ ਲੋੜ ਹੁੰਦੀ ਹੈ; ਸਿਰਫ਼-ਸਰੋਤ ਸ਼੍ਰੇਣੀਆਂ ਉਹਨਾਂ ਟਾਰਗੇਟ ਲੋਕੇਲਾਂ ਲਈ ਵਿਕਲਪਿਕ ਹਨ ਜੋ ਉਹਨਾਂ ਦੀ ਵਰਤੋਂ ਨਹੀਂ ਕਰਦੇ ਹਨ।

<!-- internationalizer:unit markdown:style-guides -->
## ਸਟਾਈਲ ਗਾਈਡਾਂ

ਸਟਾਈਲ ਗਾਈਡਾਂ ਮਾਰਕਡਾਊਨ ਫਾਈਲਾਂ ਹਨ ਜੋ LLM ਅਨੁਵਾਦ ਪ੍ਰੋਂਪਟ ਵਿੱਚ ਸ਼ਾਮਲ ਕੀਤੀਆਂ ਜਾਂਦੀਆਂ ਹਨ। ਉਹ ਟੋਨ, ਰਸਮੀਪਣ, ਟਾਈਪੋਗ੍ਰਾਫੀ, ਅਤੇ ਹੋਰ ਭਾਸ਼ਾ-ਵਿਸ਼ੇਸ਼ ਪਰੰਪਰਾਵਾਂ ਨੂੰ ਕੰਟਰੋਲ ਕਰਦੀਆਂ ਹਨ।

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### ਸਾਂਝੀਆਂ ਪਰੰਪਰਾਵਾਂ (`_conventions.md`)

ਉਹ ਨਿਯਮ ਪਰਿਭਾਸ਼ਿਤ ਕਰੋ ਜੋ ਸਾਰੀਆਂ ਭਾਸ਼ਾਵਾਂ 'ਤੇ ਲਾਗੂ ਹੁੰਦੇ ਹਨ: ਇੰਟਰਪੋਲੇਸ਼ਨ ਸਿੰਟੈਕਸ, HTML ਸੰਭਾਲ, ਸਟ੍ਰਿੰਗ ਕਿਸਮ ਦੀਆਂ ਪਰੰਪਰਾਵਾਂ (ਬਟਨ ਬਨਾਮ ਲੇਬਲ ਬਨਾਮ ਤਰੁੱਟੀਆਂ), ਆਦਿ।

### ਪ੍ਰਤੀ-ਭਾਸ਼ਾ ਗਾਈਡਾਂ (`{locale}.md`)

ਭਾਸ਼ਾ-ਵਿਸ਼ੇਸ਼ ਨਿਯਮ ਪਰਿਭਾਸ਼ਿਤ ਕਰੋ: ਰਸਮੀਪਣ ਰਜਿਸਟਰ (tu ਬਨਾਮ vous), ਵਿਸ਼ਰਾਮ ਚਿੰਨ੍ਹ (guillemets, ਉਲਟੇ ਪ੍ਰਸ਼ਨ ਚਿੰਨ੍ਹ), ਬਹੁਵਚਨ ਰੂਪ, ਮਿਤੀ/ਸੰਖਿਆ ਫਾਰਮੈਟਿੰਗ, ਅਤੇ ਇੱਕ ਸ਼ਬਦਾਵਲੀ ਗਲੋਸਰੀ।

ਸਟਾਈਲ ਗਾਈਡਾਂ ਸਥਾਈ ਨੀਤੀ ਇਨਪੁਟ ਹਨ, ਨਾ ਕਿ ਤਿਆਰ ਕੀਤਾ ਆਉਟਪੁੱਟ। Internationalizer ਉਹਨਾਂ ਨੂੰ ਪੜ੍ਹਦਾ ਹੈ ਪਰ ਕਦੇ ਵੀ ਮੁੜ ਨਹੀਂ ਲਿਖਦਾ। ਉਹਨਾਂ ਦੀ ਸਮੱਗਰੀ ਨੂੰ ਗਲੋਸਰੀ ਅਤੇ ਪ੍ਰੋਂਪਟ ਕੰਟਰੈਕਟ ਤੋਂ ਵੱਖਰੇ ਤੌਰ 'ਤੇ ਹੈਸ਼ ਕੀਤਾ ਜਾਂਦਾ ਹੈ, ਤਾਂ ਜੋ ਇੱਕ ਐਪਲੀਕੇਸ਼ਨ ਕੋਡ ਤਬਦੀਲੀ ਅਨੁਵਾਦ ਨੂੰ ਪੁਰਾਣਾ ਨਾ ਬਣਾਏ। ਕਿਸੇ ਗਾਈਡ ਨੂੰ ਸੰਪਾਦਿਤ ਕਰਨਾ ਜਾਣਬੁੱਝ ਕੇ ਉਸ ਲੋਕੇਲ ਨੂੰ ਨੀਤੀ ਸਮੀਖਿਆ ਲਈ ਚਿੰਨ੍ਹਿਤ ਕਰਦਾ ਹੈ; ਅੰਦਰੂਨੀ ਪ੍ਰੋਂਪਟ ਸ਼ਬਦਾਵਲੀ ਨੂੰ ਬਦਲਣਾ ਅਜਿਹਾ ਨਹੀਂ ਕਰਦਾ, ਜਦੋਂ ਤੱਕ ਪ੍ਰੋਂਪਟ ਕੰਟਰੈਕਟ ਵਰਜਨ ਵੀ ਨਹੀਂ ਬਦਲਦਾ।

ਕਾਰਜਸ਼ੀਲ ਉਦਾਹਰਨ ਲਈ [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) ਦੇਖੋ।

<!-- internationalizer:unit markdown:glossary-format -->
## ਗਲੋਸਰੀ ਫਾਰਮੈਟ

ਗਲੋਸਰੀ ਫਾਈਲਾਂ `{glossary_dir}/{locale}.json` ਵਿੱਚ ਸਟੋਰ ਕੀਤੀਆਂ JSON ਐਰੇ ਹਨ:

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

`variants` ਹੋਰ ਪ੍ਰਵਾਨਿਤ ਟਾਰਗੇਟ ਰੂਪਾਂ ਨੂੰ ਸੂਚੀਬੱਧ ਕਰਦਾ ਹੈ। `enforcement` ਸ਼ਾਇਦ `error`, `warning` ਹੋ ਸਕਦਾ ਹੈ, ਜਾਂ ਡਿਫੌਲਟ ਤਰੁੱਟੀ ਵਿਵਹਾਰ ਲਈ ਛੱਡਿਆ ਜਾ ਸਕਦਾ ਹੈ। ਸ਼ਬਦਾਂ ਨੂੰ LLM ਪ੍ਰੋਂਪਟ ਵਿੱਚ ਇੱਕ ਸ਼ਬਦਾਵਲੀ ਸਾਰਣੀ ਵਜੋਂ ਸ਼ਾਮਲ ਕੀਤਾ ਜਾਂਦਾ ਹੈ, ਜੋ ਤੁਹਾਡੀ ਐਪਲੀਕੇਸ਼ਨ ਵਿੱਚ ਇਕਸਾਰ ਅਨੁਵਾਦ ਨੂੰ ਯਕੀਨੀ ਬਣਾਉਂਦਾ ਹੈ। ਇੱਕ ਸਹੀ ਐਂਟਰੀ ਜਿਵੇਂ ਕਿ `{"source":"API","target":"API"}` ਉਸ ਪੂਰੇ ਸਰੋਤ-ਸਮਾਨ ਮੁੱਲ ਨੂੰ ਸਖ਼ਤ ਅਣ-ਅਨੁਵਾਦਿਤ ਮੁੱਲ ਖੋਜਾਂ ਤੋਂ ਵੀ ਛੋਟ ਦਿੰਦੀ ਹੈ; ਇਹ ਕੇਵਲ `API` ਵਾਲੇ ਇੱਕ ਲੰਬੇ ਮੁੱਲ ਨੂੰ ਛੋਟ ਨਹੀਂ ਦਿੰਦੀ।

<!-- internationalizer:unit markdown:translation-memory -->
## ਅਨੁਵਾਦ ਮੈਮੋਰੀ

ਅਨੁਵਾਦ ਮੈਮੋਰੀ ਨੂੰ ਇੱਕ JSONL ਫਾਈਲ (ਪ੍ਰਤੀ ਲਾਈਨ ਇੱਕ JSON ਰਿਕਾਰਡ) ਵਜੋਂ ਸਟੋਰ ਕੀਤਾ ਜਾਂਦਾ ਹੈ। ਹਰੇਕ ਰਿਕਾਰਡ ਵਿੱਚ ਸ਼ਾਮਲ ਹਨ:

- ਬੰਡਲ, ਕੁੰਜੀ, ਸਰੋਤ ਮੁੱਲ, ਅਨੁਵਾਦਿਤ ਮੁੱਲ, ਅਤੇ ਪ੍ਰਮਾਣਿਕ ਟਾਰਗੇਟ ਲੋਕੇਲ
- ਸਰੋਤ, ਸਟਾਈਲ-ਗਾਈਡ, ਗਲੋਸਰੀ, ਪ੍ਰੋਂਪਟ-ਕੰਟਰੈਕਟ, ਅਤੇ ਸੰਯੁਕਤ ਨੀਤੀ ਹੈਸ਼
- ਪ੍ਰਦਾਤਾ ਅਤੇ ਮਾਡਲ ਜਿਸਨੇ ਅਨੁਵਾਦ ਤਿਆਰ ਕੀਤਾ
- ਇੱਕ ਟਾਈਮਸਟੈਂਪ

ਬਾਅਦ ਦੀਆਂ ਦੌੜਾਂ 'ਤੇ, ਇੱਕੋ ਸਰੋਤ ਅਤੇ ਨੀਤੀ ਹੈਸ਼ ਵਾਲੀਆਂ ਸਟ੍ਰਿੰਗਾਂ ਨੂੰ LLM ਨੂੰ ਕਾਲ ਕੀਤੇ ਬਿਨਾਂ ਕੈਸ਼ ਤੋਂ ਸਰਵ ਕੀਤਾ ਜਾਂਦਾ ਹੈ। ਡਿਫੌਲਟ ਮਾਰਗ ਅਣਡਿੱਠ ਕੀਤੀ ਗਈ `.internationalizer/` ਡਾਇਰੈਕਟਰੀ ਅਧੀਨ ਹੈ, ਇਸ ਲਈ ਇਹ ਇੱਕ ਸਥਾਨਕ ਕੈਸ਼ ਰਹਿੰਦਾ ਹੈ। ਜੇਕਰ ਤੁਹਾਡਾ ਪ੍ਰੋਜੈਕਟ ਜਾਣਬੁੱਝ ਕੇ ਅਨੁਵਾਦ ਮੈਮੋਰੀ ਸਾਂਝੀ ਕਰਦਾ ਹੈ ਤਾਂ `tm_path` ਨੂੰ ਇੱਕ ਟ੍ਰੈਕ ਕੀਤੇ ਸਥਾਨ 'ਤੇ ਸੈੱਟ ਕਰੋ। ਸਮੀਖਿਆ ਯੋਗ `.internationalizer.lock` ਮੈਨੀਫੈਸਟ ਨੂੰ ਵੱਖਰੇ ਤੌਰ 'ਤੇ ਵਰਜਨ ਕੀਤਾ ਗਿਆ ਹੈ।

<!-- internationalizer:unit markdown:supported-formats -->
## ਸਮਰਥਿਤ ਫਾਰਮੈਟ

| ਫਾਰਮੈਟ | ਐਕਸਟੈਂਸ਼ਨਾਂ | ਮੋਡ |
|--------|-----------|------|
| JSON | `.json` | Key-value (ਨੇਸਟਡ, ਡਾਟ-ਨੋਟੇਸ਼ਨ ਫਲੈਟਨਡ) |
| YAML | `.yml`, `.yaml` | Key-value (ਟਿੱਪਣੀਆਂ ਅਤੇ ਕ੍ਰਮ ਨੂੰ ਸੁਰੱਖਿਅਤ ਰੱਖਦਾ ਹੈ) |
| Markdown | `.md`, `.mdx` | ਪ੍ਰੀਐਂਬਲ ਅਤੇ H2-ਪੱਧਰ ਦੇ ਭਾਗ |

ਮਾਰਕਡਾਊਨ ਟਾਰਗੇਟਾਂ ਵਿੱਚ H2 ਭਾਗਾਂ ਤੋਂ ਪਹਿਲਾਂ ਅਦਿੱਖ `internationalizer:unit` ਟਿੱਪਣੀਆਂ ਹੁੰਦੀਆਂ ਹਨ। ਇਹ ਸਥਿਰ ਮਾਰਕਰ Internationalizer ਨੂੰ ਗੈਰ-ਸੰਬੰਧਿਤ ਭਾਗਾਂ ਦਾ ਮੁੜ-ਅਨੁਵਾਦ ਕੀਤੇ ਬਿਨਾਂ ਇੱਕ ਸਰੋਤ ਭਾਗ ਨੂੰ ਜੋੜਨ, ਲਿਜਾਣ ਜਾਂ ਸੰਪਾਦਿਤ ਕਰਨ ਦਿੰਦੇ ਹਨ। ਮੌਜੂਦਾ ਅਣ-ਚਿੰਨ੍ਹਿਤ ਦਸਤਾਵੇਜ਼ ਉਹਨਾਂ ਦੇ ਅਗਲੇ ਸਫਲ ਅੱਪਡੇਟ 'ਤੇ ਮਾਰਕਰ ਪ੍ਰਾਪਤ ਕਰਦੇ ਹਨ।

<!-- internationalizer:unit markdown:project-type-detection -->
## ਪ੍ਰੋਜੈਕਟ ਕਿਸਮ ਦੀ ਪਛਾਣ

`internationalizer detect` ਇਹਨਾਂ ਦੀ ਜਾਂਚ ਕਰਕੇ ਤੁਹਾਡੇ i18n ਸੈੱਟਅੱਪ ਦੀ ਪਛਾਣ ਕਰਦਾ ਹੈ:

- react-i18next, next-intl, ਜਾਂ vue-i18n ਲਈ `package.json` ਨਿਰਭਰਤਾਵਾਂ
- ਆਮ ਲੋਕੇਲ ਪੈਟਰਨਾਂ ਨਾਲ ਮੇਲ ਖਾਂਦੀਆਂ ਡਾਇਰੈਕਟਰੀ ਬਣਤਰਾਂ
- ਫਾਈਲ ਐਕਸਟੈਂਸ਼ਨਾਂ ਅਤੇ ਨਾਮਕਰਨ ਪਰੰਪਰਾਵਾਂ

<!-- internationalizer:unit markdown:architecture -->
## ਆਰਕੀਟੈਕਚਰ

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

<!-- internationalizer:unit markdown:license -->
## ਲਾਇਸੰਸ

[AGPL-3.0](../../LICENSE)

ਨਿਰਭਰਤਾ ਨੋਟਿਸਾਂ ਲਈ [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) ਦੇਖੋ।

<!-- internationalizer:unit markdown:contributing -->
## ਯੋਗਦਾਨ ਪਾਉਣਾ

ਵਿਕਾਸ ਸੈੱਟਅੱਪ ਅਤੇ ਦਿਸ਼ਾ-ਨਿਰਦੇਸ਼ਾਂ ਲਈ [CONTRIBUTING.md](../../CONTRIBUTING.md) ਦੇਖੋ। ਸਾਰੇ ਯੋਗਦਾਨਾਂ ਲਈ DCO ਸਾਈਨ-ਆਫ ਦੀ ਲੋੜ ਹੁੰਦੀ ਹੈ।
