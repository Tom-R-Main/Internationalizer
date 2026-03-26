> [English (original)](../../README.md)

<p align="center">
  <img src="../../assets/logo.svg" alt="Internationalizer" width="480">
</p>

# Internationalizer

Yazılım projeleri için AI tabanlı uluslararasılaştırma (internationalization) ardışık düzeni. LLM'leri kullanarak i18n dosyalarını çevirin, doğrulayın ve yönetin.

[![CI](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml/badge.svg)](https://github.com/Tom-R-Main/Internationalizer/actions/workflows/ci.yml)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

<p align="center">
<a href="docs/i18n/ar.md">العربية</a> · <a href="docs/i18n/bn.md">বাংলা</a> · <a href="docs/i18n/de.md">Deutsch</a> · <a href="docs/i18n/es.md">Español</a> · <a href="docs/i18n/fr.md">Français</a> · <a href="docs/i18n/hi.md">हिन्दी</a> · <a href="docs/i18n/id.md">Indonesia</a> · <a href="docs/i18n/ja.md">日本語</a> · <a href="docs/i18n/pa.md">ਪੰਜਾਬੀ</a> · <a href="docs/i18n/pt-BR.md">Português</a> · <a href="docs/i18n/ru.md">Русский</a> · <a href="docs/i18n/te.md">తెలుగు</a><br><a href="docs/i18n/th.md">ไทย</a> · <a href="docs/i18n/uk.md">Українська</a> · <a href="docs/i18n/yue.md">粵語</a> · <a href="docs/i18n/zh-CN.md">简体中文</a> · <a href="docs/i18n/zh-TW.md">繁體中文</a>
</p>

---

## Neden Internationalizer?

Çoğu i18n aracı ya çalışma zamanı kütüphaneleridir (i18next, react-intl) ya da anahtar yönetimi SaaS platformlarıdır (Crowdin, Lokalise). Hiçbiri asıl çeviri sorununu iyi bir şekilde çözmez:

- **Manuel çeviri** birkaç dilden sonra ölçeklenemez
- **Makine çevirisi API'leri** (Google Translate, DeepL) terminolojinizi, üslubunuzu ve kullanıcı arayüzü kurallarınızı göz ardı eder
- **Genel LLM çevirisi** daha iyi çalışır, ancak sözlükler ve stil kılavuzları olmadan tutarsız sonuçlar elde edersiniz

Internationalizer farklıdır. LLM çevirisini şunlarla birleştiren bir **CLI ardışık düzenidir**:

- **Dile özel sözlükler** — uygulamanız genelinde tutarlı terminolojiyi zorunlu kılar
- **Dile özel stil kılavuzları** — üslubu, resmiyeti, çoğullaştırmayı ve tipografiyi kontrol eder
- **Çeviri belleği** — değişmeyen metinleri atlar, API çağrılarında tasarruf sağlar
- **Anahtar doğrulaması** — eksik çevirileri ve interpolasyon uyuşmazlıklarını yayına girmeden önce yakalar

## Kurulum

npm üzerinden kurun:

```bash
npm install -g internationalizer
```

Veya global kurulum yapmadan çalıştırın:

```bash
npx internationalizer --help
```

npm paketi, platforma özel isteğe bağlı bağımlılıklar aracılığıyla npm'den eşleşen önceden derlenmiş ikili dosyayı kurar.

Go ile kurun:

```bash
go install github.com/Tom-R-Main/Internationalizer/cmd/internationalizer@latest
```

Veya kaynaktan derleyin:

```bash
git clone https://github.com/Tom-R-Main/Internationalizer.git
cd Internationalizer
go build -o internationalizer ./cmd/internationalizer
```

## npm paketleri

- Git etiketleri ve npm paket sürümleri eşleşmelidir, örneğin `v0.1.0` ve `0.1.0`
- Kök `internationalizer` paketi, `internationalizer-darwin-arm64` gibi platform paketlerine bağlıdır
- Desteklenen npm hedefleri: macOS arm64/x64, Linux arm64/x64, Windows x64
- CI yayınlaması için `NPM_TOKEN` adında bir GitHub gizli anahtarı (secret) gerekir

## Hızlı başlangıç

1. Projenizin kök dizininde bir yapılandırma dosyası oluşturun:

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

2. API anahtarınızı ayarlayın:

```bash
export GOOGLE_AI_STUDIO_API_KEY=your-ai-studio-key
```

3. Nelerin çevrileceğini önizleyin:

```bash
internationalizer translate --dry-run
```

4. Çeviriyi çalıştırın:

```bash
internationalizer translate
```

5. Tüm yerel ayarları (locales) doğrulayın:

```bash
internationalizer validate
```

## Komutlar

### `translate`

Eksik anahtarları bulun ve bir LLM aracılığıyla çevirin.

```bash
internationalizer translate                    # translate all locales
internationalizer translate -l fr              # translate French only
internationalizer translate --dry-run          # preview without API calls
internationalizer translate --batch-size 20    # smaller batches
internationalizer translate --concurrency 2    # fewer parallel calls
```

### `validate`

Tüm yerel ayar dosyalarını eksik anahtarlar, fazladan anahtarlar ve interpolasyon uyuşmazlıkları açısından kontrol edin.

```bash
internationalizer validate                     # human-readable output
internationalizer validate --json              # machine-readable JSON
internationalizer validate -q                  # exit code only
```

### `detect`

i18n çerçevesini otomatik olarak algılayın ve bir yapılandırma önerin.

```bash
internationalizer detect
```

Desteklenenler: react-i18next, next-intl, vue-i18n, saf JSON, markdown belgeleri.

### `glossary`

Çeviri sırasında zorunlu kılınan dile özel sözlük terimlerini yönetin.

```bash
internationalizer glossary list --locale fr
internationalizer glossary add --locale fr --source "Dashboard" --target "Tableau de bord"
internationalizer glossary remove --locale fr --source "Dashboard"
```

### `tm`

Çeviri belleğini (önceden çevrilmiş metinlerin JSONL önbelleği) yönetin.

```bash
internationalizer tm stats                     # show record counts
internationalizer tm export                    # dump as JSON
internationalizer tm clear --force             # delete all records
```

## Yapılandırma referansı

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

## Stil kılavuzları

Stil kılavuzları, LLM çeviri istemine (prompt) enjekte edilen Markdown dosyalarıdır. Üslubu, resmiyeti, tipografiyi ve dile özgü diğer kuralları kontrol ederler.

```
style-guides/
  _conventions.md    # shared rules for all languages
  fr.md              # French-specific rules
  ja.md              # Japanese-specific rules
  ar.md              # Arabic-specific rules
```

### Ortak kurallar (`_conventions.md`)

Tüm diller için geçerli olan kuralları tanımlayın: interpolasyon sözdizimi, HTML'in korunması, metin türü kuralları (butonlar, etiketler, hatalar) vb.

### Dile özel kılavuzlar (`{locale}.md`)

Dile özgü kuralları tanımlayın: resmiyet derecesi (sen/siz), noktalama işaretleri (açılı tırnaklar, ters soru işaretleri), çoğul biçimler, tarih/sayı biçimlendirmesi ve terminoloji sözlüğü.

Çalışan bir örnek için [`examples/react-app/style-guides/`](examples/react-app/style-guides/) dizinine bakın.

## Sözlük formatı

Sözlük dosyaları, `{glossary_dir}/{locale}.json` içinde saklanan JSON dizileridir:

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

Terimler, LLM istemine bir terminoloji tablosu olarak enjekte edilir ve uygulamanız genelinde önemli terimlerin tutarlı bir şekilde çevrilmesini sağlar.

## Çeviri belleği

Çeviri belleği bir JSONL dosyası olarak saklanır (her satırda bir JSON kaydı). Her kayıt şunları içerir:

- Kaynak anahtar ve değer
- Çevrilen değer
- Kaynak değerin SHA-256 özeti (hash)
- Zaman damgası

Sonraki çalıştırmalarda, değişmeyen metinler LLM çağrılmadan TM önbelleğinden sunulur; bu da hem zamandan hem de API maliyetlerinden tasarruf sağlar. TM dosyası git dostudur ve yerel ayar dosyalarınızla birlikte işlenebilir (commit).

## Desteklenen formatlar

| Format | Uzantılar | Mod |
|--------|-----------|------|
| JSON | `.json` | Anahtar-değer (iç içe, nokta notasyonuyla düzleştirilmiş) |
| YAML | `.yml`, `.yaml` | Anahtar-değer (yorumları ve sıralamayı korur) |
| Markdown | `.md`, `.mdx` | Tüm belge çevirisi |

## Proje türü algılama

`internationalizer detect` aşağıdakileri kontrol ederek i18n kurulumunuzu tanımlar:

- react-i18next, next-intl veya vue-i18n için `package.json` bağımlılıkları
- Yaygın yerel ayar desenleriyle eşleşen dizin yapıları
- Dosya uzantıları ve isimlendirme kuralları

## Mimari

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

## Alternatiflerle karşılaştırma

| Özellik | Internationalizer | i18next | Crowdin | Genel LLM |
|---------|------------------|---------|---------|-------------|
| LLM destekli çeviri | Evet | Hayır | Kısmen | Evet |
| Dile özel stil kılavuzları | Evet | Hayır | Hayır | Hayır |
| Sözlük zorunluluğu | Evet | Hayır | Evet | Hayır |
| Çeviri belleği | Evet | Hayır | Evet | Hayır |
| CLI / yerel yürütme | Evet | N/A | Hayır | Manuel |
| Git dostu dosyalar | Evet | Evet | Kısmen | Manuel |
| SaaS bağımlılığı yok | Evet | Evet | Hayır | Değişir |
| Açık kaynak (AGPL-3.0) | Evet | Evet | Hayır | Değişir |

## Lisans

[AGPL-3.0](LICENSE)

## Katkıda bulunma

Geliştirme kurulumu ve yönergeler için [CONTRIBUTING.md](CONTRIBUTING.md) dosyasına bakın. Tüm katkılar DCO onayı gerektirir.

