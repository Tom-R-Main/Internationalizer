> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Saluran paip pengantarabangsaan natif AI untuk projek perisian. Terjemah, sahkan, dan urus fail i18n menggunakan LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Mengapa Internationalizer?

Kebanyakan alat i18n sama ada pustaka masa jalanan (i18next, react-intl) atau platform SaaS pengurusan kunci (Crowdin, Lokalise). Tiada satu pun daripadanya menyelesaikan masalah terjemahan sebenar dengan baik:

- **Terjemahan manual** tidak berskala melepasi beberapa bahasa
- **API terjemahan mesin** (Google Translate, DeepL) mengabaikan terminologi, nada, dan konvensyen UI anda
- **Terjemahan LLM generik** berfungsi lebih baik, tetapi tanpa glosari dan panduan gaya, anda mendapat hasil yang tidak konsisten

Internationalizer adalah berbeza. Ia merupakan **saluran paip CLI** yang menggabungkan terjemahan LLM dengan:

- **Glosari mengikut bahasa** — kuat kuasakan terminologi yang konsisten merentas aplikasi anda
- **Panduan gaya mengikut bahasa** — kawal nada, formaliti, pempluralan, dan tipografi
- **Memori terjemahan** — langkau rentetan yang tidak berubah, jimatkan kos panggilan API
- **Pengesahan berketentuan** — tangkap kunci yang hilang atau berlebihan, anjakan struktur dilindungi, isu glosari, serta ralat jamak atau ICU sebelum ia dilancarkan
<!-- internationalizer:unit markdown:installation -->
## Pemasangan

Pasang daripada npm:

```bash
npm install -g internationalizer
```

Atau jalankan tanpa pemasangan global:

```bash
npx internationalizer --help
```

Pakej npm memasang binari prabina yang sepadan daripada npm melalui kebergantungan pilihan khusus platform.

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
<!-- internationalizer:unit markdown:npm-packages -->
## Pakej npm

- Tag Git dan versi pakej npm mesti sepadan, contohnya `v0.1.0` dan `0.1.0`
- Pakej `internationalizer` akar bergantung pada pakej platform seperti `internationalizer-darwin-arm64`
- Sasaran npm yang disokong: macOS arm64/x64, Linux arm64/x64, Windows x64
- Penerbitan CI memerlukan rahsia GitHub bernama `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Permulaan Pantas

1. Cipta fail konfigurasi dalam akar projek anda:

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

2. Tetapkan kunci API anda:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Pratonton perkara yang akan diterjemahkan:

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
<!-- internationalizer:unit markdown:commands -->
## Perintah

### `translate`

Cari kunci yang hilang atau basi dan terjemahkannya melalui LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Keadaan terjemahan melaporkan keadaan hilang, sumber basi, dasar basi, semasa, dan diedit secara manual secara berasingan, supaya pengeditan manual tidak dapat menyembunyikan perubahan sumber atau dasar. Nilai dasar basi dilaporkan tetapi hanya diterjemahkan semula dengan `--refresh-policy`. Nilai yang diedit secara manual tidak akan ditimpa secara automatik. Gunakan `--adopt-existing` apabila memperkenalkan manifes kepada terjemahan yang telah disemak atau apabila menerima pengeditan manual yang disemak secara jelas sebagai garis dasar baharu.

### `validate`

Periksa semua fail lokal terhadap himpunan sumbernya. Pengesahan lalai memeriksa liputan struktur (peratusan kunci sasaran yang diperlukan hadir), melaporkan kunci tambahan sebagai amaran, dan gagal bagi kunci yang hilang, ketidakpadanan interpolasi, atau struktur ICU MessageFormat yang tidak sah.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` juga melaporkan liputan diterjemahkan. Nilai linguistik yang serupa dengan sumbernya dianggap belum diterjemahkan melainkan glosari secara jelas mengandungi entri sumber yang sama dan sasaran yang sama yang tepat bagi keseluruhan nilai; `ignore_case` dipatuhi, tetapi istilah glosari yang tertanam dalam nilai yang lebih panjang bukan pengecualian. Mod ketat gagal sekiranya terdapat kunci tambahan, nilai yang sama dengan sumber, perubahan struktur interpolasi/HTML/kod/pautan Markdown, pelanggaran glosari, dan bentuk jamak yang dikonfigurasikan.

`--require-state` mengesahkan setiap sasaran terhadap `.internationalizer.lock`. Ia gagal apabila kunci tidak dijejaki, atau apabila sumber yang direkodkan, dasar terjemahan, atau cincangan sasaran telah basi. Pilihan ini boleh digabungkan dengan `--strict`.

Laporan boleh baca manusia dan JSON menggunakan kod penemuan yang stabil:

| Kod | Maksud |
| --- | --- |
| `missing_key` / `extra_key` | Set kunci sumber dan sasaran berbeza |
| `blank_translation` | Sumber tidak kosong mempunyai sasaran mod ketat yang kosong |
| `source_identical` | Nilai linguistik mod ketat kekal tidak diterjemahkan |
| `protected_structure_mismatch` | Struktur interpolasi, HTML, kod, atau pautan telah berubah |
| `glossary_violation` | Tiada istilah sasaran atau varian yang diluluskan dijumpai |
| `plural_form_missing` | Bentuk jamak lokal yang dikonfigurasikan tidak wujud |
| `icu_message_syntax` | Mesej ICU sumber atau sasaran tidak terbentuk dengan betul |
| `icu_argument_mismatch` | Nama argumen, jenis, atau gaya pemformat ICU berbeza |
| `icu_selector_mismatch` | Pemilih berbeza atau kategori jamak tidak sah untuk lokal sasaran |
| `untracked` | Tiada rekod manifes wujud untuk sasaran |
| `source_stale` | Kandungan sumber berubah selepas terjemahan direkodkan |
| `policy_stale` | Gesaan yang dijana atau tetapan model telah berubah |
| `target_modified` | Kandungan sasaran berbeza daripada rekod manifes |

### `detect`

Kesan rangka kerja i18n secara automatik dan cadangkan konfigurasi.

```bash
internationalizer detect
```

Menyokong: react-i18next, next-intl, vue-i18n, JSON biasa, dokumen markdown.

### `glossary`

Urus istilah glosari mengikut bahasa yang dikuatkuasakan semasa terjemahan.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Urus memori terjemahan (cache JSONL bagi rentetan yang telah diterjemahkan sebelum ini).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Rujukan Konfigurasi

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

Pengecam lokal mestilah tag BCP 47 yang dibentuk dengan betul seperti `fr`, `pt-BR`, atau `sr-Latn-RS`. Lokal sasaran yang setara secara kanonik akan ditolak sebagai pendua, dan penggantian penyedia khusus lokal dipadankan mengikut ejaan setara secara kanonik. Dalam contoh di atas, lokal tanpa penggantian—termasuk bahasa Jepun—mewarisi konfigurasi Gemini global.

Nilai ICU MessageFormat dihuraikan secara struktur. Argumen mudah, `select`, `plural`, `selectordinal`, `number`, `date`, dan `time` disokong, termasuk mesej bersarang, ofset jamak, pemilih nombor tepat, dan `#`. Pengesahan memeriksa sintaks, jenis argumen dan gaya pemformat, ofset jamak, identiti cabang pemilih, serta kategori jamak CLDR lokal sasaran. Output penyedia yang melanggar ketakberubahan ini akan ditolak sebelum fail lokal atau rekod memori terjemahan ditulis.

Dengan `i18next-v4`, keluarga jamak sumber yang dikenal pasti dikembangkan semasa terjemahan kepada kategori CLDR lokal sasaran. Kategori khusus sasaran menggunakan nilai `_other` keluarga sumber sebagai templat terjemahannya. Pengesahan ketat memerlukan kategori sasaran tersebut; kategori khusus sumber adalah pilihan bagi lokal sasaran yang tidak menggunakannya.
<!-- internationalizer:unit markdown:style-guides -->
## Panduan Gaya

Panduan gaya ialah fail Markdown yang disuntik ke dalam gesaan terjemahan LLM. Panduan ini mengawal nada, formaliti, tipografi, dan konvensyen khusus bahasa yang lain.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Konvensyen bersama (`_conventions.md`)

Menentukan peraturan yang terpakai kepada semua bahasa: sintaks interpolasi, pemeliharaan HTML, konvensyen jenis rentetan (butang lwn. label lwn. ralat), dsb.

### Panduan mengikut bahasa (`{locale}.md`)

Menentukan peraturan khusus bahasa: laras formaliti (tu lwn. vous), tanda baca (guillemet, tanda soal terbalik), bentuk jamak, pemformatan tarikh/nombor, dan glosari terminologi.

Panduan gaya ialah input dasar yang tahan lama, bukan output yang dijana. Internationalizer membaca fail ini tetapi tidak pernah menulisnya semula. Kandungannya dicincang secara berasingan daripada glosari dan kontrak gesaan, supaya perubahan kod aplikasi tidak menjadikan terjemahan basi. Mengedit panduan sengaja menandakan lokal tersebut untuk semakan dasar; menukar susunan kata gesaan dalaman tidak berbuat demikian, melainkan versi kontrak gesaan turut berubah.

Lihat [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) untuk contoh yang berfungsi.
<!-- internationalizer:unit markdown:glossary-format -->
## Format Glosari

Fail glosari ialah tatasusunan JSON yang disimpan dalam `{glossary_dir}/{locale}.json`:

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

`variants` menyenaraikan bentuk sasaran lain yang diluluskan. `enforcement` boleh menjadi `error`, `warning`, atau diabaikan untuk gelagat ralat lalai. Istilah disuntik ke dalam gesaan LLM sebagai jadual istilah, memastikan terjemahan yang konsisten merentas aplikasi anda. Entri tepat seperti `{"source":"API","target":"API"}` juga mengecualikan keseluruhan nilai sumber yang sama daripada penemuan nilai belum diterjemahkan mod ketat; ini tidak mengecualikan nilai yang lebih panjang yang hanya mengandungi `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## Memori Terjemahan

Memori terjemahan disimpan sebagai fail JSONL (satu rekod JSON setiap baris). Setiap rekod mengandungi:

- Himpunan, kunci, nilai sumber, nilai diterjemahkan, dan lokal sasaran kanonik
- Cincangan sumber, panduan gaya, glosari, kontrak gesaan, dan dasar gabungan
- Penyedia dan model yang menghasilkan terjemahan
- Cap masa

Pada larian berikutnya, rentetan dengan cincangan sumber dan dasar yang sama disajikan daripada cache tanpa memanggil LLM. Laluan lalai berada di bawah direktori `.internationalizer/` yang diabaikan, menjadikannya cache setempat. Tetapkan `tm_path` kepada lokasi yang dijejaki sekiranya projek anda sengaja berkongsi memori terjemahan. Manifes `.internationalizer.lock` yang boleh disemak diindividukan versinya secara berasingan.
<!-- internationalizer:unit markdown:supported-formats -->
## Format yang Disokong

| Format | Sambungan | Mod |
|--------|-----------|------|
| JSON | `.json` | Nilai kunci (bersarang, diratakan tatatanda titik) |
| YAML | `.yml`, `.yaml` | Nilai kunci (mengekalkan komen dan susunan) |
| Markdown | `.md`, `.mdx` | Mukadimah dan bahagian peringkat H2 |

Sasaran Markdown mengandungi komen `internationalizer:unit` halimunan sebelum bahagian H2. Penanda stabil ini membolehkan Internationalizer menambah, memindahkan, atau mengedit satu bahagian sumber tanpa menterjemah semula bahagian lain yang tidak berkaitan. Dokumen sedia ada yang belum ditandai akan menerima penanda pada kemas kini berjaya yang seterusnya.
<!-- internationalizer:unit markdown:project-type-detection -->
## Pengesanan Jenis Projek

`internationalizer detect` mengenal pasti persediaan i18n anda dengan memeriksa:

- Kebergantungan `package.json` untuk react-i18next, next-intl, atau vue-i18n
- Struktur direktori yang sepadan dengan corak lokal lazim
- Sambungan fail dan konvensyen penamaan
<!-- internationalizer:unit markdown:architecture -->
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
## Perbandingan dengan Alternatif

| Ciri | Internationalizer | i18next | Crowdin | LLM Generik |
|---------|------------------|---------|---------|-------------|
| Terjemahan dikuasakan LLM | Ya | Tidak | Separa | Ya |
| Panduan gaya mengikut bahasa | Ya | Tidak | Tidak | Tidak |
| Penguatkuasaan glosari | Ya | Tidak | Ya | Tidak |
| Memori terjemahan | Ya | Tidak | Ya | Tidak |
| Pelaksanaan CLI / tempatan | Ya | T/B | Tidak | Manual |
| Fail mesra Git | Ya | Ya | Separa | Manual |
| Tiada kebergantungan SaaS | Ya | Ya | Tidak | Berbeza-beza |
| Sumber terbuka (AGPL-3.0) | Ya | Ya | Tidak | Berbeza-beza |
<!-- internationalizer:unit markdown:license -->
## Lesen

[AGPL-3.0](../../LICENSE)

Lihat [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) untuk notis kebergantungan.
<!-- internationalizer:unit markdown:contributing -->
## Menyumbang

Lihat [CONTRIBUTING.md](../../CONTRIBUTING.md) untuk persediaan pembangunan dan garis panduan. Semua sumbangan memerlukan perakuan DCO.
