> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Pipeline internasionalisasi AI-native untuk proyek perangkat lunak. Terjemahkan, validasi, dan kelola berkas i18n menggunakan LLM.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](../../LICENSE)

<p align="center">
<a href="ar.md">العربية</a> · <a href="bn.md">বাংলা</a> · <a href="cs.md">Čeština</a> · <a href="da.md">Dansk</a> · <a href="de.md">Deutsch</a> · <a href="el.md">Ελληνικά</a> · <a href="es.md">Español</a> · <a href="fi.md">Suomi</a> · <a href="fr.md">Français</a> · <a href="he.md">עברית</a> · <a href="hi.md">हिन्दी</a> · <a href="id.md">Indonesia</a> · <a href="it.md">Italiano</a> · <a href="ja.md">日本語</a> · <a href="ko.md">한국어</a> · <a href="ms.md">Bahasa Melayu</a><br><a href="nl.md">Nederlands</a> · <a href="pa.md">ਪੰਜਾਬੀ</a> · <a href="pl.md">Polski</a> · <a href="pt-BR.md">Português</a> · <a href="ro.md">Română</a> · <a href="ru.md">Русский</a> · <a href="sv.md">Svenska</a> · <a href="te.md">తెలుగు</a> · <a href="th.md">ไทย</a> · <a href="tr.md">Türkçe</a> · <a href="uk.md">Українська</a> · <a href="vi.md">Tiếng Việt</a> · <a href="yue.md">粵語</a> · <a href="zh-CN.md">简体中文</a> · <a href="zh-TW.md">繁體中文</a>
</p>

---
<!-- internationalizer:unit markdown:why-internationalizer -->
## Mengapa Internationalizer?

Sebagian besar alat i18n berupa pustaka *runtime* (i18next, react-intl) atau platform SaaS manajemen kunci (Crowdin, Lokalise). Belum ada yang menyelesaikan masalah penerjemahan sebenarnya dengan optimal:

- **Penerjemahan manual** tidak dapat diskalakan melampaui beberapa bahasa
- **API penerjemahan mesin** (Google Translate, DeepL) mengabaikan terminologi, nada, dan konvensi UI Anda
- **Penerjemahan LLM generik** memang bekerja lebih baik, tetapi tanpa glosarium dan panduan gaya, hasilnya menjadi tidak konsisten

Internationalizer hadir dengan pendekatan berbeda. Ini adalah **pipeline CLI** yang menggabungkan kemampuan penerjemahan LLM dengan:

- **Glosarium per bahasa** — menegakkan konsistensi terminologi di seluruh aplikasi Anda
- **Panduan gaya per bahasa** — mengontrol nada, tingkat formalitas, bentuk jamak, dan tipografi
- **Memori terjemahan** — melewati string yang tidak berubah, menghemat biaya panggilan API
- **Validasi deterministik** — mendeteksi kunci yang hilang atau berlebih, pergeseran struktur terlindungi, pelanggaran glosarium, serta galat bentuk jamak atau ICU sebelum dirilis
<!-- internationalizer:unit markdown:installation -->
## Instalasi

Instal melalui npm:

```bash
npm install -g internationalizer
```

Atau jalankan tanpa penginstalan global:

```bash
npx internationalizer --help
```

Paket npm menginstal biner siap pakai yang sesuai dari npm melalui dependensi opsional khusus platform.

Instal dengan Go:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Atau bangun dari kode sumber:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```
<!-- internationalizer:unit markdown:npm-packages -->
## Paket npm

- Tag Git dan versi paket npm harus cocok, misalnya `v0.1.0` dan `0.1.0`
- Paket utama `internationalizer` bergantung pada paket platform seperti `internationalizer-darwin-arm64`
- Target npm yang didukung: macOS arm64/x64, Linux arm64/x64, Windows x64
- Publikasi CI memerlukan *secret* GitHub bernama `NPM_TOKEN`
<!-- internationalizer:unit markdown:quick-start -->
## Memulai Cepat

1. Buat berkas konfigurasi di direktori root proyek Anda:

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

2. Atur kunci API Anda:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Pratinjau apa yang akan diterjemahkan:

```bash
internationalizer translate --dry-run
```

4. Jalankan penerjemahan:

```bash
internationalizer translate
```

5. Validasi semua lokal:

```bash
internationalizer validate
```
<!-- internationalizer:unit markdown:commands -->
## Perintah

### `translate`

Temukan kunci yang hilang atau kedaluwarsa, lalu terjemahkan melalui LLM.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --adopt-existing   # baseline existing translations without API calls
internationalizer translate --refresh-policy   # refresh prompt/style/model-stale entries
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

Status penerjemahan melaporkan kondisi kunci secara independen, baik yang hilang (missing), sumber kedaluwarsa (source-stale), kebijakan kedaluwarsa (policy-stale), mutakhir (current), maupun yang diedit manual, sehingga pengeditan manual tidak dapat menyamarkan perubahan sumber atau kebijakan. Nilai dengan kebijakan kedaluwarsa akan dilaporkan tetapi hanya diterjemahkan ulang saat menggunakan opsi `--refresh-policy`. Nilai yang diedit manual tidak akan pernah ditimpa secara otomatis. Gunakan `--adopt-existing` saat memperkenalkan manifes ke terjemahan yang sudah ditinjau atau saat secara eksplisit menerima hasil suntingan manual sebagai *baseline* baru.

### `validate`

Periksa semua berkas lokal terhadap bundel sumbernya. Validasi standar memeriksa cakupan struktural (persentase kunci target yang diperlukan yang tersedia), melaporkan kunci berlebih sebagai peringatan, dan menghasilkan kegagalan jika ada kunci yang hilang, interpolasi tidak cocok, atau struktur ICU MessageFormat tidak valid.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
internationalizer validate --strict             # enforce translation quality rules
internationalizer validate --require-state      # require current manifest provenance
```

`--strict` juga melaporkan cakupan yang diterjemahkan. Nilai linguistik yang identik dengan sumbernya dianggap belum diterjemahkan kecuali jika glosarium secara eksplisit mencantumkan entri sumber dan target yang sama persis untuk nilai lengkap tersebut; pengaturan `ignore_case` dipatuhi, tetapi istilah glosarium yang disematkan di dalam nilai yang lebih panjang tidak memenuhi syarat pengecualian. Mode ketat (*strict*) akan gagal jika terdapat kunci ekstra, nilai identik dengan sumber, perubahan struktur interpolasi/HTML/kode/tautan Markdown, pelanggaran glosarium, serta bentuk jamak yang belum dikonfigurasi.

`--require-state` memverifikasi setiap target terhadap `.internationalizer.lock`. Opsi ini akan gagal jika suatu kunci tidak dilacak (*untracked*), atau ketika rekaman sumber, kebijakan penerjemahan, atau hash target telah kedaluwarsa. Opsi ini dapat digabungkan dengan `--strict`.

Laporan teks dan JSON menggunakan kode temuan yang stabil:

| Kode | Arti |
| --- | --- |
| `missing_key` / `extra_key` | Himpunan kunci sumber dan target berbeda |
| `blank_translation` | Sumber non-kosong memiliki target kosong pada mode *strict* |
| `source_identical` | Nilai linguistik mode *strict* belum diterjemahkan |
| `protected_structure_mismatch` | Struktur interpolasi, HTML, kode, atau tautan berubah |
| `glossary_violation` | Istilah target atau varian yang disetujui tidak ditemukan |
| `plural_form_missing` | Bentuk jamak lokal yang dikonfigurasi tidak ditemukan |
| `icu_message_syntax` | Pesan ICU sumber atau target salah format |
| `icu_argument_mismatch` | Nama argumen, tipe, atau gaya pemformat ICU berbeda |
| `icu_selector_mismatch` | Selektor berbeda atau kategori jamak tidak valid untuk lokal target |
| `untracked` | Tidak ada rekaman manifes untuk target |
| `source_stale` | Konten sumber berubah setelah terjemahan direkam |
| `policy_stale` | Pengaturan prompt atau model yang dihasilkan telah berubah |
| `target_modified` | Konten target berbeda dari rekaman manifes |

### `detect`

Deteksi otomatis framework i18n dan berikan saran konfigurasi.

```bash
internationalizer detect
```

Mendukung: react-i18next, next-intl, vue-i18n, JSON standar (*vanilla*), dokumen markdown.

### `glossary`

Kelola istilah glosarium per bahasa yang wajib diterapkan selama penerjemahan.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Kelola memori terjemahan (cache JSONL dari string yang telah diterjemahkan sebelumnya).

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```
<!-- internationalizer:unit markdown:configuration-reference -->
## Referensi Konfigurasi

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

Pengenal lokal harus berupa tag BCP 47 yang terformat dengan benar seperti `fr`, `pt-BR`, atau `sr-Latn-RS`. Lokal target yang setara secara kanonikal akan ditolak sebagai duplikat, dan penimpaan penyedia khusus lokal dicocokkan dengan ejaan yang setara secara kanonikal. Pada contoh di atas, lokal tanpa penimpaan—termasuk bahasa Jepang—akan mewarisi konfigurasi global Gemini.

Nilai ICU MessageFormat diurai secara struktural. Argumen sederhana, `select`, `plural`, `selectordinal`, `number`, `date`, dan `time` didukung, termasuk pesan bersarang, *offset* bentuk jamak, selektor angka persis, serta `#`. Validasi memeriksa sintaksis, tipe argumen, gaya pemformat, *offset* bentuk jamak, identitas cabang *select*, serta kategori jamak CLDR lokal target. Keluaran penyedia yang merusak invarian ini akan ditolak sebelum berkas lokal atau rekaman memori terjemahan ditulis.

Dengan `i18next-v4`, famili jamak sumber yang dikenali diperluas selama penerjemahan ke kategori CLDR lokal target. Kategori khusus target menggunakan nilai `_other` famili sumber sebagai templat penerjemahannya. Validasi *strict* mewajibkan kategori target tersebut; kategori khusus sumber bersifat opsional bagi lokal target yang tidak menggunakannya.
<!-- internationalizer:unit markdown:style-guides -->
## Panduan Gaya

Panduan gaya adalah berkas Markdown yang disuntikkan ke dalam prompt penerjemahan LLM. Panduan ini mengontrol nada, tingkat formalitas, tipografi, dan konvensi khusus bahasa lainnya.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Konvensi bersama (`_conventions.md`)

Mendefinisikan aturan yang berlaku untuk semua bahasa: sintaksis interpolasi, pelestarian HTML, konvensi tipe string (tombol vs. label vs. galat), dll.

### Panduan per bahasa (`{locale}.md`)

Mendefinisikan aturan khusus bahasa: tingkat formalitas (tu vs. vous), tanda baca (guillemet, tanda tanya terbalik), bentuk jamak, pemformatan tanggal/angka, dan glosarium terminologi.

Panduan gaya merupakan input kebijakan yang tahan lama (*durable*), bukan hasil yang di-*generate*. Internationalizer membacanya namun tidak pernah menulis ulang berkas tersebut. Isinya di-hash secara terpisah dari glosarium dan kontrak prompt, sehingga perubahan kode aplikasi tidak menyebabkan terjemahan menjadi kedaluwarsa. Mengedit panduan gaya secara sengaja menandai lokal tersebut untuk ditinjau ulang kebijakannya; mengubah susunan kata prompt internal tidak memicunya, kecuali jika versi kontrak prompt juga berubah.

Lihat [`examples/react-app/style-guides/`](../../examples/react-app/style-guides/) untuk contoh yang berfungsi.
<!-- internationalizer:unit markdown:glossary-format -->
## Format Glosarium

Berkas glosarium adalah array JSON yang disimpan di `{glossary_dir}/{locale}.json`:

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

`variants` mencantumkan bentuk target lain yang disetujui. `enforcement` dapat bernilai `error`, `warning`, atau diabaikan untuk perilaku bawaan galat. Istilah disuntikkan ke dalam prompt LLM sebagai tabel terminologi, memastikan penerjemahan yang konsisten di seluruh aplikasi Anda. Entri persis seperti `{"source":"API","target":"API"}` juga mengecualikan nilai lengkap yang identik dengan sumber tersebut dari temuan nilai belum diterjemahkan pada mode *strict*; entri ini tidak mengecualikan nilai yang lebih panjang yang hanya sekadar mengandung `API`.
<!-- internationalizer:unit markdown:translation-memory -->
## Memori Terjemahan

Memori terjemahan disimpan sebagai berkas JSONL (satu rekaman JSON per baris). Setiap rekaman berisi:

- Bundel, kunci, nilai sumber, nilai terjemahan, dan lokal target kanonikal
- Hash sumber, panduan gaya, glosarium, kontrak prompt, dan kebijakan gabungan
- Penyedia dan model yang menghasilkan terjemahan
- Stempel waktu (*timestamp*)

Pada eksekusi berikutnya, string dengan hash sumber dan kebijakan yang sama diambil langsung dari cache tanpa memanggil LLM. Jalur default berada di bawah direktori yang diabaikan `.internationalizer/`, sehingga tetap menjadi cache lokal. Atur `tm_path` ke lokasi yang dilacak Git jika proyek Anda memang ingin membagikan memori terjemahan secara sengaja. Manifes `.internationalizer.lock` yang dapat ditinjau diatur versinya secara terpisah.
<!-- internationalizer:unit markdown:supported-formats -->
## Format yang Didukung

| Format | Ekstensi | Mode |
|--------|-----------|------|
| JSON | `.json` | Kunci-nilai (bersarang, diratakan dengan notasi titik) |
| YAML | `.yml`, `.yaml` | Kunci-nilai (mempertahankan komentar dan urutan) |
| Markdown | `.md`, `.mdx` | Pembuka (*preamble*) dan bagian tingkat H2 |

Target Markdown memuat komentar tersembunyi `internationalizer:unit` sebelum bagian H2. Penanda stabil ini memungkinkan Internationalizer menambahkan, memindahkan, atau mengedit satu bagian sumber tanpa menerjemahkan ulang bagian lain yang tidak terkait. Dokumen yang belum memiliki penanda akan otomatis menerimanya pada pembaruan berhasil berikutnya.
<!-- internationalizer:unit markdown:project-type-detection -->
## Deteksi Jenis Proyek

`internationalizer detect` mengidentifikasi konfigurasi i18n Anda dengan memeriksa:

- Dependensi `package.json` untuk react-i18next, next-intl, atau vue-i18n
- Struktur direktori yang cocok dengan pola lokal umum
- Ekstensi berkas dan konvensi penamaan
<!-- internationalizer:unit markdown:architecture -->
## Arsitektur

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

| Fitur | Internationalizer | i18next | Crowdin | LLM Generik |
|---------|------------------|---------|---------|-------------|
| Penerjemahan berbasis LLM | Ya | Tidak | Sebagian | Ya |
| Panduan gaya per bahasa | Ya | Tidak | Tidak | Tidak |
| Penegakan glosarium | Ya | Tidak | Ya | Tidak |
| Memori terjemahan | Ya | Tidak | Ya | Tidak |
| Eksekusi CLI / lokal | Ya | N/A | Tidak | Manual |
| Berkas ramah Git | Ya | Ya | Sebagian | Manual |
| Tanpa dependensi SaaS | Ya | Ya | Tidak | Bervariasi |
| Sumber terbuka (AGPL-3.0) | Ya | Ya | Tidak | Bervariasi |
<!-- internationalizer:unit markdown:license -->
## Lisensi

[AGPL-3.0](../../LICENSE)

Lihat [THIRD_PARTY_NOTICES.md](../../THIRD_PARTY_NOTICES.md) untuk pemberitahuan dependensi.
<!-- internationalizer:unit markdown:contributing -->
## Berkontribusi

Lihat [CONTRIBUTING.md](../../CONTRIBUTING.md) untuk panduan dan penyiapan pengembangan. Semua kontribusi memerlukan persetujuan DCO.
