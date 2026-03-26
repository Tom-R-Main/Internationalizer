> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Saluran paip pengantarabangsaan natif AI untuk projek perisian. Terjemah, sahkan, dan urus fail i18n menggunakan LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Mengapa Internationalizer?

Kebanyakan alat i18n sama ada pustaka masa jalanan (i18next, react-intl) atau platform SaaS pengurusan kunci (Crowdin, Lokalise). Tiada satu pun daripadanya menyelesaikan masalah terjemahan sebenar dengan baik:

- **Terjemahan manual** tidak berskala melepasi beberapa bahasa
- **API terjemahan mesin** (Google Translate, DeepL) mengabaikan terminologi, nada, dan resam UI anda
- **Terjemahan LLM generik** berfungsi lebih baik, tetapi tanpa glosari dan panduan gaya, anda mendapat hasil yang tidak konsisten

Internationalizer adalah berbeza. Ia merupakan **saluran paip CLI** yang menggabungkan terjemahan LLM dengan:

- **Glosari mengikut bahasa** — menguatkuasakan terminologi yang konsisten merentas aplikasi anda
- **Panduan gaya mengikut bahasa** — mengawal nada, formaliti, pempluralan, dan tipografi
- **Memori terjemahan** — langkau rentetan yang tidak berubah, jimatkan wang untuk panggilan API
- **Pengesahan kunci** — tangkap terjemahan yang hilang dan ketidakpadanan interpolasi sebelum ia dilancarkan

## Pemasangan

Pasang daripada npm:

```bash
npm install -g internationalizer
```

Atau jalankan tanpa pemasangan global:

```bash
npx internationalizer --help
```

Pakej npm memasang binari prapelbina yang sepadan daripada npm melalui kebergantungan pilihan khusus platform.

Pasang dengan Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Atau bina daripada sumber:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## Pakej npm

- Tag Git dan versi pakej npm mesti sepadan, contohnya `v0.1.0` dan `0.1.0`
- Pakej `internationalizer` akar bergantung pada pakej platform seperti `internationalizer-darwin-arm64`
- Sasaran npm yang disokong: macOS arm64/x64, Linux arm64/x64, Windows x64
- Penerbitan CI memerlukan rahsia GitHub bernama `NPM_TOKEN`

## Mula Pantas

1. Cipta fail konfigurasi dalam akar projek anda:

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

2. Tetapkan kunci API anda:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Pratonton apa yang akan diterjemahkan:

```bash
internationalizer translate --dry-run
```

4. Jalankan terjemahan:

```bash
internationalizer translate
```

5. Sahkan semua lokal:

```bash
internationalizer validate
```

## Arahan

### `translate`

Cari kunci yang hilang dan terjemahkannya melalui LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Periksa semua fail lokal untuk kunci yang hilang, kunci berlebihan, dan ketidakpadanan interpolasi.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

Kesan secara automatik rangka kerja i18n dan cadangkan konfigurasi.

```bash
internationalizer detect
```

Menyokong: react-i18next, next-intl, vue-i18n, vanilla JSON, dokumen markdown.

### `glossary`

Urus terma glosari mengikut bahasa yang dikuatkuasakan semasa terjemahan.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Urus memori terjemahan (cache JSONL untuk rentetan yang telah diterjemahkan sebelum ini).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Rujukan Konfigurasi

```yaml
# .internationalizer.yml

# Source language (default: en)
source_locale: en

# Languages to translate into (required)
target_locales: [fr, de, es, ja, zh-CN, ar]

# Path to the source locale file (required)
source_path: locales/en.json

# LLM provider settings
llm:
  # Provider: "anthropic", "openai", "gemini", or "openrouter" (default: gemini)
  provider: gemini

  # Model name defaults by provider:
  #   anthropic:  claude-sonnet-4-6
  #   openai:     gpt-5.4
  #   gemini:     gemini-3.1-pro-preview
  #   openrouter: google/gemini-3-flash-preview
  model: gemini-3.1-pro-preview

  # Environment variable containing the API key
  api_key_env: GOOGLE_AI_STUDIO_API_KEY

  # Base URL for OpenAI-compatible endpoints (optional)
  # base_url: https://api.openai.com

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
```

## Panduan Gaya

Panduan gaya ialah fail Markdown yang disuntik ke dalam gesaan terjemahan LLM. Ia mengawal nada, formaliti, tipografi, dan resam khusus bahasa yang lain.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Resam dikongsi (`_conventions.md`)

Tentukan peraturan yang terpakai untuk semua bahasa: sintaks interpolasi, pengekalan HTML, resam jenis rentetan (butang lwn. label lwn. ralat), dsb.

### Panduan mengikut bahasa (`{locale}.md`)

Tentukan peraturan khusus bahasa: laras bahasa formaliti (tu lwn. vous), tanda baca (guillemet, tanda soal terbalik), bentuk jamak, pemformatan tarikh/nombor, dan glosari terminologi.

Lihat [`examples/react-app/style-guides/`](examples/react-app/style-guides/) untuk contoh yang berfungsi.

## Format Glosari

Fail glosari ialah tatasusunan JSON yang disimpan dalam `{glossary_dir}/{locale}.json`:

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

Terma disuntik ke dalam gesaan LLM sebagai jadual terminologi, memastikan terjemahan terma utama yang konsisten merentas aplikasi anda.

## Memori Terjemahan

Memori terjemahan disimpan sebagai fail JSONL (satu rekod JSON setiap baris). Setiap rekod mengandungi:

- Kunci dan nilai sumber
- Nilai yang diterjemahkan
- Cincangan SHA-256 bagi nilai sumber
- Cap masa

Pada larian seterusnya, rentetan yang tidak berubah disajikan daripada cache TM tanpa memanggil LLM, menjimatkan masa dan kos API. Fail TM mesra git dan boleh dikomit bersama fail lokal anda.

## Format Disokong

| Format | Sambungan | Mod |
|--------|-----------|------|
| JSON | `.json` | Nilai kunci (bersarang, diratakan dengan tatatanda titik) |
| YAML | `.yml`, `.yaml` | Nilai kunci (mengekalkan ulasan dan susunan) |
| Markdown | `.md`, `.mdx` | Terjemahan keseluruhan dokumen |

## Pengesanan Jenis Projek

`internationalizer detect` mengenal pasti persediaan i18n anda dengan memeriksa:

- Kebergantungan `package.json` untuk react-i18next, next-intl, atau vue-i18n
- Struktur direktori yang sepadan dengan corak lokal biasa
- Sambungan fail dan resam penamaan

## Seni Bina

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
  styleguide/              Style guide loader
  tm/                      JSONL translation memory
  translate/               Translation orchestrator
  validate/                Locale validation and diffing
```

## Perbandingan dengan Alternatif

| Ciri | Internationalizer | i18next | Crowdin | LLM Generik |
|---------|------------------|---------|---------|-------------|
| Terjemahan dikuasakan LLM | Ya | Tidak | Separa | Ya |
| Panduan gaya mengikut bahasa | Ya | Tidak | Tidak | Tidak |
| Penguatkuasaan glosari | Ya | Tidak | Ya | Tidak |
| Memori terjemahan | Ya | Tidak | Ya | Tidak |
| Pelaksanaan CLI / tempatan | Ya | N/A | Tidak | Manual |
| Fail mesra Git | Ya | Ya | Separa | Manual |
| Tiada kebergantungan SaaS | Ya | Ya | Tidak | Berbeza-beza |
| Sumber terbuka (AGPL-3.0) | Ya | Ya | Tidak | Berbeza-beza |

## Lesen

[AGPL-3.0](LICENSE)

## Menyumbang

Lihat [CONTRIBUTING.md](CONTRIBUTING.md) untuk persediaan pembangunan dan garis panduan. Semua sumbangan memerlukan persetujuan DCO.

